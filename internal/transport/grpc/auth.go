package grpc

import (
	"crypto/subtle"
	"net/http"
)

// basicAuth 中间件
func (s *Server) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(username), []byte(s.config.Swagger.Username)) != 1 || subtle.ConstantTimeCompare([]byte(password), []byte(s.config.Swagger.Password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Swagger UI"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
