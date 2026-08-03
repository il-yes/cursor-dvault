import { DashboardLayout } from "@/components/DashboardLayout";
import { C3BaseStyles, C3styles } from "@/components/C3/styles/styles";
import { Global, css } from "@emotion/react";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3InboxStyles } from "@/components/C3/styles/inbox";
import { LedgerLayout } from "../LedgerLayout";
import { InboxTable } from "./InboxTable";


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
      {/* <C3styles /> */}

      <LedgerLayout isNewShareOpen={false}>
        <InboxTable />
      </LedgerLayout>

    </DashboardLayout>
  );
};

export default InboxPage;
