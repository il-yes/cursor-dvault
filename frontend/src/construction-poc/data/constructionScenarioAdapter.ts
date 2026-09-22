/**
 * Presentation Data Adapter for Authoritative XS-BIM Construction Scenario.
 * 
 * Note: This file provides presentation data structures matching the authoritative
 * XS-BIM scenario (PRJ-001, REQ-STRUCT-001, DEL-1042, etc.) for UI views whose backend
 * endpoints are not yet exposed in AppAPI.
 * 
 * This is NOT a second domain model or persistence engine.
 */

export interface ConstructionProject {
  id: string;
  code: string;
  name: string;
  type: string;
  status: "ACTIVE" | "COMPLETED" | "ON_HOLD";
  location: string;
  budgetSpentPercent: number;
  scheduleDay: number;
  scheduleStatus: "On Track" | "Delayed" | "Ahead";
  complianceScore: number;
  description: string;
}

export interface MaterialRequirement {
  id: string;
  projectId: string;
  code: string;
  materialId: string;
  materialName: string;
  specification: string;
  quantity: number;
  unit: string;
  targetPhase: string;
  status: "fulfilled_delayed" | "pending" | "fulfilled";
  requiredDate: string;
}

export interface SupplierOffer {
  id: string;
  requirementId: string;
  supplierId: string;
  supplierName: string;
  offerReference: string;
  totalPrice: string;
  unitPrice: string;
  promisedDeliveryDate: string;
  status: "ACCEPTED" | "PENDING" | "REJECTED";
}

export interface Delivery {
  id: string;
  reference: string;
  projectId: string;
  requirementId: string;
  materialId: string;
  supplierId: string;
  offerId: string;
  quantity: number;
  unit: string;
  status: "delayed" | "in_transit" | "arrived" | "received" | "inspected";
  plannedDeliveryDate: string;
  eta: string;
  actualDeliveryDate?: string;
  siteId: string;
  siteName: string;
  transportId: string;
  deliveryNotes: string;
}

export interface Transport {
  id: string;
  deliveryId: string;
  vehicle: string;
  driver: string;
  origin: string;
  destination: string;
  route: string;
  plannedDeparture: string;
  actualDeparture: string;
  eta: string;
  constraints: string;
  delayReason: string;
  status: "delayed" | "rerouted" | "in_transit" | "delivered";
}

export interface ConstructionIssue {
  id: string;
  reference: string;
  deliveryId: string;
  transportId: string;
  title: string;
  description: string;
  severity: "CRITICAL" | "HIGH" | "MEDIUM" | "LOW";
  status: "OPEN" | "RESOLVED";
  reportedAt: string;
  evidenceDocId: string;
}

export interface EvidenceDocument {
  id: string;
  code: string;
  title: string;
  fileName: string;
  fileSize: string;
  hash: string;
  uploadedAt: string;
  author: string;
}

export interface ConstructionDecision {
  id: string;
  reference: string;
  issueId: string;
  title: string;
  proposedAction: string;
  impactSummary: string;
  status: "PROPOSED" | "APPROVED" | "EXECUTED";
  proposedBy: string;
  approvedBy?: string;
  approvedAt?: string;
  actionTaken: string;
}

export interface Inspection {
  id: string;
  reference: string;
  deliveryId: string;
  materialId: string;
  title: string;
  inspector: string;
  inspectionDate: string;
  status: "PASSED" | "FAILED" | "PENDING";
  ultrasonicWeldIntegrity: string;
  dimensionalCompliance: string;
  notes: string;
}

export interface Stakeholder {
  id: string;
  name: string;
  role: string;
  organization: string;
  trustStatus: "VERIFIED" | "PENDING";
  email: string;
  phone: string;
}

export interface TraceMilestone {
  timestamp: string;
  layer: "Construction" | "C3 Collaboration" | "TraceCore";
  event: string;
  entityId: string;
  actor: string;
}

