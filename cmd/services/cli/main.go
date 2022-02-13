package main

import (
	"math"
	"os"
	"time"

	"idasen-desk/internal/desk"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type inputFlags struct {
	Verbose     bool    `json:"verbose"`
	Sit         bool    `json:"sit"`
	SitHeight   float64 `json:"sit_height"`
	Stand       bool    `json:"stand"`
	StandHeight float64 `json:"stand_height"`
	Target      float64 `json:"move"`
	Monitor     bool    `json:"monitor"`
}

func run(cliArguments inputFlags) error {
	log.SetLevel(log.InfoLevel)

	if cliArguments.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("cli_arguments", cliArguments).Debug("input arguments")

	personalDeskAddress := "C2:6D:5B:C4:17:12"
	d := desk.NewDesk(personalDeskAddress)

	if err := d.Connect(); err != nil {
		return err
	}

	height, baseHeightErr := d.GetHeight()
	if baseHeightErr != nil {
		return baseHeightErr
	}

	log.Printf("connected to %s", d.Name())

	if cliArguments.Monitor {
		return d.Monitor()
	}

	if cliArguments.Target > 0 {
		return d.MoveToTarget(cliArguments.Target / 100.0)
	}

	if cliArguments.Stand {
		return d.MoveToTarget(cliArguments.StandHeight)
	}

	if cliArguments.Sit {
		return d.MoveToTarget(cliArguments.SitHeight)
	}

	// If we got this far with no options then lets go and locate the location
	// which is the furthest away and go for that, e.g toggle between stand
	// and sit.
	sitDifference := math.Abs(cliArguments.SitHeight - height)
	standDifference := math.Abs(cliArguments.StandHeight - height)

	if sitDifference > standDifference {
		return d.MoveToTarget(cliArguments.SitHeight)
	} else {
		return d.MoveToTarget(cliArguments.StandHeight)
	}

}

func main() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	flags := inputFlags{}

	app := &cli.App{
		Name:      "Idasen CLI",
		HelpName:  "",
		Usage:     "A simple CLI to interface with the Idasen desk",
		UsageText: "",
		ArgsUsage: "",
		Version:   "",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "Enable verbose logging",
				EnvVars:     []string{"VERBOSE"},
				Destination: &flags.Verbose,
			},
			&cli.Float64Flag{
				Name:        "stand-height",
				Usage:       "The target end height for standing",
				Value:       1.12,
				Destination: &flags.StandHeight,
			},
			&cli.Float64Flag{
				Name:        "sit-height",
				Usage:       "The target end height for sitting",
				Value:       0.74,
				Destination: &flags.SitHeight,
			},

			&cli.BoolFlag{
				Name:        "sit",
				Usage:       "Put the desk into a sitting position",
				Destination: &flags.Sit,
			},
			&cli.BoolFlag{
				Name:        "stand",
				Usage:       "Put the desk into a standing position",
				Destination: &flags.Stand,
			},
			&cli.Float64Flag{
				Name:        "target",
				Aliases:     []string{"t"},
				Usage:       "Move the desk into the target position",
				Destination: &flags.Target,
			},
			&cli.BoolFlag{
				Name:        "monitor",
				Aliases:     []string{"m"},
				Usage:       "Monitor the movement of the desk during manual movement",
				Destination: &flags.Monitor,
			},
		},
		EnableBashCompletion: false,
		HideHelp:             false,
		HideHelpCommand:      false,
		HideVersion:          false,
		BashComplete:         nil,
		Before:               nil,
		After:                nil,
		Action: func(_ *cli.Context) error {
			return run(flags)
		},
		CommandNotFound:        nil,
		OnUsageError:           nil,
		Compiled:               time.Time{},
		Authors:                nil,
		Copyright:              "",
		Reader:                 nil,
		Writer:                 nil,
		ErrWriter:              nil,
		ExitErrHandler:         nil,
		Metadata:               nil,
		ExtraInfo:              nil,
		CustomAppHelpTemplate:  "",
		UseShortOptionHandling: false,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
