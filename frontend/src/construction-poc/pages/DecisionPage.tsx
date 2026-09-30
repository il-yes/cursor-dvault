import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { getDecision, getIssue, getProjectRecord, getDelivery } from "../data";
import type { ConstructionDecisionRecord } from "../data/constructionScenarioAdapter";

export const DecisionPage: React.FC = () => {
  const navigate = useNavigate();
  const [decision, setDecision] = useState<ConstructionDecisionRecord>(getDecision());
  const issue = getIssue();
  const project = getProjectRecord();
  const delivery = getDelivery();

  // Session-local approval of the canonical record. Values come from the
  // canonical decision; only the transition timestamp is generated.
  const handleApprove = () => {
    setDecision((prev) => ({
      ...prev,
      status: "approved",
      decidedBy: prev.requestedBy,
      decisionDate: new Date().toISOString(),
      consequence: prev.consequence
    }));
  };

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.projectReference}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/issues")} className="hover:underline">
          {issue.issueReference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Construction Decision</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Construction Decision {decision.decisionReference}</h1>
          <p className="text-sm text-[#44474c]">
            {decision.type} • ProposeDecision → ApproveDecision Construction Lifecycle
          </p>
        </div>
        <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase ${
          decision.status === "approved"
            ? "bg-[#6cf8bb]/20 text-[#00714d]"
            : "bg-[#ffb74d]/20 text-[#b76e00]"
        }`}>
          {decision.status}
        </span>
      </div>

      {/* Decision Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <h2 className="text-xl font-bold text-[#181c1e] mb-2">{decision.subject}</h2>

        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg mb-4">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">Context</span>
          <p className="text-sm text-[#181c1e] font-medium">{decision.context}</p>
        </div>

        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg mb-4">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
            Technical Assessment
          </span>
          <p className="text-sm text-[#181c1e]">{decision.technicalAssessment}</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
              Risks Identified
            </span>
            <ul className="space-y-1">
              {decision.risksIdentified.map((risk) => (
                <li key={risk} className="text-sm text-[#181c1e]">
                  • {risk}
                </li>
              ))}
            </ul>
          </div>
          <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
              Options Considered
            </span>
            <ul className="space-y-1">
              {decision.optionsConsidered.map((option) => (
                <li key={option} className="text-sm text-[#181c1e]">
                  • {option}
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg mb-4">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
            Participants Consulted
          </span>
          <p className="text-sm text-[#181c1e]">{decision.participantsConsulted.join(", ")}</p>
        </div>

        <div className="p-4 bg-[#fff3e0] border border-[#ffe0b2] rounded-lg mb-6">
          <span className="text-xs font-bold text-[#b76e00] uppercase tracking-wider block mb-1">
            Impact &amp; Schedule Summary
          </span>
          <p className="text-sm text-[#181c1e]">{decision.consequence}</p>
        </div>

        {/* Stakeholder Sign-Off Info */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-[#f0f3f5] mb-6 text-xs text-[#44474c]">
          <div>
            <span className="block text-[#44474c]">Requested By:</span>
            <span className="font-bold text-[#181c1e] text-sm">{decision.requestedBy}</span>
            <span className="block text-[11px] text-[#44474c] mt-0.5">{decision.requestDate}</span>
          </div>
          <div>
            <span className="block text-[#44474c]">Decided By:</span>
            <span className="font-bold text-[#006c49] text-sm">{decision.decidedBy || "Pending Sign-off"}</span>
            {decision.decisionDate && (
              <span className="block text-[11px] text-[#44474c] mt-0.5">{decision.decisionDate}</span>
            )}
          </div>
        </div>

        {/* Action Button: Approve Decision */}
        {decision.status === "proposed" ? (
          <button
            onClick={handleApprove}
            className="w-full py-3 bg-[#006c49] text-white font-bold text-sm rounded-lg hover:bg-[#006c49]/90 transition-colors shadow-md flex items-center justify-center gap-2"
          >
            <span className="material-symbols-outlined text-[20px]">check_circle</span>
            Approve Decision ({decision.decisionReference})
          </button>
        ) : (
          <div className="p-4 bg-[#e8f5e9] border border-[#c8e6c9] rounded-lg flex items-center justify-between">
            <div className="flex items-center gap-2 text-[#006c49] text-sm font-bold">
              <span className="material-symbols-outlined">task_alt</span>
              <span>Decision Approved: {decision.decision}</span>
            </div>
            <button
              onClick={() => navigate("/dashboard/construction/deliveries")}
              className="px-4 py-2 bg-[#041627] text-white font-semibold text-xs rounded-lg hover:bg-[#041627]/90 transition-colors"
            >
              Return to Delivery ({delivery.deliveryReference}) →
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default DecisionPage;
