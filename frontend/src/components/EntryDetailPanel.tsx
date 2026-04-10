import { useState, useEffect, useRef, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Eye, EyeOff, Copy, Shield, Edit, Share2, Trash2, Sparkles, Loader2, Maximize2, Cloud, HardDrive } from "lucide-react";
import { Attachment, AttachmentStorage, ENTRY_TYPE_CARD, ENTRY_TYPE_IDENTITY, ENTRY_TYPE_LOGIN, ENTRY_TYPE_SSHKEY, Folder, SettingsState, TransferStatus, UploadStorage, VaultEntry } from "@/types/vault";
import { decryptAttachment, decryptField, encryptAttachment, loadAttachment, logAuditEvent, updateEntry, uploadAttachementToIPFS, uploadToCloud } from "@/services/api";
import { toast } from "@/hooks/use-toast";
import { cn, formatFileSize, Gateways, isRenderableInBrowser } from "@/lib/utils";
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
import { devsecops_incident_v1, legal_matter_v1 } from "@/data/mockTemplate";
import { motion, AnimatePresence } from "framer-motion";
import { ExternalLink } from "lucide-react";

import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { useLocation, useNavigate } from "react-router-dom";

type CustomField = {
    id: string;
    key: string;
    value: any;
};
const EntryFooter = () => (
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
)



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
    const attachmentUrlCache = new Map<string, string>();


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
        onboarding: {
            packs: [],
            use_cases: [],
            installed_templates: [],
            completed: false
        }
    };
    const [settings, setSettings] = useState<SettingsState>(defaultSettings);
    const { vault } = useVaultStore();
    const vaultPassword = "password";

    const [selectedAttachment, setSelectedAttachment] = useState<Attachment | null>(null);
    const updateEntry = useVaultStore((state) => state.updateEntry);

    const [transferring, setTransferring] = useState<Record<string, TransferStatus>>({});

    const previousEntryId = useRef<string | null>(null);


    useEffect(() => {
        if (!vault?.Vault?.name) return
        fetchConfig(vault.Vault.name)
    }, [vault?.Vault?.name])

    const fetchConfig = async (vaultName: string) => {
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
            setCustomFields(entry.custom_fields);
        }
    }, [entry]);

    useEffect(() => {
        if (!entry?.id) return;

        // Only update when switching entry
        if (previousEntryId.current !== entry.id) {
            setAttachments(entry.attachments || []);
            previousEntryId.current = entry.id;
        }
    }, [entry?.id]);

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
            const cf = { ...editData.custom_fields, ...customFields };
            onSave({
                ...editData,
                custom_fields: cf,
            });
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
    console.log({ current })

    // const [customFields, setCustomFields] = useState<Record<string, any>>(
    //     current?.custom_fields ?? {}
    // );
    const [customFields, setCustomFields] = useState<Record<string, any>>({});
    const [tempKeys, setTempKeys] = useState<Record<string, string>>({});



    const location = useLocation();
    const navigate = useNavigate();

    const scopeFilter = new URLSearchParams(location.search).get("scope") || "entry_data";
    const [tab, setTab] = useState(scopeFilter);

    // Tabs Management
    useEffect(() => {
        setTab(scopeFilter);
    }, [scopeFilter]);

    const handleTabChange = (value: string) => {
        setTab(value);

        const params = new URLSearchParams(location.search);
        params.set("scope", value);

        navigate(`${location.pathname}?${params.toString()}`, { replace: true });
    };
    const baseList = tab === "entry_data" ? entry : entry?.attachments;



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

    // - Custom Fields

    // Update VALUE (smooth)
    const updateCustomFieldValue = (key: string, value: string) => {
        setCustomFields(prev => ({
            ...(prev ?? {}),
            [key]: value,
        }));
    };

    // Update KEY (temp state - no focus loss)
    const updateTempKey = (key: string, newKey: string) => {
        setTempKeys(prev => ({
            ...(prev ?? {}),
            [key]: newKey,
        }));
    };

    // Commit key change (on blur)
    const commitKeyChange = (oldKey: string) => {
        const newKey = tempKeys[oldKey] || oldKey;
        if (newKey.trim() && newKey !== oldKey) {
            setCustomFields(prev => {
                const fields = prev ?? {};
                const next = { ...fields };
                next[newKey] = next[oldKey];
                delete next[oldKey];
                delete tempKeys[oldKey]; // Clean up
                return next;
            });
        }
    };

    // Add
    const addCustomField = () => {
        setCustomFields(prev => {
            const fields = prev ?? {};
            let i = 1;
            let key = `new_field_${i}`;
            while (Object.prototype.hasOwnProperty.call(fields, key)) {
                i++;
                key = `new_field_${i}`;
            }
            return { ...fields, [key]: "" };
        });
    };

    // Remove
    const removeCustomField = (key: string) => {
        setCustomFields(prev => {
            const fields = prev ?? {};
            const next = { ...fields };
            delete next[key];
            return next;
        });
    };



    // - Attachements: Get raw buffer from attachment.file (no more fetch(hash))
    const fetchLocalAttachmentBuffer = async (hash: string): Promise<Uint8Array> => {
        const base64Url = await fetchAttachment(hash);
        // Extract base64 part
        const base64 = base64Url.split(",")[1];

        // Convert to bytes
        const binary = atob(base64);
        const bytes = new Uint8Array(binary.length);

        for (let i = 0; i < binary.length; i++) {
            bytes[i] = binary.charCodeAt(i);
        }

        return bytes;
    };

    const onTransferToBlockchainWithEncryption = async (attachment: Attachment) => {
        alert('--------------------- encryption path used !!!!')
        console.log('--------------------- encryption path used !!!!')
        const hash = attachment.hash;
        updateTransferStatus(hash, "uploading");

        try {
            // 1. Get RAW file bytes
            const fileBuffer = await fetchLocalAttachmentBuffer(hash);
            console.log("Uploading data preview:", fileBuffer.slice(0, 20));

            // 2. Encrypt → MUST return Uint8Array
            const encryptedBytes = await encryptAttachment(
                jwtToken,
                fileBuffer,
                vaultPassword
            );

            console.log("Encrypted size:", encryptedBytes.length);

            // 3. Upload 
            const cid = await uploadAttachementToIPFS(jwtToken, encryptedBytes);

            console.log("🌐 IPFS CID:", cid);

            updateAttachmentStorage(hash, UploadStorage.IPFS as AttachmentStorage, cid);
            toast({
                title: "✅ IPFS pinned",
                description: `CID: ${cid.slice(0, 16)}...`
            });

        } catch (error) {
            console.error("🚀 IPFS error:", error);
            updateTransferStatus(hash, "error");
            toast({ title: "❌ IPFS upload failed", variant: "destructive" });
        } finally {
            setTimeout(() => updateTransferStatus(hash, "idle"), 2000);
        }
    };
    const onTransferToBlockchain = async (attachment: Attachment) => {
        const hash = attachment.hash;
        updateTransferStatus(hash, "uploading");

        try {
            // 1. Get RAW file bytes
            const fileBuffer = await fetchLocalAttachmentBuffer(hash);

            // // 2. Encrypt → MUST return Uint8Array
            // const encryptedBase64 = await encryptAttachment(
            //     jwtToken,
            //     fileBuffer,
            //     vaultPassword
            // );

            // console.log("Encrypted size:", encryptedBase64.length);

            // 3. Upload 
            const cid = await uploadAttachementToIPFS(jwtToken, Array.from(fileBuffer));

            console.log("🌐 IPFS CID:", cid);

            updateAttachmentStorage(hash, UploadStorage.IPFS as AttachmentStorage, cid);
            toast({
                title: "✅ IPFS pinned",
                description: `CID: ${cid.slice(0, 16)}...`
            });

        } catch (error) {
            console.error("🚀 IPFS error:", error);
            updateTransferStatus(hash, "error");
            toast({ title: "❌ IPFS upload failed", variant: "destructive" });
        } finally {
            setTimeout(() => updateTransferStatus(hash, "idle"), 2000);
        }
    };

    const updateTransferStatus = (hash: string, status: TransferStatus) => {
        setTransferring(prev => ({ ...prev, [hash]: status }));
    };

    const updateAttachmentStorage = async (
        hash: string,
        storage: AttachmentStorage,
        cid?: string
    ) => {
        // 🧠 1. Optimistic UI update (instant, no flicker)
        setAttachments(prev => {
            const updated = prev.map(att =>
                att.hash === hash ? { ...att, storage, cid } : att
            );

            // 🔁 use fresh state
            syncAttachmentsAfterChange(updated);

            return updated;
        });

        // 🧠 2. Background sync (NO UI dependency)
        try {
            await withAuth((token) =>
                AppAPI.EditEntry(
                    entry!.type,
                    {
                        ...entry,
                        attachments: attachments.map(att =>
                            att.hash === hash ? { ...att, storage, cid } : att
                        ),
                    },
                    token
                )
            );
        } catch (err) {
            console.error(err);
            toast({ title: "Sync failed", variant: "destructive" });
        }
    };

    const onTransferToCloud = async (attachment: Attachment) => {
        const hash = attachment.hash;
        updateTransferStatus(hash, "uploading");

        try {
            // Your cloud upload logic
            const response = await uploadToCloud(jwtToken, attachment.hash);
            console.log("response", response);
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

    const copyCidToClipboard = async (cid: string) => {
        await navigator.clipboard.writeText(cid);
        toast({ title: "CID copied to clipboard" });
    };
    // 1. Parent component with local state
    const AttachmentsSection = ({ attachments }: { attachments: Attachment[] }) => {

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

    // 2. Updated RenderAttachements
    const RenderAttachements = ({
        attachments,
        transferring,
        onTransferToCloud,
        onTransferToBlockchain,
        onCopyCid,
    }: {
        attachments: Attachment[];
        transferring: Record<string, TransferStatus>;
        onTransferToCloud: (attachment: Attachment) => void;
        onTransferToBlockchain: (attachment: Attachment) => void;
        onCopyCid: (cid: string) => void;
    }) => {
        if (!attachments?.length) return null;
        const [deletePending, setDeletePending] = useState<string | null>(null);
        const onFullscreen = (attachment: Attachment) => {
            setSelectedAttachment(attachment);
        };

        const onDeleteAttachment = async (attachment: Attachment) => {
            if (deletePending === attachment.hash) {
                // ✨ 1. Instant UI update (no refresh, no flicker)
                setAttachments(prev => {
                    const updated = prev.filter(att => att.hash !== attachment.hash);

                    // 🔁 2. Background sync using UPDATED state
                    syncAttachmentsAfterChange(updated);

                    return updated;
                });

                toast({ title: "Attachment deleted" });
                setDeletePending(null);
            } else {
                setDeletePending(attachment.hash);
                toast({
                    title: "Double-click to delete",
                    duration: 2000
                });

                setTimeout(() => setDeletePending(null), 3000);
            }
        };

        return (
            <motion.div layout className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <AnimatePresence mode="popLayout">
                    {attachments.map((attachment) => (
                        <motion.div
                            key={attachment.hash}
                            layout
                            initial={{ opacity: 0, scale: 0.92, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.9, y: -20 }}
                            transition={{
                                type: "spring",
                                stiffness: 260,
                                damping: 22
                            }}
                        >
                            <AttachmentPreview
                                attachment={attachment}
                                transferring={transferring[attachment.hash]}
                                onTransferToCloud={onTransferToCloud}
                                onTransferToBlockchain={onTransferToBlockchain}
                                onCopyCid={onCopyCid}
                                onFullscreen={onFullscreen}
                                onDelete={onDeleteAttachment}
                            />
                        </motion.div>
                    ))}
                </AnimatePresence>
            </motion.div>
        );
    };
    const syncAttachmentsAfterChange = async (next?: Attachment[]) => {
        const updated = next ?? attachments;

        try {
            await withAuth((token) =>
                AppAPI.EditEntry(
                    entry!.type,
                    {
                        ...entry,
                        attachments: updated,
                    },
                    token
                )
            );
        } catch (err) {
            console.error(err);
        }
    };

    // 3. Final AttachmentPreview (no custom hooks needed)
    const AttachmentPreview = ({
        attachment,
        transferring,
        onTransferToCloud,
        onTransferToBlockchain,
        onCopyCid,
        onFullscreen,
        onDelete,
    }: {
        attachment: Attachment;
        transferring?: TransferStatus;
        onTransferToCloud: (attachment: Attachment) => void;
        onTransferToBlockchain: (attachment: Attachment) => void;
        onCopyCid: (cid: string) => void;
        onFullscreen: (attachment: Attachment) => void;
        onDelete: (attachment: Attachment) => void;
    }) => {
        const [src, setSrc] = useState("");
        const [showCid, setShowCid] = useState(false);

        useEffect(() => {
            let isMounted = true;

            // 🧠 1. Use cache FIRST (instant, no flicker)
            if (attachmentUrlCache.has(attachment.hash)) {
                setSrc(attachmentUrlCache.get(attachment.hash)!);
                return;
            }

            // 🧠 2. Fetch only if not cached
            fetchAttachment(attachment.hash)
                .then((url) => {
                    if (!isMounted || !url) return;

                    attachmentUrlCache.set(attachment.hash, url);
                    setSrc(url);
                })
                .catch((err) => {
                    console.error(err)
                });

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
                <div
                    className="
                        hover:shadow-[0_0_30px_rgba(201,164,74,0.25)] 
                        glass-perfect group relative overflow-hidden 
                        rounded-2xl transition-all duration-300 
                        ease-in-out hover:scale-[1.02] 
                        animate-in fade-in zoom-in-95 
                        group relative overflow-hidden 
                        rounded-2xl border border-white/20 
                        dark:border-zinc-700/20 
                        shadow-xl bg-black/5

                        hover:scale-[1.02]
                        hover:shadow-[0_0_30px_rgba(201,164,74,0.25)]
                    "
                >
                    <motion.img
                        layout
                        src={src}
                        alt={attachment.name || attachment.hash}
                        className="w-full h-auto object-cover"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ duration: 0.4 }}
                    />

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
                        <div style={{ zIndex: 1000, backgroundColor: "#000000" }} className="absolute top-14 left-3 right-3 z-30 rounded-xl border border-white/20 bg-black backdrop-blur-xl shadow-2xl p-3 text-xs text-white">
                            <div className="flex items-start justify-between gap-3" >
                                <div className="min-w-0 flex-1">
                                    <p className="font-semibold text-white/90 mb-2">IPFS Content ID</p>
                                    <div className="flex items-center gap-2">
                                        <code className="break-all font-mono text-[11px] text-white/95 flex-1 bg-black/20 px-2 py-1 rounded">
                                            {attachment.cid}
                                        </code>

                                        {/* 👁 VIEW ON IPFS */}
                                        <button
                                            onClick={() => openIpfsInBrowser(attachment)}
                                            className="
                                                rounded-full 
                                                bg-gradient-to-r from-[#C9A44A]/30 to-amber-500/30 
                                                px-3 py-1.5 
                                                text-[11px] font-semibold text-[#F3DFA6] 
                                                backdrop-blur-md 
                                                hover:from-[#C9A44A]/50 hover:to-amber-500/50
                                                transition-all 
                                                shadow-lg
                                                border border-[#C9A44A]/30
                                                hover:scale-105 active:scale-95
                                            "
                                            title="View on IPFS"
                                        >
                                            👁 View
                                        </button>

                                        {/* COPY */}
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
        attachment: Attachment;
        onClose: () => void;
        src: string;
    }) => (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/90 backdrop-blur-sm p-4">
            <div className="flex flex-col lg:flex-row gap-8 max-w-6xl max-h-screen overflow-hidden">
                {/* Fullscreen image */}
                <div className="flex-1 flex items-center justify-center min-h-[60vh] lg:min-h-[70vh]">
                    <motion.img
                        layoutId={attachment.hash}
                        src={src} // ← pass src from parent or refetch
                        alt={attachment.name}
                        className="max-w-full max-h-full object-contain rounded-2xl shadow-2xl"
                    ></motion.img>
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

    const openPrivateIpfsInBrowser = async (attachment: Attachment) => {
        try {
            console.log("Fetching CID from backend:", attachment.cid);

            // ✅ DIRECT BACKEND CALL (NO FETCH)
            const base64 = await AppAPI.GetIPFSFile(jwtToken, attachment.cid);

            // decode → bytes
            const binaryString = atob(base64);
            const bytes = new Uint8Array(binaryString.length);

            for (let i = 0; i < binaryString.length; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }

            const decryptedBuffer = await decryptAttachment(
                jwtToken,
                bytes,
                vaultPassword
            );

            if (!decryptedBuffer || decryptedBuffer.length === 0) {
                throw new Error("Decryption returned empty buffer");
            }

            const blob = new Blob([decryptedBuffer], {
                type: "image/jpeg"
            });

            console.log("Blob size:", blob.size);

            const objectUrl = URL.createObjectURL(blob);
            console.log({ objectUrl })

            AppAPI.OpenURL(objectUrl);

        } catch (err) {
            console.error("Decrypt view failed:", err);
        }
    };
    const openIpfsInBrowser = async (attachment: Attachment) => {
        try {
            const url = `${Gateways.local}/${attachment.cid}`;
            AppAPI.OpenURL(url);
        } catch (err) {
            console.error("Decrypt view failed:", err);
        }
    };
    function base64ToUint8Array(base64: string): Uint8Array {
        const binaryString = atob(base64);
        const len = binaryString.length;
        const bytes = new Uint8Array(len);

        for (let i = 0; i < len; i++) {
            bytes[i] = binaryString.charCodeAt(i);
        }

        return bytes;
    }

    // - Helpers - Custom Fields
    const RESERVED = new Set(["template_id", "record_type", "schema_version"]); // Ankhora templates identifiers

    function isObject(v: any) {
        return v && typeof v === "object" && !Array.isArray(v);
    }
    function isPlainObject(v: any) {
        return v && typeof v === "object" && !Array.isArray(v);
    }

    function isArrayOfObjects(v: any) {
        return Array.isArray(v) && v.some(isPlainObject);
    }

    function formatValue(v: any) {
        if (v === null || v === undefined) return "";
        if (typeof v === "string") return v;
        if (typeof v === "number" || typeof v === "boolean") return String(v);

        if (Array.isArray(v)) {
            // array of primitives
            if (!isArrayOfObjects(v)) return v.map((x) => String(x)).join(", ");
            // array of objects -> pretty JSON
            return JSON.stringify(v, null, 2);
        }

        // object
        return JSON.stringify(v, null, 2);
    }


    if (current) {
        // For testing
        // current.custom_fields = devsecops_incident_v1   // legal_matter_v1
    }

    // - Helpers - Attachements Cleanup
    const cleanUpAttachementsList = async () => {
        attachments.map(async (attachement) => {

            fetchAttachment(attachement.hash)
                .then((url) => {
                    if (!url) return;
                })
                .catch((err) => {
                    // if (err == "Enregistrement Introuvable") {
                    console.error(err)
                    setAttachments(prev => prev.filter(att => att.hash !== attachement.hash));
                    current.attachments = current.attachments.filter(att => att.hash !== attachement.hash);
                    // }
                    console.log('error fetch attachement', err)
                });
        })
    }
    const fuckingUpdate = async () => {
        updateEntry(current.id, current);
        AppAPI.EditEntry(current.type, current, jwtToken);
    }


    return (
        <div className="flex flex-col backdrop-blur-xl bg-gradient-to-br from-white/50 via-white/30 to-zinc-50/20 dark:from-zinc-900/50 dark:via-zinc-900/30 dark:to-black/20 border border-white/20 dark:border-zinc-700/20 shadow-2xl">


            {/* UI Tabs */}
            <Tabs
                value={tab}
                onValueChange={handleTabChange}
            >
                {/* Glass Header */}
                <div className="sticky top-0 z-20 backdrop-blur-2xl bg-white dark:bg-zinc-900 border-b border-white/40 dark:border-zinc-700/40 p-6">
                    <div className="flex items-start justify-between mb-6">
                        <div className="space-y-3">
                            <h2 className="text-4xl font-black bg-gradient-to-r from-foreground via-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-lg">
                                {current?.entry_name ?? ""}
                            </h2>
                            <button className="text-primary btn btn-sm" onClick={() => cleanUpAttachementsList()}>Clean up attachments list</button>
                            <button className="text-primary btn btn-sm" onClick={() => fuckingUpdate()}>Fucking update</button>
                            <div className="flex items-center justify-between gap-3">
                                <Badge variant="outline" className="backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border-primary/30 text-primary shadow-sm">
                                    {current?.type}
                                </Badge>
                                <Badge variant="outline" className="backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border-primary/30 text-primary shadow-sm">
                                    <Shield className="h-3 w-3 mr-1" />
                                    End-to-End Encrypted
                                </Badge>

                            </div>
                        </div>
                        <div className="flex items-end gap-2 backdrop-blur-sm" style={{ flexDirection: "column" }}>
                            <div className="mb-4">
                                {editMode ? (
                                    <>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={onCancel}
                                            className="ml-2 h-11 px-6 rounded-2xl backdrop-blur-sm bg-white/70 dark:bg-zinc-800/70 border-white/50 hover:bg-white/90 dark:hover:bg-zinc-800/90 shadow-sm hover:shadow-md transition-all font-semibold"
                                        >
                                            Cancel
                                        </Button>
                                        <Button
                                            size="sm"
                                            onClick={handleSaveEdit}
                                            className="ml-2 h-11 px-8 bg-gradient-to-r from-primary to-amber-500 hover:from-primary/90 hover:to-amber-500/90 shadow-xl hover:shadow-primary/30 rounded-2xl font-semibold text-lg transition-all"
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
                                                className="ml-2 h-12 w-12 rounded-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border border-white/40 hover:bg-primary/10 hover:border-primary/40 hover:shadow-md transition-all"
                                            >
                                                <Edit className="h-5 w-5" />
                                            </Button>
                                        )}
                                        {onShare && (
                                            <Button
                                                size="icon"
                                                variant="ghost"
                                                onClick={() => onShare(current)}
                                                className="ml-2 h-12 w-12 rounded-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border border-white/40 hover:bg-primary/10 hover:border-primary/40 hover:shadow-md transition-all"
                                            >
                                                <Share2 className="h-5 w-5" />
                                            </Button>
                                        )}
                                        {onDelete && (
                                            <Button
                                                size="icon"
                                                variant="ghost"
                                                onClick={onDelete}
                                                className="ml-2 h-12 w-12 rounded-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-800/60 border border-destructive/30 hover:bg-destructive/10 hover:border-destructive/40 hover:shadow-md transition-all text-destructive hover:text-destructive"
                                            >
                                                <Trash2 className="h-5 w-5" />
                                            </Button>
                                        )}
                                    </>
                                )}
                            </div>
                            <TabsList className="grid w-[190px] grid-cols-2" style={{ fontSize: "10px" }}>
                                <TabsTrigger className="text-xs" value="entry_data">Entry</TabsTrigger>
                                <TabsTrigger className="text-xs" value="entry_attachements" >Attachement{entry?.attachments?.length > 1 ? "s" : ""}</TabsTrigger>
                            </TabsList>
                        </div>
                    </div>
                </div>

                <TabsContent value="entry_data">
                    <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar p-8 space-y-8">

                        {/* Scrollable Glass Content */}
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
                    </div>
                    <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar p-8 space-y-8">
                        {/* Sensitive Fields Section */}
                        <div className="backdrop-blur-xl bg-gradient-to-br from-white/40 to-zinc-50/20 dark:from-zinc-900/40 dark:to-black/20 rounded-3xl border border-white/30 dark:border-zinc-700/30 ">
                            {(current?.type === ENTRY_TYPE_LOGIN || current?.type === ENTRY_TYPE_CARD || current?.type === ENTRY_TYPE_IDENTITY || current?.type === ENTRY_TYPE_SSHKEY) && <div className="flex items-center justify-between">
                                <h3 className="text-2xl font-bold flex items-center gap-3 bg-gradient-to-r from-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-lg">
                                    <Shield className="h-6 w-6" />
                                    Sensitive Information
                                </h3>
                                <Badge className="bg-gradient-to-r to-amber-500/20 backdrop-blur-sm border-primary/30  font-semibold px-6 py-2 ">
                                    End-to-End Encrypted
                                </Badge>
                            </div>}

                            {/* Your existing type-specific renderSensitiveField calls here - enhanced styling */}

                            {current?.type === ENTRY_TYPE_LOGIN && (
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

                            {current?.type === ENTRY_TYPE_CARD && (
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

                            {current?.type === ENTRY_TYPE_IDENTITY && (
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

                            {current?.type === ENTRY_TYPE_SSHKEY && (
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


                            {/* Timestamps */}
                            <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
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
                        </div>

                        {/* Additional Note */}
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
                                        className="min-h-[130px] h-14 text-2xl font-bold  backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner"
                                        value={editMode ? ((editData as any).additionnal_note ?? "") : (current?.additionnal_note ?? "")}
                                        onChange={(e) => editMode && handleFieldChange("additionnal_note", e.target.value)}
                                        readOnly={!editMode}
                                    />
                                </div>
                            </div>
                        </div>

                    </div>
                    {/* <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar p-8 space-y-8"> */}
                    {/* Custom Fields */}
                    <div className="backdrop-blur-xl bg-gradient-to-br from-white/40 to-zinc-50/20 dark:from-zinc-900/40 dark:to-black/20 rounded-3xl border border-white/30 dark:border-zinc-700/30 p-8">
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="text-2xl font-bold bg-gradient-to-r from-primary to-amber-500/80 bg-clip-text text-transparent">
                                Custom Fields
                            </h3>

                            <div className="flex gap-2">
                                {current?.custom_fields?.template_id && (
                                    <Badge variant="outline" className="bg-white/60 dark:bg-zinc-800/60">
                                        Template: {String(current.custom_fields.template_id)}
                                    </Badge>
                                )}
                                {current?.custom_fields?.schema_version && (
                                    <Badge variant="outline" className="bg-white/60 dark:bg-zinc-800/60">
                                        v{String(current.custom_fields.schema_version)}
                                    </Badge>
                                )}
                            </div>

                        </div>


                        {/* // ✅ Render with FIXED key (this stops focus loss) */}
                        {
                            editMode ? (
                                <div className="space-y-4">
                                    {Object.entries(customFields ?? {})
                                        .filter(([k]) => !RESERVED.has(k))
                                        .map(([key, value], index) => (
                                            <div key={index} className="grid grid-cols-1 md:grid-cols-12 gap-3 items-start rounded-2xl border p-4">
                                                <div className="md:col-span-4">
                                                    <Input
                                                        value={tempKeys[key] || key}
                                                        onChange={(e) => updateTempKey(key, e.target.value)}
                                                        onBlur={() => commitKeyChange(key)}
                                                        placeholder="Field key"
                                                    />
                                                </div>
                                                <Input
                                                    className="md:col-span-7"
                                                    value={String(value ?? "")}
                                                    onChange={(e) => updateCustomFieldValue(key, e.target.value)}
                                                    placeholder="Field value"
                                                />
                                                <Button
                                                    type="button"
                                                    variant="destructive"
                                                    className="md:col-span-1"
                                                    onClick={() => removeCustomField(key)}
                                                >
                                                    ✕
                                                </Button>
                                            </div>
                                        ))}
                                    <Button type="button" onClick={addCustomField}>
                                        Add field
                                    </Button>
                                </div>
                            ) : (
                                current?.custom_fields && Object.keys(current.custom_fields).length > 0 && (
                                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                                        {Object.entries(current.custom_fields)
                                            .filter(([k]) => !RESERVED.has(k))
                                            .map(([k, v]) => (
                                                <div key={k} className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-6 border border-white/40 dark:border-zinc-700/40">
                                                    <Label className="text-lg font-semibold mb-3 text-muted-foreground/80">
                                                        {k}
                                                    </Label>
                                                    {(isPlainObject(v) || isArrayOfObjects(v)) ? (
                                                        <pre className="text-sm whitespace-pre-wrap break-words font-mono bg-white/40 dark:bg-zinc-800/40 p-4 rounded-2xl border border-white/30 dark:border-zinc-700/30">
                                                            {formatValue(v)}
                                                        </pre>
                                                    ) : (
                                                        <Input value={formatValue(v)} readOnly className="h-12 text-lg font-semibold" />
                                                    )}
                                                </div>
                                            ))}
                                    </div>
                                )
                            )
                        }
                    </div>

                    {/* </div> */}
                </TabsContent>

                <TabsContent value="entry_attachements">
                    <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar p-8 space-y-8">
                        {/* File Upload Widget */}
                        <div className="">
                            <FileUploadWidget
                                onFileSelect={setAttachedFiles}
                                value={attachedFiles}
                                maxFiles={5}
                                entry={current}
                                vaultName={vaultContext.Vault.name}
                                attachments={attachments}
                                setAttachments={setAttachments}
                            />
                        </div>
                        {/* gallery attachements */}
                        <AttachmentsSection attachments={attachments} />
                    </div>
                </TabsContent>
            </Tabs>


            {/* Zero-Knowledge Footer */}
            <EntryFooter />

        </div >
    );
}


