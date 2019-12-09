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

package vector3f

import "math/rand"

type Cube struct {
	Min, Max Vector3f
}

// make normalized cube , if not need use Cube{Min: , Max:}
func NewCube(v1 Vector3f, v2 Vector3f) Cube {
	rtn := Cube{
		Min: Vector3f{},
		Max: Vector3f{},
	}
	for i := 0; i < 3; i++ {
		if v1[i] > v2[i] {
			rtn.Max[i] = v1[i]
			rtn.Min[i] = v2[i]
		} else {
			rtn.Max[i] = v2[i]
			rtn.Min[i] = v1[i]
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

func (h Cube) MakeCubeBy8Driect(center Vector3f, direct8 int) Cube {
	rtn := Vector3f{}
	for i := 0; i < 3; i++ {
		if direct8&(1<<uint(i)) != 0 {
			rtn[i] = h.Min[i]
		} else {
			rtn[i] = h.Max[i]
		}
	}
	return NewCube(center, rtn)
}

func (h Cube) Center() Vector3f {
	return h.Min.Add(h.Max).Idiv(2)
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
	hs := h.SizeVector().Imul(i / 2)
	hc := h.Center()
	return Cube{
		Min: hc.Sub(hs),
		Max: hc.Add(hs),
	}
}

func (h1 Cube) IsOverlap(h2 Cube) bool {
	return !((h1.Min[0] > h2.Max[0] || h1.Max[0] < h2.Min[0]) ||
		(h1.Min[1] > h2.Max[1] || h1.Max[1] < h2.Min[1]) ||
		(h1.Min[2] > h2.Max[2] || h1.Max[2] < h2.Min[2]))
}

func (h1 Cube) IsIn(h2 Cube) bool {
	for i := 0; i < 3; i++ {
		if h1.Min[i] < h2.Min[i] || h1.Max[i] > h2.Max[i] {
			return false
		}
	}
	return true
}
