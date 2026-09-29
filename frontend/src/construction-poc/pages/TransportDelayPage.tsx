import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

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
      <div className="flex flex-col w-full min-h-[50vh] items-center justify-center p-6 text-center font-[Inter]">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading transport delay overview...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex flex-col w-full max-w-md mx-auto my-12 p-6 bg-white border border-[#e0e3e5] rounded-xl text-center shadow-sm font-[Inter]">
        <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
          <span className="material-symbols-outlined text-[28px]">error</span>
        </div>
        <h2 className="text-lg font-bold text-[#041627] mb-1">Transport Details Unavailable</h2>
        <p className="text-xs text-[#44474c] mb-4">{error || "Transport delay overview data unavailable"}</p>
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

  const { transport, delivery, project, issue } = data;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 font-[Inter]">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)} className="hover:underline font-medium text-[#041627] cursor-pointer">
          {project.code || project.id}
        </button>
        <span>/</span>
        <button onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERY_DETAIL)} className="hover:underline font-medium text-[#041627] cursor-pointer">
          {delivery.reference || delivery.id}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Transport &amp; Delay</span>
      </div>

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Transport Operation {transport.reference || transport.id}</h1>
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
              <p className="text-xs text-[#181c1e] font-semibold">{transport.delayReason || "Road Restriction on route"}</p>
              <p className="text-xs text-[#44474c] mt-1">{Array.isArray(transport.constraints) ? transport.constraints.join(", ") : transport.constraints}</p>
            </div>
          </div>
        </div>

        {issue && (
          <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
            <span className="text-xs text-[#44474c]">
              Linked Issue: <span className="font-semibold text-[#181c1e]">{issue.reference || issue.id} ({issue.title})</span>
            </span>
            <button
              onClick={() => navigate(CONSTRUCTION_ROUTES.ISSUES)}
              className="px-4 py-2 bg-[#d32f2f] text-white font-semibold text-sm rounded-lg hover:bg-[#d32f2f]/90 transition-colors cursor-pointer"
            >
              View Issue &amp; Evidence ({issue.reference || issue.id}) →
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default TransportDelayPage;
