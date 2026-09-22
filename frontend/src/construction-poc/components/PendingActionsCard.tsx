import React from "react";

export interface PendingActionItem {
  id: string;
  title: string;
  badgeText: string;
  badgeClass: string; // e.g. "bg-[#3e2400]/10 text-[#ca8100]", "bg-[#e0e3e5] text-[#44474c]"
  iconName: string;
  iconBgClass: string; // e.g. "bg-[#3e2400]/10 text-[#3e2400]", "bg-[#041627]/10 text-[#041627]"
}

interface PendingActionsCardProps {
  actions?: PendingActionItem[];
  onViewAll?: () => void;
  onActionClick?: (action: PendingActionItem) => void;
}

const defaultActions: PendingActionItem[] = [
  {
    id: "approvals",
    title: "Approvals",
    badgeText: "3 Requires Attention",
    badgeClass: "bg-[#3e2400]/10 text-[#ca8100]",
    iconName: "fact_check",
    iconBgClass: "bg-[#3e2400]/10 text-[#ca8100]",
  },
  {
    id: "inspections",
    title: "Inspections",
    badgeText: "2 Scheduled",
    badgeClass: "bg-[#e0e3e5] text-[#44474c]",
    iconName: "search",
    iconBgClass: "bg-[#041627]/10 text-[#041627]",
  },
  {
    id: "reviews",
    title: "Reviews",
    badgeText: "4 Pending",
    badgeClass: "bg-[#e0e3e5] text-[#44474c]",
    iconName: "rate_review",
    iconBgClass: "bg-[#041627]/10 text-[#041627]",
  },
  {
    id: "confirmations",
    title: "Confirmations",
    badgeText: "1 Required",
    badgeClass: "bg-[#e0e3e5] text-[#44474c]",
    iconName: "check_box",
    iconBgClass: "bg-[#041627]/10 text-[#041627]",
  },
];

export const PendingActionsCard: React.FC<PendingActionsCardProps> = ({
  actions = defaultActions,
  onViewAll,
  onActionClick,
}) => {
  return (
    <div className="w-full bg-white rounded-2xl p-5 shadow-[0_4px_12px_rgba(4,22,39,0.05)] border border-[#e0e3e5]/50">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold text-[#181c1e] flex items-center gap-2">
          <span className="material-symbols-outlined text-[#041627]">
            assignment
          </span>
          Pending Actions
        </h2>
        <button
          onClick={onViewAll}
          className="text-[#041627] text-xs font-semibold hover:underline"
        >
          View All
        </button>
      </div>

      <div className="flex flex-col gap-3">
        {actions.map((item) => (
          <div
            key={item.id}
            onClick={() => onActionClick?.(item)}
            className="flex items-center justify-between p-3 rounded-xl bg-[#f1f4f6] hover:bg-[#ebeef0] transition-colors cursor-pointer group"
          >
            <div className="flex items-center gap-3">
              <div
                className={`w-8 h-8 rounded-full ${item.iconBgClass} flex items-center justify-center`}
              >
                <span className="material-symbols-outlined text-[18px]">
                  {item.iconName}
                </span>
              </div>
              <span className="text-sm text-[#181c1e] font-medium">
                {item.title}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span
                className={`${item.badgeClass} px-2 py-0.5 rounded-full text-[11px] font-medium`}
              >
                {item.badgeText}
              </span>
              <span className="material-symbols-outlined text-[#44474c] text-[18px] group-hover:translate-x-1 transition-transform">
                chevron_right
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
