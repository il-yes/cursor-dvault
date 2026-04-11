package blockchain_test

import (
	"bytes"
	"context"
	"testing"

	"vault-app/internal/blockchain"
	tracecore_types "vault-app/internal/tracecore/types"
	types "vault-app/internal/tracecore/types"
)




func TestHybridStorage_AddAndGet(t *testing.T) {
    // 1. Create local, cloud, and hybrid
    localS := blockchain.NewDirectIPFSStorage("localhost:5001")
    mockClient := &mockTracecoreClient{}
    cloudS := blockchain.NewCloudIPFSStorage(mockClient, "test_user", "test_vault")
    h := blockchain.NewHybridStorage(localS, cloudS)

    plaintext := []byte("test vault content via hybrid IPFS")

    // 2. Mock cloud Add to return a CID
    mockClient.AddToIPFSFunc = func(ctx context.Context, req types.SyncVaultStreamRequest) (*types.CloudResponse[types.SyncVaultResponse], error) {
        return &types.CloudResponse[types.SyncVaultResponse]{
            Data: tracecore_types.SyncVaultResponse{
                UserID: "test_user",
                CID:    "QmHybridCloudCID",
            },
        }, nil
    }

    // 3. Add (should go to local AND cloud)
    cid, err := h.Add(context.Background(), plaintext)
    if err != nil {
        t.Fatalf("HybridStorage.Add failed: %v", err)
    }

    // 4. Check that local received it too
    localData, err := localS.Get(context.Background(), cid)
    if err != nil {
        t.Fatalf("localS.Get failed: %v", err)
    }

    // 5. Compare
    if !bytes.Equal(localData, plaintext) {
        t.Errorf("HybridStorage (local) data mismatch: %q", string(localData))
    }

    // 6. Mock cloud Get to return data
    mockClient.GetDataFromCloudFunc = func(ctx context.Context, req types.IpfsCidRequest) (*types.IpfsCidResponse, error) {
        return &types.IpfsCidResponse{
            Data: string(plaintext),
        }, nil
    }

    // 7. Get (should hit local and cloud, returns local result)
    data, err := h.Get(context.Background(), cid)
    if err != nil {
        t.Fatalf("HybridStorage.Get failed: %v", err)
    }

    if !bytes.Equal(data, plaintext) {
        t.Errorf("HybridStorage data mismatch")
    }
}