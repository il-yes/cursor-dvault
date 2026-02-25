package utils

import "net/http"


type AppError struct {
	Code		string
	Message		string
	Status 		int
}

var (
	ErrInvalidCredentials = AppError {
		Code: "INVALID_CREDENTIALS",
		Message: "Email ou mot de passe incorrect",
		Status: http.StatusUnauthorized,
	}

	ErrAuteurNotFound = AppError {
		Code: "AUTEUR_NOT_FOUND",
		Message: "Auteur Introuvable",
		Status: http.StatusNotFound,
	}

	ErrRecordNotFound = AppError {
		Code: "RECORD_NOT_FOUND",
		Message: "Enregistrement Introuvable",
		Status: http.StatusNotFound,
	}

	ErrBadRequest = AppError {
		Code: "BAD_REQUEST",
		Message: "Requête Invalid",
		Status: http.StatusBadRequest,
	}

	ErrInternal = AppError {
		Code: "INTERNAL_ERROR",
		Message: "Une erreur interne est survenue",
		Status: http.StatusInternalServerError,
	}

	ErrTokenMissing = AppError {
		Code: "TOKEN_MISSING",
		Message: "Token manquant dans l'en-tete Authorization",
		Status: http.StatusUnauthorized,
	}

	ErrInvalidToken = AppError {
		Code: "INVALID_TOKEN",
		Message: "Token invalid ou expiré",
		Status: http.StatusUnauthorized,
	}

	ErrAccessDenied = AppError {
		Code: "ACCESS_DENIED",
		Message: "Accès Refusé",
		Status: http.StatusForbidden,
	}

	ErrValidationFailed = AppError {
		Code: "VALIDATION_FAILED",
		Message: "Echec de validation des données",
		Status: http.StatusForbidden,
	}

	ErrAlreadyExists = AppError {
		Code: "Already_Exists",
		Message: "Already exists",
		Status: http.StatusAlreadyReported,
	}
)