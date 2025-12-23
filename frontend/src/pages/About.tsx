import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Globe, Github, BookOpen, Shield } from "lucide-react";
import { DashboardLayout } from "@/components/DashboardLayout";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import SubscriptionManager from "@/components/Subscription/subscriptionManager";

export default function About() {
  return (
    <DashboardLayout>
      <div className="min-h-screen flex items-center justify-center p-8 bg-gradient-to-b from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900">
        <Card className="w-full max-w-lg border-none shadow-2xl backdrop-blur-sm bg-white/60 dark:bg-zinc-900/50">
          <CardHeader className="text-center space-y-8 pb-6">
            {/* Hero Logo */}
            <div className="mx-auto w-auto h-20 rounded-3xl to-amber-500/20 backdrop-blur-sm  flex items-center justify-center p-6">
              <img
                src={ankhoraLogo}
                alt="Ankhora Logo"
                className="w-auto h-32 object-contain drop-shadow-lg"
              />
            </div>
            {/* App Identity */}
            <div className="space-y-3">
              <div className="inline-flex items-center gap-2 px-4 py-1.5 bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm rounded-2xl border border-zinc-200/50 shadow-sm">
                <Shield className="h-4 w-4 text-primary" />
                <span className="text-sm font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                  ANKHORA
                </span>
                <span className="text-xs text-muted-foreground font-medium">1.0.0-alpha</span>
              </div>
              <CardDescription className="text-lg text-muted-foreground/80 leading-relaxed max-w-sm mx-auto">
                A secure, privacy-focused cryptographic vault for your digital identity.
              </CardDescription>
            </div>

          </CardHeader>

          <CardContent className="space-y-8">
            {/* Creators */}
            <div className="group space-y-3 p-6 rounded-2xl bg-white/30 dark:bg-zinc-800/30 backdrop-blur-sm border border-zinc-200/30 hover:bg-white/50 dark:hover:bg-zinc-800/50 transition-all duration-300">
              <h3 className="text-xl font-semibold bg-gradient-to-r from-foreground to-primary/70 bg-clip-text text-transparent">
                Creators
              </h3>
              <p className="text-sm text-muted-foreground leading-relaxed">
                Built by security & privacy enthusiasts dedicated to open-source tools that empower users to own their digital identity.
              </p>
            </div>

            <Separator className="border-zinc-200/50 dark:border-zinc-700/50 mx-6" />

            {/* Premium Links Grid */}
            <div className="space-y-4">
              <h3 className="text-xl font-semibold text-center bg-gradient-to-r from-foreground to-primary/70 bg-clip-text text-transparent">
                Explore
              </h3>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 px-3">
                <Button
                  variant="ghost"
                  className="group h-14 rounded-2xl border-2 border-zinc-200/50 dark:border-zinc-700/50 bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm hover:border-primary/40 hover:bg-primary/5 hover:shadow-md transition-all duration-300 justify-start gap-3 shadow-sm hover:shadow-primary/10"
                  onClick={() => window.open("https://dvault.example.com", "_blank")}
                >
                  <Globe className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
                  <span className="text-sm font-medium">Website</span>
                </Button>

                <Button
                  variant="ghost"
                  className="group h-14 rounded-2xl border-2 border-zinc-200/50 dark:border-zinc-700/50 bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm hover:border-primary/40 hover:bg-primary/5 hover:shadow-md transition-all duration-300 justify-start gap-3 shadow-sm hover:shadow-primary/10"
                  onClick={() => window.open("https://github.com/dvault/dvault", "_blank")}
                >
                  <Github className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
                  <span className="text-sm font-medium">GitHub</span>
                </Button>

                <Button
                  variant="ghost"
                  className="group h-14 rounded-2xl border-2 border-zinc-200/50 dark:border-zinc-700/50 bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm hover:border-primary/40 hover:bg-primary/5 hover:shadow-md transition-all duration-300 justify-start gap-3 shadow-sm hover:shadow-primary/10"
                  onClick={() => window.open("https://docs.dvault.example.com", "_blank")}
                >
                  <BookOpen className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
                  <span className="text-sm font-medium">Docs</span>
                </Button>
              </div>
            </div>

            <Separator className="border-zinc-200/50 dark:border-zinc-700/50 mx-6" />

            {/* Love Footer */}
            <div className="text-center pt-4 pb-8 px-6">
              <p className="text-sm text-muted-foreground/70 leading-relaxed">
                Made with <span className="text-primary font-semibold">❤️</span> for security & freedom.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  );
}
