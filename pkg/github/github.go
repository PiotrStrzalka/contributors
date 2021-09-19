package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

type Client struct {
	token   string
	client  http.Client
	baseUrl string
}

var tokenRE = regexp.MustCompile(`^[0-9a-zA-Z_]{40}$`)

func NewClient(baseURL, token string) (*Client, error) {

	if token == "" {
		return nil, errors.New("token is required")
	}

	if !tokenRE.MatchString(token) {
		return nil, errors.New("token is not correct")
	}

	return &Client{
		token:   token,
		client:  http.Client{Timeout: 5 * time.Second},
		baseUrl: baseURL,
	}, nil
}

func (c *Client) ContributorList(repo string) ([]Contributor, error) {
	// Create a request for the contributors api endpoint.
	endpoint := c.baseUrl + "/repos/" + repo + "/contributors"
	log.Printf("Endpoint: %s", endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Add the access token in the "Authorization" header.
	// The value should be in the form "token 000aa0a0..."
	req.Header.Set("Authorization", "token "+c.token)

	// Create an http.Client and make the request.
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Defer closing the response body.
	defer res.Body.Close()

	// Ensure we get a 200 OK status back.
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with a %d %s", res.StatusCode, res.Status)
	}

	// Decode the results into a []contributor.
	var con []Contributor
	if err := json.NewDecoder(res.Body).Decode(&con); err != nil {
		return nil, err
	}
	return con, nil
}
