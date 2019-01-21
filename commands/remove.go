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
		return fmt.Errorf("unable to create alternatives directory '%s'", alternativesDir)
	}
	var link string
	if alternatives, err := config.LoadAlternatives(group); err == nil {
		if len(alternatives.Alternatives) == 0 {
			return nil
		}
		link = alternatives.Link
	} else if err != nil {
		return fmt.Errorf("unable to open '%s': %s", filepath.Join(alternativesDir, group), err)
	}

	config.DeleteAlternative(group, path)
	
	alternativePath := filepath.Join(alternativesDir, group)
	linkPath, linkErr := symbolic.Readlink(alternativePath)
	if alternatives, err := config.LoadAlternatives(group); err == nil {
		if len(alternatives.Alternatives) == 0 {
			if err := symbolic.Unlink(link); err != nil {
				return fmt.Errorf("unable to remove symbolic link at '%s': %s", link, err)
			}
		} else {
			if linkPath == path {
				if err := symbolic.Unlink(alternativePath); err != nil {
					return fmt.Errorf("unable to remove symbolic link at '%s': %s", alternativePath, err)
				}

				indexHigh := 0
				priorityHigh := alternatives.Alternatives[0].Priority
				for index, alternative := range alternatives.Alternatives[1:] {
					if priorityHigh < alternative.Priority {
						indexHigh = index
						priorityHigh = alternative.Priority
					}
				}

				if err := symbolic.Ln(alternatives.Alternatives[indexHigh].Path, alternativePath); err == nil {
					fmt.Printf("update-alternatives: using %s to provide %s (%s)\n", alternatives.Alternatives[indexHigh].Path, link, group)
				} else {
					return fmt.Errorf("unable to install '%s' to '%s': %s", alternatives.Alternatives[indexHigh].Path, link, err)
				}
			}

			if linkErr != nil {
				fmt.Fprintf(ui.Stderr, "update-alternatives: warning: forcing reinstallation of alternative (%s) because link group %s is broken\n", "null", group)
			}
		}
	}

	return nil
}