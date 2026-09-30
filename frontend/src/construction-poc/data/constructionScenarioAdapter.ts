/**
 * CANONICAL PROJECTION of the authoritative XS-BIM construction scenario.
 *
 * Source of truth (do NOT diverge from):
 *   ankhora-cloud/internal/scenarios/xs-bim/
 *     - records.go            -> type aliases onto construction_domain models
 *     - simulator.go          -> the instantiated scenario data
 *     - simulator_test.go     -> the executable contract (values pinned here)
 *     - internal/construction/domain/models.go -> the struct definitions
 *
 * Every field below mirrors a Go struct field (camelCase of the json tag) and
 * carries the value produced by the scenario AFTER the full run completes, i.e.
 * the final state, not the creation state.
 *
 * Deliberate presentation transformations (currency formatting, humanised
 * labels, status capitalisation) do NOT live here. They belong in
 * scenarioMappers.ts so this file stays a 1:1 projection and any divergence
 * from the authoritative source remains reviewable.
 *
 * Known authoritative-source gaps are marked `@gap` and reported, not patched.
 */

export interface ScenarioDateWindow {
  start: string;
  end: string;
}

export interface ScenarioLocation {
  address: string;
  city: string;
  country: string;
}

export interface ScenarioParticipant {
  /** Canonical vault identity id. This is what every scenario record references. */
  vaultId: string;
  vaultAddress: string;
  organizationId: string;
  endpoint: string;
  /** Role this actor plays in the scenario, as evidenced by its trust-group role. */
  scenarioRole: "owner" | "member" | "reviewer";
}

export interface ConstructionProjectRecord {
  projectId: string;
  projectReference: string;
  projectName: string;
  projectType: string;
  status: string;
  currentPhase: string;
  startDate: string;
  plannedEndDate: string;
  actualEndDate: string;
  progressPercentage: number;
  statusSummary: string;
  location: ScenarioLocation;
  /** @gap never populated by the simulator — the project is not an aggregate root. */
  siteId: string;
  /** @gap never populated by the simulator. */
  requirementIds: string[];
  /** @gap never populated by the simulator; there is no ConstructionStakeholder aggregate. */
  stakeholderIds: string[];
  /** @gap never populated by the simulator. */
  milestones: string[];
}

export interface ConstructionMaterialRecord {
  materialId: string;
  materialReference: string;
  name: string;
  category: string;
  description: string;
  standard: string;
  specification: string;
  unit: string;
  origin: string;
  productionDate: string;
  batchReference: string;
  certificationReferences: string[];
  /** Final value after AcceptConstructionMaterial. */
  status: string;
}

export interface ConstructionRequirementRecord {
  requirementId: string;
  requirementReference: string;
  projectId: string;
  type: string;
  description: string;
  specification: string;
  materialId: string;
  quantity: number;
  unit: string;
  requiredDate: string;
  deliveryWindow: ScenarioDateWindow;
  /**
   * @gap The scenario assigns `requirement.SiteID` in memory (simulator.go:552)
   * but never persists it — there is no UpdateRequirement call. The API returns
   * "". Use `delivery.siteId` (persisted at creation) for the site relationship.
   */
  siteId: string;
  projectPhase: string;
  priority: string;
  /** Never mutated by the scenario. */
  status: string;
}

export interface ConstructionSupplierRecord {
  supplierId: string;
  supplierReference: string;
  name: string;
  type: string;
  description: string;
  /** Dangling reference — no contact aggregate exists in the domain. */
  contactId: string;
  materialCategories: string[];
  certificationReferences: string[];
  serviceAreas: string[];
  status: string;
}

export interface ConstructionOfferRecord {
  offerId: string;
  offerReference: string;
  requirementId: string;
  supplierId: string;
  materialId: string;
  quantity: number;
  unit: string;
  specification: string;
  unitPrice: number;
  currency: string;
  totalPrice: number;
  availabilityDate: string;
  proposedDeliveryDate: string;
  deliveryWindow: ScenarioDateWindow;
  validUntil: string;
  certificationReferences: string[];
  /** Final value after AcceptConstructionOffer. */
  status: string;
  notes: string;
}

