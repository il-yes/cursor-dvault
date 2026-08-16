// domain/channel.ts

import { AccessPolicy } from "../policy/policy.types";
import { EventSummary, ThreadAsset, ThreadEvent } from "../thread/asset.types";
import { ThreadAssetViewInterface } from "../thread/asset.types";

export type ChannelStatus =
    | 'pending'
    | 'active'
    | 'revoked';

export type FlowStep = {
    label: string;
    color: string;
};

export type Policy = {

    read: string[];

    write: string[];

    append?: string[];

    expiresAt?: string;

};


export type Channel = {

    id: string;

    status: ChannelStatus;


    title: string;
    subtitle: string;


    participants: FlowStep[];


    allowedEventTypes: string[];

    allowedPaths?: string[];


    policy: AccessPolicy;


    assets: ThreadAsset[];


    lastEvent?: ThreadEvent;


    createdAt: string;

    revokedAt?: string;


    stellarTx: string;


    c3Status:
    | 'extendable'
    | 'linked'
    | 'active';


    c3Label: string;
};

export type ChannelView = {

    id: string;


    title: string;

    subtitle: string;


    status:
    | 'active'
    | 'pending'
    | 'revoked'
    | 'closed'
    | 'open'


    participants: FlowStep[];


    assets: {
        total: number;

        items: ThreadAssetViewInterface[];
    };


    activity: {
        lastEvent: string;
        lastActivity: string;

        events: EventSummary[];
    };


    policy: {
        read: string[];
        write: string[];
    };


    stellarTx: string;


    c3: {
        status:
        | 'internal'
        | 'linked'
        | 'active'
        | 'closed'
        | 'open'
        | 'waiting'

        externalVault?: string;
    };

    slots: ChannelSlot[]

    defaultProperties: ChannelProperty[]
};

export interface ChannelSlot {

    id: string;

    name: string;

    role: string;

    gated: boolean;

    vault?: string

}

export interface ChannelProperty {

    key: string;

    required: boolean;

    defaultValue?: string;


}

export interface ChannelDraft {

    templateId: string;

    customProperties: PropertyValue[];

}
export interface CreateChannelPayload {
    templateId: string;
    title: string;

    slots: ChannelSlot[];

    properties: {
        key: string;
        value: string;
    }[];

    assignments: VaultAssignment[];

    policy?: any;
}
export interface PropertyValue {

    key: string;

    value: string;

}

export type Department = {
    id: string;
    name: string;
    color: string;
}

export interface VaultAssignment {

    vault: string;  // vault_finance

    owner: {

        id: string;     // alice, #54..., did:c3:alice or vault://legal/alice

        label: string; // Alice

        publicKey: string; // G...

    };

    vaultColor?: string;

}

export function createEmptyDraft() {
    return {
        templateId: "",
        customProperties: []
    }
}
export function createEmptyAssignments() {
    return {
        vault: "",
        owner: {
            id: "",
            label: "",
            publicKey: ""
        }
    }
}
export function createEmptyReview() {
    return {
        policy: createEmptyPolicy(),
        assignments: createEmptyAssignments(),
        customProperties: [],
    }
}
export function createEmptyPolicy() {
    return {
        read: [],
        write: [],
        append: [],
        expiresAt: ""
    }
}
export function createEmptyChannelDraft() {
    return {
        template: null,
        channelName: "",
        slots: [],
        properties: [],
        assignments: [],
        policy: null,
        review: null,
    }
    
}
export type ChannelRow = {
    id: string;

    status: 'active' | 'pending' | 'dispute' | 'revoked';

    type:
    | 'Contract'
    | 'Payment'
    | 'Procurement'
    | 'Payroll'
    | 'Governance'
    | 'Compliance'
    | 'Onboarding';

    title: string;
    subtitle: string;

    participants: FlowStep[];

    assetCount: number;

    lastEvent: string;

    lastActivity: string;

    stellarTx: string;

    c3Status:
    | 'internal'
    | 'linked'
    | 'active';

    c3Label: string;
};
