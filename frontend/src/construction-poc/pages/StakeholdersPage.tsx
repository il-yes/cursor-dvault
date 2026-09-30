import React from "react";
import { useNavigate } from "react-router-dom";
import { getParticipants, getProjectRecord } from "../data";

export const StakeholdersPage: React.FC = () => {
  const navigate = useNavigate();
  // The canonical domain has no ConstructionStakeholder aggregate. These are the
  // three real vault identities the scenario registers, with their trust roles.
  const participants = getParticipants();
  const project = getProjectRecord();

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Navigation Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          Projects
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">{project.projectReference}</span>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Stakeholders</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Project Stakeholders</h1>
          <p className="text-sm text-[#44474c]">
            {project.projectName} ({project.projectReference}) — Multi-party coordination team
          </p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-semibold">
          {participants.length} Active Participants
        </span>
      </div>

      {/* Stakeholders Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {participants.map((p) => (
          <div key={p.vaultId} className="bg-white border border-[#e0e3e5] rounded-xl p-5 shadow-sm flex flex-col gap-3">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-[#041627] text-white flex items-center justify-center font-bold text-sm">
                  {p.vaultId.split(/[-_]/).filter(Boolean).map(n => n[0]).join("")}
                </div>
                <div>
                  <h3 className="font-semibold text-[#181c1e] text-base">{p.vaultId}</h3>
                  <p className="text-xs text-[#44474c]">Trust role: {p.scenarioRole}</p>
                </div>
              </div>
              <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-[#006c49]/10 text-[#006c49] text-[11px] font-semibold">
                <span className="material-symbols-outlined text-[14px]">verified</span>
                {p.scenarioRole}
              </span>
            </div>

            <div className="pt-2 border-t border-[#f0f3f5] flex flex-col gap-1 text-xs text-[#44474c]">
              <div className="flex items-center justify-between">
                <span className="font-medium text-[#181c1e]">Organization:</span>
                <span>{p.organizationId}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="font-medium text-[#181c1e]">Vault Address:</span>
                <span className="text-[#041627] underline">{p.vaultAddress}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default StakeholdersPage;
