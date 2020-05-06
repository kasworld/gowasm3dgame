// Code generated by "genprotocol -ver=fd815e2fbc449528b4fb5d55480af0a03b4bfaf074ff2c5570d2e5a3ce03896b -basedir=. -prefix=w3d -statstype=int"

package w3d_error

import "fmt"

type ErrorCode uint16 // use in packet header, DO NOT CHANGE
const (
	None             ErrorCode = iota //
	ActionProhibited                  //
	ObjectNotFound                    //
	ActionChaned                      //
	ActionCanceled                    //

	ErrorCode_Count int = iota
)

var _ErrorCode2string = [ErrorCode_Count][2]string{
	None:             {"None", ""},
	ActionProhibited: {"ActionProhibited", ""},
	ObjectNotFound:   {"ObjectNotFound", ""},
	ActionChaned:     {"ActionChaned", ""},
	ActionCanceled:   {"ActionCanceled", ""},
}

func (e ErrorCode) String() string {
	if e >= 0 && e < ErrorCode(ErrorCode_Count) {
		return _ErrorCode2string[e][0]
	}
	return fmt.Sprintf("ErrorCode%d", uint16(e))
}
func (e ErrorCode) CommentString() string {
	if e >= 0 && e < ErrorCode(ErrorCode_Count) {
		return _ErrorCode2string[e][1]
	}
	return ""
}

// implement error interface
func (e ErrorCode) Error() string {
	return "w3d_error." + e.String()
}

var _string2ErrorCode = map[string]ErrorCode{
	"None":             None,
	"ActionProhibited": ActionProhibited,
	"ObjectNotFound":   ObjectNotFound,
	"ActionChaned":     ActionChaned,
	"ActionCanceled":   ActionCanceled,
}

func String2ErrorCode(s string) (ErrorCode, bool) {
	v, b := _string2ErrorCode[s]
	return v, b
}
