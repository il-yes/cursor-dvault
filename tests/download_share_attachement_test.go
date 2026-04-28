package tests

import (
	"context"
	"os"
	"testing"
	"vault-app/internal/auth"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/driver"
	"vault-app/internal/logger/logger"
	"vault-app/internal/registry"
	"vault-app/internal/tracecore"
	"vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_ui "vault-app/internal/vault/ui"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --------------------------------------------------------------------------------------------------
// Mocks
// --------------------------------------------------------------------------------------------------

// ============= MockVaultRepo =======================================================
type MockVaultRepo struct {
	Vault                 *vaults_domain.Vault
	existingVault         *vaults_domain.Vault
	saveCalled            bool
	saveError             error
	savedVault            *vaults_domain.Vault
	updateCalled          bool
	updateError           error
	deleteCalled          bool
	deleteError           error
	saveFn                func(v *vaults_domain.Vault) error
	updateFn              func(v *vaults_domain.Vault) error
	DefaultVaultNameTest  bool
	isCreate              bool
	GetLatestByUserIDFunc func(userID string) (*vaults_domain.Vault, error)

}

func (f *MockVaultRepo) CreateVault(vault *vaults_domain.Vault) error {
	return nil
}

func (m *MockVaultRepo) SaveVault(v *vaults_domain.Vault) error {
	if m.saveFn != nil {
		m.Vault = v
		return m.saveFn(v)
	}
	if m.DefaultVaultNameTest {
		utils.LogPretty("MockVaultRepo - SaveVault - Vault", v)
		m.Vault = v
		return nil
	}
	m.existingVault = v
	return nil
}

func (f *MockVaultRepo) DeleteVault(string) error {
	f.deleteCalled = true
	return f.deleteError
}
func (f *MockVaultRepo) GetByUserIDAndName(string, string) (*vaults_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vaults_domain.ErrVaultNotFound
}
func (f *MockVaultRepo) GetLatestByUserID(userID string) (*vaults_domain.Vault, error) {
	if userID == "user-1" && !f.isCreate {
		return f.existingVault, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (f *MockVaultRepo) GetVault(string) (*vaults_domain.Vault, error) {
	panic("not used")
}
func (f *MockVaultRepo) UpdateVault(v *vaults_domain.Vault) error {
	f.updateCalled = true
	// return f.updateError
	if f.updateFn != nil {
		return f.updateFn(v)
	}
	return nil
}
func (f *MockVaultRepo) UpdateVaultCID(vaultID, cid string) error {
	f.updateCalled = true
	return f.updateError
}

// ============= MockAppConfigHandler =======================================================
type MockAppConfigHandler struct {
	InitAppConfigFunc         func(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error)
	InitUserConfigFunc        func(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error)
	GetAppConfigByUserIDFunc  func(ctx context.Context, userID string) (*app_config_commands.CreateAppConfigCommandOutput, error)
	GetUserConfigByUserIDFunc func(userID string) (*app_config_domain.UserConfig, error)
}

func (f *MockAppConfigHandler) InitAppConfig(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error) {
	if f.InitAppConfigFunc != nil {
		return f.InitAppConfigFunc(input)
	}
	return &app_config_commands.CreateAppConfigCommandOutput{
		AppConfig: &app_config_domain.AppConfig{},
	}, nil
}

func (f *MockAppConfigHandler) InitUserConfig(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error) {
	if f.InitUserConfigFunc != nil {
		return f.InitUserConfigFunc(input)
	}
	return &app_config_commands.CreateUserConfigCommandOutput{
		UserConfig: &app_config_domain.UserConfig{},
	}, nil
}

func (f *MockAppConfigHandler) GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error) {
	if f.GetAppConfigByUserIDFunc != nil {
		res, err := f.GetAppConfigByUserIDFunc(ctx, userID)
		if err != nil {
			return nil, err
		}
		return res.AppConfig, nil
	}
	return &app_config_domain.AppConfig{}, nil
}

func (f *MockAppConfigHandler) GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error) {
	if f.GetUserConfigByUserIDFunc != nil {
		return f.GetUserConfigByUserIDFunc(userID)
	}
	return &app_config_domain.UserConfig{}, nil
}

// ============= MockSymKeyDecryptor =======================================================
type MockSymKeyDecryptor struct {
	DecryptPasswordWithStellarByteFunc func(encNonce, encPassword, privateKey string) ([]byte, error)
}

func (f *MockSymKeyDecryptor) DecryptPasswordWithStellarByte(encNonce, encPassword, privateKey string) ([]byte, error) {
	return f.DecryptPasswordWithStellarByteFunc(encNonce, encPassword, privateKey)
}

// ============= MockVaultHandler =======================================================
type MockVaultHandler struct {
	VaultRepository          vaults_domain.VaultRepository
	DownloadAttachmentFunc   func(ctx context.Context, req vault_ui.DownloadAttachmentRequest) (string, error)
	GetIPFSDataQuerryHandler *MockGetIPFSDataQueryHandler
}

func (v *MockVaultHandler) DownloadAttachment(ctx context.Context, req vault_ui.DownloadAttachmentRequest) (string, error) {
	if v.DownloadAttachmentFunc != nil {
		return v.DownloadAttachmentFunc(ctx, req)
	}
	return "", nil
}

// ============= MockGetIPFSDataQueryHandler =======================================================
type MockGetIPFSDataQueryHandler struct {
	ExecuteFunc func(ctx context.Context, query vault_queries.GetIPFSDataQuerry) (*vault_queries.GetIPFSDataResponse, error)
	AESDecryptFunc func( []byte, []byte) ([]byte, error)
}

func (f *MockGetIPFSDataQueryHandler) Execute(ctx context.Context, query vault_queries.GetIPFSDataQuerry) (*vault_queries.GetIPFSDataResponse, error) {
	return f.ExecuteFunc(ctx, query)
}
func (f *MockGetIPFSDataQueryHandler) AESDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	return f.AESDecryptFunc(encrypted, key)
}

