package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chenwes/licensemodule/internal/license"
	"github.com/chenwes/licensemodule/pkg/utils"
)

// 版本信息，通过 ldflags 在编译时注入
var (
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Define command line parameters
	machineID := flag.String("machine", "", "Machine ID (if empty, will use current machine's ID)")
	appID := flag.String("app", "", "Application ID")
	days := flag.Int("days", 30, "License validity period (days)")
	outFile := flag.String("out", license.DefaultLicenseFile, "Output file path")
	container := flag.Bool("container", false, "Whether to generate license for container environment")
	features := flag.String("features", "", "Optional feature list, comma separated")
	showMachineID := flag.Bool("show-id", false, "Only show current machine ID, don't generate license")
	flag.Parse()

	// Configure logging
	log.SetPrefix("[LicenseGenerator] ")

	log.Printf("CF License Generation Service Start: Version: %s, Git Commit: %s", version, gitCommit)

	// If only showing machine ID
	if *showMachineID {
		id, err := getMachineID(*container)
		if err != nil {
			log.Fatalf("Failed to get machine ID: %v", err)
		}
		fmt.Printf("Current machine ID: %s\n", id)
		return
	}

	// Get machine ID
	var id string
	var err error
	if *machineID == "" {
		// Use current machine's ID
		id, err = getMachineID(*container)
		if err != nil {
			log.Fatalf("Failed to get machine ID: %v", err)
		}
		log.Printf("Using current machine ID: %s", id)
	} else {
		// Use provided machine ID
		id = *machineID
		log.Printf("Using provided machine ID: %s", id)
	}

	// Parse feature list
	var featureList []string
	if *features != "" {
		featureList = strings.Split(*features, ",")
	}

	// Create License
	lic, err := license.NewLicense(id, *appID, *days, featureList)
	if err != nil {
		log.Fatalf("Failed to create license: %v", err)
	}

	// Display License information
	log.Printf("License created:")
	log.Printf("  Machine ID: %s", lic.MachineID)
	log.Printf("  App ID: %s", lic.AppID)
	log.Printf("  Expiry Date: %s", lic.ExpiryDate.Format(time.RFC3339))
	log.Printf("  Features: %v", lic.Features)
	log.Printf("  Creation Date: %s", lic.CreationDate.Format(time.RFC3339))

	// Save to file
	absPath, err := filepath.Abs(*outFile)
	if err != nil {
		log.Printf("Warning: Cannot get absolute path: %v", err)
		absPath = *outFile
	}

	// Ensure directory exists
	dir := filepath.Dir(absPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Cannot create directory: %v", err)
		}
	}

	if err := lic.Save(*outFile); err != nil {
		log.Fatalf("Failed to save license: %v", err)
	}
	log.Printf("License saved to: %s", absPath)

	// Print license in JSON format (optional, for debugging)
	jsonData, _ := json.MarshalIndent(lic, "", "  ")
	fmt.Println("\nLicense JSON:")
	fmt.Println(string(jsonData))
}

// Get machine ID based on environment type
func getMachineID(isContainer bool) (string, error) {
	if isContainer {
		return utils.GetContainerizedMachineID()
	}
	return utils.GetMachineID()
}
