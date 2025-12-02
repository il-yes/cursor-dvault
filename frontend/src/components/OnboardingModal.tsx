import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Check } from "lucide-react";
import { createVault } from "@/services/api";
import { toast } from "@/hooks/use-toast";
import { useConnectivity } from "@/hooks/useConnectivity";
import { useVault } from "@/hooks/useVault";
import { OfflineFallbackPanel } from "./OfflineFallbackPanel";
import { VaultContext } from "@/types/vault";
import "./contributionGraph/g-scrollbar.css";

interface OnboardingModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  upgradeMode?: boolean; // Skip step 1 and go directly to payment
}

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

export function OnboardingModal({ open, onOpenChange, upgradeMode = false }: OnboardingModalProps) {
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

    // Check connectivity
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

      // Hydrate vault context if provided
      if (response.vaultContext) {
        hydrateVault(response.vaultContext as VaultContext);
      }

      toast({
        title: "Vault created successfully 🎉",
        description: "Your sovereign vault is ready to use.",
      });

      onOpenChange(false);
      resetForm();
    } catch (error) {
      // Check if error is network-related
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

  const vaultDraft = {
    user_id: 'pending',
    Vault: {
      version: '1.0',
      name: vaultName,
      folders: [],
      entries: { login: [], card: [], note: [], sshkey: [], identity: [] },
    },
    Dirty: true,
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-5xl max-h-[90vh] overflow-y-auto scrollbar-glassmorphism thin-scrollbar">
        {step === 1 && !upgradeMode ? (
          <div className="space-y-8 py-6">
            <DialogHeader>
              <DialogTitle className="text-3xl font-light text-center">
                Own your vault. Choose your sovereignty.
              </DialogTitle>
            </DialogHeader>

            <div className="grid md:grid-cols-3 gap-6">
              {plans.map((plan) => (
                <button
                  key={plan.id}
                  onClick={() => setSelectedPlan(plan.id)}
                  className={`relative p-6 rounded-2xl border-2 transition-all text-left hover:shadow-lg ${
                    selectedPlan === plan.id
                      ? "border-[#00cfcf] bg-[#00cfcf]/5"
                      : "border-border hover:border-[#00cfcf]/50"
                  }`}
                >
                  {plan.popular && (
                    <div className="absolute -top-3 left-1/2 -translate-x-1/2 px-3 py-1 bg-[#FD871F] text-white text-xs rounded-full">
                      Popular
                    </div>
                  )}
                  
                  {selectedPlan === plan.id && (
                    <div className="absolute top-4 right-4 w-6 h-6 rounded-full bg-[#00cfcf] flex items-center justify-center">
                      <Check className="w-4 h-4 text-white" />
                    </div>
                  )}

                  <h3 className="text-xl font-semibold mb-2">{plan.name}</h3>
                  <p className="text-2xl font-light text-[#00cfcf] mb-4">{plan.price}</p>
                  
                  <ul className="space-y-2">
                    {plan.features.map((feature, idx) => (
                      <li key={idx} className="text-sm text-muted-foreground flex items-start gap-2">
                        <Check className="w-4 h-4 text-[#00cfcf] mt-0.5 flex-shrink-0" />
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>
                </button>
              ))}
            </div>

            <div className="flex justify-center pt-4">
              <Button
                onClick={handleContinue}
                disabled={!selectedPlan}
                className="px-8 h-12 text-base bg-[#00cfcf] hover:bg-[#00cfcf]/90 text-white"
              >
                Continue →
              </Button>
            </div>
          </div>
        ) : showOfflineFallback ? (
          <div className="space-y-6 py-6">
            <DialogHeader>
              <DialogTitle className="text-2xl font-light">Connection Required</DialogTitle>
            </DialogHeader>

            <OfflineFallbackPanel
              vaultDraft={vaultDraft}
              onRetry={handleRetryConnection}
              onSaveSuccess={() => {
                toast({
                  title: "Draft saved",
                  description: "You can resume onboarding when online.",
                });
                onOpenChange(false);
                resetForm();
              }}
            />

            <div className="flex justify-between pt-4 border-t">
              <Button
                variant="outline"
                onClick={() => setShowOfflineFallback(false)}
              >
                ← Back
              </Button>
            </div>
          </div>
        ) : (
          <div className="space-y-6 py-6">
            <DialogHeader>
              <DialogTitle className="text-2xl font-light">Setup Your Vault</DialogTitle>
            </DialogHeader>

            <div className="grid md:grid-cols-2 gap-8">
              {/* Left Side - Vault Setup */}
              <div className="space-y-6">
                <div className="space-y-2">
                  <Label htmlFor="vaultName">Vault Name *</Label>
                  <Input
                    id="vaultName"
                    placeholder="My Sovereign Vault"
                    value={vaultName}
                    onChange={(e) => setVaultName(e.target.value)}
                    className="h-11"
                  />
                </div>

                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div className="space-y-0.5">
                    <Label className="text-base">Secure with existing Stellar key</Label>
                    <p className="text-sm text-muted-foreground">
                      Use your own Stellar account for anchoring
                    </p>
                  </div>
                  <Switch
                    checked={useStellarKey}
                    onCheckedChange={setUseStellarKey}
                  />
                </div>

                {useStellarKey && (
                  <div className="space-y-4 pt-2">
                    <div className="space-y-2">
                      <Label htmlFor="publicKey">Public Key</Label>
                      <Input
                        id="publicKey"
                        placeholder="G..."
                        value={publicKey}
                        onChange={(e) => setPublicKey(e.target.value)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="privateKey">Private Key</Label>
                      <Input
                        id="privateKey"
                        type="password"
                        placeholder="S..."
                        value={privateKey}
                        onChange={(e) => setPrivateKey(e.target.value)}
                      />
                    </div>
                    <p className="text-xs text-muted-foreground">
                      Keys are encrypted locally and never stored in plaintext
                    </p>
                  </div>
                )}
              </div>

              {/* Right Side - Payment (only for paid plans) */}
              {selectedPlan !== "freemium" && (
                <div className="space-y-6 md:border-l md:pl-8">
                  <h3 className="text-lg font-medium">Payment Details</h3>
                  
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="paymentName">Full Name *</Label>
                      <Input
                        id="paymentName"
                        placeholder="John Doe"
                        value={paymentName}
                        onChange={(e) => setPaymentName(e.target.value)}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="paymentEmail">Email *</Label>
                      <Input
                        id="paymentEmail"
                        type="email"
                        placeholder="john@example.com"
                        value={paymentEmail}
                        onChange={(e) => setPaymentEmail(e.target.value)}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="cardNumber">Card Number *</Label>
                      <Input
                        id="cardNumber"
                        placeholder="4242 4242 4242 4242"
                        value={cardNumber}
                        onChange={(e) => setCardNumber(e.target.value)}
                      />
                    </div>

                    <p className="text-xs text-muted-foreground">
                      Secure payment processed via Stripe
                    </p>
                  </div>
                </div>
              )}
            </div>

            <div className="flex justify-between pt-4 border-t">
              <Button
                variant="outline"
                onClick={() => setStep(1)}
                disabled={isSubmitting}
              >
                ← Back
              </Button>
              <Button
                onClick={handleSubmit}
                disabled={!vaultName.trim() || isSubmitting}
                className="px-8 bg-[#00cfcf] hover:bg-[#00cfcf]/90 text-white"
              >
                {isSubmitting ? "Creating..." : "Create Vault"}
              </Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
