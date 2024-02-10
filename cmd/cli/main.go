package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type inputFlags struct {
	ConfigPath  string  `json:"config_path"`
	Verbose     bool    `json:"verbose"`
	Sit         bool    `json:"sit"`
	SitHeight   float64 `json:"sit_height"`
	Stand       bool    `json:"stand"`
	StandHeight float64 `json:"stand_height"`
	Target      float64 `json:"move"`
	Monitor     bool    `json:"monitor"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})
	flags := inputFlags{}

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

	app := &cli.App{
		Name:  "Idasen CLI",
		Usage: "A simple CLI to interface with the Idasen desk",

		Commands: []*cli.Command{{
			Name:  "configure",
			Usage: "configure the device to connect to.",
			Action: func(context *cli.Context) error {
				return configure(context, flags)
			},
			OnUsageError: nil,
			Subcommands:  nil,
			Flags:        append([]cli.Flag{}, sharedFlags...),
		}},

		Flags: append([]cli.Flag{
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
		}, sharedFlags...),
		Action: func(_ *cli.Context) error {
			log.SetLevel(log.InfoLevel)

			if flags.Verbose {
				log.SetLevel(log.DebugLevel)
			}

			log.WithField("arguments", flags).
				Debug("input arguments")

			return run(flags)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
