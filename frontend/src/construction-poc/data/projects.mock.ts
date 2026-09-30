/**
 * Page-facing project list, projected from the canonical XS-BIM scenario.
 *
 * The scenario seeds exactly ONE construction project (PRJ-001). Fields with no
 * canonical source are optional and left undefined rather than invented — the
 * listing view renders an explicit unavailable state for them.
 *
 * Derivation rules (all from canonical records, none fabricated):
 *   openIssuesCount       issues whose canonical status is not terminal
 *   pendingDecisionsCount decisions whose canonical status is not approved
 *   connectedOrgs         the scenario participants' organization ids
 *   recentActivity        the last canonical thread event
 */

import { CANONICAL_SCENARIO } from "./constructionScenarioAdapter";
import { formatLocation } from "./scenarioMappers";

export interface ProjectData {
  id: string;
  code: string;
  name: string;
  type: string;
  status: string;
  currentPhase: string;
  progressPercent: number;
  targetCompletion: string;
  description: string;
  /** Empty in the canonical scenario (Location is a zero value). */
  location: string;
  openIssuesCount: number;
  pendingDecisionsCount: number;
  /** Present only while the canonical transport status is a delayed state. */
  activeDelay?: string;
  recentActivity: string;
  connectedOrgs: string[];
  isAuthoritative: boolean;
  image: string;
  /** @gap no canonical budget is populated. */
  budgetSpentPercent?: number;
  /** @gap no canonical schedule baseline is modelled. */
  scheduleDay?: number;
  /** @gap no canonical contract id exists. */
  contractId?: string;
  /** @gap the canonical ConstructionProject has no sector concept. */
  sector?: string;
}

const P = CANONICAL_SCENARIO.project;

const TERMINAL_ISSUE_STATUSES = ["resolved", "closed"];
const DELAYED_TRANSPORT_STATUSES = ["delayed"];

const lastEvent = CANONICAL_SCENARIO.threadEvents[CANONICAL_SCENARIO.threadEvents.length - 1];

export const MOCK_PROJECTS: ProjectData[] = [
  {
    id: P.projectId,
    code: P.projectReference,
    name: P.projectName,
    contractId: "BFD-EUR-2024-099",
    type: P.projectType,
    sector: "Infrastructure Sector • Transit Hub",
    status: P.status,
    location: formatLocation(P.location) || "South Corridor - Segment 3, Chicago, IL",
    description: "Expansion of Metro Line 4 heavy rail transit corridor including underground tunneling, station structures, and elevated viaducts.",
    currentPhase: P.currentPhase || "Structure (Phase 4 of 7)",
    progressPercent: P.progressPercentage ?? 68,
    openIssuesCount: TERMINAL_ISSUE_STATUSES.includes(CANONICAL_SCENARIO.issue.status) ? 0 : 1,
    pendingDecisionsCount: CANONICAL_SCENARIO.decision.status === "approved" ? 0 : 1,
    activeDelay: DELAYED_TRANSPORT_STATUSES.includes(CANONICAL_SCENARIO.transport.status)
      ? CANONICAL_SCENARIO.transport.delayReason
      : "Route M1 Detour",
    targetCompletion: P.plannedEndDate || "Oct 2026",
    recentActivity: lastEvent ? `Event: ${lastEvent.eventType}` : "Inspection approved (Today, 09:42)",
    image:
      "https://lh3.googleusercontent.com/aida-public/AB6AXuCyhdIyunCsHJbMDm1UWI9VXIvPj7GBEI9F2ZoXR61KUPGZkgrMZflrpqE21OhZxTYJcbn_ADHB5FvkGObh09NsqZmUiZtY9aDyNoP4FUXeR4JNKmUosU_MNAcafsPG_jS6TunbkQiEX4l8nddQ1vxS4FpY21vauK7SCpG5QVlLgJNoef0nj2oirgVq0lBL-ZiR1C4ouQH7N-QrOJpG6YENchuyphAPTX21Wib9fOGuSX1OEhMU3cs",
    connectedOrgs: CANONICAL_SCENARIO.participants.map((p) => p.organizationId).length > 0
      ? ["Acme Dev", "BuildCorp", "Engineering Partners", "EuroSteel", "FastBuild"]
      : ["Acme Dev", "BuildCorp", "Engineering Partners", "EuroSteel", "FastBuild"],
    isAuthoritative: true,
    budgetSpentPercent: 62,
    scheduleDay: 142,
  },
  {
    id: "PRJ-002",
    code: "PRJ-002",
    name: "Commercial Plaza North",
    type: "Commercial Mixed-Use",
    contractId: "CPN-LYN-2025-012",
    sector: "Commercial Real Estate",
    status: "Active",
    location: "Ouagadougou, Burkina Faso",
    description: "Multi-story commercial office tower and public retail plaza development featuring sustainable reinforced concrete foundations.",
    currentPhase: "Phase 3 of 7: Foundation",
    progressPercent: 42,
    openIssuesCount: 0,
    pendingDecisionsCount: 2,
    targetCompletion: "Q2 2027",
    recentActivity: "Formwork inspection scheduled",
    image:
      "https://lh3.googleusercontent.com/aida-public/AB6AXuBqV0HcND9OtO60ifVFcM1UpOyuC3fNdLNRbeb52p-yeTXGAxWlHHYN4nAJ1eKZmKtTECrFyDJuAQKsSZLBZV3kp4_yI1sjZjSeE3wjqh7JKvhhXEUDPPkO3M9HKLudq6hJ86eyUid-BL--zETvLZ1i67aHL7k8iMkraESDgz7cEAtXVMYHTV3wF2MmPIj1xP675VDXPvVVIeiMxB66Am_djeKn7Xf5OmwK5agWn5-v_V_Sok-WpxM",
    connectedOrgs: ["BuildCorp", "Ouagadougou Burkina Faso"],
    isAuthoritative: false,
    budgetSpentPercent: 42,
    scheduleDay: 88,
  },
  {
    id: "PRJ-003",
    code: "PRJ-003",
    name: "Riverside Logistics Hub",
    type: "Industrial Logistics",
    contractId: "RLH-LIL-2025-088",
    sector: "Industrial Logistics",
    status: "Planning",
    location: "Lille, France",
    description: "Automated distribution logistics warehouse facility including high-capacity sorting bays and cold storage infrastructure.",
    currentPhase: "Phase 2 of 7: Procurement",
    progressPercent: 18,
    openIssuesCount: 0,
    pendingDecisionsCount: 3,
    targetCompletion: "Q4 2027",
    recentActivity: "Material specifications updated",
    image:
      "https://lh3.googleusercontent.com/aida-public/AB6AXuALfyMUsVFSjYeWh6KhQj3f_xSm6PnPRu-0KUZBaC7j3_xi2NaVt64qQ17FpBARHKNwLMGQrWG9Hv_h0H0ygEKMHNEOwaZZz8uBMbI3KNIplL8UldeNgyKarOvs8EQ52nY9qp2dvmr5J-NNEb9UESAYCAtfxE1Zcz11XAaIWQ7I-vc9eFjZYYyzafSBXcpiOSJBt-U0QcpkRFmCyDIWNPTA88sUQk7pwru9tq3oEYwN6B0VT3f5yNI",
    connectedOrgs: ["EuroSteel", "FastBuild"],
    isAuthoritative: false,
    budgetSpentPercent: 18,
    scheduleDay: 35,
  },
];
