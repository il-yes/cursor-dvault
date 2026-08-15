import { Global, css } from "@emotion/react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { C3BaseStyles, C3styles, PanelStyles } from "@/components/C3/styles/styles";
import { C3MenuStyles } from "@/components/C3/styles/menu";
import { C3LedgerStyles } from "@/components/C3/styles/ledger";
import { C3SharedStyles } from "@/components/C3/styles/shared";
import { useParams } from "react-router-dom";
import { LedgerLayout } from "../LedgerLayout";
import { ChannelTable } from "./ChannelTable";
import { useEffect, useState } from "react";
import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";
import { ChannelView } from "../../domain/channel/channel.types";
import { ThreadAssetSingleDrawer } from "../thread/ThreadAssetSingleDrawer";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";
import { toChannelView } from "../../domain/channel/channel.mapper";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const ChannelPage = () => {
    const { channelId } = useParams();
    const [channel, setChannel] = useState<ChannelView | null>(null);
    const [selectedAsset, setSelectedAsset] = useState<ThreadAssetViewInterface | null>(null);
    const [activating, setActivating] = useState(false);
    const [activateError, setActivateError] = useState<string | null>(null);
    const { channels, selectChannel, activateChannel: activateStoreChannel } = useC3ChannelStore();

    useEffect(() => {
        const found = channels.find((c) => c.id === channelId);
        if (found) {
            setChannel(toChannelView(found));
            selectChannel(found.id);
        } else {
            setChannel(null);
        }
    }, [channelId, channels, selectChannel]);

    const handleActivate = async () => {
        if (!channel || activating) return;

        setActivating(true);
        setActivateError(null);
        try {
            const updated = await activateStoreChannel(channel.id);
            setChannel(toChannelView(updated));
        } catch (err: any) {
            console.error("Failed to activate channel:", err);
            setActivateError(err?.message || "Activation failed.");
        } finally {
            setActivating(false);
        }
    };

    return (
        <DashboardLayout>

            <C3SharedStyles />
            <C3MenuStyles />
            <C3LedgerStyles />
            {/* <C3styles /> */}
            <PanelStyles />


            <LedgerLayout isNewShareOpen={false}>
                <ChannelTable channel={channel} onOpenAsset={(asset) => { setSelectedAsset(asset); }} onActivate={handleActivate} activating={activating} activateError={activateError} />

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


