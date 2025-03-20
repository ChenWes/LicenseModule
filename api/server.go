package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chenwes/licensemodule/internal/license"
)

type GenerateLicenseRequest struct {
	MachineID string   `json:"machine_id"`
	AppID     string   `json:"app_id"`
	SecretKey string   `json:"secret_key"`
	Days      int      `json:"days"`
	Features  []string `json:"features,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleGenerateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.MachineID == "" {
		sendError(w, "Machine ID is required", http.StatusBadRequest)
		return
	}
	if req.SecretKey == "" {
		sendError(w, "Secret key is required", http.StatusBadRequest)
		return
	}
	if req.Days <= 0 {
		sendError(w, "Days must be positive", http.StatusBadRequest)
		return
	}

	// Validate app ID
	if req.AppID == "" {
		sendError(w, "App ID is required", http.StatusBadRequest)
		return
	}

	// Verify secret key
	if req.SecretKey != string(license.SecretKey) {
		sendError(w, "Invalid secret key", http.StatusUnauthorized)
		return
	}

	// Generate license
	lic, err := license.NewLicense(req.MachineID, req.AppID, req.Days, req.Features)
	if err != nil {
		sendError(w, "Failed to generate license: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "license.dat")
	if err := lic.Save(tmpFile); err != nil {
		sendError(w, "Failed to save license: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile)

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=license.dat")
	w.Header().Set("Content-Type", "application/octet-stream")

	// Send file
	http.ServeFile(w, r, tmpFile)
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
