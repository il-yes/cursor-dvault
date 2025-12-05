import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Eye, EyeOff, Copy, Shield, Edit, Share2, Trash2, Sparkles } from "lucide-react";
import { VaultEntry } from "@/types/vault";
import { decryptField, logAuditEvent } from "@/services/api";
import { toast } from "@/hooks/use-toast";
import { cn } from "@/lib/utils";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import { Clock } from "lucide-react";
import "./contributionGraph/g-scrollbar.css";
import { Textarea } from "./ui/textarea";
import { FileUploadWidget } from "./FileUploadWidget";

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
    // keep localEntry typed explicitly
    const [localEntry, setLocalEntry] = useState<VaultEntry | null>(entry ?? null);
    const [attachedFiles, setAttachedFiles] = useState<File[]>([])

    // Sync localEntry with prop changes
    useEffect(() => {
        setLocalEntry(entry ?? null);
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

    const handleFieldChange = (fieldName: string, value: any) => {
        setEditData(prev => ({ ...prev, [fieldName]: value }));
    };

    const handleSaveEdit = () => {
        if (onSave) {
            // pass only the edited changes (editData)
            onSave(editData);
            setEditData({});
        }
    };

    const handleRevealField = async (fieldName: string) => {
        if (!entry) return;

        setIsRevealing(fieldName);
        setDecryptingField(fieldName);

        try {
            const { plaintext, expires_in } = await decryptField({
                entry_id: entry.id,
                field_name: fieldName,
            });

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
                                    'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner',
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
                            value={inputValueForEdit as string}
                            onChange={(e) => handleFieldChange(fieldName, e.target.value)}
                            readOnly={!editMode}
                            className={cn(
                                'h-14 text-2xl font-bold backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner',
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
    ];
    const cardFields = [
        { name: 'owner', label: 'Owner' },
        { name: 'number', label: 'Number' },
        { name: 'expiration', label: 'Expiration' },
        { name: 'cvc', label: 'CVC' },
    ];
    const identityFields = [
        { name: 'first_name', label: 'First name' },
        { name: 'second_name', label: 'Second name' },
        { name: 'last_name', label: 'Last name' },
        { name: 'username', label: 'Username' },
        { name: 'company', label: 'Company' },
        { name: 'social_security_number', label: 'Social security number' },
        { name: 'ID_number', label: 'ID number' },
        { name: 'driver_license', label: 'Driver license' },
        { name: 'number', label: 'Number' },
        { name: 'telephone  ', label: 'Telephone' },
        { name: 'address_one', label: 'Address one' },
        { name: 'address_two', label: 'Address two' },
        { name: 'city', label: 'City' },
        { name: 'state', label: 'State' },
        { name: 'zip', label: 'Zip' },
        { name: 'country', label: 'Country' },
    ];

    const sshkeyFields = [
        { name: 'public_key', label: 'Public key' },
        { name: 'private_key', label: 'Private key' },
        { name: 'e_fingerprint', label: 'Fingerprint' },
    ];

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
                                "h-14 text-2xl font-bold  backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl p-4 shadow-inner",
                                editMode && "border-primary/50 shadow-primary/20"
                            )}
                        />
                    </div>

                    <div className="group backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-6 border border-white/40 dark:border-zinc-700/40  hover:shadow-2xl transition-all duration-500">
                        {/* Date createdAt and updatedAt */}
                    </div>
                </div>

                {/* Sensitive Fields Section */}
                <div className="backdrop-blur-xl bg-gradient-to-br from-white/40 to-zinc-50/20 dark:from-zinc-900/40 dark:to-black/20 rounded-3xl p-8 border border-white/30 dark:border-zinc-700/30 ">
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
                    <div className="space-y-6">
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
                    />
                </div>

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


