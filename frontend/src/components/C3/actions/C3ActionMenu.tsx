import React, { useState } from "react";
import { ResourceReference } from "../domain/resource";
import { C3ShareDialog } from "./C3ShareDialog";
import { ApprovalActionDialog } from "./ApprovalActionDialog";
import { RejectActionDialog } from "./RejectActionDialog";
import { TransferActionDialog } from "./TransferActionDialog";

interface C3ActionMenuProps {
  resourceRef: ResourceReference;
  threadId?: string;
  onView?: (ref: ResourceReference) => void;
  onActionComplete?: () => void;
  className?: string;
}

export const C3ActionMenu: React.FC<C3ActionMenuProps> = ({
  resourceRef,
  threadId = "default_thread",
  onView,
  onActionComplete,
  className = "",
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [activeModal, setActiveModal] = useState<"share" | "approval" | "reject" | "transfer" | null>(null);

  const toggleMenu = () => setIsOpen((prev) => !prev);
  const closeModal = () => {
    setActiveModal(null);
    if (onActionComplete) onActionComplete();
  };

  return (
    <div className={`c3-action-menu-container ${className}`} style={{ position: "relative", display: "inline-block" }}>
      {/* Menu Trigger Button */}
      <button
        type="button"
        onClick={toggleMenu}
        aria-label="C3 Collaboration Action Menu"
        style={{
          padding: "4px 8px",
          backgroundColor: "#21262D",
          border: "1px solid #30363D",
          borderRadius: "6px",
          color: "#C9D1D9",
          fontSize: "12px",
          fontWeight: 500,
          cursor: "pointer",
          display: "flex",
          alignItems: "center",
          gap: "4px",
        }}
      >
        <span>Actions</span>
        <span style={{ fontSize: "10px" }}>▼</span>
      </button>

      {/* Dropdown Menu Overlay */}
      {isOpen && (
        <>
          <div
            onClick={() => setIsOpen(false)}
            style={{
              position: "fixed",
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              zIndex: 99,
            }}
          />
          <div
            style={{
              position: "absolute",
              right: 0,
              top: "100%",
              marginTop: "4px",
              width: "160px",
              backgroundColor: "#161B22",
              border: "1px solid #30363D",
              borderRadius: "6px",
              boxShadow: "0 10px 15px -3px rgba(0, 0, 0, 0.5)",
              zIndex: 100,
              padding: "4px 0",
              display: "flex",
              flexDirection: "column",
            }}
          >
            {onView && (
              <button
                type="button"
                onClick={() => {
                  setIsOpen(false);
                  onView(resourceRef);
                }}
                style={{
                  padding: "8px 12px",
                  textAlign: "left",
                  background: "none",
                  border: "none",
                  color: "#F0F6FC",
                  fontSize: "13px",
                  cursor: "pointer",
                }}
              >
                👁️ View
              </button>
            )}

            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                setActiveModal("share");
              }}
              style={{
                padding: "8px 12px",
                textAlign: "left",
                background: "none",
                border: "none",
                color: "#F0F6FC",
                fontSize: "13px",
                cursor: "pointer",
              }}
            >
              🤝 Create Share
            </button>

            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                setActiveModal("approval");
              }}
              style={{
                padding: "8px 12px",
                textAlign: "left",
                background: "none",
                border: "none",
                color: "#F0F6FC",
                fontSize: "13px",
                cursor: "pointer",
              }}
            >
              👍 Request Approval
            </button>

            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                setActiveModal("reject");
              }}
              style={{
                padding: "8px 12px",
                textAlign: "left",
                background: "none",
                border: "none",
                color: "#F87171",
                fontSize: "13px",
                cursor: "pointer",
              }}
            >
              🚫 Reject Resource
            </button>

            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                setActiveModal("transfer");
              }}
              style={{
                padding: "8px 12px",
                textAlign: "left",
                background: "none",
                border: "none",
                color: "#60A5FA",
                fontSize: "13px",
                cursor: "pointer",
              }}
            >
              🔄 Transfer Resource
            </button>
          </div>
        </>
      )}

      {/* Action Dialogs */}
      <C3ShareDialog
        isOpen={activeModal === "share"}
        onClose={closeModal}
        resourceRef={resourceRef}
        threadId={threadId}
      />

      <ApprovalActionDialog
        isOpen={activeModal === "approval"}
        onClose={closeModal}
        resourceRef={resourceRef}
        threadId={threadId}
      />

      <RejectActionDialog
        isOpen={activeModal === "reject"}
        onClose={closeModal}
        resourceRef={resourceRef}
        threadId={threadId}
      />

      <TransferActionDialog
        isOpen={activeModal === "transfer"}
        onClose={closeModal}
        resourceRef={resourceRef}
        threadId={threadId}
      />
    </div>
  );
};
