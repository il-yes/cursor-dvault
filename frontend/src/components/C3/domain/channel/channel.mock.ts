import { ChannelProperty, ChannelRow, ChannelSlot, ChannelView, Department, FlowStep } from "./channel.types";

export const channel: ChannelView = {

    id: 'contract-supplier',


    title:
        'Contract Execution Channel',


    subtitle:
        'Supplier X · $240,000 agreement',


    status: 'active',


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
            color: '#444'
        },
        {
            label: 'SUP',
            color: '#0891B2'
        }
    ],


    assets: {

        total: 3,


        items: [

            {
                id: 'msa',

                title: 'Master Service Agreement',

                type: 'contract',

                status: 'signed',

                lastEvent: 'contract.countersigned',
                channelId: "",
                subtitle: "",
                createdAt: ""
            },

            {
                id: 'pricing',

                title: 'Pricing Schedule',

                type: 'attachment',

                status: 'approved',

                lastEvent: 'finance.approved',
                channelId: "",
                subtitle: "",
                createdAt: ""
            }

        ]

    },


    activity: {


        lastEvent:
            'contract.countersigned',


        lastActivity:
            '2 hours ago',


        events: [

            {
                type: 'contract.created',
                actor: 'Legal',
                time: '12 Jun'
            },

            {
                type: 'finance.approved',
                actor: 'Finance',
                time: '12 Jun'
            },

            {
                type: 'contract.countersigned',
                actor: 'Supplier X',
                time: 'Today'
            }

        ]

    },



    policy: {

        read: [
            'Supplier X'
        ],

        write: [
            'Internal'
        ]

    },


    stellarTx:
        'tx_f2a9e3…',


    c3: {

        status: 'active',

        externalVault:
            'Supplier X'

    },

    slots: [],

    defaultProperties: []

};
export const channelInvoice: ChannelView = {

    id: 'invoice-cipla',


    title:
        'Invoice Processing Channel',


    subtitle:
        'Cipla India · INV-2047 · $18,450',


    status: 'active',


    participants: [
        {
            label: 'OPS',
            color: '#DC2626'
        },
        {
            label: 'FIN',
            color: '#2563EB'
        },
        {
            label: 'TRE',
            color: '#059669'
        },
        {
            label: 'CIP',
            color: '#0891B2'
        }
    ],


    assets: {

        total: 3,


        items: [

            {
                id: "invoice",

                channelId: "invoice-cipla",

                title:
                    "Vendor Invoice INV-2047",

                subtitle:
                    "Cipla India · $18,450",


                status:
                    "approved",


                createdAt:
                    "2026-06-09",


                participants: [
                    {
                        label: "OPS",
                        color: "#DC2626"
                    },
                    {
                        label: "FIN",
                        color: "#2563EB"
                    },
                    {
                        label: "TRE",
                        color: "#059669"
                    },
                    {
                        label: "CIP",
                        color: "#0891B2"
                    }
                ],


                lastEvent:
                    "receipt.issued",


                stellarTx:
                    "tx_invoice_a31",


                events: [

                    {
                        id: "evt_invoice_001",

                        type:
                            "invoice.created",

                        actor:
                            "Operations",

                        createdAt:
                            "2026-06-09T09:00",

                        time:
                            "Jun 09 · 09:00",

                        payloadCID:
                            "bafy_invoice_create",

                        receiptStatus:
                            "processed",

                        stellarTx:
                            "tx_create_01"
                    },


                    {
                        id: "evt_invoice_002",

                        type:
                            "finance.approved",

                        actor:
                            "Finance",

                        createdAt:
                            "2026-06-10T13:20",

                        time:
                            "Jun 10 · 13:20",

                        payloadCID:
                            "bafy_finance_ok",

                        receiptStatus:
                            "processed",

                        stellarTx:
                            "tx_finance_02"
                    },


                    {
                        id: "evt_invoice_003",

                        type:
                            "payment.released",

                        actor:
                            "Treasury",

                        createdAt:
                            "2026-06-11T08:45",

                        time:
                            "Jun 11 · 08:45",

                        payloadCID:
                            "bafy_payment",

                        receiptStatus:
                            "processed",

                        stellarTx:
                            "tx_payment_03"
                    },


                    {
                        id: "evt_invoice_004",

                        type:
                            "receipt.issued",

                        actor:
                            "Cipla",

                        createdAt:
                            "2026-06-12T10:00",

                        time:
                            "Today",

                        payloadCID:
                            "bafy_receipt",

                        receiptStatus:
                            "received"
                    }

                ]
            },

            {
                id: 'payment-order',

                title: 'Treasury Payment Order',

                type: 'payment',

                status: 'processing',

                lastEvent: 'payment.released',
                channelId: "",
                subtitle: "",
                createdAt: ""
            },

            {
                id: 'receipt',

                title: 'Supplier Receipt',

                type: 'attestation',

                status: 'received',

                lastEvent: 'receipt.issued',
                channelId: "",
                subtitle: "",
                createdAt: ""
            }

        ]

    },


    activity: {

        lastEvent:
            'receipt.issued',


        lastActivity:
            '3 days ago',


        events: [

            {
                type: 'invoice.created',
                actor: 'Operations',
                time: '09 Jun'
            },

            {
                type: 'finance.approved',
                actor: 'Finance',
                time: '10 Jun'
            },

            {
                type: 'payment.released',
                actor: 'Treasury',
                time: '11 Jun'
            },

            {
                type: 'receipt.issued',
                actor: 'Cipla',
                time: 'Today'
            }

        ]

    },


    policy: {

        read: [
            'Cipla India',
            'Finance'
        ],

        write: [
            'Treasury'
        ]

    },


    stellarTx:
        'tx_a3f1c2…',


    c3: {

        status: 'linked',

        externalVault:
            'Cipla India'

    },


    slots: [

        {
            id: "draft",
            name: "contract_draft",
            role: "Author",
            gated: false
        },

        {
            id: "finance",
            name: "financial_clearance",
            role: "Reviewer",
            gated: true
        },

        {
            id: "signature",
            name: "executive_signature",
            role: "Approver",
            gated: true
        }

    ],

    defaultProperties: [

        {
            key: "counterparty",
            required: true
        },

        {
            key: "contract_value",
            required: true
        }

    ]

};
export const channelBudget: ChannelView = {

    id: 'budget-q3',


    title:
        'Budget Allocation Q3',


    subtitle:
        'period: Q3 2026 · scope: all departments',


    status: 'pending',


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


    assets: {

        total: 2,


        items: [

            {
                id: 'budget-plan',

                title: 'FY26 Budget Proposal',

                type: 'planning',

                status: 'draft',

                lastEvent: 'budget.created',
                channelId: "",
                subtitle: "",
                createdAt: ""
            },

            {
                id: 'budget-review',

                title: 'Departmental Budget Review',

                type: 'attachment',

                status: 'draft',

                lastEvent: 'budget.review-requested',
                channelId: "",
                subtitle: "",
                createdAt: ""
            }

        ]

    },


    activity: {

        lastEvent:
            'budget.review-requested',


        lastActivity:
            '1 day ago',


        events: [

            {
                type: 'budget.created',
                actor: 'Finance',
                time: '01 Jun'
            },

            {
                type: 'budget.review-requested',
                actor: 'Finance',
                time: 'Today'
            }

        ]

    },



    policy: {

        read: [
            'Internal'
        ],

        write: [
            'Finance',
            'Board'
        ]

    },


    stellarTx:
        'tx_b8c2d1…',


    c3: {

        status: 'internal',

        externalVault:
            null

    },

    slots: [],

    defaultProperties: []

};
export const channelAudit: ChannelView = {

    id: 'compliance-audit-h1',


    title:
        'Compliance Audit Channel',


    subtitle:
        'External Auditor · H1 2026',


    status: 'active',


    participants: [
        {
            label: 'ALL',
            color: '#888'
        },
        {
            label: 'CMP',
            color: '#7C3AED'
        },
        {
            label: 'AUD',
            color: '#0891B2'
        }
    ],


    assets: {

        total: 4,


        items: [

            {
                id: 'attestations',

                title: 'Department Attestations',

                type: 'attestation',

                status: 'approved',

                lastEvent: 'attestation.completed',
                channelId: "",
                subtitle: "",
                createdAt: ""
            },

            {
                id: 'summary',

                title: 'Compliance Summary',

                type: 'report',

                status: 'pending',

                lastEvent: 'audit.review',
                channelId: "",
                subtitle: "",
                createdAt: ""
            },

            {
                id: 'findings',

                title: 'Audit Findings',

                type: 'report',

                status: 'open',

                lastEvent: 'finding.created',
                channelId: "",
                subtitle: "",
                createdAt: ""
            },

            {
                id: 'closure',

                title: 'Audit Closure',

                type: 'attestation',

                status: 'waiting',

                lastEvent: 'audit.close.pending',
                channelId: "",
                subtitle: "",
                createdAt: ""
            }

        ]

    },


    activity: {

        lastEvent:
            'audit.review',


        lastActivity:
            'Today',


        events: [

            {
                type: 'attestation.completed',
                actor: 'Departments',
                time: '05 Jun'
            },

            {
                type: 'audit.review',
                actor: 'Compliance',
                time: 'Today'
            }

        ]

    },


    policy: {

        read: [
            'External Auditor'
        ],

        write: [
            'Compliance'
        ]

    },


    stellarTx:
        'tx_g3h5k7…',


    c3: {

        status: 'active',

        externalVault:
            'Auditor Vault'

    },

    slots: [],

    defaultProperties: []

};
export const channels: ChannelView[] = [
    channel,
    channelInvoice,
    channelBudget,
    channelAudit,
]

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

