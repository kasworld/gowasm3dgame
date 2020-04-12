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
	"github.com/kasworld/gowasm2dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm2dgame/game/gameconst"
	"github.com/kasworld/gowasm2dgame/lib/quadtreef"
	"github.com/kasworld/gowasm2dgame/lib/vector2f"
)

func (stg *Stage) newQtree() *quadtreef.QuadTree {
	maxr := 32.0
	qtree := quadtreef.New(vector2f.Rect{
		vector2f.Vector2f{0 - maxr, 0 - maxr},
		vector2f.Vector2f{gameconst.StageW + maxr, gameconst.StageH + maxr},
	})
	return qtree
}

func (stg *Stage) checkCollision() ([]*GameObj, *quadtreef.QuadTree) {
	toDeleteList := make([]*GameObj, 0)
	qtree := stg.newQtree()
	obj2check := make([]*GameObj, 0)
	for _, bt := range stg.Teams {
		if qtree.Insert(bt.Ball) {
			obj2check = append(obj2check, bt.Ball)
		}
		for _, v := range bt.Objs {
			if v.toDelete {
				continue
			}
			if qtree.Insert(v) {
				obj2check = append(obj2check, v)
			}
		}
	}
	for _, v := range obj2check {
		if v.toDelete {
			continue
		}
		qtree.QueryByRect(
			func(qo quadtreef.QuadTreeObjI) bool {
				targetObj := qo.(*GameObj)
				if targetObj.toDelete {
					return false
				}
				_ = targetObj
				if targetObj.teamType == v.teamType {
					return false
				}
				if !v.toDelete && !targetObj.toDelete {
					if gameobjtype.InteractTo(v.GOType, targetObj.GOType) {
						toDeleteList = append(toDeleteList, v)
						return true
					}
				}
				return false
			},
			v.GetRect(),
		)
	}
	for _, v := range toDeleteList {
		v.toDelete = true
	}
	return toDeleteList, qtree
}
