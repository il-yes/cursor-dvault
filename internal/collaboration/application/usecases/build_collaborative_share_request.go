package collaboration_usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	"vault-app/internal/tracecore"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	vaults_domain "vault-app/internal/vault/domain"
)

type BuildCollaborativeShareRequestInput struct {
	AssetID      string
	TrustGroupID string
	KEKVersion   uint64
	RawPayload   []byte
	CreatedBy    string
	SenderEmail  string
	SenderID     string
	Keyring      *vaults_domain.VaultKeyring
	Title        string
	EntryType    string
	AccessMode   string
	ExpiresAt    *time.Time
}

type BuildCollaborativeShareRequestOutput struct {
	ProdRequest *tracecore.ProdCreateCryptoShareRequest
	ShareEntry  *c3_asset_domain.ShareEntry
	Prepared    *trustgroup_orchestrator.PreparedCollaborativeAsset
}

type CollaborativeShareRequestBuilder struct {
	orchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	tcClient     *tracecore.TracecoreClient
}

func NewCollaborativeShareRequestBuilder(
	orchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
	tcClient *tracecore.TracecoreClient,
) *CollaborativeShareRequestBuilder {
	return &CollaborativeShareRequestBuilder{
		orchestrator: orchestrator,
		tcClient:     tcClient,
	}
}

func (b *CollaborativeShareRequestBuilder) BuildAndPersist(
	ctx context.Context,
	input BuildCollaborativeShareRequestInput,
) (*BuildCollaborativeShareRequestOutput, error) {
	// 1. Single DEK cryptographic preparation (Encrypt payload with DEK, wrap DEK with TrustGroup KEK)
	prepared, err := b.orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      input.AssetID,
		TrustGroupID: input.TrustGroupID,
		KEKVersion:   input.KEKVersion,
		RawPayload:   input.RawPayload,
		Keyring:      input.Keyring,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to prepare collaborative asset: %w", err)
	}

	wrappedDEKStr := string(prepared.WrappedDEK)

	// 2. Build C3 ShareEntry aggregate
	shareEntry, err := c3_asset_domain.NewShareEntry(
		prepared.AssetID,
		input.TrustGroupID,
		wrappedDEKStr,
		input.KEKVersion,
		input.CreatedBy,
		map[string]string{"title": input.Title, "entry_type": input.EntryType},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create C3 share entry: %w", err)
	}

	// 3. Build Cloud ProdCreateCryptoShareRequest carrying TrustGroupID, WrappedDEK, KEKVersion
	recipients := map[string]tracecore.CryptoRecipient{
		input.TrustGroupID: {
			ID:            input.TrustGroupID,
			EncryptedKeys: wrappedDEKStr,
			Role:          "editor",
			TrustGroupID:  input.TrustGroupID,
			RecipientType: "trust_group",
		},
	}

	prodReq := &tracecore.ProdCreateCryptoShareRequest{
		SenderID:      input.SenderID,
		SenderEmail:   input.SenderEmail,
		Recipients:    recipients,
		VaultPayload:  base64.StdEncoding.EncodeToString(prepared.EncryptedData),
		EncryptedKeys: map[string]string{input.TrustGroupID: wrappedDEKStr},
		TrustGroupID:  input.TrustGroupID,
		WrappedDEK:    wrappedDEKStr,
		KEKVersion:    input.KEKVersion,
		Title:         input.Title,
		EntryType:     input.EntryType,
		AccessMode:    input.AccessMode,
		ExpiresAt:     input.ExpiresAt,
	}

	// 4. Send to Cloud tracecore client if available
	if b.tcClient != nil {
		_, err := b.tcClient.CreateShare(ctx, *prodReq)
		if err != nil {
			return nil, fmt.Errorf("failed to persist C3 share request to cloud: %w", err)
		}
	}

	return &BuildCollaborativeShareRequestOutput{
		ProdRequest: prodReq,
		ShareEntry:  &shareEntry,
		Prepared:    prepared,
	}, nil
}
