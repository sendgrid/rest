[![Build Status](https://travis-ci.org/sendgrid/rest.svg?branch=master)](https://travis-ci.org/sendgrid/rest) [![GoDoc](https://godoc.org/github.com/sendgrid/rest?status.png)](http://godoc.org/github.com/sendgrid/rest)

**HTTP REST client, simplified for Go**

Here is a quick example:

`GET /your/api/{param}/call`

```go
package main

import "github.com/sendgrid/rest"
import "fmt"

func main() {
	const host = "https://api.example.com"
	param := "myparam"
	endpoint := "/your/api/" + param + "/call"
	baseURL := host + endpoint
	method := rest.Get
	request := rest.Request{
		Method:         method,
		BaseURL:        baseURL,
	}
	response, err := rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}
}
```

`POST /your/api/{param}/call` with headers, query parameters and a request body.

```go
package main

import "github.com/sendgrid/rest"
import "fmt"

func main() {
	const host = "https://api.example.com"
	param := "myparam"
	endpoint := "/your/api/" + param + "/call"
	baseURL := host + endpoint
	requestHeaders := make(map[string]string)
	key := os.Getenv("API_KEY")
	requestHeaders["Authorization"] = "Bearer " + key
	requestHeaders["X-Test"] = "Test"
	var requestBody = []byte(`{"some": 0, "awesome": 1, "data": 3}`)
	queryParams := make(map[string]string)
	queryParams["hello"] = "0"
	queryParams["world"] = "1"
	method := rest.Post
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL,
		RequestHeaders: requestHeaders,
		QueryParams:    queryParams,
		RequestBody:    requestBody,
	}
	response, err := rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.ResponseBody)
		fmt.Println(response.ResponseHeaders)
	}
}
```

# Installation

`go get github.com/sendgrid/rest`

## Usage ##

Following is an example using SendGrid. You can get your free account [here](https://sendgrid.com/free?source=python-http-client).

First, update your environment with your [SENDGRID_API_KEY](https://app.sendgrid.com/settings/api_keys).

```bash
echo "export SENDGRID_API_KEY='YOUR_API_KEY'" > sendgrid.env
echo "sendgrid.env" >> .gitignore
source ./sendgrid.env
```

Following is an abridged example, here is the [full working code](https://github.com/sendgrid/rest/tree/master/examples).

```bash
go run examples/example.go
```

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/rest"
	"os"
)

func main() {

	// Build the URL
	const host = "https://api.sendgrid.com"
	endpoint := "/v3/api_keys"
	baseURL := host + endpoint

	// Build the request headers
	key := os.Getenv("SENDGRID_API_KEY")
	requestHeaders := make(map[string]string)
	requestHeaders["Content-Type"] = "application/json"
	requestHeaders["Authorization"] = "Bearer " + key

	// GET Collection
	method := rest.Get

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
	response, err := rest.API(request)

	// POST
	method = rest.Post

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
	response, err = rest.API(request)

	// Get a particular return value.
	// Note that you can unmarshall into a struct if
	// you know the JSON structure in advance.
	b := []byte(response.ResponseBody)
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println(err)
	}
	m := f.(map[string]interface{})
	apiKey := m["api_key_id"].(string)

	// GET Single
	method = rest.Get

	// Make the API call
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
	}
	response, err = rest.API(request)

	// PATCH
	method = rest.Patch

	requestBody = []byte(`{
        "name": "A New Hope"
    }`)
	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
		RequestBody:    requestBody,
	}
	response, err = rest.API(request)

	// PUT
	method = rest.Put

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
	response, err = rest.API(request)

	// DELETE
	method = rest.Delete

	request = rest.Request{
		Method:         method,
		BaseURL:        baseURL + "/" + apiKey,
		RequestHeaders: requestHeaders,
		QueryParams:    queryParams,
		RequestBody:    requestBody,
	}
	response, err = rest.API(request)
}

```

# Announcements

[2016.04.05] - We hit version 1!

# Roadmap

[Milestones](https://github.com/sendgrid/rest/milestones)

# How to Contribute

We encourage contribution to our libraries, please see our [CONTRIBUTING](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md) guide for details.

* [Feature Request](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#feature_request)
* [Bug Reports](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#submit_a_bug_report)
* [Improvements to the Codebase](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#improvements_to_the_codebase)

# About

![SendGrid Logo]
(https://assets3.sendgrid.com/mkt/assets/logos_brands/small/sglogo_2015_blue-9c87423c2ff2ff393ebce1ab3bd018a4.png)

rest is guided and supported by the SendGrid [Developer Experience Team](mailto:dx@sendgrid.com).

rest is maintained and funded by SendGrid, Inc. The names and logos for rest are trademarks of SendGrid, Inc.
