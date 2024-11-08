package oreillyapi

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	APIEndpoint = "https://learning.oreilly.com/api/v2/search"
)

func New() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

type Client struct {
	httpClient *http.Client
}

func (c *Client) Search(opt SearchOption) (*SearchResponse, error) {
	u, err := url.Parse(APIEndpoint)
	if err != nil {
		return nil, err
	}
	u.RawQuery = opt.queryParams().Encode()

	req := &http.Request{
		Method: http.MethodGet,
		URL:    u,
	}

	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)

	var res SearchResponse
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
