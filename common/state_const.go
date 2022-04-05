package common

type StateType byte

const (

	//The pool is initializing......
	INIT StateType = iota

	//wallet generating keys....
	KEYS

	//The local storage is corrupted. Resetting blocks engine.
	REST

	//Loading blocks from the local storage.
	LOAD

	//Blocks loaded. Waiting for 'run' command.
	STOP

	//Trying to connect to the  dev network.
	WDST

	//Trying to connect to the test network.
	WTST

	//Trying to connect to the main network.
	WAIT

	//Connected to the  dev network. Synchronizing.
	CDST

	//Connected to the test network. Synchronizing.
	CTST

	//Connected to the main network. Synchronizing.
	CONN

	//Synchronized with the  dev network. Normal testing.
	SDST

	//Synchronized with the test network. Normal testing.
	STST

	//Synchronized with the main network. Normal operation.
	SYNC

	//Waiting for transfer to complete.
	XFER
)
