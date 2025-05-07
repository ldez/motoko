package internal

import (
	"fmt"
	"net/http"
	"path"

	"github.com/andybalholm/cascadia"
	"github.com/ldez/grignotin/goproxy"
	"golang.org/x/net/html"
)

// FindHighestFromGoPkg finds the highest version of a module.
func FindHighestFromGoPkg(lib string) (string, error) {
	licenseURL := fmt.Sprintf("https://pkg.go.dev/%s", lib)

	req, err := http.NewRequest(http.MethodGet, licenseURL, http.NoBody)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	fmt.Println(resp.StatusCode)

	defer func() { _ = resp.Body.Close() }()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	compile := cascadia.MustCompile("div.go-Main-banner div.go-Message.go-Message--notice a")

	node := cascadia.Query(doc, compile)
	if node != nil && node.FirstChild != nil {
		client := goproxy.NewClient("")

		latest, err := client.GetLatest(path.Join(lib, node.FirstChild.Data))
		if err != nil {
			return "", err
		}

		return latest.Version, nil
	}

	return "", &MajorNotFoundError{}
}

// MajorNotFoundError returned when there are no major versions.
type MajorNotFoundError struct{}

func (m *MajorNotFoundError) Error() string {
	return "highest major version not found"
}
