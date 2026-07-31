package vaults_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"

	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

type EntryUpdate struct {
	ID         string
	CID        string
	TotalBytes int
	IsDirty    bool
	Reused     bool
}
type FolderUpdate struct {
	ID  string
	CID string
}
type SyncPlan struct {
	EntryUpdates  []EntryUpdate
	FolderUpdates []FolderUpdate

	ReusedLinks []vaults_domain.Link
	NewLinks    []vaults_domain.Link

	RootChanged bool
}
type LoadAttachmentResponse struct {
	File []byte
	Hash string
}
type SyncMode int

const (
	FullSync SyncMode = iota
	IncrementalSync
)

type SyncPolicy struct {
	AllowReuse bool
}

type IPFSClient interface {
	Add(data []byte) (string, error)
	Get(cid string) ([]byte, error)
}

type Encryptor interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
}

type VaultRepository interface {
	GetVault(vaultID string) (*vaults_domain.Vault, error)
	GetVaultByCID(vaultID string) (*vaults_domain.Vault, error)
	UpdateVaultCID(vaultID, cid string) error
}
type VaultNodeRepository interface {
	GetEntries(session vault_session.Session) (*vaults_domain.Entries, error)
	GetFolders(session vault_session.Session) ([]vaults_domain.Folder, error)
	GetAttachements(session vault_session.Session) ([]vaults_domain.Attachment, error)
}

type VaultHandlerInterface interface {
	LoadAttachment(userID string, vaultName string, hash string, formatReturned string) (*LoadAttachmentResponse, error)
}

type VaultService struct {
	VaultHandler VaultHandlerInterface
	Encryptor    Encryptor
	Repo         VaultRepository
	NodeRepo     VaultNodeRepository
	Password     string
	VaultKey     []byte
	IPFSHandler  vault_commands.CreateIPFSPayloadHandler
	VaultCtx     app_config_domain.VaultContext
	IsDraftMode  bool
	DraftStorage DraftStorage
	Personal     string
	C3           string
}

func NewVaultServiceDryRun(
	repo VaultRepository,
	nodeRepo VaultNodeRepository,
	encryptor Encryptor,
	vc app_config_domain.VaultContext,
	ih vault_commands.CreateIPFSPayloadHandler,
) *VaultService {

	return &VaultService{
		Encryptor:    encryptor,
		Repo:         repo,
		NodeRepo:     nodeRepo,
		VaultCtx:     vc,
		IPFSHandler:  ih,
		IsDraftMode:  true,
		DraftStorage: *NewDraftStorage(),
	}
}

func NewVaultServiceReal(
	vh VaultHandlerInterface,
	encryptor Encryptor,
	ipfs vault_commands.CreateIPFSPayloadCommandHandler,
	repo VaultRepository,
	nodeRepo VaultNodeRepository,
	vc app_config_domain.VaultContext,
) *VaultService {
	return &VaultService{
		VaultHandler: vh,
		Encryptor:    encryptor,
		Repo:         repo,
		NodeRepo:     nodeRepo,
		VaultCtx:     vc,
		IPFSHandler:  &ipfs,
		IsDraftMode:  false,
	}
}

type SyncResult struct {
	NewCIDs    []string
	ReusedCIDs []string
	TotalBytes int64
}

func (s *VaultService) CommitVaultLegacy(session vault_session.Session, mode SyncMode) (string, []EntryUpdate, int, int, error) {
	utils.LogPretty("VaultService - CommitVault - ", "starting....")

	vp, err := vault_session.DecodeSessionVault(session.Vault)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	// 0. BUILD ATTACHEMENTS
	// =========================
	// attachements := vp.GetAttachments()
	// attachementLinks, err := s.buildAttachmentLinks(session.UserID, session.Runtime.VaultName, attachements, mode)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }

	// // =========================
	// // 6. ATTACHEMENTS ROOT
	// // =========================
	// attachementCIDs, _, err := s.BuildAttachmentsRoot(attachementLinks)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }

	attachementCIDs, err := s.BuildAttachmentsBranch(session, *vp, mode)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	// 1. BUILD ENTRIES
	// =========================
	// entries := vp.Entries
	// entryLinks, indexByType, indexByFolder, entryUpdates, err := s.BuildEntries(entries, mode)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }

	// // =========================
	// // 4. ENTRIES ROOT
	// // =========================
	// entriesCID, _, err := s.BuildEntriesRoot(entryLinks, mode)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }

	entriesCID, indexByType, indexByFolder, entryUpdates, err := s.BuildEntriesBranch(session, *vp, mode)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	// 2. FOLDERS
	// =========================
	// folders := vp.Folders
	// folderLinks, err := s.BuildFolders(folders)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }
	// 	// =========================
	// // 5. FOLDERS ROOT
	// // =========================
	// foldersCID, _, err := s.BuildFoldersRoot(folderLinks)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }

	foldersCID, err := s.BuildFoldersBranch(session, *vp, mode)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	// 3. INDEX
	// =========================
	indexCID, _, err := s.buildIndex(indexByType, indexByFolder)
	if err != nil {
		return "", nil, 0, 0, err
	}

	return s.SaveVaultRoot(attachementCIDs, foldersCID, entriesCID, indexCID, entryUpdates, session)
}

