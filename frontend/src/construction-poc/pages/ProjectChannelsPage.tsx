import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useRoleContext, UserRole } from "../hooks/useRoleContext";
import { getProjectRecord, appendThreadEvent } from "../data";
import { ContextualInviteModal } from "../components/ContextualInviteModal";

export interface ThreadEvent {
  id: string;
  sender: string;
  role: string;
  time: string;
  content: string;
  locationPin?: string;
  photos?: Array<{ id: string; alt: string }>;
  attachments?: Array<{
    id: string;
    filename: string;
    size: string;
    icon: string;
    type: "pdf" | "dwg";
  }>;
  isPrimary?: boolean;
}

export interface ThreadConfig {
  id: string;
  title: string;
  businessObjectId: string;
  businessObjectType: string;
  description: string;
  statusTag?: string;
  statusSubtitle?: string;
  events: ThreadEvent[];
}

export interface ChannelConfig {
  id: string;
  name: string;
  description: string;
  icon: string;
  roles: UserRole[];
  threads: ThreadConfig[];
}

const PREDEFINED_CHANNELS: ChannelConfig[] = [
  {
    id: "general",
    name: "General",
    description: "Project-wide coordination & general updates",
    icon: "forum",
    roles: ["PM", "ARCHITECT", "ENGINEER", "SITE_SUPERINTENDENT", "QA"],
    threads: [
      {
        id: "THREAD-GEN-001",
        title: "PRJ-001 Weekly Operations Alignment",
        businessObjectId: "PRJ-001",
        businessObjectType: "Project",
        description: "General project coordination thread across all lead stakeholders.",
        statusTag: "In Progress",
        statusSubtitle: "Weekly sync meeting agenda",
        events: [
          {
            id: "evt-gen-1",
            sender: "David K.",
            role: "Project Manager",
            time: "09:00",
            content: "Welcome team. Weekly operations review starting. Please review schedule deltas."
          },
          {
            id: "evt-gen-2",
            sender: "Marcus V.",
            role: "Site Superintendent",
            time: "09:15",
            content: "Site delivery gate clear. Foundation excavation Zone A completed."
          }
        ]
      },
      {
        id: "THREAD-GEN-002",
        title: "Site Access & Safety Briefing #14",
        businessObjectId: "PRJ-001",
        businessObjectType: "Project",
        description: "Safety protocols and crane hoisting windows.",
        events: [
          {
            id: "evt-gen-3",
            sender: "Elena R.",
            role: "Safety Officer",
            time: "10:30",
            content: "PPE protocol enforced across Zone A & B. Trench shoring verified."
          }
        ]
      }
    ]
  },
  {
    id: "design",
    name: "Design",
    description: "Architectural definitions & technical specifications",
    icon: "architecture",
    roles: ["PM", "ARCHITECT", "ENGINEER"],
    threads: [
      {
        id: "THREAD-DES-001",
        title: "MAT-STRUCT-001 Precast Beam EN 13369 Technical Specs",
        businessObjectId: "REQ-STRUCT-001",
        businessObjectType: "Requirement",
        description: "C50/60 concrete specification and 12m precast beam requirements.",
        events: [
          {
            id: "evt-des-1",
            sender: "Sarah Jenkins",
            role: "Lead Architect",
            time: "11:00",
            content: "Specification confirmed: C50/60 concrete grade required per EN 13369."
          }
        ]
      },
      {
        id: "THREAD-DES-002",
        title: "Zone A Foundation Structural Plan Review (S-201)",
        businessObjectId: "DOC-S201",
        businessObjectType: "Document",
        description: "Rebar clearance and formwork tolerances per S-201 drawings.",
        events: [
          {
            id: "evt-des-2",
            sender: "Robert Chen",
            role: "Structural Engineer",
            time: "14:20",
            content: "Drawing S-201 updated with rebar spacing tolerances (+/- 5mm)."
          }
        ]
      }
    ]
  },
  {
    id: "procurement",
    name: "Procurement",
    description: "Supplier sourcing, offers & requirement fulfillment",
    icon: "inventory_2",
    roles: ["PM", "ARCHITECT", "SUPPLIER"],
    threads: [
      {
        id: "THREAD-PROC-001",
        title: "REQ-STRUCT-001 Supplier Sourcing & Offer OFF-1042",
        businessObjectId: "REQ-STRUCT-001",
        businessObjectType: "Requirement",
        description: "Procurement thread for precast concrete beams with EuroSteel.",
        events: [
          {
            id: "evt-proc-1",
            sender: "EuroSteel Rep",
            role: "Supplier",
            time: "08:45",
            content: "Offer OFF-1042 submitted: €44,500 total, 12 units C50/60 Precast Beams."
          }
        ]
      }
    ]
  },
  {
    id: "logistics",
    name: "Logistics",
    description: "Freight transport, route constraints & delay management",
    icon: "local_shipping",
    roles: ["PM", "SUPPLIER", "LOGISTICS", "SITE_SUPERINTENDENT"],
    threads: [
      {
        id: "THREAD-LOG-001",
        title: "DEL-1042 Transport TR-1042 Route Restriction",
        businessObjectId: "DEL-1042",
        businessObjectType: "Delivery",
        description: "Logistics coordination for heavy hauler rerouting via Route-B.",
        statusTag: "Route Delayed",
        statusSubtitle: "Bypass Route-B Active",
        events: [
          {
            id: "evt-log-1",
            sender: "FastBuild Logistics",
            role: "Logistics",
            time: "07:30",
            locationPin: "Route M1 Km 42",
            content: "Route M1 blocked due to bridge load restriction. Initiating reroute request."
          },
          {
            id: "evt-log-2",
            sender: "David K.",
            role: "Project Manager",
            time: "09:15",
            content: "Decision DEC-1042 approved Route-B detour. ETA updated to 2026-08-16 07:30."
          }
        ]
      }
    ]
  },
  {
    id: "site-operations",
    name: "Site Operations",
    description: "On-site receiving, equipment hoisting & site Superintendent ops",
    icon: "construction",
    roles: ["PM", "SITE_SUPERINTENDENT", "QA"],
    threads: [
      {
        id: "THREAD-SITE-001",
        title: "North Hub Station Site Receiving",
        businessObjectId: "SITE-HUB-NORTH",
        businessObjectType: "Site",
        description: "Receiving schedule post-reroute.",
        events: [
          {
            id: "evt-site-1",
            sender: "Marcus V.",
            role: "Site Superintendent",
            time: "06:45",
            content: "Crane team prepped for 44t precast beam hoisting at Bay 2."
          }
        ]
      }
    ]
  },
  {
    id: "quality",
    name: "Quality",
    description: "Inspection sign-offs, material tests & compliance",
    icon: "fact_check",
    roles: ["PM", "ENGINEER", "QA", "SITE_SUPERINTENDENT"],
    threads: [
      {
        id: "THREAD-QUAL-001",
        title: "INSP-1042 Precast Concrete Inspection Sign-Off",
        businessObjectId: "INSP-1042",
        businessObjectType: "Inspection",
        description: "QA inspection thread for material acceptance at the hub site.",
        statusTag: "Inspection Passed",
        statusSubtitle: "INSP-1042 Material Acceptance",
        events: [
          {
            id: "evt-qual-1",
            sender: "Bureau Inspection",
            role: "Inspector",
            time: "14:30",
            content: "Receiving inspection passed for MAT-STRUCT-001. Material verified for assembly."
          }
        ]
      },
      {
        id: "THREAD-QUAL-002",
        title: "Foundation Concrete Inspection",
        businessObjectId: "INSP-FOUNDATION-A",
        businessObjectType: "Inspection",
        description: "Thread ID: #FI-2023-04-12 — Foundation Zone A pour sign-off.",
        statusTag: "Awaiting Approval",
        statusSubtitle: "Foundation concrete inspection",
        events: [
          {
            id: "evt-fi-1",
            sender: "Contractor",
            role: "Contractor",
            time: "08:10",
            content: "Inspection requested for pour zone A3. Ready for sign-off."
          },
          {
            id: "evt-fi-2",
            sender: "Inspector",
            role: "Inspector",
            time: "08:24",
            locationPin: "On site. Commencing inspection.",
            content: "Arrived at Foundation Zone A. Initial formwork check underway."
          },
          {
            id: "evt-fi-3",
            sender: "Inspector",
            role: "Inspector",
            time: "09:02",
            content: "Concrete samples recorded. Slump test within tolerance (4.5 inches).",
            photos: [
              { id: "p1", alt: "Slump Test Cone 4.5 in" },
              { id: "p2", alt: "Rebar Grid Inspection" }
            ]
          },
          {
            id: "evt-fi-4",
            sender: "Inspector",
            role: "Inspector",
            time: "09:40",
            isPrimary: true,
            content: "Inspection report generated and uploaded. Awaiting structural engineer review.",
            attachments: [
              { id: "a1", filename: "Foundation_Report.pdf", size: "2.4 MB • Uploaded 09:40", icon: "picture_as_pdf", type: "pdf" },
              { id: "a2", filename: "Drawing_S-204.pdf", size: "5.1 MB • Attached ref", icon: "architecture", type: "dwg" }
            ]
          }
        ]
      }
    ]
  },
  {
    id: "decisions",
    name: "Decisions",
    description: "Formal decision approvals & action execution",
    icon: "gavel",
    roles: ["PM", "ARCHITECT", "ENGINEER", "SUPPLIER", "LOGISTICS", "SITE_SUPERINTENDENT", "QA"],
    threads: [
      {
        id: "THREAD-DEC-001",
        title: "DEC-1042 Route-B Reroute Approval",
        businessObjectId: "DEC-1042",
        businessObjectType: "Decision",
        description: "Decision approval thread for alternative route execution.",
        statusTag: "Approved",
        statusSubtitle: "Authorized by Project Manager",
        events: [
          {
            id: "evt-dec-1",
            sender: "David K.",
            role: "Project Manager",
            time: "12:00",
            content: "Decision DEC-1042 finalized: Approve Route B alternative transport detour."
          }
        ]
      }
    ]
  }
];

