/**
 * Scenario -> view-model mappers for the mock provider (Phase 2).
 *
 * These translate the authoritative XS-BIM presentation scenario into the
 * ProcurementData / LogisticsOverviewData shapes the Cloud DTO mappers produce,
 * so the presentation layer renders identically on either data source.
 *
 * Provenance rules applied here:
 *   - direct    : copied verbatim from a scenario field
 *   - derived   : computed from scenario fields by a rule stated at the call site
 *   - reference : scenario lacks the field; the view model uses a related value
 *
 * No value is invented. Where the scenario genuinely carries no datum, the field
 * is omitted (optional) rather than filled with a placeholder.
 */

import { SCENARIO_DATA } from "./constructionScenarioAdapter";
import type { LogisticsOverviewData, ProcurementData } from "./constructionTypes";

const {
  project,
  requirement,
  supplier,
  offer,
  delivery,
  transport,
  issue,
  evidenceDoc,
  decision,
  inspection,
  stakeholders,
  traceMilestones,
} = SCENARIO_DATA;

function matches(candidate: string, requested: string): boolean {
  return candidate.toLowerCase() === requested.toLowerCase();
}

/** The scenario carries a single actor attribution per entity, via the trace log. */
function actorFor(entityId: string): string {
  return traceMilestones.find((m) => matches(m.entityId, entityId))?.actor || "";
}

/** A supplier is "invited" when at least one scenario stakeholder represents its org. */
function supplierContacts(): typeof stakeholders {
  return stakeholders.filter((s) => matches(s.organization, offer.supplierName));
}

/** derived: fulfilment lateness implies procurement urgency. */
function priorityFor(status: string): string {
  if (status === "fulfilled_delayed") return "HIGH";
  if (status === "pending") return "MEDIUM";
  return "LOW";
}

export function scenarioHasRequirement(requirementId?: string): boolean {
  if (!requirementId) return false;
  return matches(requirement.id, requirementId) || matches(requirement.code, requirementId);
}

export function scenarioHasDelivery(deliveryId?: string): boolean {
  if (!deliveryId) return false;
  return matches(delivery.id, deliveryId) || matches(delivery.reference, deliveryId);
}

export function mapScenarioToProcurementData(): ProcurementData {
  const contacts = supplierContacts();

  return {
    id: requirement.id,
    requirementId: requirement.id,
    projectId: requirement.projectId,
    projectCode: project.code,
    projectName: project.name,
    code: requirement.code,
    materialId: requirement.materialId,
    materialName: requirement.materialName,
    specification: requirement.specification,
    quantity: requirement.quantity,
    unit: requirement.unit,
    targetPhase: requirement.targetPhase,
    status: requirement.status,
    requiredDate: requirement.requiredDate,
    priority: priorityFor(requirement.status),
    // reference: the requirement itself carries no site; its delivery does.
    siteId: delivery.siteId,
    siteName: delivery.siteName,
    // derived: the scenario records exactly one offer against this requirement.
    invitedSuppliersCount: contacts.length,
    offersReceivedCount: 1,
    offerId: offer.id,
    offerReference: offer.offerReference,
    supplierId: offer.supplierId,
    supplierName: offer.supplierName,
    totalPrice: offer.totalPrice,
    unitPrice: offer.unitPrice,
    promisedDeliveryDate: offer.promisedDeliveryDate,
    offerStatus: offer.status,
    // derived: verification is recorded per stakeholder, not per supplier.
    isVerifiedSupplier: contacts.some((c) => c.trustStatus === "VERIFIED"),
  };
}

export function mapScenarioToLogisticsOverviewData(): LogisticsOverviewData {
  return {
    delivery: {
      id: delivery.id,
      reference: delivery.reference,
      status: delivery.status,
      plannedDeliveryDate: delivery.plannedDeliveryDate,
      eta: delivery.eta,
      actualDeliveryDate: delivery.actualDeliveryDate,
      quantity: delivery.quantity,
      unit: delivery.unit,
      deliveryNotes: delivery.deliveryNotes,
    },
    project: {
      id: project.id,
      code: project.code,
      name: project.name,
      type: project.type,
    },
    requirement: {
      id: requirement.id,
      code: requirement.code,
      targetPhase: requirement.targetPhase,
    },
    // reference: the material lives on the requirement in the scenario.
    material: {
      id: requirement.materialId,
      name: requirement.materialName,
      specification: requirement.specification,
    },
    supplier: {
      id: supplier.id,
      name: supplier.name,
    },
    offer: {
      id: offer.id,
      reference: offer.offerReference,
      totalPrice: offer.totalPrice,
      unitPrice: offer.unitPrice,
      promisedDeliveryDate: offer.promisedDeliveryDate,
    },
    site: {
      id: delivery.siteId,
      name: delivery.siteName,
    },
    transport: {
      id: transport.id,
      // reference: the scenario transport has no separate reference; id is canonical.
      reference: transport.id,
      vehicle: transport.vehicle,
      driver: transport.driver,
      status: transport.status,
      origin: transport.origin,
      destination: transport.destination,
      route: transport.route,
      plannedDeparture: transport.plannedDeparture,
      actualDeparture: transport.actualDeparture,
      // the scenario records one arrival estimate, which serves both roles.
      plannedArrival: transport.eta,
      eta: transport.eta,
      // reference: the transport leg completed when the delivery was received.
      actualArrival: delivery.actualDeliveryDate,
      // the view model takes a list; the scenario states a single constraint.
      constraints: [transport.constraints],
      delayReason: transport.delayReason,
    },
    issue: {
      id: issue.id,
      reference: issue.reference,
      title: issue.title,
      description: issue.description,
      severity: issue.severity,
      status: issue.status,
      // derived: the scenario attributes the issue in the trace log.
      reportedBy: actorFor(issue.id),
      reportedAt: issue.reportedAt,
      // reference: the recorded impact is the quantified consequence of this issue.
      impact: decision.impactSummary,
      evidenceReferences: [issue.evidenceDocId],
    },
    decision: {
      id: decision.id,
      reference: decision.reference,
      subject: decision.title,
      // reference: the decision is recorded against this issue, whose description is its context.
      context: issue.description,
      decision: decision.proposedAction,
      decidedBy: decision.approvedBy || decision.proposedBy,
      decisionDate: decision.approvedAt || "",
      consequence: decision.impactSummary,
      status: decision.status,
    },
    inspection: {
      id: inspection.id,
      reference: inspection.reference,
      type: inspection.title,
      inspector: inspection.inspector,
      inspectionDate: inspection.inspectionDate,
      status: inspection.status,
      // the inspection records its dimensional verdict explicitly.
      result: inspection.dimensionalCompliance,
      criteria: [inspection.ultrasonicWeldIntegrity, inspection.dimensionalCompliance],
      notes: inspection.notes,
      evidenceReferences: [evidenceDoc.id],
    },
  };
}
