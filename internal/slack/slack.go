package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/navikt/dependabotnotifier/internal/httputils"

	"github.com/sirupsen/logrus"
)

const baseUrl = "https://slack.com/api/chat.postMessage"

type Message struct {
	Channel string  `json:"channel"`
	Blocks  []Block `json:"blocks"`
}

type Block struct {
	Type string      `json:"type"`
	Text interface{} `json:"text,omitempty"`
}

type Text struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type Response struct {
	Ok       bool   `json:"ok"`
	ErrorMsg string `json:"error"`
}

type Client struct {
	apiToken string
	log      logrus.FieldLogger
}

func NewClient(apiToken string, log logrus.FieldLogger) *Client {
	return &Client{
		apiToken: apiToken,
		log:      log,
	}
}

func (c *Client) SendMessage(ctx context.Context, channelName, teamSlug, repoName string) error {
	heading := fmt.Sprintf(`:wave: Hei, %s :github2:`, teamSlug)
	text := fmt.Sprintf(`Dere har knyttet GitHub-repoet <https://github.com/%[1]s|%[1]s> opp til teamet deres via <https://console.nav.cloud.nais.io/team/%[2]s/repositories|Console>. Dette repoet har ikke Dependabot alerts aktivert. Dependabot hjelper deg å oppdage biblioteker med kjente sårbarheter i appene dine. Du kan sjekke status og enable Dependabot <https://github.com/%[1]s/security|her>. Hvis repoet ikke er i bruk, vurder å arkivere det. Det kan gjøres nederst på <https://github.com/%[1]s/settings|denne siden>.`, repoName, teamSlug)

	toSend := Message{
		Channel: channelName,
		Blocks: []Block{
			{
				Type: "header",
				Text: Text{
					Type: "plain_text",
					Text: heading,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: Text{
					Type: "mrkdwn",
					Text: text,
				},
			},
		},
	}

	serialized := new(bytes.Buffer)
	if err := json.NewEncoder(serialized).Encode(toSend); err != nil {
		return err
	}

	return c.postRequest(ctx, baseUrl, serialized)
}

func (c *Client) postRequest(ctx context.Context, rawUrl string, requestBody *bytes.Buffer) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), requestBody)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"User-Agent":    {httputils.UserAgent},
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + c.apiToken},
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status code %d from %q: %v", res.StatusCode, rawUrl, res)
	}
	var parsedResponse Response
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		return err
	}
	if !parsedResponse.Ok {
		return fmt.Errorf("failed to send message: %s", parsedResponse.ErrorMsg)
	}

	return nil
}
