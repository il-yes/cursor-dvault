import React, { useState } from "react";
import { ResourceReference, C3ResourceTypes } from "../domain/resource";
import { TrustGroupSelect } from "./TrustGroupSelect";
import { C3ActionStatus } from "./C3ActionStatus";
import {
  CreateTransfer,
  ApproveTransferAction,
  RejectTransferAction,
  CompleteTransferAction,
} from "../../../../wailsjs/go/main/App";

import { useAuthStore } from "@/store/useAuthStore";

interface TransferActionDialogProps {
  isOpen: boolean;
  onClose: () => void;
  resourceRef?: ResourceReference;
  threadId?: string;
  transferId?: string;
  existingStatus?: string;
  targetTrustGroupId?: string;
  onSuccess?: () => void;
}

export const TransferActionDialog: React.FC<TransferActionDialogProps> = ({
  isOpen,
  onClose,
  resourceRef,
  threadId = "default_thread",
  transferId,
  existingStatus,
  targetTrustGroupId,
  onSuccess,
}) => {
  const [selectedTrustGroup, setSelectedTrustGroup] = useState(targetTrustGroupId || "");
  const [message, setMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const isExisting = Boolean(transferId);
  const normStatus = (existingStatus || "").toLowerCase();
  const resType = resourceRef?.resourceType || C3ResourceTypes.VaultEntry;
  const resId = resourceRef?.resourceId || "resource_id_placeholder";

  const handleRequestTransfer = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!selectedTrustGroup) {
      setError("Target Trust Group is required.");
      return;
    }

    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await CreateTransfer(token, {
        resource_type: resType,
        resource_id: resId,
        source_event_id: resourceRef?.sourceEventId || "",
        target_trust_group_id: selectedTrustGroup,
        message: message.trim(),
        thread_id: threadId,
      });

      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] CreateTransfer Error:", errMsg, err);
      setError(errMsg);
    }
  };

  const handleApproveTransfer = async () => {
    if (!transferId) return;
    setError(null);
    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await ApproveTransferAction(token, {
        transfer_id: transferId,
        thread_id: threadId,
      });
      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] ApproveTransferAction Error:", errMsg, err);
      setError(errMsg);
    }
  };

  const handleRejectTransfer = async () => {
    if (!transferId) return;
    setError(null);
    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await RejectTransferAction(token, {
        transfer_id: transferId,
        thread_id: threadId,
      });
      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] RejectTransferAction Error:", errMsg, err);
      setError(errMsg);
    }
  };

  const handleCompleteTransfer = async () => {
    if (!transferId) return;
    setError(null);
    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await CompleteTransferAction(token, {
        transfer_id: transferId,
        thread_id: threadId,
      });
      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] CompleteTransferAction Error:", errMsg, err);
      setError(errMsg);
    }
  };

  return (
    <div
      role="dialog"
      aria-labelledby="transfer-dialog-title"
      style={{
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: "rgba(0, 0, 0, 0.75)",
        backdropFilter: "blur(4px)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 1000,
      }}
    >
      <div
        style={{
          backgroundColor: "#161B22",
          border: "1px solid #30363D",
          borderRadius: "8px",
          width: "100%",
          maxWidth: "480px",
          boxShadow: "0 20px 25px -5px rgba(0, 0, 0, 0.5)",
          display: "flex",
          flexDirection: "column",
          overflow: "hidden",
        }}
      >
        {/* Header */}
        <div
          style={{
            padding: "16px 20px",
            borderBottom: "1px solid rgba(255, 255, 255, 0.08)",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <h3 id="transfer-dialog-title" style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#60A5FA" }}>
            🔄 {isExisting ? "Transfer Lifecycle Workflow" : "Transfer Resource"}
          </h3>
          <button
            type="button"
            onClick={onClose}
            disabled={isLoading}
            aria-label="Close dialog"
            style={{
              background: "none",
              border: "none",
              color: "#8B949E",
              fontSize: "18px",
              cursor: "pointer",
            }}
          >
            ✕
          </button>
        </div>

        {/* Body */}
        <div style={{ padding: "20px", display: "flex", flexDirection: "column", gap: "16px" }}>
          {error && (
            <div
              style={{
                padding: "10px 12px",
                backgroundColor: "rgba(239, 68, 68, 0.1)",
                border: "1px solid rgba(239, 68, 68, 0.3)",
                borderRadius: "6px",
                color: "#F87171",
                fontSize: "13px",
              }}
            >
              ⚠️ {error}
            </div>
          )}

          {/* Resource & Workflow Banner */}
          <div
            style={{
              padding: "10px 12px",
              backgroundColor: "rgba(59, 130, 246, 0.08)",
              border: "1px solid rgba(59, 130, 246, 0.2)",
              borderRadius: "6px",
              fontSize: "12px",
              color: "#C9D1D9",
              display: "flex",
              flexDirection: "column",
              gap: "4px",
            }}
          >
            <div>
              <span style={{ color: "#8B949E" }}>Referenced Resource:</span>{" "}
              <strong>{resType}</strong> ({resId})
            </div>
            <div style={{ color: "#8B949E", fontSize: "11px", marginTop: "2px" }}>
              🛡️ Transfer is an explicit workflow boundary targeting a TrustGroup. It does <strong>not</strong> grant immediate plaintext access without TrustGroup authorization.
            </div>
          </div>

          {isExisting ? (
            <div style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
              <div style={{ display: "flex", alignItems: "center", justifyBetween: "space-between", gap: "10px" }}>
                <span style={{ color: "#8B949E", fontSize: "13px" }}>Workflow State:</span>
                <C3ActionStatus status={normStatus} />
              </div>

              {targetTrustGroupId && (
                <div style={{ fontSize: "13px", color: "#C9D1D9" }}>
                  <span style={{ color: "#8B949E" }}>Target TrustGroup:</span> <code>{targetTrustGroupId}</code>
                </div>
              )}

              {/* Lifecycle Actions based on Current State */}
              <div style={{ marginTop: "12px", display: "flex", justifyContent: "flex-end", gap: "10px" }}>
                <button
                  type="button"
                  onClick={onClose}
                  style={{
                    padding: "8px 16px",
                    backgroundColor: "transparent",
                    border: "1px solid #30363D",
                    borderRadius: "6px",
                    color: "#C9D1D9",
                    fontSize: "13px",
                    cursor: "pointer",
                  }}
                >
                  Close
                </button>

                {normStatus === "transfer_requested" && (
                  <>
                    <button
                      type="button"
                      onClick={handleRejectTransfer}
                      disabled={isLoading}
                      style={{
                        padding: "8px 16px",
                        backgroundColor: "#DA3633",
                        border: "none",
                        borderRadius: "6px",
                        color: "#FFFFFF",
                        fontSize: "13px",
                        fontWeight: 600,
                        cursor: isLoading ? "not-allowed" : "pointer",
                      }}
                    >
                      Reject Transfer
                    </button>
                    <button
                      type="button"
                      onClick={handleApproveTransfer}
                      disabled={isLoading}
                      style={{
                        padding: "8px 16px",
                        backgroundColor: "#238636",
                        border: "none",
                        borderRadius: "6px",
                        color: "#FFFFFF",
                        fontSize: "13px",
                        fontWeight: 600,
                        cursor: isLoading ? "not-allowed" : "pointer",
                      }}
                    >
                      Approve Transfer
                    </button>
                  </>
                )}

                {normStatus === "transfer_approved" && (
                  <button
                    type="button"
                    onClick={handleCompleteTransfer}
                    disabled={isLoading}
                    style={{
                      padding: "8px 16px",
                      backgroundColor: "#1F6FEB",
                      border: "none",
                      borderRadius: "6px",
                      color: "#FFFFFF",
                      fontSize: "13px",
                      fontWeight: 600,
                      cursor: isLoading ? "not-allowed" : "pointer",
                    }}
                  >
                    Complete Transfer
                  </button>
                )}
              </div>
            </div>
          ) : (
            <form onSubmit={handleRequestTransfer} style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
              {/* Target Trust Group Selection */}
              <div>
                <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
                  Target Trust Group *
                </label>
                <TrustGroupSelect
                  value={selectedTrustGroup}
                  onChange={setSelectedTrustGroup}
                  disabled={isLoading}
                />
              </div>

              {/* Message */}
              <div>
                <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
                  Transfer Request Message
                </label>
                <textarea
                  placeholder="Please take responsibility for this entry..."
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  disabled={isLoading}
                  rows={3}
                  style={{
                    width: "100%",
                    padding: "8px 12px",
                    backgroundColor: "#0D1117",
                    border: "1px solid #30363D",
                    borderRadius: "6px",
                    color: "#F0F6FC",
                    fontSize: "13px",
                    boxSizing: "border-box",
                    resize: "vertical",
                  }}
                />
              </div>

              {/* Footer */}
              <div
                style={{
                  marginTop: "8px",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "flex-end",
                  gap: "10px",
                }}
              >
                <button
                  type="button"
                  onClick={onClose}
                  disabled={isLoading}
                  style={{
                    padding: "8px 16px",
                    backgroundColor: "transparent",
                    border: "1px solid #30363D",
                    borderRadius: "6px",
                    color: "#C9D1D9",
                    fontSize: "13px",
                    cursor: "pointer",
                  }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isLoading || !selectedTrustGroup}
                  style={{
                    padding: "8px 16px",
                    backgroundColor: "#1F6FEB",
                    border: "none",
                    borderRadius: "6px",
                    color: "#FFFFFF",
                    fontSize: "13px",
                    fontWeight: 600,
                    cursor: isLoading || !selectedTrustGroup ? "not-allowed" : "pointer",
                  }}
                >
                  {isLoading ? "Requesting..." : "Request Transfer"}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};
