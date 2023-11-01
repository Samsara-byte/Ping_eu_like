package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

)

type PingResult struct {

	Result    string `json:"result"`
}

func savePingResult(ip, result string) error {
	var existingResult string
	err := db.QueryRow("SELECT ping_result FROM ping_result WHERE ip_address=?", ip).Scan(&existingResult)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		_, err = db.Exec("INSERT INTO ping_result (ip_address, ping_result) VALUES (?, ?)", ip, result)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE ping_result SET ping_result=? WHERE ip_address=?", existingResult+"\n"+result, ip)
		if err != nil {
			return err
		}
	}

	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ip := r.FormValue("ip")
		if ip == "" {
			http.Error(w, "Please enter an IP address", http.StatusBadRequest)
			return
		}

		// Perform ping.
		pingResult, err := runPing(ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

            

		err = savePingResult(ip, pingResult)
		if err != nil {
			http.Error(w, "Failed to save ping result", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := PingResult{
			
			Result:    pingResult,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.ServeFile(w, r, "website/ping.html")
	}
}

func runPing(ip string) (string, error) {
	cmd := exec.Command("/bin/ping", "-c", "4", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing ping: %v. Output: %s", err, output)
	}
	return string(output), nil
}
