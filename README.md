[![Build Status](https://travis-ci.org/sendgrid/rest.svg?branch=master)](https://travis-ci.org/sendgrid/rest) [![GoDoc](https://godoc.org/github.com/sendgrid/rest?status.png)](http://godoc.org/github.com/sendgrid/rest)

**Quickly and easily access any RESTful or RESTful-like API.**

If you are looking for the SendGrid API client library, please see [this repo](https://github.com/sendgrid/sendgrid-go).

# Announcements

All updates to this library is documented in our [CHANGELOG](https://github.com/sendgrid/rest/blob/master/CHANGELOG.md).

# Installation

## Prerequisites

- Go version 1.6

## Install Package

```bash
go get github.com/sendgrid/rest
```

# Quick Start

`GET /your/api/{param}/call`

```go
package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sendgrid/rest"
)

func main() {
	const host = "https://api.example.com"
	param := "myparam"
	endpoint := "/your/api/" + param + "/call"

	baseURL, _ := url.Parse(host + endpoint)

	request, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	response, err := rest.API(request)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
```

`POST /your/api/{param}/call` with headers, query parameters and a request body.

```go
package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sendgrid/rest"
)

func main() {
	const host = "https://api.example.com"
	param := "myparam"
	endpoint := "/your/api/" + param + "/call"

	baseURL, _ := url.Parse(host + endpoint)

	key := os.Getenv("API_KEY")
	var body = []byte(`{"some": 0, "awesome": 1, "data": 3}`)

	params := url.Values{}
	params.Add("hello", "0")
	params.Add("world", "1")
	baseURL.RawQuery = params.Encode()

	request, err = http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+key)
	request.Header.Set("X-Test", "Test")

	response, err := rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
```

# Usage

- [Example Code](https://github.com/sendgrid/rest/tree/master/examples)

## Roadmap

If you are interested in the future direction of this project, please take a look at our [milestones](https://github.com/sendgrid/rest/milestones). We would love to hear your feedback.

## How to Contribute

We encourage contribution to our projects, please see our [CONTRIBUTING](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md) guide for details.

Quick links:

- [Feature Request](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#feature_request)
- [Bug Reports](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#submit_a_bug_report)
- [Sign the CLA to Create a Pull Request](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#cla)
- [Improvements to the Codebase](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#improvements_to_the_codebase)

# About

rest is guided and supported by the SendGrid [Developer Experience Team](mailto:dx@sendgrid.com).

rest is maintained and funded by SendGrid, Inc. The names and logos for rest are trademarks of SendGrid, Inc.

![SendGrid Logo]
(https://uiux.s3.amazonaws.com/2016-logos/email-logo%402x.png)
