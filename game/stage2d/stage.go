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

package stage2d

import (
	"context"
	"math/rand"
	"time"

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/config/serverconfig"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/game/background"
	"github.com/kasworld/gowasm3dgame/lib/idu64str"
	"github.com/kasworld/gowasm3dgame/lib/vector2f"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connbytemanager"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/htmlcolors"
)

var G_Stage2DID = idu64str.New("Stage2D")

type Stage struct {
	rnd    *rand.Rand      `prettystring:"hide"`
	log    *w3dlog.LogBase `prettystring:"hide"`
	config serverconfig.Config

	UUID  string
	Conns *w3d_connbytemanager.Manager

	BorderBounce vector3f.Cube
	BorderOctree vector3f.Cube

	Teams      []*Team
	Background *background.Background
}

func New(l *w3dlog.LogBase, config serverconfig.Config) *Stage {
	stg := &Stage{
		UUID:   G_Stage2DID.New(),
		config: config,
		log:    l,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Conns:  w3d_connbytemanager.New(),
	}

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
		stg.Teams = append(stg.Teams, NewTeam(l, v, stg.BorderBounce))
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
	if stg.rnd.Intn(100) == 0 {
		randvt := vector3f.Vector3f{
			stg.rnd.Float64() * gameconst.StageSize,
			stg.rnd.Float64() * gameconst.StageSize,
			gameconst.MaxRadius,
		}
		bt.HomeMark.AccelTo(randvt)
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
	rtn.BackgroundPos = [2]float32{
		float32(stg.Background.PosVt[0]),
		float32(stg.Background.PosVt[1]),
	}
	return rtn
}
