package vaults_storage_tests

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"

	"gorm.io/gorm"

	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vaults_service "vault-app/internal/vault/infrastructure/service"
)

// ------------------------------------------------------------------------------------------------------------
// MOCKS
// ------------------------------------------------------------------------------------------------------------
type mockQueryExecutor struct {
	data map[string][]byte
}

func (m *mockQueryExecutor) Execute(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
) (*vault_queries.GetIPFSDataResponse, error) {

	raw, ok := m.data[cmd.CID]
	if !ok {
		return nil, fmt.Errorf(
			"mockQueryExecutor - Execute - cid not found: %s",
			cmd.CID,
		)
	}

	// 1. Vault root
	var beta vaults_domain.VaultNodeBeta
	if err := json.Unmarshal(raw, &beta); err == nil &&
		beta.Personal.CID != "" {

		return &vault_queries.GetIPFSDataResponse{
			NodeBeta: beta,
			Raw:      raw,
		}, nil
	}

	// 2. Collaborative root
	var collaborative vaults_domain.CollaborativeNode
	if err := json.Unmarshal(raw, &collaborative); err == nil &&
		collaborative.Type == "collaborative_vault" {

		return &vault_queries.GetIPFSDataResponse{
			CollaborativeNode: collaborative,
			Raw:               raw,
		}, nil
	}

	// 3. Personal root
	var personal vaults_domain.PersonalNode
	if err := json.Unmarshal(raw, &personal); err == nil &&
		personal.Version == "personal_vault" {

		return &vault_queries.GetIPFSDataResponse{
			PersonalNode: personal,
			Raw:          raw,
		}, nil
	}

	return &vault_queries.GetIPFSDataResponse{
		Raw: raw,
	}, nil
}

func (m *mockQueryExecutor) AddData(key string, value string) {
	m.data[key] = []byte(value)
}

type memoryIPFS struct {
	store map[string][]byte
}

func newMemoryIPFS() *memoryIPFS {
	return &memoryIPFS{
		store: make(map[string][]byte),
	}
}

func (m *memoryIPFS) Add(data []byte) (string, error) {
	cid := fmt.Sprintf("cid-%x", sha256.Sum256(data))
	m.store[cid] = data
	return cid, nil
}

func (m *memoryIPFS) Get(_ context.Context, cid string) ([]byte, error) {
	data, ok := m.store[cid]
	if !ok {
		return nil, fmt.Errorf("memoryIPFS - Get - cid not found: %s", cid)
	}
	return data, nil
}

type memoryCIDBuilder struct {
	ipfs *memoryIPFS
}

func (b *memoryCIDBuilder) Build(data []byte) (string, error) {
	return b.ipfs.Add(data)
}

type mockUnlockHandler struct {
	key []byte
}

func (m *mockUnlockHandler) Execute(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
	return &vault_dto.UnlockVaultResult{
		VaultKey: vaults_domain.VaultKey{
			Key: m.key,
		},
	}, nil
}

type e2eQuery struct {
	ipfs   *memoryIPFS
	crypto vaults_domain.VaultCrypto
	key    []byte
}

func (q *e2eQuery) Execute(ctx context.Context, cmd vault_queries.GetIPFSDataQuerry) (*vault_queries.GetIPFSDataResponse, error) {

	raw, err := q.ipfs.Get(ctx, cmd.CID)
	if err != nil {
		return nil, err
	}

	plain, err := q.crypto.Decrypt(raw, q.key)
	if err != nil {
		return nil, err
	}

	var node vaults_domain.VaultNode
	_ = json.Unmarshal(plain, &node)

	return &vault_queries.GetIPFSDataResponse{
		Node: node,
		Raw:  plain,
	}, nil
}

type mockNodeRepo struct{}

func (m *mockNodeRepo) GetEntries(session vault_session.Session) (*vaults_domain.Entries, error) {
	return &vaults_domain.Entries{
		Login: []vaults_domain.LoginEntry{
			{
				BaseEntry: vaults_domain.BaseEntry{
					ID:        "1",
					EntryName: "gmail",
					Type:      "login",
				},
			},
		},
	}, nil
}

