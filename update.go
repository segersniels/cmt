package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-version"
)

func update() error {
	version, err := fetchLatestVersion()
	if err != nil {
		return err
	}

	ldflags := fmt.Sprintf("-w -s -X main.AppVersion=%s -X main.AppName=cmt", version)
	origin := fmt.Sprintf("github.com/segersniels/cmt@%s", version)
	args := []string{"install", "-ldflags", ldflags, origin}
	cmd := exec.Command("go", args...)

	// If the GOBIN environment variable is not set, set it to `/usr/local/bin/`.
	// Users can override it by setting GOBIN in their environment.
	if os.Getenv("GOBIN") == "" {
		cmd.Env = append(os.Environ(), "GOBIN=/usr/local/bin/")
	}

	// Capture the output of the command
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

type Failure struct {
	Message string `json:"message"`
}

type Success struct {
	TagName string `json:"tag_name"`
}

func fetchLatestVersion() (*version.Version, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/segersniels/cmt/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the response is a rate limit error
	if resp.StatusCode != http.StatusOK {
		var result Failure
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		return nil, errors.New(result.Message)
	}

	var result Success
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	latestVersion, err := version.NewVersion(result.TagName)
	if err != nil {
		return nil, err
	}

	return latestVersion, nil
}

func checkIfNewVersionIsAvailable() error {
	// If the AppVersion is not set, we don't need to check for an update.
	// This is the case when running `go install` or `go build` without
	// specifying the LDFLAGS properly.
	if AppVersion == "" {
		return nil
	}

	currentVersion, err := version.NewVersion(AppVersion)
	if err != nil {
		return err
	}

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		return err
	}

	if latestVersion.GreaterThan(currentVersion) {
		fmt.Printf("A new version of %s is available (%s). Run `cmt update` to update.\n\n", AppName, latestVersion)
	}

	return nil
}
