/*
	ValueToTypeValue Function.
*/

package arithmetic

import (
	"../../grammar"
)

// ValueToTypeValue Value to type value by limit checks.
// _type Type to parse.
// value Value to parse.
func ValueToTypeValue(_type string, value string) (string, string) {
	switch _type {
	case grammar.DtInt8:
		result, _ := ToInt8(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtInt16:
		result, _ := ToInt16(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtInt32:
		result, _ := ToInt32(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtInt64:
		result, _ := ToInt64(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt8:
		result, _ := ToUInt8(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt16:
		result, _ := ToUInt16(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt32:
		result, _ := ToUInt32(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtUInt64:
		result, _ := ToUInt64(value)
		rresult := IntToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtFloat32:
		result, _ := ToFloat32(value)
		rresult := FloatToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	case grammar.DtFloat64:
		result, _ := ToFloat64(value)
		rresult := FloatToString(result)
		if rresult != value {
			return "", "The value data type was out of range!"
		}
		return rresult, ""
	default:
		return "", "Data type is not found!"
	}
}
