import React from "react";
import { LoginRequest } from "@/types/vault";
import { Keypair } from "stellar-sdk";
import { Buffer } from "buffer";
import * as AppAPI from "../../wailsjs/go/main/App";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Shield } from "lucide-react";
import { toast } from "@/components/ui/use-toast";

export function StellarLoginForm({ onLogin }: { onLogin: (req: LoginRequest) => void }) {
    const [stellarKey, setStellarKey] = React.useState("");
    const [isLoading, setIsLoading] = React.useState(false);

    const handleStellarLogin = async () => {
        // Basic validation
        if (!stellarKey.startsWith("S")) {
            toast({
                title: "Invalid Stellar key",
                description: "A valid Stellar **secret key** starts with 'S'.",
                variant: "destructive",
            });
            return;
        }

        setIsLoading(true);

        try {
            // 1. Generate keypair
            const keypair = Keypair.fromSecret(stellarKey);
            const publicKey = keypair.publicKey();

            // 2. Request challenge from backend
            const { challenge } = await AppAPI.RequestChallenge({ public_key: publicKey });

            // 3. Sign challenge
            const signature = Buffer.from(
                keypair.sign(Buffer.from(challenge))
            ).toString("base64");

            // 4. Notify parent → do full login
            onLogin({
                email: "",
                password: "",
                publicKey,
                signedMessage: challenge,
                signature,
            });

            toast({
                title: "Stellar login successful",
                description: "Signature accepted. Logging you in...",
            });
        } catch (err: any) {
            toast({
                title: "Stellar login failed",
                description: err?.message || "Could not complete login",
                variant: "destructive",
            });
        }

        setIsLoading(false);
    };

    return (
        <>
            <CardHeader>
                <CardTitle className="flex items-center gap-2 text-xl font-semibold">
                    <Shield className="h-5 w-5 text-primary" /> Sign in with Stellar
                </CardTitle>
                <p className="text-sm text-muted-foreground">
                    Use your Stellar private key to authenticate securely.
                </p>
            </CardHeader>

            <CardContent className="space-y-4">
                <Input
                    type="password"
                    value={stellarKey}
                    placeholder="Stellar Secret Key (starts with S...)"
                    onChange={(e) => setStellarKey(e.target.value)}
                    className="font-mono"
                />
            </CardContent>

            <CardFooter>
                <Button
                    onClick={handleStellarLogin}
                    disabled={isLoading}
                    className="w-full"
                >
                    {isLoading ? "Signing..." : "Sign in with Stellar"}
                </Button>
            </CardFooter>
        </>
    );
}
