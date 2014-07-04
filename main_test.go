package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"testing"
)

func TestParseFlagMappingsValidData(t *testing.T) {
	input := []string{"/api/*. -> example.com", "/admin*. -> a.example.com"}
	mappings, err := parseFlagMappings(input)

	assert.Equal(t, 2, len(mappings), "2 slice mappings should end up as 2 PathMapping structs")
	assert.Equal(t, nil, err, "no error should be returned when mappings are parsed successfully")
}

func TestParseFlagMappingsWhiteSpaceIgnored(t *testing.T) {
	input := []string{"      /api/.*             ->           example.com          "}
	mappings, _ := parseFlagMappings(input)

	assert.Equal(t, mappings[0].HostPort, "example.com", "whitespace should be trimmed from the host name")
	assert.Equal(t, mappings[0].Regex.String(), "/api/.*", "whitespace should be trimmed from the regex")
}

func TestParseFlagMappingsInvalidRegex(t *testing.T) {
	input := []string{"[[[[[ -> example.com"}
	_, err := parseFlagMappings(input)

	assert.Equal(t, "Invalid regex: [[[[[", err.Error(), "an error should be returned when regex is invalid")
}

func TestParseFlagMappingsMissingArrow(t *testing.T) {
	input := []string{"/api/*. = example.com"} //intentionally using = instead of -> here as invalid input
	_, err := parseFlagMappings(input)

	assert.Equal(t, "Invalid mapping syntax for: /api/*. = example.com", err.Error(), "an error should be returned when mapping syntax is invalid")
}

func TestParseFlagMappingsArrowInRegex(t *testing.T) {
	input := []string{"/api->/*. -> example.com"}
	mappings, _ := parseFlagMappings(input)

	assert.NotEmpty(t, mappings, "Arrows should be allowed in regex and parse normally")
}

func TestDirectorMappingMatch(t *testing.T) {
	req := makeRequest("/api/users/19/details")
	proxy := makeBasicProxy()
	proxy.Director(req)

	assert.Equal(t, "u.example.com:8080", req.URL.Host, "successful match should change the request.URL.Host")
}

func TestDirectorMappingNoMatch(t *testing.T) {
	req := makeRequest("/no/match")
	proxy := makeBasicProxy()
	proxy.Director(req)

	assert.Equal(t, "v.example.com:8080", req.URL.Host, "last mapping should be used when no other mapping matches")
}

func makeRequest(path string) *http.Request {
	u := &url.URL{Path: path}
	return &http.Request{URL: u}
}

func makeBasicProxy() *httputil.ReverseProxy {
	r1, _ := regexp.Compile("/api/users/[0-9]+/details")
	r2, _ := regexp.Compile("/api/version.*")
	mappings := []PathMapping{PathMapping{Regex: r1, HostPort: "u.example.com:8080"}, PathMapping{Regex: r2, HostPort: "v.example.com:8080"}}
	return NewRegexRP(mappings)
}
