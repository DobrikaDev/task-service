package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"DobrikaDev/task-service/utils/config"

	"go.uber.org/zap"
)

const (
	indexEndpoint  = "/index"
	searchEndpoint = "/search"
)

var (
	ErrInvalidBaseURL = errors.New("search client: invalid base url")
	ErrUnexpectedCode = errors.New("search client: unexpected response code")
)

type Client struct {
	baseURL       *url.URL
	httpClient    *http.Client
	indexTimeout  time.Duration
	searchTimeout time.Duration
	logger        *zap.Logger
}

type Option func(*Client)

func WithLogger(logger *zap.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

func New(cfg config.SearchConfig, httpClient *http.Client, opts ...Option) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("search client: base url is required")
	}

	parsed, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidBaseURL, err)
	}

	client := &Client{
		baseURL:       parsed,
		httpClient:    httpClient,
		indexTimeout:  cfg.IndexTimeout,
		searchTimeout: cfg.SearchTimeout,
		logger:        zap.NewNop(),
	}

	if client.indexTimeout <= 0 {
		client.indexTimeout = 3 * time.Second
	}
	if client.searchTimeout <= 0 {
		client.searchTimeout = 2 * time.Second
	}
	if client.httpClient == nil {
		client.httpClient = http.DefaultClient
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

type IndexTask struct {
	TaskName string `json:"task_name"`
	TaskDesc string `json:"task_desc"`
	GeoData  string `json:"geo_data,omitempty"`
	TaskID   string `json:"task_id"`
	TaskType string `json:"task_type,omitempty"`
}

func (c *Client) IndexTask(ctx context.Context, task IndexTask) error {
	if task.TaskID == "" {
		return errors.New("search client: task_id is required")
	}

	return c.doRequest(ctx, c.indexTimeout, indexEndpoint, task, nil)
}

type SearchRequest struct {
	UserQuery string   `json:"user_query"`
	GeoData   string   `json:"geo_data,omitempty"`
	QueryType string   `json:"query_type,omitempty"`
	UserTags  []string `json:"user_tags,omitempty"`
}

type SearchResponse struct {
	TaskIDs []string `json:"task_id"`
	Status  string   `json:"status"`
}

func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	if req.UserQuery == "" {
		return nil, errors.New("search client: user_query is required")
	}

	var response SearchResponse
	if err := c.doRequest(ctx, c.searchTimeout, searchEndpoint, req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) doRequest(ctx context.Context, timeout time.Duration, endpoint string, payload any, out any) error {
	buf := &bytes.Buffer{}
	if payload != nil {
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return fmt.Errorf("search client: encode payload: %w", err)
		}
	}

	reqURL := *c.baseURL
	reqURL.Path = path.Join(reqURL.Path, endpoint)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), buf)
	if err != nil {
		return fmt.Errorf("search client: create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(request.Context(), timeout)
		defer cancel()
		request = request.WithContext(ctx)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("search client: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.Warn("search request failed", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("%w: %d", ErrUnexpectedCode, resp.StatusCode)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("search client: decode response: %w", err)
		}
	}

	return nil
}
