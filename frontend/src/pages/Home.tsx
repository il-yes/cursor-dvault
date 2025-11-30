import { useState, useEffect } from "react";
import { ArrowRight, Github, Twitter, Linkedin, ChevronDown, UserCircle } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { OnboardingModal } from "@/components/OnboardingModal";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import localSovereignty from "@/assets/local-sovereignty.jpg";
import decentralizedSync from "@/assets/decentralized-sync.jpg";
import sovereignCloud from "@/assets/sovereign-cloud.jpg";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import ankhoraLogoColored from "@/assets/ankhora-logo-colored-latest.png";
import { ThemeToggle } from "@/components/ThemeToggle";

const Home = () => {
  const navigate = useNavigate();
  const [inputValue, setInputValue] = useState("");
  const [scrolled, setScrolled] = useState(false);
  const [onboardingOpen, setOnboardingOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const handleGetStarted = () => {
    navigate('/login/email');
  };

  const scrollToSection = (id: string) => {
    const element = document.getElementById(id);
    element?.scrollIntoView({ behavior: "smooth" });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-zinc-50 via-white to-zinc-50 dark:from-zinc-950 dark:via-zinc-900 dark:to-black">
      {/* Glass Navbar */}
      <nav
        className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 backdrop-blur-xl ${
          scrolled
            ? "bg-white/80 dark:bg-zinc-900/80 shadow-xl border-b border-zinc-200/50 dark:border-zinc-700/50"
            : "bg-white/40 dark:bg-zinc-900/40"
        }`}
      >
        <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
          <button
            onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
            className="flex items-center gap-3 group hover:scale-105 transition-all duration-300"
          >
            <div className="w-10 h-10 backdrop-blur-sm flex items-center justify-center group-hover:shadow-xl">
              <img src={ankhoraLogo} alt="Ankhora Logo" className="h-7 w-7 object-contain drop-shadow-md" />
            </div>
            <span className="text-2xl font-semibold bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent tracking-tight">
              <small>ANKHORA</small>
            </span>
          </button>

          <div className="flex items-center gap-8">
            <button
              onClick={() => scrollToSection("features")}
              className="text-sm font-medium text-muted-foreground/80 px-4 py-2 rounded-xl hover:bg-white/50 dark:hover:bg-zinc-800/50 hover:text-primary transition-all duration-300 group"
            >
              Services
            </button>
            <button
              onClick={() => scrollToSection("pricing")}
              className="text-sm font-medium text-muted-foreground/80 px-4 py-2 rounded-xl hover:bg-white/50 dark:hover:bg-zinc-800/50 hover:text-primary transition-all duration-300 group"
            >
              Pricing
            </button>

            <ThemeToggle />

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="h-12 w-12 rounded-2xl border border-zinc-200/50 dark:border-zinc-700/50 bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm hover:bg-white/70 dark:hover:bg-zinc-800/70 shadow-sm hover:shadow-md transition-all"
                >
                  <UserCircle className="h-6 w-6 text-muted-foreground" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="backdrop-blur-sm bg-white/80 dark:bg-zinc-900/80 border-zinc-200/50 dark:border-zinc-700/50 shadow-2xl">
                <DropdownMenuItem onClick={() => navigate("/login/email")} className="rounded-xl hover:bg-primary/10">
                  Sign In
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => navigate("/vault/offline")} className="rounded-xl hover:bg-primary/10">
                  Offline Mode
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative min-h-screen flex items-center justify-center px-4 py-32 overflow-hidden">
        {/* Enhanced background orbs */}
        <div className="absolute top-1/4 left-10 w-80 h-80 bg-gradient-to-br from-primary/15 to-amber-500/10 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-1/4 right-10 w-96 h-96 bg-gradient-to-br from-accent/10 to-primary/15 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '2s' }} />
        
        <div className="relative z-10 max-w-5xl mx-auto text-center">
          {/* Hero Logo */}
          <div className="flex justify-center mb-12">
            <div className="relative group">
              <div className="absolute inset-0 bg-gradient-to-r from-primary via-primary/50 to-amber-500 rounded-3xl blur-3xl animate-pulse opacity-40 group-hover:opacity-70 transition-all" />
                <img 
                  src={ankhoraLogo} 
                  alt="Ankhora Logo" 
                  className="h-32 w-auto drop-shadow-2xl"
                />
            </div>
          </div>
          
          {/* Tagline */}
          <div className="mb-8">
            <div className="inline-flex items-center gap-3 px-6 py-2 bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm rounded-2xl border border-zinc-200/50 shadow-lg mb-4">
              <div className="w-2 h-2 bg-gradient-to-r from-primary to-amber-500 rounded-full" />
              <span className="text-lg uppercase tracking-widest font-semibold bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                Self-Sovereign Digital Vault
              </span>
            </div>
          </div>
          
          {/* Hero Typography */}
          <h1 className="text-7xl md:text-9xl font-light mb-12 leading-[0.9] tracking-tight text-foreground">
            Your Data. <br className="md:hidden" />
            <span className="font-black bg-gradient-to-r from-primary via-amber-500 to-primary/80 bg-clip-text text-transparent drop-shadow-2xl">
              Your Control.
            </span> <br />
            <span className="text-5xl md:text-7xl font-light text-muted-foreground/80">Forever.</span>
          </h1>
          
          {/* Hero Description */}
          <p className="text-xl md:text-2xl text-muted-foreground/90 mb-20 max-w-3xl mx-auto leading-relaxed font-light backdrop-blur-sm">
            Store sensitive information with zero-trust architecture. Encrypted on your device, 
            backed up to IPFS, and verifiable on the blockchain.
          </p>
          
          {/* Premium CTA */}
          <div className="max-w-2xl mx-auto">
            <div className="group relative">
              <div className="absolute -inset-2 bg-gradient-to-r from-primary via-amber-500 to-primary rounded-3xl blur-xl opacity-30 group-hover:opacity-50 transition-all" />
              <div className="relative bg-white/60 dark:bg-zinc-900/50 backdrop-blur-xl rounded-3xl p-1 shadow-2xl border border-white/40 dark:border-zinc-700/40 group-hover:shadow-primary/20">
                <div className="flex bg-white/0 rounded-2xl p-1">
                  <Input
                    type="text"
                    placeholder="Store anything to start your vault…"
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    className="flex-1  backdrop-blur-sm border-0 focus-visible:ring-2 focus-visible:ring-primary/40 rounded-2xl text-lg px-6 h-16 font-medium"
                    onKeyDown={(e) => e.key === "Enter" && handleGetStarted()}
                  />
                  <Button 
                    onClick={handleGetStarted}
                    className="bg-gradient-to-r from-primary ml-3 to-amber-500 hover:from-primary/90 hover:to-amber-500/90 shadow-2xl hover:shadow-primary/30 text-lg font-semibold px-10 h-16 rounded-2xl transition-all group-hover:scale-[1.02]"
                  >
                    Encrypt & Begin
                    <ArrowRight className="ml-3 h-6 w-6 group-hover:translate-x-1 transition-transform" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
          
          {/* Scroll Indicator */}
          <div className="mt-32">
            <ChevronDown className="h-8 w-8 mx-auto text-primary/60 animate-bounce group-hover:text-primary/80 transition-colors" />
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-32 px-6 bg-gradient-to-b from-white/50 via-transparent to-zinc-50/50 dark:from-zinc-900/50 dark:to-zinc-950/50">
        <div className="max-w-7xl mx-auto">
          <div className="grid md:grid-cols-3 gap-8">
            {[
              { title: "Local Sovereignty", img: localSovereignty, desc: "Start secure, stay private. Ankhora Shield gives individuals full local control with zero cloud dependency — encrypt and store your sensitive data directly on your device or sync it manually to IPFS." },
              { title: "Decentralized Sync", img: decentralizedSync, desc: "Designed for professionals and teams, Ankhora Sentinel automatically syncs data across devices using IPFS and anchors every change on Stellar for transparent immutability — without sacrificing usability." },
              { title: "Sovereign Cloud", img: sovereignCloud, desc: "For organizations that demand compliance and control. Ankhora Fortress delivers full audit trails, managed IPFS nodes, key rotation, and dedicated Stellar anchoring — your data, verifiable and sovereign." }
            ].map((feature, index) => (
              <div key={feature.title} className="group">
                <div className="text-center mb-8">
                  <h3 className="text-3xl font-semibold mb-6 bg-gradient-to-r from-foreground to-primary/70 bg-clip-text text-transparent">
                    {feature.title}
                  </h3>
                  <p className="text-lg text-muted-foreground/80 leading-relaxed max-w-md mx-auto">
                    {feature.desc}
                  </p>
                </div>
                <div className="group relative overflow-hidden rounded-3xl shadow-xl backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60 border border-white/40 dark:border-zinc-700/40 hover:shadow-2xl hover:shadow-primary/20 transition-all duration-500">
                  <div className="absolute inset-0 bg-gradient-to-t from-primary/5 to-transparent" />
                  <img 
                    src={feature.img} 
                    alt={feature.title}
                    className="w-full h-80 object-cover group-hover:scale-110 transition-transform duration-700"
                  />
                  <div className="absolute bottom-6 left-6 right-6">
                    <div className="bg-white/90 dark:bg-zinc-800/90 backdrop-blur-sm px-6 py-3 rounded-2xl shadow-lg">
                      <span className="text-sm font-semibold text-primary/90">Learn More →</span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Why Ankhora */}
      <section id="pricing" className="py-32 px-6">
        <div className="max-w-7xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-20 items-start">
            <div>
              <h2 className="text-6xl md:text-7xl font-semibold mb-8 bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent tracking-tight leading-tight">
                why <small>ANKHORA</small>?
              </h2>
              <p className="text-xl text-muted-foreground/90 mb-12 leading-relaxed max-w-lg backdrop-blur-sm">
                <small>ANKHORA</small> combines the distributed permanence of IPFS with the immutable proof 
                layer of Stellar blockchain. Your data stays encrypted and private, while every 
                change is cryptographically anchored — giving you verifiable privacy without 
                compromise.
              </p>
              <Button 
                onClick={() => window.location.hash = '#pricing'}
                size="lg"
                className="h-14 px-10 text-lg font-semibold bg-gradient-to-r from-primary to-amber-500 hover:from-primary/90 hover:to-amber-500/90 shadow-xl hover:shadow-primary/25 rounded-2xl"
              >
                Explore Plans
              </Button>
              <p className="text-2xl font-bold mt-8 bg-gradient-to-r from-primary to-amber-500 bg-clip-text text-transparent">
                Your Data, Your Keys, Your Cloud
              </p>
            </div>

            <div className="backdrop-blur-sm bg-white/40 dark:bg-zinc-900/40 rounded-3xl p-8 border border-white/30 dark:border-zinc-700/30 shadow-2xl">
              <Accordion type="single" collapsible className="w-full">
                {[
                  { value: "item-1", title: "Decentralized by Design", desc: "Built on peer-to-peer IPFS infrastructure, Ankhora eliminates central servers and single points of failure. Your data is distributed across a global network, accessible only to you, with no intermediaries." },
                  { value: "item-2", title: "Blockchain Anchored Integrity", desc: "Every vault entry is anchored to the Stellar blockchain, creating an immutable proof of existence and integrity. Verify any record at any time without exposing the underlying data." },
                  { value: "item-3", title: "Tailored for Growth", desc: "Start with local-only storage for personal use, scale to team collaboration with automatic sync, or deploy enterprise-grade infrastructure with dedicated nodes and compliance tools — all on the same sovereign architecture." },
                  { value: "item-4", title: "AI Interpreter for Compliance", desc: "Define compliance rules in plain language. Our AI interpreter translates your requirements into enforceable policies, with autonomous governance and audit trails anchored directly on-chain." }
                ].map((item) => (
                  <AccordionItem key={item.value} value={item.value} className="border-b border-white/20 dark:border-zinc-700/30 hover:bg-white/20 rounded-2xl mx-1 my-2 backdrop-blur-sm transition-all">
                    <AccordionTrigger className="text-xl font-semibold px-6 py-6 hover:text-primary hover:no-underline data-[state=open]:text-primary">
                      {item.title}
                    </AccordionTrigger>
                    <AccordionContent className="px-6 pb-8 text-lg text-muted-foreground/90 leading-relaxed backdrop-blur-sm">
                      {item.desc}
                    </AccordionContent>
                  </AccordionItem>
                ))}
              </Accordion>
            </div>
          </div>
        </div>
      </section>

      {/* Glass Footer */}
      <footer className="py-16 px-6 backdrop-blur-xl bg-white/40 dark:bg-zinc-900/40 border-t border-zinc-200/30 dark:border-zinc-700/30">
        <div className="max-w-7xl mx-auto">
          <div className="flex flex-col lg:flex-row justify-between items-center gap-8">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary/20 to-amber-500/20 backdrop-blur-sm flex items-center justify-center p-3 shadow-lg">
                <img src={ankhoraLogoColored} alt="Ankhora Logo" className="h-8 w-8 object-contain" />
              </div>
              <p className="text-lg text-muted-foreground/80 font-medium backdrop-blur-sm">
                Ankhora © 2025 — Built for the Self-Sovereign Web.
              </p>
            </div>
            
            <div className="flex flex-wrap gap-8 text-sm font-medium">
              <a href="#docs" className="text-muted-foreground/80 hover:text-primary transition-all px-4 py-2 rounded-xl hover:bg-white/50 backdrop-blur-sm">
                Docs
              </a>
              <a href="#api" className="text-muted-foreground/80 hover:text-primary transition-all px-4 py-2 rounded-xl hover:bg-white/50 backdrop-blur-sm">
                API
              </a>
              <a href="#privacy" className="text-muted-foreground/80 hover:text-primary transition-all px-4 py-2 rounded-xl hover:bg-white/50 backdrop-blur-sm">
                Privacy
              </a>
              <a href="https://github.com" className="text-muted-foreground/80 hover:text-primary transition-all px-4 py-2 rounded-xl hover:bg-white/50 backdrop-blur-sm">
                GitHub
              </a>
            </div>
            
            <div className="flex gap-4">
              {[
                { href: "https://github.com", icon: Github },
                { href: "https://twitter.com", icon: Twitter },
                { href: "https://linkedin.com", icon: Linkedin }
              ].map(({ href, icon: Icon }) => (
                <a 
                  key={href}
                  href={href}
                  className="group p-3 rounded-2xl bg-white/50 dark:bg-zinc-800/50 backdrop-blur-sm border border-zinc-200/50 hover:bg-white/70 dark:hover:bg-zinc-800/70 shadow-sm hover:shadow-md hover:scale-110 transition-all duration-300"
                >
                  <Icon className="h-5 w-5 text-muted-foreground group-hover:text-primary" />
                </a>
              ))}
            </div>
          </div>
        </div>
      </footer>

      <OnboardingModal open={onboardingOpen} onOpenChange={setOnboardingOpen} />
    </div>
  );
};

export default Home;
