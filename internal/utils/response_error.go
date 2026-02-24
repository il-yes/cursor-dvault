package utils


import (
    "encoding/json"
    "net/http"
)

// AppError is assumed to be something like:
// type AppError struct {
//     Status  int
//     Message string
//     Code    string
// }

type errorResponse struct {
    Status   int         `json:"status"`
    Message  string      `json:"message"`
    Error    interface{} `json:"error,omitempty"`
    Code     string      `json:"code"`
    Success  bool        `json:"success"`
}

func JSONAppError(w http.ResponseWriter, appErr AppError, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(appErr.Status)

    resp := errorResponse{
        Status:  appErr.Status,
        Message: appErr.Message,
        Code:    appErr.Code,
        Success: false,
    }
    if err != nil {
        resp.Error = err.Error()
    }

    _ = json.NewEncoder(w).Encode(resp)
}
