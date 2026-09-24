export const HOME = "/";
export const NOT_FOUND = "*";

export const CONSTRUCTION_ROUTES = {
  HOME: "/dashboard/construction",
  PROJECTS: "/dashboard/construction/projects",
  PROJECT: (projectId: string = ":projectId") => `/dashboard/construction/projects/${projectId}`,
  CHANNELS: "/dashboard/construction/channels",
  CHANNEL: (channelId: string = ":channelId") => `/dashboard/construction/channels/${channelId}`,
  THREAD: (channelId: string = ":channelId", threadId: string = ":threadId") => `/dashboard/construction/channels/${channelId}/threads/${threadId}`,
  STAKEHOLDERS: "/dashboard/construction/stakeholders",
  REQUIREMENTS: "/dashboard/construction/requirements",
  OFFERS: "/dashboard/construction/offers",
  DELIVERIES: "/dashboard/construction/deliveries",
  TRANSPORT: "/dashboard/construction/transport",
  ISSUES: "/dashboard/construction/issues",
  DECISIONS: "/dashboard/construction/decisions",
  INSPECTIONS: "/dashboard/construction/inspections",
  FIELD: "/dashboard/construction/field",
  PROVENANCE: "/dashboard/construction/provenance",
  HISTORY: "/dashboard/construction/history",
  ACTIVITY: "/dashboard/construction/activity",
  DOCUMENTS: "/dashboard/construction/documents",
  PROFILE: "/dashboard/construction/profile",
};
