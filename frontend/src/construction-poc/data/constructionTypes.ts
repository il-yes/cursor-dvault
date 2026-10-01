/**
 * View-model types for the Construction Functional Data Boundary.
 *
 * These are the contracts the presentation layer already depends on. They are
 * intentionally expressed in UI terms, not as a second domain model, and are
 * shared by every ConstructionDataProvider implementation.
 */

export interface CreateProjectParams {
  code: string;
  name: string;
  type: string;
  location: string;
}

export interface SupplierOfferItem {
  offerId: string;
  offerReference: string;
  supplierId: string;
  supplierName: string;
  isVerified: boolean;
  supplierStatus: string;
  totalPrice: string;
  unitPrice: string;
  currency: string;
  specification: string;
  quantity: number;
  unit: string;
  deliveryWindow: string;
  offerValidity: string;
  leadTime: string;
  complianceScore: number;
  complianceSpecVersion: string;
  certifications: string[];
  offerStatus: string;
  statusBanner: string;
  notes?: string;
  isPrimary: boolean;
}

export interface ProcurementData {
  id: string;
  requirementId: string;
  projectId: string;
  projectCode: string;
  projectName: string;
  /** Canonical `requirementReference` */
  code: string;
  materialId: string;
  materialName: string;
  materialStandard: string;
  specification: string;
  quantity: number;
  unit: string;
  targetPhase: string;
  status: string;
  requiredDate: string;
  priority: string;
  siteId: string;
  siteName: string;
  
  /* Summary Metrics */
  lowestBidTotal: string;
  lowestBidSavings: string;
  fastestDeliveryDate: string;
  fastestDeliveryStatus: string;

  /* Primary Accepted Offer (OFF-1042) */
  offerId: string;
  offerReference: string;
  supplierId: string;
  supplierName: string;
  supplierStatus: string;
  supplierCertifications: string[];
  totalPrice: string;
  unitPrice: string;
  currency: string;
  proposedDeliveryDate: string;
  availabilityDate: string;
  validUntil: string;
  offerStatus: string;

  /* Primary & Secondary Offers for Stitch layout */
  primaryOffer: SupplierOfferItem;
  secondaryOffer: SupplierOfferItem;
  offers: SupplierOfferItem[];

  /* Stitch Material Requirement presentation fields */
  projectDisplayCode?: string;
  materialSubtitle?: string;
  materialDescription?: string;
  specificationGrade?: string;
  specificationType?: string;
  quantityBeamsText?: string;
  quantityTonsText?: string;
  requiredDateSlotText?: string;
  siteCityText?: string;
  siteLocationText?: string;
  constructionPhaseText?: string;
  invitedCount?: number;
  offersReceivedCount?: number;
  selectedCount?: number;
  contractValueText?: string;
  unitPriceText?: string;
  proposedDeliveryText?: string;
}

/* DecisionData                                                         */

/**
 * A labelled fact inside an option card.
 *
 * `value` is deliberately optional: several slots the Stitch composition shows
 * (per-option delay, detour surcharge, commercial term) have no canonical field
 * behind them. Leaving `value` undefined is how the mapper says "the scenario
 * carries no such datum", and it lets the view render an explicit unavailable
 * state instead of a plausible-looking fabrication.
 */
export interface DecisionFact {
  label: string;
  value?: string;
}

export interface DecisionOptionData {
  /** 1-based ordinal, as displayed in the Stitch option badge. */
  position: number;
  /** Canonical `optionsConsidered` entry, verbatim. */
  title: string;
  isSelected: boolean;
  /** Derived: a finalized decision selects exactly one option. */
  statusLabel: string;
  /** The two-cell impact grid shown on a non-selected option. */
  impact: DecisionFact[];
  /** Highlight rows shown on the selected option. */
  highlights: DecisionFact[];
  /** Canonical `technicalAssessment` — only the selected option carries one. */
  rationale?: string;
}

