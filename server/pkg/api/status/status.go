package status

import (
	"net/http"
)

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func ServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func Success(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func Created(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
