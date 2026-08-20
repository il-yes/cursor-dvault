import React, { useState } from "react";
import { desktopResourceService, OpenedResourceResult } from "@/services/desktopResourceService";

interface C3ResourceCardProps {
  refType: "share_entry" | "storage_asset" | string;
  shareEntryId?: string;
  trustGroupId?: string;
  cid?: string;
  author?: string;
  createdAt?: string;
}

export const C3ResourceCard: React.FC<C3ResourceCardProps> = ({
  refType,
  shareEntryId,
  trustGroupId,
  cid,
  author,
  createdAt,
}) => {
  const [status, setStatus] = useState<"idle" | "loading" | "resolved" | "error">("idle");
  const [result, setResult] = useState<OpenedResourceResult | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleOpen = async () => {
    setStatus("loading");
    setErrorMessage(null);

    try {
      const res = await desktopResourceService.openResource({
        refType,
        shareEntryId,
        trustGroupId,
        cid,
      });
      setResult(res);
      setStatus("resolved");
    } catch (err: any) {
      setErrorMessage(err?.message || "Unable to open protected resource");
      setStatus("error");
    }
  };

  return (
    <div
      className="c3-resource-card"
      style={{
        marginTop: "8px",
        padding: "12px 14px",
        backgroundColor: "rgba(16, 185, 129, 0.06)",
        border: "1px solid rgba(16, 185, 129, 0.2)",
        borderRadius: "6px",
        fontFamily: "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
      }}
    >
      {/* Header Info */}
      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "8px" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "6px" }}>
          <span style={{ fontSize: "14px" }}>🔐</span>
          <strong style={{ fontSize: "12px", color: "#34D399", fontWeight: 600 }}>
            Protected Resource
          </strong>
        </div>
        {author && (
          <span style={{ fontSize: "11px", color: "#8B949E" }}>
            Shared by {author}
          </span>
        )}
      </div>

      {/* ID & Context Details */}
      <div style={{ fontSize: "11px", color: "#8B949E", marginBottom: "10px" }}>
        {shareEntryId && <div>Ref ID: <span style={{ color: "#C9D1D9", fontFamily: "monospace" }}>{shareEntryId}</span></div>}
        {trustGroupId && <div>TrustGroup: <span style={{ color: "#C9D1D9", fontFamily: "monospace" }}>{trustGroupId}</span></div>}
      </div>

      {/* State Renderers */}
      {status === "idle" && (
        <button
          type="button"
          onClick={handleOpen}
          style={{
            padding: "6px 14px",
            fontSize: "12px",
            fontWeight: 600,
            color: "#FFFFFF",
            backgroundColor: "#10B981",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer",
            transition: "background 0.15s ease",
          }}
        >
          Open Resource
        </button>
      )}

      {status === "loading" && (
        <div style={{ fontSize: "12px", color: "#58A6FF", display: "flex", alignItems: "center", gap: "6px" }}>
          <span>⏳ Opening protected resource…</span>
        </div>
      )}

      {status === "resolved" && result && (
        <div
          style={{
            padding: "10px",
            backgroundColor: "#161B22",
            border: "1px solid rgba(255, 255, 255, 0.08)",
            borderRadius: "4px",
            marginTop: "6px",
          }}
        >
          <div style={{ fontSize: "12px", fontWeight: 600, color: "#F0F6FC", marginBottom: "4px" }}>
            {result.title || "Decrypted Application Payload"}
          </div>
          <pre
            style={{
              margin: 0,
              padding: "8px",
              backgroundColor: "#0D1117",
              borderRadius: "4px",
              fontSize: "11px",
              color: "#A5D6FF",
              overflowX: "auto",
              whiteSpace: "pre-wrap",
              wordBreak: "break-all",
            }}
          >
            {result.content}
          </pre>
        </div>
      )}

      {status === "error" && errorMessage && (
        <div
          style={{
            padding: "8px 10px",
            backgroundColor: "rgba(239, 68, 68, 0.1)",
            border: "1px solid rgba(239, 68, 68, 0.25)",
            borderRadius: "4px",
            fontSize: "12px",
            color: "#EF4444",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <span>⛔ {errorMessage}</span>
          <button
            type="button"
            onClick={handleOpen}
            style={{
              background: "transparent",
              border: "1px solid rgba(239, 68, 68, 0.4)",
              color: "#EF4444",
              borderRadius: "3px",
              padding: "2px 8px",
              fontSize: "11px",
              cursor: "pointer",
            }}
          >
            Retry
          </button>
        </div>
      )}
    </div>
  );
};
