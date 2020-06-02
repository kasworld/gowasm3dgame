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

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/acttype"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/octree"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

func (stg *Stage) SelectRandomTeam(me *Team) *Team {
	for i := 0; i < gameconst.TeamPerStage; i++ {
		dstteam := stg.Teams[stg.rnd.Intn(len(stg.Teams))]
		if dstteam != me && dstteam.IsAlive {
			return dstteam
		}
	}
	return nil
}

func (stg *Stage) FindDangerObj(
	me *Team, aienv *octree.Octree) *GameObj {

	searchRect := vector3f.NewCubeByCR(
		me.Ball.PosVt,
		gameconst.StageSize/4,
	)
	var findObj *GameObj
	aienv.QueryByCube(
		func(qo octree.OctreeObjI) bool {
			targetObj := qo.(*GameObj)
			if targetObj.toDelete {
				return false
			}
			if targetObj.TeamUUID == me.UUID {
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

func (stg *Stage) TryEvade(me *Team, now int64, dsto *GameObj) *w3d_obj.Act {
	actt := acttype.Accel
	objt := gameobjtype.Ball
	if !me.CanAct(actt) {
		return nil
	}
	maxv := Attrib[objt].SpeedLimit
	return &w3d_obj.Act{
		Act: actt,
		Vt:  dsto.VelVt.NormalizedTo(maxv),
	}

}

func (stg *Stage) AI(me *Team, now int64, aienv *octree.Octree) *w3d_obj.Act {
	dangerObj := stg.FindDangerObj(me, aienv)
	if dangerObj != nil {
		acto := stg.TryEvade(me, now, dangerObj)
		if acto != nil {
			return acto
		}
	}
	switch me.rnd.Intn(10) {
	// switch 4 {
	default:
		//pass
	case 0:
		actt := acttype.Bullet
		objt := gameobjtype.Bullet
		if !me.CanAct(actt) {
			break
		}
		if !me.CanHave(objt) {
			break
		}
		dstteam := stg.SelectRandomTeam(me)
		if dstteam == nil {
			break
		}
		_, estpos, _ := stg.calcAims(me, dstteam.Ball,
			Attrib[objt].SpeedLimit)
		vt := stg.AimAdjedIntoCube(me, estpos, dstteam.Ball, objt)
		return &w3d_obj.Act{
			Act: actt,
			Vt:  vt,
		}
	case 1: // 10 random bullet
		actt := acttype.BurstBullet
		objt := gameobjtype.BurstBullet
		if !me.CanAct(actt) {
			break
		}
		if !me.CanHave(objt) {
			break
		}
		return &w3d_obj.Act{
			Act: actt,
		}
	case 2:
		actt := acttype.SuperBullet
		objt := gameobjtype.SuperBullet
		if !me.CanAct(actt) {
			break
		}
		if !me.CanHave(objt) {
			break
		}
		dstteam := stg.SelectRandomTeam(me)
		if dstteam == nil {
			break
		}
		_, estpos, _ := stg.calcAims(me, dstteam.Ball,
			Attrib[objt].SpeedLimit)
		vt := stg.AimAdjedIntoCube(me, estpos, dstteam.Ball, objt)
		return &w3d_obj.Act{
			Act: actt,
			Vt:  vt,
		}
	case 3:
		actt := acttype.HommingBullet
		objt := gameobjtype.HommingBullet
		if !me.CanAct(actt) {
			break
		}
		if !me.CanHave(objt) {
			break
		}
		dstteam := stg.SelectRandomTeam(me)
		if dstteam == nil {
			break
		}
		maxv := Attrib[objt].SpeedLimit
		if dstteam != me && dstteam.IsAlive {
			return &w3d_obj.Act{
				Act:      actt,
				Vt:       me.Ball.VelVt.NormalizedTo(maxv),
				DstObjID: dstteam.Ball.UUID,
			}
		}
	case 4:
		actt := acttype.HommingShield
		objt := gameobjtype.HommingShield
		if !me.CanAct(actt) {
			break
		}
		if !me.CanHave(objt) {
			break
		}
		dstteam := me
		if !dstteam.IsAlive {
			break
		}
		maxv := Attrib[objt].SpeedLimit
		return &w3d_obj.Act{
			Act: actt,
			Vt:  me.Ball.VelVt.NormalizedTo(maxv).Neg(),
		}

	case 5:
		actt := acttype.Accel
		// objt := gameobjtype.Ball
		if !me.CanAct(actt) {
			break
		}
		// maxv := Attrib[objt].SpeedLimit
		return &w3d_obj.Act{
			Act: actt,
			Vt: vector3f.Vector3f{
				me.rnd.Float64() * gameconst.StageSize / 10,
				me.rnd.Float64() * gameconst.StageSize / 10,
				me.rnd.Float64() * gameconst.StageSize / 10,
			},
		}
	case 6:
		actt := acttype.Shield
		objt := gameobjtype.Shield
		if !me.CanAct(actt) {
			break
		}
		if !me.CanHave(objt) {
			break
		}
		vt := vector3f.Vector3f{
			me.rnd.Float64() * gameconst.StageSize,
			me.rnd.Float64() * gameconst.StageSize,
			me.rnd.Float64() * gameconst.StageSize,
		}
		return &w3d_obj.Act{
			Act: actt,
			Vt:  vt,
		}
	}
	return &w3d_obj.Act{
		Act: acttype.Nothing,
	}
}

func (stg *Stage) AimAdjedIntoCube(
	tm *Team,
	estpos vector3f.Vector3f,
	o *GameObj,
	bulletType gameobjtype.GameObjType) vector3f.Vector3f {

	if !estpos.IsIn(stg.BorderOctree) && o.GOType != gameobjtype.Ball {
		return vector3f.VtZero
	}
	lenori := tm.Ball.PosVt.LenTo(estpos)
	if o.GOType == gameobjtype.Ball {
		var changed int
		estpos, changed = estpos.MakeIn(stg.BorderBounce)
		if changed != 0 {
			//log.Printf("target %v bounce %b", o.ID, changed)
		}
	}
	lennew := tm.Ball.PosVt.LenTo(estpos)
	lenrate := lennew / lenori
	vt := estpos.Sub(tm.Ball.PosVt).NormalizedTo(
		Attrib[bulletType].SpeedLimit).MulF(lenrate)
	return vt
}

func (stg *Stage) calcAims(
	tm *Team,
	t *GameObj,
	projectilemovelimit float64) (float64, vector3f.Vector3f, float64) {

	dur := tm.Ball.PosVt.CalcAimAheadDur(t.PosVt, t.VelVt, projectilemovelimit)
	if math.IsInf(dur, 1) {
		return math.Inf(1), vector3f.VtZero, 0
	}
	if t.VelVt.Abs() < 0.1 {
		return dur, t.PosVt, 0
	}
	estpos := t.PosVt.Add(t.VelVt.MulF(dur))
	estangle := t.PosVt.Sub(tm.Ball.PosVt).Angle(estpos.Sub(tm.Ball.PosVt))
	return dur, estpos, estangle
}
