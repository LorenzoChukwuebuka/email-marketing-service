package v1

import (
	"context"
	cronjobs "email-marketing-service/api/v1/jobs"
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
	"strconv"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"gorm.io/gorm"
)

type Server struct {
	router  *mux.Router
	db      *gorm.DB
	logRepo *repository.LogRepository
}

type SystemData struct {
	OperatingSystem   string `json:"operating_system"`
	Platform          string `json:"platform"`
	Hostname          string `json:"hostname"`
	NumProcesses      uint64 `json:"num_processes"`
	TotalMemory       string `json:"total_memory"`
	FreeMemory        string `json:"free_memory"`
	UsedMemoryPercent string `json:"used_memory_percent"`
}

type DiskData struct {
	TotalDiskSpace  string `json:"total_disk_space"`
	UsedDiskSpace   string `json:"used_disk_space"`
	FreeDiskSpace   string `json:"free_disk_space"`
	UsedDiskPercent string `json:"used_disk_percent"`
}

type CpuData struct {
	ModelName string   `json:"model_name"`
	Family    string   `json:"family"`
	Speed     string   `json:"speed"`
	CpuUsage  []string `json:"cpu_usage"`
}

const (
	megabyteDiv uint64 = 1024 * 1024
	gigabyteDiv uint64 = megabyteDiv * 1024
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Be careful with this in production
	},
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
		"":               routes.NewAuthRoute(s.db),
		"admin":          routes.NewAdminRoute(s.db),
		"templates":      routes.NewTemplateRoute(s.db),
		"contact":        routes.NewContactRoute(s.db),
		"smtpkey":        routes.NewSMTPKeyRoute(s.db),
		"apikey":         routes.NewAPIKeyRoute(s.db),
		"campaigns":      routes.NewCampaignRoute(s.db),
		"domain":         routes.NewDomainRoute(s.db),
		"sender":         routes.NewSenderRoute(s.db),
		"transaction":    routes.NewTransactionRoute(s.db),
		"support":        routes.NewSupportRoute(s.db),
		"admin/users":    routes.NewAdminUsersRoute(s.db),
		"admin/support":  routes.NewAdminSupportRoute(s.db),
		"admin/campaign": routes.NewAdminCampaignRoute(s.db),
		"system":         routes.NewSystemRoute(s.db),
	}

	for path, route := range routeMap {
		route.InitRoutes(apiV1Router.PathPrefix("/" + path).Subrouter())
	}

	s.router.Use(recoveryMiddleware)
	s.router.Use(enableCORS)
	s.router.Use(methodNotAllowedMiddleware)
	s.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	uploadsDir := filepath.Join(".", "uploads", "tickets")
	s.router.PathPrefix("/uploads/tickets/").Handler(http.StripPrefix("/uploads/tickets/", http.FileServer(http.Dir(uploadsDir))))

	// mode := os.Getenv("SERVER_MODE")

	// var staticDir string

	// if mode == "" {
	// 	staticDir = "./client"
	// } else if mode == "test" {
	// 	staticDir = "/app/client"
	// } else {
	// 	staticDir = "/app/client"
	// }

	// // Handle static files using Gorilla Mux
	// s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// // Handle all other routes by serving index.html for the SPA
	// s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	path := filepath.Join(staticDir, r.URL.Path)

	// 	// If the requested file exists, serve it
	// 	if _, err := os.Stat(path); err == nil {
	// 		http.ServeFile(w, r, path)
	// 		return
	// 	}

	// 	// Otherwise, serve index.html
	// 	http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	// })

	s.router.HandleFunc("/api/v1/system/monitor", s.handleWebSocket)

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

	c := cronjobs.SetupCronJobs(s.db)

	go func() {
		c.Start()
		defer c.Stop()
		select {}
	}()

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
		// Check if the request is coming from the same origin as the server
		if r.Header.Get("Origin") == "" {
			// Same-origin request, no need for CORS headers
			handler.ServeHTTP(w, r)
			return
		}

		// For different-origin requests (e.g., during development)
		allowedOrigins := []string{"http://localhost:5054", "http://localhost:5000", "https://crabmailer.com", "https://beta.crabmailer.com", "*"}
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

// WebSocket handler method
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	// Create a done channel to signal goroutine termination
	done := make(chan bool)
	defer close(done)

	// Handle client disconnection
	go func() {
		_, _, err := ws.ReadMessage()
		if err != nil {
			done <- true
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			systemData, err := getSystemSection()
			if err != nil {
				log.Printf("Error getting system data: %v", err)
				continue
			}

			diskData, err := getDiskSection()
			if err != nil {
				log.Printf("Error getting disk data: %v", err)
				continue
			}

			cpuData, err := getCpuSection()
			if err != nil {
				log.Printf("Error getting CPU data: %v", err)
				continue
			}

			message := map[string]interface{}{
				"system": systemData,
				"disk":   diskData,
				"cpu":    cpuData,
			}

			if err := ws.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// Helper functions for getting system information
func getSystemSection() (SystemData, error) {
	runTimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return SystemData{}, err
	}

	hostStat, err := host.Info()
	if err != nil {
		return SystemData{}, err
	}

	return SystemData{
		OperatingSystem:   runTimeOS,
		Platform:          hostStat.Platform,
		Hostname:          hostStat.Hostname,
		NumProcesses:      hostStat.Procs,
		TotalMemory:       strconv.FormatUint(vmStat.Total/megabyteDiv, 10) + " MB",
		FreeMemory:        strconv.FormatUint(vmStat.Free/megabyteDiv, 10) + " MB",
		UsedMemoryPercent: strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%",
	}, nil
}

func getDiskSection() (DiskData, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskData{}, err
	}

	return DiskData{
		TotalDiskSpace:  strconv.FormatUint(diskStat.Total/gigabyteDiv, 10) + " GB",
		UsedDiskSpace:   strconv.FormatUint(diskStat.Used/gigabyteDiv, 10) + " GB",
		FreeDiskSpace:   strconv.FormatUint(diskStat.Free/gigabyteDiv, 10) + " GB",
		UsedDiskPercent: strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%",
	}, nil
}

func getCpuSection() (CpuData, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return CpuData{}, err
	}

	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return CpuData{}, err
	}

	var cpuUsage []string
	for idx, cpupercent := range percentage {
		cpuUsage = append(cpuUsage, "CPU ["+strconv.Itoa(idx)+"]: "+strconv.FormatFloat(cpupercent, 'f', 2, 64)+"%")
	}

	data := CpuData{
		ModelName: cpuStat[0].ModelName,
		Family:    cpuStat[0].Family,
		Speed:     strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz",
		CpuUsage:  cpuUsage,
	}

	return data, nil
}
