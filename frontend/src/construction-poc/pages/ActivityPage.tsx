import React from "react";

export const ActivityPage: React.FC = () => {
  const deliveries = [
    {
      id: "DEL-402",
      material: "Ready-Mix Concrete Type C3",
      supplier: "LafargeHolcim",
      status: "DELIVERED",
      statusColor: "bg-[#6cf8bb]/20 text-[#00714d]",
      time: "Today, 09:18 AM",
      details: "12 m³ poured at Foundation Zone A.",
    },
    {
      id: "DEL-401",
      material: "Reinforced Steel Rebar #8",
      supplier: "ArcelorMittal",
      status: "IN TRANSIT",
      statusColor: "bg-[#ffddb8] text-[#ca8100]",
      time: "Expected 02:00 PM",
      details: "Truck #FR-882-AB en route to Gate 3.",
    },
    {
      id: "INS-108",
      material: "Rebar Density Scan - Zone B",
      supplier: "Bureau Veritas",
      status: "SCHEDULED",
      statusColor: "bg-[#e0e3e5] text-[#44474c]",
      time: "Tomorrow, 08:30 AM",
      details: "Pre-pour ultrasound inspection.",
    },
  ];

  return (
    <div className="flex flex-col w-full gap-4 px-4 py-4 max-w-2xl mx-auto pb-8">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-[#181c1e]">Delivery & Activity Log</h1>
        <span className="text-xs font-semibold text-[#006c49] bg-[#006c49]/10 px-2.5 py-1 rounded-full">
          Live Sync
        </span>
      </div>

      <div className="flex flex-col gap-3">
        {deliveries.map((item) => (
          <div
            key={item.id}
            className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col gap-2"
          >
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-[#44474c]">
                {item.id}
              </span>
              <span
                className={`${item.statusColor} px-2 py-0.5 rounded text-[11px] font-semibold`}
              >
                {item.status}
              </span>
            </div>
            <h3 className="text-base font-semibold text-[#181c1e]">
              {item.material}
            </h3>
            <p className="text-sm text-[#44474c]">{item.details}</p>
            <div className="flex items-center justify-between text-xs text-[#74777d] pt-1 border-t border-[#f1f4f6]">
              <span>Supplier: {item.supplier}</span>
              <span>{item.time}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ActivityPage;
