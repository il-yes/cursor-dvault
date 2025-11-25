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
import ankhoraLogoColored from "@/assets/ankhora-logo-colored.png";
import ankhoraLogo from "@/assets/ankhora-logo-transparent.png";
import localSovereignty from "@/assets/local-sovereignty.jpg";
import decentralizedSync from "@/assets/decentralized-sync.jpg";
import sovereignCloud from "@/assets/sovereign-cloud.jpg";
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
    <div className="min-h-screen bg-background">
      {/* Navbar */}
      <nav
        className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
          scrolled
            ? "bg-background/95 backdrop-blur-sm shadow-sm border-b border-border"
            : "bg-transparent"
        }`}
      >
        <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
          <button
            onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
            className="text-xl font-semibold text-foreground hover:text-primary transition-smooth"
          >
            D-Vault
          </button>

          <div className="flex items-center gap-6">
            <button
              onClick={() => scrollToSection("features")}
              className="text-sm text-muted-foreground hover:text-primary transition-smooth"
            >
              Services
            </button>
            <button
              onClick={() => scrollToSection("pricing")}
              className="text-sm text-muted-foreground hover:text-primary transition-smooth"
            >
              Pricing
            </button>
            <ThemeToggle />

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="rounded-full">
                  <UserCircle className="h-5 w-5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="bg-card">
                <DropdownMenuItem onClick={() => navigate("/auth/signin")}>
                  Sign In
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => navigate("/vault/offline")}>
                  Offline Mode
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative min-h-screen flex items-center justify-center px-4 py-20 overflow-hidden">
        {/* Subtle background gradient */}
        <div className="absolute inset-0 bg-gradient-to-b from-primary/5 via-transparent to-transparent dark:from-primary/10 pointer-events-none" />
        
        {/* Animated glow orbs */}
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-pulse pointer-events-none" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-accent/10 rounded-full blur-3xl animate-pulse pointer-events-none" style={{ animationDelay: '1s' }} />
        
        <div className="relative max-w-4xl mx-auto text-center">
          {/* Logo with intense glow */}
          <div className="flex justify-center mb-8 animate-fadeInUp">
            <div className="relative">
              <div className="absolute inset-0 bg-primary/30 rounded-full blur-2xl animate-pulse-glow" />
              <img 
                src={ankhoraLogo} 
                alt="Ankhora Logo" 
                className="relative h-28 md:h-32 w-auto drop-shadow-2xl" 
              />
            </div>
          </div>
          
          {/* Tagline with shimmer effect */}
          <div className="animate-fadeInStagger mb-6" style={{ animationDelay: '0.1s' }}>
            <p className="text-sm uppercase tracking-[0.2em] text-primary font-semibold mb-2">
              Self-Sovereign Digital Vault
            </p>
            <div className="w-24 h-0.5 bg-gradient-to-r from-transparent via-primary to-transparent mx-auto" />
          </div>
          
          {/* Main heading with gradient text */}
          <h1 className="text-6xl md:text-8xl font-light mb-8 leading-[1.1] text-foreground animate-fadeInStagger" style={{ animationDelay: '0.2s' }}>
            Your Data. <br />
            <span className="font-bold bg-gradient-to-r from-primary via-primary-glow to-primary bg-clip-text text-transparent">
              Your Control.
            </span> <br />
            <span className="text-5xl md:text-7xl">Forever.</span>
          </h1>
          
          {/* Description */}
          <p className="text-lg md:text-2xl text-muted-foreground mb-16 max-w-3xl mx-auto leading-relaxed font-light animate-fadeInStagger" style={{ animationDelay: '0.3s' }}>
            Store sensitive information with zero-trust architecture. Encrypted on your device, 
            backed up to IPFS, and verifiable on the blockchain.
          </p>
          
          {/* CTA Input with enhanced glow */}
          <div className="max-w-2xl mx-auto animate-fadeInStagger" style={{ animationDelay: '0.4s' }}>
            <div className="relative group">
              <div className="absolute -inset-1 bg-gradient-to-r from-primary via-primary-glow to-primary rounded-3xl blur-xl opacity-30 group-hover:opacity-50 transition-smooth" />
              <div className="relative flex gap-3 shadow-elegant rounded-2xl p-3 bg-card border border-border/50 backdrop-blur-sm">
                <Input
                  type="text"
                  placeholder="Store anything to start your vault…"
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  className="flex-1 border-0 focus-visible:ring-0 text-base md:text-lg bg-transparent px-4 h-14"
                  onKeyDown={(e) => e.key === "Enter" && handleGetStarted()}
                />
                <Button 
                  onClick={handleGetStarted}
                  className="gradient-primary text-primary-foreground px-8 h-14 hover:opacity-90 transition-smooth shadow-lg hover:shadow-xl font-semibold text-base"
                >
                  Encrypt & Begin
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </div>
            </div>
          </div>
          
          {/* Scroll indicator */}
          <div className="mt-20 animate-bounce">
            <ChevronDown className="h-7 w-7 mx-auto text-primary/60" />
          </div>
        </div>
      </section>


      {/* Feature Section */}
      <section id="features" className="py-20 px-4 bg-card">
        <div className="max-w-7xl mx-auto">
          <div className="grid md:grid-cols-3 gap-12">
            {/* Column 1 */}
            <div className="text-center group">
              <h3 className="text-2xl font-semibold mb-4 text-foreground">Local Sovereignty</h3>
              <p className="text-muted-foreground leading-relaxed">
                Start secure, stay private. D-Vault Shield gives individuals full local control 
                with zero cloud dependency — encrypt and store your sensitive data directly on 
                your device or sync it manually to IPFS.
              </p>
              <div className="mt-6 overflow-hidden rounded-2xl shadow-soft transition-smooth group-hover:shadow-elegant">
                <img 
                  src={localSovereignty} 
                  alt="Local Sovereignty" 
                  className="w-full h-64 object-cover transition-smooth group-hover:scale-105"
                />
              </div>
            </div>

            {/* Column 2 */}
            <div className="text-center group">
              <h3 className="text-2xl font-semibold mb-4 text-foreground">Decentralized Sync</h3>
              <p className="text-muted-foreground leading-relaxed">
                Designed for professionals and teams, D-Vault Sentinel automatically syncs data 
                across devices using IPFS and anchors every change on Stellar for transparent 
                immutability — without sacrificing usability.
              </p>
              <div className="mt-6 overflow-hidden rounded-2xl shadow-soft transition-smooth group-hover:shadow-elegant">
                <img 
                  src={decentralizedSync} 
                  alt="Decentralized Sync" 
                  className="w-full h-64 object-cover transition-smooth group-hover:scale-105"
                />
              </div>
            </div>

            {/* Column 3 */}
            <div className="text-center group">
              <h3 className="text-2xl font-semibold mb-4 text-foreground">Sovereign Cloud</h3>
              <p className="text-muted-foreground leading-relaxed">
                For organizations that demand compliance and control. D-Vault Fortress delivers 
                full audit trails, managed IPFS nodes, key rotation, and dedicated Stellar 
                anchoring — your data, verifiable and sovereign.
              </p>
              <div className="mt-6 overflow-hidden rounded-2xl shadow-soft transition-smooth group-hover:shadow-elegant">
                <img 
                  src={sovereignCloud} 
                  alt="Sovereign Cloud" 
                  className="w-full h-64 object-cover transition-smooth group-hover:scale-105"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Why D-Vault Section */}
      <section id="pricing" className="py-20 px-4">
        <div className="max-w-7xl mx-auto">
          <div className="grid md:grid-cols-2 gap-16 items-start">
            {/* Left Column */}
            <div>
              <h2 className="text-4xl md:text-5xl font-semibold mb-6 text-foreground">
                Why D-Vault?
              </h2>
              <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
                D-Vault combines the distributed permanence of IPFS with the immutable proof 
                layer of Stellar blockchain. Your data stays encrypted and private, while every 
                change is cryptographically anchored — giving you verifiable privacy without 
                compromise.
              </p>
              <Button 
                onClick={() => window.location.hash = '#pricing'}
                className="mb-6"
                size="lg"
              >
                Explore Plans
              </Button>
              <p className="text-accent font-medium">
                Your Data, Your Keys, Your Cloud
              </p>
            </div>

            {/* Right Column - Accordion */}
            <div>
              <Accordion type="single" collapsible className="w-full">
                <AccordionItem value="item-1" className="border-b border-border">
                  <AccordionTrigger className="text-lg font-medium hover:text-primary transition-smooth">
                    Decentralized by Design
                  </AccordionTrigger>
                  <AccordionContent className="text-muted-foreground leading-relaxed">
                    Built on peer-to-peer IPFS infrastructure, D-Vault eliminates central servers 
                    and single points of failure. Your data is distributed across a global network, 
                    accessible only to you, with no intermediaries.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-2" className="border-b border-border">
                  <AccordionTrigger className="text-lg font-medium hover:text-primary transition-smooth">
                    Blockchain Anchored Integrity
                  </AccordionTrigger>
                  <AccordionContent className="text-muted-foreground leading-relaxed">
                    Every vault entry is anchored to the Stellar blockchain, creating an immutable 
                    proof of existence and integrity. Verify any record at any time without exposing 
                    the underlying data.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-3" className="border-b border-border">
                  <AccordionTrigger className="text-lg font-medium hover:text-primary transition-smooth">
                    Tailored for Growth
                  </AccordionTrigger>
                  <AccordionContent className="text-muted-foreground leading-relaxed">
                    Start with local-only storage for personal use, scale to team collaboration 
                    with automatic sync, or deploy enterprise-grade infrastructure with dedicated 
                    nodes and compliance tools — all on the same sovereign architecture.
                  </AccordionContent>
                </AccordionItem>

                <AccordionItem value="item-4" className="border-b border-border">
                  <AccordionTrigger className="text-lg font-medium hover:text-primary transition-smooth">
                    AI Interpreter for Compliance
                  </AccordionTrigger>
                  <AccordionContent className="text-muted-foreground leading-relaxed">
                    Define compliance rules in plain language. Our AI interpreter translates your 
                    requirements into enforceable policies, with autonomous governance and audit 
                    trails anchored directly on-chain.
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 px-4 bg-card border-t border-border">
        <div className="max-w-7xl mx-auto">
          <div className="flex flex-col md:flex-row justify-between items-center gap-6">
            <p className="text-sm text-muted-foreground">
              D-Vault © 2025 — Built for the Self-Sovereign Web.
            </p>
            
            <div className="flex gap-6 text-sm">
              <a href="#docs" className="text-muted-foreground hover:text-primary transition-smooth">
                Docs
              </a>
              <a href="#api" className="text-muted-foreground hover:text-primary transition-smooth">
                API
              </a>
              <a href="#privacy" className="text-muted-foreground hover:text-primary transition-smooth">
                Privacy
              </a>
              <a href="https://github.com" className="text-muted-foreground hover:text-primary transition-smooth">
                GitHub
              </a>
            </div>
            
            <div className="flex gap-4">
              <a href="https://github.com" className="text-muted-foreground hover:text-primary transition-smooth">
                <Github className="h-5 w-5" />
              </a>
              <a href="https://twitter.com" className="text-muted-foreground hover:text-primary transition-smooth">
                <Twitter className="h-5 w-5" />
              </a>
              <a href="https://linkedin.com" className="text-muted-foreground hover:text-primary transition-smooth">
                <Linkedin className="h-5 w-5" />
              </a>
            </div>
          </div>
        </div>
      </footer>

      <OnboardingModal open={onboardingOpen} onOpenChange={setOnboardingOpen} />
    </div>
  );
}

export default Home;