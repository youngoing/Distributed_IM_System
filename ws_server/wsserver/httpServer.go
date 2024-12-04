package wsserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

type HealthResponse struct {
	Status    string  `json:"status"`
	CPUUsage  float64 `json:"cpu_usage,omitempty"`
	MemUsage  float64 `json:"mem_usage,omitempty"`
	DiskUsage float64 `json:"disk_usage,omitempty"`
}

// Configurable thresholds
var cpuThreshold = 80.0
var memThreshold = 80.0
var diskThreshold = 80.0

func (server *WebSocketServerStruct) runHttpServer(host, httpPort string) {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/health", healthHandler) // 简要健康检查
	http.HandleFunc("/status", statusHandler) // 详细健康检查
	http.HandleFunc("/detail", server.detailHandler)
	address := fmt.Sprintf("%s:%s", host, httpPort)
	logrus.Infof("Starting server on %s\n", address)
	// srv := &http.Server{
	// 	Addr: address,
	// }
	go func() {
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	select {}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong"))
}

// healthHandler returns a simple health status (up or warning)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Create a WaitGroup to synchronize the goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	var cpuUsage, memUsage, diskUsage float64
	var cpuErr, memErr, diskErr error

	// Concurrently get system stats
	go func() {
		defer wg.Done()
		cpuUsage, cpuErr = getCPUUsage()
	}()
	go func() {
		defer wg.Done()
		memUsage, memErr = getMemUsage()
	}()
	go func() {
		defer wg.Done()
		diskUsage, diskErr = getDiskUsage()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// If any error occurred, return internal server error
	if cpuErr != nil || memErr != nil || diskErr != nil {
		http.Error(w, "Error retrieving system metrics", http.StatusInternalServerError)
		return
	}

	// Determine status based on usage
	status := "healthy"
	if cpuUsage > cpuThreshold || memUsage > memThreshold || diskUsage > diskThreshold {
		status = "warning"
	}

	// Create the response
	response := HealthResponse{
		Status: status,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// statusHandler returns detailed system resource usage (CPU, memory, disk)
func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Create a WaitGroup to synchronize the goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	var cpuUsage, memUsage, diskUsage float64
	var cpuErr, memErr, diskErr error

	// Concurrently get system stats
	go func() {
		defer wg.Done()
		cpuUsage, cpuErr = getCPUUsage()
	}()
	go func() {
		defer wg.Done()
		memUsage, memErr = getMemUsage()
	}()
	go func() {
		defer wg.Done()
		diskUsage, diskErr = getDiskUsage()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// If any error occurred, return internal server error
	if cpuErr != nil || memErr != nil || diskErr != nil {
		http.Error(w, "Error retrieving system metrics", http.StatusInternalServerError)
		return
	}

	// Create the response with detailed information
	response := HealthResponse{
		Status:    "healthy",
		CPUUsage:  cpuUsage,
		MemUsage:  memUsage,
		DiskUsage: diskUsage,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getCPUUsage returns the current CPU usage percentage
func getCPUUsage() (float64, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return 0, err
	}
	return percent[0], nil
}

// getMemUsage returns the current memory usage percentage
func getMemUsage() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return 0, err
	}
	return v.UsedPercent, nil
}

// getDiskUsage returns the current disk usage percentage
func getDiskUsage() (float64, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Printf("Error getting disk usage: %v", err)
		return 0, err
	}
	return diskStat.UsedPercent, nil
}

// detailHandler returns some details about the service
func (server *WebSocketServerStruct) detailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 获取服务的详细信息
	details := map[string]interface{}{
		"service":   "websocket_server",
		"version":   "1.0.0",
		"host":      server.host,
		"http_port": server.httpPort,
		"ws_port":   server.webSocketPort,
		"ws_path":   server.WebSocketPath,
		"node_id":   server.nodeId,
		"user_list": server.userManager.userList,
	}

	// 返回 JSON 响应
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(details)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
