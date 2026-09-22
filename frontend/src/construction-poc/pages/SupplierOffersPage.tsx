import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const SupplierOffersPage: React.FC = () => {
  const navigate = useNavigate();
  const { offer, supplier, requirement, project } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/requirements")} className="hover:underline">
          {requirement.code}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Supplier Offer</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Supplier Offer {offer.offerReference}</h1>
          <p className="text-sm text-[#44474c]">Issued by {supplier.name} for Requirement {requirement.code}</p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-semibold">
          {offer.status}
        </span>
      </div>

      {/* Offer Summary Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div className="p-4 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block mb-1">Total Contract Price</span>
            <span className="text-2xl font-bold text-[#181c1e]">{offer.totalPrice}</span>
            <span className="text-xs text-[#44474c] block mt-1">{offer.unitPrice}</span>
          </div>

          <div className="p-4 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block mb-1">Promised Delivery Date</span>
            <span className="text-2xl font-bold text-[#181c1e]">{offer.promisedDeliveryDate}</span>
            <span className="text-xs text-[#006c49] font-medium block mt-1">Guaranteed SLA</span>
          </div>

          <div className="p-4 bg-[#f7fafc] rounded-lg border border-[#e0e3e5]">
            <span className="text-xs text-[#44474c] block mb-1">Supplier Trust ID</span>
            <span className="text-2xl font-bold text-[#041627]">{supplier.id}</span>
            <span className="text-xs text-[#006c49] font-medium block mt-1">Verified Supplier</span>
          </div>
        </div>

        <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
          <div className="text-xs text-[#44474c]">
            Linked Requirement: <span className="font-semibold text-[#181c1e]">{requirement.materialName} ({requirement.quantity} {requirement.unit})</span>
          </div>
          <button
            onClick={() => navigate("/dashboard/construction/deliveries")}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
          >
            View Active Delivery (DEL-1042) →
          </button>
        </div>
      </div>
    </div>
  );
};

export default SupplierOffersPage;
