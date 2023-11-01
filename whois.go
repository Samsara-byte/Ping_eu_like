package main

import (
	"encoding/json"
	"net/http"

	whois "github.com/likexian/whois"
)

type WhoisResult struct {
	Result string `json:"result"`
}

func whoisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ip := r.PostFormValue("ip")
		if ip == "" {
			http.Error(w, "Please enter an IP address or domain name", http.StatusBadRequest)
			return
		}

		result, err := whois.Whois(ip)
		if err != nil {
			http.Error(w, "Failed to perform WHOIS lookup", http.StatusInternalServerError)
			return
		}
		err = saveWhoisResult(ip, result)
		if err != nil {
			http.Error(w, "Failed to save WHOIS result", http.StatusInternalServerError)
			return
		}

		response := WhoisResult{
			Result: result,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.ServeFile(w, r, "website/WHOIS.html")
	}
}

func saveWhoisResult(ip string, whoisResult string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM whois_result WHERE ip_address = ?", ip).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO whois_result (ip_address, whois_result) VALUES (?, ?)", ip, whoisResult)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE whois_result SET whois_result = ? WHERE ip_address = ?", whoisResult, ip)
		if err != nil {
			return err
		}
	}

	return nil
}
