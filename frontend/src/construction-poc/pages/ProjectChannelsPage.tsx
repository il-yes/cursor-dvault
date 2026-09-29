import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useRoleContext, UserRole } from "../hooks/useRoleContext";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";
import { appendThreadEvent } from "../data";
import { ContextualInviteModal } from "../components/ContextualInviteModal";

interface ChannelConfig {
  id: string;
  name: string;
  description: string;
  icon: string;
  roles: UserRole[];
  threads: Array<{
    id: string;
    title: string;
    businessObjectId: string;
    businessObjectType: string;
    description: string;
  }>;
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
        description: "General project coordination thread across all lead stakeholders."
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
        title: "W18x86 Beam ASTM A992 Technical Specs",
        businessObjectId: "REQ-STRUCT-001",
        businessObjectType: "Requirement",
        description: "Architectural steel grade and corrosion coating specs."
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
        description: "Procurement thread for structural steel beams with Apex Steel."
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
        title: "DEL-1042 Transport TR-1042 Route M1 Restriction",
        businessObjectId: "DEL-1042",
        businessObjectType: "Delivery",
        description: "Logistics coordination for heavy hauler rerouting via Highway B."
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
        title: "Site Alpha South Pier Framework Receiving",
        businessObjectId: "SITE-SOUTH-01",
        businessObjectType: "Site",
        description: "Pier 4 receiving schedule post-reroute."
      }
    ]
  },
  {
    id: "quality",
    name: "Quality",
    description: "Inspection sign-offs, ultrasonic weld tests & compliance",
    icon: "fact_check",
    roles: ["PM", "ENGINEER", "QA", "SITE_SUPERINTENDENT"],
    threads: [
      {
        id: "THREAD-QUAL-001",
        title: "INSP-1042 Structural Steel Weld Inspection Sign-Off",
        businessObjectId: "INSP-1042",
        businessObjectType: "Inspection",
        description: "QA inspection thread for Pier 4 steel acceptance."
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
        title: "DEC-1042 Highway B Reroute Approval",
        businessObjectId: "DEC-1042",
        businessObjectType: "Decision",
        description: "Decision approval thread for alternative route B execution."
      }
    ]
  }
];

