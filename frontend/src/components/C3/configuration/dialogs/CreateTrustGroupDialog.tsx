import React, { useState } from "react";
import { useC3ConfigurationStore } from "../store/useC3ConfigurationStore";

interface CreateTrustGroupDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CreateTrustGroupDialog: React.FC<CreateTrustGroupDialogProps> = ({
  isOpen,
  onClose,
}) => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createTrustGroup = useC3ConfigurationStore((state) => state.createTrustGroup);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError("Trust Group identifier is required.");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      await createTrustGroup({
        name: name.trim(),
        description: description.trim(),
      });

      setName("");
      setDescription("");
      onClose();
    } catch (err: any) {
      setError(err?.message || "Failed to create Trust Group.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <div className="scrim" onClick={onClose}></div>
      <div className="modal-wrap">
        <div className="modal">
          <div className="modal-header">
            <div className="mh-top">
              <div className="vault-badge">
                <div className="vbdot"></div>vault_sovereign
              </div>
              <div className="mh-close" onClick={onClose}>✕</div>
            </div>
            <div className="mh-title">Create Trust Group</div>
            <div className="mh-subtitle">Configure sovereign access boundary parameters</div>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="modal-body">
              {error && (
                <div style={{ padding: "10px", background: "#FEF2F2", border: "1px solid #FECACA", color: "#DC2626", borderRadius: "6px", fontSize: "12px" }}>
                  ⚠️ {error}
                </div>
              )}

              <div className="field-group">
                <div className="fl">
                  Group Identifier <span className="fl-hint">Monospace identifier</span>
                </div>
                <input
                  type="text"
                  className="commit-value"
                  style={{ minHeight: "38px" }}
                  placeholder="e.g. internal-approvers"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  disabled={isSubmitting}
                  required
                />
              </div>

              <div className="field-group">
                <div className="fl">
                  Description <span className="fl-hint">Sovereign boundary purpose</span>
                </div>
                <textarea
                  className="commit-value"
                  placeholder="Named collection of vaults sharing KEK access envelopes..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  disabled={isSubmitting}
                />
              </div>
            </div>

            <div className="modal-footer">
              <button type="button" className="btn btn-ghost" onClick={onClose} disabled={isSubmitting}>
                Cancel
              </button>
              <button type="submit" className="btn-commit" disabled={isSubmitting || !name.trim()}>
                {isSubmitting ? "Creating..." : "✦ Create Group"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
};

