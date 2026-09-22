import React from "react";

export const ProjectsPage: React.FC = () => {
  return (
    <div className="flex flex-col w-full pb-8 max-w-2xl mx-auto">
      {/* Project Header Banner */}
      <div className="px-6 pt-6 pb-6 bg-[#f7fafc]">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2 text-xs text-[#44474c] font-semibold tracking-wider">
            <span>PROJ-RT-104</span>
            <span className="w-1 h-1 rounded-full bg-[#c4c6cd]"></span>
            <span>COMMERCIAL</span>
          </div>
          <span className="inline-flex items-center px-2 py-1 rounded bg-[#6cf8bb]/20 text-[#00714d] font-semibold text-[11px]">
            ACTIVE
          </span>
        </div>
        <h1 className="text-2xl font-bold text-[#181c1e] mb-2">
          Riverside Tower
        </h1>
        <div className="flex items-center gap-2 text-sm text-[#44474c]">
          <span className="material-symbols-outlined text-[18px]">
            location_on
          </span>
          <span>450 River North Blvd, Chicago</span>
        </div>
      </div>

      {/* Metrics Summary Grid */}
      <div className="px-6 py-6">
        <div className="grid grid-cols-2 gap-4">
          <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 flex flex-col gap-2">
            <div className="flex items-center gap-2 text-[#44474c]">
              <span className="material-symbols-outlined text-[20px]">
                account_balance_wallet
              </span>
              <span className="text-[11px] uppercase tracking-wider font-semibold">
                Budget
              </span>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-2xl font-semibold text-[#181c1e]">62%</span>
              <span className="text-sm text-[#44474c]">spent</span>
            </div>
            <div className="w-full h-1.5 bg-[#e5e9eb] rounded-full overflow-hidden mt-1">
              <div
                className="h-full bg-[#041627] rounded-full"
                style={{ width: "62%" }}
              ></div>
            </div>
          </div>

          <div className="bg-white border border-[#e0e3e5] rounded-xl p-4 flex flex-col gap-2">
            <div className="flex items-center gap-2 text-[#44474c]">
              <span className="material-symbols-outlined text-[20px]">
                calendar_today
              </span>
              <span className="text-[11px] uppercase tracking-wider font-semibold">
                Schedule
              </span>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-2xl font-semibold text-[#181c1e]">
                Day 142
              </span>
            </div>
            <div className="text-sm text-[#006c49] flex items-center gap-1 mt-1 font-medium">
              <span className="material-symbols-outlined text-[16px]">
                trending_up
              </span>
              <span>On Track</span>
            </div>
          </div>

          <div className="col-span-2 bg-white border border-[#e0e3e5] rounded-xl p-4 flex flex-col gap-2">
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2 text-[#44474c]">
                <span className="material-symbols-outlined text-[20px]">
                  verified
                </span>
                <span className="text-[11px] uppercase tracking-wider font-semibold">
                  Compliance Score
                </span>
              </div>
              <span className="text-lg font-semibold text-[#006c49]">
                98/100
              </span>
            </div>
            <div className="text-sm text-[#44474c]">
              All permits active. Last OSHA inspection passed without citations.
            </div>
          </div>
        </div>
      </div>

      {/* Project Lifecycle Section */}
      <div className="px-6 pt-4 pb-2 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-[#181c1e]">
          Project Lifecycle
        </h2>
        <button className="text-[#041627] text-xs font-semibold hover:underline">
          VIEW GANTT
        </button>
      </div>

      <div className="px-6 pb-8 relative">
        <div className="absolute left-[39px] top-6 bottom-6 w-0.5 bg-[#e0e3e5] rounded-full z-0" />
        <div className="flex flex-col gap-0 relative z-10">
          {/* Phase 1 */}
          <div className="group w-full flex items-start gap-4 py-4 text-left">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#006c49] text-white flex items-center justify-center border-2 border-[#f7fafc] mt-0.5 shadow-sm">
              <span className="material-symbols-outlined text-[16px]">
                check
              </span>
            </div>
            <div className="flex-1 bg-white border border-[#e0e3e5] rounded-lg p-4">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-[#44474c] tracking-wider font-semibold">
                  PHASE 1
                </span>
                <span className="text-[11px] font-medium text-[#006c49]">
                  COMPLETE
                </span>
              </div>
              <h3 className="text-base font-semibold text-[#181c1e] mb-1">
                Design &amp; Planning
              </h3>
              <p className="text-sm text-[#44474c]">
                Blueprints finalized, city zoning approved.
              </p>
            </div>
          </div>

          {/* Phase 2 */}
          <div className="group w-full flex items-start gap-4 py-4 text-left">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#006c49] text-white flex items-center justify-center border-2 border-[#f7fafc] mt-0.5 shadow-sm">
              <span className="material-symbols-outlined text-[16px]">
                check
              </span>
            </div>
            <div className="flex-1 bg-white border border-[#e0e3e5] rounded-lg p-4">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-[#44474c] tracking-wider font-semibold">
                  PHASE 2
                </span>
                <span className="text-[11px] font-medium text-[#006c49]">
                  COMPLETE
                </span>
              </div>
              <h3 className="text-base font-semibold text-[#181c1e] mb-1">
                Procurement
              </h3>
              <p className="text-sm text-[#44474c]">
                All primary materials secured.
              </p>
            </div>
          </div>

          {/* Phase 3 */}
          <div className="group w-full flex items-start gap-4 py-4 text-left">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#006c49] text-white flex items-center justify-center border-2 border-[#f7fafc] mt-0.5 shadow-sm">
              <span className="material-symbols-outlined text-[16px]">
                check
              </span>
            </div>
            <div className="flex-1 bg-white border border-[#e0e3e5] rounded-lg p-4">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-[#44474c] tracking-wider font-semibold">
                  PHASE 3
                </span>
                <span className="text-[11px] font-medium text-[#006c49]">
                  COMPLETE
                </span>
              </div>
              <h3 className="text-base font-semibold text-[#181c1e] mb-1">
                Foundation
              </h3>
              <p className="text-sm text-[#44474c]">
                Excavation and pouring finished.
              </p>
            </div>
          </div>

          {/* Phase 4 - Active */}
          <div className="group w-full flex items-start gap-4 py-4 text-left">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#041627] flex items-center justify-center border-4 border-[#f7fafc] mt-0.5 shadow-md ring-1 ring-[#041627]/20">
              <div className="w-2.5 h-2.5 bg-white rounded-full animate-pulse" />
            </div>
            <div className="flex-1 bg-[#041627] text-white rounded-lg p-4 shadow-md relative overflow-hidden">
              <div className="absolute top-0 right-0 w-24 h-24 bg-white/5 rounded-bl-full -mr-4 -mt-4" />
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-[#b7c8de] tracking-wider font-semibold">
                  PHASE 4
                </span>
                <span className="text-[11px] font-medium text-[#6ffbbe]">
                  IN PROGRESS
                </span>
              </div>
              <h3 className="text-base font-semibold text-white mb-1">
                Structure
              </h3>
              <p className="text-sm text-[#b7c8de] mb-3">
                Steel framing reaching floor 12 of 30.
              </p>
              <div className="w-full bg-[#1a2b3c] rounded-full h-1.5 mb-1">
                <div
                  className="bg-[#6ffbbe] h-1.5 rounded-full"
                  style={{ width: "40%" }}
                />
              </div>
              <div className="text-right text-[11px] text-[#b7c8de]">40%</div>
            </div>
          </div>

          {/* Phase 5 */}
          <div className="group w-full flex items-start gap-4 py-4 text-left opacity-75">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#e0e3e5] flex items-center justify-center border-2 border-[#f7fafc] mt-0.5">
              <div className="w-2 h-2 bg-[#44474c] rounded-full" />
            </div>
            <div className="flex-1 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg p-4">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-[#44474c] tracking-wider font-semibold">
                  PHASE 5
                </span>
                <span className="text-[11px] font-medium text-[#44474c]">
                  UPCOMING
                </span>
              </div>
              <h3 className="text-base font-semibold text-[#181c1e] mb-1">
                MEP Installation
              </h3>
              <p className="text-sm text-[#44474c]">
                Mechanical, electrical, plumbing.
              </p>
            </div>
          </div>

          {/* Phase 6 */}
          <div className="group w-full flex items-start gap-4 py-4 text-left opacity-60">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#e0e3e5] flex items-center justify-center border-2 border-[#f7fafc] mt-0.5">
              <div className="w-2 h-2 bg-[#44474c] rounded-full" />
            </div>
            <div className="flex-1 bg-[#f7fafc] border border-[#e0e3e5] rounded-lg p-4">
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs text-[#44474c] tracking-wider font-semibold">
                  PHASE 6
                </span>
                <span className="text-[11px] font-medium text-[#44474c]">
                  UPCOMING
                </span>
              </div>
              <h3 className="text-base font-semibold text-[#181c1e] mb-1">
                Finishing
              </h3>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProjectsPage;
