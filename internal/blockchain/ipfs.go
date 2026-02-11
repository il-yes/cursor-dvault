package blockchain

import (
	"bytes"
	"fmt"
	"github.com/ipfs/go-ipfs-api"
	"io"
)

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
func (client *IPFSClient) AddData(data []byte) (string, error) {
	fmt.Println("Adding data to IPFS...")
	reader := bytes.NewReader(data)
	cid, err := client.shell.Add(reader)
	if err != nil {
		return "", fmt.Errorf("failed to add data to IPFS: %w", err)
	}
	return cid, nil
}

// GetData retrieves data from IPFS using CID.
func (client *IPFSClient) GetData(cid string) ([]byte, error) {
	readCloser, err := client.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from IPFS: %w", err)
	}
	defer readCloser.Close()

	data, err := io.ReadAll(readCloser)
	if err != nil {
		return nil, fmt.Errorf("failed to read IPFS content: %w", err)
	}
	return data, nil
}

