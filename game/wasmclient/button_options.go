// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	"syscall/js"

	"github.com/kasworld/gowasm3dgame/lib/htmlbutton"
	"github.com/kasworld/gowasmlib/jslog"
)

var gameOptions *htmlbutton.HTMLButtonGroup

// prevent compiler initialize loop error
var _gameopt = htmlbutton.NewButtonGroup("Options",
	[]*htmlbutton.HTMLButton{
		htmlbutton.New("q", "LeftInfo", []string{"LeftInfoOff", "LeftInfoOn"},
			"show/hide left info", cmdToggleLeftInfo, 1),
		htmlbutton.New("w", "CenterInfo", []string{"CenterInfoOff", "CenterInfoOn"},
			"show/hide center info", cmdRotateCenterInfo, 0),
		htmlbutton.New("e", "RightInfo", []string{"RightInfoOff", "RightInfoOn"},
			"show/hide right info", cmdRotateRightInfo, 1),
	})

func cmdToggleLeftInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateLeftInfo()
	app.vp.Focus()
}

func (app *WasmClient) updateLeftInfo() {
	v := gameOptions.GetByIDBase("LeftInfo")
	infoobj := js.Global().Get("document").Call("getElementById", "leftinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // leftinfo
		infoobj.Set("innerHTML", app.makeServiceInfo()+app.makeDebugInfo())
	}
}

func cmdRotateCenterInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateCenterInfo()
	app.vp.Focus()
}

func (app *WasmClient) updateCenterInfo() {
	v := gameOptions.GetByIDBase("CenterInfo")
	infoobj := js.Global().Get("document").Call("getElementById", "centerinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // centerinfo
		infoobj.Set("innerHTML", app.makeTeamStatsInfo())
	}
}

func cmdRotateRightInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateRightInfo()
	app.vp.Focus()
}

func (app *WasmClient) updateRightInfo() {
	v := gameOptions.GetByIDBase("RightInfo")
	infoobj := js.Global().Get("document").Call("getElementById", "rightinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // right info
		app.systemMessage = app.systemMessage.GetLastN(50)
		infoobj.Set("innerHTML", app.systemMessage.ToHtmlStringRev())
	}
}
