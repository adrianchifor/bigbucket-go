package bigbucket

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getGcpAuthToken(receivingServiceUrl string) (string, error) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	params := map[string]string{
		"audience": receivingServiceUrl,
	}
	httpUrl, err := constructUrl("http://metadata", "/computeMetadata/v1/instance/service-accounts/default/identity", params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func jwtExpired(jwt string) (bool, error) {
	if jwt == "" {
		return true, nil
	}

	jwtArr := strings.Split(jwt, ".")
	seg := jwtArr[1]
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	segPayload, err := base64.URLEncoding.DecodeString(seg)
	if err != nil {
		return true, err
	}

	var data map[string]int
	json.NewDecoder(bytes.NewBuffer(segPayload)).Decode(&data)
	if _, exists := data["exp"]; !exists {
		return true, errors.New("bigbucket -- Failed to decode JWT JSON content, 'exp' not found")
	}
	expiry := data["exp"]
	// Add 30 seconds to allow time for requests
	now := int(time.Now().Unix()) + 30

	return expiry <= now, nil
}
