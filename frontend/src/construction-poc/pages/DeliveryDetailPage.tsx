import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const DeliveryDetailPage: React.FC = () => {
  const navigate = useNavigate();
  const { deliveryId } = useParams<{ deliveryId?: string }>();
  const [data, setData] = useState<LogisticsOverviewData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const targetDeliveryId = deliveryId || "DEL-1042";

  useEffect(() => {
    let isMounted = true;
    setLoading(true);
    setError(null);

    getLogisticsOverview(targetDeliveryId)
      .then((res) => {
        if (isMounted) {
          if (res) {
            setData(res);
          } else {
            setError("Logistics overview unavailable");
          }
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error(`Failed to load logistics overview for ${targetDeliveryId}:`, err);
          setError(err?.message || "Failed to load logistics overview");
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [targetDeliveryId]);

  if (loading) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading delivery details...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 max-w-md text-center shadow-sm">
          <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
            <span className="material-symbols-outlined text-[28px]">error</span>
          </div>
          <h2 className="text-lg font-bold text-[#041627] mb-1">Delivery Details Unavailable</h2>
          <p className="text-xs text-[#44474c] mb-4">{error || "Data unavailable"}</p>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="px-4 py-2 bg-[#041627] text-white text-xs font-semibold rounded-lg hover:bg-[#1a2b3c] transition-colors cursor-pointer"
          >
            ← Return to Projects
          </button>
        </div>
      </div>
    );
  }

  const { delivery, requirement, material, supplier, offer, site, transport, issue } = data;

  const stages = delivery.timelineStages || [
    { stageNumber: 1, title: "Planned", subtitle: "Aug 10 — Dispatched Work Order", state: "completed" },
    { stageNumber: 2, title: "Dispatched", subtitle: "Aug 15 — 06:30 AM (Depot North)", state: "completed" },
    { stageNumber: 3, title: "In Transit", subtitle: "Heavy convoy escort confirmed", state: "completed" },
    { stageNumber: 4, title: "Delayed on Route M1", subtitle: "Reported Aug 15, 11:41 AM (Bridge Clearance restriction)", state: "active", badge: "Active" },
    { stageNumber: 5, title: "Rerouted", subtitle: "Transport Rerouted — Alternative Route B in progress", state: "pending" },
    { stageNumber: 6, title: "Delivered", subtitle: "Site gate check-in", state: "future" },
    { stageNumber: 7, title: "Inspected", subtitle: "Quality control validation", state: "future" },
    { stageNumber: 8, title: "Accepted", subtitle: "Sign-off and structural release", state: "future" }
  ];

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 4. Header */}
      <header className="sticky top-0 z-30 bg-[#041627] text-white px-4 py-3 shadow-md flex items-center justify-between">
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
            className="text-white/80 hover:text-white transition-colors cursor-pointer flex items-center justify-center mr-1"
            aria-label="Back to Deliveries"
          >
            <span className="material-symbols-outlined text-[22px]">arrow_back</span>
          </button>
          <span className="font-bold text-base tracking-tight">Delivery Detail</span>
        </div>
        <div className="flex items-center gap-3">
          <button type="button" className="text-white/80 hover:text-white transition-colors relative">
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#ba1a1a]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-2xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 5. Delivery Headline */}
        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <span className="px-2.5 py-1 rounded-full bg-[#ba1a1a]/10 text-[#ba1a1a] text-xs font-extrabold tracking-wide uppercase flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-[#ba1a1a] animate-pulse" />
              {delivery.statusTag || "DELAYED"}
            </span>
            <span className="text-xs font-medium text-[#44474c]">
              {delivery.lastUpdatedAge || "Updated 4m ago"}
            </span>
          </div>

          <div>
            <h1 className="text-2xl font-bold text-[#041627] tracking-tight">
              {delivery.headlineTitle || `Delivery #${delivery.reference || delivery.id || "DEL-1042"}`}
            </h1>
            <p className="text-xs text-[#44474c] font-medium mt-0.5">
              {delivery.headlineQuantity || "Structural beams (84 Units / 120t)"}
            </p>
          </div>

          {/* Road Restriction Alert Banner */}
          <div className="bg-[#fff3e0] border border-[#ffe0b2] rounded-xl p-3.5 flex items-start gap-3">
            <div className="w-8 h-8 rounded-lg bg-[#ba1a1a] text-white flex items-center justify-center shrink-0">
              <span className="material-symbols-outlined text-[20px]">warning</span>
            </div>
            <div>
              <h2 className="text-xs font-bold text-[#ba1a1a] uppercase tracking-wide">
                {delivery.alertTitle || "Road Restriction on M1"}
              </h2>
              <p className="text-xs text-[#041627] font-semibold mt-0.5 leading-relaxed">
                {delivery.alertDescription || "Heavy transport clearance halted. Re-routing assessment currently underway."}
              </p>
            </div>
          </div>
        </div>

        {/* 6. Quick Actions */}
        <div className="grid grid-cols-2 gap-3">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.TRANSPORT)}
            className="py-3 px-4 bg-[#041627] hover:bg-[#1a2b3c] text-white text-xs font-bold rounded-xl shadow-sm transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">local_shipping</span>
            <span>Transport {transport.reference || "TR-1042"}</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.ISSUES)}
            className="py-3 px-4 bg-white border border-[#ba1a1a]/30 hover:bg-[#fff3e0] text-[#ba1a1a] text-xs font-bold rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer shadow-sm"
          >
            <span className="material-symbols-outlined text-[16px]">report_problem</span>
            <span>Issue {issue?.reference || "ISS-1042"}</span>
          </button>
        </div>

        {/* 7. Schedule Comparison */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Schedule Comparison</h3>
            <span className="px-2.5 py-0.5 rounded-full bg-[#ba1a1a]/10 text-[#ba1a1a] text-xs font-extrabold">
              {delivery.delayBadgeLabel || "+21.5h Delay"}
            </span>
          </div>

          <div className="grid grid-cols-2 gap-3 mb-3">
            <div className="p-3 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-semibold text-[#44474c] block mb-0.5">Planned Delivery</span>
              <span className="text-base font-bold text-[#041627] block">
                {delivery.plannedDateLabel || "Aug 15"}
              </span>
              <span className="text-xs text-[#74777f]">
                {delivery.plannedTimeLabel || "10:00 AM"}
              </span>
            </div>

            <div className="p-3 bg-[#fff3e0] rounded-xl border border-[#ffe0b2]">
              <span className="text-[11px] font-semibold text-[#ba1a1a] block mb-0.5">Revised ETA</span>
              <span className="text-base font-extrabold text-[#ba1a1a] block">
                {delivery.revisedDateLabel || "Aug 16"}
              </span>
              <span className="text-xs font-bold text-[#ba1a1a]">
                {delivery.revisedTimeLabel || "07:30 AM"}
              </span>
            </div>
          </div>

          <div className="p-3 bg-[#f0f4f7] rounded-xl border border-[#e0e3e5]">
            <span className="text-[11px] font-bold text-[#041627] block mb-0.5">Primary Delay Factor</span>
            <p className="text-xs text-[#44474c] leading-relaxed">
              {delivery.primaryDelayFactor || "Road restriction on planned route M1 reported at 11:41 by Fleet Telematics."}
            </p>
          </div>
        </div>

        {/* 8. Operational Specifications */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider mb-3">Operational Specifications</h3>

          <div className="grid grid-cols-2 gap-3 text-xs mb-3">
            <div className="p-2.5 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-medium text-[#74777f] block">Requirement</span>
              <span className="font-bold text-[#041627]">{requirement.code || "REQ-STRUCT-001"}</span>
            </div>

            <div className="p-2.5 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-medium text-[#74777f] block">Material Grade</span>
              <span className="font-bold text-[#041627] block">{material.id || "MAT-STRUCT-001"}</span>
              <span className="text-[11px] text-[#44474c]">{delivery.materialGradeText || "S355JR Heavy Steel"}</span>
            </div>

            <div className="p-2.5 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-medium text-[#74777f] block">Supplier</span>
              <span className="font-bold text-[#041627] block">{delivery.supplierNameText || supplier.name || "EuroSteel"}</span>
              <span className="text-[11px] text-[#44474c]">{delivery.supplierCodeText || supplier.id || "SUP-001"}</span>
            </div>

            <div className="p-2.5 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-medium text-[#74777f] block">Contract Offer</span>
              <span className="font-bold text-[#041627]">{delivery.contractOfferCodeText || offer.reference || "OFF-1042"}</span>
            </div>

            <div className="p-2.5 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-medium text-[#74777f] block">Destination Site</span>
              <span className="font-bold text-[#041627] block">{delivery.destinationSiteNameText || site.name || "Riverside Tower"}</span>
              <span className="text-[11px] text-[#44474c]">{delivery.destinationBayText || "SITE-001 (Bay 3B)"}</span>
            </div>

            <div className="p-2.5 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <span className="text-[11px] font-medium text-[#74777f] block">Carrier Agency</span>
              <span className="font-bold text-[#041627] block">{delivery.carrierAgencyText || "FastBuild Logistics"}</span>
              <span className="text-[11px] text-[#44474c]">{delivery.carrierVehicleText || transport.reference || "TR-1042"}</span>
            </div>
          </div>

          {/* Linked Critical Incident Banner */}
          <div className="bg-[#fff3e0] border border-[#ffe0b2] rounded-xl p-3 flex items-center justify-between text-xs">
            <div className="flex items-center gap-2">
              <span className="material-symbols-outlined text-[#ba1a1a] text-[18px]">report</span>
              <div>
                <span className="text-[11px] font-bold text-[#ba1a1a] uppercase block">Linked Critical Incident</span>
                <span className="font-bold text-[#041627]">
                  {delivery.linkedIncidentReference || issue?.reference || "ISS-1042"} ({delivery.linkedIncidentLabel || "Critical Delay"})
                </span>
              </div>
            </div>
            <button
              type="button"
              onClick={() => navigate(CONSTRUCTION_ROUTES.ISSUES)}
              className="px-2.5 py-1 bg-[#ba1a1a] text-white font-bold rounded-lg hover:bg-[#ba1a1a]/90 transition-colors cursor-pointer text-[11px]"
            >
              View Issue →
            </button>
          </div>
        </div>

        {/* 9. Last Known Telemetry */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#041627] text-[18px]">radar</span>
              <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Last Known Telemetry</h3>
            </div>
            <span className="px-2.5 py-0.5 rounded-full bg-[#006c49]/10 text-[#006c49] text-[11px] font-extrabold">
              {delivery.telemetryStatusText || "GPS Active"}
            </span>
          </div>

          <div className="relative w-full h-36 bg-[#0a192f] rounded-xl overflow-hidden border border-[#1e293b] p-3 flex flex-col justify-between">
            <svg className="absolute inset-0 w-full h-full pointer-events-none" xmlns="http://www.w3.org/2000/svg">
              <pattern id="telemetry-grid" width="20" height="20" patternUnits="userSpaceOnUse">
                <path d="M 20 0 L 0 0 0 20" fill="none" stroke="#1e293b" strokeWidth="0.5" />
              </pattern>
              <rect width="100%" height="100%" fill="url(#telemetry-grid)" />
              <path d="M 20 100 L 180 100 L 380 100" fill="none" stroke="#ba1a1a" strokeWidth="2.5" strokeDasharray="4 4" />
            </svg>

            <div className="relative z-10 flex justify-between items-start text-xs">
              <div className="bg-[#041627]/90 border border-white/20 text-white px-2.5 py-1 rounded-lg text-[11px] font-bold">
                {delivery.telemetryLocationText || "Stationary: M1 Northbound (KM 84.2)"}
              </div>
            </div>

            <div className="relative z-10 flex justify-between items-end">
              <div className="bg-[#ba1a1a] text-white px-2.5 py-1 rounded-lg text-xs font-extrabold flex items-center gap-1.5 shadow-md">
                <span className="material-symbols-outlined text-[16px]">pause_circle</span>
                <span>{delivery.telemetryStateText || "Stationary"}</span>
              </div>
            </div>
          </div>
        </div>

        {/* 10. Progress Timeline */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Progress Timeline</h3>
            <span className="text-xs font-bold text-[#44474c] bg-[#f0f4f7] px-2.5 py-0.5 rounded-full">
              {delivery.timelineStageText || "Stage 4 of 8"}
            </span>
          </div>

          <div className="relative pl-6 space-y-6 before:absolute before:left-[11px] before:top-2 before:bottom-2 before:w-[2px] before:bg-[#e0e3e5]">
            {stages.map((stage) => {
              const isCompleted = stage.state === "completed";
              const isActive = stage.state === "active";
              const isPending = stage.state === "pending";

              let dotClass = "bg-white border-2 border-[#a8abad]";
              if (isCompleted) dotClass = "bg-[#006c49] border-2 border-[#006c49] text-white";
              if (isActive) dotClass = "bg-[#ba1a1a] border-2 border-[#ba1a1a] text-white animate-pulse";
              if (isPending) dotClass = "bg-[#ffb74d] border-2 border-[#ffb74d] text-white";

              return (
                <div key={stage.stageNumber} className="relative flex items-start justify-between gap-2">
                  <div className={`absolute -left-[31px] top-0.5 w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold ${dotClass}`}>
                    {isCompleted ? "✓" : stage.stageNumber}
                  </div>
                  <div>
                    <h4 className={`text-xs font-bold ${isActive ? "text-[#ba1a1a]" : "text-[#041627]"}`}>
                      {stage.title}
                    </h4>
                    <p className="text-[11px] text-[#44474c] mt-0.5 leading-relaxed">
                      {stage.subtitle}
                    </p>
                  </div>
                  {stage.badge && (
                    <span className="px-2 py-0.5 rounded bg-[#ba1a1a]/10 text-[#ba1a1a] text-[10px] font-extrabold uppercase shrink-0">
                      {stage.badge}
                    </span>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        {/* 11. Final Actions */}
        <div className="flex flex-col gap-2 pt-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROVENANCE)}
            className="w-full py-3.5 bg-[#041627] hover:bg-[#1a2b3c] text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>View Provenance &amp; Root Cause Analysis</span>
            <span className="material-symbols-outlined text-[16px]">arrow_forward</span>
          </button>

          <button
            type="button"
            onClick={() => alert(`Printing consignment manifest & waybill for ${delivery.reference || "DEL-1042"}...`)}
            className="w-full py-3 bg-white border border-[#e0e3e5] hover:bg-[#f7fafc] text-[#041627] text-xs font-bold rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">print</span>
            <span>Print Consignment Manifest &amp; Waybill</span>
          </button>
        </div>
      </main>

      {/* Fixed Bottom Navigation */}
      <BottomNavigation activeItem="projects" />
    </div>
  );
};

export default DeliveryDetailPage;