func (s *VaultService) CommitVault(session vault_session.Session, mode SyncMode) (string, []EntryUpdate, int, int, error) {
	utils.LogPretty("VaultService - CommitVault - ", "starting....")

	// =========================
	// 1. Personal branch
	// =========================
	personalBranch, entryUpdates, err := s.CommitVaultPersonal(session, mode)
	if err != nil {
		return "", nil, 0, 0, fmt.Errorf("VaultService - CommitVault - failed to commit personalBranch: %v", err)
	}
	s.Personal = personalBranch

	// =========================
	// 2. Collaborative branch
	// =========================
	collaborativeBranch, err := s.CommitVaultCollaborative(session, mode)
	if err != nil {
		return "", nil, 0, 0, fmt.Errorf("VaultService - CommitVault - failed to commit collaborativeBranch: %v", err)
	}
	s.C3 = collaborativeBranch

	return s.SaveVaultNodeRoot(SaveVaultNodeRootParams{
		Personal: SaveVaultPersonalParams{
			CID:          personalBranch,
			entryUpdates: entryUpdates,
		},
		CollaborativeCID: collaborativeBranch,
		Session:          session,
	})

}

// CommitPersonal
func (s *VaultService) CommitVaultPersonal(session vault_session.Session, mode SyncMode) (string, []EntryUpdate, error) {
	utils.LogPretty("VaultService - CommitVault - ", "starting....")

	vp, err := vault_session.DecodeSessionVault(session.Vault)
	if err != nil {
		return "", nil, err
	}
	attachementCIDs, err := s.BuildAttachmentsBranch(session, *vp, mode)
	if err != nil {
		return "", nil, err
	}

	// =========================
	// 1. BUILD ENTRIES
	// =========================
	entriesCID, indexByType, indexByFolder, entryUpdates, err := s.BuildEntriesBranch(session, *vp, mode)
	if err != nil {
		return "", nil, err
	}

	// =========================
	// 2. FOLDERS
	// =========================
	foldersCID, err := s.BuildFoldersBranch(session, *vp, mode)
	if err != nil {
		return "", nil, err
	}

	// =========================
	// 3. INDEX
	// =========================
	indexCID, _, err := s.buildIndex(indexByType, indexByFolder)
	if err != nil {
		return "", nil, err
	}

	// =========================
	// 7. BUILD Personal ROOT
	// =========================
	return s.PersonalNode(PersonalNodeParams{
		FoldersCID:      foldersCID,
		EntriesCID:      entriesCID,
		IndexCID:        indexCID,
		AttachementCIDs: attachementCIDs,
		entryUpdates:    entryUpdates,
		session:         session,
	})
}

