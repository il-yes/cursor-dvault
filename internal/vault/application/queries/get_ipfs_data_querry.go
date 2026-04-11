package vault_queries

import (
	"context"
	"encoding/base64"
	"fmt"
	"unicode/utf8"
	app_config_domain "vault-app/internal/config/domain"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
)


// -------- INTERFACES --------
type CryptoServiceInterface interface {
	Decrypt(data []byte, password string) ([]byte, error)
}
type IpfsServiceInterface interface {
	Get(ctx context.Context, cid string) ([]byte, error)
}


// -------- HANDLER --------
type GetIPFSDataQuerryHandler struct {
	cryptoService CryptoServiceInterface
	IpfsService IpfsServiceInterface
}

// -------- CONSTRUCTOR --------
func NewGetIPFSDataQuerryHandler(cryptoService CryptoServiceInterface) *GetIPFSDataQuerryHandler {
	return &GetIPFSDataQuerryHandler{
		cryptoService: cryptoService,
	}
}

// -------- QUERRY --------
type GetIPFSDataQuerry struct {
	CID       string
	Password  string
	AppCfg    app_config_domain.AppConfig
	UserID    string
	VaultName string
}
// -------- RESPONSE --------
type GetIPFSDataResponse struct {
	Data vault_domain.VaultPayload
}
func (h *GetIPFSDataQuerryHandler) Execute(cmd GetIPFSDataQuerry) (*GetIPFSDataResponse, error) {
	// 1. Get raw encrypted bytes from IPFS → this must be the base64 string as bytes
	rawBytes, err := h.GetFromIpfs(context.Background(), cmd)
	if err != nil {
		return nil, err
	}

	// 2. Must be valid Base64 string
	if !utf8.Valid(rawBytes) {
		return nil, fmt.Errorf("invalid UTF‑8 in base64 input")
	}

	// 3. Decode Base64 → binary (salt + nonce + ciphertext)
	decoded, err := base64.StdEncoding.DecodeString(string(rawBytes))
	if err != nil {
		return nil, fmt.Errorf("❌ base64 decode failed: %w", err)
	}

	// 4. Decrypt binary data
	plain, err := h.cryptoService.Decrypt(decoded, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("❌ decryption failed: %w", err)
	}

	vaultPayload := vaults_domain.ParseVaultPayload(plain)
	return &GetIPFSDataResponse{
		Data:       vaultPayload,
		// DecodedLen: len(decoded),
		// PlainLen:   len(plain),
	}, nil
}

func (h *GetIPFSDataQuerryHandler) GetFromIpfs(ctx context.Context, req GetIPFSDataQuerry) ([]byte, error) {
	base64Text, err := h.IpfsService.Get(ctx, req.CID)
	if err != nil {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - failed to get ipfs: %w", err)
	}
	// // ✅ DECRYPT HERE
	// decrypted, err := h.cryptoService.Decrypt(encrypted, req.Password)
	// if err != nil {
	// 	return nil, fmt.Errorf("decrypt failed: %w", err)
	// }

	return base64Text, nil
}


func (h *GetIPFSDataQuerryHandler) SetIpfsService(ipfs IpfsServiceInterface) {
	h.IpfsService = ipfs
}