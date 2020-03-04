package middleware

import "net/http"

// BasicAuthMiddleware ...
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Enter pls"`)
			http.Error(w, "basic auth", http.StatusUnauthorized)
			return
		}
		if user != "aminoci" && pass != "password" {
			http.Error(w, "login error", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
