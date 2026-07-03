import React from 'react'

export default function C3TopToolbar() {
    return (
        <div className="topbar">
            <div className="topbar-logo">ANKHORA</div>
            <div className="workspace-pill">Workspace: Acme Corp ▾</div>
            <div className="topbar-spacer" />
            <div style={{ display: "flex", gap: 8 }}>
                <button className="btn btn-ghost">↓ Export</button>
                <button className="btn btn-primary">+ New Thread</button>
            </div>
        </div>
    )
}

export const TopToolbar = ({ setOpenNewThread }: { setOpenNewThread: (open: boolean) => void }) => {
    return (
        <div className="topbar">
            <div className="topbar-logo">C3</div>
            <div className="workspace-pill">Workspace: Acme Corp ▾</div>
            <div className="topbar-spacer"></div>
            <div style={{ display: "flex", gap: "8px" }}>
                <button className="btn btn-ghost">↓ Export</button>
                <button className="btn btn-primary" onClick={() => setOpenNewThread(true)}>+ New Thread</button>
            </div>
        </div>
    )
}
