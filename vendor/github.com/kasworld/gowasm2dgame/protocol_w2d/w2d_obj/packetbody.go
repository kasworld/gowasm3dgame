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

package w2d_obj

import (
	"github.com/kasworld/gowasm2dgame/enums/acttype"
	"github.com/kasworld/gowasm2dgame/enums/acttype_stats"
	"github.com/kasworld/gowasm2dgame/enums/teamtype"
)

type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

type ReqEnterStage_data struct {
	StageUUID string // may be not same to req stage
	NickToUse string
}
type RspEnterStage_data struct {
	StageUUID string // may be not same to req stage
	NickToUse string // may be not same to req nick
}

type ReqChatToStage_data struct {
	Chat string
}
type RspChatToStage_data struct {
	Dummy uint8
}

type ReqHeartbeat_data struct {
	Tick int64
}
type RspHeartbeat_data struct {
	Tick int64 // same req tick , to calc round trip time
}

//////////////////////////////////////////////////////////////////////////////

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiStageInfo_data struct {
	Tick       int64
	Background *Background
	Teams      []*Team
	Effects    []*Effect
	Clouds     []*Cloud
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
	UUID     string
	Alive    bool
	TeamType teamtype.TeamType
	ActStats acttype_stats.ActTypeStat
}

type Act struct {
	Act acttype.ActType

	// accel, fire bullet
	Angle  float64 // degree
	AngleV float64 // pixel /sec

	// homming
	DstObjID string
}
