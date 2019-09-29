package main

import (
	"fmt"
	"go/format"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func getNewVersion(latest bool, lib string, version string) string {
	if !latest {
		return "v" + strings.TrimPrefix(version, "v")
	}

	split := strings.Split(lib, "/")
	raw, err := getLatestVersion(split[1], split[2])
	if err != nil {
		log.Fatal(err)
	}

	vParts := strings.Split(raw, ".")
	return "v" + strings.TrimPrefix(vParts[0], "v")
}

func update(dir string, lib string, newVersion string, onlyFilename bool) error {
	config := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:   dir,
		Tests: true,
	}

	pkgs, err := packages.Load(config, "./...")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		for _, syn := range p.Syntax {
			var rewritten bool
			for _, imp := range syn.Imports {
				trim := strings.Trim(imp.Path.Value, `"`)
				parts := strings.Split(trim, "/")
				if len(parts) >= 3 && strings.Join(parts[:3], "/") == lib {
					newImp := createNewImport(parts, newVersion)
					if astutil.RewriteImport(p.Fset, syn, trim, newImp) {
						rewritten = true
					}
				}
			}

			if !rewritten {
				continue
			}

			goFileName := p.Fset.File(syn.Pos()).Name()
			if onlyFilename {
				fmt.Printf("%s: %s\n", lib, goFileName)
				return nil
			}

			f, err := os.Create(goFileName)
			if err != nil {
				return fmt.Errorf("could not create go file %s: %v", goFileName, err)
			}

			err = format.Node(f, p.Fset, syn)
			_ = f.Close()
			if err != nil {
				return fmt.Errorf("could not rewrite go file %s: %v", goFileName, err)
			}
		}
	}

	return nil
}

func createNewImport(parts []string, newVersion string) string {
	if len(parts) < 3 {
		panic(fmt.Sprintf("unsupported package format: %s", strings.Join(parts, "/")))
	}

	np := make([]string, 3)
	copy(np, parts[:3])
	np = append(np, newVersion)

	if len(parts) == 3 {
		// no version
		return strings.Join(np, "/")
	}

	if ok, _ := regexp.MatchString(`v\d+`, parts[3]); ok {
		if len(parts) > 4 {
			// version + sub-package
			np = append(np, parts[4:]...)
		}
	} else {
		// no version + sub-package
		np = append(np, parts[3:]...)
	}

	return strings.Join(np, "/")
}

func getLatestVersion(owner string, repo string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 10 * time.Second,
	}

	uri := fmt.Sprintf("https://github.com/%s/%s/releases/latest", owner, repo)

	resp, err := client.Get(uri)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("unable to find latest release URL: %d", resp.StatusCode)
	}

	u, err := resp.Location()
	if err != nil {
		return "", err
	}

	return path.Base(u.String()), nil
}