// CommitCollaborative
func (s *VaultService) CommitVaultCollaborative(session vault_session.Session, mode SyncMode) (string, error) {
	utils.LogPretty("VaultService - CommitVaultCollaborative - ", "starting....")

	vp, err := vault_session.DecodeSessionVault(session.Vault)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get vaultPayload %v", err)
	}
	// =========================
	// 0. BUILD TRUST MEMBERS
	// =========================
	memberCID, err := s.BuildTrustMembersBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get trustGroupMemberCID %v", err)
	}
	fmt.Println(memberCID)

	// =========================
	// 1. BUILD TRUST GROUP
	// =========================
	trustgroupCID, trustGroupIndexbyWorkspace, trustGroupIndexbyMember, err := s.BuildTrustGroupsBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get trustGroupCID %v", err)
	}
	fmt.Println(trustgroupCID)

	// =========================
	// 2. BUILD ShareEntries
	// =========================
	shareEntriesCID, err := s.BuildShareEntriesBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get shareEntriesCID %v", err)
	}
	fmt.Println(shareEntriesCID)

	// =========================
	// 3. BUILD Assets
	// =========================
	assetsCID, assetIndexByHash, assetIndexByType, err := s.BuildAssetsBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get assetsCID %v", err)
	}
	fmt.Println(assetsCID)

	// =========================
	// 4. BUILD Threads
	// =========================
	threadCID, indexThreadByChannel, indexThreadByStatus, err := s.BuildThreadsBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get threadCID %v", err)
	}

	// =========================
	// 5. BUILD Channels
	// =========================
	channelCID, err := s.BuildChannelsBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get channelCID %v", err)
	}
	fmt.Println(channelCID)

	// =========================
	// 6. BUILD Participants
	// =========================
	participantCID, err := s.BuildParticipantsBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get participantCID %v", err)
	}
	fmt.Println(participantCID)

	// =========================
	// 6. BUILD Workspaces
	// =========================
	workspaceCID, err := s.BuildWorkspacesBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get workspaceCID %v", err)
	}

	// =========================
	// 7. BUILD Federation
	// =========================
	federationCID, federationIndexCID, err := s.BuildFederationBranch(session, *vp, mode)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get federationCID %v", err)
	}
	// =========================
	// 7. BUILD Index
	// =========================
	indexThreadCID, _, err := s.buildThreadIndex(indexThreadByChannel, indexThreadByStatus)
	if err != nil {
		return "", err
	}

	assetCID, _, err := s.buildAssetIndex(assetIndexByHash, assetIndexByType)
	if err != nil {
		return "", err
	}

	trustGroupCID, _, err := s.buildTrustGroupIndex(trustGroupIndexbyWorkspace, trustGroupIndexbyMember)
	if err != nil {
		return "", err
	}



	indexCID, _, err := s.buildCollaborativeIndex(
		indexThreadCID,
		assetCID,
		federationIndexCID,
		trustGroupCID,
	)
	if err != nil {
		return "", fmt.Errorf("VaultService - CommitVaultCollaborative - failed to get indexCID %v", err)
	}

	// =========================
	// 8. BUILD C3 ROOT
	// =========================
	return s.CollaborativeNode(CollaborativeNodeParams{
		WorkspacesCID:   workspaceCID,
		ChannelsCID:     channelCID,
		ParticipantsCID: participantCID,
		ThreadsCID:      threadCID,
		ShareEntriesCID: shareEntriesCID,
		PayloadsCID:     assetsCID,
		TrustGroupsCID:  trustgroupCID,
		TrustMembersCID: memberCID,
		FederationCID:   federationCID,
		IndexCID:        indexCID,
	})
	// return "", nil
}

func (s *VaultService) RotateVaultKey(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) {
	//	RotateAttachmentKeys = newAttachmentsRootCID
	//	      	↓
	//	RotateEntryKeys = newEntriesRootCID
	//		   	↓
	// RotateFolderKeys = newFoldersRootCID
	//			↓
	// RotateIndexKeys = newIndexRootCID
	//			↓
	// SaveVault(newAttachmentsRootCID, newFoldersRootCID, newEntriesRootCID, newIndexRootCID, entryUpdates, session)
}

func (s *VaultService) SaveVaultRoot(
	attachementCIDs,
	foldersCID string,
	entriesCID string,
	indexCID string,
	entryUpdates []EntryUpdate,
	session vault_session.Session,
) (string, []EntryUpdate, int, int, error) {
	// =========================
	// 1. VAULT ROOT
	// =========================
	// vaultNode := vaults_domain.VaultNode{
	// 	Type:        "vault",
	// 	Version:     "1.0.0",
	// 	Folders:     vaults_domain.Link{CID: foldersCID},
	// 	Entries:     vaults_domain.Link{CID: entriesCID},
	// 	Index:       vaults_domain.Link{CID: indexCID},
	// 	Attachments: vaults_domain.Link{CID: attachementCIDs},
	// }

	// rootCID, _, err := s.putNode(vaultNode)
	// if err != nil {
	// 	return "", nil, 0, 0, err
	// }
	rootCID, err := s.Legacy(PersonalNodeParams{
		FoldersCID:      foldersCID,
		EntriesCID:      entriesCID,
		IndexCID:        indexCID,
		AttachementCIDs: attachementCIDs,
		entryUpdates:    entryUpdates,
		session:         session,
	})

	// =========================
	// 2. UPDATE REPO SESSION - TODO save in db
	// =========================
	err = s.Repo.UpdateVaultCID(session.Runtime.VaultID, rootCID)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	//	3. COST CALCULATION
	// =========================
	var totalBytes int
	var newBytes int

	for _, u := range entryUpdates {
		totalBytes += u.TotalBytes

		if !u.Reused {
			newBytes += u.TotalBytes
		}
	}

	utils.LogPretty("entryUpdates", entryUpdates)
	utils.LogPretty("totalBytes", totalBytes)
	utils.LogPretty("newBytes", newBytes)

	return rootCID, entryUpdates, totalBytes, newBytes, nil
}

