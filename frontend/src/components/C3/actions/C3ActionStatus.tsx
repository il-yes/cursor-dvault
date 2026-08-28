import React from "react";

export type ActionStatusType =
  | "requested"
  | "approved"
  | "rejected"
  | "created"
  | "transfer_requested"
  | "TRANSFER_REQUESTED"
  | "transfer_approved"
  | "TRANSFER_APPROVED"
  | "transfer_rejected"
  | "TRANSFER_REJECTED"
  | "transfer_completed"
  | "TRANSFER_COMPLETED"
  | string;

interface C3ActionStatusProps {
  status: ActionStatusType;
  actor?: string;
  timestamp?: string;
  className?: string;
}

export const C3ActionStatus: React.FC<C3ActionStatusProps> = ({
  status,
  actor,
  timestamp,
  className = "",
}) => {
  const normStatus = (status || "").toLowerCase();

  let label = "Pending";
  let bg = "rgba(139, 148, 158, 0.15)";
  let border = "rgba(139, 148, 158, 0.3)";
  let color = "#C9D1D9";
  let icon = "⏱️";

  if (normStatus === "requested" || normStatus === "approval_requested") {
    label = "Approval Requested";
    bg = "rgba(234, 179, 8, 0.15)";
    border = "rgba(234, 179, 8, 0.4)";
    color = "#FACC15";
    icon = "⏳";
  } else if (normStatus === "approved") {
    label = "Approved";
    bg = "rgba(34, 197, 94, 0.15)";
    border = "rgba(34, 197, 94, 0.4)";
    color = "#4ADE80";
    icon = "✓";
  } else if (normStatus === "rejected" || normStatus === "approval_rejected") {
    label = "Rejected";
    bg = "rgba(239, 68, 68, 0.15)";
    border = "rgba(239, 68, 68, 0.4)";
    color = "#F87171";
    icon = "✕";
  } else if (normStatus === "created" || normStatus === "reject_created") {
    label = "Rejection Recorded";
    bg = "rgba(239, 68, 68, 0.15)";
    border = "rgba(239, 68, 68, 0.4)";
    color = "#F87171";
    icon = "🚫";
  } else if (normStatus === "transfer_requested") {
    label = "Transfer Requested";
    bg = "rgba(59, 130, 246, 0.15)";
    border = "rgba(59, 130, 246, 0.4)";
    color = "#60A5FA";
    icon = "🔄";
  } else if (normStatus === "transfer_approved") {
    label = "Transfer Approved";
    bg = "rgba(168, 85, 247, 0.15)";
    border = "rgba(168, 85, 247, 0.4)";
    color = "#C084FC";
    icon = "👍";
  } else if (normStatus === "transfer_rejected") {
    label = "Transfer Rejected";
    bg = "rgba(239, 68, 68, 0.15)";
    border = "rgba(239, 68, 68, 0.4)";
    color = "#F87171";
    icon = "✕";
  } else if (normStatus === "transfer_completed") {
    label = "Transfer Completed";
    bg = "rgba(16, 185, 129, 0.15)";
    border = "rgba(16, 185, 129, 0.4)";
    color = "#34D399";
    icon = "🎉";
  }

  return (
    <div
      className={`c3-action-status-badge ${className}`}
      style={{
        display: "inline-flex",
        alignItems: "center",
        gap: "6px",
        padding: "4px 10px",
        borderRadius: "12px",
        backgroundColor: bg,
        border: `1px solid ${border}`,
        color: color,
        fontSize: "12px",
        fontWeight: 500,
        lineHeight: 1.2,
      }}
    >
      <span>{icon}</span>
      <span>{label}</span>
      {actor && <span style={{ opacity: 0.8, fontSize: "11px" }}>by {actor}</span>}
      {timestamp && <span style={{ opacity: 0.6, fontSize: "10px" }}>• {timestamp}</span>}
    </div>
  );
};
