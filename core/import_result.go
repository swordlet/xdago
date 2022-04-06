package core

import "xdago/common"

type ImportResult struct {
	Status    common.ImportStatus
	HashLow   common.Hash
	ErrorInfo string
}

func (ir ImportResult) IsNormal() bool {
	return ir.Status == common.IMPORTED_NOT_BEST ||
		ir.Status == common.IMPORTED_BEST ||
		ir.Status == common.IMPORT_EXIST
}
