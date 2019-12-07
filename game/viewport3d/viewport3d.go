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

package viewport3d

import (
	"fmt"
	"syscall/js"

	"github.com/kasworld/direction"
)

type Viewport3d struct {
	canvas js.Value
	ctx3d  js.Value

	zoomState  int
	neecRecalc bool

	ViewWidth  int
	ViewHeight int

	MoveType int
	Grid     bool

	MouseX   int
	MouseY   int
	MouseDir direction.Direction_Type

	cube js.Value
}

func New(cnvid string) *Viewport3d {
	vp := &Viewport3d{}

	vp.canvas = js.Global().Get("document").Call("getElementById", cnvid)
	if !vp.canvas.Truthy() {
		fmt.Printf("fail to get canvas\n")
	}
	vp.ctx3d = vp.canvas.Call("getContext", "webgl")
	if !vp.ctx3d.Truthy() {
		fmt.Printf("fail to get context\n")
	}
	return vp
}

func (vp *Viewport3d) Hide() {
	vp.canvas.Get("style").Set("display", "none")
}
func (vp *Viewport3d) Show() {
	vp.neecRecalc = true
	vp.canvas.Get("style").Set("display", "initial")
}

func (vp *Viewport3d) Resize() {
	vp.neecRecalc = true
}

func (vp *Viewport3d) Focus() {
	vp.canvas.Call("focus")
}

func (vp *Viewport3d) Zoom(state int) {
	vp.zoomState = state
	vp.neecRecalc = true
}

func (vp *Viewport3d) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.canvas.Call("addEventListener", evt, js.FuncOf(fn))
}

func (vp *Viewport3d) calcViewCellValue() {
	if !vp.neecRecalc {
		return
	}
	vp.neecRecalc = false
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()
	if winH > winW {
		winH /= 2
	} else {
		winW /= 2
	}

	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.canvas.Call("setAttribute", "width", vp.ViewWidth)
	vp.canvas.Call("setAttribute", "height", vp.ViewHeight)
}

func (vp *Viewport3d) Draw(tick int64) {
	vp.calcViewCellValue()

	// rot := vp.cube.Get("rotation")
	// rot.Set("x", rot.Get("x").Float()+0.01)
	// rot.Set("y", rot.Get("y").Float()+0.01)
	// vp.renderer.Call("render", vp.scene, vp.camera)
}
