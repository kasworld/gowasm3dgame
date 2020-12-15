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
	"context"
	"fmt"
	"time"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/config/serverconfig"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/enum/stagetype"
	"github.com/kasworld/gowasm3dgame/game/background"
	"github.com/kasworld/gowasm3dgame/lib/idu64str"
	"github.com/kasworld/gowasm3dgame/lib/octree"
	"github.com/kasworld/gowasm3dgame/lib/vector2f"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connbytemanager"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/htmlcolors"
)

var G_StageID = idu64str.New("Stage")

type Stage struct {
	rnd *g2rand.G2Rand  `prettystring:"hide"`
	log *w3dlog.LogBase `prettystring:"hide"`

	config serverconfig.Config

	StageType stagetype.StageType

	UUID  string
	Conns *w3d_connbytemanager.Manager

	BorderBounce vector3f.Cube
	BorderOctree vector3f.Cube

	Teams      []*Team
	Background *background.Background
}

func (stg *Stage) String() string {
	return fmt.Sprintf("Stage[%v %v Team:%v]",
		stg.StageType, stg.UUID,
		len(stg.Teams))
}

func (stg *Stage) GetUUID() string {
	return stg.UUID
}

func (stg *Stage) GetStageType() stagetype.StageType {
	return stg.StageType
}

func (stg *Stage) GetConnManager() *w3d_connbytemanager.Manager {
	return stg.Conns
}

func New(l *w3dlog.LogBase, config serverconfig.Config, seed int64, st stagetype.StageType) *Stage {
	stg := &Stage{
		StageType: st,
		UUID:      G_StageID.New(),
		config:    config,
		log:       l,
		rnd:       g2rand.NewWithSeed(seed),
		Conns:     w3d_connbytemanager.New(),
	}

	switch stg.StageType {
	case stagetype.Stage2D:
		stg.BorderBounce = vector3f.Cube{
			Min: vector3f.Vector3f{
				0,
				0,
				0,
			},
			Max: vector3f.Vector3f{
				gameconst.StageSize,
				gameconst.StageSize,
				gameconst.MaxRadius,
			},
		}
		stg.BorderOctree = vector3f.Cube{
			Min: vector3f.Vector3f{
				-gameconst.MaxRadius,
				-gameconst.MaxRadius,
				-gameconst.MaxRadius,
			},
			Max: vector3f.Vector3f{
				gameconst.StageSize + gameconst.MaxRadius,
				gameconst.StageSize + gameconst.MaxRadius,
				gameconst.MaxRadius + gameconst.MaxRadius,
			},
		}

	case stagetype.Stage3D:
		stg.BorderBounce = vector3f.Cube{
			Min: vector3f.Vector3f{
				0,
				0,
				0,
			},
			Max: vector3f.Vector3f{
				gameconst.StageSize,
				gameconst.StageSize,
				gameconst.StageSize,
			},
		}
		stg.BorderOctree = vector3f.Cube{
			Min: vector3f.Vector3f{
				-gameconst.MaxRadius,
				-gameconst.MaxRadius,
				-gameconst.MaxRadius,
			},
			Max: vector3f.Vector3f{
				gameconst.StageSize + gameconst.MaxRadius,
				gameconst.StageSize + gameconst.MaxRadius,
				gameconst.StageSize + gameconst.MaxRadius,
			},
		}
	}
	teamcolor := make([]htmlcolors.Color24, 0)
	for i := 0; i < gameconst.TeamPerStage; i++ {
		co := htmlcolors.NewColor24(
			uint8(stg.rnd.Intn(256)),
			uint8(stg.rnd.Intn(256)),
			uint8(stg.rnd.Intn(256)),
		)
		teamcolor = append(teamcolor, co)
	}
	for _, v := range teamcolor {
		tm := NewTeam(l, v, stg.BorderBounce, stg.rnd.Int63(), stg.StageType)
		stg.Teams = append(stg.Teams, tm)
	}
	stg.Background = background.New(
		time.Now().UnixNano(),
		vector2f.NewVectorLenAngle(
			stg.rnd.Float64()*300,
			stg.rnd.Float64()*360,
		),
		vector2f.Rect{
			vector2f.Vector2f{0, 0},
			vector2f.Vector2f{gameconst.StageSize, gameconst.StageSize},
		})
	return stg
}

