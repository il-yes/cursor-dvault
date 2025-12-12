package subscription_persistence

import (
	"time"

	subscription_domain "vault-app/internal/subscription/domain"
)

type UserSubscriptionMapper struct {
	ID              string       `json:"id" gorm:"primaryKey"`
	Username        string    `gorm:"column:username" json:"username"`
	Email           string    `gorm:"column:email" json:"email"`
	Password        string    `gorm:"column:password" json:"password"`
	Role            string    `gorm:"column:role" json:"role"`
	CreatedAt       time.Time `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"varchar(100)"`
	LastConnectedAt time.Time `json:"last_connected_at" gorm:"last_connected_at"`
}

func (s *UserSubscriptionMapper) ToDomain() *subscription_domain.UserSubscription {
	return &subscription_domain.UserSubscription{
		ID:              s.ID,
		Username:        s.Username,
		Email:           s.Email,
		Password:        s.Password,
		Role:            s.Role,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		LastConnectedAt: s.LastConnectedAt,
	}
}

func UserSubscriptionToDB(s *subscription_domain.UserSubscription) *UserSubscriptionMapper {
	return &UserSubscriptionMapper{
		ID:              s.ID,
		Username:        s.Username,
		Email:           s.Email,
		Password:        s.Password,
		Role:            s.Role,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		LastConnectedAt: s.LastConnectedAt,
	}
}	


type SubscriptionMapper struct {
	ID            string        `json:"id"`
	Wallet        string        `json:"wallet,omitempty"` // only filled for crypto billing
	UserID        string        `json:"user_id"`          // optional, not required for crypto validation
	Tier          string        `json:"tier"`
	ExpiresAt     int64         `json:"expires_at"`
	Rail          string        `json:"rail"`              // "traditional" or "crypto"
	TxHash        string        `json:"tx_hash,omitempty"` // Tx hash confirmation payment
	Active        bool          `json:"active"`
	ActivatedAt   int64         `json:"activated_at"`
	Months        int           `json:"months"`
	PaymentMethod subscription_domain.PaymentMethod `json:"payment_method"`
	PaymentIntent string        `json:"payment_intent"`
	StartedAt     time.Time     `json:"started_at"`

	Features              subscription_domain.SubscriptionFeatures `json:"features"`
	Ledger                int32                `json:"ledger"`
	BillingCycle          string               `json:"billing_cycle"`
	TrialEndsAt           int64                `json:"trial_ends_at"`
	NextBillingDate       int64                `json:"next_billing_date"`
	Price                 float64              `json:"price"`
	StripeSubscriptionID  string               `json:"stripe_subscription_id"`
	StellarPaymentAddress string               `json:"stellar_payment_address"`
	PaymentFlowType       string               `json:"payment_flow_type"`
	StellarScheduleID     string               `json:"stellar_schedule_id"`
	Status                string               `json:"status"`
	EndedAt               int64                `json:"ended_at"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}
func (s *SubscriptionMapper) ToDomain() *subscription_domain.Subscription {
	return &subscription_domain.Subscription{
		ID:            s.ID,
		Wallet:        s.Wallet,
		UserID:        s.UserID,
		Tier:          s.Tier,
		ExpiresAt:     s.ExpiresAt,
		Rail:          s.Rail,
		TxHash:        s.TxHash,
		Active:        s.Active,
		ActivatedAt:   s.ActivatedAt,
		Months:        s.Months,
		PaymentMethod: s.PaymentMethod,
		PaymentIntent: s.PaymentIntent,
		StartedAt:     s.StartedAt,
		Features:              s.Features,
		Ledger:                s.Ledger,
		BillingCycle:          s.BillingCycle,
		TrialEndsAt:           s.TrialEndsAt,
		NextBillingDate:       s.NextBillingDate,
		Price:                 s.Price,
		StripeSubscriptionID:  s.StripeSubscriptionID,
		StellarPaymentAddress: s.StellarPaymentAddress,
		PaymentFlowType:       s.PaymentFlowType,
		StellarScheduleID:     s.StellarScheduleID,
		Status:                s.Status,
		EndedAt:               s.EndedAt,
		CreatedAt:             s.CreatedAt,
		UpdatedAt:             s.UpdatedAt,
	}
}	
func SubscriptionToDB(s *subscription_domain.Subscription) *SubscriptionMapper {
	return &SubscriptionMapper{
		ID:            s.ID,
		Wallet:        s.Wallet,
		UserID:        s.UserID,
		Tier:          s.Tier,
		ExpiresAt:     s.ExpiresAt,
		Rail:          s.Rail,
		TxHash:        s.TxHash,
		Active:        s.Active,
		ActivatedAt:   s.ActivatedAt,
		Months:        s.Months,
		PaymentMethod: s.PaymentMethod,
		PaymentIntent: s.PaymentIntent,
		StartedAt:     s.StartedAt,
		Features:              s.Features,
		Ledger:                s.Ledger,
		BillingCycle:          s.BillingCycle,
		TrialEndsAt:           s.TrialEndsAt,
		NextBillingDate:       s.NextBillingDate,
		Price:                 s.Price,
		StripeSubscriptionID:  s.StripeSubscriptionID,
		StellarPaymentAddress: s.StellarPaymentAddress,
		PaymentFlowType:       s.PaymentFlowType,
		StellarScheduleID:     s.StellarScheduleID,
		Status:                s.Status,
		EndedAt:               s.EndedAt,
		CreatedAt:             s.CreatedAt,
		UpdatedAt:             s.UpdatedAt,
	}
}