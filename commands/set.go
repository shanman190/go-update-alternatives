package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shanman190/update-alternatives/config"
	"github.com/shanman190/update-alternatives/symbolic"
	"github.com/shanman190/update-alternatives/ui"
)

type SetCommand struct {}

func (command *SetCommand) Execute(args []string) error {
	if len(args) != 2 {
		fmt.Fprintf(ui.Stderr, "Invalid usage\n")
		os.Exit(1)
	}

	group := args[0]
	path := filepath.Clean(args[1])

	alternativesDir := config.GetAlternativesDir()
	alternativePath := filepath.Join(alternativesDir, group)

	err := symbolic.Ln(path, alternativePath)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: Unable to create symboic link from %s to %s\n\n%s\n", path, alternativePath, err)
		os.Exit(1)
	}
	
	return nil
}