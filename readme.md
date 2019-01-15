# Motoko (Major Motoko Kusanagi)

Based on Go modules, update a dependency to a major version.

## Main

```bash
Usage of motoko:

  motoko <command> [<flags>]

Commands:
  update   [<flags>]
  version  [<flags>]

Flags:
  --help,-h  Display help
```

## update

**Note**: for now, `--latest` works only with dependency on Github.

```bash
Usage of update:
  -filenames
        Only display file names.
  -latest
        Update to the latest available version.
  -lib string
        Lib to update. (Required)
  -version string
        Version to set. (Required)
```

## Examples

```bash
# update to the latest version:
motoko update --lib github.com/ldez/go-git-cmd-wrapper --latest

# update to a specific version:
motoko update --lib github.com/ldez/go-git-cmd-wrapper --version 6
```
