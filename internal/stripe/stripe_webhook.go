package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

// WebhookHandler handles Stripe webhook events
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		log.Println("❌ Missing STRIPE_WEBHOOK_SECRET")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Printf("⚠️ Invalid webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("🔔 Received Stripe event: %s\n", event.Type)

	// Handle relevant events
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("❌ Failed to parse session object: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("✅ Payment completed for session: %s\n", session.ID)

		// TODO: Update your DB/session to mark as paid
		// Example:
		// err := UpdateSessionAsPaid(session.ID)
		// if err != nil {
		//     log.Printf("❌ Failed to update session in DB: %v\n", err)
		// }

	case "invoice.paid":
		// Subscription payment succeeded
		// TODO: Update user subscription in DB
	default:
		log.Printf("ℹ️ Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
