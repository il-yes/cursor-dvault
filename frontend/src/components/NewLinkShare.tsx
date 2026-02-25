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
import { CreateLinkShareEntryPayload, CreateShareEntryPayload, LinkShareEntry, SharedEntry } from "@/types/sharing";
import { VaultEntry } from "@/types/vault";
import { getSharedEntry, listSharedEntries, createLinkShareEntry, listLinkSharesByMe } from "@/services/api";
import { Switch } from "@radix-ui/react-switch";
import { useAuthStore } from "@/store/useAuthStore";
import { useAppStore } from "@/store/appStore";
import { Shield, ArrowLeft, Eye, EyeOff } from "lucide-react";

interface NewShareModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onShareSuccess?: () => void;
}

export function NewLinkShareModal({ open, onOpenChange, onShareSuccess }: NewShareModalProps) {
    const { toast } = useToast();
    const vault = useVaultStore((state) => state.vault);
    const setLinkSharedByMe = useVaultStore((state) => state.setLinkSharedEntries);
    const addLinkSharedEntry = useVaultStore((state) => state.addLinkSharedEntry);
    const [selectedEntry, setSelectedEntry] = useState("");
    const [permission, setPermission] = useState("read");
    const [password, setPassword] = useState("");
    const [expirationDate, setExpirationDate] = useState<Date>();
    const [customMessage, setCustomMessage] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [maxViews, setMaxViews] = useState(1);
    const [allowDownload, setAllowDownload] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
    const jwtToken = useAuthStore.getState().jwtToken;
    const user = useAuthStore.getState().user

    const handleShare = async () => {
        if (!selectedEntry) {
            toast({
                title: "Missing Information",
                description: "Please select an entry.",
                variant: "destructive",
            });
            return;
        }

        setIsSubmitting(true);
        try {
            const selectedVaultEntry = vaultEntries.find(e => e.id === selectedEntry);
            if (!selectedVaultEntry) throw new Error("Entry not found");

            // 1️⃣ Build optimistic payload
            const optimisticPayload: LinkShareEntry = {
                id: "",
                entry_name: selectedVaultEntry.entry_name,
                status: "pending",
                expiry: expirationDate?.toISOString() || null,
                uses_left: maxViews,
                link: "",
                audit_log: [],
                payload: JSON.stringify(selectedVaultEntry),
                allow_download: allowDownload,
                password: password
            };

            // 2️⃣ Zustaand - Add optimistic entry and get temp ID
            const tempLinkShareEntry = addLinkSharedEntry(optimisticPayload);
            console.log("tempLinkShareEntry", optimisticPayload);


            // 3️⃣ Create real shared entry via backend
            const cloudResponse = await createLinkShareEntry({
                payload: selectedVaultEntry,
                expires_at: expirationDate?.toISOString() || null,
                max_views: maxViews,
                creator_user_id: user?.id,
                creator_email: user?.email,
                entry_type: selectedVaultEntry.type,
                title: selectedVaultEntry.entry_name,
                download_allowed: allowDownload,
                password: password
            });

            console.log("☁️ Cloud shared entry:", cloudResponse);

            // 4️⃣ Replace optimistic entry with real backend entry
            const fullEntries = await listLinkSharesByMe();
            setLinkSharedByMe(fullEntries);
            toast({
                title: "✅ Entry shared successfully",
                description: "Now visible in your Shared Entries",
            });

            onShareSuccess?.();

            // Reset UI
            setSelectedEntry("");
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

    // Get vault entries from store
    const vaultEntries: VaultEntry[] = vault?.Vault ? [
        ...(vault.Vault.entries?.login || []),
        ...(vault.Vault.entries?.card || []),
        ...(vault.Vault.entries?.note || []),
        ...(vault.Vault.entries?.sshkey || []),
        ...(vault.Vault.entries?.identity || []),
    ] : [];

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-xl">
                <DialogHeader>
                    <DialogTitle>Link Share Entry</DialogTitle>
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

                    {/* Max Views */}
                    <div className="space-y-2">
                        <Label htmlFor="maxViews">Max Views *</Label>
                        <Input
                            id="maxViews"
                            type="number"
                            value={maxViews}
                            onChange={(e) => setMaxViews(Number(e.target.value))}
                            className="w-full"
                        />
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


                    {/* Allow Download */}
                    <div className="space-y-3">
                        <Label>AllowDownload</Label>
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
                    </div>

                    {/* Password */}
                    <div className="space-y-2 rerlative">
                        <Label htmlFor="password">Password (Optional)</Label>
                        <Input
                            type={showPassword ? "text" : "password"}
                            id="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg transition-colors"
                        onClick={() => setShowPassword(!showPassword)}
                        >
                        {showPassword ? (
                            <EyeOff className="h-4 w-4 text-muted-foreground" />
                        ) : (
                            <Eye className="h-4 w-4 text-muted-foreground" />
                        )}
                        </Button>
                    </div>  

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
                        className="bg-[#C9A44A] hover:bg-[#B8934A]"
                    >
                        {isSubmitting ? "Sharing..." : "Share Now"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}