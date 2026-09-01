import React from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";
import { Shield, Radio, CheckSquare, ShieldCheck, KeyRound, ArrowRight } from "lucide-react";

export const OverviewTab: React.FC = () => {
  const {
    trustGroups,
    channels,
    metrics,
    securityState,
    setActiveTab,
  } = useC3ConfigurationStore();

  const activeTrustGroupsCount = trustGroups.filter((g) => g.status === "active").length;
  const pendingTrustGroupsCount = trustGroups.filter((g) => g.status === "pending").length;

  return (
    <div className="space-y-6">
      {/* Top Banner Executive Status */}
      <div className="border border-border p-6 rounded-xl flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 text-xs font-semibold text-primary uppercase tracking-wider mb-1">
            <ShieldCheck className="w-4 h-4 text-emerald-400" /> C3 Collaboration Infrastructure
          </div>
          <h2 className="text-2xl font-bold text-foreground">Executive Configuration Status</h2>
          <p className="text-xs text-muted-foreground mt-1 max-w-xl">
            Authoritative status overview of sovereign Trust Groups, governance Channels, collaboration primitives, and cryptographic keyring state.
          </p>
        </div>
        <div className="flex items-center gap-2 border border-border px-4 py-2 rounded-lg">
          <span className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-pulse" />
          <span className="text-xs font-semibold text-foreground">C3 System Active & Healthy</span>
        </div>
      </div>

      {/* Grid Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Trust Groups Card */}
        <div
          onClick={() => setActiveTab("trust-groups")}
          className="border border-border hover:border-primary p-5 rounded-xl cursor-pointer transition-all hover:scale-[1.01] group"
        >
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Trust Groups
            </span>
            <div className="p-2 bg-primary/10 text-primary rounded-lg group-hover:bg-primary group-hover:text-white transition-colors">
              <Shield className="w-4 h-4" />
            </div>
          </div>
          <div className="text-2xl font-bold text-foreground">{trustGroups.length}</div>
          <div className="flex items-center gap-3 mt-2 text-xs text-muted-foreground">
            <span className="text-emerald-400 font-medium">{activeTrustGroupsCount} Active</span>
            <span>•</span>
            <span className="text-amber-400 font-medium">{pendingTrustGroupsCount} Pending</span>
          </div>
        </div>

        {/* Channels Card */}
        <div
          onClick={() => setActiveTab("channels")}
          className="border border-border hover:border-primary p-5 rounded-xl cursor-pointer transition-all hover:scale-[1.01] group"
        >
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Channels
            </span>
            <div className="p-2 bg-purple-500/10 text-purple-400 rounded-lg group-hover:bg-purple-500 group-hover:text-white transition-colors">
              <Radio className="w-4 h-4" />
            </div>
          </div>
          <div className="text-2xl font-bold text-foreground">{channels.length}</div>
          <div className="text-xs text-muted-foreground mt-2">
            Governance & Communication Surface
          </div>
        </div>

        {/* Collaboration Primitives Card */}
        <div
          onClick={() => setActiveTab("collaboration")}
          className="border border-border hover:border-primary p-5 rounded-xl cursor-pointer transition-all hover:scale-[1.01] group"
        >
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Collaboration
            </span>
            <div className="p-2 bg-amber-500/10 text-amber-400 rounded-lg group-hover:bg-amber-500 group-hover:text-white transition-colors">
              <CheckSquare className="w-4 h-4" />
            </div>
          </div>
          <div className="text-2xl font-bold text-foreground">{metrics.totalShares}</div>
          <div className="flex items-center gap-2 mt-2 text-xs text-muted-foreground">
            <span className="text-amber-400 font-medium">{metrics.pendingApprovals} Pending</span>
            <span>•</span>
            <span className="text-blue-400 font-medium">{metrics.activeTransfers} Transfers</span>
          </div>
        </div>

        {/* Security & Sovereignty Card */}
        <div
          onClick={() => setActiveTab("security")}
          className="border border-border hover:border-primary p-5 rounded-xl cursor-pointer transition-all hover:scale-[1.01] group"
        >
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Security & KEK
            </span>
            <div className="p-2 bg-emerald-500/10 text-emerald-400 rounded-lg group-hover:bg-emerald-500 group-hover:text-white transition-colors">
              <KeyRound className="w-4 h-4" />
            </div>
          </div>
          <div className="text-2xl font-bold font-mono text-primary">
            v{securityState.currentKekVersion}
          </div>
          <div className="text-xs text-muted-foreground mt-2">
            {securityState.activeEnvelopes} Active Envelopes
          </div>
        </div>
      </div>

      {/* Security Health Summary Bar */}
      <div className="border border-border p-5 rounded-xl space-y-3">
        <h3 className="text-sm font-semibold text-foreground flex items-center justify-between">
          <span>Security & Sovereign Keyring Summary</span>
          <button
            onClick={() => setActiveTab("security")}
            className="text-xs text-primary hover:underline flex items-center gap-1 font-normal"
          >
            Full Security View <ArrowRight className="w-3 h-3" />
          </button>
        </h3>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-3 text-xs">
          <div className="p-3 border border-border rounded-lg flex items-center justify-between">
            <span className="text-muted-foreground">Device Identity:</span>
            <span className="font-semibold text-emerald-400 capitalize">
              {securityState.deviceIdentityStatus}
            </span>
          </div>
          <div className="p-3 border border-border rounded-lg flex items-center justify-between">
            <span className="text-muted-foreground">Cryptographic State:</span>
            <span className="font-semibold text-emerald-400 capitalize">
              {securityState.sovereignKeyringStatus}
            </span>
          </div>
          <div className="p-3 border border-border rounded-lg flex items-center justify-between">
            <span className="text-muted-foreground">Trust Group Keys:</span>
            <span className="font-semibold text-emerald-400 capitalize">
              {securityState.trustGroupKeysState}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};
