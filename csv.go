package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

//
func writeCSV(f string, repos *ghRepoInfo, adm map[string]ghCollaborators, lu map[string]ghNameDetail) error {
	// Create the CSV file
	fi, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fi.Close()

	// Setup a new CSV writer
	csvFile := csv.NewWriter(fi)
	defer csvFile.Flush()

	// Create the header row
	var header []string
	// Do header in 'long' for to make turning items off and on easy
	header = append(header,
		"Full Name",         // e.g. org/repo-name
		"Name",              // e.g. repo-name
		"Short Description", // Full description trimmed down to 46 characters
		//"Repo Type",           // 'public', 'private', 'internal' or 'public, fork'
		"Private",     // true or false
		"Fork",        // true or false
		"Visibility",  // public or private
		"Last Update", // e.g. 2022-06-13T07:59:05Z
		"Repo Admins", // List of all the admins
	)
	err = csvFile.Write(header)
	if err != nil {
		return err
	}

	// Add the collected details to the CSV
	for _, v := range *repos {
		// Slice of string for each CSV line
		var line []string

		// Setup some values first
		desc := v.Description
		if len(v.Description) > 46 {
			desc = v.Description[0:45]
		}
		admList := listAdmins(v.Name, adm, lu)

		// Create a CSV line
		line = append(line,
			v.FullName,
			v.Name,
			desc,
			strconv.FormatBool(v.Private),
			strconv.FormatBool(v.Fork),
			v.Visibility,
			v.UpdatedAt.Format(time.RFC3339),
			strings.TrimRight(admList, ","),
		)

		// Write out the current line
		err := csvFile.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}

//
func listAdmins(repo string, adm map[string]ghCollaborators, lu map[string]ghNameDetail) string {
	// Find the admins for the current repo
	var list string
	for _, v := range adm[repo] {
		list += v.Login + checkDetails(lu[v.Login].Name, lu[v.Login].Email)
	}

	return list
}

//
func checkDetails(n string, e string) string {
	// " (" + n + " - " + e + ") "
	// Check name
	gotName := false
	if len(n) > 0 {
		gotName = true
	}

	gotEmail := false
	if len(e) > 0 {
		gotName = true
	}

	// Return appropriate details
	if gotName && gotEmail {
		return " (" + n + " - " + e + "), "
	}
	if gotName && !gotEmail {
		return " (" + n + "), "
	}
	if !gotName && gotEmail {
		return " (" + e + "), "
	}

	return ", "
}
