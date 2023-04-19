package middlewares

import (
	"context"
	"go-chi-api/internal/auth"
	"net/http"
)

func BaseAuthentication(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.TokenValid(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		h.ServeHTTP(w, r.WithContext(ctx))
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
