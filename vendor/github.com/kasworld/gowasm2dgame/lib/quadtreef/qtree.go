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

package quadtreef

import (
	"bytes"
	"fmt"

	"github.com/kasworld/gowasm2dgame/lib/vector2f"
)

const (
	MaxQuadTreeData = 4
)

type QuadTreeObjI interface {
	GetRect() vector2f.Rect
	GetUUID() string
}

type QuadTree struct {
	BoundRect    vector2f.Rect
	Center       [2]float64
	DataList     []QuadTreeObjI
	Children     [4]*QuadTree
	TerminalNode bool // cannot split
}

func New(rt vector2f.Rect) *QuadTree {
	rtn := QuadTree{
		BoundRect: rt,
		Center:    rt.Center(),
		DataList:  make([]QuadTreeObjI, 0, MaxQuadTreeData),
	}
	szvt := rt.SizeVector()
	if szvt[0] < 2 || szvt[1] < 2 { // cannot divide
		rtn.TerminalNode = true
	} else {
		rtn.TerminalNode = false
	}
	return &rtn
}

func (ot *QuadTree) Insert(o QuadTreeObjI) bool {
	if !o.GetRect().IsIn(ot.BoundRect) {
		return false
	}
	if ot.Children[0] != nil { // splited
		if !ot.insertChild(o) { // append to me
			ot.DataList = append(ot.DataList, o)
		}
		return true
	} else { // not splited
		if ot.TerminalNode || len(ot.DataList) < MaxQuadTreeData { // check need split
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

func (ot *QuadTree) insertChild(o QuadTreeObjI) bool {
	for _, chot := range ot.Children { // try child
		if chot.Insert(o) {
			return true
		}
	}
	return false
}

func (ot *QuadTree) split() {
	if ot.Children[0] != nil {
		return
	}
	// split all data and make datalist nil
	for i, _ := range ot.Children {
		newbound := ot.BoundRect.MakeRectBy4Driect(ot.Center, i)
		ot.Children[i] = New(newbound)
	}
	// move this node data to child
	newDataList := make([]QuadTreeObjI, 0, len(ot.DataList))
	for _, o := range ot.DataList {
		if !ot.insertChild(o) {
			newDataList = append(newDataList, o)
		}
	}
	ot.DataList = newDataList
	return
}

func (ot *QuadTree) Remove(o QuadTreeObjI) bool {
	if !o.GetRect().IsIn(ot.BoundRect) {
		return false
	}
	for i, v := range ot.DataList {
		if v == o {
			ot.DataList = append(ot.DataList[:i], ot.DataList[i+1:]...)
			return true
		}
	}
	if ot.Children[0] != nil {
		for _, chot := range ot.Children {
			if chot.Remove(o) {
				return true
			}
		}
	}
	fmt.Printf("not found %v\n", o.GetUUID())
	return false
}

func (ot *QuadTree) RemoveByID(o QuadTreeObjI) bool {
	if !o.GetRect().IsIn(ot.BoundRect) {
		return false
	}
	for i, v := range ot.DataList {
		if v.GetUUID() == o.GetUUID() {
			ot.DataList = append(ot.DataList[:i], ot.DataList[i+1:]...)
			return true
		}
	}
	if ot.Children[0] != nil {
		for _, chot := range ot.Children {
			if chot.Remove(o) {
				return true
			}
		}
	}
	return false
}

func (ot *QuadTree) QueryByRect(fn func(QuadTreeObjI) bool, hr vector2f.Rect) bool {
	if !ot.BoundRect.IsOverlap(hr) {
		return false
	}
	for _, o := range ot.DataList {
		if !o.GetRect().IsOverlap(hr) {
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
		quit := o.QueryByRect(fn, hr)
		if quit {
			return true
		}
	}
	return false
}

// true : exist , false : not exist
func (ot *QuadTree) QueryByPos(fn func(QuadTreeObjI) bool, pos [2]float64) bool {
	if !ot.BoundRect.Contain(pos) {
		return false
	}
	for _, o := range ot.DataList {
		if !o.GetRect().Contain(pos) {
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
		quit := o.QueryByPos(fn, pos)
		if quit {
			return true
		}
	}
	return false
}

func (ot QuadTree) String() string {
	return fmt.Sprintf("QuadTree[%v]", ot.BoundRect)
}
func (ot QuadTree) DebugString() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "QuadTree[%v %v \n", ot.BoundRect, ot.Center)
	for _, d := range ot.DataList {
		fmt.Fprintf(&b, "%v ", d)
	}
	fmt.Fprintf(&b, "\n")
	if ot.Children[0] != nil {
		for _, cot := range ot.Children {
			b.WriteString(cot.String())
		}
	}
	fmt.Fprintf(&b, "]\n")
	return b.String()
}
