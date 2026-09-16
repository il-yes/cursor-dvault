import React from "react";
import { CardEntry, IdentityEntry, LoginEntry, NoteEntry, SSHKeyEntry } from "@/types/vault";

interface ClassicalEntryRendererProps {
	entry: any;
	isRevealed?: boolean;
	revealedValue?: string;
	fullData?: any;
	RenderAttachments?: React.ComponentType<{ entry?: { attachments?: any[] } }>;
}

export const ClassicalEntryRenderer: React.FC<ClassicalEntryRendererProps> = ({
	entry,
	isRevealed = false,
	revealedValue,
	fullData,
	RenderAttachments,
}) => {
	if (!entry || typeof entry !== "object") {
		return (
			<div className="text-xs font-mono">
				{isRevealed ? revealedValue ?? JSON.stringify(entry, null, 2) : "••••••••••••"}
			</div>
		);
	}

	const type = (entry?.type || entry?.entry_type || entry?.record_type || "").toLowerCase();

	switch (type) {
		case "login": {
			const login = entry as LoginEntry;
			if (fullData?.attachments) {
				login.attachments = fullData.attachments;
			}
			return (
				<div className="space-y-3">
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">Username:</span>
						<span className="text-lg font-bold">
							{isRevealed ? login.user_name : "••••••••"}
						</span>
					</div>
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">Website:</span>
						<span className="text-sm">{login.web_site || "—"}</span>
					</div>
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">Password:</span>
						<span className="text-lg font-bold">
							{isRevealed ? login.password : "••••••••••••"}
						</span>
					</div>
					{fullData?.attachements?.length > 0 && RenderAttachments && <RenderAttachments entry={fullData} />}
				</div>
			);
		}

		case "card": {
			const card = entry as CardEntry;
			const last4 = card.number?.slice(-4) ?? "••••";
			return (
				<div className="space-y-3">
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">Owner:</span>
						<span className="text-lg font-bold">{card.owner}</span>
					</div>
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">Card:</span>
						<span className="text-lg font-bold">
							•••• •••• •••• {isRevealed ? last4 : "••••"}
						</span>
					</div>
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">Expires:</span>
						<span className="text-lg font-bold">{card.expiration}</span>
					</div>
					<div className="flex items-center gap-2">
						<span className="text-sm font-semibold text-muted-foreground">CVC:</span>
						<span className="text-lg font-bold">
							{isRevealed ? card.cvc : "•••"}
						</span>
					</div>
				</div>
			);
		}

		case "identity": {
			const identity = entry as IdentityEntry;
			return (
				<div className="space-y-2">
					<div className="grid grid-cols-2 gap-4 text-sm">
						<div>
							<span className="font-semibold text-muted-foreground">First Name:</span>
							<span className="ml-2">{identity.firstname || "—"}</span>
						</div>
						<div>
							<span className="font-semibold text-muted-foreground">Last Name:</span>
							<span className="ml-2">{identity.lastname || "—"}</span>
						</div>
					</div>
					<div>
						<span className="font-semibold text-muted-foreground">Company:</span>
						<span className="ml-2">{identity.company || "—"}</span>
					</div>
					<div className="flex flex-wrap gap-4 text-sm">
						<span>
							<span className="font-semibold text-muted-foreground">Email:</span>
							<span className="ml-2">{identity.mail || "—"}</span>
						</span>
						<span>
							<span className="font-semibold text-muted-foreground">Phone:</span>
							<span className="ml-2">{identity.telephone || "—"}</span>
						</span>
					</div>
					<div>
						<span className="font-semibold text-muted-foreground">Address:</span>
						<span className="ml-2">
							{identity.address_one} {identity.address_two} {identity.city},{" "}
							{identity.state} {identity.postal_code}
						</span>
					</div>
				</div>
			);
		}

		case "note": {
			const note = entry as NoteEntry;
			const attachments = fullData?.attachments ?? note.attachments ?? [];
			const text = (note as any).note ?? (note as any).content ?? (note as any).additionnal_note;

			return (
				<div className="space-y-3">
					<div className="text-sm">
						{isRevealed ? text || "No content" : "••••••••••••"}
					</div>
					{RenderAttachments && <RenderAttachments entry={{ attachments }} />}
				</div>
			);
		}

		case "sshkey": {
			const ssh = entry as SSHKeyEntry;
			return (
				<div className="space-y-2 text-sm">
					<div>
						<span className="font-semibold text-muted-foreground">Public Key:</span>
						<span className="ml-2 block font-mono bg-muted/50 p-2 rounded text-xs">
							{ssh.public_key}
						</span>
					</div>
					<div>
						<span className="font-semibold text-muted-foreground">Fingerprint:</span>
						<span className="ml-2">{ssh.e_fingerprint}</span>
					</div>
					<div>
						<span className="font-semibold text-muted-foreground">Private Key:</span>
						<span className="ml-2 block font-mono bg-muted/50 p-2 rounded text-xs">
							{isRevealed ? ssh.private_key : "••••••••••••"}
						</span>
					</div>
				</div>
			);
		}

		default:
			return (
				<div className="text-xs font-mono">
					{isRevealed ? revealedValue ?? JSON.stringify(entry, null, 2) : "••••••••••••"}
				</div>
			);
	}
};
