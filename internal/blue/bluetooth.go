package blue

import (
	"time"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// ConnectToDevice attempts to make a direct bluetooth connection directly to
// a device based on its address. This requires having a pre-defined
// understanding of the device, e.g connected to this adapter already.
func ConnectToDevice(address bluetooth.Address) (*bluetooth.Device, error) {
	if enabledErr := adapter.Enable(); enabledErr != nil {
		return nil, enabledErr
	}

	device, err := adapter.Connect(address, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(time.Second),
	})

	return device, err
}

// UniqueScan starts scanning for devices to connect to and only passes back
// unique values, stopping the chance of duplicates by address.
//
// `done` must be closed to stop the scanning.
func UniqueScan(done chan struct{}) (chan *bluetooth.ScanResult, error) {
	if enabledErr := adapter.Enable(); enabledErr != nil {
		return nil, enabledErr
	}

	output := make(chan *bluetooth.ScanResult)
	uniqueTracker := map[string]struct{}{}

	// Adapter scan runs in the background and will result in it running
	// forever if not cancelled. This allows us to set up our own channel flow
	// and pipe the data back.
	go func() {
		_ = adapter.Scan(func(a *bluetooth.Adapter, result bluetooth.ScanResult) {
			if _, exists := uniqueTracker[result.Address.String()]; exists {
				return
			}

			uniqueTracker[result.Address.String()] = struct{}{}
			output <- &result
		})
	}()

	// Go and wait until `done` is closed and trigger the stop of the scanning,
	// This should be a clean way to stop the chance of this scanning running
	// forever.
	go func() {
		defer adapter.StopScan()
		defer close(output)
		<-done
	}()

	return output, nil
}

// MustParseUUID returns the bluetooth UUID based on the input value otherwise
// panics with the error if it faults.
func MustParseUUID(value string) bluetooth.UUID {
	v, err := bluetooth.ParseUUID(value)

	if err != nil {
		log.Fatal(err)
	}

	return v
}
