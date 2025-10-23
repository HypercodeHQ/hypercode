package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hyperstitieux/hypercode/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "hypercode",
		Usage:   "Hypercode CLI - Tools for working with Hypercode repositories",
		Version: "1.0.0",
		Authors: []any{
			"Hypercode Team",
		},
		Description: `The Hypercode CLI provides utilities for working with Hypercode repositories,
including tools for creating conventional commits, managing repositories, and more.`,
		Commands: commands.GetAllCommands(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