export interface DecisionParticipantData {
  /** Derived from the canonical role string; there is no person aggregate. */
  initials: string;
  /** Canonical `participantsConsulted` entry, verbatim. */
  role: string;
  /** True when the role matches canonical `requestedBy`. */
  isRequester: boolean;
  /**
   * @gap The domain records consultation as free-text roles. There is no
   * role -> organisation mapping, so `organization` is always undefined here.
   */
  organization?: string;
  /** @gap No per-participant decision or sign-off state is recorded. */
  signoffLabel?: string;
  /** @gap No per-participant timestamp is recorded. */
  signedAtLabel?: string;
}

/**
 * Page-facing projection of the canonical ConstructionDecision aggregate, shaped
 * to the Stitch DEC-1042 composition.
 */
export interface DecisionData {
  id: string;
  reference: string;
  type: string;
  status: string;
  statusLabel: string;
  /** Calendar date of the decision, canonical `decisionDate`. */
  decidedOnLabel: string;
  /**
   * @gap `decisionDate` is a bare calendar date, so the composition's clock time
   * has no source. False tells the view to say so instead of implying one.
   */
  decidedHasTime: boolean;
  requestedBy: string;
  decidedBy: string;
  subject: string;
  context: string;
  technicalAssessment: string;
  risksIdentified: string[];
  evidenceReferences: string[];
  related: {
    issueId: string;
    issueReference: string;
    deliveryId: string;
    deliveryReference: string;
    transportId: string;
    transportReference: string;
    /**
     * The third Stitch chip is a sector/zone facet. @gap the project carries no
     * sector taxonomy, so canonical `type` stands in for it.
     */
    qualifierLabel: string;
  };
  route: {
    reference: string;
    /** Derived from canonical planned vs actual arrival, when the two differ. */
    slipLabel?: string;
    arrivalLabel: string;
  };
  options: DecisionOptionData[];
  selectedIndex: number | null;
  participants: DecisionParticipantData[];
  participantCount: number;
  outcome: {
    title: string;
    body: string;
  };
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

