package share_domain

import (
	"context"
	"errors"
)

var (
	ErrShareNotFound       = errors.New("share not found")
	ErrRecipientNotAllowed = errors.New("recipient is not part of the share")
	ErrShareExpired        = errors.New("share invitation expired")
)

type Repository interface {
	ListByUser(userID string) ([]ShareEntry, error)
	GetByID(id string) (*ShareEntry, error)
	Save(entry *ShareEntry) error
	Delete(id string) error

	ListReceivedByUser(recipientID string) ([]ShareEntry, error) // shared WITH me
	GetShareForAccept(shareID string, recipientID string) (*ShareEntry, *Recipient, []byte, error)
	// CreateShare(ctx context.Context, s *ShareEntry) error
	// GetShareForRecipient(ctx context.Context, shareID string, recipientID uint) (*ShareEntry, error)
	GetShareAndRecipient(ctx context.Context, shareID string, recipientID string) (*ShareEntry, *Recipient, error)
	MarkRecipientAccepted(ctx context.Context, recipientID string) error
	MarkRecipientRejected(ctx context.Context, recipientID string) error // NEW
    // Recipient management
    GetShareByID(ctx context.Context, shareID string) (*ShareEntry, error) 
    CreateRecipient(ctx context.Context, rec *Recipient) error  
	// ConfirmRecipient(ctx context.Context, shareID string, recipientID uint) error
}
