import React, { useEffect } from "react";
import { useC3ConfigurationStore } from "./store/useC3ConfigurationStore";
import { TrustGroupsTab } from "./tabs/TrustGroupsTab";
import { ChannelsConfigTab } from "./tabs/ChannelsConfigTab";
import { AssignmentsTab } from "./tabs/AssignmentsTab";
import { PolicyTab } from "./tabs/PolicyTab";
import { PropertiesTab } from "./tabs/PropertiesTab";
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
    isLoading,
    error,
    trustGroups,
    draftSlots,
    draftAssignments,
    draftProperties,
    originalChannel,
  } = useC3ConfigurationStore();

  useEffect(() => {
    console.log(`[BOUNDARIES][READ] C3ConfigurationPage mounted -> loadConfiguration()`);
    useC3ConfigurationStore.getState().loadConfiguration();
  }, []);


  return (
    <DashboardLayout>
      <div style={{ flex: 1, display: "flex", flexDirection: "column", minHeight: "100%", overflow: "hidden", position: "relative" }}>
        {/* Error Banner */}
        {error && (
          <div className="c3-error-banner" style={{ padding: "12px 20px", background: "#FEF2F2", borderBottom: "1px solid #FCA5A5", color: "#991B1B", display: "flex", alignItems: "center", justifyContent: "space-between" }}>
            <div>
              <strong>Failed to load C3 Configuration:</strong> {error}
            </div>
            <button
              className="btn btn-secondary"
              onClick={() => loadConfiguration()}
              style={{ background: "#991B1B", color: "#FFF", border: "none", borderRadius: "6px", padding: "6px 14px", cursor: "pointer", fontSize: "13px" }}
            >
              Retry
            </button>
          </div>
        )}

        {/* Loading Indicator */}
        {isLoading && (
          <div className="c3-loading-bar" style={{ padding: "8px 20px", background: "#EFF6FF", borderBottom: "1px solid #BFDBFE", color: "#1E40AF", fontSize: "13px" }}>
            Loading C3 Channel Configuration from backend...
          </div>
        )}

        {/* Channel Header */}
        <div className="ch-header">
          <div className="ch-title-row">
            <div className="ch-title">{originalChannel?.title || "No Active Channel Selected"}</div>
            <span className={`badge ${originalChannel ? "badge-active" : "badge-inactive"}`}>
              {originalChannel ? "Active" : "Inactive"}
            </span>
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
              <span className="ch-meta-key">updated</span>
              <span className="ch-meta-val">2026-06-12 · 16:44</span>
            </div>
          </div>
        </div>

        {/* Configuration Tab Bar (5 Tabs) */}
        <div className="tab-bar">
          <div
            className={`tab ${activeTab === "channels" || activeTab === "slots" ? "active" : ""}`}
            onClick={() => setActiveTab("channels")}
          >
            Slots <span className="tab-ct">{draftSlots.length}</span>
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
            Assignments <span className="tab-ct">{draftAssignments.length}</span>
          </div>
          <div
            className={`tab ${activeTab === "security" || activeTab === "properties" ? "active" : ""}`}
            onClick={() => setActiveTab("security")}
          >
            Properties <span className="tab-ct">{draftProperties.length}</span>
          </div>
          <div
            className={`tab ${activeTab === "templates" || activeTab === "policy" ? "active" : ""}`}
            onClick={() => setActiveTab("templates")}
          >
            Policy
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
          ) : (
            <PolicyTab />
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

