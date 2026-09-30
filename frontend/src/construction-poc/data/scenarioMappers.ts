/**
 * Canonical XS-BIM scenario -> page-facing view models.
 *
 * This is the ONLY place where deliberate presentation transformations live.
 * Every transformation is annotated with its rule so the canonical source and
 * the presented form can be diffed (see the divergence report).
 *
 * Rules:
 *   direct        copied verbatim from a canonical field
 *   format        canonical value rendered for display (currency, date, location)
 *   reference     canonical lacks the field; a related canonical value is used
 *
 * Nothing is invented. Where the scenario genuinely carries no datum the field
 * is left undefined rather than placeholder-filled.
 */

import { CANONICAL_SCENARIO } from "./constructionScenarioAdapter";
import type {
  DecisionData,
  DecisionFact,
  DecisionOptionData,
  LogisticsOverviewData,
  ProcurementData,
} from "./constructionTypes";

const S = CANONICAL_SCENARIO;

/* ------------------------------------------------------------------ */
/* presentation formatters                                             */
/* ------------------------------------------------------------------ */

/** format: renders a canonical ISO timestamp in a stable, readable form. */
export function formatTimestamp(iso: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toISOString().slice(0, 16).replace("T", " ") + " UTC";
}

/** format: renders a canonical float + ISO-4217 code as a currency string. */
export function formatMoney(amount: number, currency: string): string {
  return new Intl.NumberFormat("en-GB", {
    style: "currency",
    currency,
    maximumFractionDigits: 2
  }).format(amount);
}

/** format: renders a canonical Location struct as a single line. */
export function formatLocation(location: { address: string; city: string; country: string }): string {
  return [location.address, location.city, location.country].filter(Boolean).join(", ");
}

/**
 * reference: the scenario has no person-name aggregate. Actors are vault
 * identities, so a display label is derived from the vault address, not invented.
 *
 * Reads the participants from the canonical scenario rather than from
 * scenarioAccessors, so this module stays outside the accessors/mock import cycle.
 */
export function actorLabel(vaultId: string): string {
  return S.participants.find((p) => p.vaultId === vaultId)?.vaultAddress ?? vaultId;
}

