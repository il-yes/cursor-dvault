/**
 * ConstructionDataProvider — the contract behind the Functional Data Boundary.
 *
 * The presentation layer talks to exactly one active provider, selected in
 * data/index.ts. Implementations are MockProvider (current demo phase) and
 * CloudProvider (preserved, inactive). Members are expressed in existing
 * view-model types; no new domain model is introduced here.
 */

import type { ProjectData } from "./projects.mock";
import type { scenarioAccessors } from "./scenarioAccessors";
import type {
  CreateProjectParams,
  DecisionData,
  LogisticsOverviewData,
  ProcurementData,
} from "./constructionTypes";

export type ConstructionDataProvider = typeof scenarioAccessors & {
  getProjects(): Promise<ProjectData[]>;
  getMockProjects(): ProjectData[];
  getProject(projectId?: string): ProjectData | undefined;
  createProject(params: CreateProjectParams): Promise<ProjectData>;
  getProjectOverview(projectId?: string): Promise<ProjectData | undefined>;

  /**
   * Projected from the presentation scenario by data/scenarioMappers.ts. The
   * mock provider resolves `undefined` for ids the scenario does not cover, and
   * never reaches Cloud.
   */
  getProcurementOverview(requirementId?: string): Promise<ProcurementData | undefined>;
  getLogisticsOverview(deliveryId?: string): Promise<LogisticsOverviewData | undefined>;

  /**
   * Synchronous page-facing projection of the canonical decision, shaped to the
   * DEC-1042 composition. Synchronous because the canonical aggregate is a
   * static scenario record on both providers — there is no AppAPI equivalent yet.
   */
  getDecisionData(): DecisionData;

  appendThreadEvent(threadId: string, kind: string, content: string): Promise<void>;

  /**
   * Resolves true when the invitation was dispatched to the Cloud/C3 backend,
   * false when it was acknowledged locally (no active session).
   */
  inviteToChannel(workspaceId: string, channelId: string, inviteeVaultId: string): Promise<boolean>;
};
