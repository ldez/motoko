package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	sampleMain20 = `package main // import "motoko.test"

import (
	"context"
	"fmt"

	"github.com/google/go-github/v20/github"
)

func main() {
	client := github.NewClient(nil)

	octocat, _, err := client.Octocat(context.Background(), "Go modules!")
	if err != nil {
		panic(err)
	}
	fmt.Println(octocat)
}
`

	sampleGoMod20 = `module motoko.test

go 1.15

// test
require (
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/go-github/v20 v20.0.0
)
`
)

func setupTestProject(t *testing.T) (string, error) {
	t.Helper()

	dir := t.TempDir()

	for _, filename := range []string{"main.go", "go.mod", "go.sum"} {
		err := copyFile(dir, filename)
		if err != nil {
			return "", err
		}
	}

	return dir, nil
}

func copyFile(dir, filename string) error {
	src, err := os.Open(filepath.FromSlash("./testdata/a/" + filename))
	if err != nil {
		return err
	}

	defer func() { _ = src.Close() }()

	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}

	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func quickDiff(got, want string) string {
	builder := &bytes.Buffer{}

	splitWant := strings.Split(want, "\n")
	splitGot := strings.Split(got, "\n")

	for i := 0; i < len(splitWant) && i < len(splitGot); i++ {
		v := splitWant[i]
		if splitGot[i] != v {
			fmt.Fprintf(builder, "Line %-4d got : %s\n", i, splitGot[i])
			fmt.Fprintf(builder, "Line %-4d want: %s\n", i, v)
		}
	}

	d := len(splitWant) - len(splitGot)
	if d > 0 {
		for i := len(splitWant) - d; i < len(splitWant); i++ {
			fmt.Fprintf(builder, "Line %-4d got : <nothing>\n", i)
			fmt.Fprintf(builder, "Line %-4d want: %s\n", i, splitWant[i])
		}
	} else if d < 0 {
		for i := len(splitGot) + d; i < len(splitGot); i++ {
			fmt.Fprintf(builder, "Line %-4d got : %s\n", i, splitGot[i])
			fmt.Fprintf(builder, "Line %-4d want: <nothing>\n", i)
		}
	}

	return builder.String()
}
