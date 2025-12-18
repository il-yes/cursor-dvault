package auth_domain



type UserRepository interface {
	Save(user *Principal) error
	FindByID(id string) (*Principal, error)
	FindByEmail(email string) (*Principal, error)	
	FindByUsername(username string) (*Principal, error)	
}

type AuthRepository interface {
	Save(token *TokenPairs) error
	FindByID(id string) (*TokenPairs, error)
	FindByEmail(email string) (*TokenPairs, error)	
	FindByUsername(username string) (*TokenPairs, error)	
}
