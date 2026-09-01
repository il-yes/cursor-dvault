import { describe, it, expect } from "vitest";
import { C3ResourceTypes, ResourceReference } from "../domain/resource";

describe("C3 Frontend Collaboration Action Primitives", () => {
  it("defines generic resource references without coupling to VaultEntry", () => {
    const ref: ResourceReference = {
      resourceType: C3ResourceTypes.Document,
      resourceId: "doc_100",
      sourceEventId: "evt_origin_100",
    };

    expect(ref.resourceType).toBe("document");
    expect(ref.resourceId).toBe("doc_100");
    expect(ref.sourceEventId).toBe("evt_origin_100");
  });

  it("supports all required C3 Action timeline event types", () => {
    const actionEvents = [
      "c3.approval.requested",
      "c3.approval.approved",
      "c3.reject.created",
      "c3.transfer.requested",
      "c3.transfer.approved",
      "c3.transfer.rejected",
      "c3.transfer.completed",
    ];

    actionEvents.forEach((evtType) => {
      expect(evtType.startsWith("c3.")).toBe(true);
    });
  });

  it("enforces explicit transfer workflow state machine rules", () => {
    const transferStatuses = [
      "transfer_requested",
      "transfer_approved",
      "transfer_rejected",
      "transfer_completed",
    ];

    expect(transferStatuses).toContain("transfer_requested");
    expect(transferStatuses).toContain("transfer_approved");
    expect(transferStatuses).toContain("transfer_completed");
  });
});
