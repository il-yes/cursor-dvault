import { useEffect, useState } from "react";
import { ENTRY_TYPE_CARD, ENTRY_TYPE_IDENTITY, ENTRY_TYPE_LOGIN, ENTRY_TYPE_NOTE, ENTRY_TYPE_SSHKEY, Folder, VaultEntry } from "@/types/vault";
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
import { useVaultStore } from "@/store/vaultStore";
import { devsecops_incident_v1 } from "@/data/mockTemplate";

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
  const [folderId, setFolderId] = useState<string | null>(null);
  const [folders, setFolders] = useState<Folder[]>([]);

  const vaultContext = useVaultStore((state) => state.vault);

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

  const [customFields, setCustomFields] = useState<Record<string, any>>({});

  useEffect(() => {
    vaultContext && setFolders(vaultContext.Vault.folders || []);
  }, [vaultContext]);

  useEffect(() => {
    // For testing
    // setCustomFields(devsecops_incident_v1);
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const base = {
      entry_name: entryName,
      type,
      custom_fields: customFields,
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
    entry.folder_id = folderId;
    console.log({ entry })

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
    setCustomFields({});

    onOpenChange(false);
  };

  const renderFields = () => {
    switch (type) {
      case ENTRY_TYPE_LOGIN:
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

      case ENTRY_TYPE_CARD:
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

      case ENTRY_TYPE_IDENTITY:
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

      case ENTRY_TYPE_NOTE:
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

      case ENTRY_TYPE_SSHKEY:
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

  // - Custom Fields
  const updateCustomFieldKey = (oldKey: string, newKey: string) => {
    setCustomFields(prev => {
      const next: Record<string, any> = {};
      for (const [k, v] of Object.entries(prev)) {
        if (k === oldKey) {
          if (newKey.trim()) next[newKey] = v;
        } else {
          next[k] = v;
        }
      }
      return next;
    });
  };

  const updateCustomFieldValue = (key: string, value: string) => {
    setCustomFields(prev => ({
      ...prev,
      [key]: value,
    }));
  };

  const addCustomField = () => {
    setCustomFields(prev => {
      let i = 1;
      let key = `new_field_${i}`;
      while (prev[key]) {
        i++;
        key = `new_field_${i}`;
      }
      return { ...prev, [key]: "" };
    });
  };

  const removeCustomField = (key: string) => {
    setCustomFields(prev => {
      const next = { ...prev };
      delete next[key];
      return next;
    });
  };
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
            <Label>Folder</Label>
            <Select value={folderId} onValueChange={setFolderId}>
              <SelectTrigger id="entry">
                <SelectValue placeholder="Choose an entry from your vault" />
              </SelectTrigger>
              <SelectContent>
                {folders && folders.map((folder) => (
                  <SelectItem key={folder.id} value={folder.id}>
                    <span className="capitalize">{folder.name}</span>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

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


        {/* Custom Fields */}
        <div className="backdrop-blur-xl bg-gradient-to-br from-white/40 to-zinc-50/20 dark:from-zinc-900/40 dark:to-black/20 rounded-3xl border border-white/30 dark:border-zinc-700/30 p-8">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-2xl font-bold bg-gradient-to-r from-primary to-amber-500/80 bg-clip-text text-transparent">
              Custom Fields
            </h3>
          </div>

          <div className="space-y-4">
            {customFields && Object.entries(customFields)
              .filter(([k]) => !RESERVED.has(k))
              .map(([key, value]) => (
                <div
                  key={key}
                  className="grid grid-cols-1 md:grid-cols-12 gap-3 items-start rounded-2xl border border-white30 dark:border-zinc-70030 bg-white40 dark:bg-zinc-90040 p-4"
                >
                  <Input
                    className="md:col-span-4"
                    value={key}
                    onChange={(e) => updateCustomFieldKey(key, e.target.value)}
                    placeholder="Field key"
                  />

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
        </div>

        {/* File Upload Widget 
          <div className="space-y-2">
            <Label>Attachments</Label>
             <FileUploadWidget
              onFileSelect={setAttachedFiles}
              value={attachedFiles}
              maxFiles={5}
            /> 
          </div> */}

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
    </Dialog >
  );
}
