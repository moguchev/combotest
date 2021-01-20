package closessl

/*
#cgo CFLAGS: -I${SRCDIR}/../../../dist/closessl
#cgo LDFLAGS: -L${SRCDIR}/../../../dist/closessl -lclosessl

#include <closessl.h>
*/
import "C"
import "unsafe"

func Encrypt(data []byte) []byte {
	// call C function:
	// void closessl_encrypt(uint8_t *buf, size_t len);
	C.closessl_encrypt((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	return data
}

func Decrypt(data []byte) []byte {
	// call C function:
	// void closessl_encrypt(uint8_t *buf, size_t len);
	C.closessl_decrypt((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	return data
}
