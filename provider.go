package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"

	log "github.com/charmbracelet/log"
)

type Provider interface {
	Type() ConventionType
	Construct() (string, error)
}

func startsWithColonOrEmoji(s string) (bool, string) {
	if strings.HasPrefix(s, ":") {
		return true, "code"
	}

	// Check if the string starts with an emoji
	firstRune, _ := utf8.DecodeRuneInString(s)

	if isEmoji(string(firstRune)) {
		return true, "emoji"
	}

	return false, ""
}

func isEmoji(s string) bool {
	emojiPattern := regexp.MustCompile(`^[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F1E0}-\x{1F1FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)
	return emojiPattern.MatchString(s)
}

func determineConventionFromCommitMessage() (Provider, error) {
	msg, err := exec.Command("git", "log", "-1", "--pretty=%B").Output()
	if err != nil {
		return nil, err
	}

	msg = bytes.TrimSpace(msg)
	log.Debug("last commit message", "message", string(msg))

	usesGitmoji, notation := startsWithColonOrEmoji(string(msg))
	if usesGitmoji {
		return NewGitmoji(notation), nil
	} else {
		return NewConventional(), nil
	}
}

func determineConvention() (Provider, error) {
	cfg, err := ReadConfig()
	if err != nil {
		log.Debug("error reading config", "error", err)

		fallback, err := determineConventionFromCommitMessage()
		if err != nil {
			return nil, err
		}

		log.Debug("falling back to last used convention", "provider", fallback.Type())

		return fallback, nil
	}

	var provider Provider
	switch cfg.Convention {
	case ConventionalCommitConvention:
		provider = NewConventional()
	case GitmojiConvention:
		provider = NewGitmoji("code")
	default:
		return nil, fmt.Errorf("unsupported convention type: %s", cfg.Convention)
	}

	return provider, nil
}
