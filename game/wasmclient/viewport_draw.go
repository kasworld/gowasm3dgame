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
	"syscall/js"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

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
	material := vp.getColorMaterial(o.Color24)
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

func (vp *Viewport) processRecvStageInfo(
	stageInfo *w3d_obj.NotiStageInfo_data) {

	bgPos := stageInfo.BackgroundPos
	vp.background.Get("position").Set("x", bgPos[0])
	vp.background.Get("position").Set("y", bgPos[1])

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

	addUUID := make(map[string]bool)
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
