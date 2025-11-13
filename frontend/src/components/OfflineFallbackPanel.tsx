import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { WifiOff, Save, RotateCw, Download } from "lucide-react";
import { localVaultService } from "@/services/localVault";
import { toast } from "@/hooks/use-toast";
import { VaultContext } from "@/types/vault";

interface OfflineFallbackPanelProps {
  vaultDraft: Partial<VaultContext>;
  onRetry: () => void;
  onSaveSuccess?: (draftId: string) => void;
}

export function OfflineFallbackPanel({ 
  vaultDraft, 
  onRetry,
  onSaveSuccess 
}: OfflineFallbackPanelProps) {
  const [isSaving, setIsSaving] = useState(false);
  const [isExporting, setIsExporting] = useState(false);

  const handleSaveLocally = async () => {
    setIsSaving(true);
    try {
      const draftId = await localVaultService.saveDraft(vaultDraft);
      toast({
        title: "Saved locally",
        description: "Your vault will sync when connection is restored.",
      });
      onSaveSuccess?.(draftId);
    } catch (error) {
      toast({
        title: "Save failed",
        description: error instanceof Error ? error.message : "Could not save locally.",
        variant: "destructive",
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleExport = async () => {
    setIsExporting(true);
    try {
      // Create a temporary draft for export
      const draftId = `export_${Date.now()}`;
      await localVaultService.saveDraft({
        ...vaultDraft,
        user_id: vaultDraft.user_id || 'unknown',
      } as VaultContext);
      
      const blob = await localVaultService.exportDraft(draftId);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `vault-backup-${Date.now()}.json`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      
      // Clean up temporary draft
      await localVaultService.deleteDraft(draftId);
      
      toast({
        title: "Backup exported",
        description: "Your encrypted vault backup has been downloaded.",
      });
    } catch (error) {
      toast({
        title: "Export failed",
        description: error instanceof Error ? error.message : "Could not export backup.",
        variant: "destructive",
      });
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <div className="space-y-6">
      <Alert className="border-[#C9A44A]/20 bg-[#C9A44A]/5">
        <WifiOff className="h-4 w-4 text-[#C9A44A]" />
        <AlertDescription>
          You are currently offline. Your vault data can be saved locally and synced when 
          connection is restored.
        </AlertDescription>
      </Alert>

      <div className="space-y-3">
        <h3 className="text-sm font-medium text-muted-foreground">
          Choose an action:
        </h3>

        <div className="grid gap-3">
          <Button
            onClick={handleSaveLocally}
            disabled={isSaving}
            className="justify-start h-auto py-4 px-4"
            variant="outline"
          >
            <Save className="mr-3 h-5 w-5 text-[#C9A44A]" />
            <div className="text-left">
              <div className="font-medium">Save locally</div>
              <div className="text-xs text-muted-foreground">
                Store encrypted vault on this device
              </div>
            </div>
          </Button>

          <Button
            onClick={onRetry}
            className="justify-start h-auto py-4 px-4"
            variant="outline"
          >
            <RotateCw className="mr-3 h-5 w-5 text-[#C9A44A]" />
            <div className="text-left">
              <div className="font-medium">Retry connection</div>
              <div className="text-xs text-muted-foreground">
                Attempt to connect to backend again
              </div>
            </div>
          </Button>

          <Button
            onClick={handleExport}
            disabled={isExporting}
            className="justify-start h-auto py-4 px-4"
            variant="outline"
          >
            <Download className="mr-3 h-5 w-5 text-[#C9A44A]" />
            <div className="text-left">
              <div className="font-medium">Export backup</div>
              <div className="text-xs text-muted-foreground">
                Download encrypted vault file
              </div>
            </div>
          </Button>
        </div>
      </div>

      <p className="text-xs text-muted-foreground">
        All data is encrypted before storage. Your vault will remain secure offline.
      </p>
    </div>
  );
}
