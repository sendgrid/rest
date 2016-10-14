package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/sendgrid/rest"
)

func main() {
	// Build the URL
	const host = "api.sendgrid.com"
	endpoint := "/v3/api_keys"
	key := os.Getenv("SENDGRID_API_KEY")

	// GET
	params := url.Values{
		"limit":  {"100"},
		"offset": {"0"},
	}

	baseURL := &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     endpoint,
		RawQuery: params.Encode(),
	}

	request, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	request.Header.Set("Authorization", "Bearer "+key)

	response, err := rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// POST
	body := `{
        "name": "My API Key",
        "scopes": [
            "mail.send",
            "alerts.create",
            "alerts.read"
        ]
    }`

	request, err = http.NewRequest(http.MethodPost, baseURL.String(), strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// Get a particular return value.
	// Note that you can unmarshall into a struct if
	// you know the JSON structure in advance.
	b := []byte(response.Body)
	var payload struct {
		APIKeyID string `json:"api_key_id"`
	}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println(err)
	}

	apiKey := payload.APIKeyID

	baseURL = &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(endpoint, apiKey),
	}

	request, err = http.NewRequest(http.MethodGet, baseURL.String(), nil)
	request.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// PATCH
	body = `{
        "name": "A New Hope"
    }`

	request, err = http.NewRequest(http.MethodPatch, baseURL.String(), strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(request)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// PUT
	body = `{
        "name": "A New Hope",
        "scopes": [
            "user.profile.read",
            "user.profile.update"
        ]
    }`

	request, err = http.NewRequest(http.MethodPut, baseURL.String(), strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	// DELETE
	request, err = http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	request.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
	}
}
