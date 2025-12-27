import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { VaultEntry } from "@/types/vault";
import { RotateCcw, Trash2, Shield, Calendar, Type } from "lucide-react";
import { toast } from "@/hooks/use-toast";
import "./contributionGraph/g-scrollbar.css";

interface TrashDetailPanelProps {
  entry: VaultEntry | null;
  onRestore: (entry: VaultEntry) => void;
  onDeletePermanently: (entry: VaultEntry) => void;
}

export function TrashDetailPanel({ entry, onRestore, onDeletePermanently }: TrashDetailPanelProps) {
  const handleRestore = () => {
    if (entry) {
      onRestore(entry);
      toast({
        title: "Entry restored",
        description: `${entry.entry_name} has been restored successfully.`,
      });
    }
  };

  const handleDeletePermanently = () => {
    if (entry) {
      if (confirm(`Are you sure you want to permanently delete "${entry.entry_name}"? This action cannot be undone.`)) {
        onDeletePermanently(entry);
        toast({
          title: "Entry deleted permanently",
          description: "The entry has been removed from your vault.",
          variant: "destructive",
        });
      }
    }
  };

  if (!entry) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-center p-8 bg-gradient-to-b from-background to-secondary/20">
        <div className="relative">
          <Trash2 className="h-20 w-20 text-muted-foreground/20 mb-4" />
        </div>
        <h3 className="text-lg font-semibold mb-2">
          Select a trashed entry
        </h3>
        <p className="text-sm text-muted-foreground max-w-xs">
          Select an entry from the trash to restore it or delete it permanently.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-gradient-to-b from-background to-secondary/20">
      {/* Header */}
      <div className="sticky top-0 z-10 border-b border-border p-6 bg-background">
        <div className="flex items-start justify-between mb-4">
          <div className="space-y-2">
            <h2 className="text-2xl font-bold">{entry.entry_name}</h2>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="capitalize">
                {entry.type}
              </Badge>
              <Badge variant="outline" className="text-xs bg-destructive/10 text-destructive border-destructive/20">
                <Trash2 className="h-3 w-3 mr-1" />
                Trashed
              </Badge>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto scrollbar-glassmorphism thin-scrollbar p-6">
        <div className="max-w-2xl mx-auto space-y-6">
          {/* Entry Info */}
          <div className="bg-card border border-border rounded-lg p-6 space-y-4">
            <div className="flex items-center gap-3">
              <Type className="h-5 w-5 text-primary" />
              <div>
                <p className="text-sm text-muted-foreground">Entry Type</p>
                <p className="font-medium capitalize">{entry.type}</p>
              </div>
            </div>

            <Separator />

            <div className="flex items-center gap-3">
              <Calendar className="h-5 w-5 text-primary" />
              <div>
                <p className="text-sm text-muted-foreground">Date Deleted</p>
                <p className="font-medium">{new Date(entry.updated_at).toLocaleString()}</p>
              </div>
            </div>

            <Separator />

            <div className="flex items-center gap-3">
              <Shield className="h-5 w-5 text-primary" />
              <div>
                <p className="text-sm text-muted-foreground">Original Creation</p>
                <p className="font-medium">{new Date(entry.created_at).toLocaleString()}</p>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="space-y-3">
            <Button
              onClick={handleRestore}
              className="w-full bg-primary hover:bg-primary/90 text-primary-foreground h-12"
            >
              <RotateCcw className="h-4 w-4 mr-2" />
              Restore Entry
            </Button>

            <Button
              onClick={handleDeletePermanently}
              variant="outline"
              className="w-full border-destructive/50 text-destructive hover:bg-destructive/10 h-12"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete Permanently
            </Button>
          </div>

          {/* Warning Notice */}
          <div className="bg-destructive/5 border border-destructive/20 p-4 rounded-lg">
            <div className="flex items-start gap-2">
              <Trash2 className="h-4 w-4 text-destructive flex-shrink-0 mt-0.5" />
              <div className="text-xs text-foreground space-y-1">
                <p className="font-medium">Permanent deletion warning</p>
                <p className="text-muted-foreground">
                  Permanently deleted entries cannot be recovered. Make sure you want to proceed before confirming.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
