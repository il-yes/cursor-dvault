import React, { useMemo, useState } from 'react';
import { ActivationChannel } from './ui/channel/ChannelActivation';
import { Link } from 'react-router-dom';
import { ChannelRow } from './domain/channel/channel.types';


export type ChannelSummary = {
    channelId: string;

    status: 'active' | 'pending' | 'dispute';

    relationshipType:
    | 'contract'
    | 'supplier_payment'
    | 'audit'
    | 'procurement'
    | 'payroll'
    | 'governance';

    title: string;

    participants: FlowStep[];

    assetCount: number;

    pendingEvents: number;

    lastEventType: string;

    lastActivity: string;

    stellarAnchor: string;

    federationStatus:
    | 'internal'
    | 'linked'
    | 'active';

    federationLabel: string;
};




const ChannelRowView = ({
    row,
    onClick
}: {
    row: ChannelRow,
    onClick: (id: string) => void;
}) => (

    <tr onClick={() => { onClick(row.id) }}>

        <td>
            <span
                className={`sdot s-${row.status}`}
            />
        </td>


        <td style={{
            cursor: 'pointer'
        }}>

            <div className={`th-line1 ${row.id === 'budget' ? 'new-thread' : ''}`}>
                <span className="th-type">
                    {row.type}
                </span>

                {row.title}

            </div>


            <div className="th-line2">

                {row.subtitle}

            </div>

        </td>


        <td>

            <div className="flow">

                {
                    row.participants.map(
                        (step, index) => (
                            <span key={index}>

                                {index > 0 &&
                                    <span className="fa">
                                        →
                                    </span>
                                }


                                <span
                                    className="vb"
                                    style={{
                                        background: step.color
                                    }}
                                >
                                    {step.label}
                                </span>


                            </span>
                        )
                    )
                }

            </div>


        </td>


        <td>

            <div className="asset-box">

                <span>
                    {row.assetCount}
                </span>

                assets

            </div>


            <div className="event">

                {row.lastEvent}

            </div>

        </td>



        <td className="ts">

            {row.lastActivity}

        </td>



        <td>

            <span className="stellar-val">

                {row.stellarTx}

            </span>


            <div className="row-hover-actions">

                <button className="rha-btn" onClick={(e) => {

                    e.stopPropagation();

                    onClick(row.id);

                }}>
                    Open
                </button>


                <button className="rha-btn">
                    Export
                </button>


            </div>

        </td>



        <td style={{ textAlign: 'center' }}>


            <button
                className={
                    `c3b ${row.c3Status === 'internal'
                        ? 'c3-internal'
                        :
                        row.c3Status === 'linked'
                            ? 'c3-linked'
                            :
                            'c3-active'
                    }`
                }

                title={row.c3Status}

            >

                {row.c3Label}

            </button>


        </td>


    </tr>

);

type FlowStep = { label: string; color: string };

export type LedgerRow = {
    id: string;
    status: 'ok' | 'pend' | 'dispute';
    type: string;
    title: string;
    subtitle: string;
    flow: FlowStep[];
    progress: ('done' | 'wait' | 'reject')[];
    lastActivity: string;
    stellarTx: string;
    c3Status: 'extendable' | 'linked' | 'active';
    c3Label: string;
};

