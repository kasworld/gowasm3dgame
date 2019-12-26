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
	"fmt"
	"syscall/js"
	"time"

	"github.com/kasworld/htmlcolors"
)

func CalcCurrentFrame(difftick int64, fps float64) int {
	diffsec := float64(difftick) / float64(time.Second)
	frame := fps * diffsec
	return int(frame)
}

func getImgWH(srcImageID string) (js.Value, float64, float64) {
	img := js.Global().Get("document").Call("getElementById", srcImageID)
	if !img.Truthy() {
		fmt.Printf("fail to get %v", srcImageID)
		return js.Null(), 0, 0
	}
	srcw := img.Get("naturalWidth").Float()
	srch := img.Get("naturalHeight").Float()
	return img, srcw, srch
}

func JsSetPos(jsobj js.Value, vt [3]float32) {
	jsobj.Get("position").Set("x", vt[0])
	jsobj.Get("position").Set("y", vt[1])
	jsobj.Get("position").Set("z", vt[2])
}

func JsSetRotation(jsobj js.Value, vt [3]float32) {
	jsobj.Get("rotation").Set("x", vt[0])
	jsobj.Get("rotation").Set("y", vt[1])
	jsobj.Get("rotation").Set("z", vt[2])
	// jsobj.Call("rotateX", vt[0])
	// jsobj.Call("rotateY", vt[1])
	// jsobj.Call("rotateZ", vt[2])
}

func (vp *Viewport3d) ThreeJsNew(name string, args ...interface{}) js.Value {
	return vp.threejs.Get(name).New(args...)
}

func (vp *Viewport3d) ToThColor(co htmlcolors.Color24) js.Value {
	return vp.ThreeJsNew("Color", int(co))
}
