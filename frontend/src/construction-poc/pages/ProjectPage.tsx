import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getProject, ProjectData, getRequirement, getSupplierOffer, getDelivery, getTransportDelay, getIssue, getDecision, getInspection, getEvidenceDocument } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { useRoleContext } from "../hooks/useRoleContext";
import { CreateProjectModal } from "../components/CreateProjectModal";

export const ProjectPage: React.FC = () => {
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  const { activeRole, roleConfig } = useRoleContext();

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const [project, setProject] = useState<ProjectData | undefined>(undefined);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!projectId) {
      setProject(undefined);
      setLoading(false);
      return;
    }

    const result = getProject(projectId);

    setProject(result);
    setLoading(false);

    if (!result) {
      setError(`Project ${projectId} not found`);
    } else {
      setError(null);
    }
  }, [projectId]);

  const requirement = getRequirement();
  const offer = getSupplierOffer();
  const delivery = getDelivery();
  const transportDelay = getTransportDelay();
  const issue = getIssue();
  const decision = getDecision();
  const inspection = getInspection();
  const evidence = getEvidenceDocument();

  // Loading State
  if (loading) {
    return (
      <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6 flex flex-col items-center justify-center">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-8 max-w-md text-center shadow-sm flex flex-col items-center gap-4">
          <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin" />
          <p className="text-sm text-[#44474c] font-semibold">Loading project overview from Cloud...</p>
        </div>
      </div>
    );
  }

  // Explicit API Error Presentation State (NO silent fallback to mock data)
  if (error) {
    return (
      <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6 flex flex-col items-center justify-center">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-8 max-w-md text-center shadow-sm">
          <div className="w-16 h-16 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-4">
            <span className="material-symbols-outlined text-[36px]">error</span>
          </div>
          <h2 className="text-2xl font-bold text-[#041627] mb-2">Cloud API Error</h2>
          <p className="text-sm text-[#44474c] mb-6">
            Failed to fetch project <code className="bg-[#f1f4f6] px-2 py-0.5 rounded font-mono text-[#041627]">{projectId}</code>: {error}
          </p>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="px-5 py-2.5 bg-[#041627] text-white text-sm font-semibold rounded-xl hover:bg-[#1a2b3c] transition-colors cursor-pointer"
          >
            ← Return to Projects Directory
          </button>
        </div>
      </div>
    );
  }

  // Explicit Not Found Presentation State for invalid or missing project IDs
  if (!project) {
    return (
      <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6 flex flex-col items-center justify-center">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-8 max-w-md text-center shadow-sm">
          <div className="w-16 h-16 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-4">
            <span className="material-symbols-outlined text-[36px]">domain_disabled</span>
          </div>
          <h2 className="text-2xl font-bold text-[#041627] mb-2">Project Not Found</h2>
          <p className="text-sm text-[#44474c] mb-6">
            The requested project reference <code className="bg-[#f1f4f6] px-2 py-0.5 rounded font-mono text-[#041627]">{projectId || "undefined"}</code> could not be found in this workspace collection.
          </p>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="px-5 py-2.5 bg-[#041627] text-white text-sm font-semibold rounded-xl hover:bg-[#1a2b3c] transition-colors cursor-pointer"
          >
            ← Return to Projects Directory
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6">
      <div className="flex flex-col w-full max-w-[1440px] mx-auto gap-6">
        {/* Top Command Context & Breadcrumbs */}
        <div className="flex flex-wrap items-center justify-between gap-4">
          <nav className="flex items-center gap-1 text-[#44474c] text-xs font-semibold uppercase tracking-wider">
            <span
              onClick={() => navigate(CONSTRUCTION_ROUTES.HOME)}
              className="hover:text-[#181c1e] transition-colors cursor-pointer"
            >
              Vault Home
            </span>
            <span className="material-symbols-outlined text-[14px]">chevron_right</span>
            <span
              onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
              className="hover:text-[#181c1e] transition-colors cursor-pointer"
            >
              Projects
            </span>
            <span className="material-symbols-outlined text-[14px]">chevron_right</span>
            <span className="text-[#041627] font-semibold">{project.code} Detail</span>
          </nav>

          <div className="flex items-center gap-2">
            <span className="px-3 py-1 bg-[#041627] text-white text-xs font-bold rounded-lg shadow-sm">
              {roleConfig.label} View
            </span>
          </div>
        </div>

        {/* Hero Header Card */}
        <div className="bg-white rounded-xl shadow-sm border border-[#e0e3e5] p-6">
          <div className="flex flex-wrap items-center justify-between gap-3 mb-2">
            <div className="flex items-center gap-2 text-xs font-semibold text-[#44474c]">
              <span className="bg-[#041627] text-white px-2.5 py-1 rounded font-bold">{project.code}</span>
              <span>•</span>
              <span className="uppercase tracking-wider">{project.sector || project.type}</span>
              {project.isAuthoritative && (
                <span className="px-2 py-0.5 bg-[#006c49]/10 text-[#006c49] text-[10px] font-bold rounded">
                  AUTHORITATIVE XS-BIM SCENARIO
                </span>
              )}
            </div>
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#006c49]/10 text-[#006c49] text-xs font-bold">
              <span className="w-1.5 h-1.5 rounded-full bg-[#006c49]" />
              {project.status.toUpperCase()}
            </span>
          </div>

          <h1 className="text-3xl font-bold text-[#041627] mb-2">{project.name}</h1>
          <div className="flex items-center gap-2 text-sm text-[#44474c] mb-4">
            <span className="material-symbols-outlined text-[18px] text-[#041627]">location_on</span>
            <span>{project.location}</span>
          </div>
          <p className="text-sm text-[#44474c] max-w-3xl">{project.description}</p>
        </div>

        {/* Summary Metric Cards Grid (3 Cards) */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-white rounded-xl border border-[#e0e3e5] p-4 shadow-sm flex flex-col gap-2">
            <div className="flex items-center gap-2 text-[#44474c]">
              <span className="material-symbols-outlined text-[20px] text-[#041627]">account_balance_wallet</span>
              <span className="text-xs uppercase tracking-wider font-semibold">Budget Allocation</span>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-2xl font-bold text-[#041627]">{project.budgetSpentPercent}%</span>
              <span className="text-xs text-[#44474c]">spent</span>
            </div>
            <div className="w-full h-1.5 bg-[#e5e9eb] rounded-full overflow-hidden mt-1">
              <div className="h-full bg-[#041627] rounded-full" style={{ width: `${project.budgetSpentPercent}%` }} />
            </div>
          </div>

          <div className="bg-white rounded-xl border border-[#e0e3e5] p-4 shadow-sm flex flex-col gap-2">
            <div className="flex items-center gap-2 text-[#44474c]">
              <span className="material-symbols-outlined text-[20px] text-[#006c49]">calendar_today</span>
              <span className="text-xs uppercase tracking-wider font-semibold">Schedule Status</span>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-2xl font-bold text-[#041627]">Day {project.scheduleDay}</span>
            </div>
            <div className="text-xs text-[#006c49] font-semibold flex items-center gap-1 mt-1">
              <span className="material-symbols-outlined text-[16px]">trending_up</span>
              <span>On Track • Target {project.targetCompletion}</span>
            </div>
          </div>

          <div className="bg-white rounded-xl border border-[#e0e3e5] p-4 shadow-sm flex flex-col gap-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-[#44474c]">
                <span className="material-symbols-outlined text-[20px] text-[#006c49]">verified</span>
                <span className="text-xs uppercase tracking-wider font-semibold">Compliance & Quality</span>
              </div>
              <span className="text-lg font-bold text-[#006c49]">98/100</span>
            </div>
            <div className="text-xs text-[#44474c] mt-1">
              All permits active. Inspection #{inspection.reference} passed without major citations.
            </div>
          </div>
        </div>

        {/* Main Content Split: Lifecycle Timeline + Operations Shortcuts */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          {/* Left Column: Project Lifecycle Stepper (7 Phases) */}
          <div className="lg:col-span-6 bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm flex flex-col gap-4">
            <div className="flex items-center justify-between pb-2 border-b border-[#f1f4f6]">
              <h2 className="text-lg font-bold text-[#041627]">Project Lifecycle Timeline</h2>
              <span className="text-xs font-semibold text-[#006c49]">{project.currentPhase}</span>
            </div>

            <div className="relative pl-6 space-y-4 before:absolute before:left-[11px] before:top-2 before:bottom-2 before:w-0.5 before:bg-[#e0e3e5]">
              {[
                { phase: "PHASE 1", title: "Prep & Planning", status: "COMPLETE", desc: "Site survey finalized, city permits approved.", complete: true },
                { phase: "PHASE 2", title: "Procurement", status: "COMPLETE", desc: `Steel requirement ${requirement.code} issued.`, complete: true },
                { phase: "PHASE 3", title: "Civil & Excavation", status: "COMPLETE", desc: "Foundation and spatial excavation finished.", complete: true },
                { phase: "PHASE 4", title: "Structure & Framework", status: "IN PROGRESS", desc: `Steel framework erection. ${project.progressPercent}% complete.`, active: true, progress: project.progressPercent },
                { phase: "PHASE 5", title: "Envelope & Fitout", status: "UPCOMING", desc: "Building enclosure and structural waterproofing.", upcoming: true },
                { phase: "PHASE 6", title: "MEP & Systems", status: "UPCOMING", desc: "Mechanical, electrical, and plumbing integration.", upcoming: true },
                { phase: "PHASE 7", title: "Handover & Comms", status: "UPCOMING", desc: "Final client and authority handover.", upcoming: true },
              ].map((item, idx) => (
                <div key={idx} className="relative flex items-start gap-4">
                  <div className="absolute -left-[30px] top-1">
                    {item.complete ? (
                      <div className="w-5 h-5 rounded-full bg-[#006c49] text-white flex items-center justify-center text-xs">
                        <span className="material-symbols-outlined text-[12px]">check</span>
                      </div>
                    ) : item.active ? (
                      <div className="w-5 h-5 rounded-full bg-[#041627] ring-4 ring-[#b7c8de] flex items-center justify-center">
                        <div className="w-2 h-2 rounded-full bg-white animate-pulse" />
                      </div>
                    ) : (
                      <div className="w-5 h-5 rounded-full bg-[#e0e3e5] flex items-center justify-center">
                        <div className="w-1.5 h-1.5 rounded-full bg-[#74777d]" />
                      </div>
                    )}
                  </div>

                  <div
                    className={`flex-1 p-4 rounded-xl border transition-all ${item.active
                        ? "bg-[#041627] text-white border-[#041627] shadow-md"
                        : "bg-[#f7fafc] border-[#e0e3e5] text-[#181c1e]"
                      }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span
                        className={`text-[11px] font-semibold tracking-wider ${item.active ? "text-[#b7c8de]" : "text-[#44474c]"
                          }`}
                      >
                        {item.phase}
                      </span>
                      <span
                        className={`text-[11px] font-bold ${item.complete
                            ? "text-[#006c49]"
                            : item.active
                              ? "text-[#6cf8bb]"
                              : "text-[#74777d]"
                          }`}
                      >
                        {item.status}
                      </span>
                    </div>
                    <h3 className={`text-sm font-bold mb-1 ${item.active ? "text-white" : "text-[#041627]"}`}>
                      {item.title}
                    </h3>
                    <p className={`text-xs ${item.active ? "text-[#e0e3e5]" : "text-[#44474c]"}`}>
                      {item.desc}
                    </p>

                    {item.progress !== undefined && (
                      <div className="mt-3 space-y-1">
                        <div className="w-full bg-[#1a2b3c] h-1.5 rounded-full overflow-hidden">
                          <div className="bg-[#6cf8bb] h-full rounded-full" style={{ width: `${item.progress}%` }} />
                        </div>
                        <div className="text-right text-[10px] text-[#6cf8bb] font-semibold">
                          {item.progress}%
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Right Column: Scenario Operations Action Grid */}
          <div className="lg:col-span-6 flex flex-col gap-4">
            <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm">
              <h2 className="text-lg font-bold text-[#041627] mb-1">Scenario Operational Surfaces</h2>
              <p className="text-xs text-[#44474c] mb-4">
                Access role-aware workspace views, delivery tracking, site inspections, decision logs, and provenance stories.
              </p>

              <div className="grid grid-cols-2 gap-3">
                {/* Project Channels */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.CHANNELS)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#041627] text-[22px]">forum</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Project Channels</h4>
                    <p className="text-[10px] text-[#44474c]">Collaboration Topology</p>
                  </div>
                </button>

                {/* Requirements */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.REQUIREMENTS)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#041627] text-[22px]">assignment</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Material Requirement</h4>
                    <p className="text-[10px] text-[#006c49] font-semibold">{requirement.code}</p>
                  </div>
                </button>

                {/* Supplier Offers */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.OFFERS)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#041627] text-[22px]">request_quote</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Supplier Offers</h4>
                    <p className="text-[10px] text-[#041627] font-semibold">{offer.offerReference}</p>
                  </div>
                </button>

                {/* Deliveries */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#ca8100] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#ca8100] text-[22px]">local_shipping</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Delivery Tracking</h4>
                    <p className="text-[10px] text-[#ca8100] font-semibold">{delivery.reference} (Delayed)</p>
                  </div>
                </button>

                {/* Transport Delay */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.TRANSPORT)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#ba1a1a] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#ba1a1a] text-[22px]">warning</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Transport Delay</h4>
                    <p className="text-[10px] text-[#ba1a1a] font-semibold">{transportDelay.id}</p>
                  </div>
                </button>

                {/* Issue */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.ISSUES)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#ba1a1a] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#ba1a1a] text-[22px]">bug_report</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Construction Issue</h4>
                    <p className="text-[10px] text-[#ba1a1a] font-semibold">{issue.reference}</p>
                  </div>
                </button>

                {/* Decision */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.DECISIONS)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#041627] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#041627] text-[22px]">gavel</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">Architect Decision</h4>
                    <p className="text-[10px] text-[#006c49] font-semibold">{decision.reference}</p>
                  </div>
                </button>

                {/* Inspection */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.INSPECTIONS)}
                  className="p-3 bg-[#f1f4f6] border border-[#e0e3e5] rounded-xl hover:border-[#006c49] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#006c49] text-[22px]">fact_check</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-[#041627] text-xs">QA Inspection</h4>
                    <p className="text-[10px] text-[#006c49] font-semibold">{inspection.reference}</p>
                  </div>
                </button>

                {/* Field Mode */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.FIELD)}
                  className="p-3 bg-[#181c1e] text-white rounded-xl hover:bg-[#2d3133] text-left transition-colors flex flex-col justify-between cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[#6cf8bb] text-[22px]">smartphone</span>
                  <div className="mt-2">
                    <h4 className="font-bold text-white text-xs">Field Mode UI</h4>
                    <p className="text-[10px] text-[#6cf8bb]">Site Superintendent</p>
                  </div>
                </button>

                {/* Provenance Why */}
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.PROVENANCE)}
                  className="p-3 bg-[#041627] text-white rounded-xl hover:bg-[#1a2b3c] text-left transition-colors flex flex-col justify-between cursor-pointer col-span-2"
                >
                  <div className="flex items-center justify-between">
                    <span className="material-symbols-outlined text-[#6cf8bb] text-[22px]">account_tree</span>
                    <span className="text-[11px] font-bold text-[#6cf8bb] bg-white/10 px-2 py-0.5 rounded">
                      Flagship Scenario Story
                    </span>
                  </div>
                  <div className="mt-3">
                    <h4 className="font-bold text-white text-sm">Why Is DEL-1042 Late?</h4>
                    <p className="text-xs text-[#b7c8de]">
                      Complete Provenance Chain: Reconstitution from TR-1042 → ISS-1042 → DEC-1042 → INSP-1042
                    </p>
                  </div>
                </button>
              </div>
            </div>

            {/* Stakeholders & Evidence Card */}
            <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 shadow-sm flex flex-col gap-3">
              <h3 className="text-sm font-bold text-[#041627]">Stakeholders & Traceability Artifacts</h3>
              <div className="flex items-center justify-between text-xs text-[#44474c] pt-1">
                <span>Verified Stakeholders: 7 Personas</span>
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.STAKEHOLDERS)}
                  className="text-[#041627] font-bold hover:underline cursor-pointer"
                >
                  Manage People →
                </button>
              </div>
              <div className="flex items-center justify-between text-xs text-[#44474c]">
                <span>Evidence Ledger: {evidence.code}</span>
                <button
                  type="button"
                  onClick={() => navigate(CONSTRUCTION_ROUTES.HISTORY)}
                  className="text-[#041627] font-bold hover:underline cursor-pointer"
                >
                  View Audit History →
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Create Project Modal */}
      <CreateProjectModal isOpen={isCreateModalOpen} onClose={() => setIsCreateModalOpen(false)} />
    </div>
  );
};

export default ProjectPage;
