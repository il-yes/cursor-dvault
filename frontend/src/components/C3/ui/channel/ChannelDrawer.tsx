import {
    Drawer,
    DrawerClose,
    DrawerContent,
    DrawerHeader,
    DrawerTitle,
} from "@/components/ui/drawer";

import { LifecycleSection } from "../../features/LifecycleSection";
import { EventTimelineSection } from "../../features/EventTimelineSection";
import { ReceiptSection } from "../../features/ReceiptSection";
import { PolicySection } from "../../features/PolicySection";
import { AnchorSection } from "../../features/AnchorSection";
import { RelationshipSection } from "../../features/RelationshipSection";
import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";



export function ChannelDrawer({
    open,
    asset,
    onClose,
}: {
    open: boolean;
    asset: ThreadAssetViewInterface | null;
    onClose: () => void;
}) {

    if (!asset) return null;


    return (
        <Drawer
            open={open}
            onOpenChange={(v) => {
                if (!v) onClose();
            }}
            direction="right"
        >

            <DrawerContent className="thread-drawer">

                
                

            </DrawerContent>

        </Drawer >
    );
}