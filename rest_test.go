package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
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

type fakeTransport struct {
	mux http.Handler
}

func (f fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := httptest.NewRecorder()
	f.mux.ServeHTTP(r, req)
	return r.Result(), nil
}

func verifyAPICall(t *testing.T, client *Client, url, expected string)  {
	resp, err := client.API(Request{Method: Get, BaseURL: url})
	if err != nil {
		t.Fatal("API called failed", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Invalid status code", resp.StatusCode)
	}
	if strings.TrimSpace(resp.Body) != expected {
		t.Error("Received invalid message:", strings.TrimSpace(resp.Body))
	}
}

func TestCustomHTTPClient(t *testing.T) {
	mux := http.NewServeMux()
	const msg = `{"message": "custom-client"}`
	mux.HandleFunc("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, msg)
	}))
	client := &Client{HTTPClient: &http.Client{Transport: fakeTransport{mux}}}
	verifyAPICall(t, client, "/test", msg)
}

func TestNilHTTPClient(t *testing.T) {
	const msg = `{"message": "nil-client"}`
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, msg)
	}))
	defer fakeServer.Close()
	verifyAPICall(t, &Client{}, fakeServer.URL + "/test", msg)
}