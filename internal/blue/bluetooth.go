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

// MustParseUUID returns the bluetooth UUID based on the input value otherwise
// panics with the error if it faults.
func MustParseUUID(value string) bluetooth.UUID {
	v, err := bluetooth.ParseUUID(value)

	if err != nil {
		log.Fatal(err)
	}

	return v
}
