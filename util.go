package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// handleStandardArgs takes 4 booleans: h and help for help argument plus
// v and version for version arguement and produces the correct output -
// either the help message or the programs version
func handleStarndarArgs(h bool, help bool, v bool, version bool) {
	// Print help for errors with command-line arguments
	if len(os.Args) == 1 || // No arguments
		len(flag.Args()) > 0 { // Unsupported or extra arguments
		help = true
		fmt.Println("Error: Incorrect command-line arguments proivded, printing help...")
	}

	// Print help when h/help argument given
	if help || h {
		printHelp()
		os.Exit(0)
	}

	// Check for version command-line argument
	if version || v {
		fmt.Printf("ghorg2csv version %s\n", ver)
		os.Exit(0)
	}

	return
}

// printHelp prints the usage help output for this program
func printHelp() {
	fmt.Println("")
	fmt.Println("Usage of ghorg2csv")
	fmt.Println("")
	fmt.Println("  -csv  string")
	fmt.Println("        REQUIRED - Provide the name of the CSV to create")
	fmt.Println("  -org  string")
	fmt.Println("        REQUIRED - Provide the name of the Github organization")
	fmt.Println("  -help, -h")
	fmt.Println("        Print this help message and exit")
	fmt.Println("  -version, -v")
	fmt.Println("        Print the version and exit, ignoring all other arguments")
	fmt.Println("")
	fmt.Println("  Note: GNU-style arguments like --name are also supported")
	fmt.Println("")
	fmt.Println("  WARNING: The token used to authenticate with the Github API must")
	fmt.Println("  be passed as an environmental variable named 'GHTOKEN'")
	fmt.Println("")
	fmt.Println("  Example:")
	fmt.Println("        $ ghorg2csv  --csv \"org-info.csv\" --org \"my-github-org\"")
	fmt.Println("")

}

// Ensure required arguments are provided by checking that csv name is at least 5 characters
// and ends in '.csv' as well as ensuring that org isn't empty (the default value)
func requiredArgs(c string, o string) {
	// Simple length check of CSV file
	if len(c) < 5 {
		fmt.Println("ERROR: CSV name is too short, smallest possible length is 5 characters e.g. a.csv")
		os.Exit(1)
	}
	// Ensure csv name ends in .csv
	if !strings.HasSuffix(c, ".csv") {
		fmt.Println("ERROR: CSV name should end in '.csv' e.g. my-GH-Org.csv")
		os.Exit(1)
	}

	// Make sure there's a Github org argument provided
	if len(o) == 0 {
		fmt.Println("Please provide a Github org with the -org argument")
		printHelp()
		os.Exit(1)
	}
}
