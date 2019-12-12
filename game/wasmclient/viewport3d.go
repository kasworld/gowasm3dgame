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

	"github.com/kasworld/gowasm3dgame/lib/vector3f"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

type Viewport3d struct {
	neecRecalc bool
	ViewWidth  int
	ViewHeight int

	stageInfo *w3d_obj.NotiStageInfo_data

	canvas   js.Value
	threejs  js.Value
	scene    js.Value
	camera   js.Value
	renderer js.Value
	light    js.Value

	cube js.Value
}

func NewViewport3d(cnvid string) *Viewport3d {
	vp := &Viewport3d{}
	vp.threejs = js.Global().Get("THREE")
	vp.scene = vp.ThreeJsNew("Scene")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.canvas = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvas3d").Call("appendChild", vp.canvas)

	geometry := vp.ThreeJsNew("BoxGeometry", 1, 1, 1)
	material := vp.ThreeJsNew("MeshBasicMaterial")
	vp.cube = vp.ThreeJsNew("Mesh", geometry, material)
	vp.scene.Call("add", vp.cube)
	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 45, 1, 1, 10000)

	vp.initGrid()
	vp.setCamera(vector3f.Vector3f{1000, 1000, 1000}, vector3f.Vector3f{0, 0, 0})
	vp.initLight()
	// js.Global().Get("console").Call("debug", vp.camera)
	return vp
}

func (vp *Viewport3d) initGrid() {
	helper := vp.ThreeJsNew("GridHelper", 1000, 100, 0x0000ff, 0x404040)
	// helper.Call("setColors", 0x0000ff, 0x404040)
	helper.Get("position").Set("y", -1000)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", 1000, 100, 0x0000ff, 0x404040)
	// helper.Call("setColors", 0x0000ff, 0x404040)
	helper.Get("position").Set("y", 1000)
	vp.scene.Call("add", helper)

	axisHelper := vp.ThreeJsNew("AxesHelper", 1000)
	vp.scene.Call("add", axisHelper)
}
func (vp *Viewport3d) initLight() {
	vp.light = vp.ThreeJsNew("PointLight", 0x808080, 1)
	vp.scene.Call("add", vp.light)
}

func (vp *Viewport3d) setCamera(vt1, vt2 vector3f.Vector3f) {
	JsSetPos(vp.camera, vt1)
	vp.camera.Call("lookAt", vp.Vt3fToThVt3(vt2))
	vp.camera.Call("updateProjectionMatrix")
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
	winH = winH * 2 / 3

	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.canvas.Call("setAttribute", "width", vp.ViewWidth)
	vp.canvas.Call("setAttribute", "height", vp.ViewHeight)

	vp.renderer.Call("setSize", vp.ViewWidth, vp.ViewHeight)
}

func (vp *Viewport3d) Draw(tick int64) {
	vp.calcViewCellValue()

	rot := vp.cube.Get("rotation")
	rot.Set("x", rot.Get("x").Float()+0.01)
	rot.Set("y", rot.Get("y").Float()+0.01)
	vp.renderer.Call("render", vp.scene, vp.camera)
}
