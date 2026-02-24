package utils


import (
    "encoding/json"
    "net/http"
)

// Example shape:
// type AppSuccessCRUD struct {
//     Status  int
//     Message string
// }

// shared response type (optional)
type successResponse struct {
    Status  int         `json:"status"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message"`
    Success bool        `json:"success,omitempty"`
}

// Methode classique pour renvoyer la response
func JSONResponseSuccess(w http.ResponseWriter, status int, data interface{}, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    resp := successResponse{
        Status:  status,
        Data:    data,
        Message: message,
    }

    _ = json.NewEncoder(w).Encode(resp)
}

// Methode recommendée pour le CRUD
func JSONResponseSuccessCRUD(w http.ResponseWriter, appSuccess AppSuccessCRUD, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(appSuccess.Status)

    resp := successResponse{
        Status:  appSuccess.Status,
        Data:    data,
        Message: appSuccess.Message,
        Success: true,
    }

    _ = json.NewEncoder(w).Encode(resp)
}

// Methode recommandée pour les autres API des fonctionnalités
func JSONAppSuccess(w http.ResponseWriter, message string, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := successResponse{
        Status:  http.StatusOK,
        Data:    data,
        Message: message,
        Success: true,
    }

    _ = json.NewEncoder(w).Encode(resp)
}
