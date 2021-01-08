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

package wasmclient

import (
	"bytes"
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/kasworld/gowasm3dgame/config/dataversion"
	"github.com/kasworld/gowasm3dgame/enum/acttype"
	"github.com/kasworld/gowasm3dgame/game/stagelist4client"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_version"
	"github.com/kasworld/gowasmlib/jslog"
)

func (app *WasmClient) makeButtons() string {
	var buf bytes.Buffer
	gameOptions.MakeButtonToolTipTop(&buf)
	return buf.String()
}

func (app *WasmClient) DisplayTextInfo() {
	app.updateLeftInfo()
	app.updateRightInfo()
	app.updateCenterInfo()
}

func (app *WasmClient) makeServiceInfo() string {
	msgCopyright := `</hr>Copyright 2019,2020 SeukWon Kang 
		<a href="https://github.com/kasworld/gowasm3dgame" target="_blank">gowasm3dgame</a>`

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "gowasm3dgame webclient<br/>")
	fmt.Fprintf(&buf, "Protocol %v<br/>", w3d_version.ProtocolVersion)
	fmt.Fprintf(&buf, "Data %v<br/>", dataversion.DataVersion)
	fmt.Fprintf(&buf, "%v<br/>", msgCopyright)
	return buf.String()
}

func (app *WasmClient) makeDebugInfo() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf,
		"%v<br/>Ping %v<br/>ServerClientTickDiff %v<br/>",
		app.DispInterDur, app.PingDur, app.ServerClientTictDiff,
	)
	fmt.Fprintf(&buf, "scene obj count %v<br/>",
		len(app.vp.jsSceneObjs),
	)
	fmt.Fprintf(&buf, "geomatry cache count %v<br/>",
		len(app.vp.geometryCache),
	)
	fmt.Fprintf(&buf, "material cache count %v<br/>",
		len(app.vp.materialCache),
	)
	fmt.Fprintf(&buf, "lighthelp cache count %v<br/>",
		len(app.vp.lightCache),
	)

	return buf.String()
}

func (app *WasmClient) makeTeamStatsInfo() string {
	stats := app.statsInfo
	if stats == nil {
		return ""
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%v %v<br/>",
		app.loginData.StageType, app.loginData.StageUUID)

	buf.WriteString(`<table border=1 style="border-collapse:collapse;">`)
	buf.WriteString(`<tr>
	<th>team \ act</th>
	<th>UUID</th>
	<th>AP</th>
	<th>Alive</th>
	<th>Score</th>
	<th>Kill/Death</th>
	`)
	for acti := 0; acti < acttype.ActType_Count; acti++ {
		fmt.Fprintf(&buf, "<th>Act %v</th>", acttype.ActType(acti))
	}
	buf.WriteString(`</tr>`)

	for ti, tv := range stats.Stats {
		fmt.Fprintf(&buf, `<tr style="background-color:#%06x">`,
			tv.Color24)
		fmt.Fprintf(&buf, "<td>%v</td>", ti)
		fmt.Fprintf(&buf, "<td>%v</td>", tv.UUID)
		fmt.Fprintf(&buf, "<td>%v</td>", tv.AP)
		fmt.Fprintf(&buf, "<td>%v</td>", tv.Alive)
		fmt.Fprintf(&buf, "<td>%v</td>", tv.Score)
		fmt.Fprintf(&buf, "<td>%v/%v</td>", tv.Kill, tv.Death)

		for acti := 0; acti < acttype.ActType_Count; acti++ {
			fmt.Fprintf(&buf, "<td>%v</td>",
				stats.Stats[ti].ActStats[acti])
		}
		buf.WriteString(`</tr>`)
	}
	buf.WriteString(`</table>`)

	return buf.String()
}

func loadStageListHTML() string {
	tlurl := ReplacePathFromHref("stagelist.json")
	aol, err := stagelist4client.LoadFromURL(tlurl)
	if err != nil {
		jslog.Errorf("stagelist load fail %v", err)
		return "fail to load stagelist"
	}
	var buf bytes.Buffer
	buf.WriteString(`
		stage list in server
		<table border=1 style="border-collapse:collapse;">
		<tr>
		<th>Number</th> <th>UUID</th> <th>Type</th> <th>Command</th> 
		</tr>	
		`)
	for i, stg := range aol {
		fmt.Fprintf(&buf, `
		<tr>
		<td>%v</td> <td>%v</td> <td>%v</td> 
		<td><button type="button" style="font-size:20px;" onclick="enterStage('%v')">Enter Stage</button></td> 
		</tr>`,
			i, stg.UUID, stg.StageType, stg.UUID,
		)
	}
	buf.WriteString(`
		<tr>
		<th>Number</th> <th>UUID</th> <th>Type</th> <th>Command</th> 
		</tr>
		</table>
		`)
	return buf.String()
}

func ReplacePathFromHref(s string) string {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		jslog.Errorf("%v", err)
		return ""
	}
	u.Path = s
	return u.String()
}
