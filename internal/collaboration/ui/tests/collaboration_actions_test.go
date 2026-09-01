package collaboration_ui_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_domain "vault-app/internal/collaboration/domain"
	collaboration_infra "vault-app/internal/collaboration/infrastructure"
	collaboration_ui "vault-app/internal/collaboration/ui"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
)

// Helper setup for actions test
func setupActionsHandler(threadRepo *stubThreadRepo) (*collaboration_ui.CollaborationHandler, *collaboration_infra.MemoryActionRepository) {
	actionRepo := collaboration_infra.NewMemoryActionRepository()
	appendEventUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo)
	actionUseCases := collaboration_usecases.NewActionUseCases(actionRepo, actionRepo, actionRepo, appendEventUC)
	handler := collaboration_ui.NewCollaborationHandlerWithActions(nil, nil, appendEventUC, actionUseCases)
	return handler, actionRepo
}

// ---------------------------------------------------------------------------
// 1. Approval Tests
// ---------------------------------------------------------------------------

func TestApproval_FullLifecycle_And_Events(t *testing.T) {
	ctx := context.Background()
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, actionRepo := setupActionsHandler(threadRepo)

	// Step 1: Create Approval
	createReq := collaboration_dtos.CreateApprovalRequest{
		ResourceType:  collaboration_domain.ResourceTypeShareEntry,
		ResourceID:    "share_100",
		SourceEventID: "evt_src_01",
		Message:       "Please approve this inspection entry",
		CreatedBy:     "user_alice",
		ThreadID:      th.ID,
	}

	createResp, err := handler.CreateApproval(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	assert.Equal(t, collaboration_domain.ApprovalRequested, createResp.Approval.Status)
	assert.Equal(t, "share_100", createResp.Approval.ResourceReference.ResourceID)
	assert.Equal(t, "evt_src_01", createResp.Approval.ResourceReference.SourceEventID)

	// Verify Event: c3.approval.requested
	require.Len(t, threadRepo.events[th.ID], 1)
	evt1 := threadRepo.events[th.ID][0]
	assert.Equal(t, thread_domain.EventC3ApprovalRequested, evt1.Type)
	assert.Equal(t, createResp.Approval.ID, evt1.Payload.ActionID)
	assert.Equal(t, "evt_src_01", evt1.Payload.SourceEventID)

	// Step 2: Approve Action
	apprReq := collaboration_dtos.ApproveRequest{
		ApprovalID: createResp.Approval.ID,
		ApprovedBy: "user_bob",
		ThreadID:   th.ID,
	}

	apprResp, err := handler.ApproveAction(ctx, apprReq)
	require.NoError(t, err)
	require.NotNil(t, apprResp)
	assert.Equal(t, collaboration_domain.Approved, apprResp.Approval.Status)
	assert.Equal(t, "user_bob", apprResp.Approval.ApprovedBy)
	assert.NotNil(t, apprResp.Approval.ApprovedAt)

	// Verify Event: c3.approval.approved
	require.Len(t, threadRepo.events[th.ID], 2)
	evt2 := threadRepo.events[th.ID][1]
	assert.Equal(t, thread_domain.EventC3ApprovalApproved, evt2.Type)

	// Step 3: Idempotent re-approval
	apprResp2, err := handler.ApproveAction(ctx, apprReq)
	require.NoError(t, err)
	assert.Equal(t, collaboration_domain.Approved, apprResp2.Approval.Status)

	// Step 4: Invalid transition from Approved to Reject
	apprEntity, err := actionRepo.GetApprovalByID(ctx, createResp.Approval.ID)
	require.NoError(t, err)
	err = apprEntity.RejectApproval("user_charlie")
	assert.ErrorIs(t, err, collaboration_domain.ErrApprovalAlreadyFinalized)
}

// ---------------------------------------------------------------------------
// 2. Reject Tests
// ---------------------------------------------------------------------------

func TestReject_Creation_And_ReasonPreservation(t *testing.T) {
	ctx := context.Background()
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, _ := setupActionsHandler(threadRepo)

	// Missing message -> error
	invalidReq := collaboration_dtos.CreateRejectRequest{
		ResourceType:  collaboration_domain.ResourceTypeShareEntry,
		ResourceID:    "share_200",
		SourceEventID: "evt_src_02",
		Message:       "",
		CreatedBy:     "user_alice",
		ThreadID:      th.ID,
	}
	_, err := handler.CreateReject(ctx, invalidReq)
	assert.ErrorIs(t, err, collaboration_domain.ErrRejectionReasonRequired)

	// Valid reject
	validReq := collaboration_dtos.CreateRejectRequest{
		ResourceType:  collaboration_domain.ResourceTypeShareEntry,
		ResourceID:    "share_200",
		SourceEventID: "evt_src_02",
		Message:       "Inspection document is incomplete",
		CreatedBy:     "user_alice",
		ThreadID:      th.ID,
	}

	resp, err := handler.CreateReject(ctx, validReq)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, collaboration_domain.RejectCreated, resp.Reject.Status)
	assert.Equal(t, "Inspection document is incomplete", resp.Reject.Message)
	assert.Equal(t, "share_200", resp.Reject.ResourceReference.ResourceID)
	assert.Equal(t, "evt_src_02", resp.Reject.ResourceReference.SourceEventID)

	// Verify Event: c3.reject.created
	require.Len(t, threadRepo.events[th.ID], 1)
	evt := threadRepo.events[th.ID][0]
	assert.Equal(t, thread_domain.EventC3RejectCreated, evt.Type)
	assert.Equal(t, resp.Reject.ID, evt.Payload.ActionID)
	assert.Equal(t, "evt_src_02", evt.Payload.SourceEventID)
}

