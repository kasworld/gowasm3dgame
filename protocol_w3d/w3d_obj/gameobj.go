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

package w3d_obj

import (
	"time"

	"github.com/kasworld/gowasm3dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
	"github.com/kasworld/htmlcolors"
)

type Team struct {
	ID       string
	Color24  htmlcolors.Color24
	Ball     *GameObj
	HomeMark *GameObj
	Objs     []*GameObj
}

type GameObj struct {
	GOType       gameobjtype.GameObjType
	UUID         string
	BirthTick    int64
	LastMoveTick int64 // time.unixnano
	PosVt        vector3f.Vector3f
	MvVt         vector3f.Vector3f
	AccVt        vector3f.Vector3f
	DstUUID      string // move to dest
}

func (o *GameObj) IsCollision(dst *GameObj) bool {
	return gameobjtype.CollisionTo(
		o.GOType, dst.GOType,
		dst.PosVt.Sqd(o.PosVt),
	)
}

func (o *GameObj) MoveStraight(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.MvVt.MulF(diff))
}

func (o *GameObj) MoveHomming(now int64, dstPosVt vector3f.Vector3f) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.MvVt.MulF(diff))

	maxv := gameobjtype.Attrib[o.GOType].SpeedLimit
	dxyVt := dstPosVt.Sub(o.PosVt)
	o.MvVt = o.MvVt.Add(dxyVt.Normalize().MulF(maxv))
}

////////////////////

func (o *GameObj) Move_accel(now int64) bool {
	dur := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.MvVt = o.MvVt.Add(o.AccVt.MulF(dur))
	if o.MvVt.Abs() > gameobjtype.Attrib[o.GOType].SpeedLimit {
		o.MvVt = o.MvVt.NormalizedTo(gameobjtype.Attrib[o.GOType].SpeedLimit)
	}
	o.PosVt = o.PosVt.Add(o.MvVt.MulF(dur))
	return true
}

func (o *GameObj) Move_rand(now int64, rndAccVt vector3f.Vector3f) bool {
	dur := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now

	o.AccVt = rndAccVt
	o.MvVt = o.MvVt.Add(o.AccVt.MulF(dur))
	if o.MvVt.Abs() > gameobjtype.Attrib[o.GOType].SpeedLimit {
		o.MvVt = o.MvVt.NormalizedTo(gameobjtype.Attrib[o.GOType].SpeedLimit)
	}
	o.PosVt = o.PosVt.Add(o.MvVt.MulF(dur))
	return true
}

func (o *GameObj) Move_shield(now int64, dstObj *GameObj) bool {
	dur := float64(now-o.LastMoveTick) / float64(time.Second)
	axis := dstObj.MvVt
	p := dstObj.MvVt.Cross(o.MvVt).NormalizedTo(20)
	o.PosVt = dstObj.PosVt.Add(p.RotateAround(axis, dur+o.AccVt.Abs()))
	return true
}

func (o *GameObj) Move_homming(now int64, dstObj *GameObj) bool {
	// how to other team obj pos? without panic
	o.AccVt = dstObj.PosVt.Sub(o.PosVt).NormalizedTo(gameobjtype.Attrib[o.GOType].SpeedLimit)
	return o.Move_accel(now)
}
