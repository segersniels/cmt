package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type Emoji struct {
	Emoji       string `json:"emoji"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Response struct {
	Gitmojis []Emoji `json:"gitmojis"`
}

func isCached(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func fetchFromCache(path string) ([]Emoji, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var response Response
	err = json.NewDecoder(file).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Gitmojis, nil
}

func writeToCache(path string, response Response) error {
	directory := filepath.Dir(path)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&response)
	if err != nil {
		return err
	}

	return nil
}

type Gitmoji struct {
	notation string
}

func NewGitmoji(notation string) *Gitmoji {
	return &Gitmoji{notation}
}

func (g *Gitmoji) fetch() ([]Emoji, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(dirname, ".config", "cmt", "gitmojis.json")
	if isCached(path) {
		return fetchFromCache(path)
	}

	res, err := http.Get("https://gitmoji.dev/api/gitmojis")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if err := writeToCache(path, response); err != nil {
		return nil, err
	}

	return response.Gitmojis, nil
}

// Remove zero-width joiners (U+200D) from emojis since they may not render properly in terminals
func parseEmoji(emoji Emoji) string {
	runes := []rune(emoji.Emoji)
	result := make([]rune, 0, len(runes))

	for _, r := range runes {
		if r != '\u200D' {
			result = append(result, r)
		}
	}

	return string(result)
}

func (g *Gitmoji) ask() (string, error) {
	emojis, err := g.fetch()
	if err != nil {
		return "", err
	}

	options := make([]huh.Option[string], 0, len(emojis))
	for _, e := range emojis {
		optionText := fmt.Sprintf("%s  - %s", parseEmoji(e), e.Description)
		optionValue := e.Emoji
		if g.notation != "emoji" {
			optionValue = e.Code
		}
		options = append(options, huh.NewOption(optionText, optionValue))
	}

	var emoji string
	if err := huh.NewSelect[string]().
		Title("Select the type of commit").
		Options(options...).
		Value(&emoji).
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

	return emoji, nil
}

func (g *Gitmoji) Construct() (string, error) {
	gitmoji, err := g.ask()
	if err != nil {
		return "", err
	}

	message, err := promptForMessage()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", gitmoji, message), nil
}

func (g *Gitmoji) Type() ConventionType {
	return GitmojiConvention
}

var _ Provider = (*Gitmoji)(nil)
