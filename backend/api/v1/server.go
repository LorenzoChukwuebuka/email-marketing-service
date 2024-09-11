package v1

import (
	"context"
	"email-marketing-service/api/v1/middleware"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/routes"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	router  *mux.Router
	db      *gorm.DB
	logRepo *repository.LogRepository
}

func NewServer(db *gorm.DB) *Server {
	logRepo := repository.NewLogRepository(db)
	return &Server{
		router:  mux.NewRouter(),
		db:      db,
		logRepo: logRepo,
	}
}

var (
	response = &utils.ApiResponse{}
	logger   = &utils.Logger{}
)

func (s *Server) setupRoutes() {

	apiV1Router := s.router.PathPrefix("/api/v1").Subrouter()
	routeMap := map[string]routes.RouteInterface{
		"":            routes.NewAuthRoute(s.db),
		"admin":       routes.NewAdminRoute(s.db),
		"templates":   routes.NewTemplateRoute(s.db),
		"contact":     routes.NewContactRoute(s.db),
		"smtpkey":     routes.NewSMTPKeyRoute(s.db),
		"apikey":      routes.NewAPIKeyRoute(s.db),
		"campaigns":   routes.NewCampaignRoute(s.db),
		"domain":      routes.NewDomainRoute(s.db),
		"sender":      routes.NewSenderRoute(s.db),
		"transaction": routes.NewTransactionRoute(s.db),
	}

	for path, route := range routeMap {
		route.InitRoutes(apiV1Router.PathPrefix("/" + path).Subrouter())
	}

	mode := os.Getenv("SERVER_MODE")

	var staticDir string

	if mode == "" {
		staticDir = "./client"
	} else {
		staticDir = "/app/client"
	}

	// Handle static files using Gorilla Mux
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Handle all other routes by serving index.html for the SPA
	s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, r.URL.Path)

		// If the requested file exists, serve it
		if _, err := os.Stat(path); err == nil {
			http.ServeFile(w, r, path)
			return
		}

		// Otherwise, serve index.html
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	s.router.Use(recoveryMiddleware)
	s.router.Use(enableCORS)
	s.router.Use(methodNotAllowedMiddleware)
	s.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

}

func (s *Server) setupLogger() (*os.File, error) {
	logFile, err := os.OpenFile("requests.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
		return nil, err
	}

	logger := log.New(logFile, "", log.LstdFlags)
	s.router.Use(middleware.LoggingMiddleware(logger))
	return logFile, nil
}

func (s *Server) Start() {
	s.setupRoutes()

	logFile, err := s.setupLogger()
	if err != nil {
		log.Fatal("Failed to set up logger:", err)
	}
	defer logFile.Close()

	server := &http.Server{
		Addr:    ":9000",
		Handler: s.router,
	}

	go func() {
		log.Println("Server started on port 9000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	//graceful shutdown

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	logger.Info("Server shut down gracefully")

}

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
	response.ErrorResponse(w, res)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				fmt.Printf("Panic Stack Trace:\n%s\n", stack)

				response.ErrorResponse(w, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func methodNotAllowedMiddleware(next http.Handler) http.Handler {
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
				response.ErrorResponse(w, res)
				return
			}
		}
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
