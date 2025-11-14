// hooks/useAuth.tsx
import React from "react";
import { useNavigate } from "react-router-dom";

import * as AppAPI from "../../wailsjs/go/main/App";
import { LoginRequest } from "../types/vault";
import { useAuthStore } from    "@/store/useAuthStore";     

export function useAuth() {
  const { setLoggedIn, setJwtToken, setVault, setUser, updateOnboarding } = useAuthStore();
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const [feedback, setFeedback] = React.useState<{ type: "success" | "error" | "info"; message: string } | null>(null);

  const loginSuccess = (response: any) => {
    setLoggedIn(true);
    localStorage.setItem("userId", JSON.stringify(response.User.id));
    updateOnboarding({ userId: response.User.id });
    setUser(response.User);
    setJwtToken(response.Tokens.access_token);
    setVault(response.Vault);
    setFeedback({ type: "success", message: "Login successful 🎉 Redirecting..." });
    setTimeout(() => navigate("/dashboard"), 1200);
  };

  const loginWithPassword = async (credentials: LoginRequest) => {
    try {
      setLoading(true);
      setFeedback({ type: "info", message: "Signing you in..." });
      const response = await AppAPI.SignIn(credentials);
      loginSuccess(response);
    } catch (err: any) {
      console.error("❌ Login error:", err);
      setFeedback({ type: "error", message: err?.message || "Login failed. Please try again." });
    } finally {
      setLoading(false);
    }
  };

  const loginWithStellar = async (loginReq: LoginRequest) => {
    try {
      setLoading(true);
      setFeedback({ type: "info", message: "Signing you in with Stellar..." });
      const response = await AppAPI.SignInWithStellar(loginReq);
      loginSuccess(response);
    } catch (err: any) {
      console.error("❌ Stellar login error:", err);
      setFeedback({ type: "error", message: "Stellar login failed. Please try again." });
    } finally {
      setLoading(false);
    }
  };

  return { loginWithPassword, loginWithStellar, feedback, setFeedback, loading };
}