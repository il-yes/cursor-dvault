/**
 * Static scenario accessors shared by every ConstructionDataProvider.
 *
 * The authoritative XS-BIM scenario (PRJ-001, REQ-STRUCT-001, DEL-1042, ...) has
 * no AppAPI equivalent yet, so both the mock and the Cloud provider serve these
 * values. Keeping them here means neither provider has to depend on the other.
 *
 * Accessors expose CANONICAL records verbatim. They perform no presentation
 * transformation — that belongs in scenarioMappers.ts.
 */

import { CANONICAL_SCENARIO } from "./constructionScenarioAdapter";
import { MOCK_ACTIVITY_LOG } from "./scenario.mock";
import type { ActivityItem } from "./scenario.mock";

const S = CANONICAL_SCENARIO;

export const scenarioAccessors = {
  getActivityFeed(): ActivityItem[] {
    return MOCK_ACTIVITY_LOG;
  },

  /** Named `*Record` to avoid colliding with the provider's ProjectData `getProject`. */
  getProjectRecord() {
    return S.project;
  },

  getMaterial() {
    return S.material;
  },

  getRequirement() {
    return S.requirement;
  },

  getSupplier() {
    return S.supplier;
  },

  getSupplierOffer() {
    return S.offer;
  },

  getSite() {
    return S.site;
  },

  getDelivery() {
    return S.delivery;
  },

  getTransportDelay() {
    return S.transport;
  },

  getIssue() {
    return S.issue;
  },

  getAction() {
    return S.action;
  },

  getEvidenceDocument() {
    return S.evidenceDoc;
  },

  getCertificateDocument() {
    return S.certDoc;
  },

  getDecision() {
    return S.decision;
  },

  getInspection() {
    return S.inspection;
  },

  /** @gap no ConstructionStakeholder aggregate exists — these are the three real vault actors. */
  getParticipants() {
    return S.participants;
  },

  getCollaboration() {
    return S.collaboration;
  },

  getShareEntry() {
    return S.shareEntry;
  },

  /** The authoritative 23-event sequence, in emission order. */
  getThreadEvents() {
    return S.threadEvents;
  },
};

/** Resolves a vault identity id to its canonical participant record. */
export function participantByVaultId(vaultId: string) {
  return S.participants.find((p) => p.vaultId === vaultId);
}

/** Resolves a document id to its canonical ConstructionDocument record. */
export function documentById(documentId: string) {
  return [S.evidenceDoc, S.certDoc].find((d) => d.documentId === documentId);
}
