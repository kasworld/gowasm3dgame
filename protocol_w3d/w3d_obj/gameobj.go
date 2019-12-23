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

package w3d_obj

import (
	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/htmlcolors"
)

type Team struct {
	ID       string
	Color24  htmlcolors.Color24
	Ball     *GameObj
	HomeMark *GameObj
	Objs     []*GameObj
}

type GameObj struct {
	GOType gameobjtype.GameObjType
	UUID   string
	PosVt  vector3f.Vector3f
	RotVt  vector3f.Vector3f
}
