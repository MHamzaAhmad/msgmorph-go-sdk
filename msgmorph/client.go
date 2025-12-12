// Package msgmorph provides a Go SDK for the MsgMorph API.
//
// The MsgMorph SDK enables you to manage contacts and feedback collection
// from your Go applications.
//
// # Installation
//
//	go get github.com/MHamzaAhmad/msgmorph-go-sdk
//
// # Quick Start
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//	    "log"
//	    "os"
//
//	    msgmorph "github.com/MHamzaAhmad/msgmorph-go-sdk/msgmorph"
//	)
//
//	func main() {
//	    client := msgmorph.NewClient(
//	        os.Getenv("MSGMORPH_API_KEY"),
//	        os.Getenv("MSGMORPH_ORGANIZATION_ID"),
//	    )
//
//	    contact, err := client.Contacts.Create(context.Background(), msgmorph.CreateContactInput{
//	        ExternalID: "user-123",
//	        Email:      "user@example.com",
//	        ProjectID:  os.Getenv("MSGMORPH_PROJECT_ID"),
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Created contact: %s\n", contact.ID)
//	}
package msgmorph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is the default MsgMorph API URL.
// Change this to your local API URL during development (e.g., "http://localhost:3001").
const DefaultBaseURL = "https://api.msgmorph.com/"

// DefaultTimeout is the default HTTP client timeout.
const DefaultTimeout = 30 * time.Second

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API.
//
// Example:
//
//	client := msgmorph.NewClient(apiKey, orgID,
//	    msgmorph.WithBaseURL("http://localhost:3001"),
//	)
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
//
// Example:
//
//	httpClient := &http.Client{
//	    Timeout: 60 * time.Second,
//	}
//	client := msgmorph.NewClient(apiKey, orgID,
//	    msgmorph.WithHTTPClient(httpClient),
//	)
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout.
//
// Example:
//
//	client := msgmorph.NewClient(apiKey, orgID,
//	    msgmorph.WithTimeout(60 * time.Second),
//	)
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// Client is the MsgMorph API client.
//
// Use NewClient to create a new client instance:
//
//	client := msgmorph.NewClient(
//	    os.Getenv("MSGMORPH_API_KEY"),
//	    os.Getenv("MSGMORPH_ORGANIZATION_ID"),
//	)
//
// The client provides access to various API resources through its fields:
//
//	// Create a contact
//	contact, err := client.Contacts.Create(ctx, input)
type Client struct {
	// apiKey is the MsgMorph API key for authentication.
	apiKey string

	// organizationID is the MsgMorph organization ID.
	organizationID string

	// baseURL is the API base URL.
	baseURL string

	// httpClient is the underlying HTTP client.
	httpClient *http.Client

	// Contacts provides access to contact management operations.
	Contacts *ContactsResource
}

// NewClient creates a new MsgMorph API client.
//
// Parameters:
//   - apiKey: Your MsgMorph API key (required)
//   - organizationID: Your MsgMorph organization ID (required)
//   - opts: Optional configuration options
//
// Example:
//
//	// Basic initialization
//	client := msgmorph.NewClient(
//	    os.Getenv("MSGMORPH_API_KEY"),
//	    os.Getenv("MSGMORPH_ORGANIZATION_ID"),
//	)
//
//	// With custom options
//	client := msgmorph.NewClient(
//	    apiKey,
//	    orgID,
//	    msgmorph.WithBaseURL("http://localhost:3001"),
//	    msgmorph.WithTimeout(60 * time.Second),
//	)
//
// Returns a configured Client ready to make API calls.
// Panics if apiKey or organizationID is empty.
func NewClient(apiKey, organizationID string, opts ...ClientOption) *Client {
	if apiKey == "" {
		panic(newError(
			"API key is required. Set the MSGMORPH_API_KEY environment variable.",
			400,
			ErrInvalidAPIKey,
			nil,
		))
	}
	if organizationID == "" {
		panic(newError(
			"Organization ID is required. Set the MSGMORPH_ORGANIZATION_ID environment variable.",
			400,
			ErrInvalidOrganizationID,
			nil,
		))
	}

	c := &Client{
		apiKey:         apiKey,
		organizationID: organizationID,
		baseURL:        DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize resources
	c.Contacts = &ContactsResource{client: c}

	return c
}

// request makes an authenticated HTTP request to the MsgMorph API.
// This is an internal method used by resource methods.
func (c *Client) request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil && method != http.MethodGet {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return newError(fmt.Sprintf("failed to marshal request body: %v", err), 0, ErrValidationError, nil)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return newNetworkError(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("X-Organization-Id", c.organizationID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return newNetworkError(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return newNetworkError(err)
	}

	if resp.StatusCode >= 400 {
		return parseErrorResponse(respBody, resp.StatusCode)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return newError(fmt.Sprintf("failed to parse response: %v", err), resp.StatusCode, ErrInternalError, nil)
		}
	}

	return nil
}

// parseErrorResponse parses an error response from the API.
func parseErrorResponse(body []byte, status int) *Error {
	var errResp struct {
		Message string                 `json:"message"`
		Error   string                 `json:"error"`
		Code    string                 `json:"code"`
		Details map[string]interface{} `json:"details"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return newError("An unexpected error occurred", status, errorCodeFromStatus(status), nil)
	}

	message := errResp.Message
	if message == "" {
		message = errResp.Error
	}
	if message == "" {
		message = "An unexpected error occurred"
	}

	code := ErrorCode(errResp.Code)
	if code == "" {
		code = errorCodeFromStatus(status)
	}

	return newError(message, status, code, errResp.Details)
}
