package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shanman190/update-alternatives/config"
	"github.com/shanman190/update-alternatives/symbolic"
	"github.com/shanman190/update-alternatives/ui"
)

type RemoveCommand struct {}

func (command *RemoveCommand) Execute(args []string) error {
	if len(args) != 2 {
		fmt.Fprintf(ui.Stderr, "Invalid usage\n")
		os.Exit(1)
	}

	group := args[0]
	path := filepath.Clean(args[1])

	alternativesDir := config.GetAlternativesDir()
	if err := os.MkdirAll(alternativesDir, os.ModePerm); err != nil {
		fmt.Fprintf(ui.Stderr, "Error: unable to create directory '%s'\n", alternativesDir)
		os.Exit(1)
	}
	config.DeleteAlternative(group, path)
	alternativePath := filepath.Join(alternativesDir, group)

	linkPath, err := symbolic.Readlink(alternativePath)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: Unable to resolve symlink %s\n\n%s\n", alternativePath, err)
		os.Exit(1)
	}

	if linkPath == path {
		err = symbolic.Unlink(alternativePath)
		if err != nil {
			fmt.Fprintf(ui.Stderr, "Error: Unable to remove symbolic link from %s\n\n%s\n", alternativePath, err)
			os.Exit(1)
		}
	}

	return nil
}