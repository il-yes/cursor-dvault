import { ReactNode, useState, useMemo, useEffect, forwardRef } from "react";
import { NavLink, useNavigate, useLocation } from "react-router-dom";
import {
  Database, Shield, Settings, LogOut, Menu, Search, User,
  Home, Rocket, Info, HelpCircle, Folder, Star, Trash2,
  LogIn, CreditCard, UserCircle, FileText, Key, ArrowLeft,
  Plus, Crown, X, Clock, Users
} from "lucide-react";
import { useVault } from "@/hooks/useVault";
import { OnboardingModal } from "@/components/OnboardingModal";
import { NewShareModal } from "@/components/NewShareModal";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { CreateEntryDialog } from "@/components/CreateEntryDialog"
import { SearchOverlay } from "./SearchOverlay";
import { VaultEntry, VaultPayload } from "@/types/vault";
import { useAuthStore } from "@/store/useAuthStore";
import { toast } from "@/hooks/use-toast";
import * as AppAPI from "../../wailsjs/go/main/App";
import { useAppStore } from "@/store/appStore";
import { useVaultStore } from "@/store/vaultStore";
import { ThemeToggle } from "@/components/ThemeToggle";
import { NavLink as ReactRouterNavLink, NavLinkProps as ReactRouterNavLinkProps } from "react-router-dom";
import "./contributionGraph/g-scrollbar.css";
import AvatarImg from '@/assets/7.jpg'
import { OnboardingModalBeta } from "./OnboardingModalBeta";
import { withAuth } from "@/hooks/withAuth";
import { auth } from "wailsjs/go/models";



interface CustomNavLinkProps extends Omit<ReactRouterNavLinkProps, 'className'> {
  children: React.ReactNode;
  className?: string;
  activeClassName?: string;
}

const CustomNavLink = forwardRef<HTMLAnchorElement, CustomNavLinkProps>(
  ({ children, className = "", activeClassName = "", ...props }, ref) => {
    return (
      <ReactRouterNavLink
        ref={ref}
        className={({ isActive }) =>
          `${className} ${isActive ? activeClassName : ""}`
        }
        {...props}
      >
        {children}
      </ReactRouterNavLink>
    );
  }
);
CustomNavLink.displayName = "CustomNavLink";


const dashboardNavItems = [
  { title: "Dashboard", url: "/dashboard", icon: Home },
  { title: "Vault", url: "/dashboard/vault", icon: Shield },
  { title: "Shares", url: "/dashboard/shared", icon: Rocket },
  { title: "Onboarding", url: "/dashboard/on-boarding", icon: Rocket },
];

const dashboardSecondaryItems = [
  { title: "About", url: "/dashboard/about", icon: Info },
  { title: "Feedback", url: "/dashboard/feedback", icon: HelpCircle },
];

const sharedEntriesItems = [
  { title: "All", filter: "all", url: "/dashboard/shared", icon: Folder },
  { title: "Sent", filter: "sent", url: "/dashboard/shared?filter=sent", icon: LogOut },
  { title: "Received", filter: "received", url: "/dashboard/shared?filter=received", icon: LogIn },
  { title: "Pending", filter: "pending", url: "/dashboard/shared?filter=pending", icon: Clock },
  { title: "Revoked", filter: "revoked", url: "/dashboard/shared?filter=revoked", icon: X },
  { title: "With me", filter: "withme", url: "/dashboard/shared?filter=withme", icon: Users },
];

const vaultMainItems = [
  { title: "All", url: "/dashboard/vault", icon: Folder },
  { title: "Favorites", url: "/dashboard/vault/favorites", icon: Star },
  { title: "Trash", url: "/dashboard/vault/trash", icon: Trash2 },
];

const vaultSecondaryItems = [
  { title: "Identifiers", type: "login", url: "/dashboard/vault/login", icon: LogIn },
  { title: "Payment card", type: "card", url: "/dashboard/vault/card", icon: CreditCard },
  { title: "Identity", type: "identity", url: "/dashboard/vault/identity", icon: UserCircle },
  { title: "Secure note", type: "note", url: "/dashboard/vault/note", icon: FileText },
  { title: "SSH key", type: "sshkey", url: "/dashboard/vault/sshkey", icon: Key },
];

