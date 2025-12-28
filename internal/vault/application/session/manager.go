package vault_session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
	utils "vault-app/internal"
	share_application_events "vault-app/internal/application/events/share"
	"vault-app/internal/blockchain"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/tracecore"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
}

type ManagerV0 struct {
	mu                sync.RWMutex
	sessions          map[string]*Session 
	SessionRepository SessionRepository
	VaultRepository   vaults_domain.VaultRepository

	pendingMu      sync.Mutex
	pendingCommits map[string][]tracecore.CommitEnvelope // optionally keep in-memory pending commits per user
	SessionsMu     sync.Mutex
	logger         Logger
	NowUTC         func() string
	IsDirty        bool

	EventDispatcher share_application_events.EventDispatcher
	IPFS            *blockchain.IPFSClient
	Ctx             context.Context
}
type Manager struct {
	mu sync.RWMutex

	sessions map[string]*Session

	SessionRepository SessionRepository
	VaultRepository   vaults_domain.VaultRepository

	logger Logger
	NowUTC func() string

	EventDispatcher share_application_events.EventDispatcher
	IPFS            *blockchain.IPFSClient
	Ctx             context.Context
	IsDirty         bool

	pendingMu      sync.Mutex
	pendingCommits map[string][]tracecore.CommitEnvelope // optionally keep in-memory pending commits per user	
}


func NewManager(sessionRepository SessionRepository, vaultRepository vaults_domain.VaultRepository, logger Logger, ctx context.Context, ipfs *blockchain.IPFSClient, sessions map[string]*Session) *Manager {
	return &Manager{
		sessions:          sessions,
		SessionRepository: sessionRepository,
		VaultRepository:   vaultRepository,
		logger:            logger,
		NowUTC:            func() string { return time.Now().Format(time.RFC3339) },
		IsDirty:           false,

		EventDispatcher: share_infrastructure.InitializeEventDispatcher(),
		Ctx:             ctx,
		IPFS:            ipfs,
	}
}


func (m *Manager) Prepare(userID string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.sessions[userID]; ok {
		return s, nil
	}

	existing, _ := m.SessionRepository.GetSession(userID)
	// if err != nil {
	// 	return nil, err
	// }

	var s *Session
	if existing != nil {
		s = existing
		utils.LogPretty("Manager - Prepare - existing session", existing)
	} else {
		s = newSession(userID)
		utils.LogPretty("Manager - Prepare - new session", userID)
	}

	// s.Normalize()
	m.sessions[userID] = s

	m.logger.Info("✅ Session prepared for user %s", userID)
	return s, nil
}
// TODO check the logic of this function	
func (m *Manager) AttachVault(
	userID string,
	vault *vaults_domain.VaultPayload,
	runtime *RuntimeContext,
	lastCID string,
) (*Session, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[userID]
	if !ok {
		return nil, fmt.Errorf("no session for user %s", userID)
	}
	if s.Vault != nil {
        // 🔥 DO NOT overwrite restored vault
		utils.LogPretty("Manager - AttachVault - vault already exists", s.Vault)
        return s, nil
    }

	// Just for monitoring
	// before := vault.Hash()
	// vault.Normalize()
	// after := vault.Hash()

	// if before != after {
	// 	m.logger.Warn("Normalize mutated vault structure")
	// }

	vault.Normalize()
	s.Vault = vault.ToBytes()
	s.Runtime = runtime
	s.LastCID = lastCID
	s.LastUpdated = m.NowUTC()
	s.Dirty = false

	utils.LogPretty("Manager - AttachVault - vault attached", s.Vault)
	// s.Normalize()
	return s, nil
}

func (m *Manager) GetSession(userID string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[userID]
	if !ok {
		return nil, errors.New("no active session")
	}
	return s, nil
}

