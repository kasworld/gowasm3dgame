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
}{
	Main:          {300, 10},
	Shield:        {200, 5},
	Bullet:        {300, 5},
	HommingBullet: {200, 7},
	SuperBullet:   {600, 15},
	Deco:          {600, 3},
	Mark:          {100, 3},
	Hard:          {0, 3},
	Food:          {0, 3},
}

var collisionRule = [GameObjType_Count][GameObjType_Count]bool{
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

func CollisionTo(srcType, dstType GameObjType) bool {
	return collisionRule[srcType][dstType]
}
