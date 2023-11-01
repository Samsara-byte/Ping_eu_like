package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error

type Location struct {
	Country string   `json:"country"`
	Region  string   `json:"region"`
	City    string   `json:"city"`
	Asn     string   `json:"asn"`
	Name    string   `json:"as_name"`
	Domain  []string `json:"domain_name"`
}

func main() {
	db, err = sql.Open("sqlite3", "db/ip-location.db")
	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/checkport", checkportHandler)
	http.HandleFunc("/ipinfo", ipHandler)
	http.HandleFunc("/traceroute", tracerouteHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/WHOIS", whoisHandler)
	http.HandleFunc("/location", getLocationHandler)
	http.HandleFunc("/get_data", getDataHandler)
           
        // Load SSL certificate and private key
    	 certFile := "/etc/letsencrypt/live/samiaslancan.com.tr/fullchain.pem"
   	 keyFile := "/etc/letsencrypt/live/samiaslancan.com.tr/privkey.pem"

  	// Configure HTTPS server
    	err := http.ListenAndServeTLS(":443", certFile, keyFile, nil)
    	if err != nil {
        log.Fatal("ListenAndServeTLS: ", err)
    	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "website/index.html")
	}
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "website/ipinfo.html")
		return
	}

}

func IPToInt(ip string) uint32 {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Fatalf("Invalid IP address: %s", ip)
	}

	ipBytes := parsedIP.To4()
	if ipBytes == nil {
		log.Fatalf("Invalid IPv4 address: %s", ip)
	}

	return (uint32(ipBytes[0]) << 24) + (uint32(ipBytes[1]) << 16) + (uint32(ipBytes[2]) << 8) + uint32(ipBytes[3])
}

func getPTRValues(ip string) []string {
	ptrValues, err := net.LookupAddr(ip)
	if err != nil {
		return []string{}
	}
	fmt.Println(ptrValues)

	return ptrValues
}

func getLocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "website/ipinfo.html")
		return
	}

	ip := r.PostFormValue("ip")
	if ip == "" {
		http.Error(w, "Please enter an IP address", http.StatusBadRequest)
		return
	}

	Domain := getPTRValues(ip)

	ipInt := IPToInt(ip)

	var location Location

	err := db.QueryRow("SELECT country_name, region_name, city_name FROM ip_locations WHERE ip_from <= ? AND ip_to >= ?", ipInt, ipInt).
		Scan(&location.Country, &location.Region, &location.City)
	if err != nil {
		http.Error(w, "Failed to fetch location information", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT asn, as_name FROM asn_info WHERE ip_from <= ? AND ip_to >= ?", ipInt, ipInt).
		Scan(&location.Asn, &location.Name)
	if err != nil {
		http.Error(w, "Failed to fetch ASN information", http.StatusInternalServerError)
		return
	}

	location.Domain = Domain

	if location.Domain == nil {
		location.Domain = []string{}
	}

	err = saveLocationInfo(ip, location)
	if err != nil {
		http.Error(w, "Failed to save location information", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(location)
	if err != nil {
		http.Error(w, "Failed to marshal JSON data", http.StatusInternalServerError)
		return
	}

	fmt.Println(string(jsonData))

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func saveLocationInfo(ip string, location Location) error {
	err := db.QueryRow("SELECT ip_address FROM ip_info_result WHERE ip_address = ?", ip).Scan(&ip)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO ip_info_result (ip_address, country_name, region_name, city_name, asn, as_name, domain_name) VALUES (?, ?, ?, ?, ?, ?, ?)", ip, location.Country, location.Region, location.City, location.Asn, location.Name, strings.Join(location.Domain, ","))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
