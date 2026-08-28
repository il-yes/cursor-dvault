import React, { useState } from "react";
import { ResourceReference, C3ResourceTypes } from "../domain/resource";
import { CreateReject } from "../../../../wailsjs/go/main/App";

import { useAuthStore } from "@/store/useAuthStore";

interface RejectActionDialogProps {
  isOpen: boolean;
  onClose: () => void;
  resourceRef?: ResourceReference;
  threadId?: string;
  onSuccess?: () => void;
}

export const RejectActionDialog: React.FC<RejectActionDialogProps> = ({
  isOpen,
  onClose,
  resourceRef,
  threadId = "default_thread",
  onSuccess,
}) => {
  const [reason, setReason] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const resType = resourceRef?.resourceType || C3ResourceTypes.VaultEntry;
  const resId = resourceRef?.resourceId || "resource_id_placeholder";

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!reason.trim()) {
      setError("Rejection reason is required.");
      return;
    }

    setIsLoading(true);
    try {
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      await CreateReject(token, {
        resource_type: resType,
        resource_id: resId,
        source_event_id: resourceRef?.sourceEventId || "",
        message: reason.trim(),
        thread_id: threadId,
      });

      setIsLoading(false);
      if (onSuccess) onSuccess();
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] CreateReject Error:", errMsg, err);
      setError(errMsg);
    }
  };

  return (
    <div
      role="dialog"
      aria-labelledby="reject-dialog-title"
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
          <h3 id="reject-dialog-title" style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#F87171" }}>
            🚫 Reject Resource
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
        <form onSubmit={handleSubmit} style={{ padding: "20px", display: "flex", flexDirection: "column", gap: "16px" }}>
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

          {/* Resource Context & Non-Destructive Banner */}
          <div
            style={{
              padding: "10px 12px",
              backgroundColor: "rgba(239, 68, 68, 0.08)",
              border: "1px solid rgba(239, 68, 68, 0.2)",
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
              ℹ️ Rejection records an auditable decision against this resource. It does <strong>not</strong> delete, mutate, or revoke the underlying resource.
            </div>
          </div>

          {/* Rejection Reason Input */}
          <div>
            <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
              Rejection Reason *
            </label>
            <textarea
              placeholder="Explain why this resource is being rejected..."
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              disabled={isLoading}
              rows={3}
              required
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
              disabled={isLoading || !reason.trim()}
              style={{
                padding: "8px 16px",
                backgroundColor: "#DA3633",
                border: "none",
                borderRadius: "6px",
                color: "#FFFFFF",
                fontSize: "13px",
                fontWeight: 600,
                cursor: isLoading || !reason.trim() ? "not-allowed" : "pointer",
                opacity: isLoading || !reason.trim() ? 0.6 : 1,
              }}
            >
              {isLoading ? "Submitting..." : "Reject Resource"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
