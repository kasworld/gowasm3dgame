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
	"github.com/kasworld/gowasm2dgame/enums/acttype"
	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/htmlcolors"
)

type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

type ReqMakeTeam_data struct {
	Dummy uint8
}
type RspMakeTeam_data struct {
	Dummy uint8
}

type ReqAct_data struct {
	Dummy uint8
}
type RspAct_data struct {
	Dummy uint8
}

type ReqHeartbeat_data struct {
	Tick int64
}
type RspHeartbeat_data struct {
	Tick int64
}

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiWorldInfo_data struct {
	Dummy uint8
}

type NotiNearInfo_data struct {
	Dummy uint8
}

type NotiStatsInfo_data struct {
	Dummy uint8
}

/////////////////////////////

type World struct {
	ID           string
	BorderBounce vector3f.HyperRect
	BorderOctree vector3f.HyperRect
	Teams        []*Team
}

type Team struct {
	ID      string
	Color24 htmlcolors.Color24
	Objs    []*GameObj
}

type GameObj struct {
	ObjType gameobjtype.GameObjType
	ID      string
	PosVt   vector3f.Vector3f
}

func (o *GameObj) IsCollision(dst *GameObj) bool {
	return gameobjtype.CollisionTo(
		o.ObjType, dst.ObjType,
		dst.PosVt.Sqd(o.PosVt),
	)
}

type Act struct {
	Act   acttype.ActType
	Vt    vector3f.Vector3f
	Count int
	DstID string
}
