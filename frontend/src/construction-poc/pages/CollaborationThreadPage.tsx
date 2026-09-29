import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { SCENARIO_DATA } from "../data/constructionScenarioAdapter";
import { appendThreadEvent } from "../data";

interface C3EventItem {
  id: string;
  sender: string;
  timestamp: string;
  content: string;
}

export const CollaborationThreadPage: React.FC = () => {
  const navigate = useNavigate();
  const { project, delivery, issue } = SCENARIO_DATA;

  const [messages, setMessages] = useState<C3EventItem[]>([
    {
      id: "MSG-001",
      sender: "Mark Vance (Freight Operator)",
      timestamp: "2026-08-15 11:41",
      content: "Road restriction encountered on M1 Km 42 bridge. Weight limit reduced to 35T. TR-1042 load is 48T."
    },
    {
      id: "MSG-002",
      sender: "David Chen (Logistics)",
      timestamp: "2026-08-15 12:05",
      content: "Delivery DEL-1042 ETA moved to Aug 16, 07:30. Preparing alternative route proposal via Highway B."
    },
    {
      id: "MSG-003",
      sender: "Sarah Jenkins (Site Superintendent)",
      timestamp: "2026-08-15 12:17",
      content: "Acknowledged. Foundation crew rescheduled for morning arrival on Aug 16."
    },
    {
      id: "MSG-004",
      sender: "Alex Rivera (Project Manager)",
      timestamp: "2026-08-15 14:20",
      content: "Decision DEC-1042 approved. Reroute via Highway B with secondary police escort confirmed."
    }
  ]);

  const [newMessage, setNewMessage] = useState("");

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    const newMsgObj: C3EventItem = {
      id: `MSG-${Date.now()}`,
      sender: "Current User",
      timestamp: new Date().toISOString().replace("T", " ").substring(0, 16),
      content: newMessage.trim()
    };

    setMessages(prev => [...prev, newMsgObj]);
    setNewMessage("");

    // Persist through the data boundary; the active provider decides the target.
    try {
      await appendThreadEvent("DEL-1042-THREAD", "COMMENT", newMsgObj.content);
    } catch (err) {
      console.warn("Thread event append handled locally:", err);
    }
  };

  return (
    <div className="flex flex-col w-full pb-8 max-w-4xl mx-auto px-6 pt-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#44474c] mb-4">
        <button onClick={() => navigate("/dashboard/construction/projects")} className="hover:underline">
          {project.code}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/deliveries")} className="hover:underline">
          {delivery.reference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">C3 Collaboration Thread</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">C3 Collaboration Thread</h1>
          <p className="text-sm text-[#44474c]">Evidence &amp; Coordination Log for Issue {issue.reference}</p>
        </div>
        <span className="px-3 py-1 bg-[#041627] text-white rounded-full text-xs font-bold">
          C3 ENCRYPTED
        </span>
      </div>

      {/* Messages Feed */}
      <div className="bg-white border border-[#e0e3e5] rounded-xl p-6 shadow-sm mb-6 flex flex-col gap-4">
        {messages.map((msg) => (
          <div key={msg.id} className="p-4 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg">
            <div className="flex items-center justify-between mb-1">
              <span className="font-bold text-[#181c1e] text-sm">{msg.sender}</span>
              <span className="text-xs text-[#44474c]">{msg.timestamp}</span>
            </div>
            <p className="text-sm text-[#44474c] leading-relaxed">{msg.content}</p>
          </div>
        ))}

        {/* Send Input */}
        <form onSubmit={handleSendMessage} className="flex gap-2 pt-2 border-t border-[#f0f3f5]">
          <input
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Post evidence update or comment to C3 thread..."
            className="flex-1 px-4 py-2.5 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg text-sm text-[#181c1e] focus:outline-none focus:border-[#041627]"
          />
          <button
            type="submit"
            className="px-5 py-2.5 bg-[#041627] text-white font-semibold text-sm rounded-lg hover:bg-[#041627]/90 transition-colors"
          >
            Post Event
          </button>
        </form>
      </div>
    </div>
  );
};

export default CollaborationThreadPage;
