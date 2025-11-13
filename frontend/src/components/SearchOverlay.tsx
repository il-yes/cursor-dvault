import { useState, useEffect, useMemo } from "react";
import { VaultEntry } from "@/types/vault";
import { LogIn, CreditCard, UserCircle, FileText, Key } from "lucide-react";
import { cn } from "@/lib/utils";

interface SearchOverlayProps {
  entries: VaultEntry[];
  searchQuery: string;
  onSelectEntry: (entry: VaultEntry) => void;
  onClose: () => void;
}

const entryTypeIcons = {
  login: LogIn,
  card: CreditCard,
  identity: UserCircle,
  note: FileText,
  sshkey: Key,
};

const entryTypeLabels = {
  login: "Logins",
  card: "Cards",
  identity: "Identities",
  note: "Notes",
  sshkey: "SSH Keys",
};

export function SearchOverlay({ entries, searchQuery, onSelectEntry, onClose }: SearchOverlayProps) {
  const [selectedIndex, setSelectedIndex] = useState(0);

  const filteredResults = useMemo(() => {
    if (!searchQuery.trim()) return {};

    const query = searchQuery.toLowerCase();
    const results: Record<string, VaultEntry[]> = {};

    entries.forEach((entry) => {
      const nameMatch = entry.entry_name.toLowerCase().includes(query);
      const typeMatch = entry.type.toLowerCase().includes(query);

      if (nameMatch || typeMatch) {
        if (!results[entry.type]) {
          results[entry.type] = [];
        }
        results[entry.type].push(entry);
      }
    });

    return results;
  }, [entries, searchQuery]);

  const flatResults = useMemo(() => {
    return Object.values(filteredResults).flat();
  }, [filteredResults]);

  const hasResults = Object.keys(filteredResults).length > 0;

  useEffect(() => {
    setSelectedIndex(0);
  }, [searchQuery]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      } else if (e.key === "ArrowDown") {
        e.preventDefault();
        setSelectedIndex((prev) => (prev + 1) % flatResults.length);
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setSelectedIndex((prev) => (prev - 1 + flatResults.length) % flatResults.length);
      } else if (e.key === "Enter" && flatResults[selectedIndex]) {
        e.preventDefault();
        onSelectEntry(flatResults[selectedIndex]);
        onClose();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [selectedIndex, flatResults, onSelectEntry, onClose]);

  if (!searchQuery.trim()) return null;

  return (
    <div className="absolute top-full left-0 right-0 mt-2 bg-background border border-border rounded-lg shadow-lg max-h-96 overflow-y-auto z-50 animate-fade-in">
      {!hasResults ? (
        <div className="p-4 text-center text-muted-foreground">
          No results found for "{searchQuery}"
        </div>
      ) : (
        <div className="py-2">
          {Object.entries(filteredResults).map(([type, typeEntries]) => {
            const Icon = entryTypeIcons[type as keyof typeof entryTypeIcons];
            const label = entryTypeLabels[type as keyof typeof entryTypeLabels];

            return (
              <div key={type} className="mb-2">
                <div className="px-4 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider flex items-center gap-2">
                  <Icon className="h-3 w-3" />
                  {label} ({typeEntries.length})
                </div>
                {typeEntries.map((entry, index) => {
                  const globalIndex = flatResults.indexOf(entry);
                  const isSelected = globalIndex === selectedIndex;

                  return (
                    <button
                      key={entry.id}
                      onClick={() => {
                        onSelectEntry(entry);
                        onClose();
                      }}
                      onMouseEnter={() => setSelectedIndex(globalIndex)}
                      className={cn(
                        "w-full px-4 py-2 text-left hover:bg-secondary/50 transition-colors flex items-center gap-3",
                        isSelected && "bg-primary/10 border-l-2 border-primary"
                      )}
                    >
                      <Icon className="h-4 w-4 text-muted-foreground" />
                      <div className="flex-1 min-w-0">
                        <div className="font-medium truncate">{entry.entry_name}</div>
                        {entry.type === "login" && "user_name" in entry && (
                          <div className="text-xs text-muted-foreground truncate">
                            {entry.user_name}
                          </div>
                        )}
                      </div>
                    </button>
                  );
                })}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
