package identity_ui



type IdentityHandler struct {
	LoginHandler *LoginHandler
}


func NewIdentityHandler(loginHandler *LoginHandler) *IdentityHandler {
	return &IdentityHandler{
		LoginHandler: loginHandler,
	}
}
