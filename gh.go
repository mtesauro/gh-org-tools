package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ghAPIClient struct {
	BaseURL    *url.URL
	FullURL    *url.URL
	HttpClient *http.Client
	Header     string
	Token      string
	Org        string
	Meta       ghMeta
}

// Struct to hold meta data while retrieving paginated data
type ghMeta struct {
	pagination bool
	nextPage   int
	lastPage   int
	linkHeader string
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
	g.Meta.pagination = false
	g.Meta.nextPage = 0
	g.Meta.lastPage = 0
	g.Meta.linkHeader = ""

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

// resetMeta takes a pointer to ghAPIClient and resets the
// Meta values to their defaults
func resetMeta(g *ghAPIClient) {
	g.Meta.pagination = false
	g.Meta.nextPage = 0
	g.Meta.lastPage = 0
	g.Meta.linkHeader = ""
}

// Takes a pointer to ghAPIClient and generates a CSV of the
// appropriate Github organization
func generateGhCSV(g *ghAPIClient) error {
	// Get info on the provided Github org
	oInfo := []ghOrgInfo{}
	err := getOrgInfo(g, &oInfo)
	if err != nil {
		return err
	}
	resetMeta(g)

	// Ensure we have a single results back for the org
	if len(oInfo) > 1 {
		return errors.New("Multiple Github organizations returned, which makes no sense. Exiting...")
	}

	// Use the Github org info to retrieve a list of repos for that GH org
	oRepos := ghRepoInfo{}
	err = getOrgRepos(g, &oRepos)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem preparing Request was: %v", err))
	}
	resetMeta(g)

	// For each repo in the GH org, get a list of collaborators to pull out those with admin roles
	for k := range oRepos {
		//fmt.Printf("FullName data for %+v\n", oRepos[k].FullName)
		//fmt.Printf("Name data for %+v\n", oRepos[k].Name)
		//fmt.Printf("CollaboratorsURL data for %+v\n", oRepos[k].CollaboratorsURL)
		//fmt.Printf("Type for CollaboratorsURL %T\n", oRepos[k].CollaboratorsURL)
		//
		rCollab := make(map[string]ghCollaborators)
		//current := oRepos[k].CollaboratorsURL
		//name := oRepos[k].Name
		rCollab[oRepos[k].Name], err = getCollabs(oRepos[k].CollaboratorsURL)
	}

	// For each collaborator with an admin role, determine their name (human one vs GH login name aka Github username)

	// Generate the CSV and write it out.
	return nil
}

// pagedResults takes in a http Reponse struct and looks for a "Link" header
// which Github API uses to determine if the response has been paginated.
// Returns true if there are multiple pages of results and the next and last
// page numbers from the Link header
// Example link headers:
// link: <https://api.github.com/organizations/123/repos?page=2>; rel="next", <https://api.github.com/organizations/123/repos?page=4>; rel="last"
// link: <https://api.github.com/organizations/123/repos?page=1>; rel="prev", <https://api.github.com/organizations/123/repos?page=3>; rel="next", <https://api.github.com/organizations/123/repos?page=4>; rel="last", <https://api.github.com/organizations/123/repos?page=1>; rel="first"
func pagedResults(m *ghMeta) error {
	// Look for Link header, return early if not found
	if len(m.linkHeader) == 0 {
		// Link header not found
		m.pagination = false
		m.nextPage = 0
		m.lastPage = 0
		return nil
	}

	// Split the Link header into the major parts (next, last, prev, ...)
	rawVals := strings.Split(m.linkHeader, ",")

	// Pull out the next and last values
	var nu, lu string
	foundNext := false
	for k := range rawVals {
		if strings.Contains(rawVals[k], "next") {
			nu = rawVals[k]
			foundNext = true
		}
		if strings.Contains(rawVals[k], "last") {
			lu = rawVals[k]
		}
	}

	if !foundNext {
		m.pagination = false
		m.nextPage = 0
		m.lastPage = 0
		return nil
	}

	// Get Next URL as a string and convert to a URL
	nextURL, err := onlyURL(nu)
	if err != nil {
		m.pagination = false
		m.nextPage = 0
		m.lastPage = 0
		return nil
	}

	// Get Last URL as a string and convert to a URL
	lastURL, err := onlyURL(lu)
	if err != nil {
		m.pagination = false
		m.nextPage = 0
		m.lastPage = 0
		return err
	}

	// Parse out the page values
	np, err := parsePage(nextURL)
	if err != nil {
		m.pagination = false
		m.nextPage = 0
		m.lastPage = 0
		return errors.New("Unable to parse page query parameter used in pagination")
	}
	lp, err := parsePage(lastURL)
	if err != nil {
		m.pagination = false
		m.nextPage = 0
		m.lastPage = 0
		return errors.New("Unable to parse page query parameter used in pagination")
	}

	// No issues, set the meta values
	m.pagination = true
	m.nextPage = np
	m.lastPage = lp
	return nil
}

