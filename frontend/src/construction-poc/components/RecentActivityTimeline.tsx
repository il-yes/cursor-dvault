import React from "react";

export interface TimelineItem {
  id: string;
  title: string;
  timestamp: string;
  isPrimary?: boolean;
}

interface RecentActivityTimelineProps {
  items?: TimelineItem[];
}

const defaultTimelineItems: TimelineItem[] = [
  {
    id: "act-1",
    title: "Structural inspection approved",
    timestamp: "Today, 09:42 AM",
    isPrimary: true,
  },
  {
    id: "act-2",
    title: "Concrete delivery received",
    timestamp: "Today, 09:18 AM",
    isPrimary: false,
  },
  {
    id: "act-3",
    title: "Revised drawing uploaded",
    timestamp: "Today, 08:54 AM",
    isPrimary: false,
  },
];

export const RecentActivityTimeline: React.FC<RecentActivityTimelineProps> = ({
  items = defaultTimelineItems,
}) => {
  return (
    <div className="w-full bg-white rounded-2xl p-5 shadow-[0_4px_12px_rgba(4,22,39,0.05)] border border-[#e0e3e5]/50 mb-4">
      <h2 className="text-lg font-semibold text-[#181c1e] mb-4 flex items-center gap-2">
        <span className="material-symbols-outlined text-[#041627]">
          history
        </span>
        Recent Activity
      </h2>
      <div className="relative pl-3">
        {/* Timeline Connecting Line */}
        <div className="absolute left-[15px] top-2 bottom-2 w-[2px] bg-[#e0e3e5]" />

        <div className="flex flex-col gap-6">
          {items.map((item) => (
            <div key={item.id} className="relative flex gap-4 items-start">
              <div
                className={`w-[6px] h-[6px] rounded-full ${
                  item.isPrimary ? "bg-[#041627]" : "bg-[#e0e3e5]"
                } mt-2 z-10 ring-4 ring-white`}
              />
              <div className="flex flex-col">
                <span
                  className={`text-sm ${
                    item.isPrimary
                      ? "text-[#181c1e] font-medium"
                      : "text-[#181c1e]"
                  }`}
                >
                  {item.title}
                </span>
                <span className="text-[11px] text-[#44474c] mt-0.5 font-medium">
                  {item.timestamp}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
