package pkgsite

import (
	"fmt"
	"strings"
	"time"
)

// APIError contains detailed information about an error.
type APIError struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Fixes      []string    `json:"fixes"`
	Candidates []Candidate `json:"candidates,omitempty"`
}

func (e *APIError) Error() string {
	msg := new(strings.Builder)

	_, _ = fmt.Fprintf(msg, "%d: %s", e.Code, e.Message)

	for _, fix := range e.Fixes {
		_, _ = fmt.Fprintf(msg, " [fix:  %s]", fix)
	}

	for _, candidate := range e.Candidates {
		_, _ = fmt.Fprintf(msg, " [candidate: %s, %s]", candidate.ModulePath, candidate.PackagePath)
	}

	return msg.String()
}

// A Candidate is a potential resolution for an ambiguous path.
type Candidate struct {
	ModulePath  string `json:"modulePath"`
	PackagePath string `json:"packagePath"`
}

// PaginatedResponse is a generic paginated response.
type PaginatedResponse[T any] struct {
	Items         []T    `json:"items"`
	Total         int    `json:"total"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// ModuleVersion is the response for /v1beta/versions/{path}.
type ModuleVersion struct {
	ModulePath        string    `json:"modulePath"`
	Version           string    `json:"version"`
	CommitTime        time.Time `json:"commitTime"`
	IsRedistributable bool      `json:"isRedistributable"`
	HasGoMod          bool      `json:"hasGoMod"`
	LatestVersion     string    `json:"latestVersion"`
	Deprecated        bool      `json:"deprecated"`
	DeprecationReason string    `json:"deprecationReason"`
	Retracted         bool      `json:"retracted"`
	RetractionReason  string    `json:"retractionReason"`
}

// ListParams are common pagination and filtering parameters.
type ListParams struct {
	Limit  int    `url:"limit,omitempty"`
	Token  string `url:"token,omitempty"`
	Filter string `url:"filter,omitempty"`
}

// ModuleParams are query parameters for /v1beta/module/{path}.
type ModuleParams struct {
	Version  string `url:"version,omitempty"`
	Licenses bool   `url:"licenses,omitempty"`
	Readme   bool   `url:"readme,omitempty"`
}

// Module is the response for /v1beta/module/{modulePath}.
type Module struct {
	Path              string    `json:"path"`
	Version           string    `json:"version"`
	CommitTime        time.Time `json:"commitTime"`
	IsLatest          bool      `json:"isLatest"`
	IsRedistributable bool      `json:"isRedistributable"`
	IsStandardLibrary bool      `json:"isStandardLibrary"`
	HasGoMod          bool      `json:"hasGoMod"`
	RepoURL           string    `json:"repoUrl"`
	GoModContents     string    `json:"goModContents,omitempty"`
	Readme            *Readme   `json:"readme,omitempty"`
	Licenses          []License `json:"licenses,omitempty"`
}

// Readme is a readme file.
type Readme struct {
	Filepath string `json:"filepath"`
	Contents string `json:"contents"`
}

// License is license information in API responses.
type License struct {
	Types    []string `json:"types"`
	FilePath string   `json:"filePath"`
	Contents string   `json:"contents,omitempty"`
}
