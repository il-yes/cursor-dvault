import { create } from "zustand";
import { persist } from "zustand/middleware";
import { User, VaultPayload } from "@/types/vault";

export interface FeedbackState {
  type: "success" | "error" | "info" | "warning";
  message: string;
}

interface AuthState {
  // Core user data
  user: User | null;
  isLoggedIn: boolean;

  // Tokens
  jwtToken: string | null;
  refreshToken: string | null;

  // Temp onboarding data (if needed)
  onboarding: {
    userId?: number;
    userAlias?: string;
    email?: string;
    password?: string;
    stellar?: any;
    AccountSecretKey?: string;
  };

  // Feedback system
  feedback: FeedbackState | null;

  // Vault cache (optional)
  vault: VaultPayload | null;

  // Actions
  setUser: (user: User | null) => void;
  setLoggedIn: (status: boolean) => void;

  setJwtToken: (token: string | null) => void;
  setRefreshToken: (token: string | null) => void;

  setFeedback: (f: FeedbackState | null) => void;

  updateOnboarding: (updates: Partial<AuthState["onboarding"]>) => void;

  setVault: (vault: VaultPayload | null) => void;
  clearAll: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isLoggedIn: false,

      jwtToken: null,
      refreshToken: null,

      onboarding: {},

      feedback: null,
      vault: null,

      setUser: (user) => set({ user }),
      setLoggedIn: (status) => set({ isLoggedIn: status }),

      setJwtToken: (token) => set({ jwtToken: token }),
      setRefreshToken: (token) => set({ refreshToken: token }),

      setFeedback: (f) => set({ feedback: f }),

      updateOnboarding: (updates) =>
        set({ onboarding: { ...get().onboarding, ...updates } }),

      setVault: (vault) => set({ vault }),

      clearAll: () => {
        set({
          user: null,
          jwtToken: null,
          refreshToken: null,
          isLoggedIn: false,
          onboarding: {},
          vault: null,
          feedback: null,
        });
      }
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        jwtToken: state.jwtToken,
        refreshToken: state.refreshToken,
        isLoggedIn: state.isLoggedIn,
      }),
    }
  )
);