func (stg *Stage) Run(ctx context.Context) {

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	turnDur := time.Duration(float64(time.Second) / stg.config.ActTurnPerSec)
	timerTurnTk := time.NewTicker(turnDur)
	defer timerTurnTk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerInfoTk.C:

		case <-timerTurnTk.C:
			stg.Turn()
			si := stg.ToPacket_StageInfo()
			conlist := stg.Conns.GetList()
			for _, v := range conlist {
				v.SendNotiPacket(w3d_idnoti.StageInfo,
					si,
				)
			}
		}
	}
}

func (stg *Stage) Turn() {
	now := time.Now().UnixNano()
	diag := stg.BorderBounce.DiagLen()

	stg.Background.Move(now)
	// respawn dead team
	for _, bt := range stg.Teams {
		if !bt.IsAlive && bt.RespawnTick < now {
			bt.RespawnBall(now)
		}
	}

	for _, bt := range stg.Teams {
		bt.ActPoint += bt.CalcAP(diag)
		if !bt.IsAlive {
			continue
		}
		toDelList := stg.MoveTeam(bt, now)
		_ = toDelList
	}

	collisionList, aienv := stg.checkCollision()
	for _, v := range collisionList {
		v[0].toDelete = true
	}

	for _, v := range collisionList {
		if v[0].GOType == gameobjtype.Ball {
			stg.handleBallKilled(now, v)
		}
	}

	for _, bt := range stg.Teams {
		if !bt.IsAlive {
			continue
		}
		actObj := stg.AI(bt, now, aienv)
		if bt.CanAct(actObj.Act) {
			bt.ApplyAct(actObj)
		} else {
			stg.log.Fatal("OverAct %v %v", bt, actObj)
		}
	}
}

func (stg *Stage) getTeamByUUID(id string) *Team {
	for _, bt := range stg.Teams {
		if bt.UUID == id {
			return bt
		}
	}
	return nil
}
func (stg *Stage) handleBallKilled(now int64, gobj [2]*GameObj) {
	bt := stg.getTeamByUUID(gobj[0].TeamUUID)
	if bt == nil {
		stg.log.Fatal("invalid team uuid %v", gobj)
		return
	}
	killbt := stg.getTeamByUUID(gobj[1].TeamUUID)
	if killbt == nil {
		stg.log.Fatal("invalid team uuid %v", gobj)
	}
	killbt.Kill++

	bt.IsAlive = false
	bt.Death++
	// regist respawn
	bt.RespawnTick = now + int64(time.Second)*gameconst.BallRespawnDurSec

	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		v.toDelete = true
	}

}

func (stg *Stage) MoveTeam(bt *Team, now int64) []*GameObj {
	toDeleteList := make([]*GameObj, 0)
	bt.Ball.Move_straight(now)
	bt.Ball.BounceNormalize(bt.BorderBounce)

	bt.HomeMark.Move_straight(now)
	bt.HomeMark.BounceNormalize(bt.BorderBounce)
	switch stg.StageType {
	case stagetype.Stage2D:
		if stg.rnd.Intn(100) == 0 {
			randvt := vector3f.Vector3f{
				stg.rnd.Float64() * gameconst.StageSize,
				stg.rnd.Float64() * gameconst.StageSize,
				gameconst.MaxRadius,
			}
			bt.HomeMark.AccelTo(randvt)
		}
	case stagetype.Stage3D:
		if stg.rnd.Intn(100) == 0 {
			randvt := vector3f.Vector3f{
				stg.rnd.Float64() * gameconst.StageSize,
				stg.rnd.Float64() * gameconst.StageSize,
				stg.rnd.Float64() * gameconst.StageSize,
			}
			bt.HomeMark.AccelTo(randvt)
		}
	}

	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		switch v.GOType {
		default:
		case gameobjtype.Bullet, gameobjtype.SuperBullet, gameobjtype.BurstBullet:
			v.Move_straight(now)
			if !v.PosVt.IsIn(stg.BorderBounce) {
				v.toDelete = true
				toDeleteList = append(toDeleteList, v)
			}
		case gameobjtype.Shield:
			v.Move_circular(now, bt.Ball)

		case gameobjtype.HommingShield:
			v.Move_hommingshield(now, bt.Ball)
			break

		case gameobjtype.HommingBullet:
			findDst := false
			for _, dstbt := range stg.Teams {
				if !dstbt.IsAlive {
					continue
				}
				if dstbt.Ball.UUID == v.DstUUID {
					findDst = true
					v.Move_hommingbullet(now, dstbt.Ball)
					break
				}
			}
			if !findDst {
				v.Move_straight(now)
				if !v.PosVt.IsIn(stg.BorderBounce) {
					v.toDelete = true
					toDeleteList = append(toDeleteList, v)
				}
			}
		}
		if !v.toDelete && !v.CheckLife(now) {
			v.toDelete = true
			toDeleteList = append(toDeleteList, v)
		}
	}
	return toDeleteList
}

