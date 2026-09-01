import React, { useState } from "react";
import { ResourceReference, C3ResourceTypes } from "../domain/resource";
import { CreateApproval, ApproveAction } from "../../../../wailsjs/go/main/App";
import { C3ActionStatus } from "./C3ActionStatus";

import { useAuthStore } from "@/store/useAuthStore";

interface ApprovalActionDialogProps {
  isOpen: boolean;
  onClose: () => void;
  resourceRef?: ResourceReference;
  threadId?: string;
  approvalId?: string;
  existingStatus?: string;
  onSuccess?: () => void;
}

export const ApprovalActionDialog: React.FC<ApprovalActionDialogProps> = ({
  isOpen,
  onClose,
  resourceRef,
  threadId = "default_thread",
  approvalId,
  existingStatus,
  onSuccess,
}) => {
  const [message, setMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const isExisting = Boolean(approvalId);
  const isApproved = existingStatus === "approved";
  const resType = resourceRef?.resourceType || C3ResourceTypes.VaultEntry;
  const resId = resourceRef?.resourceId || "resource_id_placeholder";

  const handleRequestApproval = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await CreateApproval(token, {
        resource_type: resType,
        resource_id: resId,
        source_event_id: resourceRef?.sourceEventId || "",
        message: message.trim(),
        thread_id: threadId,
      });
      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] CreateApproval Error:", errMsg, err);
      setError(errMsg);
    }
  };

  const handleApprove = async () => {
    if (!approvalId) return;
    setError(null);
    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await ApproveAction(token, {
        approval_id: approvalId,
        thread_id: threadId,
      });
      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] ApproveAction Error:", errMsg, err);
      setError(errMsg);
    }
  };

  return (
    <div
      role="dialog"
      aria-labelledby="approval-dialog-title"
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
          maxWidth: "460px",
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
          <h3 id="approval-dialog-title" style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#F0F6FC" }}>
            {isExisting ? "Approval Status" : "Request Approval"}
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

          {/* Resource Reference Context */}
          <div
            style={{
              padding: "10px 12px",
              backgroundColor: "rgba(37, 99, 235, 0.08)",
              border: "1px solid rgba(37, 99, 235, 0.2)",
              borderRadius: "6px",
              fontSize: "12px",
              color: "#C9D1D9",
            }}
          >
            <div>
              <span style={{ color: "#8B949E" }}>Referenced Resource:</span>{" "}
              <strong>{resType}</strong> ({resId})
            </div>
          </div>

          {isExisting ? (
            <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
              <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                <span style={{ color: "#8B949E", fontSize: "13px" }}>Current Status:</span>
                <C3ActionStatus status={existingStatus || "requested"} />
              </div>

              {isApproved ? (
                <div style={{ fontSize: "13px", color: "#4ADE80" }}>
                  ✓ Approval has been completed.
                </div>
              ) : (
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
                  <button
                    type="button"
                    onClick={handleApprove}
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
                    {isLoading ? "Approving..." : "Approve"}
                  </button>
                </div>
              )}
            </div>
          ) : (
            <form onSubmit={handleRequestApproval} style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
              <div>
                <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
                  Approval Note / Instructions
                </label>
                <textarea
                  placeholder="Please review and approve this resource for collaboration..."
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
                  }}
                />
              </div>

              <div
                style={{
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
                  {isLoading ? "Submitting..." : "Request Approval"}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};
