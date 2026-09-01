if (typeof globalThis.import === "undefined") {
  (globalThis as any).import = { meta: { env: { VITE_API_URL: "http://localhost:4001" } } };
}

if (typeof (globalThis as any).window === "undefined") {
  (globalThis as any).window = { go: { main: { App: {} } } };
}

import { useC3ConfigurationStore, SlotDraft } from "./useC3ConfigurationStore";
import { useC3ChannelStore } from "@/components/C3/infrastructure/store/useC3ChannelStore";
import { useAuthStore } from "@/store/useAuthStore";
import { ChannelResponse, UpdateChannelPayload } from "@/services/api";

async function runRegressionTest() {
  console.log("=================================================");
  console.log("RUNNING REGRESSION TEST: 3 slots → add 4th → save → reload → 4 slots → discard → 4 slots");
  console.log("=================================================");

  // Set auth token
  useAuthStore.setState({ jwtToken: "test-auth-token" });

  const dummyChannel: ChannelResponse = {
    id: "ch_reg_test_001",
    workspace_id: "ws_reg_001",
    title: "Regression Channel",
    slots: [
      { id: "s1", name: "slot_1", role: "drafter", vault_id: "v1", gated: false, order: 1 },
      { id: "s2", name: "slot_2", role: "approver", vault_id: "v2", gated: true, order: 2 },
      { id: "s3", name: "slot_3", role: "signatory", vault_id: "v3", gated: true, order: 3 },
    ],
    assignments: [],
    properties: [],
  };

  // 1. Setup existing 3 slots
  useC3ChannelStore.setState({
    channels: [dummyChannel],
    activeChannel: dummyChannel,
    activeChannelId: dummyChannel.id,
  });

  useC3ConfigurationStore.setState({
    originalChannel: dummyChannel,
    draftSlots: dummyChannel.slots.map((s) => ({ ...s })),
    isDirty: false,
  });

  console.log("1. Existing 3 slots setup complete. draftSlots count:", useC3ConfigurationStore.getState().draftSlots.length);
  if (useC3ConfigurationStore.getState().draftSlots.length !== 3) {
    throw new Error("Failed at step 1: Expected 3 initial slots");
  }

  // 2. Add 4th slot
  const s4: SlotDraft = {
    id: "s4",
    name: "slot_4_qa",
    role: "reviewer",
    vault_id: "v4",
    gated: true,
    order: 4,
  };

  const draft3 = useC3ConfigurationStore.getState().draftSlots;
  useC3ConfigurationStore.getState().setDraftSlots([...draft3, s4]);

  console.log("2. Added 4th slot via setDraftSlots. draftSlots count:", useC3ConfigurationStore.getState().draftSlots.length);
  if (useC3ConfigurationStore.getState().draftSlots.length !== 4) {
    throw new Error("Failed at step 2: Expected 4 slots in draft");
  }
  if (!useC3ConfigurationStore.getState().isDirty) {
    throw new Error("Failed at step 2: Expected isDirty to be true after mutation");
  }

  // 3. Save
  let capturedSlotsPayload: any[] = [];
  useC3ChannelStore.setState({
    updateChannelDetails: async (payload: UpdateChannelPayload): Promise<ChannelResponse> => {
      capturedSlotsPayload = payload.slots || [];
      const updatedChannel: ChannelResponse = {
        ...dummyChannel,
        slots: payload.slots?.map((s) => ({ ...s })) || [],
      };
      useC3ChannelStore.setState({ activeChannel: updatedChannel });
      return updatedChannel;
    },
  });

  const saveOk = await useC3ConfigurationStore.getState().saveChannelConfig();
  console.log("3. Save executed. Succeeded:", saveOk, "| Outbound payload slots count:", capturedSlotsPayload.length);
  if (!saveOk || capturedSlotsPayload.length !== 4) {
    throw new Error(`Failed at step 3: Outbound payload slots count is ${capturedSlotsPayload.length}, expected 4`);
  }
  if (useC3ConfigurationStore.getState().draftSlots.length !== 4) {
    throw new Error("Failed at step 3: Post-save draftSlots count must be 4");
  }
  if (useC3ConfigurationStore.getState().isDirty) {
    throw new Error("Failed at step 3: Post-save isDirty must be false");
  }

  // Mock ListTrustGroups for loadConfiguration
  (globalThis as any).window.go.main.App.ListTrustGroups = async () => [];

  // 4. Reload
  useC3ChannelStore.setState({
    refreshChannel: async () => {
      return useC3ChannelStore.getState().activeChannel!;
    },
  });

  await useC3ConfigurationStore.getState().loadConfiguration();
  console.log("4. Reload executed (loadConfiguration). draftSlots count:", useC3ConfigurationStore.getState().draftSlots.length);
  if (useC3ConfigurationStore.getState().draftSlots.length !== 4) {
    throw new Error(`Failed at step 4: Post-reload draftSlots count is ${useC3ConfigurationStore.getState().draftSlots.length}, expected 4`);
  }

  // 5. Discard (Test discardChanges on 4-slot channel)
  const current4 = useC3ConfigurationStore.getState().draftSlots;
  useC3ConfigurationStore.getState().setDraftSlots([
    ...current4,
    { id: "s5_temp", name: "slot_5_temp", role: "observer", vault_id: "v5", gated: false, order: 5 },
  ]);
  console.log("Draft mutated with temporary 5th slot. draftSlots count:", useC3ConfigurationStore.getState().draftSlots.length);

  useC3ConfigurationStore.getState().discardChanges();
  const afterDiscardCount = useC3ConfigurationStore.getState().draftSlots.length;
  console.log("5. Discard executed (discardChanges). draftSlots count:", afterDiscardCount);
  if (afterDiscardCount !== 4) {
    throw new Error(`Failed at step 5: Post-discard draftSlots count is ${afterDiscardCount}, expected 4`);
  }

  // 6. Test Trust Group CRUD & Selection
  console.log("\n=================================================");
  console.log("RUNNING TRUST GROUP STORE TESTS: Create → Select → Edit/Update → Delete");
  console.log("=================================================");

  // Mock Wails bindings on window.go.main.App
  (globalThis as any).window.go.main.App.CreateTrustGroup = async (jwtToken: string, workspaceId: string, name: string) => ({
    id: "tg_authoritative_101",
    name,
    channel_id: workspaceId,
    kek_version: 1,
    member_cids: ["vault_legal"],
    created_at: new Date().toISOString(),
  });

  (globalThis as any).window.go.main.App.UpdateTrustGroup = async (jwtToken: string, id: string, name: string) => ({
    id,
    name,
    kek_version: 2,
    created_at: new Date().toISOString(),
  });

  (globalThis as any).window.go.main.App.DeleteTrustGroup = async (jwtToken: string, id: string) => {};

  // Test 6.1: Create Trust Group
  const createdGroup = await useC3ConfigurationStore.getState().createTrustGroup({
    name: "Legal Sovereign Group",
    description: "Legal team trust group",
  });

  console.log("6.1. Created Trust Group:", createdGroup.id, createdGroup.name);
  if (createdGroup.id !== "tg_authoritative_101") {
    throw new Error(`Failed at step 6.1: Expected ID 'tg_authoritative_101', got '${createdGroup.id}'`);
  }

  const storeGroups = useC3ConfigurationStore.getState().trustGroups;
  if (storeGroups.length !== 1 || storeGroups[0].id !== "tg_authoritative_101") {
    throw new Error(`Failed at step 6.1: Trust group not immediately appended to store. count=${storeGroups.length}`);
  }

  if (useC3ConfigurationStore.getState().selectedTrustGroupId !== "tg_authoritative_101") {
    throw new Error("Failed at step 6.1: selectedTrustGroupId not automatically set to created group ID");
  }

  // Test 6.2: Selection
  useC3ConfigurationStore.getState().selectTrustGroup("tg_authoritative_101");
  const selectedId = useC3ConfigurationStore.getState().selectedTrustGroupId;
  const resolvedGroup = useC3ConfigurationStore.getState().trustGroups.find((tg) => tg.id === selectedId);

  console.log("6.2. Resolved selected Trust Group:", resolvedGroup?.id, resolvedGroup?.name);
  if (!resolvedGroup || resolvedGroup.id !== "tg_authoritative_101") {
    throw new Error("Failed at step 6.2: Selected Trust Group did not resolve correctly");
  }

  // Test 6.3: Update
  const updatedGroup = await useC3ConfigurationStore.getState().updateTrustGroup({
    id: "tg_authoritative_101",
    name: "Legal Sovereign Group V2",
  });

  console.log("6.3. Updated Trust Group:", updatedGroup.id, updatedGroup.name, "KEK:", updatedGroup.kekVersion);
  if (updatedGroup.name !== "Legal Sovereign Group V2" || updatedGroup.kekVersion !== 2) {
    throw new Error(`Failed at step 6.3: Expected updated name 'Legal Sovereign Group V2', got '${updatedGroup.name}'`);
  }

  const updatedInList = useC3ConfigurationStore.getState().trustGroups.find((tg) => tg.id === "tg_authoritative_101");
  if (updatedInList?.name !== "Legal Sovereign Group V2") {
    throw new Error("Failed at step 6.3: Updated values not reflected in store.trustGroups list");
  }

  // Test 6.4: Delete
  await useC3ConfigurationStore.getState().deleteTrustGroup("tg_authoritative_101");
  const postDeleteGroups = useC3ConfigurationStore.getState().trustGroups;
  const postDeleteSelectedId = useC3ConfigurationStore.getState().selectedTrustGroupId;

  console.log("6.4. Post-delete groups count:", postDeleteGroups.length, "| selectedTrustGroupId:", postDeleteSelectedId);
  if (postDeleteGroups.length !== 0) {
    throw new Error(`Failed at step 6.4: Expected 0 groups in store after deletion, got ${postDeleteGroups.length}`);
  }
  if (postDeleteSelectedId !== null) {
    throw new Error("Failed at step 6.4: selectedTrustGroupId was not cleared after deleting active group");
  }

  console.log("\n=================================================");
  console.log("SUCCESS: ALL REGRESSION AND TRUST GROUP TESTS PASSED CLEANLY! ✓");
  console.log("=================================================");
}

runRegressionTest().catch((err) => {
  console.error("REGRESSION TEST FAILED:", err);
  process.exit(1);
});
