package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * 10, // Timeout after 10 seconds
		},
	}
}

type CoinGeckoResponse struct {
	Prices [][]float64 `json:"prices"`
}

func (c *Client) FetchPriceHistories(ctx context.Context, symbol string, startDate time.Time, endDate time.Time) (*CoinGeckoResponse, error) {
	coinAPIUrl := os.Getenv("COINGECKO_API_URL")
	url := fmt.Sprintf("%s/coins/%s/market_chart/range?vs_currency=usd&from=%d&to=%d", coinAPIUrl, symbol, startDate.Unix(), endDate.Unix())
	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Make the HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println(url)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var cgResp CoinGeckoResponse
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &cgResp, nil
}