export const SCENARIO_DATA = {
  project: {
    id: "PRJ-001",
    code: "PRJ-001",
    name: "Metro Line 4 Expansion",
    type: "INFRASTRUCTURE",
    status: "ACTIVE",
    location: "South Corridor - Segment 3, Chicago, IL",
    budgetSpentPercent: 62,
    scheduleDay: 142,
    scheduleStatus: "Delayed",
    complianceScore: 98,
    description: "Expansion of Metro Line 4 heavy rail transit corridor including underground tunneling, station structures, and elevated viaducts."
  } as ConstructionProject,

  requirement: {
    id: "REQ-STRUCT-001",
    projectId: "PRJ-001",
    code: "REQ-STRUCT-001",
    materialId: "MAT-STRUCT-001",
    materialName: "High-Strength Structural Steel Beams (Grade A992)",
    specification: "W18x86 Heavy Flange Beams, ASTM A992 certified, anti-corrosion primer coated.",
    quantity: 120,
    unit: "Metric Tons",
    targetPhase: "Phase 4 - Structure Pier Support",
    status: "fulfilled_delayed",
    requiredDate: "2026-08-15"
  } as MaterialRequirement,

  supplier: {
    id: "SUP-001",
    name: "Apex Steel Fabrication Ltd."
  },

  offer: {
    id: "OFF-1042",
    requirementId: "REQ-STRUCT-001",
    supplierId: "SUP-001",
    supplierName: "Apex Steel Fabrication Ltd.",
    offerReference: "OFF-1042",
    totalPrice: "$142,500.00",
    unitPrice: "$1,187.50 / Ton",
    promisedDeliveryDate: "2026-08-15",
    status: "ACCEPTED"
  } as SupplierOffer,

  delivery: {
    id: "DEL-1042",
    reference: "DEL-1042",
    projectId: "PRJ-001",
    requirementId: "REQ-STRUCT-001",
    materialId: "MAT-STRUCT-001",
    supplierId: "SUP-001",
    offerId: "OFF-1042",
    quantity: 120,
    unit: "Metric Tons",
    status: "delayed",
    plannedDeliveryDate: "2026-08-15",
    eta: "2026-08-16 07:30",
    actualDeliveryDate: "2026-08-16 07:51",
    siteId: "SITE-SOUTH-01",
    siteName: "Site Alpha - South Pier Foundation",
    transportId: "TR-1042",
    deliveryNotes: "Critical structural load for Pier 4 framework. Delayed en route due to bridge load restriction on M1."
  } as Delivery,

  transport: {
    id: "TR-1042",
    deliveryId: "DEL-1042",
    vehicle: "Heavy Hauler Truck #88 (48-Ton Multi-Axle Trailer)",
    driver: "Mark Vance (CDL Class-A)",
    origin: "Apex Plant #2, Gary, IN",
    destination: "Site Alpha - South Pier, Chicago, IL",
    route: "M1 Highway Northbound → Rerouted via Highway B",
    plannedDeparture: "2026-08-15 06:00",
    actualDeparture: "2026-08-15 06:15",
    eta: "2026-08-16 07:30",
    constraints: "Axle weight limit max 35T on M1 Km 42 bridge structure.",
    delayReason: "Road Restriction on M1 Highway - Weight limit enforced.",
    status: "rerouted"
  } as Transport,

  issue: {
    id: "ISS-1042",
    reference: "ISS-1042",
    deliveryId: "DEL-1042",
    transportId: "TR-1042",
    title: "Road Restriction on M1 Highway - Axle Weight Limit",
    description: "Illinois DOT issued emergency weight reduction to 35T on M1 Km 42 bridge structure. Transport TR-1042 total load is 48T, blocking primary transit route.",
    severity: "CRITICAL",
    status: "RESOLVED",
    reportedAt: "2026-08-15 11:41",
    evidenceDocId: "DOC-EVID-001"
  } as ConstructionIssue,

  evidenceDoc: {
    id: "DOC-EVID-001",
    code: "DOC-EVID-001",
    title: "DOT Emergency Road Restriction Notice & Permit Report",
    fileName: "DOT_M1_Restriction_Permit_2026.pdf",
    fileSize: "2.4 MB",
    hash: "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    uploadedAt: "2026-08-15 11:45",
    author: "Mark Vance (Freight Operator)"
  } as EvidenceDocument,

  decision: {
    id: "DEC-1042",
    reference: "DEC-1042",
    issueId: "ISS-1042",
    title: "Approve Alternative Route B & Heavy Load Permit Reroute",
    proposedAction: "Reroute transport via Highway B with state police escort, splitting secondary axle load at intermediate weigh station.",
    impactSummary: "Foundation phase delayed by 1 day. Delivery ETA moved from Aug 15 to Aug 16, 07:30.",
    status: "EXECUTED",
    proposedBy: "David Chen (Logistics Coordinator)",
    approvedBy: "Alex Rivera (Project Manager)",
    approvedAt: "2026-08-15 14:20",
    actionTaken: "Alternative Route B approved and executed. Transport rerouted successfully."
  } as ConstructionDecision,

  inspection: {
    id: "INSP-1042",
    reference: "INSP-1042",
    deliveryId: "DEL-1042",
    materialId: "MAT-STRUCT-001",
    title: "Structural Steel Weld & Dimensional Acceptance Inspection",
    inspector: "Sarah Jenkins, PE (Lead Site Superintendent)",
    inspectionDate: "2026-08-16 09:20",
    status: "PASSED",
    ultrasonicWeldIntegrity: "100% Pass - Zero structural defects detected",
    dimensionalCompliance: "Conforms strictly to ASTM A992 & REQ-STRUCT-001 specs",
    notes: "Material received in prime condition post-reroute. Approved for immediate hoisting on Pier 4 framework."
  } as Inspection,

  stakeholders: [
    {
      id: "STK-001",
      name: "Alex Rivera",
      role: "Project Manager",
      organization: "Metro Transit Authority",
      trustStatus: "VERIFIED",
      email: "arivera@metrotransit.org",
      phone: "+1 (312) 555-0142"
    },
    {
      id: "STK-002",
      name: "David Chen",
      role: "Logistics Coordinator",
      organization: "Apex Steel Fabrication Ltd.",
      trustStatus: "VERIFIED",
      email: "dchen@apexsteel.com",
      phone: "+1 (312) 555-0188"
    },
    {
      id: "STK-003",
      name: "Mark Vance",
      role: "Freight Driver / Operator",
      organization: "Vance Logistics Group",
      trustStatus: "VERIFIED",
      email: "mvance@vancelogistics.com",
      phone: "+1 (312) 555-0204"
    },
    {
      id: "STK-004",
      name: "Sarah Jenkins, PE",
      role: "Lead Site Superintendent & QA",
      organization: "Ankhora Construction Management",
      trustStatus: "VERIFIED",
      email: "sjenkins@ankhora-cm.com",
      phone: "+1 (312) 555-0199"
    }
  ] as Stakeholder[],

  traceMilestones: [
    {
      timestamp: "2026-08-15 08:03",
      layer: "Construction",
      event: "Transport TR-1042 accepted by carrier",
      entityId: "TR-1042",
      actor: "Mark Vance"
    },
    {
      timestamp: "2026-08-15 11:41",
      layer: "C3 Collaboration",
      event: "Construction issue ISS-1042 logged: M1 Road restriction",
      entityId: "ISS-1042",
      actor: "Mark Vance"
    },
    {
      timestamp: "2026-08-15 11:45",
      layer: "C3 Collaboration",
      event: "Evidence document DOC-EVID-001 attached to thread",
      entityId: "DOC-EVID-001",
      actor: "Mark Vance"
    },
    {
      timestamp: "2026-08-15 12:05",
      layer: "Construction",
      event: "Delivery DEL-1042 status updated to delayed. ETA: 2026-08-16 07:30",
      entityId: "DEL-1042",
      actor: "David Chen"
    },
    {
      timestamp: "2026-08-15 12:17",
      layer: "C3 Collaboration",
      event: "Site Superintendent acknowledged delivery delay notice",
      entityId: "DEL-1042",
      actor: "Sarah Jenkins"
    },
    {
      timestamp: "2026-08-15 13:02",
      layer: "Construction",
      event: "Alternative Route B proposal DEC-1042 created",
      entityId: "DEC-1042",
      actor: "David Chen"
    },
    {
      timestamp: "2026-08-15 14:20",
      layer: "Construction",
      event: "Decision DEC-1042 approved by Project Manager",
      entityId: "DEC-1042",
      actor: "Alex Rivera"
    },
    {
      timestamp: "2026-08-16 07:51",
      layer: "TraceCore",
      event: "Delivery DEL-1042 received at Site Alpha South Pier",
      entityId: "DEL-1042",
      actor: "Sarah Jenkins"
    },
    {
      timestamp: "2026-08-16 09:20",
      layer: "TraceCore",
      event: "Inspection INSP-1042 completed: Structural Steel Accepted",
      entityId: "INSP-1042",
      actor: "Sarah Jenkins"
    }
  ] as TraceMilestone[]
};
