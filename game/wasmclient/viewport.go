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
	"syscall/js"

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
)

type Viewport struct {
	neecRecalc bool
	ViewWidth  int
	ViewHeight int
	RefSize    int

	Canvas   js.Value
	threejs  js.Value
	scene    js.Value
	camera   js.Value
	renderer js.Value
	light    js.Value

	jsSceneObjs   map[string]js.Value
	geometryCache map[gameobjtype.GameObjType]js.Value
	materialCache map[uint32]js.Value
}

func NewViewport(cnvid string) *Viewport {
	vp := &Viewport{
		jsSceneObjs:   make(map[string]js.Value),
		geometryCache: make(map[gameobjtype.GameObjType]js.Value),
		materialCache: make(map[uint32]js.Value),
	}

	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.Canvas = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvas3d").Call("appendChild", vp.Canvas)

	vp.scene = vp.ThreeJsNew("Scene")

	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 75, 1, gameobjtype.MaxRadius,
		gameconst.StageSize*10)

	vp.initGrid()
	vp.initLight()
	return vp
}

func (vp *Viewport) Hide() {
	vp.Canvas.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.neecRecalc = true
	vp.Canvas.Get("style").Set("display", "initial")
}

func (vp *Viewport) Resize() {
	// win := js.Global().Get("window")
	// winW := win.Get("innerWidth").Int()
	// winH := win.Get("innerHeight").Int()
	// if winW > winH {
	// 	vp.RefSize = winH / 32
	// } else {
	// 	vp.RefSize = winW / 32
	// }
	// if vp.RefSize < 32 {
	// 	vp.RefSize = 32
	// }

	vp.neecRecalc = true
}

func (vp *Viewport) Focus() {
	vp.Canvas.Call("focus")
}

func (vp *Viewport) DrawTitle() {
	// win := js.Global().Get("window")
	// winW := win.Get("innerWidth").Int()
	// winH := win.Get("innerHeight").Int()

}

func (vp *Viewport) Zoom(state int) {
	vp.neecRecalc = true
}

func (vp *Viewport) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.Canvas.Call("addEventListener", evt, js.FuncOf(fn))
}

func (vp *Viewport) calcResize() {
	if !vp.neecRecalc {
		return
	}
	vp.neecRecalc = false
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()
	size := winW
	if size > winH {
		size = winH
	}
	// size -= 20
	vp.ViewWidth = size
	vp.ViewHeight = size

	vp.Canvas.Call("setAttribute", "width", vp.ViewWidth)
	vp.Canvas.Call("setAttribute", "height", vp.ViewHeight)

	vp.renderer.Call("setSize", vp.ViewWidth, vp.ViewHeight)
}

func (vp *Viewport) Draw(tick int64) {
	vp.calcResize()

	vp.renderer.Call("render", vp.scene, vp.camera)
}