// onlyURL takes a raw string from the Link header used by the Github API
// and returns a pointer to url.URL and an error if unable to successfully
// parse the provided string. The string is create by spitting the link
// header by ',' and example follows:
// <https://api.github.com/organizations/123/repos?page=3>; rel="next"
func onlyURL(raw string) (*url.URL, error) {
	// Split into url and rel attribute
	t := strings.Split(raw, ";")

	// Pull out just the URL
	u := strings.Trim(strings.Replace(strings.Replace(t[0], ">", "", -1), "<", "", -1), " ")

	// Parse into url.URL
	final, err := url.Parse(u)
	if err != nil {
		return final, err
	}

	return final, nil
}

// pargePage takes a url.URL and returns the value for the page query value
// in the URL. This is meant for Github API paginated results where the
// link header provides next and last URLs with query strings with values
// like '?page=3'. In that example, an int of 3 would be returned or an
// error if there was problems pasring the URL
func parsePage(u *url.URL) (int, error) {
	pRaw := u.Query()

	// Ensure page query parameter exists
	if len(pRaw["page"]) == 0 {
		return 0, errors.New("Page query parameter missing but link header present. Exiting...")
	}

	// Ensure there's only 1 page paramter
	if len(pRaw["page"]) != 1 {
		return 0, errors.New("Multiple values for page query parameter which can't be right. Exiting...")
	}

	// Convert to page parameter to an int
	p, err := strconv.Atoi(pRaw["page"][0])
	if err != nil {
		return 0, errors.New("Unable to convert page query parameter to an int.  Exiting...")
	}

	return p, nil
}

// getOrgInfo takes pointers to ghAPIClient and ghOrgInfo and retrieves the Github org's repos
// (based on the Org field of ghAPIClient) to fill the ghOrgInfo struct
func getOrgInfo(g *ghAPIClient, oInfo *[]ghOrgInfo) error {
	// Add the URI for the Get organization call
	// see https://docs.github.com/en/rest/orgs/orgs#get-an-organization
	addURI(g, "/orgs/"+g.Org)

	// Setup meta data struct and temp data struct
	tempOrg := ghOrgInfo{}

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

	// Check the response code
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("API response code for Org Info was: %v", resp.StatusCode))
	}

	// Save the link header
	g.Meta.linkHeader = resp.Header.Get("link")

	// Unmarshall data to struct
	err = json.Unmarshal(body, &tempOrg)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem unmarshalling JSON was: %v", err))
	}

	// Append the data collected so far
	*oInfo = append(*oInfo, tempOrg)

	// Check for pagination
	err = pagedResults(&g.Meta)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem determining pagination was: %v", err))
	}

	if g.Meta.pagination {
		// Construct the new URL for the next page of results
		addURI(g, "/orgs/"+g.Org+"?page="+strconv.Itoa(g.Meta.nextPage))
		if g.Meta.nextPage <= g.Meta.lastPage {
			// Increament the next page
			err := getOrgInfo(g, oInfo)
			if err != nil {
				return errors.New(fmt.Sprintf("Problem requesting addition data pages was: %v", err))
			}
		}
	}

	return nil
}

// getOrgRepos takes pointers to ghAPIClient and ghRepoInfo plus an int and retrieves all the orgs repos
// (based on the ReposURL field of ghOrgInfo) to fill the ghRepoInfo struct.
func getOrgRepos(g *ghAPIClient, oRepos *ghRepoInfo) error {
	// Add the URI for the Get organization call
	// see https://docs.github.com/en/rest/orgs/orgs#get-an-organization
	if !g.Meta.pagination {
		addURI(g, "/orgs/"+g.Org+"/repos")
	}

	// Setup meta data struct and temp data struct
	tempRepos := ghRepoInfo{}

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

	// Check the response code
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("API response code for Org Repos was: %v", resp.StatusCode))
	}

	// Save the link header
	g.Meta.linkHeader = resp.Header.Get("link")

	// Unmarshall data to struct
	err = json.Unmarshal(body, &tempRepos)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem unmarshalling JSON was: %v", err))
	}

	// Append the data collected so far
	for k := range tempRepos {
		*oRepos = append(*oRepos, tempRepos[k])
	}

	// Check for pagination
	err = pagedResults(&g.Meta)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem determining pagination was: %v", err))
	}

	if g.Meta.pagination {
		// Construct the new URL for the next page of results
		addURI(g, "/orgs/"+g.Org+"/repos?page="+strconv.Itoa(g.Meta.nextPage))
		if g.Meta.nextPage <= g.Meta.lastPage {
			// Increament the next page
			err := getOrgRepos(g, oRepos)
			if err != nil {
				return errors.New(fmt.Sprintf("Problem requesting addition data pages was: %v", err))
			}
		}
	}

	return nil
}

func getCollabs(raw string) (ghCollaborators, error) {
	temp := make(ghCollaborators, 2)
	link := strings.ReplaceAll(raw, "{/collaborator}", "")
	fmt.Printf("Fixed link is %+v\n", link)

	return temp, nil
}
