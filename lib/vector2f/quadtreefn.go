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

package vector2f

func (rt Rect) MakeRectBy4Driect(center Vector2f, direct4 int) Rect {
	rtn := Vector2f{}
	for i := 0; i < 2; i++ {
		if direct4&(1<<uint(i)) != 0 {
			rtn[i] = rt.Min[i]
		} else {
			rtn[i] = rt.Max[i]
		}
	}
	return NewRect(center, rtn)

	// w1, w2 := c[0]-rt.X1(), rt.X2()-c[0]
	// h1, h2 := c[1]-rt.Y1(), rt.Y2()-c[1]
	// switch direct4 {
	// case 0:
	// 	return Rect{rt.X, rt.Y, w1, h1}
	// case 1:
	// 	return Rect{c[0], rt.Y, w2, h1}
	// case 2:
	// 	return Rect{c[0], c[1], w2, h2}
	// case 3:
	// 	return Rect{rt.X, c[1], w1, h2}
	// }
	// return Rect{}
}
