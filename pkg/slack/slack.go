package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func SendMessage(channel, heading, text, authToken string) error {
	toSend := Message{
		Channel: channel,
		Blocks: []Block{
			{
				Type: "section",
				Text: Text{
					Type: "mrkdwn",
					Text: fmt.Sprintf("%s", heading),
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
	serialized, err := json.Marshal(toSend)
	extraHeaders := http.Header{
		"User-Agent":    {"NAV IT McBotFace"},
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
	}
	if err != nil {
		return err
	}
	err = postRequest(baseUrl, serialized, extraHeaders)
	if err != nil {
		return err
	}
	return nil
}

func postRequest(rawUrl string, body []byte, extraHeaders http.Header) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header = extraHeaders
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got a %d from %s: %v", res.StatusCode, rawUrl, res)
	}
	return nil
}
