import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const MaterialRequirementPage: React.FC = () => {
  const navigate = useNavigate();
  const { requirement, project, offer } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Navigation Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Requirements</span>
        <span>/</span>
        <span className="font-semibold text-[#041627]">{requirement.code}</span>
      </div>

      {/* Requirement Header */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2 text-xs font-semibold text-[#44474c]">
            <span className="bg-[#f0f3f5] px-2.5 py-1 rounded text-[#041627]">{requirement.code}</span>
            <span>•</span>
            <span>{requirement.targetPhase}</span>
          </div>
          <span className="px-3 py-1 bg-[#ffb74d]/20 text-[#b76e00] rounded-full text-xs font-semibold uppercase">
            {requirement.status.replace("_", " ")}
          </span>
        </div>

        <h1 className="text-2xl font-bold text-[#181c1e] mb-2">{requirement.materialName}</h1>
        <p className="text-sm text-[#44474c] mb-4">{requirement.specification}</p>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-[#f0f3f5]">
          <div>
            <span className="text-xs text-[#44474c] block">Required Quantity</span>
            <span className="text-lg font-bold text-[#181c1e]">{requirement.quantity} {requirement.unit}</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] block">Target Date</span>
            <span className="text-lg font-bold text-[#181c1e]">{requirement.requiredDate}</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] block">Material ID</span>
            <span className="text-lg font-bold text-[#041627]">{requirement.materialId}</span>
          </div>
          <div>
            <span className="text-xs text-[#44474c] block">Accepted Supplier</span>
            <span className="text-lg font-bold text-[#006c49]">Apex Steel</span>
          </div>
        </div>
      </div>

      {/* Linked Actions / Next Steps */}
      <div className="bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-6 flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <h3 className="font-semibold text-[#181c1e] text-base mb-1">Supplier Offer Secured</h3>
          <p className="text-sm text-[#44474c]">Offer {offer.offerReference} accepted for {offer.totalPrice}. View delivery tracking.</p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate("/dashboard/construction/offers")}
            className="px-4 py-2 bg-white border border-[#e0e3e5] text-[#041627] font-semibold text-sm rounded-lg hover:bg-[#f0f3f5] transition-colors"
          >
            View Offer ({offer.offerReference})
          </button>
          <button
            onClick={() => navigate("/dashboard/construction/deliveries")}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
          >
            Track Delivery →
          </button>
        </div>
      </div>
    </div>
  );
};

export default MaterialRequirementPage;