// ============= mockIPFS =======================================================
type mockIPFS struct {
	cid  string
	data []byte
}

func (m *mockIPFS) StoreOnIpfs(ctx interface{}, params interface{}) (string, error) {
	return m.cid, nil
}
func (m *mockIPFS) Add(context.Context, []byte) (string, error) {
	return m.cid, nil
}

func GetFakeSession() vault_session.Session {
	return vault_session.Session{
		UserID:      "c5f83bd9-8a2a-40e9-a590-e7bf3ccc8859",
		Vault:       []byte("eyJ2ZXJzaW9uIjoiIiwibmFtZSI6IiIsImZvbGRlcnMiOltdLCJlbnRyaWVzIjp7ImxvZ2luIjpbXSwiY2FyZCI6W10sImlkZW50aXR5IjpbXSwibm90ZSI6W3siaWQiOiIwZDJjYzAxOC02Y2ViLTRmZGUtYjNkMS1mODNiNTllMmE1NzUiLCJlbnRyeV9uYW1lIjoiYW1hem9uIGIiLCJmb2xkZXJfaWQiOiIiLCJ0eXBlIjoibm90ZSIsImFkZGl0aW9ubmFsX25vdGUiOiJOb3Rl4oCmIHN0b3J5LCBlcGljIiwidHJhc2hlZCI6ZmFsc2UsImlzX2RyYWZ0Ijp0cnVlLCJpc19kaXJ0eSI6dHJ1ZSwiaXNfZmF2b3JpdGUiOmZhbHNlLCJjcmVhdGVkX2F0IjoiIiwidXBkYXRlZF9hdCI6IiIsImF0dGFjaG1lbnRzIjpbeyJpZCI6IjZiNWY3MGY0LTY3ZWUtNDZmOS05NTUwLTA4ZTFiMzJhODM4YSIsImVudHJ5X2lkIjoiMGQyY2MwMTgtNmNlYi00ZmRlLWIzZDEtZjgzYjU5ZTJhNTc1IiwiaGFzaCI6ImJhM2FjNDg0NTIzYTUwYjM0NjEwMjk5ZmNkNzFjZjEyNWRkOWJjZGQ2Yzk2OGRkYjZlZjM4NjljOGIxZjUwYWQiLCJuYW1lIjoiNjQ5MDgxMjUtZTI5ZS00NzllLWJiODUtM2VhM2I4ZjUxZGJjLmpwZWciLCJzaXplIjo0MjkzOTUsImNpZCI6IlFtUXJScDNyYVF0TUprdWQ1N2l5NW1yQWIyTXBwMVp0eldSOThqaXJybUxCVzgiLCJzdG9yYWdlIjoiaXBmcyIsImV4dCI6ImpwZWciLCJkb3dubG9hZGVkX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJoYXNoX2xvY2FsIjoiIiwiaGFzaF9zaGFyZSI6IiIsInJlY2lwaWVudF9jaWRzIjpudWxsfV19LHsiaWQiOiJkMzA1OGFkMS1kYThiLTRlZGYtYTZiMi1jODE3NjUxZWNkNjAiLCJlbnRyeV9uYW1lIjoiR2l0SHViIiwiZm9sZGVyX2lkIjoiIiwidHlwZSI6Im5vdGUiLCJhZGRpdGlvbm5hbF9ub3RlIjoiTm90ZSIsInRyYXNoZWQiOmZhbHNlLCJpc19kcmFmdCI6dHJ1ZSwiaXNfZGlydHkiOnRydWUsImlzX2Zhdm9yaXRlIjpmYWxzZSwiY3JlYXRlZF9hdCI6IjIwMjYtMDQtMjFUMTc6NDE6NDMuNTg3WiIsInVwZGF0ZWRfYXQiOiIyMDI2LTA0LTIxVDEwOjQyOjAwLTA3OjAwIiwiYXR0YWNobWVudHMiOlt7ImlkIjoiNGEwZThmNDUtNjlmNS00OTI0LWIzYjktMDYxYTM3Yzg2NTA1IiwiZW50cnlfaWQiOiJkMzA1OGFkMS1kYThiLTRlZGYtYTZiMi1jODE3NjUxZWNkNjAiLCJoYXNoIjoiZGIyOWYwZDRhNzk1ZTViNmUwYTc5Y2QwMWYyYzAxOGQwNDkxZTcwMDhmMTY0MzQ3MzBhZjQ3ZDI4Zjg4NzliYSIsIm5hbWUiOiJiYXNoZW5nYS5qcGciLCJzaXplIjozOTkxMCwiY2lkIjoiUW1haDFGQnFKTGpmdGVqVXhCM2E0Wkd6WndobWFGZmtOWWlNc2RHcVdTTW5nViIsInN0b3JhZ2UiOiJpcGZzIiwiZXh0IjoianBnIiwiZG93bmxvYWRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwiaGFzaF9sb2NhbCI6IiIsImhhc2hfc2hhcmUiOiIiLCJyZWNpcGllbnRfY2lkcyI6bnVsbH1dfSx7ImlkIjoiN2ZjMjkxMDktYjEzMS00OGU4LWFiOTYtODNkMDRkZTYyMmY0IiwiZW50cnlfbmFtZSI6IlJlZGRpdCIsImZvbGRlcl9pZCI6IiIsInR5cGUiOiJub3RlIiwiYWRkaXRpb25uYWxfbm90ZSI6IlJlZGl0IG5vdGUiLCJ0cmFzaGVkIjpmYWxzZSwiaXNfZHJhZnQiOnRydWUsImlzX2RpcnR5Ijp0cnVlLCJpc19mYXZvcml0ZSI6ZmFsc2UsImNyZWF0ZWRfYXQiOiIyMDI2LTA0LTIxVDE5OjAwOjUzLjEzNFoiLCJ1cGRhdGVkX2F0IjoiMjAyNi0wNC0yMVQxMjowMToxMS0wNzowMCIsImF0dGFjaG1lbnRzIjpbeyJpZCI6IjY3MzhhNGYyLThlYTItNDdiMS05MDc5LTJiNTNjOTI0NThhYiIsImVudHJ5X2lkIjoiN2ZjMjkxMDktYjEzMS00OGU4LWFiOTYtODNkMDRkZTYyMmY0IiwiaGFzaCI6IjU4MTIxNWFjYmM0YTJmYzhiZTEyMDQyODFjMjUyMTI2NDA2MDBjMWJhNjQ3MWZjOGE0Y2Q3Y2E0YzE0NGJiYWMiLCJuYW1lIjoiYXZhdGFyNS5qcGciLCJzaXplIjo4MjE5LCJleHQiOiJqcGciLCJkb3dubG9hZGVkX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJoYXNoX2xvY2FsIjoiIiwiaGFzaF9zaGFyZSI6IiIsInJlY2lwaWVudF9jaWRzIjpudWxsfV19LHsiaWQiOiJhMTNhMzk5MC04MzU3LTQzMTEtODRkZC04ZGQ3YjU2MTYyMWEiLCJlbnRyeV9uYW1lIjoiQ2xvdWQgZmFyZSIsImZvbGRlcl9pZCI6IiIsInR5cGUiOiJub3RlIiwiYWRkaXRpb25uYWxfbm90ZSI6Ik5vdGUiLCJ0cmFzaGVkIjpmYWxzZSwiaXNfZHJhZnQiOnRydWUsImlzX2RpcnR5Ijp0cnVlLCJpc19mYXZvcml0ZSI6ZmFsc2UsImNyZWF0ZWRfYXQiOiIyMDI2LTA0LTIxVDIwOjUxOjE3LjcwMloiLCJ1cGRhdGVkX2F0IjoiMjAyNi0wNC0yMVQxMzo1MTo0MS0wNzowMCIsImF0dGFjaG1lbnRzIjpbeyJpZCI6ImFlODE0NjBkLTY1MTItNDgwYy04MGIwLTVjZDNlYzE3MTU4NiIsImVudHJ5X2lkIjoiYTEzYTM5OTAtODM1Ny00MzExLTg0ZGQtOGRkN2I1NjE2MjFhIiwiaGFzaCI6ImU5NTVjNWMxNzE5OTUzZTViZjM3YjI4NWFmZDdkNmQ1NTQxOTJhNTk4NzQ2OTU5NzM1MjdhZGNlMTM3YmRkNmMiLCJuYW1lIjoiTGVnYWwgU2xpZGUgMTEgLSBPdXRjb21lIENUQSAoMSkucG5nIiwic2l6ZSI6ODAzNDc4LCJleHQiOiJwbmciLCJkb3dubG9hZGVkX2F0IjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJoYXNoX2xvY2FsIjoiIiwiaGFzaF9zaGFyZSI6IiIsInJlY2lwaWVudF9jaWRzIjpudWxsfV19XSwic3Noa2V5IjpbXX0sImNyZWF0ZWRfYXQiOiIiLCJ1cGRhdGVkX2F0IjoiIn0="),
		LastCID:     "QmeiRSzJCXYcNPgbR3e5Kj8zaQxnPomkHgD33emCwpH2D9",
		LastSynced:  "2026-04-21T10:40:26-07:00",
		LastUpdated: "2026-04-22T22:14:07-07:00",
		Runtime: &vault_session.RuntimeContext{
			AppConfig: *mockAppConfig(),
		},
	}
}

