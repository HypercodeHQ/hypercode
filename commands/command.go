package commands

import (
	"github.com/urfave/cli/v3"
)

// NewCommand creates a new CLI command with the given configuration
func NewCommand(name, usage, description string, action cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:        name,
		Usage:       usage,
		Description: description,
		Action:      action,
	}
}

// NewCommandWithFlags creates a new CLI command with flags
func NewCommandWithFlags(name, usage, description string, flags []cli.Flag, action cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:        name,
		Usage:       usage,
		Description: description,
		Flags:       flags,
		Action:      action,
	}
}

// GetAllCommands returns all available CLI commands
func GetAllCommands() []*cli.Command {
	return []*cli.Command{
		LoginCommand(),
		CommitCommand(),
	}
}
