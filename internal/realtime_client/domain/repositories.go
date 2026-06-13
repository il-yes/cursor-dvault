package realtime_client_domain


type OffsetRepository interface {
	GetLastSeq(userID string) (uint64, error)
	SaveLastSeq(userID string, seq uint64) error
}

