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

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	response, err := rest.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer func() {
			err = response.Body.Close()
			if err != nil {
				fmt.Println("encountered an error closing the response body:", err)
			}
		}()

		var b []byte

		b, err = ioutil.ReadAll(response.Body)
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

	req, err = http.NewRequest(http.MethodPost, baseURL.String(), strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer func() {
			err = response.Body.Close()
			if err != nil {
				fmt.Println("encountered an error closing the response body:", err)
			}
		}()
		var b []byte
		b, err = ioutil.ReadAll(response.Body)
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

	defer func() {
		err = response.Body.Close()
		if err != nil {
			fmt.Println("encountered an error closing the response body:", err)
		}
	}()
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

	req, err = http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer func() {
			err = response.Body.Close()
			if err != nil {
				fmt.Println("encountered an error closing the response body:", err)
			}
		}()

		var b []byte

		b, err = ioutil.ReadAll(response.Body)
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

	req, err = http.NewRequest(http.MethodPatch, baseURL.String(), strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer func() {
			err = response.Body.Close()
			if err != nil {
				fmt.Println("encountered an error closing the response body:", err)
			}
		}()

		var b []byte
		b, err = ioutil.ReadAll(response.Body)
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

	req, err = http.NewRequest(http.MethodPut, baseURL.String(), strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		defer func() {
			err = response.Body.Close()
			if err != nil {
				fmt.Println("encountered an error closing the response body:", err)
			}
		}()

		var b []byte
		b, err = ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Status)
		fmt.Println(string(b))
		fmt.Println(response.Header)
	}

	// DELETE
	req, err = http.NewRequest(http.MethodDelete, baseURL.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	response, err = rest.API(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Status)
		fmt.Println(response.Header)
	}
}
