package onboarding_ui_wails

import (
	"context"
	"encoding/json"
	"net/http"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
)

// HTTP handler adapter — thin controller returning an http.Handler (router)
func NewOnboardingHandler(uc *onboarding_usecase.OnboardUseCase) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/onboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// For brevity: decode minimal JSON payload
		var req onboarding_usecase.OnboardRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		res, err := uc.Execute(context.Background(), req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})
	return mux
}