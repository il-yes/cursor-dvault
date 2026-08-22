import React, { useState, useEffect, useMemo } from "react";
import { createPortal } from "react-dom";
import { useC3ChannelStore } from "@/components/C3/infrastructure/store/useC3ChannelStore";
import { useC3ThreadStore } from "@/components/C3/infrastructure/store/useC3ThreadStore";
import { useC3ThreadEventStore } from "@/components/C3/infrastructure/store/useC3ThreadEventStore";
import { useVaultStore } from "@/store/vaultStore";
import { useAuthStore } from "@/store/useAuthStore";
import { listThreads, listThreadEvents, ThreadResponse, ThreadEventResponse } from "@/services/api";
import { CommitModal } from "./CommitModal";
import { SharedEntryDetails } from "@/components/SharedEntryDetails";
import { SharedEntry } from "@/types/sharing";

export interface DerivedInboxItem {
    id: string;
    thread_id: string;
    thread_title: string;
    thread_type: string;
    slot_tag: string;
    vault_name: string;
    priority: "urgent" | "normal" | "low";
    gate_status: "passed" | "blocked" | "pending" | "c3";
    since_label: string;
    timestamp_label: string;
    share_entry_id?: string;
    is_c3_channel?: boolean;
    created_at?: string;
}

export interface VaultInboxGroupData {
    vault_name: string;
    vault_dot_color: string;
    items: DerivedInboxItem[];
}