func (m *Manager) EndSession(userID string) error {
	m.mu.Lock()
	s, ok := m.sessions[userID]
	if ok {
		err := m.SessionRepository.SaveSession(userID, s)
		if err != nil {
			m.logger.Error("❌ Failed to save session for user %s: %v", userID, err)
			return err
		}
		utils.LogPretty("💾 EndSession - Session saved and closed", s)
		delete(m.sessions, userID)
	}
	m.mu.Unlock()

	return nil
}

func (s *Session) Normalize() {
	v := vaults_domain.InitEmptyVaultPayload("", "")
	if s.Vault == nil {
		s.Vault = v.ToBytes()
	}
	if s.Runtime == nil {
		s.Runtime = NewRuntimeContext()
	}
	if s.LastSynced == "" {
		s.LastSynced = time.Now().Format(time.RFC3339)
	}
	if s.LastUpdated == "" {
		s.LastUpdated = time.Now().Format(time.RFC3339)
	}
}

func newSession(userID string) *Session {
	s := &Session{
		UserID: userID,
		Dirty:  false,
	}
	s.Normalize()
	return s
}

func (m *Manager) MarkDirty(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.sessions[userID]; ok {
		s.Dirty = true
		s.LastUpdated = m.NowUTC()
	}
}

func (m *Manager) IsMarkedDirty(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[userID]
	return ok && s.Dirty
}
func (m *Manager) LogoutUser(userID string) error {
	m.logger.Info("👋 User %s logging out", userID)

	m.mu.Lock()
	session, ok := m.sessions[userID]
	m.mu.Unlock()

	if !ok {
		m.logger.Warn("⚠️ No active session for user %s", userID)
		return nil // logout is idempotent
	}

	// 🔒 Persist session snapshot
	if err := m.SessionRepository.SaveSession(userID, session); err != nil {
		m.logger.Error(
			"❌ Failed to save session for user %s: %v",
			userID, err,
		)
		return err
	}

	// 🧹 Cleanup memory
	m.pendingMu.Lock()
	delete(m.pendingCommits, userID)
	m.pendingMu.Unlock()

	m.mu.Lock()
	delete(m.sessions, userID)
	m.mu.Unlock()

	utils.LogPretty("💾 LogoutUser - Session saved & removed", session)
	return nil
}
// Set vault from CRUD entries
func (m *Manager) SetVault(userID string, vault *vaults_domain.VaultPayload) error {
	// 1. ---------- Get user session ----------
	m.mu.Lock()
	s, ok := m.sessions[userID]
	m.mu.Unlock()

	if !ok {
		return errors.New("no active session")
	}

	// 2. ---------- Set vault ----------
	s.Vault = vault.ToBytes()
	s.LastUpdated = m.NowUTC()
	// 3. ---------- Mark session as dirty ----------	
	s.Dirty = true
	utils.LogPretty("Manager - SetVault - vault set", vault)
	return nil
}	









// --------------- LEGAGIES

func (m *Manager) Get(userID string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[userID]
	return s, ok
}

func (m *Manager) Prepare0(userID string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 1. If session already exist return it
	if s, ok := m.sessions[userID]; ok {
		return s, nil
	}

	// 2. Fetch in db to find existing session
	existing, err := m.SessionRepository.GetSession(userID)
	if err != nil {
		m.logger.Error("❌ Manager V1 - Prepare - failed to get session for user %s: %v", userID, err)
		return nil, err
	}

	// 3. Hydrate session with default values
	m.sessions[userID] = &Session{
		UserID: userID,
		Dirty:  false,
	}
	// 4. Hydrate session with existing data if exist
	if existing != nil {
		utils.LogPretty("Manager V1 - Prepare - existing", existing)
		m.sessions[userID].Vault = existing.Vault
		m.sessions[userID].Dirty = existing.Dirty
		m.sessions[userID].Runtime = existing.Runtime
		m.sessions[userID].LastCID = existing.LastCID
		m.sessions[userID].LastSynced = existing.LastSynced
		m.sessions[userID].LastUpdated = existing.LastUpdated
	}	
	utils.LogPretty("Manager V1 - Prepare - sessions", m.sessions)

	return m.sessions[userID], nil
}

