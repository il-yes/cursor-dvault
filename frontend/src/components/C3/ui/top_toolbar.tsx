import React, { useState, useRef, useEffect } from 'react';
import { useC3WorkspaceStore } from '../infrastructure/store/useC3WorkspaceStore';

interface TopToolbarProps {
    setOpenNewThread?: (open: boolean) => void;
    setOpenCreateWorkspace?: (open: boolean) => void;
    setOpenCreateChannel?: (open: boolean) => void;
}

export default function C3TopToolbar({
    setOpenCreateWorkspace,
    setOpenCreateChannel
}: TopToolbarProps) {
    return (
        <TopToolbar setOpenCreateWorkspace={setOpenCreateWorkspace} setOpenCreateChannel={setOpenCreateChannel} />
    );
}

export const TopToolbar = ({
    setOpenNewThread,
    setOpenCreateWorkspace,
    setOpenCreateChannel
}: TopToolbarProps) => {
    const { workspaces, activeWorkspace, activeWorkspaceId, selectWorkspace, isLoading, error, fetchWorkspaces } = useC3WorkspaceStore();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    // Close dropdown on outside click
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsDropdownOpen(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const handleSelect = (id: string) => {
        selectWorkspace(id);
        setIsDropdownOpen(false);
    };

    const handleOpenCreateModal = () => {
        setIsDropdownOpen(false);
        if (setOpenCreateWorkspace) {
            setOpenCreateWorkspace(true);
        }
    };

    return (
        <div className="topbar">
            <div className="topbar-logo">C3</div>

            {/* Workspace Selector Dropdown Pill */}
            <div ref={dropdownRef} style={{ position: 'relative', display: 'inline-block' }}>
                <button
                    className="workspace-pill"
                    onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                    style={{
                        cursor: 'pointer',
                        userSelect: 'none',
                        border: '1px solid var(--c3-border)',
                        background: 'rgba(245, 245, 245, 0.92)',
                        padding: '5px 12px',
                        borderRadius: '999px',
                        fontSize: '12px',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '6px',
                        color: '#333',
                        fontWeight: 500,
                    }}
                >
                    <span style={{ color: '#888', fontWeight: 400 }}>Workspace:</span>
                    <strong style={{ fontWeight: 600, color: '#1a1a1a' }}>
                        {isLoading ? 'Loading...' : activeWorkspace ? activeWorkspace.name : 'No Workspace Selected'}
                    </strong>
                    <span style={{ fontSize: '10px', color: '#666', marginLeft: '2px' }}>
                        {isDropdownOpen ? '▲' : '▾'}
                    </span>
                </button>

                {/* Dropdown Menu */}
                {isDropdownOpen && (
                    <div
                        style={{
                            position: 'absolute',
                            top: 'calc(100% + 6px)',
                            left: 0,
                            minWidth: '220px',
                            backgroundColor: '#FFFFFF',
                            borderRadius: '8px',
                            border: '1px solid #E5E7EB',
                            boxShadow: '0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.05)',
                            zIndex: 100,
                            overflow: 'hidden',
                            fontFamily: 'Inter, system-ui, sans-serif',
                        }}
                    >
                        <div style={{ padding: '6px 12px', fontSize: '10px', fontWeight: 700, color: '#9CA3AF', textTransform: 'uppercase', letterSpacing: '0.05em', borderBottom: '1px solid #F3F4F6' }}>
                            Workspaces ({workspaces.length})
                        </div>

                        <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
                            {isLoading ? (
                                <div style={{ padding: '10px 12px', fontSize: '12px', color: '#6B7280' }}>
                                    Loading workspaces…
                                </div>
                            ) : error ? (
                                <div style={{ padding: '10px 12px', fontSize: '12px', color: '#EF4444' }}>
                                    <span>Failed to load workspaces</span>
                                    <button
                                        onClick={() => fetchWorkspaces()}
                                        style={{ display: 'block', marginTop: '4px', background: 'none', border: 'none', color: '#2563EB', fontSize: '11px', cursor: 'pointer', padding: 0 }}
                                    >
                                        ↻ Retry
                                    </button>
                                </div>
                            ) : workspaces.length === 0 ? (
                                <div style={{ padding: '10px 12px', fontSize: '12px', color: '#6B7280' }}>
                                    No workspaces available.
                                </div>
                            ) : (
                                workspaces.map((ws) => {
                                    const isSelected = ws.id === activeWorkspaceId;
                                    return (
                                        <div
                                            key={ws.id}
                                            onClick={() => handleSelect(ws.id)}
                                            style={{
                                                padding: '8px 12px',
                                                fontSize: '13px',
                                                color: isSelected ? '#2563EB' : '#1F2937',
                                                backgroundColor: isSelected ? '#EFF6FF' : 'transparent',
                                                fontWeight: isSelected ? 600 : 400,
                                                cursor: 'pointer',
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'space-between',
                                                transition: 'background-color 120ms ease',
                                            }}
                                            onMouseEnter={(e) => {
                                                if (!isSelected) e.currentTarget.style.backgroundColor = '#F9FAFB';
                                            }}
                                            onMouseLeave={(e) => {
                                                if (!isSelected) e.currentTarget.style.backgroundColor = 'transparent';
                                            }}
                                        >
                                            <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                                                {ws.name}
                                            </span>
                                            {isSelected && <span style={{ fontSize: '12px', color: '#2563EB' }}>✓</span>}
                                        </div>
                                    );
                                })
                            )}
                        </div>

                        <div style={{ borderTop: '1px solid #F3F4F6', padding: '4px' }}>
                            <button
                                onClick={handleOpenCreateModal}
                                style={{
                                    width: '100%',
                                    padding: '6px 10px',
                                    background: 'transparent',
                                    border: 'none',
                                    borderRadius: '4px',
                                    color: '#C8922A',
                                    fontSize: '12px',
                                    fontWeight: 600,
                                    cursor: 'pointer',
                                    textAlign: 'left',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: '6px',
                                }}
                                onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#FDF8F0')}
                                onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
                            >
                                <span>+</span> Create Workspace
                            </button>
                        </div>
                    </div>
                )}
            </div>

            <div className="topbar-spacer"></div>

            <div style={{ display: 'flex', gap: '8px' }}>
                <button className="btn btn-ghost">↓ Export</button>
                <button
                    className="btn btn-ghost"
                    onClick={() => setOpenCreateWorkspace && setOpenCreateWorkspace(true)}
                >
                    + Workspace
                </button>
                {setOpenCreateChannel && (
                    <button
                        className="btn btn-primary"
                        onClick={() => setOpenCreateChannel(true)}
                        disabled={!activeWorkspaceId}
                        style={{
                            opacity: activeWorkspaceId ? 1 : 0.6,
                            cursor: activeWorkspaceId ? 'pointer' : 'not-allowed',
                        }}
                    >
                        + New Channel
                    </button>
                )}
                {setOpenNewThread && (
                    <button
                        className="btn btn-ghost"
                        onClick={() => setOpenNewThread(true)}
                        disabled={!activeWorkspaceId}
                        style={{
                            opacity: activeWorkspaceId ? 1 : 0.6,
                            cursor: activeWorkspaceId ? 'pointer' : 'not-allowed',
                        }}
                    >
                        + New Thread
                    </button>
                )}
            </div>
        </div>
    );
};
