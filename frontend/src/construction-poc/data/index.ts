/**
 * Functional Data Boundary
 *
 * The single public data surface for the Construction POC. Views import from
 * here exactly as before; this module delegates to one active provider.
 *
 *   MockProvider  (default, current demo phase) — in-memory scenario, no backend
 *   CloudProvider (preserved)                    — authoritative AppAPI calls
 *
 * Switching is one environment variable:
 *   VITE_CONSTRUCTION_DATA_SOURCE=cloud   (restart the dev server)
 *
 * Delegation below is unconditional: the active provider is the only thing this
 * module ever calls, so mock mode cannot reach the Cloud implementation.
 */

import { cloudProvider } from "./cloudProvider";
import { mockProvider } from "./mockProvider";
import type { ConstructionDataProvider } from "./provider";

export * from "./constructionScenarioAdapter";
export * from "./projects.mock";
export * from "./scenario.mock";
export * from "./constructionTypes";
export {
  mapProjectOverviewDTOToProjectData,
  mapProcurementOverviewDTOToData,
  mapLogisticsOverviewDTOToData,
} from "./cloudProvider";
export type { ConstructionDataProvider } from "./provider";

const configuredSource = (import.meta.env.VITE_CONSTRUCTION_DATA_SOURCE || "mock")
  .toString()
  .trim()
  .toLowerCase();

export const constructionDataSource: "mock" | "cloud" =
  configuredSource === "cloud" ? "cloud" : "mock";

export const activeProvider: ConstructionDataProvider =
  constructionDataSource === "cloud" ? cloudProvider : mockProvider;

/* Projects */

export const getProjects = (): ReturnType<ConstructionDataProvider["getProjects"]> =>
  activeProvider.getProjects();

export const getMockProjects = (): ReturnType<ConstructionDataProvider["getMockProjects"]> =>
  activeProvider.getMockProjects();

export const getProject = (projectId?: string): ReturnType<ConstructionDataProvider["getProject"]> =>
  activeProvider.getProject(projectId);

export const createProject = (params: Parameters<ConstructionDataProvider["createProject"]>[0]) =>
  activeProvider.createProject(params);

export const getProjectOverview = (projectId?: string) =>
  activeProvider.getProjectOverview(projectId);

/* Procurement and logistics — projected from the scenario on mock */

export const getProcurementOverview = (requirementId?: string) =>
  activeProvider.getProcurementOverview(requirementId);

export const getLogisticsOverview = (deliveryId?: string) =>
  activeProvider.getLogisticsOverview(deliveryId);

export const getDecisionData = () => activeProvider.getDecisionData();

/* Scenario — canonical record access */

export const getActivityFeed = () => activeProvider.getActivityFeed();
export const getProjectRecord = () => activeProvider.getProjectRecord();
export const getMaterial = () => activeProvider.getMaterial();
export const getRequirement = () => activeProvider.getRequirement();
export const getSupplier = () => activeProvider.getSupplier();
export const getSupplierOffer = () => activeProvider.getSupplierOffer();
export const getSite = () => activeProvider.getSite();
export const getDelivery = () => activeProvider.getDelivery();
export const getTransportDelay = () => activeProvider.getTransportDelay();
export const getIssue = () => activeProvider.getIssue();
export const getAction = () => activeProvider.getAction();
export const getEvidenceDocument = () => activeProvider.getEvidenceDocument();
export const getCertificateDocument = () => activeProvider.getCertificateDocument();
export const getDecision = () => activeProvider.getDecision();
export const getInspection = () => activeProvider.getInspection();
/** The three real vault actors. There is no ConstructionStakeholder aggregate. */
export const getParticipants = () => activeProvider.getParticipants();
export const getCollaboration = () => activeProvider.getCollaboration();
export const getShareEntry = () => activeProvider.getShareEntry();
/** The authoritative 23-event sequence. */
export const getThreadEvents = () => activeProvider.getThreadEvents();

/* Collaboration */

export const appendThreadEvent = (threadId: string, kind: string, content: string) =>
  activeProvider.appendThreadEvent(threadId, kind, content);

export const inviteToChannel = (workspaceId: string, channelId: string, inviteeVaultId: string) =>
  activeProvider.inviteToChannel(workspaceId, channelId, inviteeVaultId);
