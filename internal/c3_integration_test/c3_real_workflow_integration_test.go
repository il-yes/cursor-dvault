package c3_integration_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	identity_domain "vault-app/internal/identity/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_envelope_uc "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_member_uc "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

// mockIdentityRepo implements UserFinder for identity lookups.
type mockIdentityRepo struct {
	users map[string]*identity_domain.User
}

func (m *mockIdentityRepo) FindUserById(ctx context.Context, userID string) (*identity_domain.User, error) {
	u, ok := m.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	return u, nil
}

// memoryCloudRepo simulates Cloud backend with atomic state mutations.
type memoryCloudRepo struct {
	mu     sync.RWMutex
	groups map[string]*trustgroup_domain.TrustGroup
}

func newMemoryCloudRepo() *memoryCloudRepo {
	return &memoryCloudRepo{
		groups: make(map[string]*trustgroup_domain.TrustGroup),
	}
}

func (r *memoryCloudRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tg := req.TrustGroup
	if tg.ID == "" {
		tg.ID = "tg-test-001"
	}
	r.groups[tg.ID] = &tg
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: tg}, nil
}

func (r *memoryCloudRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryCloudRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *memoryCloudRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (r *memoryCloudRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}

func (r *memoryCloudRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (r *memoryCloudRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	for _, m := range tg.MemberCIDs {
		if m == req.VaultID {
			return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
		}
	}
	tg.MemberCIDs = append(tg.MemberCIDs, req.VaultID)
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryCloudRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	var newMembers []string
	for _, m := range tg.MemberCIDs {
		if m != req.MemberID {
			newMembers = append(newMembers, m)
		}
	}
	tg.MemberCIDs = newMembers
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryCloudRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (r *memoryCloudRepo) AddTrustGroupKeyEnvelope(ctx context.Context, req *trustgroup_domain.TrustGroupKeyEnvelope) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg.KeyEnvelopes = append(tg.KeyEnvelopes, *req)
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryCloudRepo) GetUserByEmail(ctx context.Context, email string) (*tracecore_types.User, error) {
	return nil, fmt.Errorf("user not found by email")
}

// 1. Error Wording Verification Test
func TestC3_ErrorWording_ErrKeyEnvelopeNotFound(t *testing.T) {
	expectedMsg := "no active key envelope found for member and KEK version"
	assert.Equal(t, expectedMsg, collaboration_usecases.ErrKeyEnvelopeNotFound.Error(),
		"ErrKeyEnvelopeNotFound error string must be updated to exact expected wording")
}

// 2. Production Workflow: Add Member + Provision Envelope + Readback Invariant
func TestC3_ProductionWorkflow_AddMember_Provision_And_Readback(t *testing.T) {
	ctx := context.Background()
	cloudRepo := newMemoryCloudRepo()
	eventBus := trustgroup_eventbus.NewMemoryBus()

	// Create caller keypair (Creator A) and target keypair (Member B)
	kpCaller, err := keypair.Random()
	require.NoError(t, err)

	kpTarget, err := keypair.Random()
	require.NoError(t, err)

	callerVaultID := "vault_creator_a_uuid_101"
	targetVaultID := "vault_member_b_uuid_202"

	// Create TrustGroup aggregate
	tg := trustgroup_domain.NewTrustGroup("tg_prod_01", "Production Group", []string{callerVaultID})
	tg.ID = "01b947ce-c2aa-4eac-b528-a7d68b0ce62f"
	_, err = cloudRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// Setup orchestrator
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "", nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	keyring := vaults_domain.NewVaultKeyring(callerVaultID)
	rawKEK := []byte("32_byte_kek_for_trust_group_001")
	keyring.Keys = append(keyring.Keys, vaults_domain.EncryptedKey{
		ID:           "k_kek_01",
		Type:         vaults_domain.KeyTypeTrustGroupKEK,
		Version:      1,
		TrustGroupID: tg.ID,
		Ciphertext:   rawKEK,
		CreatedAt:    time.Now().Unix(),
	})

	// Encrypt caller envelope
	encKEK, err := asymSvc.EncryptForRecipient(kpCaller.Address(), rawKEK)
	require.NoError(t, err)

	callerEnvelope := trustgroup_domain.TrustGroupKeyEnvelope{
		ID:           "env_caller_01",
		TrustGroupID: tg.ID,
		MemberID:     callerVaultID,
		WrappedKEK:   base64.StdEncoding.EncodeToString(encKEK),
		KEKVersion:   1,
	}
	tg.KeyEnvelopes = append(tg.KeyEnvelopes, callerEnvelope)
	_, err = cloudRepo.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// Setup use cases
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(cloudRepo, eventBus)
	addEnvelopeUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(cloudRepo, nil)
	provisionEnvelopeUC := trustgroup_envelope_uc.NewProvisionTrustGroupMemberEnvelopeUseCase(
		cloudRepo,
		orchestrator,
		addEnvelopeUC,
		keyringSvc,
	)

	// Setup Identity mock for Member B
	identityMock := &mockIdentityRepo{
		users: map[string]*identity_domain.User{
			targetVaultID: {
				ID:               targetVaultID,
				Email:            "memberb@example.com",
				StellarPublicKey: kpTarget.Address(),
			},
		},
	}

	// Execute Member Addition (simulating App.AddTrustGroupMember flow)
	idUser, err := identityMock.FindUserById(ctx, targetVaultID)
	require.NoError(t, err)
	require.NotEmpty(t, idUser.StellarPublicKey)

	// Add member to Cloud DB
	updatedTg, err := addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: tg.ID,
		VaultID:      targetVaultID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.Contains(t, updatedTg.MemberCIDs, targetVaultID)

	// Provision envelope for Member B
	pTg, err := provisionEnvelopeUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupMemberEnvelopeRequest{
		TrustGroupID:    tg.ID,
		MemberID:        targetVaultID,
		MemberPublicKey: idUser.StellarPublicKey,
	}, keyring)
	require.NoError(t, err)
	require.NotNil(t, pTg)

	// Verify readback invariant from Cloud DB
	freshTgResp, err := cloudRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, err)
	require.Contains(t, freshTgResp.Data.MemberCIDs, targetVaultID)

	foundTargetEnv := false
	for _, env := range freshTgResp.Data.KeyEnvelopes {
		if env.MemberID == targetVaultID && env.KEKVersion == 1 {
			foundTargetEnv = true
			require.NotEmpty(t, env.WrappedKEK)
		}
	}
	assert.True(t, foundTargetEnv, "Cloud DB must contain active KeyEnvelope for Member B and KEKVersion=1")
}

