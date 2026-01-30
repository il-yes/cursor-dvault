// SubscriptionManager.tsx - React version with Wails + Tailwind glassmorphism

import React, { useCallback, useEffect, useState } from 'react';
// Adjust import path to your generated Wails bindings
import {
    GetSubscriptionDetails,
    GetStorageUsage,
    UpgradeSubscription,
    CancelSubscription,
    GetBillingHistory,
} from '../../services/api';
import { cn } from '@/lib/utils'; // Adjust import path
import { Button } from '@/components/ui/button'; // shadcn/ui or your button component
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { DashboardLayout } from "@/components/DashboardLayout";
import PaymentRequestModal from '../Payments/PaymentRequestModal';

interface Subscription {
    tier: string;
    price: number;
    status: string;
    payment_method: string;
    next_billing_date: string;
    trial_ends_at?: string;
    stellar_tx_hash?: string;
    features: Record<string, any>;
}

interface StorageUsage {
    used_gb: number;
    quota_gb: number;
    percentage: number;
}

interface BillingHistoryItem {
    id: string;
    created_at: string;
    description: string;
    amount: number;
    status: string;
    stellar_tx_hash?: string;
    stripe_intent_id?: string;
}

const SubscriptionManager: React.FC = () => {
    const [subscription, setSubscription] = useState<Subscription | null>(null);
    const [storageUsage, setStorageUsage] = useState<StorageUsage | null>(null);
    const [billingHistory, setBillingHistory] = useState<BillingHistoryItem[]>([]);
    const [loading, setLoading] = useState(false);
    const [showUpgradeModal, setShowUpgradeModal] = useState(false);
    const [showCancelModal, setShowCancelModal] = useState(false);
    const [cancelReason, setCancelReason] = useState('');
    const userId = localStorage.getItem('user_id');
    const [useMocks, setUseMocks] = useState(true);
    const [isPaymentRequestOpen, setIsPaymentRequestOpen] = useState(false);  


    // MOCK VERSION - Replace loadSubscriptionData with this:
    const loadSubscriptionData = useCallback(async () => {
        setLoading(true);

        if (useMocks) {

            // Simulate API delay
            await new Promise(resolve => setTimeout(resolve, 1500));

            // RICH MOCK DATA - Shows ALL features
            setSubscription({
                tier: 'pro_plus', // Change to 'free', 'pro', 'business' to test different states
                price: 25.00,
                status: 'active',
                payment_method: 'card_encrypted',
                next_billing_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(), // 30 days
                trial_ends_at: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(), // 7-day trial
                stellar_tx_hash: 'abc123def456ghi789jkl012mno345pqr678stu901vwx234yz',
                features: {
                    storage_gb: 200,
                    cloud_backup: true,
                    mobile_apps: true,
                    encrypted_payments: true,
                    zero_telemetry: true,
                    anonymous_account: true,
                    version_history: 365, // days
                    team_collaboration: 5, // members
                    compliance_reports: true,
                },
            });

            setStorageUsage({
                used_gb: 145.3,
                quota_gb: 200,
                percentage: 72.65, // Triggers "low storage" warning (>80%)
            });

            setBillingHistory([
                {
                    id: 'inv_001',
                    created_at: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
                    description: 'Pro Plus Monthly',
                    amount: 25.00,
                    status: 'succeeded',
                    stellar_tx_hash: 'tx_123abc456def789ghi012jkl',
                },
                {
                    id: 'inv_002',
                    created_at: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000).toISOString(),
                    description: 'Pro Monthly',
                    amount: 15.00,
                    status: 'succeeded',
                    stripe_intent_id: 'pi_abc123',
                },
                {
                    id: 'inv_003',
                    created_at: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000).toISOString(),
                    description: 'Trial Activation',
                    amount: 0.00,
                    status: 'succeeded',
                    stellar_tx_hash: 'tx_def456',
                },
                {
                    id: 'inv_004',
                    created_at: new Date(Date.now() - 120 * 24 * 60 * 60 * 1000).toISOString(),
                    description: 'Pro Plus Upgrade',
                    amount: 10.00, // Upgrade fee
                    status: 'succeeded',
                },
            ]);

            setLoading(false);
        } else {
            setLoading(true);
            try {
                const [sub, usage, historyResponse] = await Promise.all([
                    GetSubscriptionDetails(),
                    GetStorageUsage(),
                    GetBillingHistory(),
                ]);
                setSubscription(sub as Subscription);
                setStorageUsage(usage as StorageUsage);
                setBillingHistory(historyResponse.history as BillingHistoryItem[]);
            } catch (err) {
                console.error('Failed to load subscription data:', err);
            } finally {
                setLoading(false);
            }
        }
    }, [useMocks]);

    // Load data on mount
    useEffect(() => {
         loadSubscriptionData();
    }, []);

    const getTierColor = useCallback((tier: string) => {
        const colors: Record<string, string> = {
            free: '#999',
            pro: '#4A90E2',
            pro_plus: '#D4AF37',
            business: '#FF8C00',
        };
        return colors[tier] || '#999';
    }, []);

    const getTierName = useCallback((tier: string) => {
        const names: Record<string, string> = {
            free: 'Free',
            pro: 'Pro',
            pro_plus: 'Pro Plus',
            business: 'Business',
        };
        return names[tier] || tier;
    }, []);

    const handleUpgrade = useCallback(async (newTier: string, paymentMethod: string) => {
        try {
            await UpgradeSubscription({
                user_id: userId,
                tier: newTier,
                payment_method: paymentMethod,
            });
            setShowUpgradeModal(false);
            await loadSubscriptionData();
            // Replace with your toast system (Sonner, shadcn toast, etc.)
            console.log('✅ Subscription upgraded successfully!');
        } catch (err: any) {
            console.error('❌ Upgrade failed:', err.message || 'Upgrade failed');
        }
    }, [loadSubscriptionData]);

    const handleCancel = useCallback(async () => {
        if (!cancelReason.trim()) {
            console.error('❌ Please provide a reason for cancellation');
            return;
        }

        try {
            await CancelSubscription({
                user_id: userId,
                reason: cancelReason,
            });
            setShowCancelModal(false);
            setCancelReason('');
            await loadSubscriptionData();
            console.log('✅ Subscription cancelled. You will be downgraded to Free at the end of your billing period.');
        } catch (err: any) {
            console.error('❌ Cancellation failed:', err.message || 'Cancellation failed');
        }
    }, [cancelReason]);

    if (loading) {
        return (
            <DashboardLayout>
                <div className="flex items-center justify-center py-12 backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-yellow-400 mr-2" />
                    <span className="text-muted-foreground">Loading subscription details...</span>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout>
            <div className="max-w-5xl mx-auto p-8 space-y-8">
                {/* Hero Header */}
                <div className="text-center backdrop-blur-xl bg-white/40 dark:bg-zinc-900/40 rounded-3xl p-12 border border-white/30 dark:border-zinc-700/30 shadow-2xl">
                    <h1 className="text-6xl font-black bg-gradient-to-r from-foreground via-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-2xl mb-4">
                        Account    
                    </h1>
                    <p className="text-2xl text-muted-foreground/90 max-w-2xl mx-auto leading-relaxed backdrop-blur-sm">
                        Manage your subscription and vault preferences
                    </p>
                </div>
                {/* <Button onClick={loadSubscriptionData} variant="outline">
                    {useMocks ? '🔌 Use Real Data' : '🧪 Use Mocks'}
                </Button> */}

                <div className="max-w-5xl mx-auto p-8 space-y-8">

                    {/* Current Plan Card */}
                    {subscription && (
                        <div
                            className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-8 border border-white/40 dark:border-zinc-700/40 shadow-xl hover:shadow-2xl transition-all duration-500"
                            style={{ borderColor: getTierColor(subscription.tier) }}
                        >
                            <div className="flex justify-between items-start mb-6">
                                <div>
                                    <h2 className="text-3xl font-bold text-foreground mb-2">
                                        {getTierName(subscription.tier)}
                                    </h2>
                                    <div className="text-2xl font-bold text-yellow-600">
                                        ${subscription.price.toFixed(2)}<span className="text-lg font-normal text-muted-foreground ml-1">/month</span>
                                    </div>
                                </div>
                                <div
                                    className="px-4 py-2 rounded-full text-sm font-semibold uppercase tracking-wide"
                                    style={{ backgroundColor: getTierColor(subscription.tier) }}
                                >
                                    {subscription.status.toUpperCase()}
                                </div>
                            </div>

                            {subscription.trial_ends_at && new Date(subscription.trial_ends_at) > new Date() && (
                                <div className="bg-gradient-to-r from-emerald-500/20 to-emerald-600/20 border border-emerald-500/30 rounded-2xl p-4 mb-6">
                                    <span className="inline-block w-6 h-6 bg-emerald-500 rounded-full flex items-center justify-center text-white text-sm mr-3">🎉</span>
                                    <span className="font-semibold">Trial ends on {new Date(subscription.trial_ends_at).toLocaleDateString()}</span>
                                </div>
                            )}

                            <div className="grid md:grid-cols-2 gap-6 mb-8">
                                <div className="space-y-3">
                                    <span className="text-sm font-medium text-muted-foreground">Payment Method:</span>
                                    <span className="font-semibold">{subscription.payment_method.replace('_', ' ')}</span>
                                </div>
                                <div className="space-y-3">
                                    <span className="text-sm font-medium text-muted-foreground">Next Billing:</span>
                                    <span className="font-semibold">{new Date(subscription.next_billing_date).toLocaleDateString()}</span>
                                </div>
                                {subscription.stellar_tx_hash && (
                                    <div className="space-y-3 col-span-2">
                                        <span className="text-sm font-medium text-muted-foreground">Blockchain Verified:</span>
                                        <a
                                            href={`https://stellar.expert/explorer/public/tx/${subscription.stellar_tx_hash}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="inline-flex items-center text-yellow-600 hover:text-yellow-700 font-semibold underline"
                                        >
                                            View on Stellar →
                                        </a>
                                    </div>
                                )}
                            </div>

                            <div className="flex flex-wrap gap-3">
                                {subscription.tier !== 'business' && subscription.tier !== 'pro_plus' && (
                                    <Button
                                        onClick={() => setShowUpgradeModal(true)}
                                        className="bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black font-semibold shadow-lg hover:shadow-xl transition-all"
                                    >
                                        Upgrade Plan
                                    </Button>
                                )}
                                {subscription.tier !== 'free' && (
                                    <Button
                                        variant="outline"
                                        onClick={() => setShowCancelModal(true)}
                                        className="border-red-300 hover:border-red-400 text-red-700 hover:bg-red-50 transition-all"
                                    >
                                        Cancel Subscription
                                    </Button>
                                )}
                            </div>
                        </div>
                    )}

                    {/* Storage Usage */}
                    {storageUsage && (
                        <div className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-8 border border-white/40 dark:border-zinc-700/40 shadow-xl">
                            <h3 className="text-2xl font-bold mb-6 flex items-center gap-2">
                                Storage Usage
                            </h3>
                            <div className="space-y-4">
                                <div className="storage-bar-container">
                                    <div className="relative h-4 bg-white/30 dark:bg-zinc-800/50 rounded-full overflow-hidden">
                                        <div
                                            className={cn(
                                                "h-full rounded-full transition-all duration-1000 shadow-inner",
                                                storageUsage.percentage > 90 ? 'bg-gradient-to-r from-red-500 to-red-600' : 'bg-gradient-to-r from-yellow-400 to-yellow-500'
                                            )}
                                            style={{ width: `${storageUsage.percentage}%` }}
                                        />
                                    </div>
                                </div>
                                <div className="text-center font-mono text-lg font-bold text-foreground">
                                    {storageUsage.used_gb.toFixed(2)} GB / {storageUsage.quota_gb} GB
                                    <span className="text-sm font-normal text-muted-foreground ml-2">
                                        ({storageUsage.percentage.toFixed(1)}%)
                                    </span>
                                </div>
                                {storageUsage.percentage > 80 && (
                                    <div className="bg-gradient-to-r from-orange-500/20 to-red-500/20 border border-orange-500/30 rounded-2xl p-4 flex items-center gap-3">
                                        <span className="text-2xl">⚠️</span>
                                        <span className="font-semibold">You're running low on storage. Consider upgrading your plan.</span>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {/* Features */}
                    {subscription && (
                        <div className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-8 border border-white/40 dark:border-zinc-700/40 shadow-xl">
                            <h3 className="text-2xl font-bold mb-6">Your Plan Includes:</h3>
                            <div className="grid md:grid-cols-2 gap-4">
                                {Object.entries(subscription.features).map(([key, value]) => {
                                    if (value === true || (typeof value === 'number' && value > 0) || (typeof value === 'string' && value !== '')) {
                                        const label = key.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase());
                                        return (
                                            <div key={key} className="flex items-center gap-3 p-3 bg-emerald-50/50 dark:bg-emerald-900/20 rounded-xl border border-emerald-200/50">
                                                <span className="text-emerald-600 text-lg">✅</span>
                                                <span className="font-medium">{label}</span>
                                                {(typeof value === 'number' || typeof value === 'string') && (
                                                    <span className="text-sm text-muted-foreground ml-auto">{value}</span>
                                                )}
                                            </div>
                                        );
                                    }
                                    return null;
                                })}
                            </div>
                        </div>
                    )}

                    {/* Billing History */}
                    {billingHistory.length > 0 && (
                        <div className="backdrop-blur-xl bg-white/60 dark:bg-zinc-900/60 rounded-3xl p-8 border border-white/40 dark:border-zinc-700/40 shadow-xl overflow-x-auto">
                            <h3 className="text-2xl font-bold mb-6">Billing History</h3>
                            <div className="overflow-x-auto">
                                <table className="w-full text-sm">
                                    <thead>
                                        <tr className="border-b border-white/20">
                                            <th className="text-left pb-4 font-semibold text-muted-foreground">Date</th>
                                            <th className="text-left pb-4 font-semibold text-muted-foreground">Description</th>
                                            <th className="text-left pb-4 font-semibold text-muted-foreground">Amount</th>
                                            <th className="text-left pb-4 font-semibold text-muted-foreground">Status</th>
                                            <th className="text-left pb-4 font-semibold text-muted-foreground">Receipt</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {billingHistory.map((payment) => (
                                            <tr key={payment.id} className="border-b border-white/10 hover:bg-white/20 transition">
                                                <td className="py-4 font-medium">{new Date(payment.created_at).toLocaleDateString()}</td>
                                                <td className="py-4">{payment.description}</td>
                                                <td className="py-4 font-semibold text-yellow-600">${payment.amount.toFixed(2)}</td>
                                                <td>
                                                    <span className={cn(
                                                        'px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wide',
                                                        payment.status === 'succeeded' ? 'bg-emerald-500/20 text-emerald-600 border border-emerald-500/30' :
                                                            payment.status === 'failed' ? 'bg-red-500/20 text-red-600 border border-red-500/30' :
                                                                'bg-yellow-500/20 text-yellow-600 border border-yellow-500/30'
                                                    )}>
                                                        {payment.status}
                                                    </span>
                                                </td>
                                                <td className="py-4">
                                                    {payment.stellar_tx_hash ? (
                                                        <a
                                                            href={`https://stellar.expert/explorer/public/tx/${payment.stellar_tx_hash}`}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                            className="text-yellow-600 hover:text-yellow-700 font-semibold text-sm underline"
                                                        >
                                                            View Receipt
                                                        </a>
                                                    ) : payment.stripe_intent_id ? (
                                                        <Button variant="ghost" size="sm" className="h-8 px-3">
                                                            Download
                                                        </Button>
                                                    ) : (
                                                        <span className="text-muted-foreground text-sm">—</span>
                                                    )}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}

                    {/* Upgrade Modal */}
                    {showUpgradeModal && (
                        <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-6">
                            <div className="backdrop-blur-xl bg-white/90 dark:bg-zinc-900/90 rounded-3xl p-8 w-full max-w-2xl max-h-[90vh] overflow-y-auto border border-white/40 shadow-2xl">
                                <div className="flex justify-between items-center mb-6">
                                    <h2 className="text-2xl font-bold text-foreground">Upgrade Your Plan</h2>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-8 w-8 p-0"
                                        onClick={() => setShowUpgradeModal(false)}
                                    >
                                        ×
                                    </Button>
                                </div>
                                <div className="grid md:grid-cols-2 gap-6">
                                    {subscription?.tier === 'free' && (
                                        <>
                                            <div className="upgrade-card border border-blue-300 hover:border-blue-400 p-6 rounded-2xl hover:shadow-xl transition-all cursor-pointer"
                                                onClick={() => handleUpgrade('pro', 'standard')}>
                                                <h3 className="text-xl font-bold mb-2">Pro</h3>
                                                <div className="text-2xl font-bold text-blue-600 mb-4">$15/month</div>
                                                <ul className="space-y-2 mb-6 text-sm text-muted-foreground">
                                                    <li>100GB storage</li>
                                                    <li>Cloud backup</li>
                                                    <li>Mobile apps</li>
                                                </ul>
                                                <Button className="w-full bg-blue-500 hover:bg-blue-600">Upgrade to Pro</Button>
                                            </div>
                                            <div className="upgrade-card featured border-2 border-yellow-400 bg-gradient-to-br from-yellow-50 to-yellow-100 p-6 rounded-2xl hover:shadow-2xl transition-all cursor-pointer relative"
                                                onClick={() => handleUpgrade('pro_plus', 'encrypted')}>
                                                <div className="absolute -top-3 left-1/2 transform -translate-x-1/2 bg-white px-4 py-1 rounded-full border border-yellow-400">
                                                    <span className="text-yellow-600 font-semibold">⭐ Featured</span>
                                                </div>
                                                <h3 className="text-xl font-bold mb-2">Pro Plus</h3>
                                                <div className="text-2xl font-bold text-yellow-600 mb-4">$25/month</div>
                                                <ul className="space-y-2 mb-6 text-sm text-muted-foreground">
                                                    <li>200GB storage</li>
                                                    <li>Encrypted payments</li>
                                                    <li>Zero telemetry</li>
                                                    <li>Anonymous account</li>
                                                </ul>
                                                <Button className="w-full bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black shadow-lg">Upgrade to Pro Plus</Button>
                                            </div>
                                        </>
                                    )}
                                    {subscription?.tier === 'pro' && (
                                        <div className="upgrade-card featured border-2 border-yellow-400 bg-gradient-to-br from-yellow-50 to-yellow-100 p-6 rounded-2xl hover:shadow-2xl transition-all cursor-pointer col-span-full md:col-span-1 mx-auto max-w-md"
                                            onClick={() => handleUpgrade('pro_plus', 'encrypted')}>
                                            <div className="absolute -top-3 left-1/2 transform -translate-x-1/2 bg-white px-4 py-1 rounded-full border border-yellow-400">
                                                <span className="text-yellow-600 font-semibold">⭐ Pro Plus</span>
                                            </div>
                                            <h3 className="text-xl font-bold mb-2">Pro Plus</h3>
                                            <div className="text-2xl font-bold text-yellow-600 mb-4">$25/month (+$10)</div>
                                            <ul className="space-y-2 mb-6 text-sm text-muted-foreground">
                                                <li>2x storage (200GB)</li>
                                                <li>Encrypted payments</li>
                                                <li>Zero telemetry</li>
                                                <li>Anonymous account</li>
                                                <li>Version history</li>
                                            </ul>
                                            <Button className="w-full bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black shadow-lg">Upgrade to Pro Plus</Button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Cancel Modal */}
                    {showCancelModal && (
                        <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-6">
                            <div className="backdrop-blur-xl bg-white/90 dark:bg-zinc-900/90 rounded-3xl p-8 w-full max-w-lg border border-white/40 shadow-2xl">
                                <div className="flex justify-between items-center mb-6">
                                    <h2 className="text-2xl font-bold text-foreground">Cancel Subscription</h2>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-8 w-8 p-0"
                                        onClick={() => setShowCancelModal(false)}
                                    >
                                        ×
                                    </Button>
                                </div>
                                <div className="space-y-6">
                                    <div className="bg-gradient-to-r from-red-500/20 to-red-600/20 border border-red-500/30 rounded-2xl p-6 flex items-start gap-4">
                                        <span className="text-3xl mt-1">⚠️</span>
                                        <div>
                                            <p className="font-semibold text-foreground">
                                                Are you sure you want to cancel your {getTierName(subscription?.tier || '')} subscription?
                                            </p>
                                        </div>
                                    </div>

                                    <div className="bg-gradient-to-r from-yellow-500/10 to-yellow-600/10 border border-yellow-500/20 rounded-2xl p-6">
                                        <h4 className="font-semibold mb-4 text-foreground">What happens when you cancel:</h4>
                                        <ul className="space-y-2 text-sm text-muted-foreground">
                                            <li>• You'll keep access until {new Date(subscription?.next_billing_date || '').toLocaleDateString()}</li>
                                            <li>• After that, you'll be downgraded to Free tier (5GB storage)</li>
                                            <li>• Your data will remain encrypted and accessible</li>
                                            <li>• You can reactivate anytime</li>
                                        </ul>
                                    </div>

                                    <div>
                                        <label className="block text-sm font-medium mb-2 text-muted-foreground">
                                            Why are you cancelling? (optional)
                                        </label>
                                        <Textarea
                                            value={cancelReason}
                                            onChange={(e) => setCancelReason(e.target.value)}
                                            placeholder="Too expensive, missing features, switching to competitor, etc."
                                            rows={4}
                                            className="resize-none"
                                        />
                                    </div>
                                </div>

                                <div className="flex gap-3 pt-6 border-t border-white/20">
                                    <Button
                                        variant="outline"
                                        onClick={() => setShowCancelModal(false)}
                                        className="flex-1"
                                    >
                                        Keep Subscription
                                    </Button>
                                    <Button
                                        onClick={handleCancel}
                                        className="flex-1 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white shadow-lg hover:shadow-xl"
                                    >
                                        Confirm Cancellation
                                    </Button>
                                </div>
                            </div>
                        </div>
                    )}
                </div>

            </div>
        </DashboardLayout>
    );
};

export default SubscriptionManager;
