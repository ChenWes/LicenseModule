package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/chenwes/licensemodule/internal/license"
	"github.com/chenwes/licensemodule/pkg/utils"
)

type Response struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GenerateRequest struct {
	MachineID string   `json:"machine_id"`
	AppID     string   `json:"app_id"`
	Days      int      `json:"days"`
	Features  []string `json:"features,omitempty"`
}

type VerifyRequest struct {
	LicenseFile   string `json:"license_file"`
	TimestampFile string `json:"timestamp_file"`
	MachineID     string `json:"machine_id"`
	AppID         string `json:"app_id"`
}

func handleGetMachineID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isContainer := r.URL.Query().Get("container") == "true"

	var id string
	var err error

	if isContainer {
		id, err = utils.GetContainerizedMachineID()
	} else {
		id, err = utils.GetMachineID()
	}

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    id,
	})
}

func handleGenerateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lic, err := license.NewLicense(req.MachineID, req.AppID, req.Days, req.Features)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert license to JSON string
	licenseData, err := json.Marshal(lic)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    string(licenseData),
	})
}

func handleVerifyLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := license.VerifyAndUpdate(req.LicenseFile, req.TimestampFile, req.MachineID, req.AppID)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    "License verified successfully",
	})
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	log.SetPrefix("[LicenseHTTPServer] ")

	http.HandleFunc("/machine-id", handleGetMachineID)
	http.HandleFunc("/generate", handleGenerateLicense)
	http.HandleFunc("/verify", handleVerifyLicense)

	log.Printf("Starting HTTP server on port %s...", *port)
	log.Printf("Endpoints:")
	log.Printf("  GET  /machine-id?container=false")
	log.Printf("  POST /generate")
	log.Printf("  POST /verify")

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
