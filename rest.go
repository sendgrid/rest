// Package rest allows for quick and easy access any REST or REST-like API.
package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Method contains the supported HTTP verbs.
type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

var (
	// DefaultHTTPClient is used when no custom HTTP client is specified.
	DefaultHTTPClient = &http.Client{Transport: http.DefaultTransport}
	// DefaultClient is used as the default for package level API calls.
	DefaultClient     = &Client{DefaultHTTPClient}
)

// Request holds the request to an API Call.
type Request struct {
	Method      Method
	BaseURL     string // e.g. https://api.sendgrid.com
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

// Response holds the response from an API call.
type Response struct {
	StatusCode int                 // e.g. 200
	Body       string              // e.g. {"result: success"}
	Headers    map[string][]string // e.g. map[X-Ratelimit-Limit:[600]]
}

// AddQueryParameters adds query paramaters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	baseURL += "?"
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + params.Encode()
}

// BuildRequestObject creates the HTTP request object.
func BuildRequestObject(request Request) (*http.Request, error) {
	req, err := http.NewRequest(string(request.Method), request.BaseURL, bytes.NewBuffer(request.Body))
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}
	if len(request.Body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, err
}

// MakeRequest makes the API call.
func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	client := c.HTTPClient
	if client == nil {
		client = DefaultHTTPClient
	}
	res, err := client.Do(req)
	return res, err
}

// BuildResponse builds the response struct.
func BuildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	response := Response{
		StatusCode: res.StatusCode,
		Body:       string(body),
		Headers:    res.Header,
	}
	return &response, nil
}

// Client provides support for custom API call configurations.
type Client struct {
	// HTTPClient is a custom *http.Client. If nil, DefaultHTTPClient will be used.
	HTTPClient *http.Client
}

// API is the main interface to the API.
func (c *Client) API(request Request) (*Response, error) {
	// Add any query parameters to the URL.
	if len(request.QueryParams) != 0 {
		request.BaseURL = AddQueryParameters(request.BaseURL, request.QueryParams)
	}

	// Build the HTTP request object.
	req, err := BuildRequestObject(request)
	if err != nil {
		return nil, err
	}

	// Build the HTTP client and make the request.
	res, err := c.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	// Build Response object.
	response, err := BuildResponse(res)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// API is the main interface to the API. DefaultClient is used to perform
// API requests.
func API(request Request) (*Response, error) {
	return DefaultClient.API(request)
}

// MakeRequest makes the API call. DefaultClient is used to perform the API call.
func MakeRequest(req *http.Request) (*http.Response, error) {
	return DefaultClient.MakeRequest(req)
}
