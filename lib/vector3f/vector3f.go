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
	"math/rand"
)

type Vector3f [3]float64

func (v Vector3f) String() string {
	return fmt.Sprintf("[%5.2f,%5.2f,%5.2f]", v[0], v[1], v[2])
}

var V3DZero = Vector3f{0, 0, 0}
var V3DUnitX = Vector3f{1, 0, 0}
var V3DUnitY = Vector3f{0, 1, 0}
var V3DUnitZ = Vector3f{0, 0, 1}

// func (p Vector3f) Copy() Vector3f {
// 	return Vector3f{p[0], p[1], p[2]}
// }
func (p Vector3f) Eq(other Vector3f) bool {
	return p == other
	//return p[0] == other[0] && p[1] == other[1] && p[2] == other[2]
}
func (p Vector3f) Ne(other Vector3f) bool {
	return !p.Eq(other)
}
func (p Vector3f) IsZero() bool {
	return p.Eq(V3DZero)
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

// func (p *Vector3f) Normalize() {
// 	d := p.Abs()
// 	if d > 0 {
// 		p[0] /= d
// 		p[1] /= d
// 		p[2] /= d
// 	}
// }
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

// for serialize
func (v Vector3f) NewInt32Vector() [3]int32 {
	return [3]int32{int32(v[0]), int32(v[1]), int32(v[2])}
}

func FromInt32Vector(s [3]int32) Vector3f {
	return Vector3f{float64(s[0]), float64(s[1]), float64(s[2])}
}

func RandVector3D(st, end float64) Vector3f {
	return Vector3f{
		rand.Float64()*(end-st) + st,
		rand.Float64()*(end-st) + st,
		rand.Float64()*(end-st) + st,
	}
}

// func RandVector(st, end Vector3f) Vector3f {
// 	return Vector3f{
// 		rand.Float64()*(end[0]-st[0]) + st[0],
// 		rand.Float64()*(end[1]-st[1]) + st[1],
// 		rand.Float64()*(end[2]-st[2]) + st[2],
// 	}
// }

func (center Vector3f) To8Direct(v2 Vector3f) int {
	rtn := 0
	for i := 0; i < 3; i++ {
		if center[i] > v2[i] {
			rtn += 1 << uint(i)
		}
	}
	return rtn
}

func (h HyperRect) MakeCubeBy8Driect(center Vector3f, direct8 int) HyperRect {
	rtn := Vector3f{}
	for i := 0; i < 3; i++ {
		if direct8&(1<<uint(i)) != 0 {
			rtn[i] = h.Min[i]
		} else {
			rtn[i] = h.Max[i]
		}
	}
	return NewHyperRect(center, rtn)
}

type HyperRect struct {
	Min, Max Vector3f
}

func (h HyperRect) Center() Vector3f {
	return h.Min.Add(h.Max).Idiv(2)
}

func (h HyperRect) DiagLen() float64 {
	return h.Min.LenTo(h.Max)
}

func (h HyperRect) SizeVector() Vector3f {
	return h.Max.Sub(h.Min)
}

func (h HyperRect) IsContact(c Vector3f, r float64) bool {
	hc := h.Center()
	hl := h.DiagLen()
	return hl/2+r >= hc.LenTo(c)
}

func NewHyperRectByCR(c Vector3f, r float64) HyperRect {
	return HyperRect{
		Vector3f{c[0] - r, c[1] - r, c[2] - r},
		Vector3f{c[0] + r, c[1] + r, c[2] + r},
	}
}

func (h HyperRect) RandVector() Vector3f {
	return Vector3f{
		rand.Float64()*(h.Max[0]-h.Min[0]) + h.Min[0],
		rand.Float64()*(h.Max[1]-h.Min[1]) + h.Min[1],
		rand.Float64()*(h.Max[2]-h.Min[2]) + h.Min[2],
	}
}

func (h HyperRect) Move(v Vector3f) HyperRect {
	return HyperRect{
		Min: h.Min.Add(v),
		Max: h.Max.Add(v),
	}
}

func (h HyperRect) IMul(i float64) HyperRect {
	hs := h.SizeVector().Imul(i / 2)
	hc := h.Center()
	return HyperRect{
		Min: hc.Sub(hs),
		Max: hc.Add(hs),
	}
}

// make normalized hyperrect , if not need use HyperRect{Min: , Max:}
func NewHyperRect(v1 Vector3f, v2 Vector3f) HyperRect {
	rtn := HyperRect{
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

func (h1 HyperRect) IsOverlap(h2 HyperRect) bool {
	return !((h1.Min[0] > h2.Max[0] || h1.Max[0] < h2.Min[0]) ||
		(h1.Min[1] > h2.Max[1] || h1.Max[1] < h2.Min[1]) ||
		(h1.Min[2] > h2.Max[2] || h1.Max[2] < h2.Min[2]))
}

func (h1 HyperRect) IsIn(h2 HyperRect) bool {
	for i := 0; i < 3; i++ {
		if h1.Min[i] < h2.Min[i] || h1.Max[i] > h2.Max[i] {
			return false
		}
	}
	return true
}

func (p Vector3f) IsIn(hr HyperRect) bool {
	return hr.Min[0] <= p[0] && p[0] <= hr.Max[0] &&
		hr.Min[1] <= p[1] && p[1] <= hr.Max[1] &&
		hr.Min[2] <= p[2] && p[2] <= hr.Max[2]
}

func (p Vector3f) MakeIn(hr HyperRect) (Vector3f, int) {
	changed := 0
	var i uint
	for i = 0; i < 3; i++ {
		if p[i] > hr.Max[i] {
			p[i] = hr.Max[i]
			changed += 1 << (i*2 + 1)
		}
		if p[i] < hr.Min[i] {
			p[i] = hr.Min[i]
			changed += 1 << (i * 2)
		}
	}
	return p, changed
}
