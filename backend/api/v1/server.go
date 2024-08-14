package v1

import (
	"context"
	"email-marketing-service/api/v1/routes"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"path/filepath"
	"runtime"
	"syscall"
	"time"
)

type Server struct {
	router *mux.Router
	db     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		router: mux.NewRouter(),
		db:     db,
	}
}

var (
	response = &utils.ApiResponse{}
	logger   = &utils.Logger{}
)

func (s *Server) setupRoutes() {
	apiV1Router := s.router.PathPrefix("/api/v1").Subrouter()
	routeMap := map[string]routes.RouteInterface{
		"":          routes.NewUserRoute(s.db),
		"admin":     routes.NewAdminRoute(s.db),
		"templates": routes.NewTemplateRoute(s.db),
		"contact":   routes.NewContactRoute(s.db),
	}

	for path, route := range routeMap {
		route.InitRoutes(apiV1Router.PathPrefix("/" + path).Subrouter())
	}

	// frontendDir := "./../frontend/dist" // Adjust if needed
	// absPath, err := filepath.Abs(frontendDir)
	// if err != nil {
	// 	log.Printf("Error getting absolute path: %v", err)
	// } else {
	// 	log.Printf("Attempting to serve frontend from: %s", absPath)
	// 	if _, err := os.Stat(absPath); os.IsNotExist(err) {
	// 		logger.Error("Frontend directory does not exist: %s", absPath)
	// 	} else {
	// 		fs := http.FileServer(http.Dir(absPath))
	// 		s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// 		s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			indexFile := filepath.Join(absPath, "index.html")
	// 			http.ServeFile(w, r, indexFile)
	// 		})

	// 		logger.Info("Frontend serving set up successfully with SPA support")

	// 	}
	// }

	s.router.Use(recoveryMiddleware)
	s.router.Use(enableCORS)
	s.router.Use(methodNotAllowedMiddleware)
	s.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}

func (s *Server) setupLogger() {
	logger, err := utils.NewLogger("app.log", utils.INFO)
	if err != nil {
		log.Fatal("logger failed to load")
	}
	defer logger.Close()
}

func (s *Server) Start() {
	s.setupRoutes()
	s.setupLogger()

	server := &http.Server{
		Addr:    ":9000",
		Handler: s.router,
	}

	go func() {
		fmt.Println("Server started on port 9000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

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
