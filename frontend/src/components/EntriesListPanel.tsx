import { VaultEntry } from "@/types/vault";
import { LogIn, CreditCard, FileText, Key, UserCircle, Star, Edit, Trash2, Share2, Folder } from "lucide-react";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { useVault } from "@/hooks/useVault";
import "./contributionGraph/g-scrollbar.css";

interface EntriesListPanelProps {
  entries: VaultEntry[];
  selectedEntryId: string | null;
  onSelectEntry: (entry: VaultEntry) => void;
  onToggleFavorite: (entryId: string) => void;
}

const entryTypeIcons = {
  login: LogIn,
  card: CreditCard,
  note: FileText,
  sshkey: Key,
  identity: UserCircle,
};

export function EntriesListPanel({
  entries,
  selectedEntryId,
  onSelectEntry,
  onToggleFavorite,
}: EntriesListPanelProps) {
  const { vaultContext } = useVault();
  
  const getFolderName = (folderId?: string) => {
    if (!folderId || !vaultContext?.Vault.folders) return null;
    return vaultContext.Vault.folders.find(f => f.id === folderId)?.name;
  };

  return (
    <div className="flex flex-col h-full">
      {/* Entries List */}
      <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar">
        {entries.length === 0 ? (
          <div className="flex items-center justify-center h-full text-muted-foreground p-8">
            <div className="text-center">
              <FileText className="h-12 w-12 mx-auto mb-3 opacity-20" />
              <p className="font-medium">No entries found</p>
              <p className="text-xs mt-1">Try adjusting your filters or search</p>
            </div>
          </div>
        ) : (
          <div className="p-2 space-y-1">
            {entries.map((entry, index) => {
              const Icon = entryTypeIcons[entry.type as keyof typeof entryTypeIcons] || FileText;
              const isSelected = entry.id === selectedEntryId;
              const folderName = getFolderName(entry.folder_id);
              const isTrashed = entry.trashed;
              
              // Get display info based on entry type
              const getSecondaryInfo = () => {
                if (entry.type === 'login') return (entry as any).user_name || (entry as any).web_site;
                if (entry.type === 'card') return (entry as any).owner;
                if (entry.type === 'identity') return (entry as any).mail;
                return entry.type;
              };

              return (
                <div
                  key={entry.id}
                  className={cn(
                    "group px-3 py-3 rounded-lg cursor-pointer transition-all duration-300 gold-border-glow relative",
                    "hover:bg-primary/5 hover:translate-x-1",
                    isSelected && "bg-primary/10 shadow-soft border-primary/30"
                  )}
                  style={{
                    animationDelay: `${index * 50}ms`,
                  }}
                >
                  <div onClick={() => onSelectEntry(entry)} className="flex items-start gap-3">
                    <div className={cn(
                      "p-2 rounded-md transition-colors",
                      isSelected ? "bg-primary/20 text-primary" : "bg-secondary text-muted-foreground"
                    )}>
                      <Icon className="h-4 w-4" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h3 className={cn(
                          "font-medium text-sm truncate transition-colors",
                          isSelected && "text-primary"
                        )}>
                          {entry.entry_name}
                        </h3>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            onToggleFavorite(entry.id);
                          }}
                          className="ml-auto flex-shrink-0 hover:scale-110 transition-transform"
                        >
                          <Star
                            className={cn(
                              "h-4 w-4 transition-all duration-300",
                              entry.is_favorite
                                ? "fill-primary text-primary drop-shadow-[0_0_6px_hsl(43_65%_54%/0.5)]"
                                : "text-muted-foreground hover:text-primary"
                            )}
                          />
                        </button>
                      </div>
                      <div className="flex items-center gap-2 mt-1">
                        <p className="text-xs text-muted-foreground truncate">
                          {getSecondaryInfo()}
                        </p>
                        {folderName && (
                          <Badge variant="outline" className="text-[10px] py-0 h-4 flex items-center gap-1">
                            <Folder className="h-2.5 w-2.5" />
                            {folderName}
                          </Badge>
                        )}
                        {isTrashed && (
                          <Badge variant="outline" className="text-[10px] py-0 h-4 bg-destructive/10 text-destructive border-destructive/20">
                            Trashed
                          </Badge>
                        )}
                      </div>
                      <p className="text-[10px] text-muted-foreground/70 mt-1">
                        Updated {new Date(entry.updated_at).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}

