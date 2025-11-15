import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, Shield } from "lucide-react";

import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
<<<<<<< Updated upstream
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/use-toast";

import { useAuth } from "@/hooks/useAuth";         // <-- REAL login logic
import { StellarLoginForm } from "@/components/StellarLoginForm";
import { LoginRequest } from "@/types/vault";
=======
import { useToast } from "@/hooks/use-toast";
import { useVaultStore } from "@/store/vaultStore";
import { useAuthStore } from "@/store/useAuthStore";
import { useVault } from "@/hooks/useVault";
import * as AppAPI from "../../wailsjs/go/main/App";
>>>>>>> Stashed changes

const SignIn = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const { setVault } = useVaultStore();
  const { hydrateVault } = useVault();
  const { setJwtToken, setRefreshToken, setUser, setLoggedIn, updateOnboarding } = useAuthStore();

  // Local UI state
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);

<<<<<<< Updated upstream
  // v0 REAL logic
  const { loginWithPassword, loginWithStellar } = useAuth();

=======
>>>>>>> Stashed changes
  const handleSignIn = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
<<<<<<< Updated upstream
      // 🔐 REAL BACKEND CALL
      // Note: loginWithPassword handles user update and navigation internally
      await loginWithPassword({ email, password });

      toast({
        title: "Authentication successful",
        description: "Redirecting to your vault...",
      });
    } catch (err: any) {
      console.error("Login failed:", err);

      toast({
        title: "Authentication failed",
        description: err?.message ?? "Invalid credentials or connection issue.",
=======
      // Call backend SignIn
      const result = await AppAPI.SignIn({ email, password });

      console.log("🔐 Login Response:", result);
      console.log("🔑 Access Token:", result.Tokens?.access_token);

      // If user exists, load vault directly into vaultStore
      if (result && result.Vault) {
        // Save authentication tokens to auth store
        if (result.Tokens) {
          setJwtToken(result.Tokens.access_token);
          setRefreshToken(result.Tokens.refresh_token);
          console.log("✅ JWT Token saved to auth store");
        }

        // Save user info
        setUser(result.User);
        setLoggedIn(true);
        updateOnboarding({ userId: result.User.id });
        localStorage.setItem("userId", JSON.stringify(result.User.id));

        // Transform LoginResponse to VaultContext format
        const vaultContext = {
          user_id: result.User.id.toString(),
          role: result.User.role,
          Vault: result.Vault as any, // Type mismatch between Wails models and VaultContext
          Dirty: false,
          LastUpdated: new Date().toISOString(),
          vault_runtime_context: {
            CurrentUser: {
              id: result.User.id.toString(),
              role: result.User.role,
              name: result.User.username,
              last_name: "",
              email: result.User.email,
              stellar_account: {
                public_key: "",
              },
            },
            AppSettings: {} as any,
            WorkingBranch: "main",
            LoadedEntries: [],
          },
        };
        console.log("📦 Vault Context:", vaultContext);

        // ✅ Save to vaultStore (for persistence)
        setVault(vaultContext as any);

        // ✅ Hydrate VaultContextProvider (for live CRUD operations)
        hydrateVault(vaultContext as any);

        toast({
          title: "Authentication successful",
          description: "Vault loaded! Redirecting to dashboard...",
        });
        setTimeout(() => navigate("/dashboard"), 500);
      } else {
        // Fallback: unknown user → optional onboarding
        toast({
          title: "User not found",
          description: "Please sign up or upgrade your account",
          variant: "destructive",
        });
        // Navigate to onboarding / upgrade page
        setTimeout(() => navigate("/onboarding"), 500);
      }
    } catch (err: any) {
      console.error("SignIn failed:", err);
      toast({
        title: "Authentication failed",
        description: err.message || "Please enter valid credentials.",
>>>>>>> Stashed changes
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
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

<<<<<<< Updated upstream
              {/* Submit */}
=======
>>>>>>> Stashed changes
              <Button
                type="submit"
                className="w-full"
                disabled={isLoading}
              >
                {isLoading ? "Signing in..." : "Sign In"}
              </Button>

<<<<<<< Updated upstream
              {/* Divider */}
=======
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream

			<StellarLoginForm onLogin={loginWithStellar} />


              {/* Offline mode (still useful) */}
=======
>>>>>>> Stashed changes
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
<<<<<<< Updated upstream
              Backend integration active — real login required.
=======
              Keycloak integration in progress. Onboarding will trigger when upgrading an account.
>>>>>>> Stashed changes
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default SignIn;
