package share_domain

func CanAddRecipient(share *ShareEntry, requesterID string) bool {
    return share.OwnerID == requesterID
}
