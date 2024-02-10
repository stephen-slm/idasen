package desk

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"tinygo.org/x/bluetooth"

	"idasen-desk/internal/blue"
)

var (
	UuidHeight         = blue.MustParseUUID("99fa0021-338a-1024-8a49-009c0215f78a")
	UuidCommand        = blue.MustParseUUID("99fa0002-338a-1024-8a49-009c0215f78a")
	UuidReferenceInput = blue.MustParseUUID("99fa0031-338a-1024-8a49-009c0215f78a")

	// UuidAdvSvc - Not currently used but can be used to determine if the given device is a
	// desk or not. If it is a desk then the deskService (services_uuid) list will contain this uuid.
	UuidAdvSvc = blue.MustParseUUID("99fa0001-338a-1024-8a49-009c0215f78a")
)

const (
	MaxHeight = 1.27
	MinHeight = 0.62
)

type Direction int

const (
	UNKNOWN Direction = iota
	UP
	DOWN
)

type Desk struct {
	name    string
	address string

	device                 *bluetooth.Device
	deskService            []bluetooth.DeviceService
	serviceCharacteristics []bluetooth.DeviceCharacteristic
}

func NewDesk(address string) *Desk {
	return &Desk{
		name:                   "Desk",
		address:                address,
		device:                 nil,
		deskService:            nil,
		serviceCharacteristics: nil,
	}
}

// Connect will attempt to connect to the desk via bluetooth.
func (d *Desk) Connect() (err error) {
	mac, _ := bluetooth.ParseMAC(d.address)
	address := bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: mac}}

	if d.device, err = blue.ConnectToDevice(address); err != nil {
		return err
	}

	if d.deskService, err = d.device.DiscoverServices(nil); err != nil {
		return err
	}

	for _, service := range d.deskService {
		characteristics, _ := service.DiscoverCharacteristics(nil)
		d.serviceCharacteristics = append(d.serviceCharacteristics, characteristics...)
	}

	_, err = d.GetHeight()
	return err
}

func (d *Desk) Name() string {
	if d.name == "" {
		return "desk"
	}

	return d.name
}

// GetHeight returns the current height of the desk by direct 1:1 communication
// and no by a notification. This includes some delay.
func (d *Desk) GetHeight() (float64, error) {
	characteristic := d.getCharacteristic(UuidHeight)

	if characteristic == nil {
		return 0, fmt.Errorf("does not have required characteristic: %s", UuidHeight.String())
	}

	data := make([]byte, 4)
	_, err := characteristic.Read(data)

	return bytesToMeters(data), err
}

// Stop tells the desk to stop moving.
//
// The desk does not stop automatically unless the safety kicks in.
func (d *Desk) Stop() error {
	commandChar := d.getCharacteristic(UuidCommand)
	referenceChar := d.getCharacteristic(UuidReferenceInput)

	commandStop := []byte{0xFF, 0x00}
	commandRefInput := []byte{0x01, 0x80}

	var eg errgroup.Group

	eg.Go(func() error {
		_, err := commandChar.WriteWithoutResponse(commandStop)
		return err
	})

	eg.Go(func() error {
		_, err := referenceChar.WriteWithoutResponse(commandRefInput)
		return err
	})

	return eg.Wait()
}

// Monitor purely listens to the notification events fired by the desk and
// prints them to the display. Existing only on the control+c calls or hard
// exits.
func (d *Desk) Monitor() error {
	heightCharacteristic := d.getCharacteristic(UuidHeight)

	if err := heightCharacteristic.EnableNotifications(func(buf []byte) {
		log.Infof("%f", bytesToMeters(buf))
	}); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}

