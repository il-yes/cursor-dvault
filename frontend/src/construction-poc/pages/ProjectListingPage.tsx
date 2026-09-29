import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { getProjects, ProjectData, getDecision, getInspection, getIssue } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { useRoleContext } from "../hooks/useRoleContext";
import { CreateProjectModal } from "../components/CreateProjectModal";

export const ProjectListingPage: React.FC = () => {
  const navigate = useNavigate();
  const { activeRole } = useRoleContext();

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<"ALL" | "ACTIVE" | "ON_HOLD" | "ARCHIVED">("ALL");
  const [selectedPhase, setSelectedPhase] = useState("ALL");
  const [selectedOrg, setSelectedOrg] = useState("ALL");
  const [sortBy, setSortBy] = useState("ACTIVITY");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [showEmptyState, setShowEmptyState] = useState(false);
  const [filterMyAssigned, setFilterMyAssigned] = useState(false);

  const [allProjects, setAllProjects] = useState<ProjectData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchProjects = () => {
    setLoading(true);
    getProjects()
      .then((data) => {
        console.log("[BOUNDARY 6][ProjectListingPage] received allProjects count=", data.length, data);
        setAllProjects(data);
        setError(null);
        setLoading(false);
      })
      .catch((err) => {
        console.error("[BOUNDARY 6][ProjectListingPage] Failed to fetch project list:", err);
        setError(err instanceof Error ? err.message : String(err));
        setAllProjects([]);
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const authoritativeProject = allProjects.length > 0 ? (allProjects.find((p) => p.isAuthoritative) || allProjects[0]) : undefined;
  const secondaryPreviewProjects = allProjects.filter((p) => p !== authoritativeProject);

  const decision = getDecision();
  const inspection = getInspection();
  const issue = getIssue();

  const filteredProjects = allProjects.filter((p) => {
    if (filterMyAssigned && !p.isAuthoritative) return false;
    const matchesSearch =
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.location.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus =
      statusFilter === "ALL" ||
      (statusFilter === "ACTIVE" && p.status.includes("Active")) ||
      (statusFilter === "ON_HOLD" && p.status.includes("Hold")) ||
      (statusFilter === "ARCHIVED" && p.status.includes("Archived"));
    return matchesSearch && matchesStatus;
  });

  const avgProgress = allProjects.length > 0 ? Math.round(allProjects.reduce((sum, p) => sum + (p.progressPercent || 0), 0) / allProjects.length) : 0;
  const totalIssues = allProjects.reduce((sum, p) => sum + (p.openIssuesCount || 0), 0);
  const totalDecisions = allProjects.reduce((sum, p) => sum + (p.pendingDecisionsCount || 0), 0);

  if (loading) {
    return (
      <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6 flex flex-col items-center justify-center">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-8 max-w-md text-center shadow-sm flex flex-col items-center gap-4">
          <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin" />
          <p className="text-sm text-[#44474c] font-semibold">Loading projects directory...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6 flex flex-col items-center justify-center">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-8 max-w-md text-center shadow-sm">
          <div className="w-16 h-16 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-4">
            <span className="material-symbols-outlined text-[36px]">error</span>
          </div>
          <h2 className="text-2xl font-bold text-[#041627] mb-2">Unable to Load Projects</h2>
          <p className="text-sm text-[#44474c] mb-6">
            Failed to fetch construction projects list: {error}
          </p>
          <button
            type="button"
            onClick={() => window.location.reload()}
            className="px-5 py-2.5 bg-[#041627] text-white text-sm font-semibold rounded-xl hover:bg-[#1a2b3c] transition-colors cursor-pointer"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  const isActuallyEmpty = allProjects.length === 0 || showEmptyState;

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] p-6">
      <div className="flex flex-col w-full max-w-[1440px] mx-auto">
        {/* Top Command Context Bar */}
        <div className="flex flex-col gap-3 mb-6">
          {/* Breadcrumb & Persona Ribbon */}
          <div className="flex flex-wrap items-center justify-between gap-4">
            <nav className="flex items-center gap-1 text-[#44474c] text-xs font-semibold uppercase tracking-wider">
              <span
                onClick={() => navigate(CONSTRUCTION_ROUTES.HOME)}
                className="hover:text-[#181c1e] transition-colors cursor-pointer"
              >
                Vault Home
              </span>
              <span className="material-symbols-outlined text-[14px]">chevron_right</span>
              <span className="text-[#041627] font-semibold">Workspace</span>
              <span className="material-symbols-outlined text-[14px]">chevron_right</span>
              <span className="text-[#181c1e] font-semibold">Projects</span>
            </nav>

            {/* Role-Awareness Micro-Badge */}
            <div className="flex items-center gap-2 bg-white px-3 py-1.5 rounded-xl shadow-sm border border-[#e0e3e5]">
              <div className="w-2 h-2 rounded-full bg-[#006c49]" />
              <span className="text-xs text-[#44474c]">Logged in as</span>
              <span className="text-sm font-semibold text-[#041627]">Manuel Vincent</span>
              <span className="text-xs text-[#44474c] bg-[#e5e9eb] px-2 py-0.5 rounded font-medium">
                Structural Engineer • Engineering Partners
              </span>
              <button
                type="button"
                onClick={() => setFilterMyAssigned(!filterMyAssigned)}
                className={`flex items-center gap-1 text-xs font-semibold transition-colors pl-2 cursor-pointer ${
                  filterMyAssigned
                    ? "bg-[#006c49] text-white px-2 py-0.5 rounded"
                    : "text-[#006c49] hover:text-[#005236]"
                }`}
              >
                <span className="material-symbols-outlined text-[14px]">assignment_ind</span>
                My Assigned Projects ({allProjects.length})
              </button>
            </div>
          </div>

          {/* Page Title & Primary Actions */}
          <div className="flex flex-wrap items-end justify-between gap-4 pt-2">
            <div>
              <h1 className="text-3xl font-bold text-[#041627] tracking-tight">Projects</h1>
              <p className="text-base text-[#44474c] mt-1">
                Your construction projects across connected engineering, contracting, and client organizations.
              </p>
            </div>
            <div className="flex items-center gap-3">
              <button
                type="button"
                className="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-white text-[#041627] text-sm font-semibold hover:bg-[#ebeef0] transition-colors shadow-sm border border-[#e0e3e5]"
              >
                <span className="material-symbols-outlined text-[20px]">file_download</span>
                Import Project
              </button>

              {activeRole === "PM" && (
                <button
                  type="button"
                  onClick={() => setIsCreateModalOpen(true)}
                  className="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-[#041627] text-white text-sm font-semibold hover:bg-[#1a2b3c] transition-colors shadow-sm"
                >
                  <span className="material-symbols-outlined text-[20px]">add</span>
                  Create Project
                </button>
              )}
            </div>
          </div>
        </div>

        {/* Operational Metric Snapshot (4-Column Grid) */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white rounded-xl p-4 shadow-sm border border-[#e0e3e5] flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#44474c]">
              <span className="text-xs uppercase tracking-wider font-semibold">Active Workspace Projects</span>
              <span className="material-symbols-outlined text-[#041627] text-[20px]">domain</span>
            </div>
            <div className="flex items-baseline gap-2 mt-3">
              <span className="text-3xl font-bold text-[#041627]">{allProjects.length}</span>
              <span className="text-xs text-[#006c49] font-medium">Data Source Active</span>
            </div>
            <div className="w-full bg-[#e5e9eb] h-1.5 rounded-full mt-3 overflow-hidden">
              <div className="bg-[#041627] h-full rounded-full" style={{ width: allProjects.length > 0 ? "100%" : "0%" }} />
            </div>
          </div>

          <div className="bg-white rounded-xl p-4 shadow-sm border border-[#e0e3e5] flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#44474c]">
              <span className="text-xs uppercase tracking-wider font-semibold">Average Milestone Progress</span>
              <span className="material-symbols-outlined text-[#006c49] text-[20px]">pie_chart</span>
            </div>
            <div className="flex items-baseline gap-2 mt-3">
              <span className="text-3xl font-bold text-[#041627]">{avgProgress}%</span>
              <span className="text-xs text-[#006c49] font-medium">On Track</span>
            </div>
            <div className="w-full bg-[#e5e9eb] h-1.5 rounded-full mt-3 overflow-hidden">
              <div className="bg-[#006c49] h-full rounded-full" style={{ width: `${avgProgress}%` }} />
            </div>
          </div>

          <div className="bg-white rounded-xl p-4 shadow-sm border border-[#e0e3e5] flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#44474c]">
              <span className="text-xs uppercase tracking-wider font-semibold">Open Critical Issues</span>
              <span className="material-symbols-outlined text-[#ba1a1a] text-[20px]">report_problem</span>
            </div>
            <div className="flex items-baseline gap-2 mt-3">
              <span className="text-3xl font-bold text-[#ba1a1a]">{totalIssues}</span>
              <span className="text-xs text-[#44474c]">Across projects</span>
            </div>
            <div className="w-full bg-[#e5e9eb] h-1.5 rounded-full mt-3 overflow-hidden">
              <div className="bg-[#ba1a1a] h-full rounded-full" style={{ width: totalIssues > 0 ? "50%" : "0%" }} />
            </div>
          </div>

          <div className="bg-white rounded-xl p-4 shadow-sm border border-[#e0e3e5] flex flex-col justify-between">
            <div className="flex items-center justify-between text-[#44474c]">
              <span className="text-xs uppercase tracking-wider font-semibold">Pending Decisions (DEC)</span>
              <span className="material-symbols-outlined text-[#ca8100] text-[20px]">assignment_turned_in</span>
            </div>
            <div className="flex items-baseline gap-2 mt-3">
              <span className="text-3xl font-bold text-[#041627]">{totalDecisions}</span>
              <span className="text-xs text-[#44474c]">Requires Signoff</span>
            </div>
            <div className="w-full bg-[#e5e9eb] h-1.5 rounded-full mt-3 overflow-hidden">
              <div className="bg-[#ca8100] h-full rounded-full" style={{ width: totalDecisions > 0 ? "30%" : "0%" }} />
            </div>
          </div>
        </div>

        {/* Toolbar & Filter Controls */}
        <div className="bg-white rounded-xl p-4 shadow-sm border border-[#e0e3e5] mb-6 flex flex-col gap-4">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
            {/* Search Input Field */}
            <div className="relative flex-1 min-w-[280px]">
              <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#44474c] text-[20px]">
                search
              </span>
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search projects by name, reference, or location..."
                className="w-full pl-10 pr-12 py-2.5 bg-[#f1f4f6] rounded-lg text-sm text-[#181c1e] placeholder:text-[#44474c] outline-none focus:bg-white focus:ring-1 focus:ring-[#041627] transition-all"
              />
              <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[11px] font-semibold text-[#44474c] bg-[#e5e9eb] px-1.5 py-0.5 rounded">
                ⌘K
              </span>
            </div>

            {/* Right Toolbar Selectors */}
            <div className="flex flex-wrap items-center gap-3">
              {/* Dropdown: Phase */}
              <select
                value={selectedPhase}
                onChange={(e) => setSelectedPhase(e.target.value)}
                className="bg-[#f1f4f6] text-[#181c1e] text-sm py-2 px-3 rounded-lg border-none outline-none cursor-pointer hover:bg-[#e5e9eb] transition-colors"
              >
                <option value="ALL">Phase: All Phases</option>
                <option value="P1">Phase 1 • Planning</option>
                <option value="P2">Phase 2 • Procurement</option>
                <option value="P3">Phase 3 • Foundation</option>
                <option value="P4">Phase 4 • Structure</option>
                <option value="P5">Phase 5 • Envelope</option>
                <option value="P6">Phase 6 • MEP Fitout</option>
                <option value="P7">Phase 7 • Handover</option>
              </select>

              {/* Dropdown: Org */}
              <select
                value={selectedOrg}
                onChange={(e) => setSelectedOrg(e.target.value)}
                className="bg-[#f1f4f6] text-[#181c1e] text-sm py-2 px-3 rounded-lg border-none outline-none cursor-pointer hover:bg-[#e5e9eb] transition-colors"
              >
                <option value="ALL">Organization: All Orgs</option>
                <option value="Acme">Acme Dev</option>
                <option value="BuildCorp">BuildCorp</option>
                <option value="Engineering">Engineering Partners</option>
                <option value="EuroSteel">EuroSteel</option>
              </select>

              {/* Dropdown: Sort */}
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className="bg-[#f1f4f6] text-[#181c1e] text-sm py-2 px-3 rounded-lg border-none outline-none cursor-pointer hover:bg-[#e5e9eb] transition-colors"
              >
                <option value="ACTIVITY">Sort by: Last Activity</option>
                <option value="PROGRESS">Sort by: Progress (High to Low)</option>
                <option value="DATE">Sort by: Target Completion Date</option>
                <option value="REF">Sort by: Reference ID</option>
              </select>

              {/* View Mode Toggle */}
              <div className="flex items-center bg-[#f1f4f6] p-1 rounded-lg">
                <button
                  type="button"
                  onClick={() => setViewMode("grid")}
                  className={`p-1.5 rounded transition-all ${
                    viewMode === "grid"
                      ? "bg-white text-[#041627] shadow-sm"
                      : "text-[#44474c] hover:text-[#181c1e]"
                  }`}
                  title="Grid View"
                >
                  <span className="material-symbols-outlined text-[18px]">grid_view</span>
                </button>
                <button
                  type="button"
                  onClick={() => setViewMode("list")}
                  className={`p-1.5 rounded transition-all ${
                    viewMode === "list"
                      ? "bg-white text-[#041627] shadow-sm"
                      : "text-[#44474c] hover:text-[#181c1e]"
                  }`}
                  title="List View"
                >
                  <span className="material-symbols-outlined text-[18px]">view_list</span>
                </button>
              </div>

              {/* Empty State Toggle Demo Pill */}
              <button
                type="button"
                onClick={() => setShowEmptyState(!showEmptyState)}
                className={`flex items-center gap-1 px-3 py-2 rounded-lg text-xs font-semibold uppercase tracking-wider transition-colors cursor-pointer ${
                  showEmptyState
                    ? "bg-[#041627] text-white"
                    : "text-[#44474c] hover:text-[#041627] hover:bg-[#f1f4f6]"
                }`}
              >
                <span className="material-symbols-outlined text-[16px]">visibility</span>
                <span>Toggle Empty State</span>
              </button>
            </div>
          </div>

          {/* Status Category Tabs */}
          <div className="flex flex-wrap items-center gap-2 pt-1 border-t border-[#f1f4f6]">
            {(
              [
                { id: "ALL", label: `All Projects (${allProjects.length})` },
                { id: "ACTIVE", label: `Active (${allProjects.filter(p => p.status.includes("Active")).length})` },
                { id: "ON_HOLD", label: `On Hold (${allProjects.filter(p => p.status.includes("Hold")).length})` },
                { id: "ARCHIVED", label: `Archived (${allProjects.filter(p => p.status.includes("Archived")).length})` },
              ] as const
            ).map((tab) => (
              <button
                key={tab.id}
                type="button"
                onClick={() => setStatusFilter(tab.id)}
                className={`px-3.5 py-1.5 rounded-lg text-sm font-semibold transition-all ${
                  statusFilter === tab.id
                    ? "bg-[#041627] text-white"
                    : "text-[#44474c] hover:bg-[#f1f4f6]"
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>

        {/* Toggleable Empty State Section */}
        {isActuallyEmpty ? (
          <div className="flex flex-col items-center justify-center p-12 bg-white rounded-xl border border-[#e0e3e5] shadow-sm text-center my-6">
            <div className="w-20 h-20 rounded-full bg-[#e5e9eb] flex items-center justify-center text-[#44474c] mb-4">
              <span className="material-symbols-outlined text-[44px]">domain_disabled</span>
            </div>
            <h3 className="text-2xl font-bold text-[#041627]">No active projects in this workspace yet</h3>
            <p className="text-sm text-[#44474c] max-w-md mt-2 mb-6">
              This sovereign workspace does not currently have imported projects or active site links. Start by drafting a new commercial building scope or importing existing plans.
            </p>
            <div className="flex items-center gap-3">
              <button
                type="button"
                onClick={() => {
                  setShowEmptyState(false);
                  setIsCreateModalOpen(true);
                }}
                className="px-4 py-2.5 rounded-xl bg-[#041627] text-white text-sm font-semibold hover:bg-[#1a2b3c] transition-colors cursor-pointer"
              >
                + Create your first project
              </button>
              {showEmptyState && (
                <button
                  type="button"
                  onClick={() => setShowEmptyState(false)}
                  className="px-4 py-2.5 rounded-xl bg-[#f1f4f6] text-[#041627] text-sm font-semibold hover:bg-[#e5e9eb] transition-colors cursor-pointer"
                >
                  Restore Projects View
                </button>
              )}
            </div>
          </div>
        ) : (
          /* Projects Grid Container */
          <div className="flex flex-col gap-6">
            {/* AUTHORITATIVE HERO CARD (Primary Layout) */}
            {authoritativeProject && (
              <div className="bg-white rounded-xl shadow-sm border border-[#e0e3e5] overflow-hidden flex flex-col transition-all hover:shadow-md">
                {/* Card Top Structural Bar */}
                <div className="bg-[#f1f4f6] px-6 py-3 flex flex-wrap items-center justify-between gap-3 border-b border-[#e5e9eb]">
                  <div className="flex items-center gap-3">
                    <span className="text-lg font-bold text-[#041627] tracking-tight">{authoritativeProject.code}</span>
                    {authoritativeProject.location && (
                      <>
                        <span className="text-[#44474c]">•</span>
                        <div className="flex items-center gap-1 text-[#44474c] text-sm">
                          <span className="material-symbols-outlined text-[18px]">location_on</span>
                          <span>{authoritativeProject.location}</span>
                        </div>
                      </>
                    )}
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#006c49]/10 text-[#006c49] text-xs font-semibold">
                      <span className="w-1.5 h-1.5 rounded-full bg-[#006c49]" />
                      {authoritativeProject.status || "Active"}
                    </span>
                    {authoritativeProject.isAuthoritative && (
                      <span className="text-xs font-semibold text-[#44474c] bg-[#e0e3e5] px-2.5 py-1 rounded">
                        Primary Authority
                      </span>
                    )}
                  </div>
                </div>

                {/* Main Body Split */}
                <div className="p-6 grid grid-cols-1 lg:grid-cols-12 gap-6">
                  {/* Left: Spatial Context & Blueprint Thumbnail */}
                  <div className="lg:col-span-4 flex flex-col gap-3">
                    <div className="relative w-full h-52 rounded-xl overflow-hidden shadow-inner bg-[#1a2b3c]">
                      {authoritativeProject.image ? (
                        <img
                          src={authoritativeProject.image}
                          alt={`${authoritativeProject.name} Context`}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-[#b7c8de]">
                          <span className="material-symbols-outlined text-[48px]">domain</span>
                        </div>
                      )}
                      <div className="absolute inset-0 bg-gradient-to-t from-[#041627]/80 via-transparent to-transparent" />
                      <div className="absolute bottom-3 left-3 right-3 text-white flex items-center justify-between">
                        <div>
                          <p className="text-[11px] uppercase tracking-wider text-[#e5e9eb] font-medium">
                            {authoritativeProject.sector || authoritativeProject.type || "Construction Project"}
                          </p>
                          <p className="text-sm font-semibold">{authoritativeProject.name}</p>
                        </div>
                        <div className="w-8 h-8 rounded-lg bg-white/20 backdrop-blur-md flex items-center justify-center">
                          <span className="material-symbols-outlined text-[20px] text-white">domain</span>
                        </div>
                      </div>
                    </div>

                    {/* Connected Organizations Chips */}
                    {authoritativeProject.connectedOrgs.length > 0 && (
                      <div className="flex flex-col gap-1 pt-1">
                        <span className="text-[11px] text-[#44474c] uppercase tracking-wider font-semibold">
                          Connected Organizations
                        </span>
                        <div className="flex items-center gap-1.5 flex-wrap">
                          {authoritativeProject.connectedOrgs.map((org) => (
                            <span
                              key={org}
                              className="px-2.5 py-0.5 rounded bg-[#ebeef0] text-xs text-[#181c1e] font-medium"
                            >
                              {org}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Center & Right: Phase Tracking, Indicators & Details */}
                  <div className="lg:col-span-8 flex flex-col justify-between gap-4">
                    <div>
                      <div className="flex flex-wrap items-baseline justify-between gap-3 mb-1">
                        <h2 className="text-2xl font-bold text-[#041627]">{authoritativeProject.name}</h2>
                        {authoritativeProject.contractId && (
                          <span className="text-xs text-[#44474c]">
                            Contract ID: {authoritativeProject.contractId}
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-[#44474c] line-clamp-2">
                        {authoritativeProject.description}
                      </p>
                    </div>

                    {/* Phase & Progress Ribbon */}
                    <div className="bg-[#f1f4f6] rounded-xl p-4 flex flex-col gap-3">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-1.5">
                          <span className="text-xs uppercase font-bold text-[#181c1e] tracking-wider">
                            Current Phase:
                          </span>
                          <span className="text-sm font-semibold text-[#041627]">
                            {authoritativeProject.currentPhase || "Active Operations"}
                          </span>
                        </div>
                        <span className="text-lg font-bold text-[#041627]">
                          {authoritativeProject.progressPercent}% Complete
                        </span>
                      </div>

                      {/* Progress Bar */}
                      <div className="w-full bg-[#e0e3e5] h-2.5 rounded-full overflow-hidden">
                        <div
                          className="bg-[#006c49] h-full rounded-full transition-all duration-500"
                          style={{ width: `${authoritativeProject.progressPercent}%` }}
                        />
                      </div>
                    </div>

                    {/* Key Operational Indicators */}
                    <div className="flex flex-wrap items-center gap-2">
                      {authoritativeProject.activeDelay && (
                        <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-[#ffddb8] text-[#2a1700] text-xs font-semibold">
                          <span className="material-symbols-outlined text-[16px]">schedule</span>
                          Active Delay: {authoritativeProject.activeDelay}
                        </span>
                      )}
                      <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-[#ebeef0] text-[#181c1e] text-xs font-semibold">
                        <span className="material-symbols-outlined text-[16px]">pending_actions</span>
                        {authoritativeProject.pendingDecisionsCount} Pending Decision{authoritativeProject.pendingDecisionsCount === 1 ? "" : "s"}
                      </span>
                      <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-[#ffdad6] text-[#93000a] text-xs font-semibold">
                        <span className="material-symbols-outlined text-[16px]">bug_report</span>
                        {authoritativeProject.openIssuesCount} Open Issue{authoritativeProject.openIssuesCount === 1 ? "" : "s"}
                      </span>
                      {authoritativeProject.targetCompletion && (
                        <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-[#e0e3e5] text-[#44474c] text-xs">
                          <span className="material-symbols-outlined text-[16px]">event_available</span>
                          Target: {authoritativeProject.targetCompletion}
                        </span>
                      )}
                    </div>

                    {/* Bottom Activity & Action Ribbon */}
                    <div className="pt-3 flex flex-wrap items-center justify-between gap-3 border-t border-[#e5e9eb]">
                      <div className="flex items-center gap-1.5 text-[#44474c] text-xs">
                        <span className="material-symbols-outlined text-[18px] text-[#006c49]">check_circle</span>
                        <span className="font-semibold text-[#181c1e]">Recent Activity:</span>
                        <span className="truncate max-w-xs">{authoritativeProject.recentActivity || "No recent activity"}</span>
                      </div>

                      <button
                        type="button"
                        onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECT(authoritativeProject.id))}
                        className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#041627] text-white text-sm font-semibold hover:bg-[#1a2b3c] transition-all hover:translate-x-0.5 shadow-sm cursor-pointer"
                      >
                        <span>Open Project</span>
                        <span className="material-symbols-outlined text-[18px]">arrow_forward</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Secondary Project Cards Grid */}
            {secondaryPreviewProjects.length > 0 && (
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {secondaryPreviewProjects.map((p) => (
                  <div
                    key={p.id}
                    className="bg-white rounded-xl shadow-sm border border-[#e0e3e5] overflow-hidden flex flex-col justify-between hover:shadow-md transition-all"
                  >
                    <div>
                      {/* Header Bar */}
                      <div className="bg-[#f1f4f6] px-6 py-3 flex items-center justify-between border-b border-[#e5e9eb]">
                        <div className="flex items-center gap-2">
                          <span className="text-base font-bold text-[#041627]">{p.code}</span>
                          {p.location && (
                            <>
                              <span className="text-[#44474c]">•</span>
                              <span className="text-xs text-[#44474c] flex items-center gap-1">
                                <span className="material-symbols-outlined text-[16px]">location_on</span>
                                {p.location}
                              </span>
                            </>
                          )}
                        </div>
                        <span
                          className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                            p.status.includes("Active")
                              ? "bg-[#006c49]/10 text-[#006c49]"
                              : "bg-[#b7c8de]/30 text-[#0b1d2d]"
                          }`}
                        >
                          <span
                            className={`w-1.5 h-1.5 rounded-full ${
                              p.status.includes("Active") ? "bg-[#006c49]" : "bg-[#041627]"
                            }`}
                          />
                          {p.status || "Active"}
                        </span>
                      </div>

                      <div className="p-6 flex flex-col gap-4">
                        {/* Thumbnail & Title Group */}
                        <div className="flex gap-4">
                          <div className="w-24 h-24 rounded-lg overflow-hidden flex-shrink-0 relative bg-[#e5e9eb]">
                            {p.image ? (
                              <img src={p.image} alt={p.name} className="w-full h-full object-cover" />
                            ) : (
                              <div className="w-full h-full flex items-center justify-center text-[#44474c]">
                                <span className="material-symbols-outlined text-[32px]">domain</span>
                              </div>
                            )}
                          </div>
                          <div className="flex flex-col justify-center">
                            <span className="text-[11px] uppercase tracking-wider text-[#44474c] font-semibold">
                              {p.type || "Project"}
                            </span>
                            <h3 className="text-lg font-bold text-[#041627]">{p.name}</h3>
                            {p.connectedOrgs.length > 0 && (
                              <span className="text-xs text-[#44474c]">
                                Connected: {p.connectedOrgs.join(", ")}
                              </span>
                            )}
                          </div>
                        </div>

                        {/* Phase & Metric Block */}
                        <div className="bg-[#f1f4f6] p-3 rounded-lg flex flex-col gap-1.5">
                          <div className="flex justify-between text-xs">
                            <span className="text-[#44474c]">
                              {p.currentPhase || "Active Phase"}
                            </span>
                            <span className="font-bold text-[#041627]">{p.progressPercent}% Complete</span>
                          </div>
                          <div className="w-full bg-[#e0e3e5] h-2 rounded-full overflow-hidden">
                            <div
                              className="bg-[#041627] h-full rounded-full"
                              style={{ width: `${p.progressPercent}%` }}
                            />
                          </div>
                        </div>

                        {/* Indicators */}
                        <div className="flex flex-wrap gap-1.5">
                          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded bg-[#ebeef0] text-[#006c49] text-xs font-semibold">
                            <span className="material-symbols-outlined text-[14px]">verified</span>
                            {p.openIssuesCount} Critical Issues
                          </span>
                          <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded bg-[#ebeef0] text-[#181c1e] text-xs font-semibold">
                            <span className="material-symbols-outlined text-[14px]">rate_review</span>
                            {p.pendingDecisionsCount} Pending Reviews
                          </span>
                          {p.targetCompletion && (
                            <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded bg-[#ebeef0] text-[#44474c] text-xs">
                              Target: {p.targetCompletion}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Footer Strip */}
                    <div className="px-6 py-3 bg-white flex items-center justify-between border-t border-[#e5e9eb]">
                      <div className="text-[#44474c] text-xs truncate max-w-[220px]">
                        <span className="font-semibold text-[#181c1e]">Recent:</span> {p.recentActivity || "Updated"}
                      </div>
                      <button
                        type="button"
                        onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECT(p.id))}
                        className="inline-flex items-center gap-1 font-body-md text-xs font-semibold text-[#041627] hover:text-[#006c49] transition-colors cursor-pointer"
                      >
                        <span>Open Project</span>
                        <span className="material-symbols-outlined text-[16px]">arrow_forward</span>
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Sovereign Workspace Security Banner */}
        <div className="mt-8 bg-white p-4 rounded-xl border border-[#e0e3e5] shadow-sm flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-[#f1f4f6] flex items-center justify-center text-[#041627]">
              <span className="material-symbols-outlined text-[24px]">verified_user</span>
            </div>
            <div>
              <h4 className="text-sm font-bold text-[#041627]">BuildFlow Sovereign Workspace Protection</h4>
              <p className="text-xs text-[#44474c]">
                Access to construction projects is scoped under your validated structural engineer authority credentials.
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs text-[#44474c] bg-[#f1f4f6] px-3 py-1.5 rounded font-mono">
              Workspace ID: SW-FR-7501
            </span>
            <span className="text-xs text-[#006c49] bg-[#006c49]/10 px-3 py-1.5 rounded font-semibold">
              Audit Status: Synchronized
            </span>
          </div>
        </div>
      </div>

      {/* Create Project Modal */}
      <CreateProjectModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onProjectCreated={fetchProjects}
      />
    </div>
  );
};

export default ProjectListingPage;
