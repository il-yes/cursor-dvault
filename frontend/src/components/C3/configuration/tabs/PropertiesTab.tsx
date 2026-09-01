import React, { useState } from "react";
import { useC3ConfigurationStore, PropertyDraft } from "../store/useC3ConfigurationStore";

export const PropertiesTab: React.FC = () => {
  const {
    draftProperties,
    setDraftProperties,
    isDirty,
    isSaving,
    saveFeedback,
    saveChannelConfig,
    discardChanges,
  } = useC3ConfigurationStore();

  const [isAddOpen, setIsAddOpen] = useState(false);
  const [editingKey, setEditingKey] = useState<string | null>(null);
  const [newKey, setNewKey] = useState("");
  const [newType, setNewType] = useState<"String" | "Number" | "Enum" | "Boolean">("String");
  const [newDefault, setNewDefault] = useState("");
  const [isRequired, setIsRequired] = useState(false);
  const [isLocked, setIsLocked] = useState(false);
  const [isOverridable, setIsOverridable] = useState(true);

  const handleEditProperty = (prop: PropertyDraft) => {
    if (prop.locked) return;
    setEditingKey(prop.key);
    setNewKey(prop.key);
    setNewType(prop.type || "String");
    setNewDefault(prop.defaultValue || prop.value || "");
    setIsRequired(!!prop.required);
    setIsLocked(!!prop.locked);
    setIsOverridable(prop.overridable !== false);
    setIsAddOpen(true);
  };

  const handleAddPropertySubmit = () => {
    if (!newKey.trim()) return;

    const newProp: PropertyDraft = {
      key: newKey.trim(),
      value: newDefault.trim(),
      desc: "Custom metadata property",
      type: newType,
      defaultValue: newDefault.trim() || undefined,
      required: isRequired,
      locked: isLocked,
      overridable: isOverridable,
    };

    if (editingKey) {
      setDraftProperties(
        draftProperties.map((p) => (p.key === editingKey ? newProp : p))
      );
    } else {
      setDraftProperties([...draftProperties, newProp]);
    }

    setEditingKey(null);
    setNewKey("");
    setNewDefault("");
    setIsAddOpen(false);
  };

  const handleDeleteProperty = (key: string) => {
    setDraftProperties(draftProperties.filter((p) => p.key !== key));
  };

  const getTypeBadgeClass = (type: string = "String") => {
    switch (type) {
      case "String":
        return "tb-string";
      case "Number":
        return "tb-number";
      case "Enum":
        return "tb-enum";
      case "Boolean":
        return "tb-bool";
      default:
        return "tb-string";
    }
  };

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

        {/* Dynamic Property Rows */}
        {draftProperties.length === 0 && (
          <div style={{ padding: "24px 0", textAlign: "center", color: "#888", fontSize: "12px" }}>
            No properties defined for this channel yet. Click <strong>+ Add Property</strong> below to define custom metadata.
          </div>
        )}
        {draftProperties.map((prop, idx) => (
          <div
            className="pg prop-row"
            key={prop.key || idx}
            style={prop.locked ? { background: "#fafafa" } : {}}
          >
            <div style={{ padding: "0 8px" }}>
              <div className="key-tag">{prop.key}</div>
              <div className="key-desc">{prop.desc || "Channel property"}</div>
            </div>
            <div style={{ padding: "0 8px" }}>
              <span className={`type-badge ${getTypeBadgeClass(prop.type)}`}>{prop.type || "String"}</span>
            </div>
            <div style={{ padding: "0 8px" }}>
              {prop.defaultValue || prop.value ? (
                <span className="default-val">{prop.defaultValue || prop.value}</span>
              ) : (
                <span className="default-empty">— (set at thread time)</span>
              )}
            </div>
            <div className="pr-flag">
              {prop.required ? <span className="flag-yes"></span> : <span className="flag-no">—</span>}
            </div>
            <div className="pr-flag">
              {prop.locked ? <span className="flag-lock">🔒</span> : <span className="flag-no">—</span>}
            </div>
            <div style={{ padding: "0 8px" }}>
              {prop.overridable !== false ? (
                <span className="override-yes">✓ Yes</span>
              ) : (
                <span className="override-no">✕ No</span>
              )}
            </div>
            <div style={{ display: "flex", gap: "3px", padding: "0 3px" }}>
              <div
                className="icon-btn"
                style={prop.locked ? { color: "#e0e0e0", cursor: "not-allowed" } : { cursor: "pointer" }}
                onClick={() => handleEditProperty(prop)}
              >
                ✎
              </div>
              <div
                className="icon-btn del"
                style={prop.locked ? { color: "#e0e0e0", cursor: "not-allowed" } : { cursor: "pointer" }}
                onClick={() => !prop.locked && handleDeleteProperty(prop.key)}
              >
                ✕
              </div>
            </div>
          </div>
        ))}

        {/* ADD / EDIT PROPERTY ROW */}
        {isAddOpen ? (
          <div className="add-row">
            <div className="add-inner-prop">
              <div style={{ padding: "0 6px" }}>
                <input
                  className="slot-input"
                  type="text"
                  placeholder="new_property_key"
                  value={newKey}
                  onChange={(e) => setNewKey(e.target.value)}
                  disabled={!!editingKey}
                  autoFocus
                />
              </div>
              <div style={{ padding: "0 6px" }}>
                <select
                  className="slot-select"
                  value={newType}
                  onChange={(e) => setNewType(e.target.value as any)}
                >
                  <option value="String">String</option>
                  <option value="Number">Number</option>
                  <option value="Enum">Enum</option>
                  <option value="Boolean">Boolean</option>
                </select>
              </div>
              <div style={{ padding: "0 6px" }}>
                <input
                  className="slot-input"
                  type="text"
                  placeholder="default (optional)"
                  value={newDefault}
                  onChange={(e) => setNewDefault(e.target.value)}
                />
              </div>
              <div className="add-req">
                <input
                  type="checkbox"
                  checked={isRequired}
                  onChange={(e) => setIsRequired(e.target.checked)}
                  style={{ accentColor: "#C8922A", width: "13px", height: "13px" }}
                />
              </div>
              <div className="add-req">
                <input
                  type="checkbox"
                  checked={isLocked}
                  onChange={(e) => setIsLocked(e.target.checked)}
                  style={{ accentColor: "#C8922A", width: "13px", height: "13px" }}
                />
              </div>
              <div style={{ padding: "0 8px", fontSize: "11px", color: "#bbb" }}>per-thread</div>
              <div className="add-actions">
                <button className="btn-confirm" onClick={handleAddPropertySubmit}>
                  {editingKey ? "✓ Update" : "+ Add"}
                </button>
                <button
                  className="btn-cancel-add"
                  onClick={() => {
                    setIsAddOpen(false);
                    setEditingKey(null);
                    setNewKey("");
                    setNewDefault("");
                  }}
                >
                  ✕
                </button>
              </div>
            </div>
          </div>
        ) : null}
      </div>

      {(isDirty || saveFeedback) && (
        <div className="dirty-bar">
          <span className="db-icon">⚠</span>
          <div className="db-msg">
            {saveFeedback ? (
              <span>{saveFeedback}</span>
            ) : (
              <span>
                <strong>Unsaved changes</strong> — property updates take effect on the next thread instantiation.
              </span>
            )}
          </div>
          <button className="btn-discard" onClick={discardChanges} disabled={isSaving}>
            Discard
          </button>
          <button className="btn-save" onClick={saveChannelConfig} disabled={isSaving}>
            {isSaving ? "Saving..." : "Save Changes"}
          </button>
        </div>
      )}

      <div className="prop-footer">
        <span className="pf-pill">{draftProperties.length} properties defined</span>
        <span className="pf-pill">{draftProperties.filter((p) => p.locked).length} locked from template</span>
        <span style={{ color: "#aaa", fontSize: "11px" }}>
          Locked properties are inherited from the template and enforced at the AFP layer.
        </span>
        <div style={{ flex: 1 }}></div>
        <button className="btn-add-prop" onClick={() => setIsAddOpen(true)}>
          + Add Property
        </button>
      </div>
    </div>
  );
};


