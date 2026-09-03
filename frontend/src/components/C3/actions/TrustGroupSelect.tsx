import React, { useEffect, useState } from "react";
import { useC3ConfigurationStore } from "@/components/C3/configuration/store/useC3ConfigurationStore";
import { useC3WorkspaceStore } from "@/components/C3/infrastructure/store/useC3WorkspaceStore";
import { listTrustGroups } from "@/services/api";

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

  const trustGroupsFromStore = useC3ConfigurationStore((state) => state.trustGroups);
  const activeWorkspaceId = useC3WorkspaceStore((state) => state.activeWorkspaceId);

  useEffect(() => {
    let isMounted = true;

    if (!activeWorkspaceId) {
      setOptions([]);
      return;
    }

    const loadAuthoritativeTrustGroups = async () => {
      setLoading(true);
      try {
        let fetched: any[] = [...trustGroupsFromStore];

        // Fetch trust groups directly for activeWorkspaceId without invoking channel refresh
        if (fetched.length === 0 && activeWorkspaceId) {
          try {
            const raw = await listTrustGroups(activeWorkspaceId);
            const rawArray = Array.isArray(raw) ? raw : (raw as any)?.data || (raw as any)?.Data || [];
            if (Array.isArray(rawArray)) {
              fetched = rawArray;
            }
          } catch (apiErr) {
            console.warn("[TrustGroupSelect] listTrustGroups direct API fetch failed:", apiErr);
          }
        }

        if (isMounted) {
          // Strictly map authoritative TrustGroup fields (id/ID, name/Name)
          // NEVER map channels or mock fallback options!
          const mapped: TrustGroupOption[] = fetched
            .filter((tg: any) => Boolean(tg.id || tg.ID))
            .map((tg: any) => ({
              id: tg.id || tg.ID,
              name: tg.name || tg.Name || `Trust Group (${tg.id || tg.ID})`,
              description: tg.description || tg.Description || "Authoritative C3 Trust Group",
            }));

          setOptions(mapped);

          // Auto-select first available real Trust Group if current value is empty or invalid
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

    loadAuthoritativeTrustGroups();

    return () => {
      isMounted = false;
    };
  }, [trustGroupsFromStore, activeWorkspaceId]);

  if (!activeWorkspaceId) {
    return (
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
        ⚠️ Select an active workspace to view Trust Groups.
      </div>
    );
  }

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
          {loading ? "Loading authoritative Trust Groups..." : "⚠️ No Trust Groups available in current workspace. Please create a Trust Group in C3 Configuration first."}
        </div>
      )}
    </div>
  );
};
