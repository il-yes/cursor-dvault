import { useEffect, useState } from "react";
import { DashboardLayout } from "@/components/DashboardLayout";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Shield, Cloud, Palette, Trash2, Download, Lock, Activity } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { useAppStore } from "@/store/appStore";

const Settings = () => {
  const { toast } = useToast();
  const [autoLockTimeout, setAutoLockTimeout] = useState();
  const [remaskDelay, setRemaskDelay] = useState("8");
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false);
  const [ipfsPinning, setIpfsPinning] = useState(true);
  const [stellarSyncFrequency, setStellarSyncFrequency] = useState("auto");
  const [theme, setTheme] = useState("system");
  const [animationsEnabled, setAnimationsEnabled] = useState(true);
  const session = useAppStore((s) => s.session);

  const handleSaveSettings = () => {
    toast({
      title: "Settings updated",
      description: "Your preferences have been saved successfully.",
    });
  };

  const handleExportVault = () => {
    toast({
      title: "Export initiated",
      description: "Your vault is being prepared for export...",
    });
  };

  const handleDeleteVault = () => {
    toast({
      title: "Vault deletion",
      description: "This action will permanently delete your vault.",
      variant: "destructive",
    });
  };

  useEffect(() => {
    // Load settings from local storage or API on component mount
    setAutoLockTimeout(session?.AppSettings?.auto_lock_timeout || "5");
  }, []);

  return (
    <DashboardLayout>
      <div className="h-full overflow-y-auto bg-gradient-to-b from-background to-secondary/20">
        <div className="max-w-4xl mx-auto p-6 space-y-6">
          {/* Header */}
          <div className="space-y-2">
            <h1 className="text-3xl font-bold">Settings</h1>
            <p className="text-muted-foreground">Manage your vault preferences and security</p>
          </div>

          {/* Security Settings */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Shield className="h-5 w-5 text-primary" />
                <CardTitle>Security</CardTitle>
              </div>
              <CardDescription>Configure security and privacy options</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label htmlFor="auto-lock">Auto-lock Timeout</Label>
                    <p className="text-sm text-muted-foreground">
                      Automatically lock vault after period of inactivity
                    </p>
                  </div>
                  <Select value={autoLockTimeout} onValueChange={() => setAutoLockTimeout}>
                    <SelectTrigger className="w-32">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="1">1 min</SelectItem>
                      <SelectItem value="5">5 min</SelectItem>
                      <SelectItem value="15">15 min</SelectItem>
                      <SelectItem value="30">30 min</SelectItem>
                      <SelectItem value="never">Never</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <Separator />

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5 flex-1">
                    <Label htmlFor="remask">Re-mask Delay</Label>
                    <p className="text-sm text-muted-foreground">
                      Seconds before sensitive fields are hidden again
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Input
                      id="remask"
                      type="number"
                      min="1"
                      max="60"
                      value={remaskDelay}
                      onChange={(e) => setRemaskDelay(e.target.value)}
                      className="w-20 text-center"
                    />
                    <span className="text-sm text-muted-foreground">sec</span>
                  </div>
                </div>

                <Separator />

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label htmlFor="2fa">Two-Factor Authentication</Label>
                    <p className="text-sm text-muted-foreground">
                      Add an extra layer of security
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant="outline" className="text-xs">
                      Coming Soon
                    </Badge>
                    <Switch
                      id="2fa"
                      checked={twoFactorEnabled}
                      onCheckedChange={setTwoFactorEnabled}
                      disabled
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Blockchain Sync */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Cloud className="h-5 w-5 text-primary" />
                <CardTitle>Blockchain Sync</CardTitle>
              </div>
              <CardDescription>Manage decentralized storage and synchronization</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="ipfs">IPFS Pinning</Label>
                  <div className="flex items-center gap-2 mt-1">
                    <Activity className="h-3 w-3 text-green-600" />
                    <p className="text-sm text-muted-foreground">
                      Status: {ipfsPinning ? "Active" : "Inactive"}
                    </p>
                  </div>
                </div>
                <Switch
                  id="ipfs"
                  checked={ipfsPinning}
                  onCheckedChange={setIpfsPinning}
                />
              </div>

              <Separator />

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="stellar-sync">Stellar Sync Frequency</Label>
                  <p className="text-sm text-muted-foreground">
                    How often to sync with Stellar blockchain
                  </p>
                </div>
                <Select value={stellarSyncFrequency} onValueChange={setStellarSyncFrequency}>
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="auto">Auto</SelectItem>
                    <SelectItem value="hourly">Hourly</SelectItem>
                    <SelectItem value="daily">Daily</SelectItem>
                    <SelectItem value="manual">Manual</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>

          {/* UI Preferences */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Palette className="h-5 w-5 text-primary" />
                <CardTitle>UI Preferences</CardTitle>
              </div>
              <CardDescription>Customize your vault appearance</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="theme">Theme</Label>
                  <p className="text-sm text-muted-foreground">
                    Choose your preferred color scheme
                  </p>
                </div>
                <Select value={theme} onValueChange={setTheme}>
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="light">Light</SelectItem>
                    <SelectItem value="dark">Dark</SelectItem>
                    <SelectItem value="system">System</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <Separator />

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="animations">Animations</Label>
                  <p className="text-sm text-muted-foreground">
                    Enable smooth transitions and effects
                  </p>
                </div>
                <Switch
                  id="animations"
                  checked={animationsEnabled}
                  onCheckedChange={setAnimationsEnabled}
                />
              </div>
            </CardContent>
          </Card>

          {/* Account Management */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Lock className="h-5 w-5 text-primary" />
                <CardTitle>Account Management</CardTitle>
              </div>
              <CardDescription>Manage your vault data</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Export Vault</Label>
                  <p className="text-sm text-muted-foreground">
                    Download an encrypted backup of your vault
                  </p>
                </div>
                <Button variant="outline" onClick={handleExportVault}>
                  <Download className="h-4 w-4 mr-2" />
                  Export
                </Button>
              </div>

              <Separator />

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Delete Vault</Label>
                  <p className="text-sm text-muted-foreground">
                    Permanently delete your vault and all data
                  </p>
                </div>
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button variant="destructive">
                      <Trash2 className="h-4 w-4 mr-2" />
                      Delete
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                      <AlertDialogDescription>
                        This action cannot be undone. This will permanently delete your vault
                        and remove all your data from our servers and the blockchain.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        onClick={handleDeleteVault}
                        className="bg-destructive hover:bg-destructive/90"
                      >
                        Delete Vault
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </CardContent>
          </Card>

          {/* Save Button */}
          <div className="flex justify-end">
            <Button
              onClick={handleSaveSettings}
              className="bg-[#C9A44A] hover:bg-[#B8934A]"
            >
              Save All Settings
            </Button>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default Settings;
