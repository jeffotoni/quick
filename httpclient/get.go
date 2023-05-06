package httpclient

import (
	"fmt"
	"io"
	"net/http"
)

func (hc *HttpClient) Get(url string) *ClientResponse {

	req, err := http.NewRequestWithContext(hc.Ctx, "GET", url, nil)

	if err != nil {
		return &ClientResponse{Error: err}
	}

	if err != nil {
		return &ClientResponse{Error: err}
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.ClientHttp.Do(req)
	if err != nil {
		fmt.Printf("Error %s", err)
		return &ClientResponse{Error: err}

	}
	defer resp.Body.Close()
	code := resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error %s", err)
		return &ClientResponse{StatusCode: code, Error: err}

	}

	return &ClientResponse{Body: body, StatusCode: code, Error: err}
}
