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
  return {
    id: S.requirement.requirementId,
    requirementId: S.requirement.requirementId,
    projectId: S.project.projectId,
    projectCode: S.project.projectReference,
    projectName: S.project.projectName,

    // direct: canonical references
    code: S.requirement.requirementReference,
    materialId: S.material.materialId,
    materialName: S.material.name,
    materialStandard: S.material.standard,
    specification: S.material.specification,
    quantity: S.requirement.quantity,
    unit: S.requirement.unit,
    targetPhase: S.requirement.projectPhase,
    status: S.requirement.status,
    requiredDate: S.requirement.requiredDate,
    priority: S.requirement.priority,

    // reference: requirement.siteId is assigned in memory but never persisted;
    // the delivery carries the persisted site link.
    siteId: S.delivery.siteId,
    siteName: S.site.name,

    // direct: the scenario seeds exactly one offer for this requirement.
    offersReceivedCount: 1,
    offerId: S.offer.offerId,
    offerReference: S.offer.offerReference,
    supplierId: S.supplier.supplierId,
    supplierName: S.supplier.name,
    supplierStatus: S.supplier.status,
    supplierCertifications: S.supplier.certificationReferences,

    // format: canonical floats + currency -> display strings
    totalPrice: formatMoney(S.offer.totalPrice, S.offer.currency),
    unitPrice: formatMoney(S.offer.unitPrice, S.offer.currency),
    currency: S.offer.currency,

    proposedDeliveryDate: S.offer.proposedDeliveryDate,
    availabilityDate: S.offer.availabilityDate,
    validUntil: S.offer.validUntil,
    offerStatus: S.offer.status
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
      deliveryNotes: S.delivery.deliveryNotes
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
      origin: formatLocation(S.transport.origin),
      // reference: transport.destinationSiteId is a site id.
      destination: S.site.name,
      route: S.transport.routeReference,
      plannedDeparture: S.transport.plannedDeparture,
      actualDeparture: S.transport.actualDeparture,
      plannedArrival: S.transport.plannedArrival,
      eta: S.transport.eta,
      actualArrival: S.transport.actualArrival || undefined,
      constraints: S.transport.constraints,
      delayReason: S.transport.delayReason
    },
    issue: {
      id: S.issue.issueId,
      reference: S.issue.issueReference,
      title: S.issue.title,
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
      evidenceReferences: S.issue.evidenceReferences
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
    }
  };
}
