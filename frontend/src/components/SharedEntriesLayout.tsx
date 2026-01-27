import { useState, useMemo, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { SharedEntry, ShareFilter, DetailView } from "@/types/sharing";
import { SharedEntriesList } from "@/components/SharedEntriesList";
import { SharedEntryOverview } from "@/components/SharedEntryOverview";
import { SharedEntryDetails } from "@/components/SharedEntryDetails";
import { useVaultStore } from "@/store/vaultStore";
import "./contributionGraph/g-scrollbar.css";

import QRCode from "react-qr-code"; // npm install qrcode.react
import { NewLinkShareModal } from "./NewLinkShare";

const linkSharesMocked = [
	{
		id: "1",
		entry_name: "aws keys",
		status: "active", // or "expired", "revoked"
		expiry: "2026-01-28",
		uses_left: 2,
		link: "https://ankhora.app/share/abcdef123456",
		audit_log: [/* ... */],
	},
	{
		id: "2",
		entry_name: "google keys",
		status: "active", // or "expired", "revoked"
		expiry: "2026-01-28",
		uses_left: 2,
		link: "https://ankhora.app/share/abcdef123456",
		audit_log: [/* ... */],
	},
	{
		id: "3",
		entry_name: "github keys",
		status: "expired", // or "expired", "revoked"
		expiry: "2025-12-28",
		uses_left: 0,
		link: "https://ankhora.app/share/abcdef123456",
		audit_log: [/* ... */],
	},
	// ...
];

export function SharedEntriesLayout() {
	// const sharedEntries = useVaultStore((state) => state.shared.items);
	const [isNewShareOpen, setIsNewShareOpen] = useState(false);

	const sharedByMe = useVaultStore(state => state.shared.items); 	// Cryptographic Share, by me
	const sharedWithMe = useVaultStore(state => state.sharedWithMe.items);
	const [searchParams] = useSearchParams();
	const filterParam = (searchParams.get("filter") || "all") as ShareFilter;
	const [filter, setFilter] = useState<ShareFilter>(filterParam);
	const [detailView, setDetailView] = useState<DetailView>("recipients");
	const [refreshKey, setRefreshKey] = useState(0);
	const [sharedEntriesRefreshKey, setSharedEntriesRefreshKey] = useState(0);

	const shareTypeParam = (searchParams.get("type") || "linkshare") as "linkshare" | "cryptographicshare";
	const scopeParam = (searchParams.get("scope") || "byme") as "byme" | "withme";

	// Data sources
	const sharedByMeLink = useVaultStore(state => state.linkSharedByMe.items); // for Link Share, by me
	const sharedWithMeLink = useVaultStore(state => state.linkSharedWithMe.items); // for Link Share, with me
	const sharedByMeCrypto = useVaultStore(state => state.shared.items); // for Cryptographic Share, by me
	const sharedWithMeCrypto = useVaultStore(state => state.sharedWithMe.items); // for Cryptographic Share, with me

	// Pick the correct entries based on params
	const sharedEntries = useMemo(() => {
		if (shareTypeParam !== "linkshare") {
			return scopeParam === "byme" ? sharedByMeCrypto : sharedWithMeCrypto;
		}
	}, [shareTypeParam, scopeParam, sharedByMeLink, sharedWithMeLink, sharedByMeCrypto, sharedWithMeCrypto]);

	useEffect(() => {
		setFilter(filterParam);
	}, [filterParam]);

	useEffect(() => {
		const handleRefresh = () => {
			setRefreshKey(prev => prev + 1);
		};

		window.addEventListener('shareEntriesRefresh', handleRefresh);
		return () => window.removeEventListener('shareEntriesRefresh', handleRefresh);
	}, []);

	useEffect(() => {
		if (!selectedEntry) return;

		const fresh = sharedByMe.find(e => e.id === selectedEntry.id);
		if (fresh) {
			setSelectedEntry(fresh);
		}
	}, [sharedByMe]);

	useEffect(() => {
		const handleRefresh = () => {
			setRefreshKey(prev => prev + 1);
		};

		return () => window.removeEventListener('shareEntriesRefresh', handleRefresh);
	}, [detailView]);

	// Filter entries
	const filteredEntries = useMemo(() => {
		let filtered = sharedEntries;

		switch (filterParam) {
			case "sent":
				// In real app, filter by entries shared by current user
				filtered = filtered.filter(e => e.status === "active" || e.status === "pending");
				break;
			case "received":
				// In real app, filter by entries received by current user
				filtered = filtered.filter(e => e.status === "active");
				break;
			case "pending":
				filtered = filtered.filter(e => e.status === "pending");
				break;
			case "revoked":
				filtered = filtered.filter(e => e.status === "revoked");
				break;
			case "all":
			default:
				// Show all entries
				break;
		}

		return filtered;
	}, [filterParam, sharedEntries, scopeParam, shareTypeParam]);

	// Reset selectedEntry when data source changes
	const [selectedEntry, setSelectedEntry] = useState<SharedEntry | null>(null);
	useEffect(() => {
		setSelectedEntry(null);
	}, [shareTypeParam, scopeParam, filterParam]);

	return (
		<div className="flex h-full" key={shareTypeParam + scopeParam + filterParam}>
			{shareTypeParam === "cryptographicshare" &&
				<>
					{/* Column 2: Shared Entries List (Column 1 is the main sidebar) */}
					<div className="w-full md:w-80 lg:w-96 flex flex-col border-r border-border bg-secondary/30 overflow-hidden">
						<SharedEntriesList
							entries={filteredEntries}
							selectedEntryId={selectedEntry?.id || null}
							onSelectEntry={setSelectedEntry}
						/>
					</div>

					{/* Column 3: Entry Overview */}
					<div className="hidden md:flex flex-col w-80 lg:w-96 border-r border-border overflow-hidden">
						<div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar">
							<SharedEntryOverview
								entry={selectedEntry}
								onViewChange={setDetailView}
							/>
						</div>
					</div>

					{/* Column 4: Detail Panel (but visually column 3) */}
					<div className="flex-1 hidden lg:flex flex-col overflow-hidden">
						<div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar">
							<SharedEntryDetails
								entry={selectedEntry}
								view={detailView}
							/>
						</div>
					</div>
				</>
			}

			{shareTypeParam === "linkshare" && <LinkShareContent linkShares={sharedByMeLink} onAddLinkShare={() => setIsNewShareOpen(true)} />}

			<NewLinkShareModal
				open={isNewShareOpen}
				onOpenChange={setIsNewShareOpen}
				onShareSuccess={() => {
					setSharedEntriesRefreshKey(prev => prev + 1);
					window.dispatchEvent(new CustomEvent('shareEntriesRefresh'));
				}}
			/>
		</div>
	);
}


const LinkShareContent = ({ linkShares, onAddLinkShare }) => {
	const [qrShare, setQrShare] = useState(null);

	const onCopyLink = (share) => {
		navigator.clipboard.writeText(share.link);
		setQrShare(share);
	};

	const onRevoke = (share) => {
		alert(`Revoke link: ${share.link}`);
	};

	const onViewAuditLog = (share) => {
		alert(`View audit log for: ${share.entry_name}`);
	};


	return (
		<div className="w-full bg-secondary/30 p-7 ">
			{/* Banner */}
			<div className="mb-10 px-4 py-3 rounded-xl bg-gradient-to-r from-[#23211b]/80 to-[#c9a44a]/10 text-[13px] text-[#C9A44A] font-medium shadow">
				Links are cryptographically signed and tracked. For sensitive data, use Cryptographic Shares for enhanced protection.
			</div>
			{/* Title + CTA – moved ABOVE banner so it’s always visible */}
			<div className="mb-3 flex flex-wrap items-center justify-between gap-4">
				<div className="mb-10">
					<h2 className="text-xl font-semibold bg-gradient-to-r from-[#23211b] via-[#C9A44A] to-amber-400 bg-clip-text text-transparent">
						Link Shares
					</h2>
					<p className="mt-1 text-xs text-slate-400">
						Create time‑boxed links to entries while preserving full cryptographic traceability.
					</p>
				</div>
				<button
					onClick={onAddLinkShare}
					className="inline-flex items-center gap-2 rounded-full bg-gradient-to-r from-[#C9A44A] via-amber-400 to-[#B8934A] px-4 py-2 text-xs font-semibold text-black shadow-[0_14px_40px_rgba(0,0,0,0.8)] hover:shadow-[0_18px_60px_rgba(0,0,0,0.9)] hover:-translate-y-0.5 transition-all"
				>
					+ Add Link Share
				</button>
			</div>


			{/* Table */}
			<table className="min-w-full table-auto border-separate border-spacing-y-3">
				<thead>
					<tr className="text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-300">
						<th className="px-4 py-2 text-left w-[28%]">Entry</th>
						<th className="px-4 py-2 text-left w-[14%]">Status</th>
						<th className="px-4 py-2 text-left w-[18%]">Expiry</th>
						<th className="px-4 py-2 text-left w-[12%]">Uses Left</th>
						<th className="px-4 py-2 text-left w-[14%]">Audit Log</th>
						<th className="px-4 py-2 text-left w-[14%]">Actions</th>
					</tr>
				</thead>
				<tbody>
					{linkShares.map((share, i) => (
						<tr key={share.id}>
							<td colSpan={6}>
								<div className="px-2 flex items-center rounded-2xl border border-white/15 bg-white/10 dark:bg-zinc-900/40 backdrop-blur-2xl shadow transition-all hover:shadow-lg">
									<div className="flex w-full items-center">

										{/* Entry */}
										<div className="w-[28%] px-0 py-3 text-sm font-medium">
											{share.entry_name}
										</div>

										{/* Status */}
										<div className="w-[14%] px-0 py-3">
											<span
												className={`inline-flex items-center rounded-full px-3 py-1 text-[11px] font-semibold backdrop-blur-md border ${share.status === "active"
														? "bg-emerald-500/15 text-emerald-300 border-emerald-400/40"
														: share.status === "expired"
															? "bg-slate-500/15 text-slate-300 border-slate-400/40"
															: "bg-red-500/15 text-red-300 border-red-400/40"
													}`}
											>
												{share.status}
											</span>
										</div>

										{/* Expiry */}
										<div className="w-[18%] px-0 py-3 text-xs">
											{share.expiry && share.expiry !== "never"
												? new Date(share.expiry).toLocaleDateString()
												: "Never"}
										</div>

										{/* Uses left */}
										<div className="w-[12%] px-0 py-3 text-xs text-center">
											{share.uses_left === -1 ? "∞" : share.uses_left}
										</div>

										{/* Audit log */}
										<div className="w-[14%] px-0 py-3 pl-4">
											<button
												className="text-xs font-semibold text-[#C9A44A] underline-offset-2 hover:underline hover:text-amber-300"
												onClick={() => onViewAuditLog(share)}
											>
												View
											</button>
										</div>

										{/* Actions */}
										<div className="w-[14%] px-0 py-3 flex gap-2 justify-end">
											<button
												className="px-3 py-1.5 rounded-full bg-white/10 text-[11px] font-semibold text-[#C9A44A] border border-[#C9A44A]/40 backdrop-blur-md hover:bg-[#C9A44A]/15 hover:border-[#C9A44A]/70 transition-all"
												onClick={() => onCopyLink(share)}
											>
												Copy Link
											</button>

											<button
												className="px-3 py-1.5 rounded-full bg-red-500/15 text-[11px] font-semibold text-red-300 border border-red-400/50 backdrop-blur-md hover:bg-red-500/25 transition-all"
												onClick={() => onRevoke(share)}
											>
												Revoke
											</button>
										</div>

									</div>
								</div>
							</td>
						</tr>
					))}

				</tbody>
			</table>




			{/* QR Code Modal */}
			{qrShare && (
				<div className="fixed inset-0 flex items-center justify-center z-50 bg-black/40">
					<div className="bg-white rounded-xl p-6 shadow-xl flex flex-col items-center">
						<QRCode value={qrShare.link} size={180} />
						<div className="mt-3 text-xs font-semibold">
							{qrShare.entry_name}
						</div>
						<button
							className="mt-4 px-4 py-2 bg-[#C9A44A] text-white rounded-lg font-medium"
							onClick={() => setQrShare(null)}
						>
							Close
						</button>
					</div>
				</div>
			)}
		</div>
	);
};

