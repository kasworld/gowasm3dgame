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

package gameobjtype

import (
	"math"
	"time"
)

const LongLife = 3600 * 24 * 365

var Attrib = [GameObjType_Count]struct {
	Size           float64
	RadiusToCenter float64 // from main ball center
	SpeedLimit     float64 // speed pixel/sec or rad/sec
	FramePerSec    float64 // animation speed
	LifeTick       int64
}{
	Ball:          {32, 0, 300, 0, int64(time.Second) * LongLife},
	Shield:        {16, 28, math.Pi, 0, int64(time.Second) * LongLife},
	SuperShield:   {16, 48, math.Pi, 30, int64(time.Second) * 60},
	HommingShield: {16, 0, 50, 30, int64(time.Second) * 60},
	Bullet:        {16, 0, 500, 0, int64(time.Second) * LongLife},
	SuperBullet:   {32, 0, 600, 30, int64(time.Second) * LongLife},
	HommingBullet: {16, 0, 300, 30, int64(time.Second) * 60},
}

var collisionRule = [GameObjType_Count][GameObjType_Count]bool{
	Ball: {
		Ball:          true,
		Shield:        true,
		SuperShield:   true,
		Bullet:        true,
		HommingShield: true,
		SuperBullet:   true,
		HommingBullet: true,
	},
	Shield: {
		Ball:          true,
		Shield:        true,
		SuperShield:   true,
		Bullet:        true,
		HommingShield: true,
		SuperBullet:   true,
		HommingBullet: true,
	},
	SuperShield: {
		SuperShield:   true,
		SuperBullet:   true,
		HommingBullet: true,
	},
	Bullet: {
		Ball:          true,
		Shield:        true,
		SuperShield:   true,
		Bullet:        true,
		HommingShield: true,
		SuperBullet:   true,
		HommingBullet: true,
	},
	HommingShield: {
		Ball:          true,
		Shield:        true,
		SuperShield:   true,
		Bullet:        true,
		HommingShield: true,
		SuperBullet:   true,
		HommingBullet: true,
	},
	SuperBullet: {
		SuperShield:   true,
		SuperBullet:   true,
		HommingBullet: true,
	},
	HommingBullet: {
		SuperShield:   true,
		SuperBullet:   true,
		HommingBullet: true,
	},
}

func InteractTo(srcType, dstType GameObjType) bool {
	return collisionRule[srcType][dstType]
}
