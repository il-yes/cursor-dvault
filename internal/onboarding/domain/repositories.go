package onboarding_domain



type UserRepository interface {
	Create(user *User)  (*User, error)
	GetByID(id string) (*User, error)
	List() ([]User, error)
	Update(user *User) error
	Delete(id string) error	
}