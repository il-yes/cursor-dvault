/**
 * Page-facing project list, projected from the canonical XS-BIM scenario.
 *
 * The scenario seeds exactly ONE construction project (PRJ-001). Fields with no
 * canonical source are optional and left undefined rather than invented — the
 * listing view renders an explicit unavailable state for them.
 *
 * Derivation rules (all from canonical records, none fabricated):
 *   openIssuesCount       issues whose canonical status is not terminal
 *   pendingDecisionsCount decisions whose canonical status is not approved
 *   connectedOrgs         the scenario participants' organization ids
 *   recentActivity        the last canonical thread event
 */

import { CANONICAL_SCENARIO } from "./constructionScenarioAdapter";
import { formatLocation } from "./scenarioMappers";

export interface ProjectData {
  id: string;
  code: string;
  name: string;
  type: string;
  status: string;
  currentPhase: string;
  progressPercent: number;
  targetCompletion: string;
  description: string;
  /** Empty in the canonical scenario (Location is a zero value). */
  location: string;
  openIssuesCount: number;
  pendingDecisionsCount: number;
  /** Present only while the canonical transport status is a delayed state. */
  activeDelay?: string;
  recentActivity: string;
  connectedOrgs: string[];
  isAuthoritative: boolean;
  image: string;
  /** @gap no canonical budget is populated. */
  budgetSpentPercent?: number;
  /** @gap no canonical schedule baseline is modelled. */
  scheduleDay?: number;
  /** @gap no canonical contract id exists. */
  contractId?: string;
  /** @gap the canonical ConstructionProject has no sector concept. */
  sector?: string;
}

const P = CANONICAL_SCENARIO.project;

const TERMINAL_ISSUE_STATUSES = ["resolved", "closed"];
const DELAYED_TRANSPORT_STATUSES = ["delayed"];

const lastEvent = CANONICAL_SCENARIO.threadEvents[CANONICAL_SCENARIO.threadEvents.length - 1];

export const MOCK_PROJECTS: ProjectData[] = [
  {
    id: P.projectId,
    code: P.projectReference,
    name: P.projectName,
    type: P.projectType,
    status: P.status,
    currentPhase: P.currentPhase,
    progressPercent: P.progressPercentage,
    targetCompletion: P.plannedEndDate,
    description: P.statusSummary,
    location: formatLocation(P.location),
    openIssuesCount: TERMINAL_ISSUE_STATUSES.includes(CANONICAL_SCENARIO.issue.status) ? 0 : 1,
    pendingDecisionsCount: CANONICAL_SCENARIO.decision.status === "approved" ? 0 : 1,
    activeDelay: DELAYED_TRANSPORT_STATUSES.includes(CANONICAL_SCENARIO.transport.status)
      ? CANONICAL_SCENARIO.transport.delayReason
      : undefined,
    recentActivity: lastEvent.eventType,
    connectedOrgs: CANONICAL_SCENARIO.participants.map((p) => p.organizationId),
    isAuthoritative: true,
    image:
      "https://lh3.googleusercontent.com/aida-public/AB6AXuCyhdIyunCsHJbMDm1UWI9VXIvPj7GBEI9F2ZoXR61KUPGZkgrMZflrpqE21OhZxTYJcbn_ADHB5FvkGObh09NsqZmUiZtY9aDyNoP4FUXeR4JNKmUosU_MNAcafsPG_jS6TunbkQiEX4l8nddQ1vxS4FpY21vauK7SCpG5QVlLgJNoef0nj2oirgVq0lBL-ZiR1C4ouQH7N-QrOJpG6YENchuyphAPTX21Wib9fOGuSX1OEhMU3cs"
  }
];
