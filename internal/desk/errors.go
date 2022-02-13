package desk

var targetHeightTooHigh = &deskError{msg: "target height too high"}
var targetHeightTooLow = &deskError{msg: "target height too low"}
var deskMoveSafetyKickIn = &deskError{msg: "desk move safety kicked in."}
var bluetoothError = &deskError{msg: "bluetooth error"}

// circuitError is used for internally generated errors
type deskError struct {
	msg string
}

func (m *deskError) Error() string {
	return m.msg
}
