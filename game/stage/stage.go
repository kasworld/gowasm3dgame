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
	"math/rand"
	"time"

	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/game/gameconst"
	"github.com/kasworld/gowasm3dgame/game/serverconfig"
	"github.com/kasworld/gowasm3dgame/lib/octree"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connmanager"
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
			-gameconst.WorldSize / 2,
			-gameconst.WorldSize / 2,
			-gameconst.WorldSize / 2,
		},
		Max: vector3f.Vector3f{
			gameconst.WorldSize / 2,
			gameconst.WorldSize / 2,
			gameconst.WorldSize / 2,
		},
	}
	wd.BorderOctree = vector3f.Cube{
		Min: vector3f.Vector3f{
			-gameconst.WorldSize/2 - gameobjtype.MaxRadius,
			-gameconst.WorldSize/2 - gameobjtype.MaxRadius,
			-gameconst.WorldSize/2 - gameobjtype.MaxRadius,
		},
		Max: vector3f.Vector3f{
			gameconst.WorldSize/2 + gameobjtype.MaxRadius,
			gameconst.WorldSize/2 + gameobjtype.MaxRadius,
			gameconst.WorldSize/2 + gameobjtype.MaxRadius,
		},
	}
	return wd
}

func (wd *Stage) MakeOctree() *octree.Octree {
	rtn := octree.New(wd.BorderOctree)
	for _, v := range wd.Teams {
		for _, o := range v.Objs {
			if o != nil && gameobjtype.Attrib[o.ObjType].AddOctree {
				rtn.Insert(o)
			}
		}
	}
	return rtn
}
