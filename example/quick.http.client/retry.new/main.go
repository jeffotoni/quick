package main

import (
    "fmt"
    "log"
    "time"

    "github.com/jeffotoni/quick/http/client"
)

func main() {
    cC := client.New(
        client.WithRetry(client.RetryConfig{
            MaxRetries:   3,
            Delay:        1 * time.Second,
            UseBackoff:   false,
            Statuses:     []int{502, 503, 504, 403},
            FailoverURLs: []string{"http://backup1", "https://reqres.in/api/users", "https://httpbin.org/post"},
            EnableLog:    true,
        }),
        client.WithHeaders(map[string]string{
            "Authorization": "Bearer token",
        }),
    )

    // Perform the POST request
    resp, err := cC.Post("http://localhost:3000/v1/user", map[string]string{
        "name":  "Jefferson",
        "email": "jeff@example.com",
    })
    if err != nil {
        log.Fatal("Error making POST request:", err)
    }

    // Print the response body and status code
    fmt.Println("POST Response Status:", resp.StatusCode)
    fmt.Println("POST Response Body:", string(resp.Body))

}
