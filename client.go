package clarifai2

import (
//	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"errors"
	"strings"
	"strconv"
)

const (
	version = "v2"
	rootURL = "https://api.clarifai.com"
)

type Client struct {
	APIKey string
	APIRoot string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey, rootURL}
}

func (client *Client) commonHTTPRequest(jsonBody interface{}, endpoint, verb string, retry bool) ([]byte, error) {
	if jsonBody == nil {
		jsonBody = struct{}{}
	}

	body, err := json.Marshal(jsonBody)

	if err != nil {
		return nil, err
	}

//	fmt.Println(client.BuildURL(endpoint))
//	s := string(body)
//	fmt.Println(s)
	
	req, err := http.NewRequest(verb, client.BuildURL(endpoint), bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Key "+client.APIKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200, 201:
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		return body, err
	case 401:
		if !retry {
			return client.commonHTTPRequest(jsonBody, endpoint, verb, true)
		}
		return nil, errors.New("APIKEY_INVALID")
	case 400:
		return nil, errors.New("ALL_ERROR")
	case 500:
		return nil, errors.New("CLARIFAI_ERROR")
	default:
		return nil, errors.New("UNEXPECTED_STATUS_CODE: "+strconv.Itoa(res.StatusCode))
	}
}

// SetAccessToken will set accessToken to a new value
func (client *Client) setAPIKey(apiKey string) {
	client.APIKey = apiKey
}

func (client *Client) setAPIRoot(root string) {
	client.APIRoot = root
}

// Helper function to build URLs
func (client *Client) BuildURL(endpoint string) string {
	parts := []string{client.APIRoot, version, endpoint}
	return strings.Join(parts, "/")
}
