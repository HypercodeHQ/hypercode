package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hypercommithq/hypercommit/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "hypercommit",
		Usage:   "Hypercommit CLI - Tools for working with Hypercommit repositories",
		Version: "1.0.0",
		Authors: []any{
			"Hypercommit Team",
		},
		Description: `The Hypercommit CLI provides utilities for working with Hypercommit repositories,
including tools for creating conventional commits, managing repositories, and more.`,
		Commands: commands.GetAllCommands(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
