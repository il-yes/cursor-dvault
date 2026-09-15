package trustgroup_usecases_envelope

import (
	"context"
	"fmt"
	"strings"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	"vault-app/internal/utils"
)

type AddTrustGroupKeyEnvelopeUseCase struct {
	repo trustgroup_domain.TrustGroupRepository
}

func NewAddTrustGroupKeyEnvelopeUseCase(
	repo trustgroup_domain.TrustGroupRepository,
	_ ...interface{},
) *AddTrustGroupKeyEnvelopeUseCase {
	return &AddTrustGroupKeyEnvelopeUseCase{
		repo: repo,
	}
}

func (uc *AddTrustGroupKeyEnvelopeUseCase) ValidateDependencies() error {
	if uc.repo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}
	return nil
}

func (uc *AddTrustGroupKeyEnvelopeUseCase) ValidateRequest(req trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return trustgroup_domain.ErrTrustGroupIDRequired
	}
	if strings.TrimSpace(req.MemberID) == "" {
		return trustgroup_domain.ErrMemberIDRequired
	}
	if strings.TrimSpace(req.WrappedKEK) == "" {
		return trustgroup_domain.ErrWrappedKEKRequired
	}
	if req.KEKVersion == 0 {
		return trustgroup_domain.ErrKEKVersionRequired
	}
	return nil
}

func (uc *AddTrustGroupKeyEnvelopeUseCase) Execute(
	ctx context.Context,
	req trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest,
) (*trustgroup_domain.TrustGroup, error) {
	utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - Execute - req - ", req)

	if err := uc.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := uc.ValidateRequest(req); err != nil {
		return nil, err
	}

	// 1. Resolve TrustGroup
	resp, err := uc.repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: req.TrustGroupID,
	})
	if err != nil {
		utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - Execute - GetTrustGroup - error:", err)
		return nil, err
	}
	if resp == nil || resp.Data.ID == "" {
		utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - Execute - GetTrustGroup - error:", trustgroup_domain.ErrTrustGroupNotFound)
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}

	tg := &resp.Data

	// 2. Validate KEKVersion matches current TrustGroup version
	if req.KEKVersion != tg.KEKVersion {
		return nil, trustgroup_domain.ErrStaleKEKVersion
	}

	// 3. Verify Member belongs to TrustGroup
	memberFound := false
	for _, cid := range tg.MemberCIDs {
		if cid == req.MemberID {
			memberFound = true
			break
		}
	}
	if !memberFound {
		utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - Execute - MemberNotInTrustGroup - error:", trustgroup_domain.ErrMemberNotInTrustGroup)
		return nil, trustgroup_domain.ErrMemberNotInTrustGroup
	}

	// 4. Add Key Envelope to TrustGroup aggregate (Member-level envelope identity)
	fmt.Printf("[ENVELOPE][GENERATED]\ntrustGroupID=%s\nmemberID=%s\nkekVersion=%d\n",
		tg.ID, req.MemberID, req.KEKVersion)

	beforeCount := len(tg.KeyEnvelopes)
	envelope := trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: tg.ID,
		MemberID:     req.MemberID,
		KEKVersion:   req.KEKVersion,
		WrappedKEK:   req.WrappedKEK,
	}

	if err := tg.AddEnvelope(envelope); err != nil {
		utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - Execute - AddEnvelope - error:", err)
		return nil, err
	}

	afterCount := len(tg.KeyEnvelopes)
	lastEnvID := ""
	if afterCount > 0 {
		lastEnvID = tg.KeyEnvelopes[afterCount-1].ID
	}
	fmt.Printf("[TRUSTGROUP][ENVELOPE][AGGREGATE]\nbeforeCount=%d\nafterCount=%d\nenvelopeID=%s\n",
		beforeCount, afterCount, lastEnvID)

	fmt.Printf("[ENVELOPE][API_UPDATE]\ntrustGroupID=%s envelopes=%d\n", tg.ID, len(tg.KeyEnvelopes))

	// 5. Persist updated TrustGroup
	updatedResp, err := uc.repo.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{
		TrustGroup: *tg,
	})
	if err != nil {
		return nil, err
	}
	if updatedResp == nil {
		return nil, trustgroup_domain.ErrRepositoryResponse
	}

	fmt.Printf("[ENVELOPE][API_UPDATE][SUCCESS]\n")

	utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - Execute - ", envelope)
	utils.LogPretty("AddTrustGroupKeyEnvelopeUseCase - updatedResp - ", updatedResp)

	return &updatedResp.Data, nil
}
