package vaults_service

import (
	"context"
	"encoding/json"
	app_config "vault-app/internal/config"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_dto "vault-app/internal/vault/application/dto"
	vaults_domain "vault-app/internal/vault/domain"
)


type CIDResolver struct {
    ipfs app_config.StorageProvider
    crypto vaults_domain.VaultCrypto
    unlock vault_commands.UnlockVaultHandlerInterface
}
type TypedNode struct {
    Type string
    VaultNode *vaults_domain.VaultNode
    Entry     any
    Folder    *vaults_domain.Folder
    Index     *vaults_domain.Index
    Raw       []byte
}
func (r *CIDResolver) Resolve(ctx context.Context, cid string, password string) (*TypedNode, error) {

    raw, err := r.ipfs.Get(ctx, cid)
    if err != nil {
        return nil, err
    }

    unlockRes, err := r.unlock.Execute(vault_dto.UnlockVaultCommand{
        Password: password,
    })
    if err != nil {
        return nil, err
    }

    plain, err := r.crypto.Decrypt(raw, unlockRes.VaultKey.Key)
    if err != nil {
        return nil, err
    }

    // ---- detect type ONCE ----
    var vaultNode vaults_domain.VaultNode
    if json.Unmarshal(plain, &vaultNode) == nil && vaultNode.Type == "vault" {
        return &TypedNode{
            Type:      "vault",
            VaultNode: &vaultNode,
            Raw:       plain,
        }, nil
    }

    var entry struct {
        Type string `json:"type"`
    }
    if json.Unmarshal(plain, &entry) == nil && entry.Type != "" {
        return &TypedNode{
            Type: entry.Type,
            Raw:  plain,
        }, nil
    }

    return &TypedNode{
        Type: "unknown",
        Raw:  plain,
    }, nil
}