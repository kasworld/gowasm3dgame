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
	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/stagetype"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
)

var BoundCube = [stagetype.StageType_Count]vector3f.Cube{
	stagetype.Stage2D: {
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
	},
	stagetype.Stage3D: {
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
	},
}

var OctTreeCube = [stagetype.StageType_Count]vector3f.Cube{
	stagetype.Stage2D: {
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
	},
	stagetype.Stage3D: {
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
	},
}
