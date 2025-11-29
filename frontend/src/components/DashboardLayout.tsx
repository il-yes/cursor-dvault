import { ReactNode, useState, useMemo, useEffect } from "react";
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

const dashboardNavItems = [
  { title: "Dashboard", url: "/dashboard", icon: Home },
  { title: "Vault", url: "/dashboard/vault", icon: Shield },
  { title: "Shares", url: "/dashboard/shared", icon: Rocket },
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
  const { vaultContext, addEntry, clearVault: clearVaultContext } = useVault();
  const { clearVault: clearVaultStore } = useVaultStore();

  const { updateOnboarding, onboarding, jwtToken } = useAuthStore();
  const [encryptedVault, setEncryptedVault] = useState<string | null>(null);
  const [vaultDecryptionPassword, setVaultDecryptionPassword] = useState<string | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [vault, setVault] = useState<VaultPayload | null>(null);
  const [decryptionPassword, setDecryptionPassword] = useState<string | null>(null);

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
      const rawEntry = await AppAPI.AddEntry(entry.type, entryPayload, jwtToken);

      // 3️⃣ Convert backend response if needed
      const newEntry: VaultEntry = {
        ...rawEntry,
        created_at: rawEntry.created_at || new Date().toISOString(),
        updated_at: rawEntry.updated_at || new Date().toISOString(),
      };

      // 4️⃣ Update Zustand store
      addEntry(newEntry);

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
          <div className="h-8 w-8 rounded-lg bg-gradient-primary flex items-center justify-center">
            <Shield className="h-5 w-5 text-primary-foreground" />
          </div>
          <span className="text-lg font-bold text-foreground">VaultCore</span>
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
            <DropdownMenuItem onClick={() => navigate("/dashboard/settings")}>
              <Settings className="mr-2 h-4 w-4" />
              Settings
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

function AppSidebar() {
  const location = useLocation();
  const navigate = useNavigate();
  const isVaultContext = location.pathname.startsWith("/dashboard/vault");
  const isSharedContext = location.pathname.startsWith("/dashboard/shared");
  const { vaultContext, addFolder } = useVault();
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


  return (
    <Sidebar className="border-r border-border w-[220px]">
      <SidebarContent>
        {(isVaultContext || isSharedContext) && (
          <div className="p-4 border-b border-border">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate("/dashboard")}
              className="w-full justify-start"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Dashboard
            </Button>
          </div>
        )}

        <SidebarGroup>
          <SidebarGroupLabel className="text-muted-foreground uppercase tracking-wider text-xs">
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
                        `flex items-center gap-3 px-3 py-2 rounded-lg transition-all ${isActive
                          ? "bg-primary/10 text-primary font-medium"
                          : "text-muted-foreground hover:bg-secondary hover:text-foreground"
                        }`
                      }
                    >
                      <item.icon className="h-5 w-5" />
                      <span>{item.title}</span>
                    </NavLink>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {secondaryItems.length > 0 && (
          <SidebarGroup>
            <SidebarGroupLabel className="text-muted-foreground uppercase tracking-wider text-xs">
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
                          `flex items-center gap-3 px-3 py-2 rounded-lg transition-all ${isActive
                            ? "bg-primary/10 text-primary font-medium"
                            : "text-muted-foreground hover:bg-secondary hover:text-foreground"
                          }`
                        }
                      >
                        <item.icon className="h-5 w-5" />
                        <span>{item.title}</span>
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
          <SidebarGroup>
            <SidebarGroupLabel className="text-muted-foreground uppercase tracking-wider text-xs flex justify-between">
              <span>Folders</span>
              <SidebarMenuButton asChild>
                <button
                  onClick={() => setIsNewFolderOpen(true)}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg transition-all text-muted-foreground hover:bg-secondary hover:text-foreground text-left"
                  style={{ width: "auto" }}
                >
                  <Plus className="h-5 w-5" />
                </button>
              </SidebarMenuButton>
            </SidebarGroupLabel>
            {/* {isVaultContext && vaultContext?.Vault.folders && vaultContext.Vault.folders.length > 0 && ( */}
            <SidebarGroupContent>
              <SidebarMenu>
                {isVaultContext && vaultContext?.Vault.folders && vaultContext.Vault.folders.length > 0 && vaultContext.Vault.folders.map((folder) => (
                  <SidebarMenuItem key={folder.id}>
                    <SidebarMenuButton asChild>
                      <NavLink
                        to={`/dashboard/vault/folder/${folder.id}`}
                        className={({ isActive }) =>
                          `flex items-center gap-3 px-3 py-2 rounded-lg transition-all ${isActive
                            ? "bg-primary/10 text-primary font-medium"
                            : "text-muted-foreground hover:bg-secondary hover:text-foreground"
                          }`
                        }
                      >
                        <Folder className="h-5 w-5" />
                        <span>{folder.name}</span>
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
            {/* )} */}
          </SidebarGroup>
        )}

        {/* New Share Button (Shared Entries only) */}
        {isSharedContext && (
          <div className="mt-auto p-4">
            <Button
              onClick={() => setIsNewShareOpen(true)}
              className="w-full bg-[#C9A44A] hover:bg-[#B8934A]"
            >
              <Plus className="h-4 w-4 mr-2" />
              New Share
            </Button>
          </div>
        )}
      </SidebarContent>

      {/* Upgrade Button at Bottom */}
      <div className="mt-auto p-4 border-t border-border space-y-3">
        <Button
          onClick={() => setIsUpgradeOpen(true)}
          className=" py-3 px-4 bg-[#C9A44A] hover:bg-[#B8934A] text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
        >
          <Crown className="h-4 w-4" />
          Upgrade to Premium
        </Button>
        {/* User Info at Bottom */}
          <div className="flex items-center gap-3 px-2 py-3 rounded-lg bg-secondary/30 border-t border border-border">
            <Avatar className="h-9 w-9">
              <AvatarFallback className="bg-primary/10 text-sm">
                <User className="h-5 w-5 text-primary" />
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-foreground truncate">{user && user?.username}</p>
              <p className="text-xs text-muted-foreground truncate">{user && user?.email}</p>
            </div>
          </div>
      </div>

      {/* Upgrade Modal */}
      <OnboardingModal
        open={isUpgradeOpen}
        onOpenChange={setIsUpgradeOpen}
      />

      {/* New Folder Dialog */}
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
              />
            </div>
            <div className="flex justify-end gap-3">
              <Button variant="outline" onClick={() => setIsNewFolderOpen(false)}>
                Cancel
              </Button>
              <Button
                onClick={handleCreateFolder}
                disabled={!newFolderName.trim()}
                className="bg-[#C9A44A] hover:bg-[#B8934A]"
              >
                Create Folder
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {/* New Share Modal */}
      <NewShareModal
        open={isNewShareOpen}
        onOpenChange={setIsNewShareOpen}
        onShareSuccess={() => {
          setSharedEntriesRefreshKey(prev => prev + 1);
          // Trigger a custom event to notify SharedEntriesLayout
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
          <main className="flex-1 overflow-auto">
            {children}
          </main>
        </div>
      </div>
    </SidebarProvider>
  );
}
