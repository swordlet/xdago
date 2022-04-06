package common

type ImportStatus byte

const (
	IMPORT_ERROR ImportStatus = iota
	IMPORT_EXIST
	NO_PARENT
	INVALID_BLOCK
	IMPORTED_NOT_BEST
	IMPORTED_BEST
)
