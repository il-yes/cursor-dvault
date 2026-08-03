import {
    Drawer,
    DrawerContent,
} from "@/components/ui/drawer";

import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";
import { useEffect, useState } from "react";
import { fetchChannel } from "../../domain/channel/channel.repository";
import { Channel } from "../../domain/channel/channel.types";
import ThreadAssetPanel, { ThreadAssetView } from "./AssetTable";


export function ThreadAssetSingleDrawer({
    open,
    asset,
    onClose,
}: {
    open: boolean;
    asset: ThreadAssetViewInterface | null;
    onClose: () => void;
}) {
    const [channel, setChannel] = useState(null)

    if (!asset) return null;

    if (asset)  {
            const load = async () => {
                const result = await fetchChannel(asset?.channelId || "");
                setChannel(result);
                console.log("channel", result);
            };
    
            load();
        }

    console.log({asset})

    return (
        <Drawer
            open={open}
            onOpenChange={(v) => {
                if (!v) onClose();
            }}
            direction="right"
        >
            <DrawerContent className="thread-drawer">
                {/* <ThreadDetailSlidingPanel channel={channel} asset={asset} hasConflict={false} /> */}
                {/* <ThreadAssetPanel asset={asset} onClose={() => console.log('close')} /> */}
                <ThreadAssetView channel={channel} asset={asset} hasConflict={false} />
            </DrawerContent>

        </Drawer >
    );
}

