// OnboardingWizardBeta.tsx (React version)

import React, { useCallback, useEffect, useRef, useState } from 'react';
// Adjust import path to your generated Wails bindings
import {
    GetRecommendedTier,
    CreateAccount,
    SetupPaymentAndActivate,
    GetTierFeatures,
} from '../services/api';
import StellarKeyImport from './ImportStellarKey';
import { CardElement, useElements, useStripe } from '@stripe/react-stripe-js';
import StripePayButton from './StripePayButton';
import * as AppAPI from "../../wailsjs/go/main/App";
import { Button } from './ui/button';


type Tier = 'free' | 'pro' | 'pro_plus' | 'business' | string;
type PaymentMethod = 'card' | 'stellar' | string;

interface TierFeatures {
    [tier: string]: {
        name?: string;
        description?: string;
        features?: string[];
    };
}

interface SecurityStats {
    // extend as needed
}

interface CreateAccountResponse {
    user_id: string;
    secret_key?: string;
}

interface OnboardingWizardBetaProps {
    onComplete?: () => void; // called when step 6 is reached (vault ready)
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
    const [stellarSecretKey, setStellarSecretKey] = useState<string>(''); // must be shown & saved
    const [userId, setUserId] = useState<string>('');
    const [tierFeatures, setTierFeatures] = useState<TierFeatures>({});
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>('');
    const [stellarKeyImported, setStellarKeyImported] = useState<boolean>(false);
    const [importedStellarKey, setImportedStellarKey] = useState<string | null>(null);
    const [cardNumber, setCardNumber] = useState<string>('');
    const [exp, setExp] = useState<string>('');
    const [cvc, setCvc] = useState<string>('');
    const [firstName, setFirstName] = useState<string>('');
    const [lastName, setLastName] = useState<string>('');

    const [cardError, setCardError] = useState<string | null>(null);
    const [cardLoading, setCardLoading] = useState(false);


    // ✅ Add these refs + types at top of component (after hooks):
    const formRef = useRef<HTMLFormElement>(null);
    const emailRef = useRef<HTMLInputElement>(null);
    const firstNameRef = useRef<HTMLInputElement>(null);
    const lastNameRef = useRef<HTMLInputElement>(null);

    // const cardElement = elements?.getElement(CardElement);


    // Load tier features on mount
    useEffect(() => {
        let isMounted = true;
        (async () => {
            try {
                const features = await GetTierFeatures('free');
                if (isMounted) {
                    console.log("Tier features:", { features });
                    setTierFeatures(features);
                }
            } catch (err: any) {
                console.error('Failed to load tier features', err);
                if (isMounted) setTierFeatures({});
            }
        })();
        return () => {
            isMounted = false;
        };
    }, []);

    // Helper: tier price
    const getTierPrice = useCallback((tier: Tier): number => {
        const prices: Record<string, number> = {
            free: 0,
            pro: 15,
            pro_plus: 25,
            business: 59,
        };
        return prices[tier] ?? 0;
    }, []);

    // Step 1: Identity selection
    const selectIdentity = useCallback(
        async (choice: string) => {
            setIdentity(choice);
            try {
                const recommended = await GetRecommendedTier(choice);
                setSelectedTier(recommended);
            } catch (err: any) {
                console.error('GetRecommendedTier failed', err);
            }

            if (choice === 'anonymous') {
                setIsAnonymous(true);
                setStep(3); // skip use cases, go to tier selection
            } else {
                setIsAnonymous(false);
                setStep(1.5);
            }
        },
        [],
    );

    // Step 1.5: Stellar key backup
    const handleStellarKeyComplete = (data) => {
        if (data.stellar_key_imported) {
            setStellarKeyImported(true);
            setImportedStellarKey(data.stellar_secret_key);

            // If importing for Pro Plus anonymous account, skip email/password
            if (identity === 'anonymous') {
                setIsAnonymous(true);
            }
        }

        // Continue to use case selection
        setStep(2);
    }

