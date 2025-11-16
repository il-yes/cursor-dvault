import { create } from "zustand";

interface AppStore {
  user: any | null;
  loggedIn: boolean | false;
  session: any | null;

  setUser: (user: any | null) => void;
  setLoggedIn: (v: boolean | false) => void;

  setSessionData: (payload: any) => void;
  reset: () => void;
}

export const useAppStore = create<AppStore>((set) => ({
  user: null,
  loggedIn: false,
  session: null,

  setUser: (user) => set({ user }),
  setLoggedIn: (v) => set({ loggedIn: v }),

  setSessionData: (sessionData) =>
    set(() => ({
      session: {
        ...sessionData,
      },
    })),

  reset: () =>
    set({
      user: null,
      loggedIn: false,
      session: null,
    }),
}));
