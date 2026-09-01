import React from "react";

export interface RemoteVaultItem {
  id: string;
  endpoint: string;
  status: "trusted" | "pending" | "observer";
  lastSeen: string;
  cursor: string;
  proto: string;
  pendingItems?: string;
  alert?: boolean;
}

interface RemoteVaultCardProps {
  vault: RemoteVaultItem;
  isExpanded: boolean;
  onToggleExpand: () => void;
}

export const RemoteVaultCard: React.FC<RemoteVaultCardProps> = ({
  vault,
  isExpanded,
  onToggleExpand,
}) => {
  const isTrusted = vault.status === "trusted";

  return (
    <div className={`vault-card ${isTrusted ? "trusted" : "pending-trust"}`}>
      <div
        className={`vc-header ${isExpanded ? "expanded" : ""}`}
        onClick={onToggleExpand}
      >
        <span className={`vc-trust ${isTrusted ? "vct-trusted" : "vct-pending"}`}>
          {isTrusted ? "✓ Trusted" : "⧖ Pending"}
        </span>
        <span className="vc-id">{vault.id}</span>
        <span className="vc-endpoint">{vault.endpoint}</span>
        <span className="vc-last-seen" style={!isTrusted ? { color: "#D97706" } : {}}>
          last seen {vault.lastSeen}
        </span>
        <span className="vc-cursor" style={!isTrusted ? { color: "#D97706" } : {}}>
          cursor: {vault.cursor}
        </span>
        <span
          className="vc-proto"
          style={!isTrusted ? { background: "#FFF7ED", borderColor: "#FED7AA", color: "#D97706" } : {}}
        >
          {vault.proto}
        </span>
      </div>

      {isExpanded && vault.id === "vault:supplier-x" && (
        <div style={{ padding: 0, background: "#fafafa" }}>
          <div className="vi-grid vi-col-hdrs">
            <div className="vi-col-lbl" style={{ paddingLeft: "12px" }}>Exchange ID</div>
            <div className="vi-col-lbl">Status</div>
            <div className="vi-col-lbl">Cursor</div>
            <div className="vi-col-lbl">Retries</div>
            <div className="vi-col-lbl">Updated</div>
            <div className="vi-col-lbl">Actions</div>
          </div>

          <div className="vi-grid vi-row">
            <div className="vi-id" style={{ paddingLeft: "12px" }}>exch_a8f2e4b1c3d9f7a2e5…</div>
            <div style={{ padding: "0 6px" }}>
              <span className="sync-badge sb-completed">
                <span className="vcdot" style={{ background: "#059669" }}></span>completed
              </span>
            </div>
            <div className="vi-cursor">1247</div>
            <div style={{ fontSize: "11px", color: "#888", padding: "0 6px" }}>0</div>
            <div style={{ fontSize: "11px", color: "#bbb", padding: "0 6px" }}>2 min ago</div>
            <div style={{ display: "flex", gap: "5px", padding: "0 6px" }}>
              <button className="vi-btn vi-btn-resend">↺ Resend</button>
            </div>
          </div>

          <div className="vi-grid vi-row">
            <div className="vi-id" style={{ paddingLeft: "12px" }}>exch_b3c9d7e2f4a8b1c6d0…</div>
            <div style={{ padding: "0 6px" }}>
              <span className="sync-badge sb-waiting">
                <span className="vcdot" style={{ background: "#D97706" }}></span>waiting_ack
              </span>
            </div>
            <div className="vi-cursor">1246</div>
            <div style={{ fontSize: "11px", color: "#888", padding: "0 6px" }}>0</div>
            <div style={{ fontSize: "11px", color: "#bbb", padding: "0 6px" }}>3 min ago</div>
            <div style={{ display: "flex", gap: "5px", padding: "0 6px" }}></div>
          </div>

          <div className="vi-grid vi-row" style={{ background: "#FFFBF5" }}>
            <div className="vi-id" style={{ paddingLeft: "12px" }}>exch_c4d8f3a1e7b2c9d5f0…</div>
            <div style={{ padding: "0 6px" }}>
              <span className="sync-badge sb-failed">
                <span className="vcdot" style={{ background: "#DC2626" }}></span>failed
              </span>
            </div>
            <div className="vi-cursor">1243</div>
            <div style={{ fontSize: "11px", color: "#DC2626", fontWeight: 600, padding: "0 6px" }}>3 / 5</div>
            <div style={{ fontSize: "11px", color: "#bbb", padding: "0 6px" }}>1 hr ago</div>
            <div style={{ display: "flex", gap: "5px", padding: "0 6px" }}>
              <button className="vi-btn vi-btn-retry">↺ Retry</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
