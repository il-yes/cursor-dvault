import React from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";

export const TrustGroupsTab: React.FC = () => {
  const { trustGroups, setCreateTrustGroupOpen } = useC3ConfigurationStore();

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
          <div className="group-card" key={group.id}>
            <div className="gc-top">
              <div className="gc-name">{group.name}</div>
              <div className="gc-top-right">
                {group.status === "pending" && <span className="gc-ext-badge">C3 External</span>}
                <div className={`level-badge ${group.status === "active" ? "lvl-commit" : "lvl-read"}`}>
                  {group.status === "active" ? "Commit" : "Read"}
                </div>
              </div>
            </div>
            <div>
              <div className="gc-members-label">Members ({group.memberCount})</div>
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
                <span className="gc-add-member">+ Add member</span>
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
                <div className="icon-btn">✎</div>
                <div className="icon-btn del">✕</div>
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

