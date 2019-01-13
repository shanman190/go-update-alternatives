package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shanman190/update-alternatives/config"
	"github.com/shanman190/update-alternatives/symbolic"
	"github.com/shanman190/update-alternatives/ui"
)

type ConfigCommand struct {}

func (command *ConfigCommand) Execute(args []string) error {
	fmt.Printf("--config %s\n", args)

	if len(args) != 1 {
		fmt.Fprintf(ui.Stderr, "Invalid usage")
		os.Exit(1)
	}

	group := args[0]

	alternatives, err := config.LoadAlternatives(group)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Could not load configuration for group %s\n", group)
		os.Exit(1)
	}

	alternativesLen := len(alternatives.Alternatives)
	if alternativesLen == 1 {
		fmt.Printf("There is only one alternative in link group %s.\n: %s\nNothing to configure.\n", group, alternatives.Alternatives[0])
		return nil
	}

	alternativePath := filepath.Join(config.GetAlternativesDir(), group)
	currentLink, err := symbolic.Readlink(alternativePath)
	if err != nil {
		fmt.Fprintf(ui.Stderr, "Error: unable to read symbolic link %s", err)
		os.Exit(1)
	}
	
	fmt.Printf("There are %d choices for the alternative %s.\nSelection\tPath\n", alternativesLen, group)
	fmt.Printf("%s\n", strings.Repeat("-", 80))

	for index, alternative := range alternatives.Alternatives {
		if currentLink == alternative {
			fmt.Printf("* %d\t\t%s\n", index, alternative)
		} else {
			fmt.Printf("  %d\t\t%s\n", index, alternative)
		}
	}

	fmt.Println()
	fmt.Print("Press <enter> to keep the current choice[*], or type selection number: ")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	if text != "" {
		choice, err := strconv.Atoi(text)
		if err != nil {
			fmt.Fprintf(ui.Stderr, "Unable to get selection: %s\n", err)
			os.Exit(1)
		}

		err = symbolic.Ln(alternatives.Alternatives[choice], alternativePath)
		if err != nil {
			fmt.Fprintf(ui.Stderr, "could not create symbolic link '%s' from choice '%s'", alternativePath, alternatives.Alternatives[choice])
			os.Exit(1)
		}
	}

	return nil
}