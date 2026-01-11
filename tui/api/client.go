package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tact-tui/model"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) FetchTimeCodes() ([]model.TimeCode, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/time-codes")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch time codes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var timeCodes []model.TimeCode
	if err := json.NewDecoder(resp.Body).Decode(&timeCodes); err != nil {
		return nil, fmt.Errorf("failed to decode time codes: %w", err)
	}

	return timeCodes, nil
}

func (c *Client) FetchWorkTypes() ([]model.WorkType, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/work-types")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch work types: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var workTypes []model.WorkType
	if err := json.NewDecoder(resp.Body).Decode(&workTypes); err != nil {
		return nil, fmt.Errorf("failed to decode work types: %w", err)
	}

	return workTypes, nil
}
