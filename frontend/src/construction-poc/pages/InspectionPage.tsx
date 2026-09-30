import React from "react";
import { useNavigate } from "react-router-dom";
import { getCertificateDocument, getInspection, getDelivery, getProjectRecord, getRequirement } from "../data";
import { actorLabel } from "../data/scenarioMappers";

export const InspectionPage: React.FC = () => {
  const navigate = useNavigate();
  const inspection = getInspection();
  const delivery = getDelivery();
  const project = getProjectRecord();
  const requirement = getRequirement();
  const certDoc = getCertificateDocument();

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
        <span className="font-semibold text-[#041627]">Inspection</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">Quality Inspection {inspection.inspectionReference}</h1>
          <p className="text-sm text-[#44474c]">
            Inspector: {actorLabel(inspection.inspectorId)} • Date: {inspection.completedAt} • Zone:{" "}
            {inspection.zone}
          </p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#00714d] rounded-full text-xs font-bold uppercase flex items-center gap-1">
          <span className="material-symbols-outlined text-[16px]">verified</span>
          {inspection.status} · {inspection.result}
        </span>
      </div>

      {/* Inspection Details Card */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6">
        <h2 className="text-lg font-bold text-[#181c1e] mb-4">{inspection.type}</h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          {inspection.criteria.map((criterion) => (
            <div key={criterion} className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
              <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
                {criterion}
              </span>
              <span className="text-sm font-bold text-[#006c49] flex items-center gap-1">
                <span className="material-symbols-outlined text-[18px]">check_circle</span>
                {inspection.result}
              </span>
            </div>
          ))}
        </div>

        <div className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg mb-6">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider block mb-1">
            Inspector Field Notes
          </span>
          <p className="text-sm text-[#181c1e] leading-relaxed">{inspection.notes}</p>
          <ul className="mt-3 space-y-1">
            {inspection.findings.map((finding) => (
              <li key={finding} className="text-sm text-[#181c1e] leading-relaxed">
                • {finding}
              </li>
            ))}
          </ul>
          <p className="text-xs text-[#44474c] font-mono mt-3">
            Certificate: {certDoc.documentReference} ({certDoc.cid})
          </p>
        </div>

        {/* Action Link to Provenance */}
        <div className="pt-4 border-t border-[#f0f3f5] flex items-center justify-between">
          <div className="text-xs text-[#44474c]">
            Requirement Verified:{" "}
            <span className="font-semibold text-[#181c1e]">{requirement.requirementReference}</span>
          </div>
          <button
            onClick={() => navigate("/dashboard/construction/provenance")}
            className="px-4 py-2 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
          >
            View Full Provenance Story →
          </button>
        </div>
      </div>
    </div>
  );
};

export default InspectionPage;
