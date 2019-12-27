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
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/htmlcolors"
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

	jsSceneObjs   map[string]js.Value
	geometryCache map[gameobjtype.GameObjType]js.Value
	materialCache map[htmlcolors.Color24]js.Value
}

func NewViewport3d(cnvid string) *Viewport3d {
	vp := &Viewport3d{
		jsSceneObjs:   make(map[string]js.Value),
		geometryCache: make(map[gameobjtype.GameObjType]js.Value),
		materialCache: make(map[htmlcolors.Color24]js.Value),
	}

	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.canvas = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvas3d").Call("appendChild", vp.canvas)

	vp.scene = vp.ThreeJsNew("Scene")

	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 75, 1, gameobjtype.MaxRadius,
		gameconst.StageSize*10)

	vp.initGrid()
	vp.initLight()
	return vp
}

func (vp *Viewport3d) initGrid() {
	helper := vp.ThreeJsNew("GridHelper",
		gameconst.StageSize, 100, 0x0000ff, 0x404040)

	helper.Get("position").Set("x", gameconst.StageSize/2)
	helper.Get("position").Set("y", 0)
	helper.Get("position").Set("z", gameconst.StageSize/2)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper",
		gameconst.StageSize, 100, 0x00ff00, 0x404040)
	helper.Get("position").Set("x", gameconst.StageSize/2)
	helper.Get("position").Set("y", 0)
	helper.Get("position").Set("z", gameconst.StageSize/2)
	vp.scene.Call("add", helper)

	box3 := vp.ThreeJsNew("Box3",
		vp.ThreeJsNew("Vector3",
			0-gameobjtype.MaxRadius,
			0-gameobjtype.MaxRadius,
			0-gameobjtype.MaxRadius,
		),
		vp.ThreeJsNew("Vector3",
			gameconst.StageSize+gameobjtype.MaxRadius,
			gameconst.StageSize+gameobjtype.MaxRadius,
			gameconst.StageSize+gameobjtype.MaxRadius,
		),
	)
	helper = vp.ThreeJsNew("Box3Helper", box3, 0xffffff)
	vp.scene.Call("add", helper)

	axisHelper := vp.ThreeJsNew("AxesHelper", gameconst.StageSize)
	vp.scene.Call("add", axisHelper)
}
func (vp *Viewport3d) initLight() {
	vp.light = vp.ThreeJsNew("PointLight", 0x808080, 1)
	vp.scene.Call("add", vp.light)
	// vp.light.Get("position").Set("x", vt[0])
	// vp.light.Get("position").Set("y", vt[1])
	// vp.light.Get("position").Set("z", vt[2])
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
	size := winW
	if size > winH {
		size = winH
	}
	// size -= 20
	vp.ViewWidth = size
	vp.ViewHeight = size

	vp.canvas.Call("setAttribute", "width", vp.ViewWidth)
	vp.canvas.Call("setAttribute", "height", vp.ViewHeight)

	vp.renderer.Call("setSize", vp.ViewWidth, vp.ViewHeight)
}

func (vp *Viewport3d) Draw(tick int64) {
	vp.calcResize()

	vp.renderer.Call("render", vp.scene, vp.camera)
}

func (vp *Viewport3d) getGeometry(gotype gameobjtype.GameObjType) js.Value {
	geo, exist := vp.geometryCache[gotype]
	if !exist {
		radius := gameobjtype.Attrib[gotype].Radius
		switch gotype {
		default:
			geo = vp.ThreeJsNew("SphereGeometry", radius, 32, 16)
		case gameobjtype.Ball:
			geo = vp.ThreeJsNew("TorusGeometry", radius, radius/2, 16, 64)
		case gameobjtype.Shield:
			geo = vp.ThreeJsNew("IcosahedronGeometry", radius)
		case gameobjtype.HommingShield:
			geo = vp.ThreeJsNew("OctahedronGeometry", radius)
			// geo = vp.ThreeJsNew("IcosahedronGeometry", radius)
		case gameobjtype.Bullet:
			geo = vp.ThreeJsNew("DodecahedronGeometry", radius)
		case gameobjtype.HommingBullet:
			geo = vp.ThreeJsNew("TetrahedronGeometry", radius)
			//geo = vp.ThreeJsNew("OctahedronGeometry", radius)
		case gameobjtype.SuperBullet:
			geo = vp.ThreeJsNew("ConeGeometry", radius, radius*2, 16)
			// geo = vp.ThreeJsNew("TetrahedronGeometry", radius)
		case gameobjtype.BurstBullet:
			geo = vp.ThreeJsNew("DodecahedronGeometry", radius)
		case gameobjtype.HomeMark:
			geo = vp.ThreeJsNew("BoxGeometry", radius*2, radius*2, radius*2)
		case gameobjtype.Deco:
			geo = vp.ThreeJsNew("SphereGeometry", radius, 32, 16)
		case gameobjtype.Hard:
			geo = vp.ThreeJsNew("SphereGeometry", radius, 32, 16)
		case gameobjtype.Food:
			geo = vp.ThreeJsNew("SphereGeometry", radius, 32, 16)
		}
		vp.geometryCache[gotype] = geo
	}
	return geo
}

func (vp *Viewport3d) getMaterial(co htmlcolors.Color24) js.Value {
	mat, exist := vp.materialCache[co]
	if !exist {
		mat = vp.ThreeJsNew("MeshStandardMaterial")
		// material.Set("color", vp.ToThColor(htmlcolors.Gray))
		mat.Set("emissive", vp.ToThColor(co))
		mat.Set("shininess", 30)
		vp.materialCache[co] = mat
	}
	return mat
}

func (vp *Viewport3d) add2Scene(o *w3d_obj.GameObj, co htmlcolors.Color24) js.Value {
	if jso, exist := vp.jsSceneObjs[o.UUID]; exist {
		jso.Get("position").Set("x", o.PosVt[0])
		jso.Get("position").Set("y", o.PosVt[1])
		jso.Get("position").Set("z", o.PosVt[2])
		jso.Get("rotation").Set("x", o.RotVt[0])
		jso.Get("rotation").Set("y", o.RotVt[1])
		jso.Get("rotation").Set("z", o.RotVt[2])
		return jso
	}
	geometry := vp.getGeometry(o.GOType)
	material := vp.getMaterial(co)
	jso := vp.ThreeJsNew("Mesh", geometry, material)
	jso.Get("position").Set("x", o.PosVt[0])
	jso.Get("position").Set("y", o.PosVt[1])
	jso.Get("position").Set("z", o.PosVt[2])
	jso.Get("rotation").Set("x", o.RotVt[0])
	jso.Get("rotation").Set("y", o.RotVt[1])
	jso.Get("rotation").Set("z", o.RotVt[2])
	vp.scene.Call("add", jso)
	vp.jsSceneObjs[o.UUID] = jso
	return jso
}

func (vp *Viewport3d) processRecvStageInfo(stageInfo *w3d_obj.NotiStageInfo_data) {
	setCamera := false
	addUUID := make(map[string]bool)
	for _, tm := range stageInfo.Teams {
		if tm == nil {
			continue
		}
		if !setCamera {
			setCamera = true

			vt1 := tm.Objs[0].PosVt
			vp.camera.Get("position").Set("x", vt1[0])
			vp.camera.Get("position").Set("y", vt1[1])
			vp.camera.Get("position").Set("z", vt1[2])

			vt2 := tm.Objs[1].PosVt
			vp.camera.Call("lookAt",
				vp.ThreeJsNew("Vector3",
					vt2[0], vt2[1], vt2[2],
				),
			)
			vp.camera.Call("updateProjectionMatrix")
		}
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
