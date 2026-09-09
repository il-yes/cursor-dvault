import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ActivationChannel } from "../channel/ChannelActivation";
import { ChannelRow } from "../../domain/channel/channel.types";
import { DisputeRow, LedgerRowView } from "./LedgerRow";
import { useC3DialogStore } from "../../infrastructure/store/c3DialogStore";
import { InviteModal } from "./InviteModal";
import * as ROUTES from "@/constants/routes";


export const LedgerTable = ({
    ledgerRow,
    newCreated,
    hasConflict,
    onOpenChannel,
}: {
    ledgerRow: ChannelRow[];
    newCreated: boolean;
    hasConflict: boolean;
    onOpenChannel: (channelId: string) => void;
}) => {
    const { openC3CreateDialog } = useC3DialogStore();

    return ledgerRow.length > 0 
            ? <List ledgerRow={ledgerRow} newCreated={newCreated} hasConflict={hasConflict} onOpenChannel={onOpenChannel} /> 
            : <EmptyList openC3CreateDialog={openC3CreateDialog} />;
};

export const List = ({
    ledgerRow,
    newCreated,
    hasConflict,
    onOpenChannel,
}: {
    ledgerRow: ChannelRow[];
    newCreated: boolean;
    hasConflict: boolean;
    onOpenChannel: (channelId: string) => void;
}) => {
    const navigate = useNavigate();
    const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
    const visibleRows = useMemo(() => ledgerRow, []);

    return (
        <div className="ledger-area">
            <div className="ledger-topbar">
                <div className="ledger-title" style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                    <span>LEDGER</span>
                    <button
                        title="Invite Member to Workspace / Channel"
                        style={{
                            cursor: "pointer",
                            fontSize: "11px",
                            padding: "3px 8px",
                            borderRadius: "4px",
                            background: "#10b981",
                            border: "1px solid #059669",
                            color: "#ffffff",
                            fontWeight: 600,
                            display: "inline-flex",
                            alignItems: "center",
                            gap: "4px",
                        }}
                        onClick={() => setIsInviteModalOpen(true)}
                    >
                        ✉ Invite
                    </button>
                    <button
                        title="Workspace / Ledger Federation"
                        style={{
                            cursor: "pointer",
                            fontSize: "11px",
                            padding: "3px 8px",
                            borderRadius: "4px",
                            background: "#fafafa",
                            border: "1px solid #e0e0e0",
                            color: "#555",
                            fontWeight: 600,
                            display: "inline-flex",
                            alignItems: "center",
                            gap: "4px",
                        }}
                        onClick={() => navigate(ROUTES.FEDERATION)}
                    >
                        ⚙ Federation
                    </button>
                </div>
                <div className="ledger-controls">
                    <button className="ctrl-btn">Filter ▾</button>
                    <button className="ctrl-btn">📅 June 2026 ▾</button>
                </div>
            </div>

            <InviteModal
                isOpen={isInviteModalOpen}
                onClose={() => setIsInviteModalOpen(false)}
            />

            <div className="table-wrap">
                <table>
                    <thead>
                        <tr>
                            <th style={{ width: '2%' }} />
                            <th style={{ width: '30%' }}>Channels</th>
                            <th style={{ width: '27%' }}>Flow</th>
                            <th style={{ width: '14%' }}>Threads</th>
                            <th style={{ width: '10%' }}>Last activity</th>
                            <th style={{ width: '12%' }}>Stellar</th>
                            <th style={{ width: '5%', textAlign: 'center' }}>C3</th>
                        </tr>
                    </thead>

                    <tbody>
                        {newCreated && <ActivationChannel />}
                        {visibleRows.map((row: ChannelRow) =>
                            row.id === 'contract' && hasConflict ? (
                                <DisputeRow key={row.id} />
                            ) : (
                                <LedgerRowView key={row.id} row={row} onClick={onOpenChannel} />
                            ),
                        )}
                    </tbody>
                </table>
            </div>

            <div className="ledger-footer">
                8 threads · 4 complete · 4 pending · Period: June 2026
            </div>
        </div>
    );
};