// MoveToTarget move the desk to the specified target float value. Within the
// constraints of the device min value and max value.
func (d *Desk) MoveToTarget(target float64) error {
	if target > MaxHeight {
		return targetHeightTooHigh
	} else if target < MinHeight {
		return targetHeightTooLow
	}

	heightCharacteristic := d.getCharacteristic(UuidHeight)
	currentHeight, err := d.GetHeight()

	if err != nil {
		return fmt.Errorf("failed to get desk height, %w", err)
	}

	previousHeight := currentHeight
	willMoveUp := target > previousHeight

	log.Infof("moving desk from %.2f to %.2f", previousHeight, target)

	var mu sync.RWMutex
	getHeight := func() float64 {
		mu.RLock()
		defer mu.RUnlock()
		return currentHeight
	}

	setHeight := func(value float64) {
		mu.Lock()
		defer mu.Unlock()
		currentHeight = value
	}

	// Use the implemented notification characteristics to get real time
	// updates on the position of the desk. Allowing the loop iteration to only
	// care about directional control.
	if err = heightCharacteristic.EnableNotifications(func(buf []byte) {
		log.Debugf("desk height notification: %f", bytesToMeters(buf))
		setHeight(bytesToMeters(buf))
	}); err != nil {
		return fmt.Errorf("failed to configure desk hight notifications, %w", err)
	}

	for {
		loopHeight := getHeight()

		differenceRaw := target - loopHeight
		differenceAbs := math.Abs(differenceRaw)

		log.Debugf("target=%f, current_height=%f previous_height=%f, difference=%f",
			target, loopHeight, previousHeight, differenceRaw)

		// The device has a moving action to protect the user if it applies
		// pressure to something when moving. This will result in the desk
		// moving in the opposite direction when the device detects something.
		// Moving out th way. If we detect this, stop.
		//
		// Only if our difference is not nothing, meaning we are not doing a
		// minor correction.
		if (loopHeight < previousHeight && willMoveUp ||
			loopHeight > previousHeight && !willMoveUp) &&
			differenceAbs > 0.010 {
			log.Errorf("stopped moving because desk safety feature kicked in.")
			return deskMoveSafetyKickIn
		}

		// If we're either less than 10mm then we need to stop every iteration
		// so that we don't overshoot
		if differenceAbs < 0.010 {
			log.Debugf("hit differnce made: height: %f - difference: %f",
				loopHeight, differenceRaw)

			if stopErr := d.Stop(); stopErr != nil {
				return stopErr
			}
		}

		// If we are within our tolerance for moving the desk then we can go and
		// stop. Additionally ensure to stop first to keep in line with our
		// tolerance. Otherwise, a shift in the difference could occur when
		// pulling the final destination.
		//
		// within 5mm
		if differenceAbs <= 0.005 {
			if stopErr := d.Stop(); stopErr != nil {
				return stopErr
			}

			// Sleep for the duration of a possible upper limit of a step
			// duration. This duration was determined from a single `MOVE`
			// operation.
			time.Sleep(time.Millisecond * 100)
			log.Infof("reached target of %.3f, actual: %.3f", target, getHeight())
			return nil
		}

		operation := UP
		if differenceRaw < 0.0 {
			operation = DOWN
		}

		// Attempt to move into the correct direction, if it faults, attempt to
		// stop and return the errors.
		if err = d.MoveDirection(operation); err != nil {
			return errors.Join(err, d.Stop())
		}

		previousHeight = loopHeight
	}
}

// MoveDirection Based on the provided direction, the desk will be told to start
// moving up or start moving down. A move action will only occur for a 1-second
// interval, which is configured by the desk.
func (d *Desk) MoveDirection(direction Direction) error {
	actionArgs := []uint8{0x47, 0x00}

	if direction == DOWN {
		actionArgs = []uint8{0x46, 0x00}
	}

	if _, err := d.getCharacteristic(UuidCommand).WriteWithoutResponse(actionArgs); err != nil {
		return fmt.Errorf("%s: %w", bluetoothError.Error(), err)
	}

	return nil
}

func (d *Desk) getCharacteristic(uuid bluetooth.UUID) *bluetooth.DeviceCharacteristic {
	for _, characteristic := range d.serviceCharacteristics {
		if characteristic.UUID() == uuid {
			return &characteristic
		}
	}

	return nil
}

// Converts the raw height response from the desk into meters.
func bytesToMeters(raw []uint8) float64 {
	var highByte int
	var lowByte int

	highByte = int(raw[1])
	lowByte = int(raw[0])

	number := (highByte << 8) + lowByte
	return (float64(number) / 10000.0) + MinHeight
}
