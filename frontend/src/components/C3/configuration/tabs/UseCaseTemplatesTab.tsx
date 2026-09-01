import React, { useState } from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";

export const UseCaseTemplatesTab: React.FC = () => {
  const { channels } = useC3ConfigurationStore();
  const [isAddOpen, setIsAddOpen] = useState(false);

  return (
    <div className="tab-body">
      <div className="tab-desc">
        Channel-level property definitions. These become the available metadata fields when instantiating a thread. Locked properties inherit from the template and cannot be overridden. Required properties must be filled at thread creation time.
      </div>

      <div className="prop-area">
        {/* Column headers */}
        <div className="pg col-hdrs">
          <div className="ch-lbl">Key</div>
          <div className="ch-lbl">Type</div>
          <div className="ch-lbl">Default value</div>
          <div className="ch-lbl" style={{ textAlign: "center" }}>Req.</div>
          <div className="ch-lbl" style={{ textAlign: "center" }}>Locked</div>
          <div className="ch-lbl">Overridable</div>
          <div className="ch-lbl"></div>
        </div>

        {/* vendor */}
        <div className="pg prop-row">
          <div style={{ padding: "0 8px" }}>
            <div className="key-tag">vendor</div>
            <div className="key-desc">Name of the counterparty or supplier</div>
          </div>
          <div style={{ padding: "0 8px" }}><span className="type-badge tb-string">String</span></div>
          <div style={{ padding: "0 8px" }}><span className="default-empty">— (set at thread time)</span></div>
          <div className="pr-flag"><span className="flag-yes"></span></div>
          <div className="pr-flag"><span className="flag-no">—</span></div>
          <div style={{ padding: "0 8px" }}><span className="override-yes">✓ Yes</span></div>
          <div style={{ display: "flex", gap: "3px", padding: "0 3px" }}>
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* contract_value */}
        <div className="pg prop-row">
          <div style={{ padding: "0 8px" }}>
            <div className="key-tag">contract_value</div>
            <div className="key-desc">Total contract amount in base currency</div>
          </div>
          <div style={{ padding: "0 8px" }}><span className="type-badge tb-number">Number</span></div>
          <div style={{ padding: "0 8px" }}><span className="default-empty">— (set at thread time)</span></div>
          <div className="pr-flag"><span className="flag-yes"></span></div>
          <div className="pr-flag"><span className="flag-no">—</span></div>
          <div style={{ padding: "0 8px" }}><span className="override-yes">✓ Yes</span></div>
          <div style={{ display: "flex", gap: "3px", padding: "0 3px" }}>
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* currency — locked */}
        <div className="pg prop-row" style={{ background: "#fafafa" }}>
          <div style={{ padding: "0 8px" }}>
            <div className="key-tag">currency</div>
            <div className="key-desc">Currency code (ISO 4217). Locked by template.</div>
          </div>
          <div style={{ padding: "0 8px" }}><span className="type-badge tb-string">String</span></div>
          <div style={{ padding: "0 8px" }}><span className="default-val">USD</span></div>
          <div className="pr-flag"><span className="flag-no">—</span></div>
          <div className="pr-flag"><span className="flag-lock">🔒</span></div>
          <div style={{ padding: "0 8px" }}><span className="override-no">✕ No</span></div>
          <div style={{ display: "flex", gap: "3px", padding: "0 3px" }}>
            <div className="icon-btn" style={{ color: "#e0e0e0", cursor: "not-allowed" }}>✎</div>
            <div className="icon-btn" style={{ color: "#e0e0e0", cursor: "not-allowed" }}>✕</div>
          </div>
        </div>

        {/* contract_type */}
        <div className="pg prop-row">
          <div style={{ padding: "0 8px" }}>
            <div className="key-tag">contract_type</div>
            <div className="key-desc">Classification of contract scope</div>
          </div>
          <div style={{ padding: "0 8px" }}><span className="type-badge tb-enum">Enum</span></div>
          <div style={{ padding: "0 8px" }}><span className="default-val">external</span></div>
          <div className="pr-flag"><span className="flag-no">—</span></div>
          <div className="pr-flag"><span className="flag-no">—</span></div>
          <div style={{ padding: "0 8px" }}><span className="override-yes">✓ Yes</span></div>
          <div style={{ display: "flex", gap: "3px", padding: "0 3px" }}>
            <div className="icon-btn">✎</div>
            <div className="icon-btn del">✕</div>
          </div>
        </div>

        {/* requires_board_approval — locked bool */}
        <div className="pg prop-row" style={{ background: "#fafafa" }}>
          <div style={{ padding: "0 8px" }}>
            <div className="key-tag">requires_board_approval</div>
            <div className="key-desc">Whether board sign-off is mandatory. Locked by template.</div>
          </div>
          <div style={{ padding: "0 8px" }}><span className="type-badge tb-bool">Boolean</span></div>
          <div style={{ padding: "0 8px" }}><span className="default-val">false</span></div>
          <div className="pr-flag"><span className="flag-no">—</span></div>
          <div className="pr-flag"><span className="flag-lock">🔒</span></div>
          <div style={{ padding: "0 8px" }}><span className="override-no">✕ No</span></div>
          <div style={{ display: "flex", gap: "3px", padding: "0 3px" }}>
            <div className="icon-btn" style={{ color: "#e0e0e0", cursor: "not-allowed" }}>✎</div>
            <div className="icon-btn" style={{ color: "#e0e0e0", cursor: "not-allowed" }}>✕</div>
          </div>
        </div>

        {/* ADD PROPERTY ROW */}
        {isAddOpen ? (
          <div className="add-row">
            <div className="add-inner">
              <div style={{ padding: "0 6px" }}>
                <input className="slot-input" type="text" placeholder="new_property_key" />
              </div>
              <div style={{ padding: "0 6px" }}>
                <select className="slot-select">
                  <option>string</option>
                  <option>number</option>
                  <option>enum</option>
                  <option>boolean</option>
                  <option>date</option>
                </select>
              </div>
              <div style={{ padding: "0 6px" }}>
                <input className="slot-input" type="text" placeholder="default (optional)" />
              </div>
              <div className="add-req">
                <input type="checkbox" style={{ accentColor: "#C8922A", width: "13px", height: "13px" }} />
              </div>
              <div className="add-req">
                <input type="checkbox" style={{ accentColor: "#C8922A", width: "13px", height: "13px" }} />
              </div>
              <div style={{ padding: "0 8px", fontSize: "11px", color: "#bbb" }}>per-thread</div>
              <div className="add-actions">
                <button className="btn-confirm" onClick={() => setIsAddOpen(false)}>+ Add</button>
                <button className="btn-cancel-add" onClick={() => setIsAddOpen(false)}>✕</button>
              </div>
            </div>
          </div>
        ) : null}
      </div>

      <div className="prop-footer">
        <span className="pf-pill">5 properties defined</span>
        <span className="pf-pill">2 locked from template</span>
        <span style={{ color: "#aaa", fontSize: "11px" }}>
          Locked properties are inherited from the template and enforced at the AFP layer.
        </span>
        <div style={{ flex: 1 }}></div>
        <button className="btn-add-prop" onClick={() => setIsAddOpen(true)}>+ Add Property</button>
      </div>
    </div>
  );
};

