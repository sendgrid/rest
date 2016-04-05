package main

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/rest"
	"os"
)

func main() {

	// Build the URL
	host := "https://api.sendgrid.com"
	endpoint := "/v3/api_keys"
	baseURL := host + endpoint

	// Build the request headers
	key := os.Getenv("SENDGRID_API_KEY")
	requestHeaders := make(map[string]string)
	requestHeaders["Content-Type"] = "application/json"
	requestHeaders["Authorization"] = "Bearer " + key

	// GET Collection
	method := "GET"

	// Build the query parameters
	queryParams := make(map[string]string)
	queryParams["limit"] = "100"
	queryParams["offset"] = "0"

	// Make the API call
	request := rest.Request{
		Method:         method,
		BaseURL:        baseURL,
		RequestHeaders: requestHeaders,
		QueryParams:    queryParams,
	}
	response, e := rest.API(request)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}

	// POST
	method = "POST"

	var requestBody = []byte(` {
        "name": "My API Key",
        "scopes": [
            "mail.send",
            "alerts.create",
            "alerts.read"
        ]
    }`)
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL,
		RequestHeaders: requestHeaders,
		QueryParams:    queryParams,
		RequestBody:    requestBody,
	}
	response, e = rest.API(request)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}

	// Get a particular return value.
	// Note that you can unmarshall into a struct if
	// you know the JSON structure in advance.
	b := []byte(response.ResponseBody)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println(err)
	}
	m := f.(map[string]interface{})
	apiKey := m["api_key_id"].(string)

	// GET Single
	method = "GET"

	// Make the API call
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
	}
	response, e = rest.API(request)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}

	// PATCH
	method = "PATCH"

	requestBody = []byte(`{
        "name": "A New Hope"
    }`)
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
		RequestBody:    requestBody,
	}
	response, e = rest.API(request)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}

	// PUT
	method = "PUT"

	requestBody = []byte(`{
        "name": "A New Hope",
        "scopes": [
            "user.profile.read",
            "user.profile.update"
        ]
    }`)
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
		RequestBody:    requestBody,
	}
	response, e = rest.API(request)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}

	// DELETE
	method = "DELETE"

	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
		QueryParams:    queryParams,
		RequestBody:    requestBody,
	}
	response, e = rest.API(request)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseHeaders)
	}
}
