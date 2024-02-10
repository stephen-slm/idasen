package commands

import (
	"fmt"
	"idasen-desk/internal/config"
	"idasen-desk/internal/desk"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Stand(_ *cli.Context, args InputFlags) (err error) {
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

	standHeight := configuration.StandHeight
	if args.StandHeight > 0 {
		standHeight = args.StandHeight
	}

	log.Printf("connected to %s", d.Name())
	return d.MoveToTarget(standHeight)
}
