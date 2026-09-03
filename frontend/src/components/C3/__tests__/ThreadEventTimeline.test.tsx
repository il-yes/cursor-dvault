import { describe, it, expect } from "vitest";
import { extractThreadEventResourceRef } from "../ThreadEventTimeline";

describe("ThreadEventTimeline Payload Extraction & Invariants", () => {
  it("extracts resource reference from Go capital Payload JSON response", () => {
    const event = {
      id: "evt_100",
      thread_id: "th_channel_1",
      type: "entry.shared",
      cursor: 1,
      created_at: "2026-09-01T12:00:00Z",
      Payload: {
        ref_type: "share_entry",
        share_entry_id: "share-123",
        trust_group_id: "trust-456",
      },
    };

    const extracted = extractThreadEventResourceRef(event);
    expect(extracted.refType).toBe("share_entry");
    expect(extracted.shareEntryId).toBe("share-123");
    expect(extracted.trustGroupId).toBe("trust-456");
  });

  it("ensures ChannelID is NOT substituted for TrustGroupID", () => {
    const channelId = "ch_finance_99";
    const trustGroupId = "tg_board_88";

    const event = {
      id: "evt_101",
      thread_id: channelId,
      type: "entry.shared",
      cursor: 1,
      created_at: "2026-09-01T12:00:00Z",
      payload: {
        ref_type: "share_entry",
        share_entry_id: "share-999",
        trust_group_id: trustGroupId,
      },
    };

    const extracted = extractThreadEventResourceRef(event);
    expect(extracted.trustGroupId).toBe(trustGroupId);
    expect(extracted.trustGroupId).not.toBe(channelId);
  });
});
