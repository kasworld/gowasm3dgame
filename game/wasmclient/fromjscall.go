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
	"strings"
	"syscall/js"

	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/lib/jskeypressmap"
	"github.com/kasworld/gowasm3dgame/lib/jsobj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (app *WasmClient) registerJSButton() {
	js.Global().Set("sendChat", js.FuncOf(app.jsSendChat))
}

func (app *WasmClient) jsSendChat(this js.Value, args []js.Value) interface{} {
	msg := getChatMsg()
	go app.sendPacket(w3d_idcmd.Chat,
		&w3d_obj.ReqChat_data{Chat: msg})
	app.vp.Focus()
	return nil
}

func getChatMsg() string {
	msg := jsobj.GetTextValueFromInputText("chattext")
	msg = strings.TrimSpace(msg)
	if len(msg) > gameconst.MaxChatLen {
		msg = msg[:gameconst.MaxChatLen]
	}
	return msg
}

func (app *WasmClient) registerKeyboardMouseEvent() {
	app.vp.AddEventListener("click", app.jsHandleMouseClick)
	app.vp.AddEventListener("mousemove", app.jsHandleMouseMove)
	app.vp.AddEventListener("mousedown", app.jsHandleMouseDown)
	app.vp.AddEventListener("mouseup", app.jsHandleMouseUp)
	app.vp.AddEventListener("wheel", app.jsHandleMouseWheel)
	app.vp.AddEventListener("contextmenu", app.jsHandleContextMenu)
	app.vp.AddEventListener("keydown", app.jsHandleKeyDown)
	app.vp.AddEventListener("keypress", app.jsHandleKeyPress)
	app.vp.AddEventListener("keyup", app.jsHandleKeyUp)
}

func (app *WasmClient) jsHandleMouseClick(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")

	mouseX, mouseY := evt.Get("offsetX").Int(), evt.Get("offsetY").Int()
	btn := evt.Get("button").Int()

	_ = mouseX
	_ = mouseY
	switch btn {
	case 0: // left
		// app.makePathToMouseClick(mouseX, mouseY)
		return nil
	case 1: // wheel

	case 2: // right
	}
	return nil
}

func (app *WasmClient) jsHandleMouseWheel(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	// Never call,  relate focus , prevent key event listen
	// evt.Call("stopPropagation")
	evt.Call("preventDefault")

	dx := evt.Get("deltaX").Float()
	dy := evt.Get("deltaY").Float()
	jslog.Infof("wheel %v %v", dx, dy)
	if dy > 0 {
	} else if dy < 0 {
	}
	app.ResizeCanvas()
	return nil
}

func (app *WasmClient) jsHandleMouseDown(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	// Never call,  relate focus , prevent key event listen
	// evt.Call("stopPropagation")
	// evt.Call("preventDefault")
	btn := evt.Get("button").Int()
	// jslog.Infof("mousedown %v", btn)

	switch btn {
	case 0: // left
	case 1: // wheel
	case 2: // right click
		// app.actByMouseRightDown()
	}
	return nil
}
func (app *WasmClient) jsHandleMouseUp(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")
	return nil
}
func (app *WasmClient) jsHandleContextMenu(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")
	return nil
}

func (app *WasmClient) jsHandleMouseMove(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")

	mouseX, mouseY := evt.Get("offsetX").Int(), evt.Get("offsetY").Int()
	_, _ = mouseX, mouseY
	return nil
}

func (app *WasmClient) jsHandleKeyDown(this js.Value, args []js.Value) interface{} {
	// jslog.Info("jsHandleKeyDownVP")
	evt := args[0]
	if evt.Get("target").Equal(app.vp.canvas) {
		evt.Call("stopPropagation")
		evt.Call("preventDefault")

		kcode := evt.Get("key").String()
		// jslog.Infof("%v %v", evt, kcode)
		if kcode != "" {
			app.KeyboardPressedMap.KeyDown(kcode)
		}
		app.actByKeyPressMap(kcode)
	}
	return nil
}
func (app *WasmClient) jsHandleKeyPress(this js.Value, args []js.Value) interface{} {
	// jslog.Info("jsHandleKeyPressVP")
	evt := args[0]
	if evt.Get("target").Equal(app.vp.canvas) {
		evt.Call("stopPropagation")
		evt.Call("preventDefault")

		kcode := evt.Get("key").String()
		app.actByKeyPressMap(kcode)
	}
	return nil
}
func (app *WasmClient) jsHandleKeyUp(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	if evt.Get("target").Equal(app.vp.canvas) {
		evt.Call("stopPropagation")
		evt.Call("preventDefault")

		kcode := evt.Get("key").String()
		// jslog.Infof("%v %v", evt, kcode)
		if kcode != "" {
			app.KeyboardPressedMap.KeyUp(kcode)
		}
		app.processKeyUpEvent(kcode)
	}
	return nil
}

func (app *WasmClient) actByKeyPressMap(kcode string) bool {
	dx, dy := app.KeyboardPressedMap.SumMoveDxDy(jskeypressmap.Key2Dir)
	_ = dx
	_ = dy
	return false
}

func (app *WasmClient) processKeyUpEvent(kcode string) bool {
	// jslog.Errorf("keyup %v", kcode)
	if kcode == "Escape" {
		// reset to default
		app.ResizeCanvas()
		return true
	}

	if btn := gameOptions.GetByKeyCode(kcode); btn != nil {
		btn.JSFn(js.Null(), nil)
		return true
	}

	return false
}
