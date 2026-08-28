import React, { useState } from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";

export const SecuritySovereigntyTab: React.FC = () => {
  const { securityState } = useC3ConfigurationStore();
  const [expandedVault, setExpandedVault] = useState<string | null>("vault:supplier-x");

  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
      <div className="fed-header">
        <div className="ch-title-row">
          <div className="fed-title">FEDERATION</div>
          <span className="badge badge-active">Active</span>
        </div>
        <div className="fed-meta">
          <div className="fed-meta-item">
            <span className="ch-meta-key">channel</span>
            <span className="ch-meta-val">contract-execution</span>
          </div>
          <div className="fed-meta-item">
            <span className="ch-meta-key">protocol</span>
            <span className="ch-meta-val">AFP v1.2</span>
          </div>
          <div className="fed-meta-item">
            <span className="ch-meta-key">workspace</span>
            <span className="ch-meta-val">acme-corp</span>
          </div>
          <div className="fed-meta-item">
            <span className="ch-meta-key">device_id</span>
            <span className="ch-meta-val">{securityState.localDeviceId}</span>
          </div>
        </div>
      </div>

      <div className="vault-list-header">
        <span className="vlh-title">Remote Vaults</span>
        <span className="vlh-count">3</span>
        <div style={{ flex: 1 }}></div>
        <button className="btn-add-vault">+ Add Remote Vault</button>
      </div>

      <div className="vault-list">
        {/* Vault 1: supplier-x (trusted, expanded) */}
        <div className="vault-card trusted">
          <div
            className={`vc-header ${expandedVault === "vault:supplier-x" ? "expanded" : ""}`}
            onClick={() => setExpandedVault(expandedVault === "vault:supplier-x" ? null : "vault:supplier-x")}
          >
            <span className="vc-trust vct-trusted">✓ Trusted</span>
            <span className="vc-id">vault:supplier-x</span>
            <span className="vc-endpoint">https://afp.supplier-x.io/v1/sync</span>
            <span className="vc-last-seen">last seen 2 min ago</span>
            <span className="vc-cursor">cursor: 1247</span>
            <span className="vc-proto">AFP v1.2</span>
          </div>

          {expandedVault === "vault:supplier-x" && (
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

        {/* Vault 2: ey-auditors */}
        <div className="vault-card trusted">
          <div
            className={`vc-header ${expandedVault === "vault:ey-auditors" ? "expanded" : ""}`}
            onClick={() => setExpandedVault(expandedVault === "vault:ey-auditors" ? null : "vault:ey-auditors")}
          >
            <span className="vc-trust vct-trusted">✓ Trusted</span>
            <span className="vc-id">vault:ey-auditors</span>
            <span className="vc-endpoint">https://afp.ey.com/ankhora/v1/sync</span>
            <span className="vc-last-seen">last seen 1 hr ago</span>
            <span className="vc-cursor">cursor: 892</span>
            <span className="vc-proto">AFP v1.1</span>
          </div>
        </div>

        {/* Vault 3: partner-bank (pending trust) */}
        <div className="vault-card pending-trust">
          <div
            className={`vc-header ${expandedVault === "vault:partner-bank" ? "expanded" : ""}`}
            onClick={() => setExpandedVault(expandedVault === "vault:partner-bank" ? null : "vault:partner-bank")}
          >
            <span className="vc-trust vct-pending">⧖ Pending</span>
            <span className="vc-id">vault:partner-bank</span>
            <span className="vc-endpoint">https://afp.partnerbank.io/sync</span>
            <span className="vc-last-seen" style={{ color: "#D97706" }}>last seen 1 day ago</span>
            <span className="vc-cursor" style={{ color: "#D97706" }}>cursor: 0</span>
            <span className="vc-proto" style={{ background: "#FFF7ED", borderColor: "#FED7AA", color: "#D97706" }}>AFP v1.0</span>
          </div>
        </div>
      </div>

      <div className="dirty-bar">
        <span className="db-icon">⚠</span>
        <div className="db-msg">
          <strong>Sovereign status:</strong> Keyring health: {securityState.sovereignKeyringStatus}. Local vault address: <code>{securityState.localVaultAddress}</code>.
        </div>
        <button className="btn-save">Publish Snapshot</button>
      </div>
    </div>
  );
};

