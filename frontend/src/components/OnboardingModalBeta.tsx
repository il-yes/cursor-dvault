import { useState } from "react";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Check, WifiOff, ExternalLink } from "lucide-react";
import { cn } from "@/lib/utils";
import { createVault } from "@/services/api";
import { toast } from "@/hooks/use-toast";
import { useConnectivity } from "@/hooks/useConnectivity";
import { useVault } from "@/hooks/useVault";
import { OfflineFallbackPanel } from "./OfflineFallbackPanel";
import { VaultContext } from "@/types/vault";
import "./contributionGraph/g-scrollbar.css";

type Plan = "freemium" | "pro" | "organization";

interface PlanOption {
  id: Plan;
  name: string;
  price: string;
  features: string[];
  popular?: boolean;
}

const plans: PlanOption[] = [
  {
    id: "freemium",
    name: "Shield",
    price: "Free",
    features: ["Local encryption", "Basic IPFS sync", "5 GB storage", "Single device"],
  },
  {
    id: "pro",
    name: "Sentinel",
    price: "$12/month",
    features: ["Multi-device sync", "Auto IPFS backup", "50 GB storage", "Stellar anchoring"],
    popular: true,
  },
  {
    id: "organization",
    name: "Fortress",
    price: "$49/month",
    features: ["Team collaboration", "Audit logs", "Unlimited storage", "Dedicated nodes", "Priority support"],
  },
];

interface OnboardingModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  upgradeMode?: boolean;
}

