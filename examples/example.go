package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/sendgrid/rest"
)

func main() {
	// Build the URL
	const host = "https://api.sendgrid.com"
	endpoint := "/v3/api_keys"
	key := os.Getenv("SENDGRID_API_KEY")

	// GET
	baseURL, _ := url.Parse(host + endpoint)

	params := url.Values{}
	params.Add("limit", "100")
	params.Add("offset", "0")
	baseURL.RawQuery = params.Encode()

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
	body := []byte(` {
        "name": "My API Key",
        "scopes": [
            "mail.send",
            "alerts.create",
            "alerts.read"
        ]
    }`)

	request, err = http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
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
	var f struct {
		APIKeyID string `json:"api_key_id"`
	}
	err = json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println(err)
	}

	apiKey := f.APIKeyID

	baseURL, err = url.Parse(host + endpoint + "/" + apiKey)

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
	body = []byte(`{
        "name": "A New Hope"
    }`)

	baseURL, err = url.Parse(host + endpoint + "/" + apiKey)
	request, err = http.NewRequest(http.MethodPatch, baseURL.String(), bytes.NewReader(body))
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
	body = []byte(`{
        "name": "A New Hope",
        "scopes": [
            "user.profile.read",
            "user.profile.update"
        ]
    }`)

	baseURL, err = url.Parse(host + endpoint + "/" + apiKey)
	request, err = http.NewRequest(http.MethodPut, baseURL.String(), bytes.NewReader(body))
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
	baseURL, err = url.Parse(host + endpoint + "/" + apiKey)
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
