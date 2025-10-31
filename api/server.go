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
	Days      int      `json:"days"`
	Features  []string `json:"features,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// 生成License文件接口
func HandleGenerateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求参数
	var req GenerateLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证请求参数
	if req.MachineID == "" {
		sendError(w, "Machine ID is required", http.StatusBadRequest)
		return
	}
	if req.AppID == "" {
		sendError(w, "App ID is required", http.StatusBadRequest)
		return
	}
	if req.Days <= 0 {
		sendError(w, "Days must be positive", http.StatusBadRequest)
		return
	}

	// 生成License文件
	lic, err := license.NewLicense(req.MachineID, req.AppID, req.Days, req.Features)
	if err != nil {
		sendError(w, "Failed to generate license: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create temporary file
	// 创建临时文件
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "license.dat")
	if err := lic.Save(tmpFile); err != nil {
		sendError(w, "Failed to save license: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile)

	// 设置文件下载头
	w.Header().Set("Content-Disposition", "attachment; filename=license.dat")
	w.Header().Set("Content-Type", "application/octet-stream")

	// 发送文件
	http.ServeFile(w, r, tmpFile)
}

// 发送错误响应
func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
