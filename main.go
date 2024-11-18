package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	log "github.com/charmbracelet/log"
	updater "github.com/segersniels/updater"
	"github.com/urfave/cli/v2"
)

var (
	AppName    string
	AppVersion string
)

func main() {
	debug := os.Getenv("DEBUG")
	if debug != "" {
		log.SetLevel(log.DebugLevel)
	}

	upd := updater.NewUpdater(AppName, AppVersion, "segersniels")
	err := upd.CheckIfNewVersionIsAvailable()
	if err != nil {
		log.Debug("Failed to check for latest release", "error", err)
	}

	app := &cli.App{
		Name:    AppName,
		Usage:   "Write commit messages independent of convention",
		Version: AppVersion,
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "Update to the latest version",
				Action: func(ctx *cli.Context) error {
					return upd.Update()
				},
			},
			{
				Name:  "init",
				Usage: "Initialize a new project",
				Action: func(ctx *cli.Context) error {
					var (
						convention ConventionType
						uppercase  bool
					)

					form := huh.NewForm(
						huh.NewGroup(
							huh.NewSelect[ConventionType]().
								Title("Which convention do you want to use?").
								Description("A lot of projects use Conventional Commits, but Gitmoji is also a popular choice.").
								Options(
									huh.NewOption("Conventional Commits", ConventionalCommitConvention),
									huh.NewOption("Gitmoji", GitmojiConvention),
								).Value(&convention),
						),
						huh.NewGroup(
							huh.NewConfirm().
								Title("Uppercase first letter of commit message?").
								Description("This will automatically uppercase the first letter of your commit message.").
								Value(&uppercase),
						),
					)

					if err := form.Run(); err != nil {
						return err
					}

					return WriteConfig(Config{Convention: convention, Uppercase: uppercase})
				},
			},
			{
				Name:    "commit",
				Aliases: []string{"c"},
				Usage:   "Create a new commit",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "allow-empty",
						Usage: "Allow the commit message to be empty",
					},
				},
				Action: func(ctx *cli.Context) error {
					convention, err := determineConvention()
					if err != nil {
						return err
					}

					msg, err := convention.Construct()
					if err != nil {
						return err
					}

					args := []string{"commit", "-m", msg}
					if ctx.Bool("allow-empty") {
						args = append(args, "--allow-empty")
					}

					cmd := exec.Command("git", args...)
					output, err := cmd.CombinedOutput()
					if err != nil {
						return fmt.Errorf("%s: %s", err, string(output))
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
