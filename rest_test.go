package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestBuildResponse(t *testing.T) {
	t.Parallel()
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
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

	response, e := makeRequest(request)

	if response.StatusCode != 200 {
		t.Error("Invalid status code in BuildResponse")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal("unable to read response body:", err)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			t.Fatal("encountered an error closing the response body:", err)
		}
	}()

	if len(body) == 0 {
		t.Error("Invalid response body in BuildResponse")
	}

	if len(response.Header) == 0 {
		t.Error("Invalid response headers in BuildResponse")
	}

	if e != nil {
		t.Errorf("Rest failed to make a valid API request. Returned error: %v", e)
	}

	//Start Print Request
	requestDump, err := httputil.DumpRequest(request, true)
	if err != nil {
		t.Errorf("Error : %v", err)
	}
	fmt.Println("Request :", string(requestDump))
	//End Print Request
}

func TestRest(t *testing.T) {
	t.Parallel()
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
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
	if err != nil {
		t.Fatal("invalid request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	response, e := API(req)

	if response.StatusCode != 200 {
		t.Error("Invalid status code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal("unable to read response body:", err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Fatal("encountered an error closing the response body:", err)
		}
	}()

	if len(body) == 0 {
		t.Error("Invalid response body in BuildResponse")
	}

	if len(response.Header) == 0 {
		t.Error("Invalid response headers")
	}

	if e != nil {
		t.Errorf("Rest failed to make a valid API request. Returned error: %v", e)
	}
}

func TestDefaultContentTypeWithBody(t *testing.T) {
	t.Parallel()

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
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
	if err != nil {
		t.Fatal("encoutered unexpected api error:", err)
	}

	if request.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set to the correct default value when a body is set")
	}

	//Start Print Request
	requestDump, err := httputil.DumpRequest(request, true)
	if err != nil {
		t.Errorf("Error : %v", err)
	}
	fmt.Println("Request :", string(requestDump))
	//End Print Request
}

func TestCustomContentType(t *testing.T) {
	t.Parallel()
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
	}))
	defer fakeServer.Close()

	baseURL, err := url.Parse(fakeServer.URL)
	if err != nil {
		t.Fatal("invalid url:", err)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), strings.NewReader("Hello World"))
	if err != nil {
		t.Error("invalid request:", err)
	}

	req.Header.Set("Content-Type", "custom")

	_, e := API(req)
	if e != nil {
		t.Error("encountered an unexpected error:", e)
	}

	if req.Header.Get("Content-Type") != "custom" {
		t.Error("Content-Type not modified correctly")
	}

	//Start Print Request
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Errorf("Error : %v", err)
	}
	fmt.Println("Request :", string(requestDump))
	//End Print Request
}

func TestCustomHTTPClient(t *testing.T) {
	t.Parallel()
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 20)
		fmt.Fprintln(w, `{"message": "success"}`)
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
	if err != nil {
		t.Error("invalid request:", err)
	}

	_, err = customClient.API(req)
	if err == nil {
		t.Error("A timeout did not trigger as expected")
	}

	if !strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
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

func TestRepoFiles(t *testing.T) {
	// Note: commenting out some of these files as they're causing the test to fail and don't exist upstream.
	files := []string{
		//		"docker/Dockerfile",
		//		"docker/docker-compose.yml",
		".env_sample",
		".gitignore",
		".travis.yml",
		//		".codeclimate.yml",
		"CHANGELOG.md",
		"CODE_OF_CONDUCT.md",
		"CONTRIBUTING.md",
		".github/ISSUE_TEMPLATE",
		"LICENSE.txt",
		".github/PULL_REQUEST_TEMPLATE",
		"README.md",
		"TROUBLESHOOTING.md",
		"USAGE.md",
		//		"USE_CASES.md",
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Repo file does not exist: %v", file)
		}
	}
}

func TestLicenseYear(t *testing.T) {
	t.Parallel()
	dat, err := ioutil.ReadFile("LICENSE.txt")

	currentYear := time.Now().Year()
	r := fmt.Sprintf("%d", currentYear)
	match, _ := regexp.MatchString(r, string(dat))

	if err != nil {
		t.Error("License File Not Found")
	}
	if !match {
		t.Error("Incorrect Year in License Copyright")
	}
}
