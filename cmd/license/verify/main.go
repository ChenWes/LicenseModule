package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chenwes/licensemodule/internal/license"
	"github.com/chenwes/licensemodule/pkg/utils"
)

// This program is used to verify the license

func main() {
	// Define command line parameters
	licFile := flag.String("license", license.DefaultLicenseFile, "License file path")
	timeFile := flag.String("timestamp", license.TimeStampFile, "Timestamp file path")
	container := flag.Bool("container", false, "Whether running in container environment")
	appID := flag.String("app", "", "Application ID")
	flag.Parse()

	// Configure logging
	log.SetPrefix("[LicenseVerifier] ")

	// Get current machine ID
	machineID, err := getMachineID(*container)
	if err != nil {
		log.Fatalf("Failed to get machine ID: %v", err)
	}
	log.Printf("Current machine ID: %s", machineID)

	// Ensure file paths are absolute
	licFilePath, err := filepath.Abs(*licFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path for license file: %v", err)
	}

	timeFilePath, err := filepath.Abs(*timeFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path for timestamp file: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(licFilePath); os.IsNotExist(err) {
		log.Fatalf("License file does not exist: %s", licFilePath)
	}

	// Perform verification
	log.Printf("Starting license verification...")
	err = license.VerifyAndUpdate(licFilePath, timeFilePath, machineID, *appID)
	if err != nil {
		log.Fatalf("License verification failed: %v", err)
	}

	// Load license to display more information
	lic, err := license.Load(licFilePath)
	if err != nil {
		log.Fatalf("Failed to load license: %v", err)
	}

	log.Printf("License verification successful!")
	log.Printf("License details:")
	log.Printf("  Machine ID: %s", lic.MachineID)
	log.Printf("  App ID: %s", lic.AppID)
	log.Printf("  Expiry Date: %s", lic.ExpiryDate.Format("2006-01-02 15:04:05"))
	log.Printf("  Features: %v", lic.Features)
	log.Printf("  Creation Date: %s", lic.CreationDate.Format("2006-01-02 15:04:05"))

	fmt.Println("\nApplication can continue running...")
}

// Get machine ID based on environment type
func getMachineID(isContainer bool) (string, error) {
	if isContainer {
		return utils.GetContainerizedMachineID()
	}
	return utils.GetMachineID()
}
