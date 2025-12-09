// StellarKeyImport.tsx - React version with glassmorphism styling

import React, { useCallback, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { CheckStellarKeyForVault, RecoverVault, ImportStellarKey, ConnectWithStellar, StellarAsksForChallenge } from '../services/api'; // Adjust path
import { cn } from '@/lib/utils';
import { LoginRequest } from '@/types/vault';

interface VaultFoundData {
  id: string;
  created_at: string;
  subscription_tier?: string;
  storage_used_gb?: number;
  last_synced_at?: string;
}

interface StellarKeyImportProps {
  onComplete?: (data: { stellar_key_imported?: boolean; stellar_secret_key?: string }) => void;
  onSkip?: () => void;
}

const StellarKeyImport: React.FC<StellarKeyImportProps> = ({ onComplete, onSkip }) => {
  const [importMethod, setImportMethod] = useState<'new' | 'import'>('new');
  const [stellarSecretKey, setStellarSecretKey] = useState('');
  const [checking, setChecking] = useState(false);
  const [importing, setImporting] = useState(false);
  const [error, setError] = useState('');
  const [vaultFound, setVaultFound] = useState<VaultFoundData | null>(null);
  const [keyValidated, setKeyValidated] = useState(false);

  // Validate Stellar key format & check vault
  const validateAndCheckKey = useCallback(async () => {
    if (!stellarSecretKey.trim()) {
      setError('Please enter your Stellar secret key');
      return;
    }

    // Validate Stellar key format (starts with S, 56 chars)
    if (!stellarSecretKey.startsWith('S') || stellarSecretKey.length !== 56) {
      setError('Invalid Stellar secret key format. Must start with "S" and be 56 characters.');
      return;
    }

    setChecking(true);
    setError('');
    setVaultFound(null);

    try {
      // Ask for challenge First for privacy  
      const { publicKey, signature, challenge } = await StellarAsksForChallenge(stellarSecretKey);
      const payload: LoginRequest = {
        publicKey: publicKey,
        signature: signature,
        signedMessage: challenge
      } 
      const result = await CheckStellarKeyForVault(publicKey);
      // const result = await ConnectWithStellar(payload);
      setKeyValidated(true);
      
      if (result && result.ok) {
        setVaultFound(result);
      } else {
        setVaultFound(null);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to validate Stellar key');
      setKeyValidated(false);
    } finally {
      setChecking(false);
    }
  }, [stellarSecretKey]);

  // Recover existing vault
  const handleRecoverVault = useCallback(async () => {
    setImporting(true);
    setError('');

    try {
      await RecoverVault(stellarSecretKey);
      // Redirect to vault (user logged in)
      window.location.href = '/dashboard/vault';
    } catch (err: any) {
      setError(err.message || 'Failed to recover vault');
    } finally {
      setImporting(false);
    }
  }, [stellarSecretKey]);

  // Import key for new vault
  const handleImportKey = useCallback(async () => {
    setImporting(true);
    setError('');

    try {
      await ImportStellarKey(stellarSecretKey);
      onComplete?.({ stellar_key_imported: true, stellar_secret_key: stellarSecretKey });
    } catch (err: any) {
      setError(err.message || 'Failed to import Stellar key');
    } finally {
      setImporting(false);
    }
  }, [stellarSecretKey, onComplete]);

  const handleSkip = useCallback(() => {
    onSkip?.();
  }, [onSkip]);

  // Reset import form
  const resetKey = useCallback(() => {
    setKeyValidated(false);
    setStellarSecretKey('');
    setVaultFound(null);
    setError('');
  }, []);

  return (
    <div className="w-full max-w-2xl mx-auto p-8 backdrop-blur-xl bg-white/70 dark:bg-zinc-900/70 rounded-3xl ">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold bg-gradient-to-r from-yellow-400 via-yellow-300 to-yellow-500 bg-clip-text text-transparent mb-4">
          Do You Have an Existing Stellar Key?
        </h1>
        <p className="text-xl text-muted-foreground">
          If you've used Ankhora before or have a Stellar account, you can import your key.
        </p>
      </div>

      {/* Method Selector */}
      <div className="grid md:grid-cols-2 gap-6 mb-8">
        <label className={cn(
          "group cursor-pointer p-6 rounded-2xl border-2 border-white/50 bg-white/80 hover:bg-white hover:border-yellow-300/70 backdrop-blur-sm shadow-lg hover:shadow-xl transition-all duration-300 flex flex-col items-center text-center",
          importMethod === 'new' && "ring-2 ring-yellow-400/50 bg-gradient-to-br from-yellow-50/80"
        )}>
          <input
            type="radio"
            name="importMethod"
            value="new"
            checked={importMethod === 'new'}
            onChange={() => setImportMethod('new')}
            className="sr-only"
          />
          <span className="text-3xl mb-4 group-hover:scale-110 transition-transform">✨</span>
          <h3 className="text-xl font-bold mb-2">Create New Account</h3>
          <p className="text-muted-foreground">Generate a new Stellar key and vault</p>
        </label>

        <label className={cn(
          "group cursor-pointer p-6 rounded-2xl border-2 border-white/50 bg-white/80 hover:bg-white hover:border-amber-300/70 backdrop-blur-sm shadow-lg hover:shadow-xl transition-all duration-300 flex flex-col items-center text-center",
          importMethod === 'import' && "ring-2 ring-amber-400/50 bg-gradient-to-br from-amber-50/80"
        )}>
          <input
            type="radio"
            name="importMethod"
            value="import"
            checked={importMethod === 'import'}
            onChange={() => setImportMethod('import')}
            className="sr-only"
          />
          <span className="text-3xl mb-4 group-hover:scale-110 transition-transform">🔑</span>
          <h3 className="text-xl font-bold mb-2">Import Existing Key</h3>
          <p className="text-muted-foreground">Use your existing Stellar secret key</p>
        </label>
      </div>

      {importMethod === 'import' ? (
        /* IMPORT FORM */
        <div className="space-y-6">
          {/* Info Box */}
          <div className="bg-gradient-to-r from-blue-500/10 to-blue-600/10 border border-blue-500/30 rounded-2xl p-6">
            <h4 className="font-bold text-lg mb-3 flex items-center gap-2">
              📌 What We'll Check:
            </h4>
            <ul className="grid grid-cols-1 gap-2 text-sm">
              <li className="flex items-center gap-2">✅ Validate your Stellar key format</li>
              <li className="flex items-center gap-2">✅ Check if an Ankhora vault exists</li>
              <li className="flex items-center gap-2">✅ Recover your vault if found</li>
            </ul>
          </div>

          {/* Key Input */}
          <div className="space-y-2">
            <label className="text-lg font-semibold block mb-2">Your Stellar Secret Key:</label>
            <Input
              type="password"
              value={stellarSecretKey}
              onChange={(e) => setStellarSecretKey(e.target.value)}
              placeholder="S... (56 characters)"
              disabled={checking || importing}
              className="h-14 text-lg font-mono backdrop-blur-sm bg-white/80 border-2 border-white/50 focus:border-yellow-400 focus:ring-yellow-400/50 shadow-inner rounded-2xl px-4 py-3"
            />
            <small className="text-xs text-muted-foreground block">
              Your key starts with "S" and is 56 characters long
            </small>
          </div>

          {error && (
            <div className="bg-gradient-to-r from-red-500/20 to-red-600/20 border border-red-500/30 rounded-2xl p-4">
              <span className="font-medium text-red-700">{error}</span>
            </div>
          )}

          {/* Check Key Button */}
          {!keyValidated ? (
            <Button
              onClick={validateAndCheckKey}
              disabled={checking || !stellarSecretKey.trim() || importing}
              className="w-full h-14 text-lg font-semibold bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black shadow-lg hover:shadow-xl transition-all duration-300"
            >
              {checking ? (
                <>
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2" />
                  Checking...
                </>
              ) : (
                'Check Key'
              )}
            </Button>
          ) : vaultFound ? (
            /* VAULT FOUND */
            <div className="bg-gradient-to-r from-emerald-500/20 to-emerald-600/20 border-2 border-emerald-500/50 rounded-3xl p-8 shadow-2xl">
              <div className="text-center mb-6">
                <div className="w-16 h-16 bg-emerald-500 rounded-full flex items-center justify-center mx-auto mb-4 text-white text-2xl">
                  🎉
                </div>
                <h3 className="text-2xl font-bold text-emerald-800 mb-2">Vault Found!</h3>
              </div>

              <div className="grid md:grid-cols-2 gap-4 mb-8 text-sm">
                <div className="space-y-2 p-4 bg-white/60 rounded-2xl backdrop-blur-sm">
                  <span className="text-muted-foreground font-medium">Vault ID:</span>
                  <span className="font-mono font-semibold">{vaultFound.id.slice(0, 8)}...</span>
                </div>
                <div className="space-y-2 p-4 bg-white/60 rounded-2xl backdrop-blur-sm">
                  <span className="text-muted-foreground font-medium">Created:</span>
                  <span>{new Date(vaultFound.created_at).toLocaleDateString()}</span>
                </div>
                <div className="space-y-2 p-4 bg-white/60 rounded-2xl backdrop-blur-sm md:col-span-2">
                  <span className="text-muted-foreground font-medium">Subscription:</span>
                  <span className="font-semibold text-lg">{vaultFound.subscription_tier || 'Free'}</span>
                </div>
                <div className="space-y-2 p-4 bg-white/60 rounded-2xl backdrop-blur-sm md:col-span-2">
                  <span className="text-muted-foreground font-medium">Storage Used:</span>
                  <span className="font-semibold">{vaultFound.storage_used_gb?.toFixed(2) || 0} GB</span>
                </div>
                {vaultFound.last_synced_at && (
                  <div className="space-y-2 p-4 bg-white/60 rounded-2xl backdrop-blur-sm md:col-span-2">
                    <span className="text-muted-foreground font-medium">Last Synced:</span>
                    <span>{new Date(vaultFound.last_synced_at).toLocaleDateString()}</span>
                  </div>
                )}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Button
                  onClick={handleRecoverVault}
                  disabled={importing}
                  className="h-14 text-lg bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-600 hover:to-emerald-700 text-white shadow-lg hover:shadow-xl transition-all font-semibold"
                >
                  {importing ? 'Recovering...' : 'Recover This Vault'}
                </Button>
                <Button
                  onClick={resetKey}
                  variant="outline"
                  className="h-14 border-zinc-300 hover:border-zinc-400"
                >
                  Try Different Key
                </Button>
              </div>
            </div>
          ) : (
            /* NO VAULT FOUND */
            <div className="bg-gradient-to-r from-blue-500/10 to-blue-600/10 border-2 border-blue-500/30 rounded-3xl p-8 shadow-xl">
              <div className="text-center mb-6">
                <div className="w-16 h-16 bg-blue-500 rounded-full flex items-center justify-center mx-auto mb-4 text-white text-2xl">
                  ℹ️
                </div>
                <h3 className="text-2xl font-bold text-blue-800 mb-2">No Vault Found</h3>
              </div>

              <p className="text-lg text-muted-foreground mb-6 text-center">
                We didn't find an existing Ankhora vault for this Stellar key.
              </p>
              <p className="text-center mb-8 text-muted-foreground">
                You can create a new vault using this key, or try a different key.
              </p>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Button
                  onClick={handleImportKey}
                  disabled={importing}
                  className="h-14 text-lg bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black shadow-lg hover:shadow-xl transition-all font-semibold"
                >
                  {importing ? 'Importing...' : 'Create Vault with This Key'}
                </Button>
                <Button
                  onClick={resetKey}
                  variant="outline"
                  className="h-14 border-zinc-300 hover:border-zinc-400"
                >
                  Try Different Key
                </Button>
              </div>
            </div>
          )}
        </div>
      ) : (
        /* NEW ACCOUNT */
        <div className="text-center space-y-8">
          <div className="bg-gradient-to-r from-emerald-500/10 to-emerald-600/10 border border-emerald-500/30 rounded-3xl p-8">
            <span className="text-4xl mb-4 block">✨</span>
            <h3 className="text-2xl font-bold mb-4">We'll Generate a New Stellar Key for You</h3>
            <p className="text-lg text-muted-foreground mb-4">
              A secure Stellar keypair will be created during vault setup.
            </p>
            <p className="text-muted-foreground">
              You'll be able to save your secret key after account creation.
            </p>
          </div>

          <Button
            onClick={handleSkip}
            size="lg"
            className="h-14 px-12 text-lg bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-black shadow-lg hover:shadow-xl transition-all font-semibold"
          >
            Continue with New Account
          </Button>
        </div>
      )}
    </div>
  );
};

export default StellarKeyImport;
