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
import { ThreadAssetHeader } from "./AssetHeader";
import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";
import { C3SharedStyles } from "../../styles/shared";
import { C3MenuStyles } from "../../styles/menu";
import { C3LedgerStyles } from "../../styles/ledger";


export function ThreadAssetDrawer_ALPHA({
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
        <>
            <C3SharedStyles />
            <C3MenuStyles />
            <C3LedgerStyles />

            <Drawer
                open={open}
                onOpenChange={(value) => {
                    if (!value) onClose();
                }}
                direction="right"
            >

                <DrawerContent
                    className="
                    h-screen
                    w-[520px]
                    ml-auto
                    rounded-none
                    bg-background
                "
                >

                    <DrawerHeader>

                        <DrawerTitle>

                            {asset?.title}

                        </DrawerTitle>

                    </DrawerHeader>


                    <div className="overflow-y-auto px-6 pb-8">


                        <ThreadAssetHeader
                            asset={asset}
                        />


                        <LifecycleSection
                            asset={asset}
                        />


                        <EventTimelineSection
                            asset={asset}
                        />


                        <ReceiptSection
                            asset={asset}
                        />


                        <PolicySection
                            asset={asset}
                        />


                        <AnchorSection
                            asset={asset}
                        />


                        <RelationshipSection
                            asset={asset}
                        />


                    </div>


                </DrawerContent>

            </Drawer>
        </>
    );
}

export function ThreadAssetDrawer({
    open,
    asset,
    onClose,
}: {
    open: boolean;
    asset: ThreadAssetViewInterface | null;
    onClose: () => void;
}) {

    if (!asset) return null;

    console.log({asset})


    return (
        <Drawer
            open={open}
            onOpenChange={(v) => {
                if (!v) onClose();
            }}
            direction="right"
        >

            <DrawerContent
                className="drawer-reset"
            >


                <DrawerTitle className="sr-only">
                    {asset.title}
                </DrawerTitle>
                <div className="slide-panel">

                    <div className="sp-header">

                        <div className="sp-header-row">

                            <div>

                                <div className="sp-title">
                                    {asset?.title}
                                </div>

                                <div className="sp-subtitle">
                                    {asset?.subtitle}
                                </div>

                            </div>


                            <DrawerClose className="sp-close">
                                ✕
                            </DrawerClose>


                        </div>

                    </div>



                    <div className="sp-body">


                        <div className="fl">
                            Lifecycle
                        </div>


                        <div className="channel-flow-box">


                            {
                                asset?.events?.map(event => (

                                    <div
                                        key={event?.id}
                                        className="pipeline-step"
                                    >

                                        <div className="ps-icon">

                                            {event?.actor
                                                .slice(0, 3)
                                                .toUpperCase()
                                            }

                                        </div>


                                        <div className="ps-content">


                                            <div className="ps-label">

                                                {event?.type}

                                            </div>


                                            <div className="ps-sublabel">

                                                {event?.actor}

                                            </div>


                                            <div className="ps-ts">

                                                {event?.time}

                                            </div>


                                            <div>

                                                CID:
                                                {" "}
                                                {event?.payloadCID}

                                            </div>


                                        </div>



                                        <div className="ps-check done">

                                            ✓

                                        </div>


                                    </div>


                                ))
                            }


                        </div>



                        <div className="stellar-info">

                            <span className="si-icon">
                                ✦
                            </span>


                            <span className="si-text">

                                Stellar anchored thread asset

                            </span>


                            <div className="si-status">

                                <div className="si-dot" />

                                Active

                            </div>


                        </div>


                    </div>



                    <div className="sp-footer">


                        <button className="start-btn">

                            ↓ Export thread

                        </button>


                        <div className="footer-note">

                            {asset?.events?.length}
                            {" "}
                            events ·
                            {" "}
                            {asset?.status}

                        </div>


                    </div>



                </div>


            </DrawerContent>


        </Drawer>
    )
}