import React, { useState } from "react";
// import { CreateStellarPayment } from "../wailsjs/go/app/BillingApp"; // adapt to your API

interface StellarPayFormProps {
  tier: string;
  userId: string;
  onComplete?: () => void;
}

const StellarPayForm: React.FC<StellarPayFormProps> = ({ tier, userId, onComplete }) => {
  const [publicKey, setPublicKey] = useState("");
  const [amount, setAmount] = useState("");
  const [memo, setMemo] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!publicKey || !amount) return;

    setLoading(true);
    setError(null);

    try {
      // Call your backend to build/sign/submit the Stellar payment
      // Example payload – adapt field names to your Go/Wails method:
      /*
      await CreateStellarPayment({
        user_id: userId,
        tier,
        destination: publicKey,
        amount,        // as string, e.g. "10"
        memo: memo || undefined,
      });
      */
      onComplete?.();
    } catch (err: any) {
      setError(err?.message || "Stellar payment failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className="space-y-3" onSubmit={handleSubmit}>
      <div className="space-y-1">
        <label className="text-xs font-medium text-slate-500">Destination address</label>
        <input
          type="text"
          placeholder="G... (Stellar public key)"
          className="w-full h-11 rounded-xl border border-amber-200/60 bg-amber-50/40 px-3 text-sm
                     placeholder:text-amber-900/50 focus:outline-none focus:ring-2 focus:ring-amber-300
                     focus:border-transparent"
          value={publicKey}
          onChange={(e) => setPublicKey(e.target.value.trim())}
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <label className="text-xs font-medium text-slate-500">Amount (XLM / USDC)</label>
          <input
            type="number"
            min="0"
            step="0.000001"
            placeholder="e.g. 10"
            className="w-full h-11 rounded-xl border border-amber-200/60 bg-amber-50/40 px-3 text-sm
                       placeholder:text-amber-900/50 focus:outline-none focus:ring-2 focus:ring-amber-300
                       focus:border-transparent"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
          />
        </div>

        <div className="space-y-1">
          <label className="text-xs font-medium text-slate-500">Memo (optional)</label>
          <input
            type="text"
            placeholder="Subscription ref"
            className="w-full h-11 rounded-xl border border-amber-200/60 bg-amber-50/40 px-3 text-sm
                       placeholder:text-amber-900/50 focus:outline-none focus:ring-2 focus:ring-amber-300
                       focus:border-transparent"
            value={memo}
            onChange={(e) => setMemo(e.target.value)}
          />
        </div>
      </div>

      {error && (
        <div className="rounded-xl border border-red-300/60 bg-red-50/80 px-3 py-2 text-xs text-red-700">
          {error}
        </div>
      )}

      <button
        type="submit"
        disabled={loading || !publicKey || !amount}
        className="w-full h-11 rounded-xl text-[13px] font-semibold
                   bg-gradient-to-r from-amber-400 to-amber-500 text-black
                   shadow-md shadow-amber-500/30
                   hover:from-amber-300 hover:to-amber-400 hover:shadow-lg
                   disabled:opacity-60 disabled:cursor-not-allowed
                   transition-all"
      >
        {loading ? "Submitting Stellar payment..." : "Pay with Stellar"}
      </button>
    </form>
  );
};

export default StellarPayForm;
