import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const InspectionPage: React.FC = () => {
  const navigate = useNavigate();
  const { inspection, delivery, project, requirement } = SCENARIO_DATA;

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
        <span className="font-semibold text-[#041627]">Inspection</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Quality Inspection {inspection.reference}</h1>
          <p className="text-sm text-[#44474c]">Inspector: {inspection.inspector} • Date: {inspection.inspectionDate}</p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-bold uppercase flex items-center gap-1">
          <span className="material-symbols-outlined text-[16px]">verified</span>
          {inspection.status}
        </span>
      </div>

      {/* Inspection Details Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <h2 className="text-lg font-bold text-[#181c1e] mb-4">{inspection.title}</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
              Ultrasonic Weld Integrity
            </span>
            <span className="text-sm font-bold text-[#006c49] flex items-center gap-1">
              <span className="material-symbols-outlined text-[18px]">check_circle</span>
              {inspection.ultrasonicWeldIntegrity}
            </span>
          </div>

          <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
              Dimensional Compliance
            </span>
            <span className="text-sm font-bold text-[#006c49] flex items-center gap-1">
              <span className="material-symbols-outlined text-[18px]">check_circle</span>
              {inspection.dimensionalCompliance}
            </span>
          </div>
        </div>

        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg mb-6">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">Inspector Field Notes</span>
          <p className="text-sm text-[#181c1e] leading-relaxed">{inspection.notes}</p>
        </div>

        {/* Action Link to Provenance */}
        <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
          <div className="text-xs text-[#44474c]">
            Requirement Verified: <span className="font-semibold text-[#181c1e]">{requirement.code}</span>
          </div>
          <button
            onClick={() => navigate("/dashboard/construction/provenance")}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
          >
            View Full Provenance Story →
          </button>
        </div>
      </div>
    </div>
  );
};

export default InspectionPage;
