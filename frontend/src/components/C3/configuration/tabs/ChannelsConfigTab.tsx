import React, { useState, useEffect } from "react";
import { useC3ConfigurationStore, SlotDraft } from "../store/useC3ConfigurationStore";

export const ChannelsConfigTab: React.FC = () => {
  useEffect(() => {
    console.log(`[BOUNDARIES][READ] ChannelsConfigTab mounted draftSlotsCount=${draftSlots.length}`);
  }, []);

  const {
    draftSlots,
    setDraftSlots,
    isDirty,
    isSaving,
    saveFeedback,
    saveChannelConfig,
    discardChanges,
  } = useC3ConfigurationStore();

  const [isAddOpen, setIsAddOpen] = useState(false);
  const [editingSlotId, setEditingSlotId] = useState<string | null>(null);
  const [newSlotName, setNewSlotName] = useState("");
  const [assignedVault, setAssignedVault] = useState("vault_finance");
  const [gateCondition, setGateCondition] = useState("executive_signature");
  const [isRequired, setIsRequired] = useState(true);

  const handleMoveSlotUp = (index: number) => {
    if (index <= 0) return;
    const updated = [...draftSlots];
    const temp = updated[index - 1];
    updated[index - 1] = updated[index];
    updated[index] = temp;
    const reordered = updated.map((s, idx) => ({ ...s, order: idx + 1 }));
    setDraftSlots(reordered);
  };

  const handleMoveSlotDown = (index: number) => {
    if (index >= draftSlots.length - 1) return;
    const updated = [...draftSlots];
    const temp = updated[index + 1];
    updated[index + 1] = updated[index];
    updated[index] = temp;
    const reordered = updated.map((s, idx) => ({ ...s, order: idx + 1 }));
    setDraftSlots(reordered);
  };

  const handleEditSlot = (slot: SlotDraft) => {
    setEditingSlotId(slot.id);
    setNewSlotName(slot.name);
    setAssignedVault(slot.vault_id || "vault_finance");
    setGateCondition(slot.gated ? "executive_signature" : "No gate (first slot)");
    setIsAddOpen(true);
  };

  const handleAddSlotSubmit = () => {
    if (!newSlotName.trim()) return;

    const existingSlot = editingSlotId ? draftSlots.find((s) => s.id === editingSlotId) : null;

    const slotToSave: SlotDraft = {
      id: editingSlotId || `slot_${Date.now()}_${Math.random().toString(36).substring(2, 6)}`,
      name: newSlotName.trim(),
      role: existingSlot?.role || "participant",
      vault_id: assignedVault || "vault_finance",
      gated: gateCondition !== "No gate (first slot)",
      order: existingSlot?.order || draftSlots.length + 1,
    };

    console.log(`[BOUNDARIES][① ChannelsConfigTab.handleAddSlotSubmit] slotToSave=`, slotToSave, `editingSlotId=${editingSlotId}`);

    if (editingSlotId) {
      setDraftSlots(
        draftSlots.map((s) => (s.id === editingSlotId ? slotToSave : s))
      );
    } else {
      setDraftSlots([...draftSlots, slotToSave]);
    }

    setEditingSlotId(null);
    setNewSlotName("");
    setIsAddOpen(false);
  };


  const handleDeleteSlot = (id: string) => {
    const remaining = draftSlots.filter((s) => s.id !== id);
    const reordered = remaining.map((s, idx) => ({ ...s, order: idx + 1 }));
    setDraftSlots(reordered);
  };

  const getVaultBadgeStyle = (vault: string) => {
    switch (vault) {
      case "vault_legal":
        return { bgClass: "vc-pur", color: "#7C3AED" };
      case "vault_finance":
        return { bgClass: "vc-blu", color: "#2563EB" };
      case "vault_direction":
        return { bgClass: "vc-dar", color: "#444" };
      default:
        return { bgClass: "vc-grn", color: "#059669" };
    }
  };

  return (
    <div className="tab-body">
      <div className="tab-desc">
        Ordered commit steps for this channel. Each slot is committed by its assigned vault in sequence. Drag to reorder. Gate conditions are validated at runtime by AFP and are not retroactively applied to running threads.
      </div>

      <div className="slot-area">
        {/* Column headers */}
        <div className="sg col-hdrs">
          <div className="ch-lbl"></div>
          <div className="ch-lbl">#</div>
          <div className="ch-lbl" style={{ paddingLeft: "14px" }}>Slot name</div>
          <div className="ch-lbl">Assigned vault</div>
          <div className="ch-lbl">Gate condition</div>
          <div className="ch-lbl" style={{ textAlign: "center" }}>Req.</div>
          <div className="ch-lbl"></div>
        </div>

        {/* Dynamic Slot Rows */}
        {draftSlots.length === 0 && !isAddOpen && (
          <div className="no-slots-msg" style={{ padding: "32px 16px", textAlign: "center", color: "#888", fontSize: "14px" }}>
            No slots configured
          </div>
        )}
        {draftSlots.map((slot, index) => {

          const badge = getVaultBadgeStyle(slot.vault_id);
          return (
            <div className="sg slot-row" key={slot.id || index}>
              <div className="sr-drag">⠿</div>
              <div className="sr-num">{index + 1}</div>
              <div className="sr-name">
                <span className="slot-tag">{slot.name}</span>
              </div>
              <div className="sr-vault">
                <span className={`vault-chip ${badge.bgClass}`}>
                  <span className="vcdot" style={{ background: badge.color }}></span>
                  {slot.vault_id}
                </span>
              </div>
              <div className="sr-gate">
                {index === 0 ? (
                  <span className="gate-none">No gate — first slot</span>
                ) : (
                  <span className="gate-chip">
                    <span className="gate-arr">→</span>
                    {draftSlots[index - 1]?.name}
                  </span>
                )}
              </div>
              <div className="sr-req">
                <span className="req-dot"></span>
              </div>
              <div className="sr-acts" style={{ display: "flex", gap: "4px", alignItems: "center" }}>
                <div
                  className="icon-btn"
                  onClick={() => handleMoveSlotUp(index)}
                  title="Move Up"
                  style={{ cursor: index === 0 ? "default" : "pointer", opacity: index === 0 ? 0.3 : 1 }}
                >
                  ↑
                </div>
                <div
                  className="icon-btn"
                  onClick={() => handleMoveSlotDown(index)}
                  title="Move Down"
                  style={{ cursor: index === draftSlots.length - 1 ? "default" : "pointer", opacity: index === draftSlots.length - 1 ? 0.3 : 1 }}
                >
                  ↓
                </div>
                <div className="icon-btn" onClick={() => handleEditSlot(slot)} title="Edit Slot" style={{ cursor: "pointer" }}>
                  ✎
                </div>
                <div className="icon-btn del" onClick={() => handleDeleteSlot(slot.id)} title="Delete Slot" style={{ cursor: "pointer" }}>
                  ✕
                </div>
              </div>
            </div>
          );
        })}

        {/* Add / Edit Slot Row */}
        {isAddOpen ? (
          <div className="add-row">
            <div className="add-inner">
              <div></div>
              <div className="add-num">{editingSlotId ? draftSlots.findIndex(s => s.id === editingSlotId) + 1 : draftSlots.length + 1}</div>
              <div style={{ padding: "0 6px" }}>
                <input
                  className="slot-input"
                  type="text"
                  placeholder="new_slot_name"
                  value={newSlotName}
                  onChange={(e) => setNewSlotName(e.target.value)}
                  autoFocus
                />
              </div>
              <div style={{ padding: "0 6px" }}>
                <select
                  className="slot-select"
                  value={assignedVault}
                  onChange={(e) => setAssignedVault(e.target.value)}
                >
                  <option value="vault_finance">vault_finance</option>
                  <option value="vault_legal">vault_legal</option>
                  <option value="vault_direction">vault_direction</option>
                  <option value="vault_treasury">vault_treasury</option>
                  <option value="vault_hr">vault_hr</option>
                  <option value="vault_ops">vault_ops</option>
                </select>
              </div>
              <div style={{ padding: "0 6px" }}>
                <select
                  className="slot-select"
                  value={gateCondition}
                  onChange={(e) => setGateCondition(e.target.value)}
                >
                  <option value="No gate (first slot)">No gate (first slot)</option>
                  <option value="executive_signature">executive_signature</option>
                  <option value="financial_clearance">financial_clearance</option>
                  <option value="contract_draft">contract_draft</option>
                </select>
              </div>
              <div className="add-req">
                <input
                  type="checkbox"
                  checked={isRequired}
                  onChange={(e) => setIsRequired(e.target.checked)}
                  style={{ accentColor: "#C8922A", width: "13px", height: "13px" }}
                />
              </div>
              <div className="add-actions">
                <button className="btn-confirm" onClick={handleAddSlotSubmit}>
                  {editingSlotId ? "✓ Update" : "+ Add"}
                </button>
                <button
                  className="btn-cancel-add"
                  onClick={() => {
                    setIsAddOpen(false);
                    setEditingSlotId(null);
                    setNewSlotName("");
                  }}
                >
                  ✕
                </button>
              </div>
            </div>
          </div>
        ) : (
          <div className="add-trigger" onClick={() => setIsAddOpen(true)}>
            <span className="add-trigger-plus">+</span> Add another slot
          </div>
        )}
      </div>

      {(isDirty || saveFeedback) && (
        <div className="dirty-bar">
          <span className="db-icon">⚠</span>
          <div className="db-msg">
            {saveFeedback ? (
              <span>{saveFeedback}</span>
            ) : (
              <span>
                <strong>Unsaved changes</strong> — slot modifications take effect on the next thread instantiation.
              </span>
            )}
          </div>
          <button className="btn-discard" onClick={discardChanges} disabled={isSaving}>
            Discard
          </button>
          <button className="btn-save" onClick={() => {
            const corrId = `save_${Date.now()}_${Math.random().toString(36).substring(2,6)}`;
            console.log(`[SLOTS][SAVE][id=${corrId}] STEP=01 EVENT=BUTTON_CLICK channelId=active`);
            console.log(`[SLOTS][SAVE][id=${corrId}] STEP=02 EVENT=STORE_INPUT slotsCount=${draftSlots.length} slots=`, JSON.stringify(draftSlots.map(s => ({ id: s.id, name: s.name, role: s.role, vault_id: s.vault_id, gated: s.gated, order: s.order }))));
            saveChannelConfig(corrId).then(() => {
              console.log(`[SLOTS][UI][id=${corrId}] STEP=21 EVENT=RENDER_STATE slots=`, JSON.stringify(useC3ConfigurationStore.getState().draftSlots));
            });
          }} disabled={isSaving}>
            {isSaving ? "Saving..." : "Save Changes"}
          </button>
        </div>
      )}
    </div>
  );
};



