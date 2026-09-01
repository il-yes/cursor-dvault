import React, { useState } from "react";

interface AddRemoteVaultPanelProps {
  isOpen: boolean;
  onClose: () => void;
  onAddVault: (vault: {
    id: string;
    endpoint: string;
    status: "trusted" | "pending" | "observer";
    proto: string;
    cursor: string;
  }) => void;
}

export const AddRemoteVaultPanel: React.FC<AddRemoteVaultPanelProps> = ({
  isOpen,
  onClose,
  onAddVault,
}) => {
  const [vaultId, setVaultId] = useState("vault:central-bank-gh");
  const [endpoint, setEndpoint] = useState("https://afp.centralbank.gov.gh/sync");
  const [trustState, setTrustState] = useState<"trusted" | "pending" | "observer">("trusted");
  const [protocolVersion, setProtocolVersion] = useState("AFP v1.2 (recommended)");
  const [startCursor, setStartCursor] = useState("0");

  if (!isOpen) return null;

  const handleSubmit = () => {
    if (!vaultId.trim() || !endpoint.trim()) return;

    onAddVault({
      id: vaultId.trim(),
      endpoint: endpoint.trim(),
      status: trustState,
      proto: protocolVersion.split(" ")[0],
      cursor: startCursor.trim() || "0",
    });

    onClose();
  };

  return (
    <>
      <div className="panel-scrim" onClick={onClose}></div>
      <div className="vault-panel">
        <div className="vp-header">
          <div className="vp-hdr-top">
            <div className="vp-title">Add Remote Vault</div>
            <div className="vp-close" onClick={onClose}>✕</div>
          </div>
          <div className="vp-subtitle">
            Register an external vault to participate in this federation snapshot via AFP.
          </div>
        </div>

        <div className="vp-body">
          {/* Vault ID */}
          <div className="vp-field">
            <div className="vp-label">
              Vault ID <span className="req">*</span>
            </div>
            <div className="vp-input-wrap">
              <input
                className="vp-input mono valid"
                type="text"
                value={vaultId}
                onChange={(e) => setVaultId(e.target.value)}
              />
              <span className="valid-icon">✓</span>
            </div>
          </div>

          {/* Endpoint */}
          <div className="vp-field">
            <div className="vp-label">
              AFP Endpoint <span className="req">*</span>
            </div>
            <div className="vp-input-wrap">
              <input
                className="vp-input valid"
                type="text"
                value={endpoint}
                onChange={(e) => setEndpoint(e.target.value)}
              />
              <span className="valid-icon">✓</span>
            </div>
          </div>

          {/* Trust state */}
          <div className="vp-field">
            <div className="vp-label">Initial Trust State</div>
            <select
              className="vp-select"
              value={trustState}
              onChange={(e) => setTrustState(e.target.value as any)}
            >
              <option value="trusted">Trusted (bidirectional sync)</option>
              <option value="pending">Pending Verification</option>
              <option value="observer">Observer Only (read-only)</option>
            </select>
          </div>

          {/* Protocol version */}
          <div className="vp-field">
            <div className="vp-label">Protocol Negotiated</div>
            <select
              className="vp-select"
              value={protocolVersion}
              onChange={(e) => setProtocolVersion(e.target.value)}
            >
              <option value="AFP v1.2 (recommended)">AFP v1.2 (recommended)</option>
              <option value="AFP v1.1">AFP v1.1</option>
              <option value="AFP v1.0">AFP v1.0</option>
            </select>
          </div>

          {/* Initial cursor */}
          <div className="vp-field">
            <div className="vp-label">Start Sync Cursor</div>
            <input
              className="vp-input mono"
              type="text"
              value={startCursor}
              onChange={(e) => setStartCursor(e.target.value)}
            />
          </div>
        </div>

        <div className="vp-footer">
          <button className="btn-cancel-add" onClick={onClose}>
            Cancel
          </button>
          <button className="btn-add" onClick={handleSubmit}>
            + Add Remote Vault
          </button>
        </div>
      </div>
    </>
  );
};
