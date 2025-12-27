import { useState } from "react";
import { VaultEntry } from "@/types/vault";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { FileUploadWidget } from "@/components/FileUploadWidget";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import "./ContributionGraph/g-scrollbar.css";

interface CreateEntryDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (
    entry: Omit<VaultEntry, "id" | "created_at" | "updated_at">
  ) => void;
}

export function CreateEntryDialog({
  open,
  onOpenChange,
  onSubmit,
}: CreateEntryDialogProps) {
  const [entryName, setEntryName] = useState("");
  const [type, setType] = useState<VaultEntry["type"]>("login");

  // LOGIN
  const [loginUsername, setLoginUsername] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [loginSite, setLoginSite] = useState("");

  // CARD
  const [cardOwner, setCardOwner] = useState("");
  const [cardNumber, setCardNumber] = useState("");
  const [cardExp, setCardExp] = useState("");
  const [cardCVC, setCardCVC] = useState("");

  // IDENTITY
  const [identity, setIdentity] = useState<Record<string, string>>({});

  // NOTE
  const [noteText, setNoteText] = useState("");

  // SSH KEY
  const [sshPrivate, setSSHPrivate] = useState("");
  const [sshPublic, setSSHPublic] = useState("");
  const [attachedFiles, setAttachedFiles] = useState<File[]>([])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const base = {
      entry_name: entryName,
      type,
      trashed: false,
      is_draft: true,
      is_favorite: false,
      created_at: new Date(),
      updated_at: new Date(),
    };

    let entry: any;

    switch (type) {
      case "login":
        entry = {
          ...base,
          user_name: loginUsername,
          password: loginPassword,
          web_site: loginSite || undefined,
        };
        break;

      case "card":
        entry = {
          ...base,
          owner: cardOwner,
          number: cardNumber,
          expiration: cardExp,
          cvc: cardCVC,
        };
        break;

      case "identity":
        entry = {
          ...base,
          ...identity,
        };
        break;

      case "note":
        entry = {
          ...base,
          additionnal_note: noteText,
        };
        break;

      case "sshkey":
        entry = {
          ...base,
          private_key: sshPrivate,
          public_key: sshPublic,
          e_fingerprint: "",
        };
        break;
    }

    onSubmit(entry);

    // Reset
    setEntryName("");
    setType("login");
    setNoteText("");
    setLoginUsername("");
    setLoginPassword("");
    setLoginSite("");
    setCardOwner("");
    setCardNumber("");
    setCardExp("");
    setCardCVC("");
    setIdentity({});
    setSSHPrivate("");
    setSSHPublic("");
    setAttachedFiles([]);

    onOpenChange(false);
  };

  const renderFields = () => {
    switch (type) {
      case "login":
        return (
          <>
            <div className="space-y-2">
              <Label>Username</Label>
              <Input
                value={loginUsername}
                onChange={(e) => setLoginUsername(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>Password</Label>
              <Input
                value={loginPassword}
                onChange={(e) => setLoginPassword(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>Website</Label>
              <Input
                value={loginSite}
                onChange={(e) => setLoginSite(e.target.value)}
              />
            </div>
          </>
        );

      case "card":
        return (
          <>
            <div className="space-y-2">
              <Label>Owner</Label>
              <Input
                value={cardOwner}
                onChange={(e) => setCardOwner(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>Card Number</Label>
              <Input
                value={cardNumber}
                onChange={(e) => setCardNumber(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>Expiration</Label>
              <Input
                placeholder="MM/YY"
                value={cardExp}
                onChange={(e) => setCardExp(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>CVC</Label>
              <Input
                value={cardCVC}
                onChange={(e) => setCardCVC(e.target.value)}
              />
            </div>
          </>
        );

      case "identity":
        const identityFields = [
          "firstname",
          "second_firstname",
          "lastname",
          "username",
          "company",
          "genre",
          "social_security_number",
          "ID_number",
          "driver_license",
          "mail",
          "telephone",
          "address_one",
          "address_two",
          "address_three",
          "city",
          "state",
          "postal_code",
          "country",
        ];

        return (
          <div className="grid grid-cols-2 gap-4">
            {identityFields.map((f) => (
              <div className="space-y-2" key={f}>
                <Label className="capitalize">{f.replace(/_/g, " ")}</Label>
                <Input
                  value={identity[f] || ""}
                  onChange={(e) =>
                    setIdentity((prev) => ({
                      ...prev,
                      [f]: e.target.value,
                    }))
                  }
                />
              </div>
            ))}
          </div>
        );

      case "note":
        return (
          <div className="space-y-2">
            <Label>Secure Note</Label>
            <Textarea
              className="min-h-[120px]"
              value={noteText}
              onChange={(e) => setNoteText(e.target.value)}
            />
          </div>
        );

      case "sshkey":
        return (
          <>
            <div className="space-y-2">
              <Label>Private Key</Label>
              <Textarea
                className="min-h-[120px]"
                value={sshPrivate}
                onChange={(e) => setSSHPrivate(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label>Public Key</Label>
              <Textarea
                className="min-h-[120px]"
                value={sshPublic}
                onChange={(e) => setSSHPublic(e.target.value)}
              />
            </div>
          </>
        );
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] bg-card border-border">
        <DialogHeader>
          <DialogTitle>Create New Entry</DialogTitle>
          <DialogDescription>
            Add a new encrypted entry to your sovereign vault.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} style={{ padding: "30px" }} className="space-y-4 max-h-[70vh] overflow-y-auto scrollbar-glassmorphism thin-scrollbar pr-2">
          <div className="space-y-2">
            <Label>Entry Name</Label>
            <Input
              required
              value={entryName}
              onChange={(e) => setEntryName(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label>Entry Type</Label>
            <Select value={type} onValueChange={(v) => setType(v as any)}>
              <SelectTrigger>
                <SelectValue placeholder="Select type" />
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

          {renderFields()}


            {/* File Upload Widget */}
            <div className="space-y-2">
              <Label>Attachments</Label>
              <FileUploadWidget
                onFileSelect={setAttachedFiles}
                value={attachedFiles}
                maxFiles={5}
              />
            </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit">Create</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
