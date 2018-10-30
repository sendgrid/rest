package rest

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestBuildURL(t *testing.T) {
	t.Parallel()
	host := "http://api.test.com"
	queryParams := map[string]string{
		"test":  "1",
		"test2": "2",
	}
	testURL := AddQueryParameters(host, queryParams)
	if testURL != "http://api.test.com?test=1&test2=2" {
		t.Error("Bad BuildURL result")
	}
}

func TestBuildRequest(t *testing.T) {
	t.Parallel()
	method := Get
	baseURL := "http://api.test.com"
	key := "API_KEY"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + key,
	}
	queryParams := map[string]string{
		"test":  "1",
		"test2": "2",
	}
	request := Request{
		Method:      method,
		BaseURL:     baseURL,
		Headers:     headers,
		QueryParams: queryParams,
	}
	req, err := BuildRequestObject(request)
	if err != nil {
		t.Errorf("Rest failed to BuildRequest. Returned error: %v", err)
	}
	if req == nil {
		t.Errorf("Failed to BuildRequest.")
	}

	//Start PrintRequest
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Request: %s", requestDump)
	//End Print Request
}

func TestBuildBadRequest(t *testing.T) {
	t.Parallel()
	request := Request{
		Method: Method("@"),
	}
	req, err := BuildRequestObject(request)
	if err == nil {
		t.Errorf("Expected an error for a bad HTTP Method")
	}
	if req != nil {
		t.Errorf("If there's an error there shouldn't be a Request.")
	}
}

func TestBuildBadAPI(t *testing.T) {
	t.Parallel()
	request := Request{
		Method: Method("@"),
	}
	res, err := API(request)
	if err == nil {
		t.Errorf("Expected an error for a bad HTTP Method")
	}
	if res != nil {
		t.Errorf("If there's an error there shouldn't be a Response.")
	}
}

func TestBuildResponse(t *testing.T) {
	t.Parallel()
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
	}))
	defer fakeServer.Close()
	baseURL := fakeServer.URL
	method := Get
	request := Request{
		Method:  method,
		BaseURL: baseURL,
	}
	req, err := BuildRequestObject(request)
	if err != nil {
		t.Error("Failed to BuildRequestObject", err)
	}
	res, err := MakeRequest(req)
	if err != nil {
		t.Error("Failed to MakeRequest", err)
	}
	response, err := BuildResponse(res)
	if response.StatusCode != 200 {
		t.Error("Invalid status code in BuildResponse")
	}
	if len(response.Body) == 0 {
		t.Error("Invalid response body in BuildResponse")
	}
	if len(response.Headers) == 0 {
		t.Error("Invalid response headers in BuildResponse")
	}
	if err != nil {
		t.Errorf("Rest failed to make a valid API request. Returned error: %v", err)
	}

	//Start Print Request
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Request: %s", requestDump)
	//End Print Request

}

type panicResponse struct{}

func (*panicResponse) Read(p []byte) (n int, err error) {
	panic(bytes.ErrTooLarge)
}

func (*panicResponse) Close() error {
	return nil
}

func TestBuildBadResponse(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Body: new(panicResponse),
	}
	_, err := BuildResponse(res)
	if err == nil {
		t.Errorf("This was a bad response and error should be returned")
	}
}

func TestRest(t *testing.T) {
	t.Parallel()
	testingAPI(t, Send)
	testingAPI(t, API)
}

func testingAPI(t *testing.T, fn func(request Request) (*Response, error)) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
	}))
	defer fakeServer.Close()

	host := fakeServer.URL
	endpoint := "/test_endpoint"
	baseURL := host + endpoint
	key := "API_KEY"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + key,
	}
	method := Get
	queryParams := map[string]string{
		"test":  "1",
		"test2": "2",
	}
	request := Request{
		Method:      method,
		BaseURL:     baseURL,
		Headers:     headers,
		QueryParams: queryParams,
	}

	//Start Print Request
	req, err := BuildRequestObject(request)
	if err != nil {
		t.Errorf("Error during BuildRequestObject: %v", err)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Request: %s", requestDump)
	//End Print Request

	response, err := fn(request)

	if response.StatusCode != 200 {
		t.Error("Invalid status code")
	}
	if len(response.Body) == 0 {
		t.Error("Invalid response body")
	}
	if len(response.Headers) == 0 {
		t.Error("Invalid response headers")
	}
	if err != nil {
		t.Errorf("Rest failed to make a valid API request. Returned error: %v", err)
	}
}

func TestDefaultContentTypeWithBody(t *testing.T) {
	t.Parallel()
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

	//Start Print Request
	t.Logf("Request Body: %s", request.Body)

	requestDump, err := httputil.DumpRequest(response, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Request: %s", requestDump)
	//End Print Request
}

func TestCustomContentType(t *testing.T) {
	t.Parallel()
	host := "http://localhost"
	headers := map[string]string{
		"Content-Type": "custom",
	}
	method := Get
	request := Request{
		Method:  method,
		BaseURL: host,
		Headers: headers,
		Body:    []byte("Hello World"),
	}
	response, _ := BuildRequestObject(request)
	if response.Header.Get("Content-Type") != "custom" {
		t.Error("Content-Type not modified correctly")
	}

	//Start Print Request
	requestDump, err := httputil.DumpRequest(response, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Request: %s", requestDump)
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
	baseURL := host + endpoint
	method := Get
	request := Request{
		Method:  method,
		BaseURL: baseURL,
	}

	customClient := &Client{&http.Client{Timeout: time.Millisecond * 10}}
	_, err := customClient.Send(request)
	if err == nil {
		t.Error("A timeout did not trigger as expected")
	}
	if !strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
		t.Error("We did not receive the Timeout error")
	}
}

func TestRestError(t *testing.T) {
	t.Parallel()
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	response := &Response{
		StatusCode: 400,
		Body:       `{"result": "failure"}`,
		Headers:    headers,
	}

	var err error = &RestError{Response: response}

	if err.Error() != `{"result": "failure"}` {
		t.Error("Invalid error message.")
	}
}

func TestRepoFiles(t *testing.T) {
	files := []string{
		".env_sample",
		".gitignore",
		".travis.yml",
		"CHANGELOG.md",
		"CODE_OF_CONDUCT.md",
		"CONTRIBUTING.md",
		".github/ISSUE_TEMPLATE",
		"LICENSE.txt",
		".github/PULL_REQUEST_TEMPLATE",
		"README.md",
		"TROUBLESHOOTING.md",
		"USAGE.md",
	}

	for _, file := range files {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			t.Errorf("Repo file does not exist: %v", file)
		}
	}
}

func TestLicenseYear(t *testing.T) {
	t.Parallel()
	f, err := os.Open("LICENSE.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)

	currentYear := time.Now().Year()
	re := regexp.MustCompile(fmt.Sprintf(`\b%d\b`, currentYear))
	if !re.MatchReader(r) {
		t.Error("Incorrect Year in License Copyright")
	}
}
