import { SharedEntry, DetailView, Recipient, AuditEvent } from "@/types/sharing";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { UserPlus, Trash2, Download, Shield, Database, Calendar, FolderOpen, Eye, EyeOff, Copy, Sparkles } from "lucide-react";
import { useEffect, useState } from "react";
import { useToast } from "@/hooks/use-toast";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import { useVaultStore } from "@/store/vaultStore";
import { cn } from "@/lib/utils";
import { decryptField, logAuditEvent } from "@/services/api";
import * as AppAPI from "../../wailsjs/go/main/App";
import { Keypair } from "stellar-sdk";
import { Buffer } from 'buffer';
import { useAppStore } from "@/store/appStore";
import { Textarea } from "./ui/textarea";
import { CardEntry, IdentityEntry, LoginEntry, NoteEntry, SSHKeyEntry } from "@/types/vault";


interface SharedEntryDetailsProps {
	entry: SharedEntry | null;
	view: DetailView;
}

interface RevealedField {
	name: string;
	value: string;
	timeout: NodeJS.Timeout;
}

const DEFAULT_REVEAL_TIMEOUT = 15;

export function SharedEntryDetails({ entry, view }: SharedEntryDetailsProps) {
	const { toast } = useToast();
	const [showAddForm, setShowAddForm] = useState(false);
	const [newRecipient, setNewRecipient] = useState<{ name: string; email: string; role: "viewer" | "editor" | "owner" | "read" }>({
		name: "",
		email: "",
		role: "viewer"
	});
	const updateRecipients = useVaultStore((state) => state.updateSharedEntryRecipients);

	const [recipients, setRecipients] = useState<Recipient[]>([]);
	const [revealedFields, setRevealedFields] = useState<Map<string, RevealedField>>(new Map());
	const [isRevealing, setIsRevealing] = useState<string | null>(null);
	const [decryptingField, setDecryptingField] = useState<string | null>(null);
	const { vault, clearVault: clearVaultStore } = useVaultStore();
	const stellar = vault?.vault_runtime_context?.UserConfig?.stellar_account

	const session = useAppStore.getState().session;
	const user = session?.user

	useEffect(() => {
		setRecipients(entry?.recipients || []);
	}, [entry?.id]);

	const getPublicKey = async (email: string) => {
		const response = await fetch(`http://164.90.213.173:4001/api/check-email?email=${email}`, {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
		});
		const data = await response.json();
		console.log({ data });
		return data.public_key;
	};

	const handleAddRecipient = async () => {
		// TODO: get public key from cloud backend
		const publicKey = await getPublicKey(newRecipient.email);
		if (!publicKey) {
			toast({
				title: "Error",
				description: "No public key found for this email",
			});
			return;
		}

		const recipient: Recipient = {
			id: `rec-${Date.now()}`,
			share_id: entry.id,
			public_key: publicKey,
			...newRecipient,
			joined_at: new Date().toISOString(),
		};

		// send AddRecipientRequest payload to Ankhora cloud backend

		// zustand store update
		setRecipients(prev => {
			const updated = [...prev, recipient];
			updateRecipients(entry.id, updated);   // correct sync
			return updated;
		});

		setNewRecipient({ name: "", email: "", role: "viewer" });
		setShowAddForm(false);

		toast({
			title: "Recipient added",
			description: `${recipient.name} has been added as a ${recipient.role}`,
		});

	};

	if (!entry) {
		return (
			<div className="flex flex-col items-center justify-center h-full text-center p-8 bg-gradient-to-b from-background to-secondary/20">
				<div className="relative">
					<img src={ankhoraLogo} alt="Ankhora Logo" className="w-32 h-auto mx-auto" />
					{/* <Shield className="h-20 w-20 text-primary/20 mb-4" /> */}
					{/* <Sparkles className="h-8 w-8 text-primary/40 absolute -top-2 -right-2 animate-pulse" /> */}
				</div>
				<p className="text-sm mt-4 text-muted-foreground max-w-xs italic" style={{ opacity: "0.8" }}>
					Select an action to view details.
				</p>
			</div>
		);
	}

	const handleRemoveRecipient = (id: string) => {
		const recipient = recipients.find(r => r.id === id);

		// call Ankhora cloud backend to remove recipient
		// send RevokeShareRequest payload to Ankhora cloud backend

		// zustand store update
		setRecipients(recipients.filter(r => r.id !== id));

		toast({
			title: "Recipient removed",
			description: `${recipient?.name} has been removed`,
		});
	};

	const handleChangeRole = (email: string, newRole: "viewer" | "read" | "editor" | "owner") => {
		setRecipients(recipients.map(r => r.email === email ? { ...r, role: newRole } : r));
		// call Ankhora cloud backend to update role
		// send UpdateRecipientRequest payload to Ankhora cloud backend

		// zustand store update
		updateRecipients(entry.id, recipients);

		toast({
			title: "Role updated",
			description: "Recipient role has been changed",
		});
	};

	const handleRevealField = async (fieldName: string) => {
		if (!entry) return;

		setIsRevealing(fieldName);
		setDecryptingField(fieldName);

		const rowRecipient = getRowRecipient();

		if (!rowRecipient) return;

		try {
			// 1. Generate keypair
			const keypair = Keypair.fromSecret(stellar.private_key);

			// 2. Request challenge from backend
			const { challenge } = await AppAPI.RequestChallenge({ public_key: stellar.public_key });

			// 3. Sign challenge
			const signature = Buffer.from(
				keypair.sign(Buffer.from(challenge))
			).toString("base64");

			const { plaintext, expires_in } = await decryptField({
				entry_id: entry.id,
				field_name: fieldName,
				challenge,
				signature,
			}); // send AccessCryptoShareRequest payload to Ankhora cloud backend

			await logAuditEvent({
				event_type: 'decrypt',
				entry_id: entry.id,
				field_name: fieldName,
				timestamp: new Date().toISOString(),
				user_id: 'current_user',
			});

			const timeout = setTimeout(() => {
				handleMaskField(fieldName);
			}, (expires_in || DEFAULT_REVEAL_TIMEOUT) * 1000);

			setRevealedFields(prev => {
				const newMap = new Map(prev);
				newMap.set(fieldName, { name: fieldName, value: plaintext, timeout });
				return newMap;
			});

			toast({
				title: "Field revealed",
				description: `Will auto-mask in ${expires_in || DEFAULT_REVEAL_TIMEOUT}s`,
			});
		} catch (error) {
			toast({
				title: "Decryption failed",
				description: error instanceof Error ? error.message : "Could not decrypt field.",
				variant: "destructive",
			});
		} finally {
			setIsRevealing(null);
			setDecryptingField(null);
		}
	};

	const handleMaskField = (fieldName: string) => {
		const field = revealedFields.get(fieldName);
		if (field) {
			clearTimeout(field.timeout);
			setRevealedFields(prev => {
				const newMap = new Map(prev);
				newMap.delete(fieldName);
				return newMap;
			});
		}
	};

	const handleCopyField = (fieldName: string) => {
		const field = revealedFields.get(fieldName);
		if (field) {
			navigator.clipboard.writeText(field.value);
			toast({
				title: "Copied to clipboard",
				description: "Field value copied securely.",
			});
		}
	};
	const renderSensitiveField = (
		fieldName: string,
		label: string,
		isSensitive: boolean = true,
		entry: any, // original encrypted/metadata entry
	) => {
		const revealed = revealedFields.get(fieldName);
		const isRevealed = !!revealed;
		const isDecrypting = decryptingField === fieldName;
		console.log({entry})

		// Parse decrypted value when present; fall back to original entry shape for type
		let decryptedEntry: LoginEntry | CardEntry | IdentityEntry | NoteEntry | SSHKeyEntry | null = null;
		if (isRevealed && typeof revealed!.value === "string") {
			try {
				decryptedEntry = JSON.parse(revealed!.value);
			} catch {
				decryptedEntry = null;
			}
		}

		const effectiveEntry = (decryptedEntry as any) || entry;

		const renderEntryContent = () => {
			switch (effectiveEntry.type) {
				case "login": {
					const login = effectiveEntry as LoginEntry;
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
						</div>
					);
				}

				case "card": {
					const card = effectiveEntry as CardEntry;
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
					const identity = effectiveEntry as IdentityEntry;
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
					const note = effectiveEntry as NoteEntry;
					const text = (note as any).note ?? (note as any).content;
					return (
						<div className="text-sm">
							{isRevealed ? text || "No content" : "••••••••••••"}
						</div>
					);
				}

				case "sshkey": {
					const ssh = effectiveEntry as SSHKeyEntry;
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
							{isRevealed ? revealed?.value ?? JSON.stringify(entry, null, 2) : "••••••••••••"}
						</div>
					);
			}
		};

		return (
			<div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500">
				<Label
					htmlFor={fieldName}
					className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all"
				>
					{label}
					<Shield className="h-3 w-3 text-primary" />
				</Label>
				<div className="flex gap-2">
					<div className="relative flex-1">
						<div className="h-40 backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl shadow-inner p-4 overflow-auto">
							{renderEntryContent()}
						</div>
						{isDecrypting && (
							<div className="absolute inset-0 pointer-events-none">
								<Sparkles className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary animate-pulse" />
							</div>
						)}
					</div>

					<Button
						size="icon"
						variant="outline"
						onClick={() => (isRevealed ? handleMaskField(fieldName) : handleRevealField(fieldName))}
						disabled={isRevealing === fieldName}
						className={cn(
							"transition-all",
							isRevealed && "border-primary/50 text-primary hover:bg-primary/10",
						)}
					>
						{isRevealed ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
					</Button>
					{isRevealed && (
						<Button
							size="icon"
							variant="outline"
							onClick={() => handleCopyField(fieldName)}
							className="border-primary/50 text-primary hover:bg-primary/10 transition-all"
						>
							<Copy className="h-4 w-4" />
						</Button>
					)}
				</div>
			</div>
		);
	};



	const getRowRecipient = () => {
		return recipients.find((r) => r.email === user?.Email)
	}

	return (
		<div className="flex-1 flex flex-col h-full bg-gradient-to-b from-background to-secondary/20 overflow-hidden">
			{/* Fixed Header */}
			<div className="sticky top-0 z-10 border-b border-border p-6 bg-background">
				<h2 className="text-2xl font-bold">
					{view === "recipients" && "Recipients"}
					{view === "audit" && "Audit Log"}
					{view === "encryption" && "Encryption Policy"}
					{view === "metadata" && "Metadata"}
				</h2>
			</div>

			{/* Scrollable Content */}
			<ScrollArea className="flex-1">
				<div className="p-6">
					{/* Recipients View */}
					{view === "recipients" && (
						<div className="space-y-4">
							{/* Add Recipient Button */}
							{!showAddForm && (
								<Button onClick={() => setShowAddForm(true)} className="w-full">
									<UserPlus className="h-4 w-4 mr-2" />
									Add Recipient
								</Button>
							)}

							{/* Add Recipient Form */}
							{showAddForm && (
								<div className="p-4 border border-border rounded-lg space-y-4 bg-background">
									<div>
										<Label htmlFor="name">Name</Label>
										<Input
											id="name"
											value={newRecipient.name}
											onChange={(e) => setNewRecipient({ ...newRecipient, name: e.target.value })}
											placeholder="Enter recipient name"
										/>
									</div>
									<div>
										<Label htmlFor="email">Email</Label>
										<Input
											id="email"
											type="email"
											value={newRecipient.email}
											onChange={(e) => setNewRecipient({ ...newRecipient, email: e.target.value })}
											placeholder="Enter recipient email"
										/>
									</div>
									<div>
										<Label htmlFor="role">Role</Label>
										<Select
											value={newRecipient.role}
											onValueChange={(value: "viewer" | "editor" | "owner") =>
												setNewRecipient({ ...newRecipient, role: value })
											}
										>
											<SelectTrigger>
												<SelectValue />
											</SelectTrigger>
											<SelectContent>
												<SelectItem value="viewer">Viewer</SelectItem>
												<SelectItem value="editor">Editor</SelectItem>
												<SelectItem value="owner">Owner</SelectItem>
											</SelectContent>
										</Select>
									</div>
									<div className="flex gap-2">
										<Button onClick={handleAddRecipient} className="flex-1">Add</Button>
										<Button onClick={() => setShowAddForm(false)} variant="outline" className="flex-1">
											Cancel
										</Button>
									</div>
								</div>
							)}

							<Separator />

							{/* Recipients List */}
							<div className="space-y-3">
								{recipients.map((recipient) => (
									<div
										key={recipient.id}
										className="p-4 border border-border rounded-lg hover:bg-accent/50 transition-colors"
									>
										<div className="flex items-start justify-between">
											<div className="flex-1">
												<h4 className="font-medium">{recipient.name}</h4>
												<p className="text-sm text-muted-foreground">{recipient.email}</p>
												<p className="text-xs text-muted-foreground mt-1">
													Joined: {new Date(recipient?.joined_at).toLocaleDateString()}
												</p>
											</div>
											<div className="flex items-center gap-2">
												<Select
													value={recipient?.role}
													onValueChange={(value: "viewer" | "editor" | "owner" | "read") =>
														handleChangeRole(recipient.email, value)
													}
												>
													<SelectTrigger className="w-24">
														<SelectValue />
													</SelectTrigger>
													<SelectContent>
														<SelectItem value="viewer">Viewer</SelectItem>
														<SelectItem value="read">Viewer</SelectItem>
														<SelectItem value="editor">Editor</SelectItem>
														<SelectItem value="owner">Owner</SelectItem>
													</SelectContent>
												</Select>
												<Button
													variant="ghost"
													size="icon"
													onClick={() => handleRemoveRecipient(recipient.id)}
												>
													<Trash2 className="h-4 w-4 text-destructive" />
												</Button>
											</div>
										</div>
									</div>
								))}
							</div>

							<p className="text-xs text-muted-foreground text-center mt-6">
								Recipients must have verified sovereign identity.
							</p>
						</div>
					)}

					{/* Audit Log View */}
					{view === "audit" && (
						<div className="space-y-4">
							<Button variant="outline" className="w-full">
								<Download className="h-4 w-4 mr-2" />
								Export Log
							</Button>

							<Separator />

							<div className="space-y-3">
								{entry.access_log.length > 0 && entry.access_log.map((event: AuditEvent) => (
									<div
										key={event.EventID}
										className="p-4 border border-border rounded-lg hover:bg-accent/50 transition-colors"
									>
										<div className="flex items-start justify-between mb-2">
											<Badge variant="outline">{event.EventType}</Badge>
											<span className="text-xs text-muted-foreground">
												{new Date(event.Timestamp).toLocaleString()}
											</span>
										</div>

										{/* Primary line: IP / actor */}
										<p className="text-sm font-medium">
											{event.IPAddress || "Unknown source"}
										</p>

										{/* Details line: human description + Stellar tx if present */}
										<p className="text-sm text-muted-foreground mt-1">
											{event.Event}
											{event.StellarTx && (
												<>
													{" · "}
													<span className="font-mono text-[11px]">
														Tx: {event.StellarTx}
													</span>
												</>
											)}
										</p>

										{/* Optional: ShareID for cross‑reference */}
										{event.ShareID && (
											<p className="text-[11px] text-muted-foreground mt-1">
												Share ID: <span className="font-mono">{event.ShareID}</span>
											</p>
										)}
									</div>
								))}
							</div>
						</div>
					)}


					{/* Encryption View */}
					{view === "encryption" && (
						<div className="space-y-6">
							<div className="p-4 border border-border rounded-lg space-y-3 bg-background">
								<div className="flex items-center gap-2">
									<Shield className="h-5 w-5 text-primary" />
									<h3 className="font-semibold">Encryption Algorithm</h3>
								</div>
								<Badge variant="default" className="text-sm">
									{entry.encryption}
								</Badge>
							</div>

							<div className="p-4 border border-border rounded-lg space-y-3 bg-background">
								<div className="flex items-center gap-2">
									<Database className="h-5 w-5 text-primary" />
									<h3 className="font-semibold">Blockchain Verification</h3>
								</div>
								<div className="space-y-2">
									<div>
										<p className="text-xs text-muted-foreground">Blockchain Hash</p>
										<p className="text-sm font-mono break-all"></p>
									</div>
									<div>
										<p className="text-xs text-muted-foreground">IPFS Anchor</p>
										<p className="text-sm font-mono break-all"></p>
									</div>
								</div>
								<div className="flex gap-2 mt-4">
									<Badge variant="outline" className="text-xs">
										<Shield className="h-3 w-3 mr-1" />
										Stellar Verified
									</Badge>
									<Badge variant="outline" className="text-xs">
										<Database className="h-3 w-3 mr-1" />
										IPFS Audit
									</Badge>
								</div>
							</div>

							<div className="p-4 border border-border rounded-lg space-y-3 bg-background">
								<h3 className="font-semibold">Key Exchange</h3>
								<p className="text-sm text-muted-foreground">
									Last key rotation: {new Date(entry.updated_at).toLocaleString()}
								</p>
								<Button variant="outline" size="sm" className="w-full">
									Force Key Rotation
								</Button>
							</div>
						</div>
					)}

					{/* Metadata View */}
					{view === "metadata" && (
						<div className="space-y-4">
							<div className="grid grid-cols-2 gap-4">
								<div>
									<Label className="text-muted-foreground">Entry Name</Label>
									<p className="font-medium mt-1">{entry.entry_name}</p>
								</div>
								<div>
									<Label className="text-muted-foreground">Type</Label>
									<Badge variant="outline" className="mt-1">{entry.entry_type}</Badge>
								</div>
								<div>
									<Label className="text-muted-foreground">Folder</Label>
									<p className="font-medium mt-1 flex items-center gap-1">
										<FolderOpen className="h-4 w-4" />
									</p>
								</div>
								<div>
									<Label className="text-muted-foreground">Status</Label>
									<Badge variant="outline" className="mt-1">{entry.status}</Badge>
								</div>
								<div>
									<Label className="text-muted-foreground">Created</Label>
									<p className="font-medium mt-1 flex items-center gap-1">
										<Calendar className="h-4 w-4" />
										{new Date(entry.created_at).toLocaleString()}
									</p>
								</div>
								<div>
									<Label className="text-muted-foreground">Last Updated</Label>
									<p className="font-medium mt-1 flex items-center gap-1">
										<Calendar className="h-4 w-4" />
										{new Date(entry.updated_at).toLocaleString()}
									</p>
								</div>
							</div>
							<Separator />
							<div>
								<Label className="text-muted-foreground">Description</Label>
								<p className="mt-2 text-sm"></p>
							</div>
							{/* Show Encrypted payload */}
							<div key={entry.entry_name}>
								<Label className="text-muted-foreground">Encrypted Payload</Label>
								{renderSensitiveField("EncryptedPayload", "Encrypted Payload", true, entry)}
							</div>
						</div>
					)}
				</div>
			</ScrollArea>
		</div>
	);
}
