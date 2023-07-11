package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var client = http.Client{}

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
	res, err := gqlRequest("https://teams.nav.cloud.nais.io/query", reqBody, authToken)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var gqlReaponse GQLResponse
	err = json.Unmarshal(resBody, &gqlReaponse)
	if err != nil {
		return nil, err
	}
	return gqlReaponse.Data.Teams, nil
}

func gqlRequest(rawUrl, body, authToken string) (*http.Response, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"User-Agent":    {"NAV IT McBotFace"},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a %d from %s: %v", res.StatusCode, rawUrl, res)
	}
	return res, nil
}
