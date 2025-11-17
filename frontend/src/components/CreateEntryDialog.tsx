import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { VaultEntry, OOFormDocument } from "@/types/vault";
import { FileText, Upload, X } from "lucide-react";

const MAX_OO_FORM_SIZE = 10 * 1024 * 1024; // 10 MB
const ACCEPTED_OO_FORM_TYPES = ".pdf,.doc,.docx,.odt,.txt,.png,.jpg,.jpeg";

const formatBytes = (bytes: number) => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
};

interface CreateEntryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (entry: Omit<VaultEntry, "id" | "created_at" | "updated_at">) => void;
}

export function CreateEntryDialog({ open, onOpenChange, onSubmit }: CreateEntryDialogProps) {
  const [entryName, setEntryName] = useState("");
  const [type, setType] = useState<VaultEntry['type']>("login");
  const [username, setUsername] = useState("");
  const [url, setUrl] = useState("");
  const [note, setNote] = useState("");
  const [ooFormDocument, setOoFormDocument] = useState<OOFormDocument | null>(null);
  const [ooFormError, setOoFormError] = useState<string | null>(null);
  const ooFormInputRef = useRef<HTMLInputElement | null>(null);

  const resetOoFormInput = () => {
    if (ooFormInputRef.current) {
      ooFormInputRef.current.value = "";
    }
  };

  const handleOoFormRemove = () => {
    setOoFormDocument(null);
    setOoFormError(null);
    resetOoFormInput();
  };

  const handleOoFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];

    if (!file) {
      return;
    }

    if (file.size > MAX_OO_FORM_SIZE) {
      setOoFormError("The OO form document exceeds the 10 MB limit.");
      setOoFormDocument(null);
      resetOoFormInput();
      return;
    }

    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result;
      if (typeof result === "string") {
        const base64 = result.includes(",") ? result.split(",")[1] : result;
        setOoFormDocument({
          name: file.name,
          size: file.size,
          type: file.type || "application/octet-stream",
          lastModified: file.lastModified,
          base64,
        });
        setOoFormError(null);
      } else {
        setOoFormError("Unable to read the OO form document. Please try again.");
        setOoFormDocument(null);
      }
    };
    reader.onerror = () => {
      setOoFormError("Failed to process the OO form document.");
      setOoFormDocument(null);
      resetOoFormInput();
    };

    reader.readAsDataURL(file);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const baseEntry = {
      entry_name: entryName,
      type,
      trashed: false,
      is_draft: false,
      is_favorite: false,
      oo_form_document: ooFormDocument,
    };

    let newEntry: Omit<VaultEntry, "id" | "created_at" | "updated_at">;
    
    switch (type) {
      case 'login':
        newEntry = {
          ...baseEntry,
          type: 'login' as const,
          user_name: username,
          password: "",
          web_site: url,
        } as Omit<VaultEntry, "id" | "created_at" | "updated_at">;
        break;
      case 'note':
        newEntry = {
          ...baseEntry,
          type: 'note' as const,
          additionnal_note: note,
        } as Omit<VaultEntry, "id" | "created_at" | "updated_at">;
        break;
      default:
        newEntry = {
          ...baseEntry,
        } as Omit<VaultEntry, "id" | "created_at" | "updated_at">;
    }

    onSubmit(newEntry);
    
    // Reset form
    setEntryName("");
    setType("login");
    setUsername("");
    setUrl("");
    setNote("");
    handleOoFormRemove();
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px] bg-card border-border">
        <DialogHeader>
          <DialogTitle className="text-foreground">Create New Entry</DialogTitle>
          <DialogDescription className="text-muted-foreground">
            Add a new encrypted entry to your sovereign vault.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="entryName">Entry Name</Label>
              <Input
                id="entryName"
                placeholder="e.g., GitHub Account"
                value={entryName}
                onChange={(e) => setEntryName(e.target.value)}
                required
                className="bg-background border-border"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="type">Entry Type</Label>
              <Select value={type} onValueChange={(val) => setType(val as VaultEntry['type'])} required>
                <SelectTrigger className="bg-background border-border">
                  <SelectValue placeholder="Select entry type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="login">Login</SelectItem>
                  <SelectItem value="card">Payment Card</SelectItem>
                  <SelectItem value="identity">Identity</SelectItem>
                  <SelectItem value="note">Secure Note</SelectItem>
                  <SelectItem value="sshkey">SSH Key</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {type === 'login' && (
              <>
                <div className="space-y-2">
                  <Label htmlFor="username">Username / Email</Label>
                  <Input
                    id="username"
                    placeholder="user@example.com"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="bg-background border-border"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="url">Website URL</Label>
                  <Input
                    id="url"
                    placeholder="https://example.com"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    className="bg-background border-border"
                  />
                </div>
              </>
            )}

            {type === 'note' && (
              <div className="space-y-2">
                <Label htmlFor="note">Note Content</Label>
                <Textarea
                  id="note"
                  placeholder="Enter your secure note..."
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  className="bg-background border-border min-h-[100px]"
                />
              </div>
            )}

              <div className="space-y-3">
                <div className="space-y-1">
                  <Label htmlFor="ooForm">OO Form Document</Label>
                  <p className="text-sm text-muted-foreground">
                    Upload the signed OO (Operational Order) form to keep a verifiable document with this entry.
                  </p>
                </div>

                <input
                  ref={ooFormInputRef}
                  id="ooForm"
                  type="file"
                  accept={ACCEPTED_OO_FORM_TYPES}
                  className="hidden"
                  onChange={handleOoFormChange}
                />

                <div className="rounded-lg border border-dashed border-border bg-muted/20 p-4 text-center">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="border-border"
                    onClick={() => ooFormInputRef.current?.click()}
                  >
                    <Upload className="mr-2 h-4 w-4" />
                    Select OO Form Document
                  </Button>
                  <p className="mt-2 text-xs text-muted-foreground">
                    Accepted: PDF, DOC/DOCX, ODT, TXT, JPG/PNG · Max 10 MB
                  </p>
                </div>

                {ooFormDocument ? (
                  <div className="flex items-center justify-between rounded-md border border-border bg-background/80 px-3 py-2">
                    <div className="flex items-center gap-3">
                      <FileText className="h-5 w-5 text-primary" />
                      <div className="text-left">
                        <p className="text-sm font-medium">{ooFormDocument.name}</p>
                        <p className="text-xs text-muted-foreground">
                          {formatBytes(ooFormDocument.size)} · {ooFormDocument.type || "Unknown type"}
                        </p>
                      </div>
                    </div>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="text-muted-foreground hover:text-destructive"
                      onClick={handleOoFormRemove}
                    >
                      <X className="mr-1 h-4 w-4" />
                      Remove
                    </Button>
                  </div>
                ) : (
                  <p className="text-sm italic text-muted-foreground">No OO form attached yet.</p>
                )}

                {ooFormError && <p className="text-sm text-destructive">{ooFormError}</p>}
              </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              className="border-border"
            >
              Cancel
            </Button>
            <Button type="submit" className="shadow-glow bg-primary hover:bg-primary/90">
              Create Entry
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
