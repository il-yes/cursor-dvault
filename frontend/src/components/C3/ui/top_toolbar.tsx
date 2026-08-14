import React, { useState } from "react";
import { useC3WorkspaceStore } from "@/components/C3/infrastructure/store/useC3WorkspaceStore";

interface TopToolbarProps {
    setOpenNewThread: (open: boolean) => void;
    setOpenCreateWorkspace?: (open: boolean) => void;
    setOpenCreateChannel?: (open: boolean) => void;
}

export const TopToolbar = ({
    setOpenNewThread,
    setOpenCreateWorkspace,
    setOpenCreateChannel,
}: TopToolbarProps) => {
    const [dropdownOpen, setDropdownOpen] = useState(false);

    const {
        workspaces,
        activeWorkspace,
        activeWorkspaceId,
        isLoading,
        selectWorkspace,
    } = useC3WorkspaceStore();

    const handleSelectWorkspace = (workspaceId: string) => {
        selectWorkspace(workspaceId);
        setDropdownOpen(false);
    };

    const handleCreateWorkspace = () => {
        setDropdownOpen(false);
        if (setOpenCreateWorkspace) {
            setOpenCreateWorkspace(true);
        }
    };

    console.log({workspaces})

    const displayName = activeWorkspace?.name || "Select Workspace";

    return (
        <div className="topbar">
            <div className="topbar-logo">C3</div>

            {/* Workspace Selector */}
            <div className="workspace-selector">
                <div
                    id="workspace-selector-trigger"
                    className={`workspace-pill${dropdownOpen ? " active" : ""}`}
                    onClick={() => setDropdownOpen(!dropdownOpen)}
                >
                    <span>{displayName}</span>
                    <span className="workspace-pill-chevron">▾</span>
                </div>

                {dropdownOpen && (
                    <>
                        {/* Click-outside backdrop */}
                        <div
                            className="workspace-dropdown-backdrop"
                            onClick={() => setDropdownOpen(false)}
                        />

                        <div className="workspace-dropdown" id="workspace-dropdown-panel">
                            <div className="workspace-dropdown-header">Workspaces</div>

                            <div className="workspace-dropdown-list">
                                {isLoading ? (
                                    <div className="workspace-loading">
                                        <span className="workspace-loading-spinner" />
                                        <span>Loading…</span>
                                    </div>
                                ) : workspaces.length === 0 ? (
                                    <div className="workspace-empty">
                                        No workspaces yet.<br />
                                        Create your first workspace below.
                                    </div>
                                ) : (
                                    workspaces.map((ws) => (
                                        <button
                                            key={ws.id}
                                            className={`workspace-option${ws.id === activeWorkspaceId ? " selected" : ""}`}
                                            onClick={() => handleSelectWorkspace(ws.id)}
                                            id={`workspace-option-${ws.id}`}
                                        >
                                            <span className="workspace-option-dot" />
                                            <span className="workspace-option-name">{ws.name}</span>
                                            {ws.id === activeWorkspaceId && (
                                                <span className="workspace-option-check">✓</span>
                                            )}
                                        </button>
                                    ))
                                )}
                            </div>

                            <div className="workspace-dropdown-divider" />

                            <button
                                className="workspace-create-action"
                                onClick={handleCreateWorkspace}
                                id="workspace-create-action-btn"
                            >
                                <span>+</span>
                                <span>New Workspace</span>
                            </button>
                        </div>
                    </>
                )}
            </div>

            <div className="topbar-spacer"></div>
            <div style={{ display: "flex", gap: "8px" }}>
                <button className="btn btn-ghost">↓ Export</button>
                <button className="btn btn-primary" onClick={() => setOpenNewThread(true)}>+ New Thread</button>
            </div>
        </div>
    );
};
