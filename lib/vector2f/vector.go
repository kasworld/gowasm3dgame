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

package vector2f

import "math"

type Vector2f [2]float64

var VtZero = Vector2f{0, 0}

func NewVectorLenAngle(l, a float64) Vector2f {
	return Vector2f{
		l * math.Cos(a),
		l * math.Sin(a),
	}
}

func (vt Vector2f) Abs() float64 {
	return math.Sqrt(vt[0]*vt[0] + vt[1]*vt[1])
}

func (vt Vector2f) Add(v2 Vector2f) Vector2f {
	return Vector2f{
		vt[0] + v2[0],
		vt[1] + v2[1],
	}
}

func (vt Vector2f) Sub(v2 Vector2f) Vector2f {
	return Vector2f{
		vt[0] - v2[0],
		vt[1] - v2[1],
	}
}

func (vt Vector2f) MulF(f float64) Vector2f {
	return Vector2f{
		vt[0] * f,
		vt[1] * f,
	}
}

func (vt Vector2f) DivF(f float64) Vector2f {
	return Vector2f{
		vt[0] / f,
		vt[1] / f,
	}
}

func (vt Vector2f) Normalize() Vector2f {
	return vt.DivF(vt.Abs())
}

func (vt Vector2f) Neg() Vector2f {
	return Vector2f{
		-vt[0],
		-vt[1],
	}
}

func (vt Vector2f) NegX() Vector2f {
	return Vector2f{
		-vt[0],
		vt[1],
	}
}

func (vt Vector2f) NegY() Vector2f {
	return Vector2f{
		vt[0],
		-vt[1],
	}
}

func (vt Vector2f) LenTo(v2 Vector2f) float64 {
	return v2.Sub(vt).Abs()
}

func (vt Vector2f) Phase() float64 {
	return math.Atan2(vt[1], vt[0])
}

func (vt Vector2f) AddAngle(angle float64) Vector2f {
	return NewVectorLenAngle(vt.Abs(), vt.Phase()+angle)
}

func (vt Vector2f) RotateBy(center Vector2f, angle float64) Vector2f {
	return vt.Sub(center).AddAngle(angle).Add(center)
}

func (vt Vector2f) Dot(v2 Vector2f) float64 {
	return vt[0]*v2[0] + vt[1]*v2[1]
}

func (vt Vector2f) Cross() Vector2f {
	return Vector2f{
		vt[1],
		-vt[0],
	}
}

func (vt Vector2f) IsIn(rt Rect) bool {
	return rt.Min[0] <= vt[0] && vt[0] <= rt.Max[0] &&
		rt.Min[1] <= vt[1] && vt[1] <= rt.Max[1]
}
