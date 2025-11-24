package share_domain

import "time"

type DomainEvent interface {
    EventName() string
    OccurredAt() time.Time
}

type BaseEvent struct {
    Name string
    Time time.Time
}

func (e BaseEvent) EventName() string  { return e.Name }
func (e BaseEvent) OccurredAt() time.Time { return e.Time }

// ---------------------------
// SHARE DOMAIN EVENTS
// ---------------------------

type ShareCreated struct {
    BaseEvent
    ShareID   string
    OwnerID   uint
}

type RecipientAdded struct {
    BaseEvent
    ShareID     string
    RecipientID string
    Email       string
}

type ShareAccepted struct {
    BaseEvent
    ShareID     string
    RecipientID string
}

type ShareRejected struct {
    BaseEvent
    ShareID     string
    RecipientID string
}

type AccessRevoked struct {
    BaseEvent
    ShareID     string
    RecipientID string
}

type AccessRenewalRequested struct {
    BaseEvent
    ShareID     string
    RecipientID string
}

type AccessRenewalApproved struct {
    BaseEvent
    ShareID     string
    RecipientID string
}
