import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import Main from "@/components/C3/Main";
import { C3BaseStyles,  C3styles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const LedgerPage = () => {
  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3LedgerStyles />
      {/* <C3styles /> */}
      <Main />
    </DashboardLayout>
  );
};

export default LedgerPage;


