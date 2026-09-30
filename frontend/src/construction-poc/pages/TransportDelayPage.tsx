import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const TransportDelayPage: React.FC = () => {
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
            setError("Transport logistics overview unavailable");
          }
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error(`Failed to load transport delay for ${targetDeliveryId}:`, err);
          setError(err?.message || "Failed to load transport delay");
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
        <p className="text-sm font-semibold text-[#44474c]">Loading transport &amp; delay details...</p>
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
          <h2 className="text-lg font-bold text-[#041627] mb-1">Transport Details Unavailable</h2>
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

  const { transport, delivery } = data;

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 1. Header */}
      <header className="sticky top-0 z-30 bg-[#041627] text-white px-4 py-3 shadow-md flex items-center justify-between">
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERY_DETAIL)}
            className="text-white/80 hover:text-white transition-colors cursor-pointer flex items-center justify-center mr-1"
            aria-label="Back to Delivery"
          >
            <span className="material-symbols-outlined text-[22px]">arrow_back</span>
          </button>
          <h1 className="font-bold text-base tracking-tight text-white">Transport &amp; Delay</h1>
        </div>
        <div className="flex items-center gap-3">
          <button type="button" className="text-white/80 hover:text-white transition-colors relative">
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#ffb74d]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-2xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 2. Reference Ribbon */}
        <div className="flex items-center justify-between bg-white border border-[#e0e3e5] rounded-xl px-4 py-3 shadow-sm">
          <div className="flex items-center gap-2">
            <span className="text-xs font-bold text-[#74777f] uppercase tracking-wider">Reference</span>
            <span className="text-sm font-bold text-[#041627]">{delivery.reference || delivery.id || "DEL-1042"}</span>
          </div>
          <span className="px-2.5 py-1 rounded-full bg-[#ba1a1a]/10 text-[#ba1a1a] text-xs font-bold tracking-wide uppercase flex items-center gap-1">
            <span className="w-1.5 h-1.5 rounded-full bg-[#ba1a1a] animate-pulse" />
            {transport.statusTag || "ACTIVE ALERT"}
          </span>
        </div>

        {/* 3. Transport Delay Alert Card */}
        <div className="bg-[#041627] text-white rounded-2xl p-5 shadow-md relative overflow-hidden">
          <div className="flex items-start justify-between mb-3">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 rounded-lg bg-[#ba1a1a] text-white flex items-center justify-center shadow-sm">
                <span className="material-symbols-outlined text-[20px]">warning</span>
              </div>
              <div>
                <h2 className="text-base font-bold tracking-tight text-white">{transport.alertTitle || "Transport Delay Detected"}</h2>
                <p className="text-xs font-medium text-[#ffb4ab]">{transport.alertConstraint || "Route M1 Blocked"}</p>
              </div>
            </div>
            <span className="px-2.5 py-1 rounded-full bg-[#ba1a1a]/30 border border-[#ffb4ab]/30 text-[#ffb4ab] text-xs font-extrabold tracking-wide">
              {transport.delayDurationLabel || "+21h 30m"}
            </span>
          </div>

          <div className="grid grid-cols-2 gap-3 mt-4 pt-4 border-t border-white/10">
            <div className="bg-white/5 rounded-xl p-3 border border-white/5">
              <span className="text-[11px] font-semibold text-white/60 block mb-0.5">Original ETA</span>
              <span className="text-sm font-bold text-white/50 line-through">
                {transport.originalEtaLabel || "Aug 15 • 10:00"}
              </span>
            </div>
            <div className="bg-[#ba1a1a]/20 rounded-xl p-3 border border-[#ba1a1a]/40">
              <span className="text-[11px] font-semibold text-[#ffb4ab] block mb-0.5">Revised ETA</span>
              <span className="text-sm font-extrabold text-white">
                {transport.revisedEtaLabel || "Aug 16 • 07:30"}
              </span>
            </div>
          </div>
        </div>

        {/* 4. Live Corridor Geometry */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#041627] text-[18px]">map</span>
              <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Live Corridor Geometry</h3>
            </div>
            <span className="text-[11px] font-semibold text-[#44474c] bg-[#f0f4f7] px-2 py-0.5 rounded-md">
              {transport.gpsSyncAge || "GPS Feed: Sync 12s ago"}
            </span>
          </div>

          {/* Map visualization canvas */}
          <div className="relative w-full h-48 bg-[#0a192f] rounded-xl overflow-hidden border border-[#1e293b] p-3 flex flex-col justify-between">
            {/* SVG Corridor Route Geometry Overlay */}
            <svg className="absolute inset-0 w-full h-full pointer-events-none" xmlns="http://www.w3.org/2000/svg">
              {/* Grid Lines */}
              <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
                <path d="M 20 0 L 0 0 0 20" fill="none" stroke="#1e293b" strokeWidth="0.5" />
              </pattern>
              <rect width="100%" height="100%" fill="url(#grid)" />

              {/* Planned Blocked Route (M1) - Red dashed path */}
              <path
                d="M 40 140 L 160 140 L 260 140 L 360 140 L 460 140"
                fill="none"
                stroke="#ba1a1a"
                strokeWidth="3"
                strokeDasharray="6 4"
                className="opacity-80"
              />

              {/* Approved Bypass Route (A86 / D914) - Green path */}
              <path
                d="M 120 140 C 140 60, 280 60, 360 140"
                fill="none"
                stroke="#006c49"
                strokeWidth="3.5"
                strokeDasharray="none"
              />
            </svg>

            {/* Top Corridor Nodes */}
            <div className="relative z-10 flex justify-between items-start text-xs font-bold">
              <div className="bg-[#041627]/90 border border-white/20 text-white px-2.5 py-1 rounded-lg backdrop-blur-sm shadow-sm flex items-center gap-1.5">
                <span className="w-2 h-2 rounded-full bg-[#006c49]" />
                <span>{transport.origin || "Origin Hub"}</span>
              </div>
              <div className="bg-[#041627]/90 border border-white/20 text-white px-2.5 py-1 rounded-lg backdrop-blur-sm shadow-sm flex items-center gap-1.5">
                <span className="w-2 h-2 rounded-full bg-[#3b82f6]" />
                <span>{transport.destination || "Viaduct Site"}</span>
              </div>
            </div>

            {/* Center Incident Marker */}
            <div className="relative z-10 mx-auto bg-[#ba1a1a] text-white px-3 py-1 rounded-full text-[11px] font-bold shadow-lg flex items-center gap-1.5 border border-white/30 animate-pulse">
              <span className="material-symbols-outlined text-[14px]">block</span>
              <span>{transport.clearanceConstraint || "M1 Clearance < 4.1m"}</span>
            </div>

            {/* Bottom Vehicle Marker & Bypass Pin */}
            <div className="relative z-10 flex justify-between items-end">
              <div className="bg-[#006c49] text-white px-2.5 py-1 rounded-lg text-xs font-extrabold flex items-center gap-1.5 shadow-md border border-white/30">
                <span className="material-symbols-outlined text-[16px]">local_shipping</span>
                <span>{transport.reference || "TR-1042"}</span>
              </div>
              <div className="flex gap-2 text-[10px] font-bold">
                <div className="bg-black/60 text-[#ffb4ab] px-2 py-0.5 rounded border border-[#ba1a1a]/40 flex items-center gap-1">
                  <span className="w-2 h-0.5 bg-[#ba1a1a]" />
                  <span>{transport.blockedRouteName || "Planned M1 (Blocked)"}</span>
                </div>
                <div className="bg-black/60 text-[#6cf8bb] px-2 py-0.5 rounded border border-[#006c49]/40 flex items-center gap-1">
                  <span className="w-2 h-0.5 bg-[#006c49]" />
                  <span>{transport.bypassRouteName || "Bypass A86 / D914"}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* 5. Constraint Details */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#041627] text-[18px]">report_problem</span>
              <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Constraint Details</h3>
            </div>
            <span className="text-xs font-semibold text-[#44474c]">11:41 AM</span>
          </div>

          <div className="space-y-3">
            <div className="p-3 rounded-xl bg-[#fff3e0] border border-[#ffe0b2]">
              <span className="text-[11px] font-bold text-[#b76e00] uppercase tracking-wide block mb-1">
                Primary Incident
              </span>
              <p className="text-xs text-[#041627] font-semibold leading-relaxed">
                {transport.primaryIncidentText || "Emergency road restriction & oversized load weight limit on planned corridor route M1."}
              </p>
            </div>

            <div className="grid grid-cols-2 gap-3 text-xs">
              <div className="p-2.5 rounded-xl bg-[#f7fafc] border border-[#e0e3e5]">
                <span className="text-[11px] font-medium text-[#44474c] block">Site Phase Impact</span>
                <span className="font-bold text-[#ba1a1a]">
                  {transport.sitePhaseImpact || "Viaduct Deck Assembly Stalled"}
                </span>
              </div>
              <div className="p-2.5 rounded-xl bg-[#f7fafc] border border-[#e0e3e5]">
                <span className="text-[11px] font-medium text-[#44474c] block">Reported By</span>
                <span className="font-bold text-[#041627]">
                  {transport.reportedByText || "FastBuild Dispatch"}
                </span>
              </div>
            </div>

            {/* Driver / Carrier Contact Card */}
            <div className="flex items-center justify-between p-3 rounded-xl bg-[#f7fafc] border border-[#e0e3e5]">
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 rounded-full bg-[#041627] text-white flex items-center justify-center font-bold text-xs">
                  {transport.driverInitials || "TS"}
                </div>
                <div>
                  <span className="text-xs font-bold text-[#041627] block">
                    {transport.driverName || "Tom Smith"}
                  </span>
                  <span className="text-[11px] text-[#44474c]">
                    {transport.carrierName || "FastBuild Logistics"} • {transport.vehicleReference || "FR-920-TG"}
                  </span>
                </div>
              </div>
              <button
                type="button"
                onClick={() => alert(`Calling driver ${transport.driverName || "Tom Smith"} (${transport.carrierName || "FastBuild Logistics"})...`)}
                className="w-8 h-8 rounded-full bg-[#006c49]/10 text-[#006c49] hover:bg-[#006c49] hover:text-white transition-colors flex items-center justify-center cursor-pointer"
                title="Contact Carrier / Driver"
              >
                <span className="material-symbols-outlined text-[18px]">call</span>
              </button>
            </div>
          </div>
        </div>

        {/* 6. Detour Evaluation */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#006c49] text-[18px]">alt_route</span>
              <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Detour Evaluation</h3>
            </div>
            <span className="px-2.5 py-0.5 rounded-full bg-[#006c49]/10 text-[#006c49] text-[11px] font-extrabold tracking-wide uppercase">
              {transport.detourStatusTag || "ROUTE B APPROVED"}
            </span>
          </div>

          <div className="space-y-3 text-xs">
            <div className="p-3 rounded-xl bg-[#f7fafc] border border-[#e0e3e5] space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-[#44474c] font-medium">Corridor Alignment</span>
                <span className="font-bold text-[#041627]">
                  {transport.detourCorridorText || "A86 Bypass > D914 Industrial Link"}
                </span>
              </div>
              <div className="flex justify-between items-center pt-2 border-t border-[#e0e3e5]">
                <span className="text-[#44474c] font-medium">Delta Distance</span>
                <span className="font-bold text-[#041627]">
                  {transport.deltaDistanceText || "+24.0 km"}
                </span>
              </div>
              <div className="flex justify-between items-center pt-2 border-t border-[#e0e3e5]">
                <span className="text-[#44474c] font-medium">Schedule Modification</span>
                <span className="font-bold text-[#ba1a1a]">
                  {transport.scheduleModificationText || "+3h 15m driving + overnight staging"}
                </span>
              </div>
            </div>

            <div className="p-3 rounded-xl bg-[#f0f4f7] border border-[#e0e3e5]">
              <span className="text-[11px] font-bold text-[#041627] block mb-0.5">Formal Clearance Authority</span>
              <p className="text-[11px] text-[#44474c] leading-relaxed">
                {transport.clearanceAuthorityText || "Dual-signed by Structural Lead & Transport Coordinator under protocol DEC-1042."}
              </p>
            </div>
          </div>
        </div>

        {/* 7. Actions */}
        <div className="flex flex-col gap-2 pt-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DECISIONS)}
            className="w-full py-3.5 bg-[#041627] hover:bg-[#1a2b3c] text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>View Decision {transport.decisionReference || "DEC-1042"}</span>
            <span className="material-symbols-outlined text-[16px]">arrow_forward</span>
          </button>
          <button
            type="button"
            onClick={() => alert(`Reported constraint to carrier ${transport.carrierName || "FastBuild Logistics"}`)}
            className="w-full py-3 bg-white border border-[#e0e3e5] hover:bg-[#f7fafc] text-[#041627] text-xs font-bold rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">campaign</span>
            <span>Report Constraint / Contact Carrier</span>
          </button>
        </div>
      </main>

      {/* 8. Fixed Bottom Navigation */}
      <BottomNavigation activeItem="projects" />
    </div>
  );
};

export default TransportDelayPage;

