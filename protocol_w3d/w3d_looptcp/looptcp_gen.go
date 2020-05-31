// Code generated by "genprotocol -ver=69ee1a9a014411856ae5de810035b3c3478b3ee188bf040ff46c3ec85f20aba5 -basedir=. -prefix=w3d -statstype=int"

package w3d_looptcp

import (
	"context"
	"net"
	"time"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_const"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

var bufPool = w3d_packet.NewPool(w3d_const.PacketBufferPoolSize)

func SendPacket(conn *net.TCPConn, buf []byte) error {
	toWrite := len(buf)
	for l := 0; l < toWrite; {
		n, err := conn.Write(buf[l:toWrite])
		if err != nil {
			return err
		}
		l += n
	}
	return nil
}

func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	SendCh chan w3d_packet.Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleSentPacketFn func(header w3d_packet.Header) error,
) error {

	defer SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		case pk := <-SendCh:
			if err = tcpConn.SetWriteDeadline(time.Now().Add(timeOut)); err != nil {
				break loop
			}
			oldbuf := bufPool.Get()
			sendBuffer, err := w3d_packet.Packet2Bytes(&pk, marshalBodyFn, oldbuf)
			if err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = SendPacket(tcpConn, sendBuffer); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = handleSentPacketFn(pk.Header); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			bufPool.Put(oldbuf)
		}
	}
	return err
}

func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	HandleRecvPacketFn func(header w3d_packet.Header, body []byte) error,
) error {

	defer SendRecvStop()

	pb := w3d_packet.NewRecvPacketBuffer()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			return nil

		default:
			if pb.IsPacketComplete() {
				header, rbody, lerr := pb.GetHeaderBody()
				if lerr != nil {
					err = lerr
					break loop
				}
				if err = HandleRecvPacketFn(header, rbody); err != nil {
					break loop
				}
				pb = w3d_packet.NewRecvPacketBuffer()
				if err = tcpConn.SetReadDeadline(time.Now().Add(timeOut)); err != nil {
					break loop
				}
			} else {
				err := pb.Read(tcpConn)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}
