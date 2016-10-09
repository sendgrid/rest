package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBuildResponse(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"message\": \"success\"}")
	}))
	defer fakeServer.Close()

	baseURL, err := url.Parse(fakeServer.URL)
	if err != nil {
		t.Fatal("invalid url:", err)
	}

	request, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		t.Fatal("unable to formulate request:", err)
	}

	res, e := makeRequest(request)

	response, e := buildResponse(res)

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
	key := "API_KEY"

	baseURL, err := url.Parse(host + endpoint)

	if err != nil {
		t.Fatal("invalid url:", err)
	}

	params := url.Values{}
	params.Add("test", "1")
	params.Add("test2", "2")
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	response, e := API(req)

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
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"message\": \"success\"}")
	}))
	defer fakeServer.Close()

	baseURL, err := url.Parse(fakeServer.URL)
	if err != nil {
		t.Fatal("invalid url:", err)
	}

	request, err := http.NewRequest(http.MethodGet, baseURL.String(), bytes.NewReader([]byte(`Hello World`)))
	if err != nil {
		t.Fatal("unable to formulate request:", err)
	}

	_, err = API(request)

	if request.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set to the correct default value when a body is set.")
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

	customClient := &Client{&http.Client{Timeout: time.Millisecond * 10}}

	baseURL, err := url.Parse(host + endpoint)

	if err != nil {
		t.Fatal("invalid url:", err)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)

	_, err = customClient.API(req)
	if err == nil {
		t.Error("A timeout did not trigger as expected")
	}

	if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") == false {
		t.Error("We did not receive the Timeout error")
	}
}

func TestUnsupportedMethodsError(t *testing.T) {
	baseURL, err := url.Parse("http://localhost")
	if err != nil {
		t.Fatal("invalid url:", err)
	}

	request, err := http.NewRequest(http.MethodHead, baseURL.String(), bytes.NewReader([]byte(`Hello World`)))
	if err != nil {
		t.Fatal("unable to formulate request:", err)
	}

	_, err = API(request)

	if err == nil {
		t.Error("requesting an unsupported method should result in an error!")
	}
}

func TestNilRequestErrors(t *testing.T) {
	_, err := API(nil)

	if err == nil {
		t.Error("providing a nil request should result in an error!")
	}
}
