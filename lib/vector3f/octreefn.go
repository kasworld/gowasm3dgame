// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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

func (h Cube) MakeCubeBy8Driect(center Vector3f, direct8 int) Cube {
	rtn := Vector3f{}
	for i := 0; i < 3; i++ {
		if direct8&(1<<uint(i)) != 0 {
			rtn[i] = h.Min[i]
		} else {
			rtn[i] = h.Max[i]
		}
	}
	return NewCube(center, rtn)
}

func (center Vector3f) To8Direct(v2 Vector3f) int {
	rtn := 0
	for i := 0; i < 3; i++ {
		if center[i] > v2[i] {
			rtn += 1 << uint(i)
		}
	}
	return rtn
}
