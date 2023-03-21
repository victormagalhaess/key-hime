package controllers

import (
	"net/http"

	"github.com/victormagalhaess/key-hime/server/pkg/api/status"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	status.Success(w, []byte("OK"))
}
