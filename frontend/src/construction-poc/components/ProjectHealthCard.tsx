import React from "react";

export interface HealthMetric {
  label: string;
  value: string;
  iconName: string;
  colorClass: string; // e.g. "text-[#006c49]", "text-[#041627]", "text-[#ba1a1a]"
}

interface ProjectHealthCardProps {
  metrics?: HealthMetric[];
}

const defaultMetrics: HealthMetric[] = [
  {
    label: "SCHEDULE",
    value: "On Track",
    iconName: "check_circle",
    colorClass: "text-[#006c49]",
  },
  {
    label: "BUDGET",
    value: "Healthy",
    iconName: "account_balance_wallet",
    colorClass: "text-[#006c49]",
  },
  {
    label: "COMPLIANCE",
    value: "94%",
    iconName: "verified",
    colorClass: "text-[#041627]",
  },
  {
    label: "OPEN ISSUES",
    value: "7 Active",
    iconName: "warning",
    colorClass: "text-[#ba1a1a]",
  },
];

export const ProjectHealthCard: React.FC<ProjectHealthCardProps> = ({
  metrics = defaultMetrics,
}) => {
  return (
    <div className="w-full bg-white rounded-2xl p-5 shadow-[0_4px_12px_rgba(4,22,39,0.05)] border border-[#e0e3e5]/50">
      <h2 className="text-lg font-semibold text-[#181c1e] mb-4 flex items-center gap-2">
        <span className="material-symbols-outlined text-[#041627]">
          health_and_safety
        </span>
        Project Health
      </h2>
      <div className="grid grid-cols-2 gap-4">
        {metrics.map((item, idx) => (
          <div key={idx} className="flex flex-col gap-1">
            <span className="text-[11px] font-medium text-[#44474c] uppercase tracking-wider">
              {item.label}
            </span>
            <div className={`flex items-center gap-1.5 ${item.colorClass}`}>
              <span className="material-symbols-outlined text-[18px]">
                {item.iconName}
              </span>
              <span className="text-sm font-medium">{item.value}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
