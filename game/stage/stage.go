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
	"context"
	"math/rand"
	"time"

	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/game/gameconst"
	"github.com/kasworld/gowasm3dgame/game/serverconfig"
	"github.com/kasworld/gowasm3dgame/lib/octree"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connmanager"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/uuidstr"
)

type Stage struct {
	rnd    *rand.Rand      `prettystring:"hide"`
	log    *w3dlog.LogBase `prettystring:"hide"`
	config serverconfig.Config

	UUID  string
	Conns *w3d_connmanager.Manager

	BorderBounce vector3f.Cube
	BorderOctree vector3f.Cube

	Teams   []*Team
	Deco    []*GameObj
	Food    []*GameObj
	Terrain []*GameObj
}

func New(l *w3dlog.LogBase, config serverconfig.Config) *Stage {
	wd := &Stage{
		UUID:   uuidstr.New(),
		config: config,
		log:    l,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Conns:  w3d_connmanager.New(),
	}

	wd.BorderBounce = vector3f.Cube{
		Min: vector3f.Vector3f{
			-gameconst.StageSize / 2,
			-gameconst.StageSize / 2,
			-gameconst.StageSize / 2,
		},
		Max: vector3f.Vector3f{
			gameconst.StageSize / 2,
			gameconst.StageSize / 2,
			gameconst.StageSize / 2,
		},
	}
	wd.BorderOctree = vector3f.Cube{
		Min: vector3f.Vector3f{
			-gameconst.StageSize/2 - gameobjtype.MaxRadius,
			-gameconst.StageSize/2 - gameobjtype.MaxRadius,
			-gameconst.StageSize/2 - gameobjtype.MaxRadius,
		},
		Max: vector3f.Vector3f{
			gameconst.StageSize/2 + gameobjtype.MaxRadius,
			gameconst.StageSize/2 + gameobjtype.MaxRadius,
			gameconst.StageSize/2 + gameobjtype.MaxRadius,
		},
	}
	return wd
}

func (wd *Stage) MakeOctree() *octree.Octree {
	rtn := octree.New(wd.BorderOctree)
	for _, v := range wd.Teams {
		for _, o := range v.Objs {
			if o != nil && gameobjtype.Attrib[o.GOType].AddOctree {
				rtn.Insert(o)
			}
		}
	}
	return rtn
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
			si := stg.ToStatsInfo()
			conlist := stg.Conns.GetList()
			for _, v := range conlist {
				v.SendNotiPacket(w3d_idnoti.StatsInfo,
					si,
				)
			}
		case <-timerTurnTk.C:
			stg.Turn()
			si := stg.ToStageInfo()
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

	// respawn dead team
	for _, bt := range stg.Teams {
		if !bt.IsAlive && bt.RespawnTick < now {
			bt.RespawnBall(now)
		}
	}

	// aienv := stg.move(now)
	// for _, bt := range stg.Teams {
	// 	if !bt.IsAlive {
	// 		continue
	// 	}
	// 	actObj := stg.AI(bt, now, aienv)
	// 	if bt.GetRemainAct(now, actObj.Act) > 0 {
	// 		bt.ApplyAct(actObj)
	// 	} else {
	// 		stg.log.Fatal("OverAct %v %v", bt, actObj)
	// 	}
	// }
}

func (stg *Stage) move(now int64) *octree.Octree {

	for _, bt := range stg.Teams {
		if !bt.IsAlive {
			continue
		}
		toDelList := stg.MoveTeam(bt, now)
		_ = toDelList
	}
	aienv := stg.MakeOctree()
	// toDelList, aienv := stg.checkCollision()
	// for _, v := range toDelList {
	// 	stg.AddEffectByGameObj(v)
	// 	if v.GOType == gameobjtype.Ball {
	// 		stg.handleBallKilled(now, v)
	// 	}
	// }

	return aienv
}

func (stg *Stage) handleBallKilled(now int64, gobj *GameObj) {
	for _, bt := range stg.Teams {
		// find ballteam
		if bt.Ball.UUID == gobj.UUID {
			bt.IsAlive = false
			// regist respawn
			bt.RespawnTick = now + int64(time.Second)*gameconst.BallRespawnDurSec

			// add effect
			for _, v := range bt.Objs {
				if v.toDelete {
					continue
				}
				v.toDelete = true
			}
			return
		}
	}
	stg.log.Fatal("ball not in ballteam? %v", gobj)
}

func (stg *Stage) MoveTeam(bt *Team, now int64) []*GameObj {
	toDeleteList := make([]*GameObj, 0)
	bt.Ball.Move_accel(now)
	bt.Ball.BounceNormalize(gameconst.StageSize)
	for _, v := range bt.Objs {
		if v.toDelete {
			continue
		}
		switch v.GOType {
		default:
		case gameobjtype.Bullet, gameobjtype.SuperBullet:
			v.Move_accel(now)
			if !v.PosVt.IsIn(stg.BorderBounce) {
				v.toDelete = true
				toDeleteList = append(toDeleteList, v)
			}
		case gameobjtype.Shield:
			v.Move_circular(now, bt.Ball)
		case gameobjtype.HommingBullet:
			findDst := false
			for _, dstbt := range stg.Teams {
				if !dstbt.IsAlive {
					continue
				}
				if dstbt.Ball.UUID == v.DstUUID {
					findDst = true
					v.Move_homming(now, dstbt.Ball)
					break
				}
			}
			if !findDst {
				v.Move_accel(now)
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

func (stg *Stage) ToStageInfo() *w3d_obj.NotiStageInfo_data {
	now := time.Now().UnixNano()
	rtn := &w3d_obj.NotiStageInfo_data{
		Tick:         now,
		ID:           stg.UUID,
		BorderBounce: stg.BorderBounce,
		BorderOctree: stg.BorderOctree,
	}
	for _, bt := range stg.Teams {
		if !bt.IsAlive {
			continue
		}
		rtn.Teams = append(rtn.Teams, bt.ToPacket())
	}
	return rtn
}

func (stg *Stage) ToStatsInfo() *w3d_obj.NotiStatsInfo_data {
	rtn := &w3d_obj.NotiStatsInfo_data{}
	for _, bt := range stg.Teams {
		rtn.ActStats = append(rtn.ActStats, bt.ActStats)
	}
	return rtn
}
