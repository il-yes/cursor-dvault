import React, { useState, useEffect } from "react";
import { ThreadAssetViewInterface } from "./domain/thread/asset.types";
import { Tabs } from "@radix-ui/react-tabs";
import { Channel } from "diagnostics_channel";
// import '../C3/ui/styles/event.thread-style.css'


export default function ThreadDetailSlidingPanel({ channel, asset, hasConflict }: { channel: Channel, asset: ThreadAssetViewInterface, hasConflict: boolean }) {
    const hasC3Extension = true;
    

    return (
        <div className="detail-panel">
            {hasConflict && <DisputeBanner />}
            <div className="dp-header">
                <div className="dp-close-row">
                    <span className="dp-badge">Invoice</span>
                    <div className="dp-close">✕</div>
                </div>
                <div className="dp-title">Invoice Processing — Cipla India</div> {/* channel.ID   */}
                <div className="dp-meta">
                    <span className="dp-meta-chip">t_sha256_INV_cipla</span> {/* asset.title */}
                    <span className="dp-meta-chip">created: 2026-06-01</span> {/* asset.createdAt */}
                    <span className="dp-meta-chip">by: vault_ops</span> {/* asset.actor */}
                </div>
                <div className="dp-custom">
                    <span style={{ color: "#bbb" }}>vendor:</span>{" "} {/* asset.participantB.type */}
                    <span>Cipla India</span> {/* asset.participantB */}
                    &nbsp;·&nbsp;
                    <span style={{ color: "#bbb" }}>ref:</span> <span>INV-2047</span>
                    &nbsp;·&nbsp;
                    <span style={{ color: "#bbb" }}>amount:</span> <span>$18,400</span>
                </div>
            </div>
            <div className="dp-body">
                {/* {asset.events.map} */}
                {/* PIPELINE /Lyfecycle */}
                <div className="dp-section">
                    <div className="dp-section-title">Pipeline</div>
                    <div className="pipeline-step">
                        <div className="ps-icon" style={{ background: "#DC2626" }}>
                            O
                        </div>
                        <div className="ps-content">
                            <div className="ps-label">invoice_submitted</div> {/* event.type */}
                            <div className="ps-sublabel">vault_ops</div> {/* event.actor */}
                            <div className="ps-ts">Jun 01 · 09:14</div> {/* event.time */}
                        </div>
                        <div className="ps-check done">✓</div> {/* event.receiptStatus */}
                    </div>
                    <div className="ps-connector" />
                    <div className="pipeline-step">
                        <div className="ps-icon" style={{ background: "#2563EB" }}>
                            F
                        </div>
                        <div className="ps-content">
                            <div className="ps-label">finance_approved</div>
                            <div className="ps-sublabel">vault_finance</div>
                            <div className="ps-ts">Jun 02 · 11:30</div>
                        </div>
                        <div className="ps-check done">✓</div>
                    </div>
                    <div className="ps-connector" />
                    <div className="pipeline-step">
                        <div className="ps-icon" style={{ background: "#059669" }}>
                            T
                        </div>
                        <div className="ps-content">
                            <div className="ps-label">payment_released</div>
                            <div className="ps-sublabel">vault_treasury</div>
                            <div className="ps-ts">Jun 03 · 08:45</div>
                        </div>
                        <div className="ps-check done">✓</div>
                    </div>

                    <div className="ps-connector" />
                    <div className="pipeline-step">
                        <div className="ps-icon ps-icon-wait">Dir</div>
                        <div className="ps-content">
                            <div className="ps-label wait">executive_signature</div>
                            <div className="ps-sublabel">vault_direction · waiting</div>
                        </div>
                        <div className="ps-check wait">○</div>
                    </div>

                    {hasC3Extension && <C3PipelineCommit />}
                    {hasConflict && <RejectedStep />}

                </div>
                {/* COMMITS */}
                <div className="dp-section">
                    <div className="dp-section-title">Commits</div>
                    <div className="commit-row">
                        <div className="commit-ts">Jun 01 · 09:14</div>
                        <div className="commit-body">
                            <div>
                                <span className="commit-actor">vault_ops</span>
                                <span className="commit-action">invoice_submitted</span>
                            </div>
                            <div className="commit-cid">CID: bafyreibcx3…8f2a</div>
                        </div>
                        <div className="commit-verify">verify ↗</div>
                    </div>
                    <div className="commit-row">
                        <div className="commit-ts">Jun 02 · 11:30</div>
                        <div className="commit-body">
                            <div>
                                <span className="commit-actor">vault_finance</span>
                                <span className="commit-action">finance_approved</span>
                            </div>
                            <div className="commit-cid">CID: bafyreihzm7…3d9c</div>
                        </div>
                        <div className="commit-verify">verify ↗</div>
                    </div>
                    <div className="commit-row">
                        <div className="commit-ts">Jun 03 · 08:45</div>
                        <div className="commit-body">
                            <div>
                                <span className="commit-actor">vault_treasury</span>
                                <span className="commit-action">payment_released</span>
                            </div>
                            <div className="commit-cid">CID: bafyreia4k9…7f1e</div>
                        </div>
                        <div className="commit-verify">verify ↗</div>
                    </div>
                    {hasConflict && <RejectedPipelineCommit />}
                </div>
                {/* STELLAR */}
                <div className="dp-section">
                    <div className="dp-section-title">Stellar Reference</div>
                    <div className="stellar-ref-row">
                        <div className="stellar-hash-full">
                            tx_a3f1c2b9e8d4f7a2c5b1e9d3f6a8c2b4e7f1d9a3c8b5e2f7a4d1c9b8e3f6a2d5
                        </div>
                        <div className="copy-btn">Copy ⧉</div>
                    </div>
                </div>
                {/* C3 EXTENSION */}
                {hasConflict ? <ResolutionPath /> :<div className="dp-section">
                    <div className="dp-section-title">C3 Extension</div>
                    {!hasC3Extension && <EmptyC3Extension />}
                    {hasC3Extension && <WithC3Extension />}
                </div>}
            </div>
            {/* ACTIONS */}
            {hasConflict ? <DisputeActionButtons />  : <div className="dp-actions">
                <button className="action-btn">↓ Export thread</button>
                <button className="action-btn">🔗 Share read access</button>
                <button className="action-btn danger">✕ Close thread</button>
            </div>}
        </div>
    );
}

