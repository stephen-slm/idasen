package main

import (
	"fmt"
	"idasen-desk/internal/config"
	"idasen-desk/internal/desk"
	"math"

	log "github.com/sirupsen/logrus"
)

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
