/**
 * Wails bridge service for desktop app integration
 * Falls back to REST API when running in browser
 */

import { VaultContext } from '@/types/vault';

declare global {
  interface Window {
    backend?: {
      Vault?: {
        SendContext?: (context: VaultContext) => void;
        GetContext?: () => Promise<VaultContext>;
      };
    };
    runtime?: {
      EventsOn?: (eventName: string, callback: (data: any) => void) => void;
      EventsOff?: (eventName: string) => void;
    };
  }
}

type VaultContextCallback = (context: VaultContext) => void;

class WailsBridgeService {
  private listeners: VaultContextCallback[] = [];
  private isWailsEnvironment: boolean;

  constructor() {
    this.isWailsEnvironment = typeof window !== 'undefined' && !!window.backend;
    this.setupListeners();
  }

  private setupListeners(): void {
    if (!this.isWailsEnvironment) return;

    // Listen for vault context updates from Wails
    if (window.runtime?.EventsOn) {
      window.runtime.EventsOn('vault:context:updated', (data: VaultContext) => {
        this.notifyListeners(data);
      });
    }
  }

  /**
   * Register a callback to receive vault context updates
   */
  onContextReceived(callback: VaultContextCallback): () => void {
    this.listeners.push(callback);
    
    // Return unsubscribe function
    return () => {
      this.listeners = this.listeners.filter(cb => cb !== callback);
    };
  }

  /**
   * Manually trigger context hydration from Wails
   */
  async requestContext(): Promise<VaultContext | null> {
    if (this.isWailsEnvironment && window.backend?.Vault?.GetContext) {
      try {
        return await window.backend.Vault.GetContext();
      } catch (error) {
        console.error('Failed to get context from Wails:', error);
      }
    }

    // Fallback to REST API
    return this.fetchContextFromAPI();
  }

  /**
   * REST API fallback for vault context
   */
  private async fetchContextFromAPI(): Promise<VaultContext | null> {
    try {
      const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const userId = localStorage.getItem('user_id'); // TODO: Replace with actual auth
      
      if (!userId) return null;

      const response = await fetch(`${API_BASE_URL}/api/vault/context?user_id=${userId}`);
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch vault context from API:', error);
      return null;
    }
  }

  /**
   * Check if running in Wails environment
   */
  isWails(): boolean {
    return this.isWailsEnvironment;
  }

  private notifyListeners(context: VaultContext): void {
    this.listeners.forEach(callback => {
      try {
        callback(context);
      } catch (error) {
        console.error('Error in vault context listener:', error);
      }
    });
  }

  cleanup(): void {
    if (window.runtime?.EventsOff) {
      window.runtime.EventsOff('vault:context:updated');
    }
    this.listeners = [];
  }
}

export const wailsBridge = new WailsBridgeService();
