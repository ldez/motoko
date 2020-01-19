# Motoko (Major Motoko Kusanagi)


[![release](https://img.shields.io/github/tag/ldez/motoko.svg)](https://github.com/ldez/motoko/releases)
[![Build Status](https://travis-ci.com/ldez/motoko.svg?branch=master)](https://travis-ci.com/ldez/motoko)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/motoko)](https://goreportcard.com/report/github.com/ldez/motoko)

[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/ldez)

Based on Go modules, update a dependency to a major version.

## How to Use

### Main

```bash
Usage of motoko:

  motoko <command> [<flags>]

Commands:
  update   [<flags>]
  version  [<flags>]

Flags:
  --help,-h  Display help
```

### Update

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
        Version to set.
```

### Examples

```bash
# update to the latest version:
motoko update --lib github.com/ldez/go-git-cmd-wrapper --latest

# update to a specific version:
motoko update --lib github.com/ldez/go-git-cmd-wrapper --version 6
```

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
