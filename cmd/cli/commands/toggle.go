package commands

import (
	"fmt"
	"idasen-desk/internal/config"
	"idasen-desk/internal/desk"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Toggle(_ *cli.Context, args InputFlags) (err error) {
	configuration, err := config.Load(args.ConfigPath)
	if err != nil {
		return err
	}

	var d *desk.Desk
	if d, err = desk.NewDesk(
		configuration.LocalName,
		configuration.ConnectionAddress,
		true,
	); err != nil {
		return fmt.Errorf("failed to create new desk instance, %w", err)
	}

	height, baseHeightErr := d.GetHeight()
	if baseHeightErr != nil {
		return baseHeightErr
	}

	sitHeight := configuration.SitHeight
	if args.SitHeight > 0 {
		sitHeight = args.SitHeight
	}

	standHeight := configuration.StandHeight
	if args.StandHeight > 0 {
		standHeight = args.StandHeight
	}

	log.Printf("connected to %s", d.Name())

	// If we got this far with no options then lets go and locate the location,
	// which is the furthest away and go for that, e.g., toggle between standing
	// and or sitting.
	sitDifference := math.Abs(sitHeight - height)
	standDifference := math.Abs(standHeight - height)

	if sitDifference > standDifference {
		return d.MoveToTarget(sitHeight)
	} else {
		return d.MoveToTarget(standHeight)
	}

}
