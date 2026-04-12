package blockchain_test

import (
	"bytes"
	"context"
	"testing"

	"vault-app/internal/blockchain"
)

// -----------------------------------------------------------------------------------------------
// 	TESTS
// -----------------------------------------------------------------------------------------------
func TestDirectIPFSStorage_AddAndGet(t *testing.T) {
    // 1. DirectIPFSStorage against localhost:5001
    s := blockchain.NewDirectIPFSStorage("localhost:5001")

    // 2. Test data
    plaintext := []byte("test vault content on desktop IPFS")

    // 3. Add
    cid, err := s.Add(context.Background(), plaintext)
    if err != nil {
        t.Fatalf("Add failed: %v", err)
    }

    // 4. Get
    data, err := s.Get(context.Background(), cid)
    if err != nil {
        t.Fatalf("Get failed: %v", err)
    }

    // 5. Compare
    if !bytes.Equal(data, plaintext) {
        t.Errorf("DirectIPFSStorage data mismatch:\nexpected:%q\ngot:%q",
            string(plaintext), string(data))
    }
}