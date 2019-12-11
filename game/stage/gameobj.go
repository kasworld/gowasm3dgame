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
	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
)

type GameObj struct {
	ObjType  gameobjtype.GameObjType
	UUID     string
	TeamUUID string
	PosVt    vector3f.Vector3f
	MvVt     vector3f.Vector3f

	BirthTick    int64
	LastMoveTick int64
	toDelete     bool

	dstUUID string
}

func (o *GameObj) Pos() vector3f.Vector3f {
	return o.PosVt
}

func (o *GameObj) IsCollision(dst *GameObj) bool {
	return gameobjtype.CollisionTo(
		o.ObjType, dst.ObjType,
		dst.PosVt.Sqd(o.PosVt),
	)
}
