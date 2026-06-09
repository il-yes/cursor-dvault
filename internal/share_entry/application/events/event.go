package share_entry_application_events

import share_entry_domain "vault-app/internal/share_entry/domain"



type EventHandler func(event share_entry_domain.DomainEvent)

type EventDispatcher interface {
    Register(eventName string, handler EventHandler)
    Dispatch(event share_entry_domain.DomainEvent)
}
