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

package w3d_obj

import (
	"github.com/kasworld/gowasm3dgame/enums/acttype"
	"github.com/kasworld/gowasm3dgame/enums/acttype_stats"
	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/htmlcolors"
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
	Tick int64
}

//////////////////////////////////////////////////////////////////////////////

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiStageInfo_data struct {
	Tick  int64
	Teams []*Team
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

	Color24  htmlcolors.Color24
	ActStats acttype_stats.ActTypeStat
}

type Act struct {
	Act      acttype.ActType
	Vt       [3]float64
	Count    int
	DstObjID string
}

type Team struct {
	ID      string
	Color24 htmlcolors.Color24
	Objs    []*GameObj // 0 : homemark, 1: ball
}

type GameObj struct {
	GOType gameobjtype.GameObjType
	UUID   string
	PosVt  [3]float64
	RotVt  [3]float64
}
