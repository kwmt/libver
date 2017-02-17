package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// HTTP Client with user and password
type Client struct {
	httpClient httpRunner
	user       string
	password   string
}

type httpRunner interface {
	Do(*http.Request) (*http.Response, error)
}

// See https://bintray.com/docs/api/#_maven_package_search
type ResponseMavenPackageSearch struct {
	MavenPackage
}

// See https://bintray.com/docs/api/#_maven_package_search
type MavenPackage struct {
	Name           string   `json:"name,omitempty"`
	Repository     string   `json:"repo,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	Description    string   `json:"desc,omitempty"`
	SystemIDs      []string `json:"system_ids,omitempty"`
	Versions       []string `json:"versions,omitempty"`
	LatestVesrsion string   `json:"latest_version,omitempty"`
}

// Create new *Client instance.
func NewClient(user, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: time.Duration(30) * time.Second},
		user:       user,
		password:   apiKey,
	}
}

// Search maven package
// groupID: ex, com.google.code.gson
// artifactID: ex, gson
// see https://bintray.com/docs/api/#_maven_package_search
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

	//	printCurl(req)

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
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, errors.Wrap(err, "json unmarshal error")
	}
	return &res, nil
}

func printCurl(req *http.Request) {
	fmt.Println("curl  -X " + req.Method + " \\\n -d '" + toString(req.Body) + "' \\\n" + " -H 'Content-Type:" + req.Header["Content-Type"][0] + "'" + " -H 'Authorization:" + req.Header["Authorization"][0] + "'" + " \\\n '" + req.URL.String() + "'")
}

func toString(rc io.ReadCloser) string {
	if rc == nil {
		return ""
	}
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return ""
	}
	return string(b)
}