export const rows: LedgerRow[] = [
    {
        id: 'inv',
        status: 'ok',
        type: 'Invoice',
        title: 'Invoice Processing — Cipla India',
        subtitle: 'vendor: Cipla · ref: INV-2047',
        flow: [
            { label: 'O', color: '#DC2626' },
            { label: 'F', color: '#2563EB' },
            { label: 'T', color: '#059669' },
        ],
        progress: ['done', 'done', 'done'],
        lastActivity: '3 days ago',
        stellarTx: 'tx_a3f1c2…',
        c3Status: 'extendable',
        c3Label: '⛓+',
    },
    {
        id: 'budget',
        status: 'pend',
        type: 'Budget',
        title: 'Budget Allocation Q3',
        subtitle: 'period: Q3 2026 · scope: all depts',
        flow: [
            { label: 'All', color: '#888888' },
            { label: 'Dir', color: '#444444' },
        ],
        progress: ['done', 'wait'],
        lastActivity: '1 day ago',
        stellarTx: 'tx_b8c2d1…',
        c3Status: 'extendable',
        c3Label: '⛓+',
    },
    {
        id: 'payroll',
        status: 'ok',
        type: 'Payroll',
        title: 'Payroll Run — June',
        subtitle: 'period: June 2026 · employees: 47',
        flow: [
            { label: 'HR', color: '#D97706' },
            { label: 'F', color: '#2563EB' },
            { label: 'T', color: '#059669' },
        ],
        progress: ['done', 'done', 'done'],
        lastActivity: '5 days ago',
        stellarTx: 'tx_e7d4f9…',
        c3Status: 'extendable',
        c3Label: '⛓+',
    },
    {
        id: 'contract',
        status: 'pend',
        type: 'Contract',
        title: 'Contract — Supplier X',
        subtitle: 'vendor: Supplier X · value: $240,000',
        flow: [
            { label: 'L', color: '#7C3AED' },
            { label: 'F', color: '#2563EB' },
            { label: 'Dir', color: '#444444' },
        ],
        progress: ['done', 'done', 'wait'],
        lastActivity: '2 days ago',
        stellarTx: 'tx_f2a9e3…',
        c3Status: 'linked',
        c3Label: '⛓✓',
    },
    {
        id: 'po',
        status: 'ok',
        type: 'Purchase',
        title: 'Purchase Order #047',
        subtitle: 'vendor: Office Depot · amount: $3,200',
        flow: [
            { label: 'O', color: '#DC2626' },
            { label: 'F', color: '#2563EB' },
            { label: 'Dir', color: '#444444' },
        ],
        progress: ['done', 'done', 'done'],
        lastActivity: '6 days ago',
        stellarTx: 'tx_c5b7a1…',
        c3Status: 'extendable',
        c3Label: '⛓+',
    },
    {
        id: 'report',
        status: 'pend',
        type: 'Report',
        title: 'Board Report — June',
        subtitle: 'period: June 2026 · type: board quarterly',
        flow: [
            { label: 'All', color: '#888888' },
            { label: 'Dir', color: '#444444' },
        ],
        progress: ['done', 'wait', 'wait'],
        lastActivity: '1 day ago',
        stellarTx: 'tx_d9e1b8…',
        c3Status: 'extendable',
        c3Label: '⛓+',
    },
    {
        id: 'audit',
        status: 'pend',
        type: 'Compliance',
        title: 'Compliance Audit H1',
        subtitle: 'auditor: EY · period: H1 2026',
        flow: [
            { label: 'All', color: '#888888' },
            { label: 'C', color: '#0891B2' },
        ],
        progress: ['done', 'done', 'wait'],
        lastActivity: 'Today',
        stellarTx: 'tx_g3h5k7…',
        c3Status: 'linked',
        c3Label: '⛓✓',
    },
    {
        id: 'onboard',
        status: 'ok',
        type: 'Onboarding',
        title: 'Onboarding — EMP-2047',
        subtitle: 'employee: Marie Chen · dept: Engineering',
        flow: [
            { label: 'HR', color: '#D97706' },
            { label: 'IT', color: '#4F46E5' },
            { label: 'F', color: '#2563EB' },
            { label: 'Dir', color: '#444444' },
        ],
        progress: ['done', 'done', 'done', 'done'],
        lastActivity: '1 week ago',
        stellarTx: 'tx_h8k2m4…',
        c3Status: 'extendable',
        c3Label: '⛓+',
    },
];

export const LedgerList = ({
    ledgerRow,
    newCreated,
    hasConflict,
    onOpenChannel
}: {
    ledgerRow: ChannelRow[];
    newCreated: boolean;
    hasConflict: boolean;
    onOpenChannel: (channelId: string) => void;
}) => {
    return ledgerRow.length > 0 ? <FakeList ledgerRow={ledgerRow} newCreated={newCreated} hasConflict={hasConflict} onOpenChannel={onOpenChannel} /> : <EmptyList />;
};

