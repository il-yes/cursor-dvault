package realtime_client_worker

import (
	"context"
	"fmt"
	"log"
	"time"

	realtime_client_infrastructure_persistence "vault-app/internal/realtime_client/infrastructure/persistence"
)

type OutboundDispatcherFunc func(ctx context.Context, item *realtime_client_infrastructure_persistence.OutboundDeliveryModel) error

type QueueRepository interface {
	RecoverPendingDeliveries(ctx context.Context) ([]realtime_client_infrastructure_persistence.OutboundDeliveryModel, error)
	MarkDelivering(ctx context.Context, envelopeID string, leaseDuration time.Duration) error
	MarkAcknowledged(ctx context.Context, eventID string) error
	RecordAttemptFailure(ctx context.Context, envelopeID string, lastErr string, maxAttempts int) error
}

type OutboundWorker struct {
	repo        QueueRepository
	dispatcher  OutboundDispatcherFunc
	interval    time.Duration
	maxAttempts int
	autoAck     bool
}

func NewOutboundWorker(repo QueueRepository, dispatcher OutboundDispatcherFunc) *OutboundWorker {
	return &OutboundWorker{
		repo:        repo,
		dispatcher:  dispatcher,
		interval:    2 * time.Second,
		maxAttempts: 5,
		autoAck:     false,
	}
}

func (w *OutboundWorker) SetInterval(interval time.Duration) *OutboundWorker {
	w.interval = interval
	return w
}

func (w *OutboundWorker) WithAutoAck(autoAck bool) *OutboundWorker {
	w.autoAck = autoAck
	return w
}

func (w *OutboundWorker) ProcessQueue(ctx context.Context) error {
	if w.repo == nil {
		return fmt.Errorf("queue repository is nil")
	}

	items, err := w.repo.RecoverPendingDeliveries(ctx)
	if err != nil {
		return fmt.Errorf("failed to recover pending deliveries: %w", err)
	}

	for _, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := w.repo.MarkDelivering(ctx, item.EnvelopeID, 30*time.Second); err != nil {
			log.Printf("[OUTBOX_WORKER] failed to mark delivering envelopeID=%s: %v", item.EnvelopeID, err)
			continue
		}

		var dispatchErr error
		if w.dispatcher != nil {
			dispatchErr = w.dispatcher(ctx, &item)
		}

		if dispatchErr == nil {
			if w.autoAck {
				if err := w.repo.MarkAcknowledged(ctx, item.EventID); err != nil {
					log.Printf("[OUTBOX_WORKER] failed to mark acknowledged eventID=%s: %v", item.EventID, err)
				} else {
					log.Printf("[OUTBOX_WORKER] successfully processed and acknowledged eventID=%s", item.EventID)
				}
			} else {
				log.Printf("[OUTBOX_WORKER] dispatch succeeded for envelopeID=%s eventID=%s, awaiting notification.ack", item.EnvelopeID, item.EventID)
			}
		} else {
			log.Printf("[OUTBOX_WORKER] dispatch failed envelopeID=%s err=%v", item.EnvelopeID, dispatchErr)
			_ = w.repo.RecordAttemptFailure(ctx, item.EnvelopeID, dispatchErr.Error(), w.maxAttempts)
		}
	}

	return nil
}

func (w *OutboundWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[OUTBOX_WORKER] stopping outbound worker")
			return
		case <-ticker.C:
			if err := w.ProcessQueue(ctx); err != nil {
				log.Printf("[OUTBOX_WORKER] error processing queue: %v", err)
			}
		}
	}
}
