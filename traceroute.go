package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func tracerouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		iptrace := r.FormValue("ip")
		if iptrace == "" {
			http.Error(w, "Please enter an IP address", http.StatusBadRequest)
			return
		}

		output, err := runTraceroute(iptrace)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute traceroute: %v", err), http.StatusInternalServerError)
			return
		}

		err = saveTracerouteResult(iptrace, output)
		if err != nil {
			http.Error(w, "Failed to save traceroute result", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := io.WriteString(w, output); err != nil {
			http.Error(w, "Failed to send traceroute output", http.StatusInternalServerError)
		}
	} else {
		http.ServeFile(w, r, "website/traceroute.html")
	}
}

func runTraceroute(ip string) (string, error) {
	cmd := exec.Command("/usr/bin/traceroute", "-m", "20", ip)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error creating StdoutPipe for traceroute command: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("error executing traceroute command: %v", err)
	}

	lines := make([]string, 0)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning traceroute output for IP %s: %v", ip, err)
	}

	return strings.Join(lines, "\n"), nil
}

func saveTracerouteResult(ip string, tracerouteResult string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM traceroute_result WHERE ip_address = ?", ip).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO traceroute_result (ip_address, traceroute_result) VALUES (?, ?)", ip, tracerouteResult)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE traceroute_result SET traceroute_result = ? WHERE ip_address = ?", tracerouteResult, ip)
		if err != nil {
			return err
		}
	}

	return nil
}
