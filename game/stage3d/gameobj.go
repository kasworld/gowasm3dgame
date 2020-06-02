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

package stage3d

import (
	"time"

	"github.com/kasworld/go-abs"
	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
)

type GameObj struct {
	GOType   gameobjtype.GameObjType
	UUID     string
	TeamUUID string

	PosVt vector3f.Vector3f
	RotVt vector3f.Vector3f

	VelVt    vector3f.Vector3f
	RotVelVt vector3f.Vector3f

	BirthTick    int64
	LastMoveTick int64
	toDelete     bool

	DstUUID string
}

func (o *GameObj) Pos() vector3f.Vector3f {
	return o.PosVt
}

func (o *GameObj) GetCube() vector3f.Cube {
	r := gameobjtype.Attrib[o.GOType].Radius
	return vector3f.NewCubeByCR(
		o.PosVt, r,
	)
}

func (o *GameObj) IsCollision(dst *GameObj) bool {
	return gameobjtype.CollisionTo(
		o.GOType, dst.GOType,
		dst.PosVt.Sqd(o.PosVt),
	)
}

func (o *GameObj) ToPacket(co uint32) *w3d_obj.GameObj {
	return &w3d_obj.GameObj{
		GOType:  o.GOType,
		UUID:    o.UUID,
		Color24: co,
		PosVt: [3]float32{
			float32(o.PosVt[0]),
			float32(o.PosVt[1]),
			float32(o.PosVt[2]),
		},
		RotVt: [3]float32{
			float32(o.RotVt[0]),
			float32(o.RotVt[1]),
			float32(o.RotVt[2]),
		},
	}
}

func (o *GameObj) CheckLife(now int64) bool {
	lifetick := gameobjtype.Attrib[o.GOType].LifeTick
	return now-o.BirthTick < lifetick
}

func (o *GameObj) BounceNormalize(border vector3f.Cube) {
	for i := 0; i < 3; i++ {
		if o.PosVt[i] < border.Min[i] {
			o.PosVt[i] = border.Min[i]
			o.VelVt[i] = abs.Absf(o.VelVt[i])
		}
		if o.PosVt[i] > border.Max[i] {
			o.PosVt[i] = border.Max[i]
			o.VelVt[i] = -abs.Absf(o.VelVt[i])
		}
	}
}

// CalcLenChange calc two gameobj change of len with time
// return current len , len change with time
// currentlen adjust with obj size
func (o *GameObj) CalcLenChange(dsto *GameObj) (float64, float64) {
	r1 := gameobjtype.Attrib[o.GOType].Radius / 2
	r2 := gameobjtype.Attrib[dsto.GOType].Radius / 2
	curLen := dsto.PosVt.Sub(o.PosVt).Abs()
	nextLen := dsto.PosVt.Add(dsto.VelVt).Sub(
		o.PosVt.Add(o.VelVt),
	).Abs()
	lenChange := nextLen - curLen
	return curLen - r1 - r2, lenChange
}

/////////////////

func (o *GameObj) AccelTo(dstPosVt vector3f.Vector3f) {
	mvLimit := gameobjtype.Attrib[o.GOType].SpeedLimit
	diff := dstPosVt.Sub(o.PosVt)
	if diff.Abs() > mvLimit {
		diff = diff.NormalizedTo(mvLimit)
	}
	o.VelVt = o.VelVt.Add(diff)
	if o.VelVt.Abs() > mvLimit {
		o.VelVt = o.VelVt.NormalizedTo(mvLimit)
	}
}

func (o *GameObj) Move_straight(now int64) {
	dur := float64(now-o.LastMoveTick) / float64(time.Second)
	o.RotVt = o.RotVt.Add(o.RotVelVt.MulF(dur))
	o.LastMoveTick = now

	mvLimit := gameobjtype.Attrib[o.GOType].SpeedLimit
	if o.VelVt.Abs() > mvLimit {
		o.VelVt = o.VelVt.NormalizedTo(mvLimit)
	}
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(dur))
}

func (o *GameObj) Move_circular(now int64, dstObj *GameObj) {
	lifedur := float64(now-o.BirthTick) / float64(time.Second)
	orbitR := gameobjtype.Attrib[gameobjtype.Ball].Radius * 4

	rotAxis := dstObj.VelVt.NormalizedTo(1).Add(o.RotVt.NormalizedTo(1))
	o.RotVt = rotAxis
	refPos := rotAxis.Cross(o.VelVt).NormalizedTo(orbitR)
	shieldPosDiff := refPos.RotateAround(rotAxis, lifedur)
	dstPos := shieldPosDiff.Add(dstObj.PosVt)
	o.PosVt = dstPos
}

func (o *GameObj) Move_hommingshield(now int64, dstObj *GameObj) {
	lifedur := float64(now-o.BirthTick) / float64(time.Second)
	orbitR := gameobjtype.Attrib[gameobjtype.Ball].Radius * 4
	p := dstObj.VelVt.Cross(o.VelVt).NormalizedTo(orbitR)
	axis := dstObj.VelVt
	diffVt := p.RotateAround(axis, lifedur)
	dstPos := dstObj.PosVt.Add(diffVt)
	o.AccelTo(dstPos)
	o.Move_straight(now)
	// o.PosVt = dstPos
}

func (o *GameObj) Move_hommingbullet(now int64, dstObj *GameObj) {
	o.Move_straight(now)
	o.AccelTo(dstObj.PosVt)
}
