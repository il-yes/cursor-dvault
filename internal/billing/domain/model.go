package billing_domain

import "time"

type PaymentMethod string

const (
	PaymentCard PaymentMethod = "card"
	PaymentCrypto PaymentMethod = "crypto"
	PaymentEncryptedCard PaymentMethod = "encrypted_card"
)

type BillingInstrument struct {
	ID string
	UserID string
	Type PaymentMethod
	EncryptedPayload string // when applicable
	CreatedAt time.Time
}