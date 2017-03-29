package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shimt/go-simplecli"
)

var cli = simplecli.NewCLI()

func initCLI() {
	cli.CommandLine.String("envname", "PATH", "Target environment variable")
	cli.CommandLine.String("separator", string(os.PathListSeparator), "Path separator")
	cli.CommandLine.String("setenv-sh", "", "Set environment style output(sh)")
	cli.CommandLine.Bool("all", false, "Output all candidates")

	cli.BindSameName("envname")
	cli.BindSameName("separator")
	cli.BindSameName("setenv-sh")
}

func init() {
	initCLI()
}

func outputPath(path string) {
	setenvSh := cli.Config.GetString("setenv-sh")

	if setenvSh != "" {
		fmt.Printf("%s=\"%s\"\n", setenvSh, path)
	} else {
		fmt.Print(path)
	}
}

func main() {
	err := cli.Setup()
	cli.Exit1IfError(err)

	if cli.CommandLine.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s ENVNAME\n", os.Args[0])
		cli.CommandLine.PrintDefaults()
	}

	env := os.Getenv(cli.Config.GetString("envname"))
	all := cli.Config.GetBool("all")
	filename := cli.CommandLine.Arg(0)

	for _, directory := range strings.Split(env, cli.Config.GetString("separator")) {
		path := filepath.Join(directory, filename)

		if stat, err := os.Stat(path); err == nil {
			if !stat.IsDir() {
				outputPath(path)
				if !all {
					break
				}
			}
		}
	}

	cli.Exit(0)
}
