import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA, ConstructionDecision } from "../data/constructionScenarioAdapter";

export const DecisionPage: React.FC = () => {
  const navigate = useNavigate();
  const [decision, setDecision] = useState<ConstructionDecision>(SCENARIO_DATA.decision);
  const { issue, project, delivery } = SCENARIO_DATA;

  const handleApprove = () => {
    setDecision(prev => ({
      ...prev,
      status: "APPROVED",
      approvedBy: "Alex Rivera (Project Manager)",
      approvedAt: new Date().toISOString().replace("T", " ").substring(0, 16),
      actionTaken: "Alternative Route B approved. Transport rerouted."
    }));
  };

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/issues")} className="hover:underline">
          {issue.reference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Construction Decision</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Construction Decision {decision.reference}</h1>
          <p className="text-sm text-[#44474c]">ProposeDecision → ApproveDecision Construction Lifecycle</p>
        </div>
        <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase ${
          decision.status === "APPROVED" || decision.status === "EXECUTED"
            ? "bg-[#6cf8bb]/20 text-[#00714d]"
            : "bg-[#ffb74d]/20 text-[#b76e00]"
        }`}>
          {decision.status}
        </span>
      </div>

      {/* Decision Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <h2 className="text-xl font-bold text-[#181c1e] mb-2">{decision.title}</h2>
        
        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg mb-4">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">Proposed Action</span>
          <p className="text-sm text-[#181c1e] font-medium">{decision.proposedAction}</p>
        </div>

        <div className="p-4 bg-[#fff3e0] border border-[#ffe0b2] rounded-lg mb-6">
          <span className="text-xs font-bold text-[#b76e00] uppercase tracking-wider block mb-1">Impact &amp; Schedule Summary</span>
          <p className="text-sm text-[#181c1e]">{decision.impactSummary}</p>
        </div>

        {/* Stakeholder Sign-Off Info */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-[#f0f3f5] mb-6 text-xs text-[#44474c]">
          <div>
            <span className="block text-[#44474c]">Proposed By:</span>
            <span className="font-bold text-[#181c1e] text-sm">{decision.proposedBy}</span>
          </div>
          <div>
            <span className="block text-[#44474c]">Approved By:</span>
            <span className="font-bold text-[#006c49] text-sm">{decision.approvedBy || "Pending Sign-off"}</span>
            {decision.approvedAt && <span className="block text-[11px] text-[#44474c] mt-0.5">{decision.approvedAt}</span>}
          </div>
        </div>

        {/* Action Button: Approve Decision */}
        {decision.status === "PROPOSED" ? (
          <button
            onClick={handleApprove}
            className="w-full py-3 bg-[#006c49] text-white font-bold text-sm rounded-lg hover:bg-[#006c49]/90 transition-colors shadow-md flex items-center justify-center gap-2"
          >
            <span className="material-symbols-outlined text-[20px]">check_circle</span>
            Approve Decision ({decision.reference})
          </button>
        ) : (
          <div className="p-4 bg-[#e8f5e9] border border-[#c8e6c9] rounded-lg flex items-center justify-between">
            <div className="flex items-center gap-2 text-[#006c49] text-sm font-bold">
              <span className="material-symbols-outlined">task_alt</span>
              <span>Decision Executed: {decision.actionTaken}</span>
            </div>
            <button
              onClick={() => navigate("/dashboard/construction/deliveries")}
              className="px-4 py-2 bg-[#041627] text-white font-semibold text-xs rounded-lg hover:bg-[#041627]/90 transition-colors"
            >
              Return to Delivery ({delivery.reference}) →
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default DecisionPage;
