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

package wasmclient

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"syscall/js"
	"time"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connwasm"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_gob"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_pid2rspfn"
)

func getConnURL() string {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}
	u.Path = "ws"
	u.Scheme = "ws"
	return u.String()
}

func (app *WasmClient) NetInit(ctx context.Context) error {
	app.wsConn = w3d_connwasm.New(
		getConnURL(),
		w3d_gob.MarshalBodyFn,
		app.handleRecvPacket,
		app.handleSentPacket)

	fmt.Println(getConnURL())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := app.wsConn.Connect(ctx, &wg)
		if err != nil {
			fmt.Printf("wsConn.Connect err %v\n", err)
			app.DoClose()
		}
	}()
	wg.Wait()

	return nil
}

func (app *WasmClient) Cleanup() {
	app.wsConn.SendRecvStop()
}

func (app *WasmClient) handleSentPacket(header w3d_packet.Header) error {
	return nil
}

func (app *WasmClient) handleRecvPacket(header w3d_packet.Header, body []byte) error {
	robj, err := w3d_gob.UnmarshalPacket(header, body)
	if err != nil {
		return err
	}
	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, robj)
	case w3d_packet.Response:
		if err := app.pid2recv.HandleRsp(header, robj); err != nil {
			return err
		}
	case w3d_packet.Notification:
		fn := DemuxNoti2ObjFnMap[header.Cmd]
		if err := fn(app, header, robj); err != nil {
			return err
		}
	}
	return nil
}

func (app *WasmClient) ReqWithRspFn(cmd w3d_idcmd.CommandID, body interface{},
	fn w3d_pid2rspfn.HandleRspFn) error {

	pid := app.pid2recv.NewPID(fn)
	spk := w3d_packet.Packet{
		Header: w3d_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: w3d_packet.Request,
		},
		Body: body,
	}
	if err := app.wsConn.EnqueueSendPacket(spk); err != nil {
		app.wsConn.SendRecvStop()
		return fmt.Errorf("Send fail %s %v:%v %v", app, cmd, pid, err)
	}
	return nil
}

func (app *WasmClient) reqHeartbeat() error {
	return app.ReqWithRspFn(
		w3d_idcmd.Heartbeat,
		&w3d_obj.ReqHeartbeat_data{
			Tick: time.Now().UnixNano(),
		},
		func(hd w3d_packet.Header, rsp interface{}) error {
			rpk := rsp.(*w3d_obj.RspHeartbeat_data)
			pingDur := time.Now().UnixNano() - rpk.Tick
			app.PingDur = (app.PingDur + pingDur) / 2
			return nil
		},
	)
}

var DemuxNoti2ObjFnMap = [...]func(me interface{}, hd w3d_packet.Header, body interface{}) error{
	w3d_idnoti.Invalid:   objRecvNotiFn_Invalid,
	w3d_idnoti.StageInfo: objRecvNotiFn_StageInfo,
	w3d_idnoti.StatsInfo: objRecvNotiFn_StatsInfo,
	w3d_idnoti.StageChat: objRecvNotiFn_StageChat,
}

func objRecvNotiFn_Invalid(me interface{}, hd w3d_packet.Header, body interface{}) error {
	robj, ok := body.(*w3d_obj.NotiInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

func objRecvNotiFn_StageInfo(me interface{}, hd w3d_packet.Header, body interface{}) error {
	robj, ok := body.(*w3d_obj.NotiStageInfo_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app, ok := me.(*WasmClient)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app.vp.processRecvStageInfo(robj)

	app.ServerClientTictDiff = robj.Tick - time.Now().UnixNano()
	return nil
}

func objRecvNotiFn_StatsInfo(me interface{}, hd w3d_packet.Header, body interface{}) error {
	robj, ok := body.(*w3d_obj.NotiStatsInfo_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app, ok := me.(*WasmClient)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app.statsInfo = robj
	return nil
}

func objRecvNotiFn_StageChat(me interface{}, hd w3d_packet.Header, body interface{}) error {
	robj, ok := body.(*w3d_obj.NotiStageChat_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app, ok := me.(*WasmClient)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app.systemMessage.Appendf("%v : %v", robj.SenderNick, robj.Chat)
	return nil
}
