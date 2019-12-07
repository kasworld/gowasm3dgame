// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_error"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_gob"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

///////////////////////////////////////////////////////////////

func (svr *Server) bytesAPIFn_ReqInvalid(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*w3d_obj.ReqInvalid_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspInvalid_data{}
	return sendHeader, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqMakeTeam(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*w3d_obj.ReqMakeTeam_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspMakeTeam_data{}
	return sendHeader, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqAct(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*w3d_obj.ReqAct_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspAct_data{}
	return sendHeader, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqHeartbeat(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	robj, err := w3d_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*w3d_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspHeartbeat_data{
		Tick: recvBody.Tick,
	}
	return sendHeader, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqNearInfo(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*w3d_obj.ReqNearInfo_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspNearInfo_data{}
	return sendHeader, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqWorldInfo(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*w3d_obj.ReqWorldInfo_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspWorldInfo_data{}
	return sendHeader, sendBody, nil
}
