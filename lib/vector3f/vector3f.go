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

import (
	"fmt"
	"math"
)

type Vector3f [3]float64

// func RandVector3D(st, end float64) Vector3f {
// 	return Vector3f{
// 		rand.Float64()*(end-st) + st,
// 		rand.Float64()*(end-st) + st,
// 		rand.Float64()*(end-st) + st,
// 	}
// }

// func RandVector(st, end Vector3f) Vector3f {
// 	return Vector3f{
// 		rand.Float64()*(end[0]-st[0]) + st[0],
// 		rand.Float64()*(end[1]-st[1]) + st[1],
// 		rand.Float64()*(end[2]-st[2]) + st[2],
// 	}
// }

func (v Vector3f) String() string {
	return fmt.Sprintf("[%5.2f,%5.2f,%5.2f]", v[0], v[1], v[2])
}

var VtZero = Vector3f{0, 0, 0}
var VtUnitX = Vector3f{1, 0, 0}
var VtUnitY = Vector3f{0, 1, 0}
var VtUnitZ = Vector3f{0, 0, 1}

func (p Vector3f) Eq(other Vector3f) bool {
	return p == other
	//return p[0] == other[0] && p[1] == other[1] && p[2] == other[2]
}
func (p Vector3f) Ne(other Vector3f) bool {
	return !p.Eq(other)
}
func (p Vector3f) IsZero() bool {
	return p.Eq(VtZero)
}
func (p Vector3f) Add(other Vector3f) Vector3f {
	return Vector3f{p[0] + other[0], p[1] + other[1], p[2] + other[2]}
}
func (p Vector3f) Neg() Vector3f {
	return Vector3f{-p[0], -p[1], -p[2]}
}
func (p Vector3f) Sub(other Vector3f) Vector3f {
	return Vector3f{p[0] - other[0], p[1] - other[1], p[2] - other[2]}
}
func (p Vector3f) Mul(other Vector3f) Vector3f {
	return Vector3f{p[0] * other[0], p[1] * other[1], p[2] * other[2]}
}
func (p Vector3f) Imul(other float64) Vector3f {
	return Vector3f{p[0] * other, p[1] * other, p[2] * other}
}
func (p Vector3f) Idiv(other float64) Vector3f {
	return Vector3f{p[0] / other, p[1] / other, p[2] / other}
}
func (p Vector3f) Abs() float64 {
	return math.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])
}
func (p Vector3f) Sqd(q Vector3f) float64 {
	var sum float64
	for dim, pCoord := range p {
		d := pCoord - q[dim]
		sum += d * d
	}
	return sum
}

func (p Vector3f) LenTo(other Vector3f) float64 {
	return math.Sqrt(p.Sqd(other))
}

func (p Vector3f) Normalized() Vector3f {
	d := p.Abs()
	if d > 0 {
		return p.Idiv(d)
	}
	return p
}
func (p Vector3f) NormalizedTo(l float64) Vector3f {
	d := p.Abs() / l
	if d != 0 {
		return p.Idiv(d)
	}
	return p
}
func (p Vector3f) Dot(other Vector3f) float64 {
	return p[0]*other[0] + p[1]*other[1] + p[2]*other[2]
}
func (p Vector3f) Cross(other Vector3f) Vector3f {
	return Vector3f{
		p[1]*other[2] - p[2]*other[1],
		-p[0]*other[2] + p[2]*other[0],
		p[0]*other[1] - p[1]*other[0],
	}
}

// reflect plane( == normal vector )
func (p Vector3f) Reflect(normal Vector3f) Vector3f {
	d := 2 * (p[0]*normal[0] + p[1]*normal[1] + p[2]*normal[2])
	return Vector3f{p[0] - d*normal[0], p[1] - d*normal[1], p[2] - d*normal[2]}
}
func (p Vector3f) RotateAround(axis Vector3f, theta float64) Vector3f {
	// Return the vector rotated around axis through angle theta. Right hand rule applies
	// Adapted from equations published by Glenn Murray.
	// http://inside.mines.edu/~gmurray/ArbitraryAxisRotation/ArbitraryAxisRotation.html
	x, y, z := p[0], p[1], p[2]
	u, v, w := axis[0], axis[1], axis[2]

	// Extracted common factors for simplicity and efficiency
	r2 := u*u + v*v + w*w
	r := math.Sqrt(r2)
	ct := math.Cos(theta)
	st := math.Sin(theta) / r
	dt := (u*x + v*y + w*z) * (1 - ct) / r2
	return Vector3f{
		(u*dt + x*ct + (-w*y+v*z)*st),
		(v*dt + y*ct + (w*x-u*z)*st),
		(w*dt + z*ct + (-v*x+u*y)*st),
	}
}
func (p Vector3f) Angle(other Vector3f) float64 {
	// Return the angle to the vector other
	return math.Acos(p.Dot(other) / (p.Abs() * other.Abs()))
}
func (p Vector3f) Project(other Vector3f) Vector3f {
	// Return one vector projected on the vector other
	n := other.Normalized()
	return n.Imul(p.Dot(n))
}

// for aim ahead target with projectile
// return time dur
func (srcpos Vector3f) CalcAimAheadDur(
	dstpos Vector3f, dstmv Vector3f, bulletspeed float64) float64 {
	totargetvt := dstpos.Sub(srcpos)
	a := dstmv.Dot(dstmv) - bulletspeed*bulletspeed
	b := 2 * dstmv.Dot(totargetvt)
	c := totargetvt.Dot(totargetvt)
	p := -b / (2 * a)
	q := math.Sqrt((b*b)-4*a*c) / (2 * a)
	t1 := p - q
	t2 := p + q

	var rtn float64
	if t1 > t2 && t2 > 0 {
		rtn = t2
	} else {
		rtn = t1
	}
	if rtn < 0 || math.IsNaN(rtn) {
		return math.Inf(1)
	}
	return rtn
}

func (center Vector3f) To8Direct(v2 Vector3f) int {
	rtn := 0
	for i := 0; i < 3; i++ {
		if center[i] > v2[i] {
			rtn += 1 << uint(i)
		}
	}
	return rtn
}

func (p Vector3f) IsIn(hr Cube) bool {
	return hr.Min[0] <= p[0] && p[0] <= hr.Max[0] &&
		hr.Min[1] <= p[1] && p[1] <= hr.Max[1] &&
		hr.Min[2] <= p[2] && p[2] <= hr.Max[2]
}

func (p Vector3f) MakeIn(hr Cube) (Vector3f, int) {
	rtn := p
	changed := 0
	var i uint
	for i = 0; i < 3; i++ {
		if rtn[i] > hr.Max[i] {
			rtn[i] = hr.Max[i]
			changed |= 1 << (i*2 + 1)
		}
		if rtn[i] < hr.Min[i] {
			rtn[i] = hr.Min[i]
			changed |= 1 << (i * 2)
		}
	}
	return rtn, changed
}

// for serialize
func (v Vector3f) ToInt32Vector() [3]int32 {
	return [3]int32{int32(v[0]), int32(v[1]), int32(v[2])}
}

func NewByInt32Vector(s [3]int32) Vector3f {
	return Vector3f{float64(s[0]), float64(s[1]), float64(s[2])}
}
