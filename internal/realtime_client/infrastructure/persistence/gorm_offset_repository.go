package realtime_client_infrastructure_persistence

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"vault-app/internal/realtime_client/domain"
)

type RealtimeOffsetModel struct {
	UserID    string `gorm:"primaryKey"`
	LastSeq   uint64
	UpdatedAt time.Time
}

func (RealtimeOffsetModel) TableName() string {
	return "realtime_offsets"
}

type gormOffsetRepository struct {
	db *gorm.DB
}

func NewGORMOffsetRepository(db *gorm.DB) realtime_client_domain.OffsetRepository {
	_ = db.AutoMigrate(&RealtimeOffsetModel{})
	return &gormOffsetRepository{db: db}
}

// GetLastSeq implements realtime_client_domain.OffsetRepository
func (r *gormOffsetRepository) GetLastSeq(userID string) (uint64, error) {
	var offset RealtimeOffsetModel

	result := r.db.Where("user_id = ?", userID).First(&offset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return offset.LastSeq, nil
}

// SaveLastSeq implements realtime_client_domain.OffsetRepository
func (r *gormOffsetRepository) SaveLastSeq(userID string, seq uint64) error {
	now := time.Now()
	offset := RealtimeOffsetModel{
		UserID:    userID,
		LastSeq:   seq,
		UpdatedAt: now,
	}

	return r.db.Save(&offset).Error
}
