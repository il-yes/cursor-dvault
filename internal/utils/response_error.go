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

type ErrorResponse struct {
    Status   int         `json:"status"`
    Message  string      `json:"message"`
    Err      interface{} `json:"error,omitempty"`
    Code     string      `json:"code"`
    Success  bool        `json:"success"`
}

func (e ErrorResponse) Error() string {
    return e.Message
}

func JSONAppError(w http.ResponseWriter, appErr AppError, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(appErr.Status)

    resp := ErrorResponse{
        Status:  appErr.Status,
        Message: appErr.Message,
        Code:    appErr.Code,
        Success: false,
    }
    if err != nil {
        resp.Err = err.Error()
    }

    _ = json.NewEncoder(w).Encode(resp)
}

func JSONAppErrorWails(appErr AppError, err error) ErrorResponse {

    resp := ErrorResponse{
        Status:  appErr.Status,
        Message: appErr.Message,
        Code:    appErr.Code,
        Success: false,
    }
    if err != nil {
        resp.Err = err.Error()
    }

    return resp
}