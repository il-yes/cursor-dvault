import { DashboardLayout } from "@/components/DashboardLayout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { User, Copy, RefreshCw, HardDrive, Clock, LogOut, Key, Upload } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { useNavigate } from "react-router-dom";
import { useVaultStore } from "@/store/vaultStore";
import { formatDistanceToNow } from "date-fns";
import { useAppStore } from "@/store/appStore";
import { formatMonthYear } from "@/services/utils";
import { useAuthStore } from "@/store/useAuthStore";
import { useVault } from "@/hooks/useVault";
import * as AppAPI from "../../wailsjs/go/main/App";
import { useState, useRef, useEffect } from "react";
import { Progress } from "@radix-ui/react-progress";
import { GlassProgressBar } from "@/components/GlassProgressBar";
import "../components/contributionGraph/g-scrollbar.css";
import EncryptionVerificationModal from "@/components/EncryptionVerificationModal";
import EncryptionVerificationModalBeta from "@/components/EncryptionVerificationModalBeta";
import { EditUserInfos, GenerateApiKey } from "@/services/api";


const mockFile = {
	name: "contracts.zip",
	size: 1420000,
	commitHash: "d1b4c0ffee-example-stellar-tx-hash",
};

const Profile = () => {
	const { toast } = useToast();
	const navigate = useNavigate();
	const { vault, lastSyncTime, loadVault, clearVault: clearVaultStore } = useVaultStore();
	const { clearVault: clearVaultContext, vaultContext } = useVault();
	const [avatarUrl, setAvatarUrl] = useState<string | null>(null);
	const fileInputRef = useRef<HTMLInputElement>(null);
	const vaultPassword = "vaultPassword";
	const totalEntries = Object.values(vault?.Vault?.entries || {}).flat().length;
	const maxEntries = vault?.vault_runtime_context?.AppSettings?.vault_settings?.max_entries || 1000;
	const usagePercent = Math.ceil((totalEntries / maxEntries) * 100);
	const lastSync = vault?.LastSynced ? formatDistanceToNow(new Date(vault.LastSynced), { addSuffix: true }) : 'Never';
	const [userName, setUserName] = useState();
	const [firstName, setFirstName] = useState();
	const [lastName, setLastName] = useState();
	const [publicKey, setPublicKey] = useState("");
	const [privateKey, setPrivateKey] = useState("");
	
	const [progressVisible, setProgressVisible] = useState(false);
	const [showModal, setShowModal] = useState(false);

	const { encryptFile, uploadToIPFS, createStellarCommit, syncVault, refreshVault } = useVault();

	const [progress, setProgress] = useState(0);
	const [stage, setStage] = useState('encrypting'); // encrypting | uploading | complete

	const session = useAppStore.getState().session;
	const user = session?.user	
	user && console.log(user)

	useEffect(() => {
		const callback = (payload: { percent: number, stage: string } | number) => {
			if (typeof payload === 'object' && payload.percent !== undefined) {
				setProgress(payload.percent);
				setStage(payload.stage);
			} else {
				setProgress(payload as number);
			}
		};

		window.runtime?.EventsOn('progress-update', callback);
		return () => window.runtime?.EventsOff('progress-update');
	}, []);

	useEffect(() => {
		setUserName(user?.user_name);
		setFirstName(user?.first_name);
		setLastName(user?.last_name);
	}, [user]);

	useEffect(() => {
		const stellar = vault?.vault_runtime_context?.UserConfig?.stellar_account
		setPublicKey(stellar?.public_key);
		setPrivateKey(stellar?.private_key);
	}, [vault]);


	// Updated handleAvatarUpload
	const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
		setProgressVisible(true);
		const file = e.target.files?.[0];
		if (file) {
			if (file.size > 2 * 1024 * 1024) {
				toast({
					title: "File too large",
					description: "Please select an image smaller than 2MB.",
					variant: "destructive",
				});
				return;
			}
			const { jwtToken } = useAuthStore.getState();
			setProgress(0); // Start at 0

			try {
				const reader = new FileReader();
				reader.onload = (e) => {
					setAvatarUrl(e.target?.result as string);
				};
				reader.readAsDataURL(file);

				setStage('encrypting...');
				// Pass file buffer/path to backend (adjust encrypt to accept File or ArrayBuffer)
				const filePath = await readFileAsBuffer(file); // Helper to get buffer
				const encryptedData = await encryptFile(jwtToken, filePath, vaultPassword); // Now async with progress events

				setStage('uploading...');
				const cid = await uploadToIPFS(jwtToken, encryptedData); // Progress events update UI
				console.log({ cid });
				setStage('committing...');
				const stellarOp = await createStellarCommit(jwtToken, cid); // Final progress to 100
				console.log({ stellarOp });

				toast({
					title: "Avatar updated",
					description: "Your profile picture has been changed.",
				});
				setProgress(100);
				setStage('complete');
				setTimeout(() => setProgressVisible(false), 2000);
			} catch (error) {
				setProgressVisible(false);
				toast({
					title: "Upload failed",
					description: "Please try again.",
					variant: "destructive",
				});
				setProgress(0);
			}
		}
	};

	// Helper to read file as buffer for Go
	const readFileAsBuffer = (file: File): Promise<Uint8Array> => {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = () => resolve(new Uint8Array(reader.result as ArrayBuffer));
			reader.onerror = reject;
			reader.readAsArrayBuffer(file);
		});
	};

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

	const handleGenerateApiKey = async () => {
		const { jwtToken } = useAuthStore.getState();
		const payload = {
			password: "password",
			jwtToken
		}
		const keypair = await GenerateApiKey(payload);
		console.log(keypair);
		setPublicKey(keypair.public_key);
		setPrivateKey(keypair.private_key);

		refreshVault();
		const vault = useVaultStore.getState().vault;
		console.log('========= Post response api generated ==============')
		console.log(vault.vault_runtime_context)
		console.log('=========. ==============')
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
		AppAPI.SignOut(useAuthStore.getState().user?.id);
		navigate("/login/email");
	};

	const handleSyncBeta = async () => {
		setProgressVisible(true);
		const vaultPassword = "password";
		const { jwtToken } = useAuthStore.getState();
		setProgress(0);
		setStage("starting...");

		try {
			toast({
				title: "Syncing...",
				description: "Synchronizing your vault with the blockchain.",
			});
			setStage("syncing vault...");

			// Call new sync method that internally runs all steps with progress emits
			const response = await syncVault(jwtToken, vaultPassword);
			console.log("Vault sync success:", response);
			setProgress(100);
			setStage("complete");

			toast({ title: "Success", description: "Vault synced successfully!" });

			// Refresh vault state on frontend after sync finishes
			await loadVault();
			await refreshVault();
			// 3. Give Zustand time to propagate (usually instant)
			await new Promise(resolve => setTimeout(resolve, 100));

			// TODO: fit the reesponse backend to fill the verif object
			// setVerif(stellarOp);

			// After some delay, hide the progress bar like a toast
			setTimeout(() => setProgressVisible(false), 2000);

			

		} catch (error) {
			setProgressVisible(false);

			toast({ title: "Error", description: `Failed to sync vault: ${(error as Error).message}` });
			setProgress(0);
			setStage("error");
		}
	};

	const handleEditUserInfos = async () => {
		const { jwtToken } = useAuthStore.getState();
		const payload = {
			user_name: "",
			last_name: "",
			first_name: "",
		}
		const response = await EditUserInfos(jwtToken, payload);
		console.log(response);
		toast({ title: "Success", description: "User info updated successfully!" });
	}

	return (
		<DashboardLayout>
			<div className="h-full overflow-y-auto scrollbar-glassmorphism thin-scrollbar bg-gradient-to-br from-white/50 via-white/30 to-zinc-50/20 dark:from-zinc-900/50 dark:via-zinc-900/30 dark:to-black/20 backdrop-blur-xl">
				<div className="max-w-5xl mx-auto p-8 space-y-8">
					{/* Hero Header */}
					<div className="text-center backdrop-blur-xl bg-white/40 dark:bg-zinc-900/40 rounded-3xl p-12 border border-white/30 dark:border-zinc-700/30 shadow-2xl">
						<h1 className="text-6xl font-black bg-gradient-to-r from-foreground via-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-2xl mb-4">
							Profile
						</h1>
						<p className="text-2xl text-muted-foreground/90 max-w-2xl mx-auto leading-relaxed backdrop-blur-sm">
							Manage your sovereign identity and vault preferences
						</p>
					</div>

					<GlassProgressBar value={progress} label={`${stage} (${Math.round(progress)}%)`} visible={progressVisible} />

					<div className="grid lg:grid-cols-2 gap-8">
						{/* User Info Glass Card */}
						<div className="group backdrop-blur-2xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-8 border border-white/40 dark:border-zinc-700/40 shadow-2xl hover:shadow-primary/20 hover:shadow-3xl transition-all duration-500">
							<div className="flex items-start justify-between mb-8">
								<div className="space-y-2">
									<h2 className="text-3xl font-bold bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent">
										User Information
									</h2>
									<p className="text-xl text-muted-foreground/80">Your personal details</p>
								</div>
								<div className="w-16 h-16 bg-gradient-to-r from-primary/20 to-amber-500/20 rounded-2xl flex items-center justify-center backdrop-blur-sm border border-primary/30">
									<User className="h-8 w-8 text-primary" />
								</div>
							</div>

							{/* Avatar Upload */}
							<div className="flex flex-col md:flex-row items-center gap-6 mb-8 p-6 backdrop-blur-xl bg-white/50 dark:bg-zinc-800/50 rounded-2xl border border-white/40 group-hover:border-primary/40 transition-all">
								<div className="relative group/avatar">
									<div className="absolute inset-0 bg-gradient-to-r from-primary via-amber-500/30 to-primary rounded-3xl blur-xl opacity-0 group-hover/avatar:opacity-50 transition-all" />
									<Avatar className="h-28 w-28 border-4 border-white/50 shadow-2xl group-hover/avatar:shadow-primary/30 transition-all">
										{avatarUrl ? (
											<AvatarImage src={avatarUrl} />
										) : (
											<AvatarFallback className="bg-gradient-to-br from-primary/20 to-amber-500/20 text-3xl font-bold border-4 border-white/60 backdrop-blur-sm">
												{user?.username?.charAt(0).toUpperCase()}
											</AvatarFallback>
										)}
									</Avatar>
								</div>
								<div className="flex-1 text-center md:text-left">
									<input
										ref={fileInputRef}
										type="file"
										accept="image/jpeg,image/jpg,image/png,image/gif"
										className="hidden"
										onChange={handleAvatarUpload}
									/>
									<Button
										variant="outline"
										size="lg"
										onClick={() => fileInputRef.current?.click()}
										className="h-14 px-8 rounded-2xl backdrop-blur-sm bg-white/70 dark:bg-zinc-800/70 border-white/50 hover:bg-white/90 shadow-lg hover:shadow-xl font-semibold group-hover:border-primary/50 transition-all"
									>
										<Upload className="h-5 w-5 mr-2" />
										Change Avatar
									</Button>
									<p className="text-sm text-muted-foreground/80 mt-3">
										JPG, PNG or GIF. Max size 2MB.
									</p>
								</div>
							</div>

							<Separator className="bg-white/40 dark:bg-zinc-700/40 my-6" />

							{/* Form Fields */}
							<div className="grid md:grid-cols-2 gap-6">
								<div className="space-y-3">
									<Label className="text-lg font-semibold text-muted-foreground/90">Name</Label>
									<Input
										id="name"
										defaultValue={userName}
										className="h-14 text-xl rounded-2xl backdrop-blur-sm bg-white/50 dark:bg-zinc-800/50 border-white/40 hover:border-primary/40 shadow-inner font-semibold"
									/>
								</div>
								<div className="space-y-3">
									<Label className="text-lg font-semibold text-muted-foreground/90">Last Name</Label>
									<Input
										id="last_name"
										defaultValue={lastName}
										className="h-14 text-xl rounded-2xl backdrop-blur-sm bg-white/50 dark:bg-zinc-800/50 border-white/40 hover:border-primary/40 shadow-inner font-semibold"
									/>
								</div>
								<div className="md:col-span-2 space-y-3">
									<Label className="text-lg font-semibold text-muted-foreground/90">Email Address</Label>
									<Input
										id="email"
										type="email"
										defaultValue={user && user?.Email}
										className="h-14 text-xl rounded-2xl backdrop-blur-sm bg-white/50 dark:bg-zinc-800/50 border-white/40 hover:border-primary/40 shadow-inner font-semibold"
									/>
								</div>
							</div>

							<Button
								className="mt-8 h-16 px-12 text-xl font-bold bg-gradient-to-r from-[#C9A44A] to-[#B8934A] hover:from-[#C9A44A]/90 hover:to-[#B8934A]/90 shadow-2xl hover:shadow-[#C9A44A]/40 rounded-3xl w-full group hover:scale-[1.02] transition-all"
							>
								Save Changes
							</Button>
						</div>

						{/* Blockchain Connection */}
						<div className="space-y-6">
							<div className="group backdrop-blur-2xl bg-gradient-to-br from-primary/10 to-amber-500/5 rounded-3xl p-8 border border-primary/20 shadow-xl hover:shadow-2xl hover:shadow-primary/30 transition-all duration-500">
								<div className="flex items-start justify-between mb-6">
									<div className="space-y-2">
										<h3 className="text-3xl font-bold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
											Stellar Connection
										</h3>
										<p className="text-xl text-primary/90 font-semibold">Your sovereign identity</p>
									</div>
									<div className="w-14 h-14 bg-gradient-to-r from-primary to-amber-500 rounded-2xl flex items-center justify-center shadow-lg">
										<Key className="h-7 w-7 text-white" />
									</div>
								</div>

								<div className="space-y-4">
									<div className="space-y-2">
										<Label className="text-lg font-semibold flex items-center gap-2 text-muted-foreground/90">
											<span>Public Key</span>
										</Label>
										<div className="flex gap-3">
											<Input
												value={publicKey}
												readOnly
												className="flex-1 h-14 font-mono text-lg bg-white/40 dark:bg-zinc-800/40 backdrop-blur-sm rounded-2xl border-primary/30 shadow-inner font-semibold"
											/>
											<Button
												variant="outline"
												size="icon"
												className="h-14 w-14 rounded-2xl backdrop-blur-sm border-primary/40 hover:bg-primary/10 shadow-md hover:shadow-lg"
												onClick={() => handleCopyStellarAddress(publicKey)}
											>
												<Copy className="h-5 w-5" />
											</Button>
										</div>
									</div>

									<div className="space-y-2">
										<Label className="text-lg font-semibold flex items-center gap-2 text-muted-foreground/90">
											<span>Secret Key</span>
										</Label>
										<div className="flex gap-3">
											<Input
												value={privateKey}
												readOnly
												type="password"
												className="flex-1 h-14 font-mono text-lg bg-white/40 dark:bg-zinc-800/40 backdrop-blur-sm rounded-2xl border-destructive/30 shadow-inner font-semibold text-destructive/80"
											/>
											<Button
												variant="outline"
												size="icon"
												className="h-14 w-14 rounded-2xl backdrop-blur-sm border-destructive/40 hover:bg-destructive/10 shadow-md hover:shadow-lg"
												onClick={() => handleCopyStellarAddress(privateKey)}
											>
												<Copy className="h-5 w-5" />
											</Button>
										</div>
									</div>

									<div className="flex items-center gap-4 pt-4 border-t border-primary/20">
										<Badge className="bg-green-500/20 text-green-600 border-green-500/30 px-4 py-2 text-lg font-semibold backdrop-blur-sm">
											● Connected
										</Badge>
										<Button
											variant="outline"
											size="lg"
											onClick={handleRotateKey}
											className="h-12 px-6 rounded-2xl backdrop-blur-sm border-primary/40 hover:bg-primary/10 font-semibold shadow-md hover:shadow-lg"
										>
											<RefreshCw className="h-4 w-4 mr-2 animate-spin" />
											Rotate Keys
										</Button>
									</div>
								</div>
							</div>

							{/* Quick Actions */}
							<div className="grid grid-cols-2 gap-4">
								<Button
									variant="outline"
									size="lg"
									onClick={handleSyncBeta}
									className="h-20 rounded-2xl backdrop-blur-xl bg-white/60 dark:bg-zinc-800/60 border-white/40 hover:bg-white/80 shadow-xl hover:shadow-2xl font-semibold group"
								>
									<RefreshCw className="h-6 w-6 mr-3 group-hover:rotate-180 transition-all" />
									Sync Now
								</Button>

								<div className="p-4">
									<button onClick={() => setShowModal(true)} className="open-btn">
										Verify Encryption
									</button>

									{showModal && (
										<div className="modal-overlay">
											<EncryptionVerificationModalBeta
												file={mockFile}
												onClose={() => setShowModal(false)}
											/>
										</div>
									)}
								</div>
								<Button
									variant="outline"
									size="lg"
									onClick={handleGenerateApiKey}
									className="h-20 rounded-2xl backdrop-blur-xl bg-white/60 dark:bg-zinc-800/60 border-white/40 hover:bg-white/80 shadow-xl hover:shadow-2xl font-semibold"
								>
									<Key className="h-6 w-6 mr-3" />
									Generate API Key
								</Button>
							</div>
						</div>
					</div>

					{/* Device Stats */}
					<div className="group backdrop-blur-2xl bg-white/50 dark:bg-zinc-900/50 rounded-3xl p-10 border border-white/30 dark:border-zinc-700/30 shadow-2xl hover:shadow-xl transition-all">
						<div className="flex items-start justify-between mb-8">
							<div className="space-y-2">
								<h3 className="text-3xl font-bold bg-gradient-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
									Device Status
								</h3>
								<p className="text-xl text-muted-foreground/80">Local storage & sync metrics</p>
							</div>
							<div className="w-16 h-16 bg-gradient-to-r from-muted-foreground/20 to-muted-foreground/10 rounded-2xl flex items-center justify-center backdrop-blur-sm">
								<HardDrive className="h-8 w-8 text-muted-foreground" />
							</div>
						</div>

						<div className="grid md:grid-cols-4 gap-8 text-center">
							<div className="space-y-3">
								<div className="text-lg font-black text-muted-foreground mb-6">{vault?.Vault?.name}</div>
								<div className="text-sm text-muted-foreground/70 uppercase tracking-wider">Name</div>
							</div>
							<div className="space-y-3">
								<div className="text-3xl font-black text-muted-foreground">{totalEntries}</div>
								<div className="text-sm text-muted-foreground/70 uppercase tracking-wider">Entries</div>
							</div>
							<div className="space-y-3">
								<div className="w-full bg-gradient-to-r from-secondary/30 to-secondary/10 rounded-2xl p-4 backdrop-blur-sm border border-white/40">
									<div className="flex items-center gap-2 text-sm font-mono text-muted-foreground mb-2">
										<HardDrive className="h-4 w-4" />
										device-7a8b9c
									</div>
									<div className="space-y-1">
										<div className="text-xs text-muted-foreground/70">Storage</div>
										<Progress value={usagePercent} className="h-2 [&>div]:bg-gradient-to-r [&>div]:from-[#C9A44A] [&>div]:to-[#B8934A] shadow-lg" />
										<div className="text-xs font-mono text-right">{totalEntries} / {maxEntries}</div>
									</div>
								</div>
							</div>
							<div className="space-y-3">
								<div className="text-lg font-bold text-muted-foreground">{lastSync}</div>
								<div className="flex items-center justify-center gap-2 text-sm text-muted-foreground/70">
									<Clock className="h-4 w-4" />
									Last Sync
								</div>
							</div>
						</div>
					</div>

					{/* Footer */}
					<div className="backdrop-blur-xl bg-gradient-to-r from-white/70 to-zinc-100/50 dark:from-zinc-900/70 dark:to-zinc-800/50 rounded-3xl p-8 border border-white/40 shadow-2xl">
						<div className="flex flex-col sm:flex-row items-center justify-between gap-6">
							<div className="text-center sm:text-left">
								<p className="text-xl font-bold text-muted-foreground">Ankhora v{vault?.Vault?.version || "1.2.0"}</p>
								<p className="text-lg text-muted-foreground/70">
									Last updated: {formatMonthYear(vault?.LastUpdated)}
								</p>
							</div>
							<Button
								variant="destructive"
								size="lg"
								onClick={handleLogout}
								className="h-16 px-12 text-xl font-bold bg-gradient-to-r from-destructive to-destructive/80 hover:from-destructive/90 shadow-2xl hover:shadow-destructive/30 rounded-3xl w-full sm:w-auto group hover:scale-[1.02]"
							>
								<LogOut className="h-5 w-5 mr-3 group-hover:translate-x-1 transition-all" />
								Sign Out
							</Button>
						</div>
					</div>
				</div>
			</div>
		</DashboardLayout>
	);
};

export default Profile;

