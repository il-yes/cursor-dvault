import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

export type NavTabPath = "home" | "projects" | "activity" | "documents" | "profile";

interface BottomNavigationProps {
  activeTab?: NavTabPath;
  onTabChange?: (tab: NavTabPath) => void;
}

interface NavItem {
  id: NavTabPath;
  label: string;
  iconName: string;
  route: string;
}

const navItems: NavItem[] = [
  { id: "home", label: "Home", iconName: "home", route: CONSTRUCTION_ROUTES.HOME },
  { id: "projects", label: "Projects", iconName: "architecture", route: CONSTRUCTION_ROUTES.PROJECTS },
  { id: "activity", label: "Activity", iconName: "pending_actions", route: CONSTRUCTION_ROUTES.ACTIVITY },
  { id: "documents", label: "Docs", iconName: "description", route: CONSTRUCTION_ROUTES.DOCUMENTS },
  { id: "profile", label: "Profile", iconName: "person_outline", route: CONSTRUCTION_ROUTES.PROFILE },
];

export const BottomNavigation: React.FC<BottomNavigationProps> = ({
  activeTab,
  onTabChange,
}) => {
  const location = useLocation();
  const navigate = useNavigate();

  const getIsActive = (item: NavItem) => {
    if (activeTab) return activeTab === item.id;
    const pathname = location.pathname;
    if (item.id === "home") {
      return pathname === CONSTRUCTION_ROUTES.HOME || pathname === `${CONSTRUCTION_ROUTES.HOME}/`;
    }
    return pathname.startsWith(`${CONSTRUCTION_ROUTES.HOME}/${item.id}`);
  };

  const handleClick = (item: NavItem) => {
    if (onTabChange) {
      onTabChange(item.id);
    }
    navigate(item.route);
  };

  return (
    <nav className="fixed bottom-0 w-full z-50 bg-[#f7fafc]/80 backdrop-blur-xl shadow-[0_-1px_8px_rgba(0,0,0,0.04)] pb-safe">
      <div className="flex h-16 items-center justify-around px-2 max-w-7xl mx-auto">
        {navItems.map((item) => {
          const isActive = getIsActive(item);
          return (
            <button
              key={item.id}
              onClick={() => handleClick(item)}
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
