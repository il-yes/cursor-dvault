import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getLogisticsOverview, LogisticsOverviewData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";
import { BottomNavigation } from "../components/BottomNavigation";

export const ProvenanceWhyPage: React.FC = () => {
  const navigate = useNavigate();
  const [data, setData] = useState<LogisticsOverviewData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  /* Interactive presentation states */
  const [showDetails, setShowDetails] = useState<boolean>(true);
  const [activeLayerId, setActiveLayerId] = useState<"construction" | "collaboration" | "milestones">("construction");

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
            setError("Provenance overview data unavailable");
          }
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error("Failed to load provenance data:", err);
          setError(err?.message || "Failed to load provenance data");
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
        <p className="text-sm font-semibold text-[#44474c]">Loading provenance reconstruction...</p>
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
          <h2 className="text-lg font-bold text-[#041627] mb-1">Provenance Data Unavailable</h2>
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

  const { delivery, provenance } = data;

  const milestones = provenance?.milestones || [
    {
      stepNumber: 1,
      category: "Physical Cause",
      timestamp: "Aug 15 • 11:15 AM",
      title: "Road restriction on Route M1",
      description: "Emergency overpass load-limit closure published by National Transport Authority.",
      nodeType: "cause" as const
    },
    {
      stepNumber: 2,
      category: "Disruption Ticket",
      timestamp: "Aug 15 • 11:41 AM",
      title: "ISS-1042 — Critical Beam Delivery Delayed",
      description: "Automated delay flag triggered by telematics geo-fence on carrier vehicle.",
      nodeType: "ticket" as const
    },
    {
      stepNumber: 3,
      category: "Verification",
      timestamp: "Aug 15 • 12:05 PM",
      title: "Road Restriction Report (#RD-9942)",
      description: "Verified civil traffic bulletin & real-time detour axle-clearance analysis.",
      nodeType: "verification" as const,
      evidenceChip: {
        filename: "Corridor_M1_Closure_Order.pdf",
        badge: "Signed"
      }
    },
    {
      stepNumber: 4,
      category: "Multi-Org Consensus",
      timestamp: "Aug 15 • 12:45 PM",
      title: "DEC-1042 — Alternative Route B Approved",
      description: "Joint electronic sign-off completed within 44 minutes of issue creation.",
      nodeType: "consensus" as const,
      signatories: ["BuildCorp (PM)", "EuroSteel (Eng)", "FastBuild (Carrier)"]
    },
    {
      stepNumber: 5,
      category: "Field Execution",
      timestamp: "Aug 15 • 01:00 PM",
      title: "Transport Rerouted via Bypass",
      description: "Convoy departed onto regional detour corridor B-88 with highway escort.",
      nodeType: "execution" as const
    },
    {
      stepNumber: 6,
      category: "Site Receipt",
      timestamp: "Aug 16 • 07:32 AM",
      title: "Delivery Received at Site-001",
      description: "Weighbridge scan valid. Transferred to laydown yard sector North-B.",
      nodeType: "receipt" as const
    },
    {
      stepNumber: 7,
      category: "Quality Control",
      timestamp: "Aug 16 • 08:45 AM",
      title: "Inspection #INSP-1042 Passed",
      description: "Flange camber, weld ultrasonic testing, and steel mill certifications verified 100%.",
      nodeType: "qc" as const
    },
    {
      stepNumber: 8,
      category: "Construction Handover",
      timestamp: "Aug 16 • 09:15 AM",
      title: "Material Accepted for Assembly",
      description: "84 structural steel beams officially released to Viaduct Erection Gang 03.",
      nodeType: "handover" as const
    }
  ];

  const auditLayers = provenance?.auditLayers || [
    {
      id: "construction" as const,
      tabLabel: "Construction",
      title: "Physical Site Data & Field Delivery",
      description: "Captured from gate scale telematics, site manager daily logs, and physical acceptance notes signed at Viaduct Sector 4.",
      footerLeft: "Authority: Site Super (BuildCorp)",
      statusRight: "Synced in Realtime"
    },
    {
      id: "collaboration" as const,
      tabLabel: "Collaboration & Evidence",
      title: "Multi-Organization Consensus Protocol",
      description: "3-party signed evidence exchange between BuildCorp, EuroSteel Fabricators, and FastBuild Haulage with route liability approvals.",
      footerLeft: "Threads: 14 exchanged records",
      statusRight: "Consensus 100%"
    },
    {
      id: "milestones" as const,
      tabLabel: "Historical Milestones",
      title: "Historical Milestone Record",
      description: "Chronological audit trail cross-referenced against the master project schedule baseline. Provides transparent, verifiable milestone delivery history.",
      footerLeft: "Ref: TRACE-MILESTONE-001",
      statusRight: "Audited Record"
    }
  ];

  const activeLayer = auditLayers.find((l) => l.id === activeLayerId) || auditLayers[0];

  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: "Provenance / Why? DEL-1042",
        text: `Causal reconstruction for Delivery ${delivery.reference || "DEL-1042"}`,
        url: window.location.href,
      }).catch(() => {});
    } else {
      navigator.clipboard?.writeText(window.location.href);
      alert("Provenance link copied to clipboard");
    }
  };

  const getNodeStyle = (nodeType: string) => {
    switch (nodeType) {
      case "cause":
        return "bg-[#ba1a1a] text-white border-[#ba1a1a]";
      case "ticket":
      case "execution":
        return "bg-[#b76e00] text-white border-[#b76e00]";
      case "verification":
      case "consensus":
        return "bg-[#041627] text-white border-[#041627]";
      case "receipt":
      case "qc":
      case "handover":
      default:
        return "bg-[#006c49] text-white border-[#006c49]";
    }
  };

  return (
    <div className="w-full bg-[#f7fafc] min-h-screen text-[#181c1e] font-[Inter] pb-24">
      {/* 5. Header */}
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
          <h1 className="font-bold text-base tracking-tight text-white">Provenance / Why?</h1>
        </div>
        <div className="flex items-center gap-3">
          <button type="button" className="text-white/80 hover:text-white transition-colors relative">
            <span className="material-symbols-outlined text-[20px]">notifications</span>
            <span className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-[#ffb74d]" />
          </button>
          <div className="w-7 h-7 rounded-full bg-white/20 flex items-center justify-center text-xs font-bold">
            MV
          </div>
        </div>
      </header>

      <main className="w-full max-w-2xl mx-auto px-4 pt-4 flex flex-col gap-4">
        {/* 6. Context Micro-Bar */}
        <div className="flex items-center justify-between bg-white border border-[#e0e3e5] rounded-xl px-4 py-3 shadow-sm">
          <div className="flex items-center gap-2">
            <span className="w-2.5 h-2.5 rounded-full bg-[#b76e00]" />
            <span className="text-xs font-bold text-[#041627]">
              Construction Intelligence • Causal Chain &amp; Provenance
            </span>
          </div>
          <button
            type="button"
            onClick={handleShare}
            className="p-1.5 rounded-lg text-[#74777f] hover:text-[#041627] transition-colors cursor-pointer"
            title="Share Provenance"
          >
            <span className="material-symbols-outlined text-[20px]">share</span>
          </button>
        </div>

        {/* 7. Root-Cause Synthesis Card */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm flex flex-col gap-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-[#b76e00]/10 text-[#b76e00] text-xs font-bold uppercase tracking-wide">
              <span className="material-symbols-outlined text-[16px]">report_problem</span>
              <span>{provenance?.rootCauseTitle || "Root-Cause Diagnosis"}</span>
            </div>
            <span className="px-2.5 py-0.5 rounded-full bg-[#ba1a1a]/10 text-[#ba1a1a] text-xs font-extrabold">
              {provenance?.delayDurationLabel || "+21h 30m Delay"}
            </span>
          </div>

          <p className="text-sm font-semibold text-[#041627] leading-relaxed">
            {provenance?.rootCauseDescription ||
              `Delivery #${delivery.reference || "DEL-1042"} was delayed by 21h 30m due to an unforeseen physical constraint on transport corridor M1, resolved via collaborative consensus on Route B.`}
          </p>

          {/* 8. Key Impact Indicators (3-column compact block) */}
          <div className="grid grid-cols-3 gap-2 bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-3 text-center">
            <div>
              <span className="text-[10px] font-bold text-[#74777f] uppercase block mb-0.5">Original ETA</span>
              <span className="text-xs font-extrabold text-[#041627]">
                {provenance?.originalEtaLabel || "Aug 15, 10:00"}
              </span>
            </div>
            <div className="border-x border-[#e0e3e5] px-1">
              <span className="text-[10px] font-bold text-[#74777f] uppercase block mb-0.5">Actual Site Gate</span>
              <span className="text-xs font-extrabold text-[#041627]">
                {provenance?.actualSiteGateLabel || "Aug 16, 07:32"}
              </span>
            </div>
            <div>
              <span className="text-[10px] font-bold text-[#74777f] uppercase block mb-0.5">Schedule Delta</span>
              <span className="text-xs font-extrabold text-[#006c49]">
                {provenance?.scheduleDeltaLabel || "Absorbed (0d)"}
              </span>
            </div>
          </div>
        </div>

        {/* 9. Causal Progression Header */}
        <div className="flex items-center justify-between pt-1">
          <div>
            <h2 className="text-lg font-bold text-[#041627] tracking-tight">Causal Progression</h2>
            <p className="text-xs text-[#74777f]">8 verifiable milestones recorded across stakeholders</p>
          </div>
          <button
            type="button"
            onClick={() => setShowDetails(!showDetails)}
            className="flex items-center gap-1 px-3 py-1.5 rounded-lg border border-[#e0e3e5] bg-white text-xs font-bold text-[#041627] hover:bg-[#f7fafc] transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[16px]">
              {showDetails ? "visibility_off" : "visibility"}
            </span>
            <span>{showDetails ? "Compact Flow" : "View Details"}</span>
          </button>
        </div>

        {/* 10. 8-Step Vertical Causal Flow */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm relative">
          {/* Vertical connecting line */}
          <div className="absolute left-[27px] top-6 bottom-6 w-0.5 bg-[#e0e3e5]" />

          <div className="space-y-4 relative z-10">
            {milestones.map((m) => (
              <div key={m.stepNumber} className="flex items-start gap-3">
                {/* Circular Node */}
                <div
                  className={`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 border shadow-sm ${getNodeStyle(
                    m.nodeType
                  )}`}
                >
                  {m.stepNumber}
                </div>

                {/* Content Card */}
                <div className="flex-1 bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-3 space-y-1">
                  <div className="flex items-center justify-between text-xs font-bold">
                    <span className="text-[#041627] uppercase tracking-wide text-[11px]">{m.category}</span>
                    <span className="text-[#74777f] text-[11px] font-semibold">{m.timestamp}</span>
                  </div>

                  <h3 className="text-xs font-bold text-[#041627] leading-snug">{m.title}</h3>

                  {showDetails && <p className="text-xs text-[#44474c] leading-relaxed">{m.description}</p>}

                  {/* Milestone 3: Evidence attachment chip */}
                  {m.evidenceChip && (
                    <div className="pt-1 flex items-center gap-2">
                      <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-white border border-[#e0e3e5] text-xs font-semibold text-[#041627]">
                        <span className="material-symbols-outlined text-[14px] text-[#041627]">picture_as_pdf</span>
                        <span>{m.evidenceChip.filename}</span>
                        <span className="px-1.5 py-0.2 rounded bg-[#006c49]/10 text-[#006c49] text-[10px] font-extrabold uppercase ml-1">
                          {m.evidenceChip.badge}
                        </span>
                      </div>
                    </div>
                  )}

                  {/* Milestone 4: Signatories chips */}
                  {m.signatories && m.signatories.length > 0 && (
                    <div className="pt-1 flex flex-wrap gap-1.5">
                      {m.signatories.map((sig, idx) => (
                        <span
                          key={idx}
                          className="px-2 py-0.5 rounded-md bg-white border border-[#e0e3e5] text-[11px] font-bold text-[#041627]"
                        >
                          {sig}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 19. Material Acceptance Card */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-3">
          <div className="flex items-center justify-between">
            <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Material Acceptance</h3>
            <span className="px-2 py-0.5 rounded bg-[#006c49]/10 text-[#006c49] text-xs font-extrabold uppercase">
              Status: {provenance?.materialAcceptance?.status || "Accepted"}
            </span>
          </div>

          <p className="text-xs text-[#74777f] font-medium">
            Material ID: {provenance?.materialAcceptance?.materialCode || "MAT-STRUCT-001"} • Inspection:{" "}
            {provenance?.materialAcceptance?.inspectionCode || "INSP-1042"}
          </p>

          {/* 20. Material Image Area */}
          <div className="relative w-full h-44 rounded-xl overflow-hidden bg-[#041627] flex flex-col justify-between p-4 shadow-inner border border-[#041627]/20">
            {/* Background pattern effect */}
            <div className="absolute inset-0 bg-gradient-to-tr from-[#041627] via-[#0b2b48] to-[#041627] opacity-90" />
            <div className="absolute inset-0 opacity-20 bg-[radial-gradient(#6cf8bb_1px,transparent_1px)] [background-size:16px_16px]" />

            <div className="relative z-10 flex items-center justify-between">
              <span className="px-2.5 py-1 rounded-lg bg-white/10 backdrop-blur-md text-white text-xs font-bold flex items-center gap-1">
                <span className="material-symbols-outlined text-[14px] text-[#6cf8bb]">verified</span>
                <span>Verified Batch</span>
              </span>
              <div className="w-8 h-8 rounded-lg bg-white/10 backdrop-blur-md text-white flex items-center justify-center">
                <span className="material-symbols-outlined text-[20px]">qr_code_2</span>
              </div>
            </div>

            <div className="relative z-10 flex items-end justify-between">
              <div>
                <span className="text-[10px] font-bold text-white/70 uppercase block tracking-wider">
                  Structural Steel Freight
                </span>
                <h4 className="text-sm font-extrabold text-white">
                  {provenance?.materialAcceptance?.lotLabel || "Lot #STM-88219 (84 Beams)"}
                </h4>
              </div>
              <span className="px-2.5 py-1 rounded-lg bg-[#006c49] text-white text-xs font-bold">
                100% Passed
              </span>
            </div>
          </div>
        </div>

        {/* 21. Sources & Audit Integrity */}
        <div className="bg-white border border-[#e0e3e5] rounded-2xl p-4 shadow-sm space-y-3">
          <div>
            <h3 className="text-xs font-bold text-[#041627] uppercase tracking-wider">Sources &amp; Audit Integrity</h3>
            <p className="text-xs text-[#74777f]">Federated evidence underlying this causal resolution</p>
          </div>

          {/* 25. Layer Selector Tabs */}
          <div className="flex bg-[#e5e9eb] p-1 rounded-xl gap-1 text-xs font-bold">
            {auditLayers.map((layer) => (
              <button
                key={layer.id}
                type="button"
                onClick={() => setActiveLayerId(layer.id)}
                className={`flex-1 py-2 rounded-lg transition-all text-center cursor-pointer ${
                  activeLayerId === layer.id
                    ? "bg-[#041627] text-white shadow-sm"
                    : "text-[#44474c] hover:text-[#041627]"
                }`}
              >
                {layer.tabLabel}
              </button>
            ))}
          </div>

          {/* Selected Layer Details Box */}
          <div className="bg-[#f7fafc] border border-[#e0e3e5] rounded-xl p-3 space-y-2">
            <h4 className="text-xs font-bold text-[#041627]">{activeLayer.title}</h4>
            <p className="text-xs text-[#44474c] leading-relaxed">{activeLayer.description}</p>
            <div className="flex items-center justify-between text-[11px] font-semibold text-[#74777f] pt-1 border-t border-[#e0e3e5]">
              <span>{activeLayer.footerLeft}</span>
              <span className="text-[#006c49] font-bold">{activeLayer.statusRight}</span>
            </div>
          </div>
        </div>

        {/* 26. Bottom Actions */}
        <div className="flex flex-col gap-2 pt-2">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.HISTORY)}
            className="w-full py-3.5 bg-[#041627] hover:bg-[#1a2b3c] text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>Open Full Project History (PRJ-001)</span>
            <span className="material-symbols-outlined text-[16px]">arrow_forward</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERY_DETAIL)}
            className="w-full py-3 bg-white border border-[#e0e3e5] hover:bg-[#f7fafc] text-[#041627] text-xs font-bold rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer shadow-sm"
          >
            <span className="material-symbols-outlined text-[16px]">arrow_back</span>
            <span>Back to Delivery ({delivery.reference || "DEL-1042"})</span>
          </button>
        </div>
      </main>

      {/* Fixed Bottom Navigation */}
      <BottomNavigation activeTab="projects" />
    </div>
  );
};

export default ProvenanceWhyPage;