func (stg *Stage) ToPacket_StageInfo() *w3d_obj.NotiStageInfo_data {
	now := time.Now().UnixNano()
	rtn := &w3d_obj.NotiStageInfo_data{
		Tick: now,
	}
	for _, bt := range stg.Teams {
		if !bt.IsAlive {
			continue
		}
		rtn.ObjList = append(rtn.ObjList, bt.ToPacket()...)
		ltPos := bt.HomeMark.PosVt
		rtn.Lights = append(rtn.Lights, &w3d_obj.Light{
			UUID: bt.UUID,
			PosVt: [3]float32{
				float32(ltPos[0]),
				float32(ltPos[1]),
				float32(ltPos[2]),
			},
			Color: uint32(bt.Color24),
		})
	}
	switch stg.StageType {
	case stagetype.Stage2D:
		rtn.CameraPos = [3]float32{
			gameconst.StageSize / 2,
			gameconst.StageSize / 2,
			gameconst.StageSize,
		}
		rtn.CameraLookAt = [3]float32{
			gameconst.StageSize / 2,
			gameconst.StageSize / 2,
			0,
		}
	case stagetype.Stage3D:
		if len(rtn.Lights) > 0 {
			rtn.CameraPos = rtn.Lights[0].PosVt
		}
		if len(rtn.ObjList) > 0 {
			rtn.CameraLookAt = rtn.ObjList[0].PosVt
		}
	}
	rtn.BackgroundPos = [2]float32{
		float32(stg.Background.PosVt[0]),
		float32(stg.Background.PosVt[1]),
	}
	return rtn
}

func (stg *Stage) ToPacket_StatsInfo() *w3d_obj.RspStatsInfo_data {
	rtn := &w3d_obj.RspStatsInfo_data{}
	for _, bt := range stg.Teams {
		teamStats := w3d_obj.TeamStat{
			UUID:     bt.UUID,
			Alive:    bt.IsAlive,
			AP:       int(bt.ActPoint),
			Score:    int(bt.Score),
			Kill:     bt.Kill,
			Death:    bt.Death,
			Color24:  uint32(bt.Color24),
			ActStats: bt.ActStats,
		}
		rtn.Stats = append(rtn.Stats, teamStats)
	}
	return rtn
}

func (stg *Stage) checkCollision() ([][2]*GameObj, *octree.Octree) {
	collisionList := make([][2]*GameObj, 0)
	ot := octree.New(stg.BorderOctree)
	obj2check := make([]*GameObj, 0)
	for _, bt := range stg.Teams {
		if !bt.IsAlive {
			continue
		}
		if ot.Insert(bt.Ball) {
			obj2check = append(obj2check, bt.Ball)
		}
		for _, v := range bt.Objs {
			if v.toDelete {
				continue
			}
			if ot.Insert(v) {
				obj2check = append(obj2check, v)
			}
		}
	}
	for _, v := range obj2check {
		if v.toDelete {
			continue
		}
		ot.QueryByCube(
			func(qo octree.OctreeObjI) bool {
				targetObj := qo.(*GameObj)
				if targetObj.toDelete {
					return false
				}
				_ = targetObj
				if targetObj.TeamUUID == v.TeamUUID {
					return false
				}
				if !v.toDelete && !targetObj.toDelete {
					if gameobjtype.InteractTo(v.GOType, targetObj.GOType) {
						collisionList = append(collisionList,
							[2]*GameObj{v, targetObj},
						)
						return true
					}
				}
				return false
			},
			v.GetCube(),
		)
	}
	return collisionList, ot
}
