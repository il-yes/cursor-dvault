package realtime_client_infrastructure_persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OutboundDeliveryModel struct {
	EnvelopeID       string `gorm:"primaryKey"`
	EventID          string `gorm:"index"`
	EventType        string
	SenderVaultID    string
	RecipientVaultID string
	Sequence         uint64
	Status           string `gorm:"index"` // "PENDING", "DELIVERING", "RETRY_PENDING", "ACKNOWLEDGED", "FAILED"
	PayloadJSON      []byte
	Attempts         int
	MaxAttempts      int
	LeaseExpiresAt   *time.Time `gorm:"index"`
	AcknowledgedAt   *time.Time
	LastError        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (OutboundDeliveryModel) TableName() string {
	return "federation_outbound_queue"
}

type GormOutboundQueueRepository struct {
	db *gorm.DB
}

func NewGormOutboundQueueRepository(db *gorm.DB) *GormOutboundQueueRepository {
	return &GormOutboundQueueRepository{db: db}
}

func (r *GormOutboundQueueRepository) SaveItem(ctx context.Context, item *OutboundDeliveryModel) error {
	if item == nil || item.EnvelopeID == "" {
		return errors.New("invalid outbound delivery item")
	}
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GormOutboundQueueRepository) GetItem(ctx context.Context, envelopeID string) (*OutboundDeliveryModel, error) {
	var model OutboundDeliveryModel
	if err := r.db.WithContext(ctx).First(&model, "envelope_id = ?", envelopeID).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *GormOutboundQueueRepository) MarkDelivering(ctx context.Context, envelopeID string, leaseDuration time.Duration) error {
	now := time.Now()
	leaseExpires := now.Add(leaseDuration)
	return r.db.WithContext(ctx).Model(&OutboundDeliveryModel{}).
		Where("envelope_id = ?", envelopeID).
		Updates(map[string]any{
			"status":           "DELIVERING",
			"lease_expires_at": leaseExpires,
			"updated_at":       now,
		}).Error
}

func (r *GormOutboundQueueRepository) MarkAcknowledged(ctx context.Context, eventID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&OutboundDeliveryModel{}).
		Where("event_id = ?", eventID).
		Updates(map[string]any{
			"status":           "ACKNOWLEDGED",
			"acknowledged_at":  now,
			"lease_expires_at": nil,
			"updated_at":       now,
		}).Error
}

func (r *GormOutboundQueueRepository) RecordAttemptFailure(ctx context.Context, envelopeID string, lastErr string, maxAttempts int) error {
	var item OutboundDeliveryModel
	if err := r.db.WithContext(ctx).First(&item, "envelope_id = ?", envelopeID).Error; err != nil {
		return err
	}

	item.Attempts++
	item.LastError = lastErr
	item.UpdatedAt = time.Now()
	item.LeaseExpiresAt = nil

	if maxAttempts > 0 && item.Attempts >= maxAttempts {
		item.Status = "FAILED" // Terminal failure state
	} else {
		item.Status = "RETRY_PENDING"
	}

	return r.db.WithContext(ctx).Save(&item).Error
}

// RecoverPendingDeliveries queries entries where:
// Status = 'PENDING' OR Status = 'RETRY_PENDING' OR (Status = 'DELIVERING' AND (lease_expires_at IS NULL OR lease_expires_at < NOW()))
func (r *GormOutboundQueueRepository) RecoverPendingDeliveries(ctx context.Context) ([]OutboundDeliveryModel, error) {
	now := time.Now()
	var items []OutboundDeliveryModel

	err := r.db.WithContext(ctx).
		Where("status IN ('PENDING', 'RETRY_PENDING') OR (status = 'DELIVERING' AND (lease_expires_at IS NULL OR lease_expires_at < ?))", now).
		Order("sequence ASC, created_at ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
