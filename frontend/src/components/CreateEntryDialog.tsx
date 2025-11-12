import { useState } from "react";
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
import { VaultEntry } from "@/types/vault";

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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const baseEntry = {
      entry_name: entryName,
      type,
      trashed: false,
      is_draft: false,
      is_favorite: false,
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
