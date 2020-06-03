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
	"math"
	"syscall/js"

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

func (vp *Viewport) ThreeJsNew(name string, args ...interface{}) js.Value {
	return vp.threejs.Get(name).New(args...)
}

func (vp *Viewport) initGrid() {
	outerStageSize := gameconst.StageSize + gameconst.MaxRadius*2
	innerStageSize := gameconst.StageSize

	helper := vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0x0000ff, 0x404040)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", -gameconst.MaxRadius)
	helper.Get("position").Set("z", innerStageSize/2)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0xffff00, 0x404040)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", gameconst.StageSize+gameconst.MaxRadius)
	helper.Get("position").Set("z", innerStageSize/2)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0xff0000, 0x404040)
	helper.Get("rotation").Set("z", math.Pi/2)
	helper.Get("position").Set("x", -gameconst.MaxRadius)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", innerStageSize/2)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0x00ffff, 0x404040)
	helper.Get("rotation").Set("z", math.Pi/2)
	helper.Get("position").Set("x", gameconst.StageSize+gameconst.MaxRadius)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", innerStageSize/2)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0x00ff00, 0x404040)
	helper.Get("rotation").Set("x", math.Pi/2)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", -gameconst.MaxRadius)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0xff00ff, 0x404040)
	helper.Get("rotation").Set("x", math.Pi/2)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", gameconst.StageSize+gameconst.MaxRadius)
	vp.scene.Call("add", helper)

	box3 := vp.ThreeJsNew("Box3",
		vp.ThreeJsNew("Vector3", 0, 0, 0),
		vp.ThreeJsNew("Vector3", innerStageSize, innerStageSize, innerStageSize),
	)
	helper = vp.ThreeJsNew("Box3Helper", box3, 0xffffff)
	vp.scene.Call("add", helper)

	// axisHelper := vp.ThreeJsNew("AxesHelper", gameconst.StageSize)
	// vp.scene.Call("add", axisHelper)
}

func (vp *Viewport) getGeometry(gotype gameobjtype.GameObjType) js.Value {
	geo, exist := vp.geometryCache[gotype]
	if !exist {
		radius := vp.Type2Radius[gotype]
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

func (vp *Viewport) getMaterial(co uint32) js.Value {
	mat, exist := vp.materialCache[co]
	if !exist {
		mat = vp.ThreeJsNew("MeshStandardMaterial")
		// material.Set("color", vp.ToThColor(htmlcolors.Gray))
		mat.Set("emissive", vp.ThreeJsNew("Color", co))
		mat.Set("shininess", 30)
		vp.materialCache[co] = mat
	}
	return mat
}

func (vp *Viewport) getLightNHelper(o *w3d_obj.Light) *LightNHelper {
	lt, exist := vp.lightCache[o.UUID]
	if !exist {
		lt = &LightNHelper{}
		lt.Light = vp.ThreeJsNew("PointLight", o.Color, 1)
		lt.Helper = vp.ThreeJsNew("PointLightHelper", lt.Light, 1)
		vp.lightCache[o.UUID] = lt
	}
	return lt
}

func (vp *Viewport) addLight2Scene(o *w3d_obj.Light) {
	if jso, exist := vp.jsSceneObjs[o.UUID]; exist {
		jso.Get("position").Set("x", o.PosVt[0])
		jso.Get("position").Set("y", o.PosVt[1])
		jso.Get("position").Set("z", o.PosVt[2])
		return
	}
	lt := vp.getLightNHelper(o)
	lt.Light.Get("position").Set("x", o.PosVt[0])
	lt.Light.Get("position").Set("y", o.PosVt[1])
	lt.Light.Get("position").Set("z", o.PosVt[2])
	vp.scene.Call("add", lt.Light)
	vp.scene.Call("add", lt.Helper)
	vp.jsSceneObjs[o.UUID] = lt.Light
}

func (vp *Viewport) add2Scene(o *w3d_obj.GameObj) js.Value {
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
	material := vp.getMaterial(o.Color24)
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

func (vp *Viewport) processRecvStageInfo(stageInfo *w3d_obj.NotiStageInfo_data) {
	addUUID := make(map[string]bool)
	vt1 := stageInfo.CameraPos
	vp.camera.Get("position").Set("x", vt1[0])
	vp.camera.Get("position").Set("y", vt1[1])
	vp.camera.Get("position").Set("z", vt1[2])

	vt2 := stageInfo.CameraLookAt
	vp.camera.Call("lookAt",
		vp.ThreeJsNew("Vector3",
			vt2[0], vt2[1], vt2[2],
		),
	)
	vp.camera.Call("updateProjectionMatrix")

	for _, o := range stageInfo.Lights {
		vp.addLight2Scene(o)
		addUUID[o.UUID] = true
	}

	for _, o := range stageInfo.ObjList {
		vp.add2Scene(o)
		addUUID[o.UUID] = true
	}
	for id, jso := range vp.jsSceneObjs {
		if !addUUID[id] {
			vp.scene.Call("remove", jso)
			delete(vp.jsSceneObjs, id)
			if lt, exist := vp.lightCache[id]; exist {
				vp.scene.Call("remove", lt.Helper)
			}
		}
	}
}
