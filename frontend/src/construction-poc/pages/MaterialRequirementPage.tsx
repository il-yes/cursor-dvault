import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getProcurementOverview, ProcurementData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

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
      <div className="flex flex-col w-full min-h-[50vh] items-center justify-center p-6 text-center">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading procurement overview...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex flex-col w-full max-w-md mx-auto my-12 p-6 bg-white border border-[#e0e3e5] rounded-xl text-center shadow-sm">
        <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
          <span className="material-symbols-outlined text-[28px]">error</span>
        </div>
        <h2 className="text-lg font-bold text-[#041627] mb-1">Procurement Details Unavailable</h2>
        <p className="text-xs text-[#44474c] mb-4">{error || "Procurement overview data unavailable"}</p>
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="px-4 py-2 bg-[#041627] text-white text-xs font-semibold rounded-lg hover:bg-[#1a2b3c] transition-colors"
        >
          ← Return to Projects
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 font-[Inter]">
      {/* Navigation Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="hover:underline font-medium text-[#041627]"
        >
          {data.projectCode || "PRJ-001"}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Requirements</span>
        <span>/</span>
        <span className="font-semibold text-[#041627]">{data.code}</span>
      </div>

      {/* Requirement Header */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2 text-xs font-semibold text-[#44474c]">
            <span className="bg-[#f0f3f5] px-2.5 py-1 rounded text-[#041627]">{data.code}</span>
            <span>•</span>
            <span>{data.targetPhase}</span>
          </div>
          <span className="px-3 py-1 bg-[#ffb74d]/20 text-[#b76e00] rounded-full text-xs font-semibold uppercase">
            {data.status.replace("_", " ")}
          </span>
        </div>

        <h1 className="text-2xl font-bold text-[#181c1e] mb-2">{data.materialName}</h1>
        <p className="text-sm text-[#44474c] mb-4">{data.specification}</p>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-[#f0f3f5]">
          <div>
            <span className="text-xs text-[#44474c] block">Required Quantity</span>
            <span className="text-lg font-bold text-[#181c1e]">{data.quantity} {data.unit}</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] block">Target Date</span>
            <span className="text-lg font-bold text-[#181c1e]">{data.requiredDate}</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] block">Material ID</span>
            <span className="text-lg font-bold text-[#041627]">{data.materialId}</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] block">Accepted Supplier</span>
            <span className="text-lg font-bold text-[#006c49]">{data.supplierName || "Apex Steel"}</span>
          </div>
        </div>
      </div>

      {/* Linked Actions / Next Steps */}
      <div className="bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-6 flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <h3 className="font-semibold text-[#181c1e] text-base mb-1">Supplier Offer Secured</h3>
          <p className="text-sm text-[#44474c]">Offer {data.offerReference} accepted for {data.totalPrice}. View delivery tracking.</p>
        </div>
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.OFFERS)}
            className="px-4 py-2 bg-white border border-[#e0e3e5] text-[#041627] font-semibold text-sm rounded-lg hover:bg-[#f0f3f5] transition-colors cursor-pointer"
          >
            View Offer ({data.offerReference})
          </button>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors cursor-pointer"
          >
            Track Delivery →
          </button>
        </div>
      </div>
    </div>
  );
};

export default MaterialRequirementPage;

