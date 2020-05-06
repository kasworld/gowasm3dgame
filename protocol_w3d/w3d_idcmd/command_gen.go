// Code generated by "genprotocol -ver=fd815e2fbc449528b4fb5d55480af0a03b4bfaf074ff2c5570d2e5a3ce03896b -basedir=. -prefix=w3d -statstype=int"

package w3d_idcmd

import "fmt"

type CommandID uint16 // use in packet header, DO NOT CHANGE
const (
	Invalid    CommandID = iota //
	Login                       //
	Heartbeat                   //
	Chat                        //
	EnterStage                  //

	CommandID_Count int = iota
)

var _CommandID2string = [CommandID_Count][2]string{
	Invalid:    {"Invalid", ""},
	Login:      {"Login", ""},
	Heartbeat:  {"Heartbeat", ""},
	Chat:       {"Chat", ""},
	EnterStage: {"EnterStage", ""},
}

func (e CommandID) String() string {
	if e >= 0 && e < CommandID(CommandID_Count) {
		return _CommandID2string[e][0]
	}
	return fmt.Sprintf("CommandID%d", uint16(e))
}

func (e CommandID) CommentString() string {
	if e >= 0 && e < CommandID(CommandID_Count) {
		return _CommandID2string[e][1]
	}
	return ""
}

var _string2CommandID = map[string]CommandID{
	"Invalid":    Invalid,
	"Login":      Login,
	"Heartbeat":  Heartbeat,
	"Chat":       Chat,
	"EnterStage": EnterStage,
}

func String2CommandID(s string) (CommandID, bool) {
	v, b := _string2CommandID[s]
	return v, b
}
