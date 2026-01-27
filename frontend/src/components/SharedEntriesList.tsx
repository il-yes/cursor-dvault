import { SharedEntry, ShareFilter } from "@/types/sharing";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CreditCard, FileText, Key, Shield, User, Users } from "lucide-react";
import { cn } from "@/lib/utils";
import "./contributionGraph/g-scrollbar.css";
import { useState } from "react";
import { NewShareModal } from "./NewCryptoShareModal";

interface SharedEntriesListProps {
  entries: SharedEntry[];
  selectedEntryId: string | null;
  onSelectEntry: (entry: SharedEntry) => void;
}

const getEntryIcon = (type: string) => {
  switch (type) {
    case "login":
      return Shield;
    case "card":
      return CreditCard;
    case "note":
      return FileText;
    case "identity":
      return User;
    case "sshkey":
      return Key;
    default:
      return Shield;
  }
};

const getStatusVariant = (status: string) => {
  switch (status) {
    case "active":
      return "default";
    case "pending":
      return "secondary";
    case "expired":
      return "outline";
    case "revoked":
      return "destructive";
    default:
      return "default";
  }
};

export function SharedEntriesList({
  entries,
  selectedEntryId,
  onSelectEntry,
}: SharedEntriesListProps) {
  const [isNewShareOpen, setIsNewShareOpen] = useState(false);
  const [sharedEntriesRefreshKey, setSharedEntriesRefreshKey] = useState(0);

  return (
    <>
      {/* Fixed Header */}
      <div className="sticky top-0 z-10 border-b border-border p-4 bg-background relative">
        <h2 className="text-lg font-semibold mb-3">Shared Entries</h2>
        <button
          onClick={() => setIsNewShareOpen(true)}
          style={{ position: "absolute", right: "5px", top: "40px" }}
          className="inline-flex items-center gap-2 rounded-full bg-gradient-to-r from-[#C9A44A] via-amber-400 to-[#B8934A] px-4 py-2 text-xs font-semibold text-black shadow-[0_14px_40px_rgba(0,0,0,0.8)] hover:shadow-[0_18px_60px_rgba(0,0,0,0.9)] hover:-translate-y-0.5 transition-all"
        >
          + Add Cryptographic Share
        </button>

        <div className="text-xs text-muted-foreground">
          {entries.length} {entries.length === 1 ? "entry" : "entries"}
        </div>
      </div>
      <NewShareModal
        open={isNewShareOpen}
        onOpenChange={setIsNewShareOpen}
        onShareSuccess={() => {
          setSharedEntriesRefreshKey(prev => prev + 1);
          window.dispatchEvent(new CustomEvent('shareEntriesRefresh'));
        }}
      />

      {/* Scrollable List */}
      <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar">
        <div className="p-2 space-y-1">
          {entries.map((entry) => {
            const Icon = getEntryIcon(entry.entry_type);
            const isSelected = selectedEntryId === entry.id;

            return (
              <button
                key={entry.id}
                onClick={() => onSelectEntry(entry)}
                className={cn(
                  "w-full text-left p-3 rounded-lg transition-all hover:bg-accent/50",
                  isSelected && "bg-accent border border-primary/20"
                )}
              >
                <div className="flex items-start gap-3">
                  <div className={cn(
                    "mt-1 p-2 rounded-md",
                    isSelected ? "bg-primary/10 text-primary" : "bg-muted"
                  )}>
                    <Icon className="h-4 w-4" />
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="font-medium text-sm truncate">
                        {entry.entry_name}
                      </h3>
                      <Badge
                        variant={getStatusVariant(entry.status)}
                        className="text-xs shrink-0"
                      >
                        {entry.status}
                      </Badge>
                    </div>

                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Users className="h-3 w-3" />
                      <span>
                        Shared with {entry?.recipients ? entry?.recipients.length : 0}{" "}
                        {entry?.recipients && entry?.recipients.length > 0 && entry?.recipients.length === 1 ? "recipient" : "recipients"}
                      </span>
                    </div>

                    {entry.expires_at && (
                      <div className="text-xs text-muted-foreground mt-1">
                        Expires: {new Date(entry.expires_at).toLocaleDateString()}
                      </div>
                    )}
                  </div>
                </div>
              </button>
            );
          })}

          {entries.length === 0 && (
            <div className="text-center py-12 text-muted-foreground">
              <Users className="h-12 w-12 mx-auto mb-3 opacity-20" />
              <p className="text-sm">No shared entries found</p>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
