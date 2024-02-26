package teams

import (
	"dependabotnotifier/pkg/httputils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GQLResponse struct {
	Data GQLResponseData `json:"data"`
}

type GQLResponseData struct {
	Teams    []Team   `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Team struct {
	Slug         string `json:"slug"`
	SlackChannel string `json:"slackChannel"`
}

type PageInfo struct {
	TotalCount  int  `json:"totalCount"`
	HasNext     bool `json:"hasNextPage"`
	HasPrevious bool `json:"hasPreviousPage"`
}

func AdminsFor(repo, authToken string) ([]Team, error) {
	offset := 0
	limit := 100
	var allTeams []Team
	done := false
	for done != true {
		response, err := singleQuery(repo, authToken, offset, limit)
		if err != nil {
			return []Team{}, err
		}
		allTeams = append(allTeams, response.Data.Teams...)
		done = !response.Data.PageInfo.HasNext
		offset += response.Data.PageInfo.TotalCount
	}
	return allTeams, nil
}

func singleQuery(repo, authToken string, offset, limit int) (GQLResponse, error) {
	queryStr := fmt.Sprintf(`
query($filter: TeamsFilter) { teams(filter: $filter, offset: $offset, limit: $limit) { slug, slackChannel } }", 
"variables": { 
  "filter": { 
    "github": { 
      "repoName": "%s", 
      "permissionName": "admin"
    }
  }, 
  "offset": %d, 
  "limit": %d 
}`, repo, offset, limit)
	reqBody := fmt.Sprintf(`{"query": "%s"}`, strings.ReplaceAll(queryStr, "\n", " "))
	extraHeaders := http.Header{
		"User-Agent":    {"NAV IT McBotFace"},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
	}
	resBody, err := httputils.GQLRequest("https://apiserver.prod-gcp.nav.cloud.nais.io/query", reqBody, extraHeaders)
	if err != nil {
		return GQLResponse{}, err
	}
	var gqlResponse GQLResponse
	err = json.Unmarshal(resBody, &gqlResponse)
	if err != nil {
		return GQLResponse{}, err
	}
	return gqlResponse, nil
}