export interface ConstructionSiteRecord {
  siteId: string;
  siteReference: string;
  projectId: string;
  name: string;
  /** Time-of-day window, not timestamps. */
  accessWindow: ScenarioDateWindow;
  accessConstraints: string[];
  storageCapacity: { value: number; unit: string };
  receivingRequirements: string[];
  inspectionRequired: boolean;
  acceptanceRequired: boolean;
  status: string;
}

export interface ConstructionDeliveryRecord {
  deliveryId: string;
  deliveryReference: string;
  projectId: string;
  requirementId: string;
  materialId: string;
  supplierId: string;
  offerId: string;
  quantity: number;
  unit: string;
  /** Final value: planned -> in_transit -> delayed -> in_transit -> received -> accepted */
  status: string;
  plannedDeliveryDate: string;
  eta: string;
  actualDeliveryDate: string;
  siteId: string;
  transportId: string;
  driverId: string;
  reference: string;
  deliveryNotes: string;
}

export interface ConstructionTransportRecord {
  transportId: string;
  transportReference: string;
  deliveryId: string;
  vehicleReference: string;
  vehicleType: string;
  /** An identity id, not a person's name — there is no driver aggregate. */
  driverId: string;
  origin: ScenarioLocation;
  destinationSiteId: string;
  /** Final value after UpdateRouteAndResumeDelivery. */
  routeReference: string;
  plannedDeparture: string;
  actualDeparture: string;
  plannedArrival: string;
  eta: string;
  actualArrival: string;
  constraints: string[];
  status: string;
  delayReason: string;
  notes: string;
}

export interface ConstructionIssueRecord {
  issueId: string;
  issueReference: string;
  projectId: string;
  title: string;
  description: string;
  type: string;
  severity: string;
  /** Final value after ResolveConstructionIssue. */
  status: string;
  reportedBy: string;
  reportedAt: string;
  /** Polymorphic reference pair, per domain convention. */
  affectedResourceType: string;
  affectedResourceId: string;
  siteId: string;
  impact: string;
  evidenceReferences: string[];
  /** Free-text role, not a participant id. */
  assignedTo: string;
  resolution: string;
  resolvedAt: string;
}

export interface ConstructionActionRecord {
  actionId: string;
  actionReference: string;
  projectId: string;
  type: string;
  resourceType: string;
  resourceId: string;
  assignedTo: string;
  status: string;
  priority: string;
  description: string;
  dueDate: string;
  completedAt: string;
  result: string;
  evidenceReferences: string[];
}

export interface ConstructionDecisionRecord {
  decisionId: string;
  decisionReference: string;
  projectId: string;
  type: string;
  requestedBy: string;
  requestDate: string;
  subject: string;
  context: string;
  evidenceReferences: string[];
  technicalAssessment: string;
  risksIdentified: string[];
  participantsConsulted: string[];
  optionsConsidered: string[];
  decision: string;
  decidedBy: string;
  decisionDate: string;
  consequence: string;
  status: string;
}

export interface ConstructionInspectionRecord {
  inspectionId: string;
  inspectionReference: string;
  projectId: string;
  deliveryId: string;
  materialId: string;
  type: string;
  zone: string;
  inspectorId: string;
  scheduledAt: string;
  completedAt: string;
  criteria: string[];
  status: string;
  result: string;
  findings: string[];
  evidenceReferences: string[];
  notes: string;
}

export interface ConstructionDocumentRecord {
  documentId: string;
  documentReference: string;
  projectId: string;
  name: string;
  category: string;
  description: string;
  cid: string;
  contentHash: string;
  version: string;
  status: string;
  uploadedBy: string;
  issuedBy: string;
  issuedAt: string;
  effectiveFrom: string;
  expiresAt: string;
  relatedResourceType: string;
  relatedResourceId: string;
}

export interface ScenarioShareEntry {
  title: string;
  entryType: string;
  downloadAllowed: boolean;
  senderUserId: string;
  senderEmail: string;
  signature: string;
  c3Cid: string;
}

export interface ScenarioThreadEvent {
  /** Monotonic ordering only. The scenario records no wall-clock time per event. */
  cursor: number;
  eventType: string;
  /** Vault identity id of the actor that appended the event. */
  actorId: string;
  idempotencyKey: string;
  /** Only `construction.action.completed` carries a payload. */
  payload?: { cid: string; contentHash: string; size: number };
}

