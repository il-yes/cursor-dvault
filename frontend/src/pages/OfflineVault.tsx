import { useState } from "react";
import { ArrowLeft, Wifi, WifiOff, Upload } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

const OfflineVault = () => {
  const navigate = useNavigate();
  const [isOnline] = useState(false); // Mock offline state

  // Mock local vault entries
  const mockEntries = [
    {
      id: "1",
      title: "Personal Documents",
      description: "Birth certificate, passport, ID cards",
      lastModified: "2025-01-15",
      encrypted: true,
    },
    {
      id: "2",
      title: "Financial Records",
      description: "Tax documents, bank statements",
      lastModified: "2025-01-10",
      encrypted: true,
    },
    {
      id: "3",
      title: "Medical Records",
      description: "Health insurance, prescriptions",
      lastModified: "2024-12-28",
      encrypted: true,
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => navigate("/")}
              className="rounded-full"
            >
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <div>
              <h1 className="text-2xl font-semibold text-foreground">Offline Vault Mode</h1>
              <p className="text-sm text-muted-foreground">
                Access your locally encrypted vault. Changes will sync once you're back online.
              </p>
            </div>
          </div>

          <div className="flex items-center gap-4">
            {/* Connection Status */}
            <Badge
              variant={isOnline ? "default" : "secondary"}
              className="flex items-center gap-2"
            >
              {isOnline ? (
                <>
                  <Wifi className="h-3 w-3" />
                  Online
                </>
              ) : (
                <>
                  <WifiOff className="h-3 w-3" />
                  Offline
                </>
              )}
            </Badge>

            <Button
              disabled={!isOnline}
              className="gradient-primary text-primary-foreground"
            >
              <Upload className="mr-2 h-4 w-4" />
              Sync to IPFS
            </Button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 py-12">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {mockEntries.map((entry) => (
            <Card
              key={entry.id}
              className="hover:shadow-elegant transition-smooth cursor-pointer group"
            >
              <CardHeader>
                <div className="flex items-start justify-between">
                  <CardTitle className="text-lg group-hover:text-primary transition-smooth">
                    {entry.title}
                  </CardTitle>
                  {entry.encrypted && (
                    <Badge variant="outline" className="text-xs">
                      🔒 Encrypted
                    </Badge>
                  )}
                </div>
                <CardDescription>{entry.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-xs text-muted-foreground">
                  Last modified: {entry.lastModified}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>

        {mockEntries.length === 0 && (
          <div className="text-center py-20">
            <p className="text-muted-foreground text-lg">
              No entries in your offline vault yet.
            </p>
            <Button className="mt-4" onClick={() => navigate("/dashboard")}>
              Create Your First Entry
            </Button>
          </div>
        )}
      </main>
    </div>
  );
};

export default OfflineVault;
