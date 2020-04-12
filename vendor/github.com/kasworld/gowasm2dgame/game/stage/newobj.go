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
	"time"

	"github.com/kasworld/gowasm2dgame/enums/effecttype"
	"github.com/kasworld/gowasm2dgame/game/gameconst"
	"github.com/kasworld/gowasm2dgame/lib/vector2f"
	"github.com/kasworld/gowasm2dgame/protocol_w2d/w2d_obj"
)

func (stg *Stage) NewBackground() *w2d_obj.Background {
	vt := vector2f.NewVectorLenAngle(
		stg.rnd.Float64()*300,
		stg.rnd.Float64()*360,
	)
	return &w2d_obj.Background{
		LastMoveTick: time.Now().UnixNano(),
		VelVt:         vt,
	}
}

func (stg *Stage) NewCloud(i int) *w2d_obj.Cloud {
	vt := vector2f.NewVectorLenAngle(
		stg.rnd.Float64()*300,
		stg.rnd.Float64()*360,
	)
	return &w2d_obj.Cloud{
		SpriteNum: i,
		PosVt: vector2f.Vector2f{
			stg.rnd.Float64() * gameconst.StageW,
			stg.rnd.Float64() * gameconst.StageH,
		},
		VelVt:         vt,
		LastMoveTick: time.Now().UnixNano(),
	}
}

func (stg *Stage) AddEffect(
	et effecttype.EffectType, posVt, mvVt vector2f.Vector2f) {
	now := time.Now().UnixNano()
	o := &w2d_obj.Effect{
		EffectType:   et,
		BirthTick:    now,
		LastMoveTick: now,
		PosVt:        posVt,
		VelVt:         mvVt,
	}
	for i, v := range stg.Effects {
		if !v.CheckLife(now) {
			stg.Effects[i] = o
			return
		}
	}
	stg.Effects = append(stg.Effects, o)
}
