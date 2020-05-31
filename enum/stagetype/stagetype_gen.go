// Code generated by "genenum -typename=StageType -packagename=stagetype -basedir=enum -vectortype=int"

package stagetype

import "fmt"

type StageType uint8

const (
	Invalid StageType = iota // unset or invalid
	Stage2D                  // 2d stage
	Stage3D                  // 3d stage

	StageType_Count int = iota
)

var _StageType2string = [StageType_Count][2]string{
	Invalid: {"Invalid", "unset or invalid"},
	Stage2D: {"Stage2D", "2d stage"},
	Stage3D: {"Stage3D", "3d stage"},
}

func (e StageType) String() string {
	if e >= 0 && e < StageType(StageType_Count) {
		return _StageType2string[e][0]
	}
	return fmt.Sprintf("StageType%d", uint8(e))
}

func (e StageType) CommentString() string {
	if e >= 0 && e < StageType(StageType_Count) {
		return _StageType2string[e][1]
	}
	return ""
}

var _string2StageType = map[string]StageType{
	"Invalid": Invalid,
	"Stage2D": Stage2D,
	"Stage3D": Stage3D,
}

func String2StageType(s string) (StageType, bool) {
	v, b := _string2StageType[s]
	return v, b
}
