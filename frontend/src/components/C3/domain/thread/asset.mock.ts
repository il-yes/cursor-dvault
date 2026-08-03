import { ThreadAssetViewInterface, ThreadEventView } from "./asset.types";


export const threadAssets1: ThreadAssetViewInterface[] = [  
    {
        id: 'invoice',

        title: 'Vendor Invoice INV-2047',

        type: 'invoice',

        status: 'approved',

        lastEvent: 'invoice.approved',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'payment-order',

        title: 'Treasury Payment Order',

        type: 'payment',

        status: 'processing',

        lastEvent: 'payment.released',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'receipt',

        title:
            'Supplier Receipt',

        type: 'attestation',

        status: 'received',

        lastEvent: 'receipt.issued',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'budget-plan',

        title:
            'FY26 Budget Proposal',

        type: 'planning',

        status: 'draft',

        lastEvent: 'budget.created',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'budget-review',

        title:
            'Departmental Budget Review',

        type: 'attachment',

        status: 'draft',

        lastEvent: 'budget.review-requested',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },


    {
        id: 'attestations',

        title:
            'Department Attestations',

        type: 'attestation',

        status: 'approved',

        lastEvent: 'attestation.completed',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'summary',

        title:
            'Compliance Summary',

        type: 'report',

        status: 'pending',

        lastEvent: 'audit.review',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'findings',

        title:
            'Audit Findings',

        type: 'report',

        status: 'open',

        lastEvent: 'finding.created',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    },

    {
        id: 'closure',

        title:
            'Audit Closure',

        type: 'attestation',

        status: 'waiting',

        lastEvent: 'audit.close.pending',
        channelId: '',
        subtitle: '',
        createdAt: '',
        participants: [],
        lifecycle: [],
        events: [],
        payloads: [],
        receipts: [],
        policy: { read: [], write: [] },
        stellarTx: "",
        relationship: { status: "internal" },
        externalVault: "",
        direction: "",
        role: ""
    }
]

export const threadAssetEvents: ThreadEventView[] = [
    {
        id: "evt_001",

        type: "asset.created",

        actor: "vault_legal",

        createdAt: "2026-06-10T09:14:00Z",

        time: "Jun 10 · 09:14",

        payloadCID:
            "bafybeigdyr-contract-draft-82a91",

        receiptStatus: "processed",

        stellarTx:
            "tx_71a82c..."
    },


    {
        id: "evt_002",

        type: "finance.approved",

        actor: "vault_finance",

        createdAt: "2026-06-10T11:42:00Z",

        time: "Jun 10 · 11:42",

        payloadCID:
            "bafybeifinance-clearance-91b2",

        receiptStatus: "processed",

        stellarTx:
            "tx_82bf91..."
    },


    {
        id: "evt_003",

        type: "direction.signed",

        actor: "vault_direction",

        createdAt: "2026-06-11T08:30:00Z",

        time: "Jun 11 · 08:30",

        payloadCID:
            "bafybeiexecutive-signature-44cd",

        receiptStatus: "processed",

        stellarTx:
            "tx_992ad1..."
    },


    {
        id: "evt_004",

        type: "counterparty.invited",

        actor: "vault_legal",

        createdAt: "2026-06-11T10:00:00Z",

        time: "Jun 11 · 10:00",

        payloadCID:
            "bafybeic3-channel-invite-a82f",

        receiptStatus:
            "received",

        stellarTx:
            "tx_a882fe..."
    },


    {
        id: "evt_005",

        type: "contract.countersigned",

        actor: "vault_supplier_x",

        createdAt: "2026-06-12T14:20:00Z",

        time: "Today · 14:20",

        payloadCID:
            "bafybeicountersignature-771a",

        receiptStatus:
            "processed",

        stellarTx:
            "tx_b88291..."
    },


    {
        id: "evt_006",

        type: "receipt.issued",

        actor: "vault_supplier_x",

        createdAt: "2026-06-12T14:25:00Z",

        time: "Today · 14:25",

        payloadCID:
            "bafybeireceipt-proof-912e",

        receiptStatus:
            "processed",

        stellarTx:
            "tx_c921fa..."
    }

];

export const threadAssets: ThreadAssetViewInterface[] = [
    {
        id: "msa",

        channelId: "contract-supplier",

        title: "Master Service Agreement",

        subtitle: "Supplier X · $240,000",


        status: "open",


        createdAt: "2026-06-10",


        participants: [
            {
                label: "L",
                color: "#7C3AED"
            },
            {
                label: "F",
                color: "#2563EB"
            },
            {
                label: "SUP",
                color: "#0891B2"
            }
        ],


        lastEvent:
            "counterparty.signed",


        stellarTx:
            "tx_contract_a11",



        events: [

            {
                id: "evt_001",

                type: "asset.created",

                actor: "vault_legal",

                createdAt:
                    "2026-06-10T09:00",

                time:
                    "Jun 10 · 09:00",

                payloadCID:
                    "bafy_create",

                receiptStatus:
                    "processed",

                stellarTx:
                    "tx_a11"
            },


            {
                id: "evt_002",

                type: "finance.approved",

                actor: "vault_finance",

                createdAt:
                    "2026-06-10T11:20",

                time:
                    "Jun 10 · 11:20",

                payloadCID:
                    "bafy_finance",

                receiptStatus:
                    "processed",

                stellarTx:
                    "tx_b22"
            },


            {
                id: "evt_003",

                type: "counterparty.signed",

                actor: "vault_supplier",

                createdAt:
                    "2026-06-12T14:00",

                time:
                    "Today",

                payloadCID:
                    "bafy_signature",

                receiptStatus:
                    "received"
            }

        ]
    },
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
    }

];