// ============= TestApp =======================================================
type TestApp struct {
	NowUTC                func() string
	VaultKeyringDecryptor *MockSymKeyDecryptor
	Logger                *logger.Logger

	// Core handlers
	AppConfigHandler *MockAppConfigHandler
	Vault            *MockVaultHandler

	// New: Global state
	RuntimeContext  *vault_session.RuntimeContext
	cancel          context.CancelFunc
	RequireAuthFunc func(token string) (*auth.Claims, error)
}

func NewTestApp() *TestApp {
	return &TestApp{
		AppConfigHandler: &MockAppConfigHandler{},
		Vault:            &MockVaultHandler{},
	}
}

func (a *TestApp) RequireAuth(token string) (*auth.Claims, error) {
	return a.RequireAuthFunc(token)
}

func (a *TestApp) DownloadShareAttachement(req vault_dto.DownloadShareAttachmentRequest, jwtToken string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	userConfig, err := a.AppConfigHandler.GetUserConfigByUserID(claims.UserID)
	if err != nil {
		return "", err
	}

	stellarAccount := userConfig.StellarAccount
	symKey, err := a.VaultKeyringDecryptor.DecryptPasswordWithStellarByte(
		string(stellarAccount.EncNonce),
		string(stellarAccount.EncPassword),
		stellarAccount.PrivateKey,
	)
	if err != nil {
		return "", err
	}

	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		return "", err
	}

	return a.Vault.DownloadAttachment(context.Background(), vault_ui.DownloadAttachmentRequest{
		UserID:       claims.UserID,
		Vault:        *vault,
		CID:          req.AttachmentCID,
		Ext:          req.FileExtension,
		PrivateKey:   stellarAccount.PrivateKey,
		EncryptedKey: req.EncryptedKey,
		SymKey:       symKey,
	})
}

