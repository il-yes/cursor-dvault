import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const InspectionPage: React.FC = () => {
  const navigate = useNavigate();
  const [data, setData] = useState<LogisticsOverviewData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

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
            setError("Inspection data unavailable");
          }
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error("Failed to load inspection data:", err);
          setError(err?.message || "Failed to load inspection data");
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  if (loading) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="w-8 h-8 border-4 border-[#041627] border-t-transparent rounded-full animate-spin mb-3" />
        <p className="text-sm font-semibold text-[#44474c]">Loading quality inspection...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="w-full min-h-screen bg-[#f7fafc] text-[#181c1e] font-[Inter] flex flex-col items-center justify-center p-6">
        <div className="bg-white rounded-xl border border-[#e0e3e5] p-6 max-w-md text-center shadow-sm">
          <div className="w-12 h-12 rounded-full bg-[#ffdad6] text-[#ba1a1a] flex items-center justify-center mx-auto mb-3">
            <span className="material-symbols-outlined text-[28px]">error</span>
          </div>
          <h2 className="text-lg font-bold text-[#041627] mb-1">Inspection Details Unavailable</h2>
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

  const inspectionPage = data.inspectionPageData || {
    id: "INSP-1042",
    reference: "INSP-REF-1042",
    statusTag: "INSPECTION COMPLETE",
    title: "Foundation — Zone A",
    inspectorText: "Inspector: Bureau Inspection",
    checklist: [
      {
        id: 1,
        title: "Reinforcement Layout",
        description: "Spacing and bar size verified per structural plans (S-201).",
        isPassed: true,
      },
      {
        id: 2,
        title: "Formwork Integrity",
        description: "Bracing and dimensions confirmed against formwork design (F-10).",
        isPassed: true,
      },
      {
        id: 3,
        title: "Concrete Quality",
        description: "Slump test passed. Mix design matched specifications.",
        isPassed: true,
      },
      {
        id: 4,
        title: "Overall Dimensions",
        description: "Tolerances within acceptable limits (+/- 5mm).",
        isPassed: true,
      },
      {
        id: 5,
        title: "Site Safety Protocol",
        description: "All workers wearing appropriate PPE. Trench shoring intact.",
        isPassed: true,
      },
    ],
    evidencePhotos: [
      { id: "photo-1", label: "Rebar Layout", imageAlt: "Rebar Layout Photo", src: "https://cloud.squidex.io/api/assets/matmatch-cms/cd36431e-2c51-4122-90ad-5f01559e3596/steel.jpg" },
      { id: "photo-2", label: "Formwork", imageAlt: "Formwork Photo", src: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRaucyCynfd1bJZNHR0KdHgrNp5yLjgmYp7DBMwvh6OZRHJXOOLmNLJkINX&s=10" },
    ],
    evidenceDocuments: [
      { id: "doc-1", icon: "science", title: "Concrete Test Results", metadata: "PDF • 1.2 MB" },
      { id: "doc-2", icon: "description", title: "Official Inspection Report", metadata: "PDF • 3.4 MB" },
    ],
    timeline: [
      {
        id: "timeline-1",
        category: "Final Approval",
        title: "Approved by Bureau Inspection",
        timestamp: "Oct 24, 2023 - 14:30",
        isTerminal: true,
      },
      {
        id: "timeline-2",
        category: "Review",
        title: "Site walk-through completed",
        timestamp: "Oct 24, 2023 - 10:15",
      },
      {
        id: "timeline-3",
        category: "Initiation",
        title: "Inspection requested by Contractor",
        timestamp: "Oct 23, 2023 - 09:00",
      },
    ],
  };

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 4. HEADER */}
      <header className="sticky top-0 z-30 bg-white/80 backdrop-blur-xl border-b border-[#e0e3e5] text-[#041627] px-4 py-3 shadow-[0_1px_4px_rgba(0,0,0,0.03)] flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="material-symbols-outlined text-[#006c49] text-[22px]">domain</span>
          <span className="font-bold text-base tracking-tight">BuildFlow</span>
        </div>
        <div className="flex items-center gap-3">
          <button
            type="button"
            className="text-[#44474c] hover:text-[#041627] transition-colors relative cursor-pointer"
            aria-label="Notifications"
          >
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#006c49]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-[#041627] text-white flex items-center justify-center text-xs font-bold shadow-sm">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-2xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 5. STATUS / INSPECTION HEADER */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm flex flex-col gap-2">
          <div className="flex items-center gap-1.5 text-xs font-extrabold text-[#006c49] uppercase tracking-wide">
            <span className="material-symbols-outlined text-[16px]">verified</span>
            <span>{inspectionPage.statusTag}</span>
          </div>

          <h1 className="text-2xl font-bold text-[#041627] tracking-tight">
            {inspectionPage.title}
          </h1>

          <p className="text-xs font-semibold text-[#74777f]">
            {inspectionPage.inspectorText}
          </p>
        </div>

        {/* 6. VERIFICATION CHECKLIST */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-3">
          <h2 className="text-xs font-bold text-[#041627] uppercase tracking-wider">
            Verification Checklist
          </h2>

          <div className="space-y-3 pt-1">
            {inspectionPage.checklist.map((item) => (
              <div key={item.id} className="flex items-start gap-3">
                <div className="w-5 h-5 rounded-full bg-[#006c49]/10 text-[#006c49] flex items-center justify-center shrink-0 mt-0.5">
                  <span className="material-symbols-outlined text-[16px]">check_circle</span>
                </div>
                <div className="flex-1 border-b border-[#f0f4f7] pb-2 last:border-0 last:pb-0">
                  <h3 className="text-xs font-bold text-[#041627]">{item.title}</h3>
                  <p className="text-xs text-[#44474c] leading-normal mt-0.5">{item.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 7. SUPPORTING EVIDENCE */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-3">
          <h2 className="text-xs font-bold text-[#041627] uppercase tracking-wider">
            Supporting Evidence
          </h2>

          {/* Photos 2-column Grid */}
          <div className="grid grid-cols-2 gap-3">
            {inspectionPage.evidencePhotos.map((photo) => (
              <div
                key={photo.id}
                className="relative h-36 rounded-xl overflow-hidden bg-[#041627] border border-[#e0e3e5] group shadow-sm flex flex-col justify-end p-2.5"
              >
                {/* Styled structural graphic container */}
                <div className="absolute inset-0 bg-gradient-to-t from-[#041627] via-[#0b2b48]/80 to-[#1e3a5f]/60 transition-transform duration-300 group-hover:scale-105" />
                <div className="absolute inset-0 opacity-20 bg-[radial-gradient(#6cf8bb_1px,transparent_1px)] [background-size:12px_12px]"
                  style={{ backgroundImage: `url("${photo?.src}")` }} />

                <div className="relative z-10">
                  <span className="inline-block px-2 py-0.5 rounded-md bg-[#041627]/70 backdrop-blur-md text-white text-[11px] font-bold">
                    {photo.label}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {/* 8. SUPPORTING DOCUMENTS */}
          <div className="space-y-2 pt-2">
            {inspectionPage.evidenceDocuments.map((doc) => (
              <div
                key={doc.id}
                className="bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-3 flex items-center justify-between gap-3 hover:border-[#041627] transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-[#041627] text-white flex items-center justify-center shrink-0">
                    <span className="material-symbols-outlined text-[18px]">{doc.icon}</span>
                  </div>
                  <div>
                    <h3 className="text-xs font-bold text-[#041627]">{doc.title}</h3>
                    <p className="text-[11px] text-[#74777f]">{doc.metadata}</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => alert(`Downloading ${doc.title}...`)}
                  className="p-1.5 rounded-lg text-[#041627] hover:bg-[#e0e3e5] transition-colors cursor-pointer"
                  title="Download Document"
                >
                  <span className="material-symbols-outlined text-[18px]">download</span>
                </button>
              </div>
            ))}
          </div>
        </div>

        {/* 9. WORKFLOW TIMELINE */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-3 relative">
          <h2 className="text-xs font-bold text-[#041627] uppercase tracking-wider">
            Workflow Timeline
          </h2>

          {/* Vertical connecting line */}
          <div className="absolute left-[27px] top-12 bottom-6 w-0.5 bg-[#e0e3e5]" />

          <div className="space-y-4 pt-1 relative z-10">
            {inspectionPage.timeline.map((entry) => (
              <div key={entry.id} className="flex items-start gap-3">
                {/* Timeline node */}
                <div
                  className={`w-6 h-6 rounded-full flex items-center justify-center shrink-0 border ${
                    entry.isTerminal
                      ? "bg-[#006c49] text-white border-[#006c49]"
                      : "bg-[#e5e9eb] text-[#74777f] border-[#e0e3e5]"
                  }`}
                >
                  <span className="material-symbols-outlined text-[14px]">
                    {entry.isTerminal ? "check" : "schedule"}
                  </span>
                </div>

                <div className="flex-1 bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-3 space-y-0.5">
                  <span className="text-[10px] font-bold text-[#74777f] uppercase tracking-wider block">
                    {entry.category}
                  </span>
                  <h3 className="text-xs font-bold text-[#041627]">{entry.title}</h3>
                  <span className="text-[11px] font-medium text-[#74777f] block pt-0.5">
                    {entry.timestamp}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>

      {/* 10. BOTTOM NAVIGATION */}
      <BottomNavigation activeTab="activity" />
    </div>
  );
};

export default InspectionPage;

