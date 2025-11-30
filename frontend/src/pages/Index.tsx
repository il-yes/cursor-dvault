import { DashboardLayout } from "@/components/DashboardLayout";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Shield, Database, Activity, Users } from "lucide-react";
import { GlobalSecurityInsight } from "@/components/GlobalSecurityInsight";
import { useVaultStore } from "@/store/vaultStore";
import ContributionGraph from "@/components/ContributionGraph"


const Index = () => {
  const { vault, lastSyncTime, loadVault } = useVaultStore();
  const sharedEntries = useVaultStore((state) => state.shared.items);
  const totalEntries = Object.values(vault?.Vault?.entries || {}).flat().length;

  const sampleData = {
    "2025-01-01": 2,
    "2025-01-02": 4,
    "2025-01-03": 10,
  };

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div className="space-y-2">
          <h1 className="text-4xl font-semibold tracking-tight bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent">
            
          </h1>
          <p className="text-lg text-muted-foreground max-w-md pt-3">
            <small><em>Overview of your vault activity and statistics</em></small>
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          <Card className="group border-none shadow-lg hover:shadow-xl backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60 hover:bg-white/80 dark:hover:bg-zinc-900/70 transition-all duration-300 border-zinc-100 dark:border-zinc-800">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <CardTitle className="text-base font-medium text-foreground group-hover:text-primary transition-colors">
                Total Entries
              </CardTitle>
              <Database className="h-5 w-5 text-muted-foreground group-hover:text-primary/70 transition-colors" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                {totalEntries}
              </div>
              <p className="text-xs text-muted-foreground mt-1">+2 from last month</p>
            </CardContent>
          </Card>

          <Card className="group border-none shadow-lg hover:shadow-xl backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60 hover:bg-white/80 dark:hover:bg-zinc-900/70 transition-all duration-300 border-zinc-100 dark:border-zinc-800">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <CardTitle className="text-base font-medium text-foreground group-hover:text-primary transition-colors">
                Encrypted
              </CardTitle>
              <Shield className="h-5 w-5 text-muted-foreground group-hover:text-primary/70 transition-colors" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                100%
              </div>
              <p className="text-xs text-muted-foreground mt-1">All data secured</p>
            </CardContent>
          </Card>

          <Card className="group border-none shadow-lg hover:shadow-xl backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60 hover:bg-white/80 dark:hover:bg-zinc-900/70 transition-all duration-300 border-zinc-100 dark:border-zinc-800">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <CardTitle className="text-base font-medium text-foreground group-hover:text-primary transition-colors">
                Shared
              </CardTitle>
              <Users className="h-5 w-5 text-muted-foreground group-hover:text-primary/70 transition-colors" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                {sharedEntries.length}
              </div>
              <p className="text-xs text-muted-foreground mt-1">Active shares</p>
            </CardContent>
          </Card>

          <Card className="group border-none shadow-lg hover:shadow-xl backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60 hover:bg-white/80 dark:hover:bg-zinc-900/70 transition-all duration-300 border-zinc-100 dark:border-zinc-800">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <CardTitle className="text-base font-medium text-foreground group-hover:text-primary transition-colors">
                Activity
              </CardTitle>
              <Activity className="h-5 w-5 text-muted-foreground group-hover:text-primary/70 transition-colors" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                24
              </div>
              <p className="text-xs text-muted-foreground mt-1">Actions this week</p>
            </CardContent>
          </Card>
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <Card className="border-none shadow-xl backdrop-blur-sm bg-white/60 dark:bg-zinc-900/50 lg:col-span-1">
            <CardHeader>
              <CardTitle className="text-xl font-medium">Recent Activity</CardTitle>
              <CardDescription>Your latest vault operations</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="group flex items-start gap-3 p-3 rounded-xl hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
                  <div className="h-2 w-2 rounded-full bg-gradient-to-r from-primary to-amber-500 mt-2 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-foreground group-hover:text-primary">New entry created</p>
                    <p className="text-xs text-muted-foreground mt-0.5">2 hours ago</p>
                  </div>
                </div>
                <div className="group flex items-start gap-3 p-3 rounded-xl hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
                  <div className="h-2 w-2 rounded-full bg-gradient-to-r from-primary to-amber-500 mt-2 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-foreground group-hover:text-primary">Entry shared with team</p>
                    <p className="text-xs text-muted-foreground mt-0.5">5 hours ago</p>
                  </div>
                </div>
                <div className="group flex items-start gap-3 p-3 rounded-xl hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
                  <div className="h-2 w-2 rounded-full bg-muted mt-2 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-foreground group-hover:text-primary">Password updated</p>
                    <p className="text-xs text-muted-foreground mt-0.5">1 day ago</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-none shadow-xl backdrop-blur-sm bg-white/60 dark:bg-zinc-900/50 lg:col-span-1">
            <CardHeader>
              <CardTitle className="text-xl font-medium">Your Contributions</CardTitle>
              <CardDescription>Activity over time</CardDescription>
            </CardHeader>
            <CardContent className="p-6">
              <ContributionGraph contributions={sampleData} />
            </CardContent>
          </Card>

          <GlobalSecurityInsight />
        </div>
      </div>
    </DashboardLayout>
  );
};


export default Index;
