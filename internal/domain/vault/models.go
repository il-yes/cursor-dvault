package vault_domain

import (
	"time"

)



type Recipient struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    ShareID   uint      `json:"share_id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    JoinedAt  time.Time `json:"joined_at"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
