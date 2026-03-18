package onboarding_persistence

import (
	"vault-app/internal/onboarding/domain"
	"gorm.io/gorm"
)

type AppStateRepository struct {
	db *gorm.DB
}

func NewAppStateRepository(db *gorm.DB) *AppStateRepository {
	return &AppStateRepository{db: db}
}

func (r *AppStateRepository) Get() (*onboarding_domain.AppState, error) {
	var appState onboarding_domain.AppState
	err := r.db.First(&appState).Error
	return &appState, err
}

func (r *AppStateRepository) Update(appState *onboarding_domain.AppState) error {
	return r.db.Save(appState).Error
}

func (r *AppStateRepository) Save(appState *onboarding_domain.AppState) error {
	return r.db.Create(appState).Error
}