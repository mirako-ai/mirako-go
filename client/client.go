package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mirako-ai/mirako-go/gen"
)

const (
	DefaultBaseURL = "https://mirako.co"
	DefaultTimeout = 60 * time.Second
)

type Client struct {
	gen.ClientInterface
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	retryConfig *RetryConfig
	logger      Logger
	tracer      Tracer
}

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		logger: &noopLogger{},
		tracer: &noopTracer{},
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	genClient, err := gen.NewClient(
		c.baseURL,
		gen.WithHTTPClient(c.httpClient),
		gen.WithRequestEditorFn(c.authRequestEditor),
		gen.WithRequestEditorFn(c.loggingRequestEditor),
		gen.WithRequestEditorFn(c.tracingRequestEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create generated client: %w", err)
	}

	c.ClientInterface = genClient

	return c, nil
}

func (c *Client) authRequestEditor(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return nil
}

func (c *Client) loggingRequestEditor(ctx context.Context, req *http.Request) error {
	c.logger.Logf("Request: %s %s", req.Method, req.URL.String())
	return nil
}

func (c *Client) tracingRequestEditor(ctx context.Context, req *http.Request) error {
	c.tracer.TraceRequest(ctx, req)
	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.retryConfig == nil {
		return c.httpClient.Do(req)
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := c.retryConfig.calculateBackoff(attempt)
			c.logger.Logf("Retrying request after %v (attempt %d/%d)", backoff, attempt, c.retryConfig.MaxRetries)
			time.Sleep(backoff)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			continue
		}

		if !c.shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		if resp.Body != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}

	return resp, err
}

func (c *Client) shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode >= 500 && statusCode < 600
}

type Logger interface {
	Logf(format string, args ...any)
}

type Tracer interface {
	TraceRequest(ctx context.Context, req *http.Request)
	TraceResponse(ctx context.Context, resp *http.Response)
}

type noopLogger struct{}

func (l *noopLogger) Logf(format string, args ...any) {}

type noopTracer struct{}

func (t *noopTracer) TraceRequest(ctx context.Context, req *http.Request)    {}
func (t *noopTracer) TraceResponse(ctx context.Context, resp *http.Response) {}