func (m *mockNodeRepo) GetFolders(session vault_session.Session) ([]vaults_domain.Folder, error) {
	return []vaults_domain.Folder{}, nil
}

type fakeVaultRepo struct {
	saveCalled       bool
	savedVault       *vaults_domain.Vault
	saveError        error
	updateCalled     bool
	updateError      error
	existingVault    *vaults_domain.Vault
	ErrVaultNotFound error
	deleteCalled     bool
	deleteError      error
	v                vaults_domain.Vault
}

func (f *fakeVaultRepo) CreateVault(vault *vaults_domain.Vault) error {
	return nil
}

// func (f *fakeVaultRepo) SaveVault(v *vaults_domain.Vault) error {
// 	f.saveCalled = true
// 	f.savedVault = v
// 	return f.saveError
// }

func (f *fakeVaultRepo) GetLatestByUserID(userID string) (*vaults_domain.Vault, error) {
	if userID == "test_user" {
		return &vaults_domain.Vault{
			UserID: userID,
			Name:   "test_vault_name",
		}, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeVaultRepo) GetVault(string) (*vaults_domain.Vault, error) {
	panic("not used")
}
func (f *fakeVaultRepo) UpdateVault(*vaults_domain.Vault) error {
	f.updateCalled = true
	return f.updateError
}
func (f *fakeVaultRepo) DeleteVault(string) error {
	f.deleteCalled = true
	return f.deleteError
}
func (f *fakeVaultRepo) GetByUserIDAndName(string, string) (*vaults_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vaults_domain.ErrVaultNotFound
}
func (f *fakeVaultRepo) UpdateVaultCID(vaultID, cid string) error {
	f.updateCalled = true
	return f.updateError
}
func (f *fakeVaultRepo) GetVaultByCID(vaultID string) (*vaults_domain.Vault, error) {
	return &f.v, nil

}

func GetConfig(userID string, vaultName string) (*app_config_domain.Config, error) {
	res, err := app_config_domain.InitConfigFromVault(userID, vaultName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetCIDRoot(name string) string {
	return "cid-" + name + "-root"
}
func GetThread() string {
	return `{
		"id":"thread1",
		"asset_type": "test",
		"title":"Inspection",
		"status":"open"
	}`
}

func GetThreadRoot() string {
	return `{
		"items":[
			{"/":"cid-thread-#1"}
		]
	}`
}
func GetTrustGroup() string {
	return `{
		"id":"cid-trustgroup-#1",
		"workspace_id":"ws1",
		"name":"Engineering",
		"member_cids": "cid-trustgroupmember-#1"
	}`
}

func GetTrustGroupRoot() string {
	return `{
    "items":[
        {"/":"cid-trustgroup-#1"}
    ]
}`
}
func GetTrustGroupMembers() string {
	return `[{
    "id":"member1",
    "vault_id":"vaultA",
    "role":"owner"
}]`
}
func GetTrustGroupMembersRoot() string {
	return `{
    "items":[
        {"/":"cid-trustgroupmember-#1"}
    ]
}`
}
func GetParticipants() string {
	return `[{
    "channel_id":"channel1",
    "vault_id":"vaultA",
    "public_key":"pk",
    "role":"owner"
}]`
}
func GetParticipantRoot() string {
	return `{
    "items":[
        {"/":"cid-participant-#33"}
    ]
}`
}
func GetAssets() string {
	return `{
    "cid":"QmAsset",
    "content_hash":"sha256",
    "size":120
}`
}
func GetAssetRoot() string {
	return `{
    "items":[
        {"/":"cid-asset-#15"}
    ]
}`
}
func GetShareEntries() string {
	return `{
    "id":"share1",
    "asset_cid":"QmAsset",
    "trust_group_id":"tg1"
}`
}
func GetShareEntriesRoot() string {
	return `{
    "items":[
        {"/":"cid-shareentry-#87"}
    ]
}`
}
func GetWorkspaces() string {
	return `{
    "id":"cid-workspace-#18",
    "name":"Mecha_design",
    "status":"active",
    "owner_id":"user_#64",
	"vault_id": "vault_user_0039"
}`
}
func GetWorkspacesRoot() string {
	return `{
    "items":[
        {"/":"cid-workspace-#18"}
    ]
}`
}
func GetChannel() string {
	return `{
		"type":"channel",
		"version":"1.0.0",
		"id":"channel#69",
		"workspace_id":"cid-workspace-#18",
		"template_id":"template#1",
		"title":"Mecha_design",
		"participants":"cid-participants-root",
		"status":"active",
		"policy_cid":{
			"/":"cid-empty"
		},
		"slots":{
			"/":"cid-slots-root"
		},
		"assignments":{
			"/":"cid-empty"
		},
		"federation_cid":{
			"/":"cid-empty"
		}
	}`
}

func GetChannelsRoot() string {
	return `{
		"items":[
			{"/":"cid-channel-#69"}
		]
	}`
}

func GetSlot(id string) string {
	if id == "23" {
		return `{
			"id":"cid-slot-#23",
			"name":"engineering_design",
			"role":"member",
			"vault_id":"user_#64",
			"gated":false,
			"order":1
		}`
	}
	if id == "26" {
		return `{
			"id":"cid-slot-#26",
			"name":"financial_review",
			"role":"finance",
			"vault_id":"user_#4",
			"gated":true,
			"order":2
		}`
	}
	return ""
}
func GetSlotsRoot() string {
	return `{
		"items":[
			{"/":"cid-slot-#23"},
			{"/":"cid-slot-#26"}
		]
	}`
}

// ------------------------------------------------------------------------------------------------------------
// TEST
// ------------------------------------------------------------------------------------------------------------
func TestVaultReconstructor_BuildFromRoot(t *testing.T) {

	ctx := context.Background()

	// -----------------------------
	// Mock Query
	// -----------------------------
	mock := &mockQueryExecutor{
		data: map[string][]byte{},
	}

	personalRoot := `{
		"type":"personal_vault",
		"version":"1.0.0",
		"entries":{
			"/":"cid-entries-root"
		},
		"folders":{
			"/":"cid-folders-root"
		},
		"attachments":{
			"/":"cid-attachments-root"
		},
		"index":{
			"/":"cid-index"
		}
	}`

	collaborativeRoot := `{
		"type":"collaborative_vault",
		"version":"1.0.0",
		"workspaces":{
			"\/":"cid-workspaces-root"
		},
		"participants":{"\/":"cid-empty"},
		"channels":{
			"\/":"cid-channels-root"
		},
		"threads":{
			"\/":"cid-threads-root"
		},
		"slots":{
			"\/":"cid-slots-root"
		},
		"share_entries":{"\/":"cid-empty"},
		"assets":{
			"\/":"cid-assets-root"
		},
		"trust_groups":{
			"\/":"cid-trustgroups-root"
		},
		"trust_members":{"\/":"cid-empty"},
		"federation":{"\/":"cid-empty"},
		"index":{"\/":"cid-empty"}
	}`

	// ---- Entry node (login)
	entry := `{
		"type": "login",
		"entry_name": "gmail",
		"username": "user",
		"password": "pass"
	}`

	// ---- EntriesRoot
	entriesRoot := `{
		"items": [
			{"/": "cid-entry-1"}
		]
	}`

	attachment := `{
		"cid":"Qm123",
		"filename":"contract.pdf"
	}`
	attachmentsRoot := `{
		"items":[
			{"/":"cid-attachment-1"}
		]
	}`

	// ---- FoldersRoot
	folder := `{
		"id":"folder1",
		"name":"Work"
	}`
	foldersRoot := `{
		"items":[
			{"/":"cid-folder-1"}
		]
	}`
	// ---- indexsRoot
	index := `{
		"items": []
	}`

	// ---- Root VaultNode
	root := `{
		"type":"vault",
		"version":"1.0.0",
		"personal":{
			"/":"cid-personal-root"
		},
		"collaborative":{
			"/":"cid-collaborative-root"
		}
	}`

	trustGroup := GetTrustGroup()
	trustGroupsRoot := GetTrustGroupRoot()
	mock.AddData("cid-trustgroup-#1", trustGroup)
	mock.AddData("cid-trustgroups-root", trustGroupsRoot)

	thread := GetThread()
	threadRoot := GetThreadRoot()
	mock.AddData("cid-thread-#1", thread)
	mock.AddData("cid-threads-root", threadRoot)

	asset := GetAssets()
	assetRoot := GetAssetRoot()
	mock.AddData("cid-asset-#15", asset)
	mock.AddData("cid-assets-root", assetRoot)

	workspace := GetWorkspaces()
	workspaceRoot := GetWorkspacesRoot()
	mock.AddData("cid-workspace-#18", workspace)
	mock.AddData("cid-workspaces-root", workspaceRoot)

	// shareEntries := GetShareEntries()
	// shareEntriesRoot := GetThreadRoot()
	// mock.AddData("cid-shareentry-#87", shareEntries)
	// mock.AddData("cid-shareentries-root", shareEntriesRoot)

	// participants := GetParticipants()
	// participantsRoot := GetParticipantRoot()
	// mock.AddData("cid-participant-#33", participants)
	// mock.AddData("cid-participants-root", participantsRoot)

	trustGroupMember := GetTrustGroupMembers()
	trustGroupMembersRoot := GetTrustGroupMembersRoot()
	mock.AddData("cid-trustgroupmember-#1", trustGroupMember)
	mock.AddData("cid-trustgroupmembers-root", trustGroupMembersRoot)

	channel := GetChannel()
	channelsRoot := GetChannelsRoot()
	mock.AddData("cid-channel-#69", channel)
	mock.AddData("cid-channels-root", channelsRoot)

	slot23 := GetSlot("23")
	slot26 := GetSlot("26")
	slotsRoot := GetSlotsRoot()
	mock.AddData("cid-slot-#23", slot23)
	mock.AddData("cid-slot-#26", slot26)
	mock.AddData("cid-slots-root", slotsRoot)
	// -----------------------------
	// Fill mock storage
	// -----------------------------
	mock.data["rootCID"] = []byte(root)

	mock.data["cid-personal-root"] = []byte(personalRoot)
	mock.data["cid-collaborative-root"] = []byte(collaborativeRoot)

	mock.data["cid-entries-root"] = []byte(entriesRoot)
	mock.data["cid-folders-root"] = []byte(foldersRoot)
	mock.data["cid-attachments-root"] = []byte(attachmentsRoot)
	mock.data["cid-entry-1"] = []byte(entry)
	mock.data["cid-index"] = []byte(index)
	mock.data["cid-folder-1"] = []byte(folder)
	mock.data["cid-attachment-1"] = []byte(attachment)

	mock.data["cid-empty"] = []byte(`{"items":[]}`)
	// -----------------------------
	// Reconstructor
	// -----------------------------
	reconstructor := &vaults_service.VaultReconstructor{Query: mock}

	cmd := vault_queries.GetIPFSDataQuerry{CID: "rootCID"}

	// -----------------------------
	// Act
	// -----------------------------
	res, err := reconstructor.BuildFromRoot(ctx, cmd)
	utils.LogPretty("TestVaultReconstructor_BuildFromRoot - res", res)

	// -----------------------------
	// Assert
	// -----------------------------
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// - Personal / Entries
	if res.Version != "1.0.0" {
		t.Fatalf("expected version 1.0.0")
	}

	if len(res.Personal.Entries.Login) != 1 {
		t.Fatalf("expected 1 login entry, got %d", len(res.Personal.Entries.Login))
	}

	if res.Personal.Entries.Login[0].EntryName != "gmail" {
		t.Fatalf("wrong entry content")
	}

	// - Personal / Folders + Attachments
	if len(res.Personal.Folders) != 1 {
		t.Fatalf("expected 1 folder, got %d", len(res.Personal.Folders))
	}

	if res.Personal.Folders[0].Name != "Work" {
		t.Fatalf("expected folder name, got %s", res.Personal.Folders[0].Name)
	}

	if len(res.Personal.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(res.Personal.Attachments))
	}

	// - Collaborative / Trust Groups
	if len(res.Collaborative.TrustGroups) != 1 {
		t.Fatalf("expected 1 trustgroup, got %d", len(res.Collaborative.TrustGroups))
	}

	if res.Collaborative.TrustGroups[0].Name != "Engineering" {
		t.Fatalf("expected trustgroup name, got %s", res.Collaborative.TrustGroups[0].Name)
	}

	// - Collaborative / Channels
	if len(res.Collaborative.Channels) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(res.Collaborative.Channels))
	}

	ch := res.Collaborative.Channels[0]

	if ch.ID != "channel#69" {
		t.Fatalf("expected channel id channel#69, got %s", ch.ID)
	}

	if ch.WorkspaceID != "cid-workspace-#18" {
		t.Fatalf("expected cid-workspace-#18, got %s", ch.WorkspaceID)
	}

	if ch.Title != "Mecha_design" {
		t.Fatalf("expected title Mecha_design, got %s", ch.Title)
	}

	if ch.Status != vaults_domain.ChannelStatus("active") {
		t.Fatalf("expected status active, got %v", ch.Status)
	}

	if len(ch.Slots) != 2 {
		t.Fatalf("expected 2 slots")
	}

	if ch.Slots[0].Name != "engineering_design" {
		t.Fatalf("wrong first slot")
	}

	if ch.Slots[1].Name != "financial_review" {
		t.Fatalf("wrong second slot")
	}

	// // - Participants + Trust Members
	// if len(res.Collaborative.Participants) == 1 {
	// 	t.Fatalf("expected 1 Participants, got %d", len(res.Collaborative.Participants))
	// }
	// if len(res.Collaborative.TrustGroupMembers) == 1 {
	// 	t.Fatalf("expected 1 TrustGroupMembers, got %d", len(res.Collaborative.TrustGroupMembers))
	// }
	// if res.Collaborative.Participants[0].VaultID == "vaultA" {
	// 	t.Fatalf("expected Participants vaultID, got %s", res.Collaborative.Participants[0].VaultID)
	// }
	// if res.Collaborative.TrustGroupMembers[0].Role == "owner" {
	// 	t.Fatalf("expected TrustGroupMembers vaultID, got %s", res.Collaborative.TrustGroupMembers[0].Role)
	// }

	// // - Assets + ShareEntries + Threads
	// if len(res.Collaborative.Assets) == 1 {
	// 	t.Fatalf("expected 1 Assets, got %d", len(res.Collaborative.Assets))
	// }
	// if len(res.Collaborative.ShareEntries) == 1 {
	// 	t.Fatalf("expected 1 ShareEntries, got %d", len(res.Collaborative.ShareEntries))
	// }
	// if len(res.Collaborative.Threads) == 1 {
	// 	t.Fatalf("expected 1 Threads, got %d", len(res.Collaborative.Threads))
	// }
	// if res.Collaborative.ShareEntries[0].AssetCID == res.Collaborative.Assets[0].CID {
	// 	t.Fatalf("expected 1 same ShareEntries, got %s", res.Collaborative.ShareEntries[0].AssetCID)
	// }

	// - Full Vault
	// Personal
	// len(res.Personal.Entries.Login) == 3
	// len(res.Personal.Folders) == 2
	// len(res.Personal.Attachments) == 5

	// // Collaborative
	// len(res.Collaborative.Workspaces) == 2
	// len(res.Collaborative.Channels) == 4
	// len(res.Collaborative.Threads) == 8
	// len(res.Collaborative.Assets) == 12
	// len(res.Collaborative.ShareEntries) == 12
	// len(res.Collaborative.TrustGroups) == 2
	// len(res.Collaborative.TrustGroupMembers) == 5
	// len(res.Collaborative.Participants) == 5
}
func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestVault_EndToEnd_Commit_Then_Reconstruct(t *testing.T) {
	ctx := context.Background()

	attachmentHash := "f648264d19025c76671a970cdf3196cbf16698a94817ebc3ccf5bb2dba35ec7a"

	attachmentNodeCID := "node-att-1"
	attachment := vaults_domain.Attachment{
		ID:      "att-1",
		NodeCID: attachmentNodeCID,
		Hash:    attachmentHash,
		Name:    "background_medium.jpg",

		RecipientCIDs: map[string]string{},
	}

	// -------------------------------------------------
	// 🔥 SINGLE SHARED STORAGE (source of truth)
	// -------------------------------------------------
	store := make(map[string][]byte)

	memoryIPFS := &memoryIPFS{
		store: store,
	}

	// -------------------------------------------------
	// Crypto
	// -------------------------------------------------
	crypto := &vault_infrastructure_crypto.AESService{}

	vaultKey := make([]byte, 32)
	_, _ = rand.Read(vaultKey)

	unlock := &mockUnlockHandler{key: vaultKey}

	// -------------------------------------------------
	// Mock storage provider (WRITE path uses SAME store)
	// -------------------------------------------------
	mockStorage := &mockStorageProvider{
		AddFunc: func(ctx context.Context, data []byte) (string, error) {
			cid := fmt.Sprintf("cid-%d", len(store)+1)
			store[cid] = data
			return cid, nil
		},
		GetFunc: func(ctx context.Context, cid string) ([]byte, error) {
			data, ok := store[cid]
			if !ok {
				return nil, fmt.Errorf("not found: %s", cid)
			}
			return data, nil
		},
	}

	mockFactory := &mockStorageFactory{
		NewFunc: func(ctx *app_config_domain.VaultContext) app_config.StorageProvider {
			return mockStorage
		},
	}

	// -------------------------------------------------
	// IPFS write handler (commit side)
	// -------------------------------------------------
	ipfsCreateHandler := &vault_commands.CreateIPFSPayloadCommandHandler{
		UnlockVaultHandler: unlock,
		CryptoService:      &vault_infrastructure_crypto.AESService{},
		StorageFactory:     mockFactory,
	}

	// -------------------------------------------------
	// Vault service (WRITE)
	// -------------------------------------------------
	userID := "user_1"
	userPassword := "password"
	vaultName := "ocean"
	cfgs, err := GetConfig(userID, vaultName)
	if err != nil {
		utils.LogPretty("Vault service (WRITE) error", err)
	}

	vc := app_config_domain.VaultContext{
		Configs:       *cfgs,
		UserID:        userID,
		VaultName:     vaultName,
		StorageConfig: cfgs.App.Storage,
	}

	nr := &MockNodeRepo{
		Entries: vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry1",
						Type:      "login",
						EntryName: "GitHub",
						IsDraft:   false,
					},
				},
			},
			Card: []vaults_domain.CardEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry3",
						Type:      "card",
						EntryName: "Visa",
						IsDraft:   false,
					},
				},
			},
		},
		Attachements: []vaults_domain.Attachment{
			attachment,
		},
	}

	service := &vaults_service.VaultService{
		VaultCtx:    vc,
		Encryptor:   &vaults_service.AESEncryptor{},
		Repo:        &fakeVaultRepo{},
		NodeRepo:    nr,
		Password:    userPassword,
		IPFSHandler: ipfsCreateHandler,
		IsDraftMode: false, // IMPORTANT: full flow
	}

	vp := fakeVaultPayload(userID, vaultName)
	session := GetSession(userID, vp)

	// -------------------------------------------------
	// ACT 1: Commit
	// -------------------------------------------------
	rootCID, _, _, _, err := service.CommitVault(session, vaults_service.FullSync)
	if err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	if rootCID == "" {
		t.Fatalf("empty root CID")
	}

	// -------------------------------------------------
	// Query (READ) — IMPORTANT: SAME STORAGE
	// -------------------------------------------------
	query := &e2eQuery{
		ipfs:   memoryIPFS, // 🔥 SAME BACKING STORE
		crypto: crypto,
		key:    vaultKey,
	}

	reconstructor := &vaults_service.VaultReconstructor{
		Query: query,
	}

	cmd := vault_queries.GetIPFSDataQuerry{
		CID: rootCID,
	}

	// -------------------------------------------------
	// ACT 2: Reconstruct
	// -------------------------------------------------
	vaultPayload, err := reconstructor.BuildFromRoot(ctx, cmd)
	if err != nil {
		t.Fatalf("reconstruct failed: %v", err)
	}
	utils.LogPretty("vault", vaultPayload)

	// -------------------------------------------------
	// ASSERT
	// -------------------------------------------------
	if vaultPayload.Version != "1.0.0" {
		t.Fatalf("wrong version: %s", vaultPayload.Version)
	}

	if len(vaultPayload.Personal.Entries.Login) == 0 {
		t.Fatalf("entries not reconstructed")
	}

	t.Logf("✅ E2E success: rootCID=%s", rootCID)
}
