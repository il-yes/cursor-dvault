package trustgroup_member_usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_events "vault-app/internal/trust_group/application/events"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type AddMemberToTrustGroupUsecase struct {
	repo     trustgroup_domain.TrustGroupRepository
	eventBus trustgroup_events.TrustGroupEventBus
}

func NewAddMemberToTrustGroupUsecase(
	repo trustgroup_domain.TrustGroupRepository,
	eventBus trustgroup_events.TrustGroupEventBus,
) *AddMemberToTrustGroupUsecase {
	return &AddMemberToTrustGroupUsecase{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (u *AddMemberToTrustGroupUsecase) ValidateDependencies() error {
	if u.repo == nil {
		return errors.New("trust group repository is required")
	}

	if u.eventBus == nil {
		return errors.New("trust group event bus is required")
	}

	return nil
}

func (u *AddMemberToTrustGroupUsecase) ValidateRequest(req trustgroup_dtos.AddMemberToTrustGroupRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return errors.New("trust group id is required")
	}

	if strings.TrimSpace(req.VaultID) == "" {
		return errors.New("vault id is required")
	}

	return nil
}

func (u *AddMemberToTrustGroupUsecase) Execute(
	ctx context.Context,
	req trustgroup_dtos.AddMemberToTrustGroupRequest,
) (*trustgroup_domain.TrustGroup, error) {
	fmt.Printf(
		"[C3][ADD_MEMBER][ENTRY] trustGroupID=%s target=%s role=%s\n",
		req.TrustGroupID,
		req.VaultID,
		req.Role,
	)
	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][04] layer=DESKTOP_USECASE file=internal/trust_group/application/usecases/member/create_usecase.go function=AddMemberToTrustGroupUsecase.Execute input.TrustGroupID=%s input.VaultID=%s input.Role=%s\n",
		req.TrustGroupID, req.VaultID, req.Role)

	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = "member"
	}

	updated, err := u.repo.AddMemberToTrustGroup(
		ctx,
		&trustgroup_domain.AddMemberToTrustGroupRequest{
			TrustGroupID: req.TrustGroupID,
			VaultID:      req.VaultID,
			Role:         role,
		},
	)
	if err != nil {
		fmt.Printf("[C3][MEMBERSHIP][WRITE] trustGroupID=%s vaultID=%s UPDATE_FAILED=%v\n",
			req.TrustGroupID, req.VaultID, err)
		return nil, err
	}

	memberCount := len(updated.Data.MemberCIDs)
	memberCIDsStr := strings.Join(updated.Data.MemberCIDs, ", ")
	fmt.Printf("[C3][MEMBERSHIP][WRITE] trustGroupID=%s vaultID=%s memberCountAfter=%d MemberCIDs=[%s]\n",
		req.TrustGroupID, req.VaultID, memberCount, memberCIDsStr)

	event := trustgroup_domain.MemberAddedToTrustGroup{
		EventID:        uuid.NewString(),
		EventTimestamp: time.Now().UTC(),
		TrustGroupID:   req.TrustGroupID,
		MemberID:       req.VaultID,
	}

	if u.eventBus != nil {
		if err := u.eventBus.PublishMemberAddedToTrustGroup(ctx, event); err != nil {
			return nil, err
		}
	}

	return &updated.Data, nil
}