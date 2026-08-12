package trustgroup_member_usecases

import (
	"context"
	"errors"
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

	if strings.TrimSpace(req.ChannelID) == "" {
		return errors.New("channel id is required")
	}

	if strings.TrimSpace(req.MemberID) == "" {
		return errors.New("member id is required")
	}

	return nil
}

func (u *AddMemberToTrustGroupUsecase) Execute(
	ctx context.Context,
	req trustgroup_dtos.AddMemberToTrustGroupRequest,
) (*trustgroup_domain.TrustGroup, error) {

	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	updated, err := u.repo.AddMemberToTrustGroup(
		ctx,
		&trustgroup_domain.AddMemberToTrustGroupRequest{
			TrustGroupID: req.TrustGroupID,
			MemberID:     req.MemberID,
		},
	)
	if err != nil {
		return nil, err
	}

	event := trustgroup_domain.MemberAddedToTrustGroup{
		EventID:        uuid.NewString(),
		EventTimestamp: time.Now().UTC(),

		TrustGroupID: req.TrustGroupID,
		ChannelID:      req.ChannelID,
		MemberID:       req.MemberID,
	}

	if err := u.eventBus.PublishMemberAddedToTrustGroup(ctx, event); err != nil {
		return nil, err
	}

	return &updated.Data, nil
}
