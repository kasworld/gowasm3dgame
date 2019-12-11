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

package stage

import (
	"math/rand"
	"time"

	"github.com/kasworld/gowasm3dgame/enums/acttype"
	"github.com/kasworld/gowasm3dgame/enums/acttype_stats"
	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/game/gameconst"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/uuidstr"
)

type Team struct {
	rnd *rand.Rand      `prettystring:"hide"`
	log *w3dlog.LogBase `prettystring:"hide"`

	ActStats acttype_stats.ActTypeStat
	Color24  htmlcolors.Color24

	IsAlive     bool
	RespawnTick int64

	UUID     string
	Ball     *GameObj
	HomeMark *GameObj
	Objs     []*GameObj

	ActPoint int
	Score    int
}

func NewTeam(l *w3dlog.LogBase, TeamType, color htmlcolors.Color24) *Team {
	nowtick := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	bt := &Team{
		rnd:     rnd,
		log:     l,
		IsAlive: true,
		Color24: color,
		Ball: &GameObj{
			GOType:       gameobjtype.Ball,
			UUID:         uuidstr.New(),
			BirthTick:    nowtick,
			LastMoveTick: nowtick,
			PosVt: vector3f.Vector3f{
				rnd.Float64() * gameconst.StageSize,
				rnd.Float64() * gameconst.StageSize,
				rnd.Float64() * gameconst.StageSize,
			},
		},
		Objs: make([]*GameObj, 0),
	}
	maxv := gameobjtype.Attrib[gameobjtype.Ball].SpeedLimit

	bt.Ball.MvVt = vector3f.Vector3f{
		bt.rnd.Float64() * maxv,
		bt.rnd.Float64() * maxv,
		bt.rnd.Float64() * maxv,
	}.NormalizedTo(maxv)
	return bt
}

func (bt *Team) RespawnBall(now int64) {
	bt.IsAlive = true
	bt.Ball.toDelete = false
	bt.Ball.PosVt = vector3f.Vector3f{
		bt.rnd.Float64() * gameconst.StageSize,
		bt.rnd.Float64() * gameconst.StageSize,
		bt.rnd.Float64() * gameconst.StageSize,
	}
	bt.Ball.MvVt = vector3f.Vector3f{
		0, 0, 0,
	}
	bt.Ball.LastMoveTick = now
	// bt.Ball.BirthTick = now
}

func (bt *Team) ToPacket() *w3d_obj.Team {
	rtn := &w3d_obj.Team{
		Ball: bt.Ball.ToPacket(),
		Objs: make([]*w3d_obj.GameObj, 0),
	}
	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		rtn.Objs = append(rtn.Objs, v.ToPacket())
	}
	return rtn
}

func (bt *Team) Count(ot gameobjtype.GameObjType) int {
	rtn := 0
	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		if v.GOType == ot {
			rtn++
		}
	}
	return rtn
}

func (bt *Team) addGObj(o *GameObj) {
	for i, v := range bt.Objs {
		if v.toDelete {
			bt.Objs[i] = o
			return
		}
	}
	bt.Objs = append(bt.Objs, o)
}

// 0(outer max) ~ GameConst.APIncFrame( 0,0,0)
func (t *Team) CalcAP(stageDiag float64) int {
	homepos := t.HomeMark.PosVt
	lenToHomepos := t.Ball.PosVt.LenTo(homepos)
	rtn := int((stageDiag - lenToHomepos) / stageDiag * float64(gameconst.APIncPerFrame))
	//log.Printf("ap:%v", rtn)
	return rtn
}

func (bt *Team) ApplyAct(actObj *w3d_obj.Act) {
	bt.ActStats.Inc(actObj.Act)
	switch actObj.Act {
	default:
		bt.log.Fatal("unknown act %+v %v", actObj, bt)
	case acttype.Nothing:
	case acttype.Bullet:
		bt.AddBullet(actObj.Vt)
	case acttype.SuperBullet:
		bt.AddSuperBullet(actObj.Vt)
	case acttype.HommingBullet:
		bt.AddHommingBullet(actObj.Vt, actObj.DstObjID)
	case acttype.Accel:
		bt.Ball.AccVt = actObj.Vt
	}
}

func (bt *Team) AddShield(vt vector3f.Vector3f) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gameobjtype.Shield,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		MvVt:         vt,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddBullet(vt vector3f.Vector3f) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gameobjtype.Bullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		MvVt:         vt,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddSuperBullet(vt vector3f.Vector3f) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gameobjtype.SuperBullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		MvVt:         vt,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddHommingBullet(vt vector3f.Vector3f, dstid string) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gameobjtype.HommingBullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		MvVt:         vt,
		DstUUID:      dstid,
	}
	bt.addGObj(o)
	return o
}