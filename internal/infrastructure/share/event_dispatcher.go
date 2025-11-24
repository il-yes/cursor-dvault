package share_infrastructure

import (
	"log"
	"sync"
	share_application "vault-app/internal/application/events/share"
	share_domain "vault-app/internal/domain/shared"
)

type InMemoryEventDispatcher struct {
	handlers map[string][]share_application.EventHandler
	mu       sync.RWMutex
}

func NewInMemoryEventDispatcher() *InMemoryEventDispatcher {
	return &InMemoryEventDispatcher{
		handlers: make(map[string][]share_application.EventHandler),
	}
}

func (d *InMemoryEventDispatcher) Register(eventName string, handler share_application.EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *InMemoryEventDispatcher) Dispatch(event share_domain.DomainEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers := d.handlers[event.EventName()]
	for _, h := range handlers {
		go h(event) // async
	}
}

func InitializeEventDispatcher() share_application.EventDispatcher {
	dispatcher := NewInMemoryEventDispatcher()
	dispatcher.Register("RecipientAdded", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.RecipientAdded)
		log.Printf("📩 New recipient added to share %d: %s\n", e.ShareID, e.Email)
	})
	dispatcher.Register("ShareAccepted", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.ShareAccepted)
		log.Printf("✔ Recipient %d accepted shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("ShareRejected", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.ShareRejected)
		log.Printf("❌ Recipient %d rejected shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("AccessRevoked", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.AccessRevoked)
		log.Printf("❌ Access revoked for recipient %d on shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("AccessRenewalRequested", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.AccessRenewalRequested)
		log.Printf("🔄 Access renewal requested for recipient %d on shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("AccessRenewalApproved", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.AccessRenewalApproved)
		log.Printf("✅ Access renewal approved for recipient %d on shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("ShareCreated", func(evt share_domain.DomainEvent) {
		e := evt.(share_domain.ShareCreated)
		log.Printf("🎉 Share created: ShareID=%d, OwnerID=%d\n", e.ShareID, e.OwnerID)
	})
	return dispatcher
}
