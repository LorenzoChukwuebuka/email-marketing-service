package middleware 

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func LoggingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			// Log the request details
			logger.Printf("%s %s %s %s\n", r.Method, r.RequestURI, r.Proto, r.RemoteAddr)

			next.ServeHTTP(wrappedWriter, r)

			// Log the response details
			logger.Printf("Status: %d, Duration: %s\n", wrappedWriter.statusCode, time.Since(start))
		})

	}

}
