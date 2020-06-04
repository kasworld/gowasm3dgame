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

package vector2f

type Rect struct {
	Min Vector2f
	Max Vector2f

	// X, Y float64
	// W, H float64
}

func NewRect(v1, v2 Vector2f) Rect {
	rtn := Rect{
		Min: v1,
		Max: v2,
	}
	for i := 0; i < 2; i++ {
		if rtn.Min[i] > rtn.Max[i] {
			rtn.Max[i], rtn.Min[i] = rtn.Min[i], rtn.Max[i]
		}
	}
	return rtn
}

func NewRectCenterWH(ctVt, whVt Vector2f) Rect {
	rtn := Rect{
		Min: ctVt,
		Max: ctVt.Add(whVt),
	}
	for i := 0; i < 2; i++ {
		if rtn.Min[i] > rtn.Max[i] {
			rtn.Max[i] = rtn.Min[i]
			rtn.Min[i] = rtn.Max[i]
		}
	}
	return rtn
}

func (rt Rect) Center() Vector2f {
	return rt.Min.Add(rt.Max).DivF(2)
}

func (rt Rect) DiagLen() float64 {
	return rt.Min.LenTo(rt.Max)
}

func (rt Rect) SizeVector() Vector2f {
	return rt.Max.Sub(rt.Min)
}

// func (rt Rect) Enlarge(size Vector2f) Rect {
// 	return Rect{
// 		rt.X - size[0], rt.Y - size[1],
// 		rt.W + size[0]*2, rt.H + size[1]*2,
// 	}
// }
// func (rt Rect) Shrink(size Vector2f) Rect {
// 	return Rect{
// 		rt.X + size[0], rt.Y + size[1],
// 		rt.W - size[0]*2, rt.H - size[1]*2,
// 	}
// }
// func (rt Rect) ShrinkSym(n float64) Rect {
// 	return Rect{rt.X + n, rt.Y + n, rt.W - n*2, rt.H - n*2}
// }

// func (rt Rect) X1() float64 {
// 	return rt.X
// }
// func (rt Rect) X2() float64 {
// 	return rt.X + rt.W
// }
// func (rt Rect) Y1() float64 {
// 	return rt.Y
// }
// func (rt Rect) Y2() float64 {
// 	return rt.Y + rt.H
// }

func (r1 Rect) IsOverlap(r2 Rect) bool {
	return !((r1.Min[0] >= r2.Max[0] || r1.Max[0] <= r2.Min[0]) ||
		(r1.Min[1] >= r2.Max[1] || r1.Max[1] <= r2.Min[1]))
}

func (r1 Rect) IsIn(r2 Rect) bool {
	if r1.Min[0] < r2.Min[0] || r1.Max[0] > r2.Max[0] {
		return false
	}
	if r1.Min[1] < r2.Min[1] || r1.Max[1] > r2.Max[1] {
		return false
	}
	return true
}

func min(i, j float64) float64 {
	if i < j {
		return i
	}
	return j
}
func max(i, j float64) float64 {
	if i > j {
		return i
	}
	return j
}

func (r1 Rect) Union(r2 Rect) Rect {
	rt := Rect{
		Vector2f{
			min(r1.Min[0], r2.Min[0]),
			min(r1.Min[1], r2.Min[1])},
		Vector2f{
			max(r1.Max[0], r2.Max[0]),
			max(r1.Max[1], r2.Max[1]),
		},
	}
	return rt
}

func (r1 Rect) Intersection(r2 Rect) Rect {
	rt := Rect{
		Vector2f{
			max(r1.Min[0], r2.Min[0]),
			max(r1.Min[1], r2.Min[1]),
		},
		Vector2f{
			min(r1.Max[0], r2.Max[0]),
			min(r1.Max[1], r2.Max[1]),
		},
	}
	return rt
}

func (rt Rect) Contain(vt Vector2f) bool {
	return rt.Min[0] <= vt[0] && vt[0] <= rt.Max[0] &&
		rt.Min[1] <= vt[1] && vt[1] <= rt.Max[1]
}

// func (rt Rect) RelPos(x, y float64) (float64, float64) {
// 	return x - rt.X, y - rt.Y
// }

func (rt Rect) WrapVector(vt Vector2f) Vector2f {
	if vt[0] < rt.Min[0] {
		vt[0] = rt.Max[0]
	}
	if vt[1] < rt.Min[1] {
		vt[1] = rt.Max[1]
	}

	if vt[0] > rt.Max[0] {
		vt[0] = rt.Min[0]
	}
	if vt[1] > rt.Max[1] {
		vt[1] = rt.Min[1]
	}
	return vt
}
