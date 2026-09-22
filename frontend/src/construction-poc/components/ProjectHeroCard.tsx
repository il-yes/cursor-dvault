import React from "react";

interface ProjectHeroCardProps {
  location?: string;
  projectName?: string;
  progressPercent?: number;
  statusBadgeText?: string;
  targetDateText?: string;
}

export const ProjectHeroCard: React.FC<ProjectHeroCardProps> = ({
  location = "Paris, France",
  projectName = "Riverside Tower",
  progressPercent = 68,
  statusBadgeText = "Construction in Progress",
  targetDateText = "Target: Oct 2026",
}) => {
  const dashArray = `${progressPercent}, 100`;

  return (
    <div className="w-full bg-[#ebeef0] rounded-2xl p-6 relative overflow-hidden shadow-sm">
      {/* Background blueprint pattern SVG overlay */}
      <div
        className="absolute inset-0 opacity-10 pointer-events-none"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23041627' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
        }}
      />
      <div className="relative z-10 flex flex-col gap-4">
        <div className="flex justify-between items-start">
          <div>
            <div className="flex items-center gap-2 text-[#041627] mb-1">
              <span className="material-symbols-outlined text-[16px]">
                location_on
              </span>
              <span className="text-xs uppercase font-semibold tracking-wider">
                {location}
              </span>
            </div>
            <h1 className="text-2xl font-bold text-[#181c1e] tracking-tight">
              {projectName}
            </h1>
          </div>

          {/* SVG Progress Ring */}
          <div className="relative w-16 h-16 flex items-center justify-center shrink-0">
            <svg
              className="w-full h-full transform -rotate-90"
              viewBox="0 0 36 36"
            >
              <path
                className="text-[#d7dadc]"
                d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                fill="none"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="text-[#006c49]"
                d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                fill="none"
                stroke="currentColor"
                strokeDasharray={dashArray}
                strokeWidth="4"
              />
            </svg>
            <div className="absolute flex flex-col items-center justify-center">
              <span className="text-xs font-semibold text-[#181c1e]">
                {progressPercent}%
              </span>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2 mt-2">
          <span className="inline-flex items-center gap-1 bg-[#006c49]/10 text-[#006c49] px-2.5 py-1 rounded-full text-xs font-medium">
            <span className="material-symbols-outlined text-[14px]">
              construction
            </span>
            {statusBadgeText}
          </span>
          <span className="text-sm text-[#44474c] ml-auto">
            {targetDateText}
          </span>
        </div>
      </div>
    </div>
  );
};
