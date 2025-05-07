package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/ldez/grignotin/goenv"
	"github.com/ldez/motoko/internal"
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

	resultsChan := make(chan updateModule)
	requiresChan := make(chan *modfile.Require, maxWorkers)

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(maxWorkers)

		for range maxWorkers {
			go func() {
				workerFind(ctx, requiresChan, resultsChan)
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

func workerFind(ctx context.Context, inChan <-chan *modfile.Require, outChan chan<- updateModule) {
	exp := regexp.MustCompile(`(.+)([/.])v\d+$`)

	for require := range inChan {
		modPath := require.Mod.Path

		latestVersion, err := internal.FindHighestFromDepsDev(ctx, modPath)
		if err != nil {
			var mnf *internal.MajorNotFoundError
			if !errors.As(err, &mnf) {
				if _, ok := os.LookupEnv("MOTOKO_DEBUG"); ok {
					log.Println(err)
				}
			}

			continue
		}

		newVersion := path.Join(modPath, semver.Major(latestVersion))

		if exp.MatchString(require.Mod.Path) {
			newVersion = exp.FindStringSubmatch(modPath)[1] + exp.FindStringSubmatch(modPath)[2] + semver.Major(latestVersion)
		}

		if newVersion == require.Mod.Path {
			continue
		}

		outChan <- updateModule{OriginPath: require.Mod.Path, NewPath: newVersion, LatestVersion: latestVersion}
	}
}
