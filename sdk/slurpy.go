package slurpy

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bobby/slurpy/pkg/models"
	"github.com/bobby/slurpy/pkg/storage"
)

// Client wraps http.Client with logging capabilities
type Client struct {
	*http.Client
	namespace string
	enabled   bool
	storage   *storage.Storage
}

// Config holds configuration for the Slurpy client
type Config struct {
	Namespace string // Unique identifier for this project/program
	Enabled   bool   // Whether to enable request logging
}

// New creates a new Slurpy client
func New(config Config) (*Client, error) {
	if config.Namespace == "" {
		config.Namespace = "default"
	}

	var store *storage.Storage
	var err error

	if config.Enabled {
		store, err = storage.New()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage: %w", err)
		}
	}

	return &Client{
		Client:    &http.Client{},
		namespace: config.Namespace,
		enabled:   config.Enabled,
		storage:   store,
	}, nil
}

// Do executes an HTTP request with logging
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if !c.enabled {
		return c.Client.Do(req)
	}

	startTime := time.Now()
	reqID := generateID()

	// Create logged request
	loggedReq := &models.LoggedRequest{
		ID:        reqID,
		Timestamp: startTime,
		Method:    req.Method,
		URL:       req.URL.String(),
		Headers:   models.HeadersFromHTTP(req.Header),
		Namespace: c.namespace,
	}

	// Capture request body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			loggedReq.Body = string(bodyBytes)
			// Restore body for the actual request
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}

	// Execute the request
	resp, err := c.Client.Do(req)
	duration := time.Since(startTime)
	loggedReq.Duration = duration

	if err != nil {
		loggedReq.Error = err.Error()
	} else {
		// Capture response
		loggedResp := &models.LoggedResponse{
			StatusCode: resp.StatusCode,
			Headers:    models.HeadersFromHTTP(resp.Header),
		}

		// Capture response body
		if resp.Body != nil {
			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr == nil {
				loggedResp.Body = string(bodyBytes)
				loggedResp.Size = int64(len(bodyBytes))
				// Restore body for the caller
				resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		loggedReq.Response = loggedResp
	}

	// Save the logged request
	if saveErr := c.storage.SaveRequest(loggedReq); saveErr != nil {
		// Don't fail the original request if logging fails
		fmt.Printf("Warning: failed to save request log: %v\n", saveErr)
	}

	return resp, err
}

// Get executes a GET request with logging
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post executes a POST request with logging
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.Do(req)
}

// Put executes a PUT request with logging
func (c *Client) Put(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.Do(req)
}

// Delete executes a DELETE request with logging
func (c *Client) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// SetNamespace updates the namespace for future requests
func (c *Client) SetNamespace(namespace string) {
	c.namespace = namespace
}

// GetNamespace returns the current namespace
func (c *Client) GetNamespace() string {
	return c.namespace
}

// SetEnabled enables or disables request logging
func (c *Client) SetEnabled(enabled bool) error {
	if enabled && c.storage == nil {
		store, err := storage.New()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		c.storage = store
	}
	c.enabled = enabled
	return nil
}

// IsEnabled returns whether request logging is enabled
func (c *Client) IsEnabled() bool {
	return c.enabled
}

// generateID creates a random hex ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// WrapDefaultClient replaces the default http client methods with logged versions
func WrapDefaultClient(config Config) error {
	client, err := New(config)
	if err != nil {
		return err
	}

	// Override package-level functions
	http.DefaultClient = client.Client
	return nil
}
