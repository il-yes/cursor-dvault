package share_application_use_cases

import (
	"context"
	"fmt"
	"strconv"
	"time"
	utils "vault-app/internal"
	share_application_events "vault-app/internal/application/events/share"
	share_domain "vault-app/internal/domain/shared"
)

type TracecoreClientInterface interface {
	CreateShare(ctx context.Context, s share_domain.ShareEntry) (*share_domain.ShareEntry, error)
	AcceptShare(ctx context.Context, shareID string) error
	RejectShare(ctx context.Context, shareID string) error
	GetShareByMe(ctx context.Context) ([]share_domain.ShareEntry, error)
	GetShareWithMe(ctx context.Context) ([]share_domain.ShareEntry, error)
	SetToken(token string)
}
type ShareUseCase struct {
	repo       share_domain.Repository
	dispatcher share_application_events.EventDispatcher
	tc         TracecoreClientInterface // new cloud client

}

func NewShareUseCase(repo share_domain.Repository, tc TracecoreClientInterface, d share_application_events.EventDispatcher) *ShareUseCase {
	return &ShareUseCase{
		repo:       repo,
		tc:         tc,
		dispatcher: d,
	}
}

// ---------------------------------------------------------
// Create Share
// ---------------------------------------------------------
func (uc *ShareUseCase) CreateShare(ctx context.Context, s share_domain.ShareEntry) (*share_domain.ShareEntry, error) {
	utils.LogPretty("share - ShareUseCase", s)
	// Mirror to cloud if client available
	createdRes, err := uc.tc.CreateShare(ctx, s); 
	if err != nil {
		return nil, fmt.Errorf("cloud CreateShare failed: %w", err)
	}
	utils.LogPretty("share - ShareUseCase - createdRes", createdRes)	

	// Dispatch event if dispatcher present
	if uc.dispatcher != nil {
		uc.dispatcher.Dispatch(share_domain.ShareCreated{
			BaseEvent: share_domain.BaseEvent{
				Name: "ShareCreated",
				Time: time.Now(),
			},
			ShareID: s.ID,
			OwnerID: s.OwnerID,
		})
	}

	return createdRes, nil
}
// ------------------------------------------------
// Use case: list shared entries for a user
// ------------------------------------------------
func (s *ShareUseCase) ListSharedEntries(ctx context.Context, userID uint, cloudToken string) ([]share_domain.ShareEntry, error) {
	utils.LogPretty("share - ListSharedEntries", userID)
	s.tc.SetToken(cloudToken)
	// Mirror to cloud if client available
	cloudShares, err := s.tc.GetShareByMe(ctx)
	if err != nil {
		return nil, fmt.Errorf("dvault ListReceivedShares failed: %w", err)
	}
	utils.LogPretty("share - ListSharedEntries - cloudShares", cloudShares)	
	
	return cloudShares, nil	
}

// ------------------------------------------------
// Use case: fetch shares *received* by the user
// ------------------------------------------------
func (s *ShareUseCase) ListReceivedShares(ctx context.Context, userID uint, cloudToken string) ([]share_domain.ShareEntry, error) {
	utils.LogPretty("share - ListSharedEntries", userID)
	s.tc.SetToken(cloudToken)
	// Mirror to cloud if client available
	cloudShares, err := s.tc.GetShareWithMe(ctx)
	if err != nil {
		return nil, fmt.Errorf("dvault ListReceivedShares failed: %w", err)
	}
	utils.LogPretty("share - ListSharedEntries - cloudShares", cloudShares)	
	
	return cloudShares, nil	
}

func (uc *ShareUseCase) GetShareForAccept(
	ctx context.Context,
	shareID string,
	recipientUserID uint,
) (*share_domain.ShareAcceptData, error) {

	share, recipient, blob, err :=
		uc.repo.GetShareForAccept(shareID, recipientUserID)

	if err != nil {
		return nil, fmt.Errorf("cannot get share for accept: %w", err)
	}

	return &share_domain.ShareAcceptData{
		Share:     *share,
		Recipient: *recipient,
		Blob:      blob,
	}, nil
}

