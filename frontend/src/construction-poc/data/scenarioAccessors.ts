/**
 * Static scenario accessors shared by every ConstructionDataProvider.
 *
 * The XS-BIM presentation scenario (PRJ-001, REQ-STRUCT-001, DEL-1042, ...) has
 * no AppAPI equivalent yet, so both the mock and the Cloud provider serve these
 * values. Keeping them here means neither provider has to depend on the other.
 */

import { SCENARIO_DATA } from "./constructionScenarioAdapter";
import { MOCK_ACTIVITY_LOG } from "./scenario.mock";
import type { ActivityItem } from "./scenario.mock";

export const scenarioAccessors = {
  getActivityFeed(): ActivityItem[] {
    return MOCK_ACTIVITY_LOG;
  },

  getRequirement(): typeof SCENARIO_DATA.requirement {
    return SCENARIO_DATA.requirement;
  },

  getSupplierOffer(): typeof SCENARIO_DATA.offer {
    return SCENARIO_DATA.offer;
  },

  getDelivery(): typeof SCENARIO_DATA.delivery {
    return SCENARIO_DATA.delivery;
  },

  getTransportDelay(): typeof SCENARIO_DATA.transport {
    return SCENARIO_DATA.transport;
  },

  getIssue(): typeof SCENARIO_DATA.issue {
    return SCENARIO_DATA.issue;
  },

  getEvidenceDocument(): typeof SCENARIO_DATA.evidenceDoc {
    return SCENARIO_DATA.evidenceDoc;
  },

  getDecision(): typeof SCENARIO_DATA.decision {
    return SCENARIO_DATA.decision;
  },

  getInspection(): typeof SCENARIO_DATA.inspection {
    return SCENARIO_DATA.inspection;
  },

  getStakeholders(): typeof SCENARIO_DATA.stakeholders {
    return SCENARIO_DATA.stakeholders;
  },

  getTraceMilestones(): typeof SCENARIO_DATA.traceMilestones {
    return SCENARIO_DATA.traceMilestones;
  },
};
