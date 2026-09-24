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

export interface ProcurementData {
  id: string;
  requirementId: string;
  projectId: string;
  projectCode: string;
  projectName: string;
  code: string;
  materialId: string;
  materialName: string;
  specification: string;
  quantity: number;
  unit: string;
  targetPhase: string;
  status: string;
  requiredDate: string;
  priority: string;
  siteId: string;
  siteName: string;
  invitedSuppliersCount: number;
  offersReceivedCount: number;
  offerId: string;
  offerReference: string;
  supplierId: string;
  supplierName: string;
  totalPrice: string;
  unitPrice: string;
  promisedDeliveryDate: string;
  offerStatus: string;
  isVerifiedSupplier: boolean;
}

export function mapProcurementOverviewDTOToData(dto: tracecore_types.ProcurementOverviewDTO): ProcurementData {
  const requirementId = dto.requirement_id || dto.id || dto.code || "";
  return {
    id: requirementId,
    requirementId,
    projectId: dto.project_id || "",
    projectCode: dto.project_code || "",
    projectName: dto.project_name || "",
    code: dto.code || requirementId,
    materialId: dto.material_id || "",
    materialName: dto.material_name || "",
    specification: dto.specification || "",
    quantity: dto.quantity ?? 0,
    unit: dto.unit || "",
    targetPhase: dto.target_phase || "",
    status: dto.status || "",
    requiredDate: dto.required_date || "",
    priority: dto.priority || "",
    siteId: dto.site_id || "",
    siteName: dto.site_name || "",
    invitedSuppliersCount: dto.invited_suppliers_count ?? 0,
    offersReceivedCount: dto.offers_received_count ?? 0,
    offerId: dto.offer_id || "",
    offerReference: dto.offer_reference || "",
    supplierId: dto.supplier_id || "",
    supplierName: dto.supplier_name || "",
    totalPrice: dto.total_price || "",
    unitPrice: dto.unit_price || "",
    promisedDeliveryDate: dto.promised_delivery_date || "",
    offerStatus: dto.offer_status || "",
    isVerifiedSupplier: Boolean(dto.is_verified_supplier),
  };
}

export async function getProcurementOverview(requirementId?: string): Promise<ProcurementData | undefined> {
  if (!requirementId) return undefined;

  const jwtToken = useAuthStore.getState().jwtToken || "";

  try {
    const dto = await AppAPI.GetProcurementOverview(jwtToken, requirementId);
    if (!dto) return undefined;
    return mapProcurementOverviewDTOToData(dto);
  } catch (err) {
    console.error(`[AppAPI] GetProcurementOverview failed for requirementId=${requirementId}:`, err);
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
