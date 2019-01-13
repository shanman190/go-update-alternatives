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

	alternativesDir := config.GetAlternativesDir()
	if err := os.MkdirAll(alternativesDir, os.ModePerm); err != nil {
		fmt.Fprintf(ui.Stderr, "Error: unable to create directory '%s'\n", alternativesDir)
		os.Exit(1)
	}
	config.SaveAlternative(group, path)
	alternativePath := filepath.Join(alternativesDir, group)

	err := symbolic.Ln(path, alternativePath)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: Unable to create symboic link from %s to %s\n\n%s\n", path, alternativePath, err)
		os.Exit(1)
	}
	err = symbolic.Ln(alternativePath, link)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: Unable to create symbolic link from %s to %s\n\n%s\n", alternativePath, link, err)
		os.Exit(1)
	}
	
	return nil
}