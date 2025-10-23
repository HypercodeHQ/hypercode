package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

// CommitCommand creates the commit command for creating conventional commits
func CommitCommand() *cli.Command {
	return NewCommandWithFlags(
		"commit",
		"Create a conventional commit",
		`Interactive tool to help you create conventional commits.

Conventional Commits format: <type>(<scope>): <description>

Types:
  feat     - A new feature
  fix      - A bug fix
  docs     - Documentation only changes
  style    - Changes that don't affect code meaning (formatting, etc)
  refactor - Code change that neither fixes a bug nor adds a feature
  perf     - Performance improvements
  test     - Adding or updating tests
  build    - Changes to build system or dependencies
  ci       - Changes to CI configuration
  chore    - Other changes that don't modify src or test files
  revert   - Reverts a previous commit`,
		[]cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Stage all changes before committing (git add -A)",
			},
			&cli.BoolFlag{
				Name:    "push",
				Aliases: []string{"p"},
				Usage:   "Push after committing",
			},
		},
		runCommit,
	)
}

func runCommit(ctx context.Context, cmd *cli.Command) error {
	// Check if we're in a git repository
	if !isGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	// Stage all changes if -a flag is set
	if cmd.Bool("all") {
		if err := gitAddAll(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
		fmt.Println("✓ Staged all changes")
	}

	// Check if there are staged changes
	hasChanges, err := hasStagedChanges()
	if err != nil {
		return fmt.Errorf("failed to check for staged changes: %w", err)
	}
	if !hasChanges {
		return fmt.Errorf("no changes staged for commit. Use 'git add' or the -a flag")
	}

	// Collect commit information via interactive form
	var (
		commitType        string
		scope             string
		description       string
		body              string
		breaking          bool
		breakingDesc      string
		includeBodyPrompt bool
	)

	// Type selection
	typeForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select commit type").
				Options(
					huh.NewOption("feat - A new feature", "feat"),
					huh.NewOption("fix - A bug fix", "fix"),
					huh.NewOption("docs - Documentation changes", "docs"),
					huh.NewOption("style - Code style changes (formatting)", "style"),
					huh.NewOption("refactor - Code refactoring", "refactor"),
					huh.NewOption("perf - Performance improvements", "perf"),
					huh.NewOption("test - Add or update tests", "test"),
					huh.NewOption("build - Build system or dependencies", "build"),
					huh.NewOption("ci - CI configuration changes", "ci"),
					huh.NewOption("chore - Other changes", "chore"),
					huh.NewOption("revert - Revert a previous commit", "revert"),
				).
				Value(&commitType),
		),
	)

	if err := typeForm.Run(); err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	// Scope and description
	mainForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Scope (optional)").
				Description("Component or module affected (e.g., api, ui, auth)").
				Placeholder("Leave empty if not applicable").
				Value(&scope),

			huh.NewInput().
				Title("Short description").
				Description("A brief description of the change").
				Placeholder("e.g., add user authentication").
				CharLimit(72).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("description is required")
					}
					return nil
				}).
				Value(&description),

			huh.NewConfirm().
				Title("Add detailed description?").
				Description("Include a longer description of the changes").
				Value(&includeBodyPrompt),
		),
	)

	if err := mainForm.Run(); err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	// Body (if requested)
	if includeBodyPrompt {
		bodyForm := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Detailed description").
					Description("Explain what and why (press ESC then Enter when done)").
					Placeholder("Describe the changes in detail...").
					CharLimit(500).
					Value(&body),
			),
		)

		if err := bodyForm.Run(); err != nil {
			return fmt.Errorf("form cancelled: %w", err)
		}
	}

	// Breaking changes
	breakingForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Does this introduce breaking changes?").
				Description("Changes that require updates to dependent code").
				Value(&breaking),
		),
	)

	if err := breakingForm.Run(); err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	if breaking {
		breakingDescForm := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Describe the breaking changes").
					Description("Explain what breaks and how to migrate").
					Placeholder("BREAKING CHANGE: ...").
					CharLimit(500).
					Validate(func(s string) error {
						if strings.TrimSpace(s) == "" {
							return fmt.Errorf("breaking change description is required")
						}
						return nil
					}).
					Value(&breakingDesc),
			),
		)

		if err := breakingDescForm.Run(); err != nil {
			return fmt.Errorf("form cancelled: %w", err)
		}
	}

	// Build the commit message
	commitMsg := buildCommitMessage(commitType, scope, description, body, breaking, breakingDesc)

	// Show preview and confirm
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("Commit message preview:")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(commitMsg)
	fmt.Println(strings.Repeat("─", 50))

	var confirm bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create this commit?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	)

	if err := confirmForm.Run(); err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	if !confirm {
		fmt.Println("Commit cancelled")
		return nil
	}

	// Create the commit
	if err := gitCommit(commitMsg); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	fmt.Println("✓ Commit created successfully")

	// Push if -p flag is set
	if cmd.Bool("push") {
		fmt.Println("Pushing to remote...")
		if err := gitPush(); err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
		fmt.Println("✓ Pushed successfully")
	}

	return nil
}

func buildCommitMessage(commitType, scope, description, body string, breaking bool, breakingDesc string) string {
	var msg strings.Builder

	// Header
	if scope != "" {
		scope = strings.TrimSpace(scope)
		msg.WriteString(fmt.Sprintf("%s(%s): %s", commitType, scope, description))
	} else {
		msg.WriteString(fmt.Sprintf("%s: %s", commitType, description))
	}

	// Body
	if body != "" {
		msg.WriteString("\n\n")
		msg.WriteString(strings.TrimSpace(body))
	}

	// Breaking changes footer
	if breaking && breakingDesc != "" {
		msg.WriteString("\n\n")
		msg.WriteString("BREAKING CHANGE: ")
		msg.WriteString(strings.TrimSpace(breakingDesc))
	}

	return msg.String()
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func hasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err == nil {
		return false, nil // No changes
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return true, nil // Has changes
	}
	return false, err
}

func gitAddAll() error {
	cmd := exec.Command("git", "add", "-A")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
