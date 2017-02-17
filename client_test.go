package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type testClient struct {
	req  *http.Request
	resp *http.Response
	err  error
}

func (c *testClient) Do(req *http.Request) (*http.Response, error) {
	c.req = req
	return c.resp, c.err
}

func TestSearchMavenPackage(t *testing.T) {
	const dummyKey = "dummyKey"

	testResponseJSON := []byte(`[{
    "name": "test-package",
    "repo": "jcenter",
    "owner": "bintray",
    "desc": "This package....",
    "system_ids": [
      "groupid:artifactid"
    ],
    "versions": [
        "1.0",
        "2.0"
    ],
    "latest_version": "2.0"
}]`)

	tc := &testClient{
		resp: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(testResponseJSON)),
		},
		err: nil,
	}

	type userPass struct {
		user string
		pass string
	}
	up := userPass{"user", "pass"}

	c := NewClient(up.user, up.pass)
	c.httpClient = tc

	actual, err := c.SearchMavenPackage("dummyGroupID", "dummyArtifactID")

	if username, password, ok := tc.req.BasicAuth(); !ok {
		t.Errorf("actual username: %s password:%s,\nexpect username: %s password:%s", username, password, up.user, up.pass)
	}

	if err != nil {
		t.Error(err)
	}

	var expect []ResponseMavenPackageSearch
	err = json.Unmarshal(testResponseJSON, &expect)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*actual, expect) {
		t.Errorf("got %v, expect %v", *actual, expect)
	}

}
