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
        // TODO: Add signing logic when backend is ready
        // payload.signedMessage = signedMessage;
        // payload.signature = signature;
      }

      const response = await login(payload);
      // handleSuccessfulAuth(response);

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

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-md space-y-6">
        <Button
          variant="ghost"
          onClick={() => navigate("/login/email")}
          className="mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>

        <div className="text-center space-y-2">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10 mb-4">
            <Shield className="h-8 w-8 text-primary" />
          </div>
          <h1 className="text-3xl font-bold">Welcome Back</h1>
          <p className="text-muted-foreground">
            Sign in to {email}
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Enter Your Credentials</CardTitle>
            <CardDescription>
              Use your authentication method to sign in
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {hasPassword && (
                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <div className="relative">
                    <Input
                      id="password"
                      type={showPassword ? "text" : "password"}
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
                      className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
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
                <div className="space-y-2">
                  <Label htmlFor="publicKey">Stellar Public Key</Label>
                  <Input
                    id="publicKey"
                    type="text"
                    placeholder="GXXXXXX..."
                    value={publicKey}
                    onChange={(e) => setPublicKey(e.target.value)}
                    required
                    disabled={isLoading}
                  />
                  <p className="text-xs text-muted-foreground">
                    You'll be prompted to sign a message with your wallet
                  </p>
                </div>
              )}

              <Button
                type="submit"
                className="w-full"
                disabled={isLoading}
              >
                {isLoading ? "Signing in..." : "Sign In"}
              </Button>

              <Button
                type="button"
                variant="link"
                className="w-full text-sm"
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
