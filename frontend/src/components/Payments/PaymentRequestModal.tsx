import React, { useEffect, useState, useCallback } from "react";
import * as AppAPI from "../../../wailsjs/go/main/App";
import { useAuthStore } from "@/store/useAuthStore";
import { main } from "@/wailsjs/go/models";
import { Dialog, DialogContent, DialogDescription, DialogTitle } from "@radix-ui/react-dialog";
import { DialogHeader } from "../ui/dialog";

type PendingRequest = main.PaymentRequest;

type DecryptedCardData = {
    number: string;
    exp_month: number;
    exp_year: number;
    cvc: string;
};

interface PaymentRequestModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const PaymentRequestModal: React.FC<PaymentRequestModalProps> = ({ open, onOpenChange }) => {
    const [pendingRequests, setPendingRequests] = useState<PendingRequest[]>([]);
    const [selectedRequest, setSelectedRequest] = useState<PendingRequest | null>(null);
    const [processing, setProcessing] = useState(false);
    const [error, setError] = useState("");
    const [decryptedCardData, setDecryptedCardData] = useState<DecryptedCardData | null>(null);
    const { updateOnboarding, onboarding, jwtToken } = useAuthStore();

    // Load pending requests on mount
    useEffect(() => {
        loadPendingRequests();
    }, []);

    const loadPendingRequests = useCallback(async () => {
        try {
            const result = await AppAPI.GetPendingPaymentRequests(jwtToken);
            setPendingRequests(result || []);
        } catch (err) {
            console.error(err);
            setError("Failed to load payment requests");
        }
    }, []);

    const handlePaymentRequest = useCallback(async (request: PendingRequest) => {
        setSelectedRequest(request);
        setError("");
        setDecryptedCardData(null);

        try {
            const vaultKey = await AppAPI.GetUserVaultKey(jwtToken);
            const encryptedEntry = await AppAPI.DecryptVaultEntry(
                request.encrypted_payment_entry_id,
                vaultKey,
            );
            const parsed = JSON.parse(encryptedEntry) as DecryptedCardData;
            setDecryptedCardData(parsed);
        } catch (err) {
            console.error(err);
            setError("Failed to decrypt payment credentials");
        }
    }, []);

    const authorizePayment = useCallback(async () => {
        if (!selectedRequest || !decryptedCardData) return;

        setProcessing(true);
        setError("");

        try {
            const stripe = (window as any).Stripe(
                (import.meta as any).env.VITE_STRIPE_PUBLIC_KEY,
            );
            if (!stripe) {
                throw new Error("Stripe is not available");
            }

            const { paymentMethod, error: stripeError } =
                await stripe.createPaymentMethod({
                    type: "card",
                    card: {
                        number: decryptedCardData.number,
                        exp_month: decryptedCardData.exp_month,
                        exp_year: decryptedCardData.exp_year,
                        cvc: decryptedCardData.cvc,
                    },
                });

            if (stripeError) {
                throw new Error(stripeError.message);
            }

            await AppAPI.ProcessEncryptedPayment({
                payment_request_id: selectedRequest.id,
                stripe_payment_method_id: paymentMethod.id,
            });

            setSelectedRequest(null);
            setDecryptedCardData(null);
            await loadPendingRequests();
            showSuccessMessage("Payment processed successfully!");
        } catch (err: any) {
            console.error(err);
            setError(err?.message || "Payment failed");
        } finally {
            setProcessing(false);
        }
    }, [selectedRequest, decryptedCardData, loadPendingRequests]);

    const cancelPayment = useCallback(() => {
        setSelectedRequest(null);
        setDecryptedCardData(null);
        setError("");
    }, []);

    const showSuccessMessage = (message: string) => {
        console.log(message);
        // Plug your toast system here
    };

    const hasPending = pendingRequests.length > 0;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>

            <DialogContent>  
                <DialogHeader>
                    <DialogTitle>Payment Required</DialogTitle>
                    <DialogDescription>
                        You have {pendingRequests.length} pending payment
                        {pendingRequests.length > 1 ? "s" : ""}
                    </DialogDescription>
                </DialogHeader>   
            {hasPending && (
                <div className="payment-requests-banner">
                    <div className="banner-content">
                        <span className="icon">💳</span>
                        <div className="text">
                            <strong>Payment Required</strong>
                            <p>
                                You have {pendingRequests.length} pending payment
                                {pendingRequests.length > 1 ? "s" : ""}
                            </p>
                        </div>
                        <button
                            className="btn-review"
                            onClick={() => handlePaymentRequest(pendingRequests[0])}
                        >
                            Review Payment
                        </button>
                    </div>
                </div>
            )}

            {selectedRequest && (
                <div className="modal-overlay" onClick={cancelPayment}>
                    <div
                        className="modal-content"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <div className="modal-header">
                            <h2>🔒 Encrypted Payment Authorization</h2>
                            <button className="btn-close" onClick={cancelPayment}>
                                ×
                            </button>
                        </div>

                        <div className="modal-body">
                            <div className="payment-details">
                                <div className="detail-row">
                                    <span className="label">Amount:</span>
                                    <span className="value">
                                        ${selectedRequest.amount.toFixed(2)}
                                    </span>
                                </div>
                                <div className="detail-row">
                                    <span className="label">Reason:</span>
                                    <span className="value">
                                        {selectedRequest.reason.replace("_", " ")}
                                    </span>
                                </div>
                                <div className="detail-row">
                                    <span className="label">Due Date:</span>
                                    <span className="value">
                                        {new Date(selectedRequest.due_date).toLocaleDateString()}
                                    </span>
                                </div>
                            </div>

                            {decryptedCardData && (
                                <div className="card-preview">
                                    <div className="card-visual">
                                        <div className="card-number">
                                            •••• •••• •••• {decryptedCardData.number.slice(-4)}
                                        </div>
                                        <div className="card-expiry">
                                            Expires: {decryptedCardData.exp_month}/
                                            {decryptedCardData.exp_year}
                                        </div>
                                    </div>
                                    <div className="security-notice">
                                        <span className="icon">🔐</span>
                                        <p>
                                            Your card details were decrypted locally and never leave
                                            your device in plain text.
                                        </p>
                                    </div>
                                </div>
                            )}

                            {error && <div className="error-message">{error}</div>}
                        </div>

                        <div className="modal-footer">
                            <button
                                className="btn-cancel"
                                onClick={cancelPayment}
                                disabled={processing}
                            >
                                Cancel
                            </button>
                            <button
                                className="btn-authorize"
                                onClick={authorizePayment}
                                disabled={processing || !decryptedCardData}
                            >
                                {processing
                                    ? "Processing..."
                                    : `Authorize $${selectedRequest.amount.toFixed(2)}`}
                            </button>
                        </div>
                    </div>
                </div>
            )}
            </DialogContent>
        </Dialog>    
    );
};

export default PaymentRequestModal;
