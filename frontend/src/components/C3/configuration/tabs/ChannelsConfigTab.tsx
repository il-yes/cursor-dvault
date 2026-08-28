import React, { useState } from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";

export const ChannelsConfigTab: React.FC = () => {
  const { channels } = useC3ConfigurationStore();
  const [isAddOpen, setIsAddOpen] = useState(false);

  return (
    <div className="tab-body">
      <div className="tab-desc">
        Ordered commit steps for this channel. Each slot is committed by its assigned vault in sequence. Drag to reorder. Gate conditions are validated at runtime by AFP and are not retroactively applied to running threads.
      </div>

      <div className="slot-area">
        {/* Column headers */}
        <div className="sg col-hdrs">
          <div className="ch-lbl"></div>
          <div className="ch-lbl">#</div>
          <div className="ch-lbl" style={{ paddingLeft: "14px" }}>Slot name</div>
          <div className="ch-lbl">Assigned vault</div>
          <div className="ch-lbl">Gate condition</div>
          <div className="ch-lbl" style={{ textAlign: "center" }}>Req.</div>
          <div className="ch-lbl"></div>
        </div>

        {/* Slot 1 */}
        <div className="sg slot-row">
          <div className="sr-drag">⠿</div>
          <div className="sr-num">1</div>
          <div className="sr-name"><span className="slot-tag">contract_draft</span></div>
          <div className="sr-vault">
            <span className="vault-chip vc-pur">
              <span className="vcdot" style={{ background: "#7C3AED" }}></span>vault_legal
            </span>
          </div>
          <div className="sr-gate"><span className="gate-none">No gate — first slot</span></div>
          <div className="sr-req"><span className="req-dot"></span></div>
          <div className="sr-acts">
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* Slot 2 */}
        <div className="sg slot-row">
          <div className="sr-drag">⠿</div>
          <div className="sr-num">2</div>
          <div className="sr-name"><span className="slot-tag">financial_clearance</span></div>
          <div className="sr-vault">
            <span className="vault-chip vc-blu">
              <span className="vcdot" style={{ background: "#2563EB" }}></span>vault_finance
            </span>
          </div>
          <div className="sr-gate">
            <span className="gate-chip"><span className="gate-arr">→</span>contract_draft</span>
          </div>
          <div className="sr-req"><span className="req-dot"></span></div>
          <div className="sr-acts">
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* Slot 3 */}
        <div className="sg slot-row">
          <div className="sr-drag">⠿</div>
          <div className="sr-num">3</div>
          <div className="sr-name"><span className="slot-tag">executive_signature</span></div>
          <div className="sr-vault">
            <span className="vault-chip vc-dar">
              <span className="vcdot" style={{ background: "#444" }}></span>vault_direction
            </span>
          </div>
          <div className="sr-gate">
            <span className="gate-chip"><span className="gate-arr">→</span>financial_clearance</span>
          </div>
          <div className="sr-req"><span className="req-dot"></span></div>
          <div className="sr-acts">
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* Add Slot Row */}
        {isAddOpen ? (
          <div className="add-row">
            <div className="add-inner">
              <div></div>
              <div className="add-num">4</div>
              <div style={{ padding: "0 6px" }}>
                <input className="slot-input" type="text" placeholder="new_slot_name" />
              </div>
              <div style={{ padding: "0 6px" }}>
                <select className="slot-select">
                  <option value="">Select vault…</option>
                  <option>vault_finance</option>
                  <option>vault_legal</option>
                  <option>vault_direction</option>
                  <option>vault_treasury</option>
                  <option>vault_hr</option>
                  <option>vault_ops</option>
                  <option>vault_compliance</option>
                </select>
              </div>
              <div style={{ padding: "0 6px" }}>
                <select className="slot-select">
                  <option>No gate (first slot)</option>
                  <option>→ after: contract_draft</option>
                  <option>→ after: financial_clearance</option>
                  <option defaultValue="executive_signature">→ after: executive_signature</option>
                </select>
              </div>
              <div className="add-req">
                <input type="checkbox" defaultChecked style={{ accentColor: "#C8922A", width: "13px", height: "13px" }} />
              </div>
              <div className="add-actions">
                <button className="btn-confirm" onClick={() => setIsAddOpen(false)}>+ Add</button>
                <button className="btn-cancel-add" onClick={() => setIsAddOpen(false)}>✕</button>
              </div>
            </div>
          </div>
        ) : (
          <div className="add-trigger" onClick={() => setIsAddOpen(true)}>
            <span className="add-trigger-plus">+</span> Add another slot
          </div>
        )}
      </div>

      <div className="dirty-bar">
        <span className="db-icon">⚠</span>
        <div className="db-msg">
          <strong>Unsaved changes</strong> — slot modifications take effect on the next thread instantiation. Active threads using this channel are not affected.
        </div>
        <button className="btn-discard">Discard</button>
        <button className="btn-save">Save Changes</button>
      </div>
    </div>
  );
};

