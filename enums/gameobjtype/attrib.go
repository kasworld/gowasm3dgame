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

var Attrib = [GameObjType_Count]struct {
	SpeedLimit float64
	Radius     float64
	AddOctree  bool
}{
	Main:          {300, 10, true},
	Shield:        {200, 5, true},
	Bullet:        {300, 5, true},
	HommingBullet: {200, 7, true},
	SuperBullet:   {600, 15, true},
	Deco:          {600, 3, false},
	Mark:          {100, 3, false},
	Hard:          {0, 3, false},
	Food:          {0, 3, true},
}

const (
	// MaxRadius need oct tree boundary
	MaxRadius = 15
)

var interactRule = [GameObjType_Count][GameObjType_Count]bool{
	Main: {
		Main:          true,
		Shield:        true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		Hard:          true,
	},
	Shield: {
		Main:          true,
		Shield:        true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		Hard:          true,
	},
	Bullet: {
		Main:          true,
		Shield:        true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		Hard:          true,
	},
	HommingBullet: {
		Main:          true,
		Shield:        true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		Hard:          true,
	},
	SuperBullet: {
		Main:          true,
		Shield:        true,
		Bullet:        true,
		HommingBullet: true,
		SuperBullet:   true,
		Hard:          true,
	},
	Deco: {},
	Mark: {},
	Hard: {},
	Food: {
		Main: true,
	},
}

func InteractTo(srcType, dstType GameObjType) bool {
	return interactRule[srcType][dstType]
}

func CollisionTo(srcType, dstType GameObjType, sqd float64) bool {
	l := Attrib[srcType].Radius + Attrib[dstType].Radius
	return interactRule[srcType][dstType] && sqd <= l*l
}
