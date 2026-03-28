import { useState, useEffect, useRef, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Eye, EyeOff, Copy, Shield, Edit, Share2, Trash2, Sparkles, Loader2, Maximize2, Cloud, HardDrive } from "lucide-react";
import { Attachment, Folder, SettingsState, VaultEntry } from "@/types/vault";
import { decryptField, loadAttachment, logAuditEvent } from "@/services/api";
import { toast } from "@/hooks/use-toast";
import { cn, formatFileSize } from "@/lib/utils";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import { Clock } from "lucide-react";
import "./contributionGraph/g-scrollbar.css";
import { Textarea } from "./ui/textarea";
import { FileUploadWidget } from "./FileUploadWidget";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useVaultStore } from "@/store/vaultStore";
import * as AppAPI from "../../wailsjs/go/main/App";
import { Keypair } from "stellar-sdk";
import { Buffer } from 'buffer';
import { useAuthStore } from "@/store/useAuthStore";
import { useVault } from "@/hooks/useVault";
import { withAuth } from "@/hooks/withAuth";
import { ToastAction } from "@radix-ui/react-toast";


interface EntryDetailPanelProps {
    entry: VaultEntry | null;
    editMode?: boolean;
    onEdit?: () => void;
    onSave?: (updates: Partial<VaultEntry>) => void;
    onCancel?: () => void;
    onShare?: (entry: VaultEntry) => void;
    onDelete?: () => void;
}

interface RevealedField {
    name: string;
    value: string;
    timeout: NodeJS.Timeout;
}

const DEFAULT_REVEAL_TIMEOUT = 15;

