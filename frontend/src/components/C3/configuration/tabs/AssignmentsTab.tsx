import React, { useState } from "react";
import { useC3ConfigurationStore, AssignmentDraft } from "../store/useC3ConfigurationStore";

export const AssignmentsTab: React.FC = () => {
  const {
    draftAssignments,
    setDraftAssignments,
    isDirty,
    isSaving,
    saveFeedback,
    saveChannelConfig,
    discardChanges,
  } = useC3ConfigurationStore();

  const [isPanelOpen, setIsPanelOpen] = useState(false);
  const [editingAssignmentKey, setEditingAssignmentKey] = useState<string | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<number>(3);
  const [selectedVault, setSelectedVault] = useState<string>("vault_finance");
  const [isVaultDropdownOpen, setIsVaultDropdownOpen] = useState(false);
  const [selectedRole, setSelectedRole] = useState<"primary" | "backup" | "observer">("primary");
  const [note, setNote] = useState("");

  const getAssignmentKey = (a: AssignmentDraft) =>
    `${a.slot_id}::${a.vault_address || a.owner_id}`;

  const handleOpenAddPanel = (slotNum?: number) => {
    setEditingAssignmentKey(null);
    if (slotNum) setSelectedSlot(slotNum);
    setNote("");
    setIsPanelOpen(true);
  };

  const handleEditAssignment = (assign: AssignmentDraft) => {
    setEditingAssignmentKey(getAssignmentKey(assign));
    setSelectedSlot(Number(assign.slot_id) || 1);
    setSelectedVault(assign.vault_address || assign.owner_id || "vault_finance");
    setSelectedRole((assign.role as any) || "primary");
    setNote(assign.by || "");
    setIsPanelOpen(true);
  };

  const handleSaveAssignment = () => {
    const updatedAssignment: AssignmentDraft = {
      slot_id: String(selectedSlot),
      owner_id: selectedVault,
      public_key: `pk_${selectedVault}_01`,
      vault_address: selectedVault,
      role: selectedRole,
      since: "just now",
      by: note.trim() || "admin@acme.io",
    };

    if (editingAssignmentKey) {
      setDraftAssignments(
        draftAssignments.map((a) =>
          getAssignmentKey(a) === editingAssignmentKey ? updatedAssignment : a
        )
      );
    } else {
      setDraftAssignments([...draftAssignments, updatedAssignment]);
    }

    setEditingAssignmentKey(null);
    setIsPanelOpen(false);
    setNote("");
  };

  const handleDeleteAssignment = (assignToDelete: AssignmentDraft) => {
    const targetKey = getAssignmentKey(assignToDelete);
    setDraftAssignments(
      draftAssignments.filter((a) => getAssignmentKey(a) !== targetKey)
    );
  };

  const getVaultBadgeStyle = (vault: string) => {
    switch (vault) {
      case "vault_legal":
        return { bgClass: "vc-pur", color: "#7C3AED" };
      case "vault_finance":
        return { bgClass: "vc-blu", color: "#2563EB" };
      case "vault_direction":
        return { bgClass: "vc-dar", color: "#444" };
      case "vault_treasury":
        return { bgClass: "", color: "#059669", style: { background: "#F0FDF4", border: "1px solid #86EFAC", color: "#059669" } };
      default:
        return { bgClass: "vc-ext", color: "#C8922A" };
    }
  };

  return (
    <div className="tab-body" style={{ position: "relative" }}>
      <div className="tab-desc">
        Default vault-to-slot bindings for this channel. These are the template-level assignments — they can be overridden per-thread at instantiation time.
      </div>

      <div className="assign-area">
        {/* Slot Group 1 */}
        <div className="assign-row-group">
          <div className="slot-group-header">
            <div className="sgh-num">1</div>
            <div className="sgh-slot">contract_draft</div>
            <span className="sgh-gate" style={{ color: "#bbb", background: "#f8f8f8", borderColor: "#ebebeb" }}>
              No gate
            </span>
            <div className="sgh-spacer"></div>
            <span className="sgh-add" style={{ cursor: "pointer" }} onClick={() => handleOpenAddPanel(1)}>
              + Add assignee
            </span>
          </div>
          {draftAssignments
            .filter((a) => a.slot_id === "1")
            .map((assign, idx) => {
              const badge = getVaultBadgeStyle(assign.vault_address || assign.owner_id);
              return (
                <div className="assign-row" key={idx}>
                  <span className={`vault-chip ${badge.bgClass}`} style={badge.style}>
                    <span className="vcdot" style={{ background: badge.color }}></span>
                    {assign.vault_address || assign.owner_id}
                  </span>
                  <div style={{ flex: 1 }}></div>
                  <span className={`role-badge ${assign.role === "backup" ? "rb-backup" : "rb-primary"}`}>
                    {(assign.role || "primary").toUpperCase()}
                  </span>
                  <div className="ar-since">{assign.since || "since Jun 10"}</div>
                  <div className="ar-by">
                    <code>{assign.by || "ch_a8f2e4b1 (template)"}</code>
                  </div>
                  <div className="ar-acts">
                    <div className="icon-btn" onClick={() => handleEditAssignment(assign)} style={{ cursor: "pointer" }}>✎</div>
                    <div className="icon-btn" onClick={() => handleDeleteAssignment(assign)} style={{ cursor: "pointer" }}>✕</div>
                  </div>
                </div>
              );
            })}
        </div>

        {/* Slot Group 2 */}
        <div className="assign-row-group">
          <div className="slot-group-header">
            <div className="sgh-num">2</div>
            <div className="sgh-slot">financial_clearance</div>
            <span className="sgh-gate">→ contract_draft</span>
            <div className="sgh-spacer"></div>
            <span className="sgh-add" style={{ cursor: "pointer" }} onClick={() => handleOpenAddPanel(2)}>
              + Add assignee
            </span>
          </div>
          {draftAssignments
            .filter((a) => a.slot_id === "2")
            .map((assign, idx) => {
              const badge = getVaultBadgeStyle(assign.vault_address || assign.owner_id);
              return (
                <div className="assign-row" key={idx}>
                  <span className={`vault-chip ${badge.bgClass}`} style={badge.style}>
                    <span className="vcdot" style={{ background: badge.color }}></span>
                    {assign.vault_address || assign.owner_id}
                  </span>
                  <div style={{ flex: 1 }}></div>
                  <span className={`role-badge ${assign.role === "backup" ? "rb-backup" : "rb-primary"}`}>
                    {(assign.role || "primary").toUpperCase()}
                  </span>
                  <div className="ar-since">{assign.since || "since Jun 10"}</div>
                  <div className="ar-by">
                    <code>{assign.by || "ch_a8f2e4b1 (template)"}</code>
                  </div>
                  <div className="ar-acts">
                    <div className="icon-btn" onClick={() => handleEditAssignment(assign)} style={{ cursor: "pointer" }}>✎</div>
                    <div className="icon-btn" onClick={() => handleDeleteAssignment(assign)} style={{ cursor: "pointer" }}>✕</div>
                  </div>
                </div>
              );
            })}
        </div>

        {/* Slot Group 3 */}
        <div className="assign-row-group">
          <div className="slot-group-header">
            <div className="sgh-num">3</div>
            <div className="sgh-slot">executive_signature</div>
            <span className="sgh-gate">→ financial_clearance</span>
            <div className="sgh-spacer"></div>
            <span className="sgh-add" style={{ cursor: "pointer" }} onClick={() => handleOpenAddPanel(3)}>
              + Add assignee
            </span>
          </div>
          {draftAssignments
            .filter((a) => a.slot_id === "3")
            .map((assign, idx) => {
              const badge = getVaultBadgeStyle(assign.vault_address || assign.owner_id);
              return (
                <div className="assign-row" key={idx}>
                  <span className={`vault-chip ${badge.bgClass}`} style={badge.style}>
                    <span className="vcdot" style={{ background: badge.color }}></span>
                    {assign.vault_address || assign.owner_id}
                  </span>
                  <div style={{ flex: 1 }}></div>
                  <span className={`role-badge ${assign.role === "backup" ? "rb-backup" : "rb-primary"}`}>
                    {(assign.role || "primary").toUpperCase()}
                  </span>
                  <div className="ar-since">{assign.since || "since Jun 10"}</div>
                  <div className="ar-by">
                    <code>{assign.by || "ch_a8f2e4b1 (template)"}</code>
                  </div>
                  <div className="ar-acts">
                    <div className="icon-btn" onClick={() => handleEditAssignment(assign)} style={{ cursor: "pointer" }}>✎</div>
                    <div className="icon-btn" onClick={() => handleDeleteAssignment(assign)} style={{ cursor: "pointer" }}>✕</div>
                  </div>
                </div>
              );
            })}
        </div>
      </div>

      <div className="assign-footer">
        <span className="af-pill ok">3 slots — all assigned</span>
        <span className="af-pill">{draftAssignments.length} vault assignments total</span>
        <div className="af-sp"></div>
        <button className="btn-add-assign" onClick={() => handleOpenAddPanel(3)}>
          + Add Assignment
        </button>
      </div>

      {(isDirty || saveFeedback) && (
        <div className="dirty-bar" style={{ margin: "0 -28px -12px" }}>
          <span className="db-icon">⚠</span>
          <div className="db-msg">
            {saveFeedback ? (
              <span>{saveFeedback}</span>
            ) : (
              <span>
                <strong>Unsaved changes</strong> — assignment updates take effect on the next thread instantiation.
              </span>
            )}
          </div>
          <button className="btn-discard" onClick={discardChanges} disabled={isSaving}>
            Discard
          </button>
          <button className="btn-save" onClick={saveChannelConfig} disabled={isSaving}>
            {isSaving ? "Saving..." : "Save Changes"}
          </button>
        </div>
      )}

      {/* Sliding Side Drawer Panel */}
      {isPanelOpen && (
        <>
          <div className="panel-scrim" onClick={() => setIsPanelOpen(false)}></div>
          <div className="assign-panel">
            <div className="ap-header">
              <div className="ap-hdr-top">
                <div className="ap-title">Add Assignment</div>
                <div className="ap-close" onClick={() => setIsPanelOpen(false)}>✕</div>
              </div>
              <div className="ap-subtitle">Bind a vault to a slot in this channel as primary, backup, or observer.</div>
              <div className="ap-ch-ctx">
                <span className="ap-ch-type">Channel</span>
                <span className="ap-ch-name">contract-execution</span>
              </div>
            </div>

            <div className="ap-body">
              {/* Slot selector */}
              <div className="ap-field">
                <div className="ap-label">
                  Slot <span style={{ color: "#EF4444" }}>*</span>
                </div>
                <div className="slot-selector">
                  <div
                    className={`slot-opt ${selectedSlot === 1 ? "selected" : ""}`}
                    onClick={() => setSelectedSlot(1)}
                  >
                    <div className={`so-radio ${selectedSlot === 1 ? "on" : ""}`}></div>
                    <div className="so-num">1</div>
                    <div className="so-slot">contract_draft</div>
                    <div className="so-current">vault_legal</div>
                  </div>
                  <div
                    className={`slot-opt ${selectedSlot === 2 ? "selected" : ""}`}
                    onClick={() => setSelectedSlot(2)}
                  >
                    <div className={`so-radio ${selectedSlot === 2 ? "on" : ""}`}></div>
                    <div className="so-num">2</div>
                    <div className="so-slot">financial_clearance</div>
                    <div className="so-current">vault_finance</div>
                    <div className="so-gate">→ contract_draft</div>
                  </div>
                  <div
                    className={`slot-opt ${selectedSlot === 3 ? "selected" : ""}`}
                    onClick={() => setSelectedSlot(3)}
                  >
                    <div className={`so-radio ${selectedSlot === 3 ? "on" : ""}`}></div>
                    <div className="so-num">3</div>
                    <div className="so-slot">executive_signature</div>
                    <div className="so-current">vault_direction</div>
                    <div className="so-gate">→ financial_clearance</div>
                  </div>
                </div>
              </div>

              {/* Vault selector */}
              <div className="ap-field">
                <div className="ap-label">
                  Vault <span style={{ color: "#EF4444" }}>*</span>
                </div>
                <div className="vault-select-wrap">
                  <div
                    className="vault-select-display"
                    onClick={() => setIsVaultDropdownOpen(!isVaultDropdownOpen)}
                  >
                    <div className="vsd-chip">
                      <span className="vsd-dot" style={{ background: "#2563EB" }}></span>
                      {selectedVault}
                    </div>
                    <div className="vsd-arr">▾</div>
                  </div>

                  {isVaultDropdownOpen && (
                    <div className="vault-dropdown">
                      {[
                        { name: "vault_finance", color: "#2563EB" },
                        { name: "vault_legal", color: "#7C3AED" },
                        { name: "vault_direction", color: "#444444" },
                        { name: "vault_treasury", color: "#059669" },
                        { name: "vault_hr", color: "#D97706" },
                        { name: "vault_ops", color: "#DC2626" },
                      ].map((v) => (
                        <div
                          key={v.name}
                          className={`vd-item ${selectedVault === v.name ? "selected" : ""}`}
                          onClick={() => {
                            setSelectedVault(v.name);
                            setIsVaultDropdownOpen(false);
                          }}
                        >
                          <span className="vd-dot" style={{ background: v.color }}></span>
                          {v.name}
                          {selectedVault === v.name && <span className="vd-check">✓</span>}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>

              {/* Role Cards */}
              <div className="ap-field">
                <div className="ap-label">
                  Role <span style={{ color: "#EF4444" }}>*</span>
                </div>
                <div className="role-cards">
                  <div
                    className={`role-card ${selectedRole === "primary" ? "selected-primary" : ""}`}
                    onClick={() => setSelectedRole("primary")}
                  >
                    <div className="rc-name primary">Primary</div>
                    <div className="rc-desc">Commits the slot by default. Required for thread progress.</div>
                  </div>
                  <div
                    className={`role-card ${selectedRole === "backup" ? "selected-backup" : ""}`}
                    onClick={() => setSelectedRole("backup")}
                  >
                    <div className="rc-name backup">Backup</div>
                    <div className="rc-desc">Can commit if primary is unavailable or delegates.</div>
                  </div>
                  <div
                    className={`role-card ${selectedRole === "observer" ? "selected-observer" : ""}`}
                    onClick={() => setSelectedRole("observer")}
                  >
                    <div className="rc-name observer">Observer</div>
                    <div className="rc-desc">Receives read-only telemetry and Stellar commit logs.</div>
                  </div>
                </div>
              </div>

              {/* Reason / Override note */}
              <div className="ap-field">
                <div className="ap-label">Assignment note</div>
                <textarea
                  className="ap-textarea"
                  placeholder="Reason for assignment change (recorded in channel log)..."
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                />
              </div>
            </div>

            <div className="ap-footer">
              <div className="ap-note">
                Template-level assignment. Applies to all future threads instantiated on this channel.
              </div>
              <button className="btn-cancel" onClick={() => setIsPanelOpen(false)}>
                Cancel
              </button>
              <button
                className="btn-save-assign"
                onClick={handleSaveAssignment}
              >
                Save Assignment
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
};


