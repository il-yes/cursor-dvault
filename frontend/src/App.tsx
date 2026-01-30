import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { useCallback, useEffect, useState } from "react";
import { useVaultStore } from "@/store/vaultStore";
import Home from "./pages/Home";
import Index from "./pages/Index";
import Vault from "./pages/Vault";
import OfflineVault from "./pages/OfflineVault";
import SignIn from "./pages/SignIn";
import NotFound from "./pages/NotFound";
import ShareEntries from "./pages/ShareEntries";
import Profile from "./pages/Profile";
import Settings from "./pages/Settings";
import EmailLookup from "./pages/EmailLookup";
import LoginStep2 from "./pages/LoginStep2";
import { ThemeProvider } from "@/components/theme-provider";
import Feedback from "./pages/Feedback";
import About from "./pages/About";
import ProfileBeta from "./pages/ProfileBeta";
import SettingsBeta from "./pages/SettingsBeta";
import OnboardingWizard from "./components/OnboardingWizard";
import OnboardingWizardBeta from "@/components/onBoardingWizardBeta";
import { Elements } from '@stripe/react-stripe-js';
import { stripePromise } from '@/lib/stripe';
import PaymentSuccess from "./components/PaymentSuccess";
import SubscriptionManager from "./components/Subscription/subscriptionManager";


const queryClient = new QueryClient();

function AppContent() {
	const loadVault = useVaultStore((state) => state.loadVault);

	useEffect(() => {
		// loadVault();
	}, []);
	// to define
	const [isOnboarded, setIsOnboarded] = useState(false);
	const [walletStatus, setWalletStatus] = useState('disconnected');
	const [ipfsStatus, setIpfsStatus] = useState('idle');
	const [isWailsReady, setIsWailsReady] = useState(false);

	// Safe Wails backend check
	const checkWailsBackend = useCallback(async () => {
		if (typeof window !== 'undefined' && window.go && window.go.Stellar) {
			try {
				const status = await window.go.Stellar.CheckWalletStatus();
				setWalletStatus(status || 'disconnected');
				setIsWailsReady(true);
			} catch (error) {
				console.warn('Wails backend not ready:', error);
				setWalletStatus('disconnected');
			}

			try {
				const ipfsStatus = await window.go.IPFS.CheckNodeStatus();
				setIpfsStatus(ipfsStatus || 'idle');
			} catch (error) {
				console.warn('IPFS backend not ready:', error);
			}
		} else {
			setWalletStatus('mock-connected');
			setIpfsStatus('mock-ready');
		}
	}, []);

	useEffect(() => {
		checkWailsBackend();

		// Poll for Wails readiness (handles hot reload)
		const interval = setInterval(checkWailsBackend, 1000);
		return () => clearInterval(interval);
	}, [checkWailsBackend]);

	if (!isOnboarded) {
		return (
			<Elements stripe={stripePromise}>
				<OnboardingWizardBeta
					onComplete={() => {
					setIsOnboarded(true);
					localStorage.setItem('ankhora-onboarded', 'true');
					}}
				/>
			</Elements>
		);
	}		

	return (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/dashboard" element={<Index />} />
			<Route path="/dashboard/vault" element={<Vault />} />
			<Route path="/dashboard/vault/:filter" element={<Vault />} />
			<Route path="/dashboard/vault/folder/:folderId" element={<Vault />} />
			<Route path="/dashboard/shared" element={<ShareEntries />} />
			<Route path="/dashboard/profile" element={<Profile />} />
			<Route path="/dashboard/profile-beta" element={<ProfileBeta />} />
			<Route path="/dashboard/settings" element={<Settings />} />
			<Route path="/dashboard/settings-beta" element={<SettingsBeta />} />
			<Route path="/vault/offline" element={<OfflineVault />} />
			<Route path="/auth/signin" element={<SignIn />} />
			<Route path="/login/email" element={<EmailLookup />} />
			<Route path="/login/step2" element={<LoginStep2 />} />
			<Route path="/dashboard/feedback" element={<Feedback />} />
			<Route path="/dashboard/about" element={<About />} />
			<Route path="/payment/success" element={<PaymentSuccess />} />
			<Route path="/dashboard/subscription" element={<SubscriptionManager />} />
			{/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
			<Route path="*" element={<NotFound />} />
		</Routes>
	);
}

const App = () => (
	<ThemeProvider defaultTheme="light" storageKey="ankhora-theme">
		<QueryClientProvider client={queryClient}>
			<Toaster />
			<Sonner />
			<BrowserRouter>
				<AppContent />
			</BrowserRouter>
		</QueryClientProvider>
	</ThemeProvider>
);

export default App;


