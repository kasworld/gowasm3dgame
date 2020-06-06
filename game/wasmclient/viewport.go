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

package wasmclient

import (
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
)

type LightNHelper struct {
	Light  js.Value
	Helper js.Value
}

type Viewport struct {
	rnd        *rand.Rand
	ViewWidth  int
	ViewHeight int
	RefSize    int

	Canvas   js.Value
	threejs  js.Value
	scene    js.Value
	camera   js.Value
	renderer js.Value

	// title
	fontLoader              js.Value
	font_helvetiker_regular js.Value
	jsoTitle                js.Value
	lightTitle              js.Value

	// background
	textureLoader js.Value
	background    js.Value

	Type2Radius [gameobjtype.GameObjType_Count]float64

	jsSceneObjs   map[string]js.Value
	geometryCache map[gameobjtype.GameObjType]js.Value
	materialCache map[uint32]js.Value

	lightCache map[string]*LightNHelper
}

func NewViewport() *Viewport {
	vp := &Viewport{
		rnd:           rand.New(rand.NewSource(time.Now().UnixNano())),
		jsSceneObjs:   make(map[string]js.Value),
		geometryCache: make(map[gameobjtype.GameObjType]js.Value),
		materialCache: make(map[uint32]js.Value),
		lightCache:    make(map[string]*LightNHelper),
	}

	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.Canvas = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvas3dholder").Call("appendChild", vp.Canvas)
	vp.Canvas.Set("tabindex", "1")

	vp.scene = vp.ThreeJsNew("Scene")

	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 75, 1, gameconst.MaxRadius,
		gameconst.StageSize*10)

	vp.textureLoader = vp.ThreeJsNew("TextureLoader")
	vp.fontLoader = vp.ThreeJsNew("FontLoader")
	vp.initGrid()
	vp.initBackground()
	vp.initTitle()
	return vp
}

func (vp *Viewport) Hide() {
	vp.Canvas.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.Canvas.Get("style").Set("display", "initial")
}

func (vp *Viewport) ResizeCanvas(title bool) {
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()
	if title {
		winH /= 3
	}
	vp.Canvas.Call("setAttribute", "width", winW)
	vp.Canvas.Call("setAttribute", "height", winH)
	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.camera.Set("aspect", float64(winW)/float64(winH))
	vp.camera.Call("updateProjectionMatrix")

	vp.Canvas.Call("setAttribute", "width", vp.ViewWidth)
	vp.Canvas.Call("setAttribute", "height", vp.ViewHeight)
	vp.renderer.Call("setSize", vp.ViewWidth, vp.ViewHeight)
}

func (vp *Viewport) Focus() {
	vp.Canvas.Call("focus")
}

func (vp *Viewport) Zoom(state int) {
}

func (vp *Viewport) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.Canvas.Call("addEventListener", evt, js.FuncOf(fn))
}

func (vp *Viewport) Draw() {
	vp.renderer.Call("render", vp.scene, vp.camera)
}

func (vp *Viewport) ThreeJsNew(name string, args ...interface{}) js.Value {
	return vp.threejs.Get(name).New(args...)
}
