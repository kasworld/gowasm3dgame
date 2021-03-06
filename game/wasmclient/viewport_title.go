// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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
)

func (vp *Viewport) setTitleCamera() {
	// set title camera pos
	SetPosition(vp.camera,
		gameconst.StageSize/2,
		gameconst.StageSize/2,
		gameconst.StageSize,
	)
	vp.camera.Call("lookAt",
		vp.ThreeJsNew("Vector3",
			gameconst.StageSize/2,
			gameconst.StageSize/2,
			0,
		),
	)
	vp.camera.Call("updateProjectionMatrix")
}

func (vp *Viewport) hideTitle() {
	vp.scene.Call("remove", vp.jsoTitle)
	vp.scene.Call("remove", vp.lightTitle)
}

func (vp *Viewport) initTitle() {
	vp.lightTitle = vp.ThreeJsNew("PointLight", 0xffffff, 1)
	SetPosition(vp.lightTitle,
		gameconst.StageSize,
		gameconst.StageSize,
		gameconst.StageSize,
	)
	vp.scene.Call("add", vp.lightTitle)
	vp.setTitleCamera()
	vp.fontLoader.Call("load", "/fonts/helvetiker_regular.typeface.json",
		js.FuncOf(vp.fontLoaded),
	)
}

func (vp *Viewport) fontLoaded(this js.Value, args []js.Value) interface{} {
	vp.font_helvetiker_regular = args[0]
	str := "gowasm3dgame"

	ftGeo := vp.ThreeJsNew("TextGeometry", str, map[string]interface{}{
		"font":           vp.font_helvetiker_regular,
		"size":           80,
		"height":         5,
		"curveSegments":  12,
		"bevelEnabled":   true,
		"bevelThickness": 10,
		"bevelSize":      8,
		"bevelOffset":    0,
		"bevelSegments":  5,
	})
	ftGeo.Call("computeBoundingBox")
	geoMax := ftGeo.Get("boundingBox").Get("max").Get("x").Float()
	geoMin := ftGeo.Get("boundingBox").Get("min").Get("x").Float()

	co := vp.rnd.Uint32() & 0x00ffffff
	ftMat := vp.getColorMaterial(co)

	vp.jsoTitle = vp.ThreeJsNew("Mesh", ftGeo, ftMat)
	SetPosition(vp.jsoTitle,
		gameconst.StageSize/2-(geoMax-geoMin)/2,
		gameconst.StageSize/2,
		gameconst.StageSize/2,
	)
	vp.scene.Call("add", vp.jsoTitle)

	return nil
}
