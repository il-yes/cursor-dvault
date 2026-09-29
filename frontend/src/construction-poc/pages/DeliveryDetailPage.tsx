import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

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
      <div className="flex flex-col w-full min-h-[50vh] items-center justify-center p-6 text-center">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading delivery logistics overview...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex flex-col w-full max-w-md mx-auto my-12 p-6 bg-white border border-[#e0e3e5] rounded-xl text-center shadow-sm font-[Inter]">
        <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
          <span className="material-symbols-outlined text-[28px]">error</span>
        </div>
        <h2 className="text-lg font-bold text-[#041627] mb-1">Delivery Details Unavailable</h2>
        <p className="text-xs text-[#44474c] mb-4">{error || "Delivery logistics overview data unavailable"}</p>
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="px-4 py-2 bg-[#041627] text-white text-xs font-semibold rounded-lg hover:bg-[#1a2b3c] transition-colors cursor-pointer"
        >
          ← Return to Projects
        </button>
      </div>
    );
  }

  const { delivery, project, transport, site, issue, decision, inspection } = data;
  const siteName = site?.name || "Site Alpha - South Pier Foundation";

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 font-[Inter]">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="hover:underline font-medium text-[#041627] cursor-pointer"
        >
          {project.code || project.id}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Deliveries</span>
        <span>/</span>
        <span className="font-semibold text-[#041627]">{delivery.reference || delivery.id}</span>
      </div>

      {/* Main Delivery Banner */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2 text-xs text-[#44474c] font-semibold">
            <span className="bg-[#041627] text-white px-2.5 py-1 rounded">{delivery.reference || delivery.id}</span>
            <span>•</span>
            <span>{siteName}</span>
          </div>
          <span className="px-3 py-1 bg-[#ffb74d]/20 text-[#b76e00] rounded-full text-xs font-semibold uppercase flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-[#b76e00] animate-pulse"></span>
            {delivery.status}
          </span>
        </div>

        <h1 className="text-2xl font-bold text-[#181c1e] mb-2">Delivery {delivery.reference || delivery.id} Tracking</h1>
        <p className="text-sm text-[#44474c] mb-4">{delivery.deliveryNotes || "Critical structural load for Pier 4 framework."}</p>

        {/* Date Comparison Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t border-[#f0f3f5]">
          <div className="p-3 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block">Planned Delivery</span>
            <span className="text-base font-bold text-[#181c1e]">{delivery.plannedDeliveryDate}</span>
          </div>

          <div className="p-3 bg-[#fff3e0] rounded-lg border border-[#ffe0b2]">
            <span className="text-xs text-[#b76e00] font-semibold block">Current ETA (Delayed)</span>
            <span className="text-base font-bold text-[#b76e00]">{delivery.eta}</span>
          </div>

          <div className="p-3 bg-[#e8f5e9] rounded-lg border border-[#c8e6c9]">
            <span className="text-xs text-[#006c49] font-semibold block">Actual Arrival</span>
            <span className="text-base font-bold text-[#006c49]">{delivery.actualDeliveryDate || "Pending"}</span>
          </div>
        </div>
      </div>

      {/* Connected Story / Navigation Flow */}
      <h2 className="text-lg font-bold text-[#181c1e] mb-3">Delivery Journey & Story Navigation</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        {/* Transport Link */}
        <div 
          onClick={() => navigate(CONSTRUCTION_ROUTES.TRANSPORT)}
          className="bg-white border border-[#e0e3e5] rounded-xl p-5 hover:border-[#041627] cursor-pointer transition-all shadow-sm flex items-start gap-4"
        >
          <div className="w-10 h-10 rounded-lg bg-[#041627]/5 text-[#041627] flex items-center justify-center">
            <span className="material-symbols-outlined">local_shipping</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Transport Operation</span>
            <h3 className="font-bold text-[#181c1e] text-base">{transport.reference || transport.id} ({transport.vehicle})</h3>
            <p className="text-xs text-[#44474c] mt-1">{transport.delayReason || "Road Restriction on route"}</p>
          </div>
        </div>

        {/* Issue Link */}
        <div 
          onClick={() => navigate(CONSTRUCTION_ROUTES.ISSUES)}
          className="bg-white border border-[#ffcdd2] rounded-xl p-5 hover:border-[#d32f2f] cursor-pointer transition-all shadow-sm flex items-start gap-4"
        >
          <div className="w-10 h-10 rounded-lg bg-[#ffebee] text-[#d32f2f] flex items-center justify-center">
            <span className="material-symbols-outlined">warning</span>
          </div>
          <div>
            <span className="text-xs text-[#d32f2f] font-semibold uppercase tracking-wider">Construction Issue</span>
            <h3 className="font-bold text-[#181c1e] text-base">{issue?.reference || issue?.id || "ISS-1042"} ({issue?.title || "Road Restriction"})</h3>
            <p className="text-xs text-[#44474c] mt-1">Status: {issue?.status || "OPEN"} • Evidence attached</p>
          </div>
        </div>

        {/* Decision Link */}
        <div 
          onClick={() => navigate(CONSTRUCTION_ROUTES.DECISIONS)}
          className="bg-white border border-[#e0e3e5] rounded-xl p-5 hover:border-[#041627] cursor-pointer transition-all shadow-sm flex items-start gap-4"
        >
          <div className="w-10 h-10 rounded-lg bg-[#041627]/5 text-[#041627] flex items-center justify-center">
            <span className="material-symbols-outlined">gavel</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Approved Decision</span>
            <h3 className="font-bold text-[#181c1e] text-base">{decision?.reference || decision?.id || "DEC-1042"} ({decision?.subject || "Approve Route B"})</h3>
            <p className="text-xs text-[#006c49] font-medium mt-1">Status: {decision?.status || "APPROVED"} by {decision?.decidedBy || "Alex Rivera"}</p>
          </div>
        </div>

        {/* Inspection Link */}
        <div 
          onClick={() => navigate(CONSTRUCTION_ROUTES.INSPECTIONS)}
          className="bg-white border border-[#e0e3e5] rounded-xl p-5 hover:border-[#006c49] cursor-pointer transition-all shadow-sm flex items-start gap-4"
        >
          <div className="w-10 h-10 rounded-lg bg-[#e8f5e9] text-[#006c49] flex items-center justify-center">
            <span className="material-symbols-outlined">fact_check</span>
          </div>
          <div>
            <span className="text-xs text-[#006c49] font-semibold uppercase tracking-wider">Quality Inspection</span>
            <h3 className="font-bold text-[#181c1e] text-base">{inspection?.reference || inspection?.id || "INSP-1042"} ({inspection?.status || "PASSED"})</h3>
            <p className="text-xs text-[#44474c] mt-1">Inspector: {inspection?.inspector || "Sarah Jenkins, PE"}</p>
          </div>
        </div>
      </div>

      {/* Provenance Button Banner */}
      <div className="bg-[#041627] text-white rounded-xl p-6 flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <h3 className="font-bold text-lg mb-1">Why is {delivery.reference || delivery.id} late?</h3>
          <p className="text-sm text-[#b7c8de]">View full multi-layer provenance story across Construction, C3, and TraceCore.</p>
        </div>
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROVENANCE)}
          className="px-5 py-2.5 bg-[#6cf8bb] text-[#041627] font-bold text-sm rounded-lg hover:bg-[#6cf8bb]/90 transition-colors whitespace-nowrap cursor-pointer"
        >
          View Provenance Story →
        </button>
      </div>
    </div>
  );
};

export default DeliveryDetailPage;
