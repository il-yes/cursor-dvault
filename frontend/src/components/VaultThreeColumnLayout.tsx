import { useState, useMemo, useEffect } from "react";
import { useVault } from "@/hooks/useVault";
import { VaultEntry } from "@/types/vault";
import { useSearchParams } from "react-router-dom";
import { EntriesListPanel } from "@/components/EntriesListPanel";
import { EntryDetailPanel } from "@/components/EntryDetailPanel";
import { TrashDetailPanel } from "@/components/TrashDetailPanel";
import { CreateEntryDialog } from "@/components/CreateEntryDialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus, Search, Shield, Radio, Database, CheckCircle2, Clock, Star } from "lucide-react";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import { useIsMobile } from "@/hooks/use-mobile";
import { Badge } from "@/components/ui/badge";

interface VaultThreeColumnLayoutProps {
  filter: string;
}

export function VaultThreeColumnLayout({ filter }: VaultThreeColumnLayoutProps) {
  const { vaultContext, addEntry, updateEntry, deleteEntry, toggleFavorite } = useVault();
  const [selectedEntry, setSelectedEntry] = useState<VaultEntry | null>(null);
  const [searchParams] = useSearchParams();
  const [editMode, setEditMode] = useState(false);
  const isMobile = useIsMobile();
  
  const allEntries = useMemo(() => {
    if (!vaultContext?.Vault) return [];

    const entries: VaultEntry[] = [
      ...(vaultContext.Vault.entries?.login || []),
      ...(vaultContext.Vault.entries?.card || []),
      ...(vaultContext.Vault.entries?.note || []),
      ...(vaultContext.Vault.entries?.sshkey || []),
      ...(vaultContext.Vault.entries?.identity || []),
    ];

    return entries;
  }, [vaultContext]);
  
  // Get entry ID from URL params
  const entryIdFromUrl = searchParams.get("entry");
  
  useEffect(() => {
    if (entryIdFromUrl && allEntries.length > 0) {
      const entry = allEntries.find(e => e.id === entryIdFromUrl);
      if (entry) {
        setSelectedEntry(entry);
      }
    }
  }, [entryIdFromUrl, allEntries]);

  const filteredEntries = useMemo(() => {
    let filtered = [...allEntries];

    // Apply category/folder filter
    if (filter === "favorites") {
      filtered = filtered.filter(e => e.is_favorite && !e.trashed);
    } else if (filter === "trash") {
      filtered = filtered.filter(e => e.trashed === true);
    } else if (filter.startsWith("folder:")) {
      const folderId = filter.replace("folder:", "");
      filtered = filtered.filter(e => e.folder_id === folderId && !e.trashed);
    } else if (filter !== "all") {
      filtered = filtered.filter(e => e.type === filter && !e.trashed);
    } else {
      // "all" filter - exclude trashed entries
      filtered = filtered.filter(e => !e.trashed);
    }

    return filtered;
  }, [allEntries, filter]);

  const handleToggleFavorite = (entryId: string) => {
    toggleFavorite(entryId);
  };

  const handleCreateEntry = (entry: Omit<VaultEntry, "id" | "created_at" | "updated_at">) => {
    const newEntry: VaultEntry = {
      ...entry,
      id: `${entry.type}-${Date.now()}`,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      trashed: false,
      is_draft: false,
    } as VaultEntry;
    addEntry(newEntry);
  };

  const handleEditEntry = (updates: Partial<VaultEntry>) => {
    if (selectedEntry) {
      updateEntry(selectedEntry.id, updates);
      setSelectedEntry({ ...selectedEntry, ...updates } as VaultEntry);
      setEditMode(false);
    }
  };

  const handleDeleteEntry = (entryId: string) => {
    deleteEntry(entryId);
    if (selectedEntry?.id === entryId) {
      setSelectedEntry(null);
    }
  };

  const handleRestoreEntry = (entryId: string) => {
    updateEntry(entryId, { trashed: false });
    if (selectedEntry?.id === entryId) {
      setSelectedEntry(null);
    }
  };

  const handleDeletePermanently = (entryId: string) => {
    // For now, just mark as deleted (in real app, would actually remove from state)
    updateEntry(entryId, { trashed: true });
    if (selectedEntry?.id === entryId) {
      setSelectedEntry(null);
    }
  };

  const syncStatus = vaultContext?.Dirty ? "unsynced" : "synced";
  
  const metrics = useMemo(() => {
    const total = filteredEntries.length;
    const synced = filteredEntries.filter(e => !vaultContext?.Dirty).length;
    const unsynced = total - synced;
    const favorites = filteredEntries.filter(e => e.is_favorite).length;
    
    return { total, synced, unsynced, favorites };
  }, [filteredEntries, vaultContext]);

  return (
    <div className="flex h-full">
      {/* Entries List - Center Column */}
      <div className="w-full md:w-80 lg:w-96 flex flex-col border-r border-border bg-secondary/30 overflow-hidden">
        {/* Header with Status */}
        <div className="sticky top-0 z-10 border-b border-border p-4 bg-background">
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-lg font-semibold flex items-center gap-2">
              {filter === "all" && "All Items"}
              {filter === "favorites" && "Favorites"}
              {filter === "trash" && "Trash"}
              {filter === "login" && "Identifiers"}
              {filter === "card" && "Payment Cards"}
              {filter === "identity" && "Identities"}
              {filter === "note" && "Secure Notes"}
              {filter === "sshkey" && "SSH Keys"}
              {filter.startsWith("folder:") && 
                vaultContext?.Vault.folders?.find(f => f.id === filter.replace("folder:", ""))?.name || "Folder"
              }
            </h2>
            
            {/* Sync Status Badge */}
            <Badge 
              variant={syncStatus === "synced" ? "default" : "secondary"}
              className={`flex items-center gap-1 ${
                syncStatus === "synced" 
                  ? "bg-primary/10 text-primary border-primary/20" 
                  : "animate-pulse-glow bg-yellow-500/10 text-yellow-600 border-yellow-500/20"
              }`}
            >
              <Radio className="h-3 w-3" />
              {syncStatus === "synced" ? "Synced" : "Pending"}
            </Badge>
          </div>

          {/* Metrics Row */}
          <div className="grid grid-cols-4 gap-2 mt-3">
            <Badge variant="outline" className="justify-center py-1.5">
              <Database className="h-3 w-3 mr-1 text-primary" />
              <span className="text-xs font-medium">{metrics.total}</span>
            </Badge>
            <Badge variant="outline" className="justify-center py-1.5">
              <CheckCircle2 className="h-3 w-3 mr-1 text-green-600" />
              <span className="text-xs font-medium">{metrics.synced}</span>
            </Badge>
            <Badge variant="outline" className="justify-center py-1.5">
              <Clock className="h-3 w-3 mr-1 text-yellow-600" />
              <span className="text-xs font-medium">{metrics.unsynced}</span>
            </Badge>
            <Badge variant="outline" className="justify-center py-1.5">
              <Star className="h-3 w-3 mr-1 text-primary fill-primary" />
              <span className="text-xs font-medium">{metrics.favorites}</span>
            </Badge>
          </div>
        </div>

        <EntriesListPanel
          entries={filteredEntries}
          selectedEntryId={selectedEntry?.id || null}
          onSelectEntry={setSelectedEntry}
          onToggleFavorite={handleToggleFavorite}
        />

      </div>

      {/* Detail Panel - Right Column (Desktop) / Drawer (Mobile) */}
      {isMobile ? (
        <Sheet open={!!selectedEntry} onOpenChange={(open) => !open && setSelectedEntry(null)}>
          <SheetContent side="right" className="w-full p-0">
            <div className="h-full overflow-y-auto">
              {filter === "trash" ? (
                <TrashDetailPanel
                  entry={selectedEntry}
                  onRestore={handleRestoreEntry}
                  onDeletePermanently={handleDeletePermanently}
                />
              ) : (
                <EntryDetailPanel 
                  entry={selectedEntry}
                  editMode={editMode}
                  onEdit={() => setEditMode(true)}
                  onSave={handleEditEntry}
                  onCancel={() => setEditMode(false)}
                  onDelete={() => selectedEntry && handleDeleteEntry(selectedEntry.id)}
                />
              )}
            </div>
          </SheetContent>
        </Sheet>
      ) : (
        <div className="flex-1 hidden md:flex flex-col overflow-hidden">
          <div className="flex-1 overflow-y-auto">
            {filter === "trash" ? (
              <TrashDetailPanel
                entry={selectedEntry}
                onRestore={handleRestoreEntry}
                onDeletePermanently={handleDeletePermanently}
              />
            ) : (
              <EntryDetailPanel 
                entry={selectedEntry}
                editMode={editMode}
                onEdit={() => setEditMode(true)}
                onSave={handleEditEntry}
                onCancel={() => setEditMode(false)}
                onDelete={() => selectedEntry && handleDeleteEntry(selectedEntry.id)}
              />
            )}
          </div>
        </div>
      )}
    </div>
  );
}