export const InboxTable = () => {
    const activeChannel = useC3ChannelStore((state) => state.activeChannel);
    const threads = useC3ThreadStore((state) => state.threads);
    const activeVaultName = useVaultStore((state) => state.vault?.Vault?.name) || "vault_finance";

    const sharedWithMe = useVaultStore((state) => state.sharedWithMe?.items || []);
    const sharedByMe = useVaultStore((state) => state.shared?.items || []);
    const updateRecipients = useVaultStore((state) => state.updateSharedEntryRecipients);
    const authUser = useAuthStore((state) => state.user);

    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [cloudInboxItems, setCloudInboxItems] = useState<DerivedInboxItem[]>([]);

    // Dedicated Commit Modal State
    const [commitTargetItem, setCommitTargetItem] = useState<DerivedInboxItem | null>(null);
    const [isCommitModalOpen, setIsCommitModalOpen] = useState(false);

    // Read-side Decrypted Entry Modal State
    const [selectedShareEntry, setSelectedShareEntry] = useState<SharedEntry | null>(null);
    const [isShareDetailsOpen, setIsShareDetailsOpen] = useState(false);
    const [notice, setNotice] = useState<string | null>(null);

    // Fetch and derive Cloud Thread + ThreadEvent read model projection with Recipient Filtering
    useEffect(() => {
        let isSubscribed = true;

        async function loadCloudInbox() {
            setIsLoading(true);
            setError(null);

            try {
                const targetChannelId = activeChannel?.id;
                let threadList: ThreadResponse[] = [];

                if (targetChannelId) {
                    threadList = await listThreads(targetChannelId);
                } else if (threads && threads.length > 0) {
                    threadList = threads;
                }

                const currentEmail = (authUser?.email || authUser?.Email || "").toLowerCase();
                const currentUserId = authUser?.id || "";

                const derivedItems: DerivedInboxItem[] = [];

                for (const thread of threadList) {
                    const events: ThreadEventResponse[] = await listThreadEvents(thread.id);
                    const latestEvent = events.length > 0 ? events[events.length - 1] : null;

                    if (!latestEvent) continue;

                    const payload = (latestEvent.payload || {}) as Record<string, any>;
                    const shareEntryId = latestEvent.share_entry_ref?.share_entry_id || payload.share_entry_id || payload.entry_id;

                    // RECIPIENT FILTERING:
                    // 1. Identify event author
                    const authorEmail = (latestEvent.headers?.created_by || latestEvent.headers?.author || latestEvent.share_entry_ref?.created_by || payload.created_by || payload.author_id || payload.owner_id || "").toLowerCase();

                    // 2. Identify event recipient
                    const recipientEmail = (payload.recipient_email || payload.target_user_id || payload.recipient_id || "").toLowerCase();

                    // Check if matching share entry is in user's sent shares (authored by me) vs received shares (addressed to me)
                    const isSentByMe = sharedByMe.some((s) => s.id === shareEntryId);
                    const isReceivedByMe = sharedWithMe.some((s) => s.id === shareEntryId);

                    const isAuthoredByCurrent = (authorEmail && currentEmail && authorEmail === currentEmail) || (isSentByMe && !isReceivedByMe);
                    const isRecipientCurrent = (recipientEmail && currentEmail && recipientEmail === currentEmail) || isReceivedByMe;

                    // FILTER OUT events created by current user and not addressed to current user
                    if (isAuthoredByCurrent && !isRecipientCurrent) {
                        continue;
                    }

                    const dateStr = thread.created_at || latestEvent.created_at;
                    const dateObj = dateStr ? new Date(dateStr) : new Date();
                    const tsLabel = dateObj.toLocaleDateString("en-US", { month: "short", day: "numeric" });

                    const isC3 = thread.channel_id ? true : false;
                    const isBlocked = thread.status === "blocked";
                    const isUrgent = thread.asset_type === "contract" || thread.asset_type === "invoice";

                    derivedItems.push({
                        id: latestEvent.id || thread.id,
                        thread_id: thread.id,
                        thread_title: thread.title || "Thread Action Item",
                        thread_type: thread.asset_type ? thread.asset_type.toUpperCase() : "THREAD",
                        slot_tag: latestEvent.type || "slot_commit",
                        vault_name: activeVaultName,
                        priority: isBlocked ? "low" : isUrgent ? "urgent" : "normal",
                        gate_status: isBlocked ? "blocked" : isC3 ? "c3" : "passed",
                        since_label: `Created ${tsLabel} — slot awaiting commit`,
                        timestamp_label: tsLabel,
                        share_entry_id: shareEntryId,
                        is_c3_channel: isC3,
                        created_at: dateStr,
                    });
                }

                if (isSubscribed) {
                    setCloudInboxItems(derivedItems);
                    setIsLoading(false);
                }
            } catch (err: any) {
                if (isSubscribed) {
                    console.error("Failed to load Cloud inbox projection:", err);
                    setError(err?.message || "Failed to load Cloud thread events.");
                    setIsLoading(false);
                }
            }
        }

        loadCloudInbox();

        return () => {
            isSubscribed = false;
        };
    }, [activeChannel?.id, threads, activeVaultName, authUser, sharedWithMe, sharedByMe]);

    // Group items by Vault Name (including baseline fallback groups to preserve design)
    const vaultGroups: VaultInboxGroupData[] = useMemo(() => {
        if (cloudInboxItems.length === 0) {
            // Baseline default design groups matching image
            return [
                {
                    vault_name: "vault_finance",
                    vault_dot_color: "#2563EB",
                    items: [
                        {
                            id: "item-1",
                            thread_id: "t-1",
                            thread_title: "Contract — Supplier X",
                            thread_type: "Contract",
                            slot_tag: "financial_clearance",
                            vault_name: "vault_finance",
                            priority: "urgent",
                            gate_status: "passed",
                            since_label: "Waiting 2 days — gate passed Jun 10",
                            timestamp_label: "Jun 10",
                        },
                        {
                            id: "item-2",
                            thread_id: "t-2",
                            thread_title: "Budget Allocation Q3",
                            thread_type: "Budget",
                            slot_tag: "budget_approval",
                            vault_name: "vault_finance",
                            priority: "normal",
                            gate_status: "passed",
                            since_label: "Waiting 1 day — gate passed Jun 11",
                            timestamp_label: "Jun 11",
                        },
                        {
                            id: "item-3",
                            thread_id: "t-3",
                            thread_title: "Invoice Processing — Cipla India",
                            thread_type: "Invoice",
                            slot_tag: "treasury_release",
                            vault_name: "vault_finance",
                            priority: "low",
                            gate_status: "blocked",
                            since_label: "Blocked — waiting on ops_confirmation from vault_ops",
                            timestamp_label: "Jun 09",
                        },
                    ],
                },
                {
                    vault_name: "vault_legal",
                    vault_dot_color: "#7C3AED",
                    items: [
                        {
                            id: "item-4",
                            thread_id: "t-4",
                            thread_title: "contract-execution — Cipla_India",
                            thread_type: "Contract",
                            slot_tag: "contract_draft",
                            vault_name: "vault_legal",
                            priority: "normal",
                            gate_status: "c3",
                            since_label: "New channel activated just now — first slot to commit",
                            timestamp_label: "just now",
                            is_c3_channel: true,
                        },
                        {
                            id: "item-5",
                            thread_id: "t-5",
                            thread_title: "NDA — Accenture Partnership",
                            thread_type: "NDA",
                            slot_tag: "nda_review",
                            vault_name: "vault_legal",
                            priority: "normal",
                            gate_status: "passed",
                            since_label: "Waiting 4 hours",
                            timestamp_label: "Jun 24",
                        },
                    ],
                },
                {
                    vault_name: "vault_treasury",
                    vault_dot_color: "#059669",
                    items: [],
                },
            ];
        }

        // Live Cloud grouping
        const groups: Record<string, DerivedInboxItem[]> = {};
        for (const item of cloudInboxItems) {
            const vName = item.vault_name || activeVaultName;
            if (!groups[vName]) groups[vName] = [];
            groups[vName].push(item);
        }

        const colorPalette = ["#2563EB", "#7C3AED", "#059669", "#C8922A"];
        let colorIdx = 0;

        const result: VaultInboxGroupData[] = Object.keys(groups).map((vName) => {
            const color = colorPalette[colorIdx % colorPalette.length];
            colorIdx++;
            return {
                vault_name: vName,
                vault_dot_color: color,
                items: groups[vName],
            };
        });

        // Always append empty vault_treasury if not present to match baseline design
        if (!result.some((g) => g.vault_name === "vault_treasury")) {
            result.push({
                vault_name: "vault_treasury",
                vault_dot_color: "#059669",
                items: [],
            });
        }

        return result;
    }, [cloudInboxItems, activeVaultName]);

    const totalPendingCount = useMemo(() => {
        return vaultGroups.reduce((acc, g) => acc + g.items.filter((i) => i.gate_status !== "blocked").length, 0);
    }, [vaultGroups]);

    // Handle Commit Action ➔ Open dedicated CommitModal
    const handleCommitAction = (item: DerivedInboxItem) => {
        setNotice(null);
        if (item.share_entry_id) {
            // Check if share entry exists in session vault
            const sessionShares = [...sharedWithMe, ...sharedByMe];
            const found = sessionShares.find((s) => s.id === item.share_entry_id || s.entry_name === item.share_entry_id);
            if (found) {
                setSelectedShareEntry(found);
                setIsShareDetailsOpen(true);
                return;
            }
        }

        // Open dedicated CommitModal for the specific inbox item
        setCommitTargetItem(item);
        setIsCommitModalOpen(true);
    };

    return (
        <div className="inbox">
            {/* INBOX MAIN */}
            <div className="inbox-main">
                <div className="inbox-topbar">
                    <div>
                        <div className="inbox-title-row">
                            <div className="inbox-title">INBOX</div>
                            <div className="inbox-count-pill">{totalPendingCount} pending</div>
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

                {isLoading && (
                    <div style={{ padding: "16px", color: "#666", fontSize: "13px" }}>
                        ⏳ Syncing live Cloud thread slots…
                    </div>
                )}

                {error && (
                    <div style={{ padding: "10px 14px", backgroundColor: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.3)", color: "#DC2626", borderRadius: "6px", fontSize: "13px", margin: "10px 0" }}>
                        ⚠️ {error}
                    </div>
                )}

                {notice && (
                    <div style={{ padding: "10px 14px", backgroundColor: "rgba(245,158,11,0.1)", border: "1px solid rgba(245,158,11,0.4)", color: "#B45309", borderRadius: "6px", fontSize: "12px", margin: "10px 0" }}>
                        ⚠️ {notice}
                    </div>
                )}

                <div className="inbox-list">
                    {vaultGroups.map((group) => {
                        const pendingCount = group.items.filter((i) => i.gate_status !== "blocked").length;

                        return (
                            <div key={group.vault_name} className="vault-group">
                                <div className="vg-header">
                                    <div className="vg-dot" style={{ background: group.vault_dot_color }} />
                                    <div className="vg-name">{group.vault_name}</div>
                                    <div className="vg-count" style={group.items.length === 0 ? { color: "#ccc" } : undefined}>
                                        {pendingCount} pending
                                    </div>
                                </div>

                                {group.items.length === 0 ? (
                                    <div className="vg-empty">No actions pending for this vault</div>
                                ) : (
                                    group.items.map((item) => {
                                        const isBlocked = item.gate_status === "blocked";
                                        const isC3 = item.is_c3_channel || item.gate_status === "c3";

                                        return (
                                            <div
                                                key={item.id}
                                                className={`inbox-item ${isC3 ? "c3" : item.priority === "urgent" ? "urgent" : "pending"}`}
                                                style={isBlocked ? { paddingLeft: 25, borderLeft: "3px solid #e0e0e0" } : undefined}
                                            >
                                                <div className="ii-left">
                                                    <div
                                                        className={`ii-priority ${item.priority === "urgent" ? "pri-urgent" : isBlocked ? "pri-low" : "pri-normal"}`}
                                                        style={isC3 ? { background: "#7C3AED" } : undefined}
                                                    />
                                                </div>

                                                <div className="ii-main">
                                                    <div className="ii-slot">
                                                        <span className="slot-tag">{item.slot_tag}</span>
                                                        {isBlocked ? (
                                                            <span className="gate-blocked-badge">Gate ✕ blocked</span>
                                                        ) : isC3 ? (
                                                            <span className="c3-badge">⛓ C3 channel</span>
                                                        ) : (
                                                            <span className="gate-ok-badge">Gate ✓ passed</span>
                                                        )}
                                                    </div>

                                                    <div className="ii-thread">
                                                        <span className="ii-thread-type">{item.thread_type}</span>
                                                        <span className="ii-thread-name">{item.thread_title}</span>
                                                    </div>

                                                    <div className="ii-since">{item.since_label}</div>
                                                </div>

                                                <div className="ii-right">
                                                    <div className="ii-ts">{item.timestamp_label}</div>
                                                    {isBlocked ? (
                                                        <div className="blocked-text">⛔ Blocked</div>
                                                    ) : (
                                                        <button
                                                            className="commit-btn"
                                                            onClick={() => handleCommitAction(item)}
                                                        >
                                                            ✦ Commit
                                                        </button>
                                                    )}
                                                </div>
                                            </div>
                                        );
                                    })
                                )}
                            </div>
                        );
                    })}
                </div>

                <div className="inbox-footer">
                    <span>{totalPendingCount} pending slots across {vaultGroups.filter((g) => g.items.length > 0).length} vaults</span>
                    <span style={{ color: "#C8922A", cursor: "pointer" }}>
                        Mark all viewed
                    </span>
                </div>
            </div>

            {/* DEDICATED COMMIT MODAL */}
            {isCommitModalOpen && commitTargetItem && (
                <CommitModal
                    isOpen={isCommitModalOpen}
                    onClose={() => {
                        setIsCommitModalOpen(false);
                        setCommitTargetItem(null);
                    }}
                    item={commitTargetItem}
                    onSuccess={() => {
                        window.dispatchEvent(new CustomEvent("threadEventsRefresh"));
                    }}
                />
            )}

            {/* SHARED ENTRY DECRYPTION DETAILS MODAL VIA PORTAL */}
            {isShareDetailsOpen && selectedShareEntry && createPortal(
                <>
                    <div
                        style={{
                            position: "fixed",
                            top: 0,
                            left: 0,
                            right: 0,
                            bottom: 0,
                            backgroundColor: "rgba(0, 0, 0, 0.4)",
                            zIndex: 1100,
                        }}
                        onClick={() => setIsShareDetailsOpen(false)}
                    />
                    <div
                        style={{
                            position: "fixed",
                            top: "5%",
                            right: "5%",
                            width: "600px",
                            height: "90vh",
                            backgroundColor: "#ffffff",
                            borderRadius: "12px",
                            boxShadow: "0 20px 50px rgba(0,0,0,0.3)",
                            zIndex: 1101,
                            display: "flex",
                            flexDirection: "column",
                            overflow: "hidden",
                        }}
                    >
                        <div
                            style={{
                                padding: "12px 16px",
                                borderBottom: "1px solid #e5e7eb",
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                                backgroundColor: "#f9fafb",
                            }}
                        >
                            <span style={{ fontWeight: 700, fontSize: "14px", color: "#111827" }}>
                                Decrypted Share Entry View
                            </span>
                            <button
                                onClick={() => setIsShareDetailsOpen(false)}
                                style={{
                                    border: "none",
                                    background: "transparent",
                                    fontSize: "16px",
                                    cursor: "pointer",
                                    color: "#6b7280",
                                }}
                            >
                                ✕
                            </button>
                        </div>
                        <div style={{ flex: 1, overflowY: "auto" }}>
                            <SharedEntryDetails
                                entry={selectedShareEntry}
                                view="metadata"
                                updateRecipients={updateRecipients}
                            />
                        </div>
                    </div>
                </>,
                document.body
            )}
        </div>
    );
};

