// Code generated by "genenum.exe -typename=GameObjType -packagename=gameobjtype -basedir=enum -vectortype=int"

package gameobjtype_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/gowasm3dgame/enum/gameobjtype"
)

type GameObjTypeVector [gameobjtype.GameObjType_Count]int

func (es GameObjTypeVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "GameObjTypeVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			gameobjtype.GameObjType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *GameObjTypeVector) Dec(e gameobjtype.GameObjType) {
	es[e] -= 1
}
func (es *GameObjTypeVector) Inc(e gameobjtype.GameObjType) {
	es[e] += 1
}
func (es *GameObjTypeVector) Add(e gameobjtype.GameObjType, v int) {
	es[e] += v
}
func (es *GameObjTypeVector) SetIfGt(e gameobjtype.GameObjType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es GameObjTypeVector) Get(e gameobjtype.GameObjType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es GameObjTypeVector) Iter(fn func(i gameobjtype.GameObjType, v int) bool) bool {
	for i, v := range es {
		if fn(gameobjtype.GameObjType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es GameObjTypeVector) VectorAdd(arg GameObjTypeVector) GameObjTypeVector {
	var rtn GameObjTypeVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es GameObjTypeVector) VectorSub(arg GameObjTypeVector) GameObjTypeVector {
	var rtn GameObjTypeVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *GameObjTypeVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>GameObjType statistics</title>
		</head>
		<body>
		<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table>
	
		<br/>
		</body>
		</html>
		`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, es); err != nil {
		return err
	}
	return nil
}

func Index(i int) string {
	return gameobjtype.GameObjType(i).String()
}

var IndexFn = template.FuncMap{
	"GameObjTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{GameObjTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
