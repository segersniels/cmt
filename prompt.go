package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/charmbracelet/log"
)

func promptForMessage() (string, error) {
	var msg string
	if err := huh.NewInput().Title("Enter your commit message").Value(&msg).Run(); err != nil {
		return "", err
	}

	if len(msg) == 0 {
		fmt.Println("Message cannot be empty")
		os.Exit(0)
	}

	cfg, err := ReadConfig()
	if err != nil {
		log.Debug("error reading config", "error", err)
		return msg, nil
	}

	// Ensure the first letter of the message is uppercase if requested by user
	if cfg.Uppercase && len(msg) > 0 {
		msg = strings.ToUpper(msg[:1]) + msg[1:]
	}

	return msg, nil
}
