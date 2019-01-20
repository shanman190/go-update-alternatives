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
		return fmt.Errorf("alternative path %s doesn't exist", path)
	} else if err != nil {
		return fmt.Errorf("unknown error encountered: %s", err)
	}

	alternativesDir := config.GetAlternativesDir()
	if err := os.MkdirAll(alternativesDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create alternatives directory '%s'", alternativesDir)
	}

	alternativePath := filepath.Join(alternativesDir, group)

	if _, err := os.Stat(alternativePath); os.IsNotExist(err) {
		fmt.Printf("update-alternatives: using %s to provide %s (%s)\n", path, link, group)
	}

	config.SaveAlternative(link, group, path)

	if err := symbolic.Ln(path, alternativePath); err != nil {
		return fmt.Errorf("unable to install '%s' to '%s': %s", path, alternativePath, err)
	}
	if err := symbolic.Ln(alternativePath, link); err != nil {
		return fmt.Errorf("unable to install '%s' to '%s': %s", alternativePath, link, err)
	}
	
	return nil
}