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
	"github.com/kasworld/gowasmlib/jslog"
)

func (vp *Viewport) setTitleCamera() {
	// set title camera pos
	vp.camera.Get("position").Set("x", gameconst.StageSize/2)
	vp.camera.Get("position").Set("y", gameconst.StageSize/2)
	vp.camera.Get("position").Set("z", gameconst.StageSize)
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
	vp.lightTitle.Get("position").Set("x", gameconst.StageSize)
	vp.lightTitle.Get("position").Set("y", gameconst.StageSize)
	vp.lightTitle.Get("position").Set("z", gameconst.StageSize)
	vp.scene.Call("add", vp.lightTitle)
	vp.setTitleCamera()
	vp.fontLoader.Call("load", "/fonts/helvetiker_regular.typeface.json",
		js.FuncOf(vp.fontLoaded),
	)
}

func (vp *Viewport) fontLoaded(this js.Value, args []js.Value) interface{} {
	vp.fontTitle = args[0]
	jslog.Info(vp.fontLoader, vp.fontTitle)
	str := "gowasm3dgame"
	ftGeo := vp.ThreeJsNew("TextGeometry", str, map[string]interface{}{
		"font":           vp.fontTitle,
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
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	co := rnd.Uint32() & 0x00ffffff
	ftMat := vp.ThreeJsNew("MeshPhongMaterial",
		map[string]interface{}{
			"color": co,
			// "flatShading": true,
		},
	)
	vp.jsoTitle = vp.ThreeJsNew("Mesh", ftGeo, ftMat)
	vp.jsoTitle.Get("position").Set("x", gameconst.StageSize/2-(geoMax-geoMin)/2)
	vp.jsoTitle.Get("position").Set("y", gameconst.StageSize/2)
	vp.jsoTitle.Get("position").Set("z", gameconst.StageSize/2)
	vp.scene.Call("add", vp.jsoTitle)

	return nil
}
