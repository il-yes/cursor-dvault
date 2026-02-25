package identity_dtos


type EditUserInfosRequest struct {
	UserName string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}