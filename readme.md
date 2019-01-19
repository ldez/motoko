# Motoko (Major Motoko Kusanagi)

[![Build Status](https://travis-ci.org/ldez/motoko.svg?branch=master)](https://travis-ci.org/ldez/motoko)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/motoko)](https://goreportcard.com/report/github.com/ldez/motoko)

Based on Go modules, update a dependency to a major version.

## How to Install

### Binaries

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/motoko/releases)
* Unzip the archive.
* Add `motoko` in your `PATH`.

Available for: Linux, MacOS, Windows, FreeBSD, OpenBSD.

### From a package manager

- [ArchLinux (AUR)](https://aur.archlinux.org/packages/motoko/):
```bash
yay -S motoko
```

- [Homebrew Taps](https://github.com/ldez/homebrew-tap)
```bash
brew tap ldez/tap
brew update
brew install motoko
```

- [Scoop Bucket](https://github.com/ldez/scoop-bucket)
```bash
scoop bucket add motoko https://github.com/ldez/scoop-bucket.git
scoop install motoko
```

### From sources

```bash
go get -u github.com/ldez/motoko
```

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
