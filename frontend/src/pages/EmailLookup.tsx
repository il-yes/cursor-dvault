import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { checkEmail } from "@/services/api";
import { Shield, ArrowLeft } from "lucide-react";
import { StellarLoginForm } from "@/components/StellarLoginForm";
import { useAuth } from "@/hooks/useAuth";

const EmailLookup = () => {
  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();
  const { loginWithStellar } = useAuth();

  const validateEmail = (email: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateEmail(email)) {
      toast({
        title: "Invalid email",
        description: "Please enter a valid email address.",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);
    try {
      const response = await checkEmail(email);
      if (response.status === "NEW_USER") {
        navigate(`/signup?email=${encodeURIComponent(email)}`);
      } else if (response.status === "EXISTS") {
        navigate(`/login/step2?email=${encodeURIComponent(email)}&methods=${response.auth_methods?.join(",")}`);
      }
    } catch (error) {
      console.error(error);
      toast({
        title: "Error",
        description: "Failed to check email. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900 p-6">
      <div className="w-full max-w-md space-y-8 animate-fadeIn">
        <Button
          variant="ghost"
          onClick={() => navigate("/")}
          className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </Button>

        <div className="text-center space-y-3">
          <div className="mx-auto flex items-center justify-center w-16 h-16 rounded-2xl bg-white/60 dark:bg-zinc-800/50 shadow-sm backdrop-blur-sm">
            <Shield className="h-8 w-8 text-primary" />
          </div>
          {/* <h1 className="text-3xl font-semibold tracking-tight">Welcome to Ankhora</h1> */}
          <p className="text-sm text-muted-foreground">Your cryptographic vault</p>
        </div>

        <Card className="border-none shadow-md backdrop-blur-sm bg-white/70 dark:bg-zinc-900/60">
          <CardHeader className="text-center">
            <CardTitle className="text-lg font-medium">Unlock with your email</CardTitle>
            <CardDescription>Access your secure vault</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-5" style={{padding: "0px 20px"}}>
              <div className="space-y-2">
                <Label htmlFor="email" className="text-sm font-medium">
                  Email Address
                </Label>
                <Input
                  id="email"
                  type="email"
                  className="h-11 rounded-xl border-zinc-200 dark:border-zinc-700 focus:ring-2 focus:ring-primary/30 transition-all"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  disabled={isLoading}
                />
              </div>

              <Button
                type="submit"
                className="w-full h-11 rounded-xl text-[15px] font-medium transition-all hover:scale-[1.01] active:scale-[0.99]"
                disabled={isLoading}
              >
                {isLoading ? "Checking..." : "Continue"}
              </Button>
            </form>

            <div className="relative mt-7">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-zinc-200 dark:border-zinc-700" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-white dark:bg-zinc-900 px-2 text-muted-foreground">
                  or continue with
                </span>
              </div>
            </div>

            <div className="mt-6">
              <StellarLoginForm onLogin={loginWithStellar} />
            </div>
          </CardContent>
        </Card>

        <p className="text-center text-xs text-muted-foreground leading-relaxed max-w-sm mx-auto">
          By continuing, you agree to our Terms of Service and Privacy Policy.
        </p>
      </div>
    </div>
  );
};

export default EmailLookup;
