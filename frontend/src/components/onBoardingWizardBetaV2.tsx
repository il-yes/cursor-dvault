import React, { useCallback, useEffect, useState } from 'react';
import { CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import {
  GetRecommendedTier,
  CreateAccount,
  SetupPaymentAndActivate,
  GetTierFeatures,
} from '../services/api';
import StellarKeyImport from './ImportStellarKey';

type Tier = 'free' | 'pro' | 'pro_plus' | 'business' | string;
type PaymentMethod = 'card' | 'stellar' | string;

interface TierFeatures {
  [tier: string]: {
    name: string;
    description?: string;
    features?: string[];
  };
}

interface CreateAccountResponse {
  user_id: string;
  secret_key?: string;
}

interface OnboardingWizardBetaProps {
  onComplete?: () => void;
}

const OnboardingWizardBeta: React.FC<OnboardingWizardBetaProps> = ({ onComplete }) => {
  const stripe = useStripe();
  const elements = useElements();

  const [step, setStep] = useState<number>(1);
  const [identity, setIdentity] = useState<string>('');
  const [useCases, setUseCases] = useState<string[]>([]);
  const [selectedTier, setSelectedTier] = useState<Tier>('');
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('');
  const [isAnonymous, setIsAnonymous] = useState<boolean>(false);
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [stellarSecretKey, setStellarSecretKey] = useState<string>('');
  const [userId, setUserId] = useState<string>('');
  const [tierFeatures, setTierFeatures] = useState<TierFeatures>({});
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [stellarKeyImported, setStellarKeyImported] = useState<boolean>(false);
  const [importedStellarKey, setImportedStellarKey] = useState<string | null>(null);
  const [cardError, setCardError] = useState<string | null>(null);
  const [cardLoading, setCardLoading] = useState(false);

  // UI-only: first/last name pour le payload
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');

  // Tier features
  useEffect(() => {
    let isMounted = true;
    (async () => {
      try {
        const features = await GetTierFeatures('free');
        if (isMounted) setTierFeatures(features);
      } catch (err: any) {
        console.error('Failed to load tier features', err);
      }
    })();
    return () => { isMounted = false; };
  }, []);

  const getTierPrice = useCallback((tier: Tier): number => {
    const prices: Record<string, number> = {
      free: 0,
      pro: 15,
      pro_plus: 25,
      business: 59,
    };
    return prices[tier] ?? 0;
  }, []);

  const selectIdentity = useCallback(async (choice: string) => {
    setIdentity(choice);
    try {
      const recommended = await GetRecommendedTier(choice);
      setSelectedTier(recommended);
    } catch (err: any) {
      console.error('GetRecommendedTier failed', err);
    }

    if (choice === 'anonymous') {
      setIsAnonymous(true);
      setStep(3);
    } else {
      setIsAnonymous(false);
      setStep(1.5);
    }
  }, []);

  const handleStellarKeyComplete = (data: any) => {
    if (data.stellar_key_imported) {
      setStellarKeyImported(true);
      setImportedStellarKey(data.stellar_secret_key);
      if (identity === 'anonymous') setIsAnonymous(true);
    }
    setStep(2);
  };

  const confirmUseCases = useCallback(() => {
    setStep(3);
  }, []);

  const selectTierAndPayment = useCallback((tier: Tier, method: PaymentMethod) => {
    setSelectedTier(tier);
    setPaymentMethod(method);
    setStep(4);
  }, []);

  const createAccount = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const response = (await CreateAccount({
        email: isAnonymous ? '' : email,
        name: isAnonymous ? '' : email.split('@')[0],
        password: isAnonymous ? '' : password,
        tier: selectedTier,
        is_anonymous: isAnonymous,
      })) as unknown as CreateAccountResponse;

      setUserId(response.user_id);

      if (isAnonymous) {
        if (response.secret_key) setStellarSecretKey(response.secret_key);
        setStep(4.5);
      } else {
        if (selectedTier === 'free') {
          setStep(6);
          onComplete?.();
        } else {
          setStep(5);
        }
      }
    } catch (err: any) {
      console.error(err);
      setError(err?.message || 'Failed to create account');
    } finally {
      setLoading(false);
    }
  }, [email, password, isAnonymous, selectedTier, onComplete]);

  const confirmSecretKeyBackup = useCallback(() => {
    if (selectedTier === 'free') {
      setStep(6);
      onComplete?.();
    } else {
      setStep(5);
    }
  }, [selectedTier, onComplete]);

  const setupPayment = useCallback(
    async (paymentData: any) => {
      setLoading(true);
      setError('');
      try {
        await SetupPaymentAndActivate({
          user_id: userId,
          tier: selectedTier,
          ...paymentData,
        });
        setStep(6);
        onComplete?.();
      } catch (err: any) {
        console.error(err);
        setError(err?.message || 'Payment setup failed');
      } finally {
        setLoading(false);
      }
    },
    [userId, selectedTier, onComplete],
  );

  const trialEndDate = new Date(
    Date.now() + 14 * 24 * 60 * 60 * 1000,
  ).toLocaleDateString();

  function prevStep() {
    if (step > 1) setStep(step - 1);
  }

  function handleStellarKeySkip() {
    setStellarKeyImported(false);
    setStep(2);
  }

  // -------- STRIPE FLOW (unique) --------

  const stripePaymentMethodHandler = async (result: any) => {
    if (result.error) {
      setCardError(result.error.message);
      return;
    }

    const payload = {
      product_id: 'pro',
      plan: selectedTier,
      payment_method: result.paymentMethod.id, // vrai pm_xxx Stripe
      email,
      last_four: result.paymentMethod.card.last4,
      card_brand: result.paymentMethod.card.brand,
      exp_month: result.paymentMethod.card.exp_month,
      exp_year: result.paymentMethod.card.exp_year,
      first_name: firstName,
      last_name: lastName,
      amount: getTierPrice(selectedTier),
    };

    try {
      const res = await fetch('/api/create-customer-and-subscribe-to-plan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      const data = await res.json();
      if (!data.error) {
        setStep(6);
        onComplete?.();
      } else {
        setCardError(data.message || 'Payment failed');
      }
    } catch (err: any) {
      setCardError(err?.message || 'Network error');
    }
  };

  const val = () => {
    const form = document.getElementById('charge_form') as HTMLFormElement | null;
    if (!form) return;

    if (!form.checkValidity()) {
      form.classList.add('was-validated');
      return;
    }

    if (!stripe || !elements) {
      setCardError('Stripe not loaded');
      return;
    }

    const cardElement = elements.getElement(CardElement);
    if (!cardElement) {
      setCardError('Card field not ready');
      return;
    }

    stripe
      .createPaymentMethod({
        type: 'card',
        card: cardElement,
        billing_details: { email },
      })
      .then(stripePaymentMethodHandler);
  };

  // -------- RENDER --------

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-white/90 backdrop-blur-3xl p-6">
      <div className="w-full max-w-4xl rounded-3xl bg-white/70 shadow-2xl border border-white/40 p-8 flex flex-col gap-6">
        {/* Header */}
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-3xl font-extrabold bg-gradient-to-r from-yellow-400 via-yellow-300 to-yellow-500 bg-clip-text text-transparent">
            Welcome to Ankhora
          </h1>
          <span className="text-sm text-muted-foreground">Step {step} of 6</span>
        </div>

        {error && (
          <div className="rounded-2xl bg-red-50 border border-red-200 text-red-700 px-4 py-3 text-sm">
            {error}
          </div>
        )}

        {/* STEP 1, 1.5, 2, 3, 4, 4.5 identiques à ta version,
            je ne les répète pas ici pour rester lisible.
            Tu peux garder exactement tes blocs existants pour ces steps. */}

        {/* STEP 5: Payment setup */}
        {step === 5 && (
          <div className="space-y-6">
            <h2 className="text-2xl font-semibold text-foreground">
              Set up your subscription
            </h2>
            <p className="text-sm text-muted-foreground">
              14-day free trial • Cancel anytime. You will not be charged until {trialEndDate}.
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {/* CARD PAYMENT */}
              <div className="space-y-4 p-6 rounded-3xl bg-white/80 border border-white/40 shadow-xl">
                <h3 className="text-xl font-bold mb-4">💳 Card Payment</h3>

                <form id="charge_form" noValidate className="space-y-4">
                  <input
                    id="cardholder-email"
                    type="email"
                    required
                    placeholder="Email"
                    className="w-full p-4 border rounded-xl bg-white"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                  <div className="grid grid-cols-2 gap-4">
                    <input
                      id="first_name"
                      required
                      placeholder="First name"
                      className="w-full p-4 border rounded-xl bg-white"
                      value={firstName}
                      onChange={(e) => setFirstName(e.target.value)}
                    />
                    <input
                      id="last_name"
                      required
                      placeholder="Last name"
                      className="w-full p-4 border rounded-xl bg-white"
                      value={lastName}
                      onChange={(e) => setLastName(e.target.value)}
                    />
                  </div>

                  <div className="p-4 border rounded-xl bg-white min-h-[60px]">
                    <CardElement />
                  </div>

                  {cardError && (
                    <div className="p-3 rounded-xl bg-red-50 border border-red-200 text-red-700 text-sm">
                      {cardError}
                    </div>
                  )}

                  <button
                    type="button"
                    onClick={val}
                    disabled={cardLoading}
                    className="w-full h-14 bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-bold shadow-xl hover:shadow-2xl transition-all disabled:opacity-60"
                  >
                    {cardLoading ? 'Processing…' : `Pay $${getTierPrice(selectedTier)}`}
                  </button>
                </form>
              </div>

              {/* STELLAR: garde ton widget ici */}
              <div className="space-y-4 p-6 rounded-3xl bg-white/80 border border-white/40 shadow-xl">
                <h3 className="text-xl font-bold mb-4">⭐ Stellar Wallet</h3>
                <button className="w-full h-14 border border-amber-300 bg-amber-50 hover:bg-amber-100">
                  Connect Stellar Wallet
                </button>
              </div>
            </div>
          </div>
        )}

        {/* STEP 6: Done */}
        {step === 6 && (
          <div className="space-y-4 text-center">
            <h2 className="text-2xl font-semibold text-foreground">Your vault is ready</h2>
            <p className="text-sm text-muted-foreground">
              Your zero-knowledge vault is now active and secured with Stellar verification.
            </p>
            <button
              onClick={onComplete}
              className="mt-2 px-6 py-2 rounded-xl bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-semibold shadow hover:shadow-lg transition"
            >
              Enter my vault
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default OnboardingWizardBeta;
