import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/sidebar-style.css"
import { useC3DialogStore } from "../../infrastructure/store/c3DialogStore";
import { Department } from "../../domain/channel/channel.types";
import { fetchDepartments } from "../../domain/channel/channel.repository";
import * as ROUTES from '@/constants/routes';

type FullC3Props = {
	setOpen: (open: boolean) => void;
	vaults: Department[];
};

export const C3SidebarMenu = () => {
	const navigate = useNavigate();
	const [departments, setDepartments] = useState<Department[]>([])
	const { openC3CreateDialog, channelId } = useC3DialogStore();
	const isC3ContextLedger = location.pathname.startsWith(ROUTES.LEDGER);
	const isC3ContextInbox = location.pathname.startsWith(ROUTES.INBOX);

	useEffect(() => {
		 getDepartments()
	}, [])

	const getDepartments = async () => {
		try {
			const res = await fetchDepartments()
			console.log({ res })
			setDepartments(res)
		} catch (err) {
			console.error("Failed to fetch departments", err)
		}
	}

	return (
		<div>
			<div className="sidebar-section-label">Navigation</div>
			<div className={`nav-row ${isC3ContextLedger ? "active" : ""}`} onClick={() => navigate(ROUTES.LEDGER)}>
				<span className="nav-icon">≡</span>
				<span className="nav-label">Ledger</span>
				<span className="nav-badge-gray">8</span>
			</div>
			<div className={`nav-row ${isC3ContextInbox ? "active" : ""}`} onClick={() => navigate(ROUTES.INBOX)}>
				<span className="nav-icon">📥</span>
				<span className="nav-label">Inbox</span>
				<span className="nav-badge">5</span>
			</div>
			<hr className="sidebar-divider" />
			<div className="sidebar-section-label">Vaults</div>
			{departments.length > 0 ? <FullC3 setOpen={openC3CreateDialog} vaults={departments} /> : <EmptyC3 /> }
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
const FullC3 = ({ setOpen, vaults }: FullC3Props) => {
	const isC3ContextInbox = location.pathname.startsWith(ROUTES.INBOX);
	return (
		<div>
			{vaults.map((vault, index) => (
				<div className="vault-row" key={index}>
					<span className="vault-dot" style={{ background: vault.color}} />
					{vault.name}
					<span className="unread-pip" /> {/* TODO: check user connexion sight */}
				</div>
			))}
			{!isC3ContextInbox && <div className="new-channel-btn" onClick={() => setOpen(true)}>+ New Channel</div>}

		</div>
	)
}