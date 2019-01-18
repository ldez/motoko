package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	libPtr := updateCommand.String("lib", "", "Lib to update. (Required)")
	versionPtr := updateCommand.String("version", "", "Version to set. (Required)")
	latestPtr := updateCommand.Bool("latest", false, "Update to the latest available version.")
	filenamePtr := updateCommand.Bool("filenames", false, "Only display file names.")

	cmds := []*flag.FlagSet{updateCommand, versionCommand}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "A subcommand is required.")
		commandsUsage(cmds)
		os.Exit(1)
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
				os.Exit(1)
			}
		} else {
			commandsUsage(cmds)
			os.Exit(1)
		}
	}

	switch {
	case versionCommand.Parsed():
		displayVersion()
		os.Exit(0)
	case updateCommand.Parsed():
		updateCmd(*latestPtr, *filenamePtr, *libPtr, *versionPtr)
	}
}

func updateCmd(latest, filename bool, lib, version string) {
	if len(lib) == 0 {
		log.Fatal("--lib is required")
	}

	if len(strings.Split(lib, "/")) != 3 {
		log.Fatal("--lib: invalid format:", lib)
	}

	if len(version) == 0 && !latest {
		log.Fatal("--version or --latest are required")
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(path.Join(dir, "go.mod"))
	if err != nil && os.IsNotExist(err) {
		log.Fatal("Unable to find 'go.mod':", dir)
	}

	fmt.Println(dir)

	v := getNewVersion(latest, lib, version)

	if err := update(dir, lib, v, filename); err != nil {
		log.Fatal(err)
	}
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
