package v1

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	megabyteDiv uint64 = 1024 * 1024
	gigabyteDiv uint64 = megabyteDiv * 1024

	// WebSocket configuration
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

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

// Global connection manager
type ConnectionManager struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (cm *ConnectionManager) Add(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[conn] = true
}

func (cm *ConnectionManager) Remove(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, conn)
}

func (cm *ConnectionManager) Count() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return len(cm.connections)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connManager = NewConnectionManager()
)

type Client struct {
	conn      *websocket.Conn
	done      chan struct{}
	closeOnce sync.Once
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		done: make(chan struct{}),
	}
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		c.conn.Close()
		connManager.Remove(c.conn)
		log.Printf("Client disconnected. Active connections: %d", connManager.Count())
	})
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(conn)
	connManager.Add(conn)
	log.Printf("New client connected. Active connections: %d", connManager.Count())

	// Configure WebSocket connection
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start ping routine
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			case <-client.done:
				return
			}
		}
	}()

	// Handle reading (just for connection monitoring)
	go func() {
		defer client.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}
		}
	}()

	// Start data sending routine
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer func() {
			ticker.Stop()
			client.Close()
		}()

		for {
			select {
			case <-client.done:
				return
			case <-ticker.C:
				if err := sendSystemData(conn); err != nil {
					log.Printf("Error sending system data: %v", err)
					return
				}
			}
		}
	}()
}

func sendSystemData(conn *websocket.Conn) error {
	systemData, err := getSystemSection()
	if err != nil {
		return fmt.Errorf("error getting system data: %v", err)
	}

	diskData, err := getDiskSection()
	if err != nil {
		return fmt.Errorf("error getting disk data: %v", err)
	}

	cpuData, err := getCpuSection()
	if err != nil {
		return fmt.Errorf("error getting CPU data: %v", err)
	}

	message := map[string]interface{}{
		"system": systemData,
		"disk":   diskData,
		"cpu":    cpuData,
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteJSON(message); err != nil {
		return fmt.Errorf("error writing message: %v", err)
	}

	return nil
}

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
		TotalMemory:       fmt.Sprintf("%d MB", vmStat.Total/megabyteDiv),
		FreeMemory:        fmt.Sprintf("%d MB", vmStat.Free/megabyteDiv),
		UsedMemoryPercent: fmt.Sprintf("%.2f%%", vmStat.UsedPercent),
	}, nil
}

func getDiskSection() (DiskData, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskData{}, err
	}

	return DiskData{
		TotalDiskSpace:  fmt.Sprintf("%d GB", diskStat.Total/gigabyteDiv),
		UsedDiskSpace:   fmt.Sprintf("%d GB", diskStat.Used/gigabyteDiv),
		FreeDiskSpace:   fmt.Sprintf("%d GB", diskStat.Free/gigabyteDiv),
		UsedDiskPercent: fmt.Sprintf("%.2f%%", diskStat.UsedPercent),
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
	for idx, cpuPercent := range percentage {
		cpuUsage = append(cpuUsage, fmt.Sprintf("CPU [%d]: %.2f%%", idx, cpuPercent))
	}

	data := CpuData{
		ModelName: cpuStat[0].ModelName,
		Family:    cpuStat[0].Family,
		Speed:     fmt.Sprintf("%.2f MHz", cpuStat[0].Mhz),
		CpuUsage:  cpuUsage,
	}

	return data, nil
}

func StartSocketServer() {
	// Set up routes
	http.HandleFunc("/ws", handleWebSocket)

	// Set up server with timeouts
	server := &http.Server{
		Addr:              ":9001",
		Handler:           http.DefaultServeMux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server
	go func() {
		log.Println("Socket Server started on port 9001")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	// Close all active WebSocket connections
	connManager.mu.Lock()
	for conn := range connManager.connections {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"))
		conn.Close()
	}
	connManager.mu.Unlock()

	// Shutdown the HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shut down gracefully")
}
