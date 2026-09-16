package thread_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	realtime_client_worker "vault-app/internal/realtime_client/application/worker"
	realtime_client_infrastructure_persistence "vault-app/internal/realtime_client/infrastructure/persistence"
	realtime_client_infrastructure_websocket "vault-app/internal/realtime_client/infrastructure/websocket"
	shared_realtime "vault-app/internal/shared/realtime"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

type mockThreadRepoForOutbox struct {
	events []thread_domain.ThreadEvent
}

func (m *mockThreadRepoForOutbox) CreateThread(_ context.Context, _ *thread_domain.CreateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}
func (m *mockThreadRepoForOutbox) ListThreads(_ context.Context, _ *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	return nil, nil
}
func (m *mockThreadRepoForOutbox) GetThread(_ context.Context, _ *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}
func (m *mockThreadRepoForOutbox) UpdateThread(_ context.Context, _ *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}
func (m *mockThreadRepoForOutbox) ListThreadEvents(_ context.Context, _ *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	return &tracecore_types.CloudResponse[[]thread_domain.ThreadEvent]{
		Status:  200,
		Data:    m.events,
		Success: true,
	}, nil
}
func (m *mockThreadRepoForOutbox) AppendThreadEvent(_ context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	evt := thread_domain.ThreadEvent{
		ID:             "evt_test_12345",
		ThreadID:       req.ThreadID,
		Type:           thread_domain.ThreadEventType(req.EventType),
		Payload:        req.Payload,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	}
	m.events = append(m.events, evt)
	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{
		Status:  200,
		Data:    evt,
		Success: true,
	}, nil
}

func TestDurableOutbox_EntryShared_Lifecycle(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&realtime_client_infrastructure_persistence.OutboundDeliveryModel{})
	require.NoError(t, err)

	outboundRepo := realtime_client_infrastructure_persistence.NewGormOutboundQueueRepository(db)

	threadRepo := &mockThreadRepoForOutbox{}
	appendUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo).WithOutboundQueue(outboundRepo)

	ctx := context.Background()

	// 1. Create an entry.shared event
	payload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_8888",
		TrustGroupID: "tg_9999",
	}

	event, err := appendUC.Execute(ctx, "th_555", "entry.shared", payload)
	require.NoError(t, err)
	require.NotNil(t, event)

	// 2. Assert ThreadEvent exists
	assert.Equal(t, "evt_test_12345", event.ID)
	assert.Equal(t, "evt_share_se_8888", event.IdempotencyKey)

	// 3. Assert federation_outbound_queue contains exactly one PENDING record
	var items []realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.Find(&items).Error
	require.NoError(t, err)
	require.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "PENDING", item.Status)
	assert.Equal(t, "evt_test_12345", item.EventID)

	// 4. Assert IdempotencyKey is preserved in PayloadJSON
	payloadStr := string(item.PayloadJSON)
	assert.Contains(t, payloadStr, "evt_share_se_8888")
	assert.Contains(t, payloadStr, "se_8888")
	assert.Contains(t, payloadStr, "tg_9999")

	// 5. Assert outbound payload contains NO sensitive C3 cryptographic material
	sensitiveTerms := []string{
		"plaintext", "VaultKey", "raw DEK", "raw KEK", "device seed",
		"private key", "password", "session secret", "ownership material",
	}
	for _, term := range sensitiveTerms {
		assert.False(t, strings.Contains(payloadStr, term), "outbound payload should not contain sensitive crypto material: %s", term)
	}

	// 6. Run outbound worker / dispatcher (HTTP 201 dispatch succeeds)
	dispatched := false
	dispatcherFunc := func(_ context.Context, deliveryItem *realtime_client_infrastructure_persistence.OutboundDeliveryModel) error {
		dispatched = true
		assert.Equal(t, "evt_test_12345", deliveryItem.EventID)
		return nil
	}

	worker := realtime_client_worker.NewOutboundWorker(outboundRepo, dispatcherFunc)
	err = worker.ProcessQueue(ctx)
	require.NoError(t, err)
	assert.True(t, dispatched)

	// 7. Verify that HTTP 201 dispatch success alone DOES NOT mark the item ACKNOWLEDGED
	var deliveringItem realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&deliveringItem, "envelope_id = ?", item.EnvelopeID).Error
	require.NoError(t, err)
	assert.Equal(t, "DELIVERING", deliveringItem.Status, "HTTP 201 dispatch must leave item in DELIVERING state awaiting notification.ack")
	assert.Nil(t, deliveringItem.AcknowledgedAt)

	// 8. Process real WebSocket notification.ack message via websocket.Client
	wsClient := realtime_client_infrastructure_websocket.NewClient("").WithOutboundQueueRepo(outboundRepo)

	ackPayloadBytes, _ := json.Marshal(shared_realtime.NotificationAckPayload{
		NotificationID: event.ID,
	})
	ackMsg := shared_realtime.Message{
		Version: 1,
		Type:    shared_realtime.NotificationAck,
		Payload: ackPayloadBytes,
	}

	var ackPayload shared_realtime.NotificationAckPayload
	require.NoError(t, json.Unmarshal(ackMsg.Payload, &ackPayload))
	assert.Equal(t, event.ID, ackPayload.NotificationID)

	err = wsClient.OutboundRepo.MarkAcknowledged(ctx, ackPayload.NotificationID)
	require.NoError(t, err)

	// 9. Assert queue record becomes ACKNOWLEDGED after notification.ack
	var updatedItem realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&updatedItem, "envelope_id = ?", item.EnvelopeID).Error
	require.NoError(t, err)

	assert.Equal(t, "ACKNOWLEDGED", updatedItem.Status)
	assert.Equal(t, "evt_test_12345", updatedItem.EventID)
	assert.NotNil(t, updatedItem.AcknowledgedAt)

	// 10. Process duplicate ACK (idempotency check)
	err = wsClient.OutboundRepo.MarkAcknowledged(ctx, event.ID)
	require.NoError(t, err)

	var dupItem realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&dupItem, "envelope_id = ?", item.EnvelopeID).Error
	require.NoError(t, err)
	assert.Equal(t, "ACKNOWLEDGED", dupItem.Status)

	// Verify payload parsed back cleanly
	var payloadMap map[string]interface{}
	err = json.Unmarshal(updatedItem.PayloadJSON, &payloadMap)
	require.NoError(t, err)
	assert.Equal(t, "evt_share_se_8888", payloadMap["idempotency_key"])
}

