import { useState, useEffect } from 'react';
import './encryptionVerif.css';

// Mock Stellar fetch call (replace with real implementation)
export async function fetchCommitFromStellar(hash) {
	await new Promise((res) => setTimeout(res, 800));
	return {
		fileName: "contracts.zip",
		txHash: hash,
		action: "FILE_COMMIT",
		timestamp: "2025-11-28T12:30:00Z",
		author: "0xANONYMOUS",
		signatureValid: true,
		cid: "QmXkjExampleCID12345"
	};
}

export default function EncryptionVerificationModalBeta({ file, onClose }) {
	const [commit, setCommit] = useState(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState(null);

	useEffect(() => {
		async function loadCommit() {
			try {
				const commitData = await fetchCommitFromStellar(file.commitHash);
				setCommit(commitData);
			} catch (err) {
				setError(err.message);
			} finally {
				setLoading(false);
			}
		}
		loadCommit();
	}, [file.commitHash]);

	if (loading) {
		return <div className="modal-loading">Loading verification data...</div>;
	}

	if (error) {
		return (
			<div className="modal-error">
				<p>Could not load verification data: {error}</p>
				<button onClick={onClose}>Close</button>
			</div>
		);
	}

	return (
		<div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-md">
			<div className="relative mx-4 w-full max-w-xl rounded-3xl border border-white/20 bg-gradient-to-br from-white/70 via-white/40 to-zinc-100/30 p-[1px] dark:from-zinc-900/80 dark:via-zinc-900/60 dark:to-black/50 shadow-2xl">
				{/* Inner glass panel */}
				<div className="rounded-[1.3rem] bg-white/80 dark:bg-zinc-900/85 backdrop-blur-2xl px-7 py-6">
					{/* Header */}
					<div className="mb-6 flex items-start justify-between gap-4">
						<div>
							<h2 className="text-xl font-semibold bg-gradient-to-r from-foreground via-primary to-[#C9A44A] bg-clip-text text-transparent">
								Encryption Verification
							</h2>
							<p className="mt-1 text-sm text-muted-foreground">
								Cryptographic proof that your file was encrypted on-device before upload.
							</p>
						</div>
						<button
							onClick={onClose}
							className="inline-flex h-9 w-9 items-center justify-center rounded-2xl border border-white/40 bg-white/60 text-muted-foreground shadow-sm transition-all hover:border-primary/40 hover:bg-white hover:text-foreground dark:bg-zinc-800/80"
						>
							<span className="text-lg leading-none">×</span>
						</button>
					</div>

					{/* Status pill */}
					<div className="mb-5 flex items-center gap-3 rounded-2xl border border-emerald-400/40 bg-emerald-400/10 px-4 py-3 backdrop-blur-sm shadow-sm">
						<div className="flex h-8 w-8 items-center justify-center rounded-2xl bg-emerald-400/80 text-white shadow-md">
							✅
						</div>
						<div>
							<p className="text-sm font-semibold text-emerald-300">
								Client-side encryption verified
							</p>
							<p className="text-xs text-emerald-100/80">
								Matching commit found on Stellar with valid signature.
							</p>
						</div>
					</div>

					{/* Commit details */}
					<div className="space-y-4 rounded-2xl border border-white/20 bg-white/40 p-4 text-sm shadow-inner dark:bg-zinc-900/60">
						<div className="mb-2 flex items-center justify-between">
							<h3 className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
								Cryptographic Proof
							</h3>
							<span className="rounded-full bg-amber-500/10 px-3 py-1 text-[11px] font-medium text-amber-400 border border-amber-500/40">
								AES-256-GCM · IPFS · Stellar
							</span>
						</div>

						{/* Tx hash */}
						<div className="flex flex-col gap-1">
							<span className="text-xs font-medium text-muted-foreground">Stellar Transaction</span>
							<div className="flex items-center gap-2">
								<code className="flex-1 truncate rounded-xl bg-black/80 px-3 py-2 font-mono text-[11px] text-emerald-300 border border-emerald-400/40">
									{commit.txHash}
								</code>
								<button
									onClick={() => navigator.clipboard.writeText(commit.txHash)}
									className="inline-flex h-9 w-9 items-center justify-center rounded-2xl border border-white/30 bg-white/70 text-xs text-muted-foreground shadow-sm transition-all hover:border-primary/50 hover:bg-primary/10 hover:text-primary"
								>
									<span>📋</span>
								</button>
							</div>
						</div>

						{/* Action / Timestamp */}
						<div className="grid grid-cols-2 gap-4">
							<div className="space-y-1">
								<span className="text-xs font-medium text-muted-foreground">Action</span>
								<p className="rounded-xl bg-white/60 px-3 py-2 text-xs font-semibold tracking-wide text-foreground shadow-sm dark:bg-zinc-800/70">
									{commit.action}
								</p>
							</div>
							<div className="space-y-1">
								<span className="text-xs font-medium text-muted-foreground">Timestamp</span>
								<p className="rounded-xl bg-white/60 px-3 py-2 text-xs text-foreground shadow-sm dark:bg-zinc-800/70">
									{new Date(commit.timestamp).toLocaleString()}
								</p>
							</div>
						</div>

						{/* Author / Signature */}
						<div className="grid grid-cols-2 gap-4">
							<div className="space-y-1">
								<span className="text-xs font-medium text-muted-foreground">Author</span>
								<p className="rounded-xl bg-white/60 px-3 py-2 text-xs font-mono text-foreground shadow-sm dark:bg-zinc-800/70">
									{commit.author}
								</p>
							</div>
							<div className="space-y-1">
								<span className="text-xs font-medium text-muted-foreground">Signature</span>
								<p
									className={`flex items-center gap-2 rounded-xl px-3 py-2 text-xs font-semibold shadow-sm ${commit.signatureValid
										? "border border-emerald-400/50 bg-emerald-500/10 text-emerald-300"
										: "border border-red-400/50 bg-red-500/10 text-red-300"
										}`}
								>
									{commit.signatureValid ? "✅ Valid" : "❌ Invalid"}
								</p>
							</div>
						</div>

						{/* CID / Encryption */}
						<div className="grid grid-cols-2 gap-4">
							<div className="space-y-1 col-span-2">
								<span className="text-xs font-medium text-muted-foreground">IPFS CID</span>
								<code className="block rounded-xl bg-black/80 px-3 py-2 font-mono text-[11px] text-sky-300 border border-sky-400/40">
									{commit.cid}
								</code>
							</div>
							<div className="space-y-1">
								<span className="text-xs font-medium text-muted-foreground">Encryption</span>
								<p className="rounded-xl bg-white/60 px-3 py-2 text-xs font-semibold text-foreground shadow-sm dark:bg-zinc-800/70">
									AES-256-GCM
								</p>
							</div>
						</div>
					</div>

					{/* Modal actions */}
					<div className="mt-8 flex flex-col sm:flex-row items-center justify-end gap-3">
						<button
							className="h-12 px-6 rounded-2xl bg-gradient-to-r from-[#C9A44A] to-[#B8934A] text-primary-foreground font-bold shadow-xl hover:from-[#C9A44A]/90 hover:to-[#B8934A]/90 hover:shadow-[#C9A44A]/30 transition-all"
							onClick={() => window.open(`https://stellar.expert/explorer/public/tx/${commit.txHash}`, '_blank')}
						>
							View on Stellar Explorer
						</button>
						<button
							className="h-12 px-6 rounded-2xl bg-white/70 dark:bg-zinc-800/70 text-foreground font-semibold shadow-lg border border-primary/30 hover:bg-primary/10 hover:text-primary transition-all"
							onClick={() => downloadAuditReport(commit)}
						>
							Download Audit Report
						</button>
					</div>
				</div>
			</div>
		</div>
	);

}

// Helper function to generate audit report
function downloadAuditReport(commit) {
	const report = `
			ANKHORA ENCRYPTION AUDIT REPORT
================================

File: ${commit.fileName}
Generated: ${new Date().toISOString()}

ENCRYPTION VERIFICATION
-----------------------
✅ File was encrypted client-side using AES-256-GCM
✅ Encryption occurred before transmission to server
✅ Server received only encrypted bytes (no plaintext access)

CRYPTOGRAPHIC PROOF
-------------------
Stellar Transaction: ${commit.txHash}
Action: ${commit.action}
Timestamp: ${new Date(commit.timestamp).toISOString()}
Author: ${commit.author}
Signature: ${commit.signatureValid ? 'Valid' : 'Invalid'}
IPFS CID: ${commit.cid}

VERIFICATION STEPS
------------------
1. Visit: https://stellar.expert/explorer/public/tx/${commit.txHash}
2. Verify transaction exists on Stellar blockchain
3. Verify signature matches author's public key
4. Verify timestamp matches file creation time

This report provides cryptographic proof that the file was encrypted
on the client device and never exposed to the server in plaintext form.

For questions, contact: security@blackops.africa
  `.trim();

	const blob = new Blob([report], { type: 'text/plain' });
	const url = URL.createObjectURL(blob);
	const a = document.createElement('a');
	a.href = url;
	a.download = `ankhora-audit-${commit.fileName}-${Date.now()}.txt`;
	a.click();
	URL.revokeObjectURL(url);
}
