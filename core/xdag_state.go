package core

import "xdago/common"

type XdagState struct {
	cmd  common.StateType
	temp int
}

func (x *XdagState) SetState(s common.StateType) {
	x.cmd = s
}

func (x *XdagState) TempState(s common.StateType) {
	x.temp = int(x.cmd)
	x.cmd = s
}

func (x *XdagState) Rollback() {
	if x.temp != -1 {
		x.cmd = common.StateType(x.temp)
		x.temp = -1
	}
}

func (x XdagState) ToString() string {
	switch x.cmd {
	case common.INIT:
		return "Pool Initializing...."
	case common.KEYS:
		return "Generating keys..."

	case common.REST:
		return "The local storage is corrupted. Resetting blocks engine."

	case common.LOAD:
		return "Loading blocks from the local storage."

	case common.STOP:
		return "locks loaded. Waiting for 'run' command."

	case common.WDST:
		return "Trying to connect to the  dev network."

	case common.WTST:
		return "Trying to connect to the test network."

	case common.WAIT:
		return "Trying to connect to the main network"

	case common.CDST:
		return "Connected to the  dev network. Synchronizing."

	case common.CTST:
		return "Connected to the test network. Synchronizing."

	case common.CONN:
		return "Connected to the main network. Synchronizing."

	case common.SDST:
		return "Synchronized with the  dev network. Normal testing."

	case common.STST:
		return "Synchronized with the test network. Normal testing."

	case common.SYNC:
		return "Synchronized with the main network. Normal operation."

	case common.XFER:
		return "Waiting for transfer to complete."
	default:
		return "Abnormal State"
	}
}
