import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const ProjectHistoryPage: React.FC = () => {
  const navigate = useNavigate();
  const { project, traceMilestones } = SCENARIO_DATA;

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Project History</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">TraceCore Historical Milestones</h1>
          <p className="text-sm text-[#44474c]">{project.name} ({project.code}) — Chronological state record</p>
        </div>
        <span className="px-3 py-1 bg-[#041627] text-white rounded-full text-xs font-bold font-mono">
          TRACECORE AUDIT
        </span>
      </div>

      {/* Milestones List */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <div className="flex flex-col gap-4">
          {traceMilestones.map((m, idx) => (
            <div key={idx} className="flex items-start gap-4 pb-4 border-b border-[#f0f3f5] last:border-b-0 last:pb-0">
              <div className="text-xs font-mono text-[#44474c] whitespace-nowrap min-w-[120px] pt-1">
                {m.timestamp}
              </div>
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                    m.layer === "Construction" ? "bg-[#041627] text-white" :
                    m.layer === "C3 Collaboration" ? "bg-[#e8f5e9] text-[#006c49]" :
                    "bg-[#fff3e0] text-[#b76e00]"
                  }`}>
                    {m.layer}
                  </span>
                  <span className="text-xs font-mono text-[#44474c]">Entity: {m.entityId}</span>
                </div>
                <h4 className="text-sm font-bold text-[#181c1e]">{m.event}</h4>
                <p className="text-xs text-[#44474c]">Actor: {m.actor}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProjectHistoryPage;
