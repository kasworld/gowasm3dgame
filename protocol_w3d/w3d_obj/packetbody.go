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

package w3d_obj

import (
	"github.com/kasworld/gowasm3dgame/enum/acttype"
	"github.com/kasworld/gowasm3dgame/enum/acttype_vector"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
)

type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

type ReqLogin_data struct {
	SessionKey string
	NickName   string
	AuthKey    string
}
type RspLogin_data struct {
	Version         string
	ProtocolVersion string
	DataVersion     string

	SessionKey string
	NickName   string
	CmdList    [w3d_idcmd.CommandID_Count]bool
}

type ReqHeartbeat_data struct {
	Tick int64
}
type RspHeartbeat_data struct {
	Tick int64
}

type ReqChat_data struct {
	Chat string
}
type RspChat_data struct {
	Dummy uint8
}

//////////////////////////////////////////////////////////////////////////////

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiStageInfo_data struct {
	Tick         int64
	CameraPos    [3]float32
	CameraLookAt [3]float32
	ObjList      []*GameObj
}

type NotiStatsInfo_data struct {
	UUID  string
	Stats []TeamStat
}

type NotiStageChat_data struct {
	SenderNick string
	Chat       string
}

//////////////////////////////////////////////////////////////////////////////

type TeamStat struct {
	UUID  string
	Alive bool
	AP    int
	Score int
	Kill  int
	Death int

	Color24  uint32 // from htmlcolors.Color24
	ActStats acttype_vector.ActTypeVector
}

type Act struct {
	Act      acttype.ActType
	Vt       [3]float64
	Count    int
	DstObjID string
}

type GameObj struct {
	UUID    string
	GOType  gameobjtype.GameObjType
	Color24 uint32 // from htmlcolors.Color24
	PosVt   [3]float32
	RotVt   [3]float32
}
