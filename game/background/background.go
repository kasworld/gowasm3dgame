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

package background

import (
	"time"

	"github.com/kasworld/gowasm3dgame/lib/vector2f"
)

type Background struct {
	LastMoveTick int64 // time.unixnano
	PosVt        vector2f.Vector2f
	VelVt        vector2f.Vector2f
	Bound        vector2f.Rect
}

func (o *Background) Move(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))
	o.PosVt = o.Bound.WrapVector(o.PosVt)
}

func New(
	now int64,
	velvt vector2f.Vector2f,
	boundrt vector2f.Rect,
) *Background {
	return &Background{
		LastMoveTick: now,
		VelVt:        velvt,
		Bound:        boundrt,
	}
}
