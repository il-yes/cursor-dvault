import { useLocation, useNavigate } from "react-router-dom";
import { useVaultStore } from "@/store/vaultStore";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Share2 } from "lucide-react";
import { Card } from "./ui/card";

export default function ShareEntries() {
  const location = useLocation();
  const navigate = useNavigate();

  const sharedByMe = useVaultStore(state => state.shared.items);
  const sharedWithMe = useVaultStore(state => (state as any).sharedWithMe?.items || []); // existing store slice
  const urlFilter = new URLSearchParams(location.search).get("filter") || "all";

  // Choose data source based on tab
  const tab = urlFilter === "withme" ? "withme" : "byme";

  const baseList = tab === "byme" ? sharedByMe : sharedWithMe;

  // Apply your existing filters on the selected slice
  let entries = baseList;
  if (urlFilter === "sent") entries = sharedByMe;
  if (urlFilter === "pending") entries = sharedByMe.filter(e => e.status === "pending");
  if (urlFilter === "revoked") entries = sharedByMe.filter(e => e.status === "revoked");
  if (urlFilter === "withme") entries = sharedWithMe;

  return (
    <div className="p-6 space-y-6">

      <div className="flex items-center gap-3">
        <Share2 className="h-6 w-6 text-primary" />
        <h1 className="text-xl font-bold">Shared Entries</h1>
      </div>

      {/* UI Tabs */}
      <Tabs
        value={tab}
        onValueChange={(t) => {
          // Only switch data source, do NOT break filters
          if (t === "withme") {
            navigate("/dashboard/shared?filter=withme");
          } else {
            navigate("/dashboard/shared");
          }
        }}
      >
        <TabsList className="grid w-[300px] grid-cols-2">
          <TabsTrigger value="byme">By me</TabsTrigger>
          <TabsTrigger value="withme">With me</TabsTrigger>
        </TabsList>

        <TabsContent value="byme">
          <SharesList entries={entries} emptyLabel="No shares sent by you yet." />
        </TabsContent>

        <TabsContent value="withme">
          <SharesList entries={entries} emptyLabel="No shares received for you yet." />
        </TabsContent>
      </Tabs>

    </div>
  );
}


function SharesList({ entries, emptyLabel }: { entries: any[]; emptyLabel: string }) {
  if (!entries.length) {
    return <p className="text-sm text-muted-foreground">{emptyLabel}</p>;
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {entries.map((s) => (
        <Card key={s.id} className="p-3">
          <div className="text-sm font-medium">{s.entry_name}</div>
          <div className="text-xs text-muted-foreground">
            shared with {s.recipients?.length || 0} people
          </div>
        </Card>
      ))}
    </div>
  );
}


