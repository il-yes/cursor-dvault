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
import { useC3ThreadStore } from "../../infrastructure/store/useC3ThreadStore";
import { toChannelView } from "../../domain/channel/channel.mapper";
import { RevokeChannelConfirmModal } from "./RevokeChannelConfirmModal";
import { ParticipantsPanel } from "./ParticipantsPanel";
import { InvitationsPanel } from "./InvitationsPanel";


export const c3BaseStyles = css`
  /* your global C3 styles here */
`;

export const C3GlobalStyles = () => <Global styles={C3BaseStyles} />;

const ChannelPage = () => {
    const { channelId } = useParams();
    const [selectedAsset, setSelectedAsset] = useState<ThreadAssetViewInterface | null>(null);
    const [activating, setActivating] = useState(false);
    const [activateError, setActivateError] = useState<string | null>(null);
    const [revokeConfirmOpen, setRevokeConfirmOpen] = useState(false);
    const [revoking, setRevoking] = useState(false);
    const [revokeError, setRevokeError] = useState<string | null>(null);
    const { activeChannel, activeChannelId, selectChannel, activateChannel: activateStoreChannel, revokeChannel: revokeStoreChannel } = useC3ChannelStore();
    const threads = useC3ThreadStore((state) => state.threads);

    useEffect(() => {
        if (channelId && activeChannelId !== channelId) {
            selectChannel(channelId);
        }
    }, [channelId, activeChannelId, selectChannel]);

    // Pure derived view — zero local state duplication, zero synchronization lag
    const channel = activeChannel ? toChannelView(activeChannel, threads) : null;

    const handleActivate = async () => {
        if (!channel || activating) return;

        setActivating(true);
        setActivateError(null);
        try {
            await activateStoreChannel(channel.id);
        } catch (err: any) {
            console.error("Failed to activate channel:", err);
            setActivateError(err?.message || "Activation failed.");
        } finally {
            setActivating(false);
        }
    };

    const handleRevoke = async () => {
        if (!channel || revoking) return;

        setRevoking(true);
        setRevokeError(null);
        try {
            // The Cloud response carries no Channel data; the store refreshes
            // the workspace channel list and this component re-derives the
            // channel view (revoked) from the updated store.
            await revokeStoreChannel(channel.id);
            setRevokeConfirmOpen(false);
        } catch (err: any) {
            console.error("Failed to revoke channel:", err);
            setRevokeError(err?.message || "Revocation failed.");
        } finally {
            setRevoking(false);
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
                <ChannelTable channel={channel} onOpenAsset={(asset) => { setSelectedAsset(asset); }} onActivate={handleActivate} activating={activating} activateError={activateError} onRevoke={() => setRevokeConfirmOpen(true)} revoking={revoking} revokeError={revokeError}>

                    {channelId && <ParticipantsPanel channelId={channelId} />}

                    {channelId && <InvitationsPanel channelId={channelId} />}

                </ChannelTable>

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

            <RevokeChannelConfirmModal
                isOpen={revokeConfirmOpen}
                channelTitle={channel?.title || ""}
                isRevoking={revoking}
                onCancel={() => setRevokeConfirmOpen(false)}
                onConfirm={handleRevoke}
            />
        </DashboardLayout>
    );
};

export default ChannelPage;


