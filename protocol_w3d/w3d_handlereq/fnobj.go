// Code generated by "genprotocol.exe -ver=69ee1a9a014411856ae5de810035b3c3478b3ee188bf040ff46c3ec85f20aba5 -basedir=protocol_w3d -prefix=w3d -statstype=int"

package w3d_handlereq

import (
	"fmt"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_error"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

// obj base demux fn map
var DemuxReq2ObjAPIFnMap = [...]func(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error){
	w3d_idcmd.Invalid:   Req2ObjAPI_Invalid,   // Invalid not used, make empty packet error
	w3d_idcmd.Login:     Req2ObjAPI_Login,     // Login make session with nickname and enter stage
	w3d_idcmd.Heartbeat: Req2ObjAPI_Heartbeat, // Heartbeat prevent connection timeout
	w3d_idcmd.Chat:      Req2ObjAPI_Chat,      // Chat chat to stage
	w3d_idcmd.Act:       Req2ObjAPI_Act,       // Act send user action
	w3d_idcmd.StatsInfo: Req2ObjAPI_StatsInfo, // StatsInfo game stats info

} // DemuxReq2ObjAPIFnMap

// Invalid not used, make empty packet error
func Req2ObjAPI_Invalid(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error) {
	req, ok := robj.(*w3d_obj.ReqInvalid_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqInvalid(me, hd, req)
	return rhd, rsp, err
}

// Invalid not used, make empty packet error
func objAPIFn_ReqInvalid(
	me interface{}, hd w3d_packet.Header, robj *w3d_obj.ReqInvalid_data) (
	w3d_packet.Header, *w3d_obj.RspInvalid_data, error) {
	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspInvalid_data{}
	return sendHeader, sendBody, nil
}

// Login make session with nickname and enter stage
func Req2ObjAPI_Login(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error) {
	req, ok := robj.(*w3d_obj.ReqLogin_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqLogin(me, hd, req)
	return rhd, rsp, err
}

// Login make session with nickname and enter stage
func objAPIFn_ReqLogin(
	me interface{}, hd w3d_packet.Header, robj *w3d_obj.ReqLogin_data) (
	w3d_packet.Header, *w3d_obj.RspLogin_data, error) {
	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspLogin_data{}
	return sendHeader, sendBody, nil
}

// Heartbeat prevent connection timeout
func Req2ObjAPI_Heartbeat(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error) {
	req, ok := robj.(*w3d_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqHeartbeat(me, hd, req)
	return rhd, rsp, err
}

// Heartbeat prevent connection timeout
func objAPIFn_ReqHeartbeat(
	me interface{}, hd w3d_packet.Header, robj *w3d_obj.ReqHeartbeat_data) (
	w3d_packet.Header, *w3d_obj.RspHeartbeat_data, error) {
	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspHeartbeat_data{}
	return sendHeader, sendBody, nil
}

// Chat chat to stage
func Req2ObjAPI_Chat(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error) {
	req, ok := robj.(*w3d_obj.ReqChat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqChat(me, hd, req)
	return rhd, rsp, err
}

// Chat chat to stage
func objAPIFn_ReqChat(
	me interface{}, hd w3d_packet.Header, robj *w3d_obj.ReqChat_data) (
	w3d_packet.Header, *w3d_obj.RspChat_data, error) {
	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspChat_data{}
	return sendHeader, sendBody, nil
}

// Act send user action
func Req2ObjAPI_Act(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error) {
	req, ok := robj.(*w3d_obj.ReqAct_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqAct(me, hd, req)
	return rhd, rsp, err
}

// Act send user action
func objAPIFn_ReqAct(
	me interface{}, hd w3d_packet.Header, robj *w3d_obj.ReqAct_data) (
	w3d_packet.Header, *w3d_obj.RspAct_data, error) {
	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspAct_data{}
	return sendHeader, sendBody, nil
}

// StatsInfo game stats info
func Req2ObjAPI_StatsInfo(
	me interface{}, hd w3d_packet.Header, robj interface{}) (
	w3d_packet.Header, interface{}, error) {
	req, ok := robj.(*w3d_obj.ReqStatsInfo_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqStatsInfo(me, hd, req)
	return rhd, rsp, err
}

// StatsInfo game stats info
func objAPIFn_ReqStatsInfo(
	me interface{}, hd w3d_packet.Header, robj *w3d_obj.ReqStatsInfo_data) (
	w3d_packet.Header, *w3d_obj.RspStatsInfo_data, error) {
	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspStatsInfo_data{}
	return sendHeader, sendBody, nil
}