export function EntryDetailPanel({ entry, editMode, onEdit, onSave, onCancel, onShare, onDelete }: EntryDetailPanelProps) {
    const [revealedFields, setRevealedFields] = useState<Map<string, RevealedField>>(new Map());
    const [isRevealing, setIsRevealing] = useState<string | null>(null);
    const [decryptingField, setDecryptingField] = useState<string | null>(null);
    const [editData, setEditData] = useState<Partial<VaultEntry>>({});
    const [folderId, setFolderId] = useState<string | null>(null);
    const [folders, setFolders] = useState<Folder[]>([]);
    // keep localEntry typed explicitly
    const [localEntry, setLocalEntry] = useState<VaultEntry | null>(entry ?? null);
    const [attachedFiles, setAttachedFiles] = useState<File[]>([])

    const vaultContext = useVaultStore((state) => state.vault);
    const stellar = vaultContext?.vault_runtime_context?.UserConfig?.stellar_account
    const { jwtToken } = useAuthStore.getState();
    const { user } = useAuthStore();

    const [progressVisible, setProgressVisible] = useState(false);
    const [showModal, setShowModal] = useState(false);

    const { encryptFile, uploadToIPFS, createStellarCommit, syncVault, refreshVault } = useVault();

    const [progress, setProgress] = useState(0);
    const [stage, setStage] = useState('encrypting'); // encrypting | uploading | complete


    const defaultSettings: SettingsState = {
        security: {
            autoLockSeconds: 300,
            clearClipboardAfter: 60,
            twoFactorEnabled: false
        },
        sync: {
            stellarFrequency: "manual",
            ipfsPinning: false,
            syncIntervalSeconds: 60,
            maxRetries: 3
        },
        ui: {
            theme: "system",
            animationsEnabled: false
        },
        features: {
            tracecoreEnabled: false,
            cloudBackupEnabled: false,
            threatDetectionEnabled: false,
            browserExtensionEnabled: false,
            gitCLIEnabled: false
        },
        backup: {
            enabled: false,
            schedule: "daily",
            retentionDays: 30,
            encryption: false
        },
        device: {
            user_id: user.id,
            vault_name: "",
            device_id: "0",
            device_name: "",
            last_synced: 0,
        },
        subscription: {
            plan: "free",
            features: {
                tracecoreEnabled: false,
                cloudBackupEnabled: false,
                threatDetectionEnabled: false,
                browserExtensionEnabled: false,
                gitCLIEnabled: false

            },
            limits: {
                maxVaults: 1,
                maxUsers: 1,
                maxDevices: 1,
                maxShares: 1
            }
        },
        sharing: {
            allowExternalSharing: false,
            defaultExpiryHours: 60,
            requirePassword: false,
            maxSharesPerEntry: 3,
        },

        privacy: {
            telemetryEnabled: false,
            anonymousMode: false,
        },
    };
    const [settings, setSettings] = useState<SettingsState>(defaultSettings);
    const { vault } = useVaultStore();
    const vaultPassword = "vaultPassword";

    const [selectedAttachment, setSelectedAttachment] = useState<ExtendedAttachment | null>(null);

    const [transferring, setTransferring] = useState<Record<string, TransferStatus>>({});



    useEffect(() => {
        if (!vault?.Vault?.name) return
        fetchConfig(vault.Vault.name, jwtToken)
    }, [vault])

    const fetchConfig = async (vaultName: string, jwtToken: string) => {
        try {
            const response = await withAuth((token) => {
                return AppAPI.GetConfig(vaultName, token)
            });

            console.log("fetchConfig response", response)

            setSettings(response as unknown as SettingsState)

        } catch (err) {
            console.error("fetchConfig failed", err)
        }
    }

    // just for dev checks
    const [attachments, setAttachments] = useState<Attachment[]>([])

    useEffect(() => {
        vaultContext && setFolders(vaultContext.Vault.folders || []);
    }, [vaultContext]);

    // Sync localEntry with prop changes
    useEffect(() => {
        setLocalEntry(entry ?? null);
    }, [entry]);

    useEffect(() => {
        if (entry) {
            setFolderId(entry.folder_id);
        }
    }, [entry]);

    // Cleanup timeouts on unmount / when revealedFields changes
    useEffect(() => {
        return () => {
            revealedFields.forEach(field => clearTimeout(field.timeout));
        };
    }, [revealedFields]);

    // Clear revealed fields when entry changes
    useEffect(() => {
        // Clear all timeouts and reset revealed fields when switching entries
        revealedFields.forEach(field => clearTimeout(field.timeout));
        setRevealedFields(new Map());
    }, [entry?.id]);

    // Initialize editData when switching to editMode
    useEffect(() => {
        if (entry && editMode) {
            const initialData: Partial<VaultEntry> = { entry_name: entry.entry_name };

            // Copy all fields except immutable ones
            Object.keys(entry).forEach((key) => {
                if (!["id", "created_at", "updated_at"].includes(key)) {
                    // keep type-safety quiet here: entry has many optional fields depending on type
                    // we intentionally copy everything that might be editable
                    // @ts-ignore
                    initialData[key] = (entry as any)[key];
                }
            });

            setEditData(initialData);
        }

        // if leaving editMode, clear local editData (but don't mutate outside)
        if (!editMode) {
            setEditData({});
        }
    }, [entry, editMode]);


    const fetchAttachment = async (hash: string) => {
        const attachmentHash = await loadAttachment(jwtToken, vaultContext.Vault.name, hash);
        return attachmentHash;
    };

    useEffect(() => {
        if (entry) {
            console.log(vaultContext.Vault)
            const attachements = entry.attachments || [];
            setAttachments(attachements)
        }
    }, [entry?.id]);
    useEffect(() => {
        return () => {
            // Clear timeout on unmount
            if (currentTimeoutRef.current) {
                clearTimeout(currentTimeoutRef.current);
            }
        };
    }, []);


    const handleFieldChange = (fieldName: string, value: any) => {
        setEditData(prev => ({ ...prev, [fieldName]: value }));
    };

    const handleSaveEdit = () => {
        if (onSave) {
            // pass only the edited changes (editData)
            if (folderId != null) {
                // find folder by id
                const f = folders.find(f => f.id === folderId);
                if (f) {
                    editData.folder_id = f.id;
                }
            }
            console.log({ editData });
            onSave(editData);
            setEditData({});
        }
    };

    const currentTimeoutRef = useRef<NodeJS.Timeout | null>(null);

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

    const handleRevealField = (fieldName: string) => {
        if (!entry) return;

        // Clear existing timeout
        if (currentTimeoutRef.current) {
            clearTimeout(currentTimeoutRef.current);
            currentTimeoutRef.current = null;
        }

        const clearClipboardAfter = settings?.Vaults?.security?.ClearClipboardAfter;

        // Create timeout
        const timeout = setTimeout(() => {
            handleMaskField(fieldName);
            currentTimeoutRef.current = null;
        }, clearClipboardAfter * 1000);

        currentTimeoutRef.current = timeout;

        // Use entry data directly (no decryption)
        setRevealedFields(prev => {
            const newMap = new Map(prev);
            newMap.set(fieldName, {
                name: fieldName,
                value: entry[fieldName] as string, // ← Direct from entry
                timeout
            });
            return newMap;
        });

        toast({
            title: "Field revealed",
            description: `Will auto-mask in ${clearClipboardAfter}s`,
        });
    };

    const handleMaskField = useCallback((fieldName: string) => {
        setRevealedFields(prev => {
            const field = prev.get(fieldName);
            if (!field) return prev;

            clearTimeout(field.timeout);

            // ← CRITICAL: Create COMPLETELY NEW Map every time
            const newMap = new Map(prev);
            newMap.delete(fieldName);
            return newMap; // ← React sees new reference → re-renders
        });
    }, []);

    /**
     * Renders a sensitive or non-sensitive field.
     * Behavior:
     *  - In view mode: shows masked value and reveal/copy controls (uses revealedFields for plaintext)
     *  - In edit mode: shows editable Input bound to editData; for sensitive fields the input type will be password unless field is revealed
     *
     * Note: Because encrypted fields are not always available in plaintext until decrypted,
     * edit mode will allow user to overwrite the field even if not previously revealed.
     */
    const renderSensitiveField = (fieldName: string, label: string, isSensitive: boolean = true) => {
        const revealed = revealedFields.get(fieldName);
        const isRevealed = !!revealed;
        const isDecrypting = decryptingField === fieldName;

        // Value to display when in view-mode
        const viewPlaintext = isRevealed ? revealed!.value : undefined;

        // Value to display when in edit-mode (prefer editData if present)
        const editValue = (editData as any)[fieldName] ?? "";

        // In non-sensitive fields we still support reveal (because original code had that),
        // but in edit mode non-sensitive fields are simple editable text.
        if (!isSensitive) {
            if (!editMode) {
                return (
                    <div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500">
                        <Label htmlFor={fieldName} className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                            {label}
                        </Label>
                        <div className="flex gap-2">
                            <Input
                                id={fieldName}
                                type="text"
                                // ensure controlled string value (avoid undefined)
                                value={(viewPlaintext ?? (localEntry ? ((localEntry as any)[fieldName] ?? "") : "")) as string}
                                readOnly
                                className="transition-all duration-300 border-border/50"
                            />
                            <Button
                                size="icon"
                                variant="outline"
                                onClick={() => isRevealed ? handleMaskField(fieldName) : handleRevealField(fieldName)}
                                disabled={isRevealing === fieldName}
                                className="transition-all"
                            >
                                <Eye className="h-4 w-4" />
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
            }

            // edit mode - plain editable input
            return (
                <div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500">
                    <Label htmlFor={fieldName} className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                        {label}
                    </Label>
                    <Input
                        id={fieldName}
                        type="text"
                        value={editValue as string}
                        onChange={(e) => handleFieldChange(fieldName, e.target.value)}
                        readOnly={!editMode}
                        className="border-primary/50"
                    />
                </div>
            );
        }

        // SENSITIVE FIELD rendering
        if (!editMode) {
            // view mode for sensitive field
            return (
                <div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500">
                    <Label htmlFor={fieldName} className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                        {label}
                        <Shield className="h-3 w-3 text-primary" />
                    </Label>
                    <div className="flex gap-2">
                        <div className="relative flex-1">
                            <Input
                                id={fieldName}
                                type={isRevealed ? "text" : "password"}
                                value={isRevealed ? revealed!.value : "••••••••••••"}
                                readOnly
                                className={cn(
                                    'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2ll shadow-inner',
                                    editMode && 'border-primary/50 shadow-primary/20',
                                )}
                            />
                            {isDecrypting && (
                                <div className="absolute inset-0 pointer-events-none">
                                    <Sparkles className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary animate-pulse" />
                                </div>
                            )}
                        </div>
                        <Button
                            size="icon"
                            variant="outline"
                            onClick={() => isRevealed ? handleMaskField(fieldName) : handleRevealField(fieldName)}
                            disabled={isRevealing === fieldName}
                            className={cn(
                                "transition-all",
                                isRevealed && "border-primary/50 text-primary hover:bg-primary/10"
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
        }

        // editMode for sensitive field
        // - Prefer editData if user typed something
        // - If the field has been revealed, we can show revealed value as initial value
        const inputValueForEdit = editValue !== "" ? editValue : (isRevealed ? revealed!.value : "");

        return (
            <div key={fieldName} className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500">
                <Label htmlFor={fieldName} className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                    {label}
                    <Shield className="h-3 w-3 text-primary" />
                </Label>
                <div className="flex gap-2">
                    <div className="relative flex-1">
                        <Input
                            id={fieldName}
                            type={isRevealed ? "text" : "password"}
                            value={inputValueForEdit as string}
                            onChange={(e) => handleFieldChange(fieldName, e.target.value)}
                            readOnly={!editMode}
                            className={cn(
                                'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl shadow-inner',
                                editMode && 'border-primary/50 shadow-primary/20',
                            )}
                        />
                        {isDecrypting && (
                            <div className="absolute inset-0 pointer-events-none">
                                <Sparkles className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-primary animate-pulse" />
                            </div>
                        )}
                    </div>
                    <Button
                        size="icon"
                        variant="outline"
                        onClick={() => isRevealed ? handleMaskField(fieldName) : handleRevealField(fieldName)}
                        disabled={isRevealing === fieldName}
                        className={cn(
                            "transition-all",
                            isRevealed && "border-primary/50 text-primary hover:bg-primary/10"
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

    // Use localEntry for guard (but we still use 'entry' prop for other logic)
    const current = localEntry ?? entry;


    if (!current) {
        return (
            <div className="flex flex-col items-center justify-center h-full text-center p-12 backdrop-blur-xl bg-white/40 dark:bg-zinc-900/40 border border-white/30 dark:border-zinc-700/30 rounded-3xl shadow-2xl">
                <div className="relative group mb-8">
                    <div className="absolute inset-0 bg-gradient-to-r  via-amber-500/30 to-primary rounded-3xl blur-xl opacity-40 group-hover:opacity-70 transition-all" />
                    <img src={ankhoraLogo} alt="Ankhora Logo" className="w-32 h-auto mx-auto" />
                    <Sparkles className="h-10 w-10 text-primary/50 absolute -top-4 -right-4 animate-pulse" />
                </div>
                <h3 className="text-2xl font-bold mb-4 bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent">
                    Select an entry
                </h3>
                <p className="text-lg text-muted-foreground/90 mb-8 max-w-md leading-relaxed backdrop-blur-sm">
                    Your vault entries are encrypted and protected. Choose an entry from the list to securely view its contents.
                </p>
                <Badge className="bg-white/70 dark:bg-zinc-800/70 backdrop-blur-sm border-white/50 text-primary font-semibold px-4 py-2">
                    <Shield className="h-4 w-4 mr-2" />
                    AES-256-GCM Encryption
                </Badge>
            </div>
        );
    }


    {/* ... other types with same glass treatment */ }
    const loginFields = [
        { name: 'user_name', label: 'Username', isSensitive: false },
        { name: 'password', label: 'Password', isSensitive: true },
        // { name: 'additionnal_note', label: 'Additionnal note', isSensitive: false },
    ];
    const cardFields = [
        { name: 'owner', label: 'Owner' },
        { name: 'number', label: 'Number' },
        { name: 'expiration', label: 'Expiration' },
        { name: 'cvc', label: 'CVC' },
        // { name: 'additionnal_note', label: 'Additionnal note', isSensitive: false },
    ];
    const identityFields = [
        { name: 'firstname', label: 'First name' },
        { name: 'secondname', label: 'Second name' },
        { name: 'lastname', label: 'Last name' },
        { name: 'username', label: 'Username' },
        { name: 'company', label: 'Company' },
        { name: 'social_security_number', label: 'Social security number' },
        { name: 'ID_number', label: 'ID number' },
        { name: 'driver_license', label: 'Driver license' },
        { name: 'mail', label: 'Email' },
        { name: 'telephone', label: 'Telephone' },
        { name: 'address_one', label: 'Address one' },
        { name: 'address_two', label: 'Address two' },
        { name: 'city', label: 'City' },
        { name: 'state', label: 'State' },
        { name: 'zip', label: 'Zip' },
        { name: 'country', label: 'Country' },
        // { name: 'additionnal_note', label: 'Additionnal note', isSensitive: false },
    ];
    const sshkeyFields = [
        { name: 'public_key', label: 'Public key' },
        { name: 'private_key', label: 'Private key' },
        { name: 'e_fingerprint', label: 'Fingerprint' },
        // { name: 'additionnal_note', label: 'Additionnal note', isSensitive: false },
    ];

    const readFileAsBuffer = (file: File): Promise<Uint8Array> => {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(new Uint8Array(reader.result as ArrayBuffer));
            reader.onerror = reject;
            reader.readAsArrayBuffer(file);
        });
    };

    /** Current version */
    type AttachmentStorage = "local" | "cloud" | "ipfs";

    interface Attachment {
        hash: string;
        cid?: string;
        storage?: AttachmentStorage;
        name?: string;
        transferStatus?: TransferStatus;
    }


    /** Beta version */
    type TransferStatus = "idle" | "uploading" | "success" | "error";
    const uploadToCloud = async (hash: string) => {
        const { jwtToken } = useAuthStore.getState();
        const buffer = await fetchLocalAttachmentBuffer(hash);
        const encryptedData = await encryptFile(jwtToken, buffer, vaultPassword);
        const cid = await uploadToIPFS(jwtToken, encryptedData); // TODO: change to upload to cloud
        const stellarOp = await createStellarCommit(jwtToken, cid);
        return { cid, stellarOp };
    }

    const fetchLocalAttachmentBuffer = async (hash: string) => {
        const attachment = attachments.find(att => att.hash === hash);
        if (!attachment) throw new Error("Attachment not found");
        const response = await fetch(attachment.hash);
        if (!response.ok) throw new Error("Failed to fetch attachment");
        return Buffer.from(await response.arrayBuffer());
    };

    interface ExtendedAttachment {
        hash: string;
        cid?: string;
        storage?: AttachmentStorage;
        name?: string;
        transferStatus?: TransferStatus;
    }



    const updateTransferStatus = (hash: string, status: TransferStatus) => {
        setTransferring(prev => ({ ...prev, [hash]: status }));
    };

    const updateAttachmentStorage = (hash: string, storage: AttachmentStorage, cid?: string) => {
        // Update your main attachments array here
        console.log(`Updated ${hash} → ${storage}${cid ? ` (CID: ${cid})` : ""}`);
        // Example: setAttachments(prev => prev.map(att => att.hash === hash ? {...att, storage, cid} : att))
    };

    const onTransferToCloud = async (attachment: ExtendedAttachment) => {
        const hash = attachment.hash;
        updateTransferStatus(hash, "uploading");

        try {
            // Your cloud upload logic
            await uploadToCloud(attachment.hash);
            updateAttachmentStorage(hash, "cloud");
            toast({ title: "Cloud upload complete" });
        } catch (error) {
            console.error("🚀 ~ onTransferToCloud ~ error:", error)
            updateTransferStatus(hash, "error");
            toast({ title: "Upload failed", variant: "destructive" });
        } finally {
            setTimeout(() => updateTransferStatus(hash, "idle"), 2000);
        }
    };

    const onTransferToBlockchain = async (attachment: ExtendedAttachment) => {
        const hash = attachment.hash;
        updateTransferStatus(hash, "uploading");

        try {
            const fileBuffer = await fetchLocalAttachmentBuffer(hash);
            const { jwtToken } = useAuthStore.getState();

            const encryptedData = await encryptFile(jwtToken, fileBuffer, vaultPassword);
            const cid = await uploadToIPFS(jwtToken, encryptedData);

            updateAttachmentStorage(hash, "ipfs", cid);
            toast({ title: "IPFS pinned", description: `CID: ${cid.slice(0, 16)}...` });
        } catch (error) {
            updateTransferStatus(hash, "error");
            console.error("🚀 ~ onTransferToBlockchain ~ error:", error)
            toast({ title: "IPFS upload failed", variant: "destructive" });
        } finally {
            setTimeout(() => updateTransferStatus(hash, "idle"), 2000);
        }
    };

    const copyCidToClipboard = async (cid: string) => {
        await navigator.clipboard.writeText(cid);
        toast({ title: "CID copied to clipboard" });
    };
    // 2. Parent component with local state
    const AttachmentsSection = ({ attachments }: { attachments: ExtendedAttachment[] }) => {

        return (
            <RenderAttachements
                attachments={attachments}
                transferring={transferring}
                onTransferToCloud={onTransferToCloud}
                onTransferToBlockchain={onTransferToBlockchain}
                onCopyCid={copyCidToClipboard}
            />
        );
    };

    // 3. Updated RenderAttachements
    const RenderAttachements = ({
        attachments,
        transferring,
        onTransferToCloud,
        onTransferToBlockchain,
        onCopyCid,
    }: {
        attachments: ExtendedAttachment[];
        transferring: Record<string, TransferStatus>;
        onTransferToCloud: (attachment: ExtendedAttachment) => void;
        onTransferToBlockchain: (attachment: ExtendedAttachment) => void;
        onCopyCid: (cid: string) => void;
    }) => {
        if (!attachments?.length) return null;


        const onFullscreen = (attachment: ExtendedAttachment) => {
            setSelectedAttachment(attachment);
        };


        const [deletePending, setDeletePending] = useState<string | null>(null);

        const onDeleteAttachment = (attachment: ExtendedAttachment) => {
            alert(`Delete "${attachment.name || attachment.hash}"?`); // Forces gesture
            const confirmed = window.confirm(`Confirm delete?`);
            if (deletePending === attachment.hash) {
                // Double click = delete
                setAttachments(prev => prev.filter(att => att.hash !== attachment.hash));
                toast({ title: "Attachment deleted" });
                setDeletePending(null);
            } else {
                // First click = preview
                setDeletePending(attachment.hash);
                toast({
                    title: "Double-click to delete",
                    duration: 2000
                });
                setTimeout(() => setDeletePending(null), 3000);
            }
        };



        return (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {attachments.map((attachment, index) => (
                    <AttachmentPreview
                        key={attachment.hash}
                        attachment={attachment}
                        transferring={transferring[attachment.hash]}
                        onTransferToCloud={onTransferToCloud}
                        onTransferToBlockchain={onTransferToBlockchain}
                        onCopyCid={onCopyCid}
                        onFullscreen={onFullscreen}
                        onDelete={onDeleteAttachment}
                    />
                ))}
            </div>
        );
    };

    // 4. Final AttachmentPreview (no custom hooks needed)
    const AttachmentPreview = ({
        attachment,
        transferring,
        onTransferToCloud,
        onTransferToBlockchain,
        onCopyCid,
        onFullscreen,
        onDelete,
    }: {
        attachment: ExtendedAttachment;
        transferring?: TransferStatus;
        onTransferToCloud: (attachment: ExtendedAttachment) => void;
        onTransferToBlockchain: (attachment: ExtendedAttachment) => void;
        onCopyCid: (cid: string) => void;
        onFullscreen: (attachment: ExtendedAttachment) => void;
        onDelete: (attachment: ExtendedAttachment) => void;
    }) => {
        const [src, setSrc] = useState("");
        const [showCid, setShowCid] = useState(false);

        // ... your existing useEffect for src ...
        useEffect(() => {
            let isMounted = true;
            fetchAttachment(attachment.hash)
                .then((url) => {
                    if (isMounted && url) setSrc(url as string);
                })
                .catch(console.error);

            return () => {
                isMounted = false;
            };
        }, [attachment.hash]);

        const isIpfs = attachment.storage === "ipfs" || !!attachment.cid;
        const isLocal = !attachment.storage || attachment.storage === "local";
        const isTransferring = transferring === "uploading";

        if (!src) {
            return (
                <div className="animate-pulse bg-zinc-200 dark:bg-zinc-800 rounded-2xl w-full h-32 flex items-center justify-center text-sm text-zinc-500">
                    Loading image...
                </div>
            );
        }

        return (
            <>
                {selectedAttachment && selectedAttachment.hash === attachment.hash && (
                    <FullscreenAttachmentModal
                        attachment={attachment}
                        onClose={() => setSelectedAttachment(null)}
                        src={src}
                    />
                )}
                <div className="group relative overflow-hidden rounded-2xl border border-white/20 dark:border-zinc-700/20 shadow-xl bg-black/5">
                    <img src={src} alt={attachment.name || attachment.hash} className="w-full h-auto object-cover" />

                    {/* ← TRASH ICON - CENTERED, HOVER ONLY */}
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            e.preventDefault(); // ← Extra safety
                            onDelete(attachment);
                        }}
                        className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-40 opacity-0 group-hover:opacity-100 transition-all duration-300 scale-0 group-hover:scale-100 pointer-events-auto"
                        title="Delete attachment"
                    >
                        <Trash2
                            className="h-8 w-8 text-red-400 bg-red-500/20 hover:bg-red-500/30 rounded-2xl p-2 backdrop-blur-md shadow-2xl border border-red-500/30 transition-all hover:scale-110 hover:text-red-300"
                        />
                    </button>

                    {/* ← ADD FULLSCREEN BUTTON top-left */}
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            onFullscreen(attachment);
                        }}
                        className="absolute right-3 top-3 z-20 rounded-full bg-white/15 p-2 backdrop-blur-md hover:bg-white/25 transition-all text-white shadow-lg"
                        title="Fullscreen preview"
                    >
                        <Maximize2 className="h-4 w-4" />
                    </button>

                    {/* Status badges - top left */}
                    <div className="absolute left-3 top-3 z-20 flex gap-2">
                        {isIpfs && (
                            <button
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setShowCid(true);
                                }}
                                className="inline-flex items-center gap-2 rounded-full border border-cyan-400/30 bg-cyan-500/20 px-3 py-1 text-[11px] font-semibold text-cyan-100 backdrop-blur-md hover:bg-cyan-500/30 transition-all"
                                title="View IPFS details"
                            >
                                <span className="h-2 w-2 rounded-full bg-cyan-300" />
                                IPFS
                            </button>
                        )}
                        {isLocal && (
                            <span className="inline-flex items-center rounded-full border border-white/20 bg-white/15 px-3 py-1 text-[11px] font-semibold text-white/80 backdrop-blur-md">
                                Local
                            </span>
                        )}
                        {isTransferring && (
                            <div className="inline-flex items-center rounded-full bg-primary/20 px-3 py-1 text-[11px] font-semibold text-primary backdrop-blur-md animate-pulse">
                                <Loader2 className="h-3 w-3 animate-spin mr-1" />
                                Transfer
                            </div>
                        )}
                    </div>

                    {/* Hover overlay - bottom */}
                    <div className="pointer-events-none absolute inset-0 opacity-0 group-hover:opacity-100 transition-all duration-300 bg-gradient-to-t from-black/80 via-black/30 to-transparent">
                        <div className="pointer-events-auto absolute inset-x-0 bottom-0 p-4">
                            <div className="flex items-end justify-between gap-3">
                                <div className="min-w-0">
                                    <p className="text-xs font-medium text-white/90 truncate">{attachment.name || attachment.hash}</p>
                                    <p className="mt-1 text-[11px] text-white/60 truncate font-mono">{attachment.hash.slice(0, 12)}...</p>
                                </div>
                                <div className="flex gap-2 flex-shrink-0">
                                    {isLocal && !isTransferring && (
                                        <>
                                            <button
                                                disabled={isTransferring}
                                                className="rounded-full bg-white/15 px-3 py-1.5 text-[11px] font-semibold text-white backdrop-blur-md hover:bg-white/25 transition-all disabled:opacity-50"
                                                onClick={() => onTransferToCloud(attachment)}
                                            >
                                                Cloud
                                            </button>
                                            <button
                                                disabled={isTransferring}
                                                className="rounded-full bg-[#C9A44A]/25 px-3 py-1.5 text-[11px] font-semibold text-[#F3DFA6] backdrop-blur-md hover:bg-[#C9A44A]/40 transition-all disabled:opacity-50"
                                                onClick={() => onTransferToBlockchain(attachment)}
                                            >
                                                IPFS
                                            </button>
                                        </>
                                    )}
                                    {isIpfs && (
                                        <button
                                            className="rounded-full bg-cyan-500/25 px-3 py-1.5 text-[11px] font-semibold text-cyan-100 backdrop-blur-md hover:bg-cyan-500/40 transition-all"
                                            onClick={() => onCopyCid(attachment.cid!)}
                                        >
                                            Copy CID
                                        </button>
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* CID details panel */}
                    {showCid && isIpfs && (
                        <div className="absolute top-14 left-3 right-3 z-30 rounded-xl border border-white/20 bg-black/80 backdrop-blur-xl shadow-2xl p-3 text-xs text-white">
                            <div className="flex items-start justify-between gap-3">
                                <div className="min-w-0 flex-1">
                                    <p className="font-semibold text-white/90 mb-2">IPFS Content ID</p>
                                    <div className="flex items-center gap-2">
                                        <code className="break-all font-mono text-[11px] text-white/95 flex-1 bg-black/20 px-2 py-1 rounded">
                                            {attachment.cid}
                                        </code>
                                        <button
                                            onClick={() => onCopyCid(attachment.cid!)}
                                            className="rounded-full bg-white/15 px-3 py-1.5 text-[11px] font-semibold text-white hover:bg-white/30 transition-all whitespace-nowrap"
                                            title="Copy to clipboard"
                                        >
                                            Copy
                                        </button>
                                    </div>
                                </div>
                                <button
                                    onClick={() => setShowCid(false)}
                                    className="rounded-full border border-white/20 bg-white/10 px-2 py-1 text-[11px] text-white hover:bg-white/25 transition-all flex-shrink-0"
                                >
                                    ✕
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </>
        );
    };


    const FullscreenAttachmentModal = ({
        attachment,
        onClose,
        src
    }: {
        attachment: ExtendedAttachment;
        onClose: () => void;
        src: string;
    }) => (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/90 backdrop-blur-sm p-4">
            <div className="flex flex-col lg:flex-row gap-8 max-w-6xl max-h-screen overflow-hidden">
                {/* Fullscreen image */}
                <div className="flex-1 flex items-center justify-center min-h-[60vh] lg:min-h-[70vh]">
                    <img
                        src={src} // ← pass src from parent or refetch
                        alt={attachment.name}
                        className="max-w-full max-h-full object-contain rounded-2xl shadow-2xl"
                    />
                </div>

                {/* ← METADATA PANEL */}
                <div className="hidden w-full lg:w-80 lg:max-w-sm flex flex-col justify-between bg-white/10 dark:bg-zinc-900/60 backdrop-blur-xl rounded-3xl border border-white/20 p-6 shadow-2xl max-h-[70vh] overflow-y-auto">
                    {/* Header */}
                    <div className="mb-6">
                        <h3 className="text-xl font-bold text-white mb-2 truncate">
                            {attachment.name || "Untitled"}
                        </h3>
                        <div className="flex items-center gap-4 text-sm text-white/70">
                            <span>Size: {formatFileSize(attachment?.size || 0)}</span>
                            <span>Hash: {attachment.hash.slice(0, 16)}...</span>
                        </div>
                    </div>

                    {/* Storage status */}
                    <div className="space-y-4 mb-8">
                        <div>
                            <span className="text-xs font-semibold text-white/60 uppercase tracking-wide mb-2 block">
                                Storage
                            </span>
                            <div className="flex items-center gap-3">
                                {attachment.storage === "ipfs" && (
                                    <div className="inline-flex items-center gap-2 rounded-full border border-cyan-400/40 bg-cyan-500/10 px-3 py-1.5 text-xs font-semibold text-cyan-200">
                                        <span className="h-2 w-2 rounded-full bg-cyan-300" />
                                        Decentralized (IPFS)
                                    </div>
                                )}
                                {attachment.storage === "cloud" && (
                                    <div className="inline-flex items-center gap-2 rounded-full border border-blue-400/40 bg-blue-500/10 px-3 py-1.5 text-xs font-semibold text-blue-200">
                                        <Cloud className="h-3 w-3" />
                                        Cloud Storage
                                    </div>
                                )}
                                {attachment.storage === "local" && (
                                    <div className="inline-flex items-center gap-2 rounded-full border border-orange-400/40 bg-orange-500/10 px-3 py-1.5 text-xs font-semibold text-orange-200">
                                        <HardDrive className="h-3 w-3" />
                                        Local
                                    </div>
                                )}
                            </div>
                        </div>

                        {attachment.cid && (
                            <div>
                                <span className="text-xs font-semibold text-white/60 uppercase tracking-wide mb-2 block">
                                    IPFS CID
                                </span>
                                <div className="flex items-center gap-2">
                                    <code className="flex-1 break-all font-mono text-sm text-white/90 bg-black/20 px-3 py-2 rounded-xl">
                                        {attachment.cid}
                                    </code>
                                    <Button
                                        size="sm"
                                        variant="ghost"
                                        onClick={() => copyCidToClipboard(attachment.cid!)}
                                        className="h-8 w-8 p-0"
                                    >
                                        <Copy className="h-3 w-3" />
                                    </Button>
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Footer actions */}
                    <div className="flex flex-col sm:flex-row gap-3 pt-4 border-t border-white/10">
                        <Button
                            variant="outline"
                            className="flex-1 bg-white/20 text-white hover:bg-white/30 border-white/30"
                            onClick={onClose}
                        >
                            Close
                        </Button>
                        {attachment.storage === "local" && (
                            <div className="flex gap-2">
                                <Button
                                    size="sm"
                                    variant="secondary"
                                    className="bg-white/20 text-white hover:bg-white/30"
                                    onClick={() => onTransferToCloud(attachment)}
                                >
                                    Upload Cloud
                                </Button>
                                <Button
                                    size="sm"
                                    className="bg-gradient-to-r from-[#C9A44A] to-amber-500 text-black hover:from-[#B8934A]"
                                    onClick={() => onTransferToBlockchain(attachment)}
                                >
                                    Pin IPFS
                                </Button>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* ESC key handler */}
            <button
                className="fixed inset-0 z-40"
                onClick={onClose}
                tabIndex={-1}
            />
        </div>
    );




    return (
        <div className="flex flex-col h-full backdrop-blur-xl bg-gradient-to-br from-white/50 via-white/30 to-zinc-50/20 dark:from-zinc-900/50 dark:via-zinc-900/30 dark:to-black/20 border border-white/20 dark:border-zinc-700/20 shadow-2xl">
            {/* Glass Header */}
            <div className="sticky top-0 z-20 backdrop-blur-2xl bg-white/70 dark:bg-zinc-900/70 border-b border-white/40 dark:border-zinc-700/40 p-6">
                <div className="flex items-start justify-between mb-6">
                    <div className="space-y-3">
                        <h2 className="text-4xl font-black bg-gradient-to-r from-foreground via-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-lg">
                            {current?.entry_name ?? ""}
                        </h2>
                        <div className="flex items-center gap-3">
                            <Badge variant="outline" className="backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border-primary/30 text-primary shadow-sm">
                                {current?.type}
                            </Badge>
                            <Badge variant="outline" className="backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border-primary/30 text-primary shadow-sm">
                                <Shield className="h-3 w-3 mr-1" />
                                End-to-End Encrypted
                            </Badge>
                        </div>
                    </div>
                    <div className="flex gap-2 backdrop-blur-sm">
                        {editMode ? (
                            <>
                                <Button
                                    size="sm"
                                    variant="outline"
                                    onClick={onCancel}
                                    className="h-11 px-6 rounded-2xl backdrop-blur-sm bg-white/70 dark:bg-zinc-800/70 border-white/50 hover:bg-white/90 dark:hover:bg-zinc-800/90 shadow-sm hover:shadow-md transition-all font-semibold"
                                >
                                    Cancel
                                </Button>
                                <Button
                                    size="sm"
                                    onClick={handleSaveEdit}
                                    className="h-11 px-8 bg-gradient-to-r from-primary to-amber-500 hover:from-primary/90 hover:to-amber-500/90 shadow-xl hover:shadow-primary/30 rounded-2xl font-semibold text-lg transition-all"
                                >
                                    Save Changes
                                </Button>
                            </>
                        ) : (
                            <>
                                {onEdit && (
                                    <Button
                                        size="icon"
                                        variant="ghost"
                                        onClick={onEdit}
                                        className="h-12 w-12 rounded-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border border-white/40 hover:bg-primary/10 hover:border-primary/40 hover:shadow-md transition-all"
                                    >
                                        <Edit className="h-5 w-5" />
                                    </Button>
                                )}
                                {onShare && (
                                    <Button
                                        size="icon"
                                        variant="ghost"
                                        onClick={() => onShare(current)}
                                        className="h-12 w-12 rounded-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border border-white/40 hover:bg-primary/10 hover:border-primary/40 hover:shadow-md transition-all"
                                    >
                                        <Share2 className="h-5 w-5" />
                                    </Button>
                                )}
                                {onDelete && (
                                    <Button
                                        size="icon"
                                        variant="ghost"
                                        onClick={onDelete}
                                        className="h-12 w-12 rounded-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border border-destructive/30 hover:bg-destructive/10 hover:border-destructive/40 hover:shadow-md transition-all text-destructive hover:text-destructive"
                                    >
                                        <Trash2 className="h-5 w-5" />
                                    </Button>
                                )}
                            </>
                        )}
                    </div>
                </div>
            </div>

            {/* Scrollable Glass Content */}
            <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar p-8 space-y-8">
                {/* Metadata Glass Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500">
                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                            Entry Name
                        </Label>
                        <Input
                            value={editMode ? ((editData as any).entry_name ?? "") : (current?.entry_name ?? "")}
                            onChange={(e) => editMode && handleFieldChange("entry_name", e.target.value)}
                            readOnly={!editMode}
                            className={cn(
                                "h-14 text-2xl font-bold  backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl shadow-inner",
                                editMode && "border-primary/50 shadow-primary/20"
                            )}
                        />
                    </div>

                    <div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40  hover:shadow-2xl transition-all duration-500">

                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                            {!editMode ? "Folder" : "Select Folder *"}
                        </Label>
                        {!editMode &&
                            <>
                                <Input
                                    value={folders.find(f => f.id === current?.folder_id)?.name ?? ""}
                                    readOnly={!editMode}
                                    className={cn(
                                        "h-14 text-2xl font-bold  backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner",
                                        editMode && "border-primary/50 shadow-primary/20"
                                    )}
                                />
                            </>
                        }
                        {editMode && (
                            <>
                                <Select value={folderId} onValueChange={setFolderId}>
                                    <SelectTrigger id="entry">
                                        <SelectValue placeholder="Choose an entry from your vault" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {folders.map((folder) => (
                                            <SelectItem key={folder.id} value={folder.id}>
                                                <span className="capitalize">{folder.name}</span>
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </>)}
                    </div>
                </div>

                {/* Sensitive Fields Section */}
                <div className="backdrop-blur-xl bg-gradient-to-br from-white/40 to-zinc-50/20 dark:from-zinc-900/40 dark:to-black/20 rounded-3xl border border-white/30 dark:border-zinc-700/30 ">
                    {current?.type === 'login' || current?.type === 'card' || current?.type === 'identity' || current?.type === 'sshkey' && <div className="flex items-center justify-between">
                        <h3 className="text-2xl font-bold flex items-center gap-3 bg-gradient-to-r from-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-lg">
                            <Shield className="h-6 w-6" />
                            Sensitive Information
                        </h3>
                        <Badge className="bg-gradient-to-r to-amber-500/20 backdrop-blur-sm border-primary/30  font-semibold px-6 py-2 ">
                            End-to-End Encrypted
                        </Badge>
                    </div>}

                    {/* Your existing type-specific renderSensitiveField calls here - enhanced styling */}

                    {current?.type === 'login' && (
                        <div className="grid gap-6 lg:grid-cols-2">
                            {loginFields.map((field) => {
                                const value = editMode
                                    ? (editData as any)?.[field.name] ?? ''
                                    : (current as any)?.[field.name] ?? '';

                                return (
                                    <div key={field.name}>
                                        {renderSensitiveField(field.name, field.label, field.isSensitive)}
                                    </div>
                                );
                            })}
                        </div>
                    )}

                    {current?.type === 'card' && (
                        <div className="grid gap-6 lg:grid-cols-2">
                            {cardFields.map((field) => {
                                const value = editMode
                                    ? (editData as any)?.[field.name] ?? ''
                                    : (current as any)?.[field.name] ?? '';

                                return (
                                    <div
                                        key={field.name}
                                        className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500"
                                    >
                                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                                            {field.label}
                                        </Label>
                                        <Input
                                            value={value}
                                            onChange={(e) =>
                                                editMode && handleFieldChange(field.name, e.target.value)
                                            }
                                            readOnly={!editMode}
                                            className={cn(
                                                'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner',
                                                editMode && 'border-primary/50 shadow-primary/20',
                                            )}
                                        />
                                    </div>
                                );
                            })}
                        </div>
                    )}

                    {current?.type === 'identity' && (
                        <div className="grid gap-6 lg:grid-cols-2">
                            {identityFields.map((field) => {
                                const value = editMode
                                    ? (editData as any)?.[field.name] ?? ''
                                    : (current as any)?.[field.name] ?? '';

                                return (
                                    <div
                                        key={field.name}
                                        className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500"
                                    >
                                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                                            {field.label}
                                        </Label>
                                        <Input
                                            value={value}
                                            onChange={(e) =>
                                                editMode && handleFieldChange(field.name, e.target.value)
                                            }
                                            readOnly={!editMode}
                                            className={cn(
                                                'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner',
                                                editMode && 'border-primary/50 shadow-primary/20',
                                            )}
                                        />
                                    </div>
                                );
                            })}
                        </div>
                    )}

                    {current?.type === 'sshkey' && (
                        <div className="space-y-6">
                            {sshkeyFields.map((field) => {
                                const value = editMode
                                    ? (editData as any)?.[field.name] ?? ''
                                    : (current as any)?.[field.name] ?? '';

                                return (
                                    <div
                                        key={field.name}
                                        className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-2 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl transition-all duration-500"
                                    >
                                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                                            {field.label}
                                        </Label>
                                        <Input
                                            value={value}
                                            onChange={(e) =>
                                                editMode && handleFieldChange(field.name, e.target.value)
                                            }
                                            readOnly={!editMode}
                                            className={cn(
                                                'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner',
                                                editMode && 'border-primary/50 shadow-primary/20',
                                            )}
                                        />
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </div>
                <div className="backdrop-blur-xl bg-gradient-to-br from-white/40 to-zinc-50/20 dark:from-zinc-900/40 dark:to-black/20 rounded-3xl  border border-white/30 dark:border-zinc-700/30 ">
                    <div className="">
                        <div className={cn(
                            "group relative rounded-3xl p-10 transition-all cursor-pointer backdrop-blur-xl",
                            "border-2 border-white/30 bg-white/50 dark:bg-zinc-900/50 shadow-2xl hover:shadow-primary/20 hover:shadow-3xl",
                            "hover:border-primary/50 hover:bg-white/70 dark:hover:bg-zinc-900/70"
                        )}>
                            <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80 group-hover:text-foreground transition-all">
                                Note
                            </Label>
                            <Textarea
                                className="min-h-[120px] h-14 text-2xl font-bold  backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner"
                                value={editMode ? ((editData as any).additionnal_note ?? "") : (current?.additionnal_note ?? "")}
                                onChange={(e) => editMode && handleFieldChange("additionnal_note", e.target.value)}
                                readOnly={!editMode}
                            />
                        </div>
                    </div>
                </div>

                {/* File Upload Widget */}
                <div className="">
                    <FileUploadWidget
                        onFileSelect={setAttachedFiles}
                        value={attachedFiles}
                        maxFiles={5}
                        entry={current}
                        vaultName={vaultContext.Vault.name}
                    />
                </div>
                {/* gallery attachements */}
                <AttachmentsSection attachments={attachments} />


                {/* Timestamps */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="group backdrop-blur-xl bg-white/50 dark:bg-zinc-900/50 rounded-3xl p-6 border border-white/30 dark:border-zinc-700/30   transition-all">
                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80">
                            <Clock className="h-5 w-5" />
                            Created
                        </Label>
                        <Input
                            value={current?.created_at ? new Date(current?.created_at).toLocaleString() : ""}
                            readOnly
                            className="h-14 font-mono text-lg  backdrop-blur-sm border-0 rounded-2xl p-4 shadow-inner"
                        />
                    </div>
                    <div className="group backdrop-blur-xl bg-white/50 dark:bg-zinc-900/50 rounded-3xl p-6 border border-white/30 dark:border-zinc-700/30   transition-all">
                        <Label className="text-lg font-semibold mb-4 flex items-center gap-2 text-muted-foreground/80">
                            <Clock className="h-5 w-5" />
                            Last Updated
                        </Label>
                        <Input
                            value={current?.updated_at ? new Date(current?.updated_at).toLocaleString() : ""}
                            readOnly
                            className="h-14 font-mono text-lg  backdrop-blur-sm border-0 rounded-2xl p-4 shadow-inner"
                        />
                    </div>
                </div>

                {/* Zero-Knowledge Footer */}
                <div className="backdrop-blur-xl bg-gradient-to-r from-primary/10 to-amber-500/5 border border-primary/20 p-8 rounded-3xl shadow-xl">
                    <div className="flex items-start gap-4 p-6 rounded-2xl bg-white/40 dark:bg-zinc-900/40 backdrop-blur-sm border border-white/30">
                        <div className="w-12 h-12 rounded-2xl bg-gradient-to-r from-primary to-amber-500/80 flex items-center justify-center flex-shrink-0 shadow-lg mt-1">
                            <Shield className="h-6 w-6 text-white drop-shadow-md" />
                        </div>
                        <div className="space-y-2">
                            <h4 className="text-xl font-bold text-foreground">Zero-Knowledge Encryption</h4>
                            <p className="text-lg text-muted-foreground/90 leading-relaxed">
                                Sensitive fields remain encrypted at rest. Decryption happens only on-demand in your browser.
                                Every view action is cryptographically logged for complete auditability.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

}

// (rest of interfaces unchanged — you already have them below in your file)


