import { MOCK_PROJECTS, ProjectData } from "./projects.mock";
import { MOCK_ACTIVITY_LOG, ActivityItem } from "./scenario.mock";
import { SCENARIO_DATA } from "./constructionScenarioAdapter";
import * as AppAPI from "../../../wailsjs/go/main/App";
import { useAuthStore } from "@/store/useAuthStore";
import { tracecore_types } from "../../../wailsjs/go/models";

export * from "./constructionScenarioAdapter";
export * from "./projects.mock";
export * from "./scenario.mock";

/**
 * Functional Data Boundary
 * 
 * Exposes accessors for UI views. Later, this boundary will delegate to AppAPI
 * and Go backend endpoints without changing UI presentation components.
 */

export function getProjects(): ProjectData[] {
  return MOCK_PROJECTS;
}

export function mapProjectOverviewDTOToProjectData(dto: tracecore_types.ProjectOverviewDTO): ProjectData {
  const id = dto.id || dto.project_id || "";
  const code = dto.code || dto.project_reference || id;
  const name = dto.name || dto.project_name || "";

  return {
    id,
    code,
    name,
    contractId: dto.contract_id || "",
    type: dto.type || dto.project_type || "",
    sector: dto.sector || "",
    status: dto.status || "",
    location: dto.location || "",
    description: dto.description || "",
    currentPhase: dto.current_phase || "",
    progressPercent: dto.progress_percent ?? 0,
    openIssuesCount: dto.open_issues_count ?? 0,
    pendingDecisionsCount: dto.pending_decisions_count ?? 0,
    activeDelay: dto.active_delay || undefined,
    targetCompletion: dto.target_completion || "",
    recentActivity: dto.recent_activity || "",
    image: dto.image || "",
    connectedOrgs: dto.connected_orgs || [],
    isAuthoritative: Boolean(dto.is_authoritative),
    budgetSpentPercent: dto.budget_spent_percent ?? 0,
    scheduleDay: dto.schedule_day ?? 0,
  };
}

export function getProject(projectId?: string): ProjectData | undefined {
  if (!projectId) return undefined;
  return MOCK_PROJECTS.find(
    (p) => p.id.toLowerCase() === projectId.toLowerCase() || p.code.toLowerCase() === projectId.toLowerCase()
  );
}

export async function getProjectOverview(projectId?: string): Promise<ProjectData | undefined> {
  if (!projectId) return undefined;

  const jwtToken = useAuthStore.getState().jwtToken || "";

  try {
    const dto = await AppAPI.GetProjectOverview(jwtToken, projectId);
    if (!dto) return undefined;
    return mapProjectOverviewDTOToProjectData(dto);
  } catch (err) {
    console.error(`[AppAPI] GetProjectOverview failed for projectId=${projectId}:`, err);
    throw err;
  }
}

export function getActivityFeed(): ActivityItem[] {
  return MOCK_ACTIVITY_LOG;
}

export function getRequirement() {
  return SCENARIO_DATA.requirement;
}

export function getSupplierOffer() {
  return SCENARIO_DATA.offer;
}

export function getDelivery() {
  return SCENARIO_DATA.delivery;
}

export function getTransportDelay() {
  return SCENARIO_DATA.transport;
}

export function getIssue() {
  return SCENARIO_DATA.issue;
}

export function getEvidenceDocument() {
  return SCENARIO_DATA.evidenceDoc;
}

export function getDecision() {
  return SCENARIO_DATA.decision;
}

export function getInspection() {
  return SCENARIO_DATA.inspection;
}

export function getStakeholders() {
  return SCENARIO_DATA.stakeholders;
}

export function getTraceMilestones() {
  return SCENARIO_DATA.traceMilestones;
}
