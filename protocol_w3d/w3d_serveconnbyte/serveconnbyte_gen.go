// Code generated by "genprotocol -ver=213afa194ef0e682076c6a0cbf801946c13d343cc54330be7c4557e46057a498 -basedir=. -prefix=w3d -statstype=int"

package w3d_serveconnbyte

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_authorize"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_const"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_looptcp"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_loopwsgorilla"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_statapierror"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_statnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_statserveapi"
)

type CounterI interface {
	Inc()
}

func (scb *ServeConnByte) String() string {
	return fmt.Sprintf("ServeConnByte[SendCh:%v/%v]",
		len(scb.sendCh), cap(scb.sendCh))
}

type ServeConnByte struct {
	connData       interface{} // custom data for this conn
	sendCh         chan w3d_packet.Packet
	sendRecvStop   func()
	authorCmdList  *w3d_authorize.AuthorizedCmds
	pid2ApiStatObj *w3d_statserveapi.PacketID2StatObj
	apiStat        *w3d_statserveapi.StatServeAPI
	notiStat       *w3d_statnoti.StatNotification
	errorStat      *w3d_statapierror.StatAPIError
	sendCounter    CounterI
	recvCounter    CounterI

	demuxReq2BytesAPIFnMap [w3d_idcmd.CommandID_Count]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error)
}

// New with stats local
func New(
	connData interface{},
	sendBufferSize int,
	authorCmdList *w3d_authorize.AuthorizedCmds,
	sendCounter, recvCounter CounterI,
	demuxReq2BytesAPIFnMap [w3d_idcmd.CommandID_Count]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error),
) *ServeConnByte {
	scb := &ServeConnByte{
		connData:               connData,
		sendCh:                 make(chan w3d_packet.Packet, sendBufferSize),
		pid2ApiStatObj:         w3d_statserveapi.NewPacketID2StatObj(),
		apiStat:                w3d_statserveapi.New(),
		notiStat:               w3d_statnoti.New(),
		errorStat:              w3d_statapierror.New(),
		sendCounter:            sendCounter,
		recvCounter:            recvCounter,
		authorCmdList:          authorCmdList,
		demuxReq2BytesAPIFnMap: demuxReq2BytesAPIFnMap,
	}
	scb.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", scb)
	}
	return scb
}

// NewWithStats with stats global
func NewWithStats(
	connData interface{},
	sendBufferSize int,
	authorCmdList *w3d_authorize.AuthorizedCmds,
	sendCounter, recvCounter CounterI,
	apiStat *w3d_statserveapi.StatServeAPI,
	notiStat *w3d_statnoti.StatNotification,
	errorStat *w3d_statapierror.StatAPIError,
	demuxReq2BytesAPIFnMap [w3d_idcmd.CommandID_Count]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error),
) *ServeConnByte {
	scb := &ServeConnByte{
		connData:               connData,
		sendCh:                 make(chan w3d_packet.Packet, sendBufferSize),
		pid2ApiStatObj:         w3d_statserveapi.NewPacketID2StatObj(),
		apiStat:                apiStat,
		notiStat:               notiStat,
		errorStat:              errorStat,
		sendCounter:            sendCounter,
		recvCounter:            recvCounter,
		authorCmdList:          authorCmdList,
		demuxReq2BytesAPIFnMap: demuxReq2BytesAPIFnMap,
	}
	scb.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", scb)
	}
	return scb
}

