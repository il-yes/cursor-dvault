import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { C3BaseStyles, C3styles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { LedgerTable } from "./LedgerTable";
import { useMemo, useState } from "react";
import { LedgerLayout } from "../LedgerLayout";
import { useNavigate } from "react-router-dom";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";
import { toChannelRows } from "../../domain/channel/channel.mapper";
import * as ROUTES from '@/constants/routes';


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const LedgerPage = () => {
  const [isNewShareOpen, setIsNewShareOpen] = useState(false);
  const navigate = useNavigate();
  const { channels, isLoading } = useC3ChannelStore();

  const ledgerRow = useMemo(() => toChannelRows(channels), [channels]);

  const onOpen = (id: string) => {
    navigate(ROUTES.CHANNEL.replace(':channelId', id));
  };

  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3LedgerStyles />
      <C3styles />

      {/* Main layout */}
      <LedgerLayout isNewShareOpen={isNewShareOpen}>
        {isLoading && channels.length === 0 ? (
          <div className="ledger-area">
            <div className="ledger-topbar">
              <div className="ledger-title">LEDGER</div>
            </div>
            <div className="empty-area">
              <div className="empty-subtext">Loading channels…</div>
            </div>
          </div>
        ) : (
          <LedgerTable ledgerRow={ledgerRow} newCreated={false} hasConflict={false} onOpenChannel={onOpen} />
        )}
      </LedgerLayout>
    </DashboardLayout>

  );
};

export default LedgerPage;


