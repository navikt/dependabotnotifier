package github

import (
	"dependabotshamer/pkg/httputils"
	"encoding/json"
	"fmt"
	"github.com/tomnomnom/linkheader"
	"io"
	"net/http"
)

const baseUrl = "https://api.github.com"

type RestRepo struct {
	Name     string    `json:"name"`
	Owner    RepoOwner `json:"owner"`
	Archived bool      `json:"archived"`
}

type RepoOwner struct {
	Name string `json:"login"`
}

type GQLResponse struct {
	Data GQLData `json:"data"`
}

type GQLData struct {
	Repository GQLRepository `json:"repository"`
}

type GQLRepository struct {
	HasVulnerabilityAlertsEnabled bool `json:"hasVulnerabilityAlertsEnabled"`
}

func (res GQLResponse) HasDependabotAlertsEnabled() bool {
	return res.Data.Repository.HasVulnerabilityAlertsEnabled
}

func ReposWithDependabotAlertsDisabled(org, authToken string) ([]RestRepo, error) {
	allRepos, err := allReposFor(org, authToken)
	if err != nil {
		return nil, err
	}
	filteredRepos := []RestRepo{}
	for _, repo := range allRepos {
		if repo.Archived {
			continue
		}
		hasAlerts, err := hasDependabotAlerts(repo.Owner.Name, repo.Name, authToken)
		if err != nil {
			return nil, err
		}
		if !hasAlerts {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

func allReposFor(org, authToken string) ([]RestRepo, error) {
	url := fmt.Sprintf("%s/orgs/%s/repos?per_page=100", baseUrl, org)
	var allRepos []RestRepo
	for url != "" {
		fmt.Printf("Retrieving repos from %s\n", url)
		var reposPart []RestRepo
		extraHeaders := http.Header{
			"Accept":        {"application/vnd.github.v3+json"},
			"User-Agent":    {"NAV IT McBotFace"},
			"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
		}
		res, err := httputils.GetRequest(url, authToken, extraHeaders)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("got a %d from GitHub", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &reposPart)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, reposPart...)
		linkHeader := res.Header.Get("Link")
		url = nextUrl(linkHeader)
	}
	return allRepos, nil
}

func hasDependabotAlerts(owner, repo, authToken string) (bool, error) {
	fmt.Printf("Checking Dependabot alert status for %s/%s\n", owner, repo)
	queryStr := fmt.Sprintf(`query { repository(name: \"%s\", owner: \"%s\") { name hasVulnerabilityAlertsEnabled } }"`, repo, owner)
	reqBody := fmt.Sprintf(`{"query": "%s"}`, queryStr)
	u := fmt.Sprintf("%s/graphql", baseUrl)
	extraHeaders := http.Header{
		"User-Agent":    {"NAV IT McBotFace"},
		"Content-Type":  {"application/json"},
		"Accept":        {"application/vnd.github.v4.idl"},
		"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
	}
	resBody, err := httputils.GQLRequest(u, reqBody, authToken, extraHeaders)
	if err != nil {
		return false, err
	}
	var gqlResponse GQLResponse
	err = json.Unmarshal(resBody, &gqlResponse)
	if err != nil {
		return false, err
	}
	return gqlResponse.HasDependabotAlertsEnabled(), nil
}

func nextUrl(linkHeader string) string {
	links := linkheader.Parse(linkHeader)
	for _, link := range links {
		if link.Rel == "next" {
			return link.URL
		}
	}
	return ""
}
