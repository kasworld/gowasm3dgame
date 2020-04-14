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

package gameobjtype

import "time"

const LongLife = 3600 * 24 * 365

var Attrib = [GameObjType_Count]struct {
	SpeedLimit float64
	Radius     float64
	MaxInTeam  int
	AddOctree  bool
	LifeTick   int64
}{
	Ball:          {300, 10, 1, true, int64(time.Second) * LongLife},
	Shield:        {400, 5, 32, true, int64(time.Second) * LongLife},
	HommingShield: {400, 7, 16, true, int64(time.Second) * 60},
	Bullet:        {300, 5, 100, true, int64(time.Second) * 10},
	HommingBullet: {300, 7, 10, true, int64(time.Second) * 60},
	SuperBullet:   {600, 15, 10, true, int64(time.Second) * 10},
	BurstBullet:   {300, 5, 100, true, int64(time.Second) * 10},
	HomeMark:      {100, 3, 1, false, int64(time.Second) * LongLife},
	Deco:          {600, 3, 100, false, int64(time.Second) * LongLife},
	Hard:          {0, 3, 1, false, int64(time.Second) * LongLife},
	Food:          {0, 3, 1, true, int64(time.Second) * LongLife},
}

const (
	// MaxRadius need oct tree boundary
	MaxRadius = 15
)

var interactRule = [GameObjType_Count][GameObjType_Count]bool{
	Ball: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	Shield: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	HommingShield: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	Bullet: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	HommingBullet: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	SuperBullet: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	BurstBullet: {
		Ball:          true,
		Shield:        true,
		HommingShield: true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		BurstBullet:   true,
		Hard:          true,
	},
	HomeMark: {},
	Deco:     {},
	Hard:     {},
	Food: {
		Ball: true,
	},
}

func InteractTo(srcType, dstType GameObjType) bool {
	return interactRule[srcType][dstType]
}

func CollisionTo(srcType, dstType GameObjType, sqd float64) bool {
	l := Attrib[srcType].Radius + Attrib[dstType].Radius
	return interactRule[srcType][dstType] && sqd <= l*l
}
