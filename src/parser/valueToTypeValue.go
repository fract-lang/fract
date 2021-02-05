/*
	ValueToTypeValue Function.
*/

package parser

import (
	"../fract/arithmetic"
	"../grammar"
)

// ValueToTypeValue Value to type value by limit checks.
// _type Type to parse.
// value Value to parse.
func ValueToTypeValue(_type string, value string) (string, string) {
	switch _type {
	case grammar.DtInt8:
		result, err := arithmetic.ToInt8(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtInt16:
		result, err := arithmetic.ToInt16(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtInt32:
		result, err := arithmetic.ToInt32(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtInt64:
		result, err := arithmetic.ToInt64(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtUInt8:
		result, err := arithmetic.ToUInt8(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtUInt16:
		result, err := arithmetic.ToUInt16(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtUInt32:
		result, err := arithmetic.ToUInt32(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtUInt64:
		result, err := arithmetic.ToUInt64(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.IntToString(result), ""
	case grammar.DtFloat32:
		result, err := arithmetic.ToFloat32(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.FloatToString(result), ""
	case grammar.DtFloat64:
		result, err := arithmetic.ToFloat64(value)
		if err != nil {
			return "", "Value out of range!"
		}
		return arithmetic.FloatToString(result), ""
	case grammar.DtBoolean:
		if value != grammar.KwTrue && value != grammar.KwFalse &&
			value != "0" && value != "1" {
			return "", "Boolean value is not valid!"
		}
		result := grammar.KwFalse
		if value == "1" || value == grammar.KwTrue {
			result = grammar.KwTrue
		}
		return result, ""
	default:
		return "", "Data type is not found!"
	}
}