type AcceptShareResult struct {
	Share     share_domain.ShareEntry
	Recipient share_domain.Recipient
	Blob      []byte // encrypted payload for this user
}

// ---------------------------------------------------------
// Accept Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) AcceptShare(ctx context.Context, shareID uint, userID uint) (*AcceptShareResult, error) {

	// 1. Load share entry + recipient-specific data
	share, recipient, err := uc.repo.GetShareAndRecipient(ctx, strconv.FormatUint(uint64(shareID), 10), userID)
	if err != nil {
		return nil, err
	}

	// 2. Check expiration
	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		return nil, share_domain.ErrShareExpired
	}

	// 3. Mark recipient accepted
	if err := uc.repo.MarkRecipientAccepted(ctx, recipient.ID); err != nil {
		return nil, fmt.Errorf("failed to accept share: %w", err)
	}

	// 4. Return data to caller (VaultHandler → frontend)
	return &AcceptShareResult{
		Share:     *share,
		Recipient: *recipient,
		Blob:      recipient.EncryptedBlob,
	}, nil
}

type RejectShareResult struct {
	ShareID     uint
	RecipientID uint
	Message     string
}

// ---------------------------------------------------------
// Reject Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) RejectShare(ctx context.Context, shareID uint, userID uint) (*RejectShareResult, error) {

	// Load share + recipient
	_, recipient, err := uc.repo.GetShareAndRecipient(ctx, strconv.FormatUint(uint64(shareID), 10), userID)
	if err != nil {
		return nil, err
	}

	// Mark the recipient invitation as "rejected"
	if err := uc.repo.MarkRecipientRejected(ctx, recipient.ID); err != nil {
		return nil, fmt.Errorf("failed to reject share: %w", err)
	}

	return &RejectShareResult{
		ShareID:     shareID,
		RecipientID: stringToUint(recipient.ID),
		Message:     "Share invitation rejected",
	}, nil
}

// ---------------------------------------------------------
// Add Receiver
// ---------------------------------------------------------
type AddReceiverInput struct {
	ShareID uint
	Name    string
	Email   string
	Role    string
}

type AddReceiverResult struct {
	ShareID     uint
	RecipientID uint
	Message     string
}

func (uc *ShareUseCase) AddReceiver(ctx context.Context, requesterID uint, in AddReceiverInput) (*AddReceiverResult, error) {

	// Load share
	share, err := uc.repo.GetShareByID(ctx, strconv.FormatUint(uint64(in.ShareID), 10))

	if err != nil {
		return nil, fmt.Errorf("share not found: %w", err)
	}

	// Domain rule: only owner can add recipients
	if !share_domain.CanAddRecipient(share, requesterID) {
		return nil, fmt.Errorf("permission denied: not share owner")
	}

	// Create new recipient
	newRecipient := &share_domain.Recipient{
		ShareID:   strconv.FormatUint(uint64(in.ShareID	), 10),
		Name:      in.Name,
		Email:     in.Email,
		Role:      in.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		JoinedAt:  time.Now(),
	}

	if err := uc.repo.CreateRecipient(ctx, newRecipient); err != nil {
		return nil, fmt.Errorf("failed to add recipient: %w", err)
	}
	u64, err := strconv.ParseUint(newRecipient.ID, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	wd := uint(u64)
	// share_domain.RecipientAdded event
	uc.dispatcher.Dispatch(share_domain.RecipientAdded{
		BaseEvent: share_domain.BaseEvent{
			Name: "RecipientAdded",
			Time: time.Now(),
		},
		ShareID:     share.ID,
		RecipientID: newRecipient.ID,
		Email:       newRecipient.Email,
	})

	return &AddReceiverResult{
		ShareID:     in.ShareID,
		RecipientID: wd,
		Message:     "Recipient added successfully",
	}, nil
}
func stringToUint(str string) uint {
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	return uint(u64)
}