export const ProjectChannelsPage: React.FC = () => {
  const navigate = useNavigate();
  const { channelId } = useParams<{ channelId?: string }>();
  const { activeRole, roleConfig } = useRoleContext();
  const project = getProjectRecord();

  const visibleChannels = PREDEFINED_CHANNELS.filter((c) => c.roles.includes(activeRole));
  const activeChannel = PREDEFINED_CHANNELS.find((c) => c.id === channelId) || visibleChannels[0] || PREDEFINED_CHANNELS[0];

  const [selectedThreadIdState, setSelectedThreadIdState] = useState<string | null>(null);

  const activeThread = activeChannel.threads.find((t) => t.id === selectedThreadIdState) || activeChannel.threads[0];

  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
  const [threadEventsState, setThreadEventsState] = useState<Record<string, ThreadEvent[]>>({});
  const [newEventText, setNewEventText] = useState("");

  const currentThreadEvents = threadEventsState[activeThread.id] || activeThread.events;

  const handlePostEvent = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newEventText.trim()) return;

    const newEvt: ThreadEvent = {
      id: `EVT-${Date.now()}`,
      sender: roleConfig.label,
      role: roleConfig.label,
      time: new Date().toISOString().replace("T", " ").substring(11, 16),
      content: newEventText.trim()
    };

    setThreadEventsState((prev) => ({
      ...prev,
      [activeThread.id]: [...(prev[activeThread.id] || activeThread.events), newEvt]
    }));
    setNewEventText("");

    try {
      await appendThreadEvent(activeThread.id, "COMMENT", newEvt.content);
    } catch (err) {
      console.warn("Thread event append handled locally:", err);
    }
  };

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 gap-6">
      {/* Breadcrumb Header */}
      <div className="flex items-center gap-2 text-xs text-[#44474c]">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.projectReference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Project Channels</span>
      </div>

      {/* Main Channel Topology Header */}
      <div className="bg-[#041627] text-white rounded-xl p-6 shadow-md flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <span className="text-xs font-mono text-[#6cf8bb] uppercase tracking-wider font-semibold block mb-1">
            PROJECT COLLABORATION TOPOLOGY
          </span>
          <h1 className="text-2xl font-bold text-white mb-1">Channels &amp; Threads</h1>
          <p className="text-xs text-[#b7c8de]">
            Role-scoped channels &amp; business object collaboration threads.
          </p>
        </div>
        <span className="px-3 py-1 bg-[#6cf8bb]/20 text-[#6cf8bb] text-xs font-bold rounded">
          Role: {roleConfig.label} ({visibleChannels.length} Visible Channels)
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Left Sidebar: Predefined Channel Topology (Restored Channel List) */}
        <div className="flex flex-col gap-2">
          <span className="text-xs font-bold text-[#44474c] uppercase tracking-wider px-1">
            Predefined Channels
          </span>

          {PREDEFINED_CHANNELS.map((ch) => {
            const isVisible = ch.roles.includes(activeRole);
            const isSelected = ch.id === activeChannel.id;

            return (
              <button
                key={ch.id}
                type="button"
                onClick={() => {
                  navigate(`/dashboard/construction/channels/${ch.id}`);
                  setSelectedThreadIdState(null);
                }}
                className={`p-3 rounded-xl border text-left transition-all flex items-center justify-between cursor-pointer ${
                  isSelected
                    ? "bg-[#041627] text-white border-[#041627] shadow-sm"
                    : isVisible
                    ? "bg-white text-[#181c1e] border-[#e0e3e5] hover:border-[#041627]"
                    : "bg-[#f0f3f5] text-[#a0a5aa] border-transparent opacity-60"
                }`}
              >
                <div className="flex items-center gap-2.5">
                  <span className="material-symbols-outlined text-[20px]">{ch.icon}</span>
                  <div>
                    <h4 className="font-bold text-xs">{ch.name}</h4>
                    <p className={`text-[10px] ${isSelected ? "text-[#b7c8de]" : "text-[#44474c]"}`}>
                      {ch.threads.length} Thread{ch.threads.length !== 1 ? "s" : ""}
                    </p>
                  </div>
                </div>
                {!isVisible && (
                  <span className="text-[9px] bg-[#e0e3e5] text-[#44474c] px-1.5 py-0.5 rounded font-mono">
                    Restricted
                  </span>
                )}
              </button>
            );
          })}
        </div>

        {/* Right Main Area: Active Channel Banner, Native Thread Select & Stitch Event Feed */}
        <div className="md:col-span-2 flex flex-col gap-4">
          {/* Active Channel Info Card & Contextual Invite Action */}
          <div className="bg-white border border-[#e0e3e5] rounded-xl p-5 shadow-sm">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <span className="material-symbols-outlined text-[#041627]">{activeChannel.icon}</span>
                <h3 className="font-bold text-[#181c1e] text-base">{activeChannel.name} Channel</h3>
              </div>

              {/* Contextual Invite Button */}
              <button
                type="button"
                onClick={() => setIsInviteModalOpen(true)}
                className="px-3 py-1.5 bg-[#041627] text-white font-bold text-xs rounded-lg hover:bg-[#1a2b3c] transition-colors flex items-center gap-1 shadow-sm cursor-pointer"
              >
                <span className="material-symbols-outlined text-[16px]">person_add</span>
                Invite Stakeholder
              </button>
            </div>
            <p className="text-xs text-[#44474c]">{activeChannel.description}</p>
          </div>

          {/* Native HTML <select> for Thread Selection */}
          <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <div className="flex-1 w-full flex items-center gap-2">
              <label htmlFor="thread-select" className="text-xs font-bold text-[#041627] uppercase tracking-wider shrink-0">
                Thread:
              </label>
              <select
                id="thread-select"
                aria-label="Select Thread"
                value={activeThread.id}
                onChange={(e) => setSelectedThreadIdState(e.target.value)}
                className="w-full bg-[#f7fafc] border border-[#e0e3e5] text-[#041627] text-xs font-bold rounded-lg px-3 py-2 outline-none focus:border-[#041627] cursor-pointer"
              >
                {activeChannel.threads.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.title} ({t.businessObjectId})
                  </option>
                ))}
              </select>
            </div>
            <span className="text-[11px] text-[#74777f] font-mono shrink-0">
              {activeChannel.threads.length} Thread{activeChannel.threads.length !== 1 ? "s" : ""} in {activeChannel.name}
            </span>
          </div>

          {/* Stitch Thread Presentation Container (Without Header) */}
          <div className="flex flex-col gap-4">
            {/* Status Ribbon (if statusTag exists) */}
            {activeThread.statusTag && (
              <div className="bg-[#ffddb8] text-[#653e00] px-4 py-2.5 rounded-xl flex items-center justify-between shadow-sm">
                <div className="flex items-center gap-1.5">
                  <span className="material-symbols-outlined text-[18px]">warning</span>
                  <span className="text-xs font-bold uppercase tracking-wider">{activeThread.statusTag}</span>
                </div>
                <span className="text-[11px] font-semibold opacity-90">{activeThread.statusSubtitle || activeThread.title}</span>
              </div>
            )}

            {/* Selected Thread Meta & Participants Stack */}
            <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 shadow-sm space-y-2">
              <h2 className="text-lg font-bold text-[#041627] tracking-tight">{activeThread.title}</h2>
              <p className="text-xs text-[#74777f] font-mono">Thread ID: #{activeThread.id} • Object: {activeThread.businessObjectId} ({activeThread.businessObjectType})</p>
              <div className="flex -space-x-2 pt-1">
                <div className="w-8 h-8 rounded-full border-2 border-white bg-[#041627] text-white flex items-center justify-center text-xs font-bold z-30" title="Engineer">EN</div>
                <div className="w-8 h-8 rounded-full border-2 border-white bg-[#006c49] text-white flex items-center justify-center text-xs font-bold z-20" title="Contractor">CO</div>
                <div className="w-8 h-8 rounded-full border-2 border-white bg-[#221200] text-white flex items-center justify-center text-xs font-bold z-10" title="Inspector">IN</div>
              </div>
            </div>

            {/* Stitch Thread Event Feed */}
            <div className="bg-white border border-[#e0e3e5] rounded-xl p-5 shadow-sm space-y-4">
              <div className="relative pl-7 pb-4">
                {/* Vertical connecting line */}
                <div className="absolute left-[11px] top-2 bottom-6 w-[2px] bg-[#e0e3e5]" />

                <div className="space-y-6">
                  {currentThreadEvents.map((evt) => (
                    <div key={evt.id} className="relative">
                      {/* Timeline dot node */}
                      <div className={`absolute -left-[31px] mt-1.5 w-4 h-4 rounded-full border-4 border-white ${evt.isPrimary ? "bg-[#041627]" : "bg-[#74777d]"}`} />
                      
                      <div className="flex flex-col gap-1">
                        <div className="flex items-center gap-2 text-xs">
                          <span className="font-bold text-[#041627]">{evt.role || evt.sender}</span>
                          <span className="text-[11px] text-[#74777f] font-mono">{evt.time}</span>
                        </div>
                        
                        <div className={`rounded-xl rounded-tl-none p-3.5 shadow-sm border ${
                          evt.isPrimary
                            ? "bg-[#041627] text-white border-[#041627]"
                            : "bg-[#f7fafc] text-[#041627] border-[#e0e3e5]"
                        }`}>
                          {evt.locationPin && (
                            <div className="flex items-center gap-1.5 mb-1.5 text-xs text-[#006c49] font-semibold">
                              <span className="material-symbols-outlined text-[16px]">location_on</span>
                              <span>{evt.locationPin}</span>
                            </div>
                          )}

                          <p className="text-xs leading-relaxed">{evt.content}</p>

                          {/* Photos Grid */}
                          {evt.photos && evt.photos.length > 0 && (
                            <div className="flex gap-2 overflow-x-auto pt-3 pb-1">
                              {evt.photos.map((p) => (
                                <div key={p.id} className="min-w-[120px] h-[80px] rounded-lg bg-[#041627] overflow-hidden relative border border-[#e0e3e5] flex flex-col justify-end p-1.5">
                                  <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent" />
                                  <span className="relative z-10 text-[10px] font-bold text-white leading-tight">{p.alt}</span>
                                </div>
                              ))}
                            </div>
                          )}

                          {/* Attachments (PDF / DWG) */}
                          {evt.attachments && evt.attachments.length > 0 && (
                            <div className="flex flex-col gap-2 pt-3">
                              {evt.attachments.map((att) => (
                                <div key={att.id} className="flex items-center gap-3 p-2.5 bg-white rounded-lg border border-[#e0e3e5] text-[#041627]">
                                  <div className={`w-9 h-9 rounded flex items-center justify-center shrink-0 ${att.type === "pdf" ? "bg-[#ffdad6] text-[#ba1a1a]" : "bg-[#6cf8bb]/20 text-[#006c49]"}`}>
                                    <span className="material-symbols-outlined text-[20px]">
                                      {att.type === "pdf" ? "picture_as_pdf" : "architecture"}
                                    </span>
                                  </div>
                                  <div className="flex-1 min-w-0">
                                    <p className="text-xs font-bold truncate">{att.filename}</p>
                                    <p className="text-[10px] text-[#74777f]">{att.size}</p>
                                  </div>
                                  <button type="button" className="p-1 text-[#44474c] hover:text-[#041627] cursor-pointer" title="Download">
                                    <span className="material-symbols-outlined text-[18px]">download</span>
                                  </button>
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Bottom Pulse Dot */}
                <div className="absolute -left-[30px] bottom-0 w-4 h-4 rounded-full bg-white border-4 border-[#e0e3e5] flex items-center justify-center">
                  <span className="material-symbols-outlined text-[#74777f] text-[14px] animate-pulse">more_horiz</span>
                </div>
              </div>

              {/* Post Event Input Form */}
              <form onSubmit={handlePostEvent} className="flex gap-2 pt-2 border-t border-[#f0f3f5]">
                <input
                  type="text"
                  value={newEventText}
                  onChange={(e) => setNewEventText(e.target.value)}
                  placeholder={`Post update to ${activeThread.title}...`}
                  className="flex-1 px-3.5 py-2.5 bg-[#f7fafc] border border-[#e0e3e5] rounded-xl text-xs text-[#041627] outline-none focus:border-[#041627]"
                />
                <button
                  type="submit"
                  className="px-4 py-2.5 bg-[#041627] text-white font-bold text-xs rounded-xl hover:bg-[#1a2b3c] transition-colors cursor-pointer"
                >
                  Post Event
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>

      {/* Contextual Invitation Modal */}
      <ContextualInviteModal
        isOpen={isInviteModalOpen}
        onClose={() => setIsInviteModalOpen(false)}
        workspaceId={project.projectId}
        workspaceName={project.projectName}
        channelId={activeChannel.name}
        channelName={activeChannel.name}
        threadId={activeThread.id}
        threadTitle={activeThread.title}
        defaultRole={activeRole === "PM" ? "SUPPLIER" : activeRole}
      />
    </div>
  );
};

export default ProjectChannelsPage;

