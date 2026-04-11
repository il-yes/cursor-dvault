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
import { SettingsState, Vault } from "@/types/vault";
import * as AppAPI from "../../wailsjs/go/main/App";
import { useAuthStore } from "@/store/useAuthStore";
import { useVaultStore } from "@/store/vaultStore";
import { EditConfig, GetConfig } from "@/services/api";
import { withAuth } from "@/hooks/withAuth";



const SettingsBeta = () => {
	const { toast } = useToast();
	const [ipfsPinning, setIpfsPinning] = useState(true);
	const { vault, loadVault, clearVault: clearVaultStore } = useVaultStore();
	const jwtToken = useAuthStore.getState().jwtToken;
	const { user } = useAuthStore();

	const defaultSettings: SettingsState = {
		security: {
			autoLockSeconds: 300,
			clearClipboardAfter: 60,
			twoFactorEnabled: false
		},
		sync: {
			stellarFrequency: "manual",
			ipfsPinning: false,
			syncIntervalSeconds: 60,
			maxRetries: 3
		},
		ui: {
			theme: "system",
			animationsEnabled: false
		},
		features: {
			tracecoreEnabled: false,
			cloudBackupEnabled: false,
			threatDetectionEnabled: false,
			browserExtensionEnabled: false,
			gitCLIEnabled: false
		},
		backup: {
			enabled: false,
			schedule: "daily",
			retentionDays: 30,
			encryption: false
		},
		device: {
			user_id: user.id,
			vault_name: "",
			device_id: "0",
			device_name: "",
			last_synced: 0,
		},
		subscription: {
			plan: "free",
			features: {
				tracecoreEnabled: false,
				cloudBackupEnabled: false,
				threatDetectionEnabled: false,
				browserExtensionEnabled: false,
				gitCLIEnabled: false

			},
			limits: {
				maxVaults: 1,
				maxUsers: 1,
				maxDevices: 1,
				maxShares: 1
			}
		},
		sharing: {
			allowExternalSharing: false,
			defaultExpiryHours: 60,
			requirePassword: false,
			maxSharesPerEntry: 3,
		},

		privacy: {
			telemetryEnabled: false,
			anonymousMode: false,
		},
		onboarding: {
			packs: [],
			use_cases: [],
			installed_templates: [],
			completed: false
		},
		storage: {
			mode: "local",
			localIPFS: {
				api_endpoint: "http://localhost:5001",
			},
			cloud: {
				base_url: "http://localhost:4001/api",
			},
		},
	};

	const [settings, setSettings] = useState<SettingsState>(defaultSettings);

	const handleSaveSettings = async (vault: Vault) => {
		try {
			const response = await EditConfig(user, vault, settings, jwtToken)

			console.log("EditConfig response", response)

			await fetchConfig(vault.name, jwtToken)

			toast({
				title: "Settings updated",
				description: "Your preferences have been saved successfully.",
			})

		} catch (err) {
			console.error("Save settings failed", err)

			toast({
				title: "Error",
				description: "Failed to save settings",
				variant: "destructive",
			})
		}
	}

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
	}

	useEffect(() => {
		if (!vault?.Vault?.name) return
		fetchConfig(vault.Vault.name, jwtToken)
	}, [vault])

	const fetchConfig = async (vaultName, jwtToken) => {
		try {

			const response = await withAuth((token) => {
				return GetConfig(vaultName, token)
			});

			console.log("fetchConfig response", response)

			setSettings(response)

		} catch (err) {
			console.error("fetchConfig failed", err)
		}
	}

	const syncMap = {
		auto: 60,
		hourly: 3600,
		daily: 86400,
		manual: 0,
	};

	const mapSync = (seconds: number) => {
		if (seconds === 0) return "manual"
		if (seconds < 60) return "auto"
		if (seconds < 300) return "medium"
		return "slow"
	}

	const mapSyncToSeconds = (sync: string) => {
		if (sync === "auto") return 0
		if (sync === "fast") return 30
		if (sync === "medium") return 180
		return 600
	}
	const updateSettings = (path: string, value: any) => {
		setSettings((prev) => {
			const keys = path.split(".");
			const newState = { ...prev };

			let obj: any = newState;

			for (let i = 0; i < keys.length - 1; i++) {
				obj[keys[i]] = { ...obj[keys[i]] };
				obj = obj[keys[i]];
			}

			obj[keys[keys.length - 1]] = value;

			return newState;
		});
	};

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
									<Select value={String((settings?.security?.autoLockSeconds ?? 300) / 60)}
										onValueChange={(v) =>
											updateSettings("security.autoLockSeconds", Number(v) * 60)
										}>
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
										<Input
											id="remask"
											type="number"
											min="1"
											max="60"
											value={settings?.security?.clearClipboardAfter}
											onChange={e => updateSettings("security.clearClipboardAfter", Number(e.target.value))}
											className="w-20 text-center bg-white/60 dark:bg-zinc-900/60 border border-primary/30 rounded-xl shadow" />
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
										<Switch id="2fa" checked={settings?.security?.twoFactorEnabled} onCheckedChange={() => updateSettings("security.twoFactorEnabled", !settings.security.twoFactorEnabled)} disabled />
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
									<Switch id="ipfs" checked={ipfsPinning} onCheckedChange={() => updateSettings("sync.ipfsPinning", !ipfsPinning)} />
								</div>
								<Separator className="bg-white/30 dark:bg-zinc-700/30" />

								{/* Stellar Sync Frequency */}
								<div className="flex items-center justify-between">
									<div>
										<Label htmlFor="stellar-sync" className="text-lg font-semibold">Stellar Sync Frequency</Label>
										<p className="text-sm text-muted-foreground">How often to sync with Stellar blockchain</p>
									</div>
									<Select value={settings?.sync?.stellarFrequency} onValueChange={(v) => updateSettings("sync.stellarFrequency", v)}>
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
									<Select value={settings?.ui?.theme} onValueChange={(v) => updateSettings("ui.theme", v)}>
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
									<Switch id="animations" checked={settings?.ui?.animationsEnabled} onCheckedChange={(v) => updateSettings("ui.animationsEnabled", v)} />
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

						{/* Storage Settings */}
						<Card className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl border border-white/40 dark:border-zinc-700/40 shadow-2xl">
							<CardHeader>
								<div className="flex items-center gap-2">
									<Cloud className="h-6 w-6 text-primary" />
									<CardTitle className="text-2xl font-bold">Storage</CardTitle>
								</div>
								<CardDescription className="text-base text-muted-foreground/80 mt-2">Configure storage and backup options</CardDescription>
							</CardHeader>
							<CardContent className="space-y-7">
								{/* Storage Mode */}
								<div className="flex items-center justify-between">
									<div>
										<Label htmlFor="storage-mode" className="text-lg font-semibold">Storage Mode</Label>
										<p className="text-sm text-muted-foreground">Choose where your vault is stored</p>
									</div>
									<Select value={settings?.storage?.mode} onValueChange={(value) => updateSettings("storage.mode", value)}>
										<SelectTrigger className="w-[200px]">
											<SelectValue placeholder="Select storage mode" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="local">Local</SelectItem>
											<SelectItem value="cloud">Cloud</SelectItem>
										</SelectContent>
									</Select>
								</div>
							</CardContent>
						</Card>
					</div>

					{/* Save Button */}
					<div className="flex justify-end pt-7">
						<Button
							onClick={() => handleSaveSettings(vault?.Vault)}
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
