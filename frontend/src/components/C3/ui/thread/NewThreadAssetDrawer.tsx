import {
    Drawer,
    DrawerContent,
} from "@/components/ui/drawer";
// import '../styles/asset-thread.css';



export function NewThreadAssetDrawer({
    open,
    onClose,
}: {
    open: boolean;
    onClose: () => void;
}) {

    return (
        <Drawer
            open={open}
            onOpenChange={(v) => {
                if (!v) onClose();
            }}
            direction="right"
        >
            <DrawerContent className="thread-drawer">
                <AddThreadSlidingPanel />
            </DrawerContent>

        </Drawer >
    );
}
export const AddThreadSlidingPanel = () => {
    return (
            <div className="slide-panel">
                <div className="sp-header">
                    <div className="sp-header-row">
                        <div>
                            <div className="sp-title">New Thread</div>
                            <div className="sp-subtitle">
                                Instantiate a channel into a new thread
                            </div>
                        </div>
                        <div className="sp-close">✕</div>
                    </div>
                </div>
                <div className="sp-body">
                    {/* 1. Channel selector */}
                    <div>
                        <div className="fl">
                            Channel{" "}
                            <span className="fl-hint">defines slots, gates, vault flow</span>
                        </div>
                        <div className="channel-selected">
                            <div className="cs-icon-wrap">📄</div>
                            <div className="cs-content">
                                <div className="cs-name">Contract Execution</div>
                                <div className="cs-desc">
                                    contract-execution · 3 slots · 2 gated
                                </div>
                            </div>
                            <span className="cs-arrow">▾</span>
                        </div>
                    </div>
                    {/* 2. Channel flow preview */}
                    <div>
                        <div className="fl">Flow</div>
                        <div className="channel-flow-box">
                            <div className="cfb-row">
                                <div className="cfb-vault">
                                    <div className="cfb-dot" style={{ background: "#7C3AED" }} />
                                    vault_legal
                                </div>
                                <div className="cfb-arrow">→</div>
                                <div className="cfb-vault">
                                    <div className="cfb-dot" style={{ background: "#2563EB" }} />
                                    vault_finance
                                </div>
                                <div className="cfb-arrow">→</div>
                                <div className="cfb-vault">
                                    <div className="cfb-dot" style={{ background: "#444" }} />
                                    vault_direction
                                </div>
                            </div>
                            <div className="cfb-meta">
                                <div className="cfb-metaitem">
                                    <strong>3</strong> slots
                                </div>
                                <div className="cfb-metaitem">
                                    <strong>2</strong> gated
                                </div>
                                <div className="cfb-metaitem">
                                    first slot: <strong>contract_draft</strong>
                                </div>
                            </div>
                        </div>
                    </div>
                    {/* 3. Thread name */}
                    <div>
                        <div className="fl">Thread name</div>
                        <div className="thread-name-wrap">
                            <div className="thread-name-prefix">contract-execution —</div>
                            <input
                                className="thread-name-input"
                                type="text"
                                defaultValue="Accenture Partnership"
                                placeholder="descriptor"
                            />
                        </div>
                    </div>
                    {/* 4. Properties */}
                    <div>
                        <div className="fl">
                            Properties{" "}
                            <span className="fl-hint">
                                key : value pairs attached to this thread
                            </span>
                        </div>
                        <div className="props-grid">
                            <div className="prop-row">
                                <div className="prop-key-wrap">
                                    <div className="prop-label">Key</div>
                                    <input
                                        className="prop-input prefilled"
                                        type="text"
                                        defaultValue="counterparty"
                                        readOnly
                                    />
                                </div>
                                <div className="prop-val-wrap">
                                    <div className="prop-label">Value</div>
                                    <input
                                        className="prop-input"
                                        type="text"
                                        defaultValue="Accenture"
                                        placeholder="value"
                                    />
                                </div>
                                <div className="prop-remove">✕</div>
                            </div>
                            <div className="prop-row">
                                <div className="prop-key-wrap">
                                    <div className="prop-label">Key</div>
                                    <input
                                        className="prop-input"
                                        type="text"
                                        placeholder="e.g. value"
                                    />
                                </div>
                                <div className="prop-val-wrap">
                                    <div className="prop-label">Value</div>
                                    <input
                                        className="prop-input"
                                        type="text"
                                        placeholder="e.g. $180,000"
                                    />
                                </div>
                                <div className="prop-remove">✕</div>
                            </div>
                        </div>
                        <div className="add-prop-btn" style={{ marginTop: 8 }}>
                            + Add property
                        </div>
                    </div>
                    {/* 5. Vault overrides (optional) */}
                    <div>
                        <div className="fl">
                            Vault assignment{" "}
                            <span className="fl-hint">
                                pre-set by channel — override if needed
                            </span>
                        </div>
                        <div className="vault-overrides">
                            <div className="vault-override-row">
                                <div className="vor-role">Author</div>
                                <div className="vor-select">
                                    <div className="vor-dot" style={{ background: "#7C3AED" }} />
                                    <div className="vor-name">vault_legal</div>
                                    <div className="vor-arrow">▾</div>
                                </div>
                            </div>
                            <div className="vault-override-row">
                                <div className="vor-role">Reviewer</div>
                                <div className="vor-select">
                                    <div className="vor-dot" style={{ background: "#2563EB" }} />
                                    <div className="vor-name">vault_finance</div>
                                    <div className="vor-arrow">▾</div>
                                </div>
                            </div>
                            <div className="vault-override-row">
                                <div className="vor-role">Approver</div>
                                <div className="vor-select">
                                    <div className="vor-dot" style={{ background: "#444" }} />
                                    <div className="vor-name">vault_direction</div>
                                    <div className="vor-arrow">▾</div>
                                </div>
                            </div>
                        </div>
                    </div>
                    {/* 6. Stellar anchor info */}
                    <div className="stellar-info">
                        <span className="si-icon">✦</span>
                        <span className="si-text">
                            A <strong>genesis transaction</strong> will be anchored on Stellar
                            the moment this thread starts. Every subsequent commit is
                            automatically anchored.
                        </span>
                        <div className="si-status">
                            <div className="si-dot" />
                            <span className="si-label">Active</span>
                        </div>
                    </div>
                </div>
                {/* Footer */}
                <div className="sp-footer">
                    <button className="start-btn">▶ Start Thread</button>
                    <div className="footer-note">
                        Thread appears in the ledger immediately.{" "}
                        <strong>vault_legal</strong> can commit{" "}
                        <strong>contract_draft</strong> right away. C3 extension can be
                        added at any time.
                    </div>
                </div>
            </div>
    )
}