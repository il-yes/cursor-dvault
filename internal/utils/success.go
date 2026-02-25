package utils

import "net/http"


type AppSuccessCRUD struct {
	Code		string
	Message		string
	Status 		int
}

type AppSuccess struct {
	Message		string
	Status 		int
	Data		any
}

var (
	SuccessRecordCreated = AppSuccessCRUD{
		Code: "RECORD_CREATED",
		Message: "Enregistrement créé avec succès",
		Status: http.StatusCreated,
	}
	SuccessRecordUpdated = AppSuccessCRUD{
		Code: "RECORD_UPDATED",
		Message: "Enregistrement mis à jour avec succès",
		Status: http.StatusOK,
	}
	SuccessRecordDeleted = AppSuccessCRUD{
		Code: "RECORD_DELETED",
		Message: "Enregistrement supprimé avec succès",
		Status: http.StatusOK,
	}
	SuccessRecordFetched = AppSuccessCRUD{
		Code: "RECORD_FETCHED",
		Message: "Enregistrement recupéré avec succès",
		Status: http.StatusOK,
	}
	SuccessLogin = AppSuccessCRUD{
		Code: "LOGIN_SUCCESS",
		Message: "Connexion Réussie",
		Status: http.StatusOK,
	}
)