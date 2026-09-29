/**
 * CloudProvider — the authoritative AppAPI implementation.
 *
 * Preserved unchanged from the previous data/index.ts boundary and currently
 * inactive (VITE_CONSTRUCTION_DATA_SOURCE defaults to "mock"). All Wails
 * bindings and DTO mappers live in this file, so the mock path can never reach
 * the backend.
 */

import * as AppAPI from "../../../wailsjs/go/main/App";
import { tracecore_types } from "../../../wailsjs/go/models";
import { useAuthStore } from "@/store/useAuthStore";
import { useVaultStore } from "@/store/vaultStore";
import { listWorkspaces } from "@/services/api";
import { MOCK_PROJECTS } from "./projects.mock";
import type { ProjectData } from "./projects.mock";
import { scenarioAccessors } from "./scenarioAccessors";
import type {
  CreateProjectParams,
  LogisticsOverviewData,
  ProcurementData,
} from "./constructionTypes";
import type { ConstructionDataProvider } from "./provider";

export const cloudProvider: ConstructionDataProvider = {
  ...scenarioAccessors,

  async getProjects(): Promise<ProjectData[]> {
    const jwtToken = useAuthStore.getState().jwtToken || "";
    const vaultId = useVaultStore.getState().vault?.vault_runtime_context?.VaultID || "";
    try {
      const dtos = await AppAPI.ListConstructionProjects(jwtToken, vaultId);
      console.log("[BOUNDARY 4][AppAPI.ListConstructionProjects] raw dtos:", dtos, "isArray:", Array.isArray(dtos), "count:", Array.isArray(dtos) ? dtos.length : 0);
      if (!dtos || !Array.isArray(dtos)) return [];
      const mapped = dtos.map(mapProjectOverviewDTOToProjectData);
      console.log("[BOUNDARY 5][getProjects mapper] mapped count:", mapped.length);
      return mapped;
    } catch (err) {
      console.error("[AppAPI] ListConstructionProjects failed:", err);
      throw err;
    }
  },

  getMockProjects(): ProjectData[] {
    return MOCK_PROJECTS;
  },

  getProject(projectId?: string): ProjectData | undefined {
    if (!projectId) return undefined;
    return MOCK_PROJECTS.find(
      (p) => p.id.toLowerCase() === projectId.toLowerCase() || p.code.toLowerCase() === projectId.toLowerCase()
    );
  },

  async createProject(params: CreateProjectParams): Promise<ProjectData> {
    const jwtToken = useAuthStore.getState().jwtToken || "";

    let workspaceId = "";
    try {
      const workspaces = await listWorkspaces();
      if (workspaces && workspaces.length > 0) {
        workspaceId = workspaces[0].id;
      }
    } catch (err) {
      console.warn("[createProject] Could not resolve workspaceId:", err);
    }

    const code = params.code.trim() || "PRJ-002";
    const name = params.name.trim() || "Commercial Plaza North";
    const type = params.type.trim() || "COMMERCIAL";
    const loc = params.location.trim() || "Oaugadougou, Burkina Faso";
    let city = loc;
    let country = "";
    if (loc.includes(",")) {
      const parts = loc.split(",");
      city = parts[0].trim();
      country = parts.slice(1).join(",").trim();
    }

    const dto: tracecore_types.ProjectOverviewDTO = new tracecore_types.ProjectOverviewDTO({
      id: code,
      project_id: code,
      workspace_id: workspaceId,
      code: code,
      project_reference: code,
      name: name,
      project_name: name,
      type: type,
      project_type: type,
      status: "active",
      location: {
        address: "",
        city: city,
        country: country,
      } as any,
      progress_percent: 42,
      progress_percentage: 42,
      current_phase: "Phase 3 of 7: Foundation",
      recent_activity: "Formwork inspection scheduled",
      status_summary: "Formwork inspection scheduled",
    });

    console.log("[BOUNDARY 4][AppAPI.CreateConstructionProject] sending dto:", dto);
    const createdDto = await AppAPI.CreateConstructionProject(jwtToken, dto);
    console.log("[BOUNDARY 4][AppAPI.CreateConstructionProject] returned dto:", createdDto);
    return mapProjectOverviewDTOToProjectData(createdDto);
  },

  async getProjectOverview(projectId?: string): Promise<ProjectData | undefined> {
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
  },

  async getProcurementOverview(requirementId?: string): Promise<ProcurementData | undefined> {
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
  },

  async getLogisticsOverview(deliveryId?: string): Promise<LogisticsOverviewData | undefined> {
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
  },

  async appendThreadEvent(threadId: string, kind: string, content: string): Promise<void> {
    const jwtToken = useAuthStore.getState().jwtToken || "";
    if (!jwtToken) return;
    await AppAPI.AppendThreadEvent(jwtToken, threadId, kind, content);
  },

  async inviteToChannel(workspaceId: string, channelId: string, inviteeVaultId: string): Promise<boolean> {
    const jwtToken = useAuthStore.getState().jwtToken || "";
    if (!jwtToken) return false;
    // Go InviteToChannel signature: (token, workspaceId, channelId, inviteeVaultId)
    await AppAPI.InviteToChannel(jwtToken, workspaceId, channelId, inviteeVaultId);
    return true;
  },
};

export function mapProjectOverviewDTOToProjectData(dto: tracecore_types.ProjectOverviewDTO): ProjectData {
  const id = dto.id || dto.project_id || "";
  const code = dto.code || dto.project_reference || id;
  const name = dto.name || dto.project_name || "";
  const rawLoc = dto.location as any;
  const locStr = typeof rawLoc === "string" ? rawLoc : (rawLoc?.city || rawLoc?.address || "");

  return {
    id,
    code,
    name,
    contractId: dto.contract_id || "",
    type: dto.type || dto.project_type || "",
    sector: dto.sector || "",
    status: dto.status || "",
    location: locStr,
    description: dto.description || "",
    currentPhase: dto.current_phase || "",
    progressPercent: dto.progress_percent ?? (dto as any).progress_percentage ?? 0,
    openIssuesCount: dto.open_issues_count ?? 0,
    pendingDecisionsCount: dto.pending_decisions_count ?? 0,
    activeDelay: dto.active_delay || undefined,
    targetCompletion: dto.target_completion || "",
    recentActivity: dto.recent_activity || (dto as any).status_summary || "",
    image: dto.image || "",
    connectedOrgs: dto.connected_orgs || [],
    isAuthoritative: Boolean(dto.is_authoritative),
    budgetSpentPercent: dto.budget_spent_percent ?? 0,
    scheduleDay: dto.schedule_day ?? 0,
  };
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
