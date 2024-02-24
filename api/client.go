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
	if response, err := c.makeRequest(
		ctx,
		http.MethodGet,
		c.appendGetParameters("Sessions", getParameters),
		nil,
		nil,
	); err != nil {
		return nil, err
	} else {
		defer response.Body.Close()

		if jsonBytes, err := io.ReadAll(response.Body); err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		} else {
			return buildModels[Session](jsonBytes, NewSession)
		}
	}
}

func (c *jellyfinApiClientImpl) ListLibraryFolders(ctx context.Context, getParameters map[string]string) ([]LibraryFolder, error) {
	if response, err := c.makeRequest(
		ctx,
		http.MethodGet,
		c.appendGetParameters("Library/VirtualFolders", getParameters),
		nil,
		nil,
	); err != nil {
		return nil, err
	} else {
		defer response.Body.Close()

		if jsonBytes, err := io.ReadAll(response.Body); err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		} else {
			return buildModels[LibraryFolder](jsonBytes, NewLibraryFolder)
		}
	}
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

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = []string{"application/json"}
	}

	if _, ok := headers["X-Emby-Token"]; !ok {
		headers["X-Emby-Token"] = []string{c.config.Token}
	}

	if url, err := c.buildUrl(endpoint); err != nil {
		return nil, err
	} else if request, err := http.NewRequestWithContext(
		ctx,
		method,
		url.String(),
		bytes.NewBuffer(body),
	); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	} else {
		request.Header = headers

		startedAt := time.Now()

		if response, err := c.httpClient.Do(request); err != nil {
			return response, fmt.Errorf("failed to make request: %w", err)
		} else {
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
	}
}

func (c *jellyfinApiClientImpl) buildUrl(endpoint string) (*url.URL, error) {
	baseUrl := c.config.BaseUrl
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	baseUrl += endpoint

	if url, err := url.Parse(baseUrl); err != nil {
		return nil, fmt.Errorf("failed to parse URL from %s: %w", baseUrl, err)
	} else {
		return url, nil
	}
}
