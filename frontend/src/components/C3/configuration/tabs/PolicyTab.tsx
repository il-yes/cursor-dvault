import React, { useState, useRef } from "react";
import { useC3ConfigurationStore, PolicyDraft } from "../store/useC3ConfigurationStore";

export const PolicyTab: React.FC = () => {
  const {
    draftPolicy,
    setDraftPolicy,
    isDirty,
    isSaving,
    saveFeedback,
    saveChannelConfig,
    discardChanges,
  } = useC3ConfigurationStore();

  const [activeSection, setActiveSection] = useState<"execution" | "receipts" | "retention" | "disputes" | "stellar">("execution");

  const containerRef = useRef<HTMLDivElement>(null);
  const executionRef = useRef<HTMLDivElement>(null);
  const receiptsRef = useRef<HTMLDivElement>(null);
  const retentionRef = useRef<HTMLDivElement>(null);
  const disputesRef = useRef<HTMLDivElement>(null);
  const stellarRef = useRef<HTMLDivElement>(null);

  const scrollToSection = (section: "execution" | "receipts" | "retention" | "disputes" | "stellar") => {
    setActiveSection(section);
    const refMap = {
      execution: executionRef,
      receipts: receiptsRef,
      retention: retentionRef,
      disputes: disputesRef,
      stellar: stellarRef,
    };

    const targetRef = refMap[section];
    if (targetRef.current && containerRef.current) {
      containerRef.current.scrollTo({
        top: targetRef.current.offsetTop - 16,
        behavior: "smooth",
      });
    }
  };

  const handleScroll = () => {
    if (!containerRef.current) return;
    const scrollPos = containerRef.current.scrollTop + 60;

    const sections = [
      { id: "execution" as const, ref: executionRef },
      { id: "receipts" as const, ref: receiptsRef },
      { id: "retention" as const, ref: retentionRef },
      { id: "disputes" as const, ref: disputesRef },
      { id: "stellar" as const, ref: stellarRef },
    ];

    for (let i = sections.length - 1; i >= 0; i--) {
      const sec = sections[i];
      if (sec.ref.current && sec.ref.current.offsetTop <= scrollPos) {
        setActiveSection(sec.id);
        break;
      }
    }
  };

  const handleToggle = (key: keyof PolicyDraft) => {
    setDraftPolicy({ [key]: !draftPolicy[key] });
  };

  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
      <div className="policy-body">
        {/* Left section navigation (190px width) */}
        <div className="policy-nav">
          <div
            className={`pnav-item ${activeSection === "execution" ? "active" : ""}`}
            onClick={() => scrollToSection("execution")}
          >
            <div className="pnav-dot" style={{ background: "#2563EB" }}></div>Execution
          </div>
          <div
            className={`pnav-item ${activeSection === "receipts" ? "active" : ""}`}
            onClick={() => scrollToSection("receipts")}
          >
            <div className="pnav-dot" style={{ background: "#C8922A" }}></div>Receipts
          </div>
          <div
            className={`pnav-item ${activeSection === "retention" ? "active" : ""}`}
            onClick={() => scrollToSection("retention")}
          >
            <div className="pnav-dot" style={{ background: "#059669" }}></div>Retention
          </div>
          <div
            className={`pnav-item ${activeSection === "disputes" ? "active" : ""}`}
            onClick={() => scrollToSection("disputes")}
          >
            <div className="pnav-dot" style={{ background: "#DC2626" }}></div>Disputes
          </div>
          <div
            className={`pnav-item ${activeSection === "stellar" ? "active" : ""}`}
            onClick={() => scrollToSection("stellar")}
          >
            <div className="pnav-dot" style={{ background: "#7C3AED" }}></div>Stellar
          </div>
        </div>

        {/* Right settings form with smooth scroll container */}
        <div className="policy-settings" ref={containerRef} onScroll={handleScroll} style={{ overflowY: "auto" }}>
          {/* EXECUTION */}
          <div className="policy-section" id="execution" ref={executionRef}>
            <div className="ps-title">Execution</div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Require all slots before completion</div>
                <div className="sr-desc">Thread cannot be marked complete unless every required slot has been committed. Partial completion is not allowed.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.requireAllSlots ? "on" : "off"}`}
                  onClick={() => handleToggle("requireAllSlots")}
                ></div>
              </div>
            </div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Allow parallel slot commits</div>
                <div className="sr-desc">When enabled, vaults can commit their slots without waiting for gate conditions. Gate ordering is ignored.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.allowParallelCommits ? "on" : "off"}`}
                  onClick={() => handleToggle("allowParallelCommits")}
                ></div>
              </div>
            </div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Maximum thread duration</div>
                <div className="sr-desc">Threads exceeding this duration are flagged as stale and vaults receive an overdue notification.</div>
              </div>
              <div className="sr-control">
                <select
                  className="policy-select"
                  value={draftPolicy.maxThreadDuration}
                  onChange={(e) => setDraftPolicy({ maxThreadDuration: e.target.value })}
                >
                  <option>7 days</option>
                  <option>14 days</option>
                  <option value="30 days">30 days</option>
                  <option>60 days</option>
                  <option>90 days</option>
                  <option>No limit</option>
                </select>
              </div>
            </div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Allow commit amendments</div>
                <div className="sr-desc">Vaults can resubmit a slot commit after it has been recorded. Each amendment is versioned and Stellar-anchored.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.allowAmendments ? "on toggle-gold" : "off"}`}
                  onClick={() => handleToggle("allowAmendments")}
                ></div>
              </div>
            </div>
          </div>

          {/* RECEIPTS */}
          <div className="policy-section" id="receipts" ref={receiptsRef}>
            <div className="ps-title">Receipts &amp; Acknowledgment</div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Require receipt acknowledgment</div>
                <div className="sr-desc">After a vault commits a slot, the next vault in the pipeline must explicitly acknowledge: received | processed | rejected.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.requireReceiptAck ? "on" : "off"}`}
                  onClick={() => handleToggle("requireReceiptAck")}
                ></div>
              </div>
            </div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Receipt timeout</div>
                <div className="sr-desc">If a vault does not acknowledge within this window, the thread is flagged as awaiting acknowledgment.</div>
              </div>
              <div className="sr-control">
                <select
                  className="policy-select"
                  value={draftPolicy.receiptTimeout}
                  onChange={(e) => setDraftPolicy({ receiptTimeout: e.target.value })}
                >
                  <option>24 hours</option>
                  <option>48 hours</option>
                  <option value="72 hours">72 hours</option>
                  <option>7 days</option>
                  <option>No timeout</option>
                </select>
              </div>
            </div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Allow rejection with reason code</div>
                <div className="sr-desc">Vaults can submit a structured reason code (policy_violation | decrypt_failed | schema_invalid | vault_suspended) alongside a free-text explanation.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.allowRejectionWithReason ? "on" : "off"}`}
                  onClick={() => handleToggle("allowRejectionWithReason")}
                ></div>
              </div>
            </div>
          </div>

          {/* RETENTION */}
          <div className="policy-section" id="retention" ref={retentionRef}>
            <div className="ps-title">Retention</div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Retain thread data for</div>
                <div className="sr-desc">How long completed or archived thread data is kept on Ankhora-hosted nodes before expiry. Stellar anchors are permanent regardless.</div>
              </div>
              <div className="sr-control">
                <select
                  className="policy-select"
                  value={draftPolicy.retainThreadData}
                  onChange={(e) => setDraftPolicy({ retainThreadData: e.target.value })}
                >
                  <option>90 days</option>
                  <option>180 days</option>
                  <option value="365 days">365 days</option>
                  <option>3 years</option>
                  <option>7 years</option>
                  <option>Indefinite</option>
                </select>
              </div>
            </div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Auto-archive on completion</div>
                <div className="sr-desc">Move threads to archive status automatically once all slots are committed and receipted.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.autoArchive ? "on toggle-gold" : "off"}`}
                  onClick={() => handleToggle("autoArchive")}
                ></div>
              </div>
            </div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Delete on archive</div>
                <div className="sr-desc">Permanently remove thread data when it enters archived state. Cannot be undone. Stellar anchors are unaffected.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.deleteOnArchive ? "on" : "off"}`}
                  onClick={() => handleToggle("deleteOnArchive")}
                ></div>
              </div>
            </div>
          </div>

          {/* DISPUTES */}
          <div className="policy-section" id="disputes" ref={disputesRef}>
            <div className="ps-title">Dispute Resolution</div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Dispute auto-escalation timeout</div>
                <div className="sr-desc">If a disputed thread is not resolved within this window, it is escalated to the designated escalation vault.</div>
              </div>
              <div className="sr-control">
                <select
                  className="policy-select"
                  value={draftPolicy.disputeTimeout}
                  onChange={(e) => setDraftPolicy({ disputeTimeout: e.target.value })}
                >
                  <option>3 days</option>
                  <option>7 days</option>
                  <option value="14 days">14 days</option>
                  <option>30 days</option>
                  <option>No auto-escalation</option>
                </select>
              </div>
            </div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Notify all vaults on dispute</div>
                <div className="sr-desc">When a rejection is submitted and the thread enters disputed state, all assigned vaults in the channel receive a notification.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.notifyAllOnDispute ? "on" : "off"}`}
                  onClick={() => handleToggle("notifyAllOnDispute")}
                ></div>
              </div>
            </div>
          </div>

          {/* STELLAR */}
          <div className="policy-section" id="stellar" ref={stellarRef}>
            <div className="ps-title">Stellar Anchoring</div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Anchor all commits</div>
                <div className="sr-desc">Every slot commit generates a Stellar transaction. Commit hash, vault ID, slot name, and CID are anchored immutably.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.anchorAllCommits ? "on" : "off"}`}
                  onClick={() => handleToggle("anchorAllCommits")}
                ></div>
                <span style={{ fontFamily: "'SF Mono', monospace", fontSize: "10px", color: "#22C55E", fontWeight: 600 }}>Required</span>
              </div>
            </div>

            <div className="setting-row active-row">
              <div className="sr-main">
                <div className="sr-label">Anchor all receipts</div>
                <div className="sr-desc">Every receipt acknowledgment (received | processed | rejected) is signed and anchored on Stellar alongside the commit it references.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.anchorAllReceipts ? "on" : "off"}`}
                  onClick={() => handleToggle("anchorAllReceipts")}
                ></div>
              </div>
            </div>

            <div className="setting-row">
              <div className="sr-main">
                <div className="sr-label">Anchor trust group changes</div>
                <div className="sr-desc">Additions and removals of TrustGroupMembers are anchored — creating a verifiable chain of custody for access changes.</div>
              </div>
              <div className="sr-control">
                <div
                  className={`toggle ${draftPolicy.anchorTrustGroupChanges ? "on toggle-gold" : "off"}`}
                  onClick={() => handleToggle("anchorTrustGroupChanges")}
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="policy-save-bar">
        <div className="psb-note">
          {saveFeedback ? (
            <span style={{ color: saveFeedback.includes("successfully") ? "#059669" : "#DC2626", fontWeight: 600 }}>
              {saveFeedback}
            </span>
          ) : (
            "Policy changes apply to new threads only. Active threads continue under the policy in effect at their creation time."
          )}
        </div>
        <button className="btn-discard" onClick={discardChanges} disabled={isSaving || !isDirty}>
          Discard
        </button>
        <button className="btn-save" onClick={saveChannelConfig} disabled={isSaving}>
          {isSaving ? "Saving..." : "Save Policy"}
        </button>
      </div>
    </div>
  );
};



