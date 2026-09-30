import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getProcurementOverview, ProcurementData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const SupplierOffersPage: React.FC = () => {
  const navigate = useNavigate();
  const [data, setData] = useState<ProcurementData | undefined>(undefined);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<"all" | "eval" | "audit">("all");

  useEffect(() => {
    let isMounted = true;
    setLoading(true);
    setError(null);

    getProcurementOverview("REQ-STRUCT-001")
      .then((res) => {
        if (isMounted) {
          setData(res);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error("Failed to load procurement overview:", err);
          setError(err?.message || "Failed to load procurement overview");
          setData(undefined);
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  if (loading) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading Supplier Offers...</p>
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
          <h2 className="text-lg font-bold text-[#041627] mb-1">Supplier Offers Unavailable</h2>
          <p className="text-xs text-[#44474c] mb-4">{error || "Data unavailable"}</p>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="px-4 py-2 bg-[#041627] text-white text-xs font-semibold rounded-lg hover:bg-[#1a2b3c] transition-colors"
          >
            ← Return to Projects
          </button>
        </div>
      </div>
    );
  }

  const primary = data.primaryOffer || data.offers[0];
  const secondary = data.secondaryOffer || data.offers[1];

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 1. Fixed BuildFlow Header */}
      <header className="sticky top-0 z-30 bg-[#041627] text-white px-4 py-3 shadow-md flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="material-symbols-outlined text-[#6cf8bb] text-[22px]">domain</span>
          <span className="font-bold text-base tracking-tight">BuildFlow</span>
        </div>
        <div className="flex items-center gap-3">
          <button type="button" className="text-white/80 hover:text-white transition-colors relative">
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#6cf8bb]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-2xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 2. Context Navigation */}
        <div className="flex items-center justify-between">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.REQUIREMENTS)}
            className="flex items-center gap-1 text-xs font-bold text-[#041627] hover:text-[#006c49] transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">arrow_back</span>
            <span>Requirements</span>
          </button>
          <span className="px-2.5 py-1 rounded-full bg-[#006c49]/10 text-[#006c49] text-xs font-bold tracking-wide uppercase">
            Procurement Active
          </span>
        </div>

        {/* 3. Page Heading */}
        <div>
          <h1 className="text-2xl font-bold text-[#041627] tracking-tight">Supplier Offers</h1>
          <p className="text-xs text-[#44474c] font-medium mt-0.5">
            {data.code} • {data.materialName}
          </p>
        </div>

        {/* 4. Segmented Tabs */}
        <div className="flex bg-[#e5e9eb] p-1 rounded-xl gap-1 text-xs font-bold">
          <button
            type="button"
            onClick={() => setActiveTab("all")}
            className={`flex-1 py-2 rounded-lg transition-all ${
              activeTab === "all"
                ? "bg-[#041627] text-white shadow-sm"
                : "text-[#44474c] hover:text-[#041627]"
            }`}
          >
            All Offers ({data.offers?.length || 2})
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("eval")}
            className={`flex-1 py-2 rounded-lg transition-all ${
              activeTab === "eval"
                ? "bg-[#041627] text-white shadow-sm"
                : "text-[#44474c] hover:text-[#041627]"
            }`}
          >
            Evaluation
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("audit")}
            className={`flex-1 py-2 rounded-lg transition-all ${
              activeTab === "audit"
                ? "bg-[#041627] text-white shadow-sm"
                : "text-[#44474c] hover:text-[#041627]"
            }`}
          >
            Audit Trail
          </button>
        </div>

        {/* 5. Summary Metrics (2 Cards Grid) */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-white rounded-xl border border-[#e0e3e5] p-3.5 shadow-sm flex flex-col justify-between">
            <span className="text-[10px] uppercase tracking-wider text-[#44474c] font-bold">
              Lowest Bid Total
            </span>
            <div className="my-1">
              <span className="text-xl font-bold text-[#041627] block">{data.lowestBidTotal}</span>
              <span className="text-[11px] font-bold text-[#006c49]">{data.lowestBidSavings}</span>
            </div>
          </div>

          <div className="bg-white rounded-xl border border-[#e0e3e5] p-3.5 shadow-sm flex flex-col justify-between">
            <span className="text-[10px] uppercase tracking-wider text-[#44474c] font-bold">
              Fastest Delivery
            </span>
            <div className="my-1">
              <span className="text-xl font-bold text-[#041627] block">{data.fastestDeliveryDate}</span>
              <span className="text-[11px] font-bold text-[#006c49]">{data.fastestDeliveryStatus}</span>
            </div>
          </div>
        </div>

        {/* 6. Primary Accepted Offer Card */}
        {primary && (
          <div className="bg-white rounded-xl border-2 border-[#006c49] p-5 shadow-sm flex flex-col gap-4 relative overflow-hidden">
            {/* Top Status & Supplier Ribbon */}
            <div className="flex items-start justify-between gap-2 border-b border-[#f1f4f6] pb-3">
              <div>
                <span className="inline-block px-2.5 py-0.5 rounded-full bg-[#006c49]/10 text-[#006c49] text-[10px] font-bold uppercase tracking-wider mb-1">
                  {primary.offerStatus}
                </span>
                <div className="flex items-center gap-1.5">
                  <h2 className="text-lg font-bold text-[#041627]">{primary.supplierName}</h2>
                  {primary.isVerified && (
                    <span className="material-symbols-outlined text-[#006c49] text-[18px]" title="Verified Supplier">
                      verified
                    </span>
                  )}
                </div>
                <span className="text-xs font-mono text-[#44474c]">Offer {primary.offerReference}</span>
              </div>

              <div className="text-right">
                <span className="text-2xl font-bold text-[#041627] block leading-tight">{primary.totalPrice}</span>
                <span className="text-xs text-[#44474c] font-semibold">{primary.unitPrice}</span>
              </div>
            </div>

            {/* Quick Specs Matrix */}
            <div className="grid grid-cols-2 gap-3 text-xs bg-[#f7fafc] p-3.5 rounded-xl border border-[#e0e3e5]">
              <div>
                <span className="text-[10px] uppercase text-[#44474c] font-bold block">Specification</span>
                <span className="font-bold text-[#041627]">{primary.specification}</span>
              </div>
              <div>
                <span className="text-[10px] uppercase text-[#44474c] font-bold block">Quantity</span>
                <span className="font-bold text-[#041627]">{primary.quantity} {primary.unit}</span>
              </div>
              <div>
                <span className="text-[10px] uppercase text-[#44474c] font-bold block">Delivery Window</span>
                <span className="font-bold text-[#041627]">{primary.deliveryWindow}</span>
              </div>
              <div>
                <span className="text-[10px] uppercase text-[#44474c] font-bold block">Offer Validity</span>
                <span className="font-bold text-[#041627]">{primary.offerValidity}</span>
              </div>
            </div>

            {/* Technical Evaluation Score */}
            <div className="flex items-center justify-between bg-[#006c49]/5 p-3 rounded-xl border border-[#006c49]/20">
              <div className="flex items-center gap-2">
                <span className="w-8 h-8 rounded-lg bg-[#006c49] text-white flex items-center justify-center font-bold text-xs">
                  {primary.complianceScore}
                </span>
                <div>
                  <h4 className="text-xs font-bold text-[#041627]">Technical Evaluation</h4>
                  <p className="text-[10px] text-[#006c49] font-semibold">Compliant with {primary.complianceSpecVersion}</p>
                </div>
              </div>
              <span className="material-symbols-outlined text-[#006c49] text-[20px]">fact_check</span>
            </div>

            {/* Document Pills */}
            <div className="flex flex-wrap gap-1.5 pt-1">
              {primary.certifications.map((doc) => (
                <span
                  key={doc}
                  className="px-2.5 py-1 bg-[#ebeef0] text-[#041627] text-[11px] font-semibold rounded-lg flex items-center gap-1 border border-[#e0e3e5]"
                >
                  <span className="material-symbols-outlined text-[14px] text-[#006c49]">description</span>
                  {doc}
                </span>
              ))}
            </div>

            {/* Status Banner */}
            <div className="bg-[#006c49] text-white text-xs font-bold px-3 py-2 rounded-lg text-center shadow-sm">
              {primary.statusBanner}
            </div>
          </div>
        )}

        {/* 7. Secondary Offer Card */}
        {secondary && (
          <div className="bg-white rounded-xl border border-[#e0e3e5] p-5 shadow-sm flex flex-col gap-3">
            <div className="flex items-start justify-between gap-2">
              <div>
                <span className="inline-block px-2.5 py-0.5 rounded-full bg-[#ca8100]/10 text-[#ca8100] text-[10px] font-bold uppercase tracking-wider mb-1">
                  {secondary.offerStatus}
                </span>
                <h3 className="text-base font-bold text-[#041627]">{secondary.supplierName}</h3>
                <span className="text-xs font-mono text-[#44474c]">Offer {secondary.offerReference}</span>
              </div>

              <div className="text-right">
                <span className="text-xl font-bold text-[#041627] block leading-tight">{secondary.totalPrice}</span>
                <span className="text-xs text-[#44474c] font-semibold">{secondary.unitPrice}</span>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-2 text-xs bg-[#f7fafc] p-3 rounded-lg border border-[#e0e3e5]">
              <div>
                <span className="text-[10px] uppercase text-[#44474c] font-bold block">Lead Time</span>
                <span className="font-semibold text-[#041627]">{secondary.leadTime}</span>
              </div>
              <div>
                <span className="text-[10px] uppercase text-[#44474c] font-bold block">Technical Score</span>
                <span className="font-semibold text-[#041627]">{secondary.complianceScore}/100</span>
              </div>
            </div>

            <div className="flex items-center justify-between pt-2 border-t border-[#f1f4f6]">
              <span className="text-xs text-[#44474c] italic">Non-binding alternative</span>
              <button
                type="button"
                onClick={() => navigate(CONSTRUCTION_ROUTES.REQUIREMENTS)}
                className="text-xs font-bold text-[#041627] hover:text-[#006c49] flex items-center gap-0.5 cursor-pointer"
              >
                <span>Compare Details</span>
                <span className="material-symbols-outlined text-[14px]">arrow_forward</span>
              </button>
            </div>
          </div>
        )}

        {/* 8. Primary Action Area */}
        <div className="flex flex-col gap-2 pt-2">
          <button
            type="button"
            className="w-full py-3.5 bg-[#006c49] text-white font-bold text-sm rounded-xl shadow-md flex items-center justify-center gap-2"
          >
            <span className="material-symbols-outlined text-[20px]">check_circle</span>
            <span>Offer Accepted ({primary?.offerReference || "OFF-1042"})</span>
          </button>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
            className="w-full py-3 bg-[#041627] text-white font-bold text-xs rounded-xl shadow-sm hover:bg-[#1a2b3c] transition-colors flex items-center justify-center gap-2 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[18px]">local_shipping</span>
            <span>View Delivery Plan (DEL-1042)</span>
          </button>
        </div>
      </main>

      {/* 9. Fixed Bottom Navigation */}
      <BottomNavigation />
    </div>
  );
};

export default SupplierOffersPage;



