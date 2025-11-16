// hooks/useAuth.tsx
import React from "react";
import { useNavigate } from "react-router-dom";

import * as AppAPI from "../../wailsjs/go/main/App";
import { LoginRequest } from "../types/vault";
import { useAuthStore } from    "@/store/useAuthStore";
import { useVaultStore } from "@/store/vaultStore";
import { useAppStore } from "@/store/appStore";
import { normalizePreloadedVault } from "@/services/normalizeVault";

export function useAuth() {
  const { setLoggedIn, setJwtToken, setRefreshToken, setUser, updateOnboarding } = useAuthStore();
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const [feedback, setFeedback] = React.useState<{ type: "success" | "error" | "info"; message: string } | null>(null);
  const vaultStore = useVaultStore.getState();

  const loginSuccess = async (response: any) => {
    console.log("🔐 Login Response:", response);
    console.log("🔑 Access Token:", response.Tokens?.access_token);

    // Normalize backend shape into what store expects
    const normalized = normalizePreloadedVault(response);
    console.log("🎯 Normalized payload (Stellar):", normalized);

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
    console.log('💾 useAuth: Saving session data');
    useAppStore.getState().setSessionData({
      user: normalized.User,
      vault_runtime_context: normalized.vault_runtime_context || null,
      last_cid: normalized.last_cid,
      dirty: normalized.dirty,
    });
    console.log('✅ useAuth: Session data saved');

    // Load vault into zustand (pass normalized)
    console.log('🚀 useAuth: About to call vaultStore.loadVault with:', {
      hasUser: !!normalized.User,
      hasVault: !!normalized.Vault,
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
      console.log('✅ useAuth: vaultStore.loadVault completed successfully');
    } catch (error) {
      console.error('❌ useAuth: vaultStore.loadVault failed:', error);
      throw error;
    }

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