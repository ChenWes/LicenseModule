package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chenwes/licensemodule/api"
)

func main() {
	// Define command line parameters
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	// Configure logging
	log.SetPrefix("[LicenseAPI] ")

	// Register handlers
	http.HandleFunc("/api/license/generate", api.HandleGenerateLicense)

	// Start server
	log.Printf("Starting server on port %s...", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
