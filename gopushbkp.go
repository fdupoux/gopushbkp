package main

import (
	"fmt"
	"os"
)

var config configuration

func main() {

	var err error

	// Read configuration file
	config, err = readConfig()
	if err != nil {
		fmt.Printf("Failed to read configuration file\n")
		os.Exit(1)
	}

	// Create the encrypted archive locally
	archfile, err := archiveData()
	if err != nil {
		fmt.Printf("Failed to create backup archive: %s\n", err)
		os.Exit(2)
	}

	// Upload archive to the cloud
	err = uploadArchive(archfile)
	if err != nil {
		fmt.Printf("Failed to upload archive: %s\n", err)
		os.Exit(3)
	}

	// Show results
	fmt.Printf("Operation completed successfully\n")
	os.Exit(0)
}
