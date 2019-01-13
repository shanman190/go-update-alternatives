package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shanman190/update-alternatives/config"
	"github.com/shanman190/update-alternatives/symbolic"
	"github.com/shanman190/update-alternatives/ui"
)

type DisplayCommand struct {}

func (command *DisplayCommand) Execute(args []string) error {
	if len(args) != 1 {
		fmt.Fprintln(ui.Stderr, "Invalid usage")
		os.Exit(1)
	}

	group := args[0]

	alternatives, err := config.LoadAlternatives(group)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: could not load alternative config for '%s'\n%s\n", group, err)
		os.Exit(1)
	}

	alternativesDir := config.GetAlternativesDir()
	if err := os.MkdirAll(alternativesDir, os.ModePerm); err != nil {
		fmt.Fprintf(ui.Stderr, "Error: unable to create directory '%s'\n", alternativesDir)
		os.Exit(1)
	}
	alternativePath := filepath.Join(alternativesDir, group)

	currentLink, err := symbolic.Readlink(alternativePath)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: unable to read symbolic link %s", err)
		os.Exit(1)
	}

	for index, alternative := range alternatives.Alternatives {
		if currentLink == alternative {
			fmt.Printf("* %d  %s\n", index, alternative)
		} else {
			fmt.Printf("  %d  %s\n", index, alternative)
		}
	}

	return nil
}