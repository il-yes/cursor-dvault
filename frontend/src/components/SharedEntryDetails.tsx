import { SharedEntry, DetailView, Recipient } from "@/types/sharing";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { UserPlus, Trash2, Download, Shield, Database, Calendar, FolderOpen } from "lucide-react";
import { useEffect, useState } from "react";
import { useToast } from "@/hooks/use-toast";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import { useVaultStore } from "@/store/vaultStore";

interface SharedEntryDetailsProps {
	entry: SharedEntry | null;
	view: DetailView;
}

export function SharedEntryDetails({ entry, view }: SharedEntryDetailsProps) {
	const { toast } = useToast();
	const [showAddForm, setShowAddForm] = useState(false);
	const [newRecipient, setNewRecipient] = useState<{ name: string; email: string; role: "viewer" | "editor" | "owner" }>({
		name: "",
		email: "",
		role: "viewer"
	});
	const updateRecipients = useVaultStore((state) => state.updateSharedEntryRecipients);	


	const [recipients, setRecipients] = useState<Recipient[]>([]);

	useEffect(() => {
		setRecipients(entry?.recipients || []);
	}, [entry?.id]);



	const handleAddRecipient = () => {
		const recipient: Recipient = {
			id: `rec-${Date.now()}`,
			share_id: entry.id,
			...newRecipient,
			joined_at: new Date().toISOString(),
		};

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
					<img src={ankhoraLogo} alt="Ankhora Logo" className=" w-auto" style={{ width: "11.5rem" }} />
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
		setRecipients(recipients.filter(r => r.id !== id));

		toast({
			title: "Recipient removed",
			description: `${recipient?.name} has been removed`,
		});
	};

	const handleChangeRole = (id: string, newRole: "viewer" | "read" | "editor" | "owner") => {
		setRecipients(recipients.map(r => r.id === id ? { ...r, role: newRole } : r));

		toast({
			title: "Role updated",
			description: "Recipient role has been changed",
		});
	};

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
												<SelectItem  value="viewer">Viewer</SelectItem>
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
													value={recipient.role}
													onValueChange={(value: "viewer" | "editor" | "owner" | "read") =>
														handleChangeRole(recipient.id, value)
													}
												>
													<SelectTrigger className="w-24">
														<SelectValue />
													</SelectTrigger>
													<SelectContent>
														<SelectItem  value="viewer">Viewer</SelectItem>
														<SelectItem  value="read">Viewer</SelectItem>
														<SelectItem  value="editor">Editor</SelectItem>
														<SelectItem  value="owner">Owner</SelectItem>
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
								{entry.audit_log.map((event) => (
									<div
										key={event.id}
										className="p-4 border border-border rounded-lg hover:bg-accent/50 transition-colors"
									>
										<div className="flex items-start justify-between mb-2">
											<Badge variant="outline">{event.action}</Badge>
											<span className="text-xs text-muted-foreground">
												{new Date(event.timestamp).toLocaleString()}
											</span>
										</div>
										<p className="text-sm font-medium">{event.actor}</p>
										{event.details && (
											<p className="text-sm text-muted-foreground mt-1">{event.details}</p>
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
										<p className="text-sm font-mono break-all">{entry?.blockchain_hash}</p>
									</div>
									{entry?.ipfs_anchor && (
										<div>
											<p className="text-xs text-muted-foreground">IPFS Anchor</p>
											<p className="text-sm font-mono break-all">{entry?.ipfs_anchor}</p>
										</div>
									)}
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
								{entry?.folder && (
									<div>
										<Label className="text-muted-foreground">Folder</Label>
										<p className="font-medium mt-1 flex items-center gap-1">
											<FolderOpen className="h-4 w-4" />
											{entry?.folder}
										</p>
									</div>
								)}
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

							{entry?.description && (
								<>
									<Separator />
									<div>
										<Label className="text-muted-foreground">Description</Label>
										<p className="mt-2 text-sm">{entry?.description}</p>
									</div>
								</>
							)}
						</div>
					)}
				</div>
			</ScrollArea>
		</div>
	);
}
