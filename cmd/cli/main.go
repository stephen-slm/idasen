package main

import (
	"fmt"
	"idasen-desk/internal/config"
	"idasen-desk/internal/desk"
	"math"
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

func run(cliArguments inputFlags) (err error) {
	configuration := &config.Configuration{}
	if err = configuration.Load(cliArguments.ConfigPath); err != nil {
		return err
	}

	var d *desk.Desk
	if d, err = desk.NewDesk(configuration.ConnectionAddress, true); err != nil {
		return fmt.Errorf("failed to create new desk instance, %w", err)
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
		return d.MoveToTarget(cliArguments.Target)
	}

	if cliArguments.Stand {
		return d.MoveToTarget(cliArguments.StandHeight)
	}

	if cliArguments.Sit {
		return d.MoveToTarget(cliArguments.SitHeight)
	}

	// If we got this far with no options then lets go and locate the location,
	// which is the furthest away and go for that, e.g., toggle between standing
	// and or sitting.
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
		Name:  "Idasen CLI",
		Usage: "A simple CLI to interface with the Idasen desk",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Specify the path to the configuration file.",
				Value:       "./.desk.yml",
				Destination: &flags.ConfigPath,
			},
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
