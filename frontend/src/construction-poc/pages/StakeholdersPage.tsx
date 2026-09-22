import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const StakeholdersPage: React.FC = () => {
  const navigate = useNavigate();
  const { stakeholders, project } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Navigation Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          Projects
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">{project.code}</span>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Stakeholders</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Project Stakeholders</h1>
          <p className="text-sm text-[#44474c]">{project.name} ({project.code}) — Multi-party coordination team</p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-semibold">
          {stakeholders.length} Active Participants
        </span>
      </div>

      {/* Stakeholders Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {stakeholders.map((stk) => (
          <div key={stk.id} className="bg-white border border-[#e0e3e5] rounded-xl p-5 shadow-sm flex flex-col gap-3">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-[#041627] text-white flex items-center justify-center font-bold text-sm">
                  {stk.name.split(" ").map(n => n[0]).join("")}
                </div>
                <div>
                  <h3 className="font-semibold text-[#181c1e] text-base">{stk.name}</h3>
                  <p className="text-xs text-[#44474c]">{stk.role}</p>
                </div>
              </div>
              <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-[#006c49]/10 text-[#006c49] text-[11px] font-semibold">
                <span className="material-symbols-outlined text-[14px]">verified</span>
                {stk.trustStatus}
              </span>
            </div>

            <div className="pt-2 border-t border-[#f0f3f5] flex flex-col gap-1 text-xs text-[#44474c]">
              <div className="flex items-center justify-between">
                <span className="font-medium text-[#181c1e]">Organization:</span>
                <span>{stk.organization}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="font-medium text-[#181c1e]">Email:</span>
                <span className="text-[#041627] underline">{stk.email}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="font-medium text-[#181c1e]">Phone:</span>
                <span>{stk.phone}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default StakeholdersPage;
