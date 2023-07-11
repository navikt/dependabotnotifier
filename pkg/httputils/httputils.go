package httputils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GQLRequest(rawUrl, body string, extraHeaders http.Header) ([]byte, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header = extraHeaders
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a %d from %s: %v", res.StatusCode, rawUrl, res)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

func GetRequest(rawUrl string, extraHeaders http.Header) (*http.Response, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header = extraHeaders
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a %d from %s: %v", res.StatusCode, rawUrl, res)
	}
	return res, nil
}
