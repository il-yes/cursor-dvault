package shared_realtime

type NotificationPayload struct {
    ID string `json:"id"`
    Title          string `json:"title"`
    Body           string `json:"body"`
    CreatedAt      string `json:"created_at"`
	
}

type ShareInvitationNotificationPayload struct {
    NotificationPayload
	ShareID        string `json:"share_id"`
	IntentID       string `json:"intent_id"`
	RecipientEmail string `json:"recipient_email"`
}

type ShareAcceptedPayload struct {
    NotificationPayload
	ShareID        string `json:"share_id"`
	IntentID       string `json:"intent_id"`
	RecipientEmail string `json:"recipient_email"`
}
type ShareRejectedPayload struct {
    NotificationPayload
	ShareID        string `json:"share_id"`
	IntentID       string `json:"intent_id"`
	RecipientEmail string `json:"recipient_email"`
}

type ShareReadyToAcceptPayload struct {
    NotificationPayload
	ShareID        string `json:"share_id"`
	IntentID       string `json:"intent_id"`
	RecipientEmail string `json:"recipient_email"`
}

type SubscriptionActivatedPayload struct {
    NotificationPayload
	SubscriptionID string `json:"subscription_id"`
	UserID         string `json:"user_id"`
}
