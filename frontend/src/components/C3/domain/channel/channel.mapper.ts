import { ChannelResponse } from "@/services/api";
import { ChannelRow, ChannelView, FlowStep } from "./channel.types";
import { ThreadAssetViewInterface } from "../thread/asset.types";

type RowType = ChannelRow["type"];

const TEMPLATE_TYPE: Record<string, RowType> = {
    "contract-execution": "Contract",
    "contract": "Contract",
    "payment": "Payment",
    "invoice": "Payment",
    "procurement": "Procurement",
    "purchase-order": "Procurement",
    "payroll": "Payroll",
    "budget": "Governance",
    "governance": "Governance",
    "compliance": "Compliance",
    "audit": "Compliance",
    "onboarding": "Onboarding",
};

const TEMPLATE_LABEL: Record<string, string> = {
    "contract-execution": "Contract Execution",
    "payment": "Invoice Processing",
    "procurement": "Procurement",
    "payroll": "Payroll",
    "budget": "Budget Allocation",
    "governance": "Governance",
    "compliance": "Compliance Audit",
    "onboarding": "Employee Onboarding",
};

const normalizeTemplateId = (templateId?: string): string =>
    (templateId || "default").trim().toLowerCase();

const resolveType = (templateId?: string): RowType =>
    TEMPLATE_TYPE[normalizeTemplateId(templateId)] || "Contract";

const resolveLabel = (templateId?: string): string =>
    TEMPLATE_LABEL[normalizeTemplateId(templateId)] || "Channel";

const mapParticipants = (participants?: unknown[]): FlowStep[] => {
    if (!Array.isArray(participants)) return [];
    return participants
        .map((p) => {
            const rec = (p ?? {}) as { label?: unknown; name?: unknown; color?: unknown };
            return {
                label: typeof rec.label === "string" ? rec.label : typeof rec.name === "string" ? rec.name : "",
                color: typeof rec.color === "string" ? rec.color : "#666666",
            };
        })
        .filter((step) => step.label !== "");
};

const formatDate = (value?: string): string => {
    if (!value) return "";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "";
    return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
};

const mapRowStatus = (status?: string): ChannelRow["status"] => {
    if (status === "active") return "active";
    if (status === "revoked") return "revoked";
    return "pending";
};

const mapC3Status = (status?: string): ChannelRow["c3Status"] => {
    if (status === "active") return "active";
    if (status === "revoked") return "linked";
    return "internal";
};

const c3LabelFor = (status: ChannelRow["c3Status"]): string => {
    switch (status) {
        case "active":
            return "⛓●";
        case "linked":
            return "⛓✓";
        default:
            return "⛓+";
    }
};

export const toChannelRow = (channel: ChannelResponse): ChannelRow => {
    const status = mapRowStatus(channel.status);
    const c3Status = mapC3Status(channel.status);

    return {
        id: channel.id,
        status,
        type: resolveType(channel.template_id),
        title: channel.title,
        subtitle: resolveLabel(channel.template_id),
        participants: mapParticipants(channel.participants),
        assetCount: channel.asset_count ?? 0,
        lastEvent: channel.last_event || "",
        lastActivity: formatDate(channel.updated_at || channel.created_at),
        stellarTx: "",
        c3Status,
        c3Label: c3LabelFor(c3Status),
    };
};

export const toChannelRows = (channels: ChannelResponse[]): ChannelRow[] =>
    channels.map(toChannelRow);

const mapViewStatus = (status?: string): ChannelView["status"] => {
    const s = (status || "active").toLowerCase();
    if (s === "pending" || s === "revoked" || s === "closed" || s === "open") return s as ChannelView["status"];
    return "active";
};

type StatusView = { dotClass: string; label: string };

const STATUS_VIEW: Record<ChannelView["status"], StatusView> = {
    active: { dotClass: "s-ok", label: "Active" },
    pending: { dotClass: "s-pend", label: "Pending" },
    revoked: { dotClass: "s-dispute", label: "Revoked" },
    closed: { dotClass: "s-pend", label: "Closed" },
    open: { dotClass: "s-pend", label: "Open" },
};

export const toChannelStatusView = (status: ChannelView["status"]): StatusView =>
    STATUS_VIEW[status] || { dotClass: "s-pend", label: "Pending" };

export const toRowStatusView = (status: ChannelRow["status"]): StatusView => {
    if (status === "active") return { dotClass: "s-ok", label: "Active" };
    if (status === "revoked") return { dotClass: "s-dispute", label: "Revoked" };
    if (status === "dispute") return { dotClass: "s-dispute", label: "Dispute" };
    return { dotClass: "s-pend", label: "Pending" };
};


export const toThreadAssetItem = (thread: { id: string; channel_id: string; title: string; subtitle?: string; asset_type?: string; status?: string; created_at?: string }): ThreadAssetViewInterface => ({
    id: thread.id,
    channelId: thread.channel_id,
    title: thread.title,
    subtitle: thread.subtitle || thread.asset_type || "thread",
    type: thread.asset_type || "note",
    status: (thread.status as ThreadAssetViewInterface["status"]) || "open",
    createdAt: formatDate(thread.created_at),
    lastEvent: "thread.created",
});

export const toChannelView = (channel: ChannelResponse, threads: { id: string; channel_id: string; title: string; subtitle?: string; asset_type?: string; status?: string; created_at?: string }[] = []): ChannelView => {
    const lastActivity = formatDate(channel.updated_at || channel.created_at);
    const mappedItems = threads.map(toThreadAssetItem);
    console.log("[BOUNDARY LOG 3: TO_CHANNEL_VIEW MAPPER]", {
        channelId: channel?.id,
        inputThreadsCount: threads.length,
        outputItemsCount: mappedItems.length,
        outputItemIDs: mappedItems.map((i) => i.id),
    });

    return {
        id: channel.id,
        title: channel.title,
        subtitle: resolveLabel(channel.template_id),
        status: mapViewStatus(channel.status),
        participants: mapParticipants(channel.participants),
        assets: {
            total: mappedItems.length || channel.asset_count || 0,
            items: mappedItems,
        },
        activity: {
            lastEvent: channel.last_event || "",
            lastActivity,
            events: [],
        },
        policy: {
            read: [],
            write: [],
        },
        stellarTx: "",
        c3: {
            status: channel.status === "active" ? "active" : "internal",
        },
        slots: [],
        defaultProperties: [],
    };
};
