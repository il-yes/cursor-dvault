import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Trash2, UserPlus, Mail } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { SharedEntry, Recipient } from "@/types/sharing";
import { ScrollArea } from "@/components/ui/scroll-area";

interface RecipientsManagementModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  entry: SharedEntry | null;
}

export function RecipientsManagementModal({
  open,
  onOpenChange,
  entry,
}: RecipientsManagementModalProps) {
  const { toast } = useToast();
  const [recipients, setRecipients] = useState<Recipient[]>(entry?.recipients || []);
  const [isAddingNew, setIsAddingNew] = useState(false);
  const [newRecipientEmail, setNewRecipientEmail] = useState("");
  const [newRecipientRole, setNewRecipientRole] = useState<"viewer" | "editor" | "owner">("viewer");

  const handleRoleChange = (recipientId: string, newRole: "viewer" | "editor" | "owner") => {
    setRecipients(recipients.map(r =>
      r.id === recipientId ? { ...r, role: newRole } : r
    ));
  };

  const handleRemoveRecipient = (recipientId: string) => {
    setRecipients(recipients.filter(r => r.id !== recipientId));
    toast({
      title: "Recipient removed",
      description: "The recipient has been removed from this shared entry.",
    });
  };

  const handleAddRecipient = () => {
    if (!newRecipientEmail.trim()) return;

    const newRecipient: Recipient = {
      id: `recipient-${Date.now()}`,
      name: newRecipientEmail.split("@")[0],
      email: newRecipientEmail,
      role: newRecipientRole,
      joined_at: new Date().toISOString(),
    };

    setRecipients([...recipients, newRecipient]);
    setNewRecipientEmail("");
    setNewRecipientRole("viewer");
    setIsAddingNew(false);

    toast({
      title: "Recipient added",
      description: `${newRecipient.email} has been added as ${newRecipientRole}.`,
    });
  };

  const handleResendInvitation = (recipient: Recipient) => {
    toast({
      title: "Invitation sent",
      description: `Invitation resent to ${recipient.email}`,
    });
  };

  const handleSaveChanges = () => {
    // In real app, this would make an API call
    toast({
      title: "Changes saved",
      description: "Recipient permissions have been updated.",
    });
    onOpenChange(false);
  };

  const getRoleBadgeColor = (role: string) => {
    switch (role) {
      case "owner":
        return "bg-[#C9A44A]/10 text-[#C9A44A] border-[#C9A44A]/20";
      case "editor":
        return "bg-blue-500/10 text-blue-600 border-blue-500/20";
      default:
        return "bg-secondary";
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[80vh]">
        <DialogHeader>
          <DialogTitle>Manage Recipients</DialogTitle>
          {entry && (
            <p className="text-sm text-muted-foreground mt-1">
              {entry.entry_name}
            </p>
          )}
        </DialogHeader>

        <ScrollArea className="max-h-[50vh] pr-4">
          <div className="space-y-4 py-4">
            {/* Add New Recipient */}
            {isAddingNew ? (
              <div className="p-4 border border-border rounded-lg space-y-4 bg-secondary/20">
                <div className="space-y-2">
                  <Label htmlFor="new-email">Email Address</Label>
                  <Input
                    id="new-email"
                    type="email"
                    placeholder="recipient@example.com"
                    value={newRecipientEmail}
                    onChange={(e) => setNewRecipientEmail(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") handleAddRecipient();
                      if (e.key === "Escape") setIsAddingNew(false);
                    }}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="new-role">Role</Label>
                  <Select value={newRecipientRole} onValueChange={(value) => setNewRecipientRole(value as "viewer" | "editor" | "owner")}>
                    <SelectTrigger id="new-role">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="viewer">Viewer</SelectItem>
                      <SelectItem value="editor">Editor</SelectItem>
                      <SelectItem value="owner">Owner</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex gap-2">
                  <Button
                    onClick={handleAddRecipient}
                    disabled={!newRecipientEmail.trim()}
                    className="bg-[#C9A44A] hover:bg-[#B8934A]"
                  >
                    Add Recipient
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => {
                      setIsAddingNew(false);
                      setNewRecipientEmail("");
                      setNewRecipientRole("viewer");
                    }}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            ) : (
              <Button
                variant="outline"
                className="w-full"
                onClick={() => setIsAddingNew(true)}
              >
                <UserPlus className="h-4 w-4 mr-2" />
                Add Recipient
              </Button>
            )}

            {/* Recipients List */}
            {recipients.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No recipients yet</p>
                <p className="text-sm mt-1">Add recipients to share this entry</p>
              </div>
            ) : (
              <div className="space-y-3">
                {recipients.map((recipient) => (
                  <div
                    key={recipient.id}
                    className="flex items-center gap-3 p-3 border border-border rounded-lg hover:bg-secondary/50 transition-colors"
                  >
                    <Avatar className="h-10 w-10">
                      <AvatarFallback className="bg-primary/10">
                        {recipient.name.slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>

                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">{recipient.name}</p>
                      <p className="text-sm text-muted-foreground truncate">
                        {recipient.email}
                      </p>
                    </div>

                    <Select
                      value={recipient.role}
                      onValueChange={(value) => handleRoleChange(recipient.id, value as "viewer" | "editor" | "owner")}
                    >
                      <SelectTrigger className="w-32">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="viewer">Viewer</SelectItem>
                        <SelectItem value="editor">Editor</SelectItem>
                        <SelectItem value="owner">Owner</SelectItem>
                      </SelectContent>
                    </Select>

                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleResendInvitation(recipient)}
                      title="Resend invitation"
                    >
                      <Mail className="h-4 w-4" />
                    </Button>

                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleRemoveRecipient(recipient.id)}
                      className="text-destructive hover:text-destructive hover:bg-destructive/10"
                      title="Remove recipient"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}

            <p className="text-xs text-muted-foreground mt-4">
              Recipients must have verified sovereign identity to access shared entries.
            </p>
          </div>
        </ScrollArea>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleSaveChanges}
            className="bg-[#C9A44A] hover:bg-[#B8934A]"
          >
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
