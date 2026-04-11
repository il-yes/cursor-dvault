package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	app_config "vault-app/internal/config"
	tracecore_types "vault-app/internal/tracecore/types"
	utils "vault-app/internal/utils"

	shell "github.com/ipfs/go-ipfs-api"
)

type TracecoreClt interface {
	AddToIPFS(ctx context.Context, req tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error)
	GetDataFromCloudStorage(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.IpfsCidResponse, error)
	AddToS3(ctx context.Context, req tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error)
	GetDataFromS3(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.CloudResponse[tracecore_types.IpfsCidResponse], error)
}

// ---------------------------------------------------------
// Vault Storage Factory
// ---------------------------------------------------------

type Config struct {
	StorageConfig app_config.StorageConfig
	UserID        string
	VaultName     string
}

func NewStorageProvider(cfg Config, client TracecoreClt) app_config.StorageProvider {
	switch cfg.StorageConfig.Mode {

	case app_config.StorageCloud:
		utils.LogPretty("StorageCloud - Cloud.APIEndpoint", cfg.StorageConfig.Cloud.BaseURL)
		return NewCloudIPFSStorage(client, cfg.UserID, cfg.VaultName)

	case app_config.StorageLocal:
		utils.LogPretty("StorageLocal - LocalIPFS.APIEndpoint", cfg.StorageConfig.LocalIPFS.APIEndpoint)
		return NewDirectIPFSStorage(cfg.StorageConfig.LocalIPFS.APIEndpoint)

	case app_config.StoragePrivateIPFS:
		return NewDirectIPFSStorage(cfg.StorageConfig.PrivateIPFS.APIEndpoint)

	case app_config.StorageEnterpriseS3:
		return NewEnterpriseS3Storage(client, cfg.UserID, cfg.VaultName)

	case app_config.StorageHybrid:
		return NewHybridStorage(
			NewDirectIPFSStorage(cfg.StorageConfig.LocalIPFS.APIEndpoint),
			NewCloudIPFSStorage(client, cfg.UserID, cfg.VaultName),
		)

	default:
		return NewCloudIPFSStorage(client, cfg.UserID, cfg.VaultName) // ← cloud as default (production mindset)
	}
}

// ---------------------------------------------------------
// IPFS Client - to delete - Duplicata with localIPFS 
// ---------------------------------------------------------
type IPFSClient struct {
	shell *shell.Shell
}

// NewIPFSClient creates a new IPFS client connected to a node (default: localhost).
func NewIPFSClient(endpoint string) *IPFSClient {
	fmt.Println("🔧 Initializing IPFS client...")
	if endpoint == "" {
		endpoint = "localhost:5001"
	}
	return &IPFSClient{
		shell: shell.NewShell(endpoint),
	}
}

// AddData adds encrypted data to IPFS and returns the CID. 
func (client *IPFSClient) Add(ctx context.Context, data []byte) (string, error) {
	fmt.Println("Adding data to IPFS...")
	fmt.Println(len(data))
	reader := bytes.NewReader(data)
	cid, err := client.shell.Add(reader)
	if err != nil {
		return "", fmt.Errorf("failed to add data to IPFS: %w", err)
	}
	err = client.shell.Pin(cid)
	if err != nil {
		return "", fmt.Errorf("failed to pin: %w", err)
	}
	return cid, nil
}
func (client *IPFSClient) GetData(cid string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.shell.Request("cat", cid).Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("IPFS request failed: %w", err)
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, fmt.Errorf("IPFS response error: %v", resp.Error)
	}

	data, err := io.ReadAll(resp.Output) // 🔥 THIS IS THE FIX
	if err != nil {
		return nil, fmt.Errorf("failed to read IPFS output: %w", err)
	}

	utils.LogPretty("IPFS DATA SIZE:", len(data)) // 🔥 DEBUG
	utils.LogPretty("resp.Output is nil?", resp.Output == nil)

	return data, nil
}

// Legacy
func (client *IPFSClient) GetData1(cid string) ([]byte, error) {
	// readCloser, err := client.shell.Cat(cid)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve data from IPFS: %w", err)
	// }
	fmt.Print("Before ReadAll")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.shell.Request("cat", cid).Send(ctx)
	// resp, err := http.Get("http://127.0.0.1:8080/ipfs/" + cid)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	data, err := io.ReadAll(resp.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to read IPFS content: %w", err)
	}
	return data, nil
}
func (client *IPFSClient) GetData0(cid string) ([]byte, error) {
	resp, err := http.Get("http://127.0.0.1:8080/ipfs/" + cid)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
func (client *IPFSClient) GetData2(cid string) ([]byte, error) {
	resp, err := client.shell.Request("cat", cid).Send(context.Background())
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, fmt.Errorf("ipfs error: %s", resp.Error.Message)
	}

	data, err := io.ReadAll(resp.Output)
	if err != nil {
		return nil, err
	}

	// 🔥 CRITICAL: drain remaining body if any
	io.Copy(io.Discard, resp.Output)

	return data, nil
}

// ---------------------------------------------------------
// Local IPFS Storage
// ---------------------------------------------------------
type DirectIPFSStorage struct {
	shell *shell.Shell
}

