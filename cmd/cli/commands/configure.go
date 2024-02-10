package commands

import (
	"fmt"
	"idasen-desk/internal/blue"
	"idasen-desk/internal/config"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"tinygo.org/x/bluetooth"
)

func handleSelectionOfDevice(o *bluetooth.ScanResult, c *config.Configuration, path string) {
	c.ConnectionAddress = o.Address.String()
	c.LocalName = o.LocalName()

	if err := c.Save(path); err != nil {
		log.Fatal(err)
	}
}

func HeaderPrimitive(c *config.Configuration) tview.Primitive {
	name := c.ConnectionAddress

	if c.LocalName != "" {
		name = fmt.Sprintf("%s (%s)", name, c.LocalName)
	}

	value := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("current: %s", name))

	value.SetBackgroundColor(tcell.ColorDefault)
	return value
}

func Configure(_ *cli.Context, args InputFlags) (err error) {
	configuration, err := config.Load(args.ConfigPath)
	if err != nil {
		return err
	}

	done := make(chan struct{})

	values, err := blue.UniqueScan(done)
	if err != nil {
		return fmt.Errorf("failed to start unique scan, %w", err)
	}

	list := tview.NewList()

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(HeaderPrimitive(configuration), 0, 0, 1, 3, 0, 0, false)

	grid.AddItem(list, 1, 0, 1, 3, 0, 0, true)

	grid.SetBackgroundColor(tcell.ColorDefault)
	list.SetBackgroundColor(tcell.ColorDefault)

	var scanResults []*bluetooth.ScanResult
	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
		case tcell.Key('q'):
			app.Stop()
			return event
		}
		return event
	})

	go func() {
		for value := range values {
			scanResults = append(scanResults, value)

			name := value.Address.String()
			if value.LocalName() != "" {
				name = fmt.Sprintf("%s: %s", value.Address.String(),
					strings.TrimSpace(value.LocalName()))
			}

			list.AddItem(name, "", rune(96+len(scanResults)), func() {
				handleSelectionOfDevice(
					scanResults[list.GetCurrentItem()],
					configuration,
					args.ConfigPath,
				)
				defer close(done)
				defer app.Stop()
			})

			app.Draw()
		}
	}()

	return app.SetRoot(grid, true).
		EnableMouse(false).
		Run()
}
