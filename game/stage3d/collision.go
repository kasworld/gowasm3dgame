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
	"time"

	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/octree"
)

const LongLife = 3600 * 24 * 365

var Attrib = [gameobjtype.GameObjType_Count]struct {
	SpeedLimit float64
	Radius     float64
	MaxInTeam  int
	AddOctree  bool
	LifeTick   int64
}{
	gameobjtype.Ball:          {300, 10, 1, true, int64(time.Second) * LongLife},
	gameobjtype.Shield:        {400, 5, 32, true, int64(time.Second) * LongLife},
	gameobjtype.HommingShield: {400, 7, 16, true, int64(time.Second) * 60},
	gameobjtype.Bullet:        {300, 5, 100, true, int64(time.Second) * 10},
	gameobjtype.HommingBullet: {300, 7, 10, true, int64(time.Second) * 60},
	gameobjtype.SuperBullet:   {600, 15, 10, true, int64(time.Second) * 10},
	gameobjtype.BurstBullet:   {300, 5, 100, true, int64(time.Second) * 10},
	gameobjtype.HomeMark:      {100, 3, 1, false, int64(time.Second) * LongLife},
	gameobjtype.Deco:          {600, 3, 100, false, int64(time.Second) * LongLife},
	gameobjtype.Hard:          {0, 3, 1, false, int64(time.Second) * LongLife},
	gameobjtype.Food:          {0, 3, 1, true, int64(time.Second) * LongLife},
}

// for client
func (stg *Stage) MakeType2Radius() [gameobjtype.GameObjType_Count]float64 {
	var rtn [gameobjtype.GameObjType_Count]float64
	for i, v := range Attrib {
		rtn[i] = v.Radius
	}
	return rtn
}

var interactRule = [gameobjtype.GameObjType_Count][gameobjtype.GameObjType_Count]bool{
	gameobjtype.Ball: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.Shield: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.HommingShield: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.Bullet: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.HommingBullet: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.SuperBullet: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.BurstBullet: {
		gameobjtype.Ball:          true,
		gameobjtype.Shield:        true,
		gameobjtype.HommingShield: true,
		gameobjtype.Bullet:        true,
		gameobjtype.HommingBullet: true,
		gameobjtype.SuperBullet:   true,
		gameobjtype.BurstBullet:   true,
		gameobjtype.Hard:          true,
	},
	gameobjtype.HomeMark: {},
	gameobjtype.Deco:     {},
	gameobjtype.Hard:     {},
	gameobjtype.Food: {
		gameobjtype.Ball: true,
	},
}

func InteractTo(srcType, dstType gameobjtype.GameObjType) bool {
	return interactRule[srcType][dstType]
}

func CollisionTo(srcType, dstType gameobjtype.GameObjType, sqd float64) bool {
	l := Attrib[srcType].Radius + Attrib[dstType].Radius
	return interactRule[srcType][dstType] && sqd <= l*l
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
					if InteractTo(v.GOType, targetObj.GOType) {
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
	// fmt.Printf("obj check len %v, del %v\n",
	// 	len(obj2check), len(collisionList))
	return collisionList, ot
}
