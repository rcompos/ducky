package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {

    host := "https://hubtest.com"

    url := host + "/api/tokens/authenticate"

    req, _ := http.NewRequest("POST", url, nil)

    // Add 'token ' before inserting the token from the Hub
    req.Header.Add("authorization", "token ")
    req.Header.Add("cache-control", "no-cache")

    res, _ := http.DefaultClient.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    fmt.Println(res)
    fmt.Println(string(body))

}
