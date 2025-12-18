package auth_persistence

import (
	auth_domain "vault-app/internal/auth/domain"

	"gorm.io/gorm"
)

type GormAuthRepository struct {
	db *gorm.DB
}

func NewGormAuthRepository(db *gorm.DB) *GormAuthRepository {
	return &GormAuthRepository{db: db}
}

func (r *GormAuthRepository) Save(tp *auth_domain.TokenPairs) error {
	
	return r.db.Create(TokenPairsToModel(*tp)).Error
}

func (m *GormAuthRepository) SaveJwtToken(tokens auth_domain.TokenPairs) (*auth_domain.TokenPairs, error) {
	// Persist only the refresh token (access tokens are short-lived; no need to store them long-term)
	if err := m.db.Model(&auth_domain.TokenPairs{}).
		Create(&tokens).Error; err != nil {
		return nil, err
	}
	return &tokens, nil
}	
func (m *GormAuthRepository) FindByEmail(email string) (*auth_domain.TokenPairs, error) {
	var tp auth_domain.TokenPairs
	if err := m.db.Where("email = ?", email).First(&tp).Error; err != nil {
		return nil, err
	}
	return &tp, nil
}
func (m *GormAuthRepository) FindByID(id string) (*auth_domain.TokenPairs, error) {
	var tp auth_domain.TokenPairs
	if err := m.db.Where("id = ?", id).First(&tp).Error; err != nil {
		return nil, err
	}
	return &tp, nil
}
func (m *GormAuthRepository) FindByUsername(username string) (*auth_domain.TokenPairs, error) {
	var tp auth_domain.TokenPairs
	if err := m.db.Where("username = ?", username).First(&tp).Error; err != nil {
		return nil, err
	}
	return &tp, nil
}	