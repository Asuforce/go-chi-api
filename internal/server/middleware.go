package server

import (
	"fmt"
	"net/http"
)

func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "admin" {
				respondError(w, http.StatusUnauthorized, fmt.Errorf("Authorization Forbidden"))
			}
		})
	}
	return
}
