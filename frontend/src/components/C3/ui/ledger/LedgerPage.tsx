import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import Main from "@/components/C3/Main";
import { C3BaseStyles, C3styles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { LedgerTable } from "./LedgerTable";
import { useEffect, useState } from "react";
import { LedgerLayout } from "../LedgerLayout";
import { useNavigate, useParams } from "react-router-dom";
import { Dialog, DialogContent, DialogDescription, DialogOverlay, DialogPortal, DialogTitle } from "@radix-ui/react-dialog";
import { DialogHeader } from "@/components/ui/dialog";
import { useC3DialogStore } from "../../infrastructure/store/c3DialogStore";
import { ChannelRow } from "../../domain/channel/channel.types";
import { fetchChannelRows } from "../../domain/channel/channel.repository";
import * as ROUTES from '@/constants/routes';
// import '../styles/ledger-style.css'


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const LedgerPage = () => {
  const [isNewShareOpen, setIsNewShareOpen] = useState(false);
  const navigate = useNavigate();
  const { channelId } = useParams()
  const [open, setOpen] = useState(false)
  const [channels, setChannels] = useState<ChannelRow[]>([])

  console.log({ channelId })

  const onOpen = (id: string) => {
    navigate(ROUTES.CHANNEL.replace(':channelId', id));
  };

  const fetchChannels = async () => {
    const data = await fetchChannelRows()
    setChannels(data)
  }

  useEffect(() => {
    // fetchChannels()
  }, [])


  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3LedgerStyles />
      <C3styles />

      {/* Main layout */}
      <LedgerLayout isNewShareOpen={isNewShareOpen}>
        <LedgerTable ledgerRow={channels} newCreated={false} hasConflict={false} onOpenChannel={onOpen} />
      </LedgerLayout>
    </DashboardLayout>

  );
};

export default LedgerPage;


