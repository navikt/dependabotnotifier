package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type Response struct {
	Ok       bool   `json:"ok"`
	ErrorMsg string `json:"error"`
}

func SendMessage(channel, heading, text, authToken string) Response {
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
		return NewError(err)
	}
	return postRequest(baseUrl, serialized, extraHeaders)
}

func postRequest(rawUrl string, reqBody []byte, extraHeaders http.Header) Response {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return NewError(err)
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return NewError(err)
	}
	req.Header = extraHeaders
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return NewError(err)
	}
	if res.StatusCode != http.StatusOK {
		return NewError(fmt.Errorf("got a %d from %s: %v", res.StatusCode, rawUrl, res))
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return NewError(err)
	}
	var parsedResponse Response
	err = json.Unmarshal(resBody, &parsedResponse)
	if err != nil {
		return NewError(err)
	}
	return parsedResponse
}

func NewError(err error) Response {
	return Response{
		Ok:       false,
		ErrorMsg: err.Error(),
	}
}
