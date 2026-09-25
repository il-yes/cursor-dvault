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

export interface LogisticsOverviewData {
  delivery: {
    id: string;
    reference: string;
    status: string;
    plannedDeliveryDate: string;
    eta: string;
    actualDeliveryDate?: string;
    quantity: number;
    unit: string;
    deliveryNotes?: string;
  };
  project: {
    id: string;
    code: string;
    name: string;
    type: string;
  };
  requirement: {
    id: string;
    code: string;
    targetPhase: string;
  };
  material: {
    id: string;
    name: string;
    specification: string;
  };
  supplier: {
    id: string;
    name: string;
  };
  offer: {
    id: string;
    reference: string;
    totalPrice: string;
    unitPrice: string;
    promisedDeliveryDate: string;
  };
  site: {
    id: string;
    name: string;
  };
  transport: {
    id: string;
    reference: string;
    vehicle: string;
    driver: string;
    status: string;
    origin: string;
    destination: string;
    route: string;
    plannedDeparture: string;
    actualDeparture: string;
    plannedArrival: string;
    eta: string;
    actualArrival?: string;
    constraints?: string[];
    delayReason?: string;
  };
  issue?: {
    id: string;
    reference: string;
    title: string;
    description: string;
    severity: string;
    status: string;
    reportedBy: string;
    reportedAt: string;
    impact: string;
    evidenceReferences?: string[];
  };
  decision?: {
    id: string;
    reference: string;
    subject: string;
    context: string;
    technicalAssessment?: string;
    risksIdentified?: string[];
    participantsConsulted?: string[];
    optionsConsidered?: string[];
    decision: string;
    decidedBy: string;
    decisionDate: string;
    consequence: string;
    status: string;
  };
  inspection?: {
    id: string;
    reference: string;
    type: string;
    inspector: string;
    inspectionDate: string;
    status: string;
    result: string;
    criteria?: string[];
    findings?: string[];
    notes?: string;
    evidenceReferences?: string[];
  };
}

export function mapLogisticsOverviewDTOToData(dto: any): LogisticsOverviewData {
  const d = dto.delivery || {};
  const p = dto.project || {};
  const r = dto.requirement || {};
  const m = dto.material || {};
  const s = dto.supplier || {};
  const o = dto.offer || {};
  const st = dto.site || {};
  const t = dto.transport || {};
  const iss = dto.issue;
  const dec = dto.decision;
  const insp = dto.inspection;

  return {
    delivery: {
      id: d.id || d.reference || "",
      reference: d.reference || d.id || "",
      status: d.status || "",
      plannedDeliveryDate: d.planned_delivery_date || "",
      eta: d.eta || "",
      actualDeliveryDate: d.actual_delivery_date || undefined,
      quantity: d.quantity ?? 0,
      unit: d.unit || "",
      deliveryNotes: d.delivery_notes || undefined,
    },
    project: {
      id: p.id || p.code || "",
      code: p.code || p.id || "",
      name: p.name || "",
      type: p.type || "",
    },
    requirement: {
      id: r.id || r.code || "",
      code: r.code || r.id || "",
      targetPhase: r.target_phase || "",
    },
    material: {
      id: m.id || "",
      name: m.name || "",
      specification: m.specification || "",
    },
    supplier: {
      id: s.id || "",
      name: s.name || "",
    },
    offer: {
      id: o.id || o.reference || "",
      reference: o.reference || o.id || "",
      totalPrice: o.total_price || "",
      unitPrice: o.unit_price || "",
      promisedDeliveryDate: o.promised_delivery_date || "",
    },
    site: {
      id: st.id || "",
      name: st.name || "",
    },
    transport: {
      id: t.id || t.reference || "",
      reference: t.reference || t.id || "",
      vehicle: t.vehicle || "",
      driver: t.driver || "",
      status: t.status || "",
      origin: t.origin || "",
      destination: t.destination || "",
      route: t.route || "",
      plannedDeparture: t.planned_departure || "",
      actualDeparture: t.actual_departure || "",
      plannedArrival: t.planned_arrival || "",
      eta: t.eta || "",
      actualArrival: t.actual_arrival || undefined,
      constraints: t.constraints || undefined,
      delayReason: t.delay_reason || undefined,
    },
    issue: iss ? {
      id: iss.id || iss.reference || "",
      reference: iss.reference || iss.id || "",
      title: iss.title || "",
      description: iss.description || "",
      severity: iss.severity || "",
      status: iss.status || "",
      reportedBy: iss.reported_by || "",
      reportedAt: iss.reported_at || "",
      impact: iss.impact || "",
      evidenceReferences: iss.evidence_references || undefined,
    } : undefined,
    decision: dec ? {
      id: dec.id || dec.reference || "",
      reference: dec.reference || dec.id || "",
      subject: dec.subject || "",
      context: dec.context || "",
      technicalAssessment: dec.technical_assessment || undefined,
      risksIdentified: dec.risks_identified || undefined,
      participantsConsulted: dec.participants_consulted || undefined,
      optionsConsidered: dec.options_considered || undefined,
      decision: dec.decision || "",
      decidedBy: dec.decided_by || "",
      decisionDate: dec.decision_date || "",
      consequence: dec.consequence || "",
      status: dec.status || "",
    } : undefined,
    inspection: insp ? {
      id: insp.id || insp.reference || "",
      reference: insp.reference || insp.id || "",
      type: insp.type || "",
      inspector: insp.inspector || "",
      inspectionDate: insp.inspection_date || "",
      status: insp.status || "",
      result: insp.result || "",
      criteria: insp.criteria || undefined,
      findings: insp.findings || undefined,
      notes: insp.notes || undefined,
      evidenceReferences: insp.evidence_references || undefined,
    } : undefined,
  };
}

export async function getLogisticsOverview(deliveryId?: string): Promise<LogisticsOverviewData | undefined> {
  if (!deliveryId) return undefined;

  const jwtToken = useAuthStore.getState().jwtToken || "";

  try {
    const dto = await AppAPI.GetLogisticsOverview(jwtToken, deliveryId);
    if (!dto) return undefined;
    return mapLogisticsOverviewDTOToData(dto);
  } catch (err) {
    console.error(`[AppAPI] GetLogisticsOverview failed for deliveryId=${deliveryId}:`, err);
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