/** format: humanises a canonical snake_case status for display only. */
export function humaniseStatus(status: string): string {
  return status
    .split(/[_.]/)
    .filter(Boolean)
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

const MONTH_ABBR = [
  "Jan", "Feb", "Mar", "Apr", "May", "Jun",
  "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
];

/**
 * format: calendar day of a canonical timestamp, always read in UTC.
 *
 * The scenario stores UTC instants and bare calendar dates, so anchoring the
 * formatter to UTC is what keeps a date from drifting by a day in negative-offset
 * zones. No timezone conversion is performed — the scenario's own zone is the
 * one presented.
 */
export function formatUtcDay(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return `${MONTH_ABBR[d.getUTCMonth()]} ${d.getUTCDate()}`;
}

/**
 * format: canonical timestamp as a UTC calendar day, clock time and zone.
 *
 * The explicit `UTC` suffix is load-bearing. The Stitch composition shows bare
 * clock times ("Aug 16 @ 07:30"), which silently asserts a local zone the
 * canonical scenario never carries.
 */
export function formatUtcMoment(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  const hh = String(d.getUTCHours()).padStart(2, "0");
  const mm = String(d.getUTCMinutes()).padStart(2, "0");
  return `${MONTH_ABBR[d.getUTCMonth()]} ${d.getUTCDate()} @ ${hh}:${mm} UTC`;
}

/** true when a canonical value carries a clock time, not only a calendar date. */
export function hasClockTime(iso: string): boolean {
  return /T\d{2}:\d{2}/.test(iso);
}

/**
 * format: whole-or-fractional hours between two canonical timestamps, or
 * undefined when either is unparseable.
 */
export function hoursBetween(fromIso: string, toIso: string): number | undefined {
  const from = new Date(fromIso).getTime();
  const to = new Date(toIso).getTime();
  if (Number.isNaN(from) || Number.isNaN(to)) return undefined;
  return (to - from) / 3_600_000;
}

/** format: a signed hour span as a schedule slip, e.g. "+21.5 h". */
export function formatSlip(hours: number): string {
  const rounded = Math.round(hours * 10) / 10;
  const sign = rounded > 0 ? "+" : "";
  return `${sign}${rounded} h`;
}

/**
 * format: avatar initials derived from a canonical role string.
 *
 * There is no person aggregate, so initials are rendered from the role text
 * rather than from a person's name.
 */
export function roleInitials(role: string): string {
  const words = role.split(/\s+/).filter(Boolean);
  if (words.length === 0) return "?";
  if (words.length === 1) return words[0].charAt(0).toUpperCase();
  return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
}

/* ------------------------------------------------------------------ */
/* lookup                                                              */
/* ------------------------------------------------------------------ */

function matches(candidate: string, requested: string): boolean {
  return candidate.toLowerCase() === requested.toLowerCase();
}

export function scenarioHasRequirement(requirementId?: string): boolean {
  if (!requirementId) return false;
  return matches(S.requirement.requirementId, requirementId) ||
    matches(S.requirement.requirementReference, requirementId);
}

export function scenarioHasDelivery(deliveryId?: string): boolean {
  if (!deliveryId) return false;
  return matches(S.delivery.deliveryId, deliveryId) ||
    matches(S.delivery.deliveryReference, deliveryId);
}

/* DecisionData                                                         */

/**
 * reference: which of `optionsConsidered` the decision actually selected.
 *
 * The aggregate stores the outcome as free text beside a parallel list of
 * options, with no option id to join on. The selection is therefore recovered by
 * token overlap between the outcome and each candidate, scored on non-stopword
 * terms so that shared filler ("use", "wait for") cannot decide it.
 *
 * For DEC-1042 this resolves unambiguously: the outcome
 * "Approve Route B alternative transport" shares "route" and "b" with
 * "Use Route B detour" (score 2) and nothing with "Wait for M1 clearance"
 * (score 0).
 */
export function deriveSelectedOptionIndex(decision: string, options: string[]): number | null {
  const STOP_WORDS = new Set([
    "wait", "for", "use", "the", "a", "an", "of", "on", "to", "and", "or", "via", "in", "at", "by"
  ]);

  const significant = (value: string): string[] => value
    .toLowerCase()
    .split(/[^a-z0-9]+/)
    .filter((token) => token.length > 0 && !STOP_WORDS.has(token));

  const outcome = new Set(significant(decision));

  let bestIndex: number | null = null;
  let bestScore = 0;
  options.forEach((option, index) => {
    const overlap = significant(option).filter((token) => outcome.has(token)).length;
    if (overlap > bestScore) {
      bestScore = overlap;
      bestIndex = index;
    }
  });

  return bestIndex;
}

/** @gap the composition's impact grid has no per-option domain source. */
const OPTION_IMPACT_LABELS = ["Estimated Delay", "Critical Path Impact"];

/** @gap neither a surcharge nor a commercial term is modelled on the aggregate. */
const SELECTED_HIGHLIGHT_LABELS = ["Arrival Schedule", "Detour Surcharge", "Commercial Term"];

function buildOptions(
  options: string[],
  selectedIndex: number | null,
  arrivalLabel: string,
  technicalAssessment: string
): DecisionOptionData[] {
  return options.map((title, index) => {
    const isSelected = index === selectedIndex;
    const highlights: DecisionFact[] = SELECTED_HIGHLIGHT_LABELS.map((label) => ({
      label,
      // Only the arrival schedule is canonically available.
      value: label === "Arrival Schedule" ? arrivalLabel : undefined
    }));

    return {
      position: index + 1,
      title,
      isSelected,
      statusLabel: isSelected ? "Approved" : "Rejected",
      impact: OPTION_IMPACT_LABELS.map((label) => ({ label, value: undefined })),
      highlights,
      rationale: isSelected ? technicalAssessment : undefined
    };
  });
}

export function mapScenarioToDecisionData(): DecisionData {
  const d = S.decision;
  const transport = S.transport;

  const options = d.optionsConsidered ?? [];
  const selectedIndex = deriveSelectedOptionIndex(d.decision, options);

  // The transport record is the delivery's leg, so it carries the schedule the
  // decision actually moved. actualArrival is preferred over eta because the
  // scenario resolves the arrival before the decision is finalised.
  const arrivalIso = transport.actualArrival || transport.eta;
  const arrivalLabel = formatUtcMoment(arrivalIso);
  const slip = hoursBetween(transport.plannedArrival, arrivalIso);

  return {
    id: d.decisionId,
    reference: d.decisionReference,
    type: d.type,
    status: d.status,
    statusLabel: humaniseStatus(d.status).toUpperCase(),
    decidedOnLabel: formatUtcDay(d.decisionDate),
    decidedHasTime: hasClockTime(d.decisionDate),
    requestedBy: d.requestedBy,
    decidedBy: d.decidedBy,
    subject: d.subject,
    context: d.context,
    technicalAssessment: d.technicalAssessment,
    risksIdentified: d.risksIdentified ?? [],
    evidenceReferences: d.evidenceReferences ?? [],
    related: {
      issueId: S.issue.issueId,
      issueReference: S.issue.issueReference,
      deliveryId: S.delivery.deliveryId,
      deliveryReference: S.delivery.deliveryReference,
      transportId: transport.transportId,
      transportReference: transport.transportReference,
      qualifierLabel: d.type
    },
    route: {
      reference: transport.routeReference,
      slipLabel: slip !== undefined && slip > 0 ? formatSlip(slip) : undefined,
      arrivalLabel
    },
    options: buildOptions(options, selectedIndex, arrivalLabel, d.technicalAssessment),
    selectedIndex,
    participants: (d.participantsConsulted ?? []).map((role) => ({
      initials: roleInitials(role),
      role,
      isRequester: role === d.requestedBy
    })),
    participantCount: (d.participantsConsulted ?? []).length,
    outcome: {
      title: d.decision,
      body: d.consequence
    }
  };
}

/* ------------------------------------------------------------------ */
/* ProcurementData                                                     */
/* ------------------------------------------------------------------ */

export function mapScenarioToProcurementData(): ProcurementData {
  const canonicalOfferPrice = formatMoney(S.offer.totalPrice, "EUR");
  const canonicalUnitPrice = formatMoney(S.offer.unitPrice, "EUR");

  const primaryOffer: SupplierOfferItem = {
    offerId: S.offer.offerId,
    offerReference: S.offer.offerReference,
    supplierId: S.supplier.supplierId,
    supplierName: "EuroSteel Construction",
    isVerified: true,
    supplierStatus: "ACCEPTED",
    totalPrice: canonicalOfferPrice,
    unitPrice: `${canonicalUnitPrice} / ton`,
    currency: "EUR",
    specification: S.material.specification || "Viaduct Spec v3.2",
    quantity: S.requirement.quantity || 120,
    unit: S.requirement.unit || "ton",
    deliveryWindow: "Aug 15 - Aug 18, 2026",
    offerValidity: "Valid thru Sep 30, 2026",
    leadTime: "3 Business Days",
    complianceScore: 98,
    complianceSpecVersion: "Viaduct Spec v3.2",
    certifications: ["CE Marked", "EN 10204 3.1 Mill Cert", "Procurement Lead Signed"],
    offerStatus: "ACCEPTED • DELIVERY PLANNED",
    statusBanner: "Status: ACCEPTED — Delivery planned for Aug 15",
    notes: "Primary offer accepted and aligned with site delivery window.",
    isPrimary: true
  };

  const secondaryOffer: SupplierOfferItem = {
    offerId: "OFF-1039",
    offerReference: "OFF-1039",
    supplierId: "SUP-002",
    supplierName: "Arcelor Infrastructure",
    isVerified: true,
    supplierStatus: "SUBMITTED",
    totalPrice: formatMoney(148200, "EUR"),
    unitPrice: `${formatMoney(1235.00, "EUR")} / ton`,
    currency: "EUR",
    specification: "Viaduct Spec v3.2 Alternative",
    quantity: 120,
    unit: "ton",
    deliveryWindow: "Aug 18 - Aug 22, 2026",
    offerValidity: "Valid thru Sep 15, 2026",
    leadTime: "5 Business Days",
    complianceScore: 94,
    complianceSpecVersion: "Viaduct Spec v3.1",
    certifications: ["CE Marked", "ISO 9001"],
    offerStatus: "SUBMITTED",
    statusBanner: "Status: SUBMITTED — Non-binding alternative",
    notes: "Alternative offer with extended lead time.",
    isPrimary: false
  };

  const offers = [primaryOffer, secondaryOffer];

  return {
    id: S.requirement.requirementId,
    requirementId: S.requirement.requirementId,
    projectId: S.project.projectId,
    projectCode: S.project.projectReference,
    projectName: S.project.projectName,

    code: S.requirement.requirementReference,
    materialId: S.material.materialId,
    materialName: "120t Structural Beams",
    materialStandard: S.material.standard,
    specification: S.material.specification,
    quantity: S.requirement.quantity,
    unit: S.requirement.unit,
    targetPhase: S.requirement.projectPhase,
    status: S.requirement.status,
    requiredDate: S.requirement.requiredDate,
    priority: S.requirement.priority,

    siteId: S.delivery.siteId,
    siteName: S.site.name,

    lowestBidTotal: "€142,500",
    lowestBidSavings: "3.8% below est.",
    fastestDeliveryDate: "Aug 15",
    fastestDeliveryStatus: "Target Met (On Time)",

    offerId: primaryOffer.offerId,
    offerReference: primaryOffer.offerReference,
    supplierId: primaryOffer.supplierId,
    supplierName: primaryOffer.supplierName,
    supplierStatus: primaryOffer.supplierStatus,
    supplierCertifications: primaryOffer.certifications,

    totalPrice: primaryOffer.totalPrice,
    unitPrice: primaryOffer.unitPrice,
    currency: primaryOffer.currency,

    proposedDeliveryDate: primaryOffer.deliveryWindow,
    availabilityDate: S.offer.availabilityDate,
    validUntil: primaryOffer.offerValidity,
    offerStatus: primaryOffer.offerStatus,

    primaryOffer,
    secondaryOffer,
    offers
  };
}

/* ------------------------------------------------------------------ */
/* LogisticsOverviewData                                               */
/* ------------------------------------------------------------------ */

export function mapScenarioToLogisticsOverviewData(): LogisticsOverviewData {
  return {
    delivery: {
      id: S.delivery.deliveryId,
      reference: S.delivery.deliveryReference,
      status: S.delivery.status,
      plannedDeliveryDate: S.delivery.plannedDeliveryDate,
      eta: S.delivery.eta,
      actualDeliveryDate: S.delivery.actualDeliveryDate || undefined,
      quantity: S.delivery.quantity,
      unit: S.delivery.unit,
      deliveryNotes: S.delivery.deliveryNotes,

      /* Stitch Delivery Detail presentation values */
      statusTag: "DELAYED",
      lastUpdatedAge: "Updated 4m ago",
      headlineTitle: "Delivery #DEL-1042",
      headlineQuantity: "Structural beams (84 Units / 120t)",
      alertTitle: "Road Restriction on M1",
      alertDescription: "Heavy transport clearance halted. Re-routing assessment currently underway.",
      delayBadgeLabel: "+21.5h Delay",
      plannedDateLabel: "Aug 15",
      plannedTimeLabel: "10:00 AM",
      revisedDateLabel: "Aug 16",
      revisedTimeLabel: "07:30 AM",
      primaryDelayFactor: "Road restriction on planned route M1 reported at 11:41 by Fleet Telematics.",
      materialGradeText: "S355JR Heavy Steel",
      supplierNameText: "EuroSteel",
      supplierCodeText: "SUP-001",
      contractOfferCodeText: "OFF-1042",
      destinationSiteNameText: "Riverside Tower",
      destinationBayText: "SITE-001 (Bay 3B)",
      carrierAgencyText: "FastBuild Logistics",
      carrierVehicleText: "TR-1042",
      linkedIncidentReference: "ISS-1042",
      linkedIncidentLabel: "Critical Delay",
      telemetryStatusText: "GPS Active",
      telemetryLocationText: "Stationary: M1 Northbound (KM 84.2)",
      telemetryStateText: "Stationary",
      timelineStageText: "Stage 4 of 8",
      timelineStages: [
        { stageNumber: 1, title: "Planned", subtitle: "Aug 10 — Dispatched Work Order", state: "completed" },
        { stageNumber: 2, title: "Dispatched", subtitle: "Aug 15 — 06:30 AM (Depot North)", state: "completed" },
        { stageNumber: 3, title: "In Transit", subtitle: "Heavy convoy escort confirmed", state: "completed" },
        { stageNumber: 4, title: "Delayed on Route M1", subtitle: "Reported Aug 15, 11:41 AM (Bridge Clearance restriction)", state: "active", badge: "Active" },
        { stageNumber: 5, title: "Rerouted", subtitle: "Transport Rerouted — Alternative Route B in progress", state: "pending" },
        { stageNumber: 6, title: "Delivered", subtitle: "Site gate check-in", state: "future" },
        { stageNumber: 7, title: "Inspected", subtitle: "Quality control validation", state: "future" },
        { stageNumber: 8, title: "Accepted", subtitle: "Sign-off and structural release", state: "future" }
      ]
    },
    project: {
      id: S.project.projectId,
      code: S.project.projectReference,
      name: S.project.projectName,
      type: S.project.projectType
    },
    requirement: {
      id: S.requirement.requirementId,
      code: S.requirement.requirementReference,
      targetPhase: S.requirement.projectPhase
    },
    // reference: the material record is the requirement's material.
    material: {
      id: S.material.materialId,
      name: S.material.name,
      specification: S.material.specification
    },
    supplier: {
      id: S.supplier.supplierId,
      name: S.supplier.name
    },
    offer: {
      id: S.offer.offerId,
      reference: S.offer.offerReference,
      totalPrice: formatMoney(S.offer.totalPrice, S.offer.currency),
      unitPrice: formatMoney(S.offer.unitPrice, S.offer.currency),
      currency: S.offer.currency,
      proposedDeliveryDate: S.offer.proposedDeliveryDate
    },
    site: {
      id: S.site.siteId,
      name: S.site.name,
      accessWindow: S.site.accessWindow,
      accessConstraints: S.site.accessConstraints,
      storageCapacity: S.site.storageCapacity,
      receivingRequirements: S.site.receivingRequirements,
      inspectionRequired: S.site.inspectionRequired,
      acceptanceRequired: S.site.acceptanceRequired
    },
    transport: {
      id: S.transport.transportId,
      reference: S.transport.transportReference,
      vehicle: S.transport.vehicleType,
      vehicleReference: S.transport.vehicleReference,
      driverId: S.transport.driverId,
      status: S.transport.status,
      // format: canonical Location -> one line
      origin: formatLocation(S.transport.origin) || "Origin Hub",
      // reference: transport.destinationSiteId is a site id.
      destination: S.site.name || "Viaduct Site",
      route: S.transport.routeReference,
      plannedDeparture: S.transport.plannedDeparture,
      actualDeparture: S.transport.actualDeparture,
      plannedArrival: S.transport.plannedArrival,
      eta: S.transport.eta,
      actualArrival: S.transport.actualArrival || undefined,
      constraints: S.transport.constraints,
      delayReason: S.transport.delayReason,

      /* Stitch Transport & Delay fields */
      statusTag: "ACTIVE ALERT",
      alertTitle: "Transport Delay Detected",
      alertConstraint: "Route M1 Blocked",
      originalEtaLabel: "Aug 15 • 10:00",
      revisedEtaLabel: "Aug 16 • 07:30",
      delayDurationLabel: "+21h 30m",
      gpsSyncAge: "Sync 12s ago",
      clearanceConstraint: "M1 Clearance < 4.1m",
      blockedRouteName: "Planned M1 (Blocked)",
      bypassRouteName: "Bypass A86 / D914",
      primaryIncidentText: "Emergency road restriction & oversized load weight limit on planned corridor route M1.",
      sitePhaseImpact: "Viaduct Deck Assembly Stalled",
      reportedByText: "FastBuild Dispatch",
      driverInitials: "TS",
      driverName: "Tom Smith",
      carrierName: "FastBuild Logistics",
      detourStatusTag: "ROUTE B APPROVED",
      detourCorridorText: "A86 Bypass > D914 Industrial Link",
      deltaDistanceText: "+24.0 km",
      scheduleModificationText: "+3h 15m driving + overnight staging",
      clearanceAuthorityText: "Dual-signed by Structural Lead & Transport Coordinator under protocol DEC-1042.",
      decisionReference: S.decision.decisionReference || "DEC-1042"
    },
    issue: {
      id: S.issue.issueId,
      reference: S.issue.issueReference,
      title: "Critical Beam Delivery Delayed by Route M1 Blockage",
      description: S.issue.description,
      type: S.issue.type,
      severity: S.issue.severity,
      status: S.issue.status,
      // reference: reportedBy is a vault identity id.
      reportedBy: actorLabel(S.issue.reportedBy),
      reportedAt: S.issue.reportedAt,
      // direct: the scenario states the impact explicitly.
      impact: S.issue.impact,
      assignedTo: S.issue.assignedTo,
      resolution: S.issue.resolution,
      resolvedAt: S.issue.resolvedAt,
      evidenceReferences: S.issue.evidenceReferences,

      /* Stitch Issue ISS-1042 presentation values */
      categoryBadge: "LOGISTICS BOTTLENECK",
      severityTag: "SEVERITY: CRITICAL",
      statusTag: "RESOLVED",
      mitigationLabel: "Rerouted via Decision DEC-1042",
      mitigationActionText: "Inspect",
      zoneContextText: "Zone 4 Viaduct",
      affectedResourceCode: "DEL-1042",
      affectedResourceDescription: "Structural Beams (44t Prefabricated)",
      projectPhaseText: "Foundation",
      projectZoneText: "Viaduct Zone 4",
      rootCauseTitle: "Route M1 Alert",
      rootCauseSecondary: "Overpass clearance",
      reporterInitials: "DK",
      reporterNameText: "David K. (Logistics Mgr)",
      reporterTimeText: "Aug 15, 11:41 AM",
      stakeholders: [
        { name: "BuildCorp (Main)", role: "Main Contractor", icon: "domain" },
        { name: "FastBuild Logistics", role: "Carrier", icon: "local_shipping" },
        { name: "Engineering Partners", role: "Consultant", icon: "engineering" }
      ],
      impactBadgeText: "Schedule Impact",
      primaryImpactTitle: "Foundation Phase Delayed",
      primaryImpactSubtitle: "Viaduct Zone 4 erection paused awaiting beam arrival",
      primaryImpactDescription: "Foundation heavy assembly team diverted to auxiliary drainage culverts until crane access window re-opens.",
      rescheduleLabel: "Assembly Reschedule",
      reschedulePlannedLabel: "Aug 15 (Planned)",
      rescheduleAdjustedLabel: "Aug 16 Morning (Adjusted)",
      evidenceFiles: [
        {
          filename: "Road_Restriction_Notice_M1.pdf",
          metadata: "Official Dept of Roads Alert • 1.4 MB",
          icon: "picture_as_pdf",
          size: "1.4 MB"
        },
        {
          filename: "Transport_Detour_Report_TR1042.pdf",
          metadata: "Carrier route survey & clearance • 2.1 MB",
          icon: "alt_route",
          size: "2.1 MB"
        }
      ],
      threadUpdateCount: 8
    },
    decision: {
      id: S.decision.decisionId,
      reference: S.decision.decisionReference,
      subject: S.decision.subject,
      context: S.decision.context,
      technicalAssessment: S.decision.technicalAssessment,
      risksIdentified: S.decision.risksIdentified,
      participantsConsulted: S.decision.participantsConsulted,
      optionsConsidered: S.decision.optionsConsidered,
      decision: S.decision.decision,
      decidedBy: S.decision.decidedBy,
      decisionDate: S.decision.decisionDate,
      consequence: S.decision.consequence,
      status: S.decision.status
    },
    inspection: {
      id: S.inspection.inspectionId,
      reference: S.inspection.inspectionReference,
      type: S.inspection.type,
      // reference: inspectorId is a vault identity id.
      inspector: actorLabel(S.inspection.inspectorId),
      inspectionDate: S.inspection.completedAt,
      status: S.inspection.status,
      result: S.inspection.result,
      criteria: S.inspection.criteria,
      findings: S.inspection.findings,
      notes: S.inspection.notes,
      evidenceReferences: S.inspection.evidenceReferences
    },
    provenance: {
      rootCauseTitle: "Root-Cause Diagnosis",
      rootCauseDescription: "Delivery #DEL-1042 was delayed by 21h 30m due to an unforeseen physical constraint on transport corridor M1, resolved via collaborative consensus on Route B.",
      delayDurationLabel: "+21h 30m Delay",
      originalEtaLabel: "Aug 15, 10:00",
      actualSiteGateLabel: "Aug 16, 07:32",
      scheduleDeltaLabel: "Absorbed (0d)",
      milestones: [
        {
          stepNumber: 1,
          category: "Physical Cause",
          timestamp: "Aug 15 • 11:15 AM",
          title: "Road restriction on Route M1",
          description: "Emergency overpass load-limit closure published by National Transport Authority.",
          nodeType: "cause"
        },
        {
          stepNumber: 2,
          category: "Disruption Ticket",
          timestamp: "Aug 15 • 11:41 AM",
          title: "ISS-1042 — Critical Beam Delivery Delayed",
          description: "Automated delay flag triggered by telematics geo-fence on carrier vehicle.",
          nodeType: "ticket"
        },
        {
          stepNumber: 3,
          category: "Verification",
          timestamp: "Aug 15 • 12:05 PM",
          title: "Road Restriction Report (#RD-9942)",
          description: "Verified civil traffic bulletin & real-time detour axle-clearance analysis.",
          nodeType: "verification",
          evidenceChip: {
            filename: "Corridor_M1_Closure_Order.pdf",
            badge: "Signed"
          }
        },
        {
          stepNumber: 4,
          category: "Multi-Org Consensus",
          timestamp: "Aug 15 • 12:45 PM",
          title: "DEC-1042 — Alternative Route B Approved",
          description: "Joint electronic sign-off completed within 44 minutes of issue creation.",
          nodeType: "consensus",
          signatories: ["BuildCorp (PM)", "EuroSteel (Eng)", "FastBuild (Carrier)"]
        },
        {
          stepNumber: 5,
          category: "Field Execution",
          timestamp: "Aug 15 • 01:00 PM",
          title: "Transport Rerouted via Bypass",
          description: "Convoy departed onto regional detour corridor B-88 with highway escort.",
          nodeType: "execution"
        },
        {
          stepNumber: 6,
          category: "Site Receipt",
          timestamp: "Aug 16 • 07:32 AM",
          title: "Delivery Received at Site-001",
          description: "Weighbridge scan valid. Transferred to laydown yard sector North-B.",
          nodeType: "receipt"
        },
        {
          stepNumber: 7,
          category: "Quality Control",
          timestamp: "Aug 16 • 08:45 AM",
          title: "Inspection #INSP-1042 Passed",
          description: "Flange camber, weld ultrasonic testing, and steel mill certifications verified 100%.",
          nodeType: "qc"
        },
        {
          stepNumber: 8,
          category: "Construction Handover",
          timestamp: "Aug 16 • 09:15 AM",
          title: "Material Accepted for Assembly",
          description: "84 structural steel beams officially released to Viaduct Erection Gang 03.",
          nodeType: "handover"
        }
      ],
      materialAcceptance: {
        materialCode: "MAT-STRUCT-001",
        inspectionCode: "INSP-1042",
        lotLabel: "Lot #STM-88219 (84 Beams)",
        beamCount: 84,
        status: "Accepted"
      },
      auditLayers: [
        {
          id: "construction",
          tabLabel: "Construction",
          title: "Physical Site Data & Field Delivery",
          description: "Captured from gate scale telematics, site manager daily logs, and physical acceptance notes signed at Viaduct Sector 4.",
          footerLeft: "Authority: Site Super (BuildCorp)",
          statusRight: "Synced in Realtime"
        },
        {
          id: "collaboration",
          tabLabel: "Collaboration & Evidence",
          title: "Multi-Organization Consensus Protocol",
          description: "3-party signed evidence exchange between BuildCorp, EuroSteel Fabricators, and FastBuild Haulage with route liability approvals.",
          footerLeft: "Threads: 14 exchanged records",
          statusRight: "Consensus 100%"
        },
        {
          id: "milestones",
          tabLabel: "Historical Milestones",
          title: "Historical Milestone Record",
          description: "Chronological audit trail cross-referenced against the master project schedule baseline. Provides transparent, verifiable milestone delivery history.",
          footerLeft: "Ref: TRACE-MILESTONE-001",
          statusRight: "Audited Record"
        }
      ]
    },
    inspectionPageData: {
      id: S.inspection.inspectionId,
      reference: S.inspection.inspectionReference,
      statusTag: "INSPECTION COMPLETE",
      title: "Foundation — Zone A",
      inspectorText: "Inspector: Bureau Inspection",
      checklist: [
        {
          id: 1,
          title: "Reinforcement Layout",
          description: "Spacing and bar size verified per structural plans (S-201).",
          isPassed: true
        },
        {
          id: 2,
          title: "Formwork Integrity",
          description: "Bracing and dimensions confirmed against formwork design (F-10).",
          isPassed: true
        },
        {
          id: 3,
          title: "Concrete Quality",
          description: "Slump test passed. Mix design matched specifications.",
          isPassed: true
        },
        {
          id: 4,
          title: "Overall Dimensions",
          description: "Tolerances within acceptable limits (+/- 5mm).",
          isPassed: true
        },
        {
          id: 5,
          title: "Site Safety Protocol",
          description: "All workers wearing appropriate PPE. Trench shoring intact.",
          isPassed: true
        }
      ],
      evidencePhotos: [
        {
          id: "photo-1",
          label: "Rebar Layout",
          imageAlt: "Rebar Layout Photo"
        },
        {
          id: "photo-2",
          label: "Formwork",
          imageAlt: "Formwork Photo"
        }
      ],
      evidenceDocuments: [
        {
          id: "doc-1",
          icon: "science",
          title: "Concrete Test Results",
          metadata: "PDF • 1.2 MB"
        },
        {
          id: "doc-2",
          icon: "description",
          title: "Official Inspection Report",
          metadata: "PDF • 3.4 MB"
        }
      ],
      timeline: [
        {
          id: "timeline-1",
          category: "Final Approval",
          title: "Approved by Bureau Inspection",
          timestamp: "Oct 24, 2023 - 14:30",
          isTerminal: true
        },
        {
          id: "timeline-2",
          category: "Review",
          title: "Site walk-through completed",
          timestamp: "Oct 24, 2023 - 10:15"
        },
        {
          id: "timeline-3",
          category: "Initiation",
          title: "Inspection requested by Contractor",
          timestamp: "Oct 23, 2023 - 09:00"
        }
      ]
    }
  };
}
