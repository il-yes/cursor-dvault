import React from "react";
import {
	format,
	subDays,
	eachDayOfInterval,
	startOfWeek,
	startOfMonth,
	isSameMonth,
} from "date-fns";
import "./g-scrollbar.css";


const generateLastYearDays = () => {
	const end = new Date();
	const start = subDays(end, 365);
	return eachDayOfInterval({ start, end });
};

const getColor = (count) => {
	if (count === 0) return "bg-gray-200 dark:bg-gray-700";
	if (count < 3) return "bg-green-200";
	if (count < 6) return "bg-green-400";
	if (count < 10) return "bg-green-600";
	return "bg-green-800";
};

const ContributionGraph = ({ contributions }) => {
	const days = generateLastYearDays();

	// Build weeks (each column is one week)
	const weeks = [];
	let currentWeek = [];

	days.forEach((day, i) => {
		const iso = format(day, "yyyy-MM-dd");
		const count = contributions?.[iso] || 0;
		const color = getColor(count);
		currentWeek.push({ day, color, count });

		if (day.getDay() === 6 || i === days.length - 1) {
			weeks.push(currentWeek);
			currentWeek = [];
		}
	});

	// Month labels
	const monthLabels = [];
	weeks.forEach((week, i) => {
		const firstDay = week[0];
		const prevWeek = weeks[i - 1];
		if (
			(!prevWeek ||
				!isSameMonth(firstDay.day, prevWeek[0].day)) &&
			firstDay.day.getDate() <= 7
		) {
			monthLabels.push({ index: i, label: format(firstDay.day, "MMM") });
		}
	});

	const weekdayLabels = ["Mon", "ma", "Wed", "my", "Fri"];

	return Glassmorphism(weeks, weekdayLabels, monthLabels);
};

const Basic = (weeks, weekdayLabels, monthLabels) => (
	<div className="flex flex-col items-start gap-1 text-xs text-gray-500 dark:text-gray-400 gap-[2px] overflow-x-auto p-1 max-w-full box-border">
		{/* Month labels */}
		<div className="flex gap-0.5" style={{ marginLeft: "26px" }}>
			{weeks.map((_, i) => {
				const month = monthLabels.find((m) => m.index === i);
				return (
					<div key={i} className="w-3">
						{month ? <span className="text-[10px]">{month.label}</span> : ""}
					</div>
				);
			})}
		</div>

		<div className="flex">
			{/* Weekday labels */}
			<div className="flex flex-col justify-between mr-1 h-[60px] pl-[2px]">
				{weekdayLabels.map((d) => (
					<span key={d} className="text-[8px]">
						{d}
					</span>
				))}
			</div>

			{/* Graph */}
			<div className="flex gap-[2px] overflow-x-auto p-1">
				{weeks.map((week, wi) => (
					<div key={wi} className="flex flex-col gap-[2px]">
						{Array.from({ length: 7 }).map((_, di) => {
							const d = week.find((x) => x.day.getDay() === di);
							if (!d)
								return (
									<div
										key={di}
										className="w-3 h-3 rounded-sm bg-transparent"
									></div>
								);
							return (
								<div
									key={di}
									className={`w-3 h-3 rounded-sm ${d.color} hover:scale-125 transition-transform`}
									title={`${format(
										d.day,
										"MMM d, yyyy"
									)}: ${d.count} contributions`}
								></div>
							);
						})}
					</div>
				))}
			</div>
		</div>
	</div>
)
const Glassmorphism = (weeks, weekdayLabels, monthLabels) => (
  <div className="flex flex-col items-start gap-1 text-xs text-gray-500 dark:text-gray-400 gap-[2px] p-3 max-w-full box-border backdrop-blur-xl bg-white/20 dark:bg-zinc-900/20 rounded-2xl border border-white/30 dark:border-zinc-700/30 shadow-lg">
    <div className="overflow-x-auto w-full scrollbar-glassmorphism thin-scrollbar">
      {/* Month labels and graph scroll together */}
      <div className="min-w-max" style={{ cursor: "grap" }}> {/* ensures content can scroll */}
        {/* Month labels */}
        <div className="flex gap-0.5 mb-3 px-[26px]">
          {weeks.map((_, i) => {
            const month = monthLabels.find((m) => m.index === i);
            return (
              <div key={i} className="w-3">
                {month ? (
                  <span className="text-[10px] font-medium bg-white/40 dark:bg-zinc-800/40 px-1 py-0.5 rounded backdrop-blur-sm border border-white/50 dark:border-zinc-700/50">
                    {month.label}
                  </span>
                ) : ""}
              </div>
            );
          })}
        </div>

        <div className="flex">
          {/* Weekday labels (not scrolling horizontally) */}
          <div className="flex flex-col justify-between mr-1 h-[60px] pl-[2px]">
            {weekdayLabels.map((d) => (
              <span key={d} className="text-[8px] font-semibold tracking-tight">
                {d}
              </span>
            ))}
          </div>

          {/* Graph */}
          <div className="flex gap-[2px] p-1">
            {weeks.map((week, wi) => (
              <div key={wi} className="flex flex-col gap-[2px]">
                {Array.from({ length: 7 }).map((_, di) => {
                  const d = week.find((x) => x.day.getDay() === di);
                  if (!d)
                    return (
                      <div
                        key={di}
                        className="w-3 h-3 rounded-sm bg-gradient-to-br from-transparent/50 to-transparent/20 backdrop-blur-sm border border-white/30 dark:border-zinc-800/50 hover:border-primary/50 transition-all"
                      />
                    );
                  return (
                    <div
                      key={di}
                      className={`w-3 h-3 rounded-sm ${d.color} hover:scale-125 hover:shadow-lg hover:shadow-primary/50 hover:backdrop-blur-sm border border-white/40 dark:border-zinc-800/40 transition-all duration-200 group-hover:shadow-xl`}
                      title={`${format(
                        d.day,
                        "MMM d, yyyy"
                      )}: ${d.count} contributions`}
                    />
                  );
                })}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  </div>
);


export default ContributionGraph;
