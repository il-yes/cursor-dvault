package onboarding_usecase

import onboarding_domain "vault-app/internal/onboarding/domain"


type FindUsersUseCaseInterface interface {
	Execute() ([]onboarding_domain.User, error)
	FindByEmail(email string) (*onboarding_domain.User, error)
}


type FindUsersUseCase struct {
	UserRepo onboarding_domain.UserRepository
}

func NewFindUsersUseCase(userRepo onboarding_domain.UserRepository) *FindUsersUseCase {
	return &FindUsersUseCase{UserRepo: userRepo}
}

func (f *FindUsersUseCase) Execute() ([]onboarding_domain.User, error) {
	return f.UserRepo.FindAll()
}



func (f *FindUsersUseCase) FindByEmail(email string) (*onboarding_domain.User, error) {
	return f.UserRepo.FindByEmail(email)
}

