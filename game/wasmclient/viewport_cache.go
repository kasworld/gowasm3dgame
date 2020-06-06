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

	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

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

func (vp *Viewport) getColorMaterial(co uint32) js.Value {
	mat, exist := vp.materialCache[co]
	if !exist {
		mat = vp.ThreeJsNew("MeshPhongMaterial",
			map[string]interface{}{
				"color": co,
				// "flatShading": true,
			},
		)
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
