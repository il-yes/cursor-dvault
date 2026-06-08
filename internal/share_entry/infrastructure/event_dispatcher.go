package share_entry_infrastructure

import (
	"log"
	"sync"
	share_entry_application_events "vault-app/internal/share_entry/application/events"
	share_entry_domain "vault-app/internal/share_entry/domain"
	"vault-app/internal/utils"
)




type InMemoryEventDispatcher struct {
	handlers map[string][]share_entry_application_events.EventHandler
	mu       sync.RWMutex
}

func NewInMemoryEventDispatcher() *InMemoryEventDispatcher {
	return &InMemoryEventDispatcher{
		handlers: make(map[string][]share_entry_application_events.EventHandler),
	}
}

func (d *InMemoryEventDispatcher) Register(eventName string, handler share_entry_application_events.EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *InMemoryEventDispatcher) Dispatch(event share_entry_domain.DomainEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers := d.handlers[event.EventName()]
	for _, h := range handlers {
		go h(event) // async
	}
}

func InitializeEventDispatcher() share_entry_application_events.EventDispatcher {
	dispatcher := NewInMemoryEventDispatcher()
	dispatcher.Register("RecipientAdded", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.RecipientAdded)
		log.Printf("📩 New recipient added to share %d: %s\n", e.ShareID, e.Email)
	})
	dispatcher.Register("ShareAccepted", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.ShareAccepted)
		log.Printf("✔ Recipient %d accepted shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("ShareRejected", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.ShareRejected)
		log.Printf("❌ Recipient %d rejected shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("AccessRevoked", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.AccessRevoked)
		log.Printf("❌ Access revoked for recipient %d on shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("AccessRenewalRequested", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.AccessRenewalRequested)
		log.Printf("🔄 Access renewal requested for recipient %d on shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("AccessRenewalApproved", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.AccessRenewalApproved)
		log.Printf("✅ Access renewal approved for recipient %d on shared entry %d\n", e.RecipientID, e.ShareID)
	})
	dispatcher.Register("ShareCreated", func(evt share_entry_domain.DomainEvent) {
		e := evt.(share_entry_domain.ShareCreated)
		log.Printf("🎉 Share created: ShareID=%d, OwnerID=%d\n", e.ShareID, e.OwnerID)
		utils.LogPretty("InMemoryEventDispatcher - InitializeEventDispatcher - ShareCreated", e)
	})
	return dispatcher
}
