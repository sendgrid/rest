package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Request holds the request to an API Call
type Request struct {
    Method string
    BaseURL string
    RequestHeaders map[string]string
    QueryParams map[string]string
    RequestBody []byte
}

// Response holds the response from an API call
type Response struct {
	StatusCode      int                 // e.g. 200
	ResponseBody    string              // e.g. {"result: success"}
	ResponseHeaders map[string][]string // e.g. map[X-Ratelimit-Limit:[600]]
}

// BuildURL adds query paramaters to the URL
func BuildURL(baseURL string, queryParams map[string]string) (string) {
    baseURL += "?"
    params := url.Values{}
    for key, value := range queryParams {
        params.Add(key, value)
    }
    return baseURL + params.Encode()
}

// BuildRequest creates the HTTP request object
func BuildRequest(request Request) (*http.Request, error) {
    req, e := http.NewRequest(request.Method, request.BaseURL, bytes.NewBuffer(request.RequestBody))
	for key, value := range request.RequestHeaders {
		req.Header.Set(key, value)
	}
    return req, e
}

// MakeRequest makes the API call
func MakeRequest(req *http.Request) (*http.Response, error) {
    var Client = &http.Client{
		Transport: http.DefaultTransport,
	}
	res, e := Client.Do(req)
	return res, e
}

// BuildResponse builds the response struct
func BuildResponse(res *http.Response) (Response, error) {
    var response Response
    response.StatusCode = res.StatusCode
	body, e := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	response.ResponseBody = string(body)
	response.ResponseHeaders = res.Header
    return response, e
}

// API allows for quick and easy access any REST or REST-like API.
func API(request Request) (Response, error) {
	
	// Build the final URL
	if len(request.QueryParams) != 0 {
		request.BaseURL = BuildURL(request.BaseURL, request.QueryParams)
	}

	// Build the http request object
	req, e := BuildRequest(request)

	// Build the HTTP client and make the request
	res, e := MakeRequest(req)

	// Build Response object
    response, e := BuildResponse(res)

	return response, e
}