function DashboardNavbar() {
  const navigate = useNavigate();
  const location = useLocation();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [showSearchOverlay, setShowSearchOverlay] = useState(false);

  const vaultContext = useVaultStore((state) => state.vault);
  const addEntry = useVaultStore((state) => state.addEntry);
  const clearVaultStore = useVaultStore((state) => state.clearVault);

  const { updateOnboarding, onboarding, jwtToken } = useAuthStore();
  const [encryptedVault, setEncryptedVault] = useState<string | null>(null);
  const [vaultDecryptionPassword, setVaultDecryptionPassword] = useState<string | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [vault, setVault] = useState<VaultPayload | null>(null);
  const [decryptionPassword, setDecryptionPassword] = useState<string | null>(null);
  const auth = useAuthStore.getState();

  const isVaultContext = location.pathname.startsWith("/dashboard/vault");

  const allEntries = useMemo(() => {
    if (!vaultContext?.Vault) return [];
    return [
      ...(vaultContext.Vault.entries?.login || []),
      ...(vaultContext.Vault.entries?.card || []),
      ...(vaultContext.Vault.entries?.note || []),
      ...(vaultContext.Vault.entries?.sshkey || []),
      ...(vaultContext.Vault.entries?.identity || []),
    ];
  }, [vaultContext]);

  const handleLogout = () => {
    // Clear all stores and state
    clearVaultStore();                   // Clear vault store (entries, vault data)
    useAuthStore.getState().clearAll();  // Clear auth store (user, tokens)
    useAppStore.getState().reset();      // Clear app store (session)

    // Clear specific localStorage items (not all, to preserve settings)
    localStorage.removeItem('userId');
    localStorage.removeItem('vault-storage');

    toast({
      title: "Logged out",
      description: "You have been successfully logged out.",
    });
    AppAPI.SignOut(auth.user.id);
    navigate("/login/email");
  };
  const handleAddEntry = () => {
    setIsCreateDialogOpen(true);
  };
  const handleCreateEntry = async (
    entry: Omit<VaultEntry, "id" | "created_at" | "updated_at">
  ) => {
    try {
      // Debug: Check if jwtToken is available
      console.log("🔑 JWT Token:", jwtToken);


      if (!jwtToken) {
        toast({
          title: "Authentication Error",
          description: "You are not authenticated. Please log in again.",
          variant: "destructive",
        });
        return;
      }

      // 1️⃣ Prepare entry to send to backend
      const entryPayload = {
        ...entry,
        // Backend will usually generate id / timestamps
      };

      // 2️⃣ Send to backend (v0 logic)
      // const rawEntry = await AppAPI.AddEntry(entry.type, entryPayload, jwtToken);
      const rawEntry = await withAuth((token) => {
        console.log("🚀 ~ withAuth ~ token:", token)
        return AppAPI.AddEntry(entry.type, entryPayload, token)
      });
      console.log("🚀 ~ handleCreateEntry ~ rawEntry:", rawEntry)

      // 3️⃣ Convert backend response if needed
      const newEntry: VaultEntry = {
        ...rawEntry,
        created_at: rawEntry.created_at || new Date().toISOString(),
        updated_at: rawEntry.updated_at || new Date().toISOString(),
      };

      // 4️⃣ Update Zustand store
      addEntry(newEntry);
      console.log("🚀 ~ handleCreateEntry ~ newEntry:", newEntry)

      // 5️⃣ Show feedback
      toast({
        title: "Entry created",
        description: `${newEntry.entry_name} has been added to your vault.`,
      });
    } catch (err) {
      console.error("Failed to create entry:", err);
      toast({
        title: "Error",
        description: "Could not save entry. Please try again.",
        variant: "destructive",
      });
    }
  };

  const handleSelectEntry = (entry: VaultEntry) => {
    navigate(`/dashboard/vault/${entry.type}?entry=${entry.id}`);
    setSearchQuery("");
    setShowSearchOverlay(false);
  };

  useEffect(() => {
    setShowSearchOverlay(searchQuery.trim().length > 0);
  }, [searchQuery]);

  const SaveSessionTest = async () => {
    await AppAPI.SaveSessionTest(auth.jwtToken);
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-16 items-center px-4 md:px-6">
        <div className="flex items-center gap-2 md:hidden">
          <SidebarTrigger>
            <Button variant="ghost" size="icon">
              <Menu className="h-5 w-5" />
            </Button>
          </SidebarTrigger>
        </div>

        <div className="hidden md:flex items-center gap-3">
          {/* <img src={ankhoraLogo} alt="Ankhora Logo" className="h-9 w-auto" /> */}
          <span onClick={() => SaveSessionTest()} className="text-lg font-bold text-foreground"><small>ANKHORA</small></span>
        </div>

        <div className="ml-20" style={{ display: "flex", justifyContent: "center", alignItems: "center", cursor: "pointer" }}  >
          <Plus
            data-tooltip-target="tooltip-default"
            onClick={handleAddEntry}
            className="h-5 w-8"
          />
        </div>

        <div className="flex-1 flex justify-center px-4 md:px-8">
          <div className="w-full max-w-md relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search entries, identities, or cards…"
              className="pl-9 bg-secondary/50"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              disabled={!isVaultContext}
            />
            {isVaultContext && showSearchOverlay && (
              <SearchOverlay
                entries={allEntries}
                searchQuery={searchQuery}
                onSelectEntry={handleSelectEntry}
                onClose={() => {
                  setSearchQuery("");
                  setShowSearchOverlay(false);
                }}
              />
            )}
          </div>
        </div>

        <ThemeToggle />

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="rounded-full">
              <Avatar className="h-8 w-8">
                <AvatarFallback className="bg-primary/10">
                  <User className="h-4 w-4 text-primary" />
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <DropdownMenuItem onClick={() => navigate("/dashboard/profile")}>
              <User className="mr-2 h-4 w-4" />
              Profile
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => navigate("/dashboard/profile-beta")}>
              <User className="mr-2 h-4 w-4" />
              Profile (beta)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => navigate("/dashboard/settings")}>
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => navigate("/dashboard/settings-beta")}>
              <Settings className="mr-2 h-4 w-4" />
              Settings (beta)
            </DropdownMenuItem>
            <DropdownMenuItem onClick={handleLogout} className="text-destructive focus:text-destructive">
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Create Entry Dialog */}
      <CreateEntryDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSubmit={handleCreateEntry}
      />
    </header>
  );
}


