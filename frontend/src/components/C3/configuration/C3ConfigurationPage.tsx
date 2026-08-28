import React, { useEffect } from "react";
import { useC3ConfigurationStore } from "./store/useC3ConfigurationStore";
import { TrustGroupsTab } from "./tabs/TrustGroupsTab";
import { ChannelsConfigTab } from "./tabs/ChannelsConfigTab";
import { AssignmentsTab } from "./tabs/AssignmentsTab";
import { PolicyTab } from "./tabs/PolicyTab";
import { PropertiesTab } from "./tabs/PropertiesTab";
import { FederationTab } from "./tabs/FederationTab";
import { CreateTrustGroupDialog } from "./dialogs/CreateTrustGroupDialog";
import { DashboardLayout } from "@/components/DashboardLayout";
import "./styles/c3-template-theme.css";

export const C3ConfigurationPageContent: React.FC = () => {
  const {
    activeTab,
    setActiveTab,
    isCreateTrustGroupOpen,
    setCreateTrustGroupOpen,
    loadConfiguration,
    trustGroups,
  } = useC3ConfigurationStore();

  useEffect(() => {
    loadConfiguration();
  }, [loadConfiguration]);

  return (
    <DashboardLayout>
      <div style={{ flex: 1, display: "flex", flexDirection: "column", background: "#fff", minHeight: "100%", overflow: "hidden", position: "relative" }}>
        {/* Channel Header */}
        <div className="ch-header">
          <div className="ch-title-row">
            <div className="ch-title">contract-execution</div>
            <span className="badge badge-active">Active</span>
            <div style={{ flex: 1 }}></div>
            <button className="btn btn-primary" onClick={() => setCreateTrustGroupOpen(true)}>
              + Create Trust Group
            </button>
          </div>
          <div className="ch-meta">
            <div className="ch-meta-item">
              <span className="ch-meta-key">workspace</span>
              <span className="ch-meta-val">acme-corp</span>
            </div>
            <div className="ch-meta-item">
              <span className="ch-meta-key">template</span>
              <span className="ch-meta-val">contract-execution-v1</span>
            </div>
            <div className="ch-meta-item">
              <span className="ch-meta-key">federation</span>
              <span className="ch-meta-val">internal</span>
            </div>
            <div className="ch-meta-item">
              <span className="ch-meta-key">updated</span>
              <span className="ch-meta-val">2026-06-12 · 16:44</span>
            </div>
          </div>
        </div>

        {/* Configuration Tab Bar */}
        <div className="tab-bar">
          <div
            className={`tab ${activeTab === "channels" || activeTab === "slots" ? "active" : ""}`}
            onClick={() => setActiveTab("channels")}
          >
            Slots <span className="tab-ct">3</span>
          </div>
          <div
            className={`tab ${activeTab === "trust-groups" || activeTab === "trust" || activeTab === "overview" ? "active" : ""}`}
            onClick={() => setActiveTab("trust-groups")}
          >
            Trust <span className="tab-ct">{trustGroups.length}</span>
          </div>
          <div
            className={`tab ${activeTab === "collaboration" || activeTab === "assignments" ? "active" : ""}`}
            onClick={() => setActiveTab("collaboration")}
          >
            Assignments <span className="tab-ct">3</span>
          </div>
          <div
            className={`tab ${activeTab === "security" || activeTab === "properties" ? "active" : ""}`}
            onClick={() => setActiveTab("security")}
          >
            Properties <span className="tab-ct">5</span>
          </div>
          <div
            className={`tab ${activeTab === "templates" || activeTab === "policy" ? "active" : ""}`}
            onClick={() => setActiveTab("templates")}
          >
            Policy
          </div>
          <div
            className={`tab ${activeTab === "federation" ? "active" : ""}`}
            onClick={() => setActiveTab("federation" as any)}
          >
            Federation <span className="tab-ct">3</span>
          </div>
        </div>

        {/* Tab View Content */}
        <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden", position: "relative" }}>
          {activeTab === "trust-groups" || activeTab === "trust" || activeTab === "overview" ? (
            <TrustGroupsTab />
          ) : activeTab === "channels" || activeTab === "slots" ? (
            <ChannelsConfigTab />
          ) : activeTab === "collaboration" || activeTab === "assignments" ? (
            <AssignmentsTab />
          ) : activeTab === "security" || activeTab === "properties" ? (
            <PropertiesTab />
          ) : activeTab === "templates" || activeTab === "policy" ? (
            <PolicyTab />
          ) : (
            <FederationTab />
          )}
        </div>

        {/* Create Trust Group Light Modal */}
        <CreateTrustGroupDialog
          isOpen={isCreateTrustGroupOpen}
          onClose={() => setCreateTrustGroupOpen(false)}
        />
      </div>
    </DashboardLayout>
  );
};


export const C3ConfigurationPage: React.FC = () => {
  return <C3ConfigurationPageContent />;
};

export default C3ConfigurationPage;

