package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/ldez/grignotin/goenv"
	"github.com/ldez/motoko/internal/pkgsite"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

const maxWorkers = 10

type updateModule struct {
	OriginPath    string
	NewPath       string
	LatestVersion string
}

func findCmd(ctx context.Context) error {
	goModPath, err := goenv.GetOne(ctx, goenv.GOMOD)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Clean(goModPath))
	if err != nil {
		return fmt.Errorf("reading go.mod: %w", err)
	}

	mod, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return err
	}

	worker := Worker{client: pkgsite.NewClient(nil)}

	resultsChan := make(chan updateModule)
	requiresChan := make(chan *modfile.Require, maxWorkers)

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(maxWorkers)

		for range maxWorkers {
			go func() {
				worker.find(ctx, requiresChan, resultsChan)
				wg.Done()
			}()
		}

		wg.Wait()
		close(resultsChan)
	}()

	go func() {
		for _, require := range mod.Require {
			if require.Indirect {
				continue
			}

			requiresChan <- require
		}

		close(requiresChan)
	}()

	for result := range resultsChan {
		fmt.Printf("%s: %s (%s)\n", result.OriginPath, result.NewPath, result.LatestVersion)
	}

	return nil
}

type Worker struct {
	client *pkgsite.Client
}

func (w *Worker) find(ctx context.Context, inChan <-chan *modfile.Require, outChan chan<- updateModule) {
	exp := regexp.MustCompile(`(.+)([/.])v\d+$`)

	for require := range inChan {
		modPath := require.Mod.Path

		versions, err := w.client.Versions(ctx, modPath, &pkgsite.ListParams{Limit: 1})
		if err != nil {
			mnf, ok := errors.AsType[*pkgsite.APIError](err)
			if !ok || mnf.Code != http.StatusNotFound {
				if _, ok := os.LookupEnv("MOTOKO_DEBUG"); ok {
					log.Println(modPath, err)
				}
			}

			continue
		}

		if len(versions.Items) != 1 {
			continue
		}

		latestVersion := versions.Items[0].Version

		major := semver.Major(latestVersion)

		if major == "v0" || major == "v1" {
			continue
		}

		newVersion := path.Join(modPath, major)

		if exp.MatchString(require.Mod.Path) {
			newVersion = exp.FindStringSubmatch(modPath)[1] + exp.FindStringSubmatch(modPath)[2] + semver.Major(latestVersion)
		}

		if newVersion == require.Mod.Path {
			continue
		}

		outChan <- updateModule{OriginPath: require.Mod.Path, NewPath: newVersion, LatestVersion: latestVersion}
	}
}
