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

// ToInt8 String to 8bit integer.
// value Value to parse.
func ToInt8(value string) (int8, error) {
	result, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(result), err
}

// ToInt16 String to 16bit integer.
// value Value to parse.
func ToInt16(value string) (int16, error) {
	result, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(result), err
}

// ToInt32 String to 32bit integer.
// value Value to parse.
func ToInt32(value string) (int32, error) {
	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(result), err
}

// ToInt64 String to 64bit integer.
// value Value to parse.
func ToInt64(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(result), err
}

// ToUInt8 String to 8bit unsigned integer.
// value Value to parse.
func ToUInt8(value string) (uint8, error) {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(result), err
}

// ToUInt16 String to 16bit unsigned integer.
// value Value to parse.
func ToUInt16(value string) (uint16, error) {
	result, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(result), err
}

// ToUInt32 String to 32bit integer.
// value Value to parse.
func ToUInt32(value string) (uint32, error) {
	result, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(result), err
}

// ToUInt64 String to 64bit unsigned integer.
// value Value to parse.
func ToUInt64(value string) (uint64, error) {
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