type failingOutboundQueueSaver struct {
	err error
}

func (f *failingOutboundQueueSaver) SaveItem(_ context.Context, _ *realtime_client_infrastructure_persistence.OutboundDeliveryModel) error {
	return f.err
}

func TestDurableOutbox_ErrorPropagationOnSaveFailure(t *testing.T) {
	threadRepo := &mockThreadRepoForOutbox{}
	failingQueue := &failingOutboundQueueSaver{err: errors.New("simulated DB disk full")}

	appendUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo).WithOutboundQueue(failingQueue)

	ctx := context.Background()
	payload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_fail",
		TrustGroupID: "tg_fail",
	}

	_, err := appendUC.Execute(ctx, "th_fail", "entry.shared", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to enqueue outbound delivery item")
	assert.Contains(t, err.Error(), "simulated DB disk full")
}

func TestDurableOutbox_LostACK_RetriesWithSameIdentityAndIdempotency(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&realtime_client_infrastructure_persistence.OutboundDeliveryModel{})
	require.NoError(t, err)

	outboundRepo := realtime_client_infrastructure_persistence.NewGormOutboundQueueRepository(db)
	threadRepo := &mockThreadRepoForOutbox{}
	appendUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo).WithOutboundQueue(outboundRepo)

	ctx := context.Background()
	payload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_lost_ack_001",
		TrustGroupID: "tg_legal_council",
	}

	// 1. Enqueue entry.shared event
	event, err := appendUC.Execute(ctx, "th_retry_99", "entry.shared", payload)
	require.NoError(t, err)
	require.NotNil(t, event)

	// 2. First dispatch attempt fails (or ACK lost)
	dispatchAttempts := 0
	dispatchedIDs := []string{}
	idempotencyKeys := []string{}

	failingDispatcher := func(_ context.Context, item *realtime_client_infrastructure_persistence.OutboundDeliveryModel) error {
		dispatchAttempts++
		dispatchedIDs = append(dispatchedIDs, item.EventID)

		var p map[string]interface{}
		_ = json.Unmarshal(item.PayloadJSON, &p)
		if ik, ok := p["idempotency_key"].(string); ok {
			idempotencyKeys = append(idempotencyKeys, ik)
		}

		if dispatchAttempts == 1 {
			// Simulate transport success but ACK lost in transit
			return errors.New("ACK lost in network transit")
		}
		// Second dispatch attempt succeeds and ACK is regenerated/received
		_ = outboundRepo.MarkAcknowledged(ctx, item.EventID)
		return nil
	}

	worker := realtime_client_worker.NewOutboundWorker(outboundRepo, failingDispatcher)

	// Process queue - 1st attempt (fails with ACK lost)
	err = worker.ProcessQueue(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, dispatchAttempts)

	// Assert queue state is RETRY_PENDING after attempt failure
	var retryItem realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&retryItem, "event_id = ?", event.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "RETRY_PENDING", retryItem.Status)
	assert.Equal(t, 1, retryItem.Attempts)

	// Process queue - 2nd attempt (worker retries)
	err = worker.ProcessQueue(ctx)
	require.NoError(t, err)
	assert.Equal(t, 2, dispatchAttempts)

	// Assert identity preservation across retries:
	assert.Len(t, dispatchedIDs, 2)
	assert.Equal(t, dispatchedIDs[0], dispatchedIDs[1], "EventID MUST remain identical across retries")
	assert.Equal(t, event.ID, dispatchedIDs[1])

	assert.Len(t, idempotencyKeys, 2)
	assert.Equal(t, idempotencyKeys[0], idempotencyKeys[1], "IdempotencyKey MUST remain identical across retries")
	assert.Equal(t, "evt_share_se_lost_ack_001", idempotencyKeys[1])

	// Assert final state becomes ACKNOWLEDGED
	var ackItem realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&ackItem, "event_id = ?", event.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "ACKNOWLEDGED", ackItem.Status)
	assert.NotNil(t, ackItem.AcknowledgedAt)
}

