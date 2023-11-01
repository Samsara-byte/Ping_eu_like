package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ScanResult struct {
	Port      int    `json:"port"`
	IsOpen    bool   `json:"is_open"`
	Timestamp string `json:"timestamp"`
}

func scanTCP(ip string, port int) ScanResult {
	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	result := ScanResult{
		Port:      port,
		IsOpen:    false,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err == nil {
		conn.Close()
		result.IsOpen = true
	}

	return result
}

func checkportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ip := r.FormValue("ip")
		portStr := r.FormValue("port")
		if ip == "" || portStr == "" {
			http.Error(w, "Please enter both IP address and port", http.StatusBadRequest)
			return
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			http.Error(w, "Invalid port number", http.StatusBadRequest)
			return
		}

		result := scanTCP(ip, port)

		// Save scan result to the database
		err = saveScanResult(ip, result)
		if err != nil {
			http.Error(w, "Failed to save scan result", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if result.IsOpen {
			w.Write([]byte(fmt.Sprintf("Port %d is open on %s", result.Port, ip)))
		} else {
			w.Write([]byte(fmt.Sprintf("Port %d is closed on %s", result.Port, ip)))
		}
	} else {
		http.ServeFile(w, r, "website/checkport.html")
	}
}

func saveScanResult(ip string, result ScanResult) error {
	_, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal scan result: %v", err)
	}

	var existingResults string
	err = db.QueryRow("SELECT scan_result FROM port_scan_result WHERE ip_address = ?", ip).Scan(&existingResults)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to fetch existing results: %v", err)
	}

	var scanResults []ScanResult

	if existingResults != "" {
		if err := json.Unmarshal([]byte(existingResults), &scanResults); err != nil {
			return fmt.Errorf("failed to unmarshal existing results: %v", err)
		}
	}

	scanResults = append(scanResults, result)

	updatedData, err := json.Marshal(scanResults)
	if err != nil {
		return fmt.Errorf("failed to marshal updated results: %v", err)
	}

	if len(existingResults) > 0 {
		_, err = db.Exec("UPDATE port_scan_result SET scan_result = ? WHERE ip_address = ?", string(updatedData), ip)
		if err != nil {
			return fmt.Errorf("failed to update row: %v", err)
		}
	} else {
		_, err = db.Exec("INSERT INTO port_scan_result (ip_address, scan_result) VALUES (?, ?)", ip, string(updatedData))
		if err != nil {
			return fmt.Errorf("failed to insert new row: %v", err)
		}
	}

	return nil
}
