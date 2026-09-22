import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const FieldModePage: React.FC = () => {
  const navigate = useNavigate();
  const { delivery, project, inspection } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full min-h-screen bg-[#181c1e] text-white p-4 max-w-md mx-auto">
      {/* Top Field Mode Header */}
      <div className="flex items-center justify-between pb-4 border-b border-[#33373b] mb-4 pt-2">
        <div className="flex items-center gap-2">
          <span className="material-symbols-outlined text-[#6cf8bb]">smartphone</span>
          <span className="font-bold text-sm tracking-wider uppercase">Field Mode</span>
        </div>
        <button
          onClick={() => navigate("/dashboard/construction")}
          className="text-xs text-[#a0a5aa] hover:text-white px-2 py-1 bg-[#282c30] rounded"
        >
          Exit Field Mode
        </button>
      </div>

      {/* Operational Site Context */}
      <div className="bg-[#23272b] border border-[#33373b] rounded-xl p-4 mb-4">
        <div className="flex justify-between items-center mb-1 text-xs text-[#a0a5aa]">
          <span>{project.code}</span>
          <span className="text-[#6cf8bb] font-semibold">SITE ACTIVE</span>
        </div>
        <h2 className="text-lg font-bold text-white mb-1">{delivery.siteName}</h2>
        <p className="text-xs text-[#a0a5aa]">Lead Superintendent: Sarah Jenkins, PE</p>
      </div>

      {/* Operational Action Card: Material Acceptance */}
      <div className="bg-[#23272b] border border-[#33373b] rounded-xl p-4 mb-4">
        <div className="flex items-center justify-between mb-3">
          <span className="text-xs font-bold text-[#a0a5aa] uppercase tracking-wider">Active Field Task</span>
          <span className="px-2 py-0.5 bg-[#6cf8bb]/20 text-[#6cf8bb] text-[10px] font-bold rounded">INSPECTION COMPLETED</span>
        </div>

        <h3 className="font-bold text-base text-white mb-1">Delivery {delivery.reference}</h3>
        <p className="text-xs text-[#a0a5aa] mb-4">High-Strength Structural Steel Beams (Grade A992) • 120 Metric Tons</p>

        <div className="p-3 bg-[#181c1e] rounded-lg border border-[#33373b] mb-4 text-xs space-y-1">
          <div className="flex justify-between">
            <span className="text-[#a0a5aa]">Weld Integrity:</span>
            <span className="font-bold text-[#6cf8bb]">{inspection.ultrasonicWeldIntegrity}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-[#a0a5aa]">Dimensional Check:</span>
            <span className="font-bold text-[#6cf8bb]">100% Pass</span>
          </div>
        </div>

        <button
          onClick={() => navigate("/dashboard/construction/inspections")}
          className="w-full py-2.5 bg-[#6cf8bb] text-[#181c1e] font-bold text-xs rounded-lg hover:bg-[#6cf8bb]/90 transition-colors flex items-center justify-center gap-1"
        >
          <span className="material-symbols-outlined text-[16px]">fact_check</span>
          Open Quality Sign-Off Record ({inspection.reference})
        </button>
      </div>

      {/* Operational Notice / Boundary Information */}
      <div className="p-3 bg-[#23272b] border border-[#33373b] rounded-xl text-center text-xs text-[#a0a5aa]">
        <span className="material-symbols-outlined text-[20px] text-[#6cf8bb] block mx-auto mb-1">info</span>
        Field Mode presentation interface active. Live site inspections synchronize with TraceCore history.
      </div>
    </div>
  );
};

export default FieldModePage;
