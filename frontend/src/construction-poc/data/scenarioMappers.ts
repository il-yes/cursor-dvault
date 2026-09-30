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
import type { LogisticsOverviewData, ProcurementData } from "./constructionTypes";

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
