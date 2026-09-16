package thread_usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	realtime_client_infrastructure_persistence "vault-app/internal/realtime_client/infrastructure/persistence"
	thread_domain "vault-app/internal/thread/domain"
)

type OutboundQueueSaver interface {
	SaveItem(ctx context.Context, item *realtime_client_infrastructure_persistence.OutboundDeliveryModel) error
}

type ListThreadEventsUsecase struct {
	Repo thread_domain.ThreadRepository
}

func NewListThreadEventsUsecase(repo thread_domain.ThreadRepository) *ListThreadEventsUsecase {
	return &ListThreadEventsUsecase{
		Repo: repo,
	}
}

func (uc *ListThreadEventsUsecase) Execute(ctx context.Context, threadID string) ([]thread_domain.ThreadEvent, error) {
	if uc.Repo == nil {
		return nil, errors.New("repository is required")
	}
	if threadID == "" {
		return nil, errors.New("thread id is required")
	}

	resp, err := uc.Repo.ListThreadEvents(ctx, &thread_domain.ListThreadEventsRequest{
		ThreadID: threadID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return []thread_domain.ThreadEvent{}, nil
	}

	return resp.Data, nil
}

type AppendThreadEventUsecase struct {
	Repo          thread_domain.ThreadRepository
	OutboundQueue OutboundQueueSaver
}

func NewAppendThreadEventUsecase(repo thread_domain.ThreadRepository) *AppendThreadEventUsecase {
	return &AppendThreadEventUsecase{
		Repo: repo,
	}
}

func (uc *AppendThreadEventUsecase) WithOutboundQueue(queue OutboundQueueSaver) *AppendThreadEventUsecase {
	uc.OutboundQueue = queue
	return uc
}

func (uc *AppendThreadEventUsecase) Execute(
	ctx context.Context,
	threadID string,
	eventType string,
	payload thread_domain.EventResourceRef,
	idempotencyKey ...string,
) (*thread_domain.ThreadEvent, error) {
	fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute enter threadID=%s eventType=%s refType=%s shareEntryID=%s trustGroupID=%s\n", threadID, eventType, payload.RefType, payload.ShareEntryID, payload.TrustGroupID)
	if uc.Repo == nil {
		fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute error: repo nil\n")
		return nil, errors.New("repository is required")
	}
	if threadID == "" {
		fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute error: threadID empty\n")
		return nil, errors.New("thread id is required")
	}
	if eventType == "" {
		fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute error: eventType empty\n")
		return nil, errors.New("event type is required")
	}

	if payload.RefType == thread_domain.ResourceShareEntry || (eventType == string(thread_domain.EventEntryShared) && payload.RefType != thread_domain.ResourceStorageAsset) {
		if strings.TrimSpace(payload.ShareEntryID) == "" {
			fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute error: entry.shared missing share_entry_id\n")
			return nil, errors.New("entry.shared thread event requires a non-empty share_entry_id")
		}
		if strings.TrimSpace(payload.TrustGroupID) == "" {
			fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute error: entry.shared missing trust_group_id\n")
			return nil, errors.New("entry.shared thread event requires a non-empty trust_group_id")
		}
		if payload.RefType == "" {
			payload.RefType = thread_domain.ResourceShareEntry
		}
	}

	key := ""
	if len(idempotencyKey) > 0 && idempotencyKey[0] != "" {
		key = idempotencyKey[0]
	} else if payload.RefType == thread_domain.ResourceShareEntry && payload.ShareEntryID != "" {
		key = "evt_share_" + payload.ShareEntryID
	}

	fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute calling Repo.AppendThreadEvent key=%s\n", key)
	resp, err := uc.Repo.AppendThreadEvent(ctx, &thread_domain.AppendThreadEventRequest{
		ThreadID:       threadID,
		EventType:      eventType,
		Payload:        payload,
		IdempotencyKey: key,
	})
	if err != nil {
		fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute Repo.AppendThreadEvent error=%v\n", err)
		return nil, err
	}
	if resp == nil {
		fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute Repo.AppendThreadEvent returned nil\n")
		return nil, errors.New("empty repository response")
	}

	fmt.Printf("[APPEND][STEP=05] AppendThreadEventUsecase.Execute success eventID=%s\n", resp.Data.ID)

	if uc.OutboundQueue != nil {
		payloadBytes, _ := json.Marshal(map[string]interface{}{
			"ref_type":       payload.RefType,
			"share_entry_id": payload.ShareEntryID,
			"trust_group_id": payload.TrustGroupID,
			"thread_id":      threadID,
			"event_id":       resp.Data.ID,
			"idempotency_key": key,
		})
		outboundItem := &realtime_client_infrastructure_persistence.OutboundDeliveryModel{
			EnvelopeID:  "env_" + resp.Data.ID,
			EventID:     resp.Data.ID,
			EventType:   eventType,
			Status:      "PENDING",
			PayloadJSON: payloadBytes,
			Attempts:    0,
			MaxAttempts: 5,
		}
		if err := uc.OutboundQueue.SaveItem(ctx, outboundItem); err != nil {
			fmt.Printf("[APPEND][OUTBOX] failed to enqueue outbound item: %v\n", err)
			return nil, fmt.Errorf("failed to enqueue outbound delivery item: %w", err)
		}
		fmt.Printf("[APPEND][OUTBOX] successfully enqueued outbound item eventID=%s status=PENDING\n", resp.Data.ID)
	}

	return &resp.Data, nil
}
