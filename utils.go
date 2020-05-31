package bigbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func httpRequest(method string, client *Client, path string, params map[string]string, body map[string]string) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(client.timeout) * time.Second,
	}
	httpUrl, err := constructUrl(client.address, path, params)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, httpUrl, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, httpUrl, nil)
		if err != nil {
			return nil, err
		}
	}

	if client.gcpAuth {
		expired, err := jwtExpired(client.jwt)
		if err != nil {
			return nil, err
		}
		if expired {
			client.jwt, err = getGcpAuthToken(client.address)
			if err != nil {
				return nil, err
			}
		}

		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", client.jwt))
	}

	for header, value := range client.headers {
		req.Header.Set(header, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func constructUrl(address string, path string, params map[string]string) (string, error) {
	Url, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	Url.Path += path
	if params != nil {
		parameters := url.Values{}
		for k, v := range params {
			parameters.Add(k, v)
		}
		Url.RawQuery = parameters.Encode()
	}

	return Url.String(), nil
}

func constructError(response *http.Response) error {
	var data map[string]string
	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("bigbucket -- %d: %s", response.StatusCode, data["error"]))
}
