package httpclient

func Get(url string) *ClientResponse {
	return defaultClient.Get(url)
}

func (c *Client) Get(url string) *ClientResponse {
	return c.createRequest(url, "GET", nil)
}
