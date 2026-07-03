import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { C3BaseStyles, C3styles, PanelStyles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { useNavigate, useParams } from "react-router-dom";
import { LedgerLayout } from "../LedgerLayout";
import { ChannelTable } from "./ChannelTable";
import { fetchChannel } from "../../domain/channel/channel.repository";
import { useEffect, useState } from "react";
import { ThreadAssetDrawer } from "../thread/ThreadAssetDrawer";
import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";
import { fetchThreadAsset } from "../../domain/thread/asset.repository";
import { NewThreadAssetDrawer } from "../thread/NewThreadAssetDrawer";
import { ChannelView } from "../../domain/channel/channel.types";
import { ThreadAssetSingleDrawer } from "../thread/ThreadAssetSingleDrawer";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const ChannelPage = () => {
    const { channelId } = useParams();
    const [channel, setChannel] = useState<ChannelView | null>(null);
    const [selectedAsset, setSelectedAsset] = useState<ThreadAssetViewInterface | null>(null);


    useEffect(() => {
        const load = async () => {
            const result = await fetchChannel(channelId || "");
            setChannel(result);
        };

        load();
    }, [channelId]);

    return (
        <DashboardLayout>

            <C3SharedStyles />
            <C3MenuStyles />
            <C3LedgerStyles />
            {/* <C3styles /> */}
            <PanelStyles />


            <LedgerLayout isNewShareOpen={false}>
                <ChannelTable channel={channel} onOpenAsset={(asset) => { setSelectedAsset(asset); }} />

                {/* <ThreadAssetDrawer

                    open={selectedAsset !== null}

                    asset={selectedAsset}

                    onClose={() => setSelectedAsset(null)}

                /> */}
                
                {/* Thread asset detail */}
                <ThreadAssetSingleDrawer
                    open={selectedAsset !== null}
                    asset={selectedAsset}
                    onClose={() => setSelectedAsset(null)}
                />
            </LedgerLayout>
        </DashboardLayout>
    );
};

export default ChannelPage;