    /* Stitch Delivery Detail fields */
    statusTag?: string;
    lastUpdatedAge?: string;
    headlineTitle?: string;
    headlineQuantity?: string;
    alertTitle?: string;
    alertDescription?: string;
    delayBadgeLabel?: string;
    plannedDateLabel?: string;
    plannedTimeLabel?: string;
    revisedDateLabel?: string;
    revisedTimeLabel?: string;
    primaryDelayFactor?: string;
    materialGradeText?: string;
    supplierNameText?: string;
    supplierCodeText?: string;
    contractOfferCodeText?: string;
    destinationSiteNameText?: string;
    destinationBayText?: string;
    carrierAgencyText?: string;
    carrierVehicleText?: string;
    linkedIncidentReference?: string;
    linkedIncidentLabel?: string;
    telemetryStatusText?: string;
    telemetryLocationText?: string;
    telemetryStateText?: string;
    timelineStageText?: string;
    timelineStages?: Array<{
      stageNumber: number;
      title: string;
      subtitle: string;
      state: "completed" | "active" | "pending" | "future";
      badge?: string;
    }>;
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
    currency: string;
    proposedDeliveryDate: string;
  };
  site: {
    id: string;
    name: string;
    accessWindow?: { start: string; end: string };
    accessConstraints?: string[];
    storageCapacity?: { value: number; unit: string };
    receivingRequirements?: string[];
    inspectionRequired?: boolean;
    acceptanceRequired?: boolean;
  };
  transport: {
    id: string;
    reference: string;
    /** Canonical `vehicleType`. */
    vehicle: string;
    vehicleReference: string;
    /** Canonical `driverId` — an identity id, there is no driver aggregate. */
    driverId: string;
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

    /* Stitch Transport & Delay fields */
    statusTag?: string;
    alertTitle?: string;
    alertConstraint?: string;
    originalEtaLabel?: string;
    revisedEtaLabel?: string;
    delayDurationLabel?: string;
    gpsSyncAge?: string;
    clearanceConstraint?: string;
    blockedRouteName?: string;
    bypassRouteName?: string;
    primaryIncidentText?: string;
    sitePhaseImpact?: string;
    reportedByText?: string;
    driverInitials?: string;
    driverName?: string;
    carrierName?: string;
    detourStatusTag?: string;
    detourCorridorText?: string;
    deltaDistanceText?: string;
    scheduleModificationText?: string;
    clearanceAuthorityText?: string;
    decisionReference?: string;
  };
  issue?: {
    id: string;
    reference: string;
    title: string;
    description: string;
    type: string;
    severity: string;
    status: string;
    reportedBy: string;
    reportedAt: string;
    impact: string;
    assignedTo: string;
    resolution?: string;
    resolvedAt?: string;
    evidenceReferences?: string[];

    /* Stitch Issue ISS-1042 presentation fields */
    categoryBadge?: string;
    severityTag?: string;
    statusTag?: string;
    mitigationLabel?: string;
    mitigationActionText?: string;
    zoneContextText?: string;
    affectedResourceCode?: string;
    affectedResourceDescription?: string;
    projectPhaseText?: string;
    projectZoneText?: string;
    rootCauseTitle?: string;
    rootCauseSecondary?: string;
    reporterInitials?: string;
    reporterNameText?: string;
    reporterTimeText?: string;
    stakeholders?: Array<{ name: string; role: string; icon: string }>;
    impactBadgeText?: string;
    primaryImpactTitle?: string;
    primaryImpactSubtitle?: string;
    primaryImpactDescription?: string;
    rescheduleLabel?: string;
    reschedulePlannedLabel?: string;
    rescheduleAdjustedLabel?: string;
    rescheduleNotes?: string;
    evidenceFiles?: Array<{
      filename: string;
      metadata: string;
      icon: string;
      size: string;
    }>;
    threadUpdateCount?: number;
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
  provenance?: ProvenanceOverviewData;
  inspectionPageData?: InspectionPageData;
}

export interface InspectionChecklistItem {
  id: number;
  title: string;
  description: string;
  isPassed: boolean;
}

export interface InspectionEvidencePhoto {
  id: string;
  label: string;
  imageAlt: string;
}

export interface InspectionEvidenceDocument {
  id: string;
  icon: string;
  title: string;
  metadata: string;
}

export interface InspectionTimelineItem {
  id: string;
  category: string;
  title: string;
  timestamp: string;
  isTerminal?: boolean;
}

export interface InspectionPageData {
  id: string;
  reference: string;
  statusTag: string;
  title: string;
  inspectorText: string;
  checklist: InspectionChecklistItem[];
  evidencePhotos: InspectionEvidencePhoto[];
  evidenceDocuments: InspectionEvidenceDocument[];
  timeline: InspectionTimelineItem[];
}

export interface ProvenanceMilestoneItem {
  stepNumber: number;
  category: string;
  timestamp: string;
  title: string;
  description: string;
  nodeType: "cause" | "ticket" | "verification" | "consensus" | "execution" | "receipt" | "qc" | "handover";
  evidenceChip?: {
    filename: string;
    badge: string;
  };
  signatories?: string[];
}

export interface ProvenanceAuditLayer {
  id: "construction" | "collaboration" | "milestones";
  tabLabel: string;
  title: string;
  description: string;
  footerLeft: string;
  statusRight: string;
}

export interface ProvenanceOverviewData {
  rootCauseTitle: string;
  rootCauseDescription: string;
  delayDurationLabel: string;
  originalEtaLabel: string;
  actualSiteGateLabel: string;
  scheduleDeltaLabel: string;
  milestones: ProvenanceMilestoneItem[];
  materialAcceptance: {
    materialCode: string;
    inspectionCode: string;
    lotLabel: string;
    beamCount: number;
    status: string;
  };
  auditLayers: ProvenanceAuditLayer[];
}
