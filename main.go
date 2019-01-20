package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/shanman190/update-alternatives/commands"
	"github.com/shanman190/update-alternatives/ui"
)

var options struct {
	Install bool `long:"install"`
	Remove  bool `long:"remove"`
	Set     bool `long:"set"`
	Display bool `long:"display"`
	Config  bool `long:"config"`
	Help    bool `short:"h" long:"help"`
}

type Command interface {
	Execute(args []string) error
}

func main() {
	parser := flags.NewParser(&options, flags.PassDoubleDash)

	remainingArgs, err := parser.Parse()
	if err != nil {
		handleError(fmt.Errorf("unable to parse arguments: %s", err))
	}

	var command Command
	if options.Install {
		command = &commands.InstallCommand{}
	} else if options.Remove {
		command = &commands.RemoveCommand{}
	} else if options.Set {
		command = &commands.SetCommand{}
	} else if options.Display {
		command = &commands.DisplayCommand{}
	} else if options.Help {
		command = &commands.HelpCommand{}
	} else if options.Config {
		command = &commands.ConfigCommand{}
	} else {
		handleError(fmt.Errorf("unknown command: %s\n", os.Args[1:]))
	}

	handleError(command.Execute(remainingArgs))
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(ui.Stderr, "update-alternatives: error: %s\n", err)
		os.Exit(1)
	}
}
