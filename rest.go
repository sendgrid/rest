// Package rest allows for quick and easy access any REST or REST-like API.
package rest

import (
	"errors"
	"net/http"
)

// Method contains the supported HTTP verbs.
type Method string

// Supported HTTP verbs.
const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

// Request holds the request to an API Call.
type Request struct {
	Method      Method
	BaseURL     string // e.g. https://api.sendgrid.com
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

// RestError is a struct for an error handling.
type RestError struct {
	Response *Response
}

// Error is the implementation of the error interface.
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

func makeRequest(req *http.Request) (*http.Response, error) {
	return DefaultClient.HTTPClient.Do(req)
}

// API is the main interface to the API.
func API(request *http.Request) (*http.Response, error) {
	return DefaultClient.API(request)
}

// The following functions enable the ability to define a
// custom HTTP Client

// MakeRequest makes the API call.
func (c *Client) makeRequest(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

// API is the main interface to the API.
func (c *Client) API(request *http.Request) (*http.Response, error) {
	if request == nil {
		return nil, errors.New("invalid request: request cannot be nil")
	}

	if request.Header.Get("Content-Type") == "" && request.Body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return c.makeRequest(request)
}
