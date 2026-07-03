import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { C3styles } from "./styles/styles";



export const C3SidebarMenu = () => {
	const navigate = useNavigate();
	const [init, setInit] = useState(false)
	
	const isC3ContextLedger = location.pathname.startsWith("/dashboard/c3/ledger");
	const isC3ContextInbox = location.pathname.startsWith("/dashboard/c3/inbox");



	return (
		<div>
			<C3styles />
			<div className="sidebar-section-label">Navigation</div>
			<div className={`nav-row ${isC3ContextLedger ? "active" : ""}`} onClick={() => navigate("/dashboard/c3/ledger")}>
				<span className="nav-icon">≡</span>
				<span className="nav-label">Ledger</span>
				<span className="nav-badge-gray">8</span>
			</div>
			<div className={`nav-row ${isC3ContextInbox ? "active" : ""}`} onClick={() => navigate("/dashboard/c3/inbox")}>
				<span className="nav-icon">📥</span>
				<span className="nav-label">Inbox</span>
				<span className="nav-badge">5</span>
			</div>
			<hr className="sidebar-divider" />
			<div className="sidebar-section-label">Vaults</div>
			{init ? <EmptyC3 /> : <FullC3 />}
		</div>
	)
}

const EmptyC3 = () => {
	return (
		<div className="sidebar-empty-hint">
			Vaults will appear here once you start a thread.
		</div>
	)
}
const FullC3 = () => {

	const isC3ContextLedger = location.pathname.startsWith("/dashboard/c3/ledger");
	return (
		<div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#2563EB" }} />
				Finance
				<span className="unread-pip" />
			</div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#059669" }} />
				Treasury
			</div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#7C3AED" }} />
				Legal
				<span className="unread-pip" />
			</div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#D97706" }} />
				HR
			</div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#4F46E5" }} />
				IT
			</div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#DC2626" }} />
				Ops
			</div>
			<div className="vault-row">
				<span className="vault-dot" style={{ background: "#0891B2" }} />
				Compliance
			</div>
			{isC3ContextLedger && <div className="new-channel-btn">+ New Channel</div>}
		</div>
	)
}