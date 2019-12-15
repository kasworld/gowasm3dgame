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

package octree

import (
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
)

const (
	MaxOctreeData = 8
)

type OctreeObjI interface {
	Pos() vector3f.Vector3f
	GetCube() vector3f.Cube
}

type OctreeObjList []OctreeObjI

type Octree struct {
	BoundCube    vector3f.Cube
	Center       vector3f.Vector3f
	DataList     OctreeObjList
	Children     [8]*Octree
	TerminalNode bool // cannot split
}

func New(cube vector3f.Cube) *Octree {
	rtn := Octree{
		BoundCube: cube,
		DataList:  make(OctreeObjList, 0, MaxOctreeData),
		Center:    cube.Center(),
	}
	szvt := cube.SizeVector()
	for i := 0; i < 3; i++ {
		if szvt[i] < 2 { // cannot divide
			rtn.TerminalNode = true
		}
	}
	return &rtn
}

func (ot *Octree) Insert(o OctreeObjI) bool {
	//log.Printf("insert to octree obj%v %v", o.ID, o.Pos())
	if !o.Pos().IsIn(ot.BoundCube) {
		// log.Printf("invalid Insert Octree %v %v", ot.BoundCube, o.Pos())
		return false
	}

	if ot.Children[0] != nil { // splited
		if !ot.insertChild(o) { // append to me
			ot.DataList = append(ot.DataList, o)
		}
		return true
	} else { // not splited
		if ot.TerminalNode || len(ot.DataList) < MaxOctreeData { // check need split
			// simple append
			ot.DataList = append(ot.DataList, o)
			return true
		} else {
			ot.split()
			if !ot.insertChild(o) { // append to me
				ot.DataList = append(ot.DataList, o)
			}
			return true
		}
	}

}

func (ot *Octree) insertChild(o OctreeObjI) bool {
	for _, chot := range ot.Children { // try child
		if chot.Insert(o) {
			return true
		}
	}
	return false
}

func (ot *Octree) split() {
	if ot.Children[0] != nil {
		return
	}
	// split all data and make datalist nil
	//log.Printf("split octree %v %v", ot.BoundCube, ot.Center)
	for i, _ := range ot.Children {
		newbound := ot.BoundCube.MakeCubeBy8Driect(ot.Center, i)
		ot.Children[i] = New(newbound)
	}
	// move this node data to child
	newDataList := make([]OctreeObjI, 0, len(ot.DataList))
	for _, o := range ot.DataList {
		if !ot.insertChild(o) {
			newDataList = append(newDataList, o)
		}
	}
	ot.DataList = newDataList
	return

}

func (ot *Octree) QueryByCube(
	fn func(OctreeObjI) bool, hr vector3f.Cube) bool {
	if !ot.BoundCube.IsOverlap(hr) {
		return false
	}
	for _, o := range ot.DataList {
		if !o.GetCube().IsOverlap(hr) {
			continue
		}
		if fn == nil || fn(o) {
			return true
		}
	}
	if ot.Children[0] == nil {
		return false
	}
	for _, o := range ot.Children {
		quit := o.QueryByCube(fn, hr)
		if quit {
			return true
		}
	}
	return false
}
