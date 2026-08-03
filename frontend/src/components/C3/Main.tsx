import React, { useState, useEffect } from "react";
import { useVaultStore } from "@/store/vaultStore";
import { NewShareModal } from "@/components/NewCryptoShareModal";
import { LedgerList, LedgerRow } from "./ledger_list";
import { ChannelView as ChannelViewType } from "./channel-view";
import { ReceiptAckModal } from "./receipt_ack_modal";
import { AddThreadSlidingPanel } from "./add_thread";
import { Step1 } from "./channel-creation-steps/step-1.name";
import { Step2 } from "./channel-creation-steps/step-2.configure";
import { Step3 } from "./channel-creation-steps/step-3.vaults";
import { C3Step4 } from "./channel-creation-steps/step-4.activate";
import ThreadDetailSlidingPanel from "./thread_detail";
import { ChannelView } from "./channel-view";
import { ChannelRow } from "./ui/ledger/LedgerRow";


const rows: LedgerRow[] = [
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
export const channelRows: ChannelRow[] = [

    {
        id: 'invoice-cipla',

        status: 'active',

        type: 'Payment',

        title:
            'Invoice Processing — Cipla India',

        subtitle:
            'vendor: Cipla · ref: INV-2047',


        participants: [
            {
                label: 'O',
                color: '#DC2626'
            },
            {
                label: 'F',
                color: '#2563EB'
            },
            {
                label: 'T',
                color: '#059669'
            },
        ],


        assetCount: 2,

        lastEvent:
            'payment.released',

        lastActivity:
            '3 days ago',

        stellarTx:
            'tx_a3f1c2…',


        c3Status:
            'linked',

        c3Label:
            '⛓✓'
    },



    {
        id: 'budget-q3',

        status: 'pending',

        type: 'Governance',

        title:
            'Budget Allocation Q3',

        subtitle:
            'period: Q3 2026 · scope: all departments',


        participants: [
            {
                label: 'ALL',
                color: '#888888'
            },
            {
                label: 'DIR',
                color: '#444444'
            }
        ],


        assetCount: 1,

        lastEvent:
            'budget.awaiting-approval',

        lastActivity:
            '1 day ago',

        stellarTx:
            'tx_b8c2d1…',


        c3Status:
            'internal',

        c3Label:
            '⛓+'
    },



    {
        id: 'payroll-june',

        status: 'active',

        type: 'Payroll',

        title:
            'Payroll Run — June',

        subtitle:
            'employees: 47 · period: June 2026',


        participants: [
            {
                label: 'HR',
                color: '#D97706'
            },
            {
                label: 'F',
                color: '#2563EB'
            },
            {
                label: 'T',
                color: '#059669'
            }
        ],


        assetCount: 3,


        lastEvent:
            'payroll.completed',


        lastActivity:
            '5 days ago',


        stellarTx:
            'tx_e7d4f9…',


        c3Status:
            'linked',

        c3Label:
            '⛓✓'
    },



    {
        id: 'contract-supplier',


        status: 'pending',


        type: 'Contract',


        title:
            'Contract Execution — Supplier X',


        subtitle:
            'supplier: Supplier X · value: $240,000',



        participants: [
            {
                label: 'L',
                color: '#7C3AED'
            },
            {
                label: 'F',
                color: '#2563EB'
            },
            {
                label: 'DIR',
                color: '#444444'
            },
            {
                label: 'SUP',
                color: '#0891B2'
            }
        ],



        assetCount: 4,


        lastEvent:
            'contract.countersigned',


        lastActivity:
            'Today',


        stellarTx:
            'tx_f2a9e3…',


        c3Status:
            'active',


        c3Label:
            '⛓●'
    },



    {
        id: 'purchase-order-047',

        status: 'active',

        type: 'Procurement',

        title:
            'Purchase Order #047',

        subtitle:
            'vendor: Office Depot · amount: $3,200',


        participants: [
            {
                label: 'O',
                color: '#DC2626'
            },
            {
                label: 'F',
                color: '#2563EB'
            },
            {
                label: 'DIR',
                color: '#444444'
            }
        ],


        assetCount: 2,


        lastEvent:
            'purchase.approved',


        lastActivity:
            '6 days ago',


        stellarTx:
            'tx_c5b7a1…',


        c3Status:
            'internal',

        c3Label:
            '⛓+'
    },



    {
        id: 'audit-h1',


        status: 'pending',


        type: 'Compliance',


        title:
            'Compliance Audit H1',


        subtitle:
            'auditor: EY · period: H1 2026',



        participants: [
            {
                label: 'ALL',
                color: '#888888'
            },
            {
                label: 'AUD',
                color: '#0891B2'
            }
        ],



        assetCount: 8,


        lastEvent:
            'attestation.waiting',


        lastActivity:
            'Today',


        stellarTx:
            'tx_g3h5k7…',


        c3Status:
            'active',


        c3Label:
            '⛓●'
    },



    {
        id: 'employee-onboarding',


        status: 'active',


        type: 'Onboarding',


        title:
            'Onboarding — EMP-2047',


        subtitle:
            'employee: Marie Chen · Engineering',



        participants: [
            {
                label: 'HR',
                color: '#D97706'
            },
            {
                label: 'IT',
                color: '#4F46E5'
            },
            {
                label: 'F',
                color: '#2563EB'
            },
            {
                label: 'DIR',
                color: '#444444'
            }
        ],


        assetCount: 5,


        lastEvent:
            'onboarding.completed',


        lastActivity:
            '1 week ago',


        stellarTx:
            'tx_h8k2m4…',


        c3Status:
            'linked',


        c3Label:
            '⛓✓'
    }

];


export const channels: ChannelViewType[] = [
    {
        id: 'supplier-payment',

        status: 'active',

        title: 'Supplier Payment Channel — Cipla',

        subtitle:
            'Invoice settlement relationship · vendor: Cipla India',

        participants: [
            { label: 'O', color: '#DC2626' },
            { label: 'F', color: '#2563EB' },
            { label: 'T', color: '#059669' },
            { label: 'C', color: '#0891B2' },
        ],

        assets: {
            total: 4,
            items: [],
        },

        activity: {
            lastEvent: 'payment.receipt.signed',
            lastActivity: '3 days ago',
            events: [],
        },

        policy: {
            read: ['Cipla India'],
            write: ['Internal'],
        },

        stellarTx:
            'tx_a3f1c2…',

        c3: {
            status: 'active',
        },
    },


    {
        id: 'budget-governance',

        status: 'pending',

        title:
            'Budget Governance Channel — Q3',

        subtitle:
            'department budgets · finance validation',

        participants: [
            { label: 'ALL', color: '#888888' },
            { label: 'F', color: '#2563EB' },
            { label: 'DIR', color: '#444444' },
        ],

        assets: {
            total: 9,
            items: [],
        },

        activity: {
            lastEvent: 'budget.approval.requested',
            lastActivity: '1 day ago',
            events: [],
        },

        policy: {
            read: ['All Departments'],
            write: ['Finance', 'Direction'],
        },

        stellarTx:
            'tx_b8c2d1…',

        c3: {
            status: 'internal',
        },
    },


    {
        id: 'payroll',

        status:
            'active',

        title:
            'Payroll Distribution Channel',

        subtitle:
            'June 2026 payroll · employees: 47',

        participants: [
            { label: 'HR', color: '#D97706' },
            { label: 'F', color: '#2563EB' },
            { label: 'T', color: '#059669' },
        ],

        assets: {
            total: 3,
            items: [],
        },

        activity: {
            lastEvent: 'payroll.release.completed',
            lastActivity: '5 days ago',
            events: [],
        },

        policy: {
            read: ['HR', 'Finance'],
            write: ['Finance', 'Treasury'],
        },

        stellarTx:
            'tx_e7d4f9…',

        c3: {
            status: 'linked',
        },
    },


    {
        id: 'contract-supplier',

        status:
            'pending',

        title:
            'Contract Execution Channel — Supplier X',

        subtitle:
            '$240,000 agreement · counterparty connected',

        participants: [
            { label: 'L', color: '#7C3AED' },
            { label: 'F', color: '#2563EB' },
            { label: 'DIR', color: '#444444' },
            { label: 'SUP', color: '#0891B2' },
        ],

        assets: {
            total: 3,
            items: [],
        },

        activity: {
            lastEvent: 'signature.requested',
            lastActivity: '2 days ago',
            events: [],
        },

        policy: {
            read: ['Supplier X'],
            write: ['Internal'],
        },

        stellarTx:
            'tx_f2a9e3…',

        c3: {
            status: 'active',
            externalVault: 'Supplier X',
        },
    },


    {
        id: 'procurement',

        status:
            'active',

        title:
            'Procurement Channel — Office Depot',

        subtitle:
            'Purchase order lifecycle · $3,200',

        participants: [
            { label: 'O', color: '#DC2626' },
            { label: 'F', color: '#2563EB' },
            { label: 'DIR', color: '#444444' },
        ],

        assets: {
            total: 5,
            items: [],
        },

        activity: {
            lastEvent: 'purchase.approved',
            lastActivity: '6 days ago',
            events: [],
        },

        policy: {
            read: ['Operations'],
            write: ['Finance', 'Direction'],
        },

        stellarTx:
            'tx_c5b7a1…',

        c3: {
            status: 'internal',
        },
    },


    {
        id: 'board-report',

        status:
            'pending',

        title:
            'Board Reporting Channel',

        subtitle:
            'quarterly reporting package · June 2026',

        participants: [
            { label: 'ALL', color: '#888888' },
            { label: 'DIR', color: '#444444' },
        ],

        assets: {
            total: 14,
            items: [],
        },

        activity: {
            lastEvent: 'report.review.pending',
            lastActivity: 'Today',
            events: [],
        },

        policy: {
            read: ['Board', 'Direction'],
            write: ['Internal'],
        },

        stellarTx:
            'tx_d9e1b8…',

        c3: {
            status: 'linked',
        },
    },


    {
        id: 'audit',

        status:
            'pending',

        title:
            'Compliance Audit Channel — EY',

        subtitle:
            'external auditor observer access',

        participants: [
            { label: 'ALL', color: '#888888' },
            { label: 'C', color: '#0891B2' },
            { label: 'EY', color: '#7C3AED' },
        ],

        assets: {
            total: 22,
            items: [],
        },

        activity: {
            lastEvent: 'attestation.received',
            lastActivity: 'Today',
            events: [],
        },

        policy: {
            read: ['EY (Observer)'],
            write: ['Internal'],
        },

        stellarTx:
            'tx_g3h5k7…',

        c3: {
            status: 'active',
            externalVault: 'EY',
        },
    },


    {
        id: 'employee-onboarding',

        status:
            'active',

        title:
            'Employee Onboarding Channel',

        subtitle:
            'Marie Chen · Engineering',

        participants: [
            { label: 'HR', color: '#D97706' },
            { label: 'IT', color: '#4F46E5' },
            { label: 'F', color: '#2563EB' },
            { label: 'DIR', color: '#444444' },
        ],

        assets: {
            total: 6,
            items: [],
        },

        activity: {
            lastEvent: 'employee.provisioned',
            lastActivity: '1 week ago',
            events: [],
        },

        policy: {
            read: ['HR', 'IT', 'Finance'],
            write: ['HR'],
        },

        stellarTx:
            'tx_h8k2m4…',

        c3: {
            status: 'internal',
        },
    },
];

export default function C3App() {
    const [isNewShareOpen, setIsNewShareOpen] = useState(false);
    const [sharedEntriesRefreshKey, setSharedEntriesRefreshKey] = useState(0);
    const hasConflict = false;
    type Screen =
        | { type: "ledger" }
        | { type: "channel"; channelId: string }
        | { type: "asset"; assetId: string };

    const [screen, setScreen] = useState<Screen>({ type: "ledger" });

    const renderView = () => {
        switch (screen.type) {

            case "ledger":

                return (

                    <LedgerList
                        ledgerRow={channelRows}
                        newCreated={isNewShareOpen}
                        hasConflict={hasConflict}
                        onOpenChannel={(id) =>
                            setScreen({
                                type: "channel",
                                channelId: id,
                            })
                        }
                    />

                );

            case "channel":

                return (

                    <ChannelView
                        channel={channels.find(c => c.id === screen.channelId)!}
                    />

                );

        }
    }

    return (
        <>
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
                {/* <LedgerList ledgerRow={channels as any} newCreated={isNewShareOpen} hasConflict={hasConflict} onOpenChannel={(channelId) => setChannelId("2")} /> */}
                {renderView()}
                {/* <ThreadDetailSlidingPanel hasConflict={hasConflict} /> */}
                {/* <AddThreadSlidingPanel /> */}
            </div>

            <div className="modal-wrap">
                {/* <Step1 /> */}
                {/* <Step2 /> */}
                {/* <Step3 /> */}
                {/* <C3Step4 /> */}
                {/* <ReceiptAckModal /> */}
            </div>
        </>
    )
}