const Avatars = [
  { 'id': "38", src: AvatarImg },
  { 'id': "37", src: "https://i.ebayimg.com/images/g/eEEAAOSweZVjNAXV/s-l1200.jpg" },
  { 'id': "34", src: 'https://www.independent.com/wp-content/uploads/2017/08/01/raekwon.jpg' },
  { 'id': "39", src: 'https://upload.wikimedia.org/wikipedia/commons/1/1d/Teyana_Taylor_%28cropped%29.jpg' }
]
const RenderAvatar = (id: string | number) => {
  const img = Avatars.find((f) => String(f.id) === String(id));
  return img ? img.src : ""
}

function AppSidebar() {
  const location = useLocation();
  const navigate = useNavigate();
  const isVaultContext = location.pathname.startsWith("/dashboard/vault");
  const isSharedContext = location.pathname.startsWith("/dashboard/shared");

  const vaultContext = useVaultStore((state) => state.vault);
  const addFolder = useVaultStore((state) => state.addFolder);
  const [isUpgradeOpen, setIsUpgradeOpen] = useState(false);
  const [isNewFolderOpen, setIsNewFolderOpen] = useState(false);
  const [isNewShareOpen, setIsNewShareOpen] = useState(false);
  const [newFolderName, setNewFolderName] = useState("");
  const [sharedEntriesRefreshKey, setSharedEntriesRefreshKey] = useState(0);
  const { user } = useAuthStore();

  const mainItems = isVaultContext ? vaultMainItems : isSharedContext ? sharedEntriesItems : dashboardNavItems;
  const secondaryItems = isVaultContext ? vaultSecondaryItems : isSharedContext ? [] : dashboardSecondaryItems;

  const handleCreateFolder = () => {
    if (newFolderName.trim()) {
      addFolder(newFolderName);
      setNewFolderName("");
      setIsNewFolderOpen(false);
    }
  };
  const avatar = user && RenderAvatar(user.id)



  return (
    <Sidebar className="border-r border-transparent w-[240px] backdrop-blur-sm bg-white/40 dark:bg-zinc-900/40 shadow-2xl ">
      <SidebarContent className="backdrop-blur-sm bg-white/30 dark:bg-zinc-900/30 border-r border-zinc-200/30 dark:border-zinc-700/30 scrollbar-glassmorphism thin-scrollbar">
        {(isVaultContext || isSharedContext) && (
          <div className="p-4 border-b border-zinc-200/30 dark:border-zinc-700/30 bg-white/20 dark:bg-zinc-900/20 ">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate("/dashboard")}
              className="w-full justify-start h-11 rounded-xl backdrop-blur-sm bg-white/50 dark:bg-zinc-800/50 hover:bg-white/70 dark:hover:bg-zinc-800/70 border border-zinc-200/50 dark:border-zinc-700/50 transition-all group"
            >
              <ArrowLeft className="mr-3 h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              <span className="text-sm font-medium">Back to Dashboard</span>
            </Button>
          </div>
        )}

        <SidebarGroup >
          <SidebarGroupLabel style={{ marginTop: "30px" }} className="px-4 py-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground/80 bg-white/20 dark:bg-zinc-900/20 border-b border-zinc-200/20 dark:border-zinc-700/20">
            {isVaultContext ? "Vault" : isSharedContext ? "Shares" : "Main"}
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {mainItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <NavLink
                      to={item.url}
                      end
                      className={({ isActive }) =>
                        `group flex items-center gap-3 px-4 py-3 rounded-2xl mx-2 my-1 transition-all duration-200 backdrop-blur-sm border border-transparent hover:border-primary/30 hover:bg-white/50 dark:hover:bg-zinc-800/50 hover:shadow-md ${isActive
                          ? "bg-gradient-to-r from-primary/20 to-amber-500/20 text-primary font-semibold shadow-lg border-primary/40"
                          : "text-muted-foreground hover:text-foreground"
                        }`
                      }
                    >
                      <item.icon className="h-5 w-5 flex-shrink-0 group-hover:scale-110 transition-transform" />
                      <span className="text-sm font-medium">{item.title}</span>
                    </NavLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {secondaryItems.length > 0 && (
          <SidebarGroup style={{ marginTop: "30px" }} >
            <SidebarGroupLabel className="px-4 py-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground/80 bg-white/20 dark:bg-zinc-900/20 border-b border-zinc-200/20 dark:border-zinc-700/20">
              {isVaultContext ? "Entry Types" : "More"}
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {secondaryItems.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild>
                      <NavLink
                        to={item.url}
                        className={({ isActive }) =>
                          `group flex items-center gap-3 px-4 py-3 rounded-2xl mx-2 my-1 transition-all duration-200 backdrop-blur-sm border border-transparent hover:border-primary/30 hover:bg-white/50 dark:hover:bg-zinc-800/50 hover:shadow-md ${isActive
                            ? "bg-gradient-to-r from-primary/20 to-amber-500/20 text-primary font-semibold shadow-lg border-primary/40"
                            : "text-muted-foreground hover:text-foreground"
                          }`
                        }
                      >
                        <item.icon className="h-5 w-5 flex-shrink-0 group-hover:scale-110 transition-transform" />
                        <span className="text-sm font-medium">{item.title}</span>
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}

        {/* Folders Section (Vault only) */}
        {isVaultContext && (
          <SidebarGroup style={{ marginTop: "30px" }} >
            <SidebarGroupLabel className="px-4 py-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground/80 flex justify-between items-center bg-white/20 dark:bg-zinc-900/20 border-b border-zinc-200/20 dark:border-zinc-700/20">
              <span>Folders</span>
              <SidebarMenuButton asChild>
                <button
                  onClick={() => setIsNewFolderOpen(true)}
                  className="p-2 rounded-xl hover:bg-white/50 dark:hover:bg-zinc-800/50 transition-all hover:scale-105 group"
                >
                  <Plus className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
                </button>
              </SidebarMenuButton>
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {isVaultContext && vaultContext?.Vault.folders && vaultContext.Vault.folders.length > 0 && vaultContext.Vault.folders.map((folder) => (
                  <SidebarMenuItem key={folder.id}>
                    <SidebarMenuButton asChild>
                      <NavLink
                        to={`/dashboard/vault/folder/${folder.id}`}
                        className={({ isActive }) =>
                          `group flex items-center gap-3 px-4 py-3 rounded-2xl mx-2 my-1 transition-all duration-200 backdrop-blur-sm border border-transparent hover:border-primary/30 hover:bg-white/50 dark:hover:bg-zinc-800/50 hover:shadow-md ${isActive
                            ? "bg-gradient-to-r from-primary/20 to-amber-500/20 text-primary font-semibold shadow-lg border-primary/40"
                            : "text-muted-foreground hover:text-foreground"
                          }`
                        }
                      >
                        <Folder className="h-5 w-5 flex-shrink-0 group-hover:scale-110 transition-transform" />
                        <span className="text-sm font-medium">{folder.name}</span>
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}

        {/* New Share Button (Shared Entries only) */}
        {isSharedContext && (
          <div className="mt-auto p-6 bg-white/20 dark:bg-zinc-900/20 backdrop-blur-sm">
            <Button
              onClick={() => setIsNewShareOpen(true)}
              className="w-full h-12 rounded-2xl bg-gradient-to-r from-[#C9A44A] to-[#B8934A] hover:from-[#C9A44A]/90 hover:to-[#B8934A]/90 shadow-xl hover:shadow-[#C9A44A]/25 text-white font-semibold text-sm transition-all hover:scale-[1.02] active:scale-[0.98]"
            >
              <Plus className="h-4 w-4 mr-2" />
              New Share
            </Button>
          </div>
        )}
      </SidebarContent>

      {/* Premium Footer */}
      <div className="mt-auto p-6 border-t border-zinc-200/30 dark:border-zinc-700/30 bg-gradient-to-b from-white/50 to-white/30 dark:from-zinc-900/50 dark:to-zinc-900/30 backdrop-blur-sm space-y-4">
        <Button
          onClick={() => setIsUpgradeOpen(true)}
          className="w-full h-12 py-3 px-4 bg-gradient-to-r from-[#C9A44A] to-[#B8934A] hover:from-[#C9A44A]/90 hover:to-[#B8934A]/90 shadow-2xl hover:shadow-[#C9A44A]/30 text-white font-semibold rounded-2xl transition-all hover:scale-[1.02] active:scale-[0.98]"
        >
          <Crown className="h-4 w-4 mr-2" />
          Upgrade to Premium
        </Button>

        {/* User Info */}
        <div className="flex items-center gap-3 p-4 rounded-2xl bg-white/40 dark:bg-zinc-800/40 backdrop-blur-sm border border-zinc-200/30 dark:border-zinc-700/30 hover:bg-white/60 dark:hover:bg-zinc-800/60 transition-all group">
          <Avatar className="h-10 w-10 flex-shrink-0">
            <AvatarFallback className="bg-gradient-to-br from-primary/20 to-amber-500/20 backdrop-blur-sm border border-primary/20 text-sm">
              {avatar ? <img src={avatar} alt="User Avatar" className="h-5 w-5" /> : <User className="h-5 w-5 text-primary" />}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-semibold text-foreground truncate group-hover:text-primary transition-colors">{user?.username}</p>
            <p className="text-xs text-muted-foreground truncate">{user?.email}</p>
          </div>
        </div>
      </div>

      {/* Modals unchanged */}
      <OnboardingModalBeta open={isUpgradeOpen} onOpenChange={setIsUpgradeOpen} />
      <Dialog open={isNewFolderOpen} onOpenChange={setIsNewFolderOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Create New Folder</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="folder-name">Folder Name</Label>
              <Input
                id="folder-name"
                placeholder="e.g., Work Accounts"
                value={newFolderName}
                onChange={(e) => setNewFolderName(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') handleCreateFolder();
                }}
                className="rounded-xl h-11"
              />
            </div>
            <div className="flex justify-end gap-3">
              <Button variant="outline" onClick={() => setIsNewFolderOpen(false)} className="rounded-xl h-10">
                Cancel
              </Button>
              <Button
                onClick={handleCreateFolder}
                disabled={!newFolderName.trim()}
                className="bg-[#C9A44A] hover:bg-[#B8934A] rounded-xl h-10"
              >
                Create Folder
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
      <NewShareModal
        open={isNewShareOpen}
        onOpenChange={setIsNewShareOpen}
        onShareSuccess={() => {
          setSharedEntriesRefreshKey(prev => prev + 1);
          window.dispatchEvent(new CustomEvent('shareEntriesRefresh'));
        }}
      />
    </Sidebar>
  );
}

export function DashboardLayout({ children }: { children: ReactNode }) {
  const { user } = useAuthStore();

  if (!user) {
    return <div>Loading session...</div>;
  }

  return (
    <SidebarProvider>
      <div className="flex h-screen w-full bg-background overflow-hidden">
        <AppSidebar />
        <div className="flex-1 flex flex-col overflow-hidden">
          <DashboardNavbar />
          <main className="flex-1 overflow-auto scrollbar-glassmorphism thin-scrollbar" style={{ paddingRight: "25px" }}>
            {children}
          </main>
        </div>
      </div>
    </SidebarProvider>
  );
}
