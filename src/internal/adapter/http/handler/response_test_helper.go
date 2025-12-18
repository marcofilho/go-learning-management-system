package handler

import "net/http"

func HandleUseCaseErrorForTest(w http.ResponseWriter, err error) {
	handleUseCaseError(w, err)
}
