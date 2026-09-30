import React from "react";
import { useNavigate } from "react-router-dom";
import { getEvidenceDocument, getIssue, getDecision, getProjectRecord } from "../data";

export const ConstructionIssuePage: React.FC = () => {
  const navigate = useNavigate();
  const issue = getIssue();
  const evidenceDoc = getEvidenceDocument();
  const decision = getDecision();
  const project = getProjectRecord();

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.projectReference}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/transport")} className="hover:underline">
          Transport
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Construction Issue</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Construction Issue {issue.issueReference}</h1>
          <p className="text-sm text-[#44474c]">Reported at {issue.reportedAt}</p>
        </div>
        <span className="px-3 py-1 bg-[#ffebee] text-[#d32f2f] border border-[#ffcdd2] rounded-full text-xs font-bold uppercase">
          {issue.severity} severity · {issue.type}
        </span>
      </div>

      {/* Issue Details Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <h2 className="text-lg font-bold text-[#181c1e] mb-2">{issue.title}</h2>
        <p className="text-sm text-[#44474c] mb-6 leading-relaxed">{issue.description}</p>

        {/* Evidence Document Attachment Box */}
        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-xl mb-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[18px] text-[#041627]">attach_file</span>
              Attached Evidence Document
            </span>
            <span className="text-[11px] bg-[#e0e3e5] text-[#041627] px-2 py-0.5 rounded font-mono">
              {evidenceDoc.documentReference}
            </span>
          </div>

          <div className="bg-white p-3 rounded-lg border border-[#e0e3e5] flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded bg-[#d32f2f]/10 text-[#d32f2f] flex items-center justify-center font-bold text-xs">
                PDF
              </div>
              <div>
                <h4 className="text-sm font-bold text-[#181c1e]">{evidenceDoc.name}</h4>
                <p className="text-xs text-[#44474c] font-mono mt-0.5">
                  {evidenceDoc.category} • v{evidenceDoc.version}
                </p>
              </div>
            </div>
            <div className="text-right text-xs">
              <span className="text-[#006c49] font-medium block">Issued by: {evidenceDoc.issuedBy}</span>
              <span className="text-[#44474c] font-mono text-[10px] block truncate max-w-[180px]">
                {evidenceDoc.contentHash}
              </span>
            </div>
          </div>
        </div>

        {/* Action Link to Decision */}
        <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
          <div className="text-xs text-[#44474c]">
            Proposed Resolution Decision:{" "}
            <span className="font-semibold text-[#181c1e]">{decision.decisionReference}</span>
          </div>
          <button
            onClick={() => navigate("/dashboard/construction/decisions")}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
          >
            Review Decision ({decision.decisionReference}) →
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConstructionIssuePage;
