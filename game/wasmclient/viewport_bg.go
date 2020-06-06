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

	"github.com/kasworld/gowasm3dgame/config/gameconst"
)

func (vp *Viewport) initGrid() {
	outerStageSize := gameconst.StageSize + gameconst.MaxRadius*2
	innerStageSize := gameconst.StageSize

	center := vp.ThreeJsNew("Vector3",
		gameconst.StageSize/2,
		gameconst.StageSize/2,
		gameconst.StageSize/2,
	)

	helper := vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0x0000ff, 0x404040)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", -gameconst.MaxRadius)
	helper.Get("position").Set("z", innerStageSize/2)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", center)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0xffff00, 0x404040)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", gameconst.StageSize+gameconst.MaxRadius)
	helper.Get("position").Set("z", innerStageSize/2)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", center)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0xff0000, 0x404040)
	helper.Get("position").Set("x", -gameconst.MaxRadius)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", innerStageSize/2)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", center)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0x00ffff, 0x404040)
	helper.Get("position").Set("x", gameconst.StageSize+gameconst.MaxRadius)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", innerStageSize/2)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", center)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0x00ff00, 0x404040)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", -gameconst.MaxRadius)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", center)
	vp.scene.Call("add", helper)

	helper = vp.ThreeJsNew("GridHelper", outerStageSize, 10, 0xff00ff, 0x404040)
	helper.Get("position").Set("x", innerStageSize/2)
	helper.Get("position").Set("y", innerStageSize/2)
	helper.Get("position").Set("z", gameconst.StageSize+gameconst.MaxRadius)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", center)
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

func (vp *Viewport) initBackground() {
	bgMap := vp.textureLoader.Call("load", "/resource/background.png")
	bgMap.Set("wrapS", vp.threejs.Get("RepeatWrapping"))
	bgMap.Set("wrapT", vp.threejs.Get("RepeatWrapping"))
	bgMap.Get("repeat").Set("x", 25)
	bgMap.Get("repeat").Set("y", 25)
	// var groundTexture = loader.load( 'textures/terrain/grasslight-big.jpg' );
	// groundTexture.wrapS = groundTexture.wrapT = THREE.RepeatWrapping;
	// groundTexture.repeat.set( 25, 25 );
	// groundTexture.anisotropy = 16;
	// groundTexture.encoding = THREE.sRGBEncoding;
	// var groundMaterial = new THREE.MeshLambertMaterial( { map: groundTexture } );
	// var mesh = new THREE.Mesh( new THREE.PlaneBufferGeometry( 20000, 20000 ), groundMaterial );
	// mesh.position.y = - 250;
	// mesh.rotation.x = - Math.PI / 2;
	// mesh.receiveShadow = true;
	// scene.add( mesh );
	bgMaterial := vp.ThreeJsNew("MeshBasicMaterial",
		map[string]interface{}{
			"map": bgMap,
		},
	)
	bgGeo := vp.ThreeJsNew("PlaneBufferGeometry",
		gameconst.StageSize*25, gameconst.StageSize*25)
	vp.background = vp.ThreeJsNew("Mesh", bgGeo, bgMaterial)
	// jslog.Info(vp.background)
	// vp.background = vp.ThreeJsNew("Sprite", bgMaterial)
	// vp.background.Get("scale").Set("x", gameconst.StageSize)
	// vp.background.Get("scale").Set("y", gameconst.StageSize)
	// vp.background.Get("scale").Set("z", 1)
	vp.background.Get("position").Set("x", gameconst.StageSize/2)
	vp.background.Get("position").Set("y", gameconst.StageSize/2)
	vp.background.Get("position").Set("z", -gameconst.MaxRadius)
	vp.scene.Call("add", vp.background)
}