export const ProjectChannelsPage: React.FC = () => {
  const navigate = useNavigate();
  const { channelId, threadId } = useParams<{ channelId?: string; threadId?: string }>();
  const { activeRole, roleConfig } = useRoleContext();
  const { project } = SCENARIO_DATA;

  // Filter channels visible to current active role
  const visibleChannels = PREDEFINED_CHANNELS.filter(c => c.roles.includes(activeRole));
  
  // Selected channel or fallback to first visible channel
  const activeChannel = PREDEFINED_CHANNELS.find(c => c.id === channelId) || visibleChannels[0] || PREDEFINED_CHANNELS[0];
  
  // Selected thread or first thread in channel
  const activeThread = activeChannel.threads.find(t => t.id === threadId) || activeChannel.threads[0];

  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);

  // Local thread event state
  const [events, setEvents] = useState([
    {
      id: "EVT-01",
      sender: "Mark Vance (Logistics Operator)",
      time: "2026-08-15 11:41",
      content: "Road restriction reported on M1 Km 42 bridge structure. Weight limit reduced to 35T."
    },
    {
      id: "EVT-02",
      sender: "David Chen (Supplier / Apex Steel)",
      time: "2026-08-15 12:05",
      content: "Delivery DEL-1042 ETA updated to 2026-08-16 07:30. Proposing alternative Route B."
    },
    {
      id: "EVT-03",
      sender: "Alex Rivera (Project Manager)",
      time: "2026-08-15 14:20",
      content: "Decision DEC-1042 approved. Transport TR-1042 rerouted via Highway B."
    }
  ]);
  const [newEventText, setNewEventText] = useState("");

  const handlePostEvent = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newEventText.trim()) return;

    const newEvt = {
      id: `EVT-${Date.now()}`,
      sender: `Current User (${roleConfig.label})`,
      time: new Date().toISOString().replace("T", " ").substring(0, 16),
      content: newEventText.trim()
    };

    setEvents(prev => [...prev, newEvt]);
    setNewEventText("");

    if (activeThread) {
      try {
        await appendThreadEvent(activeThread.id, "COMMENT", newEvt.content);
      } catch (err) {
        console.warn("Thread event append handled locally:", err);
      }
    }
  };

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6 gap-6">
      {/* Breadcrumb Header */}
      <div className="flex items-center gap-2 text-xs text-[#44474c]">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
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
        {/* Left Sidebar: Predefined Channel Topology */}
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
                onClick={() => navigate(`/dashboard/construction/channels/${ch.id}`)}
                className={`p-3 rounded-xl border text-left transition-all flex items-center justify-between ${
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

        {/* Right Main Area: Active Channel Threads & Event Feed */}
        <div className="md:col-span-2 flex flex-col gap-4">
          {/* Active Channel Info & Contextual Invite Action */}
          <div className="bg-white border border-[#e0e3e5] rounded-xl p-5 shadow-sm">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <span className="material-symbols-outlined text-[#041627]">{activeChannel.icon}</span>
                <h3 className="font-bold text-[#181c1e] text-base">{activeChannel.name} Channel</h3>
              </div>

              {/* Contextual Invite Button */}
              <button
                onClick={() => setIsInviteModalOpen(true)}
                className="px-3 py-1.5 bg-[#041627] text-white font-bold text-xs rounded-lg hover:bg-[#041627]/90 transition-colors flex items-center gap-1 shadow-sm"
              >
                <span className="material-symbols-outlined text-[16px]">person_add</span>
                Invite Stakeholder
              </button>
            </div>
            <p className="text-xs text-[#44474c] mb-3">{activeChannel.description}</p>

            {/* Active Thread Bar */}
            {activeThread && (
              <div className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
                <div className="flex items-center justify-between text-xs mb-1">
                  <span className="font-bold text-[#006c49] flex items-center gap-1">
                    <span className="material-symbols-outlined text-[16px]">topic</span>
                    Thread: {activeThread.title}
                  </span>
                  <span className="font-mono text-[#44474c] text-[10px]">
                    Object: {activeThread.businessObjectId} ({activeThread.businessObjectType})
                  </span>
                </div>
                <p className="text-xs text-[#44474c]">{activeThread.description}</p>
              </div>
            )}
          </div>

          {/* Thread Events List */}
          <div className="bg-white border border-[#e0e3e5] rounded-xl p-5 shadow-sm flex flex-col gap-4">
            <h4 className="text-xs font-bold text-[#44474c] uppercase tracking-wider border-b border-[#f0f3f5] pb-2">
              Thread Events &amp; Evidence Messages
            </h4>

            <div className="flex flex-col gap-3">
              {events.map((evt) => (
                <div key={evt.id} className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-xs">
                  <div className="flex justify-between items-center mb-1">
                    <span className="font-bold text-[#181c1e]">{evt.sender}</span>
                    <span className="text-[#44474c] font-mono text-[10px]">{evt.time}</span>
                  </div>
                  <p className="text-[#44474c] leading-relaxed">{evt.content}</p>
                </div>
              ))}
            </div>

            {/* Post Event Input */}
            <form onSubmit={handlePostEvent} className="flex gap-2 pt-2 border-t border-[#f0f3f5]">
              <input
                type="text"
                value={newEventText}
                onChange={(e) => setNewEventText(e.target.value)}
                placeholder={`Post update to ${activeThread?.title || "thread"}...`}
                className="flex-1 px-3 py-2 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-xs text-[#181c1e] outline-none focus:border-[#041627]"
              />
              <button
                type="submit"
                className="px-4 py-2 bg-[#041627] text-white font-bold text-xs rounded-lg hover:bg-[#041627]/90"
              >
                Post Event
              </button>
            </form>
          </div>
        </div>
      </div>

      {/* Contextual Invitation Modal */}
      <ContextualInviteModal
        isOpen={isInviteModalOpen}
        onClose={() => setIsInviteModalOpen(false)}
        workspaceId={project.code}
        workspaceName={project.name}
        channelId={activeChannel.name}
        channelName={activeChannel.name}
        threadId={activeThread?.id}
        threadTitle={activeThread?.title}
        defaultRole={activeRole === "PM" ? "SUPPLIER" : activeRole}
      />
    </div>
  );
};

export default ProjectChannelsPage;