func NewDirectIPFSStorage(endpoint string) *DirectIPFSStorage {
	if endpoint == "" {
		endpoint = "localhost:5001"
	}
	return &DirectIPFSStorage{
		shell: shell.NewShell(endpoint),
	}
}

func (d *DirectIPFSStorage) Add(ctx context.Context, data []byte) (string, error) {
	shellID, err := d.shell.ID()
	if err != nil {
		utils.LogPretty("DirectIPFSStorage - Add - shell.ID()", err)
	}
	utils.LogPretty("DirectIPFSStorage - Add - shell.ID()", shellID)
	reader := bytes.NewReader(data)
	utils.LogPretty("DirectIPFSStorage - Add - reader", reader)
	cid, err := d.shell.Add(reader)
	if err != nil {
		utils.LogPretty("DirectIPFSStorage - Add - err", err)
	}
	utils.LogPretty("DirectIPFSStorage - Add - added cid", cid)
	err = d.shell.Pin(cid)
	if err != nil {
		utils.LogPretty("DirectIPFSStorage - Add - pin err", err)
	}
	return cid, nil
}

func (d *DirectIPFSStorage) Get(ctx context.Context, cid string) ([]byte, error) {
	rc, err := d.shell.Cat(cid)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// ---------------------------------------------------------
// Cloud IPFS Storage
// ---------------------------------------------------------
type CloudIPFSStorage struct {
	client TracecoreClt
	userID string
	vault  string
}

func NewCloudIPFSStorage(client TracecoreClt, userID, vault string) *CloudIPFSStorage {
	utils.LogPretty("CloudIPFSStorage - NewCloudIPFSStorage - userID", userID)

	return &CloudIPFSStorage{
		client: client,
		userID: userID,
		vault:  vault,
	}
}

func (c *CloudIPFSStorage) Add(ctx context.Context, data []byte) (string, error) {
	req := tracecore_types.SyncVaultStreamRequest{
		UserID:    c.userID,
		VaultName: c.vault,
		Stream:    data,
	}
	resp, err := c.client.AddToIPFS(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Data.CID, nil
}

func (c *CloudIPFSStorage) Get(ctx context.Context, cid string) ([]byte, error) {
	utils.LogPretty("CloudIPFSStorage - Get - userID", c.userID)
	utils.LogPretty("CloudIPFSStorage - Get - vault", c.vault)
	utils.LogPretty("CloudIPFSStorage - Get - cid", cid)
	req := tracecore_types.IpfsCidRequest{
		UserID:    c.userID,
		VaultName: c.vault,
		CID:       cid,
	}
	resp, err := c.client.GetDataFromCloudStorage(ctx, req)
	if err != nil {
		return nil, err
	}

	return []byte(resp.Data), nil
}

// ---------------------------------------------------------
// Enterprise S3 Storage
// ---------------------------------------------------------
type EnterpriseS3Storage struct {
	client TracecoreClt
	userID string
	vault  string
}

func NewEnterpriseS3Storage(client TracecoreClt, userID, vault string) *EnterpriseS3Storage {
	return &EnterpriseS3Storage{
		client: client,
		userID: userID,
		vault:  vault,
	}
}

func (e *EnterpriseS3Storage) Add(ctx context.Context, data []byte) (string, error) {
	response, err := e.client.AddToS3(ctx, tracecore_types.SyncVaultStreamRequest{
		UserID:    e.userID,
		VaultName: e.vault,
		Stream:    data,
	})
	if err != nil {
		return "", err
	}
	return response.Data.CID, nil
}

func (e *EnterpriseS3Storage) Get(ctx context.Context, cid string) ([]byte, error) {
	response, err := e.client.GetDataFromS3(ctx, tracecore_types.IpfsCidRequest{
		UserID:    e.userID,
		VaultName: e.vault,
		CID:       cid,
	})
	if err != nil {
		return nil, err
	}
	return []byte(response.Data.Data), nil
}

// ---------------------------------------------------------
// Hybrid Storage
// ---------------------------------------------------------
type HybridStorage struct {
	local *DirectIPFSStorage
	cloud *CloudIPFSStorage
}

func NewHybridStorage(local *DirectIPFSStorage, cloud *CloudIPFSStorage) *HybridStorage {
	return &HybridStorage{
		local: local,
		cloud: cloud,
	}
}

func (h *HybridStorage) Add(ctx context.Context, data []byte) (string, error) {
	res1, err1 := h.local.Add(ctx, data)
	if err1 != nil {
		return "", err1
	}
	res2, err2 := h.cloud.Add(ctx, data)
	if err2 != nil {
		return "", err2
	}
	utils.LogPretty("HybridStorage - Add - res1", res1)
	utils.LogPretty("HybridStorage - Add - res2", res2)
	return res1, nil
}

func (h *HybridStorage) Get(ctx context.Context, cid string) ([]byte, error) {
	res1, err1 := h.local.Get(ctx, cid)
	if err1 != nil {
		return nil, err1
	}
	res2, err2 := h.cloud.Get(ctx, cid)
	if err2 != nil {
		return nil, err2
	}
	utils.LogPretty("HybridStorage - Get - res1", res1)
	utils.LogPretty("HybridStorage - Get - res2", res2)
	return res1, nil
}
