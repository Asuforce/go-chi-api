package server

import (
	"fmt"
	"net/http"

	"github.com/Asuforce/go-chi-api/internal/handler"
)

func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "admin" {
				handler.RespondError(w, http.StatusUnauthorized, fmt.Errorf("Authorization Forbidden"))
			}
		})
	}
	return
}
