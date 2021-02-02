/*
	INTEGER FUNCTIONS
*/

package arithmetic

import (
	"fmt"
	"regexp"
	"strconv"
)

// IsInteger Value is an integer?
// value Value to check.
func IsInteger(value string) bool {
	state, _ := regexp.MatchString("^(-|)\\s*[0-9]+$", value)
	return state
}

// ToSByte String to 8bit integer.
// value Value to parse.
func ToSByte(value string) (int8, error) {
	result, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(result), err
}

// ToShort String to 16bit integer.
// value Value to parse.
func ToShort(value string) (int16, error) {
	result, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(result), err
}

// ToInt String to 32bit integer.
// value Value to parse.
func ToInt(value string) (int32, error) {
	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(result), err
}

// ToLong String to 64bit integer.
// value Value to parse.
func ToLong(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(result), err
}

// ToByte String to 8bit unsigned integer.
// value Value to parse.
func ToByte(value string) (uint8, error) {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(result), err
}

// ToUShort String to 16bit unsigned integer.
// value Value to parse.
func ToUShort(value string) (uint16, error) {
	result, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(result), err
}

// ToUInt String to 32bit integer.
// value Value to parse.
func ToUInt(value string) (uint32, error) {
	result, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(result), err
}

// ToULong String to 64bit unsigned integer.
// value Value to parse.
func ToULong(value string) (uint64, error) {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), err
}

// IntToString Integer to string.
// value Value to parse.
func IntToString(value interface{}) string {
	return fmt.Sprintf("%d", value)
}
