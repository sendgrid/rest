// Package rest allows for quick and easy access any REST or REST-like API.
package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Method contains the supported HTTP verbs.
type Method string

// Supported HTTP verbs.
const (
	Get    Method = http.MethodGet
	Post   Method = http.MethodPost
	Put    Method = http.MethodPut
	Patch  Method = http.MethodPatch
	Delete Method = http.MethodDelete
)

// Request holds REST API request data.
type Request struct {
	Method      Method
	BaseURL     string // e.g. https://api.sendgrid.com
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

// RestError is an error derived from an REST API response.
type RestError struct { //nolint: golint
	Response *Response
}

// Error returns an error message associated with an HTTP response.
func (e *RestError) Error() string {
	return e.Response.Body
}

// DefaultClient is used if no custom HTTP client is defined
var DefaultClient = &Client{HTTPClient: http.DefaultClient}

// Client allows modification of client headers, redirect policy
// and other settings
// See https://golang.org/pkg/net/http
type Client struct {
	HTTPClient *http.Client
}

// Response holds the response from an API call.
type Response struct {
	StatusCode int                 // e.g. 200
	Body       string              // e.g. {"result: success"}
	Headers    map[string][]string // e.g. map[X-Ratelimit-Limit:[600]]
}

// AddQueryParameters adds query parameters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return baseURL
	}
	params := make(url.Values, len(queryParams))
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + "?" + params.Encode()
}

// BuildRequestObject builds an http.Request given a rest.Request
func BuildRequestObject(request Request) (*http.Request, error) {
	// Add any query parameters to the URL.
	if len(request.QueryParams) != 0 {
		request.BaseURL = AddQueryParameters(request.BaseURL, request.QueryParams)
	}
	const contentType = "Content-Type"
	const jsonMIME = "application/json"
	_, hasContentType := request.Headers[contentType]
	var setToJSON bool
	switch {
	case !hasContentType && len(request.Body) > 0:
		setToJSON = true
		fallthrough
	case request.Headers[contentType] == jsonMIME:
		// it is semantically okay to trim leading and trailing space around JSON
		request.Body = bytes.TrimSpace(request.Body)
	}
	var body io.Reader
	if len(request.Body) > 0 {
		body = bytes.NewReader(request.Body)
	}
	method := string(request.Method)
	req, err := http.NewRequest(method, request.BaseURL, body)
	if err != nil {
		return req, err
	}
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}
	if setToJSON {
		req.Header.Set(contentType, jsonMIME)
	}
	return req, err
}

// MakeRequest makes the API call.
func MakeRequest(req *http.Request) (*http.Response, error) {
	return DefaultClient.HTTPClient.Do(req)
}

// BuildResponse builds the response struct.
func BuildResponse(resp *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}
	rr := &Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Headers:    resp.Header,
	}
	return rr, err
}

// API supports old implementation (deprecated)
func API(request Request) (*Response, error) {
	return Send(request)
}

// Send uses the DefaultClient to send a request.
func Send(request Request) (*Response, error) {
	return DefaultClient.Send(request)
}

// The following functions enable the ability to define a
// custom HTTP Client

// MakeRequest is equivalent to http.Do.
func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// API supports old implementation (deprecated)
func (c *Client) API(request Request) (*Response, error) {
	return c.Send(request)
}

// Send sends a REST request and returns a REST response.
func (c *Client) Send(request Request) (*Response, error) {
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
	return BuildResponse(res)
}
