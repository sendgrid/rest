// Package rest allows for quick and easy access any REST or REST-like API.
package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// DefaultClient is used if no custom HTTP client is defined
var DefaultClient = &Client{HTTPClient: http.DefaultClient}

var supportedMethods = map[string]struct{}{
	http.MethodGet:    struct{}{},
	http.MethodPost:   struct{}{},
	http.MethodPut:    struct{}{},
	http.MethodPatch:  struct{}{},
	http.MethodDelete: struct{}{},
}

// Client allows modification of client headers, redirect policy
// and other settings
// See https://golang.org/pkg/net/http
type Client struct {
	HTTPClient *http.Client
}

// Response holds the response from an API call.
type Response struct {
	StatusCode int         // e.g. 200
	Body       string      // e.g. {"result: success"}
	Headers    http.Header // e.g. map[X-Ratelimit-Limit:[600]]
}

func makeRequest(req *http.Request) (*http.Response, error) {
	return DefaultClient.HTTPClient.Do(req)
}

func buildResponse(res *http.Response) (*Response, error) {
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

// API is the main interface to the API.
func API(request *http.Request) (*Response, error) {
	return DefaultClient.API(request)
}

// The following functions enable the ability to define a
// custom HTTP Client

// MakeRequest makes the API call.
func (c *Client) makeRequest(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// API is the main interface to the API.
func (c *Client) API(request *http.Request) (*Response, error) {
	if request == nil {
		return nil, errors.New("invalid request: request cannot be nil")
	}
	_, supportedMethod := supportedMethods[request.Method]

	if !supportedMethod {
		return nil, errors.New("unsupported request method:" + request.Method)
	}

	if request.Header.Get("Content-Type") == "" && request.Body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	res, err := c.makeRequest(request)
	if err != nil {
		return nil, err
	}

	response, err := buildResponse(res)
	if err != nil {
		return nil, err
	}

	return response, nil
}
