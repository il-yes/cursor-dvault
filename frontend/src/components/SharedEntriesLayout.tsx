import { useState, useMemo } from "react";
import { SharedEntry, ShareFilter, DetailView } from "@/types/sharing";
import { SharedEntriesList } from "@/components/SharedEntriesList";
import { SharedEntryOverview } from "@/components/SharedEntryOverview";
import { SharedEntryDetails } from "@/components/SharedEntryDetails";
import { mockSharedEntries } from "@/data/mockSharedData";

export function SharedEntriesLayout() {
  const [selectedEntry, setSelectedEntry] = useState<SharedEntry | null>(null);
  const [filter, setFilter] = useState<ShareFilter>("all");
  const [detailView, setDetailView] = useState<DetailView>("recipients");

  const filteredEntries = useMemo(() => {
    let filtered = [...mockSharedEntries];

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
  }, [filter]);

  return (
    <div className="flex h-full overflow-hidden">
      {/* Column 1: Shared Entries List */}
      <SharedEntriesList
        entries={filteredEntries}
        selectedEntryId={selectedEntry?.id || null}
        filter={filter}
        onSelectEntry={setSelectedEntry}
        onFilterChange={setFilter}
      />

      {/* Column 2: Entry Overview */}
      <div className="hidden md:flex flex-col w-80 lg:w-96 border-r border-border overflow-hidden">
        <SharedEntryOverview
          entry={selectedEntry}
          onViewChange={setDetailView}
        />
      </div>

      {/* Column 3: Detail Panel */}
      <div className="hidden lg:flex flex-1 overflow-hidden">
        <SharedEntryDetails
          entry={selectedEntry}
          view={detailView}
        />
      </div>
    </div>
  );
}
