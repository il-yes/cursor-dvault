import { CheckCircle, Clock, XCircle, UserRound, X } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { PendingIntentShare } from "@/types/sharing";
import * as AppAPI from "../../wailsjs/go/main/App";
import { toast } from "@/hooks/use-toast";
import { useState } from "react";
import { VaultEntry } from "@/types/vault";
import { useVaultStore } from "@/store/vaultStore";
import { createSharedEntry } from "@/services/api";
import { CreateShareEntryPayload, SharedEntry } from "@/types/sharing";
import { buildEntrySnapshot, cn } from "@/lib/utils";
import { listSharedEntries } from "@/services/api";
import { useAuthStore } from "@/store/useAuthStore";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Label } from "./ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select";
import { Textarea } from "./ui/textarea";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { Calendar } from "./ui/calendar";
import { format } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";
import { Input } from "./ui/input";


interface PendingIntentShareViewProps {
    items: PendingIntentShare[];
    onReject: (intent: PendingIntentShare) => void;
    onShareSuccess: () => void;
}


export function PendingIntentShareView({
    items,
    onReject,
    onShareSuccess,
}: PendingIntentShareViewProps) {
    const vault = useVaultStore((state) => state.vault);
    const setSharedEntries = useVaultStore((state) => state.setSharedEntries);
    const addSharedEntry = useVaultStore((state) => state.addSharedEntry);
    const jwtToken = useAuthStore((state) => state.jwtToken);

    const [selectedEntry, setSelectedEntry] = useState("");
    const [recipients, setRecipients] = useState<string[]>([]);
    const [permission, setPermission] = useState<"read" | "edit" | "temporary">("read");
    const [expirationDate, setExpirationDate] = useState<Date | undefined>(undefined);
    const [customMessage, setCustomMessage] = useState("");
    const [allowDownload, setAllowDownload] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);
	const [recipientInput, setRecipientInput] = useState("");

    // Get vault entries from store
    const vaultEntries: VaultEntry[] = vault?.Vault ? [
        ...(vault.Vault.entries?.login || []),
        ...(vault.Vault.entries?.card || []),
        ...(vault.Vault.entries?.note || []),
        ...(vault.Vault.entries?.sshkey || []),
        ...(vault.Vault.entries?.identity || []),
    ] : [];

    const [activeIntent, setActiveIntent] = useState<PendingIntentShare | null>(null);


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

    const onActivate = async () => {
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

                entry_snapshot: buildEntrySnapshot(selectedVaultEntry), // TODO: fix to get the customFields too

                // IMPORTANT: OPTIMISTIC → no recipients yet
                recipients: [],

                expires_at: expirationDate?.toISOString() || null,
                download_allowed: allowDownload,
                attachmentCIDs: selectedVaultEntry?.attachmentCIDs,
            };

            // 2️⃣ Add optimistic entry and get temp ID
            const tempId = addSharedEntry(optimisticPayload);


            let pendingShareIntents: string[] = []
            const getCustomersFromCloud = async (email: string) => {
                try {
                    const response = await AppAPI.CheckUserEmail(jwtToken, email);
                    console.log('getCustomersFromCloud response', response);
                    return response;

                } catch (error) {
                    console.log('User not found with this email: ', email)
                    pendingShareIntents.push(email)
                    // build temp trace_core.User
                    const tempUser = {
                        id: 0,
                        firstName: email.split("@")[0],
                        lastName: "",
                        email: email,
                        public_key: "",
                    }
                    return tempUser
                }
            };

            // Get public key for recipients
            const customers = await Promise.all(recipients.map(email => getCustomersFromCloud(email)));
            console.log({ customers });
            const publicKeys = customers.map(customer => customer?.public_key);
            console.log({ publicKeys });

            // TODO: change no found logic, if not found, let inform that this unknow user will be notified to join Ankhora and accept the share
            if (pendingShareIntents.length > 0) {
                console.log("pendingShareIntents", pendingShareIntents);
                toast({
                    title: "Invitation to Join",
                    description: `${pendingShareIntents.toString()} doesn't have any Ankhora account yet. We'll notify them to join Ankhora and accept the share`,
                });
            }

            // 3️⃣ Create real shared entry via backend
            const cloudResponse = await createSharedEntry({
                entry_id: selectedEntry,
                recipients: customers.map((customer) => ({
                    id: customer?.id.toString(),
                    name: customer?.first_name + " " + customer?.last_name,
                    email: customer?.email,
                    publicKey: customer?.public_key,
                    role: permission == "edit" ? "editor" : "viewer",
                })),
                permission: permission as "read" | "edit" | "temporary",
                expires_at: expirationDate?.toISOString(),
                custom_message: customMessage,
                download_allowed: allowDownload,
                attachmentCIDs: selectedVaultEntry?.attachmentCIDs,
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
            setPermission("read");
            setExpirationDate(undefined);
            setCustomMessage("");

            // delete pending intent share.

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
        <div className="flex h-full flex-col">

            {/* Header */}
            <div className="border-b p-5">
                <h2 className="text-xl font-semibold">
                    Pending Share Invitations
                </h2>

                <p className="text-sm text-muted-foreground mt-1">
                    Review accepted invitations and activate cryptographic shares.
                </p>
            </div>


            {/* List */}
            <div className="flex-1 overflow-y-auto p-4 space-y-3">

                {items.length === 0 && (
                    <div className="flex flex-col items-center justify-center h-full text-muted-foreground">
                        <Clock className="h-10 w-10 mb-3 opacity-40" />
                        <p>No pending share invitations to activate</p>
                    </div>
                )}

                {items.map(intent => (
                    <div key={intent.id} className="rounded-xl border bg-secondary/30 p-4 space-y-4">
                        <div className="flex items-start justify-between">
                            <div className="flex gap-3">
                                <div className="p-2 rounded-lg bg-primary/10 text-primary">
                                    <UserRound size={18} />
                                </div>
                                <div>
                                    <h3 className="font-semibold">{intent.recipient_email}</h3>
                                    <p className="text-sm text-muted-foreground">Requested access to your encrypted share</p>
                                </div>
                            </div>

                            <Badge variant={intent.status === "accepted" ? "default" : "secondary"}>{intent.status}</Badge>
                        </div>

                        <div className="text-sm space-y-1">
                            <div>
                                Share ID:
                                <span className="ml-2 font-mono">
                                    {intent.share_id}
                                </span>
                            </div>

                            <div>
                                Created:
                                <span className="ml-2">
                                    {new Date(
                                        intent.created_at
                                    ).toLocaleDateString()}
                                </span>
                            </div>
                        </div>

                        <div className="flex justify-end gap-2">
                            <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                onClick={() => onReject(intent)}
                            >
                                <XCircle className="mr-2 h-4 w-4" />
                                Reject
                            </Button>

                            {/* Open dialog here: Intent form */}
                            <Dialog>
                                <DialogTrigger asChild>
                                    <Button
                                        size="sm"
                                        onClick={() => {
                                            setRecipientInput(intent.recipient_email);
                                        }}
                                    >
                                        <CheckCircle className="mr-2 h-4 w-4" />
                                        Activate Share
                                    </Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Activate Share</DialogTitle>
                                        <DialogDescription>
                                            Activate this share to make it active
                                        </DialogDescription>
                                    </DialogHeader>
                                    {/* TODO: Reuse the ShareViewForm but enable edit mode */}
                                    <form
                                        onSubmit={async (e) => {
                                            e.preventDefault();
                                            await onActivate();
                                        }}>
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
                                            <RadioGroup value={permission} onValueChange={(value) => setPermission(value as "read" | "edit" | "temporary")}>
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
                                        {/* role */}
                                        <Label htmlFor="role">Role</Label>
                                        <Select
                                            value={permission}
                                            onValueChange={(value) => setPermission(value as "read" | "edit")}
                                        >
                                            <SelectTrigger>
                                                <SelectValue placeholder="Select role" />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="read">Read Only</SelectItem>
                                                <SelectItem value="edit">Edit</SelectItem>
                                            </SelectContent>
                                        </Select>

                                        {/* downloadAllowed */}
                                        <Label htmlFor="downloadAllowed">Allow Download</Label>

                                        <RadioGroup
                                            value={allowDownload ? "true" : "false"}
                                            onValueChange={(value: string) => setAllowDownload(value === "true")}
                                        >
                                            <div className="flex items-center space-x-2">
                                                <RadioGroupItem value="true" id="true" />
                                                <Label htmlFor="true" className="font-normal cursor-pointer">
                                                    Yes
                                                </Label>
                                            </div>
                                            <div className="flex items-center space-x-2">
                                                <RadioGroupItem value="false" id="false" />
                                                <Label htmlFor="false" className="font-normal cursor-pointer">
                                                    No
                                                </Label>
                                            </div>
                                        </RadioGroup>


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


                                        <DialogFooter>
                                            <Button type="submit">Activate Share</Button>
                                        </DialogFooter>
                                    </form>




                                </DialogContent>
                            </Dialog>
                        </div>
                    </div>
                ))}

            </div>

        </div>
    );
}