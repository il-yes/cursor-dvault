package identity_persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StellarAuthChallenge struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	PublicKey string    `gorm:"uniqueIndex;column:public_key;not null" json:"public_key"`
	Challenge string    `gorm:"column:challenge;not null" json:"challenge"`
	ExpiresAt time.Time `gorm:"column:expires_at;index;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (StellarAuthChallenge) TableName() string {
	return "stellar_auth_challenges"
}

type GormStellarChallengeRepository struct {
	db *gorm.DB
}

func NewGormStellarChallengeRepository(db *gorm.DB) *GormStellarChallengeRepository {
	return &GormStellarChallengeRepository{db: db}
}

func (r *GormStellarChallengeRepository) CreateOrReplace(ctx context.Context, publicKey string, challengeStr string, ttl time.Duration) (*StellarAuthChallenge, error) {
	if publicKey == "" || challengeStr == "" {
		return nil, errors.New("public_key and challenge are required")
	}

	var challengeRecord *StellarAuthChallenge
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete any previous challenge for this public_key (replacement invariant)
		if err := tx.Where("public_key = ?", publicKey).Delete(&StellarAuthChallenge{}).Error; err != nil {
			return fmt.Errorf("failed to replace existing challenge: %w", err)
		}

		now := time.Now().UTC()
		record := StellarAuthChallenge{
			ID:        uuid.New().String(),
			PublicKey: publicKey,
			Challenge: challengeStr,
			ExpiresAt: now.Add(ttl),
			CreatedAt: now,
		}

		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to persist challenge: %w", err)
		}

		challengeRecord = &record
		return nil
	})

	if err != nil {
		return nil, err
	}
	return challengeRecord, nil
}

func (r *GormStellarChallengeRepository) Consume(ctx context.Context, publicKey string) (string, error) {
	if publicKey == "" {
		return "", errors.New("public_key is required")
	}

	var consumedChallenge string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record StellarAuthChallenge

		// Query for challenge matching public_key with row locking if supported
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("public_key = ?", publicKey).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("challenge not found or already consumed")
			}
			// Fallback query without locking for SQLite dialects that don't support UPDATE lock clause
			if errFallback := tx.Where("public_key = ?", publicKey).First(&record).Error; errFallback != nil {
				return errors.New("challenge not found or already consumed")
			}
		}

		// Verify 5-minute expiry
		if time.Now().UTC().After(record.ExpiresAt.UTC()) {
			tx.Delete(&record)
			return errors.New("challenge expired")
		}

		// Atomically delete the consumed challenge (single-use invariant)
		if err := tx.Delete(&record).Error; err != nil {
			return fmt.Errorf("failed to consume challenge: %w", err)
		}

		consumedChallenge = record.Challenge
		return nil
	})

	if err != nil {
		return "", err
	}
	return consumedChallenge, nil
}
