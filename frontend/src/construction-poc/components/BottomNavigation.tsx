import React from "react";

export type NavTabPath = "home" | "projects" | "activity" | "documents" | "profile";

interface BottomNavigationProps {
  activeTab: NavTabPath;
  onTabChange: (tab: NavTabPath) => void;
}

interface NavItem {
  path: NavTabPath;
  label: string;
  iconName: string;
}

const navItems: NavItem[] = [
  { path: "home", label: "Home", iconName: "home" },
  { path: "projects", label: "Projects", iconName: "architecture" },
  { path: "activity", label: "Activity", iconName: "pending_actions" },
  { path: "documents", label: "Docs", iconName: "description" },
  { path: "profile", label: "Profile", iconName: "person_outline" },
];

export const BottomNavigation: React.FC<BottomNavigationProps> = ({
  activeTab,
  onTabChange,
}) => {
  return (
    <nav className="fixed bottom-0 w-full z-50 bg-[#f7fafc]/80 backdrop-blur-xl shadow-[0_-1px_8px_rgba(0,0,0,0.04)] pb-safe">
      <div className="flex h-16 items-center justify-around px-2 max-w-7xl mx-auto">
        {navItems.map((item) => {
          const isActive = activeTab === item.path;
          return (
            <button
              key={item.path}
              onClick={() => onTabChange(item.path)}
              className={`flex flex-col items-center justify-center gap-1 min-w-[64px] transition-all ${
                isActive
                  ? "text-[#041627] font-semibold"
                  : "text-[#44474c] hover:text-[#041627]"
              }`}
            >
              <span className="material-symbols-outlined">{item.iconName}</span>
              <span className="text-[11px] font-medium">{item.label}</span>
            </button>
          );
        })}
      </div>
    </nav>
  );
};
