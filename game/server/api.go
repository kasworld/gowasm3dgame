// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	"runtime"
	"time"

	"github.com/kasworld/gowasm3dgame/config/authdata"
	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/lib/conndata"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_authorize"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_error"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_gob"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_serveconnbyte"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_version"
	"github.com/kasworld/version"
)

func (svr *Server) setFnMap() {
	svr.DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error){

		w3d_idcmd.Invalid:   svr.bytesAPIFn_ReqInvalid,   // Invalid not used, make empty packet error
		w3d_idcmd.Login:     svr.bytesAPIFn_ReqLogin,     // Login make session with nickname and enter stage
		w3d_idcmd.Heartbeat: svr.bytesAPIFn_ReqHeartbeat, // Heartbeat prevent connection timeout
		w3d_idcmd.Chat:      svr.bytesAPIFn_ReqChat,      // Chat chat to stage
		w3d_idcmd.Act:       svr.bytesAPIFn_ReqAct,       // Act send user action

	}
}

// Invalid not used, make empty packet error
func (svr *Server) bytesAPIFn_ReqInvalid(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	sendHeader := w3d_packet.Header{}
	return sendHeader, nil, fmt.Errorf("invalid packet")
}

// Login make session with nickname and enter stage
func (svr *Server) bytesAPIFn_ReqLogin(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	robj, err := w3d_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*w3d_obj.ReqLogin_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}

	c2sc, ok := me.(*w3d_serveconnbyte.ServeConnByte)
	if !ok {
		panic(fmt.Sprintf("invalid me not w3d_serveconnbyte.ServeConnByte %#v", me))
	}

	if err := authdata.UpdateByAuthKey(c2sc.GetAuthorCmdList(), recvBody.AuthKey); err != nil {
		return sendHeader, nil, err
	}
	connData := c2sc.GetConnData().(*conndata.ConnData)

	ss := svr.sessionManager.UpdateOrNew(
		recvBody.SessionKey,
		connData.RemoteAddr,
		recvBody.NickName)

	if oldc2sc := svr.connManager.Get(ss.ConnUUID); oldc2sc != nil {
		oldc2sc.Disconnect()
		// wait
		trycount := 10
		for svr.connManager.Get(ss.ConnUUID) != nil && trycount > 0 {
			runtime.Gosched()
			time.Sleep(time.Millisecond * 100)
			trycount--
		}
	}
	if svr.connManager.Get(ss.ConnUUID) != nil {
		svr.log.Fatal("old connection online %v", ss)
		return sendHeader, nil, err
	}

	ss.ConnUUID = connData.UUID
	connData.Session = ss

	// select stage to play
	stg := svr.stageManager.GetAny()
	ss.StageID = stg.GetUUID()
	stg.GetConnManager().Add(connData.UUID, c2sc)

	// user login?

	if err != nil {
		return sendHeader, nil, err
	} else {
		sendBody := &w3d_obj.RspLogin_data{
			Version:         version.GetVersion(),
			ProtocolVersion: w3d_version.ProtocolVersion,
			DataVersion:     gameconst.DataVersion,
			SessionKey:      recvBody.SessionKey,
			NickName:        recvBody.NickName,
			CmdList:         *w3d_authorize.NewAllSet(),
		}
		return sendHeader, sendBody, nil
	}
}

// Chat chat to stage
func (svr *Server) bytesAPIFn_ReqChat(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	robj, err := w3d_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*w3d_obj.ReqChat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	conn, ok := me.(*w3d_serveconnbyte.ServeConnByte)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", me)
	}
	connData, ok := conn.GetConnData().(*conndata.ConnData)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", conn.GetConnData())
	}

	stg := svr.stageManager.GetByUUID(connData.Session.StageID)
	if stg == nil {
		svr.log.Fatal("no stage to chat %v", connData)
		return hd, nil, fmt.Errorf("stage not ready %v", connData)
	}
	connList := stg.GetConnManager().GetList()
	noti := &w3d_obj.NotiStageChat_data{
		SenderNick: connData.Session.NickName,
		Chat:       recvBody.Chat,
	}
	for _, v := range connList {
		v.SendNotiPacket(w3d_idnoti.StageChat,
			noti,
		)
	}

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspChat_data{}
	return sendHeader, sendBody, nil
}

// Act send user action
func (svr *Server) bytesAPIFn_ReqAct(
	me interface{}, hd w3d_packet.Header, rbody []byte) (
	w3d_packet.Header, interface{}, error) {
	robj, err := w3d_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*w3d_obj.ReqAct_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := w3d_packet.Header{
		ErrorCode: w3d_error.None,
	}
	sendBody := &w3d_obj.RspAct_data{}
	return sendHeader, sendBody, nil
}

// Heartbeat prevent connection timeout
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
