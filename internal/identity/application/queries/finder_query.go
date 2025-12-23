package identity_queries

import (
	"context"
	utils "vault-app/internal"
	identity_domain "vault-app/internal/identity/domain"
)

type ReqQuery struct {
	Email string
}

type FinderQueryHandler struct {
	IdentityRepository identity_domain.UserRepository
}
func NewFinderQueryHandler(
	identityRepository identity_domain.UserRepository,
) *FinderQueryHandler {
	return &FinderQueryHandler{
		IdentityRepository: identityRepository,
	}
}		
func (h *FinderQueryHandler) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	return h.IdentityRepository.FindByEmail(ctx, email)
}
func (h *FinderQueryHandler) FindById(ctx context.Context, id string) (*identity_domain.User, error) {
	utils.LogPretty("FinderQueryHandler - FindById - id", id)
	return h.IdentityRepository.FindByID(ctx, id)
}