import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import Main from "@/components/C3/Main";
import { C3BaseStyles,  C3styles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { AssetCard, AssetView } from "@/components/C3/asset-view";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;
export const assets: AssetView[] = [

    {
        id: 'msa',
        channelId: 'contract-supplier',

        type: 'contract',

        title:
            'Master Service Agreement',

        subtitle:
            'Supplier X · $240,000',

        status: 'signed',

        createdAt:
            '12 Jun 2026',

        lastEvent: {
            type: 'contract.countersigned',
            at: '2 hours ago'
        },

        participants: [
            'Legal',
            'Supplier X'
        ],

        stellarTx:
            'tx_f2a9e3…',

        c3Visibility:
            'external'
    },


    {
        id: 'invoice-2047',

        channelId: 'supplier-payment',

        type: 'invoice',

        title:
            'Invoice INV-2047',

        subtitle:
            'Cipla India',

        status: 'approved',

        createdAt:
            '02 Jun 2026',

        lastEvent: {
            type: 'payment.approved',
            at: '3 days ago'
        },

        participants: [
            'Operations',
            'Finance'
        ],

        stellarTx:
            'tx_a3f1c2…',

        c3Visibility:
            'shared'
    }

];
const AssetPage = () => {
  return (
    <DashboardLayout>
      <C3SharedStyles />
      <C3MenuStyles />
      <C3LedgerStyles />
      <C3styles />
      <AssetCard asset={assets[0]} />
    </DashboardLayout>
  );
};

export default AssetPage;


