import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

export const FieldModePage: React.FC = () => {
  const navigate = useNavigate();

  // Active state for scanner / photo presentation feedback modals
  const [activeFeedback, setActiveFeedback] = useState<string | null>(null);

  const handleActionClick = (actionName: string, routeTarget?: string) => {
    if (routeTarget) {
      navigate(routeTarget);
    } else {
      setActiveFeedback(actionName);
    }
  };

  return (
    <div className="w-full min-h-screen bg-[#0f1923] text-white font-[Inter] flex flex-col items-center justify-between select-none overflow-x-hidden">
      {/* Mobile Shell Wrapper (max-w-md viewport container) */}
      <div className="w-full max-w-md min-h-screen flex flex-col bg-[#0b131c] relative shadow-2xl pb-24 border-x border-[#1a2b3c]">
        
        {/* Fixed Top Header (Stitch Spec: BuildFlow | notifications | profile) */}
        <header className="sticky top-0 z-30 w-full backdrop-blur-md bg-[#041627]/85 border-b border-[#1a2b3c] px-4 py-3 flex items-center justify-between shadow-md">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-xl bg-[#6cf8bb] text-[#041627] flex items-center justify-center font-bold text-base shadow-sm">
              <span className="material-symbols-outlined text-[20px]">architecture</span>
            </div>
            <div>
              <span className="font-extrabold text-base tracking-tight text-white block leading-tight">BuildFlow</span>
              <span className="text-[10px] text-[#6cf8bb] font-mono tracking-widest uppercase font-semibold">Field Mode</span>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => setActiveFeedback("Notifications")}
              className="relative p-2 text-[#b7c8de] hover:text-white transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined text-[22px]">notifications</span>
              <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-[#ffb74d] animate-ping" />
              <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-[#ffb74d]" />
            </button>

            <button
              type="button"
              onClick={() => navigate(CONSTRUCTION_ROUTES.PROFILE)}
              className="w-8 h-8 rounded-full bg-[#1a2b3c] border border-[#2c3d50] flex items-center justify-center text-[#6cf8bb] hover:text-white transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined text-[20px]">account_circle</span>
            </button>
          </div>
        </header>

        {/* Active Site Context Card (Stitch Spec: Site Active, Riverside Tower, Level 4 / Zone B) */}
        <div className="bg-[#041627] rounded-b-3xl p-5 border-b border-[#1a2b3c] shadow-lg relative overflow-hidden">
          {/* Ambient Background Radial Glow */}
          <div className="absolute -top-12 -right-12 w-44 h-44 rounded-full bg-[#006c49]/20 blur-2xl pointer-events-none" />
          <div className="absolute -bottom-10 -left-10 w-36 h-36 rounded-full bg-[#004b75]/25 blur-xl pointer-events-none" />

          {/* Top Status & Time Row */}
          <div className="flex items-center justify-between mb-3 relative z-10">
            <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#006c49]/30 border border-[#006c49]/50 text-[#6cf8bb] text-xs font-bold uppercase tracking-wider">
              <span className="w-2 h-2 rounded-full bg-[#6cf8bb] animate-pulse" />
              Site Active
            </div>
            <span className="text-xs font-mono font-semibold text-[#b7c8de] bg-[#102235] px-2.5 py-1 rounded-lg border border-[#1e344d]">
              09:42 CST
            </span>
          </div>

          {/* Site Title & Zone Information */}
          <div className="relative z-10 mb-2">
            <h1 className="text-2xl font-black text-white tracking-tight mb-1">Riverside Tower</h1>
            <div className="flex items-center gap-2 text-xs text-[#b7c8de]">
              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-[#102235] text-[#6cf8bb] font-bold border border-[#1e344d]">
                <span className="material-symbols-outlined text-[14px]">location_on</span>
                Level 4 / Zone B
              </span>
            </div>
          </div>
        </div>

        {/* Main Glove-Friendly Action Surface */}
        <main className="flex-1 p-4 flex flex-col gap-4">
          
          {/* Section Header */}
          <div className="flex items-center justify-between px-1">
            <h2 className="text-xs font-bold uppercase tracking-wider text-[#8a9bb0]">
              Glove-Friendly Field Actions
            </h2>
            <span className="text-[11px] text-[#6cf8bb] font-semibold">Touch Optimized</span>
          </div>

          {/* Primary Action (Full-Width): REPORT ISSUE */}
          <button
            type="button"
            onClick={() => handleActionClick("Report Issue", CONSTRUCTION_ROUTES.ISSUES)}
            className="w-full bg-gradient-to-r from-[#ba1a1a] to-[#93000a] text-white p-5 rounded-2xl shadow-xl border border-[#ff897d]/30 flex items-center justify-between cursor-pointer transition-all active:scale-[0.98] hover:shadow-red-900/40"
          >
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-white/10 text-white flex items-center justify-center font-bold shadow-inner">
                <span className="material-symbols-outlined text-[28px]">warning</span>
              </div>
              <div className="text-left">
                <h3 className="text-lg font-black text-white leading-tight">Report Issue</h3>
                <p className="text-xs text-[#ffdad6]/90 font-medium">Log a hazard or defect</p>
              </div>
            </div>
            <span className="material-symbols-outlined text-[28px] text-white/80">chevron_right</span>
          </button>

          {/* 2-Column Action Grid for Secondary Operational Cards */}
          <div className="grid grid-cols-2 gap-3.5">
            
            {/* Action 1: Photo */}
            <button
              type="button"
              onClick={() => handleActionClick("Photo")}
              className="bg-[#121e2b] hover:bg-[#1a2b3c] border border-[#233549] hover:border-[#6cf8bb]/50 p-4 rounded-2xl flex flex-col justify-between gap-3 text-left transition-all active:scale-[0.97] cursor-pointer shadow-md group min-h-[110px]"
            >
              <div className="w-10 h-10 rounded-xl bg-[#6cf8bb]/10 text-[#6cf8bb] flex items-center justify-center group-hover:scale-110 transition-transform">
                <span className="material-symbols-outlined text-[24px]">photo_camera</span>
              </div>
              <div>
                <h4 className="font-bold text-white text-base mb-0.5">Photo</h4>
                <p className="text-[11px] text-[#8a9bb0]">Photo</p>
              </div>
            </button>

            {/* Action 2: Inspect */}
            <button
              type="button"
              onClick={() => handleActionClick("Inspect", CONSTRUCTION_ROUTES.INSPECTIONS)}
              className="bg-[#121e2b] hover:bg-[#1a2b3c] border border-[#233549] hover:border-[#6cf8bb]/50 p-4 rounded-2xl flex flex-col justify-between gap-3 text-left transition-all active:scale-[0.97] cursor-pointer shadow-md group min-h-[110px]"
            >
              <div className="w-10 h-10 rounded-xl bg-[#006c49]/20 text-[#6cf8bb] flex items-center justify-center group-hover:scale-110 transition-transform">
                <span className="material-symbols-outlined text-[24px]">fact_check</span>
              </div>
              <div>
                <h4 className="font-bold text-white text-base mb-0.5">Inspect</h4>
                <p className="text-[11px] text-[#8a9bb0]">Inspect</p>
              </div>
            </button>

            {/* Action 3: Scan */}
            <button
              type="button"
              onClick={() => handleActionClick("Scan")}
              className="bg-[#121e2b] hover:bg-[#1a2b3c] border border-[#233549] hover:border-[#6cf8bb]/50 p-4 rounded-2xl flex flex-col justify-between gap-3 text-left transition-all active:scale-[0.97] cursor-pointer shadow-md group min-h-[110px]"
            >
              <div className="w-10 h-10 rounded-xl bg-[#ffb74d]/10 text-[#ffb74d] flex items-center justify-center group-hover:scale-110 transition-transform">
                <span className="material-symbols-outlined text-[24px]">qr_code_scanner</span>
              </div>
              <div>
                <h4 className="font-bold text-white text-base mb-0.5">Scan</h4>
                <p className="text-[11px] text-[#8a9bb0]">Scan</p>
              </div>
            </button>

            {/* Action 4: Delivery */}
            <button
              type="button"
              onClick={() => handleActionClick("Delivery", CONSTRUCTION_ROUTES.DELIVERIES)}
              className="bg-[#121e2b] hover:bg-[#1a2b3c] border border-[#233549] hover:border-[#6cf8bb]/50 p-4 rounded-2xl flex flex-col justify-between gap-3 text-left transition-all active:scale-[0.97] cursor-pointer shadow-md group min-h-[110px]"
            >
              <div className="w-10 h-10 rounded-xl bg-[#004b75]/30 text-[#40c4ff] flex items-center justify-center group-hover:scale-110 transition-transform">
                <span className="material-symbols-outlined text-[24px]">local_shipping</span>
              </div>
              <div>
                <h4 className="font-bold text-white text-base mb-0.5">Delivery</h4>
                <p className="text-[11px] text-[#8a9bb0]">Delivery</p>
              </div>
            </button>

          </div>

        </main>

        {/* Presentation Feedback Overlay (for photo/scan feedback) */}
        {activeFeedback && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fadeIn">
            <div className="bg-[#121e2b] border border-[#233549] rounded-2xl p-6 max-w-xs w-full text-center shadow-2xl flex flex-col items-center gap-3">
              <div className="w-12 h-12 rounded-full bg-[#6cf8bb]/20 text-[#6cf8bb] flex items-center justify-center">
                <span className="material-symbols-outlined text-[28px]">touch_app</span>
              </div>
              <h3 className="text-base font-bold text-white">{activeFeedback}</h3>
              <p className="text-xs text-[#8a9bb0]">
                Field Mode presentation surface. Action capabilities will bind to application commands when implemented.
              </p>
              <button
                type="button"
                onClick={() => setActiveFeedback(null)}
                className="mt-2 px-5 py-2 bg-[#006c49] text-white font-bold text-xs rounded-xl hover:bg-[#006c49]/90 transition-colors w-full cursor-pointer"
              >
                Acknowledge
              </button>
            </div>
          </div>
        )}

        {/* Fixed Bottom Navigation (Stitch Spec: Home | Projects | Activity | Docs | Profile) */}
        <nav className="fixed bottom-0 max-w-md w-full backdrop-blur-md bg-[#041627]/95 border-t border-[#1a2b3c] px-3 py-2 flex items-center justify-around z-40 text-[11px] font-semibold text-[#8a9bb0]">
          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.HOME)}
            className="flex flex-col items-center gap-0.5 hover:text-white transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[20px]">home</span>
            <span>Home</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
            className="flex flex-col items-center gap-0.5 hover:text-white transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[20px]">domain</span>
            <span>Projects</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.ACTIVITY)}
            className="flex flex-col items-center gap-0.5 hover:text-white transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[20px]">pulse</span>
            <span>Activity</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.DOCUMENTS)}
            className="flex flex-col items-center gap-0.5 hover:text-white transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[20px]">description</span>
            <span>Docs</span>
          </button>

          <button
            type="button"
            onClick={() => navigate(CONSTRUCTION_ROUTES.PROFILE)}
            className="flex flex-col items-center gap-0.5 hover:text-white transition-colors cursor-pointer"
          >
            <span className="material-symbols-outlined text-[20px]">person</span>
            <span>Profile</span>
          </button>
        </nav>

      </div>
    </div>
  );
};

export default FieldModePage;