func mockAppConfig() *app_config_domain.AppConfig {
	return &app_config_domain.AppConfig{
		ID:                 "910d0b25-cddd-479b-a5f5-0fe74c6e0269",
		RepoID:             "my-repo-id",
		Branch:             "main",
		TracecoreEnabled:   false,
		CommitRules:        nil,
		BranchingModel:     "",
		EncryptionPolicy:   "",
		Actors:             nil,
		FederatedProviders: nil,
		DefaultPhase:       "vault_entry",
		DefaultVaultPath:   "",
		VaultSettings: app_config_domain.VaultConfig{
			MaxEntries:       3,
			AutoSyncEnabled:  true,
			EncryptionScheme: "AES-256-GCM",
		},
		Blockchain: app_config_domain.BlockchainConfig{
			Stellar: app_config_domain.StellarConfig{
				Network:       "testnet",
				HorizonURL:    "https://horizon-testnet.stellar.org",
				Fee:           100,
				SyncFrequency: "",
			},
			IPFS: app_config_domain.IPFSConfig{
				APIEndpoint: "http://localhost:5001",
				GatewayURL:  "https://ipfs.io/ipfs/",
			},
		},
		UserID:               "c5f83bd9-8a2a-40e9-a590-e7bf3ccc8859",
		AutoLockTimeout:      "",
		AccessPolicyDuration: 0,
		RemaskDelay:          "",
		Theme:                "",
		AnimationsEnabled:    false,
		Storage: app_config.StorageConfig{
			Mode: "cloud",
			LocalIPFS: app_config.IPFSConfig{
				APIEndpoint: "http://localhost:5001",
				GatewayURL:  "https://ipfs.io/ipfs/",
			},
			PrivateIPFS: app_config.IPFSConfig{
				APIEndpoint: "http://192.168.1.10:5001",
				GatewayURL:  "http://192.168.1.10:8080/ipfs/",
			},
			Cloud: app_config.CloudConfig{
				BaseURL: "http://localhost:4001/api",
			},
			EnterpriseS3: app_config.S3Config{
				Region:   "us-east-1",
				Bucket:   "ankhora-enterprise",
				Endpoint: "https://s3.us-east-1.amazonaws.com",
			},
		},
	}
}

