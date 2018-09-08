![SendGrid Logo](https://uiux.s3.amazonaws.com/2016-logos/email-logo%402x.png)

[![Build Status](https://travis-ci.org/sendgrid/rest.svg?branch=master)](https://travis-ci.org/sendgrid/rest)
[![GoDoc](https://godoc.org/github.com/sendgrid/rest?status.png)](http://godoc.org/github.com/sendgrid/rest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sendgrid/rest)](https://goreportcard.com/report/github.com/sendgrid/rest)
[![Email Notifications Badge](https://dx.sendgrid.com/badge/go)](https://dx.sendgrid.com/newsletter/go)
[![Twitter Follow](https://img.shields.io/twitter/follow/sendgrid.svg?style=social&label=Follow)](https://twitter.com/sendgrid)
[![GitHub contributors](https://img.shields.io/github/contributors/sendgrid/rest.svg)](https://github.com/sendgrid/rest/graphs/contributors)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE.txt)

**Quickly and easily access any RESTful or RESTful-like API.**

If you are looking for the SendGrid API client library, please see [this repo](https://github.com/sendgrid/sendgrid-go).

# Announcements

All updates to this library is documented in our [CHANGELOG](https://github.com/sendgrid/rest/blob/master/CHANGELOG.md).

# Table of Contents
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [How to Contribute](#contribute)
- [About](#about)
- [License](#license)

<a name="installation"></a>
# Installation

## Prerequisites

- Go version 1.6.X, 1.7.X, 1.8.X, 1.9.X or 1.10.X

## Install Package

```bash
go get github.com/sendgrid/rest
```

## Setup Environment Variables

### Initial Setup

```bash
cp .env_sample .env
```

### Environment Variable

Update the development environment with your [SENDGRID_API_KEY](https://app.sendgrid.com/settings/api_keys), for example:

```bash
echo "export SENDGRID_API_KEY='YOUR_API_KEY'" > sendgrid.env
echo "sendgrid.env" >> .gitignore
source ./sendgrid.env
```

<a name="quick-start"></a>
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
	const host = "api.example.com"
	param := "myparam"
	endpoint := "/your/api/" + param + "/call"

	baseURL := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   endpoint,
	}

	request, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
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
	const host = "api.example.com"
	param := "myparam"
	endpoint := "/your/api/" + param + "/call"

	key := os.Getenv("API_KEY")
	var body = []byte(`{"some": 0, "awesome": 1, "data": 3}`)

	params := url.Values{}
	params.Add("hello", "0")
	params.Add("world", "1")

	baseURL := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   endpoint,
		RawQuery: params.Encode(),
	}
	request, err = http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+key)
	request.Header.Set("X-Test", "Test")

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
}
```

<a name="usage"></a>
# Usage

- [Usage Examples](USAGE.md)

<a name="roadmap"></a>
# Roadmap

If you are interested in the future direction of this project, please take a look at our [milestones](https://github.com/sendgrid/rest/milestones). We would love to hear your feedback.

<a name="contribute"></a>
# How to Contribute

We encourage contribution to our projects, please see our [CONTRIBUTING](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md) guide for details.

Quick links:

- [Feature Request](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#feature-request)
- [Bug Reports](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#submit-a-bug-report)
- [Sign the CLA to Create a Pull Request](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#cla)
- [Improvements to the Codebase](https://github.com/sendgrid/rest/blob/master/CONTRIBUTING.md#improvements-to-the-codebase)

<a name="about"></a>
# About

rest is guided and supported by the SendGrid [Developer Experience Team](mailto:dx@sendgrid.com).

rest is maintained and funded by SendGrid, Inc. The names and logos for rest are trademarks of SendGrid, Inc.

<a name="license"></a>
# License
[The MIT License (MIT)](LICENSE.txt)