export const EmptyList = () => (
    <div className="ledger-area">
        <div className="ledger-topbar">
            <div className="ledger-title">LEDGER</div>
            <div className="ledger-controls">
                <div className="ctrl-btn">Filter ▾</div>
                <div className="ctrl-btn">📅 June 2026 ▾</div>
            </div>
        </div>

        <div className="empty-area">
            <div className="empty-headline">Start your first thread</div>
            <div className="empty-subtext">
                Choose a template to create your vault network and first auditable workflow.
                Every thread is cryptographically anchored on Stellar.
            </div>

            <div className="template-grid">
                {[
                    ['🧾', 'Finance', 'Invoice Processing', 'Submit, approve, and release vendor payments with a full audit trail.', '#FBF0D8', 'Ops', 'Finance', 'Treasury'],
                    ['📊', 'Finance', 'Budget Allocation', 'Coordinate quarterly budget requests across all departments with directorial sign-off.', '#EEF2FF', 'All depts', 'Direction'],
                    ['📜', 'Legal', 'Contract Execution', 'Draft, review, and bilaterally sign contracts. Extend to supplier vaults with C3.', '#F5F3FF', 'Legal', 'Finance', 'Direction'],
                    ['🔎', 'Compliance', 'Compliance Audit Close', 'Provide auditors read-only, verifiable access — no document transfers required.', '#ECFEFF', 'All', 'Compliance'],
                    ['👤', 'HR', 'Employee Onboarding', 'Coordinate HR, IT provisioning, and Finance payroll setup in one sequential thread.', '#FFFBEB', 'HR', 'IT+Finance', 'Direction'],
                ].map(([icon, label, name, desc, bg, ...flow]) => (
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
                        <button className="tc-start">Start</button>
                    </div>
                ))}
            </div>
        </div>

        <div className="ledger-footer-hint">
            0 threads · Start a template above to populate the ledger
        </div>
    </div>
);

export const FakeList = ({
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
    const visibleRows = useMemo(() => ledgerRow, []);

    return (
        <div className="ledger-area">
            <div className="ledger-topbar">
                <div className="ledger-title">LEDGER</div>
                <div className="ledger-controls">
                    <button className="ctrl-btn">Filter ▾</button>
                    <button className="ctrl-btn">📅 June 2026 ▾</button>
                </div>
            </div>

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
                                // <ChannelRowView key={row.id} row={row} onOpen={onOpenChannel} />
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

const LedgerRowView = ({
    row,
    onClick
}: {
    row: ChannelRow;
    onClick: (id: string) => void;
}) => (

    <tr
        className="ledger-row"
        onClick={() => onClick(row.id)}
    >


        <td>
            <span className={`sdot s-${row.status}`} />
        </td>


        <td>

            <div className="th-line1">

                <span className="th-type">
                    {row.type}
                </span>

                {row.title}

            </div>


            <div className="th-line2">
                {row.subtitle}
            </div>

        </td>



        <td>

            <div className="flow">

                {
                    row.participants.map((step, idx) => (

                        <span key={idx}>

                            {idx > 0 &&
                                <span className="fa">→</span>
                            }


                            <span
                                className="vb"
                                style={{
                                    background: step.color
                                }}
                            >
                                {step.label}
                            </span>

                        </span>

                    ))
                }

            </div>

        </td>



        <td>

            <div>

                <b>{row.assetCount} </b> assets

            </div>


            <div className="event">
                {row.lastEvent}
            </div>


        </td>



        <td className="ts">
            {row.lastActivity}
        </td>



        <td>

            <span className="stellar-val">
                {row.stellarTx}
            </span>


            <button

                className="rha-btn"

                onClick={(e) => {

                    e.stopPropagation();

                    onClick(row.id);

                }}

            >
                Open
            </button>


        </td>



        <td>

            <button

                className={`c3b ${row.c3Status
                    }`}

                onClick={(e) => {

                    e.stopPropagation();

                }}

            >
                {row.c3Label}

            </button>

        </td>


    </tr>

);

const DisputeRow = () => (
    <tr className="row-sel row-new-left-border">
        <td>
            <span className="sdot s-dispute" />
        </td>

        <td>
            <div>
                <span className="th-type">Contract</span>
                Contract — Supplier X <span className="new-label">Disputed</span>
            </div>
            <div className="th-line2" style={{ marginTop: 3 }}>
                vendor: Supplier X · value: $240,000
            </div>
        </td>

        <td>
            <div className="flow">
                <span className="vb" style={{ background: '#7C3AED' }}>
                    L
                </span>
                <span className="fa">→</span>
                <span className="vb" style={{ background: '#2563EB' }}>
                    F
                </span>
                <span className="fa">→</span>
                <span className="vb" style={{ background: '#444444' }}>
                    Dir
                </span>
            </div>
        </td>

        <td>
            <div className="pipeline">
                <div className="pseg pseg-done" />
                <div className="pseg pseg-reject" />
                <div className="pseg pseg-wait" />
            </div>
        </td>

        <td className="ts" style={{ color: '#DC2626' }}>
            today
        </td>

        <td>
            <span className="stellar-val">tx_f2a9…</span>
        </td>

        <td style={{ textAlign: 'center' }}>
            <button className="c3b c3-active">⛓+</button>
        </td>
    </tr>
);

