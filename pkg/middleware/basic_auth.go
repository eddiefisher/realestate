package middleware

import "net/http"

// BasicAuthMiddleware ...
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "basic auth", http.StatusForbidden)
		}
		if user != "aminoci" && pass != "password" {
			http.Error(w, "login error", http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}