export interface ChannelTemplate {

    id: string;

    title: string;

    subtitle: string;

    description: string;

    participants: FlowStep[];

    slots: ChannelSlot[];

    defaultProperties: ChannelProperty[];

}

export const ContractExecutionTemplate: ChannelTemplate = {

    id: "contract-execution",

    title: "Contract Execution",

    subtitle: "Legal → Finance → Direction",

    description: "",

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

    slots: [

        {
            id: "draft",
            name: "contract_draft",
            role: "Author",
            gated: false
        },

        {
            id: "finance",
            name: "financial_clearance",
            role: "Reviewer",
            gated: true
        },

        {
            id: "signature",
            name: "executive_signature",
            role: "Approver",
            gated: true
        }

    ],

    defaultProperties: [
        {
            key: "counterparty",
            required: true
        },
        {
            key: "contract_value",
            required: true
        }
    ]

}

export const templates = [
    {

        id: "contract-execution",

        title: "Contract Execution",

        subtitle: "Legal → Finance → Direction",

        description: "",

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

        slots: [

            {
                id: "draft",
                name: "contract_draft",
                role: "Author",
                gated: false,
                vault: "vault_legal"
            },

            {
                id: "finance",
                name: "financial_clearance",
                role: "Reviewer",
                gated: true,
                vault: "vault_finance"
            },

            {
                id: "signature",
                name: "executive_signature",
                role: "Approver",
                gated: true,
                vault: "vault_direction"
            }

        ],

        defaultProperties: [
            {
                key: "counterparty",
                required: true
            },
            {
                key: "contract_value",
                required: true
            }
        ]

    }
]


export const departments: Department[] = [
    {
        id: "vault_legal",
        name: "Legal",
        color: "#7C3AED"
    },
    {
        id: "vault_finance",
        name: "Finance",
        color: "#2563EB"
    },
    {
        id: "vault_direction",
        name: "Direction",
        color: "#444"
    },
    {
        id: "vault_treasury",
        name: "Treasury",
        color: "#1d5c09ff"
    },
    {
        id: "vault_hr",
        name: "HR",
        color: "#D97706"
    },
    {
        id: "vault_it",
        name: "IT",
        color: "#4F46E5"
    },
    {
        id: "vault_ops",
        name: "Ops",
        color: "#d41006ff"
    },
    {
        id: "vault_compliance",
        name: "Compliance",
        color: "#0bf5e5ff"
    }
]
