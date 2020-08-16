package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

const (
	errorExitCode      = 1
	minArgNumber       = 2
	moduleNameFragment = 3
)

type config struct {
	lib      string
	version  string
	latest   bool
	filename bool
}

func main() {
	cfg := config{}

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	updateCommand.StringVar(&cfg.lib, "lib", "", "Lib to update. (Required)")
	updateCommand.StringVar(&cfg.version, "version", "", "Version to set.")
	updateCommand.BoolVar(&cfg.latest, "latest", false, "Update to the latest available version.")
	updateCommand.BoolVar(&cfg.filename, "filenames", false, "Only display file names.")

	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)

	cmds := []*flag.FlagSet{updateCommand, versionCommand}

	if len(os.Args) < minArgNumber {
		fmt.Fprintln(os.Stderr, "A subcommand is required.")
		commandsUsage(cmds)
		os.Exit(errorExitCode)
	}

	switch os.Args[1] {
	case "-h", "--help":
		flag.CommandLine.SetOutput(os.Stdout)
		commandsUsage(cmds)
		os.Exit(0)
	default:
		cmd := getCommand(cmds)
		if cmd != nil {
			err := cmd.Parse(os.Args[2:])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				cmd.Usage()
				os.Exit(errorExitCode)
			}
		} else {
			commandsUsage(cmds)
			os.Exit(errorExitCode)
		}
	}

	switch {
	case versionCommand.Parsed():
		displayVersion()
		os.Exit(0)
	case updateCommand.Parsed():
		err := updateCmd(cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func updateCmd(cfg config) error {
	if len(cfg.lib) == 0 {
		return errors.New("--lib is required")
	}

	if len(strings.Split(cfg.lib, "/")) != moduleNameFragment {
		return fmt.Errorf("--lib: invalid format: %s", cfg.lib)
	}

	if len(cfg.version) == 0 && !cfg.latest {
		return errors.New("--version or --latest are required")
	}

	if len(cfg.version) > 0 && cfg.latest {
		return errors.New("--version and --latest cannot be used at the same time")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = os.Stat(path.Join(dir, "go.mod"))
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("unable to find 'go.mod': %s", dir)
	}

	v, err := getNewVersion(cfg.latest, cfg.lib, cfg.version)
	if err != nil {
		return err
	}

	return update(dir, cfg.lib, v, cfg.filename)
}

func getCommand(cmds []*flag.FlagSet) *flag.FlagSet {
	for _, cmd := range cmds {
		if os.Args[1] == cmd.Name() {
			return cmd
		}
	}

	return nil
}

func commandsUsage(cmds []*flag.FlagSet) {
	flag.Usage()

	output := flag.CommandLine.Output()

	fmt.Fprintf(output, "\n  %s <command> [<flags>]\n\n", path.Base(os.Args[0]))
	fmt.Fprintln(output, "Commands:")

	for _, cmd := range cmds {
		fmt.Fprintf(output, "  %-8s [<flags>]\n", cmd.Name())
	}

	fmt.Fprintln(output)
	fmt.Fprintln(output, "Flags:")
	fmt.Fprintln(output, "  --help,-h  Display help")
}
