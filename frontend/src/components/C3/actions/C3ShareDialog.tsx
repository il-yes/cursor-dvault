import React, { useState } from "react";
import { ResourceReference, C3ResourceTypes } from "../domain/resource";
import { TrustGroupSelect } from "./TrustGroupSelect";
import { CreateCollaborativeShare } from "../../../../wailsjs/go/main/App";

import { useAuthStore } from "@/store/useAuthStore";

interface C3ShareDialogProps {
  isOpen: boolean;
  onClose: () => void;
  resourceRef?: ResourceReference;
  threadId?: string;
  onShareCreated?: (shareRef: any) => void;
}

export const C3ShareDialog: React.FC<C3ShareDialogProps> = ({
  isOpen,
  onClose,
  resourceRef,
  threadId = "default_thread",
  onShareCreated,
}) => {
  const [trustGroupId, setTrustGroupId] = useState("");
  const [targetVaultId, setTargetVaultId] = useState("vault_target_member");
  const [notes, setNotes] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const resType = resourceRef?.resourceType || C3ResourceTypes.VaultEntry;
  const resId = resourceRef?.resourceId || "entry_select_placeholder";

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!trustGroupId) {
      setError("Please select a Trust Group.");
      return;
    }

    setIsLoading(true);
    try {
      // Execute C3 Collaborative Share Wails API
      const token = useAuthStore.getState().jwtToken || localStorage.getItem("jwt_token") || "";
      const shareRef = await CreateCollaborativeShare(
        token,
        threadId,
        trustGroupId,
        resId, // assetCID or resourceID
        targetVaultId,
        notes,
        "desktop_orchestrated_wrapped_dek",
        1
      );

      setIsLoading(false);
      if (onShareCreated) onShareCreated(shareRef);
      onClose();
    } catch (err: any) {
      setIsLoading(false);
      const errMsg = typeof err === "string" ? err : err?.message || JSON.stringify(err);
      console.error("[C3 Diagnostic] CreateCollaborativeShare Error:", errMsg, err);
      setError(errMsg);
    }
  };

  return (
    <div
      role="dialog"
      aria-labelledby="c3-share-dialog-title"
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
          <h3 id="c3-share-dialog-title" style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#F0F6FC" }}>
            🤝 Create C3 Share
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

          {/* Trust Group Selector */}
          <div>
            <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
              Trust Group *
            </label>
            <TrustGroupSelect
              value={trustGroupId}
              onChange={setTrustGroupId}
              disabled={isLoading}
            />
          </div>

          {/* Target Vault ID */}
          <div>
            <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
              Target Member Vault ID
            </label>
            <input
              type="text"
              value={targetVaultId}
              onChange={(e) => setTargetVaultId(e.target.value)}
              disabled={isLoading}
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

          {/* Optional Message */}
          <div>
            <label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
              Optional Message
            </label>
            <textarea
              placeholder="Add optional notes for recipient Trust Group members..."
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
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
              {isLoading ? "Creating..." : "Create C3 Share"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
