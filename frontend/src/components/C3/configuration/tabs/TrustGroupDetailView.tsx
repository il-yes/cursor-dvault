import React, { useState } from "react";
import { C3TrustGroupSummary } from "../domain/configuration";
import { Shield, Users, Smartphone, KeyRound, Radio, History, ArrowLeft, CheckCircle2, Clock, AlertTriangle } from "lucide-react";

interface TrustGroupDetailViewProps {
  trustGroup: C3TrustGroupSummary;
  onBack: () => void;
}

type DetailTab = "overview" | "members" | "devices" | "crypto" | "channels" | "activity";

export const TrustGroupDetailView: React.FC<TrustGroupDetailViewProps> = ({
  trustGroup,
  onBack,
}) => {
  const [activeSubTab, setActiveSubTab] = useState<DetailTab>("overview");

  return (
    <div className="space-y-6">
      {/* Back Button & Header Banner */}
      <div className="flex items-center justify-between bg-[#161B22] border border-[#30363D] p-5 rounded-lg">
        <div className="flex items-center gap-4">
          <button
            onClick={onBack}
            className="p-2 rounded-md bg-[#21262D] border border-[#30363D] text-[#C9D1D9] hover:text-white hover:bg-[#30363D] transition-colors"
            title="Back to Trust Groups List"
          >
            <ArrowLeft className="w-4 h-4" />
          </button>
          <div>
            <div className="flex items-center gap-2">
              <Shield className="w-5 h-5 text-[#58A6FF]" />
              <h2 className="text-xl font-bold text-[#F0F6FC]">{trustGroup.name}</h2>
              <span
                className={`text-xs px-2.5 py-0.5 rounded-full font-medium border ${
                  trustGroup.status === "active"
                    ? "bg-emerald-500/10 border-emerald-500/30 text-emerald-400"
                    : "bg-amber-500/10 border-amber-500/30 text-amber-400"
                }`}
              >
                {trustGroup.status.toUpperCase()}
              </span>
            </div>
            <p className="text-xs text-[#8B949E] mt-1 font-mono">
              ID: {trustGroup.id} • KEK Version: v{trustGroup.kekVersion}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="text-right hidden sm:block">
            <div className="text-xs text-[#8B949E]">Cryptographic Envelope Health</div>
            <div className="text-sm font-semibold text-emerald-400 flex items-center gap-1 justify-end">
              <CheckCircle2 className="w-4 h-4" /> Sovereign & Active
            </div>
          </div>
        </div>
      </div>

      {/* Sub Navigation Bar */}
      <div className="flex border-b border-[#30363D] gap-1 overflow-x-auto">
        {[
          { id: "overview", label: "Overview", icon: Shield },
          { id: "members", label: `Members (${trustGroup.members.length})`, icon: Users },
          { id: "devices", label: `Devices (${trustGroup.devices.length})`, icon: Smartphone },
          { id: "crypto", label: "Cryptographic State", icon: KeyRound },
          { id: "channels", label: `Channels (${trustGroup.associatedChannelIds.length})`, icon: Radio },
          { id: "activity", label: "Activity Audit", icon: History },
        ].map((tab) => {
          const Icon = tab.icon;
          const isActive = activeSubTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveSubTab(tab.id as DetailTab)}
              className={`flex items-center gap-2 px-4 py-2.5 text-xs font-medium border-b-2 transition-colors whitespace-nowrap ${
                isActive
                  ? "border-[#58A6FF] text-[#58A6FF] bg-[#161B22]/50"
                  : "border-transparent text-[#8B949E] hover:text-[#C9D1D9] hover:border-[#30363D]"
              }`}
            >
              <Icon className="w-3.5 h-3.5" />
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Tab Contents */}
      {activeSubTab === "overview" && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-[#161B22] border border-[#30363D] p-4 rounded-lg space-y-2">
            <div className="text-xs text-[#8B949E] font-medium">Description</div>
            <p className="text-xs text-[#C9D1D9] leading-relaxed">
              {trustGroup.description || "No detailed description provided."}
            </p>
          </div>
          <div className="bg-[#161B22] border border-[#30363D] p-4 rounded-lg space-y-2">
            <div className="text-xs text-[#8B949E] font-medium">Identity Metrics</div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-[#8B949E]">Total Members:</span>
              <span className="font-semibold text-[#F0F6FC]">{trustGroup.memberCount}</span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-[#8B949E]">Registered Devices:</span>
              <span className="font-semibold text-[#F0F6FC]">{trustGroup.deviceCount}</span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-[#8B949E]">Created:</span>
              <span className="font-mono text-[#C9D1D9]">
                {new Date(trustGroup.createdDate).toLocaleDateString()}
              </span>
            </div>
          </div>
          <div className="bg-[#161B22] border border-[#30363D] p-4 rounded-lg space-y-2">
            <div className="text-xs text-[#8B949E] font-medium">Cryptographic KEK Metadata</div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-[#8B949E]">Current KEK Version:</span>
              <span className="font-mono font-bold text-[#58A6FF]">v{trustGroup.kekVersion}</span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-[#8B949E]">Active Envelopes:</span>
              <span className="font-semibold text-emerald-400">
                {trustGroup.cryptographicState.activeEnvelopes}
              </span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-[#8B949E]">Last Rotation:</span>
              <span className="font-mono text-[#C9D1D9]">
                {new Date(trustGroup.cryptographicState.lastRotation).toLocaleDateString()}
              </span>
            </div>
          </div>
        </div>
      )}

      {activeSubTab === "members" && (
        <div className="bg-[#161B22] border border-[#30363D] rounded-lg overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-[#30363D] bg-[#0D1117] text-[11px] font-semibold text-[#8B949E] uppercase tracking-wider">
                <th className="py-3 px-4">Identity</th>
                <th className="py-3 px-4">Role</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Devices</th>
                <th className="py-3 px-4">Joined</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#30363D] text-xs">
              {trustGroup.members.map((m) => (
                <tr key={m.id} className="hover:bg-[#21262D]/50 transition-colors">
                  <td className="py-3 px-4 font-medium text-[#F0F6FC]">
                    <div>{m.identityName}</div>
                    {m.email && <div className="text-[11px] text-[#8B949E]">{m.email}</div>}
                  </td>
                  <td className="py-3 px-4 text-[#C9D1D9] capitalize">{m.role}</td>
                  <td className="py-3 px-4">
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                      {m.status}
                    </span>
                  </td>
                  <td className="py-3 px-4 text-[#C9D1D9] font-mono">{m.deviceCount}</td>
                  <td className="py-3 px-4 text-[#8B949E] font-mono">
                    {new Date(m.joinedAt).toLocaleDateString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeSubTab === "devices" && (
        <div className="bg-[#161B22] border border-[#30363D] rounded-lg overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-[#30363D] bg-[#0D1117] text-[11px] font-semibold text-[#8B949E] uppercase tracking-wider">
                <th className="py-3 px-4">Device</th>
                <th className="py-3 px-4">Owner Identity</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Key State</th>
                <th className="py-3 px-4">Last Seen</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#30363D] text-xs">
              {trustGroup.devices.map((d) => (
                <tr key={d.id} className="hover:bg-[#21262D]/50 transition-colors">
                  <td className="py-3 px-4 font-medium text-[#F0F6FC]">
                    <div>{d.deviceName}</div>
                    <div className="text-[10px] text-[#8B949E] font-mono">{d.id}</div>
                  </td>
                  <td className="py-3 px-4 text-[#C9D1D9] font-mono">{d.identityId}</td>
                  <td className="py-3 px-4">
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                      {d.status}
                    </span>
                  </td>
                  <td className="py-3 px-4">
                    <span
                      className={`inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-medium border ${
                        d.keyState === "valid"
                          ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
                          : "bg-amber-500/10 text-amber-400 border-amber-500/20"
                      }`}
                    >
                      {d.keyState === "valid" ? <CheckCircle2 className="w-3 h-3" /> : <AlertTriangle className="w-3 h-3" />}
                      {d.keyState}
                    </span>
                  </td>
                  <td className="py-3 px-4 text-[#8B949E] font-mono">
                    {new Date(d.lastSeenAt).toLocaleString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeSubTab === "crypto" && (
        <div className="space-y-4">
          <div className="bg-[#161B22] border border-[#30363D] p-5 rounded-lg">
            <h3 className="text-sm font-semibold text-[#F0F6FC] mb-3 flex items-center gap-2">
              <KeyRound className="w-4 h-4 text-[#58A6FF]" /> Cryptographic Key Envelope Metadata
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md">
                <div className="text-[11px] text-[#8B949E]">KEK Version</div>
                <div className="text-lg font-bold text-[#58A6FF] font-mono">
                  v{trustGroup.cryptographicState.kekVersion}
                </div>
              </div>
              <div className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md">
                <div className="text-[11px] text-[#8B949E]">Total Envelopes</div>
                <div className="text-lg font-bold text-[#F0F6FC] font-mono">
                  {trustGroup.cryptographicState.envelopeCount}
                </div>
              </div>
              <div className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md">
                <div className="text-[11px] text-[#8B949E]">Active Envelopes</div>
                <div className="text-lg font-bold text-emerald-400 font-mono">
                  {trustGroup.cryptographicState.activeEnvelopes}
                </div>
              </div>
              <div className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md">
                <div className="text-[11px] text-[#8B949E]">Revoked Envelopes</div>
                <div className="text-lg font-bold text-red-400 font-mono">
                  {trustGroup.cryptographicState.revokedEnvelopes}
                </div>
              </div>
            </div>
          </div>

          <div className="bg-amber-500/10 border border-amber-500/30 p-4 rounded-lg text-amber-300 text-xs flex items-start gap-3">
            <Shield className="w-5 h-5 shrink-0 mt-0.5" />
            <div>
              <div className="font-semibold mb-1">Sovereign Cryptographic Principle</div>
              <div>
                Raw secret materials (DEK/KEK bytes, private keys, wrapped payloads, and device seed material) remain strictly isolated within local sovereign hardware and keyring boundaries. Only non-sensitive metadata and envelope status are rendered.
              </div>
            </div>
          </div>
        </div>
      )}

      {activeSubTab === "channels" && (
        <div className="bg-[#161B22] border border-[#30363D] p-5 rounded-lg space-y-3">
          <h3 className="text-sm font-semibold text-[#F0F6FC]">Associated Channels</h3>
          {trustGroup.associatedChannelIds.length > 0 ? (
            <div className="space-y-2">
              {trustGroup.associatedChannelIds.map((chId) => (
                <div
                  key={chId}
                  className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md flex items-center justify-between text-xs"
                >
                  <div className="flex items-center gap-2">
                    <Radio className="w-4 h-4 text-[#58A6FF]" />
                    <span className="font-medium text-[#F0F6FC]">{chId}</span>
                  </div>
                  <span className="text-[11px] text-emerald-400 font-mono bg-emerald-500/10 px-2 py-0.5 rounded">
                    Bound as Sovereign Access Boundary
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-xs text-[#8B949E] p-4 text-center border border-dashed border-[#30363D] rounded-md">
              No channels currently bound to this Trust Group.
            </div>
          )}
        </div>
      )}

      {activeSubTab === "activity" && (
        <div className="bg-[#161B22] border border-[#30363D] p-5 rounded-lg space-y-3">
          <h3 className="text-sm font-semibold text-[#F0F6FC]">Audit Trail</h3>
          <div className="space-y-2 text-xs">
            <div className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md flex items-center justify-between">
              <span className="text-[#C9D1D9]">KEK Rotation to v{trustGroup.kekVersion} completed</span>
              <span className="text-[#8B949E] font-mono">
                {new Date(trustGroup.cryptographicState.lastRotation).toLocaleString()}
              </span>
            </div>
            <div className="p-3 bg-[#0D1117] border border-[#30363D] rounded-md flex items-center justify-between">
              <span className="text-[#C9D1D9]">Trust Group established with ID {trustGroup.id}</span>
              <span className="text-[#8B949E] font-mono">
                {new Date(trustGroup.createdDate).toLocaleString()}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
