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

	"github.com/kasworld/gowasm2dgame/enums/acttype"
	"github.com/kasworld/gowasm2dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm2dgame/enums/teamtype"
	"github.com/kasworld/gowasm2dgame/game/gameconst"
	"github.com/kasworld/gowasm2dgame/lib/quadtreef"
	"github.com/kasworld/gowasm2dgame/lib/vector2f"
	"github.com/kasworld/gowasm2dgame/protocol_w2d/w2d_obj"
)

func (stg *Stage) SelectRandomTeam(me *Team) *Team {
	for i := 0; i < teamtype.TeamType_Count; i++ {
		dstteam := stg.Teams[stg.rnd.Intn(len(stg.Teams))]
		if dstteam != me && dstteam.IsAlive {
			return dstteam
		}
	}
	return nil
}

func (stg *Stage) FindDangerObj(
	me *Team, aienv *quadtreef.QuadTree) *GameObj {

	searchWHVt := vector2f.Vector2f{
		gameconst.StageW / 2,
		gameconst.StageH / 2,
	}
	searchRect := vector2f.NewRectCenterWH(
		me.Ball.PosVt,
		searchWHVt,
	)
	var findObj *GameObj
	aienv.QueryByRect(
		func(qo quadtreef.QuadTreeObjI) bool {
			targetObj := qo.(*GameObj)
			if targetObj.toDelete {
				return false
			}
			if targetObj.teamType == me.TeamType {
				return false
			}
			if _, lenchange := me.Ball.CalcLenChange(targetObj); lenchange < 0 {
				findObj = targetObj
				return true
			}
			return false
		},
		searchRect,
	)
	return findObj
}

func (stg *Stage) TryEvade(me *Team, now int64, dsto *GameObj) *w2d_obj.Act {
	actt := acttype.Accel
	objt := gameobjtype.Ball
	if me.GetRemainAct(now, actt) <= 0 {
		return nil
	}
	maxv := gameobjtype.Attrib[objt].SpeedLimit
	angle := dsto.VelVt.AddAngle(me.rnd.Float64()*math.Pi - math.Pi/2).Phase()
	return &w2d_obj.Act{
		Act:    actt,
		Angle:  angle,
		AngleV: maxv,
	}

}

func (stg *Stage) AI(me *Team, now int64, aienv *quadtreef.QuadTree) *w2d_obj.Act {
	dangerObj := stg.FindDangerObj(me, aienv)
	if dangerObj != nil {
		acto := stg.TryEvade(me, now, dangerObj)
		if acto != nil {
			return acto
		}
	}
	switch me.rnd.Intn(10) {
	default:
		//pass
	case 0:
		actt := acttype.Bullet
		objt := gameobjtype.Bullet
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		dstteam := stg.SelectRandomTeam(me)
		if dstteam == nil {
			break
		}
		angle, v := me.CalcAimAngleAndV(objt, dstteam.Ball)
		return &w2d_obj.Act{
			Act:    actt,
			Angle:  angle,
			AngleV: v,
		}
	case 1:
		actt := acttype.SuperBullet
		objt := gameobjtype.SuperBullet
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		dstteam := stg.SelectRandomTeam(me)
		if dstteam == nil {
			break
		}
		angle, v := me.CalcAimAngleAndV(objt, dstteam.Ball)
		return &w2d_obj.Act{
			Act:    actt,
			Angle:  angle,
			AngleV: v,
		}
	case 2:
		actt := acttype.HommingBullet
		objt := gameobjtype.HommingBullet
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		dstteam := stg.SelectRandomTeam(me)
		if dstteam == nil {
			break
		}
		maxv := gameobjtype.Attrib[objt].SpeedLimit
		if dstteam != me && dstteam.IsAlive {
			return &w2d_obj.Act{
				Act:      actt,
				Angle:    me.rnd.Float64() * 2 * math.Pi,
				AngleV:   maxv,
				DstObjID: dstteam.Ball.UUID,
			}
		}
	case 3:
		actt := acttype.Shield
		objt := gameobjtype.Shield
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		if me.Count(objt) < 12 {
			maxv := gameobjtype.Attrib[objt].SpeedLimit
			return &w2d_obj.Act{
				Act:    actt,
				Angle:  me.rnd.Float64() * 2 * math.Pi,
				AngleV: me.rnd.Float64() * maxv,
			}
		}
	case 4:
		actt := acttype.SuperShield
		objt := gameobjtype.SuperShield
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		if me.Count(objt) < 12 && me.rnd.Intn(10) == 0 {
			maxv := gameobjtype.Attrib[objt].SpeedLimit
			return &w2d_obj.Act{
				Act:    actt,
				Angle:  me.rnd.Float64() * 2 * math.Pi,
				AngleV: me.rnd.Float64() * maxv,
			}
		}
	case 5:
		actt := acttype.HommingShield
		objt := gameobjtype.HommingShield
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		if me.Count(objt) < 6 && me.rnd.Intn(10) == 0 {
			maxv := gameobjtype.Attrib[objt].SpeedLimit
			return &w2d_obj.Act{
				Act:    actt,
				Angle:  me.rnd.Float64() * 2 * math.Pi,
				AngleV: maxv,
			}
		}
	case 6:
		actt := acttype.Accel
		objt := gameobjtype.Ball
		if me.GetRemainAct(now, actt) <= 0 {
			break
		}
		maxv := gameobjtype.Attrib[objt].SpeedLimit
		return &w2d_obj.Act{
			Act:    actt,
			Angle:  me.rnd.Float64() * 2 * math.Pi,
			AngleV: me.rnd.Float64() * maxv,
		}
	}
	return &w2d_obj.Act{
		Act: acttype.Nothing,
	}
}
