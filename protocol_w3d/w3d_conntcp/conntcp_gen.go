// Code generated by "genprotocol.exe -ver=a99e26984c4d0465e81623a4767d3a2aa3cb4fcc9890904054c0c51f30e0b79f -basedir=protocol_w3d -prefix=w3d -statstype=int"

package w3d_conntcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_looptcp"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

type Connection struct {
	conn         *net.TCPConn
	sendCh       chan w3d_packet.Packet
	sendRecvStop func()

	readTimeoutSec     time.Duration
	writeTimeoutSec    time.Duration
	marshalBodyFn      func(interface{}, []byte) ([]byte, byte, error)
	handleRecvPacketFn func(header w3d_packet.Header, body []byte) error
	handleSentPacketFn func(header w3d_packet.Header) error
}

func New(
	readTimeoutSec, writeTimeoutSec time.Duration,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleRecvPacketFn func(header w3d_packet.Header, body []byte) error,
	handleSentPacketFn func(header w3d_packet.Header) error,
) *Connection {
	tc := &Connection{
		sendCh:             make(chan w3d_packet.Packet, 10),
		readTimeoutSec:     readTimeoutSec,
		writeTimeoutSec:    writeTimeoutSec,
		marshalBodyFn:      marshalBodyFn,
		handleRecvPacketFn: handleRecvPacketFn,
		handleSentPacketFn: handleSentPacketFn,
	}

	tc.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", tc)
	}
	return tc
}

func (tc *Connection) ConnectTo(remoteAddr string) error {
	tcpaddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return err
	}
	tc.conn, err = net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		return err
	}
	return nil
}

func (tc *Connection) Cleanup() {
	tc.sendRecvStop()
	if tc.conn != nil {
		tc.conn.Close()
	}
}

func (tc *Connection) Run(mainctx context.Context) error {
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	tc.sendRecvStop = sendRecvCancel
	var rtnerr error
	var sendRecvWaitGroup sync.WaitGroup
	sendRecvWaitGroup.Add(2)
	go func() {
		defer sendRecvWaitGroup.Done()
		err := w3d_looptcp.RecvLoop(
			sendRecvCtx,
			tc.sendRecvStop,
			tc.conn,
			tc.readTimeoutSec,
			tc.handleRecvPacketFn)
		if err != nil {
			rtnerr = err
		}
	}()
	go func() {
		defer sendRecvWaitGroup.Done()
		err := w3d_looptcp.SendLoop(
			sendRecvCtx,
			tc.sendRecvStop,
			tc.conn,
			tc.writeTimeoutSec,
			tc.sendCh,
			tc.marshalBodyFn,
			tc.handleSentPacketFn)
		if err != nil {
			rtnerr = err
		}
	}()
	sendRecvWaitGroup.Wait()
	return rtnerr
}

func (tc *Connection) EnqueueSendPacket(pk w3d_packet.Packet) error {
	select {
	case tc.sendCh <- pk:
		return nil
	default:
		return fmt.Errorf("Send channel full %v", tc)
	}
}
