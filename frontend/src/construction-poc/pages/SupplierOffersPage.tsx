import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getProcurementOverview, ProcurementData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

export const SupplierOffersPage: React.FC = () => {
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
        <p className="text-sm font-semibold text-[#44474c]">Loading supplier offer overview...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex flex-col w-full max-w-md mx-auto my-12 p-6 bg-white border border-[#e0e3e5] rounded-xl text-center shadow-sm font-[Inter]">
        <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
          <span className="material-symbols-outlined text-[28px]">error</span>
        </div>
        <h2 className="text-lg font-bold text-[#041627] mb-1">Supplier Offers Unavailable</h2>
        <p className="text-xs text-[#44474c] mb-4">{error || "Supplier offer data unavailable"}</p>
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
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="hover:underline font-medium text-[#041627]"
        >
          {data.projectCode || "PRJ-001"}
        </button>
        <span>/</span>
        <button
          type="button"
          onClick={() => navigate(CONSTRUCTION_ROUTES.REQUIREMENTS)}
          className="hover:underline font-medium text-[#041627]"
        >
          {data.code}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Supplier Offer</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Supplier Offer {data.offerReference}</h1>
          <p className="text-sm text-[#44474c]">Issued by {data.supplierName} for Requirement {data.code}</p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-semibold">
          {data.offerStatus || "ACCEPTED"}
        </span>
      </div>

      {/* Offer Summary Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div className="p-4 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block mb-1">Total Contract Price</span>
            <span className="text-2xl font-bold text-[#181c1e]">{data.totalPrice}</span>
            <span className="text-xs text-[#44474c] block mt-1">{data.unitPrice}</span>
          </div>

          <div className="p-4 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block mb-1">Promised Delivery Date</span>
            <span className="text-2xl font-bold text-[#181c1e]">{data.promisedDeliveryDate}</span>
            <span className="text-xs text-[#006c49] font-medium block mt-1">Guaranteed SLA</span>
          </div>

          <div className="p-4 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block mb-1">Supplier Trust ID</span>
            <span className="text-2xl font-bold text-[#041627]">{data.supplierId || "SUP-001"}</span>
            <span className="text-xs text-[#006c49] font-medium block mt-1">
              {data.isVerifiedSupplier ? "Verified Supplier" : "Registered Supplier"}
            </span>
          </div>
        </div>

        <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
          <div className="text-xs text-[#44474c]">
            Linked Requirement: <span className="font-semibold text-[#181c1e]">{data.materialName} ({data.quantity} {data.unit})</span>
          </div>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors cursor-pointer"
          >
            View Active Delivery (DEL-1042) →
          </button>
        </div>
      </div>
    </div>
  );
};

export default SupplierOffersPage;

