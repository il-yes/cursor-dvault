import { useState, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { login } from "@/services/api";
import { useVaultStore } from "@/store/vaultStore";
import { Shield, ArrowLeft, Eye, EyeOff } from "lucide-react";
import { handlers } from "wailsjs/go/models";
import * as AppAPI from "../../wailsjs/go/main/App";
import { normalizePreloadedVault } from "@/services/normalizeVault";
import { useAppStore } from "@/store/appStore";
import { useAuthStore } from "@/store/useAuthStore";
import { useVault } from "@/hooks/useVault";

const LoginStep2 = () => {
  const [searchParams] = useSearchParams();
  const email = searchParams.get("email") || "";
  const methods = searchParams.get("methods")?.split(",") || ["password"];

  const [password, setPassword] = useState("");
  const [publicKey, setPublicKey] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();
  const { setJwtToken, setRefreshToken, setUser, setLoggedIn, updateOnboarding } = useAuthStore();

  const vaultStore = useVaultStore.getState();

  const hasPassword = methods.includes("password");
  const hasStellar = methods.includes("stellar");

  useEffect(() => {
    if (!email) {
      navigate("/login/email");
    }
  }, [email, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (hasPassword && !password) {
      toast({
        title: "Password required",
        description: "Please enter your password.",
        variant: "destructive",
      });
      return;
    }

    if (hasStellar && !publicKey) {
      toast({
        title: "Stellar key required",
        description: "Please enter your Stellar public key.",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);

    try {
      const payload: any = { email };

      if (hasPassword && password) {
        payload.password = password;
      }

      if (hasStellar && publicKey) {
        payload.publicKey = publicKey;
        console.log("Stellar login payload:", payload);

        const res: handlers.LoginResponse = await AppAPI.SignIn(payload);

        if (!res) throw new Error("SignIn failed: empty result");

        await handleSuccessfulAuth(res);

        toast({
          title: "Welcome back!",
          description: "Successfully logged in.",
        });

        navigate("/dashboard");
      }

      const response = await login(payload);
      if (!response) throw new Error("SignIn failed: empty result");

      await handleSuccessfulAuth(response);

      toast({
        title: "Welcome back!",
        description: "Successfully logged in.",
      });

      navigate("/dashboard");
    } catch (error) {
      toast({
        title: "Login failed",
        description: error instanceof Error ? error.message : "Please check your credentials and try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleSuccessfulAuth = async (data: any) => {
    console.log("🎉 Auth Success (raw):", data);

    const normalized = normalizePreloadedVault(data);
    console.log("🎯 Normalized payload:", normalized);

    setUser(normalized.User);
    setLoggedIn(true);
    updateOnboarding({ userId: normalized.User.id });
    localStorage.setItem("userId", JSON.stringify(normalized.User.id));

    if (normalized.Tokens) {
      setJwtToken(normalized.Tokens.access_token);
      setRefreshToken(normalized.Tokens.refresh_token);
    }

    useAppStore.getState().setSessionData({
      user: normalized.User,
      vault_runtime_context: normalized.vault_runtime_context || null,
      last_cid: normalized.last_cid,
      dirty: normalized.dirty,
    });

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

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900 p-6">
      <div className="w-full max-w-md space-y-8 animate-fadeIn">
        <Button
          variant="ghost"
          onClick={() => navigate("/login/email")}
          className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </Button>

        <div className="text-center space-y-3">
          <div className="mx-auto flex items-center justify-center w-16 h-16 rounded-2xl bg-white/60 dark:bg-zinc-800/50 shadow-sm backdrop-blur-sm">
            <Shield className="h-8 w-8 text-primary" />
          </div>
          <h1 className="text-3xl font-semibold tracking-tight">Welcome Back</h1>
          <p className="text-sm text-muted-foreground max-w-sm mx-auto">
            Sign in to <span className="font-medium text-foreground">{email}</span>
          </p>
        </div>

        <Card className="border-none shadow-md backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60">
          <CardHeader className="text-center space-y-1">
            <CardTitle className="text-lg font-medium">Enter Credentials</CardTitle>
            <CardDescription>Access your secure vault</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-5">
              {hasPassword && (
                <div className="space-y-2">
                  <Label htmlFor="password" className="text-sm font-medium">
                    Password
                  </Label>
                  <div className="relative">
                    <Input
                      id="password"
                      type={showPassword ? "text" : "password"}
                      className="h-11 rounded-xl border-zinc-200 dark:border-zinc-700 pr-12 focus:ring-2 focus:ring-primary/30 transition-all"
                      placeholder="Enter your password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      required
                      disabled={isLoading}
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg transition-colors"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? (
                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                      ) : (
                        <Eye className="h-4 w-4 text-muted-foreground" />
                      )}
                    </Button>
                  </div>
                </div>
              )}

              {hasStellar && (
                <div className="space-y-3">
                  <Label htmlFor="publicKey" className="text-sm font-medium">
                    Stellar Public Key
                  </Label>
                  <Input
                    id="publicKey"
                    type="text"
                    className="h-11 rounded-xl border-zinc-200 dark:border-zinc-700 focus:ring-2 focus:ring-primary/30 transition-all"
                    placeholder="GXXXXXX..."
                    value={publicKey}
                    onChange={(e) => setPublicKey(e.target.value)}
                    required
                    disabled={isLoading}
                  />
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    You'll be prompted to sign a message with your wallet
                  </p>
                </div>
              )}

              <Button
                type="submit"
                className="w-full h-11 rounded-xl text-[15px] font-medium transition-all hover:scale-[1.01] active:scale-[0.99]"
                disabled={isLoading}
              >
                {isLoading ? "Signing in..." : "Sign In"}
              </Button>

              <Button
                type="button"
                variant="link"
                className="w-full text-xs text-muted-foreground hover:text-foreground h-9 rounded-xl border border-zinc-200 dark:border-zinc-700 transition-all hover:bg-zinc-50 dark:hover:bg-zinc-800"
                onClick={() => {
                  toast({
                    title: "Coming soon",
                    description: "Password reset functionality will be available soon.",
                  });
                }}
              >
                Forgot your password?
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default LoginStep2;