// ------------------------------------------------------------------------
//
//	TracecoreClient - AccessEncryptedEntry - cloud response
//
// ------------------------------------------------------------------------
func mockAccessEncryptedEntryResponse() map[string]any {
	return map[string]any{
		"status": 200,
		"data": map[string]any{
			"EncryptedKey":     "rmdHW9H7+uCWi8UWMQBDOb92+l5CHhYdx6BOE4V41VPCAq+LoBIXezPALrVgEAd6DS68W0QiFAw39zoMWMYbM5fIeuhdBffV+EGgLZLakb8=",
			"SenderPublicKey":  "GC64CEHLFIV4DRIRUHCQZ6WJ26HCDG5IOZSCK2VXBMKDWL7K52FCGX5X",
			"EncryptedPayload": "j0NWcIuO4RVMzGeZNEpNex60wFerqW+vRGr2210nnjQT6LQguR3ZdaeJYV4H6Wk4x1Yf04/FS+ueMZnwGQz9XZnVLwhDyg8EMR5bJoVmDQ0JZgxdd/axlCAvMEm5Fbz6PhwLAgNyzdmCtHMnyd1wFiVq8L2X6kXIejRSdL+UCIszrB54/gCASTvcrJY+xw9xjc8SPNcVxJcxvhNJmWRxepopKPlIL2/uuceXNGBy5JUi7W52bU5gzXLO0wHjxT/ji3ytu7nYk+LroEg2Yl9rg/XW6gZ/jOgw+MEh62bNOsTxPJPZtnSL2hnAVTvDTyeKIVI+QKK4giKiY3r3KrLbIRbTUmFP31oPPXoVeoRxQxUq58/owYdau0Odit4LJm0qiVMjhr685fsPM/P/BtZ7kURPvIbyfIe933oKk1Nm+m+dy77PBIGL/qfli5vZ3x5Xi+YT8Whn/IilIPrqHBYAx/2MuXYgOEDaqrBlBt2b5AbKcuHkrYTRO1VybIuFqWmlJPtJxSWYt69qAdjoCN+U98SctgFpkoh/XvpFKDRxt8quHCrbpkLEqKc6VNnLP04ukv5D0kfDgx/X+DDBLns14BqS8YecCS3OdceMOiU4BG9l07lv/uXft5MRLCCTLYBNPnDM01UlJ2JLbgjlDmLRzy0gdSMQgO2iS3nT+tJk1JUtVVhv7d3M0hkLn7MgmPTyb999ocbb84aSjXuQKRmeiABOwTabCqtc35FnyIR5BeOy9WitRGu6tNxUHv+grRv43NWnD9jBT8rcVRv7u1tnj4JZqK6Qnj+rSxsiOQkE8xxlAfGG08LB3fw6O1Q/ryOfH3yge+A2No8PvfxMtrjQgGE67NbebR4C/2Lu2k5HFJYz2ONx6knJrEn5xAUnl35cOLCvXAG0oYH4ovHqeUvzRPP46d1FCaOC8sXAF4pHl0eqnacKpGh2drPfo/nUYJjW6knoSZUFPk+sTCzjPl7u+NSXdD2q01lpDGvtbgLlJ4dlZ431k7NSoJ13BrD5KDfP0hzuZbpa+tMnqW6Ynv/T8FWwNYGQtHeBvpO5kfptGZA29mFdD6f+2YskbCj2Ldvn2/ct56UWDDPlMnGH1WTqivP9FHHgiZlVHdjMZNp6TlGVuQ6LEP5KTW7M79IoYuWNfW4culkKXu5wWywkJakgPIZpyPQRGwvfAvX7wHtA7r/Q8Y3xN4IR+wSKbvZLeWkeULRpBUfdx2c3cvxYlQnXSfV7DJsOX1rBHa58nSLwlUT1J40Y4yZc/v28WCpW8iEs+OQJxwTubWiNCG6feOtb9m2jhr+9Zp4DMtGmHpofEctDfla3fFuNtNOGsEyo+p/scmduxHFQOszGZG5HywzNZTnx1gTYBVzii2RZMzA+hqg2LWC65sPEvtupWXoFcxczGFXNpt0yYM16m4C3/z44ukS1WLj4KU/7",
			"DownloadAllowed":  true,
		},
		"message": "share loaded",
	}
}

