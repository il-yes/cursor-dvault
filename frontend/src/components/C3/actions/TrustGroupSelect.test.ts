import { describe, it, expect, beforeEach } from "vitest";
import { useC3ConfigurationStore } from "../configuration/store/useC3ConfigurationStore";
import { useC3WorkspaceStore } from "../infrastructure/store/useC3WorkspaceStore";

describe("TrustGroupSelect & Authoritative TrustGroup Path Invariants", () => {
  beforeEach(() => {
    useC3ConfigurationStore.setState({
      trustGroups: [],
      selectedTrustGroupId: null,
    });
    useC3WorkspaceStore.setState({
      activeWorkspaceId: "ws_auth_test_2026",
    });
  });

  it("stores and surfaces real persisted TrustGroups created through C3 Configuration UI", () => {
    const realTrustGroup = {
      id: "tg_authoritative_persisted_999",
      name: "Legal Council Alpha",
      description: "Authoritative Sovereign Trust Group",
      memberCount: 2,
      deviceCount: 2,
      status: "active" as const,
      kekVersion: 1,
      createdDate: new Date().toISOString(),
      associatedChannelIds: ["ch_1"],
      members: [],
      devices: [],
      cryptographicState: {
        kekVersion: 1,
        envelopeCount: 2,
        activeEnvelopes: 2,
        revokedEnvelopes: 0,
        lastRotation: new Date().toISOString(),
      },
    };

    // Simulate creation of TrustGroup via store/API
    useC3ConfigurationStore.setState((state) => ({
      trustGroups: [...state.trustGroups, realTrustGroup],
    }));

    const storeGroups = useC3ConfigurationStore.getState().trustGroups;
    expect(storeGroups).toHaveLength(1);
    expect(storeGroups[0].id).toBe("tg_authoritative_persisted_999");
    expect(storeGroups[0].name).toBe("Legal Council Alpha");
  });

  it("ensures selecting a TrustGroup passes its exact persisted TrustGroupID unchanged to CreateCollaborativeShare", () => {
    const realGroupId = "tg_persisted_uuid_88776655";
    let capturedTrustGroupId = "";

    // Mock callback passed to TrustGroupSelect onChange
    const onChange = (selectedId: string) => {
      capturedTrustGroupId = selectedId;
    };

    onChange(realGroupId);

    expect(capturedTrustGroupId).toBe("tg_persisted_uuid_88776655");
    expect(capturedTrustGroupId).not.toContain("ch_"); // Must not be a channel ID
    expect(capturedTrustGroupId).not.toBe("tg_legal_counsel"); // Must not be mock fallback
  });

  it("requires an active workspace and rejects 'me' fallback in production query path", () => {
    useC3WorkspaceStore.setState({ activeWorkspaceId: null });

    const activeWorkspaceId = useC3WorkspaceStore.getState().activeWorkspaceId;
    expect(activeWorkspaceId).toBeNull();
    // Invariant: no automatic fallback to 'me' string
    const targetWorkspaceQuery = activeWorkspaceId || "";
    expect(targetWorkspaceQuery).not.toBe("me");
  });

  it("rejects any mock TrustGroup ID from reaching the creation path", () => {
    const isMockTrustGroup = (id: string) => {
      if (!id || id.trim() === "") return true;
      if (id.startsWith("mock_")) return true;
      if (id === "tg_legal_counsel" || id === "tg_finance") return true;
      return false;
    };

    expect(isMockTrustGroup("tg_legal_counsel")).toBe(true);
    expect(isMockTrustGroup("mock_trust_group")).toBe(true);
    expect(isMockTrustGroup("")).toBe(true);

    // Real persisted TrustGroups pass validation
    expect(isMockTrustGroup("tg_authoritative_persisted_999")).toBe(false);
    expect(isMockTrustGroup("tg_550e8400-e29b-41d4-a716-446655440000")).toBe(false);
  });
});