// ---------------------------------------------------------------------------
// 3. Transfer Lifecycle Tests
// ---------------------------------------------------------------------------

func TestTransfer_FullLifecycle_And_Transitions(t *testing.T) {
	ctx := context.Background()
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, _ := setupActionsHandler(threadRepo)

	// Missing TargetTrustGroupID -> error
	_, err := handler.CreateTransfer(ctx, collaboration_dtos.CreateTransferRequest{
		ResourceType:       collaboration_domain.ResourceTypeDocument,
		ResourceID:         "doc_300",
		SourceEventID:      "evt_src_03",
		TargetTrustGroupID: "",
		CreatedBy:          "user_alice",
	})
	assert.ErrorIs(t, err, collaboration_domain.ErrTargetTrustGroupRequired)

	// Step 1: Create Transfer (REQUESTED)
	trReq := collaboration_dtos.CreateTransferRequest{
		ResourceType:       collaboration_domain.ResourceTypeDocument,
		ResourceID:         "doc_300",
		SourceEventID:      "evt_src_03",
		TargetTrustGroupID: "tg_engineering_reviewers",
		Message:            "Please take responsibility for this inspection",
		CreatedBy:          "user_alice",
		ThreadID:           th.ID,
	}

	trResp, err := handler.CreateTransfer(ctx, trReq)
	require.NoError(t, err)
	require.NotNil(t, trResp)
	assert.Equal(t, collaboration_domain.TransferRequested, trResp.Transfer.Status)
	assert.Equal(t, "tg_engineering_reviewers", trResp.Transfer.TargetTrustGroupID)

	// Attempt direct Completion before Approval -> error (ErrTransferNotApproved)
	_, err = handler.CompleteTransferAction(ctx, collaboration_dtos.CompleteTransferRequest{
		TransferID:  trResp.Transfer.ID,
		CompletedBy: "user_bob",
		ThreadID:    th.ID,
	})
	assert.ErrorIs(t, err, collaboration_domain.ErrTransferNotApproved)

	// Step 2: Approve Transfer
	apprResp, err := handler.ApproveTransferAction(ctx, collaboration_dtos.ApproveTransferRequest{
		TransferID: trResp.Transfer.ID,
		ApprovedBy: "user_lead",
		ThreadID:   th.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, collaboration_domain.TransferApproved, apprResp.Transfer.Status)
	assert.Equal(t, "user_lead", apprResp.Transfer.ApprovedBy)

	// Step 3: Complete Transfer
	compResp, err := handler.CompleteTransferAction(ctx, collaboration_dtos.CompleteTransferRequest{
		TransferID:  trResp.Transfer.ID,
		CompletedBy: "user_bob",
		ThreadID:    th.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, collaboration_domain.TransferCompleted, compResp.Transfer.Status)
	assert.Equal(t, "user_bob", compResp.Transfer.CompletedBy)
	assert.NotNil(t, compResp.Transfer.CompletedAt)

	// Step 4: Idempotent Completion
	compResp2, err := handler.CompleteTransferAction(ctx, collaboration_dtos.CompleteTransferRequest{
		TransferID:  trResp.Transfer.ID,
		CompletedBy: "user_bob",
		ThreadID:    th.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, collaboration_domain.TransferCompleted, compResp2.Transfer.Status)

	// Step 5: Invalid transition on Completed Transfer
	_, err = handler.RejectTransferAction(ctx, collaboration_dtos.RejectTransferRequest{
		TransferID: trResp.Transfer.ID,
		RejectedBy: "user_hacker",
	})
	assert.ErrorIs(t, err, collaboration_domain.ErrTransferAlreadyCompleted)
}

func TestTransfer_RejectionWorkflow(t *testing.T) {
	ctx := context.Background()
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, _ := setupActionsHandler(threadRepo)

	// Step 1: Create Transfer
	trResp, err := handler.CreateTransfer(ctx, collaboration_dtos.CreateTransferRequest{
		ResourceType:       collaboration_domain.ResourceTypeDocument,
		ResourceID:         "doc_400",
		SourceEventID:      "evt_src_04",
		TargetTrustGroupID: "tg_legal",
		Message:            "Transferring to legal",
		CreatedBy:          "user_alice",
		ThreadID:           th.ID,
	})
	require.NoError(t, err)

	// Step 2: Reject Transfer
	rejResp, err := handler.RejectTransferAction(ctx, collaboration_dtos.RejectTransferRequest{
		TransferID: trResp.Transfer.ID,
		RejectedBy: "user_legal_head",
		ThreadID:   th.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, collaboration_domain.TransferRejected, rejResp.Transfer.Status)

	// Attempt to Complete a Rejected Transfer -> error
	_, err = handler.CompleteTransferAction(ctx, collaboration_dtos.CompleteTransferRequest{
		TransferID:  trResp.Transfer.ID,
		CompletedBy: "user_alice",
	})
	assert.ErrorIs(t, err, collaboration_domain.ErrInvalidTransferTransition)
}

// ---------------------------------------------------------------------------
// 4. Security & Payload Sanity Tests
// ---------------------------------------------------------------------------

func TestActionEvents_SecurityBoundaryVerification(t *testing.T) {
	ctx := context.Background()
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, _ := setupActionsHandler(threadRepo)

	// Execute all three actions
	apprResp, _ := handler.CreateApproval(ctx, collaboration_dtos.CreateApprovalRequest{
		ResourceType:  "share_entry",
		ResourceID:    "se_999",
		SourceEventID: "evt_999",
		CreatedBy:     "user_1",
		ThreadID:      th.ID,
	})

	_, _ = handler.CreateReject(ctx, collaboration_dtos.CreateRejectRequest{
		ResourceType:  "share_entry",
		ResourceID:    "se_999",
		SourceEventID: "evt_999",
		Message:       "Rejected due to incomplete field",
		CreatedBy:     "user_1",
		ThreadID:      th.ID,
	})

	_, _ = handler.CreateTransfer(ctx, collaboration_dtos.CreateTransferRequest{
		ResourceType:       "share_entry",
		ResourceID:         "se_999",
		SourceEventID:      "evt_999",
		TargetTrustGroupID: "tg_target",
		CreatedBy:          "user_1",
		ThreadID:           th.ID,
	})

	_, _ = handler.ApproveAction(ctx, collaboration_dtos.ApproveRequest{
		ApprovalID: apprResp.Approval.ID,
		ApprovedBy: "user_2",
		ThreadID:   th.ID,
	})

	events := threadRepo.events[th.ID]
	require.NotEmpty(t, events)

	forbiddenSecrets := []string{
		"wrapped_dek",
		"wrappedDEK",
		"wrapped_kek",
		"wrappedKEK",
		"key_envelopes",
		"PrivateKey",
		"SecretKey",
		"private_key",
		"secret_key",
		"encrypted_payload",
		"decrypted_content",
	}

	for _, evt := range events {
		marshaledBytes, err := json.Marshal(evt)
		require.NoError(t, err)
		marshaledStr := string(marshaledBytes)

		for _, forbidden := range forbiddenSecrets {
			if strings.Contains(marshaledStr, forbidden) {
				t.Errorf("SECURITY VIOLATION: Action Event %s contains forbidden secret substring: %q", evt.Type, forbidden)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// 5. Provenance Survival Test
// ---------------------------------------------------------------------------

func TestAction_ProvenanceSurvival(t *testing.T) {
	ctx := context.Background()
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, _ := setupActionsHandler(threadRepo)

	trResp, err := handler.CreateTransfer(ctx, collaboration_dtos.CreateTransferRequest{
		ResourceType:       collaboration_domain.ResourceTypeDocument,
		ResourceID:         "doc_provenance_777",
		SourceEventID:      "evt_provenance_888",
		TargetTrustGroupID: "tg_audit",
		Message:            "Provenance test",
		CreatedBy:          "user_alice",
		ThreadID:           th.ID,
	})
	require.NoError(t, err)

	trAppr, err := handler.ApproveTransferAction(ctx, collaboration_dtos.ApproveTransferRequest{
		TransferID: trResp.Transfer.ID,
		ApprovedBy: "user_bob",
		ThreadID:   th.ID,
	})
	require.NoError(t, err)

	trComp, err := handler.CompleteTransferAction(ctx, collaboration_dtos.CompleteTransferRequest{
		TransferID:  trAppr.Transfer.ID,
		CompletedBy: "user_charlie",
		ThreadID:    th.ID,
	})
	require.NoError(t, err)

	// Assert provenance survival through all state transitions
	assert.Equal(t, "doc_provenance_777", trComp.Transfer.ResourceReference.ResourceID)
	assert.Equal(t, collaboration_domain.ResourceTypeDocument, trComp.Transfer.ResourceReference.ResourceType)
	assert.Equal(t, "evt_provenance_888", trComp.Transfer.ResourceReference.SourceEventID)
}
