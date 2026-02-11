package vault_session

import (
	"context"
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
	var s *Session
	if existing != nil {
		s = existing
		utils.LogPretty("Manager - Prepare - existing session", existing)
	} else {
		s = newSession(userID)
		utils.LogPretty("Manager - Prepare - new session", userID)
	}

	// TODO Initialize vault if nil
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
	if s.Vault == nil {
		vault.Normalize()
		s.Vault = vault.ToBytes()
		s.LastCID = lastCID
		s.LastUpdated = m.NowUTC()
		s.Dirty = false
    }
	
	s.Runtime = runtime
	return s, nil
}

func (m *Manager) AttachRuntime(userID string, runtime *RuntimeContext) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[userID]
	if !ok {
		return nil, fmt.Errorf("no session for user %s", userID)
	}
	s.Runtime = runtime
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
func (m *Manager) GetSessions() map[string]*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions
}
// Close session and save to db	: duplication risk with LogoutUser
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
// Logout user and remove session from memory : duplication risk with EndSession		
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

func (m *Manager) Sync(userID string, newCID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.sessions[userID]; ok {
		s.LastCID = newCID
		s.LastUpdated = m.NowUTC()	
		s.LastSynced = time.Now().Format(time.RFC3339)
		s.Dirty = false
	}
	m.logger.Info("✅ Vault synced for user %s", userID)
}