// --------------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------------

func TestDownloadAttachment_Success_ALPHA(t *testing.T) {
	// given: fake user, vault, and attachment
	userID := "c5f83bd9-8a2a-40e9-a590-e7bf3ccc8859"
	vault := vaults_domain.Vault{
		Name:               "Leeks",
		UserSubscriptionID: "sub-1",
	}
	cid := "QmaGCJpbK3oAnsBFwk6Qs9TSu3fzgwGL9UNuJCxkjHrfq4"
	ext := "png"

	// 1. Mock VaultRepository
	vaultRepo := &MockVaultRepo{}
	vaultRepo.existingVault = &vault

	// 2. Mock GetIPFSDataQueryHandler
	//    (it returns raw bytes, then CloudIPFSStorage.Get decodes)
	var executedQuery vault_queries.GetIPFSDataQuerry
	var executeCalled bool

	// 4. symKey (from DecryptPasswordWithStellarByte, but stubbed here)
	// symKey := []byte("sym-key-32-bytes-..............")
	// require.Equal(t, 32, len(symKey)

	aesService := vault_infrastructure_crypto.AESService{}
	symKey, err := aesService.DecryptPasswordWithStellarByte(
		[]byte("cPu8Oc3uPwOZXq4Q"),
		[]byte("Eo4A8jXUXZ/M9sfGHx95wrI0v3IQaAL9"),
		"SB3RRWLCYSMUZZIVYFM6PPGMXIYLRHWDTGMMGKVYWRW64LK4PJFWRLTF",
	)
	require.Equal(t, 32, len(symKey))

	// getIPFSHandler := &MockGetIPFSDataQueryHandler{
	//     ExecuteFunc: func(ctx context.Context, query vault_queries.GetIPFSDataQuerry) (*vault_queries.GetIPFSDataResponse, error) {
	//         executedQuery = query
	//         executeCalled = true

	//         require.Equal(t, cid, query.CID)
	//         require.Equal(t, symKey, query.SymKey)

	//         // return fake raw bytes
	//         return &vault_queries.GetIPFSDataResponse{
	//             Raw: []byte("fake-ipfs-data"),
	//         }, nil
	//     },
	// }

	// 3. VaultHandler that uses those mocks
	// vaultHandler1 := &MockVaultHandler{
	//     VaultRepository:          vaultRepo,
	//     GetIPFSDataQuerryHandler: getIPFSHandler,
	//     // Logger: your mocked logger
	//     // logger: logger.NewFromEnv(),
	// }
	appLogger := logger.NewFromEnv()

	ipfs := blockchain.NewIPFSClient(os.Getenv("IPFS_CLIENT"))

	db, err := driver.InitDatabase("sqlite3.db", *appLogger)
	if err != nil {
		appLogger.Error("❌ Failed to init DB: %v", err)
		os.Exit(1)
	}
	// -------------------------------------------------------------------------------------------------
	// Tracecore
	// -------------------------------------------------------------------------------------------------
	tracecoreClient := tracecore.NewTracecoreClient(os.Getenv("TRACECORE_URL"), os.Getenv("TRACECORE_TOKEN"), os.Getenv("CLOUD_FRONT_URL"), os.Getenv("CLOUD_BACK_URL"))

	// -------------------------------------------------------------------------------------------------
	// Registry
	// -------------------------------------------------------------------------------------------------
	reg := registry.NewRegistry(appLogger)
	reg.RegisterDefinitions([]registry.EntryDefinition{
		{
			Type:    "login",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.LoginEntry{} },
			Handler: vault_ui.NewLoginHandler(*db, appLogger),
		},
		{
			Type:    "card",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.CardEntry{} },
			Handler: vault_ui.NewCardHandler(*db, appLogger),
		},
		{
			Type:    "note",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.NoteEntry{} },
			Handler: vault_ui.NewNoteHandler(*db, appLogger),
		},
		{
			Type:    "identity",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.IdentityEntry{} },
			Handler: vault_ui.NewIdentityHandler(*db, appLogger),
		},
		{
			Type:    "sshkey",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.SSHKeyEntry{} },
			Handler: vault_ui.NewSSHKeyHandler(*db, appLogger),
		},
	})

	vaultHandler := vault_ui.NewVaultHandler(
		reg,
		*appLogger,
		context.Background(),
		ipfs,
		nil,
		nil,
		*tracecoreClient,
		"",
	)

	// 5. DownloadAttachmentRequest
	downloadReq := vault_ui.DownloadAttachmentRequest{
		UserID:       userID,
		Vault:        vault,
		CID:          cid,
		Ext:          ext,
		Password:     "password",
		PrivateKey:   "private-key",
		EncryptedKey: "encrypted-key",
		SymKey:       symKey,
		AppCfg:       mockAppConfig(),
	}
	// 6. Act
	path, err := vaultHandler.DownloadAttachment(context.Background(), downloadReq)

	// 7. Assert
	require.NoError(t, err)
	require.True(t, executeCalled, "GetIPFSDataQuerryHandler.Execute should be called")
	require.Equal(t, cid, executedQuery.CID)
	require.Equal(t, symKey, executedQuery.SymKey)

	// assert that path is something like ~/Downloads/VaultCore/attachments/attachment_QmCID.jpg
	require.Contains(t, path, "VaultCore/attachments")
	require.Contains(t, path, "attachment_"+cid)
	require.Contains(t, path, ext)

	// and that file exists / can be read (optional, depending on your test mode)
}

