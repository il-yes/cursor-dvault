package tracecore

import "time"

type LoginRequest struct {
    Email         string `json:"email,omitempty"`
    Password      string `json:"password,omitempty"`
    PublicKey     string `json:"public_key,omitempty"`
    SignedMessage string `json:"signed_message,omitempty"`
    Signature     string `json:"signature,omitempty"`
}

type LoginResponse struct {
    UserID      int64  `json:"user_id"`
    Email       string `json:"email"`
    Token       string `json:"token"`
    VaultCID    string `json:"vault_cid,omitempty"`
    VaultDirty  bool   `json:"vault_dirty"`
}

type VaultEntry struct {
    ID        string    `json:"id"`
    EntryName string    `json:"entry_name"`
    Type      string    `json:"type"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ShareEntry struct {
    ID        string    `json:"id"`
    EntryRef  string    `json:"entry_ref"`
    EntryName string    `json:"entry_name"`
    Status    string    `json:"status"`
    SharedAt  time.Time `json:"shared_at"`
}
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	Password  string `json:"password"` // if needed
}
