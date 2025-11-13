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

  useEffect(() => {
    return () => {
      revealedFields.forEach(field => clearTimeout(field.timeout));
    };
  }, [revealedFields]);

  useEffect(() => {
    if (entry && editMode) {
      setEditData({ entry_name: entry.entry_name });
    }
  }, [entry, editMode]);

  const handleSaveEdit = () => {
    if (onSave) {
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

      setRevealedFields(prev => new Map(prev).set(fieldName, {
        name: fieldName,
        value: plaintext,
        timeout,
      }));

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

  const renderSensitiveField = (fieldName: string, label: string, isSensitive: boolean = true) => {
    const revealed = revealedFields.get(fieldName);
    const isRevealed = !!revealed;
    const isDecrypting = decryptingField === fieldName;

    // For non-sensitive fields like username, just show as regular text
    if (!isSensitive) {
      return (
        <div className="space-y-2">
          <Label htmlFor={fieldName} className="flex items-center gap-2 text-sm font-medium">
            {label}
          </Label>
          <div className="flex gap-2">
            <Input
              id={fieldName}
              type="text"
              value={revealed?.value || "••••••••••••"}
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

    return (
      <div className="space-y-2">
        <Label htmlFor={fieldName} className="flex items-center gap-2 text-sm font-medium">
          {label}
          <Shield className="h-3 w-3 text-primary" />
        </Label>
        <div className="flex gap-2">
          <div className="relative flex-1">
            <Input
              id={fieldName}
              type={isRevealed ? "text" : "password"}
              value={isRevealed ? revealed.value : "••••••••••••"}
              readOnly
              className={cn(
                "transition-all duration-300 border-border/50",
                isRevealed && "animate-revealBlur border-primary/50 shadow-sm",
                isDecrypting && "animate-goldShimmer"
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
            {isRevealed ? (
              <EyeOff className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
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

  if (!entry) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-center p-8 bg-gradient-to-b from-background to-secondary/20">
        <div className="relative">
          <Shield className="h-20 w-20 text-primary/20 mb-4" />
          <Sparkles className="h-8 w-8 text-primary/40 absolute -top-2 -right-2 animate-pulse" />
        </div>
        <h3 className="text-lg font-semibold mb-2">
          Select an entry to view details
        </h3>
        <p className="text-sm text-muted-foreground max-w-xs">
          Your vault entries are encrypted and protected. Choose an entry from the list to securely view its contents.
        </p>
        <Badge variant="outline" className="mt-4 text-xs">
          <Shield className="h-3 w-3 mr-1 text-primary" />
          AES-256-GCM Encryption
        </Badge>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-gradient-to-b from-background to-secondary/20">
      {/* Header */}
      <div className="sticky top-0 z-10 border-border p-6 bg-background" style={{paddingBottom: "0px"}}>
        <div className="flex items-start justify-between mb-4">
          <div className="space-y-2">
            <h2 className="text-2xl font-bold">{entry.entry_name}</h2>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="capitalize">
                {entry.type}
              </Badge>
              <Badge variant="outline" className="text-xs">
                <Shield className="h-3 w-3 mr-1 text-primary" />
                Encrypted
              </Badge>
            </div>
          </div>
          <div className="flex gap-2">
            {editMode ? (
              <>
                <Button 
                  size="sm" 
                  variant="outline" 
                  onClick={onCancel}
                  className="hover:border-border transition-all"
                >
                  Cancel
                </Button>
                <Button 
                  size="sm" 
                  onClick={handleSaveEdit}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground"
                >
                  Save Changes
                </Button>
              </>
            ) : (
              <>
                {onEdit && (
                  <Button 
                    size="icon" 
                    variant="outline" 
                    onClick={onEdit}
                    className="hover:border-primary/50 hover:text-primary transition-all"
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                )}
                {onShare && (
                  <Button 
                    size="icon" 
                    variant="outline" 
                    onClick={() => onShare(entry)}
                    className="hover:border-primary/50 hover:text-primary transition-all"
                  >
                    <Share2 className="h-4 w-4" />
                  </Button>
                )}
                {onDelete && (
                  <Button 
                    size="icon" 
                    variant="outline" 
                    onClick={onDelete}
                    className="hover:border-destructive hover:text-destructive transition-all"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </>
            )}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {/* Editable fields */}
        <div className="space-y-4">
          <div className="space-y-2">
            <Label className="text-sm">Entry Name</Label>
            <Input 
              value={editMode ? (editData.entry_name || '') : entry.entry_name} 
              onChange={(e) => editMode && setEditData({ ...editData, entry_name: e.target.value })}
              readOnly={!editMode}
              className={editMode ? "border-primary/50" : ""}
            />
          </div>

          <div className="space-y-2">
            <Label className="text-sm">Entry ID</Label>
            <Input value={entry.id} readOnly className="font-mono text-xs" />
          </div>

          <div className="space-y-2">
            <Label className="text-sm">Created</Label>
            <Input value={new Date(entry.created_at).toLocaleString()} readOnly />
          </div>

          <div className="space-y-2">
            <Label className="text-sm">Last Updated</Label>
            <Input value={new Date(entry.updated_at).toLocaleString()} readOnly />
          </div>
        </div>

        {/* Sensitive fields based on entry type */}
        <div className="border-t pt-6 space-y-4">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Shield className="h-4 w-4 text-primary" />
              Sensitive Information
            </h3>
            <Badge variant="outline" className="text-xs">
              End-to-end encrypted
            </Badge>
          </div>

          {entry.type === 'login' && (
            <>
              {renderSensitiveField('user_name', 'Username', false)}
              {renderSensitiveField('password', 'Password', true)}
              {(entry as any).web_site && (
                <div className="space-y-2">
                  <Label className="text-sm">Website URL</Label>
                  <Input value={(entry as any).web_site} readOnly />
                </div>
              )}
            </>
          )}

          {entry.type === 'card' && (
            <>
              {(entry as any).owner && (
                <div className="space-y-2">
                  <Label className="text-sm">Card Owner</Label>
                  <Input value={(entry as any).owner} readOnly />
                </div>
              )}
              {renderSensitiveField('number', 'Card Number', true)}
              {renderSensitiveField('cvc', 'CVC', true)}
              {(entry as any).expiration && (
                <div className="space-y-2">
                  <Label className="text-sm">Expiration</Label>
                  <Input value={(entry as any).expiration} readOnly />
                </div>
              )}
            </>
          )}

          {entry.type === 'note' && (
            <>
              {(entry as any).additionnal_note && (
                <div className="space-y-2">
                  <Label className="text-sm">Note Content</Label>
                  <div className="p-3 bg-secondary/30 rounded-md border border-border/50 text-sm whitespace-pre-wrap">
                    {(entry as any).additionnal_note}
                  </div>
                </div>
              )}
            </>
          )}

          {entry.type === 'sshkey' && (
            <>
              {(entry as any).public_key && (
                <div className="space-y-2">
                  <Label className="text-sm">Public Key</Label>
                  <div className="p-3 bg-secondary/30 rounded-md border border-border/50 text-xs font-mono break-all">
                    {(entry as any).public_key}
                  </div>
                </div>
              )}
              {renderSensitiveField('private_key', 'Private Key', true)}
              {(entry as any).e_fingerprint && (
                <div className="space-y-2">
                  <Label className="text-sm">Fingerprint</Label>
                  <Input value={(entry as any).e_fingerprint} readOnly className="font-mono text-xs" />
                </div>
              )}
            </>
          )}

          {entry.type === 'identity' && (
            <>
              {(entry as any).firstname && (
                <div className="space-y-2">
                  <Label className="text-sm">First Name</Label>
                  <Input value={(entry as any).firstname} readOnly />
                </div>
              )}
              {(entry as any).lastname && (
                <div className="space-y-2">
                  <Label className="text-sm">Last Name</Label>
                  <Input value={(entry as any).lastname} readOnly />
                </div>
              )}
              {(entry as any).mail && (
                <div className="space-y-2">
                  <Label className="text-sm">Email</Label>
                  <Input value={(entry as any).mail} readOnly />
                </div>
              )}
              {(entry as any).telephone && (
                <div className="space-y-2">
                  <Label className="text-sm">Telephone</Label>
                  <Input value={(entry as any).telephone} readOnly />
                </div>
              )}
              {(entry as any).address_one && (
                <div className="space-y-2">
                  <Label className="text-sm">Address</Label>
                  <Input value={(entry as any).address_one} readOnly />
                </div>
              )}
              {(entry as any).city && (
                <div className="space-y-2">
                  <Label className="text-sm">City</Label>
                  <Input value={(entry as any).city} readOnly />
                </div>
              )}
              {(entry as any).postal_code && (
                <div className="space-y-2">
                  <Label className="text-sm">Postal Code</Label>
                  <Input value={(entry as any).postal_code} readOnly />
                </div>
              )}
              {(entry as any).country && (
                <div className="space-y-2">
                  <Label className="text-sm">Country</Label>
                  <Input value={(entry as any).country} readOnly />
                </div>
              )}
            </>
          )}
        </div>

        <div className="bg-primary/5 border border-primary/20 p-4 rounded-lg space-y-2">
          <div className="flex items-start gap-2">
            <Shield className="h-4 w-4 text-primary flex-shrink-0 mt-0.5" />
            <div className="text-xs text-foreground space-y-1">
              <p className="font-medium">Zero-knowledge encryption</p>
              <p className="text-muted-foreground">
                Sensitive fields are encrypted at rest and decrypted only on demand. 
                All view actions are logged for audit.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
