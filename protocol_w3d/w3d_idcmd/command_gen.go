// Code generated by "genprotocol -ver=69ee1a9a014411856ae5de810035b3c3478b3ee188bf040ff46c3ec85f20aba5 -basedir=. -prefix=w3d -statstype=int"

package w3d_idcmd

import "fmt"

type CommandID uint16 // use in packet header, DO NOT CHANGE
const (
	Invalid   CommandID = iota // not used, make empty packet error
	Login                      // make session with nickname and enter stage
	Heartbeat                  // prevent connection timeout
	Chat                       // chat to stage
	Act                        // send user action
	StatsInfo                  // game stats info

	CommandID_Count int = iota
)

var _CommandID2string = [CommandID_Count][2]string{
	Invalid:   {"Invalid", "not used, make empty packet error"},
	Login:     {"Login", "make session with nickname and enter stage"},
	Heartbeat: {"Heartbeat", "prevent connection timeout"},
	Chat:      {"Chat", "chat to stage"},
	Act:       {"Act", "send user action"},
	StatsInfo: {"StatsInfo", "game stats info"},
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
	"Invalid":   Invalid,
	"Login":     Login,
	"Heartbeat": Heartbeat,
	"Chat":      Chat,
	"Act":       Act,
	"StatsInfo": StatsInfo,
}

func String2CommandID(s string) (CommandID, bool) {
	v, b := _string2CommandID[s]
	return v, b
}
