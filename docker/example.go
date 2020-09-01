package main

import "github.com/sendgrid/rest"
import "fmt"

func main() {
    const host = "https://httpbin.org"
    param := "get"
    endpoint := "/" + param
    baseURL := host + endpoint
    method := rest.Get
    request := rest.Request{
        Method:  method,
        BaseURL: baseURL,
    }
    response, err := rest.Send(request)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(response.StatusCode)
        fmt.Println(response.Body)
        fmt.Println(response.Headers)
    }
}
