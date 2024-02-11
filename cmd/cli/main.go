package main

import (
	"idasen-desk/cmd/cli/commands"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	flags := commands.InputFlags{}

	sharedFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Usage:       "Enable verbose logging",
			EnvVars:     []string{"VERBOSE"},
			Destination: &flags.Verbose,
		},
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "Specify the path to the configuration file.",
			Value:       "./.desk.yml",
			Destination: &flags.ConfigPath,
		},
	}

	standHeightFlag := &cli.Float64Flag{
		Name:        "stand-height",
		Aliases:     nil,
		Usage:       "The optional target end height for standing",
		EnvVars:     nil,
		FilePath:    "",
		Required:    false,
		Hidden:      false,
		DefaultText: "Configuration File Value",
		Destination: &flags.StandHeight,
		HasBeenSet:  false,
	}

	sitHeightFlag := &cli.Float64Flag{
		Name:        "sit-height",
		Usage:       "The optional target end height for sitting",
		DefaultText: "Configuration File Value",
		Destination: &flags.SitHeight,
	}

	cliCommands := []*cli.Command{{
		Name:  "configure",
		Usage: "configure the device to connect to.",
		Flags: append([]cli.Flag{standHeightFlag}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Configure(context, flags)
		},
	}, {
		Name:  "stand",
		Usage: "Move the desk to the configured standing position.",
		Flags: append([]cli.Flag{standHeightFlag}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Stand(context, flags)
		},
	}, {
		Name:  "sit",
		Usage: "Move the desk to the configured sitting position.",
		Flags: append([]cli.Flag{sitHeightFlag}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Sit(context, flags)
		},
	}, {
		Name:      "position",
		Usage:     "Move the desk to the provided position value.",
		ArgsUsage: "[position]",
		Flags:     append([]cli.Flag{}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Position(context, flags)
		},
	}, {
		Name:  "toggle",
		Usage: "Toggle the desk height between standing and sitting.",
		Flags: append([]cli.Flag{sitHeightFlag, standHeightFlag}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Toggle(context, flags)
		},
	}, {
		Name:  "monitor",
		Usage: "Monitor and log the position of the desk as it moves",
		Flags: append([]cli.Flag{}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Monitor(context, flags)
		},
	}, {
		Name:  "height",
		Usage: "Get the current height of the desk.",
		Flags: append([]cli.Flag{}, sharedFlags...),
		Action: func(context *cli.Context) error {
			return commands.Height(context, flags)
		},
	}}

	app := &cli.App{
		Name:     "Idasen CLI",
		Usage:    "A simple CLI to interface with the Idasen desk",
		Commands: cliCommands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
