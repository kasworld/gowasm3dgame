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

package stage3d

import (
	"math"
	"math/rand"
	"time"

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/acttype"
	"github.com/kasworld/gowasm3dgame/enum/acttype_vector"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/idu64str"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/htmlcolors"
)

var G_Team3DID = idu64str.New("Team3D")
var G_GameObj3DID = idu64str.New("GameObj3D")

type Team struct {
	rnd *rand.Rand      `prettystring:"hide"`
	log *w3dlog.LogBase `prettystring:"hide"`

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
}

func NewTeam(l *w3dlog.LogBase, color htmlcolors.Color24, BorderBounce vector3f.Cube) *Team {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	bt := &Team{
		rnd:          rnd,
		log:          l,
		UUID:         G_Team3DID.New(),
		BorderBounce: BorderBounce,
		IsAlive:      true,
		Color24:      color,
		Objs:         make([]*GameObj, 0),
	}

	maxv := Attrib[gameobjtype.HomeMark].SpeedLimit
	bt.HomeMark = bt.NewGameObj(
		gameobjtype.HomeMark,
		bt.RandPosVt(),
		vector3f.Vector3f{
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
		}.NormalizedTo(maxv),
	)

	maxv = Attrib[gameobjtype.Ball].SpeedLimit
	bt.Ball = bt.NewGameObj(
		gameobjtype.Ball,
		bt.RandPosVt(),
		vector3f.Vector3f{
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
			bt.rnd.Float64() * maxv,
		}.NormalizedTo(maxv),
	)

	return bt
}

func (bt *Team) RespawnBall(now int64) {
	bt.log.Debug("ball respawned %v", bt)
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
	svt := bt.BorderBounce.SizeVector()
	return vector3f.Vector3f{
		bt.rnd.Float64()*svt[0] + bt.BorderBounce.Min[0],
		bt.rnd.Float64()*svt[1] + bt.BorderBounce.Min[1],
		bt.rnd.Float64()*svt[2] + bt.BorderBounce.Min[2],
	}
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
	return bt.CountByGOType(objt) <= Attrib[objt].MaxInTeam
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
		UUID:         G_GameObj3DID.New(),
		BirthTick:    nowtick,
		LastMoveTick: nowtick,
		PosVt:        at,
		VelVt:        velvt,
		RotVelVt:     bt.RandRotVt(),
	}
	return o
}

func (bt *Team) ToPacket() []*w3d_obj.GameObj {
	rtn := make([]*w3d_obj.GameObj, 0)
	co := uint32(bt.Color24)

	// replaced by lighthelper
	// rtn = append(rtn, bt.HomeMark.ToPacket(co))
	rtn = append(rtn, bt.Ball.ToPacket(co))
	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		rtn = append(rtn, v.ToPacket(co))
	}
	return rtn
}
