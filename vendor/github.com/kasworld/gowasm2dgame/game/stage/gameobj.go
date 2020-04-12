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

package stage

import (
	"math"
	"time"

	"github.com/kasworld/go-abs"
	"github.com/kasworld/gowasm2dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm2dgame/enums/teamtype"
	"github.com/kasworld/gowasm2dgame/lib/vector2f"
	"github.com/kasworld/gowasm2dgame/protocol_w2d/w2d_obj"
)

func (o *GameObj) GetUUID() string {
	return o.UUID
}
func (o *GameObj) GetRect() vector2f.Rect {
	r := gameobjtype.Attrib[o.GOType].Size
	return vector2f.NewRectCenterWH(
		o.PosVt,
		vector2f.Vector2f{r, r},
	)
}

type GameObj struct {
	teamType     teamtype.TeamType
	GOType       gameobjtype.GameObjType
	UUID         string
	BirthTick    int64
	LastMoveTick int64 // time.unixnano
	toDelete     bool

	PosVt vector2f.Vector2f
	VelVt  vector2f.Vector2f

	Angle  float64 // move circular
	AngleV float64

	DstUUID string // move to dest
}

func (o *GameObj) ToPacket() *w2d_obj.GameObj {
	return &w2d_obj.GameObj{
		GOType:       o.GOType,
		UUID:         o.UUID,
		BirthTick:    o.BirthTick,
		LastMoveTick: o.LastMoveTick,
		PosVt:        o.PosVt,
		VelVt:         o.VelVt,
		Angle:        o.Angle,
		AngleV:       o.AngleV,
		DstUUID:      o.DstUUID,
	}
}

func (o *GameObj) MoveStraight(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))
}

func (o *GameObj) MoveCircular(now int64, ctvt vector2f.Vector2f) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.Angle += o.AngleV * diff
	r := gameobjtype.Attrib[o.GOType].RadiusToCenter
	o.PosVt = o.CalcCircularPos(ctvt, r)
}

func (o *GameObj) MoveHommingShield(now int64, dstPosVt vector2f.Vector2f) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))

	maxv := gameobjtype.Attrib[o.GOType].SpeedLimit

	dxyVt := dstPosVt.Sub(o.PosVt)
	o.VelVt = o.VelVt.Add(dxyVt.Normalize().MulF(maxv))
}

func (o *GameObj) MoveHommingBullet(now int64, dstPosVt vector2f.Vector2f) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))

	maxv := gameobjtype.Attrib[o.GOType].SpeedLimit
	dxyVt := dstPosVt.Sub(o.PosVt)
	o.VelVt = o.VelVt.Add(dxyVt.Normalize().MulF(maxv))
	o.LimitDxy()
}

func (o *GameObj) CheckLife(now int64) bool {
	lifetick := gameobjtype.Attrib[o.GOType].LifeTick
	return now-o.BirthTick < lifetick
}

func (o *GameObj) IsIn(w, h float64) bool {
	return 0 <= o.PosVt[0] && o.PosVt[0] <= w && 0 <= o.PosVt[1] && o.PosVt[1] <= h
}

func (o *GameObj) SetDxy(vt vector2f.Vector2f) {
	o.VelVt = vt
	o.LimitDxy()
}

func (o *GameObj) AddDxy(vt vector2f.Vector2f) {
	o.VelVt = o.VelVt.Add(vt)
	o.LimitDxy()
}

func (o *GameObj) LimitDxy() {
	maxv := gameobjtype.Attrib[o.GOType].SpeedLimit
	if o.VelVt.Abs() > maxv {
		o.VelVt = o.VelVt.Normalize().MulF(maxv)
	}
}

func (o *GameObj) BounceNormalize(w, h float64) {
	if o.PosVt[0] < 0 {
		o.PosVt[0] = 0
		o.VelVt[0] = abs.Absf(o.VelVt[0])
	}
	if o.PosVt[1] < 0 {
		o.PosVt[1] = 0
		o.VelVt[1] = abs.Absf(o.VelVt[1])
	}

	if o.PosVt[0] > w {
		o.PosVt[0] = w
		o.VelVt[0] = -abs.Absf(o.VelVt[0])
	}
	if o.PosVt[1] > h {
		o.PosVt[1] = h
		o.VelVt[1] = -abs.Absf(o.VelVt[1])
	}
}

func (o *GameObj) CalcCircularPos(center vector2f.Vector2f, r float64) vector2f.Vector2f {
	rpos := vector2f.Vector2f{r * math.Cos(o.Angle), r * math.Sin(o.Angle)}
	return center.Add(rpos)
}

// CalcLenChange calc two gameobj change of len with time
// return current len , len change with time
// currentlen adjust with obj size
func (o *GameObj) CalcLenChange(dsto *GameObj) (float64, float64) {
	r1 := gameobjtype.Attrib[o.GOType].Size / 2
	r2 := gameobjtype.Attrib[dsto.GOType].Size / 2
	curLen := dsto.PosVt.Sub(o.PosVt).Abs()
	nextLen := dsto.PosVt.Add(dsto.VelVt).Sub(
		o.PosVt.Add(o.VelVt),
	).Abs()
	lenChange := nextLen - curLen
	return curLen - r1 - r2, lenChange
}
