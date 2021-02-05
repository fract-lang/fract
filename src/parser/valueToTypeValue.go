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
		result, _ := arithmetic.ToInt8(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtInt16:
		result, _ := arithmetic.ToInt16(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtInt32:
		result, _ := arithmetic.ToInt32(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtInt64:
		result, _ := arithmetic.ToInt64(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt8:
		result, _ := arithmetic.ToUInt8(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt16:
		result, _ := arithmetic.ToUInt16(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt32:
		result, _ := arithmetic.ToUInt32(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt64:
		result, _ := arithmetic.ToUInt64(value)
		rresult := arithmetic.IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtFloat32:
		result, _ := arithmetic.ToFloat32(value)
		rresult := arithmetic.FloatToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtFloat64:
		result, _ := arithmetic.ToFloat64(value)
		rresult := arithmetic.FloatToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	default:
		return "", "Data type is not found!"
	}
}
