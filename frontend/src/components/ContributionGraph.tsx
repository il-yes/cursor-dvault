import React from "react";
import {
	format,
	subDays,
	eachDayOfInterval,
	startOfWeek,
	startOfMonth,
	isSameMonth,
} from "date-fns";

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

	const weekdayLabels = ["Mon", "-", "Wed", "-", "Fri"];

	return (
		<div className="flex flex-col items-start gap-1 text-xs text-gray-500 dark:text-gray-400">
			{/* Month labels */}
			<div className="flex gap-0.5"  style={{marginLeft: "26px"}}>
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
	);
};

export default ContributionGraph;
