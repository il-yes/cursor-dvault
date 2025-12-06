package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"

	// identity
	identity "vault-app/internal/identity/application/usecase"
	idPersistence "vault-app/internal/identity/infrastructure/persistence"

	// subscription
	subscription_usecase "vault-app/internal/subscription/application/usecase"
	subPersistence "vault-app/internal/subscription/infrastructure/persistence"

	// billing
	billing_usecase "vault-app/internal/billing/application/usecase"
	billPersistence "vault-app/internal/billing/infrastructure/persistence"

	// onboarding
	onbUI "vault-app/internal/onboarding/ui/wails"
)

func idGen() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}


func main() {
	// identity wiring
	idRepo := idPersistence.NewMemoryUserRepository()
	// use eventbus functions from shared package (simple global memory bus)
	// create identity usecases
	regStd := identity.NewRegisterStandardUserUseCase(idRepo, sharedBusWrapper{}, idGen)
	regAnon := identity.NewRegisterAnonymousUserUseCase(idRepo, sharedBusWrapper{}, idGen)

	// subscription wiring
	subRepo := subPersistence.NewMemorySubscriptionRepository()
	subCreate := subscription_usecase.NewCreateSubscriptionUseCase(subRepo, sharedBusWrapper{}, idGen)

	// billing wiring
	billRepo := billPersistence.NewMemoryBillingRepository()
	addPayment := billing_usecase.NewAddPaymentMethodUseCase(billRepo, sharedBusWrapper{}, idGen)

	// vault: minimal stub implementing VaultPort
	vault := &localVault{}

	// onboarding usecase
	onb := onboarding_usecase.NewOnboardUseCase(regAdapter{regStd, regAnon}, billingAdapter{addPayment}, subscriptionAdapter{subCreate}, vault)

	// HTTP handler
	h := onbUI.NewOnboardingHandler(onb)
	server := &http.Server{Addr: ":8081", Handler: h}

	log.Println("starting onboarding http server :8081")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}