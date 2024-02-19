package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type JellyfinApiClient interface {
	// ListSessions return a list of active sessions
	ListSessions(context.Context) ([]Session, error)
}

type jellyfinApiClientImpl struct {
	config     JellyfinApiConfig
	httpClient *http.Client
}

func NewJellyfinApiClient(config JellyfinApiConfig) JellyfinApiClient {
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
	}
}

func (c *jellyfinApiClientImpl) ListSessions(ctx context.Context) ([]Session, error) {
	if response, err := c.makeRequest(ctx, http.MethodGet, "Sessions", nil, nil); err != nil {
		return nil, err
	} else {
		defer response.Body.Close()

		var sessions []Session
		if jsonBytes, err := io.ReadAll(response.Body); err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		} else if err = json.Unmarshal(jsonBytes, &sessions); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		} else {
			return sessions, nil
		}
	}
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
	} else if request, err := http.NewRequest(method, url.String(), bytes.NewBuffer(body)); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	} else {
		request.Header = headers

		if response, err := c.httpClient.Do(request); err != nil {
			return response, fmt.Errorf("failed to make request: %w", err)
		} else {
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
		return nil, fmt.Errorf("failed to parse URL from %s: %w", os.Getenv("HOST"), err)
	} else {
		return url, nil
	}
}
