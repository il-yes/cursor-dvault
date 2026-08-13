import React from "react";
import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { C3BaseStyles, C3styles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { LedgerTable } from "./LedgerTable";
import { LedgerLayout } from "../LedgerLayout";
import { useNavigate, useParams } from "react-router-dom";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";
import { useC3WorkspaceStore } from "../../infrastructure/store/useC3WorkspaceStore";
import { ChannelRow } from "../../domain/channel/channel.types";
import * as ROUTES from '@/constants/routes';

export const c3BaseStyles = css`
  /* global C3 styles */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const LedgerPage = () => {
  const [isNewShareOpen, setIsNewShareOpen] = React.useState(false);
  const navigate = useNavigate();
  const { channelId } = useParams();

  const { activeWorkspace } = useC3WorkspaceStore();
  const { channels, activeChannelId, selectChannel, isLoading, error } = useC3ChannelStore();

  const onOpenChannel = (id: string) => {
    selectChannel(id);
    navigate(ROUTES.CHANNEL.replace(':channelId', id));
  };

  // Map backend channel responses to C3 ChannelRow presentation model
  const channelRows: ChannelRow[] = channels.map((c) => ({
    id: c.id,
    status: c.status === 'revoked' ? 'pending' : 'active',
    type: 'Governance',
    title: c.title,
    subtitle: `Workspace: ${activeWorkspace?.name || 'C3 Substrate'}`,
    participants: [
      { label: 'Substrate', color: '#2563EB' },
      { label: 'Vault', color: '#059669' },
    ],
    assetCount: c.asset_count || 0,
    lastEvent: c.last_event || 'channel.created',
    lastActivity: c.updated_at ? new Date(c.updated_at).toLocaleDateString() : 'Just now',
    stellarTx: 'tx_anchored…',
    c3Status: 'active',
    c3Label: '⛓✓',
  }));

  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3LedgerStyles />
      <C3styles />

      {/* Main layout */}
      <LedgerLayout isNewShareOpen={isNewShareOpen}>
        {isLoading ? (
          <div style={{
            flex: 1,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '60px 20px',
            color: '#6B7280',
            fontSize: '14px',
          }}>
            <span>Loading channels for workspace…</span>
          </div>
        ) : (
          <LedgerTable
            ledgerRow={channelRows}
            newCreated={false}
            hasConflict={false}
            onOpenChannel={onOpenChannel}
          />
        )}
      </LedgerLayout>
    </DashboardLayout>
  );
};

export default LedgerPage;
