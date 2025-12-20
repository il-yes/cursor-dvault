import React, { useState, useEffect, useRef } from "react";
import * as AppAPI from "../../wailsjs/go/main/App";

export default function StripePayButton({ onComplete, plainPassword }: { onComplete: () => void, plainPassword: string }) {
  const [loading, setLoading] = useState(false);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const intervalRef = useRef<number | null>(null);

  const pay = async () => {
    if (intervalRef.current) return; // ⛔ already polling

    setLoading(true);

    const url = await AppAPI.GetCheckoutURL("bronze");
    await AppAPI.OpenURL(url.url);

    intervalRef.current = window.setInterval(async () => {
      try {
        const status = await AppAPI.PollPaymentStatus(url.sessionId, plainPassword);
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
    <>
      {!showConfirmation && <button onClick={pay} disabled={loading}>
        Stripe $5
      </button>}

      {showConfirmation && (
        <div>
          <h1>Payment Confirmed</h1>
          <p>Thank you for your payment</p>
        </div>
      )}
    </>
  );
}

