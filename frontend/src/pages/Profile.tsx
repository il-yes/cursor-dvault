import { DashboardLayout } from "@/components/DashboardLayout";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { User, Copy, RefreshCw, HardDrive, Clock, LogOut, Key } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { useNavigate } from "react-router-dom";
import { useVaultStore } from "@/store/vaultStore";
import { formatDistanceToNow } from "date-fns";
import { useAppStore } from "@/store/appStore";
import { formatMonthYear } from "@/services/utils";
import { useAuthStore } from "@/store/useAuthStore";
import { useVault } from "@/hooks/useVault";
import * as AppAPI from "../../wailsjs/go/main/App";


const Profile = () => {
  const { toast } = useToast();
  const navigate = useNavigate();
  const { vault, lastSyncTime, loadVault, clearVault: clearVaultStore } = useVaultStore();
  const { clearVault: clearVaultContext, vaultContext } = useVault();

  const totalEntries = Object.values(vault?.Vault?.entries || {}).flat().length;
  const maxEntries = vault?.vault_runtime_context?.AppSettings?.vault_settings?.max_entries || 1000;
  const usagePercent = Math.ceil((totalEntries / maxEntries) * 100);
  const lastSync = vault?.LastSynced ? formatDistanceToNow(new Date(vault.LastSynced), { addSuffix: true }) : 'Never';
  const session = useAppStore.getState().session;
  const user = session?.vault_runtime_context?.CurrentUser

  const handleCopyStellarAddress = (stellarAddress: string) => {
    navigator.clipboard.writeText(stellarAddress);
    toast({
      title: "Copied!",
      description: "Stellar address copied to clipboard",
    });
  };

  const handleRotateKey = () => {
    toast({
      title: "Key Rotation Initiated",
      description: "Your encryption key is being rotated. This may take a moment.",
    });
  };

  const handleSyncNow = () => {
    toast({
      title: "Syncing...",
      description: "Synchronizing your vault with the blockchain.",
    });
  };

  const handleGenerateApiKey = () => {
    toast({
      title: "API Key Generated",
      description: "Your new API key has been created.",
    });
  };

  const handleLogout = () => {
    // Clear all stores and state
    clearVaultStore();                   // Clear vault store (entries, vault data)
    clearVaultContext();                 // Clear vault context provider
    useAuthStore.getState().clearAll();  // Clear auth store (user, tokens)
    useAppStore.getState().reset();      // Clear app store (session)

    // Clear specific localStorage items (not all, to preserve settings)
    localStorage.removeItem('userId');
    localStorage.removeItem('vault-storage');

    toast({
      title: "Logged out",
      description: "You have been successfully logged out.",
    });

    navigate("/auth/signin");
  };

  const handleSync = async () => {
    const { jwtToken } = useAuthStore.getState();
    const { refreshVault } = useVault();

    try {
      const vaultPassword = "password"
      const cid = await AppAPI.SynchronizeVault(jwtToken, vaultPassword || "vaultPassword");
      console.log("Vault sync success:", cid);

      // 1. Reload full updated vault session from backend
      const updatedContext = await loadVault();

      // 2. Update the VaultProvider runtime with fresh session
      // 🔄 Reload fresh context from backend → pushes into Zustand → Provider picks it up
      await refreshVault();

      toast({
        title: "Success",
        description: "Vault synced successfully!",
      });
    } catch (err) {
      toast({
        title: "Error",
        description: "Failed to sync vault: " + (err instanceof Error ? err.message : "Unknown error"),
      });
    }
  };


  return (
    <DashboardLayout>
      <div className="h-full overflow-y-auto bg-gradient-to-b from-background to-secondary/20">
        <div className="max-w-4xl mx-auto p-6 space-y-6">
          {/* Header */}
          <div className="space-y-2">
            <h1 className="text-3xl font-bold">Profile</h1>
            <p className="text-muted-foreground">Manage your account and preferences</p>
          </div>

          {/* User Info Card */}
          <Card>
            <CardHeader>
              <CardTitle>User Information</CardTitle>
              <CardDescription>Your personal details</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center gap-4">
                <Avatar className="h-20 w-20">
                  <AvatarFallback className="bg-primary/10 text-2xl">
                    <User className="h-10 w-10 text-primary" />
                  </AvatarFallback>
                </Avatar>
                <div>
                  <Button variant="outline" size="sm">
                    Change Avatar
                  </Button>
                  <p className="text-xs text-muted-foreground mt-2">
                    JPG, PNG or GIF. Max size 2MB.
                  </p>
                </div>
              </div>

              <Separator />

              <div className="grid md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input id="name" defaultValue={session?.user?.username} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="last_name">Last Name</Label>
                  <Input id="last_name" defaultValue={session?.user?.last_name} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Email Address</Label>
                  <Input id="email" type="email" defaultValue={session?.user?.email} />
                </div>
              </div>

              <Button className="bg-[#C9A44A] hover:bg-[#B8934A]">
                Save Changes
              </Button>
            </CardContent>
          </Card>

          {/* Connection Card */}
          <Card>
            <CardHeader>
              <CardTitle>Blockchain Connection</CardTitle>
              <CardDescription>Your sovereign identity credentials</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Stellar Address</Label>
                <div className="flex gap-2">
                  <Input
                    value={user?.stellar_account?.public_key}
                    readOnly
                    className="font-mono bg-secondary"
                  />
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => handleCopyStellarAddress(user?.stellar_account?.public_key)}
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                </div>

                <div className="flex gap-2">
                  <Input
                    value={user?.stellar_account?.private_key}
                    readOnly
                    className="font-mono bg-secondary"
                  />
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => handleCopyStellarAddress(user?.stellar_account?.private_key)}
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  Your unique blockchain identifier
                </p>
              </div>

              <div className="flex items-center gap-2">
                <Badge variant="outline" className="bg-green-500/10 text-green-600 border-green-500/20">
                  ● Connected
                </Badge>
                <Button variant="outline" size="sm" onClick={handleRotateKey}>
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Rotate Key
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Device Info Card */}
          <Card>
            <CardHeader>
              <CardTitle>Device Information</CardTitle>
              <CardDescription>Local storage and sync status</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <HardDrive className="h-4 w-4" />
                    <span>Device ID</span>
                  </div>
                  <p className="font-mono text-sm">device-7a8b9c</p>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <HardDrive className="h-4 w-4" />
                    <span>Storage Usage</span>
                  </div>
                  <div className="space-y-1">
                    <p className="font-medium">{totalEntries} MB / 500 MB</p>
                    <div className="h-2 bg-secondary rounded-full overflow-hidden">
                      <div className={`h-full bg-[#C9A44A] w-${usagePercent}`} />
                    </div>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Clock className="h-4 w-4" />
                    <span>Last Sync</span>
                  </div>
                  <p className="font-medium">{lastSync}</p>
                </div>

                <div className="flex items-center">
                  <Button variant="outline" onClick={() => handleSync()}>
                    <RefreshCw className="h-4 w-4 mr-2" />
                    Sync Now
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Enterprise Features Card */}
          <Card>
            <CardHeader>
              <CardTitle>Enterprise Features</CardTitle>
              <CardDescription>API access and advanced options</CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="outline" onClick={handleGenerateApiKey}>
                <Key className="h-4 w-4 mr-2" />
                Generate API Key
              </Button>
              <p className="text-xs text-muted-foreground mt-2">
                API keys allow external applications to access your vault programmatically.
              </p>
            </CardContent>
          </Card>

          {/* Footer */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium">VaultCore v1.2.0</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Last updated: {formatMonthYear(vault.LastUpdated)}
                  </p>
                </div>
                <Button
                  variant="destructive"
                  onClick={handleLogout}
                >
                  <LogOut className="h-4 w-4 mr-2" />
                  Logout
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default Profile;
