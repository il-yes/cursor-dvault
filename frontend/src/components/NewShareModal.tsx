import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { X, Calendar as CalendarIcon } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { format, set } from "date-fns";
import { buildEntrySnapshot, cn } from "@/lib/utils";
import { useVaultStore } from "@/store/vaultStore";
import { createSharedEntry } from "@/services/api";
import { CreateShareEntryPayload, SharedEntry } from "@/types/sharing";
import { VaultEntry } from "@/types/vault";
import { getSharedEntry, listSharedEntries } from "@/services/api";

interface NewShareModalProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
	onShareSuccess?: () => void;
}

export function NewShareModal({ open, onOpenChange, onShareSuccess }: NewShareModalProps) {
	const { toast } = useToast();
	const vault = useVaultStore((state) => state.vault);
	const setSharedEntries = useVaultStore((state) => state.setSharedEntries);	
	const addSharedEntry = useVaultStore((state) => state.addSharedEntry);
	const updateSharedEntry = useVaultStore((state) => state.updateSharedEntry);
	const removeSharedEntry = useVaultStore((state) => state.removeSharedEntry);
	const [selectedEntry, setSelectedEntry] = useState("");
	const [recipients, setRecipients] = useState<string[]>([]);
	const [recipientInput, setRecipientInput] = useState("");
	const [permission, setPermission] = useState("read");
	const [expirationDate, setExpirationDate] = useState<Date>();
	const [customMessage, setCustomMessage] = useState("");
	const [isSubmitting, setIsSubmitting] = useState(false);

	// Get vault entries from store
	const vaultEntries: VaultEntry[] = vault?.Vault ? [
		...(vault.Vault.entries?.login || []),
		...(vault.Vault.entries?.card || []),
		...(vault.Vault.entries?.note || []),
		...(vault.Vault.entries?.sshkey || []),
		...(vault.Vault.entries?.identity || []),
	] : [];

	const handleAddRecipient = () => {
		const trimmed = recipientInput.trim();
		console.log('recipientInput', recipientInput);
		if (trimmed && !recipients.includes(trimmed)) {
			setRecipients([...recipients, trimmed]);
			setRecipientInput("");
		}
	};

	const handleRemoveRecipient = (recipient: string) => {
		setRecipients(recipients.filter(r => r !== recipient));
	};

	const handleShare2 = async () => {
		if (!selectedEntry || recipients.length === 0) {
			toast({
				title: "Missing Information",
				description: "Please select an entry and add at least one recipient.",
				variant: "destructive",
			});
			return;
		}

		setIsSubmitting(true);
		try {
			// Call API to create shared entry
			const response = await createSharedEntry({
				entry_id: selectedEntry,
				recipients: recipients.map(email => ({ name: email.split("@")[0], email, role: permission })),
				permission: permission as 'read' | 'edit' | 'temporary',
				expiration_date: expirationDate?.toISOString(),
				custom_message: customMessage,
			});
			console.log('✅ Shared entry created: --------------------', response);


			// we have to fetch the entry from the vault and fill the fields 

			// Add to local store
			const selectedVaultEntry = vaultEntries.find(e => e.id === selectedEntry);

			const newSharedEntry: CreateShareEntryPayload = {
				entry_name: selectedVaultEntry?.entry_name || "Shared Entry",
				entry_type: selectedVaultEntry?.type || "note",
				status: "active",
				access_mode: permission === "edit" ? "edit" : "read",
				encryption: "AES-256-GCM",

				entry_snapshot: buildEntrySnapshot(selectedVaultEntry!),

				recipients: recipients.map(email => ({
					name: email.split("@")[0],
					email,
					public_key: "",
					role: permission === "edit" ? "editor" : "viewer",
				})),

				expires_at: expirationDate?.toISOString() || null,
			};

			addSharedEntry(newSharedEntry);

			toast({
				title: "✅ Entry shared successfully",
				description: "Now visible in your Shared Entries",
			});

			// Notify parent to refresh the list
			onShareSuccess?.();

			// Reset and close
			setSelectedEntry("");
			setRecipients([]);
			setRecipientInput("");
			setPermission("read");
			setExpirationDate(undefined);
			setCustomMessage("");
			onOpenChange(false);
		} catch (error) {
			console.error('Failed to share entry:', error);
			toast({
				title: "Failed to share",
				description: "Could not share entry. Please try again.",
				variant: "destructive",
			});
		} finally {
			setIsSubmitting(false);
		}
	};

	const handleShare = async () => {
		if (!selectedEntry || recipients.length === 0) {
			toast({
				title: "Missing Information",
				description: "Please select an entry and add at least one recipient.",
				variant: "destructive",
			});
			return;
		}

		setIsSubmitting(true);
		try {
			const selectedVaultEntry = vaultEntries.find(e => e.id === selectedEntry);
			if (!selectedVaultEntry) throw new Error("Entry not found");

			// 1️⃣ Build optimistic payload
			const optimisticPayload: CreateShareEntryPayload = {
				entry_name: selectedVaultEntry.entry_name,
				entry_type: selectedVaultEntry.type,

				// always pending at creation
				status: "pending",

				access_mode: permission === "edit" ? "edit" : "read",
				encryption: "AES-256-GCM",

				entry_snapshot: buildEntrySnapshot(selectedVaultEntry),

				// IMPORTANT: OPTIMISTIC → no recipients yet
				recipients: [],

				expires_at: expirationDate?.toISOString() || null,
			};

			// 2️⃣ Add optimistic entry and get temp ID
			const tempId = addSharedEntry(optimisticPayload);

			
			const getPublicKey = async (email: string) => {
				const response = await fetch(`http://localhost:4001/api/check-email?email=${email}`, {
					method: "GET",
					headers: {
						"Content-Type": "application/json",
					},
				});
				const data = await response.json();
				console.log({data});
				return data.public_key;
			};
			// Get public key for recipients
			const publicKeys = await Promise.all(recipients.map(email => getPublicKey(email)));
			console.log({publicKeys});
			console.log(publicKeys[0]);
			
			if (!publicKeys) {
				toast({
					title: "Error",
					description: "No public key found for this email",
				});
				return;
			}

			// 3️⃣ Create real shared entry via backend
			const cloudResponse = await createSharedEntry({
				entry_id: selectedEntry,
				recipients: recipients.map((email, index) => ({
					name: email.split("@")[0],
					email,
					publicKey: publicKeys[index],
					role: permission,
				})),
				permission: permission as "read" | "edit" | "temporary",
				expiration_date: expirationDate?.toISOString(),
				custom_message: customMessage,
			});

			console.log("☁️ Cloud shared entry:", cloudResponse);

			// 4️⃣ Replace optimistic entry with real backend entry
			const fullEntries = await listSharedEntries();
			// updateSharedEntry(tempId, fullEntry);
			setSharedEntries(fullEntries);

			toast({
				title: "✅ Entry shared successfully",
				description: "Now visible in your Shared Entries",
			});

			onShareSuccess?.();

			// Reset UI
			setSelectedEntry("");
			setRecipients([]);
			setRecipientInput("");
			setPermission("read");
			setExpirationDate(undefined);
			setCustomMessage("");
			onOpenChange(false);

		} catch (err) {
			console.error("Failed to share entry:", err);
			toast({
				title: "Failed to share",
				description: "Could not share entry. Please try again.",
				variant: "destructive",
			});
		} finally {
			setIsSubmitting(false);
		}
	};


	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="sm:max-w-xl">
				<DialogHeader>
					<DialogTitle>Share Entry</DialogTitle>
				</DialogHeader>

				<div className="space-y-6 py-4">
					{/* Select Entry */}
					<div className="space-y-2">
						<Label htmlFor="entry">Select Entry *</Label>
						<Select value={selectedEntry} onValueChange={setSelectedEntry}>
							<SelectTrigger id="entry">
								<SelectValue placeholder="Choose an entry from your vault" />
							</SelectTrigger>
							<SelectContent>
								{vaultEntries.map((entry) => (
									<SelectItem key={entry.id} value={entry.id}>
										<span className="capitalize">{entry.type}</span> • {entry.entry_name}
									</SelectItem>
								))}
							</SelectContent>
						</Select>
					</div>

					{/* Recipients */}
					<div className="space-y-2">
						<Label htmlFor="recipients">Recipients *</Label>
						<div className="flex gap-2">
							<Input
								id="recipients"
								placeholder="Enter email or username"
								value={recipientInput}
								onChange={(e) => setRecipientInput(e.target.value)}
								onKeyDown={(e) => {
									if (e.key === "Enter") {
										e.preventDefault();
										handleAddRecipient();
									}
								}}
							/>
							<Button type="button" onClick={handleAddRecipient} variant="secondary">
								Add
							</Button>
						</div>
						{recipients.length > 0 && (
							<div className="flex flex-wrap gap-2 mt-2">
								{recipients.map((recipient) => (
									<Badge key={recipient} variant="secondary" className="pl-3 pr-1 py-1">
										{recipient}
										<button
											onClick={() => handleRemoveRecipient(recipient)}
											className="ml-2 hover:text-destructive"
										>
											<X className="h-3 w-3" />
										</button>
									</Badge>
								))}
							</div>
						)}
					</div>

					{/* Permissions */}
					<div className="space-y-3">
						<Label>Permissions</Label>
						<RadioGroup value={permission} onValueChange={setPermission}>
							<div className="flex items-center space-x-2">
								<RadioGroupItem value="read" id="read" />
								<Label htmlFor="read" className="font-normal cursor-pointer">
									Read-only
								</Label>
							</div>
							<div className="flex items-center space-x-2">
								<RadioGroupItem value="edit" id="edit" />
								<Label htmlFor="edit" className="font-normal cursor-pointer">
									Read & Edit
								</Label>
							</div>
							<div className="flex items-center space-x-2">
								<RadioGroupItem value="temporary" id="temporary" />
								<Label htmlFor="temporary" className="font-normal cursor-pointer">
									Temporary Access
								</Label>
							</div>
						</RadioGroup>
					</div>

					{/* Expiration Date */}
					{permission === "temporary" && (
						<div className="space-y-2">
							<Label>Expiration Date</Label>
							<Popover>
								<PopoverTrigger asChild>
									<Button
										variant="outline"
										className={cn(
											"w-full justify-start text-left font-normal",
											!expirationDate && "text-muted-foreground"
										)}
									>
										<CalendarIcon className="mr-2 h-4 w-4" />
										{expirationDate ? format(expirationDate, "PPP") : "Pick a date"}
									</Button>
								</PopoverTrigger>
								<PopoverContent className="w-auto p-0">
									<Calendar
										mode="single"
										selected={expirationDate}
										onSelect={setExpirationDate}
										disabled={(date) => date < new Date()}
										initialFocus
									/>
								</PopoverContent>
							</Popover>
						</div>
					)}

					{/* Custom Message */}
					<div className="space-y-2">
						<Label htmlFor="message">Custom Message (Optional)</Label>
						<Textarea
							id="message"
							placeholder="Add a note for recipients..."
							value={customMessage}
							onChange={(e) => setCustomMessage(e.target.value)}
							rows={3}
						/>
					</div>
				</div>

				<DialogFooter>
					<Button variant="outline" onClick={() => onOpenChange(false)}>
						Cancel
					</Button>
					<Button
						onClick={handleShare}
						disabled={!selectedEntry || recipients.length === 0 || isSubmitting}
						className="bg-[#C9A44A] hover:bg-[#B8934A]"
					>
						{isSubmitting ? "Sharing..." : "Share Now"}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