const EmptyC3Extension = () => (
    <div className="c3-section">
        <div className="c3-internal-label">
            This thread is internal. Add an external vault to extend it —
            giving the counterparty cryptographically-verifiable access to
            this pipeline.
        </div>
        <button className="c3-add-btn">
            ⛓ &nbsp; Add external vault
        </button>
    </div>
);

const WithC3Extension = () => (
    <div className="c3-extended-box">
        <div className="c3-ext-header">
            <span className="c3-ext-icon">⛓</span>
            <span className="c3-ext-title">vault:supplier-x</span>
            <span className="c3-ext-active-badge">Active</span>
        </div>
        <div className="c3-ext-detail">
            <span>Joined Jun 04</span> &nbsp;·&nbsp; role: observer
            &nbsp;·&nbsp; reads on: <span>financial_clearance</span>
        </div>
        <button className="c3-view-btn">↗ View external commits</button>
    </div>
);
const C3PipelineCommit = () => (
    <>
        <div className="c3-pipe-divider" style={{ marginTop: 14 }}>
            <div className="c3-pipe-line" />
            <div className="c3-pipe-label">C3 Extension</div>
            <div className="c3-pipe-line" />
        </div>
        <div className="pipeline-step">
            <div className="ps-icon ps-icon-external">S·X</div>
            <div className="ps-content">
                <div className="ps-label ps-label-external">
                    counterparty_acknowledged{" "}
                    <span className="ext-role-badge">observer</span>
                </div>
                <div className="ps-sublabel ps-sublabel-external">
                    vault:supplier-x · joined Jun 04
                </div>
                <div
                    className="ps-ts"
                    style={{ color: "#C8922A", opacity: "0.6" }}
                >
                    reads on: financial_clearance
                </div>
            </div>
            <div className="ps-check done" style={{ color: "#C8922A" }}>
                ✓
            </div>
        </div>
    </>
)

