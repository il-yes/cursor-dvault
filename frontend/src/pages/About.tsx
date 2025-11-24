import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Globe, Github, BookOpen, Shield } from "lucide-react";
import { DashboardLayout } from "@/components/DashboardLayout";

export default function About() {
  return (
      <DashboardLayout>
    <div className="min-h-screen bg-background flex items-center justify-center p-6">
      <Card className="w-full max-w-2xl shadow-soft">
        <CardHeader className="text-center space-y-6 pb-4">
          {/* App Icon */}
          <div className="flex justify-center">
            <Avatar className="h-20 w-20">
              <AvatarFallback className="bg-primary text-primary-foreground text-2xl font-bold">
                <Shield className="h-10 w-10" />
              </AvatarFallback>
            </Avatar>
          </div>

          {/* App Name & Description */}
          <div className="space-y-2">
            <CardTitle className="text-3xl font-bold text-foreground">Ankhora</CardTitle>
            <CardDescription className="text-base text-muted-foreground">
              A secure, privacy-focused password & identity vault.
            </CardDescription>
          </div>

          {/* Version */}
          <div className="inline-block px-4 py-1.5 bg-muted rounded-full">
            <span className="text-sm font-medium text-muted-foreground">Version 1.0.0-alpha</span>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          <Separator />

          {/* Creators / Credits */}
          <div className="space-y-3">
            <h3 className="text-lg font-semibold text-foreground">Creators / Credits</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Built by a team of security & privacy enthusiasts dedicated to creating open-source tools
              that empower users to take control of their digital identity.
            </p>
          </div>

          <Separator />

          {/* Social & Links */}
          <div className="space-y-3">
            <h3 className="text-lg font-semibold text-foreground">Links</h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <Button
                variant="outline"
                className="w-full justify-start gap-2 transition-smooth hover:border-primary"
                onClick={() => window.open("https://dvault.example.com", "_blank")}
              >
                <Globe className="h-4 w-4" />
                Website
              </Button>
              <Button
                variant="outline"
                className="w-full justify-start gap-2 transition-smooth hover:border-primary"
                onClick={() => window.open("https://github.com/dvault/dvault", "_blank")}
              >
                <Github className="h-4 w-4" />
                GitHub
              </Button>
              <Button
                variant="outline"
                className="w-full justify-start gap-2 transition-smooth hover:border-primary"
                onClick={() => window.open("https://docs.dvault.example.com", "_blank")}
              >
                <BookOpen className="h-4 w-4" />
                Documentation
              </Button>
            </div>
          </div>

          <Separator />

          {/* Footer */}
          <div className="text-center pt-2">
            <p className="text-sm text-muted-foreground">
              Made with <span className="text-destructive">❤️</span> for security & freedom.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>

      </DashboardLayout>
  );
}
