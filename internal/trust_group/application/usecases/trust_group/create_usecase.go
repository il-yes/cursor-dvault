package trustgroup_usecases


import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_events "vault-app/internal/trust_group/application/events"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
)

type CreateTrustGroupUsecase struct {
	repo     trustgroup_domain.TrustGroupRepository
	eventBus trustgroup_events.TrustGroupEventBus
}

func NewCreateTrustGroupUsecase(
	repo trustgroup_domain.TrustGroupRepository,
	eventBus trustgroup_events.TrustGroupEventBus,
) *CreateTrustGroupUsecase {
	return &CreateTrustGroupUsecase{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (u *CreateTrustGroupUsecase) ValidateDependencies() error {
	if u.repo == nil {
		return errors.New("trust group repository is required")
	}

	if u.eventBus == nil {
		return errors.New("trust group event bus is required")
	}

	return nil
}

func (u *CreateTrustGroupUsecase) ValidateRequest(
	req trustgroup_dtos.CreateTrustGroupRequest,
) error {
	if strings.TrimSpace(req.ChannelID) == "" {
		return trustgroup_domain.ErrChannelIDRequired
	}

	if strings.TrimSpace(req.Name) == "" {
		return trustgroup_domain.ErrTrustGroupNameRequired
	}

	return nil
}

func (u *CreateTrustGroupUsecase) Execute(
	ctx context.Context,
	req trustgroup_dtos.CreateTrustGroupRequest,
) (*trustgroup_domain.TrustGroup, error) {

	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	trustGroup := trustgroup_domain.NewTrustGroup(
		req.ChannelID,
		req.Name,
		req.MemberCIDs,
	)

	created, err := u.repo.CreateTrustGroup(
		ctx,
		&trustgroup_domain.CreateTrustGroupRequest{
			TrustGroup: *trustGroup,
		},
	)
	if err != nil {
		return nil, err
	}

	event := trustgroup_domain.TrustGroupCreated{
		EventID:        uuid.NewString(),
		EventTimestamp: time.Now().UTC(),

		ChannelID: req.ChannelID,
		VaultID: req.VaultID,
		OwnerID: req.OwnerID,
		Name: req.Name,
	}

	if err := u.eventBus.PublishTrustGroupCreated(ctx, event); err != nil {
		return nil, err
	}

	return &created.Data, nil
}