func TestDownloadAttachment_Success(t *testing.T) {
    userID := "c5f83bd9-8a2a-40e9-a590-e7bf3ccc8859"
    vaultName := "Leeks"
    cid := "QmaGCJpbK3oAnsBFwk6Qs9TSu3fzgwGL9UNuJCxkjHrfq4"
    ext := "png"
    password := "password"

    // 1. Vault
    vault := vaults_domain.Vault{
        Name:               vaultName,
        UserSubscriptionID: "sub-1",
    }

    // 2. VaultRepo mock
    /*fakeVault := &vaults_domain.Vault{
        ID:                 "vault-1",
        Name:               vaultName,
        UserSubscriptionID: vault.UserSubscriptionID,
    }
    vaultRepo := &MockVaultRepo{
        existingVault: fakeVault,
    }

    // 3. GetIPFSDataQueryHandler mock
    var executedQuery vault_queries.GetIPFSDataQuerry
    var executeCalled bool

    queryHandler := &MockGetIPFSDataQueryHandler{
        ExecuteFunc: func(ctx context.Context, query vault_queries.GetIPFSDataQuerry) (*vault_queries.GetIPFSDataResponse, error) {
            executedQuery = query
            executeCalled = true

            require.Equal(t, cid, query.CID, "CID must match")
            require.Equal(t, password, query.Password, "Password must match")
            require.NotNil(t, query.SymKey, "SymKey must be non‑nil")
            require.Equal(t, 32, len(query.SymKey), "SymKey must be 32 bytes")

            return &vault_queries.GetIPFSDataResponse{
                Raw: []byte("fake-encrypted-attachment-data"),
            }, nil
        },
    }*/

    // 4. Logger
    appLogger := logger.NewFromEnv()

    // 5. DB (if you want; otherwise, VaultRepo is mocked anyway)
    //    (if you keep DB, fine; for UT you can just mock VaultRepo and ignore DB)
    db, err := driver.InitDatabase("sqlite3.db", *appLogger)
    require.NoError(t, err)

    // 6. Tracecore (you can keep or pass nil depending on VaultHandler contract)
    tracecoreClient := tracecore.NewTracecoreClient(
        os.Getenv("TRACECORE_URL"),
        os.Getenv("TRACECORE_TOKEN"),
        os.Getenv("CLOUD_FRONT_URL"),
        os.Getenv("CLOUD_BACK_URL"),
    )

    // 7. Registry (already in your code)
    reg := registry.NewRegistry(appLogger)
    reg.RegisterDefinitions([]registry.EntryDefinition{
        {
            Type:    "login",
            Factory: func() vaults_domain.VaultEntry { return &vaults_domain.LoginEntry{} },
            Handler: vault_ui.NewLoginHandler(*db, appLogger),
        },
        {
            Type:    "card",
            Factory: func() vaults_domain.VaultEntry { return &vaults_domain.CardEntry{} },
            Handler: vault_ui.NewCardHandler(*db, appLogger),
        },
        {
            Type:    "note",
            Factory: func() vaults_domain.VaultEntry { return &vaults_domain.NoteEntry{} },
            Handler: vault_ui.NewNoteHandler(*db, appLogger),
        },
        {
            Type:    "identity",
            Factory: func() vaults_domain.VaultEntry { return &vaults_domain.IdentityEntry{} },
            Handler: vault_ui.NewIdentityHandler(*db, appLogger),
        },
        {
            Type:    "sshkey",
            Factory: func() vaults_domain.VaultEntry { return &vaults_domain.SSHKeyEntry{} },
            Handler: vault_ui.NewSSHKeyHandler(*db, appLogger),
        },
    })

    // 8. IPFS client (you can keep real IPFS if you want, or a mock)
    ipfs := blockchain.NewIPFSClient(os.Getenv("IPFS_CLIENT"))

    // 9. VaultHandler with mocks
    vaultHandler := vault_ui.NewVaultHandler(
        reg,
        *appLogger,
        context.Background(),
        ipfs,
        nil,
        nil,
        *tracecoreClient,
        "",
    )

    // 10. Compute symKey
    aesService := vault_infrastructure_crypto.AESService{}
    symKey, err := aesService.DecryptPasswordWithStellarByte(
        []byte("cPu8Oc3uPwOZXq4Q"),
        []byte("Eo4A8jXUXZ/M9sfGHx95wrI0v3IQaAL9"),
        "SB3RRWLCYSMUZZIVYFM6PPGMXIYLRHWDTGMMGKVYWRW64LK4PJFWRLTF",
    )
    require.NoError(t, err)
    require.Equal(t, 32, len(symKey))

    // 11. Download request
    downloadReq := vault_ui.DownloadAttachmentRequest{
        UserID:       userID,
        Vault:        vault,
        CID:          cid,
        Ext:          ext,
        Password:     password,
        PrivateKey:   "private-key",
        EncryptedKey: "encrypted-key",
        SymKey:       symKey,
        AppCfg:       mockAppConfig(),
    }

    // 12. Act
    path, err := vaultHandler.DownloadAttachment(context.Background(), downloadReq)

    // 13. Assert
    require.NoError(t, err)

    require.Contains(t, path, "VaultCore/attachments")
    require.Contains(t, path, "attachment_"+cid)
    require.Contains(t, path, ext)

    require.Contains(t, path, "VaultCore/attachments")
    require.Contains(t, path, "attachment_"+cid)
    require.Contains(t, path, ext)
}

