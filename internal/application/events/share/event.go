package share_application_events    

import share_domain "vault-app/internal/domain/shared"



type EventHandler func(event share_domain.DomainEvent)

type EventDispatcher interface {
    Register(eventName string, handler EventHandler)
    Dispatch(event share_domain.DomainEvent)
}
