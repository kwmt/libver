package main

import (
	"encoding/json"
	//"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	httpClient httpRunner
	user       string
	password   string
}

type httpRunner interface {
	Do(*http.Request) (*http.Response, error)
}

type ResponseMavenPackageSearch struct {
	Name           string   `json:"name,omitempty"`
	Repository     string   `json:"repo,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	Description    string   `json:"desc,omitempty"`
	SystemIDs      []string `json:"system_ids,omitempty"`
	Versions       []string `json:"versions,omitempty"`
	LatestVesrsion string   `json:"latest_version,omitempty"`
}

func NewClient(user, apiKey string) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: time.Duration(30) * time.Second},
		user:       user,
		password:   apiKey,
	}
	return c
}

// Search maven package
func (c *Client) SearchMavenPackage(groupID string, artifactID string) (*[]ResponseMavenPackageSearch, error) {
	values := url.Values{}
	values.Set("g", groupID)
	values.Set("a", artifactID)

	resp, err := c.request("GET", "https://api.bintray.com/search/packages/maven?"+values.Encode(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "request error")
	}
	defer resp.Body.Close()

	return parse(resp)
}

func (c *Client) request(httpMethod string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "response error")
	}

	return resp, nil
}

func parse(resp *http.Response) (*[]ResponseMavenPackageSearch, error) {
	var res []ResponseMavenPackageSearch
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	//fmt.Println(string(b))
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, errors.Wrap(err, "json unmarshal error")
	}
	return &res, nil
}
