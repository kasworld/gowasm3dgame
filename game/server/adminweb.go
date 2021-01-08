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

package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/version"
	"github.com/kasworld/weblib"
	"github.com/kasworld/weblib/webprofile"
)

func (svr *Server) web_FaviconIco(w http.ResponseWriter, r *http.Request) {
}

func (svr *Server) initAdminWeb() {
	authdata := weblib.NewAuthData("server")
	authdata.ReLoadUserData([][2]string{
		{svr.config.WebAdminID, svr.config.WebAdminPass},
	})
	webMux := weblib.NewAuthMux(authdata, svr.log)

	if !version.IsRelease() {
		webprofile.AddWebProfile(webMux)
	}

	webMux.HandleFunc("/favicon.ico", svr.web_FaviconIco)
	webMux.HandleFuncAuth("/", svr.web_ServerInfo)

	webMux.HandleFuncAuth("/StatServeAPI", svr.web_ProtocolStat)
	webMux.HandleFuncAuth("/StatNotification", svr.web_NotiStat)
	webMux.HandleFuncAuth("/StatAPIError", svr.web_ErrorStat)
	webMux.HandleFuncAuth("/ConnectionManager", svr.web_ConnMan)
	webMux.HandleFuncAuth("/StageManager", svr.web_StageMan)

	authdata.AddAllActionName(svr.config.WebAdminID)

	svr.adminWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", svr.config.AdminPort),
	}
}

func (svr *Server) web_ServerInfo(w http.ResponseWriter, r *http.Request) {
	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>Server admin</title>
	</head>
	<body>

	BuildDate : {{.BuildDate.Format "2006-01-02T15:04:05Z07:00"}}
	<br/>
	Version: {{.GetVersion}}
	<br/>
	ProtocolVersion : {{.GetProtocolVersion}}
	<br/>
	DataVersion : {{.GetDataVersion}}
	<br/>
	Start : {{.GetStartTime}} / {{.GetRunDur}}
	<br/>
	goroutine : {{.NumGoroutine}}	
	<br/>
	global wrapper : {{.WrapInfo}}	
	<br/>
	SendStat : {{.GetSendStat}}
	<br/>
	RecvStat : {{.GetRecvStat}}
	<br/>
    <a href="/StatServeAPI" target="_blank">{{.GetProtocolStat}}</a>
    <br/>
    <a href="/StatNotification" target="_blank">{{.GetNotiStat}}</a>
    <br/>
    <a href="/StatAPIError" target="_blank">{{.GetErrorStat}}</a>
    <br/>
    <a href="/ConnectionManager?page=0" target="_blank">{{.GetConnMan}}</a>
    <br/>
    <a href="/StageManager?page=0" target="_blank">{{.GetStageMan}}</a>
    <br/>
	<pre>{{.Config.StringForm}}</pre>
	<br/>
	</body> </html> 
	`)
	if err != nil {
		svr.log.Error("%v", err)
	}
	if err := weblib.SetFresh(w, r); err != nil {
		svr.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, svr); err != nil {
		svr.log.Error("%v", err)
	}
}

func (svr *Server) web_ProtocolStat(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		svr.log.Error("%v", err)
	}
	svr.apiStat.ToWeb(w, r)
}

func (svr *Server) web_NotiStat(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		svr.log.Error("%v", err)
	}
	svr.notiStat.ToWeb(w, r)
}

func (svr *Server) web_ErrorStat(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		svr.log.Error("%v", err)
	}
	svr.errorStat.ToWeb(w, r)
}

func (svr *Server) web_ConnMan(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		svr.log.Error("%v", err)
	}
	svr.connManager.ToWeb(w, r)
}

func (svr *Server) web_StageMan(w http.ResponseWriter, r *http.Request) {
	if err := weblib.SetFresh(w, r); err != nil {
		svr.log.Error("%v", err)
	}
	svr.stageManager.ToWeb(w, r)
}
