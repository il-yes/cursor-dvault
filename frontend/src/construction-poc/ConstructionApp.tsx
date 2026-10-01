import React from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { RoleProvider } from "./hooks/useRoleContext";
import { CONSTRUCTION_ROUTES } from "./constants/routes";
import { ConstructionHeader } from "./components/ConstructionHeader";
import { BottomNavigation } from "./components/BottomNavigation";
import { ConstructionDashboard } from "./pages/ConstructionDashboard";
import { ProjectListingPage } from "./pages/ProjectListingPage";
import { ProjectPage } from "./pages/ProjectPage";
import { StakeholdersPage } from "./pages/StakeholdersPage";
import { MaterialRequirementPage } from "./pages/MaterialRequirementPage";
import { SupplierOffersPage } from "./pages/SupplierOffersPage";
import { DeliveryDetailPage } from "./pages/DeliveryDetailPage";
import { TransportDelayPage } from "./pages/TransportDelayPage";
import { ConstructionIssuePage } from "./pages/ConstructionIssuePage";
import { DecisionPage } from "./pages/DecisionPage";
import { InspectionPage } from "./pages/InspectionPage";
import { FieldModePage } from "./pages/FieldModePage";
import { CollaborationThreadPage } from "./pages/CollaborationThreadPage";
import { ProvenanceWhyPage } from "./pages/ProvenanceWhyPage";
import { ProjectHistoryPage } from "./pages/ProjectHistoryPage";
import { ProjectChannelsPage } from "./pages/ProjectChannelsPage";
import { ActivityPage } from "./pages/ActivityPage";
import { DocumentsPage } from "./pages/DocumentsPage";
import { ProfilePage } from "./pages/ProfilePage";

export const ConstructionAppInner: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="bg-[#f7fafc] text-[#181c1e] font-sans min-h-screen flex flex-col antialiased">
      {/* BuildFlow Application Header */}
      <ConstructionHeader
        onNotificationClick={() => navigate(CONSTRUCTION_ROUTES.ACTIVITY)}
        onProfileClick={() => navigate(CONSTRUCTION_ROUTES.PROFILE)}
      />

      {/* Main Screen Workspace */}
      <main className="flex-1 pt-16 pb-24 bg-[#f7fafc]">
        <Routes>
          {/* Sovereign Vault Home Front Door */}
          <Route path="/" element={<ConstructionDashboard />} />

          {/* Projects Collection List & Project Detail Hierarchy */}
          <Route path="/projects" element={<ProjectListingPage />} />
          <Route path="/projects/:projectId" element={<ProjectPage />} />

          {/* Channels & Threads Collaboration Surface */}
          <Route path="/channels" element={<ProjectChannelsPage />} />
          <Route path="/channels/:channelId" element={<ProjectChannelsPage />} />

          {/* Construction Scenario Entity Detail Views */}
          <Route path="/stakeholders" element={<StakeholdersPage />} />
          <Route path="/requirements" element={<MaterialRequirementPage />} />
          <Route path="/offers" element={<SupplierOffersPage />} />
          <Route path="/deliveries" element={<DeliveryDetailPage />} />
          <Route path="/transport" element={<TransportDelayPage />} />
          <Route path="/issues" element={<ConstructionIssuePage />} />
          <Route path="/decisions" element={<DecisionPage />} />
          <Route path="/inspections" element={<InspectionPage />} />
          <Route path="/field" element={<FieldModePage />} />
          <Route path="/thread" element={<CollaborationThreadPage />} />
          <Route path="/provenance" element={<ProvenanceWhyPage />} />
          <Route path="/history" element={<ProjectHistoryPage />} />

          {/* Secondary Pages */}
          <Route path="/activity" element={<ActivityPage />} />
          <Route path="/documents" element={<DocumentsPage />} />
          <Route path="/profile" element={<ProfilePage />} />
        </Routes>
      </main>

      {/* Sticky Bottom Navigation Bar */}
      <BottomNavigation />
    </div>
  );
};

export const ConstructionApp: React.FC = () => (
  <RoleProvider>
    <ConstructionAppInner />
  </RoleProvider>
);

export default ConstructionApp;
