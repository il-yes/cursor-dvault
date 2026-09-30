/**
 * MockProvider — the active provider for the current Construction demo phase.
 *
 * Serves the XS-BIM presentation scenario entirely in memory. It performs no
 * network or Wails calls. Project data and the procurement/logistics view models
 * are projected from the scenario by data/scenarioMappers.ts.
 */

import { MOCK_PROJECTS } from "./projects.mock";
import type { ProjectData } from "./projects.mock";
import { scenarioAccessors } from "./scenarioAccessors";
import {
  mapScenarioToDecisionData,
  mapScenarioToLogisticsOverviewData,
  mapScenarioToProcurementData,
  scenarioHasDelivery,
  scenarioHasRequirement,
} from "./scenarioMappers";
import type {
  CreateProjectParams,
  DecisionData,
  LogisticsOverviewData,
  ProcurementData,
} from "./constructionTypes";
import type { ConstructionDataProvider } from "./provider";

/** Session-scoped project list, so a project created in the demo is listed. */
let projects: ProjectData[] = [...MOCK_PROJECTS];

export const mockProvider: ConstructionDataProvider = {
  ...scenarioAccessors,

  async getProjects(): Promise<ProjectData[]> {
    return projects;
  },

  getMockProjects(): ProjectData[] {
    return projects;
  },

  getProject(projectId?: string): ProjectData | undefined {
    if (!projectId) return undefined;
    return projects.find(
      (p) =>
        p.id.toLowerCase() === projectId.toLowerCase() ||
        p.code.toLowerCase() === projectId.toLowerCase()
    );
  },

  async createProject(params: CreateProjectParams): Promise<ProjectData> {
    const code = params.code.trim();
    const name = params.name.trim();
    const created: ProjectData = {
      id: code,
      code,
      name,
      contractId: "",
      type: params.type.trim(),
      sector: "",
      status: "Active",
      location: params.location.trim(),
      description: "",
      currentPhase: "Phase 3 of 7: Foundation",
      progressPercent: 42,
      openIssuesCount: 0,
      pendingDecisionsCount: 0,
      targetCompletion: "Q2 2027",
      recentActivity: "Formwork inspection scheduled",
      image: MOCK_PROJECTS[0]?.image ?? "",
      connectedOrgs: [],
      isAuthoritative: false,
      budgetSpentPercent: 42,
      scheduleDay: 88,
    };
    projects = [created, ...projects];
    return created;
  },

  async getProjectOverview(projectId?: string): Promise<ProjectData | undefined> {
    return mockProvider.getProject(projectId);
  },

  async getProcurementOverview(requirementId?: string): Promise<ProcurementData | undefined> {
    if (!scenarioHasRequirement(requirementId)) return undefined;
    return mapScenarioToProcurementData();
  },

  async getLogisticsOverview(deliveryId?: string): Promise<LogisticsOverviewData | undefined> {
    if (deliveryId && !scenarioHasDelivery(deliveryId)) return undefined;
    return mapScenarioToLogisticsOverviewData();
  },

  getDecisionData(): DecisionData {
    return mapScenarioToDecisionData();
  },

  async appendThreadEvent(): Promise<void> {
    // Thread events are already applied optimistically by the presentation
    // layer; nothing to persist while running on mock data.
  },

  async inviteToChannel(): Promise<boolean> {
    // Invite is acknowledged locally; nothing to dispatch while on mock data.
    return false;
  },
};
