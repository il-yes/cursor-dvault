import React, { useEffect, useState } from "react";
import { useC3ChannelStore } from "@/components/C3/infrastructure/store/useC3ChannelStore";
import { useC3WorkspaceStore } from "@/components/C3/infrastructure/store/useC3WorkspaceStore";
import { useAuthStore } from "@/store/useAuthStore";
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

export const TrustGroupSelect: React.FC<TrustGroupSelectProps> = ({
  value,
  onChange,
  disabled = false,
  className = "",
}) => {
  const [options, setOptions] = useState<TrustGroupOption[]>([]);
  const [loading, setLoading] = useState(false);

  const channels = useC3ChannelStore((state) => state.channels);
  const activeWorkspaceId = useC3WorkspaceStore((state) => state.activeWorkspaceId);

  useEffect(() => {
    let isMounted = true;

    const loadTrustGroups = async () => {
      setLoading(true);
      try {
        let fetchedChannels = [...channels];

        // If channel store is empty, attempt to fetch from active workspace
        if (fetchedChannels.length === 0 && activeWorkspaceId) {
          await useC3ChannelStore.getState().fetchChannels(activeWorkspaceId);
          fetchedChannels = useC3ChannelStore.getState().channels;
        }

        // Fallback directly to Wails API if still empty
        if (fetchedChannels.length === 0) {
          const jwtToken = useAuthStore.getState().jwtToken || "";
          if (jwtToken) {
            const res = await ListChannels(jwtToken, activeWorkspaceId || "me");
            if (res && Array.isArray(res)) {
              fetchedChannels = res;
            }
          }
        }

        if (isMounted) {
          const mapped: TrustGroupOption[] = fetchedChannels.map((ch: any) => ({
            id: ch.id,
            name: ch.title || `Trust Group (${ch.id})`,
            description: ch.subtitle || "Authoritative C3 Trust Group",
          }));

          setOptions(mapped);

          // Auto-select first available Trust Group if current value is empty or invalid
          if (mapped.length > 0) {
            const exists = mapped.some((opt) => opt.id === value);
            if (!value || !exists) {
              onChange(mapped[0].id);
            }
          }
        }
      } catch (err) {
        console.error("Failed to load authoritative Trust Groups:", err);
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    loadTrustGroups();

    return () => {
      isMounted = false;
    };
  }, [channels, activeWorkspaceId]);

  return (
    <div style={{ width: "100%" }}>
      {options.length > 0 ? (
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled || loading}
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
            cursor: disabled || loading ? "not-allowed" : "pointer",
          }}
        >
          {options.map((g) => (
            <option key={g.id} value={g.id}>
              🛡️ {g.name} ({g.id})
            </option>
          ))}
        </select>
      ) : (
        <div
          style={{
            padding: "8px 12px",
            backgroundColor: "rgba(245, 158, 11, 0.08)",
            border: "1px solid rgba(245, 158, 11, 0.3)",
            borderRadius: "6px",
            color: "#F59E0B",
            fontSize: "12px",
          }}
        >
          {loading ? "Loading authoritative Trust Groups..." : "⚠️ No Trust Groups / Channels available in current workspace."}
        </div>
      )}
    </div>
  );
};
