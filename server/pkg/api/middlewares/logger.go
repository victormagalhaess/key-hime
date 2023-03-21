package middlewares

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Logger(log *log.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
			h.ServeHTTP(w, r)
		})
	}
}
