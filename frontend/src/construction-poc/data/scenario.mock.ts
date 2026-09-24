import { SCENARIO_DATA } from "./constructionScenarioAdapter";

export interface ActivityItem {
  id: string;
  time: string;
  type: string;
  badge: string;
  title: string;
  actor: string;
  link: string;
}

export const MOCK_ACTIVITY_LOG: ActivityItem[] = [
  {
    id: "ACT-001",
    time: "08:42",
    type: "PROJECT_CREATED",
    badge: "PRJ-001",
    title: "Project created: Metro Line 4 Expansion",
    actor: "Manuel (Project Owner)",
    link: "/dashboard/construction/projects/PRJ-001",
  },
  {
    id: "ACT-002",
    time: "08:51",
    type: "STAKEHOLDER_INVITED",
    badge: "STK-001",
    title: "Stakeholder invited: Lead Architect → Project",
    actor: "Manuel (Project Owner)",
    link: "/dashboard/construction/stakeholders",
  },
  {
    id: "ACT-003",
    time: "09:14",
    type: "REQUIREMENT_CREATED",
    badge: SCENARIO_DATA.requirement.code,
    title: `Requirement created: ${SCENARIO_DATA.requirement.materialName}`,
    actor: "Alex Rivera (Project Manager)",
    link: "/dashboard/construction/requirements",
  },
  {
    id: "ACT-004",
    time: "09:42",
    type: "SUPPLIER_INVITED",
    badge: SCENARIO_DATA.supplier.id,
    title: `Supplier invitation sent: ${SCENARIO_DATA.supplier.name} (${SCENARIO_DATA.supplier.id})`,
    actor: "Alex Rivera (Project Manager)",
    link: "/dashboard/construction/offers",
  },
  {
    id: "ACT-005",
    time: "10:05",
    type: "DELIVERY_DELAYED",
    badge: SCENARIO_DATA.delivery.reference,
    title: `Delivery delayed: ${SCENARIO_DATA.delivery.reference} (ETA moved to 2026-08-16 07:30)`,
    actor: "David Chen (Logistics Coordinator)",
    link: "/dashboard/construction/deliveries",
  },
  {
    id: "ACT-006",
    time: "10:07",
    type: "ISSUE_REPORTED",
    badge: SCENARIO_DATA.issue.reference,
    title: `Issue reported: ${SCENARIO_DATA.issue.title} (${SCENARIO_DATA.issue.severity} Severity)`,
    actor: "Mark Vance (Freight Operator)",
    link: "/dashboard/construction/issues",
  },
  {
    id: "ACT-007",
    time: "10:21",
    type: "DECISION_APPROVED",
    badge: SCENARIO_DATA.decision.reference,
    title: `Decision approved: ${SCENARIO_DATA.decision.title}`,
    actor: "Alex Rivera (Project Manager)",
    link: "/dashboard/construction/decisions",
  },
  {
    id: "ACT-008",
    time: "10:48",
    type: "INSPECTION_COMPLETED",
    badge: SCENARIO_DATA.inspection.reference,
    title: `Inspection completed: ${SCENARIO_DATA.inspection.title} (${SCENARIO_DATA.inspection.status})`,
    actor: "Sarah Jenkins, PE (Site Superintendent)",
    link: "/dashboard/construction/inspections",
  },
];
