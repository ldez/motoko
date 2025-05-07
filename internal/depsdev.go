package internal

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"

	"github.com/ldez/grignotin/goproxy"
)

// InsightVersions represents a subsets of field from devs.dep.
type InsightVersions struct {
	Version struct {
		Version         string `json:"version"`
		RelatedPackages struct {
			GoModuleSuggestedName *GoModuleMajorVersion  `json:"goModuleSuggestedName"`
			GoModuleMajorVersions []GoModuleMajorVersion `json:"goModuleMajorVersions"`
		} `json:"relatedPackages"`
	} `json:"version"`
}

// GoModuleMajorVersion represents a major version of a module.
type GoModuleMajorVersion struct {
	System string `json:"system"`
	Name   string `json:"name"`
}

// DepsDevClient deps.dev client.
type DepsDevClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

// NewDepsDevClient creates a new DepsDevClient.
func NewDepsDevClient() *DepsDevClient {
	baseURL, _ := url.Parse("https://deps.dev")

	return &DepsDevClient{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}
}

// GetModuleMajorVersions returns all major versions of a module.
func (c *DepsDevClient) GetModuleMajorVersions(ctx context.Context, modPath string) ([]GoModuleMajorVersion, error) {
	endpoint := c.baseURL.JoinPath("_/s/go/p/", url.PathEscape(modPath), "v/")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var result InsightVersions

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Version.RelatedPackages.GoModuleSuggestedName != nil {
		return []GoModuleMajorVersion{*result.Version.RelatedPackages.GoModuleSuggestedName}, nil
	}

	exp := regexp.MustCompile(`.+[/.]v(\d+)$`)

	slices.SortFunc(result.Version.RelatedPackages.GoModuleMajorVersions, func(a, b GoModuleMajorVersion) int {
		switch {
		case exp.MatchString(a.Name) && exp.MatchString(b.Name):
			aVal, _ := strconv.Atoi(exp.FindStringSubmatch(a.Name)[1])
			bVal, _ := strconv.Atoi(exp.FindStringSubmatch(b.Name)[1])

			return cmp.Compare(bVal, aVal)

		case exp.MatchString(a.Name):
			return -1

		case exp.MatchString(b.Name):
			return 1

		default:
			return cmp.Compare(a.Name, b.Name)
		}
	})

	return result.Version.RelatedPackages.GoModuleMajorVersions, nil
}

// FindHighestFromDepsDev finds the highest version of a module.
func FindHighestFromDepsDev(ctx context.Context, lib string) (string, error) {
	depsClient := NewDepsDevClient()

	versions, err := depsClient.GetModuleMajorVersions(ctx, lib)
	if err != nil {
		return "", fmt.Errorf("get module major versions (%s): %w", lib, err)
	}

	if len(versions) == 0 {
		return "", &MajorNotFoundError{}
	}

	client := goproxy.NewClient("")

	latest, err := client.GetLatest(versions[0].Name)
	if err != nil {
		return "", fmt.Errorf("get latest (%s:%s): %w", lib, versions[0].Name, err)
	}

	return latest.Version, nil
}
