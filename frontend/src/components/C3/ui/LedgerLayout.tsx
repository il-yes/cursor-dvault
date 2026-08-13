import React, { useState, useEffect } from "react";
import { NewThreadAssetDrawer } from "./thread/NewThreadAssetDrawer";
import { TopToolbar } from "./top_toolbar";
import { CreateWorkspaceModal } from "../CreateWorkspaceModal";
import { CreateChannelModal } from "../CreateChannelModal";
import { CreateThreadModal } from "../CreateThreadModal";
import { WorkspaceResponse, ChannelResponse, ThreadResponse } from "@/services/api";
import { useC3WorkspaceStore } from "../infrastructure/store/useC3WorkspaceStore";
import { useC3ChannelStore } from "../infrastructure/store/useC3ChannelStore";
import { useC3ThreadStore } from "../infrastructure/store/useC3ThreadStore";

export const LedgerLayout = ({
    children,
    isNewShareOpen
}: {
    children: React.ReactNode;
    isNewShareOpen: boolean;
}) => {
    const [openNewThread, setOpenNewThread] = useState(false);
    const [openCreateWorkspace, setOpenCreateWorkspace] = useState(false);
    const [openCreateChannel, setOpenCreateChannel] = useState(false);

    const {
        workspaces,
        activeWorkspace,
        activeWorkspaceId,
        isLoading: isWorkspaceLoading,
        error: workspaceError,
        fetchWorkspaces,
        addWorkspace
    } = useC3WorkspaceStore();

    const {
        fetchChannels,
        addChannel,
        activeChannel,
        activeChannelId,
        isLoading: isChannelLoading,
        error: channelError
    } = useC3ChannelStore();

    const {
        fetchThreads,
        addThread,
        isLoading: isThreadLoading,
        error: threadError
    } = useC3ThreadStore();

    // Load workspaces from backend on mount
    useEffect(() => {
        fetchWorkspaces();
    }, [fetchWorkspaces]);

    // Automatically fetch channels whenever activeWorkspaceId changes (Workspace -> Channel Invariant)
    useEffect(() => {
        if (activeWorkspaceId) {
            fetchChannels(activeWorkspaceId);
        }
    }, [activeWorkspaceId, fetchChannels]);

    // Automatically fetch threads whenever activeChannelId changes (Channel -> Thread Invariant)
    useEffect(() => {
        if (activeChannelId) {
            fetchThreads(activeChannelId);
        }
    }, [activeChannelId, fetchThreads]);

    const handleWorkspaceCreated = (workspace: WorkspaceResponse) => {
        addWorkspace(workspace);
    };

    const handleChannelCreated = (channel: ChannelResponse) => {
        addChannel(channel);
    };

    const handleThreadCreated = (thread: ThreadResponse) => {
        addThread(thread);
    };

    const combinedError = workspaceError || channelError || threadError;

    return (
        <>
            <TopToolbar
                setOpenNewThread={setOpenNewThread}
                setOpenCreateWorkspace={setOpenCreateWorkspace}
                setOpenCreateChannel={setOpenCreateChannel}
            />

            {/* Error Notification / Retry Bar */}
            {combinedError && (
                <div style={{
                    backgroundColor: 'rgba(239, 68, 68, 0.12)',
                    borderBottom: '1px solid rgba(239, 68, 68, 0.3)',
                    padding: '8px 20px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    fontSize: '13px',
                    color: '#EF4444',
                }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                        <span>⚠️</span>
                        <span>{combinedError}</span>
                    </div>
                    <button
                        onClick={() => {
                            fetchWorkspaces();
                            if (activeWorkspaceId) fetchChannels(activeWorkspaceId);
                            if (activeChannelId) fetchThreads(activeChannelId);
                        }}
                        style={{
                            background: '#EF4444',
                            color: '#FFFFFF',
                            border: 'none',
                            borderRadius: '4px',
                            padding: '4px 10px',
                            fontSize: '12px',
                            fontWeight: 600,
                            cursor: 'pointer',
                        }}
                    >
                        Retry
                    </button>
                </div>
            )}

            {/* Activation toast */}
            {isNewShareOpen && <div className="toast-bar">
                <span className="toast-icon">⚡</span>
                <span className="toast-text">
                    Channel <strong>contract-execution</strong> activated —{" "}
                    <strong>vault_legal</strong> can now commit{" "}
                    <strong>contract_draft</strong>. Stellar genesis tx anchored.
                </span>
                <span className="toast-action">Go to inbox →</span>
                <span className="toast-close">✕</span>
            </div>}

            {/* Main layout */}
            <div className="layout">
                {!isWorkspaceLoading && workspaces.length === 0 && !workspaceError ? (
                    <div style={{
                        flex: 1,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: 'center',
                        padding: '40px 20px',
                        textAlign: 'center',
                        color: '#6B7280',
                    }}>
                        <div style={{ fontSize: '32px', marginBottom: '12px' }}>📂</div>
                        <h3 style={{ fontSize: '16px', fontWeight: 600, color: '#374151', margin: '0 0 6px 0' }}>
                            No workspaces yet
                        </h3>
                        <p style={{ fontSize: '13px', color: '#6B7280', margin: '0 0 16px 0', maxWidth: '360px' }}>
                            Create your first C3 workspace to organize channels, threads, and collaboration assets.
                        </p>
                        <button
                            className="btn btn-primary"
                            onClick={() => setOpenCreateWorkspace(true)}
                            style={{
                                padding: '8px 18px',
                                fontSize: '13px',
                                borderRadius: '6px',
                                fontWeight: 600,
                            }}
                        >
                            + Create Your First Workspace
                        </button>
                    </div>
                ) : (
                    children
                )}

                {/* New Thread Drawer */}
                <NewThreadAssetDrawer
                    open={false}
                    onClose={() => setOpenNewThread(false)}
                />
            </div>

            {/* Create Workspace Modal */}
            <CreateWorkspaceModal
                isOpen={openCreateWorkspace}
                onClose={() => setOpenCreateWorkspace(false)}
                onWorkspaceCreated={handleWorkspaceCreated}
            />

            {/* Create Channel Modal */}
            <CreateChannelModal
                isOpen={openCreateChannel}
                activeWorkspaceId={activeWorkspaceId}
                activeWorkspaceName={activeWorkspace?.name}
                onClose={() => setOpenCreateChannel(false)}
                onChannelCreated={handleChannelCreated}
            />

            {/* Create Thread Modal */}
            <CreateThreadModal
                isOpen={openNewThread}
                activeWorkspaceName={activeWorkspace?.name}
                activeChannelId={activeChannelId}
                activeChannelTitle={activeChannel?.title}
                onClose={() => setOpenNewThread(false)}
                onThreadCreated={handleThreadCreated}
            />
        </>
    );
};