func (m *Manager) AttachVault0(
	dirty bool,
	userID string,
	vault *vaults_domain.VaultPayload,
	runtime *RuntimeContext,
	lastCID string,
	lastSynced string,
	lastUpdated string,
) *Session {

	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[userID]
	if !ok {
		s = &Session{UserID: userID}
		m.sessions[userID] = s
	}
	s.Vault = vault.ToBytes()
	s.Runtime = runtime
	s.LastCID = lastCID
	s.LastSynced = m.NowUTC()
	s.LastUpdated = m.NowUTC()
	s.Dirty = dirty
	// s.Normalize()

	return s
}

func (m *Manager) MarkDirty0(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.sessions[userID]; ok {
		s.LastUpdated = m.NowUTC()
		s.Dirty = true
	}
}

func (m *Manager) Close(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, userID)
}

func (m *Manager) StartSession(userID string, vault vaults_domain.VaultPayload, lastCID string, ctx *RuntimeContext) *Session {
	now := time.Now().Format(time.RFC3339)
	session := &Session{
		UserID:      userID,
		Vault:       vault.ToBytes(),
		LastCID:     lastCID,
		LastSynced:  now,
		LastUpdated: now,
		Dirty:       false,
		Runtime:     ctx,
	}
	m.sessions[userID] = session
	return session
}

func (m *Manager) AttachSession(s *Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[s.UserID] = s
}
func (m *Manager) GetSession0(userID string) (*Session, error) {
	session, ok := m.sessions[userID]
	if !ok {
		return nil, errors.New("no vault session found")
	}
	utils.LogPretty("Manager - GetSession - userID", userID)
	return session, nil
}
func (m *Manager) HasSessionForUser(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.sessions[userID]
	return ok
}
func (m *Manager) HasSession() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions) > 0
}
func (m *Manager) GetSessions() map[string]*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions
}
func (m *Manager) EndSession0(userID string) {
	if session, ok := m.sessions[userID]; ok {
		if err := m.SessionRepository.SaveSession(userID, session); err != nil {
			m.logger.Error("❌ ManagerV1 - EndSession - failed to save session for user %s: %v", userID, err)
		} else {
			m.logger.Info("💾 ManagerV1 - EndSession - Session saved for user %s", userID)
		}
	}

	delete(m.sessions, userID)
}
func (m *Manager) LogoutUser0(userID string) error {
	m.logger.Info("👋 User %s logging out", userID)
	m.mu.Lock()
	session, ok := m.sessions[userID]
	m.mu.Unlock()
	m.logger.Info("👋 User %s session logged out", session)

	if !ok {
		m.logger.Info("👋 User %s has no active session", userID)
		return fmt.Errorf("no active session for user %s", userID)
	}

	// Persist to DB
	if err := m.SessionRepository.SaveSession(userID, session); err != nil {
		m.logger.Error("❌ Manager V1 - LogoutUser - failed to save session for user %s: %v", userID, err)
		return fmt.Errorf("failed to save session for user %s: %w", userID, err)
	}
	utils.LogPretty("💾 LogoutUser - Session saved for user", session)

	// m.pendingMu.Lock()
	// delete(m.pendingCommits, userID)
	// m.pendingMu.Unlock()

	// m.SessionsMu.Lock()
	// delete(m.sessions, userID)
	// m.SessionsMu.Unlock()
	// m.logger.Info("👋 User %s logged out and session saved", userID)
	return nil
}
// todo: to erase
func (vh *Manager) IsVaultDirty() bool {
	return vh.IsDirty
}
func (vh *Manager) IsMarkedDirty0(userID string) bool {
	session, ok := vh.sessions[userID]
	if !ok {
		return false
	}
	return session.Dirty
}

