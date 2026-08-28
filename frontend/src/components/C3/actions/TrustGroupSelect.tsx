import React, { useEffect, useState } from "react";
import { ListChannels } from "../../../../wailsjs/go/main/App";

export interface TrustGroupOption {
  id: string;
  name: string;
  description?: string;
}

interface TrustGroupSelectProps {
  value: string;
  onChange: (trustGroupId: string) => void;
  disabled?: boolean;
  className?: string;
}

const DEFAULT_TRUST_GROUPS: TrustGroupOption[] = [
  { id: "tg_engineering", name: "Engineering Trust Group", description: "Architecture & Engineering Reviewers" },
  { id: "tg_legal", name: "Legal Counsel", description: "Legal & Regulatory Reviewers" },
  { id: "tg_executive", name: "Executive Committee", description: "Executive Sign-off & Approvals" },
  { id: "tg_finance", name: "Finance & Accounting", description: "Financial Audit & Release" },
];

export const TrustGroupSelect: React.FC<TrustGroupSelectProps> = ({
  value,
  onChange,
  disabled = false,
  className = "",
}) => {
  const [groups, setGroups] = useState<TrustGroupOption[]>(DEFAULT_TRUST_GROUPS);

  useEffect(() => {
    let isMounted = true;
    const fetchChannels = async () => {
      try {
        const channels = await ListChannels("me", "all");
        if (channels && channels.length > 0 && isMounted) {
          const dynamicGroups: TrustGroupOption[] = channels.map((ch: any) => ({
            id: `tg_${ch.id}`,
            name: ch.title || `Channel ${ch.id}`,
            description: ch.subtitle || "C3 Collaboration Trust Group",
          }));
          // Combine defaults and dynamic channel groups
          const map = new Map<string, TrustGroupOption>();
          [...DEFAULT_TRUST_GROUPS, ...dynamicGroups].forEach((g) => map.set(g.id, g));
          setGroups(Array.from(map.values()));
        }
      } catch (e) {
        // Fallback to default trust groups
      }
    };
    fetchChannels();
    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      disabled={disabled}
      className={className}
      style={{
        width: "100%",
        padding: "8px 12px",
        backgroundColor: "#0D1117",
        border: "1px solid #30363D",
        borderRadius: "6px",
        color: "#F0F6FC",
        fontSize: "13px",
        boxSizing: "border-box",
        cursor: disabled ? "not-allowed" : "pointer",
      }}
    >
      <option value="" disabled>
        -- Select Trust Group --
      </option>
      {groups.map((g) => (
        <option key={g.id} value={g.id}>
          {g.name} ({g.id})
        </option>
      ))}
    </select>
  );
};
