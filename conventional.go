package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

type CommitType struct {
	Type        string
	Description string
	SubType     string
}

var CommitTypes = []CommitType{
	{Type: "chore", Description: "Changes that don't change source code or tests"},
	{Type: "feat", Description: "Adds or removes a new feature"},
	{Type: "fix", Description: "Fixes a bug"},
	{Type: "refactor", Description: "A code change that neither fixes a bug nor adds a feature, eg. renaming a variable, remove dead code, etc."},
	{Type: "docs", Description: "Documentation only changes"},
	{Type: "style", Description: "Changes the style of the code eg. linting"},
	{Type: "perf", Description: "Improves the performance of the code"},
	{Type: "test", Description: "Adding missing tests or correcting existing tests"},
	{Type: "build", Description: "Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)"},
	{Type: "ci", Description: "Changes to CI configuration files and scripts"},
	{Type: "revert", Description: "Reverts a previous commit"},
	{Type: "chore", SubType: "release", Description: "Release / Version tags"},
	{Type: "chore", SubType: "deps", Description: "Add, remove or update dependencies"},
	{Type: "chore", SubType: "dev-deps", Description: "Add, remove or update development dependencies"},
	{Type: "chore", SubType: "types", Description: "Add or update types."},
}

type Commit struct {
	Type    string
	Message string
}

type Conventional struct{}

func NewConventional() *Conventional {
	return &Conventional{}
}

// promptForScope prompts the user for the main commit type and optional sub-type
func (c *Conventional) ask() (string, error) {
	var main, opt string

	options := make([]huh.Option[string], 0, len(CommitTypes))
	for _, ct := range CommitTypes {
		optionText := fmt.Sprintf("%s: %s", ct.Type, ct.Description)
		optionValue := ct.Type

		// If there's a sub-type associated with the commit type, include it in the option text and value
		if ct.SubType != "" {
			optionValue = fmt.Sprintf("%s(%s)", ct.Type, ct.SubType)
			optionText = fmt.Sprintf("%s: %s", optionValue, ct.Description)
		}

		options = append(options, huh.NewOption(optionText, optionValue))
	}

	if err := huh.NewSelect[string]().
		Title("Select the type of commit").
		Options(options...).
		Value(&main).
		Filtering(true).
		Height(10).
		Validate(func(val string) error {
			if val == "" {
				return errors.New("type cannot be empty")
			}

			return nil
		}).Run(); err != nil {
		return "", err
	}

	// If the user didn't provide an optional sub-type, just return the main type
	if opt == "" {
		return main, nil
	}

	return fmt.Sprintf("%s(%s)", main, opt), nil
}

// Prompt user for commit type, scope, and message, then execute the commit
func (c *Conventional) Construct() (string, error) {
	// Get the commit scope (type and optional sub-type)
	scope, err := c.ask()
	if err != nil {
		return "", err
	}

	// Get the commit message
	msg, err := promptForMessage()
	if err != nil {
		return "", err
	}

	// Combine scope and message into a conventional commit format
	return fmt.Sprintf("%s: %s", scope, msg), nil
}

func (c *Conventional) Type() ConventionType {
	return ConventionalCommitConvention
}

var _ Provider = (*Conventional)(nil)
