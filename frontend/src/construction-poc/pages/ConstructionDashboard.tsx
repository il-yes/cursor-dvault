import React from "react";
import { ProjectHeroCard } from "../components/ProjectHeroCard";
import { ProjectHealthCard } from "../components/ProjectHealthCard";
import { PendingActionsCard, PendingActionItem } from "../components/PendingActionsCard";
import { RecentActivityTimeline } from "../components/RecentActivityTimeline";

interface ConstructionDashboardProps {
  onNavigateToProjects?: () => void;
  onNavigateToActivity?: () => void;
}

export const ConstructionDashboard: React.FC<ConstructionDashboardProps> = ({
  onNavigateToProjects,
  onNavigateToActivity,
}) => {
  const handleActionClick = (action: PendingActionItem) => {
    if (action.id === "inspections" || action.id === "approvals") {
      onNavigateToActivity?.();
    } else {
      onNavigateToProjects?.();
    }
  };

  return (
    <div className="flex flex-col w-full gap-4 px-4 py-4 max-w-2xl mx-auto">
      {/* Hero / Project Title Card */}
      <ProjectHeroCard />

      {/* Project Health Card */}
      <ProjectHealthCard />

      {/* Pending Actions Card */}
      <PendingActionsCard
        onViewAll={onNavigateToActivity}
        onActionClick={handleActionClick}
      />

      {/* Recent Activity Timeline */}
      <RecentActivityTimeline />
    </div>
  );
};

export default ConstructionDashboard;
