import React, { useState } from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";
import { addTrustGroupMember } from "@/services/api";
import { TrustGroupDetailView } from "./TrustGroupDetailView";

export const TrustGroupsTab: React.FC = () => {
  const {
    trustGroups,
    selectedTrustGroupId,
    selectTrustGroup,
    deleteTrustGroup,
    setCreateTrustGroupOpen,
    originalChannel,
  } = useC3ConfigurationStore();

  const activeChannelId = useC3ChannelStore((state) => state.activeChannelId);
  const [addingMemberGroupId, setAddingMemberGroupId] = useState<string | null>(null);
  const [newMemberName, setNewMemberName] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [memberError, setMemberError] = useState<string | null>(null);
  const [deletingGroupId, setDeletingGroupId] = useState<string | null>(null);

  // If a Trust Group is selected, render its detail view/editor
  const selectedGroup = trustGroups.find((tg) => tg.id === selectedTrustGroupId);
  if (selectedTrustGroupId && selectedGroup) {
    return (
      <div className="trust-body">
        <TrustGroupDetailView
          trustGroup={selectedGroup}
          onBack={() => selectTrustGroup(null)}
        />
      </div>
    );
  }

  const handleAddMember = async (groupId: string) => {
    if (!newMemberName.trim()) return;

    setIsSubmitting(true);
    setMemberError(null);

    const memberId = newMemberName.trim();
    const targetChannelId = originalChannel?.id || activeChannelId || "contract-execution";

    try {
      const res = await addTrustGroupMember(groupId, targetChannelId, memberId);

      // Backend succeeded -> update UI store state directly
      const newMemberObj = {
        id: res?.member?.id || res?.MemberID || `usr_${Math.random().toString(36).substring(2, 8)}`,
        identityName: memberId,
        role: "member" as const,
        status: "active" as const,
        deviceCount: 1,
        joinedAt: new Date().toISOString(),
      };

      useC3ConfigurationStore.setState((state) => ({
        trustGroups: state.trustGroups.map((group) => {
          if (group.id !== groupId) return group;
          const updatedMembers = [...(group.members || []), newMemberObj];
          return {
            ...group,
            memberCount: updatedMembers.length,
            members: updatedMembers,
          };
        }),
      }));

      setNewMemberName("");
      setAddingMemberGroupId(null);
    } catch (err: any) {
      setMemberError(err?.message || "Failed to add member to Trust Group.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteGroup = async (e: React.MouseEvent, groupId: string, groupName: string) => {
    e.stopPropagation();
    if (!window.confirm(`Are you sure you want to delete Trust Group "${groupName}" (${groupId})?`)) {
      return;
    }
    setDeletingGroupId(groupId);
    try {
      await deleteTrustGroup(groupId);
    } catch (err: any) {
      alert(`Failed to delete Trust Group: ${err?.message || err}`);
    } finally {
      setDeletingGroupId(null);
    }
  };

  return (
    <div className="trust-body">
      {/* ── TRUST RULES ── */}
      <div className="section-hdr">
        <div className="section-hdr-left">
          <div className="section-title">Trust Rules</div>
          <div className="section-desc">
            Explicit access grants between vaults for this channel. Defines who can read, observe, or commit beyond their assigned slot. Evaluated at the AFP boundary.
          </div>
        </div>
        <button className="btn-add-sm">+ Add Rule</button>
      </div>

      <div className="rules-list">
        {/* Rule 1 */}
        <div className="rule-row">
          <span className="vault-chip vc-pur">
            <span className="vcdot" style={{ background: "#7C3AED" }}></span>
            vault_legal
          </span>
          <span className="rule-arrow">→</span>
          <span className="rule-grants">grants</span>
          <span className="vault-chip vc-blu">
            <span className="vcdot" style={{ background: "#2563EB" }}></span>
            vault_finance
          </span>
          <div className="rule-spacer"></div>
          <div className="scope-chip">
            <code>contract_draft</code>
          </div>
          <div className="level-badge lvl-read">Read</div>
          <div className="rule-acts">
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* Rule 2 */}
        <div className="rule-row">
          <span className="vault-chip vc-dar">
            <span className="vcdot" style={{ background: "#444" }}></span>
            vault_direction
          </span>
          <span className="rule-arrow">→</span>
          <span className="rule-grants">grants</span>
          <span className="vault-chip vc-ext">
            <span className="vcdot" style={{ background: "#C8922A" }}></span>
            vault:supplier-x
          </span>
          <div className="rule-spacer"></div>
          <div className="scope-chip wildcard">
            <code>* all slots</code>
          </div>
          <div className="level-badge lvl-observe">Observe</div>
          <div className="rule-acts">
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>
      </div>

      <hr className="section-divider" style={{ marginTop: "18px" }} />

      {/* ── TRUST GROUPS ── */}
      <div className="section-hdr" style={{ paddingTop: "16px" }}>
        <div className="section-hdr-left">
          <div className="section-title">Trust Groups</div>
          <div className="section-desc">
            Named collections of vaults sharing a trust relationship. Assign trust once to a group instead of per-vault. TrustGroupMembers are resolved at thread instantiation.
          </div>
        </div>
        <button className="btn-add-sm" onClick={() => setCreateTrustGroupOpen(true)}>
          + New Group
        </button>
      </div>

      <div className="groups-grid">
        {trustGroups.map((group) => (
          <div
            className="group-card"
            key={group.id}
            onClick={() => selectTrustGroup(group.id)}
            style={{ cursor: "pointer" }}
          >
            <div className="gc-top">
              <div className="gc-name">
                {group.name}
                <div style={{ fontSize: "10px", color: "#8B949E", fontFamily: "monospace", marginTop: "2px" }}>
                  {group.id}
                </div>
              </div>
              <div className="gc-top-right">
                {group.status === "pending" && <span className="gc-ext-badge">C3 External</span>}
                <div className={`level-badge ${group.status === "active" ? "lvl-commit" : "lvl-read"}`}>
                  {group.status === "active" ? "Commit" : "Read"}
                </div>
              </div>
            </div>
            <div>
              <div className="gc-members-label">Members ({group.members?.length || 0})</div>
              <div className="gc-members">
                {group.members && group.members.length > 0 ? (
                  group.members.map((member) => (
                    <span className="vault-chip vc-blu" key={member.id}>
                      <span className="vcdot" style={{ background: "#2563EB" }}></span>
                      {member.identityName}
                    </span>
                  ))
                ) : (
                  <span className="vault-chip vc-dar">
                    <span className="vcdot" style={{ background: "#444" }}></span>
                    {group.name}
                  </span>
                )}

                {addingMemberGroupId === group.id ? (
                  <div
                    style={{ display: "inline-flex", flexDirection: "column", gap: "2px" }}
                    onClick={(e) => e.stopPropagation()}
                  >
                    <div style={{ display: "inline-flex", alignItems: "center", gap: "4px" }}>
                      <input
                        type="text"
                        placeholder="vault_identity"
                        value={newMemberName}
                        onChange={(e) => setNewMemberName(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter" && !isSubmitting) handleAddMember(group.id);
                          if (e.key === "Escape") setAddingMemberGroupId(null);
                        }}
                        disabled={isSubmitting}
                        autoFocus
                        style={{
                          padding: "2px 6px",
                          fontSize: "11px",
                          borderRadius: "4px",
                          border: "1.5px solid #C8922A",
                          outline: "none",
                          width: "120px",
                        }}
                      />
                      <button
                        onClick={() => handleAddMember(group.id)}
                        disabled={isSubmitting}
                        style={{
                          padding: "2px 7px",
                          fontSize: "10px",
                          background: isSubmitting ? "#ccc" : "#C8922A",
                          color: "#fff",
                          border: "none",
                          borderRadius: "4px",
                          cursor: isSubmitting ? "wait" : "pointer",
                          fontWeight: 600,
                        }}
                      >
                        {isSubmitting ? "..." : "+"}
                      </button>
                      <button
                        onClick={() => {
                          setAddingMemberGroupId(null);
                          setMemberError(null);
                        }}
                        disabled={isSubmitting}
                        style={{
                          padding: "2px 5px",
                          fontSize: "10px",
                          background: "#f5f5f5",
                          color: "#888",
                          border: "1px solid #e5e5e5",
                          borderRadius: "4px",
                          cursor: "pointer",
                        }}
                      >
                        ✕
                      </button>
                    </div>
                    {memberError && (
                      <span style={{ fontSize: "10px", color: "#DC2626", fontWeight: 500 }}>
                        {memberError}
                      </span>
                    )}
                  </div>
                ) : (
                  <span
                    className="gc-add-member"
                    onClick={(e) => {
                      e.stopPropagation();
                      setAddingMemberGroupId(group.id);
                      setNewMemberName("");
                      setMemberError(null);
                    }}
                  >
                    + Add member
                  </span>
                )}
              </div>
            </div>
            <div className="gc-meta-row">
              <div className="gc-meta-item">
                <span style={{ color: "#ccc" }}>scope:</span>
                <span className="gc-meta-val">* all slots</span>
              </div>
              <div className="gc-meta-item">
                <span style={{ color: "#ccc" }}>kek:</span>
                <span className="gc-meta-val">v{group.kekVersion}</span>
              </div>
              <div className="gc-meta-sp"></div>
              <div className="gc-acts">
                <div
                  className="icon-btn"
                  title="Configure / Edit Trust Group"
                  onClick={(e) => {
                    e.stopPropagation();
                    selectTrustGroup(group.id);
                  }}
                >
                  ✎
                </div>
                <div
                  className="icon-btn del"
                  title="Delete Trust Group"
                  onClick={(e) => handleDeleteGroup(e, group.id, group.name)}
                >
                  {deletingGroupId === group.id ? "..." : "✕"}
                </div>
              </div>
            </div>
          </div>
        ))}

        {/* New Group Trigger */}
        <div className="new-group-row" style={{ display: "flex", alignItems: "stretch" }}>
          <button className="new-group-btn" onClick={() => setCreateTrustGroupOpen(true)}>
            + New Trust Group
          </button>
        </div>
      </div>
    </div>
  );
};
