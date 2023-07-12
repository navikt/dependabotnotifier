package teams

import (
	"dependabotnotifier/pkg/httputils"
	"encoding/json"
	"fmt"
	"net/http"
)

type GQLResponse struct {
	Data GQLResponseData `json:"data"`
}

type GQLResponseData struct {
	Teams []Team `json:"teamsWithPermissionInGitHubRepo"`
}

type Team struct {
	Slug         string `json:"slug"`
	SlackChannel string `json:"slackChannel"`
}

func AdminsFor(repo, authToken string) ([]Team, error) {
	queryStr := fmt.Sprintf(`query { teamsWithPermissionInGitHubRepo(repoName: \"%s\", permissionName: \"admin\") { slug slackChannel } }`, repo)
	reqBody := fmt.Sprintf(`{"query": "%s"}`, queryStr)
	extraHeaders := http.Header{
		"User-Agent":    {"NAV IT McBotFace"},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
	}
	resBody, err := httputils.GQLRequest("https://teams.nav.cloud.nais.io/query", reqBody, extraHeaders)
	if err != nil {
		return nil, err
	}
	var gqlResponse GQLResponse
	err = json.Unmarshal(resBody, &gqlResponse)
	if err != nil {
		return nil, err
	}
	return gqlResponse.Data.Teams, nil
}
