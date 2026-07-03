import React from "react";

export const ReceiptAckModal = () => {
    return (
    
                <div className="modal">
                    <div className="modal-header">
                        <div className="mh-top">
                            <div className="vault-badge vb-dir">
                                <div className="vbdot" style={{ background: "#444" }} />
                                vault_direction
                            </div>
                            <div className="mh-close">✕</div>
                        </div>
                        <div className="mh-title">Receipt: financial_clearance</div>
                        <div className="mh-sub">
                            Acknowledge a commit received from another vault
                        </div>
                        <div className="thread-ctx">
                            <span className="tc-type">Contract</span>
                            <span className="tc-name">Contract — Supplier X</span>
                            <span className="tc-tid">t_sha256_CONT_supplierx</span>
                        </div>
                    </div>
                    <div className="modal-body">
                        {/* Incoming commit details */}
                        <div className="commit-box">
                            <div className="cb-header">
                                <div className="cb-from-badge">
                                    <div className="cb-from-dot" />
                                    vault_finance
                                </div>
                                <div className="cb-arrow">→ committed</div>
                                <div className="cb-slot">financial_clearance</div>
                                <div className="cb-ts">Jun 12 · 14:07</div>
                            </div>
                            <div className="cb-value">
                                Approved. Financial review complete. Contract value $240,000
                                within Q3 budget envelope. Risk: LOW. Ref: FIN-2026-0412.
                            </div>
                            <div className="cb-meta">
                                <div className="cb-cid">CID: bafyreihzm7a4kpq9x3d9c…</div>
                                <div className="cb-verify">verify on Stellar ↗</div>
                            </div>
                        </div>
                        {/* Acknowledgment options — Reject path shown active */}
                        <div className="action-options">
                            <div className="action-opt">
                                <div className="ao-radio" />
                                <div className="ao-content">
                                    <div className="ao-label ok">Received</div>
                                    <div className="ao-desc">
                                        Acknowledge you have seen this commit. No further action
                                        required from your vault.
                                    </div>
                                </div>
                            </div>
                            <div className="action-opt">
                                <div className="ao-radio" />
                                <div className="ao-content">
                                    <div className="ao-label ok">Processed</div>
                                    <div className="ao-desc">
                                        Confirm you have acted on this commit — reviewed the value,
                                        verified the data, or taken downstream action.
                                    </div>
                                </div>
                            </div>
                            <div className="action-opt selected-reject">
                                <div className="ao-radio reject-active" />
                                <div className="ao-content">
                                    <div className="ao-label reject">Reject</div>
                                    <div className="ao-desc">
                                        Dispute or refuse this commit. A signed rejection with your
                                        reason will be anchored on Stellar. The thread enters a
                                        disputed state.
                                    </div>
                                    <div className="rejection-reason">
                                        <div className="rr-label">
                                            Reason code <span style={{ color: "#EF4444" }}>*</span>
                                        </div>
                                        <select className="rr-code-select">
                                            <option value="">Select a reason code…</option>
                                            <option value="policy_violation" selected>
                                                policy_violation — commit violates an active policy rule
                                            </option>
                                            <option value="decrypt_failed">
                                                decrypt_failed — value could not be decrypted by this
                                                vault
                                            </option>
                                            <option value="schema_invalid">
                                                schema_invalid — commit structure does not match channel
                                                schema
                                            </option>
                                            <option value="vault_suspended">
                                                vault_suspended — originating vault is suspended or
                                                revoked
                                            </option>
                                        </select>
                                        <div className="rr-code-hint">
                                            Machine-actionable. Used by AFP to route and classify this
                                            rejection.
                                        </div>
                                        <div className="rr-label" style={{ marginTop: 0 }}>
                                            Rejection detail{" "}
                                            <span
                                                style={{
                                                    color: "#999",
                                                    fontWeight: 400,
                                                    textTransform: "none",
                                                    letterSpacing: 0
                                                }}
                                            >
                                                (optional)
                                            </span>
                                        </div>
                                        <textarea
                                            className="rr-input"
                                            placeholder="Add a human-readable explanation for the record…"
                                            defaultValue={
                                                "Budget envelope reference is incorrect — Q3 cap is $200,000. Contract value $240,000 exceeds approved limit. Requires re-authorization from CFO before clearance can be accepted."
                                            }
                                        />
                                        <div className="rr-note">
                                            Both fields will be signed by vault_direction, attached to
                                            the rejection receipt, and anchored on Stellar as an
                                            immutable record.
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        {/* Stellar anchor confirmation */}
                        <div className="stellar-row">
                            <span className="sr-icon">✦</span>
                            <span className="sr-text">
                                This acknowledgment will be signed by{" "}
                                <strong>vault_direction</strong> and anchored on Stellar — a
                                tamper-proof receipt record.
                            </span>
                            <div className="sr-dot" />
                        </div>
                    </div>
                    <div style={{ padding: "0 24px 12px" }}>
                        <div className="irreversible-warn">
                            <span className="iw-icon">⚠</span>
                            <span className="iw-text">
                                <strong>This action is irreversible.</strong> Once submitted, this
                                rejection is write-once — it cannot be edited or retracted. A new
                                receipt will require vault_finance to resubmit the commit.
                            </span>
                        </div>
                    </div>
                    <div className="modal-footer">
                        <button className="btn btn-ghost">Cancel</button>
                        <button className="btn-reject">✕ Submit Rejection</button>
                    </div>
                </div>
        
    );
}

