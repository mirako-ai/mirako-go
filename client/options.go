package client

import (
	"math"
	"net/http"
	"time"
)

type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

func WithBearerToken(token string) Option {
	return WithAPIKey(token)
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithRetry(config RetryConfig) Option {
	return func(c *Client) {
		c.retryConfig = &config
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithTracer(tracer Tracer) Option {
	return func(c *Client) {
		c.tracer = tracer
	}
}

type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		BackoffFactor:  2.0,
	}
}

func (r *RetryConfig) calculateBackoff(attempt int) time.Duration {
	backoff := float64(r.InitialBackoff) * math.Pow(r.BackoffFactor, float64(attempt-1))
	if backoff > float64(r.MaxBackoff) {
		backoff = float64(r.MaxBackoff)
	}
	return time.Duration(backoff)
}
