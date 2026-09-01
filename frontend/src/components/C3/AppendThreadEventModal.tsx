import React, { useState, useMemo, useEffect } from "react";
import { createPortal } from "react-dom";
import { appendThreadEvent, ThreadEventResponse } from "@/services/api";
import { useVaultStore } from "@/store/vaultStore";
import { VaultEntry } from "@/types/vault";
import { TrustGroupSelect } from "./actions/TrustGroupSelect";
import {
  CreateApproval,
  CreateReject,
  CreateTransfer,
  CreateCollaborativeShare,
} from "../../../wailsjs/go/main/App";

import { useAuthStore } from "@/store/useAuthStore";

export interface AppendThreadEventSlidingViewProps {
  isOpen: boolean;
  activeWorkspaceName?: string;
  activeChannelTitle?: string;
  activeThreadId: string | null;
  activeThreadTitle?: string;
  onClose: () => void;
  onEventAppended?: (event: ThreadEventResponse) => void;
}

export type C3ActionMode = "normal" | "c3_share" | "approval" | "reject" | "transfer";

export const AppendThreadEventSlidingView: React.FC<AppendThreadEventSlidingViewProps> = ({
  isOpen,
  activeWorkspaceName = "Active Workspace",
  activeChannelTitle = "Active Channel",
  activeThreadId,
  activeThreadTitle = "Active Thread",
  onClose,
  onEventAppended,
}) => {
  const [actionMode, setActionMode] = useState<C3ActionMode>("normal");
  const [notes, setNotes] = useState("");
  const [trustGroupId, setTrustGroupId] = useState("");
  const [targetVaultId, setTargetVaultId] = useState("vault_target_member");
  const [reason, setReason] = useState("");

  // Source of Available Entries: Session Vault (useVaultStore)
  const vaultContext = useVaultStore((state) => state.vault);

  const allVaultEntries = useMemo(() => {
    if (!vaultContext?.Vault?.entries) return [];
    const entries: VaultEntry[] = [
      ...(vaultContext.Vault.entries.login || []),
      ...(vaultContext.Vault.entries.card || []),
      ...(vaultContext.Vault.entries.note || []),
      ...(vaultContext.Vault.entries.sshkey || []),
      ...(vaultContext.Vault.entries.identity || []),
    ];
    return entries;
  }, [vaultContext]);

  const [selectedEntryId, setSelectedEntryId] = useState<string>("");

  // Auto-select first available entry if none selected
  useEffect(() => {
    if (allVaultEntries.length > 0 && !selectedEntryId) {
      setSelectedEntryId(allVaultEntries[0].id);
    }
  }, [allVaultEntries, selectedEntryId]);

  // Status States
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleClose = () => {
    if (isLoading) return;
    setError(null);
    setSuccessMsg(null);
    onClose();
  };

  const selectedEntry = allVaultEntries.find((entry) => entry.id === selectedEntryId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccessMsg(null);

    if (!activeThreadId) {
      setError("No active thread selected. Please select a thread before taking action.");
      return;
    }

    if (!selectedEntryId) {
      setError("Please select a Vault Entry from your session vault.");
      return;
    }

    setIsLoading(true);
    const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";

    const resourceType = selectedEntry?.type || "vault_entry";
    const resourceId = selectedEntryId;

    console.log(`[C3 Diagnostic] Executing ${actionMode}:`, {
      threadID: activeThreadId,
      userID: useAuthStore.getState().user?.ID || "current",
      resourceType,
      resourceID: resourceId,
      trustGroupId,
      notes: notes.trim(),
      reason: reason.trim(),
      tokenPresent: Boolean(token),
    });

    try {
      if (actionMode === "normal") {
        // Standard Append Thread Event
        const newEvent = await appendThreadEvent({
          thread_id: activeThreadId,
          type: "entry.shared",
          payload: {
            ref_type: "vault_entry",
            entry_id: selectedEntryId,
            entry_name: selectedEntry?.entry_name || "Vault Entry",
            entry_type: selectedEntry?.type || "note",
            notes: notes.trim(),
          },
        });
        setSuccessMsg("✓ Normal Event appended to timeline.");
        if (onEventAppended) onEventAppended(newEvent);
      } else if (actionMode === "c3_share") {
        // Create Collaborative C3 Share
        if (!trustGroupId) {
          setError("Please select a Trust Group.");
          setIsLoading(false);
          return;
        }
        await CreateCollaborativeShare(
          token,
          activeThreadId,
          trustGroupId,
          selectedEntryId,
          targetVaultId,
          notes.trim(),
          "desktop_orchestrated_wrapped_dek",
          1
        );
        setSuccessMsg("✓ C3 Share created and appended to timeline.");
        if (onEventAppended) {
          onEventAppended({
            id: `evt_share_${Date.now()}`,
            thread_id: activeThreadId,
            type: "entry.shared",
            cursor: 0,
            payload: {
              share_entry_id: selectedEntryId,
              trust_group_id: trustGroupId,
              asset_cid: selectedEntryId,
              notes: notes.trim(),
            },
            created_at: new Date().toISOString(),
          } as ThreadEventResponse);
        }
      } else if (actionMode === "approval") {
        // Create Approval Action
        await CreateApproval(token, {
          resource_type: resourceType,
          resource_id: resourceId,
          message: notes.trim(),
          thread_id: activeThreadId,
        });
        setSuccessMsg("✓ Approval requested successfully.");
        if (onEventAppended) {
          onEventAppended({
            id: `evt_approval_${Date.now()}`,
            thread_id: activeThreadId,
            type: "c3.approval.requested",
            cursor: 0,
            payload: {
              resource_type: resourceType,
              resource_id: resourceId,
              notes: notes.trim(),
            },
            created_at: new Date().toISOString(),
          } as ThreadEventResponse);
        }
      } else if (actionMode === "reject") {
        // Create Reject Action
        if (!reason.trim()) {
          setError("Rejection reason is required.");
          setIsLoading(false);
          return;
        }
        await CreateReject(token, {
          resource_type: resourceType,
          resource_id: resourceId,
          message: reason.trim(),
          thread_id: activeThreadId,
        });
        setSuccessMsg("✓ Resource rejection recorded.");
        if (onEventAppended) {
          onEventAppended({
            id: `evt_reject_${Date.now()}`,
            thread_id: activeThreadId,
            type: "c3.reject.created",
            cursor: 0,
            payload: {
              resource_type: resourceType,
              resource_id: resourceId,
              notes: reason.trim(),
            },
            created_at: new Date().toISOString(),
          } as ThreadEventResponse);
        }
      } else if (actionMode === "transfer") {
        // Create Transfer Action
        if (!trustGroupId) {
          setError("Please select a Target Trust Group.");
          setIsLoading(false);
          return;
        }
        await CreateTransfer(token, {
          resource_type: resourceType,
          resource_id: resourceId,
          target_trust_group_id: trustGroupId,
          message: notes.trim(),
          thread_id: activeThreadId,
        });
        setSuccessMsg("✓ Resource transfer requested.");
        if (onEventAppended) {
          onEventAppended({
            id: `evt_transfer_${Date.now()}`,
            thread_id: activeThreadId,
            type: "c3.transfer.requested",
            cursor: 0,
            payload: {
              resource_type: resourceType,
              resource_id: resourceId,
              target_trust_group_id: trustGroupId,
              notes: notes.trim(),
            },
            created_at: new Date().toISOString(),
          } as ThreadEventResponse);
        }
      }
    } catch (err: any) {
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error(`[C3 Diagnostic] Backend Error executing ${actionMode}:`, errMsg, err);
      setError(errMsg);
    } finally {
      setIsLoading(false);
    }
  };

  return createPortal(
    <>
      {/* Backdrop */}
      <div className="c3-sliding-view-backdrop" onClick={handleClose} />

      {/* Viewport-Anchored Panel */}
      <div className="c3-sliding-view-container">
        <div className="slide-panel">
          {/* Header */}
          <div className="sp-header">
            <div className="sp-header-row">
              <div>
                <div className="sp-title">Append Thread Event & C3 Actions</div>
                <div className="sp-subtitle">
                  Execute generic C3 collaboration primitives or append normal timeline events
                </div>
              </div>
              <div
                className="sp-close"
                onClick={handleClose}
                role="button"
                tabIndex={0}
              >
                ✕
              </div>
            </div>
          </div>

          {/* Form Body */}
          <form onSubmit={handleSubmit} style={{ display: "contents" }}>
            <div className="sp-body">
              {/* Success Message Alert */}
              {successMsg && (
                <div
                  style={{
                    backgroundColor: "rgba(16, 185, 129, 0.08)",
                    border: "1px solid rgba(16, 185, 129, 0.3)",
                    borderRadius: "6px",
                    padding: "10px 12px",
                    color: "#059669",
                    fontSize: "13px",
                    fontWeight: 600,
                    display: "flex",
                    alignItems: "center",
                    gap: "8px",
                  }}
                >
                  <span>{successMsg}</span>
                </div>
              )}

              {/* Error Alert */}
              {error && (
                <div
                  style={{
                    backgroundColor: "rgba(239, 68, 68, 0.08)",
                    border: "1px solid rgba(239, 68, 68, 0.3)",
                    borderRadius: "6px",
                    padding: "10px 12px",
                    color: "#DC2626",
                    fontSize: "13px",
                    display: "flex",
                    alignItems: "flex-start",
                    gap: "8px",
                  }}
                >
                  <span>⚠️</span>
                  <div style={{ flex: 1 }}>{error}</div>
                </div>
              )}

              {/* Hierarchy Context Display */}
              <div>
                <div className="fl">Target Thread Context</div>
                <div className="channel-flow-box">
                  <div className="cfb-row" style={{ fontSize: "12px", color: "#333" }}>
                    <span style={{ color: "#888" }}>Workspace:</span>
                    <strong>{activeWorkspaceName}</strong>
                    <span style={{ color: "#ccc" }}>→</span>
                    <span style={{ color: "#888" }}>Channel:</span>
                    <strong>{activeChannelTitle}</strong>
                  </div>
                  <div style={{ marginTop: "6px", fontSize: "13px", fontWeight: 600, color: "#C8922A" }}>
                    Thread: {activeThreadTitle}
                  </div>
                </div>
              </div>

              {/* Action Mode Segmented Choice */}
              <div>
                <div className="fl">Action Type</div>
                <div
                  style={{
                    display: "grid",
                    gridTemplateColumns: "repeat(auto-fit, minmax(100px, 1fr))",
                    gap: "6px",
                    backgroundColor: "#161B22",
                    padding: "4px",
                    borderRadius: "6px",
                    border: "1px solid #30363D",
                  }}
                >
                  <button
                    type="button"
                    onClick={() => setActionMode("normal")}
                    style={{
                      padding: "6px 8px",
                      fontSize: "11px",
                      fontWeight: 600,
                      borderRadius: "4px",
                      border: "none",
                      backgroundColor: actionMode === "normal" ? "#238636" : "transparent",
                      color: actionMode === "normal" ? "#FFFFFF" : "#8B949E",
                      cursor: "pointer",
                    }}
                  >
                    📄 Normal Event
                  </button>
                  <button
                    type="button"
                    onClick={() => setActionMode("c3_share")}
                    style={{
                      padding: "6px 8px",
                      fontSize: "11px",
                      fontWeight: 600,
                      borderRadius: "4px",
                      border: "none",
                      backgroundColor: actionMode === "c3_share" ? "#238636" : "transparent",
                      color: actionMode === "c3_share" ? "#FFFFFF" : "#8B949E",
                      cursor: "pointer",
                    }}
                  >
                    🤝 Create C3 Share
                  </button>
                  <button
                    type="button"
                    onClick={() => setActionMode("approval")}
                    style={{
                      padding: "6px 8px",
                      fontSize: "11px",
                      fontWeight: 600,
                      borderRadius: "4px",
                      border: "none",
                      backgroundColor: actionMode === "approval" ? "#D97706" : "transparent",
                      color: actionMode === "approval" ? "#FFFFFF" : "#8B949E",
                      cursor: "pointer",
                    }}
                  >
                    👍 Approval
                  </button>
                  <button
                    type="button"
                    onClick={() => setActionMode("reject")}
                    style={{
                      padding: "6px 8px",
                      fontSize: "11px",
                      fontWeight: 600,
                      borderRadius: "4px",
                      border: "none",
                      backgroundColor: actionMode === "reject" ? "#DA3633" : "transparent",
                      color: actionMode === "reject" ? "#FFFFFF" : "#8B949E",
                      cursor: "pointer",
                    }}
                  >
                    🚫 Reject
                  </button>
                  <button
                    type="button"
                    onClick={() => setActionMode("transfer")}
                    style={{
                      padding: "6px 8px",
                      fontSize: "11px",
                      fontWeight: 600,
                      borderRadius: "4px",
                      border: "none",
                      backgroundColor: actionMode === "transfer" ? "#1F6FEB" : "transparent",
                      color: actionMode === "transfer" ? "#FFFFFF" : "#8B949E",
                      cursor: "pointer",
                    }}
                  >
                    🔄 Transfer
                  </button>
                </div>
              </div>

              {/* Session Vault Entry Select (Resource Reference) */}
              <div>
                <div className="fl">
                  Referenced Vault Entry <span style={{ color: "#EF4444" }}>*</span>
                </div>
                {allVaultEntries.length > 0 ? (
                  <select
                    className="prop-input"
                    value={selectedEntryId}
                    onChange={(e) => setSelectedEntryId(e.target.value)}
                    disabled={isLoading}
                    style={{ height: "38px" }}
                  >
                    {allVaultEntries.map((entry) => (
                      <option key={entry.id} value={entry.id}>
                        [{entry.type.toUpperCase()}] {entry.entry_name || "Untitled Entry"} ({entry.id.slice(0, 8)}…)
                      </option>
                    ))}
                  </select>
                ) : (
                  <div style={{ fontSize: "12px", color: "#B45309", backgroundColor: "rgba(245, 158, 11, 0.08)", border: "1px solid rgba(245, 158, 11, 0.3)", borderRadius: "6px", padding: "8px 10px" }}>
                    No Vault Entries found in current session vault. Please unlock or sync your vault.
                  </div>
                )}
              </div>

              {/* Trust Group Selector for C3 Share and Transfer */}
              {(actionMode === "c3_share" || actionMode === "transfer") && (
                <div>
                  <div className="fl">
                    Trust Group <span style={{ color: "#EF4444" }}>*</span>
                  </div>
                  <TrustGroupSelect
                    value={trustGroupId}
                    onChange={setTrustGroupId}
                    disabled={isLoading}
                  />
                </div>
              )}

              {/* Mandatory Rejection Reason for Reject Action */}
              {actionMode === "reject" ? (
                <div>
                  <div className="fl">
                    Rejection Reason <span style={{ color: "#EF4444" }}>*</span>
                  </div>
                  <textarea
                    className="prop-input"
                    placeholder="Provide detailed justification for rejecting this resource..."
                    value={reason}
                    onChange={(e) => setReason(e.target.value)}
                    disabled={isLoading}
                    rows={3}
                    required
                    style={{ resize: "vertical" }}
                  />
                  <div style={{ fontSize: "11px", color: "#8B949E", marginTop: "4px" }}>
                    ℹ️ Rejection records an auditable decision in the timeline. It does not mutate or delete the resource.
                  </div>
                </div>
              ) : (
                /* Notes / Message for Other Actions */
                <div>
                  <div className="fl">
                    Notes / Message <span style={{ color: "#999", fontWeight: 400 }}>(optional)</span>
                  </div>
                  <textarea
                    className="prop-input"
                    placeholder={
                      actionMode === "c3_share"
                        ? "Notes for Trust Group recipients..."
                        : actionMode === "approval"
                        ? "Instructions for reviewers..."
                        : actionMode === "transfer"
                        ? "Transfer request message..."
                        : "e.g. Countersigned contract draft committed to vault."
                    }
                    value={notes}
                    onChange={(e) => setNotes(e.target.value)}
                    disabled={isLoading}
                    rows={3}
                    style={{ resize: "vertical" }}
                  />
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="sp-footer">
              <button
                type="submit"
                className="start-btn"
                disabled={
                  isLoading ||
                  allVaultEntries.length === 0 ||
                  (actionMode === "reject" && !reason.trim())
                }
                style={{
                  opacity:
                    isLoading ||
                    allVaultEntries.length === 0 ||
                    (actionMode === "reject" && !reason.trim())
                      ? 0.6
                      : 1,
                  cursor:
                    isLoading ||
                    allVaultEntries.length === 0 ||
                    (actionMode === "reject" && !reason.trim())
                      ? "not-allowed"
                      : "pointer",
                  backgroundColor:
                    actionMode === "reject"
                      ? "#DA3633"
                      : actionMode === "transfer"
                      ? "#1F6FEB"
                      : actionMode === "approval"
                      ? "#D97706"
                      : "#238636",
                }}
              >
                {isLoading ? (
                  <span>Executing...</span>
                ) : actionMode === "c3_share" ? (
                  "🤝 Create C3 Share"
                ) : actionMode === "approval" ? (
                  "👍 Request Approval"
                ) : actionMode === "reject" ? (
                  "🚫 Reject Resource"
                ) : actionMode === "transfer" ? (
                  "🔄 Request Transfer"
                ) : (
                  "▶ Append Normal Event"
                )}
              </button>
              <div className="footer-note">
                Action will be executed and recorded in the authoritative thread timeline.
              </div>
            </div>
          </form>
        </div>
      </div>
    </>,
    document.body
  );
};

export const AppendThreadEventModal = AppendThreadEventSlidingView;
