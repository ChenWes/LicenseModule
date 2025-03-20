package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultLicenseFile = "license.dat"
	TimeStampFile      = "timestamp.dat"
)

var (
	// Secret key used for signature verification, should be properly secured in production
	SecretKey = []byte("0aea8a18b07463ad5f5e3318db20d527c912c4ab9e7be28e94e8f486263a86fd/CF/WESCHAN")

	ErrInvalidLicense        = errors.New("invalid license")
	ErrExpiredLicense        = errors.New("license has expired")
	ErrInvalidSignature      = errors.New("invalid license signature")
	ErrSystemTimeManipulated = errors.New("system time has been manipulated")
	ErrMachineMismatch       = errors.New("license does not match current machine")
)

// License represents a software license
type License struct {
	MachineID    string    `json:"machine_id"`    // Unique machine identifier
	AppID        string    `json:"app_id"`        // Application identifier
	ExpiryDate   time.Time `json:"expiry_date"`   // Expiration time
	Features     []string  `json:"features"`      // Optional feature list
	Signature    string    `json:"signature"`     // Digital signature
	CreationDate time.Time `json:"creation_date"` // Creation time
}

// TimestampRecord used to prevent system time manipulation
type TimestampRecord struct {
	LastRun time.Time `json:"last_run"` // Last execution time
}

// NewLicense creates a new license
func NewLicense(machineID string, appID string, expiryDays int, features []string) (*License, error) {
	if machineID == "" {
		return nil, errors.New("machine ID cannot be empty")
	}
	if appID == "" {
		return nil, errors.New("app ID cannot be empty")
	}

	now := time.Now()
	expiryDate := now.AddDate(0, 0, expiryDays)

	license := &License{
		MachineID:    machineID,
		AppID:        appID,
		ExpiryDate:   expiryDate,
		Features:     features,
		CreationDate: now,
	}

	// Generate signature
	if err := license.Sign(); err != nil {
		return nil, err
	}

	return license, nil
}

// Sign adds a signature to the license
func (l *License) Sign() error {
	l.Signature = "" // Clear old signature
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	// Calculate signature using HMAC-SHA256
	h := hmac.New(sha256.New, SecretKey)
	h.Write(data)
	signature := h.Sum(nil)
	l.Signature = base64.StdEncoding.EncodeToString(signature)
	return nil
}

// Verify checks if the license is valid
func (l *License) Verify(currentMachineID string, appID string) error {
	// Verify machine ID
	if l.MachineID != currentMachineID {
		return ErrMachineMismatch
	}

	// Verify app ID
	if l.AppID != appID {
		return errors.New("license does not match application ID")
	}

	// Verify expiration time
	if time.Now().After(l.ExpiryDate) {
		return ErrExpiredLicense
	}

	// Verify signature
	signature := l.Signature
	l.Signature = ""
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	l.Signature = signature

	h := hmac.New(sha256.New, SecretKey)
	h.Write(data)
	expectedSignature := h.Sum(nil)
	actualSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	if !hmac.Equal(actualSignature, expectedSignature) {
		return ErrInvalidSignature
	}

	return nil
}

// Save saves the license to a file
func (l *License) Save(filePath string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

// Load loads a license from a file
func Load(filePath string) (*License, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return nil, ErrInvalidLicense
	}

	return &license, nil
}

// UpdateTimestamp updates the last run timestamp
func UpdateTimestamp(filePath string) error {
	record := TimestampRecord{
		LastRun: time.Now(),
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// CheckTimestamp checks the last run time to prevent system time manipulation
func CheckTimestamp(filePath string) error {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		// Create a new timestamp file if it doesn't exist
		return UpdateTimestamp(filePath)
	}
	if err != nil {
		return err
	}

	var record TimestampRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return err
	}

	// Check if current time is earlier than last run time (time rollback)
	// Add a small tolerance period (e.g., 10 minutes) to allow for time sync or minor adjustments
	if time.Now().Add(10 * time.Minute).Before(record.LastRun) {
		return ErrSystemTimeManipulated
	}

	// Update timestamp
	return UpdateTimestamp(filePath)
}

// VerifyAndUpdate verifies the license and updates the timestamp
func VerifyAndUpdate(licenseFilePath, timestampFilePath, currentMachineID, appID string) error {
	// Check if system time has been manipulated
	if err := CheckTimestamp(timestampFilePath); err != nil {
		return fmt.Errorf("timestamp check failed: %w", err)
	}

	// Load license
	license, err := Load(licenseFilePath)
	if err != nil {
		return fmt.Errorf("failed to load license: %w", err)
	}

	// Verify license
	if err := license.Verify(currentMachineID, appID); err != nil {
		return fmt.Errorf("license verification failed: %w", err)
	}

	return nil
}
