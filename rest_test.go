package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBuildURL(t *testing.T) {
	host := "http://api.test.com"
	queryParams := make(map[string]string)
	queryParams["test"] = "1"
	queryParams["test2"] = "2"
	testURL := AddQueryParameters(host, queryParams)
	if testURL != "http://api.test.com?test=1&test2=2" {
		t.Error("Bad BuildURL result")
	}
}

func TestBuildRequest(t *testing.T) {
	method := Get
	baseURL := "http://api.test.com"
	key := "API_KEY"
	Headers := make(map[string]string)
	Headers["Content-Type"] = "application/json"
	Headers["Authorization"] = "Bearer " + key
	queryParams := make(map[string]string)
	queryParams["test"] = "1"
	queryParams["test2"] = "2"
	request := Request{
		Method:      method,
		BaseURL:     baseURL,
		Headers:     Headers,
		QueryParams: queryParams,
	}
	req, e := BuildRequestObject(request)
	if e != nil {
		t.Errorf("Rest failed to BuildRequest. Returned error: %v", e)
	}
	if req == nil {
		t.Errorf("Failed to BuildRequest.")
	}
}

func TestBuildResponse(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"message\": \"success\"}")
	}))
	defer fakeServer.Close()
	baseURL := fakeServer.URL
	method := Get
	request := Request{
		Method:  method,
		BaseURL: baseURL,
	}
	req, e := BuildRequestObject(request)
	res, e := MakeRequest(req)
	response, e := BuildResponse(res)
	if response.StatusCode != 200 {
		t.Error("Invalid status code in BuildResponse")
	}
	if len(response.Body) == 0 {
		t.Error("Invalid response body in BuildResponse")
	}
	if len(response.Headers) == 0 {
		t.Error("Invalid response headers in BuildResponse")
	}
	if e != nil {
		t.Errorf("Rest failed to make a valid API request. Returned error: %v", e)
	}
}

func TestRest(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"message\": \"success\"}")
	}))
	defer fakeServer.Close()
	host := fakeServer.URL
	endpoint := "/test_endpoint"
	baseURL := host + endpoint
	key := "API_KEY"
	Headers := make(map[string]string)
	Headers["Content-Type"] = "application/json"
	Headers["Authorization"] = "Bearer " + key
	method := Get
	queryParams := make(map[string]string)
	queryParams["test"] = "1"
	queryParams["test2"] = "2"
	request := Request{
		Method:      method,
		BaseURL:     baseURL,
		Headers:     Headers,
		QueryParams: queryParams,
	}
	response, e := API(request)
	if response.StatusCode != 200 {
		t.Error("Invalid status code")
	}
	if len(response.Body) == 0 {
		t.Error("Invalid response body")
	}
	if len(response.Headers) == 0 {
		t.Error("Invalid response headers")
	}
	if e != nil {
		t.Errorf("Rest failed to make a valid API request. Returned error: %v", e)
	}
}

func TestDefaultContentTypeWithBody(t *testing.T) {
	host := "http://localhost"
	method := Get
	request := Request{
		Method:  method,
		BaseURL: host,
		Body:    []byte("Hello World"),
	}
	response, _ := BuildRequestObject(request)
	if response.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set to the correct default value when a body is set.")
	}
}

func TestCustomContentType(t *testing.T) {
	host := "http://localhost"
	Headers := make(map[string]string)
	Headers["Content-Type"] = "custom"
	method := Get
	request := Request{
		Method:  method,
		BaseURL: host,
		Headers: Headers,
		Body:    []byte("Hello World"),
	}
	response, _ := BuildRequestObject(request)
	if response.Header.Get("Content-Type") != "custom" {
		t.Error("Content-Type not modified correctly")
	}
}

func TestCustomHTTPClient(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 20)
		fmt.Fprintln(w, "{\"message\": \"success\"}")
	}))
	defer fakeServer.Close()
	host := fakeServer.URL
	endpoint := "/test_endpoint"
	baseURL := host + endpoint
	method := Get
	request := Request{
		Method:  method,
		BaseURL: baseURL,
	}
	customClient := &Client{&http.Client{Timeout: time.Millisecond * 10}}
	_, err := customClient.API(request)
	if err == nil {
		t.Error("A timeout did not trigger as expected")
	}
	if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") == false {
		t.Error("We did not receive the Timeout error")
	}
}
