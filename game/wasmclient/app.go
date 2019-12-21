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

package wasmclient

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/kasworld/actjitter"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connwasm"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_obj"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_pid2rspfn"
	"github.com/kasworld/gowasmlib/textncount"
	"github.com/kasworld/intervalduration"
)

type WasmClient struct {
	DoClose              func()
	pid2recv             *w3d_pid2rspfn.PID2RspFn
	wsConn               *w3d_connwasm.Connection
	ServerJitter         *actjitter.ActJitter
	ClientJitter         *actjitter.ActJitter
	PingDur              int64
	ServerClientTictDiff int64
	DispInterDur         *intervalduration.IntervalDuration
	systemMessage        textncount.TextNCountList

	vp        *Viewport3d
	statsInfo *w3d_obj.NotiStatsInfo_data
}

func InitApp() {
	// dst := "ws://localhost:8080/ws"
	app := &WasmClient{
		DoClose:      func() { fmt.Println("Too early DoClose call") },
		pid2recv:     w3d_pid2rspfn.New(),
		ServerJitter: actjitter.New("Server"),
		ClientJitter: actjitter.New("Client"),
	}
	app.DispInterDur = intervalduration.New("Display")
	app.vp = NewViewport3d("canvas3d")

	app.ResizeCanvas()
	win := js.Global().Get("window")
	win.Call("addEventListener", "resize", js.FuncOf(app.handleResizeCanvas))

	go app.run()
}

func (app *WasmClient) run() {
	ctx, closeCtx := context.WithCancel(context.Background())
	app.DoClose = closeCtx
	defer app.DoClose()
	if err := app.NetInit(ctx); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer app.Cleanup()

	app.updataServiceInfo()
	js.Global().Call("requestAnimationFrame", js.FuncOf(app.drawCanvas))

	timerPingTk := time.NewTicker(time.Second)
	defer timerPingTk.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerPingTk.C:
			go app.reqHeartbeat()
			app.updateDebugInfo()
			app.updateTeamStatsInfo()
			app.updateSysmsg()
		}
	}
}

func (app *WasmClient) drawCanvas(this js.Value, args []js.Value) interface{} {
	defer func() {
		js.Global().Call("requestAnimationFrame", js.FuncOf(app.drawCanvas))
	}()
	dispCount := app.DispInterDur.GetCount()
	_ = dispCount
	act := app.DispInterDur.BeginAct()
	defer act.End()

	now := app.GetEstServerTick()
	app.vp.Draw(now)

	return nil
}

func (app *WasmClient) GetEstServerTick() int64 {
	return time.Now().UnixNano() + app.ServerClientTictDiff
}

func (app *WasmClient) handleResizeCanvas(this js.Value, args []js.Value) interface{} {
	app.ResizeCanvas()
	return nil
}
func (app *WasmClient) ResizeCanvas() {
	app.vp.Resize()
}
