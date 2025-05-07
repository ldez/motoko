package main // import "motoko.test"

import (
	"context"
	"fmt"

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
