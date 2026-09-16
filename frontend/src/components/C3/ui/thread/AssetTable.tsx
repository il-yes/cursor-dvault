import { Channel } from "../../domain/channel/channel.types";
import { AssetView, ThreadAssetViewInterface } from "../../domain/thread/asset.types";
import { AssetSummary, } from "../../domain/thread/asset.types";
import { AnchorSection } from "../../features/AnchorSection";
import { EventTimelineSection } from "../../features/EventTimelineSection";
import { LifecycleSection } from "../../features/LifecycleSection";
import { PayloadSection } from "../../features/PayloadSection";
import { PolicySection } from "../../features/PolicySection";
import { ReceiptSection } from "../../features/ReceiptSection";
import { RelationshipSection } from "../../features/RelationshipSection";
import '../styles/event-thread-extended-style.css'


export const AssetCard = ({
    asset
}: {
    asset: AssetView
}) => (

    <div className="asset-card">


        <div className="asset-header">


            <div>

                <span className="asset-type">

                    {asset.type}

                </span>


                <h3>
                    {asset.title}
                </h3>


                <p>
                    {asset.subtitle}
                </p>

            </div>



            <span
                className={`asset-status ${asset.status}`}
            >

                {asset.status}

            </span>


        </div>



        <div className="asset-meta">


            <div>

                <label>
                    Last Event
                </label>

                <p>
                    {asset.lastEvent.type}
                </p>

            </div>



            <div>

                <label>
                    Activity
                </label>

                <p>
                    {asset.lastEvent.at}
                </p>

            </div>


        </div>



        <div className="asset-footer">


            <div className="participants">

                {
                    asset.participants.map(p =>
                        <span key={p}>
                            {p}
                        </span>
                    )
                }

            </div>


            <div>

                <span className="stellar-val">
                    {asset.stellarTx}
                </span>

            </div>


        </div>



        <div className="asset-actions">

            <button>
                Open
            </button>


            <button>
                Timeline
            </button>


        </div>


    </div>

);


export const AssetSummaryCard = ({
    asset
}: {
    asset: AssetSummary
}) => (

    <div className="asset-card">


        <div className="asset-header">


            <div>

                <span className="asset-type">

                    {asset.type}

                </span>


                <h3>
                    {asset.title}
                </h3>

            </div>



            <span
                className={`asset-status ${asset.status}`}
            >

                {asset.status}

            </span>


        </div>




        <div className="asset-event">

            Last event:

            <strong>
                {asset.lastEvent}
            </strong>

        </div>




        <div className="asset-actions">

            <button>
                Open
            </button>


        </div>


    </div>

);



export function ThreadAssetView_ALPHA({ hasConflict }: { hasConflict: boolean }) {
    const hasC3Extension = true;

    return (
        <div className="detail-panel">
            {hasConflict && <DisputeBanner />}
            <div className="dp-header">
                <div className="dp-close-row">
                    <span className="dp-badge">Invoice</span>
                    <div className="dp-close">✕</div>
                </div>
                <div className="dp-title">Invoice Processing — Cipla India</div>
                <div className="dp-meta">
                    <span className="dp-meta-chip">t_sha256_INV_cipla</span>
                    <span className="dp-meta-chip">created: 2026-06-01</span>
                    <span className="dp-meta-chip">by: vault_ops</span>
                </div>
                <div className="dp-custom">
                    <span style={{ color: "#bbb" }}>vendor:</span>{" "}
                    <span>Cipla India</span>
                    &nbsp;·&nbsp;
                    <span style={{ color: "#bbb" }}>ref:</span> <span>INV-2047</span>
                    &nbsp;·&nbsp;
                    <span style={{ color: "#bbb" }}>amount:</span> <span>$18,400</span>
                </div>
            </div>
            <div className="dp-body">
                {/* PIPELINE */}
                <div className="dp-section">
                    <div className="dp-section-title">Pipeline</div>
                    <div className="pipeline-step">
                        <div className="ps-icon" style={{ background: "#DC2626" }}>
                            O
                        </div>
                        <div className="ps-content">
                            <div className="ps-label">invoice_submitted</div>
                            <div className="ps-sublabel">vault_ops</div>
                            <div className="ps-ts">Jun 01 · 09:14</div>
                        </div>
                        <div className="ps-check done">✓</div>
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
                {hasConflict ? <ResolutionPath /> : <div className="dp-section">
                    <div className="dp-section-title">C3 Extension</div>
                    {!hasC3Extension && <EmptyC3Extension />}
                    {hasC3Extension && <WithC3Extension />}
                </div>}
            </div>
            {/* ACTIONS */}
            {hasConflict ? <DisputeActionButtons /> : <div className="dp-actions">
                <button className="action-btn">↓ Export thread</button>
                <button className="action-btn">🔗 Share read access</button>
                <button className="action-btn danger">✕ Close thread</button>
            </div>}
        </div>
    );
}
import { useState, useEffect, useMemo } from "react";
import { createPortal } from "react-dom";
import { useC3ThreadEventStore } from "../../infrastructure/store/useC3ThreadEventStore";
import { AppendThreadEventSlidingView } from "../../AppendThreadEventModal";
import { useVaultStore } from "@/store/vaultStore";
import { SharedEntryDetails } from "@/components/SharedEntryDetails";
import { SharedEntry } from "@/types/sharing";
import { getShareEntry, C3ShareEntryResponse } from "@/services/api";
import { C3ResourceCard } from "@/components/C3/C3ResourceCard";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ParticipantsPanel } from "../channel/ParticipantsPanel";
import { InvitationsPanel } from "../channel/InvitationsPanel";

