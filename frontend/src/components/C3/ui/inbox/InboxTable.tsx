import React from "react";


export const InboxTable = () => {
    return (
        <div className="inbox">
            {/* INBOX MAIN */}
            <div className="inbox-main">
                <div className="inbox-topbar">
                    <div>
                        <div className="inbox-title-row">
                            <div className="inbox-title">INBOX</div>
                            <div className="inbox-count-pill">5 pending</div>
                        </div>
                        <div className="inbox-subtitle">
                            Slots waiting for a commit from your vaults
                        </div>
                    </div>
                    <div className="inbox-controls">
                        <div className="ctrl-btn">All vaults ▾</div>
                        <div className="ctrl-btn">Sort: Oldest first ▾</div>
                    </div>
                </div>
                <div className="inbox-list">
                    {/* vault_finance group */}
                    <div className="vault-group">
                        <div className="vg-header">
                            <div className="vg-dot" style={{ background: "#2563EB" }} />
                            <div className="vg-name">vault_finance</div>
                            <div className="vg-count">3 pending</div>
                        </div>
                        {/* Item 1 - urgent */}
                        <div className="inbox-item urgent">
                            <div className="ii-left">
                                <div className="ii-priority pri-urgent" />
                            </div>
                            <div className="ii-main">
                                <div className="ii-slot">
                                    <span className="slot-tag">financial_clearance</span>
                                    <span className="gate-ok-badge">Gate ✓ passed</span>
                                </div>
                                <div className="ii-thread">
                                    <span className="ii-thread-type">Contract</span>
                                    <span className="ii-thread-name">Contract — Supplier X</span>
                                </div>
                                <div className="ii-since">
                                    Waiting 2 days — gate passed Jun 10
                                </div>
                            </div>
                            <div className="ii-right">
                                <div className="ii-ts">Jun 10</div>
                                <button className="commit-btn">✦ Commit</button>
                            </div>
                        </div>
                        {/* Item 2 - normal */}
                        <div className="inbox-item pending">
                            <div className="ii-left">
                                <div className="ii-priority pri-normal" />
                            </div>
                            <div className="ii-main">
                                <div className="ii-slot">
                                    <span className="slot-tag">budget_approval</span>
                                    <span className="gate-ok-badge">Gate ✓ passed</span>
                                </div>
                                <div className="ii-thread">
                                    <span className="ii-thread-type">Budget</span>
                                    <span className="ii-thread-name">Budget Allocation Q3</span>
                                </div>
                                <div className="ii-since">
                                    Waiting 1 day — gate passed Jun 11
                                </div>
                            </div>
                            <div className="ii-right">
                                <div className="ii-ts">Jun 11</div>
                                <button className="commit-btn">✦ Commit</button>
                            </div>
                        </div>
                        {/* Item 3 - blocked */}
                        <div
                            className="inbox-item"
                            style={{ paddingLeft: 25, borderLeft: "3px solid #e0e0e0" }}
                        >
                            <div className="ii-left">
                                <div className="ii-priority pri-low" />
                            </div>
                            <div className="ii-main">
                                <div className="ii-slot">
                                    <span className="slot-tag">treasury_release</span>
                                    <span className="gate-blocked-badge">Gate ✕ blocked</span>
                                </div>
                                <div className="ii-thread">
                                    <span className="ii-thread-type">Invoice</span>
                                    <span className="ii-thread-name">
                                        Invoice Processing — Cipla India
                                    </span>
                                </div>
                                <div className="ii-since">
                                    Blocked — waiting on ops_confirmation from vault_ops
                                </div>
                            </div>
                            <div className="ii-right">
                                <div className="ii-ts">Jun 09</div>
                                <div className="blocked-text">⛔ Blocked</div>
                            </div>
                        </div>
                    </div>
                    {/* vault_legal group */}
                    <div className="vault-group">
                        <div className="vg-header">
                            <div className="vg-dot" style={{ background: "#7C3AED" }} />
                            <div className="vg-name">vault_legal</div>
                            <div className="vg-count">2 pending</div>
                        </div>
                        {/* Item 4 - C3 external */}
                        <div className="inbox-item c3">
                            <div className="ii-left">
                                <div
                                    className="ii-priority"
                                    style={{ background: "#7C3AED" }}
                                />
                            </div>
                            <div className="ii-main">
                                <div className="ii-slot">
                                    <span className="slot-tag">contract_draft</span>
                                    <span className="c3-badge">⛓ C3 channel</span>
                                </div>
                                <div className="ii-thread">
                                    <span className="ii-thread-type">Contract</span>
                                    <span className="ii-thread-name">
                                        contract-execution — Cipla_India
                                    </span>
                                </div>
                                <div className="ii-since">
                                    New channel activated just now — first slot to commit
                                </div>
                            </div>
                            <div className="ii-right">
                                <div className="ii-ts">just now</div>
                                <button className="commit-btn">✦ Commit</button>
                            </div>
                        </div>
                        {/* Item 5 */}
                        <div className="inbox-item pending">
                            <div className="ii-left">
                                <div className="ii-priority pri-normal" />
                            </div>
                            <div className="ii-main">
                                <div className="ii-slot">
                                    <span className="slot-tag">nda_review</span>
                                    <span className="gate-ok-badge">Gate ✓ passed</span>
                                </div>
                                <div className="ii-thread">
                                    <span className="ii-thread-type">NDA</span>
                                    <span className="ii-thread-name">
                                        NDA — Accenture Partnership
                                    </span>
                                </div>
                                <div className="ii-since">Waiting 4 hours</div>
                            </div>
                            <div className="ii-right">
                                <div className="ii-ts">Jun 24</div>
                                <button className="commit-btn">✦ Commit</button>
                            </div>
                        </div>
                    </div>
                    {/* Other vaults — nothing pending */}
                    <div className="vault-group">
                        <div className="vg-header">
                            <div className="vg-dot" style={{ background: "#059669" }} />
                            <div className="vg-name">vault_treasury</div>
                            <div className="vg-count" style={{ color: "#ccc" }}>
                                0 pending
                            </div>
                        </div>
                        <div className="vg-empty">No actions pending for this vault</div>
                    </div>
                </div>
                <div className="inbox-footer">
                    <span>5 pending slots across 2 vaults</span>
                    <span style={{ color: "#C8922A", cursor: "pointer" }}>
                        Mark all viewed
                    </span>
                </div>
            </div>
        </div>
    )
}