func (scb *ServeConnByte) Disconnect() {
	scb.sendRecvStop()
}
func (scb *ServeConnByte) GetConnData() interface{} {
	return scb.connData
}
func (scb *ServeConnByte) GetAPIStat() *w3d_statserveapi.StatServeAPI {
	return scb.apiStat
}
func (scb *ServeConnByte) GetNotiStat() *w3d_statnoti.StatNotification {
	return scb.notiStat
}
func (scb *ServeConnByte) GetErrorStat() *w3d_statapierror.StatAPIError {
	return scb.errorStat
}
func (scb *ServeConnByte) GetAuthorCmdList() *w3d_authorize.AuthorizedCmds {
	return scb.authorCmdList
}
func (scb *ServeConnByte) StartServeWS(
	mainctx context.Context, conn *websocket.Conn,
	readTimeoutSec, writeTimeoutSec time.Duration,
	marshalfn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error),
) error {
	var returnerr error
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	scb.sendRecvStop = sendRecvCancel
	go func() {
		err := w3d_loopwsgorilla.RecvLoop(sendRecvCtx, scb.sendRecvStop, conn,
			readTimeoutSec, scb.handleRecvPacket)
		if err != nil {
			returnerr = fmt.Errorf("end RecvLoop %v", err)
		}
	}()
	go func() {
		err := w3d_loopwsgorilla.SendLoop(sendRecvCtx, scb.sendRecvStop, conn,
			writeTimeoutSec, scb.sendCh,
			marshalfn, scb.handleSentPacket)
		if err != nil {
			returnerr = fmt.Errorf("end SendLoop %v", err)
		}
	}()
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		}
	}
	return returnerr
}
func (scb *ServeConnByte) StartServeTCP(
	mainctx context.Context, conn *net.TCPConn,
	readTimeoutSec, writeTimeoutSec time.Duration,
	marshalfn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error),
) error {
	var returnerr error
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	scb.sendRecvStop = sendRecvCancel
	go func() {
		err := w3d_looptcp.RecvLoop(sendRecvCtx, scb.sendRecvStop, conn,
			readTimeoutSec, scb.handleRecvPacket)
		if err != nil {
			returnerr = fmt.Errorf("end RecvLoop %v", err)
		}
	}()
	go func() {
		err := w3d_looptcp.SendLoop(sendRecvCtx, scb.sendRecvStop, conn,
			writeTimeoutSec, scb.sendCh,
			marshalfn, scb.handleSentPacket)
		if err != nil {
			returnerr = fmt.Errorf("end SendLoop %v", err)
		}
	}()
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		}
	}
	return returnerr
}
func (scb *ServeConnByte) handleSentPacket(header w3d_packet.Header) error {
	scb.sendCounter.Inc()
	switch header.FlowType {
	default:
		return fmt.Errorf("invalid packet type %s %v", scb, header)

	case w3d_packet.Request:
		return fmt.Errorf("request packet not supported %s %v", scb, header)

	case w3d_packet.Response:
		statOjb := scb.pid2ApiStatObj.Del(header.ID)
		if statOjb != nil {
			statOjb.AfterSendRsp(header)
		} else {
			return fmt.Errorf("send StatObj not found %v", header)
		}
	case w3d_packet.Notification:
		scb.notiStat.Add(header)
	}
	return nil
}
func (scb *ServeConnByte) handleRecvPacket(rheader w3d_packet.Header, rbody []byte) error {
	scb.recvCounter.Inc()
	if rheader.FlowType != w3d_packet.Request {
		return fmt.Errorf("Unexpected rheader packet type: %v", rheader)
	}
	if int(rheader.Cmd) >= len(scb.demuxReq2BytesAPIFnMap) {
		return fmt.Errorf("Invalid rheader command %v", rheader)
	}
	if !scb.authorCmdList.CheckAuth(w3d_idcmd.CommandID(rheader.Cmd)) {
		return fmt.Errorf("Not authorized packet %v", rheader)
	}

	statObj, err := scb.apiStat.AfterRecvReqHeader(rheader)
	if err != nil {
		return err
	}
	if err := scb.pid2ApiStatObj.Add(rheader.ID, statObj); err != nil {
		return err
	}
	statObj.BeforeAPICall()

	// timeout api call
	apiResult := scb.callAPI_timed(rheader, rbody)
	sheader, sbody, apierr := apiResult.header, apiResult.body, apiResult.err

	// no timeout api call
	//fn := scb.demuxReq2BytesAPIFnMap[rheader.Cmd]
	//sheader, sbody, apierr := fn(scb, rheader, rbody)

	statObj.AfterAPICall()

	scb.errorStat.Inc(w3d_idcmd.CommandID(rheader.Cmd), sheader.ErrorCode)
	if apierr != nil {
		return apierr
	}
	if sbody == nil {
		return fmt.Errorf("Response body nil")
	}
	sheader.FlowType = w3d_packet.Response
	sheader.Cmd = rheader.Cmd
	sheader.ID = rheader.ID
	rpk := w3d_packet.Packet{
		Header: sheader,
		Body:   sbody,
	}
	return scb.EnqueueSendPacket(rpk)
}

type callAPIResult struct {
	header w3d_packet.Header
	body   interface{}
	err    error
}

func (scb *ServeConnByte) callAPI_timed(rheader w3d_packet.Header, rbody []byte) callAPIResult {
	rtnCh := make(chan callAPIResult, 1)
	go func(rtnCh chan callAPIResult, rheader w3d_packet.Header, rbody []byte) {
		fn := scb.demuxReq2BytesAPIFnMap[rheader.Cmd]
		sheader, sbody, apierr := fn(scb, rheader, rbody)
		rtnCh <- callAPIResult{sheader, sbody, apierr}
	}(rtnCh, rheader, rbody)
	timeoutTk := time.NewTicker(w3d_const.ServerAPICallTimeOutDur)
	defer timeoutTk.Stop()
	select {
	case apiResult := <-rtnCh:
		return apiResult
	case <-timeoutTk.C:
		return callAPIResult{rheader, nil, fmt.Errorf("APICall Timeout %v", rheader)}
	}
}
func (scb *ServeConnByte) EnqueueSendPacket(pk w3d_packet.Packet) error {
	select {
	case scb.sendCh <- pk:
		return nil
	default:
		return fmt.Errorf("Send channel full %v", scb)
	}
}
func (scb *ServeConnByte) SendNotiPacket(
	cmd w3d_idnoti.NotiID, body interface{}) error {
	err := scb.EnqueueSendPacket(w3d_packet.Packet{
		w3d_packet.Header{
			Cmd:      uint16(cmd),
			FlowType: w3d_packet.Notification,
		},
		body,
	})
	if err != nil {
		scb.Disconnect()
	}
	return err
}
