import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getActivityFeed, getProjects, ProjectData, getIssue, getDelivery, getDecision } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

/** Mirrors TERMINAL_ISSUE_STATUSES in projects.mock.ts: a resolved/closed issue is not open. */
const TERMINAL_ISSUE_STATUSES = ["resolved", "closed"];

export const ConstructionDashboard: React.FC = () => {
  const navigate = useNavigate();

  const activityLog = getActivityFeed();
  const issue = getIssue();
  const delivery = getDelivery();
  const decision = getDecision();

  const [projects, setProjects] = useState<ProjectData[]>([]);

  useEffect(() => {
    getProjects()
      .then(setProjects)
      .catch((err) => console.error("Dashboard failed to load projects:", err));
  }, []);

  const activeProjects = projects.filter((p) => p.status === "active");
  const openIssue = !TERMINAL_ISSUE_STATUSES.includes(issue.status);
  const pendingDecision = decision.status !== "approved" && decision.status !== "rejected";

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 gap-6 font-[Inter]">
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
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="px-4 py-2 bg-[#6cf8bb] text-[#041627] font-bold text-xs rounded-lg hover:bg-[#6cf8bb]/90 transition-colors shadow cursor-pointer"
          >
            View Projects Directory →
          </button>
        </div>
      </div>

      {/* Vault Metric Summary Cards */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Projects</span>
          <div className="text-2xl font-bold text-[#041627] mt-2">{projects.length}</div>
          <span className="text-[10px] text-[#006c49] font-medium mt-1">
            {activeProjects.length} Active{activeProjects[0] ? ` (${activeProjects[0].code})` : ""}
          </span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Open Issues</span>
          <div className="text-2xl font-bold text-[#d32f2f] mt-2">{openIssue ? 1 : 0}</div>
          <span className="text-[10px] text-[#d32f2f] font-medium mt-1">
            {openIssue
              ? `1 ${issue.severity} (${issue.issueReference})`
              : `No open issues (${issue.issueReference} ${issue.status})`}
          </span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Pending Actions</span>
          <div className="text-2xl font-bold text-[#b76e00] mt-2">{pendingDecision ? 1 : 0}</div>
          <span className="text-[10px] text-[#b76e00] font-medium mt-1">
            {pendingDecision
              ? "Approvals & Sign-offs"
              : `No approvals pending (${decision.decisionReference} ${decision.status})`}
          </span>
        </div>

        <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Deliveries</span>
          <div className="text-2xl font-bold text-[#041627] mt-2">1</div>
          <span className="text-[10px] text-[#b76e00] font-medium mt-1">
            1 {delivery.status} ({delivery.deliveryReference})
          </span>
        </div>

        <div className="col-span-2 md:col-span-1 bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col justify-between">
          <span className="text-xs text-[#44474c] font-semibold uppercase tracking-wider">Recent Events</span>
          <div className="text-2xl font-bold text-[#006c49] mt-2">{activityLog.length}</div>
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
          <span className="text-xs font-mono text-[#44474c]">Cursor-Ordered</span>
        </div>

        <div className="flex flex-col gap-0 relative">
          <div className="absolute left-[35px] top-3 bottom-3 w-0.5 bg-[#e0e3e5] rounded z-0" />
          
          {activityLog.map((act) => (
            <div key={act.id} className="flex items-start gap-4 py-3 relative z-10 hover:bg-[#f7fafc] px-2 rounded-lg transition-colors">
              <span className="text-xs font-mono text-[#44474c] min-w-[45px] pt-1">#{act.cursor}</span>
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
                  type="button"
                  onClick={() => navigate(act.link)}
                  className="text-xs font-semibold text-[#041627] hover:underline self-start md:self-auto cursor-pointer"
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
