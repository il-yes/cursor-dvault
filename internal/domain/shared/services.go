package share_domain

import "context"

func CanAddRecipient(share *ShareEntry, requesterID string) bool {
    return share.OwnerID == requesterID
}

type OnSharingService interface {
	Build(ctx context.Context, share *ShareEntry) (string, error)
}
	