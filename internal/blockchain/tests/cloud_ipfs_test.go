package blockchain_test

import (
	"bytes"
	"context"
	"testing"

	"vault-app/internal/blockchain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// -----------------------------------------------------------------------------------------------
// MOCK TRACECORE CLIENT
// -----------------------------------------------------------------------------------------------
type mockTracecoreClient struct {
    AddToIPFSFunc        func(context.Context, tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error)
    GetDataFromCloudFunc func(context.Context, tracecore_types.IpfsCidRequest) (*tracecore_types.IpfsCidResponse, error)
    AddToS3Func          func(context.Context, tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error)
    GetDataFromS3Func    func(context.Context, tracecore_types.IpfsCidRequest) (*tracecore_types.CloudResponse[tracecore_types.IpfsCidResponse], error)
}

func (m *mockTracecoreClient) AddToIPFS(
    ctx context.Context,
    req tracecore_types.SyncVaultStreamRequest,
) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error) {
    if m.AddToIPFSFunc != nil {
        return m.AddToIPFSFunc(ctx, req)
    }
    return &tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse]{Data: tracecore_types.SyncVaultResponse{
        UserID: "test",
        CID:    "QmcCID",
    }}, nil
}

func (m *mockTracecoreClient) GetDataFromCloudStorage(
    ctx context.Context,
    req tracecore_types.IpfsCidRequest,
) (*tracecore_types.IpfsCidResponse, error) {
    if m.GetDataFromCloudFunc != nil {
        return m.GetDataFromCloudFunc(ctx, req)
    }
    return &tracecore_types.IpfsCidResponse{
        Data: "encodedData",
    }, nil
}
func (m *mockTracecoreClient) AddToS3(
    ctx context.Context,
    req tracecore_types.SyncVaultStreamRequest,
) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error) {
    if m.AddToS3Func != nil {
        return m.AddToS3Func(ctx, req)
    }
    return &tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse]{Data: tracecore_types.SyncVaultResponse{
        UserID: "test",
        CID:    "QmcCID",
    }}, nil
}
func (m *mockTracecoreClient) GetDataFromS3(
    ctx context.Context,
    req tracecore_types.IpfsCidRequest,
) (*tracecore_types.CloudResponse[tracecore_types.IpfsCidResponse], error) {
    if m.GetDataFromS3Func != nil {
        return m.GetDataFromS3Func(ctx, req)
    }
    return &tracecore_types.CloudResponse[tracecore_types.IpfsCidResponse]{Data: tracecore_types.IpfsCidResponse{
        Data: "encodedData",
    }}, nil
}

// -----------------------------------------------------------------------------------------------
// TESTS
// -----------------------------------------------------------------------------------------------
func TestCloudIPFSStorage_AddAndGet(t *testing.T) {
    // 1. Mock client
    client := &mockTracecoreClient{}

    // 2. CloudIPFSStorage
    s := blockchain.NewCloudIPFSStorage(client, "test_user", "test_vault")

    plaintext := []byte("test vault content via cloud IPFS")

    // 3. Mock AddToIPFS to return a CID we can check
    client.AddToIPFSFunc = func(ctx context.Context, req tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error) {
        if string(req.Stream) != string(plaintext) {
            t.Errorf("AddToIPFS received different data:\nexpected:%q\ngot:%q",
                string(plaintext), string(req.Stream))
        }
        return &tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse]{
            Data: tracecore_types.SyncVaultResponse{
                UserID: "test_user",
                CID:    "QmSGH4oAre11ktm4DvMsQ7XT2xfyxd9ERU68cUFRrt6FB7",
            },
        }, nil
    }

    // 4. Add
    cid, err := s.Add(context.Background(), plaintext)
    if err != nil {
        t.Fatalf("CloudIPFSStorage.Add failed: %v", err)
    }
    if cid != "QmSGH4oAre11ktm4DvMsQ7XT2xfyxd9ERU68cUFRrt6FB7" {
        t.Errorf("expected CID=%s, got %s", "QmSGH4oAre11ktm4DvMsQ7XT2xfyxd9ERU68cUFRrt6FB7", cid)
    }

    // 5. Mock GetDataFromCloudStorage
    client.GetDataFromCloudFunc = func(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.IpfsCidResponse, error) {
        if req.UserID != "test_user" || req.VaultName != "test_vault" || req.CID != cid {
            t.Errorf("GetDataFromCloudStorage got wrong request: %+v", req)
        }
        return &tracecore_types.IpfsCidResponse{
            Data: string(plaintext),
        }, nil
    }

    // 6. Get
    data, err := s.Get(context.Background(), cid)
    if err != nil {
        t.Fatalf("CloudIPFSStorage.Get failed: %v", err)
    }

    // 7. Compare
    if !bytes.Equal(data, plaintext) {
        t.Errorf("CloudIPFSStorage data mismatch:\nexpected:%q\ngot:%q",
            string(plaintext), string(data))
    }
}