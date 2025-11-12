import { SharedEntry, DetailView } from "@/types/sharing";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Users, FileText, Lock, Info, Calendar, Shield } from "lucide-react";
import { useState } from "react";

interface SharedEntryOverviewProps {
  entry: SharedEntry | null;
  onViewChange: (view: DetailView) => void;
}

export function SharedEntryOverview({ entry, onViewChange }: SharedEntryOverviewProps) {
  const [accessMode, setAccessMode] = useState<"read" | "edit">(entry?.access_mode || "read");

  if (!entry) {
    return (
      <div className="flex-1 flex items-center justify-center bg-gradient-to-b from-background to-secondary/20">
        <div className="text-center text-muted-foreground">
          <Users className="h-16 w-16 mx-auto mb-4 opacity-20" />
          <p className="text-lg">Select a shared entry</p>
          <p className="text-sm mt-2">Choose an entry to view its details</p>
        </div>
      </div>
    );
  }

  const handleAccessModeToggle = (checked: boolean) => {
    setAccessMode(checked ? "edit" : "read");
    // In real app, this would trigger an API call
  };

  return (
    <div className="flex-1 flex flex-col h-full bg-gradient-to-b from-background to-secondary/20 overflow-hidden">
      {/* Fixed Header */}
      <div className="sticky top-0 z-10 border-b border-border p-6 bg-background">
        <div className="space-y-2">
          <h2 className="text-2xl font-bold">{entry.entry_name}</h2>
          {entry.description && (
            <p className="text-sm text-muted-foreground">{entry.description}</p>
          )}
          {entry.folder && (
            <Badge variant="outline" className="text-xs">
              📁 {entry.folder}
            </Badge>
          )}
        </div>
      </div>

      {/* Scrollable Content */}
      <ScrollArea className="flex-1">
        <div className="p-6 space-y-6">
          {/* Access Mode */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="access-mode" className="text-base font-medium">
                  Access Mode
                </Label>
                <p className="text-sm text-muted-foreground">
                  {accessMode === "read" ? "Read-only access" : "Editable access"}
                </p>
              </div>
              <Switch
                id="access-mode"
                checked={accessMode === "edit"}
                onCheckedChange={handleAccessModeToggle}
              />
            </div>
          </div>

          <Separator />

          {/* Quick Info */}
          <div className="space-y-3">
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Info className="h-4 w-4" />
              Quick Information
            </h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-muted-foreground">Encryption</p>
                <Badge variant="outline" className="mt-1">
                  <Shield className="h-3 w-3 mr-1" />
                  {entry.encryption}
                </Badge>
              </div>
              <div>
                <p className="text-muted-foreground">Recipients</p>
                <p className="font-medium mt-1">{entry.recipients.length}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Shared Date</p>
                <p className="font-medium mt-1">
                  {new Date(entry.shared_at).toLocaleDateString()}
                </p>
              </div>
              {entry.expires_at && (
                <div>
                  <p className="text-muted-foreground">Expires</p>
                  <div className="flex items-center gap-1 mt-1">
                    <Calendar className="h-3 w-3 text-yellow-600" />
                    <p className="font-medium text-yellow-600">
                      {new Date(entry.expires_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              )}
            </div>
          </div>

          <Separator />

          {/* Action Links */}
          <div className="space-y-3">
            <h3 className="text-sm font-semibold">Details</h3>
            <div className="space-y-2">
              <Button
                variant="outline"
                className="w-full justify-start"
                onClick={() => onViewChange("recipients")}
              >
                <Users className="h-4 w-4 mr-2" />
                Recipients
                <Badge variant="secondary" className="ml-auto">
                  {entry.recipients.length}
                </Badge>
              </Button>

              <Button
                variant="outline"
                className="w-full justify-start"
                onClick={() => onViewChange("audit")}
              >
                <FileText className="h-4 w-4 mr-2" />
                Audit Log
                <Badge variant="secondary" className="ml-auto">
                  {entry.audit_log.length}
                </Badge>
              </Button>

              <Button
                variant="outline"
                className="w-full justify-start"
                onClick={() => onViewChange("encryption")}
              >
                <Lock className="h-4 w-4 mr-2" />
                Encryption & Permissions
              </Button>

              <Button
                variant="outline"
                className="w-full justify-start"
                onClick={() => onViewChange("metadata")}
              >
                <Info className="h-4 w-4 mr-2" />
                Metadata
              </Button>
            </div>
          </div>
        </div>
      </ScrollArea>
    </div>
  );
}