export const EmptyList = ({ openC3CreateDialog }: { openC3CreateDialog: (open: boolean, channelId?: string) => void }) => {
    const navigate = useNavigate();
    const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
    return (
    <div className="ledger-area">
        <div className="ledger-topbar">
            <div className="ledger-title" style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                <span>LEDGER</span>
                <button
                    title="Invite Member to Workspace / Channel"
                    style={{
                        cursor: "pointer",
                        fontSize: "11px",
                        padding: "3px 8px",
                        borderRadius: "4px",
                        background: "#10b981",
                        border: "1px solid #059669",
                        color: "#ffffff",
                        fontWeight: 600,
                        display: "inline-flex",
                        alignItems: "center",
                        gap: "4px",
                    }}
                    onClick={() => setIsInviteModalOpen(true)}
                >
                    ✉ Invite
                </button>
                <button
                    title="Workspace / Ledger Federation"
                    style={{
                        cursor: "pointer",
                        fontSize: "11px",
                        padding: "3px 8px",
                        borderRadius: "4px",
                        background: "#fafafa",
                        border: "1px solid #e0e0e0",
                        color: "#555",
                        fontWeight: 600,
                        display: "inline-flex",
                        alignItems: "center",
                        gap: "4px",
                    }}
                    onClick={() => navigate(ROUTES.FEDERATION)}
                >
                    ⚙ Federation
                </button>
            </div>
            <div className="ledger-controls">
                <div className="ctrl-btn">Filter ▾</div>
                <div className="ctrl-btn">📅 June 2026 ▾</div>
            </div>
        </div>

        <InviteModal
            isOpen={isInviteModalOpen}
            onClose={() => setIsInviteModalOpen(false)}
        />

        <div className="empty-area">
            <div className="empty-headline">Start your first channel</div>
            <div className="empty-subtext">
                Choose a template to create your vault network and first auditable workflow.
                Every thread is cryptographically anchored on Stellar.
            </div>

            <div className="template-grid">
                {[
                    ['', '🧾', 'Finance', 'Invoice Processing', 'Submit, approve, and release vendor payments with a full audit trail.', '#FBF0D8', 'Ops', 'Finance', 'Treasury'],
                    ['', '📊', 'Finance', 'Budget Allocation', 'Coordinate quarterly budget requests across all departments with directorial sign-off.', '#EEF2FF', 'All depts', 'Direction'],
                    ['contract-execution', '📜', 'Legal', 'Contract Execution', 'Draft, review, and bilaterally sign contracts. Extend to supplier vaults with C3.', '#F5F3FF', 'Legal', 'Finance', 'Direction'],
                    ['', '🔎', 'Compliance', 'Compliance Audit Close', 'Provide auditors read-only, verifiable access — no document transfers required.', '#ECFEFF', 'All', 'Compliance'],
                    ['', '👤', 'HR', 'Employee Onboarding', 'Coordinate HR, IT provisioning, and Finance payroll setup in one sequential thread.', '#FFFBEB', 'HR', 'IT+Finance', 'Direction'],
                ].map(([id, icon, label, name, desc, bg, ...flow]) => (
                    <div className="template-card" key={name as string}>
                        <div className="tc-icon" style={{ background: bg as string }}>
                            {icon as string}
                        </div>
                        <div className="tc-type-label">{label as string}</div>
                        <div className="tc-name">{name as string}</div>
                        <div className="tc-desc">{desc as string}</div>
                        <div className="tc-flow">
                            {flow.map((item, idx) => (
                                <span key={idx} className={idx % 2 === 0 ? 'tc-vault' : 'tc-arrow'}>
                                    {item as string}
                                </span>
                            ))}
                        </div>
                        <button className="tc-start" onClick={() => openC3CreateDialog(true, id as string)}>Start</button>
                    </div>
                ))}
            </div>
        </div>

        <div className="ledger-footer-hint">
            0 threads · Start a template above to populate the ledger
        </div>
    </div>
  );
};
