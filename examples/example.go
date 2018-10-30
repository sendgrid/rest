package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/sendgrid/rest"
)

func main() {

	// Build the URL
	const host = "https://api.sendgrid.com"
	endpoint := "/v3/api_keys"
	baseURL := host + endpoint

	// Build the request headers
	key := os.Getenv("SENDGRID_API_KEY")
	Headers := make(map[string]string)
	Headers["Authorization"] = "Bearer " + key

	// GET Collection
	// Build the query parameters
	queryParams := map[string]string{
		"limit":  "100",
		"offset": "0",
	}

	// Make the API call
	req := rest.Request{
		Method:      rest.Get,
		BaseURL:     baseURL,
		Headers:     Headers,
		QueryParams: queryParams,
	}
	resp, err := rest.Send(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		fmt.Println(resp.Headers)
	}

	// POST
	body := `
		{
		    "name": "My API Key",
		    "scopes": [
		        "mail.send",
		        "alerts.create",
		        "alerts.read"
		    ]
		}`
	req = rest.Request{
		Method:      rest.Post,
		BaseURL:     baseURL,
		Headers:     Headers,
		QueryParams: queryParams,
		Body:        []byte(body),
	}
	resp, err = rest.Send(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		fmt.Println(resp.Headers)
	}

	// Get a particular return value.
	var data struct {
		APIKeyID string `json:"api_key_id"`
	}
	err = json.Unmarshal([]byte(resp.Body), &data)
	if err != nil {
		fmt.Println(err)
	}
	apiKey := data.APIKeyID

	// GET Single
	// Make the API call
	req = rest.Request{
		Method:  rest.Get,
		BaseURL: path.Join(baseURL, apiKey),
		Headers: Headers,
	}
	resp, err = rest.Send(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		fmt.Println(resp.Headers)
	}

	// PATCH
	body = `{"name": "A New Hope"}`
	req = rest.Request{
		Method:  rest.Patch,
		BaseURL: path.Join(baseURL, apiKey),
		Headers: Headers,
		Body:    []byte(body),
	}
	resp, err = rest.Send(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		fmt.Println(resp.Headers)
	}

	// PUT
	body = `
		{
		    "name": "A New Hope",
		    "scopes": [
		        "user.profile.read",
		        "user.profile.update"
		    ]
		}`
	req = rest.Request{
		Method:  rest.Put,
		BaseURL: path.Join(baseURL, apiKey),
		Headers: Headers,
		Body:    []byte(body),
	}
	resp, err = rest.Send(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		fmt.Println(resp.Headers)
	}

	// DELETE
	req = rest.Request{
		Method:      rest.Delete,
		BaseURL:     path.Join(baseURL, apiKey),
		Headers:     Headers,
		QueryParams: queryParams,
	}
	resp, err = rest.Send(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Headers)
	}
}
