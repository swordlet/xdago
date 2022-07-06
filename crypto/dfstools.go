package crypto

//#cgo CFLAGS: -I../clib/dfstools/src
//#cgo linux,amd64 LDFLAGS:-L${SRCDIR}/../clib -lxdag_crypto_Linux -lm -lstdc++
//#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/../clib -lxdag_crypto_Darwin -lm -lstdc++
//#cgo windows,amd64 LDFLAGS:-L${SRCDIR}/../clib -lxdag_crypto_Windows -static -static-libgcc -static-libstdc++
//#include <stdlib.h>
//#include "wrapper.h"
import "C"
import "unsafe"

func LoadDnetKeys(p unsafe.Pointer, length int) int {
	return C.loadDnetKeys(p, length)
}

func DnetCryptInit() int {
	return C.dnetCryptInit()
}
