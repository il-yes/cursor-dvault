import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const ConstructionIssuePage: React.FC = () => {
  const navigate = useNavigate();
  const { issueId } = useParams<{ issueId?: string }>();
  const [data, setData] = useState<LogisticsOverviewData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [isBookmarked, setIsBookmarked] = useState<boolean>(false);

  useEffect(() => {
    let isMounted = true;
    setLoading(true);
    setError(null);

    getLogisticsOverview("DEL-1042")
      .then((res) => {
        if (isMounted) {
          if (res) {
            setData(res);
          } else {
            setError("Issue overview data unavailable");
          }
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error("Failed to load issue data:", err);
          setError(err?.message || "Failed to load issue data");
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [issueId]);

  if (loading) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading construction issue...</p>
      </div>
    );
  }

  if (error || !data || !data.issue) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 max-w-md text-center shadow-sm">
          <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
            <span className="material-symbols-outlined text-[28px]">error</span>
          </div>
          <h2 className="text-lg font-bold text-[#041627] mb-1">Issue Details Unavailable</h2>
          <p className="text-xs text-[#44474c] mb-4">{error || "Data unavailable"}</p>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="px-4 py-2 bg-[#041627] text-white text-xs font-semibold rounded-lg hover:bg-[#1a2b3c] transition-colors cursor-pointer"
          >
            ← Return to Projects
          </button>
        </div>
      </div>
    );
  }

  const { issue, transport, decision, delivery } = data;

  const stakeholders = issue.stakeholders || [
    { name: "BuildCorp (Main)", role: "Main Contractor", icon: "domain" },
    { name: "FastBuild Logistics", role: "Carrier", icon: "local_shipping" },
    { name: "Engineering Partners", role: "Consultant", icon: "engineering" }
  ];

  const evidenceFiles = issue.evidenceFiles || [
    {
      filename: "Road_Restriction_Notice_M1.pdf",
      metadata: "Official Dept of Roads Alert • 1.4 MB",
      icon: "picture_as_pdf",
      size: "1.4 MB"
    },
    {
      filename: "Transport_Detour_Report_TR1042.pdf",
      metadata: "Carrier route survey & clearance • 2.1 MB",
      icon: "alt_route",
      size: "2.1 MB"
    }
  ];

  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: issue.title,
        text: `Construction Issue ${issue.reference || "ISS-1042"}`,
        url: window.location.href,
      }).catch(() => {});
    } else {
      navigator.clipboard?.writeText(window.location.href);
      alert("Issue link copied to clipboard");
    }
  };

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 4. Header */}
      <header className="sticky top-0 z-30 bg-[#041627] text-white px-4 py-3 shadow-md flex items-center justify-between">
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERY_DETAIL)}
            className="text-white/80 hover:text-white transition-colors cursor-pointer flex items-center justify-center mr-1"
            aria-label="Back to Delivery"
          >
            <span className="material-symbols-outlined text-[22px]">arrow_back</span>
          </button>
          <h1 className="font-bold text-base tracking-tight text-white">Construction Issue</h1>
        </div>
        <div className="flex items-center gap-3">
          <button type="button" className="text-white/80 hover:text-white transition-colors relative">
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#ba1a1a]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-2xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 5. Status Strip */}
        <div className="flex items-center justify-between bg-white border border-[#e0e3e5] rounded-xl px-4 py-3 shadow-sm">
          <div className="flex items-center gap-2">
            <span className="px-2.5 py-0.5 rounded-full bg-[#ba1a1a]/10 text-[#ba1a1a] text-xs font-extrabold flex items-center gap-1">
              <span className="w-1.5 h-1.5 rounded-full bg-[#ba1a1a]" />
              {issue.severityTag || "SEVERITY: CRITICAL"}
            </span>
            <span className="px-2.5 py-0.5 rounded-full bg-[#006c49]/10 text-[#006c49] text-xs font-extrabold flex items-center gap-1">
              <span className="material-symbols-outlined text-[14px]">check_circle</span>
              {issue.statusTag || "RESOLVED"}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setIsBookmarked(!isBookmarked)}
              className={`p-1.5 rounded-lg transition-colors cursor-pointer ${
                isBookmarked ? "text-[#041627] bg-[#f0f4f7]" : "text-[#74777f] hover:text-[#041627]"
              }`}
              title="Bookmark Issue"
            >
              <span className="material-symbols-outlined text-[20px]">
                {isBookmarked ? "bookmark" : "bookmark_border"}
              </span>
            </button>
            <button
              type="button"
              onClick={handleShare}
              className="p-1.5 rounded-lg text-[#74777f] hover:text-[#041627] transition-colors cursor-pointer"
              title="Share Issue"
            >
              <span className="material-symbols-outlined text-[20px]">share</span>
            </button>
          </div>
        </div>

        {/* 6. Issue Title / Resolution Card */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm flex flex-col gap-2">
          <div className="flex items-center justify-between text-xs font-bold text-[#74777f]">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#ba1a1a] text-[18px]">warning</span>
              <span>{issue.categoryBadge || "LOGISTICS BOTTLENECK"}</span>
            </div>
            <span>Ref: {issue.reference || "ISS-1042"}</span>
          </div>

          <h2 className="text-xl font-bold text-[#041627] tracking-tight leading-snug">
            {issue.title || "Critical Beam Delivery Delayed by Route M1 Blockage"}
          </h2>
        </div>

        {/* 7. Active Mitigation */}
        <div className="bg-[#e8f5e9] border border-[#c8e6c9] rounded-2xl p-4 shadow-sm flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-[#006c49] text-white flex items-center justify-center shadow-sm">
              <span className="material-symbols-outlined text-[20px]">alt_route</span>
            </div>
            <div>
              <span className="text-[11px] font-extrabold text-[#006c49] uppercase tracking-wide block">
                Active Mitigation
              </span>
              <span className="text-xs font-bold text-[#041627]">
                {issue.mitigationLabel || "Rerouted via Decision DEC-1042"}
              </span>
            </div>
          </div>
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DECISIONS)}
            className="px-3 py-1.5 bg-[#006c49] text-white text-xs font-bold rounded-lg hover:bg-[#006c49]/90 transition-colors cursor-pointer flex items-center gap-1"
          >
            <span>{issue.mitigationActionText || "Inspect"}</span>
            <span className="material-symbols-outlined text-[14px]">arrow_forward</span>
          </button>
        </div>

        {/* 8. Issue Metadata / Specifications */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-4">
          <div className="flex items-center justify-between border-b border-[#e0e3e5] pb-2.5">
            <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Specifications</h3>
            <span className="text-xs font-bold text-[#44474c] bg-[#f0f4f7] px-2.5 py-0.5 rounded-md">
              {issue.zoneContextText || "Zone 4 Viaduct"}
            </span>
          </div>

          <div className="grid grid-cols-2 gap-3 text-xs">
            {/* Affected Resource */}
            <div className="p-3 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <div className="flex items-center gap-1.5 text-[#74777f] mb-1">
                <span className="material-symbols-outlined text-[16px]">local_shipping</span>
                <span className="text-[11px] font-semibold uppercase">Affected Resource</span>
              </div>
              <span className="font-bold text-[#041627] block text-sm">
                {issue.affectedResourceCode || delivery.reference || "DEL-1042"}
              </span>
              <span className="text-[11px] text-[#44474c]">
                {issue.affectedResourceDescription || "Structural Beams (44t Prefabricated)"}
              </span>
            </div>

            {/* Project Phase */}
            <div className="p-3 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <div className="flex items-center gap-1.5 text-[#74777f] mb-1">
                <span className="material-symbols-outlined text-[16px]">engineering</span>
                <span className="text-[11px] font-semibold uppercase">Project Phase</span>
              </div>
              <span className="font-bold text-[#041627] block text-sm">
                {issue.projectPhaseText || "Foundation"}
              </span>
              <span className="text-[11px] text-[#44474c]">
                {issue.projectZoneText || "Viaduct Zone 4"}
              </span>
            </div>

            {/* Root Cause */}
            <div className="p-3 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <div className="flex items-center gap-1.5 text-[#74777f] mb-1">
                <span className="material-symbols-outlined text-[16px]">search</span>
                <span className="text-[11px] font-semibold uppercase">Root Cause</span>
              </div>
              <span className="font-bold text-[#ba1a1a] block text-sm">
                {issue.rootCauseTitle || "Route M1 Alert"}
              </span>
              <span className="text-[11px] text-[#44474c]">
                {issue.rootCauseSecondary || "Overpass clearance"}
              </span>
            </div>

            {/* Reported By */}
            <div className="p-3 bg-[#f7fafc] rounded-xl border border-[#e0e3e5]">
              <div className="flex items-center gap-1.5 text-[#74777f] mb-1">
                <span className="material-symbols-outlined text-[16px]">person</span>
                <span className="text-[11px] font-semibold uppercase">Reported By</span>
              </div>
              <div className="flex items-center gap-2 mt-0.5">
                <div className="w-6 h-6 rounded-full bg-[#041627] text-white flex items-center justify-center text-[10px] font-bold">
                  {issue.reporterInitials || "DK"}
                </div>
                <div>
                  <span className="font-bold text-[#041627] block text-xs leading-none">
                    {issue.reporterNameText || "David K. (Logistics Mgr)"}
                  </span>
                  <span className="text-[10px] text-[#74777f]">
                    {issue.reporterTimeText || "Aug 15, 11:41 AM"}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Assigned Stakeholders */}
          <div className="pt-2">
            <span className="text-[11px] font-bold text-[#74777f] uppercase tracking-wider block mb-2">
              Assigned Stakeholders
            </span>
            <div className="flex flex-wrap gap-2">
              {stakeholders.map((st, idx) => (
                <div
                  key={idx}
                  className="bg-[#f0f4f7] border border-[#e0e3e5] rounded-lg px-2.5 py-1 flex items-center gap-1.5 text-xs font-semibold text-[#041627]"
                >
                  <span className="material-symbols-outlined text-[16px] text-[#44474c]">{st.icon}</span>
                  <span>{st.name}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* 9. Impact Assessment */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-3">
          <div className="flex items-center justify-between mb-1">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#041627] text-[18px]">analytics</span>
              <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Impact Assessment</h3>
            </div>
            <span className="px-2.5 py-0.5 rounded-full bg-[#ffb74d]/20 text-[#b76e00] text-[11px] font-extrabold uppercase">
              {issue.impactBadgeText || "Schedule Impact"}
            </span>
          </div>

          <div className="p-3 bg-[#fff3e0] border border-[#ffe0b2] rounded-xl">
            <span className="text-[11px] font-bold text-[#b76e00] uppercase block mb-0.5">Primary Impact</span>
            <h4 className="text-xs font-bold text-[#ba1a1a]">
              {issue.primaryImpactTitle || "Foundation Phase Delayed"}
            </h4>
            <p className="text-xs text-[#041627] font-semibold mt-0.5">
              {issue.primaryImpactSubtitle || "Viaduct Zone 4 erection paused awaiting beam arrival"}
            </p>
          </div>

          {/* Horizontal Reschedule Visual Bar */}
          <div className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-xl space-y-2">
            <div className="flex justify-between text-xs font-bold">
              <span className="text-[#041627]">{issue.rescheduleLabel || "Assembly Reschedule"}</span>
              <span className="text-[#006c49]">Mitigated to &lt; 24h</span>
            </div>

            <div className="w-full h-3 bg-[#e0e3e5] rounded-full overflow-hidden flex">
              <div className="w-[65%] bg-[#ba1a1a] h-full" title="Delayed Window" />
              <div className="w-[35%] bg-[#006c49] h-full" title="Mitigated Reroute" />
            </div>

            <div className="flex justify-between text-[11px] font-semibold text-[#74777f]">
              <span className="text-[#ba1a1a]">{issue.reschedulePlannedLabel || "Aug 15 (Planned)"}</span>
              <span className="text-[#006c49]">{issue.rescheduleAdjustedLabel || "Aug 16 Morning (Adjusted)"}</span>
            </div>
          </div>

          <p className="text-xs text-[#44474c] leading-relaxed italic bg-[#f0f4f7] p-2.5 rounded-xl border border-[#e0e3e5]">
            {issue.primaryImpactDescription || "Foundation heavy assembly team diverted to auxiliary drainage culverts until crane access window re-opens."}
          </p>
        </div>

        {/* 10. Supporting Evidence */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-1.5">
              <span className="material-symbols-outlined text-[#041627] text-[18px]">folder_open</span>
              <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">
                Supporting Evidence ({evidenceFiles.length})
              </h3>
            </div>
            <span className="px-2 py-0.5 rounded bg-[#006c49]/10 text-[#006c49] text-[10px] font-extrabold uppercase">
              Verified
            </span>
          </div>

          <div className="space-y-2">
            {evidenceFiles.map((file, idx) => (
              <div
                key={idx}
                className="bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-3 flex items-center justify-between gap-3 hover:border-[#041627] transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-[#041627] text-white flex items-center justify-center shrink-0">
                    <span className="material-symbols-outlined text-[18px]">{file.icon}</span>
                  </div>
                  <div>
                    <h4 className="text-xs font-bold text-[#041627]">{file.filename}</h4>
                    <p className="text-[11px] text-[#74777f]">{file.metadata}</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => alert(`Downloading ${file.filename}...`)}
                  className="p-1.5 rounded-lg text-[#041627] hover:bg-[#e0e3e5] transition-colors cursor-pointer"
                  title="Download File"
                >
                  <span className="material-symbols-outlined text-[18px]">download</span>
                </button>
              </div>
            ))}
          </div>
        </div>

        {/* 11. Decision + Delivery Actions */}
        <div className="flex flex-col gap-2 pt-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DECISIONS)}
            className="w-full py-3.5 bg-[#041627] hover:bg-[#1a2b3c] text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">gavel</span>
            <span>View Decision ({decision?.reference || "DEC-1042"})</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERY_DETAIL)}
            className="w-full py-3 bg-white border border-[#e0e3e5] hover:bg-[#f7fafc] text-[#041627] text-xs font-bold rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer shadow-sm"
          >
            <span className="material-symbols-outlined text-[16px]">inventory_2</span>
            <span>View Delivery ({transport.reference || "DEL-1042"})</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.CHANNELS)}
            className="w-full py-2.5 bg-[#f0f4f7] hover:bg-[#e0e3e5] text-[#041627] text-xs font-bold rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">forum</span>
            <span>Open Collaboration Thread ({issue.threadUpdateCount || 8} Updates)</span>
          </button>
        </div>
      </main>

      {/* Fixed Bottom Navigation */}
      <BottomNavigation activeItem="projects" />
    </div>
  );
};

export default ConstructionIssuePage;

