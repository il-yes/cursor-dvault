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
	if r == nil || r.db == nil {
		return &onboarding_domain.AppState{}, nil
	}
	var appState onboarding_domain.AppState
	err := r.db.First(&appState).Error
	return &appState, err
}

func (r *AppStateRepository) Update(appState *onboarding_domain.AppState) error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Save(appState).Error
}

func (r *AppStateRepository) Save(appState *onboarding_domain.AppState) error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Create(appState).Error
}