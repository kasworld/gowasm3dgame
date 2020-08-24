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

package main

import (
	"flag"
	"fmt"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	"github.com/kasworld/go-profile"
	"github.com/kasworld/gowasm3dgame/config/gameconst"
	"github.com/kasworld/gowasm3dgame/game/server"
	"github.com/kasworld/gowasm3dgame/config/serverconfig"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_version"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/signalhandle"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func printVersion() {
	fmt.Println("gowasm3dgame")
	fmt.Println("Build     ", version.GetVersion())
	fmt.Println("Data      ", gameconst.DataVersion)
	fmt.Println("Protocol  ", w3d_version.ProtocolVersion)
	fmt.Println()
}

func main() {
	printVersion()

	configurl := flag.String("i", "", "server config file or url")
	signalhandle.AddArgs()
	profile.AddArgs()

	ads := argdefault.New(&serverconfig.Config{})
	ads.RegisterFlag()
	flag.Parse()
	config := &serverconfig.Config{
		LogLevel:      w3dlog.LL_All,
		SplitLogLevel: 0,
	}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		if err := configutil.LoadIni(*configurl, &config); err != nil {
			w3dlog.Fatal("%v", err)
		}
	}
	ads.ApplyFlagTo(config)
	if *configurl == "" {
		configutil.SaveIni("w3dserver.ini", &config)
	}
	if profile.IsCpu() {
		fn := profile.StartCPUProfile()
		defer fn()
	}

	l, err := w3dlog.NewWithDstDir(
		"",
		config.MakeLogDir(),
		logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
		config.LogLevel,
		config.SplitLogLevel,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	w3dlog.GlobalLogger = l

	svr := server.New(*config)
	if err := signalhandle.StartByArgs(svr); err != nil {
		w3dlog.Error("%v", err)
	}

	if profile.IsMem() {
		profile.WriteHeapProfile()
	}
}
