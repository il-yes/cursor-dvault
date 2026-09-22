import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const TransportDelayPage: React.FC = () => {
  const navigate = useNavigate();
  const { transport, delivery, project, issue } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/deliveries")} className="hover:underline">
          {delivery.reference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Transport &amp; Delay</span>
      </div>

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Transport Operation {transport.id}</h1>
          <p className="text-sm text-[#44474c]">Vehicle: {transport.vehicle} • Driver: {transport.driver}</p>
        </div>
        <span className="px-3 py-1 bg-[#ffb74d]/20 text-[#b76e00] rounded-full text-xs font-semibold uppercase">
          {transport.status}
        </span>
      </div>

      {/* Transport Details Grid */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div className="flex flex-col gap-3">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider">Logistics Plan</span>
            
            <div className="flex items-start gap-3 p-3 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
              <span className="material-symbols-outlined text-[#041627]">location_on</span>
              <div>
                <span className="text-xs text-[#44474c] block">Origin</span>
                <span className="text-sm font-bold text-[#181c1e]">{transport.origin}</span>
              </div>
            </div>

            <div className="flex items-start gap-3 p-3 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
              <span className="material-symbols-outlined text-[#006c49]">flag</span>
              <div>
                <span className="text-xs text-[#44474c] block">Destination</span>
                <span className="text-sm font-bold text-[#181c1e]">{transport.destination}</span>
              </div>
            </div>

            <div className="flex items-start gap-3 p-3 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
              <span className="material-symbols-outlined text-[#041627]">alt_route</span>
              <div>
                <span className="text-xs text-[#44474c] block">Active Route</span>
                <span className="text-sm font-bold text-[#181c1e]">{transport.route}</span>
              </div>
            </div>
          </div>

          <div className="flex flex-col gap-3">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider">Schedule &amp; Delay Factor</span>

            <div className="p-3 bg-[#f7fafc] rounded-lg border border-[#e0e3e5] text-xs space-y-1">
              <div className="flex justify-between">
                <span className="text-[#44474c]">Planned Departure:</span>
                <span className="font-semibold text-[#181c1e]">{transport.plannedDeparture}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-[#44474c]">Actual Departure:</span>
                <span className="font-semibold text-[#181c1e]">{transport.actualDeparture}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-[#44474c]">Updated ETA:</span>
                <span className="font-semibold text-[#b76e00]">{transport.eta}</span>
              </div>
            </div>

            <div className="p-4 bg-[#fff3e0] border border-[#ffe0b2] rounded-lg">
              <span className="text-xs text-[#b76e00] font-bold block mb-1">Delay Constraint &amp; Reason</span>
              <p className="text-xs text-[#181c1e] font-semibold">{transport.delayReason}</p>
              <p className="text-xs text-[#44474c] mt-1">{transport.constraints}</p>
            </div>
          </div>
        </div>

        <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
          <span className="text-xs text-[#44474c]">
            Linked Issue: <span className="font-semibold text-[#181c1e]">{issue.reference} ({issue.title})</span>
          </span>
          <button
            onClick={() => navigate("/dashboard/construction/issues")}
            className="px-4 py-2 bg-[#d32f2f] text-white font-semibold text-sm rounded-lg hover:bg-[#d32f2f]/90 transition-colors"
          >
            View Issue &amp; Evidence ({issue.reference}) →
          </button>
        </div>
      </div>
    </div>
  );
};

export default TransportDelayPage;
