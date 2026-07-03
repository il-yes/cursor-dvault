// domain/threadAsset.ts

import { FlowStep } from "../../channel-view";

export type ThreadEvent = {

    eventId:string;


    channelId:string;

    assetId:string;


    eventType:string;


    createdAt:string;


    actor:string;


    idempotencyKey:string;


    payloadRef:string;


    status:
    | 'received'
    | 'processed'
    | 'rejected';

};

export type EventSummary = {

    type: string;

    actor: string;

    time: string;

    status?:
    'done'
    |
    'pending'
    |
    'rejected';
};

export interface ThreadEventView {

    id: string;

    type: string;

    actor: string;

    createdAt: string;

    payloadCID: string;

    receiptStatus:
        | "received"
        | "processed"
        | "rejected";

    stellarTx?: string;

    time?: string;

}
export type ThreadAssetStatus =
    | 'open'
    | 'transferring'
    | 'closed';


export type ThreadAsset = {

    id: string;


    ledgerId: string;

    channelId: string;


    title: string;

    type:
    | 'contract'
    | 'invoice'
    | 'payment'
    | 'attestation'
    | 'report';


    status: ThreadAssetStatus;


    createdAt: string;

    expiredAt?: string;


    events: ThreadEvent[];
};

export type AssetSummary = {

    id: string;

    type:
    | 'contract'
    | 'invoice'
    | 'report'
    | 'credential'
    | 'attestation'
    | 'attachment'
    | 'planning'
    | 'payment';


    title: string;

    status:
    | 'draft'
    | 'pending'
    | 'approved'
    | 'signed'
    | 'rejected'
    | 'open'
    | 'closed'
    | 'waiting'
    | 'processing'
    | 'received'
    | 'transferred'
    | 'completed';


    lastEvent: string;

};

export type AssetStatus =
    | 'draft'
    | 'pending'
    | 'approved'
    | 'signed'
    | 'rejected'
    | 'closed'
    | 'open'
    | 'waiting';


export type AssetView = {

    id: string;

    channelId: string;

    type:
    | 'contract'
    | 'invoice'
    | 'report'
    | 'credential'
    | 'attestation';


    title: string;

    subtitle: string;


    status: AssetStatus;


    createdAt: string;


    lastEvent: {
        type: string;
        at: string;
    };


    participants: string[];


    stellarTx: string;


    c3Visibility:
    | 'private'
    | 'shared'
    | 'external';


};




export type Screen =
    {
        type: "ledger"
    }

    |
    {
        type: "channel";
        channelId: string;
    }

    |
    {
        type: "asset";
        assetId: string;
    };

export interface ThreadAssetViewInterface {

    id: string;

    channelId: string;

    title: string;

    subtitle: string;

    type?: string;

    status:
    | "open"
    | "transferring"
    | "closed"
    | "disputed"
    | 'draft'
    | 'pending'
    | 'approved'
    | 'signed'
    | 'rejected'
    | 'closed'
    | 'open'
    | 'waiting'
    | 'processing'
    | 'received'
    | 'transferred'
    | 'completed';

    createdAt: string;

    expiresAt?: string;

    participants?: FlowStep[];

    lifecycle?: LifecycleStep[];

    events?: ThreadEventView[];

    payloads?: PayloadView[];

    receipts?: ReceiptView[];

    policy?: PolicyView;

    stellarTx?: string;

    relationship?: RelationshipView;

    externalVault?: string;
    direction?: string;
    role?: string;

    lastEvent?: string;


}

export interface LifecycleStep {

    label: string;

    actor: string;

    icon: string;

    color: string;

    state:
    | "done"
    | "current"
    | "waiting"
    | "rejected";

}

export interface PayloadView {

    cid: string;

    name: string;

    size: string;

    encrypted: boolean;

}

export interface ReceiptView {

    id: string;

    vault: string;

    status:
    | "received"
    | "processed"
    | "rejected";

    time: string;

    reason?: string;

}

export interface PolicyView {

    read: string[];

    write: string[];

}

export interface RelationshipView {

    status:
    | "internal"
    | "linked"
    | "active";

    externalVault?: string;

    role?: string;

    direction?:
    | "A→B"
    | "B→A"
    | "↔";

}

export const ThreadAssetViewInterfaceType: ThreadAssetViewInterface = {
    id: "",
    channelId: "",
    title: "",
    subtitle: "",
    type: "",
    status: "open",
    createdAt: "",
    expiresAt: "",
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
    role: "",
    lastEvent: "",
}   
