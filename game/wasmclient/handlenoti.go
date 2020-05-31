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
	"fmt"
	"time"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

var DemuxNoti2ObjFnMap = [...]func(me interface{}, hd w3d_packet.Header, body interface{}) error{
	w3d_idnoti.Invalid:   objRecvNotiFn_Invalid,
	w3d_idnoti.StageInfo: objRecvNotiFn_StageInfo,
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
	app.updateLeftInfo()
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
