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

	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/conndata"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connbytemanager"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_serveconnbyte"
)

func (svr *Server) api_me2conndata(me interface{}) (*conndata.ConnData, error) {
	conn, ok := me.(*w3d_serveconnbyte.ServeConnByte)
	if !ok {
		return nil, fmt.Errorf("Packet type miss match %v", me)
	}
	connData, ok := conn.GetConnData().(*conndata.ConnData)
	if !ok {
		return nil, fmt.Errorf("Packet type miss match %v", conn.GetConnData())
	}
	return connData, nil
}

type stageApiI interface {
	MakeType2Radius() [gameobjtype.GameObjType_Count]float64
	GetConnManager() *w3d_connbytemanager.Manager
	ToPacket_StatsInfo() *w3d_obj.RspStatsInfo_data
}
