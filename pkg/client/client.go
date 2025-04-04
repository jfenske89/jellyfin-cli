package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/jfenske89/jellyfin-cli/pkg/models"
)

// Client defines the interface for interacting with the Jellyfin API
type Client interface {
	// ListSessions returns a list of active sessions
	ListSessions(ctx context.Context, params map[string]string) ([]models.Session, error)

	// ListLibraryFolders returns a list of library virtual folders
	ListLibraryFolders(ctx context.Context, params map[string]string) ([]models.LibraryFolder, error)

	// ListActivityLogs returns recent activity
	ListActivityLogs(ctx context.Context, params map[string]string) (*models.ActivityLog, error)

	// Search returns search results
	Search(ctx context.Context, term string, params map[string]string) (*models.SearchResponse, error)

	// RefreshLibrary initiates a library refresh
	RefreshLibrary(ctx context.Context) error
}

// JellyfinClient is the implementation of the Client interface
type JellyfinClient struct {
	config     models.JellyfinConfig
	httpClient *http.Client
	logger     *zap.SugaredLogger
}

// NewClient creates a new Jellyfin API client
func NewClient(config models.JellyfinConfig, logger *zap.SugaredLogger) Client {
	return &JellyfinClient{
		config: config,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.SkipSSLVerify,
				},
			},
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// ListSessions retrieves active sessions from the Jellyfin server
func (c *JellyfinClient) ListSessions(ctx context.Context, params map[string]string) ([]models.Session, error) {
	var sessions []models.Session

	err := c.doRequest(ctx, http.MethodGet, "Sessions", params, nil, &sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

// ListLibraryFolders retrieves library folders from the Jellyfin server
func (c *JellyfinClient) ListLibraryFolders(ctx context.Context, params map[string]string) ([]models.LibraryFolder, error) {
	var folders []models.LibraryFolder

	err := c.doRequest(ctx, http.MethodGet, "Library/VirtualFolders", params, nil, &folders)
	if err != nil {
		return nil, fmt.Errorf("failed to list library folders: %w", err)
	}

	return folders, nil
}

// ListActivityLogs retrieves activity logs from the Jellyfin server
func (c *JellyfinClient) ListActivityLogs(ctx context.Context, params map[string]string) (*models.ActivityLog, error) {
	var logs models.ActivityLog

	err := c.doRequest(ctx, http.MethodGet, "System/ActivityLog/Entries", params, nil, &logs)
	if err != nil {
		return nil, fmt.Errorf("failed to list activity logs: %w", err)
	}

	return &logs, nil
}

// Search searches the Jellyfin server for content
func (c *JellyfinClient) Search(ctx context.Context, term string, params map[string]string) (*models.SearchResponse, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["searchTerm"] = term

	var response models.SearchResponse

	err := c.doRequest(ctx, http.MethodGet, "Search/Hints", params, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	return &response, nil
}

// RefreshLibrary initiates a library refresh on the Jellyfin server
func (c *JellyfinClient) RefreshLibrary(ctx context.Context) error {
	err := c.doRequest(ctx, http.MethodPost, "Library/Refresh", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to refresh library: %w", err)
	}

	return nil
}

// doRequest handles the HTTP request to the Jellyfin API
func (c *JellyfinClient) doRequest(
	ctx context.Context,
	method string,
	endpoint string,
	params map[string]string,
	body interface{},
	result interface{},
) error {
	// Build the URL with any query parameters
	fullEndpoint := c.appendQueryParams(endpoint, params)

	// Get the full URL for the request
	reqURL, err := c.buildURL(fullEndpoint)
	if err != nil {
		return err
	}

	// Marshal the body if present
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Emby-Token", c.config.Token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Warnw("failed to close response body", "error", err)
		}
	}()

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// If no result is expected, return
	if result == nil {
		return nil
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// appendQueryParams appends query parameters to an endpoint
func (c *JellyfinClient) appendQueryParams(endpoint string, params map[string]string) string {
	if len(params) == 0 {
		return endpoint
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return fmt.Sprintf("%s?%s", endpoint, values.Encode())
}

// buildURL builds the full URL for an API request
func (c *JellyfinClient) buildURL(endpoint string) (*url.URL, error) {
	baseURL, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	apiPath := "/emby/"
	if strings.HasSuffix(baseURL.Path, "/") {
		apiPath = "emby/"
	}

	fullURL, err := baseURL.Parse(fmt.Sprintf("%s%s", apiPath, endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	return fullURL, nil
}
