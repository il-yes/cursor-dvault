import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";
import { useRoleContext } from "../hooks/useRoleContext";
import { CreateProjectModal } from "../components/CreateProjectModal";

export const ProjectsPage: React.FC = () => {
  const navigate = useNavigate();
  const { activeRole, roleConfig } = useRoleContext();
  const { project, requirement, delivery, issue, decision, inspection } = SCENARIO_DATA;
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 gap-6">
      {/* Header Banner & Create Project Button */}
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <h1 className="text-2xl font-bold text-[#181c1e]">Projects Directory</h1>
            <span className="px-2.5 py-0.5 bg-[#041627] text-white text-xs font-bold rounded">
              {roleConfig.label} View
            </span>
          </div>
          <p className="text-xs text-[#44474c]">
            Role-aware project entry surface for sovereign workspace.
          </p>
        </div>

        {/* PM Role Action: Create Project */}
        {activeRole === "PM" && (
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="px-4 py-2.5 bg-[#041627] text-white font-bold text-xs rounded-lg hover:bg-[#041627]/90 transition-colors shadow flex items-center gap-1.5"
          >
            <span className="material-symbols-outlined text-[18px]">add_business</span>
            + Create Project
          </button>
        )}
      </div>

      {/* Main Active Project Card: PRJ-001 */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2 text-xs font-bold text-[#44474c]">
            <span className="bg-[#041627] text-white px-2.5 py-1 rounded font-mono">{project.code}</span>
            <span>•</span>
            <span>{project.type}</span>
          </div>
          <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-semibold">
            Execution • 42% Complete
          </span>
        </div>

        <h2 className="text-2xl font-bold text-[#181c1e] mb-2">{project.name}</h2>
        <p className="text-sm text-[#44474c] mb-4">{project.description}</p>
        
        <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
          <span className="material-symbols-outlined text-[18px] text-[#041627]">location_on</span>
          <span>{project.location}</span>
        </div>

        {/* Role-Specific Actions Grid */}
        <div className="pt-4 border-t border-[#f0f3f5]">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-3">
            Available Operations ({roleConfig.label})
          </span>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            {/* General Project Channels Button */}
            <button
              onClick={() => navigate("/dashboard/construction/channels")}
              className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg hover:border-[#041627] text-left transition-colors flex flex-col justify-between"
            >
              <span className="material-symbols-outlined text-[#041627] text-[20px] mb-1">forum</span>
              <div>
                <h4 className="font-bold text-[#181c1e] text-xs">Project Channels</h4>
                <p className="text-[10px] text-[#44474c]">Collaboration Topology</p>
              </div>
            </button>

            {/* PM Specific Actions */}
            {(activeRole === "PM" || activeRole === "ARCHITECT" || activeRole === "ENGINEER") && (
              <button
                onClick={() => navigate("/dashboard/construction/requirements")}
                className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg hover:border-[#041627] text-left transition-colors flex flex-col justify-between"
              >
                <span className="material-symbols-outlined text-[#041627] text-[20px] mb-1">assignment</span>
                <div>
                  <h4 className="font-bold text-[#181c1e] text-xs">Requirements</h4>
                  <p className="text-[10px] text-[#44474c]">{requirement.code}</p>
                </div>
              </button>
            )}

            {/* Supplier / Logistics Actions */}
            {(activeRole === "PM" || activeRole === "SUPPLIER" || activeRole === "LOGISTICS") && (
              <button
                onClick={() => navigate("/dashboard/construction/deliveries")}
                className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg hover:border-[#b76e00] text-left transition-colors flex flex-col justify-between"
              >
                <span className="material-symbols-outlined text-[#b76e00] text-[20px] mb-1">local_shipping</span>
                <div>
                  <h4 className="font-bold text-[#181c1e] text-xs">Deliveries</h4>
                  <p className="text-[10px] text-[#b76e00] font-semibold">{delivery.reference} (Delayed)</p>
                </div>
              </button>
            )}

            {/* Site Superintendent Action */}
            {(activeRole === "PM" || activeRole === "SITE_SUPERINTENDENT") && (
              <button
                onClick={() => navigate("/dashboard/construction/field")}
                className="p-3 bg-[#181c1e] text-white rounded-lg hover:bg-[#181c1e]/90 text-left transition-colors flex flex-col justify-between"
              >
                <span className="material-symbols-outlined text-[#6cf8bb] text-[20px] mb-1">smartphone</span>
                <div>
                  <h4 className="font-bold text-white text-xs">Field Mode</h4>
                  <p className="text-[10px] text-[#6cf8bb]">Site Superintendent</p>
                </div>
              </button>
            )}

            {/* QA Inspector Action */}
            {(activeRole === "PM" || activeRole === "ENGINEER" || activeRole === "QA") && (
              <button
                onClick={() => navigate("/dashboard/construction/inspections")}
                className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg hover:border-[#006c49] text-left transition-colors flex flex-col justify-between"
              >
                <span className="material-symbols-outlined text-[#006c49] text-[20px] mb-1">fact_check</span>
                <div>
                  <h4 className="font-bold text-[#181c1e] text-xs">Inspections</h4>
                  <p className="text-[10px] text-[#006c49] font-semibold">{inspection.reference}</p>
                </div>
              </button>
            )}

            {/* PM / Decision Action */}
            {activeRole === "PM" && (
              <button
                onClick={() => navigate("/dashboard/construction/decisions")}
                className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg hover:border-[#041627] text-left transition-colors flex flex-col justify-between"
              >
                <span className="material-symbols-outlined text-[#041627] text-[20px] mb-1">gavel</span>
                <div>
                  <h4 className="font-bold text-[#181c1e] text-xs">Decisions</h4>
                  <p className="text-[10px] text-[#006c49] font-semibold">{decision.reference}</p>
                </div>
              </button>
            )}

            {/* Provenance Button */}
            <button
              onClick={() => navigate("/dashboard/construction/provenance")}
              className="p-3 bg-[#041627] text-white rounded-lg hover:bg-[#041627]/90 text-left transition-colors flex flex-col justify-between"
            >
              <span className="material-symbols-outlined text-[#6cf8bb] text-[20px] mb-1">account_tree</span>
              <div>
                <h4 className="font-bold text-white text-xs">Why Is It Late?</h4>
                <p className="text-[10px] text-[#6cf8bb]">Provenance Story</p>
              </div>
            </button>
          </div>
        </div>
      </div>

      {/* Project Status Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="p-4 bg-white border border-[#e0e3e5] rounded-xl shadow-sm">
          <span className="text-xs text-[#44474c] font-semibold block mb-1">Budget Allocation</span>
          <span className="text-xl font-bold text-[#181c1e]">{project.budgetSpentPercent}% Spent</span>
          <div className="w-full bg-[#e5e9eb] h-1.5 rounded-full mt-2 overflow-hidden">
            <div className="bg-[#041627] h-full" style={{ width: `${project.budgetSpentPercent}%` }}></div>
          </div>
        </div>

        <div className="p-4 bg-white border border-[#e0e3e5] rounded-xl shadow-sm">
          <span className="text-xs text-[#44474c] font-semibold block mb-1">Schedule Milestone</span>
          <span className="text-xl font-bold text-[#181c1e]">Day {project.scheduleDay}</span>
          <span className="text-xs text-[#b76e00] font-medium block mt-1">3 Issues • 2 Pending Decisions</span>
        </div>

        <div className="p-4 bg-white border border-[#e0e3e5] rounded-xl shadow-sm">
          <span className="text-xs text-[#44474c] font-semibold block mb-1">Compliance &amp; Quality</span>
          <span className="text-xl font-bold text-[#006c49]">98/100</span>
          <span className="text-xs text-[#006c49] font-medium block mt-1">Inspection {inspection.reference} Passed</span>
        </div>
      </div>

      {/* Create Project Modal */}
      <CreateProjectModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
      />
    </div>
  );
};

export default ProjectsPage;
