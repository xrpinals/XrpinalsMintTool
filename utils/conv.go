package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func Uint64Supply(supply interface{}) (uint64, error) {
	typeKind := reflect.TypeOf(supply).Kind()
	typeValue := reflect.ValueOf(supply)

	switch typeKind {
	case reflect.Float32, reflect.Float64:
		return uint64(typeValue.Float()), nil
	case reflect.String:
		u, err := strconv.ParseUint(typeValue.String(), 10, 64)
		if err != nil {
			return 0, err
		} else {
			return u, nil
		}
	}
	return 0, fmt.Errorf("not support type")
}
func ReverseBytesInPlace(input []byte) {
	length := len(input)

	for i := 0; i < length/2; i++ {
		input[i], input[length-1-i] = input[length-1-i], input[i]
	}
}
