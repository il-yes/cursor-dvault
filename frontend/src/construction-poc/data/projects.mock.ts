import { SCENARIO_DATA } from "./constructionScenarioAdapter";

export interface ProjectData {
  id: string;
  code: string;
  name: string;
  contractId: string;
  type: string;
  sector: string;
  status: string;
  location: string;
  description: string;
  currentPhase: string;
  progressPercent: number;
  openIssuesCount: number;
  pendingDecisionsCount: number;
  activeDelay?: string;
  targetCompletion: string;
  recentActivity: string;
  image: string;
  connectedOrgs: string[];
  isAuthoritative: boolean;
  budgetSpentPercent: number;
  scheduleDay: number;
}

export const MOCK_PROJECTS: ProjectData[] = [
  {
    id: SCENARIO_DATA.project.id,
    code: SCENARIO_DATA.project.code,
    name: SCENARIO_DATA.project.name,
    contractId: "BFD-EUR-2024-099",
    type: SCENARIO_DATA.project.type,
    sector: "Infrastructure Sector • Transit Hub",
    status: "Active • On Schedule",
    location: SCENARIO_DATA.project.location,
    description: SCENARIO_DATA.project.description,
    currentPhase: "Structure (Phase 4 of 7)",
    progressPercent: 68,
    openIssuesCount: 7,
    pendingDecisionsCount: 1,
    activeDelay: "Route M1 Detour",
    targetCompletion: "Oct 2026",
    recentActivity: `Inspection #${SCENARIO_DATA.inspection.reference} approved (Today, 09:42)`,
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuCyhdIyunCsHJbMDm1UWI9VXIvPj7GBEI9F2ZoXR61KUPGZkgrMZflrpqE21OhZxTYJcbn_ADHB5FvkGObh09NsqZmUiZtY9aDyNoP4FUXeR4JNKmUosU_MNAcafsPG_jS6TunbkQiEX4l8nddQ1vxS4FpY21vauK7SCpG5QVlLgJNoef0nj2oirgVq0lBL-ZiR1C4ouQH7N-QrOJpG6YENchuyphAPTX21Wib9fOGuSX1OEhMU3cs",
    connectedOrgs: ["Acme Dev", "BuildCorp", "Engineering Partners", "EuroSteel", "FastBuild"],
    isAuthoritative: true,
    budgetSpentPercent: SCENARIO_DATA.project.budgetSpentPercent,
    scheduleDay: SCENARIO_DATA.project.scheduleDay,
  },
  {
    id: "PRJ-002",
    code: "PRJ-002",
    name: "Commercial Plaza North",
    type: "Commercial Mixed-Use",
    contractId: "CPN-LYN-2025-012",
    sector: "Commercial Real Estate",
    status: "Active",
    location: "Lyon, France",
    description: "Multi-story commercial office tower and public retail plaza development featuring sustainable reinforced concrete foundations.",
    currentPhase: "Phase 3 of 7: Foundation",
    progressPercent: 42,
    openIssuesCount: 0,
    pendingDecisionsCount: 2,
    targetCompletion: "Q2 2027",
    recentActivity: "Formwork inspection scheduled",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuBqV0HcND9OtO60ifVFcM1UpOyuC3fNdLNRbeb52p-yeTXGAxWlHHYN4nAJ1eKZmKtTECrFyDJuAQKsSZLBZV3kp4_yI1sjZjSeE3wjqh7JKvhhXEUDPPkO3M9HKLudq6hJ86eyUid-BL--zETvLZ1i67aHL7k8iMkraESDgz7cEAtXVMYHTV3wF2MmPIj1xP675VDXPvVVIeiMxB66Am_djeKn7Xf5OmwK5agWn5-v_V_Sok-WpxM",
    connectedOrgs: ["BuildCorp", "Lyon Metropole"],
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
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuALfyMUsVFSjYeWh6KhQj3f_xSm6PnPRu-0KUZBaC7j3_xi2NaVt64qQ17FpBARHKNwLMGQrWG9Hv_h0H0ygEKMHNEOwaZZz8uBMbI3KNIplL8UldeNgyKarOvs8EQ52nY9qp2dvmr5J-NNEb9UESAYCAtfxE1Zcz11XAaIWQ7I-vc9eFjZYYyzafSBXcpiOSJBt-U0QcpkRFmCyDIWNPTA88sUQk7pwru9tq3oEYwN6B0VT3f5yNI",
    connectedOrgs: ["EuroSteel", "FastBuild"],
    isAuthoritative: false,
    budgetSpentPercent: 18,
    scheduleDay: 35,
  },
];
