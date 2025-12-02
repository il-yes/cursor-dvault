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

const SettingsBeta = () => {
  const { toast } = useToast();
  const [autoLockTimeout, setAutoLockTimeout] = useState<string>("5");
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
      <div className="h-full overflow-y-auto scrollbar-glassmorphism thin-scrollbar bg-gradient-to-br from-white/50 via-white/30 to-zinc-50/20 dark:from-zinc-900/50 dark:via-zinc-900/30 dark:to-black/20 backdrop-blur-xl">
        <div className="max-w-4xl mx-auto p-8 space-y-9">
          {/* Header */}
          <div className="space-y-3 text-center backdrop-blur-xl bg-white/40 dark:bg-zinc-900/40 rounded-3xl p-12 border border-white/30 dark:border-zinc-700/30 shadow-2xl">
            <h1 className="text-5xl font-black bg-gradient-to-r from-foreground via-primary to-[#C9A44A] bg-clip-text text-transparent drop-shadow-2xl mb-2">
              Settings
            </h1>
            <p className="text-lg text-muted-foreground/90 mx-auto max-w-xl leading-relaxed">
              Manage your vault preferences and security
            </p>
          </div>

          {/* Settings Panels */}
          <div className="space-y-8">
            {/* Security Settings */}
            <Card className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl border border-white/40 dark:border-zinc-700/40 shadow-2xl">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Shield className="h-6 w-6 text-primary" />
                  <CardTitle className="text-2xl font-bold">Security</CardTitle>
                </div>
                <CardDescription className="text-base text-muted-foreground/80 mt-2">Configure security and privacy options</CardDescription>
              </CardHeader>
              <CardContent className="space-y-7">
                {/* Auto-lock Timeout */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="auto-lock" className="text-lg font-semibold">Auto-lock Timeout</Label>
                    <p className="text-sm text-muted-foreground">Automatically lock vault after inactivity</p>
                  </div>
                  <Select value={autoLockTimeout} onValueChange={setAutoLockTimeout}>
                    <SelectTrigger className="w-32 rounded-xl bg-white/60 dark:bg-zinc-900/60 border border-primary/30 shadow">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent className="backdrop-blur-xl bg-white/80 dark:bg-zinc-900/80 border-white/30 dark:border-zinc-700/30 rounded-xl shadow-lg">
                      <SelectItem value="1">1 min</SelectItem>
                      <SelectItem value="5">5 min</SelectItem>
                      <SelectItem value="15">15 min</SelectItem>
                      <SelectItem value="30">30 min</SelectItem>
                      <SelectItem value="never">Never</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <Separator className="bg-white/30 dark:bg-zinc-700/30" />

                {/* Re-mask Delay */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="remask" className="text-lg font-semibold">Re-mask Delay</Label>
                    <p className="text-sm text-muted-foreground">Seconds before sensitive fields re-hide</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Input id="remask" type="number" min="1" max="60" value={remaskDelay} onChange={e => setRemaskDelay(e.target.value)} className="w-20 text-center bg-white/60 dark:bg-zinc-900/60 border border-primary/30 rounded-xl shadow" />
                    <span className="text-sm text-muted-foreground">sec</span>
                  </div>
                </div>
                <Separator className="bg-white/30 dark:bg-zinc-700/30" />

                {/* Two-Factor Authentication */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="2fa" className="text-lg font-semibold">Two-Factor Authentication</Label>
                    <p className="text-sm text-muted-foreground">Add an extra layer of security</p>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge variant="outline" className="bg-white/60 dark:bg-zinc-900/60 border-primary/30 text-xs">Coming Soon</Badge>
                    <Switch id="2fa" checked={twoFactorEnabled} onCheckedChange={setTwoFactorEnabled} disabled />
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Blockchain Sync */}
            <Card className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl border border-white/40 dark:border-zinc-700/40 shadow-2xl">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Cloud className="h-6 w-6 text-primary" />
                  <CardTitle className="text-2xl font-bold">Blockchain Sync</CardTitle>
                </div>
                <CardDescription className="text-base text-muted-foreground/80 mt-2">Manage decentralized storage and sync</CardDescription>
              </CardHeader>
              <CardContent className="space-y-7">
                {/* IPFS Pinning */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="ipfs" className="text-lg font-semibold">IPFS Pinning</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <Activity className={`h-3 w-3 ${ipfsPinning ? "text-green-600" : "text-muted-foreground"}`} />
                      <p className="text-sm text-muted-foreground">Status: {ipfsPinning ? "Active" : "Inactive"}</p>
                    </div>
                  </div>
                  <Switch id="ipfs" checked={ipfsPinning} onCheckedChange={setIpfsPinning} />
                </div>
                <Separator className="bg-white/30 dark:bg-zinc-700/30" />

                {/* Stellar Sync Frequency */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="stellar-sync" className="text-lg font-semibold">Stellar Sync Frequency</Label>
                    <p className="text-sm text-muted-foreground">How often to sync with Stellar blockchain</p>
                  </div>
                  <Select value={stellarSyncFrequency} onValueChange={setStellarSyncFrequency}>
                    <SelectTrigger className="w-32 rounded-xl bg-white/60 dark:bg-zinc-900/60 border border-primary/30 shadow">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent className="backdrop-blur-xl bg-white/80 dark:bg-zinc-900/80 border-white/30 dark:border-zinc-700/30 rounded-xl shadow-lg">
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
            <Card className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl border border-white/40 dark:border-zinc-700/40 shadow-2xl">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Palette className="h-6 w-6 text-primary" />
                  <CardTitle className="text-2xl font-bold">UI Preferences</CardTitle>
                </div>
                <CardDescription className="text-base text-muted-foreground/80 mt-2">Customize your vault appearance</CardDescription>
              </CardHeader>
              <CardContent className="space-y-7">
                {/* Theme */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="theme" className="text-lg font-semibold">Theme</Label>
                    <p className="text-sm text-muted-foreground">Choose your preferred color scheme</p>
                  </div>
                  <Select value={theme} onValueChange={setTheme}>
                    <SelectTrigger className="w-32 rounded-xl bg-white/60 dark:bg-zinc-900/60 border border-primary/30 shadow">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent className="backdrop-blur-xl bg-white/80 dark:bg-zinc-900/80 border-white/30 dark:border-zinc-700/30 rounded-xl shadow-lg">
                      <SelectItem value="light">Light</SelectItem>
                      <SelectItem value="dark">Dark</SelectItem>
                      <SelectItem value="system">System</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <Separator className="bg-white/30 dark:bg-zinc-700/30" />

                {/* Animations */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="animations" className="text-lg font-semibold">Animations</Label>
                    <p className="text-sm text-muted-foreground">Enable smooth transitions and effects</p>
                  </div>
                  <Switch id="animations" checked={animationsEnabled} onCheckedChange={setAnimationsEnabled} />
                </div>
              </CardContent>
            </Card>

            {/* Account Management */}
            <Card className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl border border-white/40 dark:border-zinc-700/40 shadow-2xl">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Lock className="h-6 w-6 text-primary" />
                  <CardTitle className="text-2xl font-bold">Account Management</CardTitle>
                </div>
                <CardDescription className="text-base text-muted-foreground/80 mt-2">Manage your vault data</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {/* Export Vault */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-lg font-semibold">Export Vault</Label>
                    <p className="text-sm text-muted-foreground">Download an encrypted backup</p>
                  </div>
                  <Button variant="outline" onClick={handleExportVault} className="rounded-2xl border-primary/30 shadow">
                    <Download className="h-5 w-5 mr-2" />
                    Export
                  </Button>
                </div>
                <Separator className="bg-white/30 dark:bg-zinc-700/30" />

                {/* Delete Vault */}
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-lg font-semibold">Delete Vault</Label>
                    <p className="text-sm text-muted-foreground">Permanently delete your vault and data</p>
                  </div>
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button variant="destructive" className="rounded-2xl shadow">
                        <Trash2 className="h-5 w-5 mr-2" />
                        Delete
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent className="backdrop-blur-xl bg-white/90 dark:bg-zinc-900/90 border border-destructive/30 rounded-2xl shadow-md">
                      <AlertDialogHeader>
                        <AlertDialogTitle className="text-xl font-bold text-destructive">
                          Are you absolutely sure?
                        </AlertDialogTitle>
                        <AlertDialogDescription className="text-base text-muted-foreground/80">
                          This action cannot be undone. This will permanently delete your vault and remove all your data from our servers and the blockchain.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleDeleteVault} className="bg-destructive hover:bg-destructive/90 rounded-2xl">
                          Delete Vault
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Save Button */}
          <div className="flex justify-end pt-7">
            <Button
              onClick={handleSaveSettings}
              className="h-14 px-10 text-lg font-bold bg-gradient-to-r from-[#C9A44A] to-[#B8934A] hover:from-[#C9A44A]/90 hover:to-[#B8934A]/90 shadow-2xl hover:shadow-[#C9A44A]/30 rounded-2xl group hover:scale-[1.03] transition-all"
            >
              Save All Settings
            </Button>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );

};

export default SettingsBeta;
