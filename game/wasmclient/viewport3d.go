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

	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/game/gameconst"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/htmlcolors"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

type Viewport3d struct {
	neecRecalc bool
	ViewWidth  int
	ViewHeight int

	canvas   js.Value
	threejs  js.Value
	scene    js.Value
	camera   js.Value
	renderer js.Value
	light    js.Value

	jsSceneObjs map[string]js.Value
}

func NewViewport3d(cnvid string) *Viewport3d {
	vp := &Viewport3d{}
	vp.jsSceneObjs = make(map[string]js.Value)
	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.canvas = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvas3d").Call("appendChild", vp.canvas)

	vp.scene = vp.ThreeJsNew("Scene")

	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 45, 1, 1,
		gameconst.StageSize*10)

	vp.initGrid()
	vp.setCamera(vector3f.Vector3f{
		gameconst.StageSize * 3 / 2,
		gameconst.StageSize * 3 / 2,
		gameconst.StageSize * 3 / 2}, vector3f.Vector3f{0, 0, 0})
	vp.initLight()

	// vp.cube = vp.newGLObj(30, htmlcolors.Red)
	// vp.scene.Call("add", vp.cube)
	return vp
}

func (vp *Viewport3d) initGrid() {
	helper := vp.ThreeJsNew("GridHelper", gameconst.StageSize, 100, 0x0000ff, 0x404040)
	helper.Get("position").Set("y", -gameconst.StageSize)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", gameconst.StageSize, 100, 0x0000ff, 0x404040)
	helper.Get("position").Set("y", gameconst.StageSize)
	vp.scene.Call("add", helper)

	axisHelper := vp.ThreeJsNew("AxesHelper", gameconst.StageSize)
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

func (vp *Viewport3d) calcResize() {
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
	vp.calcResize()

	vp.renderer.Call("render", vp.scene, vp.camera)
}

func (vp *Viewport3d) add2Scene(o *w3d_obj.GameObj, co htmlcolors.Color24) js.Value {
	if jso, exist := vp.jsSceneObjs[o.UUID]; exist {
		JsSetPos(jso, o.PosVt)
		return jso
	}
	radius := gameobjtype.Attrib[o.GOType].Radius
	geometry := vp.ThreeJsNew("SphereGeometry", radius, 32, 16)
	material := vp.ThreeJsNew("MeshPhongMaterial")
	material.Set("color", vp.ToThColor(htmlcolors.Gray))
	material.Set("emissive", vp.ToThColor(co))
	material.Set("shininess", 30)
	jso := vp.ThreeJsNew("Mesh", geometry, material)
	JsSetPos(jso, o.PosVt)
	vp.scene.Call("add", jso)
	vp.jsSceneObjs[o.UUID] = jso
	return jso
}

func (vp *Viewport3d) processRecvStageInfo(stageInfo *w3d_obj.NotiStageInfo_data) {
	addUUID := make(map[string]bool)
	for _, tm := range stageInfo.Teams {
		if tm == nil {
			continue
		}
		vp.add2Scene(tm.Ball, tm.Color24)
		addUUID[tm.Ball.UUID] = true
		vp.add2Scene(tm.HomeMark, tm.Color24)
		addUUID[tm.HomeMark.UUID] = true
		for _, v := range tm.Objs {
			if v == nil {
				continue
			}
			vp.add2Scene(v, tm.Color24)
			addUUID[v.UUID] = true
		}
	}
	for id, jso := range vp.jsSceneObjs {
		if !addUUID[id] {
			vp.scene.Call("remove", jso)
			delete(vp.jsSceneObjs, id)
		}
	}
}
