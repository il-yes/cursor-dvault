import React from "react";
import { useNavigate } from "react-router-dom";
import {
  getProjectRecord,
  getDelivery,
  getIssue,
  getEvidenceDocument,
  getDecision,
  getTransportDelay,
  getInspection,
  getMaterial,
  getSite,
} from "../data";

export const ProvenanceWhyPage: React.FC = () => {
  const navigate = useNavigate();
  const project = getProjectRecord();
  const delivery = getDelivery();
  const issue = getIssue();
  const evidenceDoc = getEvidenceDocument();
  const decision = getDecision();
  const transport = getTransportDelay();
  const inspection = getInspection();
  const material = getMaterial();
  const site = getSite();

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.projectReference}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/deliveries")} className="hover:underline">
          {delivery.deliveryReference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Provenance / Why?</span>
      </div>

      {/* Main Question Banner */}
      <div className="bg-[#041627] text-white rounded-xl p-6 shadow-md mb-8">
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs text-[#6cf8bb] font-mono uppercase tracking-wider font-semibold">
            PROVENANCE RECONSTRUCTION ENGINE
          </span>
          <span className="px-2.5 py-0.5 bg-[#6cf8bb]/20 text-[#6cf8bb] rounded text-[11px] font-bold">
            VERIFIED HISTORICAL CHAIN
          </span>
        </div>
        <h1 className="text-2xl font-bold text-white mb-2">Why is {delivery.deliveryReference} late?</h1>
        <p className="text-sm text-[#b7c8de]">
          Reconstructed operational causal chain for {material.name} on {project.projectName}.
        </p>
      </div>

      {/* Causal Chain Timeline (Step-by-Step Story) */}
      <h2 className="text-lg font-bold text-[#181c1e] mb-4 flex items-center gap-2">
        <span className="material-symbols-outlined text-[#041627]">account_tree</span>
        Causal Reconstruction Story
      </h2>

      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-8 relative">
        <div className="flex flex-col gap-6 relative z-10">
          
          {/* Step 1: Cause */}
          <div className="flex items-start gap-4">
            <div className="w-8 h-8 rounded-full bg-[#d32f2f] text-white flex items-center justify-center font-bold text-xs flex-shrink-0 mt-1">
              1
            </div>
            <div className="p-4 bg-[#fff5f5] border border-[#ffcdd2] rounded-lg flex-1">
              <span className="text-xs font-bold text-[#d32f2f] uppercase tracking-wider block mb-1">Root Cause</span>
              <h3 className="font-bold text-[#181c1e] text-base mb-1">{transport.delayReason}</h3>
              <p className="text-xs text-[#44474c]">
                Reported on route {transport.routeReference} by transport {transport.transportReference}.{" "}
                {transport.constraints.join("; ")}.
              </p>
            </div>
          </div>

          {/* Step 2: Issue */}
          <div className="flex items-start gap-4">
            <div className="w-8 h-8 rounded-full bg-[#d32f2f] text-white flex items-center justify-center font-bold text-xs flex-shrink-0 mt-1">
              2
            </div>
            <div className="p-4 bg-white border border-[#e0e3e5] rounded-lg flex-1">
              <div className="flex justify-between items-center mb-1">
                <span className="text-xs font-bold text-[#041627] uppercase tracking-wider">Construction Issue</span>
                <span className="text-xs font-mono text-[#44474c]">{issue.issueReference}</span>
              </div>
              <h3 className="font-bold text-[#181c1e] text-base mb-1">{issue.title}</h3>
              <p className="text-xs text-[#44474c]">{issue.description}</p>
            </div>
          </div>

          {/* Step 3: Evidence */}
          <div className="flex items-start gap-4">
            <div className="w-8 h-8 rounded-full bg-[#041627] text-white flex items-center justify-center font-bold text-xs flex-shrink-0 mt-1">
              3
            </div>
            <div className="p-4 bg-white border border-[#e0e3e5] rounded-lg flex-1">
              <div className="flex justify-between items-center mb-1">
                <span className="text-xs font-bold text-[#041627] uppercase tracking-wider">Evidence Document</span>
                <span className="text-xs font-mono text-[#006c49]">{evidenceDoc.documentReference}</span>
              </div>
              <h3 className="font-bold text-[#181c1e] text-base mb-1">{evidenceDoc.name}</h3>
              <p className="text-xs text-[#44474c] font-mono">
                {evidenceDoc.category} v{evidenceDoc.version} • CID: {evidenceDoc.cid} • Hash:{" "}
                {evidenceDoc.contentHash.substring(0, 24)}...
              </p>
            </div>
          </div>

          {/* Step 4: Decision */}
          <div className="flex items-start gap-4">
            <div className="w-8 h-8 rounded-full bg-[#006c49] text-white flex items-center justify-center font-bold text-xs flex-shrink-0 mt-1">
              4
            </div>
            <div className="p-4 bg-white border border-[#e0e3e5] rounded-lg flex-1">
              <div className="flex justify-between items-center mb-1">
                <span className="text-xs font-bold text-[#006c49] uppercase tracking-wider">Approved Decision</span>
                <span className="text-xs font-mono text-[#041627]">{decision.decisionReference}</span>
              </div>
              <h3 className="font-bold text-[#181c1e] text-base mb-1">{decision.subject}</h3>
              <p className="text-xs text-[#44474c]">{decision.decision}</p>
            </div>
          </div>

          {/* Step 5: Delivery & Inspection */}
          <div className="flex items-start gap-4">
            <div className="w-8 h-8 rounded-full bg-[#006c49] text-white flex items-center justify-center font-bold text-xs flex-shrink-0 mt-1">
              5
            </div>
            <div className="p-4 bg-[#e8f5e9] border border-[#c8e6c9] rounded-lg flex-1">
              <span className="text-xs font-bold text-[#006c49] uppercase tracking-wider block mb-1">Delivery Received &amp; Material Accepted</span>
              <h3 className="font-bold text-[#181c1e] text-base mb-1">
                Delivery {delivery.deliveryReference} Received at {site.name}
              </h3>
              <p className="text-xs text-[#44474c]">
                Inspection {inspection.inspectionReference} ({inspection.type}) completed {inspection.completedAt}. Result:{" "}
                {inspection.result}.
              </p>
            </div>
          </div>

        </div>
      </div>

      {/* Expose 3 Historical System Layers */}
      <h2 className="text-lg font-bold text-[#181c1e] mb-4 flex items-center gap-2">
        <span className="material-symbols-outlined text-[#041627]">layers</span>
        System Layer Architecture
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div className="p-4 bg-white border border-[#e0e3e5] rounded-xl shadow-sm">
          <span className="w-3 h-3 rounded-full bg-[#041627] inline-block mr-2"></span>
          <span className="text-xs font-bold text-[#041627] uppercase tracking-wider">Construction Layer</span>
          <p className="text-xs text-[#44474c] mt-2 leading-relaxed">
            Represents physical business state: Projects, Requirements, Materials, Deliveries, Decisions, and Inspections.
          </p>
        </div>

        <div className="p-4 bg-white border border-[#e0e3e5] rounded-xl shadow-sm">
          <span className="w-3 h-3 rounded-full bg-[#006c49] inline-block mr-2"></span>
          <span className="text-xs font-bold text-[#006c49] uppercase tracking-wider">C3 Collaboration</span>
          <p className="text-xs text-[#44474c] mt-2 leading-relaxed">
            Represents evidence sharing, thread events, multi-party messaging, and encrypted attachments.
          </p>
        </div>

        <div className="p-4 bg-white border border-[#e0e3e5] rounded-xl shadow-sm">
          <span className="w-3 h-3 rounded-full bg-[#b76e00] inline-block mr-2"></span>
          <span className="text-xs font-bold text-[#b76e00] uppercase tracking-wider">TraceCore History</span>
          <p className="text-xs text-[#44474c] mt-2 leading-relaxed">
            Represents chronological milestone records and historical audit log of state changes.
          </p>
        </div>
      </div>

      {/* Button to Trace History */}
      <div className="flex justify-end">
        <button
          onClick={() => navigate("/dashboard/construction/history")}
          className="px-5 py-2.5 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
        >
          View Full Trace History Timeline →
        </button>
      </div>
    </div>
  );
};

export default ProvenanceWhyPage;
