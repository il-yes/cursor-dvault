import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { getNotificationsData, NotificationCategory } from "../data";

export const NotificationsPage: React.FC = () => {
  const navigate = useNavigate();
  const data = getNotificationsData();

  const [activeFilter, setActiveFilter] = useState<"all" | NotificationCategory>("all");
  const [markedRead, setMarkedRead] = useState<boolean>(false);

  const filteredNotifications = data.notifications.filter((n) => {
    if (activeFilter === "all") return true;
    return n.category === activeFilter;
  });

  const filterTabs: Array<{ id: "all" | NotificationCategory; label: string; count: number }> = [
    { id: "all", label: "All", count: data.counts.all },
    { id: "critical", label: "Critical", count: data.counts.critical },
    { id: "operations", label: "Operations", count: data.counts.operations },
    { id: "decisions", label: "Decisions", count: data.counts.decisions },
    { id: "collaboration", label: "Collaboration", count: data.counts.collaboration },
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-[#e0e3e5] pb-4">
        <div>
          <h1 className="text-2xl font-bold text-[#041627]">Notifications</h1>
          <p className="text-xs text-[#44474c] font-medium mt-0.5">
            Operational alerts &amp; project coordination
          </p>
        </div>

        <button
          onClick={() => setMarkedRead(true)}
          disabled={markedRead}
          className={`px-4 py-2 text-xs font-bold rounded-lg transition-colors shadow-sm flex items-center gap-1.5 self-start sm:self-auto ${
            markedRead
              ? "bg-[#e0e3e5] text-[#707478] cursor-not-allowed"
              : "bg-white text-[#041627] border border-[#d0d4d8] hover:bg-[#f0f4f8]"
          }`}
        >
          <span className="material-symbols-outlined text-[16px]">
            {markedRead ? "check_circle" : "done_all"}
          </span>
          {markedRead ? "✓ Caught up!" : "Mark all read"}
        </button>
      </div>

      {/* Operational Health Pulse Card */}
      <div className="p-4 rounded-xl border border-[#e0e3e5] bg-white shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="relative flex items-center justify-center">
            <span className="w-3 h-3 rounded-full bg-emerald-500 animate-ping absolute" />
            <span className="w-3 h-3 rounded-full bg-emerald-500 relative" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h3 className="text-sm font-bold text-[#041627]">{data.healthPulse.statusText}</h3>
              <span className="px-2 py-0.5 bg-emerald-500/10 text-emerald-700 text-[10px] font-bold rounded">
                ACTIVE
              </span>
            </div>
            <p className="text-xs text-[#44474c] font-medium mt-0.5">
              {data.healthPulse.detailText}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-1.5 text-xs text-[#44474c] font-medium shrink-0 self-end sm:self-auto">
          <span className="material-symbols-outlined text-[16px] text-[#006c49]">sync</span>
          <span>{data.healthPulse.syncAgeText}</span>
        </div>
      </div>

      {/* Filter Tabs */}
      <div className="flex items-center gap-2 border-b border-[#e0e3e5] pb-2 overflow-x-auto no-scrollbar">
        {filterTabs.map((tab) => {
          const isActive = activeFilter === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveFilter(tab.id)}
              className={`px-3.5 py-1.5 rounded-lg text-xs font-bold transition-colors whitespace-nowrap flex items-center gap-1.5 ${
                isActive
                  ? "bg-[#041627] text-white shadow-sm"
                  : "bg-white text-[#44474c] border border-[#e0e3e5] hover:bg-[#f0f4f8] hover:text-[#041627]"
              }`}
            >
              <span>{tab.label}</span>
              <span
                className={`px-1.5 py-0.2 text-[10px] rounded-full font-extrabold ${
                  isActive ? "bg-white/20 text-white" : "bg-[#f0f4f8] text-[#041627]"
                }`}
              >
                {tab.count}
              </span>
            </button>
          );
        })}
      </div>

      {/* Notification Feed */}
      <div className="space-y-4">
        {filteredNotifications.map((n) => {
          const isUnread = !markedRead && n.isUnread;

          return (
            <div
              key={n.id}
              className={`p-5 rounded-xl bg-white border shadow-sm space-y-3 transition-all ${
                n.isCritical
                  ? "border-l-4 border-l-red-500 border-t-[#e0e3e5] border-r-[#e0e3e5] border-b-[#e0e3e5]"
                  : "border-[#e0e3e5]"
              }`}
            >
              {/* Card Top Strip: Badge, Unread Dot, Timestamp */}
              <div className="flex items-center justify-between gap-2">
                <div className="flex items-center gap-2">
                  {isUnread && (
                    <span className="w-2.5 h-2.5 rounded-full bg-blue-600 shadow-sm shrink-0" title="Unread" />
                  )}
                  <span className={`px-2.5 py-0.5 text-[10px] font-extrabold rounded border tracking-wider uppercase ${n.badgeStyle}`}>
                    {n.badge}
                  </span>
                </div>
                <span className="text-xs text-[#44474c] font-medium">{n.timestamp}</span>
              </div>

              {/* Title & Body */}
              <div>
                <h3 className="text-base font-bold text-[#041627] flex items-center gap-2">
                  {n.isCritical && (
                    <span className="material-symbols-outlined text-red-500 text-[18px]">warning</span>
                  )}
                  {n.title}
                </h3>
                <p className="text-xs text-[#44474c] leading-relaxed mt-1">{n.body}</p>
              </div>

              {/* Embedded Widgets */}
              {n.widget && n.widget.type === "route" && (
                <div className="p-3 bg-red-50 border border-red-200 rounded-lg flex items-center justify-between text-xs font-bold text-red-900">
                  <div className="flex items-center gap-2">
                    <span className="material-symbols-outlined text-red-600 text-[18px]">wrong_location</span>
                    <span>{n.widget.title}</span>
                  </div>
                  <span className="px-2 py-0.5 bg-red-600 text-white text-[10px] font-extrabold rounded">
                    {n.widget.subtitle}
                  </span>
                </div>
              )}

              {n.widget && n.widget.type === "consensus" && (
                <div className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs">
                  <div className="flex items-center gap-1.5">
                    {n.widget.chips?.map((chip, idx) => (
                      <span key={idx} className="w-6 h-6 rounded-full bg-[#041627] text-white text-[10px] font-bold flex items-center justify-center">
                        {chip}
                      </span>
                    ))}
                    <span className="font-bold text-[#041627] ml-2">{n.widget.title}</span>
                  </div>
                  <span className="text-xs font-bold text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200 self-start sm:self-auto">
                    {n.widget.subtitle}
                  </span>
                </div>
              )}

              {n.widget && n.widget.type === "metrics" && (
                <div className="p-3 bg-emerald-50 border border-emerald-200 rounded-lg flex items-center justify-between text-xs font-bold text-emerald-900">
                  <div className="flex items-center gap-2">
                    <span className="material-symbols-outlined text-emerald-600 text-[18px]">fact_check</span>
                    <span>{n.widget.title}</span>
                  </div>
                  <span className="px-2 py-0.5 bg-emerald-700 text-white text-[10px] font-extrabold rounded">
                    {n.widget.subtitle}
                  </span>
                </div>
              )}

              {n.widget && n.widget.type === "attachment" && (
                <div className="p-3 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg flex items-center gap-3 text-xs">
                  <span className="material-symbols-outlined text-[#041627] text-[22px]">picture_as_pdf</span>
                  <div>
                    <p className="font-bold text-[#041627]">{n.widget.filename}</p>
                    <p className="text-[11px] text-[#44474c]">{n.widget.filesize}</p>
                  </div>
                </div>
              )}

              {/* Card Footer: Metadata & Action Buttons */}
              <div className="pt-2 border-t border-[#e0e3e5] flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                <span className="text-xs font-semibold text-[#44474c]">{n.metadataText}</span>

                <div className="flex flex-wrap items-center gap-2">
                  {n.actions.map((act, idx) => (
                    <button
                      key={idx}
                      onClick={() => navigate(act.route)}
                      className={`px-3.5 py-1.5 text-xs font-bold rounded-lg transition-colors shadow-sm flex items-center gap-1 ${
                        act.primary
                          ? "bg-[#041627] hover:bg-[#006c49] text-white"
                          : "bg-white hover:bg-[#f0f4f8] text-[#041627] border border-[#d0d4d8]"
                      }`}
                    >
                      <span>{act.label}</span>
                      <span className="material-symbols-outlined text-[14px]">arrow_forward</span>
                    </button>
                  ))}
                </div>
              </div>
            </div>
          );
        })}

        {filteredNotifications.length === 0 && (
          <div className="p-8 text-center bg-white border border-[#e0e3e5] rounded-xl text-[#44474c] space-y-2">
            <span className="material-symbols-outlined text-4xl text-[#707478]">notifications_off</span>
            <p className="text-sm font-bold text-[#041627]">No notifications in this category</p>
            <p className="text-xs">Select another filter tab to view operational alerts.</p>
          </div>
        )}
      </div>

      {/* Bottom Preference Controls */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-4 border-t border-[#e0e3e5]">
        <div
          onClick={() => alert("Alert preferences configuration opened")}
          className="p-4 rounded-xl border border-[#e0e3e5] bg-white hover:border-[#006c49]/50 transition-colors cursor-pointer flex items-center justify-between group shadow-sm"
        >
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center group-hover:bg-[#041627] group-hover:text-white transition-colors">
              <span className="material-symbols-outlined text-[20px]">tune</span>
            </div>
            <div>
              <h4 className="text-xs font-bold text-[#041627]">Configure Alert Preferences</h4>
              <p className="text-[11px] text-[#44474c]">Mute non-critical logs or adjust SMS/Email relays</p>
            </div>
          </div>
          <span className="material-symbols-outlined text-[#44474c] group-hover:text-[#041627] transition-colors text-[20px]">
            chevron_right
          </span>
        </div>

        <div
          onClick={() => alert("Notification channels management opened")}
          className="p-4 rounded-xl border border-[#e0e3e5] bg-white hover:border-[#006c49]/50 transition-colors cursor-pointer flex items-center justify-between group shadow-sm"
        >
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-[#f0f4f8] text-[#041627] flex items-center justify-center group-hover:bg-[#041627] group-hover:text-white transition-colors">
              <span className="material-symbols-outlined text-[20px]">hub</span>
            </div>
            <div>
              <h4 className="text-xs font-bold text-[#041627]">Manage Notification Channels</h4>
              <p className="text-[11px] text-[#44474c]">Slack, Webhooks, WhatsApp Site Crews</p>
            </div>
          </div>
          <span className="material-symbols-outlined text-[#44474c] group-hover:text-[#041627] transition-colors text-[20px]">
            chevron_right
          </span>
        </div>
      </div>
    </div>
  );
};

export default NotificationsPage;
