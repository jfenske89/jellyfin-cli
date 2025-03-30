package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type JellyfinApiClient interface {
	// ListSessions return a list of active sessions
	ListSessions(context.Context, map[string]string) ([]Session, error)

	// ListLibraryFolders return a list of library virtual folders
	ListLibraryFolders(context.Context, map[string]string) ([]LibraryFolder, error)

	// ListActivityLogs return an object describing recent activity
	ListActivityLogs(ctx context.Context, getParameters map[string]string) (ActivityLog, error)

	// Search return a search response
	Search(ctx context.Context, term string, getParameters map[string]string) (SearchResponse, error)

	// RefreshLibrary initiates a library refresh
	RefreshLibrary(ctx context.Context) error
}

type jellyfinApiClientImpl struct {
	config     JellyfinApiConfig
	httpClient *http.Client
	logger     *zap.SugaredLogger
}

func NewJellyfinApiClient(config JellyfinApiConfig, logger *zap.SugaredLogger) JellyfinApiClient {
	return &jellyfinApiClientImpl{
		config: config,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.SkipSslVerify,
				},
			},
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (c *jellyfinApiClientImpl) ListSessions(ctx context.Context, getParameters map[string]string) ([]Session, error) {
	response, err := c.makeRequest(
		ctx,
		http.MethodGet,
		c.appendGetParameters("Sessions", getParameters),
		nil,
		nil,
	)

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return buildModels(jsonBytes, NewSession)
}

func (c *jellyfinApiClientImpl) ListLibraryFolders(ctx context.Context, getParameters map[string]string) ([]LibraryFolder, error) {
	response, err := c.makeRequest(
		ctx,
		http.MethodGet,
		c.appendGetParameters("Library/VirtualFolders", getParameters),
		nil,
		nil,
	)

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return buildModels(jsonBytes, NewLibraryFolder)
}

func (c *jellyfinApiClientImpl) ListActivityLogs(ctx context.Context, getParameters map[string]string) (ActivityLog, error) {
	response, err := c.makeRequest(
		ctx,
		http.MethodGet,
		c.appendGetParameters("System/ActivityLog/Entries", getParameters),
		nil,
		nil,
	)

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return buildModel(jsonBytes, NewActivityLog)
}

func (c *jellyfinApiClientImpl) Search(ctx context.Context, term string, getParameters map[string]string) (SearchResponse, error) {
	if getParameters == nil {
		getParameters = make(map[string]string)
	}
	getParameters["searchTerm"] = term

	response, err := c.makeRequest(
		ctx,
		http.MethodGet,
		c.appendGetParameters("Search/Hints", getParameters),
		nil,
		nil,
	)

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return buildModel(jsonBytes, NewSearchResponse)
}

func (c *jellyfinApiClientImpl) RefreshLibrary(ctx context.Context) error {
	response, err := c.makeRequest(
		ctx,
		http.MethodPost,
		"Library/Refresh",
		nil,
		nil,
	)

	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	return err
}

func (c *jellyfinApiClientImpl) appendGetParameters(endpoint string, getParameters map[string]string) string {
	result := endpoint

	for key, value := range getParameters {
		if result == endpoint {
			result += "?"
		} else {
			result += "&"
		}

		result += url.QueryEscape(key) + "=" + url.QueryEscape(value)
	}

	return result
}

func (c *jellyfinApiClientImpl) makeRequest(
	ctx context.Context,
	method string,
	endpoint string,
	headers map[string][]string,
	body []byte,
) (*http.Response, error) {
	if headers == nil {
		headers = make(map[string][]string)
	}

	defaultHeaders := map[string][]string{
		header_content_type: {"application/json"},
		header_emby_token:   {c.config.Token},
	}

	for key, val := range defaultHeaders {
		if _, ok := headers[key]; !ok {
			headers[key] = val
		}
	}

	url, err := c.buildUrl(endpoint)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		url.String(),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header = headers

	startedAt := time.Now()

	response, err := c.httpClient.Do(request)
	if err != nil {
		return response, fmt.Errorf("failed to make request: %w", err)
	}

	c.logger.Debugw(
		"executed request",
		"method", method,
		"url", url.String(),
		"duration", time.Since(startedAt).String(),
		"status-code", response.StatusCode,
		"content-length", response.ContentLength,
	)

	return response, nil
}

func (c *jellyfinApiClientImpl) buildUrl(endpoint string) (*url.URL, error) {
	baseUrl := c.config.BaseUrl
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	url, err := url.Parse(baseUrl + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL from %s: %w", baseUrl, err)
	}

	return url, nil
}
