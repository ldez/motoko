package main

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ldez/grignotin/goproxy"
	"github.com/ldez/motoko/internal"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func updatePackages(dir, lib, newVersion string, onlyFilename bool) error {
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedCompiledGoFiles | packages.NeedImports |
			packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:   dir,
		Tests: true,
	}

	pkgs, err := packages.Load(config, "./...")
	if err != nil {
		return err
	}

	for _, p := range pkgs {
		for _, syn := range p.Syntax {
			var rewritten bool

			for _, imp := range syn.Imports {
				trim := strings.Trim(imp.Path.Value, `"`)
				parts := strings.Split(trim, "/")

				if len(parts) >= 3 && strings.Join(parts[:3], "/") == lib {
					newImp, err := createNewImport(parts, newVersion)
					if err != nil {
						return err
					}

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
				return fmt.Errorf("could not create go file %s: %w", goFileName, err)
			}

			err = format.Node(f, p.Fset, syn)
			_ = f.Close()

			if err != nil {
				return fmt.Errorf("could not rewrite go file %s: %w", goFileName, err)
			}
		}
	}

	return nil
}

func createNewImport(parts []string, newVersion string) (string, error) {
	if len(parts) < 3 {
		return "", fmt.Errorf("unsupported package format: %s", strings.Join(parts, "/"))
	}

	np := make([]string, 3)
	copy(np, parts[:3])
	np = append(np, newVersion)

	if len(parts) == 3 {
		// no version
		return strings.Join(np, "/"), nil
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

	return strings.Join(np, "/"), nil
}

func updateModFile(dir, lib, full, major string) error {
	const filename = "go.mod"

	modPath := path.Join(dir, filename)

	content, err := os.ReadFile(filepath.Clean(modPath))
	if err != nil {
		return err
	}

	file, err := modfile.Parse(filename, content, nil)
	if err != nil {
		return err
	}

	for _, require := range file.Require {
		if !strings.HasPrefix(require.Mod.Path, lib) || require.Indirect {
			continue
		}

		// Replaces the current module path:
		// when the `AddRequire` adds the new path,
		// the existing module will be replaced instead of be added at the end of the last `require` block.
		require.Mod.Path = path.Join(lib, major)
	}

	err = file.AddRequire(path.Join(lib, major), full)
	if err != nil {
		return err
	}

	file.Cleanup()

	data, err := file.Format()
	if err != nil {
		return err
	}

	stat, err := os.Stat(modPath)
	if err != nil {
		return err
	}

	return os.WriteFile(modPath, data, stat.Mode())
}

func guessVersion(lib string, latest bool, rawVersion string) (string, string, error) {
	if ok, _ := regexp.MatchString(`^v?\d+\.\d+\.\d+.*$`, rawVersion); ok {
		return rawVersion, semver.Major(rawVersion), nil
	}

	client := goproxy.NewClient("")

	var moduleName string

	if latest || rawVersion == "latest" {
		latestVersion, err := internal.FindHighestFromGoPkg(lib)
		if err != nil {
			return "", "", err
		}

		moduleName = path.Join(lib, semver.Major(latestVersion))
	} else {
		moduleName = path.Join(lib, "v"+strings.TrimPrefix(rawVersion, "v"))
	}

	lst, err := client.GetLatest(moduleName)
	if err != nil {
		return "", "", err
	}

	return lst.Version, semver.Major(lst.Version), nil
}
