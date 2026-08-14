package pkgsite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	querystring "github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
)

const defaultBaseURL = "https://pkg.go.dev/v1beta/"

// Client is a pkg.go.dev client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewClient creates a new Client.
func NewClient(httpClient *http.Client) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	if httpClient != nil {
		retryClient.HTTPClient = httpClient
	} else {
		retryClient.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}

	retryClient.Logger = slog.New(slog.DiscardHandler)

	c := &Client{
		httpClient: retryClient.StandardClient(),
		baseURL:    baseURL,
	}

	if httpClient != nil {
		c.httpClient = httpClient
	}

	return c
}

// Module information about the module at {path}.
func (c *Client) Module(ctx context.Context, mPath string, params *ModuleParams) (*Module, error) {
	endpoint := c.baseURL.JoinPath("module", mPath)

	values, err := querystring.Values(params)
	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = values.Encode()

	result := new(Module)

	err = c.do(ctx, endpoint, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Versions versions of the module at {path}.
func (c *Client) Versions(ctx context.Context, mPath string, params *ListParams) (*PaginatedResponse[ModuleVersion], error) {
	endpoint := c.baseURL.JoinPath("versions", mPath)

	values, err := querystring.Values(params)
	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = values.Encode()

	result := new(PaginatedResponse[ModuleVersion])

	err = c.do(ctx, endpoint, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Latest gets the latest version of a module.
func (c *Client) Latest(ctx context.Context, mPath string) (string, error) {
	versions, err := c.Versions(ctx, mPath, &ListParams{Limit: 1})
	if err != nil {
		return "", fmt.Errorf("versions: %w", err)
	}

	if len(versions.Items) == 0 {
		return "", fmt.Errorf("no versions found for %s", mPath)
	}

	return versions.Items[0].LatestVersion, nil
}

func (c *Client) do(ctx context.Context, endpoint *url.URL, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return parseError(resp)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	err = json.Unmarshal(raw, result)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

func parseError(resp *http.Response) error {
	raw, _ := io.ReadAll(resp.Body)

	errAPI := new(APIError)

	err := json.Unmarshal(raw, errAPI)
	if err != nil {
		return fmt.Errorf("%d: %s", resp.StatusCode, string(raw))
	}

	return errAPI
}