// Dispute State
const DisputeBanner = () => (
    <div className="dispute-banner">
        <span className="db-icon">⚠</span>
        <span className="db-text">
            <span>Receipt rejected</span> by vault_direction — thread is
            disputed. Resolution required.
        </span>
        <span className="db-ts">Jun 12 · 16:44</span>
    </div>
)
const RejectedStep = () => (
    <>
        <div className="ps-connector red" />
        <div className="pipeline-step">
            <div className="ps-icon ps-icon-rejected">✕</div>
            <div className="ps-content">
                <div className="ps-label rejected">Receipt rejected</div>
                <div className="ps-sublabel rejected">
                    vault_direction · Jun 12 · 16:44
                </div>
                <div className="reject-reason-inline">
                    <div className="rri-label">Rejection reason</div>
                    Budget envelope reference is incorrect — Q3 cap is $200,000.
                    Contract value $240,000 exceeds approved limit. Requires
                    re-authorization from CFO before clearance can be accepted.
                </div>
            </div>
            <div className="ps-check reject">✕</div>
        </div>
        <div className="ps-connector" />
        <div className="pipeline-step">
            <div className="ps-icon ps-icon-wait">Dir</div>
            <div className="ps-content">
                <div className="ps-label wait">executive_signature</div>
                <div className="ps-sublabel">
                    vault_direction · blocked — disputed
                </div>
            </div>
            <div className="ps-check wait">○</div>
        </div>
    </>
)
const RejectedPipelineCommit = () => (
    <div className="commit-row rejected-row">
        <div className="commit-ts red">Jun 12 · 16:44</div>
        <div className="commit-body">
            <div>
                <span className="commit-actor red">vault_direction</span>
                <span className="commit-action red"> rejected receipt</span>
                <span className="reject-badge">✕ REJECTED</span>
            </div>
            <div className="commit-cid">
                receipt_tx: bafyreirej4…xt22 (Stellar anchored)
            </div>
            <div className="commit-reason">
                Budget envelope reference is incorrect — Q3 cap is $200,000.
                Contract value $240,000 exceeds approved limit. Requires
                re-authorization from CFO before clearance can be accepted.
            </div>
        </div>
        <div className="commit-verify red">verify ↗</div>
    </div>
)
const ResolutionPath = () => (
    <div className="dp-section">
        <div className="dp-section-title">Resolution path</div>
        <div className="resolution-box">
            <div className="rb-header">
                <span className="rb-icon">⚠</span>
                <span className="rb-title">
                    Thread is disputed — action required
                </span>
            </div>
            <div className="rb-steps">
                <div className="rb-step">
                    <div className="rb-step-num">1</div>vault_finance must address
                    the rejection — amend or resubmit financial_clearance with a
                    corrected budget reference.
                </div>
                <div className="rb-step">
                    <div className="rb-step-num">2</div>vault_direction must issue
                    a new receipt (received or processed) on the amended commit.
                </div>
                <div className="rb-step">
                    <div className="rb-step-num">3</div>Once acknowledged, the
                    thread resumes and executive_signature gate unlocks.
                </div>
            </div>
        </div>
    </div>

)
const DisputeActionButtons = () => (
    <div className="dp-actions">
        <button className="action-btn">↓ Export dispute log</button>
        <button className="action-btn primary">✉ Notify vault_finance</button>
    </div>
)



export { EmptyC3Extension, WithC3Extension, C3PipelineCommit, RejectedPipelineCommit, DisputeBanner, RejectedStep }
