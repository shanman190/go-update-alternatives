package commands

import (
	"errors"
	"fmt"
)

var ErrShowHelpMessage = errors.New("help command invoked")

type HelpCommand struct {}

func (command *HelpCommand) Execute(args []string) error {
	fmt.Printf(`
Usage:
	update-alternatives --install LINK NAME PATH
	update-alternatives --remove NAME PATH
	update-alternatives --set NAME PATH
	update-alternatives --display NAME
	update-alternatives --config NAME
	update-alternatives -h|--help

Commands:
	--install LINK NAME PATH
			Add a group of alternatives to the system. NAME is the 
			generic name for the master link, LINK is the name of 
			it's symlink, PATH is the alternative being introduced 
			for the master link.
	
	--remove NAME PATH
			Remove an alternative. NAME is a name in the 
			alternatives directory and PATH is an absolute filename
			to which NAME could be linked. If NAME is indeed linked
			to PATH, NAME will be updated to point to another 
			appropriate alternative or removed, correspondingly. If 
			the link is not currently pointing to PATH, no links 
			are changed; only the information about the alternative
			is removed.

	--set NAME PATH
			The symbolic link for group NAME set to those 
			configured for PATH.

	--display NAME
			Display information about the link group of which NAME
			is the master link. Information displayed includes 
			which alternative the symlink currently points to and
			other alternatives that are available.

	--config NAME
			Show available alternatives for a link group and allow
			the user to select which one to use. The link group is
			then updated with the selection.

	-h|--help	Show this help message

Examples:
	update-alternatives --install /usr/bin/java java /usr/lib/jvm/jdk-11.0.6/bin/java 0
	update-alternatives --remove java /usr/lib/jvm/jdk-11.0.6/bin/java
	update-alternatives --set java /usr/lib/jvm/jdk-11.0.6/bin/java
	update-alternatives --display java
	update-alternatives --config java
`)

	return nil
}