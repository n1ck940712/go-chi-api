package middlewares

import (
	"go-chi-api/internal/auth"
	"net/http"
)

func JwtAuthentication(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := auth.TokenValid(r); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func IsAdmin(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := auth.TokenValidIsAdmin(r); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
