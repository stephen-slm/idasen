package blue

import (
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func FindDevice(address bluetooth.Address) (*bluetooth.Device, error) {
	if enabledErr := adapter.Enable(); enabledErr != nil {
		return nil, enabledErr
	}

	device, err := adapter.Connect(address, bluetooth.ConnectionParams{
		ConnectionTimeout: 0,
		MinInterval:       0,
		MaxInterval:       0,
	})

	return device, err
}

func ParseUUID(value string) bluetooth.UUID {
	v, _ := bluetooth.ParseUUID(value)
	return v
}
