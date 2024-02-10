package commands

import (
	"fmt"
	"idasen-desk/internal/config"
	"idasen-desk/internal/desk"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Height(_ *cli.Context, args InputFlags) (err error) {
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

	height, err := d.GetHeight()
	if err != nil {
		return err
	}

	log.Printf("height: %.2f", height)
	return nil
}