export function ThreadAssetView({ channel, asset, hasConflict }: { channel: Channel | null, asset: ThreadAssetViewInterface, hasConflict: boolean }) {
    const hasC3Extension = true;
    const [isAppendOpen, setIsAppendOpen] = useState(false);

    // Read-side C3 Share Entry resolution state
    const [selectedC3ShareEntry, setSelectedC3ShareEntry] = useState<C3ShareEntryResponse | null>(null);
    const [isShareDetailsOpen, setIsShareDetailsOpen] = useState(false);
    const [isResolvingC3Share, setIsResolvingC3Share] = useState(false);
    const [shareResolutionNotice, setShareResolutionNotice] = useState<string | null>(null);

    const updateRecipients = useVaultStore((state) => state.updateSharedEntryRecipients);

    const events = useC3ThreadEventStore((state) => state.events);
    const isLoadingEvents = useC3ThreadEventStore((state) => state.isLoading);
    const eventError = useC3ThreadEventStore((state) => state.error);
    const fetchEvents = useC3ThreadEventStore((state) => state.fetchEvents);

    useEffect(() => {
        if (asset?.id) {
            fetchEvents(asset.id);
        }
    }, [asset?.id, fetchEvents]);

    const handleSeeShareEntry = async (entryId: string) => {
        setShareResolutionNotice(null);
        if (!entryId) {
            setShareResolutionNotice("Malformed event: 'entry.shared' missing authoritative share_entry_id.");
            return;
        }

        setIsResolvingC3Share(true);
        try {
            const shareEntry = await getShareEntry(entryId);
            setSelectedC3ShareEntry(shareEntry);
            setIsShareDetailsOpen(true);
        } catch (err: any) {
            const errMsg = err?.message || String(err);
            setShareResolutionNotice(`Failed to resolve C3 ShareEntry "${entryId}": ${errMsg}`);
        } finally {
            setIsResolvingC3Share(false);
        }
    };

    events.length > 0 ? console.log(events) : null

    return (
        <div className="detail-panel">
            {hasConflict && <DisputeBanner />}
            <div className="dp-header">
                <div className="dp-close-row">
                    <span className="dp-badge">{asset?.type || "Thread"}</span>
                    <div className="dp-close">✕</div>
                </div>
                <div className="dp-title">{asset?.title || "Thread Detail"}</div>
                <div className="dp-meta">
                    <span className="dp-meta-chip">ID: {asset?.id}</span>
                    <span className="dp-meta-chip">Status: {asset?.status || "open"}</span>
                    <span className="dp-meta-chip">Channel: {channel?.title || asset?.channelId}</span>
                </div>
            </div>
            <div className="dp-body">
                {/* DYNAMIC THREAD EVENT TIMELINE */}
                <div className="dp-section">
                    <div className="dp-section-title" style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <span>Authoritative Thread Events</span>
                        <span style={{ fontSize: "12px", color: "#888", fontWeight: 400 }}>
                            {events.length} event{events.length === 1 ? "" : "s"}
                        </span>
                    </div>

                    {isLoadingEvents && (
                        <div style={{ padding: "12px", color: "#666", fontSize: "13px" }}>
                            ⏳ Loading events from Cloud backend…
                        </div>
                    )}

                    {eventError && (
                        <div style={{ padding: "10px", backgroundColor: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.3)", color: "#DC2626", borderRadius: "6px", fontSize: "13px" }}>
                            ⚠️ Failed to load events: {eventError}
                        </div>
                    )}

                    {shareResolutionNotice && (
                        <div style={{ padding: "10px", backgroundColor: "rgba(245, 158, 11, 0.1)", border: "1px solid rgba(245, 158, 11, 0.4)", color: "#B45309", borderRadius: "6px", fontSize: "12px", marginBottom: "8px" }}>
                            ⚠️ {shareResolutionNotice}
                        </div>
                    )}

                    {!isLoadingEvents && !eventError && events.length === 0 && (
                        <div style={{ padding: "14px", backgroundColor: "#fafafa", border: "1px dashed #ddd", borderRadius: "6px", color: "#888", fontSize: "13px", textAlign: "center" }}>
                            No thread events recorded yet. Click <strong>+ Append Event</strong> below to append an event.
                        </div>
                    )}

                    {!isLoadingEvents && events.length > 0 && (
                        <div style={{ display: "flex", flexDirection: "column", gap: "8px", marginTop: "8px" }}>
                            {events.map((evt) => {
                                const payload = (evt.payload || {}) as Record<string, any>;
                                const shareEntryId = evt.share_entry_ref?.share_entry_id || evt.resource_ref?.share_entry_id || payload.share_entry_id || payload.entry_id;
                                const trustGroupId = evt.share_entry_ref?.trust_group_id || evt.resource_ref?.trust_group_id || payload.trust_group_id;
                                const entryName = payload.entry_name || "C3 Share Entry";
                                const entryType = payload.entry_type || payload.ref_type || evt.resource_ref?.ref_type || "share_entry";
                                const notes = payload.notes;

                                return (
                                    <div
                                        key={evt.id}
                                        style={{
                                            padding: "10px 12px",
                                            backgroundColor: "#ffffff",
                                            border: "1px solid #e5e7eb",
                                            borderRadius: "6px",
                                            fontSize: "12px",
                                            display: "flex",
                                            flexDirection: "column",
                                            gap: "4px",
                                        }}
                                    >
                                        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                                            <span style={{ fontWeight: 600, color: "#111827", backgroundColor: "#f3f4f6", padding: "2px 6px", borderRadius: "4px", fontSize: "11px" }}>
                                                {evt.type}
                                            </span>
                                            <span style={{ color: "#9ca3af", fontSize: "11px" }}>
                                                {evt.created_at ? new Date(evt.created_at).toLocaleString() : "Just now"}
                                            </span>
                                        </div>

                                        <div style={{ color: "#6b7280", fontFamily: "monospace", fontSize: "11px" }}>
                                            ID: {evt.id}
                                        </div>

                                        {evt.type === "entry.shared" && (
                                            <div style={{ marginTop: "4px", padding: "8px 10px", backgroundColor: "#f0fdf4", border: "1px solid #bbf7d0", borderRadius: "6px", color: "#166534", display: "flex", flexDirection: "column", gap: "4px" }}>
                                                {!shareEntryId ? (
                                                    <div style={{ color: "#dc2626", fontWeight: 600 }}>
                                                        ⚠️ Malformed event: 'entry.shared' missing authoritative share_entry_id.
                                                    </div>
                                                ) : (
                                                    <>
                                                        <div><strong>Referenced C3 Entry:</strong> [{entryType.toUpperCase()}] {entryName}</div>
                                                        <div style={{ fontFamily: "monospace", fontSize: "11px" }}><strong>ShareEntry ID:</strong> {shareEntryId}</div>
                                                        {trustGroupId && <div style={{ fontFamily: "monospace", fontSize: "11px" }}><strong>TrustGroup ID:</strong> {trustGroupId}</div>}
                                                        {notes && <div style={{ fontSize: "11px", color: "#15803d", fontStyle: "italic" }}>"{notes}"</div>}
                                                        
                                                        <div style={{ display: "flex", gap: "8px", alignItems: "center", marginTop: "4px" }}>
                                                            <button
                                                                onClick={() => handleSeeShareEntry(shareEntryId)}
                                                                disabled={isResolvingC3Share}
                                                                style={{
                                                                    padding: "4px 10px",
                                                                    backgroundColor: "#C8922A",
                                                                    color: "#ffffff",
                                                                    border: "none",
                                                                    borderRadius: "4px",
                                                                    fontSize: "11px",
                                                                    fontWeight: 600,
                                                                    cursor: isResolvingC3Share ? "wait" : "pointer",
                                                                    display: "inline-flex",
                                                                    alignItems: "center",
                                                                    gap: "4px",
                                                                    opacity: isResolvingC3Share ? 0.7 : 1,
                                                                }}
                                                            >
                                                                {isResolvingC3Share ? "⏳ Resolving ShareEntry…" : "👁 View C3 ShareEntry Metadata"}
                                                            </button>
                                                        </div>

                                                        <div style={{ marginTop: "6px" }}>
                                                            <C3ResourceCard
                                                                refType="share_entry"
                                                                shareEntryId={shareEntryId}
                                                                trustGroupId={trustGroupId}
                                                                cid={payload.asset_cid || payload.cid}
                                                                author={payload.created_by}
                                                                createdAt={evt.created_at}
                                                            />
                                                        </div>
                                                    </>
                                                )}
                                            </div>
                                        )}

                                        {evt.payload_ref && (
                                            <div style={{ marginTop: "4px", padding: "6px 8px", backgroundColor: "#f8fafc", border: "1px solid #e2e8f0", borderRadius: "4px", color: "#334155" }}>
                                                <div><strong>CID:</strong> {evt.payload_ref.cid}</div>
                                                <div><strong>Content Hash:</strong> {evt.payload_ref.content_hash}</div>
                                                <div><strong>Size:</strong> {evt.payload_ref.size} bytes</div>
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </div>

                {/* C3 SHARE ENTRY METADATA DETAILS MODAL VIA RADIX DIALOG */}
                <Dialog open={isShareDetailsOpen && !!selectedC3ShareEntry} onOpenChange={(open) => !open && setIsShareDetailsOpen(false)}>
                    <DialogContent className="sm:max-w-[640px] max-h-[90vh] p-0 overflow-hidden bg-background border border-border rounded-xl shadow-2xl">
                        <DialogHeader className="p-4 border-b border-border bg-muted/30 flex flex-row items-center justify-between space-y-0">
                            <DialogTitle className="text-sm font-bold text-foreground">
                                Authoritative C3 ShareEntry Metadata
                            </DialogTitle>
                        </DialogHeader>
                        <div className="flex-1 overflow-y-auto p-4 max-h-[calc(90vh-60px)] space-y-4">
                            {selectedC3ShareEntry && (
                                <div className="space-y-3 text-xs">
                                    <div className="grid grid-cols-2 gap-3 p-3 bg-muted/20 border border-border rounded-lg">
                                        <div>
                                            <div className="text-muted-foreground font-semibold">ShareEntry ID</div>
                                            <div className="font-mono text-foreground font-medium break-all">{selectedC3ShareEntry.id}</div>
                                        </div>
                                        <div>
                                            <div className="text-muted-foreground font-semibold">TrustGroup ID</div>
                                            <div className="font-mono text-foreground font-medium break-all">{selectedC3ShareEntry.trust_group_id}</div>
                                        </div>
                                        <div>
                                            <div className="text-muted-foreground font-semibold">Status</div>
                                            <div className="inline-block px-2 py-0.5 mt-0.5 rounded text-[10px] font-bold uppercase bg-emerald-500/10 text-emerald-600 border border-emerald-500/20">
                                                {selectedC3ShareEntry.status || "active"}
                                            </div>
                                        </div>
                                        <div>
                                            <div className="text-muted-foreground font-semibold">Creator</div>
                                            <div className="font-mono text-foreground">{selectedC3ShareEntry.created_by}</div>
                                        </div>
                                        <div>
                                            <div className="text-muted-foreground font-semibold">KEK Version</div>
                                            <div className="font-mono text-foreground">v{selectedC3ShareEntry.kek_version}</div>
                                        </div>
                                        <div>
                                            <div className="text-muted-foreground font-semibold">Created At</div>
                                            <div className="text-foreground">
                                                {selectedC3ShareEntry.created_at ? new Date(selectedC3ShareEntry.created_at).toLocaleString() : "N/A"}
                                            </div>
                                        </div>
                                    </div>

                                    {selectedC3ShareEntry.asset_cid && (
                                        <div className="p-3 bg-muted/20 border border-border rounded-lg space-y-1">
                                            <div className="text-muted-foreground font-semibold">Referenced Asset CID</div>
                                            <div className="font-mono text-foreground break-all">{selectedC3ShareEntry.asset_cid}</div>
                                        </div>
                                    )}

                                    {selectedC3ShareEntry.metadata && Object.keys(selectedC3ShareEntry.metadata).length > 0 && (
                                        <div className="p-3 bg-muted/20 border border-border rounded-lg space-y-2">
                                            <div className="text-muted-foreground font-semibold">Custom Metadata</div>
                                            <div className="space-y-1">
                                                {Object.entries(selectedC3ShareEntry.metadata).map(([k, v]) => (
                                                    <div key={k} className="flex justify-between border-b border-border/40 pb-1">
                                                        <span className="text-muted-foreground font-mono">{k}:</span>
                                                        <span className="text-foreground font-mono">{v}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </DialogContent>
                </Dialog>

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
                {hasConflict ? <ResolutionPath /> : <div className="dp-section">
                    <div className="dp-section-title">C3 Extension</div>
                    {!hasC3Extension && <EmptyC3Extension />}
                    {hasC3Extension && <WithC3Extension />}
                </div>}

                {/* ACCESS & MEMBERSHIP (Add Participant vs Add Member) */}
                {(asset?.channelId || channel?.id) && (
                    <div className="dp-section">
                        <div className="dp-section-title">
                            Thread Access & Membership
                        </div>
                        <div style={{ display: "flex", flexDirection: "column", gap: "12px", marginTop: "8px" }}>
                            <div style={{ padding: "10px 12px", border: "1px solid #e5e7eb", borderRadius: "8px", backgroundColor: "#f9fafb" }}>
                                <div style={{ fontSize: "12px", fontWeight: 700, color: "#1f2937", marginBottom: "6px" }}>
                                    👤 Channel Participation ("Add Participant")
                                </div>
                                <ParticipantsPanel channelId={asset?.channelId || channel?.id || ""} />
                            </div>
                            <div style={{ padding: "10px 12px", border: "1px solid #e5e7eb", borderRadius: "8px", backgroundColor: "#f9fafb" }}>
                                <div style={{ fontSize: "12px", fontWeight: 700, color: "#1f2937", marginBottom: "6px" }}>
                                    🛡️ TrustGroup & Member Invitations ("Add Member")
                                </div>
                                <InvitationsPanel channelId={asset?.channelId || channel?.id || ""} />
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {/* ACTIONS */}
            {hasConflict ? <DisputeActionButtons /> : <div className="dp-actions">
                <button
                    className="action-btn primary"
                    onClick={() => setIsAppendOpen(true)}
                    style={{ backgroundColor: "#2563eb", color: "#fff", fontWeight: 600 }}
                >
                    + Append Event
                </button>
                <button className="action-btn">↓ Export thread</button>
                <button className="action-btn danger">✕ Close thread</button>
            </div>}

            {/* APPEND THREAD EVENT SLIDING MODAL */}
            {asset?.id && (
                <AppendThreadEventSlidingView
                    isOpen={isAppendOpen}
                    activeWorkspaceName="Workspace"
                    activeChannelTitle={channel?.title || "Channel"}
                    activeThreadId={asset.id}
                    activeThreadTitle={asset.title || asset.id}
                    onClose={() => setIsAppendOpen(false)}
                    onEventAppended={(newEvent) => {
                        console.log("[THREAD UI] Event appended successfully:", newEvent);
                        if (asset?.id) {
                            fetchEvents(asset.id);
                        }
                    }}
                />
            )}
        </div>
    );
}
export default function ThreadAssetPanel({
    asset,
    onClose,
}: {
    asset: ThreadAssetViewInterface;
    onClose: () => void;
}) {
    return (
        <div className="detail-panel">

            {/* HEADER */}

            <div className="dp-header">

                <div className="dp-close-row">

                    <span className="dp-badge">
                        {asset.type}
                    </span>

                    <div
                        className="dp-close"
                        onClick={onClose}
                    >
                        ✕
                    </div>

                </div>

                <div className="dp-title">
                    {asset.title}
                </div>

                <div className="dp-meta">

                    <span className="dp-meta-chip">
                        {asset.id}
                    </span>

                    <span className="dp-meta-chip">
                        {asset.status}
                    </span>

                    <span className="dp-meta-chip">
                        created {asset.createdAt}
                    </span>

                </div>

                <div className="flow">

                    {asset.participants.map((p, i) => (

                        <span key={i}>

                            {i > 0 &&
                                <span className="fa">
                                    →
                                </span>
                            }

                            <span
                                className="vb"
                                style={{
                                    background: p.color
                                }}
                            >
                                {p.label}
                            </span>

                        </span>

                    ))}

                </div>

            </div>


            <div className="dp-body">

                <LifecycleSection asset={asset} />

                <EventTimelineSection asset={asset} />

                <PayloadSection asset={asset} />

                <ReceiptSection asset={asset} />

                <PolicySection asset={asset} />

                <AnchorSection asset={asset} />

                <RelationshipSection asset={asset} />

            </div>


            <div className="dp-actions">

                <button className="action-btn">
                    Export Timeline
                </button>

                <button className="action-btn">
                    Verify Thread
                </button>

                <button className="action-btn danger">
                    Close Thread
                </button>

            </div>

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

