import React, { useState } from "react";
import { C3TrustGroupSummary } from "../domain/configuration";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";
import { Shield, Users, Smartphone, KeyRound, Radio, History, ArrowLeft, CheckCircle2, AlertTriangle, Edit3, Trash2, Check, X } from "lucide-react";

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
  const [isEditing, setIsEditing] = useState(false);
  const [editedName, setEditedName] = useState(trustGroup.name);
  const [editedDesc, setEditedDesc] = useState(trustGroup.description || "");
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  const { updateTrustGroup, deleteTrustGroup } = useC3ConfigurationStore();

  const handleSave = async () => {
    if (!editedName.trim()) return;
    setIsSaving(true);
    setSaveError(null);
    try {
      await updateTrustGroup({
        id: trustGroup.id,
        name: editedName.trim(),
        description: editedDesc.trim(),
      });
      setIsEditing(false);
    } catch (err: any) {
      setSaveError(err?.message || "Failed to update Trust Group.");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm(`Are you sure you want to delete Trust Group "${trustGroup.name}" (${trustGroup.id})?`)) {
      return;
    }
    setIsDeleting(true);
    try {
      await deleteTrustGroup(trustGroup.id);
      onBack();
    } catch (err: any) {
      setSaveError(err?.message || "Failed to delete Trust Group.");
      setIsDeleting(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Back Button & Header Banner */}
      <div className="flex items-center justify-between border border-border p-5 rounded-lg">
        <div className="flex items-center gap-4">
          <button
            onClick={onBack}
            className="p-2 rounded-md border border-border text-foreground hover:bg-muted transition-colors"
            title="Back to Trust Groups List"
          >
            <ArrowLeft className="w-4 h-4" />
          </button>

          {isEditing ? (
            <div className="flex items-center gap-2">
              <input
                type="text"
                value={editedName}
                onChange={(e) => setEditedName(e.target.value)}
                placeholder="Trust Group Name"
                className="px-3 py-1.5 border border-[#C8922A] text-foreground bg-transparent rounded text-sm outline-none focus:ring-1 focus:ring-[#C8922A]"
                disabled={isSaving}
                autoFocus
              />
              <button
                onClick={handleSave}
                disabled={isSaving || !editedName.trim()}
                className="p-1.5 bg-[#C8922A] text-white rounded hover:bg-[#b07e22] transition-colors disabled:opacity-50"
                title="Save Changes"
              >
                <Check className="w-4 h-4" />
              </button>
              <button
                onClick={() => {
                  setEditedName(trustGroup.name);
                  setEditedDesc(trustGroup.description || "");
                  setIsEditing(false);
                  setSaveError(null);
                }}
                disabled={isSaving}
                className="p-1.5 border border-border text-muted-foreground rounded hover:text-foreground transition-colors"
                title="Cancel Edit"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          ) : (
            <div>
              <div className="flex items-center gap-2">
                <Shield className="w-5 h-5 text-primary" />
                <h2 className="text-xl font-bold text-foreground">{trustGroup.name}</h2>
                <button
                  onClick={() => {
                    setEditedName(trustGroup.name);
                    setEditedDesc(trustGroup.description || "");
                    setIsEditing(true);
                  }}
                  className="p-1 text-muted-foreground hover:text-primary transition-colors"
                  title="Edit Trust Group Name"
                >
                  <Edit3 className="w-4 h-4" />
                </button>
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
              <p className="text-xs text-muted-foreground mt-1 font-mono">
                ID: <span className="text-primary font-semibold">{trustGroup.id}</span> • KEK Version: v{trustGroup.kekVersion}
              </p>
            </div>
          )}
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={handleDelete}
            disabled={isDeleting}
            className="flex items-center gap-1.5 px-3 py-1.5 bg-red-500/10 border border-red-500/30 text-red-400 hover:bg-red-500/20 text-xs font-medium rounded-md transition-colors disabled:opacity-50"
            title="Delete Trust Group"
          >
            <Trash2 className="w-3.5 h-3.5" />
            {isDeleting ? "Deleting..." : "Delete Group"}
          </button>
          <div className="text-right hidden sm:block">
            <div className="text-xs text-muted-foreground">Cryptographic Health</div>
            <div className="text-sm font-semibold text-emerald-400 flex items-center gap-1 justify-end">
              <CheckCircle2 className="w-4 h-4" /> Sovereign & Active
            </div>
          </div>
        </div>
      </div>

      {saveError && (
        <div className="p-3 bg-red-500/10 border border-red-500/30 text-red-400 rounded-md text-xs">
          {saveError}
        </div>
      )}

      {/* Sub Navigation Bar */}
      <div className="flex border-b border-border gap-1 overflow-x-auto">
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
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground hover:border-border"
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
          <div className="border border-border p-4 rounded-lg space-y-2">
            <div className="text-xs text-muted-foreground font-medium flex justify-between items-center">
              <span>Description</span>
            </div>
            {isEditing ? (
              <textarea
                value={editedDesc}
                onChange={(e) => setEditedDesc(e.target.value)}
                placeholder="Description"
                className="w-full h-20 p-2 border border-[#C8922A] text-foreground bg-transparent rounded text-xs outline-none"
              />
            ) : (
              <p className="text-xs text-foreground leading-relaxed">
                {trustGroup.description || "No detailed description provided."}
              </p>
            )}
          </div>
          <div className="border border-border p-4 rounded-lg space-y-2">
            <div className="text-xs text-muted-foreground font-medium">Identity Metrics</div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-muted-foreground">Total Members:</span>
              <span className="font-semibold text-foreground">{trustGroup.memberCount}</span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-muted-foreground">Registered Devices:</span>
              <span className="font-semibold text-foreground">{trustGroup.deviceCount}</span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-muted-foreground">Created:</span>
              <span className="font-mono text-foreground">
                {new Date(trustGroup.createdDate).toLocaleDateString()}
              </span>
            </div>
          </div>
          <div className="border border-border p-4 rounded-lg space-y-2">
            <div className="text-xs text-muted-foreground font-medium">Cryptographic KEK Metadata</div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-muted-foreground">Current KEK Version:</span>
              <span className="font-mono font-bold text-primary">v{trustGroup.kekVersion}</span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-muted-foreground">Active Envelopes:</span>
              <span className="font-semibold text-emerald-400">
                {trustGroup.cryptographicState.activeEnvelopes}
              </span>
            </div>
            <div className="flex justify-between items-center text-xs">
              <span className="text-muted-foreground">Last Rotation:</span>
              <span className="font-mono text-foreground">
                {new Date(trustGroup.cryptographicState.lastRotation).toLocaleDateString()}
              </span>
            </div>
          </div>
        </div>
      )}

      {activeSubTab === "members" && (
        <div className="border border-border rounded-lg overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-border text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                <th className="py-3 px-4">Identity</th>
                <th className="py-3 px-4">Role</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Devices</th>
                <th className="py-3 px-4">Joined</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-xs">
              {trustGroup.members.map((m) => (
                <tr key={m.id} className="hover:bg-muted/50 transition-colors">
                  <td className="py-3 px-4 font-medium text-foreground">
                    <div>{m.identityName}</div>
                    {m.email && <div className="text-[11px] text-muted-foreground">{m.email}</div>}
                  </td>
                  <td className="py-3 px-4 text-foreground capitalize">{m.role}</td>
                  <td className="py-3 px-4">
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                      {m.status}
                    </span>
                  </td>
                  <td className="py-3 px-4 text-foreground font-mono">{m.deviceCount}</td>
                  <td className="py-3 px-4 text-muted-foreground font-mono">
                    {new Date(m.joinedAt).toLocaleDateString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeSubTab === "devices" && (
        <div className="border border-border rounded-lg overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-border text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                <th className="py-3 px-4">Device</th>
                <th className="py-3 px-4">Owner Identity</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Key State</th>
                <th className="py-3 px-4">Last Seen</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border text-xs">
              {trustGroup.devices.map((d) => (
                <tr key={d.id} className="hover:bg-muted/50 transition-colors">
                  <td className="py-3 px-4 font-medium text-foreground">
                    <div>{d.deviceName}</div>
                    <div className="text-[10px] text-muted-foreground font-mono">{d.id}</div>
                  </td>
                  <td className="py-3 px-4 text-foreground font-mono">{d.identityId}</td>
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
                  <td className="py-3 px-4 text-muted-foreground font-mono">
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
          <div className="border border-border p-5 rounded-lg">
            <h3 className="text-sm font-semibold text-foreground mb-3 flex items-center gap-2">
              <KeyRound className="w-4 h-4 text-primary" /> Cryptographic Key Envelope Metadata
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="p-3 border border-border rounded-md">
                <div className="text-[11px] text-muted-foreground">KEK Version</div>
                <div className="text-lg font-bold text-primary font-mono">
                  v{trustGroup.kekVersion}
                </div>
              </div>
              <div className="p-3 border border-border rounded-md">
                <div className="text-[11px] text-muted-foreground">Total Envelopes</div>
                <div className="text-lg font-bold text-foreground font-mono">
                  {trustGroup.cryptographicState.envelopeCount}
                </div>
              </div>
              <div className="p-3 border border-border rounded-md">
                <div className="text-[11px] text-muted-foreground">Active Envelopes</div>
                <div className="text-lg font-bold text-emerald-400 font-mono">
                  {trustGroup.cryptographicState.activeEnvelopes}
                </div>
              </div>
              <div className="p-3 border border-border rounded-md">
                <div className="text-[11px] text-muted-foreground">Revoked Envelopes</div>
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
        <div className="border border-border p-5 rounded-lg space-y-3">
          <h3 className="text-sm font-semibold text-foreground">Associated Channels</h3>
          {trustGroup.associatedChannelIds.length > 0 ? (
            <div className="space-y-2">
              {trustGroup.associatedChannelIds.map((chId) => (
                <div
                  key={chId}
                  className="p-3 border border-border rounded-md flex items-center justify-between text-xs"
                >
                  <div className="flex items-center gap-2">
                    <Radio className="w-4 h-4 text-primary" />
                    <span className="font-medium text-foreground">{chId}</span>
                  </div>
                  <span className="text-[11px] text-emerald-400 font-mono bg-emerald-500/10 px-2 py-0.5 rounded">
                    Bound as Sovereign Access Boundary
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-xs text-muted-foreground p-4 text-center border border-dashed border-border rounded-md">
              No channels currently bound to this Trust Group.
            </div>
          )}
        </div>
      )}

      {activeSubTab === "activity" && (
        <div className="border border-border p-5 rounded-lg space-y-3">
          <h3 className="text-sm font-semibold text-foreground">Audit Trail</h3>
          <div className="space-y-2 text-xs">
            <div className="p-3 border border-border rounded-md flex items-center justify-between">
              <span className="text-foreground">KEK Rotation to v{trustGroup.kekVersion} completed</span>
              <span className="text-muted-foreground font-mono">
                {new Date(trustGroup.cryptographicState.lastRotation).toLocaleString()}
              </span>
            </div>
            <div className="p-3 border border-border rounded-md flex items-center justify-between">
              <span className="text-foreground">Trust Group established with ID {trustGroup.id}</span>
              <span className="text-muted-foreground font-mono">
                {new Date(trustGroup.createdDate).toLocaleString()}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
