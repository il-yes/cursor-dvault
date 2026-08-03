import React, { useState } from "react";
import { NewThreadAssetDrawer } from "./thread/NewThreadAssetDrawer";
import { TopToolbar } from "./top_toolbar";


export const LedgerLayout = ({
    children,
    isNewShareOpen
}: {
    children: React.ReactNode;
    isNewShareOpen: boolean;
}) => {
    const [openNewThread, setOpenNewThread] = useState(false)

    return (
        <>
            <TopToolbar setOpenNewThread={setOpenNewThread} />

            {/* Activation toast */}
            {isNewShareOpen && <div className="toast-bar">
                <span className="toast-icon">⚡</span>
                <span className="toast-text">
                    Channel <strong>contract-execution</strong> activated —{" "}
                    <strong>vault_legal</strong> can now commit{" "}
                    <strong>contract_draft</strong>. Stellar genesis tx anchored.
                </span>
                <span className="toast-action">Go to inbox →</span>
                <span className="toast-close">✕</span>
            </div>}

            {/* Main layout */}
            <div className="layout">
                {children}
                {/* New Thread Drawer */}
                <NewThreadAssetDrawer
                    open={openNewThread}
                    onClose={() => setOpenNewThread(false)}
                />
            </div>           
        </>
    )
}