export function OnboardingModalBeta({ open, onOpenChange, upgradeMode = false }: OnboardingModalProps) {
  const [step, setStep] = useState<1 | 2>(upgradeMode ? 2 : 1);
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(upgradeMode ? 'pro' : null);
  const [vaultName, setVaultName] = useState("");
  const [useStellarKey, setUseStellarKey] = useState(false);
  const [publicKey, setPublicKey] = useState("");
  const [privateKey, setPrivateKey] = useState("");
  const [paymentName, setPaymentName] = useState("");
  const [paymentEmail, setPaymentEmail] = useState("");
  const [cardNumber, setCardNumber] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showOfflineFallback, setShowOfflineFallback] = useState(false);

  const { isOnline, onReconnect } = useConnectivity();
  const { hydrateVault } = useVault();

  const handleContinue = () => {
    if (selectedPlan) {
      setStep(2);
    }
  };

  const handleSubmit = async () => {
    if (!vaultName.trim() || !selectedPlan) return;

    if (selectedPlan !== "freemium" && (!paymentName || !paymentEmail || !cardNumber)) {
      toast({
        title: "Payment Required",
        description: "Please fill in all payment details.",
        variant: "destructive",
      });
      return;
    }

    if (!isOnline) {
      setShowOfflineFallback(true);
      return;
    }

    setIsSubmitting(true);
    try {
      const payload = {
        name: vaultName,
        plan: selectedPlan,
        stellarPublicKey: useStellarKey ? publicKey : undefined,
        stellarPrivateKey: useStellarKey ? privateKey : undefined,
        payment: selectedPlan !== "freemium" ? {
          name: paymentName,
          email: paymentEmail,
          cardNumber,
        } : undefined,
      };

      const response = await createVault(payload);

      if (response.vaultContext) {
        hydrateVault(response.vaultContext as VaultContext);
        toast({
          title: "Vault created successfully 🎉",
          description: "Your sovereign vault is ready to use.",
        });
        onOpenChange(false);
        resetForm();
      }
    } catch (error) {
      if (!isOnline || (error instanceof Error && error.message.includes('fetch'))) {
        setShowOfflineFallback(true);
      } else {
        toast({
          title: "Failed to create vault",
          description: error instanceof Error ? error.message : "Please try again.",
          variant: "destructive",
        });
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const resetForm = () => {
    setStep(upgradeMode ? 2 : 1);
    setSelectedPlan(upgradeMode ? 'pro' : null);
    setVaultName("");
    setUseStellarKey(false);
    setPublicKey("");
    setPrivateKey("");
    setPaymentName("");
    setPaymentEmail("");
    setCardNumber("");
    setShowOfflineFallback(false);
  };

  const handleRetryConnection = () => {
    setShowOfflineFallback(false);
    onReconnect(() => {
      toast({
        title: "Connection restored",
        description: "You can now continue with vault creation.",
      });
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl p-0 backdrop-blur-3xl border-white/20 shadow-2xl max-h-[90vh] overflow-auto scrollbar-glassmorphism thin-scrollbar">
        {/* Outer gradient shell */}
        <div className="h-full rounded-3xl p-[1px] bg-gradient-to-br from-white/30 via-white/20 to-zinc-100/10">
          {/* Inner glass container */}
          <div className="relative h-full rounded-[1.8rem] bg-gradient-to-b from-white/70 via-white/40 to-zinc-100/20 dark:from-zinc-900/70 dark:via-zinc-900/40 dark:to-black/20 backdrop-blur-3xl overflow-hidden">

            {/* STEP 1: Plan Selection */}
            {step === 1 && !upgradeMode && (
              <div className="h-full flex flex-col">
                {/* Header */}
                <div className="sticky top-0 z-10 border-b border-white/20 bg-white/20 backdrop-blur-xl px-8 py-6">
                  <h2 className="text-2xl font-black bg-gradient-to-r from-[#C9A44A] via-amber-400 to-[#B8934A] bg-clip-text text-transparent">
                    Own your vault. Choose your sovereignty.
                  </h2>
                </div>

                {/* Plans */}
                <div className="flex-1 p-8 overflow-y-auto glass-scrollbar">
                  <div className="grid md:grid-cols-3 gap-6">
                    {plans.map((plan) => (
                      <div
                        key={plan.id}
                        onClick={() => setSelectedPlan(plan.id)}
                        className={cn(
                          "group relative p-8 rounded-3xl border-2 transition-all cursor-pointer text-left backdrop-blur-xl shadow-xl hover:shadow-2xl hover:shadow-primary/30 hover:scale-[1.02]",
                          selectedPlan === plan.id
                            ? "border-[#C9A44A]/70 bg-gradient-to-br from-[#C9A44A]/20 to-[#B8934A]/10 ring-2 ring-[#C9A44A]/40"
                            : "border-white/30 bg-white/20 dark:bg-zinc-900/40 hover:border-[#C9A44A]/50"
                        )}
                      >
                        {plan.popular && (
                          <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-gradient-to-r from-[#C9A44A] to-[#B8934A] text-black text-xs font-bold px-4 py-1 rounded-full shadow-lg">
                            Popular
                          </div>
                        )}
                        {selectedPlan === plan.id && (
                          <div className="absolute top-6 right-6 w-3 h-3 rounded-full bg-[#C9A44A] ring-2 ring-white/50 animate-ping" />
                        )}
                        <h3 className="text-xl font-bold mb-3 text-white">{plan.name}</h3>
                        <p className="text-2xl font-black text-[#C9A44A] mb-6">{plan.price}</p>
                        <div className="space-y-2 mb-8">
                          {plan.features.map((feature, idx) => (
                            <div key={idx} className="flex items-center gap-2 text-sm text-zinc-200">
                              <div className="w-5 h-5 rounded-full bg-white/20 flex items-center justify-center">
                                <Check className="w-3 h-3 text-[#C9A44A]" />
                              </div>
                              {feature}
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Continue Button */}
                <div className="sticky bottom-0 border-t border-white/20 bg-white/30 backdrop-blur-xl px-8 py-6">
                  <Button
                    onClick={handleContinue}
                    disabled={!selectedPlan}
                    className="w-full h-14 rounded-2xl bg-gradient-to-r from-[#C9A44A] via-amber-400 to-[#B8934A] text-black font-bold shadow-2xl hover:shadow-[#C9A44A]/40 transition-all hover:scale-[1.02] backdrop-blur-xl"
                  >
                    Continue →
                  </Button>
                </div>
              </div>
            )}

            {/* STEP 2: Vault Setup */}
            {step === 2 && !showOfflineFallback && (
              <div className="h-full flex flex-col">
                {/* Header */}
                <div className="sticky top-0 z-10 border-b border-white/20 bg-white/20 backdrop-blur-xl px-8 py-6">
                  <div className="flex items-center justify-between">
                    <button
                      onClick={() => setStep(1)}
                      className="inline-flex items-center gap-2 text-sm font-medium text-[#C9A44A] hover:text-[#B8934A]"
                    >
                      ← Back
                    </button>
                    <h2 className="text-2xl font-black bg-gradient-to-r from-[#C9A44A] via-amber-400 to-[#B8934A] bg-clip-text text-transparent">
                      Setup Your {plans.find(p => p.id === selectedPlan)?.name} Vault
                    </h2>
                    <div />
                  </div>
                </div>

                {/* Content */}
                <div className="flex-1 p-8 overflow-y-auto glass-scrollbar space-y-8">
                  <div className="grid md:grid-cols-2 gap-8">
                    {/* Vault Setup */}
                    <div>
                      <h3 className="text-lg font-semibold text-white mb-6">Vault Configuration</h3>
                      <div className="space-y-4">
                        <div>
                          <Label className="text-sm font-medium text-zinc-300 mb-2 block">Vault Name *</Label>
                          <Input
                            value={vaultName}
                            onChange={(e) => setVaultName(e.target.value)}
                            className="h-14 rounded-2xl bg-white/40 backdrop-blur-sm border-white/30 text-lg font-semibold shadow-inner focus:ring-[#C9A44A]/40"
                            placeholder="My Sovereign Vault"
                          />
                        </div>
                        <div className="flex items-center space-x-3">
                          <Switch
                            checked={useStellarKey}
                            onCheckedChange={setUseStellarKey}
                            className="data-[state=checked]:bg-[#C9A44A]"
                          />
                          <Label className="text-sm font-medium text-zinc-300 cursor-pointer">
                            Secure with existing Stellar key
                          </Label>
                        </div>
                        {useStellarKey && (
                          <div className="space-y-4 pt-4 border-t border-white/20 p-4 rounded-2xl bg-white/10">
                            <div>
                              <Label className="text-sm font-medium text-zinc-300 mb-2 block">Public Key</Label>
                              <Input
                                value={publicKey}
                                onChange={(e) => setPublicKey(e.target.value)}
                                className="h-12 rounded-xl bg-white/30 backdrop-blur-sm border-white/30 font-mono text-sm shadow-inner"
                              />
                            </div>
                            <div>
                              <Label className="text-sm font-medium text-zinc-300 mb-2 block">Private Key</Label>
                              <Input
                                value={privateKey}
                                onChange={(e) => setPrivateKey(e.target.value)}
                                className="h-12 rounded-xl bg-white/30 backdrop-blur-sm border-white/30 font-mono text-sm shadow-inner"
                                type="password"
                              />
                            </div>
                            <p className="text-xs text-zinc-400 italic">
                              Keys are encrypted locally and never stored in plaintext
                            </p>
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Payment */}
                    {selectedPlan !== "freemium" && (
                      <div>
                        <h3 className="text-lg font-semibold text-white mb-6">Payment Details</h3>
                        <div className="space-y-4">
                          <div>
                            <Label className="text-sm font-medium text-zinc-300 mb-2 block">Full Name *</Label>
                            <Input
                              value={paymentName}
                              onChange={(e) => setPaymentName(e.target.value)}
                              className="h-14 rounded-2xl bg-white/40 backdrop-blur-sm border-white/30 text-lg font-semibold shadow-inner focus:ring-[#C9A44A]/40"
                            />
                          </div>
                          <div>
                            <Label className="text-sm font-medium text-zinc-300 mb-2 block">Email *</Label>
                            <Input
                              value={paymentEmail}
                              onChange={(e) => setPaymentEmail(e.target.value)}
                              className="h-14 rounded-2xl bg-white/40 backdrop-blur-sm border-white/30 text-lg font-semibold shadow-inner focus:ring-[#C9A44A]/40"
                            />
                          </div>
                          <div>
                            <Label className="text-sm font-medium text-zinc-300 mb-2 block">Card Number *</Label>
                            <Input
                              value={cardNumber}
                              onChange={(e) => setCardNumber(e.target.value)}
                              className="h-14 rounded-2xl bg-white/40 backdrop-blur-sm border-white/30 text-lg font-mono font-semibold shadow-inner focus:ring-[#C9A44A]/40"
                            />
                          </div>
                          <p className="text-xs text-zinc-400 italic pt-2">
                            Secure payment processed via Stripe
                          </p>
                        </div>
                      </div>
                    )}
                  </div>
                </div>

                {/* Submit Buttons */}
                <div className="sticky bottom-0 border-t border-white/20 bg-white/30 backdrop-blur-xl px-8 py-6 space-y-3">
                  <Button
                    variant="outline"
                    onClick={() => {
                      toast({
                        title: "Draft saved",
                        description: "You can resume onboarding when online.",
                      });
                      onOpenChange(false);
                      resetForm();
                    }}
                    className="w-full h-14 rounded-2xl bg-white/50 border-white/40 text-zinc-700 hover:bg-white/70 shadow-lg font-semibold backdrop-blur-xl"
                  >
                    Save Draft
                  </Button>
                  <Button
                    onClick={handleSubmit}
                    disabled={isSubmitting || !vaultName.trim()}
                    className="w-full h-14 rounded-2xl bg-gradient-to-r from-[#C9A44A] via-amber-400 to-[#B8934A] text-black font-bold shadow-2xl hover:shadow-[#C9A44A]/40 transition-all hover:scale-[1.02] backdrop-blur-xl"
                  >
                    {isSubmitting ? "Creating..." : "Create Vault"}
                  </Button>
                </div>
              </div>
            )}

            {/* Offline Fallback */}
            {showOfflineFallback && (
              <div className="h-full flex flex-col items-center justify-center p-12 text-center space-y-8 backdrop-blur-xl">
                <div className="w-24 h-24 rounded-3xl bg-gradient-to-br from-zinc-500/30 to-zinc-700/40 flex items-center justify-center shadow-2xl">
                  <WifiOff className="w-12 h-12 text-zinc-400" />
                </div>
                <div className="space-y-4">
                  <h3 className="text-2xl font-black text-white">Connection Required</h3>
                  <p className="text-lg text-zinc-400 max-w-md mx-auto">
                    Vault creation needs internet. Your draft is saved locally.
                  </p>
                </div>
                <div className="flex flex-col sm:flex-row gap-4 w-full max-w-md">
                  <Button
                    onClick={handleRetryConnection}
                    className="h-14 flex-1 rounded-2xl bg-gradient-to-r from-[#C9A44A] to-[#B8934A] text-black font-bold shadow-xl hover:shadow-[#C9A44A]/40"
                  >
                    Retry Connection
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => {
                      toast({
                        title: "Draft saved",
                        description: "You can resume onboarding when online.",
                      });
                      onOpenChange(false);
                      resetForm();
                    }}
                    className="h-14 flex-1 rounded-2xl bg-white/50 border-white/40 text-zinc-700 hover:bg-white/70"
                  >
                    Continue Offline
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
