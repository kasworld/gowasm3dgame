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

package wasmclient

import (
	"math/rand"
	"syscall/js"
	"time"
)

type Viewport3d struct {
	Canvas    js.Value
	context2d js.Value
	rnd       *rand.Rand
}

func NewViewport3d() *Viewport3d {
	vp := &Viewport3d{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	vp.Canvas, vp.context2d = getCnv2dCtx("viewport2DCanvas")
	return vp
}

func (vp *Viewport3d) draw(now int64) {
}
