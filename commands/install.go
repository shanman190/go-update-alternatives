package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shanman190/update-alternatives/config"
	"github.com/shanman190/update-alternatives/symbolic"
	"github.com/shanman190/update-alternatives/ui"
)

type InstallCommand struct {}

func (command *InstallCommand) Execute(args []string) error {
	if len(args) != 3 {
		fmt.Fprintf(ui.Stderr, "Invalid usage\n")
		os.Exit(1)
	}

	link := filepath.Clean(args[0])
	group := args[1]
	path := filepath.Clean(args[2])

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(ui.Stderr, "update-alternatives: error: alternative path %s doesn't exist\n", path)
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(ui.Stderr, "update-alternatives: error: unknown error encountered")
		os.Exit(1)
	}

	alternativesDir := config.GetAlternativesDir()
	if err := os.MkdirAll(alternativesDir, os.ModePerm); err != nil {
		fmt.Fprintf(ui.Stderr, "update-alternatives: error: unable to create alternatives directory '%s'\n", alternativesDir)
		os.Exit(1)
	}

	alternativePath := filepath.Join(alternativesDir, group)

	if _, err := os.Stat(alternativePath); os.IsNotExist(err) {
		fmt.Printf("update-alternatives: using %s to provide %s (%s)\n", path, link, group)
	}

	config.SaveAlternative(link, group, path)

	if err := symbolic.Ln(path, alternativePath); err != nil {
		fmt.Fprintf(ui.Stderr, "update-alternatives: error: unable to install '%s' to '%s': %s\n", path, alternativePath, err)
		os.Exit(1)
	}
	if err := symbolic.Ln(alternativePath, link); err != nil {
		fmt.Fprintf(ui.Stderr, "update-alternatives: error: unable to install '%s' to '%s': %s\n", alternativePath, link, err)
		os.Exit(1)
	}
	
	return nil
}