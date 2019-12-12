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

package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/gowasm3dgame/game/serverconfig"
	"github.com/kasworld/gowasm3dgame/game/stage"
	"github.com/kasworld/gowasm3dgame/lib/w3dlog"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_connmanager"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_gob"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_statapierror"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_statnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_statserveapi"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/signalhandle"
	"github.com/kasworld/weblib/retrylistenandserve"
)

type Server struct {
	rnd       *rand.Rand      `prettystring:"hide"`
	log       *w3dlog.LogBase `prettystring:"hide"`
	config    serverconfig.Config
	adminWeb  *http.Server `prettystring:"simple"`
	clientWeb *http.Server `prettystring:"simple"`
	startTime time.Time    `prettystring:"simple"`

	sendRecvStop func()
	SendStat     *actpersec.ActPerSec `prettystring:"simple"`
	RecvStat     *actpersec.ActPerSec `prettystring:"simple"`

	apiStat   *w3d_statserveapi.StatServeAPI
	notiStat  *w3d_statnoti.StatNotification
	errorStat *w3d_statapierror.StatAPIError

	marshalBodyFn          func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error)
	unmarshalPacketFn      func(h w3d_packet.Header, bodyData []byte) (interface{}, error)
	DemuxReq2BytesAPIFnMap [w3d_idcmd.CommandID_Count]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error)

	connManager *w3d_connmanager.Manager
	stage       *stage.Stage
}

func New(config serverconfig.Config) *Server {
	svr := &Server{
		config: config,
		log:    w3dlog.GlobalLogger,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),

		SendStat: actpersec.New(),
		RecvStat: actpersec.New(),

		apiStat:     w3d_statserveapi.New(),
		notiStat:    w3d_statnoti.New(),
		errorStat:   w3d_statapierror.New(),
		connManager: w3d_connmanager.New(),
	}
	svr.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call\n")
	}
	svr.marshalBodyFn = w3d_gob.MarshalBodyFn
	svr.unmarshalPacketFn = w3d_gob.UnmarshalPacket
	svr.DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error){
		w3d_idcmd.Invalid:   svr.bytesAPIFn_ReqInvalid,
		w3d_idcmd.MakeTeam:  svr.bytesAPIFn_ReqMakeTeam,
		w3d_idcmd.Act:       svr.bytesAPIFn_ReqAct,
		w3d_idcmd.Heartbeat: svr.bytesAPIFn_ReqHeartbeat,
	}
	svr.stage = stage.New(svr.log, svr.config)
	return svr
}

// called from signal handler
func (svr *Server) GetServiceLockFilename() string {
	return svr.config.MakePIDFileFullpath()
}

// called from signal handler
func (svr *Server) GetLogger() signalhandle.LoggerI {
	return w3dlog.GlobalLogger
}

// called from signal handler
func (svr *Server) ServiceInit() error {
	return nil
}

// called from signal handler
func (svr *Server) ServiceCleanup() {
}

// called from signal handler
func (svr *Server) ServiceMain(ctx context.Context) {
	fmt.Println(prettystring.PrettyString(svr.config, 4))
	svr.startTime = time.Now()

	ctx, stopFn := context.WithCancel(ctx)
	svr.sendRecvStop = stopFn
	defer svr.sendRecvStop()

	svr.initAdminWeb()
	svr.initServiceWeb(ctx)

	fmt.Printf("open admin web\nhttp://localhost%v/\n", svr.config.AdminPort)
	fmt.Printf("open client web\nhttp://localhost%v/\n", svr.config.ServicePort)

	go retrylistenandserve.RetryListenAndServe(svr.adminWeb, svr.log, "serveAdminWeb")
	go retrylistenandserve.RetryListenAndServe(svr.clientWeb, svr.log, "serveServiceWeb")

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	svr.stage = stage.New(svr.log, svr.config)
	go svr.stage.Run(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerInfoTk.C:
			svr.SendStat.UpdateLap()
			svr.RecvStat.UpdateLap()
		}
	}
}
