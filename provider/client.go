package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	json "github.com/goccy/go-json"
	"github.com/rs/zerolog"
)

// Client for chat provider
type Client struct {
	Token string

	client *http.Client
}

// NewClient for chat provider
func NewClient(token string) *Client {
	return &Client{
		Token: token,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
			Timeout: 5 * 60 * time.Second,
		},
	}
}

// makeRequest creates and executes a HTTP request
func (c *Client) makeRequest(ctx context.Context, method, path string, body io.Reader, logger zerolog.Logger) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://nano-gpt.com/api%s", path), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	// Add API token header
	if strings.HasPrefix(path, "/v1") {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	} else {
		req.Header.Add("x-api-key", c.Token)
	}

	logger.Debug().Msgf("Performing HTTP %s - %s", method, req.URL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("status code is %d: %s", resp.StatusCode, respBody)
	}

	logger.Debug().Msgf("Received status code %d from %s - %s: %s", resp.StatusCode, method, req.URL, respBody)

	return respBody, nil
}

// GetModels - Get all available models
func (c *Client) GetModels(ctx context.Context, logger zerolog.Logger) (*GetModelsResponse, error) {
	response, err := c.makeRequest(ctx, http.MethodGet, "/v1/models", nil, logger)
	if err != nil {
		return nil, err
	}

	var result GetModelsResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &result, nil
}

// PostCheckBalance - Check balance
func (c *Client) PostCheckBalance(ctx context.Context, logger zerolog.Logger) (*PostCheckBalanceResponse, error) {
	response, err := c.makeRequest(ctx, http.MethodPost, "/check-balance", nil, logger)
	if err != nil {
		return nil, err
	}

	var result PostCheckBalanceResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &result, nil
}

// PostCheckBalance - Check balance
func (c *Client) PostChatCompletions(ctx context.Context, logger zerolog.Logger, body PostChatCompletionsRequest) (*PostChatCompletionsResponse, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marhsal json: %w", err)
	}

	// fmt.Println(string(b))

	response, err := c.makeRequest(ctx, http.MethodPost, "/v1/chat/completions", bytes.NewBuffer(b), logger)
	if err != nil {
		return nil, err
	}

	var result PostChatCompletionsResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &result, nil
}
