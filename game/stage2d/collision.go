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
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/octree"
)

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
	// fmt.Printf("obj check len %v, del %v\n",
	// 	len(obj2check), len(collisionList))
	return collisionList, ot
}
