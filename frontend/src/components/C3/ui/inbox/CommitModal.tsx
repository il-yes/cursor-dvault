import React, { useState } from "react";
import { createPortal } from "react-dom";
import { appendThreadEvent } from "@/services/api";
import { DerivedInboxItem } from "./InboxTable";

interface CommitModalProps {
    isOpen: boolean;
    onClose: () => void;
    item: DerivedInboxItem | null;
    onSuccess?: () => void;
}

export const CommitModal: React.FC<CommitModalProps> = ({
    isOpen,
    onClose,
    item,
    onSuccess,
}) => {
    const [commitValue, setCommitValue] = useState(
        "Approved. Financial review complete. Contract value within Q3 budget envelope."
    );
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    if (!isOpen || !item) return null;

    const isBlocked = item.gate_status === "blocked";
    const vaultName = item.vault_name || "vault_finance";
    const slotTag = item.slot_tag || "financial_clearance";
    const threadType = item.thread_type || "Contract";
    const threadTitle = item.thread_title || "Contract — Supplier X";
    const threadId = item.thread_id || "t_sha256_CONT_supplierx";

    const handleCommitSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (isBlocked) return;

        setIsLoading(true);
        setError(null);

        try {
            await appendThreadEvent({
                thread_id: item.thread_id,
                type: slotTag,
                payload: {
                    notes: commitValue.trim(),
                    vault_name: vaultName,
                    slot_tag: slotTag,
                    entry_type: threadType,
                    status: "committed",
                },
            });

            setIsLoading(false);
            if (onSuccess) onSuccess();
            onClose();
        } catch (err: any) {
            console.error("Commit failed:", err);
            setError(err?.message || "Failed to submit commit.");
            setIsLoading(false);
        }
    };

    return createPortal(
        <>
            {/* Scrim Overlay */}
            <div
                style={{
                    position: "fixed",
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    backgroundColor: "rgba(0, 0, 0, 0.45)",
                    backdropFilter: "blur(3px)",
                    zIndex: 1100,
                }}
                onClick={onClose}
            />

            {/* Modal Wrapper */}
            <div
                style={{
                    position: "fixed",
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    display: "flex",
                    alignItems: "center",
                    justify: "center",
                    zIndex: 1101,
                    pointerEvents: "none",
                }}
            >
                <div
                    style={{
                        width: "540px",
                        maxWidth: "92vw",
                        maxHeight: "90vh",
                        backgroundColor: "#ffffff",
                        borderRadius: "12px",
                        boxShadow: "0 24px 80px rgba(0,0,0,0.30), 0 8px 24px rgba(0,0,0,0.16)",
                        overflow: "hidden",
                        display: "flex",
                        flexDirection: "column",
                        pointerEvents: "auto",
                    }}
                >
                    {/* Header */}
                    <div
                        style={{
                            padding: "20px 24px 16px",
                            borderBottom: "1px solid #ebebeb",
                        }}
                    >
                        <div
                            style={{
                                display: "flex",
                                alignItems: "flex-start",
                                justifyContent: "space-between",
                                marginBottom: "8px",
                            }}
                        >
                            <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                                <div
                                    style={{
                                        display: "inline-flex",
                                        alignItems: "center",
                                        gap: "5px",
                                        background: "#EEF2FF",
                                        border: "1px solid #C7D2FE",
                                        padding: "3px 9px",
                                        borderRadius: "5px",
                                        fontSize: "11px",
                                        fontWeight: 600,
                                        color: "#3730A3",
                                    }}
                                >
                                    <div
                                        style={{
                                            width: "7px",
                                            height: "7px",
                                            borderRadius: "50%",
                                            background: "#2563EB",
                                        }}
                                    />
                                    {vaultName}
                                </div>
                            </div>
                            <div
                                style={{
                                    width: "24px",
                                    height: "24px",
                                    background: "#f5f5f5",
                                    border: "1px solid #e5e5e5",
                                    borderRadius: "6px",
                                    display: "flex",
                                    alignItems: "center",
                                    justify: "center",
                                    fontSize: "12px",
                                    color: "#888",
                                    cursor: "pointer",
                                }}
                                onClick={onClose}
                            >
                                ✕
                            </div>
                        </div>

                        <div style={{ fontSize: "16px", fontWeight: 700, color: "#1a1a1a" }}>
                            Commit: {slotTag}
                        </div>
                        <div style={{ fontSize: "12px", color: "#aaa", marginTop: "3px" }}>
                            Write and sign a value to this property slot
                        </div>

                        {/* Thread context strip */}
                        <div
                            style={{
                                display: "flex",
                                alignItems: "center",
                                gap: "10px",
                                padding: "10px 12px",
                                background: "#fafafa",
                                border: "1px solid #ebebeb",
                                borderRadius: "7px",
                                marginTop: "14px",
                            }}
                        >
                            <span
                                style={{
                                    fontSize: "9px",
                                    fontWeight: 700,
                                    letterSpacing: "0.1em",
                                    textTransform: "uppercase",
                                    color: "#C8922A",
                                    background: "#FBF0D8",
                                    padding: "2px 6px",
                                    borderRadius: "3px",
                                    flexShrink: 0,
                                }}
                            >
                                {threadType}
                            </span>
                            <span style={{ fontSize: "12px", fontWeight: 500, color: "#333", flex: 1 }}>
                                {threadTitle}
                            </span>
                            <span style={{ fontFamily: "monospace", fontSize: "10px", color: "#bbb" }}>
                                {threadId.length > 22 ? `${threadId.substring(0, 20)}…` : threadId}
                            </span>
                        </div>
                    </div>

                    {/* Body */}
                    <div
                        style={{
                            padding: "20px 24px",
                            display: "flex",
                            flexDirection: "column",
                            gap: "16px",
                            overflowY: "auto",
                        }}
                    >
                        {error && (
                            <div
                                style={{
                                    padding: "10px 12px",
                                    backgroundColor: "rgba(239, 68, 68, 0.08)",
                                    border: "1px solid rgba(239, 68, 68, 0.3)",
                                    borderRadius: "6px",
                                    color: "#DC2626",
                                    fontSize: "12px",
                                }}
                            >
                                ⚠️ {error}
                            </div>
                        )}

                        {/* Gate status */}
                        <div
                            style={{
                                display: "flex",
                                alignItems: "center",
                                gap: "8px",
                                padding: "10px 14px",
                                borderRadius: "7px",
                                background: isBlocked ? "#FFF7ED" : "#F0FDF4",
                                border: isBlocked ? "1px solid #FED7AA" : "1px solid #86EFAC",
                            }}
                        >
                            <span style={{ fontSize: "13px" }}>{isBlocked ? "✕" : "✓"}</span>
                            <span
                                style={{
                                    fontSize: "12px",
                                    flex: 1,
                                    color: isBlocked ? "#D97706" : "#059669",
                                }}
                            >
                                {isBlocked ? (
                                    <>Gate blocked — waiting on <strong>ops_confirmation</strong></>
                                ) : (
                                    <>Gate passed — <strong>{slotTag} exists</strong> (committed by {vaultName})</>
                                )}
                            </span>
                        </div>

                        {/* Field group */}
                        <div style={{ display: "flex", flexDirection: "column", gap: "5px" }}>
                            <div
                                style={{
                                    fontSize: "10px",
                                    fontWeight: 700,
                                    letterSpacing: "0.08em",
                                    textTransform: "uppercase",
                                    color: "#bbb",
                                }}
                            >
                                Commit value{" "}
                                <span style={{ fontSize: "11px", color: "#bbb", fontWeight: 400, textTransform: "none", letterSpacing: 0 }}>
                                    What {vaultName} is certifying
                                </span>
                            </div>
                            <textarea
                                style={{
                                    width: "100%",
                                    minHeight: "80px",
                                    padding: "10px 12px",
                                    border: "1.5px solid #e5e5e5",
                                    borderRadius: "8px",
                                    fontSize: "13px",
                                    color: "#1a1a1a",
                                    outline: "none",
                                    resize: "vertical",
                                    fontFamily: "Inter, sans-serif",
                                    lineHeight: 1.5,
                                }}
                                value={commitValue}
                                onChange={(e) => setCommitValue(e.target.value)}
                                disabled={isLoading || isBlocked}
                                placeholder="Enter commit notes..."
                            />
                        </div>

                        {/* Attachment */}
                        <div style={{ display: "flex", flexDirection: "column", gap: "5px" }}>
                            <div
                                style={{
                                    fontSize: "10px",
                                    fontWeight: 700,
                                    letterSpacing: "0.08em",
                                    textTransform: "uppercase",
                                    color: "#bbb",
                                }}
                            >
                                Attachment{" "}
                                <span style={{ fontSize: "11px", color: "#bbb", fontWeight: 400, textTransform: "none", letterSpacing: 0 }}>
                                    optional file — encrypted at rest
                                </span>
                            </div>
                            <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                                <div
                                    style={{
                                        display: "inline-flex",
                                        alignItems: "center",
                                        gap: "6px",
                                        padding: "7px 12px",
                                        border: "1px dashed #d0d0d0",
                                        borderRadius: "7px",
                                        fontSize: "12px",
                                        color: "#888",
                                        cursor: "pointer",
                                        background: "#fafafa",
                                    }}
                                >
                                    📎 Attach file
                                </div>
                                <span style={{ fontSize: "11px", color: "#bbb" }}>
                                    Max 200 MB — stored encrypted on IPFS
                                </span>
                            </div>
                        </div>

                        {/* Commit preview */}
                        <div style={{ background: "#1a1a1a", borderRadius: "8px", padding: "14px 16px" }}>
                            <div
                                style={{
                                    fontSize: "10px",
                                    fontWeight: 600,
                                    letterSpacing: "0.1em",
                                    textTransform: "uppercase",
                                    color: "#666",
                                    marginBottom: "10px",
                                }}
                            >
                                Commit preview
                            </div>
                            <div style={{ display: "flex", alignItems: "baseline", gap: "8px", marginBottom: "5px" }}>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#7C3AED", width: "110px", flexShrink: 0 }}>
                                    thread_id
                                </span>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#86EFAC", wordBreak: "break-all" }}>
                                    "{threadId}"
                                </span>
                            </div>
                            <div style={{ display: "flex", alignItems: "baseline", gap: "8px", marginBottom: "5px" }}>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#7C3AED", width: "110px", flexShrink: 0 }}>
                                    slot
                                </span>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#86EFAC", wordBreak: "break-all" }}>
                                    "{slotTag}"
                                </span>
                            </div>
                            <div style={{ display: "flex", alignItems: "baseline", gap: "8px", marginBottom: "5px" }}>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#7C3AED", width: "110px", flexShrink: 0 }}>
                                    vault
                                </span>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#86EFAC", wordBreak: "break-all" }}>
                                    "{vaultName}"
                                </span>
                            </div>
                            <div style={{ display: "flex", alignItems: "baseline", gap: "8px", marginBottom: "5px" }}>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#7C3AED", width: "110px", flexShrink: 0 }}>
                                    value_hash
                                </span>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#86EFAC", wordBreak: "break-all" }}>
                                    "sha256:e3b0c44298fc1c14…"
                                </span>
                                <span style={{ fontFamily: "monospace", fontSize: "11px", color: "#555" }}>
                                    ← AES-256
                                </span>
                            </div>
                            <div
                                style={{
                                    display: "flex",
                                    alignItems: "center",
                                    gap: "6px",
                                    marginTop: "10px",
                                    paddingTop: "10px",
                                    borderTop: "1px solid #2a2a2a",
                                }}
                            >
                                <span style={{ fontSize: "11px" }}>🔒</span>
                                <span style={{ fontSize: "11px", color: "#666" }}>
                                    Value encrypted with <strong style={{ color: "#C8922A" }}>AES-256</strong> before leaving device. Only {vaultName} and permitted vaults can decrypt.
                                </span>
                            </div>
                        </div>

                        {/* Stellar preview */}
                        <div
                            style={{
                                display: "flex",
                                alignItems: "center",
                                gap: "10px",
                                padding: "10px 14px",
                                background: "#fafafa",
                                border: "1px solid #ebebeb",
                                borderRadius: "7px",
                            }}
                        >
                            <span style={{ fontSize: "13px" }}>✦</span>
                            <span style={{ fontSize: "12px", color: "#555", flex: 1 }}>
                                This commit will be <strong>anchored on Stellar</strong> — tamper-proof, permanently verifiable.
                            </span>
                            <div style={{ display: "flex", alignItems: "center", gap: "5px" }}>
                                <div style={{ width: "6px", height: "6px", borderRadius: "50%", background: "#22C55E" }} />
                                <span style={{ fontSize: "11px", color: "#22C55E", fontWeight: 500 }}>Ready</span>
                            </div>
                        </div>
                    </div>

                    {/* Footer */}
                    <div
                        style={{
                            padding: "16px 24px",
                            borderTop: "1px solid #ebebeb",
                            display: "flex",
                            alignItems: "center",
                            justify: "space-between",
                            gap: "10px",
                        }}
                    >
                        <button
                            type="button"
                            style={{
                                background: "#f5f5f5",
                                border: "1px solid #e5e5e5",
                                color: "#555",
                                display: "inline-flex",
                                alignItems: "center",
                                gap: "6px",
                                padding: "9px 18px",
                                borderRadius: "7px",
                                fontSize: "13px",
                                fontWeight: 500,
                                cursor: "pointer",
                            }}
                            onClick={onClose}
                        >
                            Cancel
                        </button>
                        <button
                            type="button"
                            style={{
                                background: isBlocked ? "#aaa" : "#C8922A",
                                color: "#fff",
                                fontSize: "13px",
                                fontWeight: 600,
                                padding: "10px 22px",
                                borderRadius: "8px",
                                border: "none",
                                cursor: isBlocked || isLoading ? "not-allowed" : "pointer",
                                display: "flex",
                                alignItems: "center",
                                gap: "8px",
                            }}
                            onClick={handleCommitSubmit}
                            disabled={isLoading || isBlocked}
                        >
                            {isLoading ? "Signing & Committing..." : "✦ Sign & Commit"}
                        </button>
                    </div>
                </div>
            </div>
        </>,
        document.body
    );
};