func (vh *Manager) SyncVault(userID string, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 10, "stage": "retrieving session"})
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("no active session: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 20, "stage": "marshalling vault"})
	vaultBytes, err := json.Marshal(session.Vault)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 40, "stage": "encrypting vault"})
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 70, "stage": "uploading to IPFS"})
	newCID, err := vh.IPFS.AddData(encrypted)
	if err != nil {
		return "", fmt.Errorf("IPFS upload failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 90, "stage": "submitting to Stellar"})
	userCfg := session.Runtime.UserConfig
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("stellar submission failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 95, "stage": "saving metadata"})
	currentMeta, err := vh.VaultRepository.GetLatestByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get vault meta: %w", err)
	}
	newVault := vaults_domain.Vault{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	saved := vh.VaultRepository.SaveVault(&newVault)
	vh.logger.Info("💾 Vault saved for user %s: %v", userID, saved)

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 100, "stage": "complete"})

	// Update session
	session.LastCID = newCID
	session.LastSynced = time.Now().Format(time.RFC3339)
	session.Dirty = false
	vh.IsDirty = false

	vh.logger.Info("✅ Vault sync complete for user %s", userID)
	return newCID, nil
}

func (vh *Manager) EncryptFile(userID string, filePath []byte, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %s: %w", userID, err)
	}
	// ✅ Removed noisy LogPretty - too verbose for production
	// 2. Marshal in-memory vault
	vaultBytes, err := json.Marshal(session.Vault) // session.Vault.Sync()
	if err != nil {
		return "", fmt.Errorf("❌ failed to marshal vault: %w", err)
	}
	vh.logger.Info("🧱 Vault marshalled (%d bytes)", len(vaultBytes))

	// 3. Encrypt
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("❌ failed to encrypt vault: %w", err)
	}
	vh.logger.Info("🔐 Vault encrypted")

	return string(encrypted), nil
}
func (vh *Manager) UploadToIPFS(userID string, encrypted string) (string, error) {
	// GetBackendPlanParamForTransaction for managing plans from remote

	// Upload to IPFS
	newCID, err := vh.IPFS.AddData([]byte(encrypted))
	if err != nil {
		return "", fmt.Errorf("❌ failed to upload to IPFS: %w", err)
	}
	vh.logger.Info("📤 Vault uploaded to IPFS (CID: %s)", newCID)
	return newCID, nil
}
func (vh *Manager) CreateStellarCommit(userID string, newCID string) (string, error) {
	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %s: %w", userID, err)
	}

	userCfg := session.Runtime.UserConfig
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("❌ failed to submit CID to Stellar: %w", err)
	}
	vh.logger.Info("🌐 CID submitted to Stellar (TX: %s)", txHash)

	// 6. Get latest metadata
	currentMeta, err := vh.VaultRepository.GetLatestByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("❌ failed to get current vault meta: %w", err)
	}
	newVault := vaults_domain.Vault{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	saved := vh.VaultRepository.SaveVault(&newVault)
	vh.logger.Info("🗃️ VaultCID saved: %v", saved)

	// 8. Update session
	session.LastCID = newCID
	session.LastSynced = time.Now().Format(time.RFC3339)
	session.Dirty = false
	vh.IsDirty = false

	vh.logger.Info("✅ Vault sync complete for user %s", userID)
	// utils.LogPretty("session after sync", session)

	return newCID, nil
}

func (vh *Manager) EncryptVault(userID string, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %s: %w", userID, err)
	}
	// ✅ Removed noisy LogPretty - too verbose for production
	// 2. Marshal in-memory vault
	vaultBytes, err := json.Marshal(session.Vault) // session.Vault.Sync()
	if err != nil {
		return "", fmt.Errorf("❌ failed to marshal vault: %w", err)
	}
	vh.logger.Info("🧱 Vault marshalled (%d bytes)", len(vaultBytes))

	// 3. Encrypt
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("❌ failed to encrypt vault: %w", err)
	}
	vh.logger.Info("🔐 Vault encrypted")

	return string(encrypted), nil
}
func (vh *Manager) SetSessions(session map[string]*Session) {
	s := session
	vh.sessions = s
}		