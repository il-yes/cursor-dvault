import { DashboardLayout } from "@/components/DashboardLayout";
import { Inbox } from "@/components/C3/inbox";
import { C3BaseStyles, C3styles } from "@/components/C3/styles/styles";
import { Global, css } from "@emotion/react";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3InboxStyles } from "@/components/C3/styles/inbox";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const InboxPage = () => {

  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3InboxStyles />
      <C3styles />
      <Inbox />
    </DashboardLayout>
  );
};

export default InboxPage;
