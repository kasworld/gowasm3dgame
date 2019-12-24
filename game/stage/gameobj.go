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

package stage

import (
	"time"

	"github.com/kasworld/go-abs"
	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
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

func (o *GameObj) ToPacket() *w3d_obj.GameObj {
	return &w3d_obj.GameObj{
		GOType: o.GOType,
		UUID:   o.UUID,
		PosVt:  o.PosVt,
		RotVt:  o.RotVt,
	}
}

func (o *GameObj) CheckLife(now int64) bool {
	lifetick := gameobjtype.Attrib[o.GOType].LifeTick
	return now-o.BirthTick < lifetick
}

func (o *GameObj) BounceNormalize(size float64) {
	for i := 0; i < 3; i++ {
		if o.PosVt[i] < 0 {
			o.PosVt[i] = 0
			o.VelVt[i] = abs.Absf(o.VelVt[i])
		}
		if o.PosVt[i] > size {
			o.PosVt[i] = size
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
	dur := float64(now-o.LastMoveTick) / float64(time.Second)
	o.RotVt = o.RotVt.Add(o.RotVelVt.MulF(dur))
	o.LastMoveTick = now

	mvLimit := gameobjtype.Attrib[o.GOType].SpeedLimit
	p := dstObj.VelVt.Cross(o.VelVt).NormalizedTo(mvLimit)
	axis := dstObj.VelVt
	diffVt := p.RotateAround(axis, lifedur)
	o.PosVt = dstObj.PosVt.Add(diffVt)
}

func (o *GameObj) Move_homming(now int64, dstObj *GameObj) {
	o.Move_straight(now)
	o.AccelTo(dstObj.PosVt)
}
