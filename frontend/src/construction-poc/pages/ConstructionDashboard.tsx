import React from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";

export const ConstructionDashboard: React.FC = () => {
  const navigate = useNavigate();
  const { project, requirement, delivery, issue, decision, inspection, supplier } = SCENARIO_DATA;

  // Sovereign Vault Home Activity Stream (Git-style chronological log)
  const activityLog = [
    {
      time: "08:42",
      type: "PROJECT_CREATED",
      badge: "PRJ-001",
      title: "Project created: Metro Line 4 Expansion",
      actor: "Manuel (Project Owner)",
      link: "/dashboard/construction/projects"
    },
    {
      time: "08:51",
      type: "STAKEHOLDER_INVITED",
      badge: "STK-001",
      title: "Stakeholder invited: Lead Architect → Project",
      actor: "Manuel (Project Owner)",
      link: "/dashboard/construction/stakeholders"
    },
    {
      time: "09:14",
      type: "REQUIREMENT_CREATED",
      badge: requirement.code,
      title: `Requirement created: ${requirement.materialName}`,
      actor: "Alex Rivera (Project Manager)",
      link: "/dashboard/construction/requirements"
    },
    {
      time: "09:42",
      type: "SUPPLIER_INVITED",
      badge: supplier.id,
      title: `Supplier invitation sent: ${supplier.name} (${supplier.id})`,
      actor: "Alex Rivera (Project Manager)",
      link: "/dashboard/construction/offers"
    },
    {
      time: "10:05",
      type: "DELIVERY_DELAYED",
      badge: delivery.reference,
      title: `Delivery delayed: ${delivery.reference} (ETA moved to 2026-08-16 07:30)`,
      actor: "David Chen (Logistics Coordinator)",
      link: "/dashboard/construction/deliveries"
    },
    {
      time: "10:07",
      type: "ISSUE_REPORTED",
      badge: issue.reference,
      title: `Issue reported: ${issue.title} (${issue.severity} Severity)`,
      actor: "Mark Vance (Freight Operator)",
      link: "/dashboard/construction/issues"
    },
    {
      time: "10:21",
      type: "DECISION_APPROVED",
      badge: decision.reference,
      title: `Decision approved: ${decision.title}`,
      actor: "Alex Rivera (Project Manager)",
      link: "/dashboard/construction/decisions"
    },
    {
      time: "10:48",
      type: "INSPECTION_COMPLETED",
      badge: inspection.reference,
      title: `Inspection completed: ${inspection.title} (${inspection.status})`,
      actor: "Sarah Jenkins, PE (Site Superintendent)",
      link: "/dashboard/construction/inspections"
    }
  ];

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 gap-6">
      {/* Vault Front Door Header */}
      <div className="bg-[#041627] text-white rounded-xl p-6 shadow-md flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <span className="text-xs font-mono text-[#6cf8bb] uppercase tracking-wider font-semibold block mb-1">
            ANKHORA SOVEREIGN VAULT HOME
          </span>
          <h1 className="text-2xl font-bold text-white mb-1">Vault Activity &amp; Overview</h1>
          <p className="text-xs text-[#b7c8de]">
            Global sovereign activity, project metrics, and multi-party coordination log.
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => navigate("/dashboard/construction/projects")}
            className="px-4 py-2 bg-[#6cf8bb] text-[#041627] font-bold text-xs rounded-lg hover:bg-[#6cf8bb]/90 transition-colors shadow"
          >
            View Projects (PRJ-001) →
          </button>
        </div>
      </div>

      {/* Vault Metric Summary Cards */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Projects</span>
          <div className="text-2xl font-bold text-[#041627] mt-2">4</div>
          <span className="text-[10px] text-[#006c49] font-medium mt-1">1 Active ({project.code})</span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Open Issues</span>
          <div className="text-2xl font-bold text-[#d32f2f] mt-2">3</div>
          <span className="text-[10px] text-[#d32f2f] font-medium mt-1">1 Critical ({issue.reference})</span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Pending Actions</span>
          <div className="text-2xl font-bold text-[#b76e00] mt-2">7</div>
          <span className="text-[10px] text-[#b76e00] font-medium mt-1">Approvals &amp; Sign-offs</span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Deliveries</span>
          <div className="text-2xl font-bold text-[#041627] mt-2">12</div>
          <span className="text-[10px] text-[#b76e00] font-medium mt-1">1 Delayed ({delivery.reference})</span>
        </div>

        <div className="col-span-2 md:col-span-1 bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Recent Events</span>
          <div className="text-2xl font-bold text-[#006c49] mt-2">24</div>
          <span className="text-[10px] text-[#006c49] font-medium mt-1">Verified Audit Log</span>
        </div>
      </div>

      {/* Git-Style Sovereign Activity Stream */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm">
        <div className="flex items-center justify-between mb-4 pb-3 border-b border-[#f0f3f5]">
          <div className="flex items-center gap-2">
            <span className="material-symbols-outlined text-[#041627] text-[20px]">history</span>
            <h2 className="text-base font-bold text-[#181c1e]">Sovereign Workspace Activity Feed</h2>
          </div>
          <span className="text-xs font-mono text-[#44474c]">Today</span>
        </div>

        <div className="flex flex-col gap-0 relative">
          <div className="absolute left-[35px] top-3 bottom-3 w-0.5 bg-[#e0e3e5] rounded z-0" />
          
          {activityLog.map((act, index) => (
            <div key={index} className="flex items-start gap-4 py-3 relative z-10 hover:bg-[#f7fafc] px-2 rounded-lg transition-colors">
              <span className="text-xs font-mono text-[#44474c] min-w-[45px] pt-1">{act.time}</span>
              <div className="w-5 h-5 rounded-full bg-[#041627] text-white flex items-center justify-center text-[10px] font-bold mt-0.5 flex-shrink-0">
                •
              </div>
              <div className="flex-1 flex flex-col md:flex-row md:items-center justify-between gap-1">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-bold text-[#181c1e]">{act.title}</span>
                    <span className="px-1.5 py-0.5 bg-[#f0f3f5] text-[#041627] text-[10px] font-mono rounded font-semibold">
                      {act.badge}
                    </span>
                  </div>
                  <span className="text-xs text-[#44474c] block">by {act.actor}</span>
                </div>
                <button
                  onClick={() => navigate(act.link)}
                  className="text-xs font-semibold text-[#041627] hover:underline self-start md:self-auto"
                >
                  View Details →
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ConstructionDashboard;
