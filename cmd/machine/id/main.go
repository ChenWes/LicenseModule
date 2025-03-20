package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chenwes/licensemodule/pkg/utils"
)

func main() {
	// Define command line parameters
	container := flag.Bool("container", false, "Whether running in container environment")
	flag.Parse()

	// Configure logging
	log.SetPrefix("[MachineID] ")

	// Get machine ID
	id, err := getMachineID(*container)
	if err != nil {
		log.Fatalf("Failed to get machine ID: %v", err)
	}
	fmt.Printf("Current machine ID: %s\n", id)
}

// Get machine ID based on environment type
func getMachineID(isContainer bool) (string, error) {
	if isContainer {
		return utils.GetContainerizedMachineID()
	}
	return utils.GetMachineID()
}
