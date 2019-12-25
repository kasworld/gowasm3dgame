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
	"math"
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
	UUID     string

	IsAlive     bool
	RespawnTick int64

	Ball     *GameObj
	HomeMark *GameObj
	Objs     []*GameObj

	ActPoint float64
	Score    float64
}

func NewTeam(l *w3dlog.LogBase, color htmlcolors.Color24) *Team {
	nowtick := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	bt := &Team{
		rnd:     rnd,
		log:     l,
		UUID:    uuidstr.New(),
		IsAlive: true,
		Color24: color,
		Objs:    make([]*GameObj, 0),
	}

	maxv := gameobjtype.Attrib[gameobjtype.HomeMark].SpeedLimit
	bt.HomeMark = &GameObj{
		GOType:       gameobjtype.HomeMark,
		UUID:         uuidstr.New(),
		TeamUUID:     bt.UUID,
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.RandPosVt(),
		RotVelVt:     bt.RandRotVt(),
		VelVt: vector3f.Vector3f{
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
		}.NormalizedTo(maxv),
	}

	maxv = gameobjtype.Attrib[gameobjtype.Ball].SpeedLimit
	bt.Ball = &GameObj{
		GOType:       gameobjtype.Ball,
		UUID:         uuidstr.New(),
		TeamUUID:     bt.UUID,
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.RandPosVt(),
		RotVelVt:     bt.RandRotVt(),
		VelVt: vector3f.Vector3f{
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
		}.NormalizedTo(maxv),
	}
	return bt
}

func (bt *Team) RespawnBall(now int64) {
	bt.IsAlive = true
	bt.Ball.toDelete = false
	bt.Ball.PosVt = bt.RandPosVt()
	bt.Ball.VelVt = vector3f.Vector3f{
		0, 0, 0,
	}
	bt.Ball.LastMoveTick = now
}

func (bt *Team) RandRotVt() vector3f.Vector3f {
	return vector3f.Vector3f{
		bt.rnd.Float64() * math.Pi,
		bt.rnd.Float64() * math.Pi,
		bt.rnd.Float64() * math.Pi,
	}
}
func (bt *Team) RandPosVt() vector3f.Vector3f {
	return vector3f.Vector3f{
		bt.rnd.Float64() * gameconst.StageSize,
		bt.rnd.Float64() * gameconst.StageSize,
		bt.rnd.Float64() * gameconst.StageSize,
	}
}

func (bt *Team) ToPacket() *w3d_obj.Team {
	rtn := &w3d_obj.Team{
		ID:       bt.UUID,
		Color24:  bt.Color24,
		Ball:     bt.Ball.ToPacket(),
		HomeMark: bt.HomeMark.ToPacket(),
		Objs:     make([]*w3d_obj.GameObj, 0),
	}
	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		rtn.Objs = append(rtn.Objs, v.ToPacket())
	}
	return rtn
}

func (bt *Team) CountByGOType(ot gameobjtype.GameObjType) int {
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
func (t *Team) CalcAP(stageDiag float64) float64 {
	homepos := t.HomeMark.PosVt
	lenToHomepos := t.Ball.PosVt.LenTo(homepos)
	lenRate := (stageDiag - lenToHomepos) / stageDiag
	rtn := lenRate * gameconst.APIncPerFrame
	return rtn
}

func (bt *Team) CanAct(act acttype.ActType) bool {
	return bt.ActPoint >= acttype.Attrib[act].AP
}

func (bt *Team) CanHave(objt gameobjtype.GameObjType) bool {
	return bt.CountByGOType(objt) <= gameobjtype.Attrib[objt].MaxInTeam
}

func (bt *Team) ApplyAct(actObj *w3d_obj.Act) {
	bt.ActStats.Inc(actObj.Act)
	bt.ActPoint -= acttype.Attrib[actObj.Act].AP

	switch actObj.Act {
	default:
		bt.log.Fatal("unknown act %+v %v", actObj, bt)
	case acttype.Nothing:
	case acttype.Bullet:
		bt.AddBullet(actObj.Vt)
	case acttype.BurstBullet:
		for i := 0; i < 10; i++ {
			vt := bt.RandPosVt()
			bt.AddBurstBullet(vt.Sub(bt.Ball.PosVt))
		}
	case acttype.SuperBullet:
		bt.AddSuperBullet(actObj.Vt)
	case acttype.HommingBullet:
		bt.AddHommingBullet(actObj.Vt, actObj.DstObjID)
	case acttype.Accel:
		bt.Ball.VelVt = bt.Ball.VelVt.Add(actObj.Vt)
	case acttype.Shield:
		bt.AddShield(actObj.Vt)
	case acttype.HommingShield:
		bt.AddHommingShield(actObj.Vt, actObj.DstObjID)
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
		PosVt:        bt.Ball.PosVt,
		VelVt:        vt,
		RotVelVt:     bt.RandRotVt(),
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddHommingShield(vt vector3f.Vector3f, dstid string) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gameobjtype.HommingShield,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		VelVt:        vt,
		RotVelVt:     bt.RandRotVt(),
		DstUUID:      dstid,
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
		VelVt:        vt,
		RotVelVt:     bt.RandRotVt(),
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddBurstBullet(vt vector3f.Vector3f) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gameobjtype.BurstBullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		VelVt:        vt,
		RotVelVt:     bt.RandRotVt(),
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
		VelVt:        vt,
		RotVelVt:     bt.RandRotVt(),
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
		VelVt:        vt,
		RotVelVt:     bt.RandRotVt(),
		DstUUID:      dstid,
	}
	bt.addGObj(o)
	return o
}
