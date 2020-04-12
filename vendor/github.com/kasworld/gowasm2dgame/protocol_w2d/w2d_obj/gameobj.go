package w2d_obj

import (
	"math"
	"time"

	"github.com/kasworld/gowasm2dgame/enums/effecttype"
	"github.com/kasworld/gowasm2dgame/enums/gameobjtype"
	"github.com/kasworld/gowasm2dgame/enums/teamtype"
	"github.com/kasworld/gowasm2dgame/lib/vector2f"
)

type Cloud struct {
	SpriteNum    int
	LastMoveTick int64 // time.unixnano
	PosVt        vector2f.Vector2f
	VelVt        vector2f.Vector2f
}

func (o *Cloud) Move(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))
}

type Background struct {
	LastMoveTick int64 // time.unixnano
	PosVt        vector2f.Vector2f
	VelVt        vector2f.Vector2f
}

func (o *Background) Move(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))
}

type Effect struct {
	EffectType   effecttype.EffectType
	BirthTick    int64
	LastMoveTick int64 // time.unixnano
	PosVt        vector2f.Vector2f
	VelVt        vector2f.Vector2f
}

func (o *Effect) CheckLife(now int64) bool {
	lifetick := effecttype.Attrib[o.EffectType].LifeTick
	return now-o.BirthTick < lifetick
}

func (o *Effect) Move(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))
	o.VelVt = o.VelVt.MulF(0.9)
}

////////////////////

type Team struct {
	TeamType teamtype.TeamType
	Ball     *GameObj
	Objs     []*GameObj
}

type GameObj struct {
	GOType       gameobjtype.GameObjType
	UUID         string
	BirthTick    int64
	LastMoveTick int64 // time.unixnano
	PosVt        vector2f.Vector2f
	VelVt        vector2f.Vector2f
	Angle        float64 // move circular
	AngleV       float64
	DstUUID      string // move to dest
}

func (o *GameObj) MoveStraight(now int64) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))
}

func (o *GameObj) MoveCircular(now int64, center vector2f.Vector2f) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.Angle += o.AngleV * diff
	r := gameobjtype.Attrib[o.GOType].RadiusToCenter
	o.PosVt = o.CalcCircularPos(center, r)
}

func (o *GameObj) CalcCircularPos(center vector2f.Vector2f, r float64) vector2f.Vector2f {
	rpos := vector2f.Vector2f{r * math.Cos(o.Angle), r * math.Sin(o.Angle)}
	return center.Add(rpos)
}

func (o *GameObj) MoveHomming(now int64, dstPosVt vector2f.Vector2f) {
	diff := float64(now-o.LastMoveTick) / float64(time.Second)
	o.LastMoveTick = now
	o.PosVt = o.PosVt.Add(o.VelVt.MulF(diff))

	maxv := gameobjtype.Attrib[o.GOType].SpeedLimit
	dxyVt := dstPosVt.Sub(o.PosVt)
	o.VelVt = o.VelVt.Add(dxyVt.Normalize().MulF(maxv))
}
