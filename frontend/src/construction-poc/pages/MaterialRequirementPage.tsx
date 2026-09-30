import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getProcurementOverview, ProcurementData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const MaterialRequirementPage: React.FC = () => {
  const navigate = useNavigate();
  const [data, setData] = useState<ProcurementData | undefined>(undefined);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

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
        <p className="text-sm font-semibold text-[#44474c]">Loading material requirement...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 max-w-md text-center shadow-sm">
          <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
            <span className="material-symbols-outlined text-[28px]">error</span>
          </div>
          <h2 className="text-lg font-bold text-[#041627] mb-1">Procurement Details Unavailable</h2>
          <p className="text-xs text-[#44474c] mb-4">{error || "Procurement overview data unavailable"}</p>
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

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 1. FIXED BUILDFLOW HEADER */}
      <header className="sticky top-0 z-30 bg-white/80 backdrop-blur-xl border-b border-[#e0e3e5] text-[#041627] px-4 py-3 shadow-[0_1px_4px_rgba(0,0,0,0.03)] flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="material-symbols-outlined text-[#006c49] text-[22px]">domain</span>
          <span className="font-bold text-base tracking-tight">BuildFlow</span>
        </div>
        <div className="flex items-center gap-3">
          <button
            type="button"
            className="text-[#44474c] hover:text-[#041627] transition-colors relative cursor-pointer"
            aria-label="Notifications"
          >
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#006c49]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-[#041627] text-white flex items-center justify-center text-xs font-bold shadow-sm">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 2. PROJECT CONTEXT ROW */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1.5 min-w-0">
            <span className="material-symbols-outlined text-[16px] text-[#44474c] shrink-0">apartment</span>
            <span className="text-xs font-semibold text-[#44474c] truncate">
              {data.projectDisplayCode || "PROJ-RT-104"} • {data.projectName || "Metro Line 4 Expansion"}
            </span>
          </div>
          <span className="inline-flex items-center px-2 py-0.5 rounded-full bg-[#6cf8bb] text-[#00714d] text-[11px] font-bold uppercase tracking-wide">
            Active
          </span>
        </div>

        {/* 3. MATERIAL REQUIREMENT TITLE SECTION */}
        <div className="space-y-1">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-[#041627] tracking-tight">Material Requirement</h1>
            <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-[#6cf8bb]/60 text-[#00714d] text-xs font-semibold">
              <span className="w-1.5 h-1.5 rounded-full bg-[#006c49]" />
              OFFER SELECTED
            </span>
          </div>
          <p className="text-xs font-bold text-[#74777f] uppercase tracking-wider">
            {data.code} • {data.materialSubtitle || "Structural beams"}
          </p>
          <p className="text-xs text-[#44474c]">
            {data.materialDescription || "Structural beams required for viaduct section 4."}
          </p>
        </div>

        {/* 4. MATERIAL INFORMATION CARD */}
        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm space-y-4">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-[#ebeef0] flex items-center justify-center text-[#041627]">
                <span className="material-symbols-outlined text-[22px]">view_in_ar</span>
              </div>
              <div>
                <span className="text-[11px] font-bold text-[#74777f] uppercase block">Material ID</span>
                <p className="text-lg font-bold text-[#041627]">{data.materialId}</p>
              </div>
            </div>
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-[#ffdad6] text-[#ba1a1a] text-[11px] font-semibold">
              <span className="material-symbols-outlined text-[14px]">warning</span>
              Critical Path
            </span>
          </div>

          <div className="grid grid-cols-2 gap-3 pt-1">
            <div className="p-2.5 rounded-lg bg-[#f1f4f6] space-y-0.5">
              <span className="text-[11px] font-bold text-[#74777f] flex items-center gap-1">
                <span className="material-symbols-outlined text-[13px]">precision_manufacturing</span>
                Specification
              </span>
              <p className="text-xs font-bold text-[#041627]">{data.specificationGrade || "Grade S355JR"}</p>
              <p className="text-[11px] text-[#74777f]">{data.specificationType || "Hot-Rolled (EN 10025-2)"}</p>
            </div>

            <div className="p-2.5 rounded-lg bg-[#f1f4f6] space-y-0.5">
              <span className="text-[11px] font-bold text-[#74777f] flex items-center gap-1">
                <span className="material-symbols-outlined text-[13px]">scale</span>
                Quantity
              </span>
              <p className="text-xs font-bold text-[#041627]">{data.quantityBeamsText || "84 Beams"}</p>
              <p className="text-[11px] text-[#74777f]">{data.quantityTonsText || `${data.quantity} Metric Tons`}</p>
            </div>

            <div className="p-2.5 rounded-lg bg-[#f1f4f6] space-y-0.5">
              <span className="text-[11px] font-bold text-[#74777f] flex items-center gap-1">
                <span className="material-symbols-outlined text-[13px]">calendar_today</span>
                Required Date
              </span>
              <p className="text-xs font-bold text-[#041627]">Aug 15, 2024</p>
              <p className="text-[11px] text-[#74777f]">{data.requiredDateSlotText || "Slot: 08:00 - 12:00"}</p>
            </div>

            <div className="p-2.5 rounded-lg bg-[#f1f4f6] space-y-0.5">
              <span className="text-[11px] font-bold text-[#74777f] flex items-center gap-1">
                <span className="material-symbols-outlined text-[13px]">location_on</span>
                Site Destination
              </span>
              <p className="text-xs font-bold text-[#041627] truncate">{data.siteCityText || "Site-001 (Paris)"}</p>
              <p className="text-[11px] text-[#74777f] truncate">{data.siteLocationText || data.siteName}</p>
            </div>
          </div>

          <div className="p-2.5 rounded-lg bg-[#f1f4f6] flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="material-symbols-outlined text-[#74777f] text-[18px]">account_tree</span>
              <div>
                <span className="text-[11px] font-bold text-[#74777f] block">Construction Phase</span>
                <span className="text-xs font-medium text-[#041627]">{data.constructionPhaseText || data.targetPhase}</span>
              </div>
            </div>
            <span className="material-symbols-outlined text-[#74777f] text-[18px]">check_circle</span>
          </div>
        </div>

        {/* 5. SUPPLIER SOURCING CARD */}
        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm space-y-3.5">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="material-symbols-outlined text-[#041627] text-[20px]">hub</span>
              <h2 className="text-lg font-bold text-[#041627]">Supplier Sourcing</h2>
            </div>
            <span className="text-xs font-semibold text-[#006c49]">100% Filled</span>
          </div>

          {/* 3 Summary Metrics */}
          <div className="grid grid-cols-3 gap-2">
            <div className="p-2 rounded-lg bg-[#f1f4f6] text-center">
              <span className="text-lg font-bold text-[#041627] block leading-tight">{data.invitedCount ?? 3}</span>
              <span className="text-[11px] text-[#74777f]">Invited</span>
            </div>
            <div className="p-2 rounded-lg bg-[#f1f4f6] text-center">
              <span className="text-lg font-bold text-[#041627] block leading-tight">{data.offersReceivedCount ?? 2}</span>
              <span className="text-[11px] text-[#74777f]">Offers Rec'd</span>
            </div>
            <div className="p-2 rounded-lg bg-[#6cf8bb]/40 text-center">
              <span className="text-lg font-bold text-[#006c49] block leading-tight">{data.selectedCount ?? 1}</span>
              <span className="text-[11px] font-semibold text-[#00714d]">Selected</span>
            </div>
          </div>

          {/* 100% Progress Bar */}
          <div className="w-full bg-[#ebeef0] rounded-full h-1.5 flex overflow-hidden">
            <div className="bg-[#006c49] h-full rounded-full transition-all duration-300 w-full" />
          </div>

          {/* 6. SELECTED SUPPLIER CARD */}
          <div className="bg-[#f1f4f6] rounded-xl p-3.5 space-y-3 border border-[#e0e3e5]">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-full bg-[#1a2b3c] text-white flex items-center justify-center text-xs font-bold">
                  ES
                </div>
                <div>
                  <div className="flex items-center gap-1">
                    <span className="text-sm font-bold text-[#041627]">{data.supplierName || "EuroSteel Construction"}</span>
                    <span className="material-symbols-outlined text-[#006c49] text-[16px]">verified</span>
                  </div>
                  <span className="text-[11px] text-[#74777f]">Ref: {data.offerReference || "OFF-1042"}</span>
                </div>
              </div>
              <span className="inline-flex items-center px-2 py-0.5 rounded-full bg-[#006c49] text-white text-[11px] font-semibold">
                Chosen
              </span>
            </div>

            <div className="flex items-baseline justify-between pt-1">
              <div>
                <span className="text-[11px] text-[#74777f] block">Total Contract Value</span>
                <span className="text-xl font-bold text-[#041627]">{data.contractValueText || data.lowestBidTotal || "€142,500"}</span>
              </div>
              <span className="text-xs text-[#74777f] font-medium">{data.unitPriceText || "€1,187.50 / ton"}</span>
            </div>

            <div className="flex items-center gap-2 pt-1 text-[11px] text-[#74777f]">
              <span className="material-symbols-outlined text-[15px] text-[#006c49]">local_shipping</span>
              <span>
                Proposed Delivery: <strong className="text-[#041627] font-medium">{data.proposedDeliveryText || "Aug 15, 2024 at 10:00 AM"}</strong>
              </span>
            </div>
          </div>
        </div>

        {/* 7. ACTION BUTTONS */}
        <div className="pt-2 flex flex-col gap-2.5">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.OFFERS)}
            className="w-full py-3 rounded-xl bg-[#041627] text-white text-xs font-bold flex items-center justify-center gap-2 shadow-sm hover:bg-[#1a2b3c] transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[20px]">visibility</span>
            <span>View Offers (2)</span>
          </button>
          <button
            type="button"
            onClick={() => alert("Invite Supplier dialog triggered.")}
            className="w-full py-2.5 rounded-xl bg-[#e5e9eb] text-[#041627] text-xs font-semibold flex items-center justify-center gap-2 hover:bg-[#d0d4d8] transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[18px]">person_add</span>
            <span>Invite Supplier</span>
          </button>
        </div>
      </main>

      {/* 8. BOTTOM NAVIGATION */}
      <BottomNavigation activeTab="projects" />
    </div>
  );
};

export default MaterialRequirementPage;