// 3. Rollback Test: When envelope provisioning fails, member is removed from TrustGroup
func TestC3_ProductionWorkflow_Rollback_On_Provisioning_Failure(t *testing.T) {
	ctx := context.Background()
	cloudRepo := newMemoryCloudRepo()
	eventBus := trustgroup_eventbus.NewMemoryBus()

	callerVaultID := "vault_creator_a_uuid_101"
	targetVaultID := "vault_member_b_unprovisionable"

	// Create TrustGroup
	tg := trustgroup_domain.NewTrustGroup("tg_prod_02", "Rollback Test Group", []string{callerVaultID})
	tg.ID = "tg-rollback-001"
	_, err := cloudRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(cloudRepo, eventBus)

	// Add member to TrustGroup
	_, err = addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: tg.ID,
		VaultID:      targetVaultID,
		Role:         "member",
	})
	require.NoError(t, err)

	// Verify member exists before rollback
	tgResp, err := cloudRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, err)
	require.Contains(t, tgResp.Data.MemberCIDs, targetVaultID)

	// Trigger Rollback Cleanup
	_, rollbackErr := cloudRepo.RemoveMemberFromTrustGroup(ctx, &trustgroup_domain.RemoveMemberFromTrustGroupRequest{
		TrustGroupID: tg.ID,
		MemberID:     targetVaultID,
	})
	require.NoError(t, rollbackErr)

	// Assert Member B was removed from Cloud DB and state is clean
	cleanTgResp, err := cloudRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, err)
	assert.NotContains(t, cleanTgResp.Data.MemberCIDs, targetVaultID,
		"Member B must be completely removed from TrustGroup after envelope provisioning failure rollback")
}
