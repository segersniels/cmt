package main

import (
	"encoding/json"
	"os"
)

type ConventionType string

const (
	ConventionalCommitConvention ConventionType = "conventional-commit"
	GitmojiConvention            ConventionType = "gitmoji"
)

type Config struct {
	Convention ConventionType `json:"convention"`
	Uppercase  bool           `json:"uppercase"`
}

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(".cmtrc.json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func WriteConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(".cmtrc.json", data, 0644)
}
