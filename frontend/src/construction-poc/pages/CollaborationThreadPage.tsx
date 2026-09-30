import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { getProjectRecord, getDelivery, getIssue, getThreadEvents, appendThreadEvent } from "../data";
import { actorLabel } from "../data/scenarioMappers";

interface C3EventItem {
  id: string;
  sender: string;
  timestamp: string;
  content: string;
}

/**
 * The canonical scenario thread is an append-only event log, not a chat log:
 * it records `cursor`, `eventType`, `actorId` and `idempotencyKey`, with no
 * wall-clock time and no free-text body. The 23 scenario events are therefore
 * rendered verbatim; only user-posted messages carry generated text/time.
 */
function canonicalEventsToMessages(): C3EventItem[] {
  return getThreadEvents().map((e) => ({
    id: e.idempotencyKey,
    sender: `${actorLabel(e.actorId)} (${e.actorId})`,
    timestamp: `#${e.cursor}`,
    content: e.payload
      ? `${e.eventType} — attachment ${e.payload.cid} (${e.payload.size} bytes)`
      : e.eventType
  }));
}

export const CollaborationThreadPage: React.FC = () => {
  const navigate = useNavigate();
  const project = getProjectRecord();
  const delivery = getDelivery();
  const issue = getIssue();

  const [messages, setMessages] = useState<C3EventItem[]>(canonicalEventsToMessages);

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
          {project.projectReference}
        </button>
        <span>/</span>
        <button onClick={() => navigate("/dashboard/construction/deliveries")} className="hover:underline">
          {delivery.deliveryReference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">C3 Collaboration Thread</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-[#181c1e]">C3 Collaboration Thread</h1>
          <p className="text-sm text-[#44474c]">Evidence &amp; Coordination Log for Issue {issue.issueReference}</p>
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
