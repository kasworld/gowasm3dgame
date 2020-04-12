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

package stage

import "fmt"

func (stg *Stage) String() string {
	return fmt.Sprintf("Stage[%v Conn%v , Team%v]",
		stg.UUID, stg.Conns, stg.Teams,
	)
}

const (
	HTML_tableheader = `
<tr>
<th>UUID</th>
<th>Conn</th>
<th>Team</th>
<th>Command</th>
</tr>`

	HTML_row = `
<tr>
<td>{{$v.UUID}}</td>
<td>{{$v.Conns}}</td>
<td>{{$v.Teams}}</td>
<td><a href="/Del?id={{$v.UUID}}" target="_blank">[Del]</a></td>
</tr>
`
)
