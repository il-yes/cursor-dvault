import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";
import { useCallback, useEffect, useState } from "react";
import Index from "./pages/Index";
import Vault from "./pages/Vault";
import OfflineVault from "./pages/OfflineVault";
import SignIn from "./pages/SignIn";
import NotFound from "./pages/NotFound";
import ShareEntries from "./pages/ShareEntries";
import EmailLookup from "./pages/EmailLookup";
import LoginStep2 from "./pages/LoginStep2";
import { ThemeProvider } from "@/components/theme-provider";
import Feedback from "./pages/Feedback";
import About from "./pages/About";
import ProfileBeta from "./pages/ProfileBeta";
import Settings from "./pages/Settings";
import OnboardingWizardBeta from "@/components/onBoardingWizardBeta";
import { Elements } from '@stripe/react-stripe-js';
import { stripePromise } from '@/lib/stripe';
import PaymentSuccess from "./components/PaymentSuccess";
import SubscriptionManager from "./components/Subscription/subscriptionManager";
import * as AppAPI from "../wailsjs/go/main/App";
import NotificationsPage from "./pages/NotificationsPage";
import C3App from "./pages/C3.ledger";
import { Inbox } from "./components/C3/inbox";
import { AssetCard } from "./components/C3/asset-view";
import LedgerPage from "./components/C3/ui/ledger/LedgerPage";
import InboxPage from "./components/C3/ui/inbox/InboxPage";
import AssetPage from "./components/C3/ui/thread/AssetPage";
import ChannelPage from "./components/C3/ui/channel/ChannelPage";
import * as ROUTES from './constants/routes';


const queryClient = new QueryClient();

function AppContent() {
	const [isLoading, setIsLoading] = useState(true);
	// to define
	const [isOnboarded, setIsOnboarded] = useState(false);
	const [walletStatus, setWalletStatus] = useState('disconnected');
	const [ipfsStatus, setIpfsStatus] = useState('idle');
	const [isWailsReady, setIsWailsReady] = useState(false);
	const navigate = useNavigate();
	const [appState, setAppState] = useState<any>(false);


	const handleOnBoardingComplete = async () => {
		await AppAPI.CompleteOnboarding();
		const fresh = await AppAPI.GetAppState();
		setAppState(fresh);
		navigate("/");
	}


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

	const init = useCallback(async () => {
		if (typeof window !== 'undefined' && window.go) {
			try {
				// Wait until Wails backend is ready
				if (!window.go) {
					setTimeout(init, 200);
					return;
				}

				const appState = await AppAPI.GetAppState();
				setAppState(appState);
				console.log("App state:", appState);

			} catch (e) {
				console.error("Failed to get app state:", e);
				// navigate("/on-boarding", { replace: true });
			} finally {
				setIsLoading(false);
			}
		} else {
			setIsLoading(false);
		}
	}, []);

	useEffect(() => {
		checkWailsBackend();

		// Poll for Wails readiness (handles hot reload)
		const interval = setInterval(checkWailsBackend, 1000);
		return () => clearInterval(interval);
	}, [checkWailsBackend]);

	useEffect(() => {
		init();
	}, [init]);

	if (!appState.has_vault) {
		return (
			<Elements stripe={stripePromise}>
				<OnboardingWizardBeta
					onComplete={() => {
						setIsOnboarded(true);
						localStorage.setItem('ankhora-onboarded', 'true');
						handleOnBoardingComplete();
					}}
				/>
			</Elements>
		);
	}


	return (
		<Routes>


			<Route path="/" element={<EmailLookup />} />
			<Route path={ROUTES.DASHBOARD} element={<Index />} />
			<Route path={ROUTES.VAULT} element={<Vault />} />
			<Route path={ROUTES.VAULT_FILTER} element={<Vault />} />
			<Route path={ROUTES.VAULT_FOLDER} element={<Vault />} />
			<Route path={ROUTES.SHARED} element={<ShareEntries />} />
			<Route path={ROUTES.PROFILE} element={<ProfileBeta />} />
			<Route path={ROUTES.SETTINGS} element={<Settings />} />
			{/* <Route path="/vault/offline" element={<OfflineVault />} /> */}
			{/* <Route path="/auth/signin" element={<SignIn />} /> */}
			<Route path={ROUTES.LOGIN_EMAIL} element={<EmailLookup />} />
			<Route path={ROUTES.LOGIN_STEP2} element={<LoginStep2 />} />
			<Route path={ROUTES.NOTIFICATIONS} element={<NotificationsPage />} />
			<Route path={ROUTES.FEEDBACK} element={<Feedback />} />
			<Route path={ROUTES.ABOUT} element={<About />} />
			<Route path={ROUTES.PAYMENT_SUCCESS} element={<PaymentSuccess />} />
			<Route path={ROUTES.SUBSCRIPTION} element={<SubscriptionManager />} />
			<Route
				path={ROUTES.ON_BOARDING}
				element={
					<Elements stripe={stripePromise}>
						<OnboardingWizardBeta onComplete={handleOnBoardingComplete} />
					</Elements>
				}
			/>
			<Route path={ROUTES.LEDGER} element={<LedgerPage />} />
			<Route path={ROUTES.INBOX} element={<InboxPage />} />
			<Route path={ROUTES.THREAD} element={<AssetPage />} />
			<Route path={ROUTES.CHANNEL} element={<ChannelPage />} />


			{/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
			<Route path={ROUTES.NOT_FOUND} element={<NotFound />} />
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


