package subscription_infrastructure

import "context"


type SubscriptionService struct {
	
}	

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{}
}
func (s *SubscriptionService) UpdateStorage(ctx context.Context, userID string, gb int) error {
	return nil
}
func (s *SubscriptionService) EnableCloudBackup(ctx context.Context, userID string, enabled bool) error {
	return nil
}
func (s *SubscriptionService) EnableVersionHistory(ctx context.Context, userID string, days int) error {
	return nil
}
func (s *SubscriptionService) EnableTracecore(ctx context.Context, userID string) error {
	return nil
}
	