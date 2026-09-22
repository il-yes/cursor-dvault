import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const ProjectsPage: React.FC = () => {
  const navigate = useNavigate();
  const { project, requirement, delivery, issue, decision, inspection } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Project Header Banner */}
      <div className="p-6 bg-white border border-[#e0e3e5] rounded-xl shadow-sm mb-6">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2 text-xs text-[#44474c] font-semibold tracking-wider">
            <span className="bg-[#041627] text-white px-2 py-0.5 rounded font-mono">{project.code}</span>
            <span className="w-1 h-1 rounded-full bg-[#c4c6cd]"></span>
            <span>{project.type}</span>
          </div>
          <span className="inline-flex items-center px-2.5 py-1 rounded bg-[#6cf8bb]/20 text-[#00714d] font-semibold text-xs">
            {project.status}
          </span>
        </div>

        <h1 className="text-2xl font-bold text-[#181c1e] mb-2">{project.name}</h1>
        <p className="text-sm text-[#44474c] mb-4">{project.description}</p>
        
        <div className="flex items-center gap-2 text-sm text-[#44474c]">
          <span className="material-symbols-outlined text-[18px]">location_on</span>
          <span>{project.location}</span>
        </div>
      </div>

      {/* Quick Navigation Cards to Project Entity Views */}
      <h2 className="text-lg font-bold text-[#181c1e] mb-3">Project Management Views</h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
        <button
          onClick={() => navigate("/dashboard/construction/stakeholders")}
          className="p-4 bg-white border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#041627] text-[24px] mb-2">groups</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Stakeholders</h3>
            <p className="text-[11px] text-[#44474c]">4 Active Partners</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/requirements")}
          className="p-4 bg-white border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#041627] text-[24px] mb-2">assignment</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Requirements</h3>
            <p className="text-[11px] text-[#44474c]">{requirement.code}</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/deliveries")}
          className="p-4 bg-white border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#b76e00] text-[24px] mb-2">local_shipping</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Deliveries</h3>
            <p className="text-[11px] text-[#b76e00] font-semibold">{delivery.reference} (Delayed)</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/issues")}
          className="p-4 bg-white border border-[#ffcdd2] rounded-xl hover:border-[#d32f2f] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#d32f2f] text-[24px] mb-2">warning</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Issues</h3>
            <p className="text-[11px] text-[#d32f2f] font-semibold">{issue.reference}</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/decisions")}
          className="p-4 bg-white border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#041627] text-[24px] mb-2">gavel</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Decisions</h3>
            <p className="text-[11px] text-[#006c49] font-semibold">{decision.reference} ({decision.status})</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/inspections")}
          className="p-4 bg-white border border-[#e0e3e5] rounded-xl hover:border-[#006c49] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#006c49] text-[24px] mb-2">fact_check</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Inspections</h3>
            <p className="text-[11px] text-[#006c49] font-semibold">{inspection.reference} ({inspection.status})</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/provenance")}
          className="p-4 bg-[#041627] text-white rounded-xl hover:bg-[#041627]/90 text-left transition-all shadow-md flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#6cf8bb] text-[24px] mb-2">account_tree</span>
          <div>
            <h3 className="font-bold text-white text-sm">Why Is It Late?</h3>
            <p className="text-[11px] text-[#6cf8bb]">Provenance Story</p>
          </div>
        </button>

        <button
          onClick={() => navigate("/dashboard/construction/history")}
          className="p-4 bg-white border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-all shadow-sm flex flex-col justify-between"
        >
          <span className="material-symbols-outlined text-[#041627] text-[24px] mb-2">history</span>
          <div>
            <h3 className="font-bold text-[#181c1e] text-sm">Trace History</h3>
            <p className="text-[11px] text-[#44474c]">Milestones</p>
          </div>
        </button>
      </div>

      {/* Metrics Summary Grid */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 flex flex-col gap-2">
          <div className="flex items-center gap-2 text-[#44474c]">
            <span className="material-symbols-outlined text-[20px]">account_balance_wallet</span>
            <span className="text-[11px] uppercase tracking-wider font-semibold">Budget</span>
          </div>
          <div className="flex items-baseline gap-1">
            <span className="text-2xl font-semibold text-[#181c1e]">{project.budgetSpentPercent}%</span>
            <span className="text-sm text-[#44474c]">spent</span>
          </div>
          <div className="w-full h-1.5 bg-[#e5e9eb] rounded-full overflow-hidden mt-1">
            <div className="h-full bg-[#041627] rounded-full" style={{ width: `${project.budgetSpentPercent}%` }}></div>
          </div>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 flex flex-col gap-2">
          <div className="flex items-center gap-2 text-[#44474c]">
            <span className="material-symbols-outlined text-[20px]">calendar_today</span>
            <span className="text-[11px] uppercase tracking-wider font-semibold">Schedule</span>
          </div>
          <div className="flex items-baseline gap-1">
            <span className="text-2xl font-semibold text-[#181c1e]">Day {project.scheduleDay}</span>
          </div>
          <div className="text-sm text-[#b76e00] flex items-center gap-1 mt-1 font-medium">
            <span className="material-symbols-outlined text-[16px]">warning</span>
            <span>{project.scheduleStatus} (DEL-1042 Reroute)</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProjectsPage;
