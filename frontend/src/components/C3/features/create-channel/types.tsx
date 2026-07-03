import { ChannelTemplate } from "../../domain/channel/channel.mock";
import { ChannelSlot, VaultAssignment } from "../../domain/channel/channel.types";

export interface CreateChannelDraft {

    template?: ChannelTemplate;

    channelName?: string;

    slots?: ChannelSlot[];

    properties?: {
        key: string;
        value: string;
    }[];

    policy?: any;

    review?: any;

    assignments?: VaultAssignment[];

}

export interface Props {

    open: boolean;

    onClose: () => void;

}
