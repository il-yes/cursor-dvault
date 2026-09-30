/**
 * Activity feed, projected from the canonical XS-BIM thread event sequence.
 *
 * The 23 events in the authoritative scenario are the scenario's real trace
 * history (they are mirrored 1:1 into TraceCore as commits), so the feed is
 * generated from them rather than hand-authored rows.
 *
 * @gap The scenario records a monotonic `cursor` but no wall-clock time per
 * event, so `time` is left empty. The listing view renders the cursor ordinal
 * instead of fabricating timestamps.
 */

import { CANONICAL_SCENARIO } from "./constructionScenarioAdapter";
import { actorLabel } from "./scenarioMappers";

export interface ActivityItem {
  /** The canonical event idempotency key. */
  id: string;
  cursor: number;
  /** Always empty — the scenario has no per-event timestamp. */
  time: string;
  type: string;
  /** The canonical event type, e.g. `construction.issue.reported`. */
  badge: string;
  title: string;
  /** Derived from the actor's canonical vault address, never invented. */
  actor: string;
  link: string;
}

function linkFor(eventType: string): string {
  if (eventType.includes("issue")) return "/dashboard/construction/issues/ISS-1042";
  if (eventType.includes("decision")) return "/dashboard/construction/decisions/DEC-1042";
  if (eventType.includes("inspection")) return "/dashboard/construction/inspections/INSP-1042";
  if (eventType.includes("delivery") || eventType.includes("transport")) {
    return "/dashboard/construction/deliveries/DEL-1042";
  }
  if (eventType.includes("supplier") || eventType.includes("offer")) {
    return "/dashboard/construction/requirements/REQ-STRUCT-001/suppliers";
  }
  return `/dashboard/construction/projects/${CANONICAL_SCENARIO.project.projectId}`;
}

export const MOCK_ACTIVITY_LOG: ActivityItem[] = CANONICAL_SCENARIO.threadEvents.map((e) => ({
  id: e.idempotencyKey,
  cursor: e.cursor,
  time: "",
  type: e.eventType,
  badge: e.eventType,
  title: e.eventType,
  actor: actorLabel(e.actorId),
  link: linkFor(e.eventType)
}));
