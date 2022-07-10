package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	ver = "v0.1.0"
)

func main() {
	// Setup command-line arguments
	var csvName, org string
	var version, help, v, h bool
	flag.StringVar(&csvName, "csv", "Findings-example.csv", "Provide the name of the CSV to create")
	flag.StringVar(&org, "org", "", "Provide the name of the Github organization to report on")
	flag.BoolVar(&version, "version", false, "Print the version and exit")
	flag.BoolVar(&v, "v", false, "Print the version and exit")
	flag.BoolVar(&help, "help", false, "Print the help message and exit")
	flag.BoolVar(&h, "h", false, "Print the help message and exit")
	flag.Parse()

	// Handle command-line arguments
	handleStarndarArgs(h, help, v, version)

	// Check required arguments
	requiredArgs(csvName, org)

	// Setup an API client to talk to Github's API
	gh := ghAPIClient{}
	err := setupClient(&gh, org)
	if err != nil {
		fmt.Printf("Error setting up the API client was %+v\n", err)
		os.Exit(1)
	}

	// Create a CSV of Github org information
	err = generateGhCSV(&gh)
	if err != nil {
		fmt.Printf("Error occured while generating CSV\n%+v\n", err)
		os.Exit(1)
	}

}
