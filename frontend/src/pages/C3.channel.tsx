import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { channel, ChannelView } from "@/components/C3/channel-view";
import { C3BaseStyles,  C3styles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { useParams } from "react-router-dom";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const ChannelPage = () => {
    const { channelId } = useParams();
    // const channel = useChannelStore((state) => state.channels[channelId || ""]);
    console.log({channelId})
  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3LedgerStyles />
      <C3styles />
      <ChannelView channel={channel} />
    </DashboardLayout>
  );
};

export default ChannelPage;


