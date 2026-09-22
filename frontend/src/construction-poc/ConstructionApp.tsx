import React, { useState } from "react";
import { ConstructionHeader } from "./components/ConstructionHeader";
import { BottomNavigation, NavTabPath } from "./components/BottomNavigation";
import { ConstructionDashboard } from "./pages/ConstructionDashboard";
import { ProjectsPage } from "./pages/ProjectsPage";
import { ActivityPage } from "./pages/ActivityPage";
import { DocumentsPage } from "./pages/DocumentsPage";
import { ProfilePage } from "./pages/ProfilePage";

export const ConstructionApp: React.FC = () => {
  const [activeTab, setActiveTab] = useState<NavTabPath>("home");

  const renderActiveTab = () => {
    switch (activeTab) {
      case "home":
        return (
          <ConstructionDashboard
            onNavigateToProjects={() => setActiveTab("projects")}
            onNavigateToActivity={() => setActiveTab("activity")}
          />
        );
      case "projects":
        return <ProjectsPage />;
      case "activity":
        return <ActivityPage />;
      case "documents":
        return <DocumentsPage />;
      case "profile":
        return <ProfilePage />;
      default:
        return (
          <ConstructionDashboard
            onNavigateToProjects={() => setActiveTab("projects")}
            onNavigateToActivity={() => setActiveTab("activity")}
          />
        );
    }
  };

  return (
    <div className="bg-[#f7fafc] text-[#181c1e] font-sans min-h-screen flex flex-col antialiased">
      {/* BuildFlow Application Header */}
      <ConstructionHeader
        onNotificationClick={() => setActiveTab("activity")}
        onProfileClick={() => setActiveTab("profile")}
      />

      {/* Main Screen Workspace */}
      <main className="flex-1 pt-16 pb-24 bg-[#f7fafc]">
        {renderActiveTab()}
      </main>

      {/* Sticky Bottom Navigation Bar */}
      <BottomNavigation
        activeTab={activeTab}
        onTabChange={(tab) => setActiveTab(tab)}
      />
    </div>
  );
};

export default ConstructionApp;
