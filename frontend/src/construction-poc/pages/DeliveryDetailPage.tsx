import React from "react";
import { useNavigate } from "react-router-dom";
import { getDelivery, getProject, getTransportDelay, getIssue, getDecision, getInspection } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

export const DeliveryDetailPage: React.FC = () => {
  const navigate = useNavigate();
  const delivery = getDelivery();
  const project = getProject("PRJ-001");
  const transport = getTransportDelay();
  const issue = getIssue();
  const decision = getDecision();
  const inspection = getInspection();

  if (!delivery || !project || !transport || !issue || !decision || !inspection) {
    return <div>Delivery details unavailable</div>;
  }

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 font-[Inter]">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="hover:underline font-medium text-[#041627] cursor-pointer"
        >
          {project.code}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Deliveries</span>
        <span>/</span>
        <span className="font-semibold text-[#041627]">{delivery.reference}</span>
      </div>

      {/* Main Delivery Banner */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2 text-xs text-[#44474c] font-semibold">
            <span className="bg-[#041627] text-white px-2.5 py-1 rounded">{delivery.reference}</span>
            <span>•</span>
            <span>{delivery.siteName}</span>
          </div>
          <span className="px-3 py-1 bg-[#ffb74d]/20 text-[#b76e00] rounded-full text-xs font-semibold uppercase flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-[#b76e00] animate-pulse"></span>
            {delivery.status}
          </span>
        </div>

        <h1 className="text-2xl font-bold text-[#181c1e] mb-2">Delivery {delivery.reference} Tracking</h1>
        <p className="text-sm text-[#44474c] mb-4">{delivery.deliveryNotes}</p>

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
            <span className="text-base font-bold text-[#006c49]">{delivery.actualDeliveryDate}</span>
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
            <h3 className="font-bold text-[#181c1e] text-base">{transport.id} ({transport.vehicle})</h3>
            <p className="text-xs text-[#44474c] mt-1">{transport.delayReason}</p>
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
            <h3 className="font-bold text-[#181c1e] text-base">{issue.reference} ({issue.title})</h3>
            <p className="text-xs text-[#44474c] mt-1">Status: {issue.status} • Evidence attached</p>
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
            <h3 className="font-bold text-[#181c1e] text-base">{decision.reference} ({decision.title})</h3>
            <p className="text-xs text-[#006c49] font-medium mt-1">Status: {decision.status} by {decision.approvedBy}</p>
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
            <h3 className="font-bold text-[#181c1e] text-base">{inspection.reference} ({inspection.status})</h3>
            <p className="text-xs text-[#44474c] mt-1">Inspector: {inspection.inspector}</p>
          </div>
        </div>
      </div>

      {/* Provenance Button Banner */}
      <div className="bg-[#041627] text-white rounded-xl p-6 flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <h3 className="font-bold text-lg mb-1">Why is DEL-1042 late?</h3>
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
