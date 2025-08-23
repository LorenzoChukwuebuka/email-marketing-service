package middleware

import (
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
)

func EnableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is coming from the same origin as the server
		if r.Header.Get("Origin") == "" {
			// Same-origin request, no need for CORS headers
			handler.ServeHTTP(w, r)
			return
		}

		// For different-origin requests (e.g., during development)
		allowedOrigins := []string{"*", "http://localhost:5054", "https://crabmailer.com", "http://staging.crabmailer.com"}
		origin := r.Header.Get("Origin")

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				break
			}
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})

}


func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	res := map[string]string{
		"error":   "Not Found",
		"message": fmt.Sprintf("The requested resource at %s was not found", r.URL.Path),
	}
	helper.ErrorResponse(w, fmt.Errorf("an error occured"), res)
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				fmt.Printf("Panic Stack Trace:\n%s\n", stack)

				helper.ErrorResponse(w, fmt.Errorf("internal server error"), nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func MethodNotAllowedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)
		if wrappedWriter.statusCode == http.StatusNotFound {
			var match mux.RouteMatch
			if mux.NewRouter().Match(r, &match) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				res := map[string]string{
					"error":   "Method Not Allowed",
					"message": fmt.Sprintf("The requested resource exists, but does not support the %s method", r.Method),
				}
				helper.ErrorResponse(w, fmt.Errorf("an error occured"), res)
				return
			}
		}
	})
}