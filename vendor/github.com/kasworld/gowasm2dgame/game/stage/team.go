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

package stage

import (
	"math"
	"math/rand"
	"time"

	"github.com/kasworld/gowasm2dgame/enums/acttype"
	"github.com/kasworld/gowasm2dgame/enums/acttype_stats"
	"github.com/kasworld/gowasm2dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm2dgame/enums/teamtype"
	"github.com/kasworld/gowasm2dgame/game/gameconst"
	"github.com/kasworld/gowasm2dgame/lib/vector2f"
	"github.com/kasworld/gowasm2dgame/lib/w2dlog"
	"github.com/kasworld/gowasm2dgame/protocol_w2d/w2d_obj"
	"github.com/kasworld/uuidstr"
)

type Team struct {
	rnd *rand.Rand      `prettystring:"hide"`
	log *w2dlog.LogBase `prettystring:"hide"`

	ActStats acttype_stats.ActTypeStat
	UUID     string

	TeamType    teamtype.TeamType
	IsAlive     bool
	RespawnTick int64

	Ball *GameObj // ball is special
	Objs []*GameObj
}

func NewTeam(l *w2dlog.LogBase, TeamType teamtype.TeamType) *Team {
	nowtick := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	bt := &Team{
		rnd:      rnd,
		log:      l,
		IsAlive:  true,
		TeamType: TeamType,
		UUID:     uuidstr.New(),
		Ball: &GameObj{
			teamType:     TeamType,
			GOType:       gameobjtype.Ball,
			UUID:         uuidstr.New(),
			BirthTick:    nowtick,
			LastMoveTick: nowtick,
			PosVt: vector2f.Vector2f{
				rnd.Float64() * gameconst.StageW,
				rnd.Float64() * gameconst.StageH,
			},
		},
		Objs: make([]*GameObj, 0),
	}
	maxv := gameobjtype.Attrib[gameobjtype.Ball].SpeedLimit

	vt := vector2f.NewVectorLenAngle(
		bt.rnd.Float64()*maxv,
		bt.rnd.Float64()*360,
	)
	bt.Ball.SetDxy(vt)
	return bt
}

func (bt *Team) RespawnBall(now int64) {
	bt.IsAlive = true
	bt.Ball.toDelete = false
	bt.Ball.PosVt = vector2f.Vector2f{
		bt.rnd.Float64() * gameconst.StageW,
		bt.rnd.Float64() * gameconst.StageH,
	}
	bt.Ball.VelVt = vector2f.Vector2f{
		0, 0,
	}
	bt.Ball.LastMoveTick = now
	// bt.Ball.BirthTick = now
}

func (bt *Team) ToPacket() *w2d_obj.Team {
	rtn := &w2d_obj.Team{
		TeamType: bt.TeamType,
		Ball:     bt.Ball.ToPacket(),
		Objs:     make([]*w2d_obj.GameObj, 0),
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

func (bt *Team) GetRemainAct(now int64, act acttype.ActType) float64 {
	durSec := float64(now-bt.Ball.BirthTick) / float64(time.Second)
	actedCount := float64(bt.ActStats[act])
	totalCanAct := durSec * acttype.Attrib[act].PerSec
	remainAct := totalCanAct - actedCount
	return remainAct
}

func (bt *Team) ApplyAct(actObj *w2d_obj.Act) {
	bt.ActStats.Inc(actObj.Act)
	switch actObj.Act {
	default:
		bt.log.Fatal("unknown act %+v %v", actObj, bt)
	case acttype.Nothing:
	case acttype.Shield:
		bt.AddShield(actObj.Angle, actObj.AngleV)
	case acttype.SuperShield:
		bt.AddSuperShield(actObj.Angle, actObj.AngleV)
	case acttype.HommingShield:
		bt.AddHommingShield(actObj.Angle, actObj.AngleV)
	case acttype.Bullet:
		bt.AddBullet(actObj.Angle, actObj.AngleV)
	case acttype.SuperBullet:
		bt.AddSuperBullet(actObj.Angle, actObj.AngleV)
	case acttype.HommingBullet:
		bt.AddHommingBullet(actObj.Angle, actObj.AngleV, actObj.DstObjID)
	case acttype.Accel:
		vt := vector2f.NewVectorLenAngle(actObj.AngleV, actObj.Angle)
		bt.Ball.AddDxy(vt)
	}
}

func (bt *Team) AddShield(angle, anglev float64) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		teamType:     bt.TeamType,
		GOType:       gameobjtype.Shield,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		Angle:        angle,
		AngleV:       anglev,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddSuperShield(angle, anglev float64) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		teamType:     bt.TeamType,
		GOType:       gameobjtype.SuperShield,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		Angle:        angle,
		AngleV:       anglev,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddBullet(angle, anglev float64) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		teamType:     bt.TeamType,
		GOType:       gameobjtype.Bullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		VelVt:         vector2f.NewVectorLenAngle(anglev, angle),
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddSuperBullet(angle, anglev float64) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		teamType:     bt.TeamType,
		GOType:       gameobjtype.SuperBullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        bt.Ball.PosVt,
		VelVt:         vector2f.NewVectorLenAngle(anglev, angle),
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddHommingShield(angle, anglev float64) *GameObj {
	nowtick := time.Now().UnixNano()
	mvvt := vector2f.NewVectorLenAngle(anglev, angle)
	o := &GameObj{
		teamType:     bt.TeamType,
		GOType:       gameobjtype.HommingShield,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		Angle:        angle,
		AngleV:       anglev,
		PosVt:        bt.Ball.PosVt.Add(mvvt),
		VelVt:         mvvt,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) AddHommingBullet(angle, anglev float64, dstid string) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		teamType:     bt.TeamType,
		GOType:       gameobjtype.HommingBullet,
		UUID:         uuidstr.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		Angle:        angle,
		AngleV:       anglev,
		PosVt:        bt.Ball.PosVt,
		VelVt:         vector2f.NewVectorLenAngle(anglev, angle),
		DstUUID:      dstid,
	}
	bt.addGObj(o)
	return o
}

func (bt *Team) CalcAimAngleAndV(
	bullet gameobjtype.GameObjType, dsto *GameObj) (float64, float64) {
	s1 := gameobjtype.Attrib[bullet].SpeedLimit
	vt := dsto.PosVt.Sub(bt.Ball.PosVt)
	s2 := dsto.VelVt.Abs()
	if s2 == 0 {
		return vt.Phase(), s1
	}
	a2 := dsto.VelVt.Phase() - vt.Phase()
	a1 := math.Asin(s2 / s1 * math.Sin(a2))

	return vt.AddAngle(a1).Phase(), s1
}
