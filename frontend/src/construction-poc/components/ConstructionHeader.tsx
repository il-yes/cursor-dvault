import React from "react";

interface ConstructionHeaderProps {
  onNotificationClick?: () => void;
  onProfileClick?: () => void;
}

export const ConstructionHeader: React.FC<ConstructionHeaderProps> = ({
  onNotificationClick,
  onProfileClick,
}) => {
  return (
    <header className="fixed top-0 w-full z-50 bg-[#f7fafc]/80 backdrop-blur-xl shadow-[0_1px_8px_rgba(0,0,0,0.04)] pt-safe">
      <div className="h-16 px-6 flex items-center justify-between max-w-7xl mx-auto">
        <div className="flex items-center gap-2">
          <span className="text-[#041627] font-semibold text-lg tracking-tight">
            BuildFlow
          </span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={onNotificationClick}
            className="w-11 h-11 flex items-center justify-center text-[#44474c] hover:text-[#041627] transition-colors rounded-full hover:bg-black/5"
            aria-label="Notifications"
          >
            <span className="material-symbols-outlined">notifications</span>
          </button>
          <div
            onClick={onProfileClick}
            className="w-8 h-8 rounded-full bg-[#041627] flex items-center justify-center ml-2 cursor-pointer hover:opacity-90 transition-opacity"
            aria-label="User Profile"
          >
            <span className="material-symbols-outlined text-white text-[18px]">
              person
            </span>
          </div>
        </div>
      </div>
    </header>
  );
};