type PersonalNodeParams struct {
	FoldersCID      string
	EntriesCID      string
	IndexCID        string
	AttachementCIDs string
	entryUpdates    []EntryUpdate
	session         vault_session.Session
}

func (s *VaultService) Legacy(args PersonalNodeParams) (string, error) {
	vaultNode := vaults_domain.VaultNode{
		Type:        "vault",
		Version:     "1.0.0",
		Folders:     vaults_domain.Link{CID: args.FoldersCID},
		Entries:     vaults_domain.Link{CID: args.EntriesCID},
		Index:       vaults_domain.Link{CID: args.IndexCID},
		Attachments: vaults_domain.Link{CID: args.AttachementCIDs},
	}

	vaultNodeRootCID, _, err := s.putNode(vaultNode)
	if err != nil {
		return "", fmt.Errorf("VaultService - V1 - failed to get vaultNodeRootCID: %v", err)
	}
	return vaultNodeRootCID, nil
}
func (s *VaultService) PersonalNode(args PersonalNodeParams) (string, []EntryUpdate, error) {
	personalNode := vaults_domain.VaultNode{
		Type:        "personal_vault",
		Version:     "1.1.0",
		Folders:     vaults_domain.Link{CID: args.FoldersCID},
		Entries:     vaults_domain.Link{CID: args.EntriesCID},
		Index:       vaults_domain.Link{CID: args.IndexCID},
		Attachments: vaults_domain.Link{CID: args.AttachementCIDs},
	}

	personalRootCID, _, err := s.putNode(personalNode)
	if err != nil {
		return "", nil, fmt.Errorf("VaultService - V1 - failed to get personalRootCID: %v", err)
	}
	return personalRootCID, args.entryUpdates, nil
}

type CollaborativeNodeParams struct {
	WorkspacesCID   string
	ChannelsCID     string
	ParticipantsCID string
	ThreadsCID      string
	ShareEntriesCID string
	PayloadsCID     string
	TrustGroupsCID  string
	TrustMembersCID string
	FederationCID   string
	IndexCID        string
}

func (s *VaultService) CollaborativeNode(args CollaborativeNodeParams) (string, error) {
	collaborativeNode := vaults_domain.CollaborativeNode{
		Type:         "collaborative_vault",
		Version:      "1.0.0",
		Workspaces:   vaults_domain.Link{CID: args.WorkspacesCID},
		Participants: vaults_domain.Link{CID: args.ParticipantsCID},
		Channels:     vaults_domain.Link{CID: args.ChannelsCID},
		Threads:      vaults_domain.Link{CID: args.ThreadsCID},
		ShareEntries: vaults_domain.Link{CID: args.ShareEntriesCID},
		Assets:       vaults_domain.Link{CID: args.PayloadsCID},
		TrustGroups:  vaults_domain.Link{CID: args.TrustGroupsCID},
		TrustMembers: vaults_domain.Link{CID: args.TrustMembersCID},
		Federation:   vaults_domain.Link{CID: args.FederationCID},
		Index:        vaults_domain.Link{CID: args.IndexCID},
	}

	collaborativeRootCID, _, err := s.putNode(collaborativeNode)
	if err != nil {
		return "", fmt.Errorf("VaultService - V1 - failed to get collaborativeRootCID: %v", err)
	}
	utils.LogPretty("collaborativeNode", collaborativeNode)
	return collaborativeRootCID, nil
}

type SaveVaultPersonalParams struct {
	CID          string
	entryUpdates []EntryUpdate
}
type SaveVaultNodeRootParams struct {
	Personal         SaveVaultPersonalParams
	CollaborativeCID string
	Session          vault_session.Session
}

