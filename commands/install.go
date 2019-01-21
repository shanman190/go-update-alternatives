package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shanman190/update-alternatives/config"
	"github.com/shanman190/update-alternatives/symbolic"
	"github.com/shanman190/update-alternatives/ui"
)

type InstallCommand struct {}

func (command *InstallCommand) Execute(args []string) error {
	if len(args) != 4 {
		fmt.Fprintf(ui.Stderr, "Invalid usage\n")
		os.Exit(1)
	}

	link := filepath.Clean(args[0])
	group := args[1]
	path := filepath.Clean(args[2])
	priority, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "update-alternatives: priority must be an integer\n\nUse 'update-alternatives --help' for program usage information.\n")
		os.Exit(1)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("alternative path %s doesn't exist", path)
	} else if err != nil {
		return fmt.Errorf("unknown error encountered: %s", err)
	}

	alternativesDir := config.GetAlternativesDir()
	if err := os.MkdirAll(alternativesDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create alternatives directory '%s'", alternativesDir)
	}

	config.SaveAlternative(link, group, path, priority)

	alternativePath := filepath.Join(alternativesDir, group)
	linkPath, linkErr := symbolic.Readlink(alternativePath)
	if linkErr == nil {
		if alternatives, err := config.LoadAlternatives(group); err == nil {
			for _, alternative := range alternatives.Alternatives {
				if linkPath == alternative.Path && priority > alternative.Priority {
					if err := symbolic.Ln(path, alternativePath); err == nil {
						fmt.Printf("update-alternatives: using %s to provide %s (%s)\n", path, link, group)
					} else {
						return fmt.Errorf("unable to install '%s' to '%s': %s", path, alternativePath, err)
					}
					break
				}
			}
		}
	} else {
		if err := symbolic.Ln(path, alternativePath); err == nil {
			fmt.Printf("update-alternatives: using %s to provide %s (%s)\n", path, link, group)
		} else {
			return fmt.Errorf("unable to install '%s' to '%s': %s", path, alternativePath, err)
		}
	}

	if _, err := symbolic.Readlink(link); os.IsNotExist(err) {
		if err := symbolic.Ln(alternativePath, link); err != nil {
			return fmt.Errorf("unable to install '%s' to '%s': %s", alternativePath, link, err)
		}
	}
	
	return nil
}