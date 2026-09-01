import React, { useEffect, useState } from "react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { RemoteVaultList } from "./RemoteVaultList";
import { RemoteVaultItem } from "./RemoteVaultCard";
import { AddRemoteVaultPanel } from "./AddRemoteVaultPanel";
import { useC3WorkspaceStore } from "../infrastructure/store/useC3WorkspaceStore";
import { useVaultStore } from "@/store/vaultStore";
import { getWorkspaceFederation, addRemoteVaultToWorkspace } from "@/services/api";
import "../configuration/styles/c3-template-theme.css";

export const FederationPage: React.FC = () => {
  const { activeWorkspace, activeWorkspaceId } = useC3WorkspaceStore();
  const vault = useVaultStore((state) => state.vault);
  const [isAddPanelOpen, setIsAddPanelOpen] = useState(false);

  const workspaceId = activeWorkspaceId || activeWorkspace?.id || "ws_default";
  const workspaceTitle = activeWorkspace?.title || activeWorkspace?.id || "acme-corp";
  const localVaultAddress = vault?.address || vault?.owner || (vault ? "local_vault_active" : null);

  const [remoteVaults, setRemoteVaults] = useState<RemoteVaultItem[]>([]);

  useEffect(() => {
    let isMounted = true;
    async function loadSnapshot() {
      try {
        const snap = await getWorkspaceFederation(workspaceId);
        if (isMounted && snap?.remote_vaults) {
          setRemoteVaults(snap.remote_vaults);
        }
      } catch (err) {
        console.error("Failed to load workspace federation snapshot:", err);
      }
    }
    loadSnapshot();
    return () => {
      isMounted = false;
    };
  }, [workspaceId]);

  const handleAddVault = async (newVaultData: {
    id: string;
    endpoint: string;
    status: "trusted" | "pending" | "observer";
    proto: string;
    cursor: string;
  }) => {
    const newVault: RemoteVaultItem = {
      id: newVaultData.id,
      endpoint: newVaultData.endpoint,
      status: newVaultData.status,
      lastSeen: "just now",
      cursor: newVaultData.cursor,
      proto: newVaultData.proto,
      pendingItems: newVaultData.status === "pending" ? "1 pending" : "0 pending",
      alert: false,
    };

    try {
      const snap = await addRemoteVaultToWorkspace(workspaceId, newVault);
      if (snap?.remote_vaults) {
        setRemoteVaults(snap.remote_vaults);
      } else {
        setRemoteVaults((prev) => [...prev, newVault]);
      }
    } catch (err) {
      console.error("Failed to add remote vault to workspace backend:", err);
      setRemoteVaults((prev) => [...prev, newVault]);
    }
  };

  return (
    <DashboardLayout>
      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          background: "#fff",
          minHeight: "100%",
          overflow: "hidden",
          position: "relative",
        }}
      >
        {/* Federation Header */}
        <div className="fed-header">
          <div className="ch-title-row">
            <div className="fed-title">WORKSPACE & LEDGER FEDERATION</div>
            <span className="badge badge-active">Active</span>
          </div>
          <div className="fed-meta">
            <div className="fed-meta-item">
              <span className="ch-meta-key">scope</span>
              <span className="ch-meta-val">Workspace / Ledger</span>
            </div>
            <div className="fed-meta-item">
              <span className="ch-meta-key">workspace</span>
              <span className="ch-meta-val">{workspaceTitle}</span>
            </div>
            <div className="fed-meta-item">
              <span className="ch-meta-key">protocol</span>
              <span className="ch-meta-val">AFP v1.2</span>
            </div>
            <div className="fed-meta-item">
              <span className="ch-meta-key">vault_status</span>
              <span className="ch-meta-val">{localVaultAddress ? "Active" : "Uninitialized"}</span>
            </div>
          </div>
        </div>

        {/* Remote Vault List & Exchanges */}
        <RemoteVaultList
          vaults={remoteVaults}
          onOpenAddPanel={() => setIsAddPanelOpen(true)}
        />

        {/* Sovereign Footer Bar */}
        <div className="dirty-bar">
          <span className="db-icon">⚠</span>
          <div className="db-msg">
            <strong>Sovereign status:</strong> Keyring health:{" "}
            {localVaultAddress ? "Operational" : "Vault uninitialized"}. Local vault address:{" "}
            <code>{localVaultAddress || "None"}</code>.
          </div>
          <button className="btn-save">Publish Snapshot</button>
        </div>

        {/* Reusable Add Remote Vault Panel */}
        <AddRemoteVaultPanel
          isOpen={isAddPanelOpen}
          onClose={() => setIsAddPanelOpen(false)}
          onAddVault={handleAddVault}
        />
      </div>
    </DashboardLayout>
  );
};

export default FederationPage;
