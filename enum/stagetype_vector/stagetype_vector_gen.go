// Code generated by "genenum.exe -typename=StageType -packagename=stagetype -basedir=enum -vectortype=int"

package stagetype_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/gowasm3dgame/enum/stagetype"
)

type StageTypeVector [stagetype.StageType_Count]int

func (es StageTypeVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "StageTypeVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			stagetype.StageType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *StageTypeVector) Dec(e stagetype.StageType) {
	es[e] -= 1
}
func (es *StageTypeVector) Inc(e stagetype.StageType) {
	es[e] += 1
}
func (es *StageTypeVector) Add(e stagetype.StageType, v int) {
	es[e] += v
}
func (es *StageTypeVector) SetIfGt(e stagetype.StageType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es StageTypeVector) Get(e stagetype.StageType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es StageTypeVector) Iter(fn func(i stagetype.StageType, v int) bool) bool {
	for i, v := range es {
		if fn(stagetype.StageType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es StageTypeVector) VectorAdd(arg StageTypeVector) StageTypeVector {
	var rtn StageTypeVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es StageTypeVector) VectorSub(arg StageTypeVector) StageTypeVector {
	var rtn StageTypeVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *StageTypeVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>StageType statistics</title>
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
	return stagetype.StageType(i).String()
}

var IndexFn = template.FuncMap{
	"StageTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{StageTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