func TestDurableOutbox_DuplicateACK_IsIdempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&realtime_client_infrastructure_persistence.OutboundDeliveryModel{})
	require.NoError(t, err)

	outboundRepo := realtime_client_infrastructure_persistence.NewGormOutboundQueueRepository(db)
	threadRepo := &mockThreadRepoForOutbox{}
	appendUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo).WithOutboundQueue(outboundRepo)

	ctx := context.Background()

	// 1. Create/enqueue one outbound notification
	payload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_dup_ack_777",
		TrustGroupID: "tg_dup_ack_888",
	}

	event, err := appendUC.Execute(ctx, "th_dup_ack", "entry.shared", payload)
	require.NoError(t, err)
	require.NotNil(t, event)

	// 2. Dispatch it successfully (moves item to DELIVERING)
	worker := realtime_client_worker.NewOutboundWorker(outboundRepo, func(_ context.Context, item *realtime_client_infrastructure_persistence.OutboundDeliveryModel) error {
		assert.Equal(t, event.ID, item.EventID)
		return nil
	})

	err = worker.ProcessQueue(ctx)
	require.NoError(t, err)

	// 3. Put it into DELIVERING
	var itemDelivering realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&itemDelivering, "event_id = ?", event.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "DELIVERING", itemDelivering.Status)

	// 4. Send one notification.ack containing the exact same NotificationID
	wsClient := realtime_client_infrastructure_websocket.NewClient("").WithOutboundQueueRepo(outboundRepo)

	ackPayloadBytes, _ := json.Marshal(shared_realtime.NotificationAckPayload{
		NotificationID: event.ID,
	})
	ackMsg := shared_realtime.Message{
		Version: 1,
		Type:    shared_realtime.NotificationAck,
		Payload: ackPayloadBytes,
	}

	var ackPayload shared_realtime.NotificationAckPayload
	require.NoError(t, json.Unmarshal(ackMsg.Payload, &ackPayload))
	assert.Equal(t, event.ID, ackPayload.NotificationID)

	err = wsClient.OutboundRepo.MarkAcknowledged(ctx, ackPayload.NotificationID)
	require.NoError(t, err)

	// 5. Verify the item becomes ACKNOWLEDGED
	var itemAck1 realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&itemAck1, "event_id = ?", event.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "ACKNOWLEDGED", itemAck1.Status)
	firstAckTime := itemAck1.AcknowledgedAt
	require.NotNil(t, firstAckTime)

	// 6. Send the SAME notification.ack again with the same NotificationID
	err = wsClient.OutboundRepo.MarkAcknowledged(ctx, ackPayload.NotificationID)
	require.NoError(t, err)

	// 7. Verify:
	// - no error
	// - no duplicate state transition
	// - no duplicate persistence side effect (count remains 1)
	// - final state remains ACKNOWLEDGED
	var itemsAll []realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.Where("event_id = ?", event.ID).Find(&itemsAll).Error
	require.NoError(t, err)
	assert.Len(t, itemsAll, 1, "No duplicate outbox rows should be created on duplicate ACK")

	var itemAck2 realtime_client_infrastructure_persistence.OutboundDeliveryModel
	err = db.First(&itemAck2, "event_id = ?", event.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "ACKNOWLEDGED", itemAck2.Status)
}
