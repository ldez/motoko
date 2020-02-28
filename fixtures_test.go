package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	sampleMain = `package main // import "motoko.test"

import (
    "fmt"
    "context"

    "github.com/google/go-github/github"
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

	sampleGoMod = `module motoko.test

require (
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
)
`

	sampleGoMod20 = `package main // import "motoko.test"

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
)

func setupTestProject(t *testing.T) (string, error) {
	dir, err := ioutil.TempDir("", "motoko")
	if err != nil {
		return "", err
	}

	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	if ioutil.WriteFile(filepath.Join(dir, "main.go"), []byte(sampleMain), 0644) != nil {
		return "", err
	}

	if ioutil.WriteFile(filepath.Join(dir, "go.mod"), []byte(sampleGoMod), 0644) != nil {
		return "", err
	}

	return dir, nil
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
