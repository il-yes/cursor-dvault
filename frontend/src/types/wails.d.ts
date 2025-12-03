// src/types/wails.d.ts
declare global {
  interface Window {
    /** The Wails‑generated Go bridge */
    go: {
      /** Stellar‑related functions */
      Stellar: {
        /** Returns the wallet connection status */
        CheckWalletStatus: () => Promise<string>;
        // add any other Stellar methods you use
      };
      /** IPFS‑related functions */
      IPFS: {
        /** Returns the IPFS node status */
        CheckNodeStatus: () => Promise<string>;
        // add any other IPFS methods you use
      };
      // If you have other Go modules, declare them here
    };
  }
}

/* This file must be a module to augment the global scope */
export {};