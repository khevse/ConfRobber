package zlibwrapper

/*
#cgo windows LDFLAGS: -L"${SRCDIR}/cpp_src" -lzlibwrapper -lzlibstatic -lgcc -lstdc++
#include "cpp_src/zlibwrapper.h"
*/
import "C"
import "unsafe"
import _ "runtime/cgo"
import "errors"

func Compress(sourceData []byte) (result []byte, err error) {

	if len(sourceData) == 0 {
		return
	}

	var data *C.BYTE
	var size C.usize

	if C.compress_data((*C.BYTE)(&sourceData[0]), C.usize(len(sourceData)), &data, &size) {
		result = C.GoBytes(unsafe.Pointer(data), C.int(size))
		free(data)
	} else {
		err = errors.New("Compress error.")
	}

	return
}

func Decompress(sourceData []byte) (result []byte, err error) {

	if len(sourceData) == 0 {
		return
	}

	var data *C.BYTE
	var size C.usize

	if C.decompress_data((*C.BYTE)(&sourceData[0]), C.usize(len(sourceData)), &data, &size) {
		result = C.GoBytes(unsafe.Pointer(data), C.int(size))
		free(data)
	} else {
		err = errors.New("Decompress error.")
	}

	return
}

func free(data *C.BYTE) {
	C.free_data(&data)
}
