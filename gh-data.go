package main

import "time"

// Response from Github API for info on an organization
// e.g. https://api.github.com/orgs/[org name]
// see https://docs.github.com/en/rest/orgs/orgs#get-an-organization
type ghOrgInfo struct {
	Login                                string    `json:"login"`
	ID                                   int       `json:"id"`
	NodeID                               string    `json:"node_id"`
	URL                                  string    `json:"url"`
	ReposURL                             string    `json:"repos_url"`
	EventsURL                            string    `json:"events_url"`
	HooksURL                             string    `json:"hooks_url"`
	IssuesURL                            string    `json:"issues_url"`
	MembersURL                           string    `json:"members_url"`
	PublicMembersURL                     string    `json:"public_members_url"`
	AvatarURL                            string    `json:"avatar_url"`
	Description                          string    `json:"description"`
	Name                                 string    `json:"name"`
	Company                              string    `json:"company"`
	Blog                                 string    `json:"blog"`
	Location                             string    `json:"location"`
	Email                                string    `json:"email"`
	TwitterUsername                      string    `json:"twitter_username"`
	IsVerified                           bool      `json:"is_verified"`
	HasOrganizationProjects              bool      `json:"has_organization_projects"`
	HasRepositoryProjects                bool      `json:"has_repository_projects"`
	PublicRepos                          int       `json:"public_repos"`
	PublicGists                          int       `json:"public_gists"`
	Followers                            int       `json:"followers"`
	Following                            int       `json:"following"`
	HTMLURL                              string    `json:"html_url"`
	CreatedAt                            time.Time `json:"created_at"`
	UpdatedAt                            time.Time `json:"updated_at"`
	Type                                 string    `json:"type"`
	TotalPrivateRepos                    int       `json:"total_private_repos"`
	OwnedPrivateRepos                    int       `json:"owned_private_repos"`
	PrivateGists                         int       `json:"private_gists"`
	DiskUsage                            int       `json:"disk_usage"`
	Collaborators                        int       `json:"collaborators"`
	BillingEmail                         string    `json:"billing_email"`
	DefaultRepositoryPermission          string    `json:"default_repository_permission"`
	MembersCanCreateRepositories         bool      `json:"members_can_create_repositories"`
	TwoFactorRequirementEnabled          bool      `json:"two_factor_requirement_enabled"`
	MembersAllowedRepositoryCreationType string    `json:"members_allowed_repository_creation_type"`
	MembersCanCreatePublicRepositories   bool      `json:"members_can_create_public_repositories"`
	MembersCanCreatePrivateRepositories  bool      `json:"members_can_create_private_repositories"`
	MembersCanCreateInternalRepositories bool      `json:"members_can_create_internal_repositories"`
	MembersCanCreatePages                bool      `json:"members_can_create_pages"`
	MembersCanForkPrivateRepositories    bool      `json:"members_can_fork_private_repositories"`
	MembersCanCreatePublicPages          bool      `json:"members_can_create_public_pages"`
	MembersCanCreatePrivatePages         bool      `json:"members_can_create_private_pages"`
	WebCommitSignoffRequired             bool      `json:"web_commit_signoff_required"`
	Plan                                 struct {
		Name         string `json:"name"`
		Space        int    `json:"space"`
		PrivateRepos int    `json:"private_repos"`
		FilledSeats  int    `json:"filled_seats"`
		Seats        int    `json:"seats"`
	} `json:"plan"`
}

// Response from Github API for info on an organization's repositories
// e.g. https://api.github.com/orgs/[org name]/repos
// see https://docs.github.com/en/rest/repos/repos#list-organization-repositories
type ghRepoInfo []struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HTMLURL          string    `json:"html_url"`
	Description      string    `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	ForksURL         string    `json:"forks_url"`
	KeysURL          string    `json:"keys_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	TeamsURL         string    `json:"teams_url"`
	HooksURL         string    `json:"hooks_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	EventsURL        string    `json:"events_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BranchesURL      string    `json:"branches_url"`
	TagsURL          string    `json:"tags_url"`
	BlobsURL         string    `json:"blobs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	TreesURL         string    `json:"trees_url"`
	StatusesURL      string    `json:"statuses_url"`
	LanguagesURL     string    `json:"languages_url"`
	StargazersURL    string    `json:"stargazers_url"`
	ContributorsURL  string    `json:"contributors_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommitsURL       string    `json:"commits_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	CommentsURL      string    `json:"comments_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	ContentsURL      string    `json:"contents_url"`
	CompareURL       string    `json:"compare_url"`
	MergesURL        string    `json:"merges_url"`
	ArchiveURL       string    `json:"archive_url"`
	DownloadsURL     string    `json:"downloads_url"`
	IssuesURL        string    `json:"issues_url"`
	PullsURL         string    `json:"pulls_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	LabelsURL        string    `json:"labels_url"`
	ReleasesURL      string    `json:"releases_url"`
	DeploymentsURL   string    `json:"deployments_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         time.Time `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	SvnURL           string    `json:"svn_url"`
	Homepage         string    `json:"homepage"`
	Size             int       `json:"size"`
	StargazersCount  int       `json:"stargazers_count"`
	WatchersCount    int       `json:"watchers_count"`
	Language         string    `json:"language"`
	HasIssues        bool      `json:"has_issues"`
	HasProjects      bool      `json:"has_projects"`
	HasDownloads     bool      `json:"has_downloads"`
	HasWiki          bool      `json:"has_wiki"`
	HasPages         bool      `json:"has_pages"`
	ForksCount       int       `json:"forks_count"`
	MirrorURL        string    `json:"mirror_url"`
	Archived         bool      `json:"archived"`
	Disabled         bool      `json:"disabled"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	AllowForking             bool     `json:"allow_forking"`
	IsTemplate               bool     `json:"is_template"`
	WebCommitSignoffRequired bool     `json:"web_commit_signoff_required"`
	Topics                   []string `json:"topics"`
	Visibility               string   `json:"visibility"`
	Forks                    int      `json:"forks"`
	OpenIssues               int      `json:"open_issues"`
	Watchers                 int      `json:"watchers"`
	DefaultBranch            string   `json:"default_branch"`
	Permissions              struct {
		Admin    bool `json:"admin"`
		Maintain bool `json:"maintain"`
		Push     bool `json:"push"`
		Triage   bool `json:"triage"`
		Pull     bool `json:"pull"`
	} `json:"permissions"`
}

// Response from Github API for info on a repo's collaborators
// e.g. https://api.github.com/repos/[org name]/[repo name//collaborators
// see https://docs.github.com/en/rest/collaborators/collaborators#list-repository-collaborators
type ghCollaborators []struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Permissions       struct {
		Admin    bool `json:"admin"`
		Maintain bool `json:"maintain"`
		Push     bool `json:"push"`
		Triage   bool `json:"triage"`
		Pull     bool `json:"pull"`
	} `json:"permissions"`
	RoleName string `json:"role_name"`
}

