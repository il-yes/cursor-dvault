import React, { useState, useEffect, useRef } from "react";
import * as AppAPI from "../../wailsjs/go/main/App";

export default function StripePayButton({ 
  onComplete, 
  plainPassword, 
  email,  
  tier, 
  isAnonymous, 
  identity 
}: { 
  onComplete: () => void, 
  plainPassword: string, 
  email: string, 
  tier: string, 
  isAnonymous: boolean, 
  identity: string 
}) {
  const [loading, setLoading] = useState(false);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const intervalRef = useRef<number | null>(null);
  const rail = "standard"

  const pay = async () => {
    if (intervalRef.current) return; // ⛔ already polling

    setLoading(true);
    const bronzePlan = "bronze"
    const url = await AppAPI.GetCheckoutURL(identity, isAnonymous, rail, email, tier, bronzePlan);
    console.log("URL:", url);
    console.log({identity, rail, isAnonymous, email, tier, bronzePlan})
    await AppAPI.OpenURL(url.url);

    intervalRef.current = window.setInterval(async () => {
      try {
        const status = await AppAPI.PollPaymentStatus(url.sessionId, email, plainPassword);
        console.log("STATUS:", status);

        if (status === "paid") {
          console.log("✅ PAYMENT CONFIRMED");

          clearInterval(intervalRef.current!);
          intervalRef.current = null; // 🔑 THIS STOPS EVERYTHING
          setShowConfirmation(true);
          onComplete();
        }
      } catch (e) {
        console.error("Polling error", e);
      }
    }, 1000);
  };

  // 🔑 CLEANUP ON UNMOUNT (THIS IS WHAT YOU WERE MISSING)
  useEffect(() => {
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, []);

  return (
  <div className="w-full max-w-sm rounded-2xl border border-white/20 bg-white/80 p-5">
    {!showConfirmation && (
      <div className="space-y-3 text-center">
        <h2 className="text-lg font-semibold text-slate-900">
          Complete your payment
        </h2>
        <p className="text-xs text-slate-500">
          You’ll be redirected to a secure Stripe page to finalize your payment.
        </p>

        <button
          onClick={pay}
          disabled={loading}
          className="
            mt-2 inline-flex items-center justify-center w-full h-11
            rounded-xl bg-gradient-to-r from-indigo-500 to-blue-500
            text-sm font-semibold text-white
            shadow-md hover:shadow-lg
            hover:from-indigo-400 hover:to-blue-400
            disabled:opacity-60 disabled:cursor-not-allowed
            transition-all
          "
        >
          {loading ? "Opening Stripe…" : "Continue with Stripe"}
        </button>

        <p className="mt-1 text-[11px] text-slate-400">
          Powered by Stripe. You’ll return here automatically once it’s done.
        </p>
      </div>
    )}

    {showConfirmation && (
      <div className="space-y-2 text-center">
        <h1 className="text-lg font-semibold text-emerald-600">
          Payment confirmed
        </h1>
        <p className="text-xs text-slate-500">
          Thank you for your payment. Your subscription is now active.
        </p>
      </div>
    )}
  </div>
);

}

