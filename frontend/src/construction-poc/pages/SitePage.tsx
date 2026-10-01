import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { getSiteData } from "../data";

export const SitePage: React.FC = () => {
  const { siteId } = useParams<{ siteId?: string }>();
  const navigate = useNavigate();
  const site = getSiteData(siteId);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-8">
      {/* Quick Context Bar / Breadcrumbs */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-[#e0e3e5] pb-4">
        <div>
          <div className="flex items-center gap-2 text-xs font-semibold text-[#44474c] uppercase tracking-wider mb-1">
            <button
              onClick={() => navigate(-1)}
              className="inline-flex items-center gap-1 hover:text-[#041627] transition-colors"
            >
              <span className="material-symbols-outlined text-[16px]">arrow_back</span>
              Back
            </button>
            <span>/</span>
            <span>Site Detail</span>
          </div>
          <h1 className="text-2xl font-bold text-[#041627]">
            {site.code} • {site.projectName}
          </h1>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
            className="px-4 py-2 bg-[#041627] text-white text-xs font-bold rounded-lg hover:bg-[#006c49] transition-colors flex items-center gap-2 shadow-sm"
          >
            <span className="material-symbols-outlined text-[18px]">local_shipping</span>
            Track Delivery #{site.incomingDelivery.reference}
          </button>
          <button
            onClick={() => navigate(CONSTRUCTION_ROUTES.INSPECTIONS)}
            className="px-4 py-2 bg-white text-[#041627] border border-[#d0d4d8] text-xs font-bold rounded-lg hover:bg-[#f0f4f8] transition-colors flex items-center gap-2 shadow-sm"
          >
            <span className="material-symbols-outlined text-[18px]">verified</span>
            View Inspection Status
          </button>
        </div>
      </div>

      {/* Large Construction Hero Card */}
      <div className="relative rounded-2xl overflow-hidden shadow-lg border border-[#e0e3e5] bg-[#041627] text-white">
        <div className="absolute inset-0 z-0">
          <img
            src={site.heroImageUrl}
            alt={site.name}
            className="w-full h-full object-cover opacity-45"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-[#041627] via-[#041627]/60 to-transparent" />
        </div>

        <div className="relative z-10 p-6 sm:p-8 flex flex-col justify-end min-h-[260px] space-y-4">
          <div className="flex flex-wrap items-center gap-2">
            <span className="px-3 py-1 bg-[#006c49] text-white text-xs font-bold rounded-full tracking-wide uppercase shadow-sm">
              {site.phaseText}
            </span>
            <span className="px-3 py-1 bg-white/20 backdrop-blur-md text-white text-xs font-medium rounded-full">
              {site.zoneText}
            </span>
          </div>

          <div>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
              {site.name}
            </h2>
            <p className="text-sm text-gray-300 mt-1">
              Active Construction Zone • Site ID: {site.code}
            </p>
          </div>

          <div className="pt-2 flex flex-wrap items-center gap-6 text-xs font-medium text-gray-300 border-t border-white/10">
            <div className="flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px] text-[#006c49]">my_location</span>
              <span>{site.coordinatesText}</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="material-symbols-outlined text-[18px] text-[#006c49]">engineering</span>
              <span>{site.inspectorName} ({site.inspectorRole})</span>
            </div>
          </div>
        </div>
      </div>

      {/* Grid: Site Information & Zone Surveillance Feeds */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Site Information Card */}
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm space-y-6 lg:col-span-1 flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-[#041627] flex items-center gap-2">
                <span className="material-symbols-outlined text-[#006c49]">domain</span>
                Site Information
              </h3>
              <span className="px-2.5 py-1 bg-[#006c49]/10 text-[#006c49] text-xs font-bold rounded-full">
                {site.completionPercent}% Completion
              </span>
            </div>

            {/* Completion Progress Bar */}
            <div className="space-y-1 mb-6">
              <div className="w-full h-2 bg-[#e5e9eb] rounded-full overflow-hidden">
                <div
                  className="h-full bg-[#006c49] rounded-full transition-all duration-500"
                  style={{ width: `${site.completionPercent}%` }}
                />
              </div>
            </div>

            {/* Specifications 2x2 Grid */}
            <div className="grid grid-cols-2 gap-4 border-t border-b border-[#e0e3e5] py-4">
              {site.specifications.map((spec, idx) => (
                <div key={idx} className="space-y-0.5">
                  <span className="text-[11px] font-semibold uppercase tracking-wider text-[#44474c]">
                    {spec.label}
                  </span>
                  <p className="text-sm font-bold text-[#041627]">{spec.value}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="space-y-3 pt-4">
            <div className="flex items-center justify-between p-3 bg-[#f7fafc] rounded-lg border border-[#e0e3e5] text-xs">
              <span className="font-semibold text-[#44474c]">Lead Inspector</span>
              <span className="font-bold text-[#041627]">{site.inspectorName}</span>
            </div>

            <button
              onClick={() => alert(`GPS Location: ${site.coordinatesText}`)}
              className="w-full py-2.5 px-4 bg-[#f0f4f8] hover:bg-[#e2e8f0] text-[#041627] text-xs font-bold rounded-lg transition-colors flex items-center justify-center gap-2"
            >
              <span className="material-symbols-outlined text-[18px]">pin_drop</span>
              Coordinates: {site.coordinatesText}
            </button>
          </div>
        </div>

        {/* Zone Surveillance Feeds */}
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm lg:col-span-2 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-bold text-[#041627] flex items-center gap-2">
              <span className="material-symbols-outlined text-red-500">videocam</span>
              Zone Surveillance Feeds
            </h3>
            <span className="text-xs text-[#44474c] font-medium flex items-center gap-1">
              <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse" />
              Live Site Feeds
            </span>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {site.surveillanceFeeds.map((feed) => (
              <div
                key={feed.id}
                className="relative rounded-lg overflow-hidden border border-[#e0e3e5] bg-[#041627] group"
              >
                <img
                  src={feed.imageUrl}
                  alt={feed.name}
                  className="w-full h-44 object-cover opacity-85 group-hover:scale-105 transition-transform duration-300"
                />
                <div className="absolute top-3 left-3 px-2.5 py-1 bg-red-600 text-white text-[10px] font-bold rounded flex items-center gap-1 shadow">
                  <span className="w-1.5 h-1.5 rounded-full bg-white animate-ping" />
                  {feed.badge}
                </div>
                <div className="absolute bottom-0 inset-x-0 bg-gradient-to-t from-[#041627] to-transparent p-3 text-white">
                  <p className="text-xs font-bold truncate">{feed.name}</p>
                  <p className="text-[10px] text-gray-300">{feed.timeLabel}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Operational Zones Section */}
      <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm space-y-4">
        <div className="flex items-center justify-between border-b border-[#e0e3e5] pb-3">
          <h3 className="text-lg font-bold text-[#041627] flex items-center gap-2">
            <span className="material-symbols-outlined text-[#006c49]">grid_view</span>
            Operational Zones
          </h3>
          <span className="text-xs font-bold text-[#44474c]">
            {site.operationalZones.length} Active Sectors
          </span>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {site.operationalZones.map((zone) => (
            <div
              key={zone.code}
              className="p-4 rounded-xl border border-[#e0e3e5] bg-[#f7fafc] space-y-3 hover:border-[#006c49]/40 transition-colors"
            >
              <div className="flex items-center justify-between">
                <span className="text-xs font-bold bg-[#041627] text-white px-2 py-0.5 rounded">
                  {zone.code}
                </span>
                <span className={`text-[10px] font-bold px-2 py-0.5 rounded border ${zone.statusColor}`}>
                  {zone.status}
                </span>
              </div>
              <h4 className="text-sm font-bold text-[#041627] line-clamp-1">{zone.name}</h4>
              <div className="space-y-1">
                <div className="flex justify-between text-[11px] font-medium text-[#44474c]">
                  <span>Progress</span>
                  <span className="font-bold text-[#041627]">{zone.progressPercent}%</span>
                </div>
                <div className="w-full h-1.5 bg-[#e5e9eb] rounded-full overflow-hidden">
                  <div
                    className="h-full bg-[#006c49] rounded-full"
                    style={{ width: `${zone.progressPercent}%` }}
                  />
                </div>
              </div>
              <p className="text-[11px] text-[#44474c] font-medium flex items-center gap-1 pt-1">
                <span className="material-symbols-outlined text-[14px]">info</span>
                {zone.activityCountText}
              </p>
            </div>
          ))}
        </div>
      </div>

      {/* Storage & Logistics */}
      <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between border-b border-[#e0e3e5] pb-3 gap-2">
          <h3 className="text-lg font-bold text-[#041627] flex items-center gap-2">
            <span className="material-symbols-outlined text-[#006c49]">inventory_2</span>
            Storage & Logistics
          </h3>
          <div className="flex items-center gap-3">
            <span className="text-xs font-bold text-[#44474c]">
              Yard Occupancy: {site.logisticsOccupancyPercent}% Occupied
            </span>
            <div className="w-24 h-2 bg-[#e5e9eb] rounded-full overflow-hidden">
              <div
                className="h-full bg-[#006c49] rounded-full"
                style={{ width: `${site.logisticsOccupancyPercent}%` }}
              />
            </div>
          </div>
        </div>

        {/* Active Inventory 2x2 / 4-card Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {site.inventoryItems.map((item, idx) => (
            <div key={idx} className="p-4 rounded-xl border border-[#e0e3e5] bg-white space-y-2">
              <span className="text-[11px] font-semibold text-[#44474c] uppercase tracking-wider block">
                {item.category}
              </span>
              <p className="text-sm font-bold text-[#041627]">{item.itemCount}</p>
              <div className="w-full h-1 bg-[#e5e9eb] rounded-full overflow-hidden">
                <div
                  className="h-full bg-[#041627] rounded-full"
                  style={{ width: `${item.occupancyPercent}%` }}
                />
              </div>
            </div>
          ))}
        </div>

        {/* Incoming Logistics Banner (DEL-1042) */}
        <div className="p-4 rounded-xl bg-amber-50 border border-amber-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-start sm:items-center gap-3">
            <span className="material-symbols-outlined text-amber-600 text-[24px]">warning</span>
            <div>
              <h4 className="text-sm font-bold text-amber-900">
                Incoming Delivery #{site.incomingDelivery.reference} — {site.incomingDelivery.statusText}
              </h4>
              <p className="text-xs text-amber-700">
                {site.incomingDelivery.description} • {site.incomingDelivery.etaText}
              </p>
            </div>
          </div>
          <button
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
            className="px-4 py-2 bg-amber-600 hover:bg-amber-700 text-white text-xs font-bold rounded-lg transition-colors shadow-sm shrink-0"
          >
            Track Delivery #{site.incomingDelivery.reference}
          </button>
        </div>
      </div>

      {/* Bottom Action Controls */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pt-4 border-t border-[#e0e3e5]">
        <button
          onClick={() => navigate(CONSTRUCTION_ROUTES.INSPECTIONS)}
          className="w-full sm:w-auto px-6 py-3 bg-[#041627] hover:bg-[#006c49] text-white text-sm font-bold rounded-xl transition-colors shadow-sm flex items-center justify-center gap-2"
        >
          <span className="material-symbols-outlined text-[20px]">fact_check</span>
          View Inspection Status ({site.inspectionStatus.reference})
        </button>

        <button
          onClick={() => alert("Site Safety & Access Manifest opened")}
          className="w-full sm:w-auto px-6 py-3 bg-white hover:bg-[#f0f4f8] text-[#041627] border border-[#d0d4d8] text-sm font-bold rounded-xl transition-colors shadow-sm flex items-center justify-center gap-2"
        >
          <span className="material-symbols-outlined text-[20px]">shield</span>
          Site Safety & Access Manifest
        </button>
      </div>
    </div>
  );
};

export default SitePage;