    // Step 2: Use cases (optional)
    const confirmUseCases = useCallback(() => {
        setStep(3);
    }, []);

    // Step 3: Tier + payment selection
    const selectTierAndPayment = useCallback((tier: Tier, method: PaymentMethod) => {
        setSelectedTier(tier);
        setPaymentMethod(method);
        setStep(4);
    }, []);

    // Step 4: Account creation
    const createAccount = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const response = (await CreateAccount({
                email: isAnonymous ? '' : email,
                name: isAnonymous ? '' : email.split('@')[0], // Use email prefix as name, or empty for anonymous
                password: isAnonymous ? '' : password,
                tier: selectedTier,
                is_anonymous: isAnonymous,
            })) as unknown as CreateAccountResponse;

            setUserId(response.user_id);

            if (isAnonymous) {
                // Anonymous: show Stellar secret key & require backup
                if (response.secret_key) {
                    setStellarSecretKey(response.secret_key);
                }
                setStep(4.5); // secret key backup step
            } else {
                // Identified account: branch based on tier
                if (selectedTier === 'free') {
                    setStep(6);
                    onComplete?.();
                } else {
                    setStep(5); // payment setup
                }
            }
        } catch (err: any) {
            console.error(err);
            setError(err?.message || 'Failed to create account');
        } finally {
            setLoading(false);
        }
    }, [email, password, isAnonymous, selectedTier, onComplete]);

    // Step 4.5: Secret key backup (anonymous only)
    const confirmSecretKeyBackup = useCallback(() => {
        if (selectedTier === 'free') {
            setStep(6);
            onComplete?.();
        } else {
            setStep(5);
        }
    }, [selectedTier, onComplete]);

    // Step 5: Payment setup
    const setupPayment = useCallback(
        async (paymentData: any) => {
            setLoading(true);
            setError('');

            console.log(({
                user_id: userId,
                tier: selectedTier,
                ...paymentData
            }))

            try {
                await SetupPaymentAndActivate({
                    user_id: userId,
                    tier: selectedTier,
                    ...paymentData
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

    const handleCardPayment = async () => {
        if (!cardNumber || cardNumber.length < 13 || !exp || !cvc) {
            setCardError('Please enter valid card details');
            return;
        }

        setCardLoading(true);
        // Mock payment method ID for testing
        await setupPayment({
            payment_method_id: 'pm_mock_' + Date.now(), // Your backend handles this
            card_number: cardNumber.replace(/\s/g, ''),
            exp_month: exp.split('/')[0],
            exp_year: exp.split('/')[1],
            cvc: cvc,
        });
    };


    // 14-day trial date string
    const trialEndDate = new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toLocaleDateString();

    function prevStep() {
        if (step > 1) {
            setStep(step - 1);
        }
    }

    function handleStellarKeySkip() {
        setStellarKeyImported(false);

        // Continue to use case selection
        setStep(2);
    }

    const cardElement = elements?.getElement(CardElement);

    // Form submit handler - IDENTICAL TO YOUR JS
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        const form = formRef.current;
        if (!form || !form.checkValidity()) {  // ✅ formRef.current = HTMLFormElement
            e.stopPropagation();
            form?.classList.add("was-validated");
            return;
        }

        form.classList.add("was-validated");

        if (!stripe || !cardElement) {
            setCardError('Stripe not loaded');
            return;
        }

        stripe.createPaymentMethod({
            type: 'card',
            card: cardElement,
            billing_details: {
                email: emailRef.current?.value || '',
            },
        }).then(stripePaymentMethodHandler);
    };


    // YOUR JS → React (IDENTICAL)
    const val = () => {
        console.log("🚀 val() appelée");

        const cardElement = elements.getElement(CardElement);
        if (!cardElement) {
            console.log("❌ CardElement vide");
            return;
        }


        const form = document.getElementById("charge_form") as HTMLFormElement;
        if (!form.checkValidity()) {
            form.classList.add("was-validated");
            return;
        }
        console.log("✅ Création payment_method...");

        stripe!.createPaymentMethod({
            type: 'card',
            card: elements!.getElement(CardElement)!,
            billing_details: {
                email: (document.getElementById("cardholder-email") as HTMLInputElement).value,
            },
        }).then(stripePaymentMethodHandler);
    };

    const stripePaymentMethodHandler = async (result) => {
        if (result.error) {
            setError(result.error.message);
            return;
        }

        // REAL STRIPE payment_method.id - NO MOCK
        const payload = {
            product_id: "pro",
            plan: selectedTier,
            payment_method: result.paymentMethod.id,  // pm_1ABC123 - REAL
            email: (document.getElementById("cardholder-email") as HTMLInputElement).value,
            last_four: result.paymentMethod.card.last4,
            card_brand: result.paymentMethod.card.brand,
            exp_month: result.paymentMethod.card.exp_month,
            exp_year: result.paymentMethod.card.exp_year,
            first_name: (document.getElementById("first_name") as HTMLInputElement).value,
            last_name: (document.getElementById("last_name") as HTMLInputElement).value,
            amount: getTierPrice(selectedTier),
        };

        const response = await fetch("/api/create-customer-and-subscribe-to-plan", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });

        const data = await response.json();
        if (!data.error) {
            setStep(6);
        }
    };


    const handleClick = async () => {
        await AppAPI.OpenGoogle();
    };

    const onPaymentSuccess = () => {
        setStep(6);
    };

    const handleSkip = () => {
        setStep(6);
    };

    const GetSession = async () => {
        const session = await AppAPI.GetSession(userId);
        console.log("Session:", session);
    };  

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

                {/* STEP 1: Identity */}
                {step === 1 && (
                    <div className="space-y-6">
                        <h2 className="text-2xl font-semibold text-foreground">
                            How do you want to use Ankhora?
                        </h2>
                        <p className="text-muted-foreground">
                            All plans include zero-knowledge encryption and Stellar verification.
                        </p>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <button
                                onClick={() => selectIdentity('personal')}
                                className="rounded-2xl p-5 border border-white/60 bg-white/70 hover:bg-white shadow-sm hover:shadow-xl transition"
                            >
                                <div className="font-semibold mb-1">Personal / Professional</div>
                                <div className="text-sm text-muted-foreground">
                                    Use Ankhora with an email, password, and traditional account.
                                </div>
                            </button>
                            <button
                                onClick={() => selectIdentity('anonymous')}
                                className="rounded-2xl p-5 border border-amber-300 bg-amber-50/80 hover:bg-amber-100 shadow-sm hover:shadow-xl transition"
                            >
                                <div className="font-semibold mb-1">Fully Anonymous</div>
                                <div className="text-sm text-muted-foreground">
                                    “Even your subscription is anonymous”. A Stellar keypair is generated for you.
                                </div>
                            </button>
                        </div>
                        <button onClick={() => setStep(1.5)}>Stellar Key backup</button>

                        <button onClick={GetSession}>Get Session</button>

                        <Button
                            onClick={handleSkip}
                            size="lg"
                            className="h-14 px-12 text-lg bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black shadow-lg hover:shadow-xl transition-all font-semibold"
                        >
                            Continue with New Account
                        </Button>
                    </div>
                )}

                {/* STEP 1.5: Stellar key backup */}
                {step === 1.5 && (
                    <StellarKeyImport
                        onComplete={handleStellarKeyComplete}
                        onSkip={handleStellarKeySkip}
                        onBack={() => setStep(1)}
                    />
                )}

                {/* STEP 2: Use cases (optional) */}
                {step === 2 && (
                    <div className="space-y-6">
                        <h2 className="text-2xl font-semibold text-foreground">
                            Choose the option that best describes your needs
                        </h2>
                        <p className="text-muted-foreground text-sm">Select all that apply.</p>

                        {/* Example use case chips */}
                        <div className="flex flex-wrap gap-3">
                            {['Personal backups', 'Team collaboration', 'Client data', 'Compliance-critical'].map(
                                (u) => (
                                    <button
                                        key={u}
                                        type="button"
                                        onClick={() =>
                                            setUseCases((prev) =>
                                                prev.includes(u) ? prev.filter((x) => x !== u) : [...prev, u],
                                            )
                                        }
                                        className={`px-4 py-2 rounded-full text-sm border transition ${useCases.includes(u)
                                            ? 'bg-amber-100 border-amber-400 text-amber-900'
                                            : 'bg-white/70 border-zinc-200 text-zinc-700 hover:bg-white'
                                            }`}
                                    >
                                        {u}
                                    </button>
                                ),
                            )}
                        </div>

                        <div className="flex justify-between">
                            {step > 1 && (
                                <button
                                    onClick={prevStep}
                                    className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
                                >
                                    Back
                                </button>
                            )}
                            <button
                                onClick={confirmUseCases}
                                className="px-6 py-2 rounded-xl bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-semibold shadow hover:shadow-lg transition"
                            >
                                Continue
                            </button>
                        </div>
                    </div>
                )}

                {/* STEP 3: Tier + payment selection */}
                {step === 3 && (
                    <div className="space-y-6">
                        <h2 className="text-2xl font-semibold text-foreground mb-2">
                            Choose your tier
                        </h2>
                        <p className="text-muted-foreground text-sm">
                            All plans include zero-knowledge encryption and Stellar verification.
                        </p>

                        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                            {(['free', 'pro', 'pro_plus', 'business'] as Tier[]).map((tier) => (
                                <button
                                    key={tier}
                                    type="button"
                                    onClick={() => selectTierAndPayment(tier, 'card')}
                                    className={`flex flex-col items-start rounded-2xl border p-4 text-left transition ${selectedTier === tier
                                        ? 'border-yellow-400 bg-yellow-50/80 shadow-lg'
                                        : 'border-zinc-200 bg-white/70 hover:bg-white shadow-sm hover:shadow-lg'
                                        }`}
                                >
                                    <div className="font-semibold capitalize">{tier.replace('_', ' ')}</div>
                                    <div className="text-sm text-muted-foreground">
                                        ${getTierPrice(tier)} / month
                                    </div>
                                    <div className="mt-2 text-xs text-muted-foreground">
                                        {(tierFeatures[tier]?.description) || 'Secure, encrypted storage.'}
                                    </div>
                                </button>
                            ))}
                        </div>

                        <div className="mt-4 text-xs text-muted-foreground">
                            14-day free trial • Cancel anytime. You will not be charged until {trialEndDate}.
                        </div>
                        {step > 1 && (
                            <button
                                onClick={prevStep}
                                className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
                            >
                                Back
                            </button>
                        )}
                    </div>
                )}

                {/* STEP 4: Account creation (email/password or anon) */}
                {step === 4 && (
                    <div className="space-y-6">
                        {!isAnonymous ? (
                            <>
                                <h2 className="text-2xl font-semibold text-foreground">
                                    Create your Ankhora account
                                </h2>
                                <div className="space-y-4">
                                    <div>
                                        <label className="block text-sm font-medium mb-1">Email</label>
                                        <input
                                            type="email"
                                            value={email}
                                            onChange={(e) => setEmail(e.target.value)}
                                            className="w-full rounded-xl border border-zinc-200 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-amber-300 bg-white/80"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium mb-1">Password</label>
                                        <input
                                            type="password"
                                            value={password}
                                            onChange={(e) => setPassword(e.target.value)}
                                            className="w-full rounded-xl border border-zinc-200 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-amber-300 bg-white/80"
                                        />
                                    </div>
                                </div>
                                <div className="flex justify-between">
                                    {step > 1 && (
                                        <button
                                            onClick={prevStep}
                                            className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
                                        >
                                            Back
                                        </button>
                                    )}
                                    <button
                                        disabled={loading}
                                        onClick={createAccount}
                                        className="px-6 py-2 rounded-xl bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-semibold shadow hover:shadow-lg disabled:opacity-60 transition"
                                    >
                                        {loading ? 'Creating account…' : 'Create account'}
                                    </button>
                                </div>
                            </>
                        ) : (
                            <>
                                <h2 className="text-2xl font-semibold text-foreground">
                                    Creating your anonymous vault
                                </h2>
                                <p className="text-sm text-muted-foreground">
                                    “Even your subscription is anonymous”. We will generate a Stellar keypair for you.
                                </p>
                                <div className="flex justify-end">
                                    {step > 1 && (
                                        <button
                                            onClick={prevStep}
                                            className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
                                        >
                                            Back
                                        </button>
                                    )}
                                    <button
                                        disabled={loading}
                                        onClick={createAccount}
                                        className="px-6 py-2 rounded-xl bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-semibold shadow hover:shadow-lg disabled:opacity-60 transition"
                                    >
                                        {loading ? 'Creating anonymous account…' : 'Generate Stellar keypair'}
                                    </button>
                                </div>
                            </>
                        )}
                    </div>
                )}

                {/* STEP 4.5: Secret key backup */}
                {step === 4.5 && (
                    <div className="space-y-6">
                        <h2 className="text-2xl font-semibold text-red-600">
                            Save your Stellar secret key
                        </h2>
                        <p className="text-sm text-red-700">
                            You <strong>must</strong> save your secret key — it cannot be recovered by Ankhora.
                            Without it, you will permanently lose access to your vault.
                        </p>
                        <div className="rounded-2xl bg-zinc-900 text-amber-300 font-mono text-sm px-4 py-3 overflow-x-auto">
                            {stellarSecretKey || '••••••••••••••••••••••••••'}
                        </div>
                        <ul className="text-xs text-muted-foreground list-disc pl-5 space-y-1">
                            <li>This is your <strong>only</strong> way to access your vault.</li>
                            <li>Store it in a safe place—preferably offline (password manager, hardware key, etc.).</li>
                        </ul>
                        <div className="flex justify-end">
                            {step > 1 && (
                                <button
                                    onClick={prevStep}
                                    className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
                                >
                                    Back
                                </button>
                            )}
                            <button
                                onClick={confirmSecretKeyBackup}
                                className="px-6 py-2 rounded-xl bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-semibold shadow hover:shadow-lg transition"
                            >
                                I have safely stored my secret key
                            </button>
                        </div>
                    </div>
                )}

                {/* STEP 5: Payment setup (card / Stellar) */}

                {step === 5 && (
                    <div className="space-y-6">
                        <h2 className="text-2xl font-semibold text-foreground">
                            Set up your subscription
                        </h2>
                        <p className="text-sm text-muted-foreground">
                            14-day free trial • Cancel anytime
                        </p>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                            {/* CARD - WRITABLE INPUTS */}
                            <div className="space-y-4 p-6 rounded-3xl bg-white/80 border border-white/40 shadow-xl">
                                <h3 className="text-xl font-bold mb-4">💳 Card Payment</h3>

                                <StripePayButton onComplete={onPaymentSuccess} plainPassword={password} />
                            </div>


                            {/* STELLAR */}
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
                        <h2 className="text-2xl font-semibold text-foreground">
                            Your vault is ready
                        </h2>
                        <p className="text-sm text-muted-foreground">
                            Your zero-knowledge vault is now active and secured with Stellar verification.
                        </p>
                        {step > 1 && (
                            <button
                                onClick={prevStep}
                                className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
                            >
                                Back
                            </button>
                        )}
                        <button
                            onClick={onComplete}
                            className="mt-2 px-6 py-2 rounded-xl bg-gradient-to-r from-yellow-400 to-yellow-500 text-black font-semibold shadow hover:shadow-lg transition"
                        >
                            Enter my vault
                        </button>
                    </div>
                )}
            </div>
        </div >
    );
};

export default OnboardingWizardBeta;
