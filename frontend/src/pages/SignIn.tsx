import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, Shield } from "lucide-react";

import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { useToast } from "@/hooks/use-toast";
import { useVaultStore } from "@/store/vaultStore";
import { useAuthStore } from "@/store/useAuthStore";
import { useAuth } from "@/hooks/useAuth";
import * as AppAPI from "../../wailsjs/go/main/App";
import { Button } from "@/components/ui/button";
import { StellarLoginForm } from "@/components/StellarLoginForm";
import { handlers } from "../../wailsjs/go/models";
import { useAppStore } from "@/store/appStore";
import { useVault } from "@/hooks/useVault";
import { normalizePreloadedVault } from "@/services/normalizeVault";


const SignIn = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const { setJwtToken, setRefreshToken, setUser, setLoggedIn, updateOnboarding } = useAuthStore();
  const { loginWithStellar } = useAuth();

  // Local UI state
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);


  const vaultStore = useVaultStore.getState();


  const handleSignIn = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const res: handlers.LoginResponse = await AppAPI.SignIn({ email, password });

      // Check if we got a valid response
      if (!res) throw new Error("SignIn failed: empty result");

      await handleSuccessfulAuth(res);

    } catch (err: any) {
      // Extract error message from various error formats (Wails, React, etc.)
      const errorMsg = String(err?.message || err?.toString?.() || err || "").toLowerCase();

      // Check if this is a "user not found" error - try signup silently
      const isUserNotFound =
        errorMsg.includes("user not found") ||
        errorMsg.includes("invalid credentials") ||
        errorMsg.includes("no rows") ||
        errorMsg.includes("record not found");

      if (isUserNotFound) {
        // Silently attempt signup without showing any error
        try {
          await handleSignupAfterSigninFailure(email, password);
          // If signup succeeds, user will be redirected automatically
          return; // Exit early, don't show any error
        } catch (signupErr: any) {
          // Even if signup fails, show a friendly message
          const signupErrorMsg = String(signupErr?.message || signupErr?.toString?.() || signupErr || "");

          toast({
            title: "Account Creation",
            description: signupErrorMsg.includes("duplicate") || signupErrorMsg.includes("already exists")
              ? "This account already exists. Please check your credentials."
              : "Unable to create account. Please try again or contact support.",
            variant: "destructive",
          });
          return;
        }
      }

      // For other errors (network, server, etc.), show the error
      toast({
        title: "Authentication failed",
        description: errorMsg || "An unexpected error occurred. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleSuccessfulAuth = async (data: any) => {
    console.log("🎉 Auth Success (raw):", data);

    // Normalize backend shape into what store expects
    const normalized = normalizePreloadedVault(data);
    console.log("🎯 Normalized payload:", normalized);

    // Save user
    setUser(normalized.User);
    setLoggedIn(true);
    updateOnboarding({ userId: normalized.User.id });
    localStorage.setItem("userId", JSON.stringify(normalized.User.id));

    // Save tokens
    if (normalized.Tokens) {
      setJwtToken(normalized.Tokens.access_token);
      setRefreshToken(normalized.Tokens.refresh_token);
    }

    // Save session data (runtime from backend)
    useAppStore.getState().setSessionData({
      user: normalized.User,
      vault_runtime_context: normalized.vault_runtime_context || null,
      last_cid: normalized.last_cid,
      dirty: normalized.dirty,
    });

    // Load vault into zustand (pass normalized)
    console.log('🚀 SignIn: About to call vaultStore.loadVault with:', {
      hasUser: !!normalized.User,
      hasVault: !!normalized.Vault,
      hasEntries: !!normalized.Vault?.entries,
      User: normalized.User,
      Vault: normalized.Vault,
    });

    try {
      await vaultStore.loadVault({
        User: normalized.User,
        Vault: normalized.Vault,
        SharedEntries: normalized.SharedEntries,
        vault_runtime_context: normalized.vault_runtime_context,
        last_cid: normalized.last_cid,
        dirty: normalized.dirty,
        Tokens: normalized.Tokens,
      });
      console.log('✅ SignIn: vaultStore.loadVault completed successfully');
    } catch (error) {
      console.error('❌ SignIn: vaultStore.loadVault failed:', error);
      throw error;
    }

    toast({
      title: "Welcome!",
      description: "Redirecting to your vault...",
    });

    setTimeout(() => navigate("/dashboard"), 500);
  };

  const handleSignupAfterSigninFailure = async (email: string, password: string) => {
    const userAlias = email.split("@")[0];

    const signupResult = await AppAPI.SignUp({
      user_id: email,
      user_alias: userAlias,
      password,
      vault_name: `${userAlias}-vault`,
      role: "user",
      repo_template: "",
      encryption_policy: "AES-256-GCM",
      federated_providers: [],
    });

    return handleSuccessfulAuth(signupResult);
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate("/")}
          className="mb-6"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Home
        </Button>

        <Card>
          <CardHeader className="text-center space-y-4">
            <div className="mx-auto h-12 w-12 rounded-xl bg-gradient-primary flex items-center justify-center">
              <Shield className="h-7 w-7 text-primary-foreground" />
            </div>
            <div>
              <CardTitle className="text-2xl">Sign In to VaultCore</CardTitle>
              <CardDescription>
                Enter your credentials to access your vault
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSignIn} className="space-y-4">

              {/* Email */}
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="your@email.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>

              {/* Password */}
              <div className="space-y-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>

              <Button
                type="submit"
                className="w-full"
                disabled={isLoading}
              >
                {isLoading ? "Signing in..." : "Sign In"}
              </Button>

              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t border-border" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-background px-2 text-muted-foreground">
                    Or continue with
                  </span>
                </div>
              </div>

              <StellarLoginForm onLogin={loginWithStellar} />


              {/* Offline mode (still useful) */}
              <Button
                type="button"
                variant="outline"
                className="w-full"
                onClick={() => navigate("/vault/offline")}
              >
                Offline Mode
              </Button>
            </form>

            <p className="text-xs text-muted-foreground text-center mt-6">
              Keycloak integration in progress. Onboarding will trigger when upgrading an account.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default SignIn;
