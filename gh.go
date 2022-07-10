package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ghAPIClient struct {
	BaseURL    *url.URL
	FullURL    *url.URL
	HttpClient *http.Client
	Header     string
	Token      string
	Org        string
}

// setupClient takes a pointer to ghAPIClient, validates that
// the required environmental variable 'GHTOKEN' exists and, if so,
// creates a ghAPIClient with default values set
func setupClient(g *ghAPIClient, o string) error {
	// Setup the necessary config from the environment
	t, present := os.LookupEnv("GHTOKEN")
	if !present {
		// If GHTOKEN isn't set, error out
		return errors.New("Required environmental variable 'GHTOKEN' not found")
	}

	// Setup base URL
	u, err := url.Parse("https://api.github.com")
	if err != nil {
		fmt.Printf("Error parsing the Github API URL was %+v\n", err)
		return err
	}

	// Setup HttpClient
	c := &http.Client{}

	// Create ghAPIClient based on config values
	g.BaseURL = u
	g.HttpClient = c
	g.Header = "Authorization"
	g.Token = "token " + t
	g.Org = o

	return nil
}

// addURL takes a pointer to a ghAPIClient and a string
// which holds a URI to append on the existing URL in
// ghAPIClient to produce a full URL for an API call
func addURI(g *ghAPIClient, u string) error {
	// Setup values for ghAPIClient URL
	x, err := url.Parse(g.BaseURL.String() + u)
	if err != nil {
		fmt.Printf("Error parsing the Github API URL was %+v\n", err)
		return err
	}

	// Set new Full URL
	g.FullURL = x

	return nil
}

// Set the organization for the ghAPIClient
func setOrg(g *ghAPIClient, o string) {
	g.Org = o
}

func generateGhCSV(g *ghAPIClient) error {
	// Get info on the provided Github org
	oInfo := ghOrgInfo{}
	err := getOrgInfo(g, &oInfo)
	if err != nil {
		return err
	}

	// DEBUG
	fmt.Printf("%+v\n", oInfo.ReposURL)
	os.Exit(0)
	// Use the Github org info to retrieve a list of repos for that GH org
	oRepos := ghRepoInfo{}
	err = getOrgRepos(g, &oRepos)
	if err != nil {
		return err
	}

	// For each repo in the GH org, get a list of collaborators to pull out those with admin roles

	// For each collaborator with an admin role, determine their name (human one vs GH login name aka Github username)

	// Generate the CSV and write it out.
	return nil
}

// getOrgInfo takes pointers to ghAPIClient and ghOrgInfo and retrieves the Github org's repos
// (based on the Org field of ghAPIClient) to fill the ghOrgInfo struct
func getOrgInfo(g *ghAPIClient, oInfo *ghOrgInfo) error {
	// Add the URI for the Get organization call
	// see https://docs.github.com/en/rest/orgs/orgs#get-an-organization
	addURI(g, "/orgs/"+g.Org)

	// Setup meta data struct
	orgMeta := ghOrgInfoMeta{}

	// Orgs are single responses - no pagination
	orgMeta.pagination = false

	// Setup the request
	rawResp := ""
	req, err := http.NewRequest(http.MethodGet, g.FullURL.String(), strings.NewReader(rawResp))
	if err != nil {
		return errors.New(fmt.Sprintf("Problem preparing Request was: %v", err))
	}
	req.Header.Add("content-type", "application/vnd.github+json")
	req.Header.Add(g.Header, g.Token)

	// Send the request
	resp, err := g.HttpClient.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem sending Request was: %v", err))
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem reading response body was: %v", err))
	}

	// Unmarshall data to struct
	err = json.Unmarshal(body, oInfo)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem unmarshalling JSON was: %v", err))
	}

	return nil
}

// getOrgRepos takes pointers to ghAPIClient and ghRepoInfo and retrieves all the orgs repos
// (based on the ReposURL field of ghOrgInfo) to fill the ghRepoInfo struct
func getOrgRepos(g *ghAPIClient, oRepos *ghRepoInfo) error {
	// DON'T FORGET TO CHECK FOR THE LINK header in the response to handle pagination
	return nil
}
