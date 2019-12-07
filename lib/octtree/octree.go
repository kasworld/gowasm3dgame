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

package octtree

import (
	"github.com/kasworld/gowasm3dgame/lib/vector3f"
)

const (
	MaxOctreeData = 8
)

type OctreeObjI interface {
	Pos() vector3f.Vector3f
}

type OctreeObjList []OctreeObjI

type Octree struct {
	BoundCube *vector3f.HyperRect
	Center    vector3f.Vector3f
	DataList  OctreeObjList
	Children  [8]*Octree
}

func NewOctree(cube *vector3f.HyperRect) *Octree {
	rtn := Octree{
		BoundCube: cube,
		DataList:  make(OctreeObjList, 0, MaxOctreeData),
		Center:    cube.Center(),
	}
	//log.Printf("new octree %v", rtn.BoundCube)
	return &rtn
}

func (ot *Octree) Split() {
	if ot.Children[0] != nil {
		return
	}
	// split all data and make datalist nil
	//log.Printf("split octree %v %v", ot.BoundCube, ot.Center)
	for i, _ := range ot.Children {
		newbound := ot.BoundCube.MakeCubeBy8Driect(ot.Center, i)
		ot.Children[i] = NewOctree(newbound)
	}
}

func (ot *Octree) Insert(o OctreeObjI) bool {
	//log.Printf("insert to octree obj%v %v", o.ID, o.Pos())
	if !o.Pos().IsIn(ot.BoundCube) {
		// log.Printf("invalid Insert Octree %v %v", ot.BoundCube, o.Pos())
		return false
	}
	if len(ot.DataList) < MaxOctreeData {
		// simple append
		ot.DataList = append(ot.DataList, o)
		return true
	} else {
		ot.Split()
		d8 := ot.Center.To8Direct(o.Pos())
		return ot.Children[d8].Insert(o)
	}
}

func (ot *Octree) QueryByHyperRect(fn func(OctreeObjI) bool, hr *vector3f.HyperRect) bool {
	if !ot.BoundCube.IsOverlap(hr) {
		return false
	}
	for _, o := range ot.DataList {
		if !o.Pos().IsIn(hr) {
			continue
		}
		if fn(o) {
			return true
		}
	}
	if ot.Children[0] == nil {
		return false
	}
	for _, o := range ot.Children {
		quit := o.QueryByHyperRect(fn, hr)
		if quit {
			return true
		}
	}
	return false
}
