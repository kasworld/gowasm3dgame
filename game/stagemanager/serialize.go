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

package stagemanager

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/kasworld/weblib"
)

func (man *Manager) String() string {
	return fmt.Sprintf(
		"StageManager[Count:%v]",
		man.Count(),
	)
}

func (man *Manager) ToWeb(w http.ResponseWriter, r *http.Request) {
	weblib.WebFormBegin("Stage list", w, r)
	man.ToWebMid(w, r)
	weblib.WebFormEnd(w, r)
}

func (man *Manager) ToWebMid(w http.ResponseWriter, r *http.Request) {

	connList := man.GetList()
	sort.Sort(stageList(connList))
	page := weblib.GetIntByName("page", -1, w, r)
	if page < 0 {
		return
	}
	pagesize := 20

	weblib.PageMid(len(connList), pagesize, "/StageManager", w, r)

	st := page * pagesize
	if st < 0 || st >= len(connList) {
		st = 0
	}

	ed := st + pagesize
	if ed > len(connList) {
		ed = len(connList)
	}

	rtn := connList[st:ed]

	tplIndex, err := template.New("index").Parse(`
	<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table>
	<br/>
	`)
	if err != nil {
		man.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, rtn); err != nil {
		man.log.Error("%v", err)
	}
}

const (
	HTML_tableheader = `
<tr>
<th>UUID</th>
<th>Conn</th>
<th>Stage</th>
<th>Command</th>
</tr>`

	HTML_row = `
<tr>
<td>{{$v.GetUUID}}</td>
<td>{{$v.GetConnManager}}</td>
<td>{{$v}}</td>
<td><a href="/Del?id={{$v.GetUUID}}" target="_blank">[Del]</a></td>
</tr>
`
)
