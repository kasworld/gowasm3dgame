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

package vector3f

import "math/rand"

type Cube struct {
	Min Vector3f
	Max Vector3f
}

// make normalized cube , if not need use Cube{Min: , Max:}
func NewCube(v1 Vector3f, v2 Vector3f) Cube {
	rtn := Cube{
		Min: v1,
		Max: v2,
	}
	for i := 0; i < 3; i++ {
		if rtn.Min[i] > rtn.Max[i] {
			rtn.Max[i], rtn.Min[i] = rtn.Min[i], rtn.Max[i]
		}
	}
	return rtn
}

func NewCubeByCR(c Vector3f, r float64) Cube {
	return Cube{
		Vector3f{c[0] - r, c[1] - r, c[2] - r},
		Vector3f{c[0] + r, c[1] + r, c[2] + r},
	}
}

func (h Cube) Center() Vector3f {
	return h.Min.Add(h.Max).DivF(2)
}

func (h Cube) DiagLen() float64 {
	return h.Min.LenTo(h.Max)
}

func (h Cube) SizeVector() Vector3f {
	return h.Max.Sub(h.Min)
}

func (h Cube) IsContact(c Vector3f, r float64) bool {
	hc := h.Center()
	hl := h.DiagLen()
	return hl/2+r >= hc.LenTo(c)
}

func (h Cube) RandVector() Vector3f {
	return Vector3f{
		rand.Float64()*(h.Max[0]-h.Min[0]) + h.Min[0],
		rand.Float64()*(h.Max[1]-h.Min[1]) + h.Min[1],
		rand.Float64()*(h.Max[2]-h.Min[2]) + h.Min[2],
	}
}

func (h Cube) Move(v Vector3f) Cube {
	return Cube{
		Min: h.Min.Add(v),
		Max: h.Max.Add(v),
	}
}

func (h Cube) IMul(i float64) Cube {
	hs := h.SizeVector().MulF(i / 2)
	hc := h.Center()
	return Cube{
		Min: hc.Sub(hs),
		Max: hc.Add(hs),
	}
}

func (h1 Cube) IsOverlap(h2 Cube) bool {
	for i := 0; i < 3; i++ {
		if h1.Max[i] < h2.Min[i] || h2.Max[i] < h1.Min[i] {
			return false
		}
	}
	return true
}

func (h1 Cube) IsIn(h2 Cube) bool {
	for i := 0; i < 3; i++ {
		if h1.Min[i] < h2.Min[i] || h1.Max[i] > h2.Max[i] {
			return false
		}
	}
	return true
}