export const CANONICAL_SCENARIO = {
  /** The three scenario actors. @gap no ConstructionStakeholder aggregate exists. */
  participants: [
    {
      vaultId: "vault_001-oem",
      vaultAddress: "oem-metro@partner.com",
      organizationId: "org_001",
      endpoint: "https://vault_oem@partner.com",
      scenarioRole: "owner"
    },
    {
      vaultId: "vault_002-michelin",
      vaultAddress: "apex-supplier@partner.com",
      organizationId: "org_002",
      endpoint: "https://vault_supplier@partner.com",
      scenarioRole: "member"
    },
    {
      vaultId: "vault_003-faa",
      vaultAddress: "regulator@faa.gov",
      organizationId: "org_003",
      endpoint: "https://vault_regulator@faa.gov",
      scenarioRole: "reviewer"
    }
  ] as ScenarioParticipant[],

  collaboration: {
    workspace: {
      name: "Metro Line 4 Expansion Program",
      description: "Civil & Structural Logistics",
      ownerId: "org_001"
    },
    channel: {
      templateId: "construction-logistics",
      title: "Site Logistics & Supply Chain",
      slotName: "Site Logistics & Supply Chain",
      slotRole: "lead-contractor",
      gated: true,
      slotOrder: 1,
      slotOwnerId: "ProjectManager-01"
    },
    thread: {
      title: "Material Delivery Tracking - Beam D-1042",
      subtitle: "Metro Line 4 Extension",
      assetType: "construction_delivery"
    },
    trustGroup: {
      name: "Logistics Coordination Group",
      groupKey: "logistics-group-key",
      keyVersion: 1
    },
    partnerTrust: [
      { vaultId: "vault_002-michelin", trustName: "supplier-logistics-trust", direction: "inbound" },
      { vaultId: "vault_003-faa", trustName: "site-inspection-trust", direction: "inbound" }
    ]
  },

  project: {
    projectId: "PRJ-001",
    projectReference: "PRJ-METRO-001",
    projectName: "Metro Line 4 Expansion",
    projectType: "Infrastructure",
    status: "active",
    currentPhase: "Phase 2 Structural",
    startDate: "2026-01-01",
    plannedEndDate: "2027-12-31",
    actualEndDate: "",
    progressPercentage: 25.0,
    statusSummary: "Civil works underway",
    location: { address: "", city: "", country: "" },
    siteId: "",
    requirementIds: [],
    stakeholderIds: [],
    milestones: []
  } as ConstructionProjectRecord,

  material: {
    materialId: "MAT-STRUCT-001",
    materialReference: "MAT-REF-88",
    name: "Precast Concrete Beam Heavy Grade",
    category: "Structural Concrete",
    description: "High load capacity precast structural beam",
    standard: "EN 13369",
    specification: "C50/60 Concrete, 12m length",
    unit: "units",
    origin: "Apex Fabrication Facility",
    productionDate: "2026-08-01",
    batchReference: "BATCH-2026-08-A",
    certificationReferences: ["CERT-MAT-881"],
    status: "accepted"
  } as ConstructionMaterialRecord,

  requirement: {
    requirementId: "REQ-STRUCT-001",
    requirementReference: "REQ-REF-1042",
    projectId: "PRJ-001",
    type: "Structural Material",
    description: "Heavy grade precast beams required for viaduct section 4",
    specification: "EN 13369 C50/60",
    materialId: "MAT-STRUCT-001",
    quantity: 12.0,
    unit: "units",
    requiredDate: "2026-08-15",
    deliveryWindow: { start: "2026-08-15T08:00:00Z", end: "2026-08-15T12:00:00Z" },
    siteId: "",
    projectPhase: "Phase 2 Structural",
    priority: "critical",
    status: "open"
  } as ConstructionRequirementRecord,

  supplier: {
    supplierId: "SUP-001",
    supplierReference: "SUP-APEX",
    name: "Apex Precast Logistics",
    type: "Manufacturer",
    description: "Specialized heavy concrete component supplier",
    contactId: "contact-apex-01",
    materialCategories: ["Structural Concrete", "Precast Beams"],
    certificationReferences: ["ISO-9001-APEX"],
    serviceAreas: ["Metropolitan Region"],
    status: "active"
  } as ConstructionSupplierRecord,

  offer: {
    offerId: "OFF-1042",
    offerReference: "OFF-REF-1042",
    requirementId: "REQ-STRUCT-001",
    supplierId: "SUP-001",
    materialId: "MAT-STRUCT-001",
    quantity: 12.0,
    unit: "units",
    specification: "C50/60 Concrete, 12m length",
    unitPrice: 4500.0,
    currency: "EUR",
    totalPrice: 54000.0,
    availabilityDate: "2026-08-10",
    proposedDeliveryDate: "2026-08-15",
    deliveryWindow: { start: "2026-08-15T08:00:00Z", end: "2026-08-15T12:00:00Z" },
    validUntil: "2026-08-30",
    certificationReferences: ["CERT-MAT-881"],
    status: "accepted",
    notes: "Includes specialized transport vehicle and driver"
  } as ConstructionOfferRecord,

  site: {
    siteId: "SITE-001",
    siteReference: "SITE-HUB-NORTH",
    projectId: "PRJ-001",
    name: "North Hub Station Site",
    accessWindow: { start: "06:00", end: "18:00" },
    accessConstraints: ["Heavy vehicle access via Gate 3 only"],
    storageCapacity: { value: 500, unit: "sqm" },
    receivingRequirements: ["Crane operator standby", "Inspector on site"],
    inspectionRequired: true,
    acceptanceRequired: true,
    status: "active"
  } as ConstructionSiteRecord,

  delivery: {
    deliveryId: "DEL-1042",
    deliveryReference: "DEL-REF-1042",
    projectId: "PRJ-001",
    requirementId: "REQ-STRUCT-001",
    materialId: "MAT-STRUCT-001",
    supplierId: "SUP-001",
    offerId: "OFF-1042",
    quantity: 12.0,
    unit: "units",
    status: "accepted",
    plannedDeliveryDate: "2026-08-15",
    eta: "2026-08-16T07:30:00Z",
    actualDeliveryDate: "2026-08-16T07:30:00Z",
    siteId: "SITE-001",
    transportId: "TR-1042",
    driverId: "driver-heavy-09",
    reference: "EXT-LOG-1042",
    deliveryNotes: "Structural beams for viaduct section 4"
  } as ConstructionDeliveryRecord,

  transport: {
    transportId: "TR-1042",
    transportReference: "TR-REF-1042",
    deliveryId: "DEL-1042",
    vehicleReference: "TRUCK-HEAVY-9",
    vehicleType: "Heavy Transporter",
    driverId: "driver-heavy-09",
    origin: { address: "Factory St 10", city: "Industrial Park", country: "DE" },
    destinationSiteId: "SITE-001",
    routeReference: "Route-B",
    plannedDeparture: "2026-08-15T06:00:00Z",
    actualDeparture: "2026-08-15T06:05:00Z",
    plannedArrival: "2026-08-15T10:00:00Z",
    eta: "2026-08-16T07:30:00Z",
    actualArrival: "2026-08-16T07:30:00Z",
    constraints: ["Road restriction on planned route M1"],
    status: "completed",
    delayReason: "Road restriction on planned route M1",
    notes: ""
  } as ConstructionTransportRecord,

  issue: {
    issueId: "ISS-1042",
    issueReference: "ISS-REF-1042",
    projectId: "PRJ-001",
    title: "Critical Beam Delivery Delayed by Route M1 Blockage",
    description:
      "Severe road restriction on highway M1 prevents heavy transport vehicle TR-1042 from reaching site on scheduled date",
    type: "logistics",
    severity: "high",
    status: "resolved",
    reportedBy: "vault_002-michelin",
    reportedAt: "2026-08-15T11:41:00Z",
    affectedResourceType: "construction_delivery",
    affectedResourceId: "DEL-1042",
    siteId: "SITE-001",
    impact: "Delay in phase 2 structural assembly",
    evidenceReferences: ["DOC-EVID-001"],
    assignedTo: "Logistics Manager",
    resolution: "Material rerouted via Route B and accepted after site inspection",
    resolvedAt: "2026-08-16T09:20:00Z"
  } as ConstructionIssueRecord,

  action: {
    actionId: "ACT-1042",
    actionReference: "ACT-REF-1042",
    projectId: "PRJ-001",
    type: "review",
    resourceType: "construction_issue",
    resourceId: "ISS-1042",
    assignedTo: "Logistics Team",
    status: "completed",
    priority: "high",
    description: "Evaluate alternative transport route B and clear axle load limits",
    dueDate: "2026-08-15",
    completedAt: "2026-08-15T13:00:00Z",
    result: "Route B detour evaluated and verified compatible with heavy vehicle constraints",
    evidenceReferences: ["DOC-EVID-001"]
  } as ConstructionActionRecord,

  decision: {
    decisionId: "DEC-1042",
    decisionReference: "DEC-REF-1042",
    projectId: "PRJ-001",
    type: "Delivery Route Change",
    requestedBy: "Logistics Manager",
    requestDate: "2026-08-15",
    subject: "Structural material delivery delay alternative route",
    context: "Route M1 blocked, ETA delayed to 2026-08-16T07:30:00Z",
    evidenceReferences: ["DOC-EVID-001"],
    technicalAssessment:
      "Alternative Route B bypasses M1 restriction; clear height and axle load requirements satisfied",
    risksIdentified: ["6-hour travel extension", "Additional transport cost"],
    participantsConsulted: ["Contractor", "Supplier", "Logistics", "Project Manager"],
    optionsConsidered: ["Wait for M1 clearance", "Use Route B detour"],
    decision: "Approve Route B alternative transport",
    decidedBy: "Project Manager",
    decisionDate: "2026-08-15",
    consequence: "Transport re-routed via Route B, delivery expected 2026-08-16T07:30:00Z",
    status: "approved"
  } as ConstructionDecisionRecord,

  inspection: {
    inspectionId: "INSP-1042",
    inspectionReference: "INSP-REF-1042",
    projectId: "PRJ-001",
    deliveryId: "DEL-1042",
    materialId: "MAT-STRUCT-001",
    type: "Receiving Inspection",
    zone: "Unloading Bay 2",
    inspectorId: "vault_003-faa",
    scheduledAt: "2026-08-16T08:00:00Z",
    completedAt: "2026-08-16T09:20:00Z",
    criteria: ["Dimensional check", "Surface crack inspection", "Mill test certificate match"],
    status: "completed",
    result: "accepted",
    findings: [
      "Beams delivered in sound structural condition",
      "No microcracks or transport damage observed"
    ],
    evidenceReferences: ["DOC-CERT-882"],
    notes: "Material verified compliant with EN 13369"
  } as ConstructionInspectionRecord,

  /** Road-restriction notice, the only C3-shared artifact. */
  evidenceDoc: {
    documentId: "DOC-EVID-001",
    documentReference: "DOC-REF-001",
    projectId: "PRJ-001",
    name: "Road Restriction Official Notice M1",
    category: "road_restriction_report",
    description: "Department of Transportation closure notice for Highway M1",
    cid: "QmRoadRestrictionReport1042",
    contentHash: "hash-restriction-notice-m1",
    version: "1.0",
    status: "issued",
    uploadedBy: "vault_002-michelin",
    issuedBy: "Highway Authority",
    issuedAt: "2026-08-15T11:00:00Z",
    effectiveFrom: "",
    expiresAt: "",
    relatedResourceType: "construction_transport",
    relatedResourceId: "TR-1042"
  } as ConstructionDocumentRecord,

  /** Inspection certificate, cited by `inspection.evidenceReferences`. */
  certDoc: {
    documentId: "DOC-CERT-882",
    documentReference: "DOC-REF-882",
    projectId: "PRJ-001",
    name: "Site Delivery Quality Inspection Certificate",
    category: "inspection_certificate",
    description: "On-site quality audit and ultrasonic testing report",
    cid: "QmInspectionCert882",
    contentHash: "hash-cert-882",
    version: "1.0",
    status: "issued",
    uploadedBy: "vault_003-faa",
    issuedBy: "Site Quality Auditor",
    issuedAt: "2026-08-16T09:15:00Z",
    effectiveFrom: "",
    expiresAt: "",
    relatedResourceType: "construction_delivery",
    relatedResourceId: "DEL-1042"
  } as ConstructionDocumentRecord,

  shareEntry: {
    title: "Road Restriction Report M1",
    entryType: "road_restriction_report",
    downloadAllowed: true,
    senderUserId: "user_supplier_01",
    senderEmail: "vault_002-michelin",
    signature: "sig-supplier-doc",
    c3Cid: "QmRoadRestrictionReport1042"
  } as ScenarioShareEntry,

  /**
   * The authoritative event sequence, in emission order. These 23 events are
   * mirrored 1:1 into TraceCore as commits (RepoID PRJ-001, Branch thread.ID),
   * so this array is the scenario's full trace history.
   */
  threadEvents: [
    { cursor: 1, eventType: "construction.project.phase.started", actorId: "vault_001-oem", idempotencyKey: "evt-proj-started" },
    { cursor: 2, eventType: "construction.requirement.created", actorId: "vault_001-oem", idempotencyKey: "evt-req-created" },
    { cursor: 3, eventType: "construction.supplier.invited", actorId: "vault_001-oem", idempotencyKey: "evt-sup-invited" },
    { cursor: 4, eventType: "construction.offer.submitted", actorId: "vault_002-michelin", idempotencyKey: "evt-off-submitted" },
    { cursor: 5, eventType: "construction.offer.accepted", actorId: "vault_001-oem", idempotencyKey: "evt-off-accepted" },
    { cursor: 6, eventType: "construction.supplier.confirmed", actorId: "vault_002-michelin", idempotencyKey: "evt-sup-confirmed" },
    { cursor: 7, eventType: "construction.delivery.requested", actorId: "vault_001-oem", idempotencyKey: "evt-del-requested" },
    { cursor: 8, eventType: "construction.transport.assigned", actorId: "vault_002-michelin", idempotencyKey: "evt-tr-assigned" },
    { cursor: 9, eventType: "construction.transport.accepted", actorId: "vault_002-michelin", idempotencyKey: "evt-tr-accepted" },
    { cursor: 10, eventType: "construction.transport.departed", actorId: "vault_002-michelin", idempotencyKey: "evt-tr-departed" },
    { cursor: 11, eventType: "construction.transport.constraint.reported", actorId: "vault_002-michelin", idempotencyKey: "evt-tr-constraint" },
    { cursor: 12, eventType: "construction.transport.eta.updated", actorId: "vault_002-michelin", idempotencyKey: "evt-tr-eta-updated" },
    { cursor: 13, eventType: "construction.delivery.delayed", actorId: "vault_001-oem", idempotencyKey: "evt-del-delayed" },
    { cursor: 14, eventType: "construction.issue.reported", actorId: "vault_002-michelin", idempotencyKey: "evt-iss-reported" },
    { cursor: 15, eventType: "construction.action.created", actorId: "vault_001-oem", idempotencyKey: "evt-act-created" },
    {
      cursor: 16,
      eventType: "construction.action.completed",
      actorId: "vault_001-oem",
      idempotencyKey: "evt-act-completed",
      payload: { cid: "QmRoadRestrictionReport1042", contentHash: "hash-restriction-notice-m1", size: 2048 }
    },
    { cursor: 17, eventType: "construction.decision.proposed", actorId: "vault_001-oem", idempotencyKey: "evt-dec-proposed" },
    { cursor: 18, eventType: "construction.decision.approved", actorId: "vault_001-oem", idempotencyKey: "evt-dec-approved" },
    { cursor: 19, eventType: "construction.delivery.updated", actorId: "vault_002-michelin", idempotencyKey: "evt-del-updated" },
    { cursor: 20, eventType: "construction.delivery.received", actorId: "vault_002-michelin", idempotencyKey: "evt-del-received" },
    { cursor: 21, eventType: "construction.inspection.completed", actorId: "vault_003-faa", idempotencyKey: "evt-insp-completed" },
    { cursor: 22, eventType: "construction.issue.resolved", actorId: "vault_001-oem", idempotencyKey: "evt-iss-resolved" },
    { cursor: 23, eventType: "construction.material.accepted", actorId: "vault_001-oem", idempotencyKey: "evt-mat-accepted" }
  ] as ScenarioThreadEvent[]
} as const;
