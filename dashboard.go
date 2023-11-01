package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type ButtonData struct {
	Data map[string][]interface{} `json:"data"`
}

type IPInfoData struct {
	ID          int    `json:"id"`
	IPAddress   string `json:"ip_address"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	CityName    string `json:"city_name"`
	ASN         string `json:"asn"`
	ASName      string `json:"as_name"`
	DomainName  string `json:"domain_name"`
}

type TCPPortData struct {
	ID        int    `json:"id"`
	IPAddress string `json:"ip_address"`
	PortData  string `json:"scan_result"`
}

type WhoIsInfoData struct {
	ID        int    `json:"id"`
	IPAddress string `json:"ip_address"`
	WhoisData string `json:"whois_result"`
}

type TraceData struct {
	ID             int    `json:"id"`
	IPAddress      string `json:"ip_address"`
	TracerouteData string `json:"traceroute_result"`
}

type PingInfoData struct {
	ID        int    `json:"id"`
	IPAddress string `json:"ip_address"`
	PingData  string `json:"ping_result"`
}

func convertToInterfaceSlice(data interface{}) []interface{} {
	s := reflect.ValueOf(data)
	if s.Kind() != reflect.Slice {
		return nil
	}
	interfaces := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		interfaces[i] = s.Index(i).Interface()
	}
	return interfaces
}

func FetchData(buttonName string, db *sql.DB) ButtonData {
	data := make(map[string][]interface{})

	switch buttonName {
	case "ipinfo":
		ipInfoData := fetchIPInfoDataFromDatabase(db)
		data["ipinfo"] = convertToInterfaceSlice(ipInfoData)
	case "tcpport":
		tcpPortData := fetchPortDataFromDatabase(db)
		data["tcpport"] = convertToInterfaceSlice(tcpPortData)
	case "whois":
		whoisData := fetchWhoisDataFromDatabase(db)
		data["whois"] = convertToInterfaceSlice(whoisData)
	case "traceroute":
		traceData := fetchTracerouteDataFromDatabase(db)
		data["traceroute"] = convertToInterfaceSlice(traceData)
	case "ping":
		pingData := fetchPingDataFromDatabase(db)
		data["ping"] = convertToInterfaceSlice(pingData)
	default:
		data["error"] = nil
	}

	return ButtonData{
		Data: data,
	}
}

func fetchIPInfoDataFromDatabase(db *sql.DB) []IPInfoData {
	rows, err := db.Query("SELECT * FROM ip_info_result")
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return nil
	}
	fmt.Println("Fetching data from the database")
	var data []IPInfoData
	for rows.Next() {
		var id int
		var IP, Country, regionName, City, DomainName, asn, asname string
		err := rows.Scan(&id, &IP, &Country, &regionName, &City, &asn, &asname, &DomainName)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		rowData := IPInfoData{
			ID:          id,
			IPAddress:   IP,
			CountryName: Country,
			RegionName:  regionName,
			CityName:    City,
			ASN:         asn,
			ASName:      asname,
			DomainName:  DomainName,
		}
		data = append(data, rowData)
		fmt.Println(data)
	}

	return data
}
func fetchPortDataFromDatabase(db *sql.DB) []TCPPortData {
	rows, err := db.Query("SELECT * FROM port_scan_result")
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return nil
	}
	fmt.Println("Fetching data from the database")
	var data []TCPPortData
	for rows.Next() {
		var id int
		var IPAddress, PortData string
		err := rows.Scan(&id, &IPAddress, &PortData)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		rowData := TCPPortData{
			ID:        id,
			IPAddress: IPAddress,
			PortData:  PortData,
		}
		data = append(data, rowData)
		fmt.Println(data)
	}

	return data
}

func fetchWhoisDataFromDatabase(db *sql.DB) []WhoIsInfoData {
	rows, err := db.Query("SELECT * FROM whois_result")
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return nil
	}
	fmt.Println("Fetching data from the database")
	var data []WhoIsInfoData
	for rows.Next() {
		var id int
		var IPAddress, WhoisData string
		err := rows.Scan(&id, &IPAddress, &WhoisData)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		rowData := WhoIsInfoData{
			ID:        id,
			IPAddress: IPAddress,
			WhoisData: WhoisData,
		}
		data = append(data, rowData)
		fmt.Println(data)
	}

	return data
}

func fetchTracerouteDataFromDatabase(db *sql.DB) []TraceData {
	rows, err := db.Query("SELECT * FROM traceroute_result")
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return nil
	}
	fmt.Println("Fetching data from the database")
	var data []TraceData
	for rows.Next() {
		var id int
		var IPAddress, TracerouteData string
		err := rows.Scan(&id, &IPAddress, &TracerouteData)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		rowData := TraceData{
			ID:             id,
			IPAddress:      IPAddress,
			TracerouteData: TracerouteData,
		}
		data = append(data, rowData)
		fmt.Println(data)
	}

	return data
}

func fetchPingDataFromDatabase(db *sql.DB) []PingInfoData {
	rows, err := db.Query("SELECT * FROM ping_result")
	if err != nil {
		fmt.Println("Error executing the query:", err)
		return nil
	}
	fmt.Println("Fetching data from the database")
	var data []PingInfoData
	for rows.Next() {
		var id int
		var IPAddress, PingData string
		err := rows.Scan(&id, &IPAddress, &PingData)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		rowData := PingInfoData{
			ID:        id,
			IPAddress: IPAddress,
			PingData:  PingData,
		}
		data = append(data, rowData)
		fmt.Println(data)
	}

	return data
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to fetch data...")
	buttonName := r.URL.Query().Get("data-name")
	//fmt.Printf("button Name: %T\n", buttonName)
	data := FetchData(buttonName, db).Data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	//fmt.Println(jsonData)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
