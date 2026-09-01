import React, { useState } from "react";
import { RemoteVaultCard, RemoteVaultItem } from "./RemoteVaultCard";

interface RemoteVaultListProps {
  vaults: RemoteVaultItem[];
  onOpenAddPanel: () => void;
}

export const RemoteVaultList: React.FC<RemoteVaultListProps> = ({
  vaults,
  onOpenAddPanel,
}) => {
  const [expandedVaultId, setExpandedVaultId] = useState<string | null>("vault:supplier-x");

  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
      <div className="vault-list-header">
        <span className="vlh-title">Remote Vaults</span>
        <span className="vlh-count">{vaults.length}</span>
        <div style={{ flex: 1 }}></div>
        <button className="btn-add-vault" onClick={onOpenAddPanel}>
          + Add Remote Vault
        </button>
      </div>

      <div className="vault-list">
        {vaults.map((vault) => (
          <RemoteVaultCard
            key={vault.id}
            vault={vault}
            isExpanded={expandedVaultId === vault.id}
            onToggleExpand={() =>
              setExpandedVaultId(expandedVaultId === vault.id ? null : vault.id)
            }
          />
        ))}
      </div>
    </div>
  );
};
