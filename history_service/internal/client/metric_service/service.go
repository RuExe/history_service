package metricservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"page/internal/domain"
)

type Client struct {
	URL        string
	httpClient http.Client
}

func NewClient(url string) *Client {
	return &Client{
		URL: url,
		httpClient: http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *Client) GetData() (*domain.MetricData, error) {
	resp, err := c.httpClient.Get(c.URL)
	if err != nil {
		return nil, fmt.Errorf("no response from request")
	}
	defer func() {
		closeBodyErr := resp.Body.Close()
		if err == nil && closeBodyErr != nil {
			err = closeBodyErr
		}
	}()

	result := new(domain.MetricData)
	if json.NewDecoder(resp.Body).Decode(result) != nil {
		return nil, fmt.Errorf("can not decode response body")
	}
	return result, nil
}
