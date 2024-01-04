package utils

// #include <stdio.h>
// #include <stdlib.h>
// #include "./x17/x17.h"
// #cgo CFLAGS: -std=gnu99 -Wall -I.
// #cgo LDFLAGS: -L./x17/ -lx17
import "C"

import (
	"encoding/hex"
	"fmt"
	"unsafe"
)

func X17_Sum256(input string) []byte {
	in, err := hex.DecodeString(input)
	if err != nil {
		fmt.Println("X17_Sum256 DecodeString error:", err)
		return nil
	}
	in1 := (*C.char)(unsafe.Pointer(&in[0]))

	output := make([]byte, 32)
	out := (*C.char)(unsafe.Pointer(&output[0]))

	C.x17_hash(unsafe.Pointer(out), unsafe.Pointer(in1))

	return output
}

func X17_Byte_Sum256(input []byte) []byte {

	in1 := (*C.char)(unsafe.Pointer(&input[0]))

	output := make([]byte, 32)
	out := (*C.char)(unsafe.Pointer(&output[0]))

	C.x17_hash(unsafe.Pointer(out), unsafe.Pointer(in1))

	return output
}
