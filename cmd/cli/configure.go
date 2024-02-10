package main

import (
	"fmt"
	"idasen-desk/internal/blue"
	"idasen-desk/internal/config"
	"strings"

	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"tinygo.org/x/bluetooth"
)

func handleSelectionOfDevice(output *bluetooth.ScanResult, configurationPath string) {
	c := config.Configuration{}
	if err := c.Load(configurationPath); err != nil {
		log.Fatal(err)
	}

	c.ConnectionAddress = output.Address.String()
	if err := c.Save(configurationPath); err != nil {
		log.Fatal(err)
	}
}

func configure(_ *cli.Context, cliArguments inputFlags) (err error) {
	done := make(chan struct{})

	values, err := blue.UniqueScan(done)
	if err != nil {
		return fmt.Errorf("failed to start unique scan, %w", err)
	}

	app := tview.NewApplication()
	list := tview.NewList().
		AddItem("Quit", "Press to exit", 'q', func() { app.Stop() })

	var scanResults []*bluetooth.ScanResult

	go func() {
		for value := range values {
			scanResults = append(scanResults, value)

			name := value.Address.String()
			if value.LocalName() != "" {
				name = fmt.Sprintf("%s: %s", value.Address.String(),
					strings.TrimSpace(value.LocalName()))
			}

			list.AddItem(name, "", rune(96+len(scanResults)), func() {
				handleSelectionOfDevice(scanResults[list.GetCurrentItem()-1], cliArguments.ConfigPath)
				defer close(done)
				defer app.Stop()
			})

			app.Draw()
		}
	}()

	return app.SetRoot(list, true).Run()
}
