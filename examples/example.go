package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Status)
		fmt.Println(string(b))
		fmt.Println(response.Header)
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
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Status)
		fmt.Println(string(b))
		fmt.Println(response.Header)
	}

	// Get a particular return value.
	// Note that you can unmarshall into a struct if
	// you know the JSON structure in advance.

	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
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
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Status)
		fmt.Println(string(b))
		fmt.Println(response.Header)
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
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Status)
		fmt.Println(string(b))
		fmt.Println(response.Header)
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
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Status)
		fmt.Println(string(b))
		fmt.Println(response.Header)
	}

	// DELETE
	request, err = http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	request.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Status)
		fmt.Println(response.Header)
	}
}
