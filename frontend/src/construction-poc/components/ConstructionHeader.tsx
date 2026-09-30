import React from "react";
import { useRoleContext, ROLES, UserRole } from "../hooks/useRoleContext";

interface ConstructionHeaderProps {
  onNotificationClick?: () => void;
  onProfileClick?: () => void;
}

export const ConstructionHeader: React.FC<ConstructionHeaderProps> = ({
  onNotificationClick,
  onProfileClick,
}) => {
  const { activeRole, setRole, roleConfig } = useRoleContext();

  return (
    <header className="fixed top-0 w-full z-50 bg-[#f7fafc]/90 backdrop-blur-xl shadow-[0_1px_8px_rgba(0,0,0,0.04)] pt-safe border-b border-[#e0e3e5]">
      <div className="h-16 px-6 flex items-center justify-between max-w-7xl mx-auto">
        <div className="flex items-center gap-3">
          <span className="text-[#041627] font-bold text-xl tracking-tight">
            BuildFlow
          </span>
          <span className="text-xs px-2 py-0.5 bg-[#041627]/10 text-[#041627] rounded font-semibold hidden md:inline-block">
            Sovereign Workspace
          </span>
        </div>

        {/* Persona / Role Selector */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 bg-white border border-[#e0e3e5] px-3 py-1.5 rounded-lg shadow-sm">
            <span className="material-symbols-outlined text-[#041627] text-[18px]">
              {roleConfig.icon}
            </span>
            <div className="flex flex-col">
              <span className="text-[10px] text-[#44474c] font-semibold uppercase tracking-wider leading-none">
                Active Persona
              </span>
              <select
                value={activeRole}
                onChange={(e) => setRole(e.target.value as UserRole)}
                className="text-xs font-bold text-[#041627] bg-transparent outline-none cursor-pointer pr-1"
                title="Frontend Persona Switcher (controls UI visibility & available actions)"
              >
                {(Object.keys(ROLES) as UserRole[]).map((r) => (
                  <option key={r} value={r}>
                    {ROLES[r].label} ({r})
                  </option>
                ))}
              </select>
            </div>
          </div>

          <button
            onClick={onNotificationClick}
            className="w-10 h-10 flex items-center justify-center text-[#44474c] hover:text-[#041627] transition-colors rounded-full hover:bg-black/5"
            aria-label="Notifications"
          >
            <span className="material-symbols-outlined text-[20px]">notifications</span>
          </button>
          
          <div
            onClick={onProfileClick}
            className="w-8 h-8 rounded-full bg-[#041627] flex items-center justify-center cursor-pointer hover:opacity-90 transition-opacity shadow-sm"
            aria-label="User Profile"
          >
            <span className="material-symbols-outlined text-white text-[16px]">
              person
            </span>
          </div>
        </div>
      </div>
    </header>
  );
};

export default ConstructionHeader;
