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
	"time"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/acttype"
	"github.com/kasworld/gowasm3dgame/enum/acttype_vector"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/enum/stagetype"
	"github.com/kasworld/gowasm3dgame/lib/idu64str"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/htmlcolors"
)

var G_TeamID = idu64str.New("Team")
var G_GameObjID = idu64str.New("GameObj")

type Team struct {
	rnd       *g2rand.G2Rand  `prettystring:"hide"`
	log       *w3dlog.LogBase `prettystring:"hide"`
	StageType stagetype.StageType

	ActStats     acttype_vector.ActTypeVector
	Color24      htmlcolors.Color24
	UUID         string
	BorderBounce vector3f.Cube

	IsAlive     bool
	RespawnTick int64

	Ball     *GameObj
	HomeMark *GameObj
	Objs     []*GameObj

	ActPoint float64
	Score    float64
	Kill     int
	Death    int

	// 2d,3d mode fn
	RandBaseVtFn func() vector3f.Vector3f
}

func NewTeam(
	l *w3dlog.LogBase,
	color htmlcolors.Color24,
	BorderBounce vector3f.Cube,
	seed int64,
	StageType stagetype.StageType,
) *Team {
	bt := &Team{
		StageType:    StageType,
		rnd:          g2rand.NewWithSeed(seed),
		log:          l,
		UUID:         G_TeamID.New(),
		BorderBounce: BorderBounce,
		IsAlive:      true,
		Color24:      color,
		Objs:         make([]*GameObj, 0),
	}

	switch bt.StageType {
	default:
		bt.log.Fatal("invalid stagetype %v", bt.StageType)
	case stagetype.Stage2D:
		bt.RandBaseVtFn = bt.RandBaseVt2D
	case stagetype.Stage3D:
		bt.RandBaseVtFn = bt.RandBaseVt3D
	}

	maxv := gameobjtype.Attrib[gameobjtype.HomeMark].SpeedLimit
	bt.HomeMark = bt.NewGameObj(
		gameobjtype.HomeMark,
		bt.RandPosVt(),
		bt.RandVelVt(maxv),
	)
	maxv = gameobjtype.Attrib[gameobjtype.Ball].SpeedLimit
	bt.Ball = bt.NewGameObj(
		gameobjtype.Ball,
		bt.RandPosVt(),
		bt.RandVelVt(maxv),
	)
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
	return bt.BorderBounce.RandVector(bt.rnd.Float64)
}

func (bt *Team) RandBaseVt2D() vector3f.Vector3f {
	return vector3f.Vector3f{
		bt.rnd.Float64(),
		bt.rnd.Float64(),
		0,
	}
}
func (bt *Team) RandBaseVt3D() vector3f.Vector3f {
	return vector3f.Vector3f{
		bt.rnd.Float64(),
		bt.rnd.Float64(),
		bt.rnd.Float64(),
	}
}

func (bt *Team) RandVelVt(maxv float64) vector3f.Vector3f {
	return bt.RandBaseVtFn().MulF(maxv).NormalizedTo(maxv)
}

func (bt *Team) RandAccelVt() vector3f.Vector3f {
	return bt.RandBaseVtFn().MulF(gameconst.StageSize / 10)
}

func (bt *Team) RandShielVelVt() vector3f.Vector3f {
	return bt.RandBaseVtFn().MulF(gameconst.StageSize)
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
		o := bt.NewGameObj(gameobjtype.Bullet,
			bt.Ball.PosVt,
			actObj.Vt,
		)
		bt.addGObj(o)
	case acttype.BurstBullet:
		for i := 0; i < 10; i++ {
			vt := bt.RandPosVt()
			o := bt.NewGameObj(gameobjtype.BurstBullet,
				bt.Ball.PosVt,
				vt.Sub(bt.Ball.PosVt),
			)
			bt.addGObj(o)
		}
	case acttype.SuperBullet:
		o := bt.NewGameObj(gameobjtype.SuperBullet,
			bt.Ball.PosVt,
			actObj.Vt,
		)
		bt.addGObj(o)
	case acttype.HommingBullet:
		o := bt.NewGameObj(gameobjtype.HommingBullet,
			bt.Ball.PosVt,
			actObj.Vt,
		)
		o.DstUUID = actObj.DstObjID
		bt.addGObj(o)
	case acttype.Accel:
		bt.Ball.VelVt = bt.Ball.VelVt.Add(actObj.Vt)
	case acttype.Shield:
		o := bt.NewGameObj(gameobjtype.Shield,
			bt.Ball.PosVt,
			actObj.Vt,
		)
		bt.addGObj(o)
	case acttype.HommingShield:
		o := bt.NewGameObj(gameobjtype.HommingShield,
			bt.Ball.PosVt,
			actObj.Vt,
		)
		bt.addGObj(o)
	}
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

func (bt *Team) NewGameObj(
	gotype gameobjtype.GameObjType,
	at vector3f.Vector3f,
	velvt vector3f.Vector3f,
) *GameObj {
	nowtick := time.Now().UnixNano()
	o := &GameObj{
		TeamUUID:     bt.UUID,
		GOType:       gotype,
		UUID:         G_GameObjID.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        at,
		VelVt:        velvt,
		RotVelVt:     bt.RandRotVt(),
	}
	switch bt.StageType {
	default:
		bt.log.Fatal("invalid stagetype %v", bt.StageType)
	case stagetype.Stage2D:
		o.Move_circularFn = o.Move_circular2D
	case stagetype.Stage3D:
		o.Move_circularFn = o.Move_circular3D
	}

	return o
}

func (bt *Team) ToPacket() []*w3d_obj.GameObj {
	rtn := make([]*w3d_obj.GameObj, 0)
	co := uint32(bt.Color24)

	rtn = append(rtn, bt.Ball.ToPacket(co))
	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		rtn = append(rtn, v.ToPacket(co))
	}
	return rtn
}