func (s *VaultService) SaveVaultNodeRoot(args SaveVaultNodeRootParams) (string, []EntryUpdate, int, int, error) {
	vaultNode := vaults_domain.VaultNodeBeta{
		Type:          "vault",
		Version:       "1.0.0",
		Personal:      vaults_domain.Link{CID: args.Personal.CID},
		Collaborative: vaults_domain.Link{CID: args.CollaborativeCID},
	}

	rootCID, _, err := s.putNode(vaultNode)
	if err != nil {
		return "", nil, 0, 0, fmt.Errorf("VaultService - V1 - failed to get rootCID: %v", err)
	}

	// =========================
	// 2. UPDATE REPO SESSION - TODO save in db
	// =========================
	err = s.Repo.UpdateVaultCID(args.Session.Runtime.VaultID, rootCID)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	//	3. COST CALCULATION
	// =========================
	var totalBytes int
	var newBytes int

	for _, u := range args.Personal.entryUpdates {
		totalBytes += u.TotalBytes

		if !u.Reused {
			newBytes += u.TotalBytes
		}
	}

	// utils.LogPretty("entryUpdates", args.Personal.entryUpdates)
	// utils.LogPretty("totalBytes", totalBytes)
	// utils.LogPretty("newBytes", newBytes)
	utils.LogPretty("vaultNode", vaultNode)

	return rootCID, args.Personal.entryUpdates, totalBytes, newBytes, nil
}

func (s *VaultService) putNode(v interface{}) (string, int, error) {
	data, err := json.Marshal(v) // switch to DAG-CBOR later
	if err != nil {
		return "", 0, err
	}

	if s.IsDraftMode {
		cid, _ := s.DraftStorage.Add(data)
		return cid, len(data), nil
	}

	if s.IPFSHandler == nil {
		return "", 0, errors.New("IPFSHandler is nil")
	}

	if s.VaultCtx.UserID == "" {
		return "", 0, errors.New("VaultCtx is nil")
	}

	res, err := s.IPFSHandler.Execute(
		context.Background(),
		s.VaultCtx,
		vault_commands.CreateIPFSPayloadCommand{
			Data:             data,
			Password:         s.Password,
			UserOnboardingID: s.VaultCtx.UserOnboarding,
		})
	if err != nil {
		utils.LogPretty("VaultService - putNode - err", err)
		return "", 0, err
	}
	log.Println("RETURNED CID:", res.CID)

	return res.CID, len(data), nil
}
func (s *VaultService) putRawFile(data []byte) (string, int, error) {

	if s.IsDraftMode {
		cid, _ := s.DraftStorage.Add(data)
		return cid, len(data), nil
	}

	res, err := s.IPFSHandler.Execute(
		context.Background(),
		s.VaultCtx,
		vault_commands.CreateIPFSPayloadCommand{
			Data:             data,
			Password:         s.Password,
			UserOnboardingID: s.VaultCtx.UserOnboarding,
		},
	)

	if err != nil {
		return "", 0, err
	}

	return res.CID, len(data), nil
}

func (s *VaultService) FindOrphans()  {}
func (s *VaultService) DeleteUnused() {}

func resolvePolicy(mode SyncMode) SyncPolicy {
	switch mode {
	case IncrementalSync:
		return SyncPolicy{
			AllowReuse: true,
		}
	default:
		return SyncPolicy{
			AllowReuse: false,
		}
	}
}

type AESEncryptor struct{}

func (e *AESEncryptor) Encrypt(data []byte, key []byte) ([]byte, error) {
	utils.LogPretty("AESEncryptor Encrypt", "ok")
	vc := &vault_infrastructure_crypto.AESService{}
	return vc.Encrypt(data, key)
}

type NoopEncryptor struct{}

func (e *NoopEncryptor) Encrypt(data []byte, b []byte) ([]byte, error) {
	utils.LogPretty("NoopEncryptor Encrypt", "ok")
	return data, nil // 👈 no encryption
}

type DraftStorage struct {
	Store map[string][]byte
	Order []string
}

func NewDraftStorage() *DraftStorage {
	return &DraftStorage{
		Store: make(map[string][]byte),
		Order: []string{},
	}
}
func (d *DraftStorage) Exists(cid string) bool {
	_, ok := d.Store[cid]
	return ok
}
func (d *DraftStorage) Add(data []byte) (string, error) {
	if d.Store == nil {
		d.Store = make(map[string][]byte)
	}

	// cid := fmt.Sprintf("cid-%d", len(d.Store)+1)
	cid, err := d.ComputeCID(data)
	if err != nil {
		return "", err
	}
	d.Store[cid] = data
	d.Order = append(d.Order, cid)
	return cid, nil
}
func (d *DraftStorage) Get(cid string) ([]byte, error) {
	v, ok := d.Store[cid]
	if !ok {
		return nil, fmt.Errorf("NewDraftStorage - Get - cid not found: %s", cid)
	}
	return v, nil
}
func (b *DraftStorage) ComputeCID(data []byte) (string, error) {
	hash, _ := mh.Sum(data, mh.SHA2_256, -1)
	c := cid.NewCidV1(cid.Raw, hash)
	return c.String(), nil
}
