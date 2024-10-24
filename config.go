package main

import (
	"encoding/json"
	"os"

	log "github.com/charmbracelet/log"
)

type ConventionType int

const (
	ConventionalCommitConvention ConventionType = iota
	GitmojiConvention
)

var conventionType = map[ConventionType]string{
	ConventionalCommitConvention: "conventional-commit",
	GitmojiConvention:            "gitmoji",
}

func (c ConventionType) String() string {
	return conventionType[c]
}

func (c *ConventionType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case conventionType[ConventionalCommitConvention]:
		*c = ConventionalCommitConvention
	case conventionType[GitmojiConvention]:
		*c = GitmojiConvention
	default:
		log.Fatal("unsupported convention", "type", s)
	}

	return nil
}

func (c ConventionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(conventionType[c])
}

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
