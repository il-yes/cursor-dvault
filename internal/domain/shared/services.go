package share_domain

func CanAddRecipient(share *ShareEntry, requesterID uint) bool {
    return share.OwnerID == requesterID
}
