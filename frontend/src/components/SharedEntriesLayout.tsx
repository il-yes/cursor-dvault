import { useState, useMemo, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { SharedEntry, ShareFilter, DetailView } from "@/types/sharing";
import { SharedEntriesList } from "@/components/SharedEntriesList";
import { SharedEntryOverview } from "@/components/SharedEntryOverview";
import { SharedEntryDetails } from "@/components/SharedEntryDetails";
import { useVaultStore } from "@/store/vaultStore";

export function SharedEntriesLayout() {
  const sharedEntries = useVaultStore((state) => state.shared.items);
  const [selectedEntry, setSelectedEntry] = useState<SharedEntry | null>(null);
  const [searchParams] = useSearchParams();
  const filterParam = (searchParams.get("filter") || "all") as ShareFilter;
  const [filter, setFilter] = useState<ShareFilter>(filterParam);
  const [detailView, setDetailView] = useState<DetailView>("recipients");
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    setFilter(filterParam);
  }, [filterParam]);

  useEffect(() => {
    const handleRefresh = () => {
      setRefreshKey(prev => prev + 1);
    };
    
    window.addEventListener('shareEntriesRefresh', handleRefresh);
    return () => window.removeEventListener('shareEntriesRefresh', handleRefresh);
  }, []);

  const filteredEntries = useMemo(() => {
    let filtered = [...sharedEntries];

    switch (filter) {
      case "sent":
        // In real app, filter by entries shared by current user
        filtered = filtered.filter(e => e.status === "active" || e.status === "pending");
        break;
      case "received":
        // In real app, filter by entries received by current user
        filtered = filtered.filter(e => e.status === "active");
        break;
      case "pending":
        filtered = filtered.filter(e => e.status === "pending");
        break;
      case "revoked":
        filtered = filtered.filter(e => e.status === "revoked");
        break;
      case "all":
      default:
        // Show all entries
        break;
    }

    return filtered;
  }, [filter, sharedEntries]);

  return (
    <div className="flex h-full" key={refreshKey}>
      {/* Column 2: Shared Entries List (Column 1 is the main sidebar) */}
      <div className="w-full md:w-80 lg:w-96 flex flex-col border-r border-border bg-secondary/30 overflow-hidden">
        <SharedEntriesList
          entries={filteredEntries}
          selectedEntryId={selectedEntry?.id || null}
          onSelectEntry={setSelectedEntry}
        />
      </div>

      {/* Column 3: Entry Overview */}
      <div className="hidden md:flex flex-col w-80 lg:w-96 border-r border-border overflow-hidden">
        <div className="flex-1 overflow-y-auto">
          <SharedEntryOverview
            entry={selectedEntry}
            onViewChange={setDetailView}
          />
        </div>
      </div>

      {/* Column 4: Detail Panel (but visually column 3) */}
      <div className="flex-1 hidden lg:flex flex-col overflow-hidden">
        <div className="flex-1 overflow-y-auto">
          <SharedEntryDetails
            entry={selectedEntry}
            view={detailView}
          />
        </div>
      </div>
    </div>
  );
}
